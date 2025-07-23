# Makefile for video-rental-api

BINARY=video-rental-api
MAIN=cmd/server/main.go
OUTPUT=bin/$(BINARY)

.PHONY: all build build-linux run test clean integration-test

all: build

## Build for current system
build:
	go build -o $(OUTPUT) $(MAIN)

## Cross-compile for Linux (e.g. for Docker deployment)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)-linux $(MAIN)

## Run the application locally
run:
	go run $(MAIN)

## Run unit tests
test:
	go test ./internal/... -v

## Run integration tests
integration-test:
	bash test/integration-base.sh
	bash test/integration-rental.sh

## Clean generated binaries
clean:
	rm -f $(OUTPUT)*
