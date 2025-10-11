# SCache - 高性能 Go 缓存框架

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/scache)](https://goreportcard.com/report/github.com/yourusername/scache)
[![GoDoc](https://godoc.org/github.com/yourusername/scache?status.svg)](https://godoc.org/github.com/yourusername/scache)

SCache 是一个用 Go 语言编写的高性能、通用的内存缓存框架，专为现代 Go 应用程序设计。它提供了丰富的功能，包括多种淘汰策略、TTL 支持、并发安全以及详细的统计信息。

## 特性

- 🚀 **高性能** - 分片设计减少锁竞争，提供出色的并发性能
- 🔄 **多种淘汰策略** - 支持 LRU、LFU、FIFO 等缓存淘汰策略
- ⏰ **TTL 支持** - 支持带过期时间的缓存项，自动清理过期数据
- 🔒 **并发安全** - 使用读写锁和分片技术确保并发安全
- 📊 **统计信息** - 提供详细的缓存命中率和性能统计
- 🔧 **高度可配置** - 灵活的配置选项，满足不同场景需求
- 📦 **易于集成** - 简洁的 API 设计，可轻松集成到现有项目中
- 🌐 **全局缓存管理** - 支持全局注册机制，便于在大型应用中管理多个缓存
- 🏗️ **模块化设计** - 清晰的项目结构，便于维护和扩展

## 快速开始

### 安装

```bash
go get github.com/yourusername/scache
```

### 项目结构

```
scache/
├── scache.go                 # 主入口文件，重新导出所有功能
├── pkg/                      # 核心包
│   ├── cache/               # 缓存核心实现
│   ├── policies/            # 淘汰策略实现
│   │   ├── lru/            # LRU 策略
│   │   ├── lfu/            # LFU 策略
│   │   └── fifo/           # FIFO 策略
│   ├── manager/            # 全局缓存管理器
│   └── global/             # 全局便捷函数
├── cmd/                     # 示例和命令
│   └── examples/
│       ├── basic/           # 基础示例
│       └── advanced/        # 高级示例
└── examples/                # 兼容旧版本的示例
```

### 两种使用方式

SCache 提供两种使用方式：传统的实例化方式和全局缓存方式。

#### 方式一：传统实例化

```go
package main

import (
	"fmt"
	"time"

	"scache"
)

func main() {
	// 创建缓存实例
	c := scache.New()
	defer c.Close()

	// 或者创建特定策略的缓存
	lruCache := scache.NewLRU(1000)
	lfuCache := scache.NewLFU(1000)
	fifoCache := scache.NewFIFO(1000)

	// 设置和获取缓存
	c.Set("key1", "value1")
	if value, exists := c.Get("key1"); exists {
		fmt.Println("找到值:", value)
	}

	// 设置带过期时间的缓存
	c.SetWithTTL("key2", "value2", 5*time.Minute)
}
```

#### 方式二：全局缓存管理

```go
package main

import (
	"fmt"

	"scache/pkg/global"
)

func main() {
	// 注册不同类型的全局缓存
	global.RegisterLRU("users", 1000)      // 用户缓存
	global.RegisterLFU("sessions", 500)   // 会话缓存
	global.RegisterFIFO("products", 2000) // 产品缓存

	// 获取并使用缓存
	usersCache, _ := global.Get("users")
	usersCache.Set("user:1", "Alice")

	// 或者使用默认缓存
	global.Set("app:version", "1.0.0")
	if value, exists := global.GetFromDefault("app:version"); exists {
		fmt.Println("应用版本:", value)
	}

	// 清理
	global.Close()
}
```

## 高级用法

### 使用不同的淘汰策略

```go
// LRU (Least Recently Used)
lruCache := cache.NewLRU(1000) // 最大 1000 项

// LFU (Least Frequently Used)
lfuCache := cache.NewLFU(1000)

// FIFO (First In First Out)
fifoCache := cache.NewFIFO(1000)
```

### 自定义配置

```go
c := cache.New(
	cache.WithMaxSize(10000),           // 最大缓存项数量
	cache.WithDefaultTTL(30*time.Minute), // 默认过期时间
	cache.WithEvictionPolicy("lru"),     // 淘汰策略
	cache.WithShards(16),               // 分片数量
	cache.WithStatistics(true),         // 启用统计
	cache.WithCleanupInterval(10*time.Minute), // 清理间隔
)
```

### 并发使用示例

```go
func handleRequest(cache cache.Cache, userID string) {
	// 尝试从缓存获取用户数据
	if userData, exists := cache.Get("user:" + userID); exists {
		// 缓存命中
		processUserData(userData)
		return
	}

	// 缓存未命中，从数据库加载
	userData := loadUserFromDB(userID)

	// 存入缓存，设置 5 分钟过期
	cache.SetWithTTL("user:"+userID, userData, 5*time.Minute)

	processUserData(userData)
}
```

## 配置选项

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `MaxSize` | `int` | `10000` | 最大缓存项数量 |
| `DefaultTTL` | `time.Duration` | `0` | 默认过期时间（0 表示永不过期） |
| `EvictionPolicy` | `string` | `"lru"` | 淘汰策略 (lru/lfu/fifo) |
| `Shards` | `int` | `16` | 分片数量，影响并发性能 |
| `CleanupInterval` | `time.Duration` | `10分钟` | 过期项清理间隔 |
| `EnableStatistics` | `bool` | `true` | 是否启用统计信息 |
| `EnableLazyExpiration` | `bool` | `true` | 是否启用懒过期检查 |

## API 参考

### 传统实例方式

#### 基本操作

- `Set(key string, value interface{}) error` - 设置缓存项
- `SetWithTTL(key string, value interface{}, ttl time.Duration) error` - 设置带过期时间的缓存项
- `Get(key string) (interface{}, bool)` - 获取缓存项
- `Delete(key string) bool` - 删除缓存项
- `Exists(key string) bool` - 检查缓存项是否存在
- `Clear() error` - 清空所有缓存

#### 批量操作

- `SetBatch(items map[string]interface{}) error` - 批量设置缓存项
- `GetBatch(keys []string) map[string]interface{}` - 批量获取缓存项
- `DeleteBatch(keys []string) map[string]bool` - 批量删除缓存项

#### 统计信息

- `Size() int` - 获取当前缓存项数量
- `Keys() []string` - 获取所有键
- `Stats() CacheStats` - 获取详细统计信息

#### 生命周期

- `Close() error` - 关闭缓存，释放资源

### 全局缓存方式

#### 缓存管理

- `Register(name string, c Cache) error` - 注册缓存
- `RegisterLRU(name string, maxSize int, opts ...Option) error` - 注册 LRU 缓存
- `RegisterLFU(name string, maxSize int, opts ...Option) error` - 注册 LFU 缓存
- `RegisterFIFO(name string, maxSize int, opts ...Option) error` - 注册 FIFO 缓存
- `Get(name string) (Cache, error)` - 获取已注册的缓存
- `GetOrDefault(name string, opts ...Option) Cache` - 获取缓存，不存在则创建默认缓存
- `Remove(name string) error` - 移除已注册的缓存
- `List() []string` - 列出所有已注册的缓存名称
- `Exists(name string) bool` - 检查缓存是否已注册

#### 默认缓存操作

- `Set(key string, value interface{}) error` - 在默认缓存中设置键值
- `SetWithTTL(key string, value interface{}, ttl time.Duration) error` - 在默认缓存中设置带过期时间的键值
- `GetFromDefault(key string) (interface{}, bool)` - 从默认缓存中获取值
- `Delete(key string) bool` - 从默认缓存中删除键
- `ExistsInKey(key string) bool` - 检查默认缓存中是否存在键
- `ClearDefault() error` - 清空默认缓存

#### 全局管理

- `Clear() error` - 清空所有缓存
- `Close() error` - 关闭所有缓存并清理管理器
- `Stats() map[string]CacheStats` - 获取所有缓存的统计信息
- `Size() int` - 获取所有缓存的总大小

## 性能

基准测试结果（Apple M1 Pro）：

```
BenchmarkMemoryCache_Set-10          3588164    327.3 ns/op
BenchmarkMemoryCache_Get-10          6415347    182.3 ns/op
BenchmarkMemoryCache_SetWithTTL-10   3879196    312.2 ns/op
BenchmarkMemoryCache_ConcurrentOps   1451967    888.5 ns/op
```

## 淘汰策略详解

### LRU (Least Recently Used)
- 最近最少使用策略
- 优先淘汰最长时间未被访问的缓存项
- 适用于访问模式有局部性的场景

### LFU (Least Frequently Used)
- 最少使用频率策略
- 优先淘汰访问次数最少的缓存项
- 适用于热点数据明显的场景

### FIFO (First In First Out)
- 先进先出策略
- 按照添加时间顺序淘汰缓存项
- 适用于缓存项访问时间均匀的场景

## 最佳实践

1. **选择合适的分片数量**：对于高并发场景，建议使用 16-64 个分片
2. **设置合理的 TTL**：避免缓存项无限期存在，设置适当的过期时间
3. **监控命中率**：定期检查缓存命中率，调整缓存策略
4. **合理设置容量**：根据内存大小和应用需求设置最大缓存数量
5. **使用批量操作**：对于多个缓存操作，优先使用批量 API

## 示例项目

查看 `examples/` 目录中的完整示例：

- [基础使用示例](examples/basic/main.go)
- [Web 服务集成示例](examples/webserver/main.go)
- [高并发场景示例](examples/concurrent/main.go)

## 贡献

欢迎提交 Issue 和 Pull Request！请确保：

1. 代码通过所有测试
2. 遵循 Go 代码规范
3. 添加必要的测试用例
4. 更新相关文档

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持 LRU、LFU、FIFO 淘汰策略
- 支持 TTL 和自动过期清理
- 分片设计，高并发性能优化
- 完整的统计信息和监控支持