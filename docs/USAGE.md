# SCache 使用指南

SCache 提供了简单易用的缓存API，支持局部缓存和全局缓存两种使用方式。

## 快速开始

### 导入包

```go
import "scache/cache"
```

### 基本使用

```go
// 设置字符串
cache.SetString("key", "value", time.Minute*10)

// 获取字符串
if value, found := cache.GetString("key"); found {
    fmt.Println(value)
}
```

## 局部缓存

局部缓存适合需要多个独立缓存实例的场景。

### 创建局部缓存

```go
import "scache/config"

// 使用默认配置
localCache := cache.NewLocalCache()

// 使用自定义配置
localCache := cache.NewLocalCache(
    config.WithMaxSize(1000),                    // 最大1000个键
    config.WithDefaultExpiration(time.Hour),      // 默认1小时过期
    config.WithMemoryThreshold(0.8),              // 80%内存阈值
)
```

### 局部缓存操作

```go
// 字符串操作
localCache.SetString("user:1", "张三", time.Minute*30)
if name, found := localCache.GetString("user:1"); found {
    fmt.Printf("用户名: %s\n", name)
}

// 列表操作
permissions := []interface{}{"read", "write", "delete"}
localCache.SetList("permissions:1", permissions, 0)
if perms, found := localCache.GetList("permissions:1"); found {
    fmt.Printf("权限: %v\n", perms)
}

// 哈希操作
profile := map[string]interface{}{
    "name":  "张三",
    "age":   30,
    "email": "zhangsan@example.com",
}
localCache.SetHash("profile:1", profile, 0)
if prof, found := localCache.GetHash("profile:1"); found {
    fmt.Printf("档案: %+v\n", prof)
}
```

## 全局缓存

全局缓存适合单例应用场景，使用简单方便。

### 初始化全局缓存

```go
// 使用默认配置（自动初始化）
cache.SetString("key", "value", 0)

// 自定义初始化
cache.InitGlobalCache(
    config.WithMaxSize(5000),
    config.WithDefaultExpiration(time.Hour*2),
)
```

### 全局缓存操作

```go
// 字符串操作
cache.SetString("app_name", "MyApp", 0)
if name, found := cache.GetString("app_name"); found {
    fmt.Printf("应用名: %s\n", name)
}

// 列表操作
cache.SetList("admin_users", []interface{}{"admin1", "admin2"}, 0)

// 哈希操作
config := map[string]interface{}{
    "debug": true,
    "port":  8080,
}
cache.SetHash("app_config", config, 0)
```

## 数据类型支持

### 字符串 (String)

```go
// 设置
cache.SetString("message", "Hello, World!", time.Minute*10)

// 获取
if msg, found := cache.GetString("message"); found {
    fmt.Println(msg)
}
```

### 列表 (List)

```go
// 设置
items := []interface{}{"item1", "item2", "item3"}
cache.SetList("mylist", items, 0)

// 获取
if list, found := cache.GetList("mylist"); found {
    fmt.Printf("列表: %v\n", list)
}
```

### 哈希 (Hash)

```go
// 设置
data := map[string]interface{}{
    "field1": "value1",
    "field2": "value2",
    "field3": 123,
}
cache.SetHash("myhash", data, 0)

// 获取
if hash, found := cache.GetHash("myhash"); found {
    fmt.Printf("哈希: %+v\n", hash)
}
```

## 过期管理

### TTL操作

```go
// 设置带TTL的键
cache.SetString("temp", "temporary data", time.Minute*5)

// 动态设置过期时间
cache.Expire("temp", time.Minute*10)

// 查看剩余时间
if ttl, exists := cache.TTL("temp"); exists {
    fmt.Printf("剩余时间: %v\n", ttl)
}
```

### 过期检查

```go
// 检查键是否存在且未过期
if cache.Exists("key") {
    fmt.Println("键存在且未过期")
}
```

## 缓存管理

### 基本操作

```go
// 删除键
cache.Delete("key")

// 检查键是否存在
if cache.Exists("key") {
    // 处理存在的键
}

// 获取所有键
keys := cache.Keys()
fmt.Printf("所有键: %v\n", keys)

// 获取缓存大小
size := cache.Size()
fmt.Printf("缓存大小: %d\n", size)

// 清空所有缓存
cache.Flush()
```

### 统计信息

```go
stats := cache.Stats()
if statsMap, ok := stats.(map[string]interface{}); ok {
    fmt.Printf("命中次数: %.0f\n", statsMap["hits"])
    fmt.Printf("未命中次数: %.0f\n", statsMap["misses"])
    fmt.Printf("命中率: %.2f%%\n", statsMap["hit_rate"].(float64)*100)
    fmt.Printf("当前大小: %.0f\n", statsMap["keys"])
}
```

## 配置选项

### 预定义配置

```go
// 小型配置（适合内存较小的环境）
cache.InitGlobalCache(config.SmallConfig...)

// 中型配置（适合一般应用）
cache.InitGlobalCache(config.MediumConfig...)

// 大型配置（适合高负载应用）
cache.InitGlobalCache(config.LargeConfig...)

// 默认配置
cache.InitGlobalCache(config.DefaultConfig...)
```

### 自定义配置

```go
cache.InitGlobalCache(
    config.WithMaxSize(10000),                    // 最大10000个键
    config.WithMemoryThreshold(0.85),             // 85%内存阈值
    config.WithDefaultExpiration(time.Hour*6),    // 默认6小时过期
    config.WithBackgroundCleanup(time.Minute*5),  // 5分钟清理间隔
)
```

## 高级用法

### 多缓存实例

```go
// 用户缓存
userCache := cache.NewLocalCache(
    config.WithMaxSize(1000),
    config.WithDefaultExpiration(time.Hour),
)

// 配置缓存
configCache := cache.NewLocalCache(
    config.WithMaxSize(100),
    config.WithDefaultExpiration(time.Hour*24),
)

// 会话缓存
sessionCache := cache.NewLocalCache(
    config.WithMaxSize(5000),
    config.WithDefaultExpiration(time.Minute*30),
)
```

### 获取底层引擎

```go
// 对于需要更复杂操作的场景
localCache := cache.NewLocalCache()
engine := localCache.GetEngine()

// 使用底层引擎的完整功能
obj, exists := engine.Get("key")
if exists {
    // 处理对象
}
```

## 最佳实践

1. **选择合适的缓存类型**
   - 简单应用使用全局缓存
   - 需要多实例或隔离时使用局部缓存

2. **合理设置过期时间**
   - 用户会话：30分钟-2小时
   - 配置信息：几小时到几天
   - 临时数据：几分钟

3. **监控缓存统计**
   - 定期检查命中率
   - 监控缓存大小
   - 根据统计调整配置

4. **避免缓存过大**
   - 设置合理的MaxSize
   - 定期清理无用数据
   - 使用适当的内存阈值

## 错误处理

```go
// 大部分操作都是安全的，只有设置操作可能返回错误
err := cache.SetString("key", "value", time.Minute)
if err != nil {
    log.Printf("设置缓存失败: %v", err)
}

// 获取操作不会返回错误，只返回是否存在
if value, found := cache.GetString("key"); found {
    // 使用value
} else {
    // 处理不存在的情况
}
```

## 示例项目

查看 `examples/` 目录下的完整示例：

- `examples/basic/` - 基本使用示例
- `examples/simple/` - 简单示例
- `examples/concurrent/` - 并发使用示例
- `examples/webserver/` - Web服务器集成示例