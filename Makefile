.DEFAULT_GOAL := build

fmt: 
	@go fmt ./...

lint: fmt
	@golint ./...

vet: fmt
	@go vet ./...

build: vet
	@go build -o bin/main main.go

clean:
	@go clean
	rm bin/main

test:
	@go test ./...

run:
	@go run main.go dns.google.com

.PHONY: fmt lint vet build clean
