#!/bin/bash

# ./create-service-accounts.sh project_id

gcloud iam service-accounts create aircraft-list --project $1
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:aircraft-list@$1.iam.gserviceaccount.com" \
    --role="roles/datastore.viewer"
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:aircraft-list@$1.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"

gcloud iam service-accounts create airspace-daily-history --project $1
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:airspace-daily-history@$1.iam.gserviceaccount.com" \
    --role="roles/datastore.viewer"
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:airspace-daily-history@$1.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"

gcloud iam service-accounts create airspace-monthly-history --project $1
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:airspace-monthly-history@$1.iam.gserviceaccount.com" \
    --role="roles/datastore.viewer"
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:airspace-monthly-history@$1.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"

gcloud iam service-accounts create airspace-history-calculator --project $1
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:airspace-history-calculator@$1.iam.gserviceaccount.com" \
    --role="roles/datastore.user"
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:airspace-history-calculator@$1.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"

gcloud iam service-accounts create flink-sa --project $1
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:flink-sa@$1.iam.gserviceaccount.com" \
    --role="roles/datastore.user"
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:flink-sa@$1.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:flink-sa@$1.iam.gserviceaccount.com" \
    --role="roles/pubsub.publisher"
