package main

import (
	"context"

	"cloud.google.com/go/firestore"
)

func FirestoreInit(projectID string) *firestore.Client {
	// Use a service account
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	checkErr(err)

	return client
}
