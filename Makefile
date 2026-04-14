.PHONY: help build test lint fmt tidy install clean tag snapshot release

BIN        := kiart
PKG        := ./cmd/kiart
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS    := -s -w -X main.version=$(VERSION)

help: ## Show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build the binary
	go build -ldflags "$(LDFLAGS)" -o bin/$(BIN) $(PKG)

test: ## Run tests
	go test ./...

lint: ## Run linter
	golangci-lint run

fmt: ## Format Go code
	gofmt -s -w $(shell go list -f '{{.Dir}}' ./...)

tidy: ## Tidy dependencies
	go mod tidy

install: ## Install binary to GOPATH
	go install -ldflags "$(LDFLAGS)" $(PKG)

clean: ## Remove build artifacts
	rm -rf bin/ dist/

snapshot: ## Build snapshot release (no tag required)
	goreleaser release --snapshot --clean

release: ## Publish release for a tag (tag=v1.0.0)
	@if [ -z "$(tag)" ]; then echo "Error: Please specify a tag, e.g., make release tag=v1.0.0"; exit 1; fi
	@if ! git rev-parse "$(tag)" >/dev/null 2>&1; then \
		git tag -a $(tag) -m "Release $(tag)"; \
		git push origin $(tag); \
	fi
	git checkout $(tag)
	goreleaser release --clean
	git checkout -
