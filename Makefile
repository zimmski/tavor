.PHONY: all binaries clean fmt install lint test tools

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: clean install test install binaries

binaries:
	go install $(ROOT_DIR)/bin/tavor.go
clean:
	go clean -i ./...
coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
debugbinaries:
	go install -race $(ROOT_DIR)/bin/tavor.go
dependencies:
	go get -d -v -u ./...
fmt:
	gofmt -l -w $(ROOT_DIR)/
install:
	go install ./...
	go install -race ./...
lint: clean install
	go tool vet -all=true -v=true $(ROOT_DIR)/
	golint $(ROOT_DIR)/
test: clean
	go test -race ./...
tools:
	go get -u code.google.com/p/go.tools/cmd/cover
	go get -u code.google.com/p/go.tools/cmd/godoc
	go get -u code.google.com/p/go.tools/cmd/vet
	go get -u github.com/golang/lint
