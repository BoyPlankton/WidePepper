.PHONY: help build test clean run install lint fmt vet coverage bench

# Variables
BINARY_NAME=widepepper
GO=go
GOFLAGS=-v
BINARY_PATH=./bin/$(BINARY_NAME)

# Default target
help:
	@echo "WidePepper Language - Build & Development Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  help        - Display this help message"
	@echo "  build       - Build the WidePepper compiler"
	@echo "  install     - Install WidePepper to $$GOPATH/bin"
	@echo "  test        - Run all tests with verbose output"
	@echo "  test-cover  - Run tests and generate coverage report"
	@echo "  test-quick  - Run tests without verbose output"
	@echo "  coverage    - Display test coverage percentage"
	@echo "  bench       - Run benchmarks"
	@echo "  clean       - Remove build artifacts and coverage files"
	@echo "  run-example - Run an example script (EXAMPLE=01_hello_world)"
	@echo "  run-tokens  - Run lexer tokens on a script (EXAMPLE=01_hello_world)"
	@echo "  lint        - Run linter (if golangci-lint installed)"
	@echo "  fmt         - Format all Go source files"
	@echo "  vet         - Run go vet for code issues"
	@echo "  deps        - Download and verify dependencies"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make run-example EXAMPLE=02_variables"
	@echo "  make test-cover"

# Build the compiler
build: clean
	@echo "Building WidePepper compiler..."
	mkdir -p bin
	$(GO) build $(GOFLAGS) -o $(BINARY_PATH) .
	@echo "Build complete: $(BINARY_PATH)"

# Install to local bin directory
install: build
	@echo "Installing WidePepper to $(BINARY_PATH)..."
	@echo "Installation complete (binary at $(BINARY_PATH))"

# Run all tests
test:
	@echo "Running all tests..."
	$(GO) test $(GOFLAGS) ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test -cover $(GOFLAGS) ./...

# Quick test run (quiet)
test-quick:
	@echo "Running tests..."
	$(GO) test ./...

# Generate detailed coverage report
coverage:
	@echo "Generating coverage report..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "For HTML coverage report, run:"
	@echo "  go tool cover -html=coverage.out -o coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	$(GO) clean
	@echo "Clean complete"

# Run an example script
run-example: build
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "Usage: make run-example EXAMPLE=<example_name>"; \
		echo "Available examples:"; \
		ls -1 examples/*.wp | sed 's/.*\///g' | sed 's/\.wp//g'; \
	else \
		if [ -f "examples/$(EXAMPLE).wp" ]; then \
			echo "Running example: $(EXAMPLE).wp"; \
			$(BINARY_PATH) examples/$(EXAMPLE).wp; \
		else \
			echo "Error: examples/$(EXAMPLE).wp not found"; \
			exit 1; \
		fi; \
	fi

# Show tokens for an example script
run-tokens: build
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "Usage: make run-tokens EXAMPLE=<example_name>"; \
		echo "Available examples:"; \
		ls -1 examples/*.wp | sed 's/.*\///g' | sed 's/\.wp//g'; \
	else \
		if [ -f "examples/$(EXAMPLE).wp" ]; then \
			echo "Tokenizing: $(EXAMPLE).wp"; \
			$(BINARY_PATH) -tokens examples/$(EXAMPLE).wp; \
		else \
			echo "Error: examples/$(EXAMPLE).wp not found"; \
			exit 1; \
		fi; \
	fi

# Format Go source code
fmt:
	@echo "Formatting Go source code..."
	$(GO) fmt ./...
	@echo "Format complete"

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "Vet complete"

# Run linter (if installed)
lint:
	@which golangci-lint > /dev/null 2>&1 || { echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	@echo "Running linter..."
	golangci-lint run ./...

# Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify
	@echo "Dependencies verified"

# Development build (faster, with debug info)
dev-build:
	@echo "Building development version..."
	mkdir -p bin
	$(GO) build -o $(BINARY_PATH) .
	@echo "Development build complete"

# Watch and rebuild on changes (requires entr or similar)
watch:
	@which entr > /dev/null 2>&1 || { echo "entr not installed. Install with: brew install entr"; exit 1; }
	@find . -name "*.go" | entr -r make dev-build

# Run a script from stdin
stdin-test:
	@echo "Enter WidePepper script (Ctrl+D when done):"
	cat | $(BINARY_PATH) -stdin

# Display version
version: build
	@$(BINARY_PATH) -version
