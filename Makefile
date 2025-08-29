# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Binary names
BINARY_NAME=task-scheduler
BINARY_UNIX=$(BINARY_NAME)_unix

# Build directory
BUILD_DIR=bin

.PHONY: all build clean test coverage deps fmt vet lint run help

# Default target
all: clean deps fmt vet test build

# Build the binary
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./cmd/scheduler

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Vet code
vet:
	@echo "Vetting code..."
	$(GOVET) ./...

# Run the application
run: build
	@echo "Running application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Cross compilation
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v ./cmd/scheduler

# Install dependencies for development
dev-deps:
	@echo "Installing development dependencies..."
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Run 'make dev-deps' first." && exit 1)
	golangci-lint run

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./pkg/... > docs/api.md

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t task-scheduler:latest .

# Show help
help:
	@echo "Available targets:"
	@echo "  all        - Clean, download deps, format, vet, test, and build"
	@echo "  build      - Build the binary"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  coverage   - Run tests with coverage report"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  fmt        - Format code"
	@echo "  vet        - Vet code"
	@echo "  lint       - Lint code (requires golangci-lint)"
	@echo "  run        - Build and run the application"
	@echo "  build-linux- Cross compile for Linux"
	@echo "  dev-deps   - Install development dependencies"
	@echo "  docs       - Generate documentation"
	@echo "  docker-build- Build Docker image"
	@echo "  help       - Show this help message"