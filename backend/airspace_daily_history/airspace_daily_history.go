package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

const env_credJson = "GOOGLE_APPLICATION_CREDENTIALS"
const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type HistoryValues struct {
	CO2       int64
	Distance  int64
	Timestamp string
}

type Interval struct {
	oneHour  string
	sixHours string
	oneDay   string
}

var interval = Interval{
	oneHour:  "1h",
	sixHours: "6h",
	oneDay:   "24h",
}

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

func main() {

	var credJson = mustGetenv(env_credJson)
	var projectID = mustGetenv(env_projectID)

	go backgroundUpdateState(credJson, projectID, "1h-history", &oneHourState)
	go backgroundUpdateState(credJson, projectID, "6h-history", &sixHoursState)
	go backgroundUpdateState(credJson, projectID, "24h-history", &oneDayState)

	router := gin.New()
	router.SetTrustedProxies(nil)

	// eg. airspace/history/realtime/:interval
	// interval: 1h, 6h, 24h
	router.GET("/airspace/history/realtime/:interval", func(c *gin.Context) {
		interval := c.Param("interval")

		switch interval {
		case "24h":
			c.String(http.StatusOK, oneDayState.Read())
		case "6h":
			c.String(http.StatusOK, sixHoursState.Read())
		case "1h":
			fallthrough
		default:
			c.String(http.StatusOK, oneHourState.Read())
		}
	})

	router.Run()
}

func backgroundUpdateState(credString string, projectId string, documentID string, historyState *HistoryState) {

	client := FirestoreInit([]byte(credString), projectId)
	defer client.Close()

	ctx := context.Background()
	snapIter := client.Collection("airspace").Doc(documentID).Snapshots(ctx)

	for {
		// Wait for new snapshots of the document
		snap, err := snapIter.Next()
		checkErr(err)

		if !snap.Exists() {
			panic("Document no longer exists.")
		}

		jsonData, err := json.Marshal(snap.Data())
		checkErr(err)

		historyState.Write(string(jsonData))
	}
}

func checkErr(err error) {
	if err != nil {
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
