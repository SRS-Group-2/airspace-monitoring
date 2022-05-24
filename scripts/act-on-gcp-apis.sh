#!/bin/bash

gcloud services $1 apigateway.googleapis.com --project $2
gcloud services $1 artifactregistry.googleapis.com --project $2
gcloud services $1 autoscaling.googleapis.com --project $2
gcloud services $1 compute.googleapis.com --project $2
gcloud services $1 container.googleapis.com --project $2
gcloud services $1 containerfilesystem.googleapis.com --project $2
gcloud services $1 datastore.googleapis.com --project $2
gcloud services $1 firestore.googleapis.com --project $2
gcloud services $1 iam.googleapis.com --project $2
gcloud services $1 iamcredentials.googleapis.com --project $2
gcloud services $1 monitoring.googleapis.com --project $2
gcloud services $1 pubsub.googleapis.com --project $2
gcloud services $1 run.googleapis.com --project $2
gcloud services $1 storage-api.googleapis.com --project $2
gcloud services $1 storage-component.googleapis.com --project $2
gcloud services $1 storage.googleapis.com --project $2
