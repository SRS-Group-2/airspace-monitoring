package main

import (
	gcp "aircraft_positions/gcp"
	"context"
	"log"
	"net/http"
	"os"
	"unicode"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/secure"
)

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const env_topicID = "GOOGLE_PUBSUB_AIRCRAFT_POSITIONS_TOPIC_ID"
const logName = "AIRCRAFT_POSITIONS_LOG"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

var posTopic *pubsub.Topic

var hub = Hub{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type LogTypes struct {
	Debug    *log.Logger
	Error    *log.Logger
	Critical *log.Logger
}

var Log = LogTypes{}

func main() {
	var projectID = mustGetenv(env_projectID)
	var topicID = mustGetenv(env_topicID)

	ctx := context.Background()
	loggerClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer loggerClient.Close()

	Log.Debug = loggerClient.Logger(logName).StandardLogger(logging.Debug)
	Log.Error = loggerClient.Logger(logName).StandardLogger(logging.Error)
	Log.Critical = loggerClient.Logger(logName).StandardLogger(logging.Critical)

	Log.Debug.Println("Starting Aircraft Positions Service.")
	defer Log.Debug.Println("Stopping Aircraft Positions Service.")

	err = gcp.Initialize(projectID)
	checkErr(err)
	defer gcp.PubsubClient.Close()

	posTopic, err = gcp.GetTopic(topicID)
	checkErr(err)

	hub.icaos = make(map[string]*PubSubListener)

	router := gin.New()

	router.SetTrustedProxies(nil)
	router.Use(secure.New(secure.Config{
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))
	router.GET("/airspace/aircraft/:icao24/position", httpRequestHandler)

	router.Run()
}

func httpRequestHandler(c *gin.Context) {
	icao24 := c.Param("icao24")

	if !validIcao24(icao24) {
		c.String(http.StatusNotAcceptable, "invalid icao24")
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// TODO check what best to do here
		c.String(http.StatusUpgradeRequired, "Websocket upgrade failed")
		Log.Error.Println("Error: websocket upgrade failed for client with icao24=", icao24, " err: ", err)
		return
	}

	cancellableCtx, cancel := context.WithCancel(context.Background())

	cl := &WSClient{
		ws:           ws,
		icao24:       icao24,
		cliId:        uuid.New().String(),
		ch:           make(chan string, 100),
		chErr:        make(chan string, 10),
		clientCtx:    cancellableCtx,
		clientCancel: cancel,
	}

	go hub.RegisterClient(cl)

	go cl.WSWriteLoop()
	go cl.WSReadLoop()
}

func validIcao24(icao24 string) bool {
	if len(icao24) != 6 || hasSymbol(icao24) {
		return false
	}
	return true
}

func hasSymbol(str string) bool {
	for _, letter := range str {
		if unicode.IsSymbol(letter) {
			return true
		}
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		Log.Critical.Println("critical error, panicking: ", err)
		panic(err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("Environment variable not set: " + k)
	}
	return v
}

func trySend(ch chan int) {
	select {
	case ch <- 1:
	default:
	}
}
