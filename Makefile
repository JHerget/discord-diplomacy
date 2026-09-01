format:
	gofmt -w .

build: format
	mkdir -p bin
	go build -o bin/bot ./cmd/bot

run: build
	source .env.local; ./bin/bot
