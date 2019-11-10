# SOURCEDIR=.
# SOURCES = $(shell find $(SOURCEDIR) -name '*.go')
# VERSION=$(git describe --always --tags)
# BINARY=bin/pd

# bin: $(BINARY)

# $(BINARY): $(SOURCES)
# 	go build -o $(BINARY) command/*

.PHONY: build

GOPATH?=$(shell go env GOPATH)
GO111MODULE=auto

build: build-deps
	go build -mod=vendor -o pd ./command
build-deps:
	go get
	go mod verify
	go mod vendor

install: build
	cp pd $(GOPATH)/bin

.PHONY: test
test:
	go test ./...

deploy:
	- curl -sL https://git.io/goreleaser | bash

