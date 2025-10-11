# SCache Makefile

.PHONY: test test-verbose benchmark clean build examples lint format help

# 默认目标
all: test benchmark

# 运行测试
test:
	go test ./...

# 运行详细测试
test-verbose:
	go test -v ./...

# 运行基准测试
benchmark:
	go test -bench=. -benchmem ./...

# 运行指定包的基准测试
benchmark-cache:
	go test -bench=. -benchmem ./cache

# 运行并发测试
test-concurrent:
	cd examples/concurrent && go run main.go

# 运行基础示例
run-example:
	cd examples/basic && go run main.go

# 启动 Web 服务示例
run-webserver:
	cd examples/webserver && go run main.go

# 清理构建文件
clean:
	go clean -testcache
	rm -f coverage.out

# 生成测试覆盖率报告
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码格式化
format:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 安装依赖
deps:
	go mod tidy
	go mod download

# 构建项目
build:
	go build ./...

# 显示帮助信息
help:
	@echo "可用的命令:"
	@echo "  test              - 运行所有测试"
	@echo "  test-verbose      - 运行详细测试"
	@echo "  benchmark         - 运行基准测试"
	@echo "  benchmark-cache   - 运行缓存包基准测试"
	@echo "  test-concurrent   - 运行并发测试"
	@echo "  run-example       - 运行基础示例"
	@echo "  run-webserver     - 启动 Web 服务示例"
	@echo "  clean             - 清理构建文件"
	@echo "  coverage          - 生成测试覆盖率报告"
	@echo "  format            - 格式化代码"
	@echo "  lint              - 代码检查"
	@echo "  deps              - 安装依赖"
	@echo "  build             - 构建项目"
	@echo "  help              - 显示此帮助信息"