.PHONY: setup build run clean lint format test

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)

# Default target
all: build

# Install development tools
setup:
	@echo "Installing development tools..."
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

# Build the application
build:
	@echo "Building mastodon-to-markdown..."
	@go build -ldflags "$(LDFLAGS)" -o mastodon-to-markdown .
	@echo "✅ Built: mastodon-to-markdown"

# Run the application
run: build
	./mastodon-to-markdown

# Clean build artifacts
clean:
	@rm -f mastodon-to-markdown
	@echo "✅ Cleaned"

# Lint code
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "❌ golangci-lint not found. Install with: make setup"; \
		exit 1; \
	}
	@echo "Running linters..."
	@golangci-lint run --timeout 5m
	@echo "✅ Lint complete"

# Format code
format:
	@go fmt ./...
	@command -v gofumpt >/dev/null 2>&1 || { \
		echo "❌ gofumpt not found. Install with: make setup"; \
		exit 1; \
	}
	@gofumpt -l -w .
	@echo "✅ Format complete"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...
	@echo "✅ Tests complete"
