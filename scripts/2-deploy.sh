#!/bin/bash
set -eo pipefail
ARTIFACT_BUCKET=$(cat bucket-name.txt)


$(cd .. && GOOS=linux go build -o ./scripts/build/lambda cmd/lambda/main.go)

aws cloudformation package --template-file template.yml --s3-bucket $ARTIFACT_BUCKET --output-template-file out.yml

export $(grep -v '^#' ../.env | xargs)
cat out.yml | envsubst > sub_out.yml
mv sub_out.yml out.yml

aws cloudformation deploy --region us-east-1 --template-file out.yml --stack-name homefinder --capabilities CAPABILITY_NAMED_IAM