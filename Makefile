SHELL := /bin/sh
BINARY := quintaedizione-mcp

.PHONY: build test format vet install clean

.DEFAULT_GOAL := build

build: format vet
	go build -o $(BINARY) .

test: format vet
	go test -race -v ./...

format:
	go fmt ./...

vet:
	go vet ./...

install:
	go install .

clean:
	rm -f $(BINARY)
