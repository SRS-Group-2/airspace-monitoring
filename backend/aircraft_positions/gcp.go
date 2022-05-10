package main

import (
	"context"

	"google.golang.org/api/option"

	"cloud.google.com/go/pubsub"
)

// PubsubClient is the GCP pubsub service client.
var PubsubClient *pubsub.Client

// Initialize initializes GCP client service using the environment.
func GcpInitialize(credJson string, projectName string) error {
	var err error
	ctx := context.Background()
	credentialOpt := option.WithCredentialsJSON([]byte(credJson))

	PubsubClient, err = pubsub.NewClient(ctx, projectName, credentialOpt)
	return err
}

// GetTopic returns the specified topic in GCP pub/sub service and create it if it not exist.
func GcpGetTopic(topicName string) (*pubsub.Topic, error) {
	topic := PubsubClient.Topic(topicName)
	ctx := context.Background()
	isTopicExist, err := topic.Exists(ctx)

	if err != nil {
		return topic, err
	}

	if !isTopicExist {
		ctx = context.Background()
		topic, err = PubsubClient.CreateTopic(ctx, topicName)
	}

	return topic, err
}

// GetSubscription returns the specified subscription in GCP pub/sub service and creates it if it not exist.
func GcpGetSubscription(subName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	sub := PubsubClient.Subscription(subName)
	ctx := context.Background()
	isSubExist, err := sub.Exists(ctx)

	if err != nil {
		return sub, err
	}

	if !isSubExist {
		ctx = context.Background()
		sub, err = PubsubClient.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{Topic: topic})
	}

	return sub, err
}
