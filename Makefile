# Makefile for GTMLP - Go HTML Parsing Library
#
# Shell Aliases (add to ~/.bashrc, ~/.zshrc, etc.):
#   alias mb='make build'           # Build with checks and warnings (default)
#   alias mq='make build-quick'     # Quick build (silent)
#   alias mv='make build-verbose'   # Verbose build
#   alias mt='make test'            # Run tests
#   alias me='make example TYPE='   # Run examples

# Default target
all: build test

# Build the library (with checks and warnings by default)
build:
	@echo "Building gtmlp library..."
	@go build -v ./...
	@echo "Build successful!"

# Quick build (silent, no checks)
build-quick:
	@echo "Building gtmlp library (quick)..."
	@go build ./internal/... .
	@echo "Build successful!"

# Build with full verbose output
build-verbose:
	@echo "Building gtmlp library (full verbose)..."
	@go vet ./internal/... . || true
	@go build -v -x ./internal/... . 2>&1 | tail -20
	@echo "Build completed!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./internal/... . -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./internal/... . -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run examples
example:
	@echo "Running GTMLP examples..."
	@echo "Usage: make example TYPE=<type>"
	@echo "Types: xpath, json, builder, url, random-ua, pattern, chainable"
	@if [ -z "$(TYPE)" ]; then \
		echo "Error: TYPE parameter is required"; \
		echo "Example: make example TYPE=xpath"; \
	else \
		go run examples/basic/*.go -type=$(TYPE); \
	fi

# Run XPath example
example-xpath:
	@echo "Running XPath example..."
	@go run examples/basic/*.go -type=xpath

# Run JSON example
example-json:
	@echo "Running JSON example..."
	@go run examples/basic/*.go -type=json

# Run builder example
example-builder:
	@echo "Running builder example..."
	@go run examples/basic/*.go -type=builder

# Run URL fetch example
example-url:
	@echo "Running URL fetch example..."
	@go run examples/basic/*.go -type=url

# Run Random User-Agent example
example-random-ua:
	@echo "Running Random User-Agent example..."
	@go run examples/basic/*.go -type=random-ua

# Run pattern example
example-pattern:
	@echo "Running pattern extraction example..."
	@go run examples/basic/*.go -type=pattern

# Run chainable pattern example
example-chainable:
	@echo "Running chainable pattern API example..."
	@go run examples/basic/*.go -type=chainable

# Run all basic examples
examples: example-xpath example-json example-builder example-url example-random-ua example-pattern example-chainable
	@echo ""
	@echo "All basic examples completed!"

# Run all advanced examples
examples-advanced:
	@echo "Running all advanced examples..."
	@go run examples/advanced/*.go -type=all

# Run specific advanced example
example-advanced:
	@echo "Running advanced example..."
	@echo "Usage: make example-advanced TYPE=<type>"
	@echo "Types: 01-09 or basic-product-scraping, xpath-validation, etc."
	@if [ -z "$(TYPE)" ]; then \
		echo "Error: TYPE parameter is required"; \
		echo "Example: make example-advanced TYPE=01"; \
	else \
		go run examples/advanced/*.go -type=$(TYPE); \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./internal/... .

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./internal/... .; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Create release build
release:
	@echo "Creating release build..."
	@go mod tidy
	@echo "Release ready"

# Show help
help:
	@echo "GTMLP Makefile Targets:"
	@echo ""
	@echo "Build:"
	@echo "  build            - Build library (with vet checks + warnings)"
	@echo "  build-quick      - Quick build (silent, no checks)"
	@echo "  build-verbose    - Build with full verbose output"
	@echo ""
	@echo "Testing:"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo ""
	@echo "Basic Examples:"
	@echo "  example          - Run example (usage: make example TYPE=<type>)"
	@echo "  example-xpath    - Run XPath example"
	@echo "  example-json     - Run JSON conversion example"
	@echo "  example-builder  - Run builder/fluent API example"
	@echo "  example-url      - Run URL fetch example"
	@echo "  example-random-ua - Run Random User-Agent example"
	@echo "  example-pattern  - Run pattern extraction example"
	@echo "  example-chainable - Run chainable pattern API example"
	@echo "  examples         - Run all basic examples"
	@echo ""
	@echo "Advanced Examples:"
	@echo "  example-advanced - Run specific advanced example (usage: make example-advanced TYPE=<type>)"
	@echo "  examples-advanced - Run all advanced examples"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean            - Clean build artifacts"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  deps             - Install dependencies"
	@echo "  release          - Prepare release build"
	@echo "  help             - Show this help message"

.PHONY: all build build-quick build-verbose test test-coverage example example-xpath example-json example-builder example-url example-random-ua example-pattern example-chainable examples example-advanced examples-advanced clean fmt lint deps release help
