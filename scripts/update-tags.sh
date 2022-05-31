#!/bin/bash

echo -n > ./terraform/image-versions.auto.tfvars

for dockerfile in `find . -iname Dockerfile -not -path *.devcontainer*`; do
    directory=${dockerfile%/*}
    service=${directory##*/}
    hash=`git log -n 1 --pretty=format:%H -- $directory`
    echo "$service"_tag = \"$hash\" >> ./terraform/image-versions.auto.tfvars
done

terraform -version 2> /dev/null > /dev/null
if [[ "$?" == "0" ]]; then
    terraform -chdir=terraform fmt
fi
