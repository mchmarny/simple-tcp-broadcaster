# Go parameters
BIN_NAME=simple
LOCAL_PORT=5050

all: test

build:
	go build -v -o ./bin/$(BIN_NAME)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(BIN_NAME)-linux

build-docker:
	docker build -t $(BIN_NAME) .

test:
	go test -v ./...

run-server:
	./bin/$(BIN_NAME) server start --port $(LOCAL_PORT)

run-client:
	./bin/$(BIN_NAME) client connect --port $(LOCAL_PORT)

clean:
	go clean
	rm -f ./bin/$(BIN_NAME)

deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure