export AWS_REGION=us-west-2
export SQS_QUEUE_URL=https://sqs.us-west-2.amazonaws.com/620486971062/diplomacy-api-v1-events

format:
	gofmt -w .

build: format
	mkdir -p bin
	go build -o bin/bot ./cmd/bot

run: build
	source .env.local; ./bin/bot
