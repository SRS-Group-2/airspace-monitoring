package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
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
	val HistoryValues
}

func (l *HistoryState) Read() HistoryValues {
	l.RLock()
	defer l.RUnlock()
	return l.val
}

func (l *HistoryState) Write(newVal HistoryValues) {
	l.Lock()
	defer l.Unlock()
	l.val = newVal
}

var oneHourState = HistoryState{
	val: HistoryValues{
		CO2:       0,
		Distance:  0,
		Timestamp: "0000-00-00-00-00",
	},
}

var sixHoursState = HistoryState{
	val: HistoryValues{
		CO2:       0,
		Distance:  0,
		Timestamp: "0000-00-00-00-00",
	},
}
var oneDayState = HistoryState{
	val: HistoryValues{
		CO2:       0,
		Distance:  0,
		Timestamp: "0000-00-00-00-00",
	},
}

func main() {

	var credString = mustGetenv(env_credJson)
	var projectID = mustGetenv(env_projectID)

	go firestoreUpdates(credString, projectID, "1h-history", oneHourState)
	go firestoreUpdates(credString, projectID, "6h-history", sixHoursState)
	go firestoreUpdates(credString, projectID, "24h-history", oneDayState)

	router := gin.New()
	router.SetTrustedProxies(nil)

	// eg. airspace/history/realtime/:interval
	// interval: 1h, 6h, 24h
	router.GET("/airspace/history/realtime/:interval", func(c *gin.Context) {
		interval := c.Param("interval")

		switch interval {
		case "24h":
			c.JSON(http.StatusOK, oneDayState.Read())
		case "6h":
			c.JSON(http.StatusOK, sixHoursState.Read())
		case "1h":
			fallthrough
		default:
			c.JSON(http.StatusOK, oneHourState.Read())
		}
	})

	router.Run()
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

func firestoreUpdates(credString string, projectID string, documentID string, historyState HistoryState) {

	client := FirestoreInit([]byte(credString), projectID)
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
		fmt.Fprintf(w, "Received document snapshot: %v\n", snap.Data())
	}
}

func pubsubHandler(topicID string, subscriptionID string) {

	topic, err := GcpGetTopic(topicID)
	checkErr(err)

	sub, err := GcpGetSubscription(subscriptionID, topic)
	checkErr(err)

	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		messageHandler(subscriptionID, ctx, msg)
	})
	checkErr(err)
}

func messageHandler(subscriptionID string, ctx context.Context, msg *pubsub.Message) {
	defer func() {
		if r := recover(); r != nil {
			msg.Ack()
		}
	}()

	msg.Ack()

	aircraftList.Write(string(msg.Data))
}
