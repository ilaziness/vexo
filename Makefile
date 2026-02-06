# Makefile for vexo SSH desktop application

.PHONY: build-windows build-darwin build-linux clean

# Default target
all: build-windows build-darwin build-linux

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@wails3 build GOOS=windows

# Build for macOS (Darwin)
build-mac:
	@echo "Building for macOS (Darwin)..."
	@wails3 build GOOS=darwin

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@wails3 build GOOS=linux

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/