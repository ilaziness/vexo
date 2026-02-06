# Makefile for vexo SSH desktop application
VERSION ?= v1.0.0
MODE ?= release
GIT_INFO := $(shell git log -1 --format="%cd %H" --date=format:"%Y-%m-%d %H:%M:%S")
BUILD_TIME := $(shell date +"%F %X")

.PHONY: build-windows build-darwin build-linux clean

# Default target
all: build-windows build-darwin build-linux

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@wails3 build GOOS=windows VERSION=$(VERSION) MODE=$(MODE) GITINFO=$(GIT_INFO) BUILDTIME=$(BUILD_TIME)

# Build for macOS ARM64 (Darwin)
build-mac:
	@echo "Building for macOS (Darwin)..."
	@wails3 build GOOS=darwin VERSION=$(VERSION) MODE=$(MODE) GITINFO=$(GIT_INFO) BUILDTIME=$(BUILD_TIME)

# Build for macOS Intel (Darwin)
build-mac-intel:
	@echo "Building for macOS (Darwin)..."
	@wails3 build GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) MODE=$(MODE) GITINFO=$(GIT_INFO) BUILDTIME=$(BUILD_TIME)

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@wails3 build GOOS=linux VERSION=$(VERSION) MODE=$(MODE) GITINFO=$(GIT_INFO) BUILDTIME=$(BUILD_TIME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/