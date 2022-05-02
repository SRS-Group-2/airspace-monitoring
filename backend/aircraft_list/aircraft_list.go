package main

import (
	"context"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
)

const env_credFile = "GOOGLE_APPLICATION_CREDENTIALS"

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const env_topicID = "GOOGLE_PUBSUB_AIRCRAFT_LIST_TOPIC_ID"
const env_subID = "GOOGLE_PUBSUB_AIRCRAFT_LIST_SUBSCRIBER_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type AircraftList struct {
	sync.RWMutex
	val string
}

func (l *AircraftList) Read() string {
	l.RLock()
	defer l.RUnlock()
	return l.val
}

func (l *AircraftList) Write(newVal string) {
	l.Lock()
	defer l.Unlock()
	l.val = newVal
}

var aircraftList AircraftList

const startVal string = `
{
	"timestamp":0,
	"list":[]
}`

func main() {

	aircraftList.Write(startVal)

	var credFile = mustGetenv(env_credFile)
	var projectID = mustGetenv(env_projectID)
	var topicID = mustGetenv(env_topicID)
	var subID = mustGetenv(env_subID)

	GcpInitialize(credFile, projectID)

	go pubsubHandler(topicID, subID)

	router := gin.New()

	router.SetTrustedProxies(nil)
	router.GET("/airspace/aircraft/list", getList)

	router.Run()
}

func getList(c *gin.Context) {

	c.String(http.StatusOK, aircraftList.Read())
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
