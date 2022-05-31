package main

import (
	"context"

	"cloud.google.com/go/firestore"
)

func FirestoreInit(projectID string) *firestore.Client {
	// Use JSON credentials
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	checkErr(err)

	return client
}
