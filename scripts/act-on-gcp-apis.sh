#!/bin/bash

if [[ "$1" == "disable" ]]; then 
    FLAGS="--force"
else
    FLAGS=""
fi

gcloud services $1 artifactregistry.googleapis.com --project $2 ${FLAGS}
gcloud services $1 run.googleapis.com --project $2 ${FLAGS}
gcloud services $1 container.googleapis.com --project $2 ${FLAGS}
gcloud services $1 compute.googleapis.com --project $2 ${FLAGS}
gcloud services $1 autoscaling.googleapis.com --project $2 ${FLAGS}
gcloud services $1 pubsub.googleapis.com --project $2 ${FLAGS}
gcloud services $1 firestore.googleapis.com --project $2 ${FLAGS}
gcloud services $1 apigateway.googleapis.com --project $2 ${FLAGS}
