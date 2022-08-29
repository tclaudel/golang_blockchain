BIN=bin

all: run-server

build:

run-server: deps build-server
	./$(BIN)/server

build-server: deps
	/home/tclaudel/go/go1.18/bin/go build -o $(BIN)/server cmd/server/main.go

deps:
	/home/tclaudel/go/go1.18/bin/go mod tidy

.PHONY: build run-server build-server deps