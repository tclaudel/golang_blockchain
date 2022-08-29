BIN=bin

all: run-server

build:

run-server: deps build-server
	./$(BIN)/server

build-server: deps
	go build -o $(BIN)/server cmd/server/main.go

deps:
	go mod tidy

.PHONY: build run-server build-server deps