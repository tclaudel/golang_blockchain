BIN=bin

all: build-server build-cli

build:

run-server: deps build-server
	./$(BIN)/server

build-server: deps
	go build -o $(BIN)/server cmd/server/main.go

build-cli: deps
	go build -o $(BIN)/cli cmd/cli/main.go

server: build-server
	./$(BIN)/server

deps:
	go mod tidy

reset:
	rm -Rf ./data/wallet/*
	rm -Rf ./data/blockchain/*

.PHONY: build run-server build-server deps