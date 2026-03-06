.PHONY: build test lint clean

BINARY_NAME=maao
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/maao

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)

install:
	go install ./cmd/maao
