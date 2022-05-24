package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func FirestoreInitWithCredentials(projectID string, credJson []byte) *firestore.Client {
	// Use JSON credentials
	ctx := context.Background()
	opt := option.WithCredentialsJSON(credJson)
	client, err := firestore.NewClient(ctx, projectID, opt)
	checkErr(err)

	return client
}

func FirestoreInit(projectID string) *firestore.Client {
	// Use JSON credentials
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	checkErr(err)

	return client
}
