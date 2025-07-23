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
	export PORT=8080
	export API_KEY="secure-dev-key-123"
	go run $(MAIN)

## Run unit tests
test:
	go test ./internal/... -v

## Run integration tests
integration-test:	
	@echo "\n# Running Customer Tests"
	python3 test/customer.py
	@echo "# Finished Customer Tests...\n"

	@echo "\n# Running Rental Tests..."
	python3 test/rental.py
	@echo "# Finished Rental Tests...\n"

	@echo "\n# Running Inventory Tests..."
	python3 test/inventory.py
	@echo "# Finished Inventory Tests...\n"

	@echo "\n# Running Store Tests..."
	python3 test/store.py
	@echo "# Finished Store Tests...\n"

	@echo "\n# Running Film Tests..."
	python3 test/film.py
	@echo "# Finished Film Tests...\n"

## Clean generated binaries
clean:
	rm -f $(OUTPUT)*
