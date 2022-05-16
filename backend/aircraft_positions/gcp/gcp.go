package gcp

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/option"

	"cloud.google.com/go/pubsub"
)

// PubsubClient is the GCP pubsub service client.
var PubsubClient *pubsub.Client

// Initialize initializes GCP client service using the environment.
func Initialize(credJson string, projectName string) error {
	var err error
	ctx := context.Background()
	credentialOpt := option.WithCredentialsJSON([]byte(credJson))

	PubsubClient, err = pubsub.NewClient(ctx, projectName, credentialOpt)
	return err
}

// GetTopic returns the specified topic in GCP pub/sub service and create it if it not exist.
func GetTopic(topicName string) (*pubsub.Topic, error) {
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
func GetSubscription(subName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
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

// CreateSubscription creates the specified subscription in GCP pub/sub service and fails it if it already exist.
func CreateSubscription(subName string, topic *pubsub.Topic, filter string) (*pubsub.Subscription, error) {
	sub := PubsubClient.Subscription(subName)
	ctx := context.Background()
	isSubExist, err := sub.Exists(ctx)

	if err != nil {
		return sub, err
	}

	if isSubExist {
		err = errors.New("subscription already exists")
		return sub, err
	}

	ctx = context.Background()
	sub, err = PubsubClient.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{Topic: topic, Filter: filter})

	return sub, err
}

func DeleteSubscription(subID string, topic *pubsub.Topic) error {
	ctx := context.Background()

	sub := PubsubClient.Subscription(subID)
	if err := sub.Delete(ctx); err != nil {
		return fmt.Errorf("delete: %v", err)
	}
	return nil
}
