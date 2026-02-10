# Makefile for vexo SSH desktop application
VERSION ?= v1.0.0
MODE ?= release

.PHONY: build-windows build-darwin build-linux clean

# Default target
all: build-windows build-darwin build-mac-intel build-linux

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@wails3 build GOOS=windows VERSION=$(VERSION) MODE=$(MODE) 
# Build for macOS ARM64 (Darwin)
build-mac:
	@echo "Building for macOS (Darwin)..."
	@wails3 build GOOS=darwin VERSION=$(VERSION) MODE=$(MODE)

# Build for macOS Intel (Darwin)
build-mac-intel:
	@echo "Building for macOS (Darwin)..."
	@wails3 build GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) MODE=$(MODE) 

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@wails3 build GOOS=linux VERSION=$(VERSION) MODE=$(MODE) 

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/