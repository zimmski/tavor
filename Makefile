.PHONY: all binaries clean fmt install lint test tools

all: clean install test install binaries

binaries:
	go install $(GOPATH)/src/github.com/zimmski/tavor/bin/tavor.go
clean:
	go clean -i ./...
coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
debugbinaries:
	go install -race $(GOPATH)/src/github.com/zimmski/tavor/bin/tavor.go
fmt:
	gofmt -l -w $(GOPATH)/src/github.com/zimmski/tavor
install:
	go install ./...
	go install -race ./...
lint: clean install
	go tool vet -all=true -v=true $(GOPATH)/src/github.com/zimmski/tavor
	golint $(GOPATH)/src/github.com/zimmski/tavor
test: clean
	go test -race ./...
tools:
	go get -u code.google.com/p/go.tools/cmd/cover
	go get -u code.google.com/p/go.tools/cmd/godoc
	go get -u code.google.com/p/go.tools/cmd/vet
	go get -u github.com/golang/lint
