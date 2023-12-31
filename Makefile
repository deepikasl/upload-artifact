SHELL := /bin/bash
GOCMD = go
install:
	go mod download

tidy:
	go mod tidy

run:
	go run main.go $(ARGS)

test:
	go test -v --cover ./...

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(ARGS)

build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(ARGS)

build-all:
	GOOS=linux GOARCH=amd64 go build $(ARGS) -o bin/upload-artifact-Linux-x86_64
	GOOS=linux GOARCH=arm64 go build $(ARGS) -o bin/upload-artifact-Linux-ARM64
	GOOS=darwin GOARCH=amd64 go build $(ARGS) -o bin/upload-artifact-Darwin-x86_64
	GOOS=darwin GOARCH=arm64 go build $(ARGS) -o bin/upload-artifact-Darwin-ARM64
	GOOS=windows GOARCH=amd64 go build $(ARGS) -o bin/upload-artifact-Windows-x86_64

do-all: install test build-all
