# Variables
BINARY_NAME=kiart
MAIN_FILE=main.go
BUILD_DIR=bin

.PHONY: all help run fmt tidy install

all: tidy fmt test build

.PHONY: help
help: ## Show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build the binary
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/$(BINARY_NAME)/$(MAIN_FILE)

run: build ## Run the binary
	@./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

test: ## Run Go tests
	@echo "Running tests..."
	@go test ./... -v

clean: ## Clean up built binaries
	@echo "Cleaning up..."
	@go clean
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Clean complete!"

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

tidy: ## Ensure dependencies are tidy
	@echo "Tidy dependencies..."
	@go mod tidy

install: build ## Install the binary to your GOPATH
	@echo "Installing $(BINARY_NAME) to GOPATH..."
	@go install
	@echo "Install complete!"
