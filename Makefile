deploy:
	cd ./scripts && ./2-deploy.sh

run-offline:
	go run ./cmd/native

test-lambda:
	cd ./scripts && ./3-invoke.sh