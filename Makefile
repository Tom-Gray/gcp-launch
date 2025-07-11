# Makefile for gcp-launch

# Variables
BINARY_NAME=gcp-launch
BIN_DIR=.bin
BUILD_DIR=$(BIN_DIR)
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(GO_FILES)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

# Run tests
.PHONY: test
test:
	go test ./...

# Install binary to system PATH
.PHONY: install
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build   - Build the binary to .bin/gcp-launch"
	@echo "  test    - Run all tests"
	@echo "  clean   - Remove build artifacts"
	@echo "  install - Install binary to /usr/local/bin"
	@echo "  help    - Show this help message"
