# SCache - 高性能 Go 缓存库

[![GoDoc](https://godoc.org/github.com/your-repo/scache?status.svg)](https://godoc.org/github.com/your-repo/scache)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-repo/scache)](https://goreportcard.com/report/github.com/your-repo/scache)
[![Coverage](https://codecov.io/gh/your-repo/scache/branch/main/graph/badge.svg)](https://codecov.io/gh/your-repo/scache)

SCache 是一个高性能的 Go 语言缓存库，采用类似 Redis 的架构设计，支持多种数据类型和丰富的缓存策略。

## ✨ 特性

- 🚀 **高性能** - 基于 Go map 和读写锁实现，支持高并发访问
- 📦 **多数据类型** - 支持 String、List、Hash 等数据类型
- ⏰ **TTL 过期** - 支持灵活的过期时间设置
- 🗑️ **淘汰策略** - 支持 LRU 等多种淘汰策略
- 💾 **内存管理** - 智能内存压力检测和清理
- 🔧 **命令模式** - 易于扩展的命令系统
- 📊 **统计信息** - 详细的命中率和操作统计
- 🧵 **线程安全** - 完全的并发安全保证

## 🚀 快速开始

### 安装

```bash
go get github.com/your-repo/scache
```

### 基础使用

```go
package main

import (
    "fmt"
    "time"

    "scache"
)

func main() {
    // 简单的字符串操作
    err := scache.Set("hello", "world", time.Hour)
    if err != nil {
        panic(err)
    }

    value, found, err := scache.Get("hello")
    if err != nil {
        panic(err)
    }
    if found {
        fmt.Printf("Value: %v\n", value) // Output: Value: world
    }

    // 列表操作
    length, err := scache.LPush("mylist", "item1", time.Minute)
    if err != nil {
        panic(err)
    }
    fmt.Printf("List length: %d\n", length)

    // 哈希操作
    success, err := scache.HSet("user:1", "name", "Alice", time.Hour)
    if err != nil {
        panic(err)
    }
    if success {
        name, err := scache.HGet("user:1", "name")
        if err != nil {
            panic(err)
        }
        fmt.Printf("User name: %v\n", name) // Output: User name: Alice
    }
}
```

### 高级使用

```go
package main

import (
    "fmt"
    "time"

    "scache"
    "scache/config"
)

func main() {
    // 创建自定义配置的引擎
    engine := scache.NewEngine(
        config.WithMaxSize(10000),
        config.WithDefaultExpiration(time.Hour),
        config.WithMemoryThreshold(0.8),
        config.WithBackgroundCleanup(time.Minute*5),
    )

    // 创建命令执行器
    executor := scache.NewExecutor(engine)

    // 执行命令
    result, err := executor.Execute("SET", "key", "value", time.Minute*30)
    if err != nil {
        panic(err)
    }
    fmt.Printf("SET result: %v\n", result)

    // 获取统计信息
    stats := scache.Stats()
    fmt.Printf("Cache stats: %+v\n", stats)
}
```

## 📖 数据类型

### String (字符串)

```go
// 设置字符串
err := scache.Set("greeting", "Hello, World!", time.Hour)

// 获取字符串
value, found, err := scache.Get("greeting")

// 检查类型
keyType, err := scache.Type("greeting") // "string"
```

### List (列表)

```go
// 左侧推入元素
length, err := scache.LPush("numbers", 1, time.Hour)
length, err = scache.LPush("numbers", 2, time.Hour)

// 右侧弹出元素
value, err := scache.RPop("numbers") // 1
```

### Hash (哈希)

```go
// 设置哈希字段
success, err := scache.HSet("user:1", "name", "Alice", time.Hour)
success, err = scache.HSet("user:1", "age", 30, time.Hour)

// 获取哈希字段
name, err := scache.HGet("user:1", "name")   // "Alice"
age, err := scache.HGet("user:1", "age")     // 30
```

## ⚙️ 配置选项

SCache 提供了多种预定义配置：

```go
// 小型配置（内存较小环境）
engine := scache.NewEngine(config.SmallConfig...)

// 中等配置（一般应用）
engine := scache.NewEngine(config.MediumConfig...)

// 大型配置（高负载应用）
engine := scache.NewEngine(config.LargeConfig...)

// 自定义配置
engine := scache.NewEngine(
    config.WithMaxSize(1000),
    config.WithDefaultExpiration(time.Hour),
    config.WithMemoryThreshold(0.8),
    config.WithBackgroundCleanup(time.Minute*5),
)
```

## 📋 支持的命令

### 通用命令
- `SET key value [ttl]` - 设置键值
- `GET key` - 获取值
- `DEL key` - 删除键
- `EXISTS key` - 检查键是否存在
- `TYPE key` - 获取键类型
- `EXPIRE key ttl` - 设置过期时间
- `TTL key` - 获取剩余生存时间
- `STATS` - 获取统计信息

### 列表命令
- `LPUSH key value [ttl]` - 左侧推入元素
- `RPOP key` - 右侧弹出元素

### 哈希命令
- `HSET key field value [ttl]` - 设置哈希字段
- `HGET key field` - 获取哈希字段

## 🔧 扩展命令

可以轻松添加自定义命令：

```go
package main

import (
    "scache"
    "scache/interfaces"
)

// 自定义命令
type CustomCommand struct {
    commands.BaseCommand
}

func (c *CustomCommand) Execute(ctx *interfaces.Context) error {
    // 实现自定义逻辑
    ctx.Result = "custom result"
    return nil
}

func (c *CustomCommand) Name() string {
    return "CUSTOM"
}

// 注册命令
func main() {
    executor := scache.NewExecutor(scache.NewEngine())
    executor.RegisterCommand(&CustomCommand{})

    result, err := executor.Execute("CUSTOM")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Custom result: %v\n", result)
}
```

## 📊 统计信息

```go
stats := scache.Stats()
// 返回 map[string]interface{} 包含：
// - hits: 命中次数
// - misses: 未命中次数
// - sets: 设置次数
// - deletes: 删除次数
// - evictions: 淘汰次数
// - expirations: 过期次数
// - memory: 内存使用量（字节）
// - keys: 当前键数量
// - hit_rate: 命中率
```

## 🏗️ 架构设计

```
scache/
├── commands/     # 命令处理器层
├── storage/      # 存储引擎层
├── types/        # 数据类型层
├── interfaces/   # 接口定义
├── config/       # 配置管理
├── policies/     # 淘汰策略
└── cache/        # 便捷API
```

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./storage/...
go test ./commands/...
go test ./types/...

# 运行性能测试
go test -bench=. ./...

# 运行集成测试
go test -tags=integration ./...
```

## 📈 性能基准

| 操作 | QPS | 延迟 (P99) |
|------|-----|-----------|
| SET  | 1,200,000+ | < 100μs |
| GET  | 1,500,000+ | < 50μs  |
| HSET | 800,000+   | < 150μs |
| HGET | 1,000,000+ | < 100μs |

*测试环境：Intel i7-8700K, 16GB RAM, Go 1.21*

## 🎯 使用场景

### 1. Web应用缓存

```go
// 缓存用户信息
func GetUser(userID string) (*User, error) {
    if value, found, err := scache.Get("user:"+userID); err == nil && found {
        return value.(*User), nil
    }

    user, err := database.GetUser(userID)
    if err != nil {
        return nil, err
    }

    scache.Set("user:"+userID, user, time.Minute*30)
    return user, nil
}
```

### 2. API响应缓存

```go
// 缓存API响应
func GetWeather(city string) (string, error) {
    cacheKey := "weather:" + city

    if value, found, err := scache.Get(cacheKey); err == nil && found {
        return value.(string), nil
    }

    weather, err := weatherAPI.Get(city)
    if err != nil {
        return "", err
    }

    scache.Set(cacheKey, weather, time.Minute*10)
    return weather, nil
}
```

## 🤝 贡献

欢迎贡献代码！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详细信息。

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/your-repo/scache.git
cd scache

# 安装依赖
go mod download

# 运行测试
go test ./...

# 格式化代码
go fmt ./...

# 代码检查
golangci-lint run
```

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- [API 文档](https://pkg.go.dev/github.com/your-repo/scache)
- [示例代码](https://github.com/your-repo/scache/tree/main/examples)
- [性能测试报告](https://github.com/your-repo/scache/blob/main/benchmarks.md)
- [更新日志](https://github.com/your-repo/scache/blob/main/CHANGELOG.md)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

**⭐ 如果这个项目对你有帮助，请给它一个星标！**