#!/bin/bash
BUCKET_NAME=byrnedo-homefinder
echo $BUCKET_NAME > bucket-name.txt
aws s3 mb --region us-east-1 s3://$BUCKET_NAME