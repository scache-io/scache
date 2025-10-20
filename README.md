# SCache - Go 结构体缓存代码生成工具

## 项目概述

SCache 是一个智能的 Go 结构体缓存代码生成工具，自动扫描项目中的结构体并生成对应的缓存操作方法。项目采用模块化设计，代码结构清晰，便于维护和扩展。

## 项目结构

```
scache/
├── README.md               # 项目文档
├── LICENSE                 # 开源协议
├── go.mod                  # Go 模块定义
├── scache.go               # 主包入口
├── scache_test.go          # 主包测试
├── errors.go               # 错误定义
├── cmd/                    # 命令行工具
│   └── scache/
│       └── main.go         # CLI 入口
├── generator/              # 代码生成器
│   ├── generator.go        # 生成器核心逻辑
│   └── cache.tpl           # 代码模板
├── cache/                  # 缓存实现
│   ├── cache.go            # 缓存核心逻辑
│   └── cache_test.go       # 缓存测试
├── storage/                # 存储引擎
│   ├── engine.go           # 存储引擎接口
│   └── engine_test.go      # 存储引擎测试
├── policies/               # 缓存策略
│   └── lru/
│       ├── lru.go          # LRU 策略实现
│       └── lru_test.go     # LRU 策略测试
├── types/                  # 数据类型
│   ├── data_objects.go     # 数据对象定义
│   └── data_objects_test.go # 数据对象测试
├── interfaces/             # 接口定义
│   └── interfaces.go       # 核心接口
├── config/                 # 配置管理
│   ├── engine_config.go    # 引擎配置
│   └── engine_config_test.go # 配置测试
└── constants/              # 常量定义
    └── constants.go        # 项目常量
```

## 快速开始

### 安装

```bash
go install github.com/scache-io/scache/cmd/scache@latest
```

### 基本用法

```bash
# 在当前目录生成所有结构体的缓存代码
scache gen

# 指定目录生成
scache gen -dir ./models

# 只生成指定结构体
scache gen -structs User,Product
```

### 每个结构体生成的方法

- `Store{Struct}` - 存储结构体
- `Load{Struct}` - 加载结构体
- `MustStore{Struct}` - 存储（失败时panic）
- `MustLoad{Struct}` - 加载（失败时panic）
- `Store{Struct}WithKey` - 格式化key存储
- `Load{Struct}WithKey` - 格式化key加载

### 使用示例

```go
package main

import (
    "yourproject/scache"
    "yourproject/scache/user"
    "yourproject/scache/product"
)

func main() {
    // 创建用户
    u := User{Name: "张三", Age: 25}

    // 使用User包的方法存储
    user.StoreUser("user:1", &u, time.Hour)

    // 加载用户
    var loadedUser User
    user.LoadUser("user:1", &loadedUser)

    // 使用主包工具函数
    stats := scache.CacheStats()
    size := scache.CacheSize()
}
```

## 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-dir` | 项目目录路径 | `.` |
| `-package` | 包名 | 目录名 |
| `-exclude` | 排除的目录 | `vendor,node_modules,.git` |
| `-structs` | 指定结构体名 | 生成所有 |
| `-h` | 显示帮助信息 | - |

## 帮助信息

使用 `-h` 参数查看完整的帮助信息：

```bash
scache -h
```

## 项目特性

- ✅ **智能扫描** - 自动发现所有Go结构体
- ✅ **模块化设计** - 清晰的目录结构和职责分离
- ✅ **类型安全** - 支持指针和值类型
- ✅ **错误处理** - 完整的错误处理机制
- ✅ **高性能** - 基于高效的存储引擎
- ✅ **可配置** - 灵活的结构体选择和配置
- ✅ **易于测试** - 完整的单元测试覆盖
- ✅ **可扩展** - 支持自定义缓存策略

## 核心组件

### 1. 代码生成器 (generator/)
- 扫描Go源文件中的结构体定义
- 基于模板生成缓存操作代码
- 支持灵活的配置和过滤选项

### 2. 缓存核心 (cache/)
- 提供统一的缓存操作接口
- 支持多种存储后端
- 线程安全的并发操作

### 3. 存储引擎 (storage/)
- 抽象的存储引擎接口
- 支持内存、文件等多种存储方式
- 可插拔的引擎设计

### 4. 缓存策略 (policies/)
- LRU (Least Recently Used) 缓存淘汰策略
- 可扩展的策略接口
- 支持自定义策略实现

### 5. 数据类型 (types/)
- 标准化的数据对象定义
- 支持复杂的数据结构
- 类型安全的操作接口

## 开发指南

### 构建项目

```bash
# 克隆项目
git clone https://github.com/scache-io/scache.git
cd scache

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 构建工具
go build -o scache ./cmd/scache
```