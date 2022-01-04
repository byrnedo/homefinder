.EXPORT_ALL_VARIABLES:
AWS_PROFILE = own

deploy:
	cd ./scripts && ./2-deploy.sh

run-offline:
	go run ./cmd/native

test-lambda:
	cd ./scripts && ./3-invoke.sh
logs:
	aws logs tail /aws/lambda/$(shell aws cloudformation describe-stack-resources  --region us-east-1 --stack-name homefinder  --logical-resource-id function |jq '.StackResources[0].PhysicalResourceId' -r) --region us-east-1