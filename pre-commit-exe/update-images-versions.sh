#!/bin/bash

if [[ -z `git status --porcelain | grep A | grep terraform/image-versions.auto.tfvars` ]]; then
    echo -n > ./terraform/image-versions.auto.tfvars

    for dockerfile in `find . -iname Dockerfile`; do
        directory=${dockerfile%*Dockerfile}
        service=${directory%*/}
        service=${service##*/}
        if [[ "$service" != ".devcontainer" ]]; then
            commit=`git log -n 1 --pretty=format:%H -- ${directory}`
            echo "${service}_tag = \"${commit}\"" >> ./terraform/image-versions.auto.tfvars
        fi
    done

    terraform -chdir=terraform fmt > /dev/null
    
    exit 1
else
    exit 0
fi
