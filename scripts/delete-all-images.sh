#!/bin/bash

# ./delete-all-images.sh project_id region docker_repo_name

for image in `gcloud artifacts docker images list $2-docker.pkg.dev/$1/$3`; do
    gcloud artifacts docker images delete $2-docker.pkg.dev/$1/$3/$image --delete-tags
done
