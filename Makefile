BIN=bin

all: gen build-server build-cli deps

build: deps

run-server: deps build-server
	./$(BIN)/server

build-server: deps
	go build -o $(BIN)/server cmd/server/main.go

build-cli: deps
	go build -o $(BIN)/cli cmd/cli/main.go

server: build-server deps
	./$(BIN)/server

deps:
	mkdir -p ./data/blockchain ./data/wallet
	go mod tidy

reset:
	rm -Rf ./data/wallet/*
	rm -Rf ./data/blockchain/*
	$(MAKE) deps

gen: deps
	oapi-codegen -package rest -generate client,types ./docs/swagger.yaml > ./pkg/interfaces/http/rest/api.gen.go

.PHONY: build run-server build-server deps reset gen server build-cli