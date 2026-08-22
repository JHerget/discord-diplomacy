format:
	gofmt -w .

build: format
	mkdir -p bin
	go build -o bin/bot ./cmd/bot

start: build
	./bin/bot
