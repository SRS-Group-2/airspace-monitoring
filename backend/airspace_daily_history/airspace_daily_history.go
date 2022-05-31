package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/secure"
)

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const logName = "REALTIME_HISTORY_LOG"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type HistoryState struct {
	sync.RWMutex
	val string
}

func (l *HistoryState) Read() string {
	l.RLock()
	defer l.RUnlock()
	return l.val
}

func (l *HistoryState) Write(newVal string) {
	l.Lock()
	defer l.Unlock()
	l.val = newVal
}

var oneHourState = HistoryState{
	val: "{}",
}

var sixHoursState = HistoryState{
	val: "{}",
}
var oneDayState = HistoryState{
	val: "{}",
}

type LogType struct {
	Debug    *log.Logger
	Error    *log.Logger
	Critical *log.Logger
}

var Log = LogType{}

func main() {

	var projectID = mustGetenv(env_projectID)

	ctx := context.Background()
	loggerClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer loggerClient.Close()

	Log.Debug = loggerClient.Logger(logName).StandardLogger(logging.Debug)
	Log.Error = loggerClient.Logger(logName).StandardLogger(logging.Error)
	Log.Critical = loggerClient.Logger(logName).StandardLogger(logging.Critical)

	Log.Debug.Print("Starting Daily History Service.")
	defer Log.Debug.Println("Stopping Daily History Service.")

	go backgroundUpdateState(projectID, "1h-history", &oneHourState)
	go backgroundUpdateState(projectID, "6h-history", &sixHoursState)
	go backgroundUpdateState(projectID, "24h-history", &oneDayState)

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(func() gin.HandlerFunc {
        return func(c *gin.Context) {
            c.Writer.Header().Set("Cache-Control", "public, max-age=300")
        }
    }())

	router.Use(secure.New(secure.Config{
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	// Interval: 1h, 6h, 24h
	router.GET("/airspace/history/realtime/:interval", func(c *gin.Context) {
		interval := c.Param("interval")

		switch interval {
		case "24h":
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(oneDayState.Read() ))
			// c.String(http.StatusOK, oneDayState.Read())
		case "6h":
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(sixHoursState.Read() ))
			// c.String(http.StatusOK, sixHoursState.Read())
		case "1h":
			fallthrough
		default:
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(oneHourState.Read() ))
			// c.String(http.StatusOK, oneHourState.Read())
		}
	})

	router.Run()
}

func backgroundUpdateState(projectId string, documentID string, historyState *HistoryState) {
	Log.Debug.Println("Starting background update thread (listening to changes to ", documentID, ")")
	defer Log.Debug.Println("Stopping background update thread (listening to changes to ", documentID, ")")

	client := FirestoreInit(projectId)
	defer client.Close()

	ctx := context.Background()
	snapIter := client.Collection("airspace").Doc(documentID).Snapshots(ctx)

	for {
		// Wait for new snapshots of the document
		snap, err := snapIter.Next()
		checkErr(err)

		if !snap.Exists() {
			Log.Critical.Println("Document no longer exists, panicking.")
			panic("Document no longer exists.")
		}

		jsonData, err := json.Marshal(snap.Data())
		checkErr(err)

		historyState.Write(string(jsonData))
	}
}

func checkErr(err error) {
	if err != nil {
		Log.Critical.Println("Critical error, panicking: ", err)
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
