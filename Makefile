.PHONY: all generate check-generate fmt vet test build run clean

BINARY := bin/server
APP_PKG := ./cmd/server

all: generate fmt vet test build

check-generate: generate
	@git diff --exit-code -- localizer/message/all_gen.go

fmt:
	@gofmt -l -w .

build: generate
	@go build -o $(BINARY) $(APP_PKG)

run: generate
	@go run $(APP_PKG)

clean:
	@rm -rf bin

vet:
	@go vet ./...

generate:
	@go generate ./...
