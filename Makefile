.PHONY: all clean coverage debug-install dependencies fmt install lint test tools

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: clean install test

clean:
	go clean -i ./...
	go clean -i -race ./...
coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
debug-install: clean
	go install -race -v ./...
dependencies:
	go get -d -t -u -v ./...
	go build -v ./...
fmt:
	gofmt -l -w $(ROOT_DIR)/
install: clean
	go install -v ./...
lint: install
	go tool vet -all=true -v=true $(ROOT_DIR)/
	golint $(ROOT_DIR)/
test: clean
	go test -race ./...
tools:
	go get -u code.google.com/p/go.tools/cmd/cover
	go get -u code.google.com/p/go.tools/cmd/godoc
	go get -u code.google.com/p/go.tools/cmd/vet
	go get -u github.com/golang/lint
	go install github.com/golang/lint/golint
