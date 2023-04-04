#!/bin/bash
set -eo pipefail

(cd .. && GOOS=linux GOARCH=amd64 go build -o ./scripts/build/lambda cmd/lambda/main.go)

./1-template.sh


aws cloudformation deploy --region us-east-1 --template-file out.yml --stack-name homefinder --capabilities CAPABILITY_NAMED_IAM