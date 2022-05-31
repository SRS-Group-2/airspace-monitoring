#!/bin/bash

# ./delete-service-accounts.sh project_id

gcloud iam service-accounts delete aircraft-info@$1.iam.gserviceaccount.com --project $1 --quiet

gcloud iam service-accounts delete aircraft-list@$1.iam.gserviceaccount.com --project $1 --quiet

gcloud iam service-accounts delete aircraft-positions@$1.iam.gserviceaccount.com --project $1 --quiet

gcloud iam service-accounts delete airspace-daily-history@$1.iam.gserviceaccount.com --project $1 --quiet

gcloud iam service-accounts delete airspace-monthly-history@$1.iam.gserviceaccount.com --project $1 --quiet

gcloud iam service-accounts delete airspace-history-calculator@$1.iam.gserviceaccount.com --project $1 --quiet

gcloud iam service-accounts delete flink-sa@$1.iam.gserviceaccount.com --project $1 --quiet
