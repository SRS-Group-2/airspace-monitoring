package main

import (
	"context"

	"cloud.google.com/go/firestore"
)

func FirestoreInit(projectID string) *firestore.Client {
	// Authentication uses Workload Identity
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	checkErr(err)

	return client
}
