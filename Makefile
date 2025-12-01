# SCache Makefile
# Copyright 2024 SCache Authors

# 变量定义
BINARY_NAME=scache
CMD_PATH=cmd/scache
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(COMMIT)"

# 默认目标
.PHONY: all
all: clean test build

# 构建二进制文件
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./$(CMD_PATH)

# 构建所有平台的二进制文件
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./$(CMD_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./$(CMD_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./$(CMD_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_PATH)

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 整理依赖
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# 安装开发依赖
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# 安装到本地
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./$(CMD_PATH)

# 运行示例
.PHONY: example
example:
	@echo "Running example..."
	cd examples && go run main.go

# 显示帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Clean, test, and build"
	@echo "  build         - Build the binary for current platform"
	@echo "  build-all     - Build binaries for all platforms"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  install-tools - Install development tools"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install to local GOPATH/bin"
	@echo "  example       - Run example"
	@echo "  help          - Show this help message"