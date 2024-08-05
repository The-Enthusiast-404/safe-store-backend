.PHONY: build run

build:
    go build -o bin/api cmd/api/main.go

run: build
    ./bin/api
