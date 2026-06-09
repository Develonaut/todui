.DEFAULT_GOAL := build
BIN := todui
PKG := ./cmd/todui
GOBIN := $(shell go env GOPATH)/bin

.PHONY: build test lint vet fmt install run tidy clean

build:
	go build ./...

test:
	go test -race ./...

vet:
	go vet ./...

lint: vet
	golangci-lint run

fmt:
	gofmt -w .

tidy:
	go mod tidy

install:
	go install $(PKG)
	ln -sf $(GOBIN)/$(BIN) $(GOBIN)/td

run:
	go run $(PKG) tui

clean:
	rm -f $(BIN) td
	rm -rf dist
