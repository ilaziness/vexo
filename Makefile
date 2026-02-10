# Makefile for vexo SSH desktop application
VERSION ?= v1.0.0
MODE ?= release

# 获取 Git 信息（完整 hash）
GIT_INFO := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
# RFC3339 格式，带时区偏移 (如 2026-02-10T18:56:45+08:00)
BUILD_TIME := $(shell date "+%Y-%m-%dT%H:%M:%S%:z" 2>/dev/null || echo "unknown")

CONFIG_FILE := services/config_service.go
CONFIG_BACKUP := services/config_service.go.bak

.PHONY: build-windows build-darwin build-linux clean all build-mac build-mac-intel

# Default target
all: build-windows build-darwin build-mac-intel build-linux

# 替换配置文件中的变量
replace-config:
	@echo "Replacing build variables in $(CONFIG_FILE)..."
	@cp $(CONFIG_FILE) $(CONFIG_BACKUP)
	@sed -i 's/Mode        = "debug"/Mode        = "$(MODE)"/' $(CONFIG_FILE)
	@sed -i 's/Version   = "v1.0.0"/Version   = "$(VERSION)"/' $(CONFIG_FILE)
	@sed -i 's|GitInfo   = ".*"|GitInfo   = "$(GIT_INFO)"|' $(CONFIG_FILE)
	@sed -i 's|BuildTime = ".*"|BuildTime = "$(BUILD_TIME)"|' $(CONFIG_FILE)

# 恢复配置文件
restore-config:
	@echo "Restoring $(CONFIG_FILE)..."
	@if [ -f $(CONFIG_BACKUP) ]; then \
		mv $(CONFIG_BACKUP) $(CONFIG_FILE); \
	fi

# Build for Windows
build-windows: replace-config
	@echo "Building for Windows..."
	@wails3 build GOOS=windows VERSION=$(VERSION) MODE=$(MODE) || true
	@$(MAKE) restore-config

# Build for macOS ARM64 (Darwin)
build-mac: replace-config
	@echo "Building for macOS (Darwin)..."
	@wails3 build GOOS=darwin VERSION=$(VERSION) MODE=$(MODE) || true
	@$(MAKE) restore-config

# Build for macOS Intel (Darwin)
build-mac-intel: replace-config
	@echo "Building for macOS Intel (Darwin)..."
	@wails3 build GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) MODE=$(MODE) || true
	@$(MAKE) restore-config

# Build for Linux
build-linux: replace-config
	@echo "Building for Linux..."
	@wails3 build GOOS=linux VERSION=$(VERSION) MODE=$(MODE) || true
	@$(MAKE) restore-config

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@if [ -f $(CONFIG_BACKUP) ]; then \
		rm $(CONFIG_BACKUP); \
	fi
