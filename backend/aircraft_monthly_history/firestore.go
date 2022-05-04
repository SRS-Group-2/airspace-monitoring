package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func FirestoreInit(credJson string, projectID string) *firestore.Client {
	// Use a service account
	ctx := context.Background()
	opt := option.WithCredentialsFile(credJson)
	client, err := firestore.NewClient(ctx, "projectID", opt)
	checkErr(err)

	return client
}
