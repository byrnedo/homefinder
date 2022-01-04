.EXPORT_ALL_VARIABLES:
AWS_PROFILE = own

deploy:
	cd ./scripts && ./2-deploy.sh

run-offline:
	go run ./cmd/native

test-lambda:
	cd ./scripts && ./3-invoke.sh
logs:
	aws logs tail /aws/lambda/homefinder-function-nENJWb1mYMzj --region us-east-1