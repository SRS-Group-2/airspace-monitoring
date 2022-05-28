package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/logging"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const logName = "AIRCRAFT_LIST_LOG"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type AircraftList struct {
	sync.RWMutex
	val []byte
}

func (l *AircraftList) Read() []byte {
	l.RLock()
	defer l.RUnlock()
	return l.val
}

func (l *AircraftList) Write(newVal []byte) {
	l.Lock()
	defer l.Unlock()
	l.val = newVal
}

var aircraftList = AircraftList{
	val: []byte(`
{
	"timestamp":0,
	"list":[]
}`),
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

	Log.Debug.Print("Starting Aircraft List Service.")
	defer Log.Debug.Println("Stopping Aircraft List Service.")

	go backgroundUpdateState(projectID, "aircraft-list", &aircraftList)

	router := gin.New()

	router.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Cache-Control", "public, max-age=120")
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

	router.SetTrustedProxies(nil)
	router.GET("/airspace/aircraft/list", getList)

	router.Run()
}

func getList(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", aircraftList.Read())
}

func checkErr(err error) {
	if err != nil {
		Log.Critical.Println("A critical error occurred causing a panic: ", err)
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

func backgroundUpdateState(projectId string, documentID string, state *AircraftList) {
	Log.Debug.Println("Starting background update thread (listen to db changes).")
	defer Log.Debug.Println("Stopping background update thread.")

	client := FirestoreInit(projectId)
	defer client.Close()

	ctx := context.Background()
	snapIter := client.Collection("airspace").Doc(documentID).Snapshots(ctx)

	for {
		// Wait for new snapshots of the document
		snap, err := snapIter.Next()
		checkErr(err)

		if !snap.Exists() {
			Log.Critical.Println("Aircraft list db document no longer exists, panicking with err: ", err)
			panic("Document no longer exists.")
		}

		jsonData, err := json.Marshal(snap.Data())
		checkErr(err)

		jsonString := []byte(jsonData)
		state.Write(jsonString)
	}
}
