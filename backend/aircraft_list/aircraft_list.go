package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
)

func main() {

	var projectID = "fake-pubsub-123"
	var subID = "aircraft_list_sub"
	var topicID = "states-topic"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	checkErr(err)
	defer client.Close()

	topic := client.Topic(mustGetenv(topicID))

	createSubscription(ctx, client, subID, topic)

	router := gin.New()

	router.SetTrustedProxies(nil)
	router.GET("/airspace/aircraft/list", getList)

	router.Run()
}

func getList(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, "Yee")
}

func createSubscription(ctx context.Context, client *pubsub.Client, subID string, topic *pubsub.Topic) error {
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic

	sub, err := client.CreateSubscription(ctx, subID,
		pubsub.SubscriptionConfig{Topic: topic})
	checkErr(err)
	log.Println("Created subscription: \n", sub)
	return nil
}

func pubsubMessageHandler(ctx context.Context, msg *pubsub.Message) {
	//TODO: Update the internal list
	msg.Ack()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}
