#!/bin/bash
set -eo pipefail
FUNCTION=$(aws cloudformation describe-stack-resource --region us-east-1 --stack-name homefinder --logical-resource-id function --query 'StackResourceDetail.PhysicalResourceId' --output text)
echo $FUNCTION
while true; do
  aws lambda invoke --function-name $FUNCTION --region us-east-1 --payload "$(base64 ./event.json)" out.json
  cat out.json
  echo ""
  sleep 2
done

