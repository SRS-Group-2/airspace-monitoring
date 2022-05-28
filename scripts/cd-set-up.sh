#!/bin/bash

# ./cd-set-up.sh project_id

gcloud services enable iam.googleapis.com --project $1
gcloud services enable cloudresourcemanager.googleapis.com --project $1
gcloud services enable iamcredentials.googleapis.com --project $1
gcloud services enable sts.googleapis.com --project $1

gcloud iam service-accounts create terraform --project $1
gcloud projects add-iam-policy-binding $1 \
    --member="serviceAccount:terraform@$1.iam.gserviceaccount.com" \
    --role="roles/editor"

gcloud iam workload-identity-pools create "github-pool" \
  --project="$1" \
  --location="global" \
  --display-name="GitHub pool"

gcloud iam workload-identity-pools providers create-oidc "github-provider" \
  --project="$1" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --display-name="Github provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.ref=assertion.ref,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"

gcloud iam service-accounts add-iam-policy-binding "terraform@$1.iam.gserviceaccount.com" \
  --project="$1" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/$1/locations/global/workloadIdentityPools/github-pool/attribute.repository/SRS-Group-2/airspace-monitoring"
