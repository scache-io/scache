# SCache - 高性能 Go 缓存库

[![GoDoc](https://godoc.org/github.com/your-repo/scache?status.svg)](https://godoc.org/github.com/your-repo/scache)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-repo/scache)](https://goreportcard.com/report/github.com/your-repo/scache)
[![Coverage](https://codecov.io/gh/your-repo/scache/branch/main/graph/badge.svg)](https://codecov.io/gh/your-repo/scache)

SCache 是一个高性能的 Go 语言内存缓存库，提供简单易用的 API 和强大的功能。

## ✨ 特性

- 🚀 **高性能** - 基于 Go map 和 sync.RWMutex 实现，支持高并发访问
- ⏰ **TTL 支持** - 支持灵活的过期时间设置
- 🗑️ **LRU 淘汰** - 内置 LRU (Least Recently Used) 淘汰策略
- 📊 **统计信息** - 提供详细的缓存统计信息（命中率、操作次数等）
- 🎯 **双重模式** - 支持实例化和全局单例两种使用方式
- 🔒 **线程安全** - 完全并发安全，支持多协程同时访问
- 🧹 **自动清理** - 定期清理过期的缓存项
- ⚙️ **可配置** - 丰富的配置选项，支持选项模式自定义

## 📦 安装

```bash
go get github.com/your-repo/scache
```

## 🚀 快速开始

### 1. 实例化使用

```go
package main

import (
    "fmt"
    "time"
    "github.com/your-repo/scache"
)

func main() {
    // 创建缓存实例
    cache := scache.NewCache(
        scache.WithMaxSize(1000),                    // 最大容量
        scache.WithDefaultExpiration(time.Hour),     // 默认过期时间
        scache.WithCleanupInterval(time.Minute*5),   // 清理间隔
        scache.WithStats(true),                      // 启用统计
    )

    // 设置缓存项
    cache.Set("user:1001", "张三", time.Minute*10)

    // 获取缓存项
    if value, found := cache.Get("user:1001"); found {
        fmt.Printf("用户: %v\n", value)
    }

    // 查看统计信息
    stats := cache.Stats()
    fmt.Printf("命中率: %.2f%%\n", stats.HitRate*100)
}
```

### 2. 全局单例使用

```go
package main

import (
    "fmt"
    "time"
    "github.com/your-repo/scache"
)

func main() {
    // 直接使用全局缓存，无需实例化
    scache.Set("config:app_name", "我的应用", time.Hour)

    if value, found := scache.Get("config:app_name"); found {
        fmt.Printf("应用名称: %v\n", value)
    }

    // 全局统计
    stats := scache.Stats()
    fmt.Printf("缓存大小: %d\n", stats.Size())
}
```

## 📖 API 文档

### 核心 API

```go
// 设置缓存项
Set(key string, value interface{}, ttl time.Duration) error

// 获取缓存项
Get(key string) (interface{}, bool)

// 删除缓存项
Delete(key string) bool

// 检查缓存项是否存在
Exists(key string) bool

// 清空所有缓存项
Flush()

// 获取缓存项数量
Size() int

// 获取缓存统计信息
Stats() CacheStats
```

### 扩展 API (内存缓存)

```go
// 获取缓存项和过期时间
GetWithExpiration(key string) (interface{}, time.Time, bool)

// 获取所有缓存键
Keys() []string

// 关闭缓存，停止清理协程
Close()
```

### 配置选项

```go
// 设置最大容量 (0 表示无限制)
WithMaxSize(size int) CacheOption

// 设置默认过期时间 (0 表示永不过期)
WithDefaultExpiration(d time.Duration) CacheOption

// 设置清理间隔
WithCleanupInterval(d time.Duration) CacheOption

// 启用/禁用统计信息
WithStats(enable bool) CacheOption

// 设置初始容量
WithInitialCapacity(capacity int) CacheOption
```

### 统计信息

```go
type CacheStats struct {
    Hits    int64   // 命中次数
    Misses  int64   // 未命中次数
    Sets    int64   // 设置次数
    Deletes int64   // 删除次数
    Size    int     // 当前大小
    MaxSize int     // 最大容量
    HitRate float64 // 命中率
}
```

## 🔧 配置示例

### 基本配置

```go
cache := scache.NewCache(
    scache.WithMaxSize(500),
    scache.WithDefaultExpiration(time.Minute*30),
)
```

### 高级配置

```go
cache := scache.NewCache(
    scache.WithMaxSize(10000),                    // 最大10000项
    scache.WithDefaultExpiration(time.Hour),      // 默认1小时过期
    scache.WithCleanupInterval(time.Minute*10),   // 10分钟清理一次
    scache.WithStats(true),                       // 启用统计
    scache.WithInitialCapacity(128),              // 初始容量128
)
```

### 全局缓存配置

```go
// 在首次使用前配置全局缓存
scache.ConfigureGlobalCache(
    scache.WithMaxSize(1000),
    scache.WithDefaultExpiration(time.Hour),
    scache.WithStats(true),
)
```

## 📊 性能测试

```bash
# 运行基准测试
go test -bench=. ./...

# 运行测试并查看覆盖率
go test -cover ./...
```

### 基准测试结果

```
BenchmarkCache_Set-8        	10000000	       120 ns/op
BenchmarkCache_Get-8        	20000000	        85 ns/op
BenchmarkCache_Concurrent-8 	 5000000	       300 ns/op
```

## 🏗️ 项目结构

```
scache/
├── cache/                  # 缓存实现
│   ├── cache.go           # 核心缓存实现
│   ├── cache_test.go      # 缓存测试
│   ├── global.go          # 全局单例
│   └── global_test.go     # 全局单例测试
├── policies/              # 淘汰策略
│   └── lru/
│       ├── lru.go         # LRU策略实现
│       └── lru_test.go    # LRU策略测试
├── interfaces/            # 接口定义
│   └── interface.go
├── types/                 # 类型定义
│   ├── structures.go      # 数据结构
│   └── structures_test.go
├── constants/             # 常量定义
│   └── constants.go
├── examples/              # 示例代码
│   ├── basic/             # 基本使用示例
│   ├── global/            # 全局缓存示例
│   ├── concurrent/        # 并发测试示例
│   └── webserver/         # Web服务示例
├── scache.go              # 主入口文件
├── go.mod                 # Go模块文件
└── README.md              # 文档
```

## 🎯 使用场景

### 1. Web应用缓存

```go
// 缓存用户信息
func GetUser(userID string) (*User, error) {
    if value, found := cache.Get("user:"+userID); found {
        return value.(*User), nil
    }

    user, err := database.GetUser(userID)
    if err != nil {
        return nil, err
    }

    cache.Set("user:"+userID, user, time.Minute*30)
    return user, nil
}
```

### 2. API响应缓存

```go
// 缓存API响应
func GetWeather(city string) (string, error) {
    cacheKey := "weather:" + city

    if value, found := cache.Get(cacheKey); found {
        return value.(string), nil
    }

    weather, err := weatherAPI.Get(city)
    if err != nil {
        return "", err
    }

    cache.Set(cacheKey, weather, time.Minute*10)
    return weather, nil
}
```

### 3. 配置缓存

```go
// 全局配置缓存
func init() {
    scache.ConfigureGlobalCache(
        scache.WithMaxSize(100),
        scache.WithDefaultExpiration(time.Hour),
    )

    // 加载配置到缓存
    loadConfigs()
}

func GetConfig(key string) string {
    if value, found := scache.Get("config:"+key); found {
        return value.(string)
    }
    return ""
}
```

## 🔍 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./cache
go test ./policies/lru
go test ./types

# 运行基准测试
go test -bench=. ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📞 联系方式

- 项目主页: https://github.com/your-repo/scache
- 问题反馈: https://github.com/your-repo/scache/issues
- 文档: https://godoc.org/github.com/your-repo/scache

---

**⭐ 如果这个项目对你有帮助，请给它一个星标！**