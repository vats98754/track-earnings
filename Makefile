SHELL := /bin/zsh

APP := track-earnings
BIN := bin/$(APP)

.PHONY: build run tidy test lint fmt mod dev

build:
	go build -o $(BIN) ./cmd/server

run: build
	./$(BIN)

fmt:
	gofmt -s -w .

mod:
	go mod tidy -v

test:
	go test ./...

dev:
	ADDR=:8080 ENV=dev go run ./cmd/server
