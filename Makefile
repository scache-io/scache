# Makefile for SCache Project

.PHONY: help build test clean coverage benchmark lint fmt vet deps example

# 变量定义
BINARY_NAME=scache
MAIN_PACKAGE=.
PACKAGES=$(shell go list ./...)
TEST_PACKAGES=$(shell go list ./...)

# 默认目标
help: ## 显示帮助信息
	@echo "SCache Makefile 命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 构建目标
build: ## 构建项目
	@echo "🔨 构建项目..."
	go build -v $(MAIN_PACKAGE)

# 测试目标
test: ## 运行所有测试
	@echo "🧪 运行测试..."
	go test -v $(TEST_PACKAGES)

# 快速测试
test-quick: ## 快速运行测试（跳过基准测试）
	@echo "⚡ 快速测试..."
	go test -short -v $(TEST_PACKAGES)

# 覆盖率测试
coverage: ## 生成测试覆盖率报告
	@echo "📊 生成覆盖率报告..."
	go test -coverprofile=coverage.out $(TEST_PACKAGES)
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告已生成: coverage.html"

# 基准测试
benchmark: ## 运行基准测试
	@echo "🏃 运行基准测试..."
	go test -bench=. -benchmem $(TEST_PACKAGES)

# 代码检查
lint: ## 运行代码检查
	@echo "🔍 运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过代码检查"; \
		echo "   安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 代码格式化
fmt: ## 格式化代码
	@echo "🎨 格式化代码..."
	go fmt $(PACKAGES)

# 代码静态分析
vet: ## 运行 go vet
	@echo "🔬 运行 go vet..."
	go vet $(PACKAGES)

# 依赖管理
deps: ## 下载和验证依赖
	@echo "📦 管理依赖..."
	go mod download
	go mod verify
	go mod tidy

# 清理目标
clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	rm -f coverage.out coverage.html
	rm -f $(BINARY_NAME)
	go clean -cache

# 运行示例
example: ## 运行所有示例
	@echo "🚀 运行示例程序..."
	@echo "1. 基本使用示例:"
	cd examples/basic && go run basic_usage.go
	@echo ""
	@echo "2. 全局缓存示例:"
	cd examples/global && go run global_usage.go

# 并发测试示例
example-concurrent: ## 运行并发测试示例
	@echo "🏃‍♂️ 运行并发测试示例..."
	cd examples/concurrent && go run concurrent_test.go

# Web服务示例
example-web: ## 运行Web服务示例
	@echo "🌐 启动Web服务示例..."
	cd examples/webserver && go run webserver.go

# 完整检查
check: fmt vet lint test ## 运行完整代码检查
	@echo "✅ 代码检查完成"

# CI/CD 流水线
ci: deps fmt vet test benchmark coverage ## 模拟CI/CD流水线
	@echo "🚀 CI/CD 流水线完成"

# 安装开发工具
install-tools: ## 安装开发工具
	@echo "🔧 安装开发工具..."
	@echo "安装 golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 生成文档
docs: ## 生成文档
	@echo "📚 生成文档..."
	godoc -http=:6060 &
	@echo "📖 文档服务已启动: http://localhost:6060/pkg/scache/"

# 项目统计
stats: ## 显示项目统计信息
	@echo "📊 项目统计信息:"
	@echo "Go 文件数量: $(shell find . -name '*.go' | wc -l)"
	@echo "代码行数: $(shell find . -name '*.go' -exec wc -l {} + | tail -1)"
	@echo "测试文件数量: $(shell find . -name '*_test.go' | wc -l)"
	@echo "测试行数: $(shell find . -name '*_test.go' -exec wc -l {} + | tail -1)"

# 版本信息
version: ## 显示版本信息
	@echo "📋 版本信息:"
	@echo "Go 版本: $(shell go version)"
	@echo "模块信息:"
	@go list -m

# 监听文件变化并运行测试
watch: ## 监听文件变化并运行测试
	@echo "👀 监听文件变化..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} make test-quick; \
	else \
		echo "⚠️  fswatch 未安装，跳过监听"; \
		echo "   安装命令: brew install fswatch (macOS)"; \
	fi

# 发布准备
release-prep: clean fmt vet test benchmark ## 准备发布
	@echo "🎯 准备发布..."
	@echo "✅ 发布准备完成"

# 显示项目信息
info: ## 显示项目信息
	@echo "📋 SCache 项目信息:"
	@echo "=================="
	@echo "描述: 高性能 Go 语言缓存库"
	@echo "作者: Your Name"
	@echo "版本: v1.0.0"
	@echo "Go 版本要求: >= 1.18"
	@echo "=================="