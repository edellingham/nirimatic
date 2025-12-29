# Nirimatic Makefile

BINARY_NAME=nirimatic
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean install run dev test lint deps

# Default target
all: build

# Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/nirimatic

# Build for release (smaller binary)
release:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/nirimatic
	upx --best --lzma $(BINARY_NAME) 2>/dev/null || true

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean

# Install to /usr/local/bin
install: build
	sudo install -Dm755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Uninstall
uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Run the application
run: build
	./$(BINARY_NAME)

# Run in development mode (with go run)
dev:
	go run ./cmd/nirimatic

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Update dependencies
update:
	go get -u ./...
	go mod tidy

# Show help
help:
	@echo "Nirimatic Makefile targets:"
	@echo "  build    - Build the binary"
	@echo "  release  - Build optimized release binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install to /usr/local/bin"
	@echo "  run      - Build and run"
	@echo "  dev      - Run with go run"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  update   - Update dependencies"
