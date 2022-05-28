#!/bin/bash

gcloud iam service-accounts delete terraform@$1.iam.gserviceaccount.com --quiet

gcloud iam workload-identity-pools providers delete "github-provider" \
  --project="$1" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --quiet

gcloud iam workload-identity-pools delete "github-pool" \
  --project="$1" \
  --location="global" \
  --quiet

gcloud services disable cloudresourcemanager.googleapis.com --project $1 --force
gcloud services disable iamcredentials.googleapis.com --project $1 --force
gcloud services disable sts.googleapis.com --project $1 --force
