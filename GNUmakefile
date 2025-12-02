# Copyright (c) Materialize Inc.

default: build

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	$(shell go env GOPATH)/bin/golangci-lint run

generate:
	go generate ./...

fmt:
	gofmt -s -w .

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

# Run acceptance tests serially (-p 1) to avoid conflicts when multiple
# test packages hit the same Ory project simultaneously
testacc:
	TF_ACC=1 go test -v -cover -timeout 30m -p 1 ./...

# Install development tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set up git hooks
hooks:
	git config core.hooksPath .githooks
	@echo "Git hooks installed. Pre-commit checks will run automatically."

.PHONY: build install lint generate fmt test testacc tools hooks
