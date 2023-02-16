.PHONY:
.SILENT:
.DEFAULT_GOAL := help

include .env
PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
#GOPATH="$(GOBASE)/vendor:$(GOBASE)"

## migrate-up: goose up migrations
migrate-up:
	goose -dir=./migrations postgres $(DATABASE_URI) up

## migrate-down: goose down migrations
migrate-down:
	goose -dir=./migrations postgres $(DATABASE_URI) down

## test: run tests
test:
	go test -count=1 -cover=1 ./...

## test-integration: run integration tests
test-integration:
	go test -tags integration -count=1 ./...

## lint: run linter checks
lint:
	golangci-lint run

## bench: run benchmark tests
bench:
	go test -bench=1 ./...

## build-server: build server app binary
build-server:
	go build -o $(GOBIN)/server cmd/server/main.go

## build-client: build client desktop application binary

build-client:
	CGO_ENABLED=1 go build -o $(GOBIN)/client cmd/client/main.go


## compile: Compiling for every OS and Platform https://developer.fyne.io/started/cross-compiling
compile:
	echo "Compiling for every OS and Platform"
	CGO_ENABLED=1 $(GOPATH)/bin/fyne-cross linux -arch=* ./cmd/client/main.go


help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
