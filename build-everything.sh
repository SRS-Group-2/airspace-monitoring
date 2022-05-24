#!/bin/bash

# ./build-everything.sh project_id region docker_repo_name

home=`pwd`

for dockerfile in `find . -iname Dockerfile -not -path *.devcontainer*`; do
    directory=${dockerfile%/*}
    hash=`git log -n 1 --pretty=format:%H -- $directory`
    service=${directory##*/}
    cd $directory
    echo building and pushing $2-docker.pkg.dev/$1/$3/$service:$hash
    docker build -t $2-docker.pkg.dev/$1/$3/$service:$hash .
    docker push $2-docker.pkg.dev/$1/$3/$service:$hash
    cd $home
done
