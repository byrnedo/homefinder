#!/bin/bash
set -eo pipefail
ARTIFACT_BUCKET=$(cat bucket-name.txt)

aws cloudformation package --template-file template.yml --s3-bucket $ARTIFACT_BUCKET --output-template-file out.yml

export $(grep -v '^#' ../.env | xargs)
cat out.yml | envsubst > sub_out.yml
mv sub_out.yml out.yml
