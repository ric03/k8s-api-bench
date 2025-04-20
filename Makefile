# Makefile for k8s-api-bench

# Variables
BINARY_NAME=k8s-api-bench
BUILD_DIR=build
GO=go
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
.PHONY: all
all: clean build

# Clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)

# Build for current platform
.PHONY: build
build:
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for all platforms
.PHONY: build-all
build-all: build-linux-amd64 build-windows-amd64

# Build for Linux AMD64 (64-bit)
.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 .


# Build for Windows AMD64 (64-bit)
.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe .


# Run the application
.PHONY: run
run:
	$(GO) run .

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all            - Clean and build for current platform"
	@echo "  clean          - Remove build directory"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms (Linux and Windows, 64-bit)"
	@echo "  build-linux-amd64  - Build for Linux 64-bit"
	@echo "  build-windows-amd64 - Build for Windows 64-bit"
	@echo "  run            - Run the application"
	@echo "  help           - Show this help message"
