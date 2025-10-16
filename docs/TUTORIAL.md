# SCache 教程

## 目录

1. [快速入门](#快速入门)
2. [基础概念](#基础概念)
3. [数据类型详解](#数据类型详解)
4. [命令系统](#命令系统)
5. [配置管理](#配置管理)
6. [高级用法](#高级用法)
7. [性能优化](#性能优化)
8. [常见问题](#常见问题)
9. [实战案例](#实战案例)

## 快速入门

### 安装

```bash
go get github.com/your-repo/scache
```

### 第一个缓存程序

```go
package main

import (
    "fmt"
    "time"
    "scache"
)

func main() {
    // 设置一个缓存项
    err := scache.Set("message", "Hello, SCache!", time.Minute)
    if err != nil {
        panic(err)
    }

    // 获取缓存项
    value, found, err := scache.Get("message")
    if err != nil {
        panic(err)
    }

    if found {
        fmt.Println("缓存值:", value)
    }

    // 检查键是否存在
    exists, err := scache.Exists("message")
    if err != nil {
        panic(err)
    }
    fmt.Println("键存在:", exists)
}
```

编译并运行：

```bash
go run main.go
```

输出：
```
缓存值: Hello, SCache!
键存在: true
```

## 基础概念

### 存储引擎

SCache 的核心是存储引擎（StorageEngine），它负责所有数据的存储和检索：

```go
import "scache/storage"

// 创建存储引擎
engine := storage.NewStorageEngine(nil)

// 直接使用存储引擎
strObj := types.NewStringObject("value", time.Hour)
engine.Set("key", strObj)

obj, exists := engine.Get("key")
if exists {
    if strObj, ok := obj.(*types.StringObject); ok {
        fmt.Println(strObj.Value())
    }
}
```

### 命令执行器

命令执行器（Executor）提供了更高级的接口：

```go
import "scache"

// 创建执行器
engine := scache.NewEngine()
executor := scache.NewExecutor(engine)

// 执行命令
result, err := executor.Execute("SET", "key", "value", time.Hour)
if err != nil {
    panic(err)
}

result, err = executor.Execute("GET", "key")
if err != nil {
    panic(err)
}
fmt.Println("结果:", result)
```

### 便捷 API

对于简单用例，SCache 提供了全局便捷函数：

```go
// 这些函数使用全局默认实例
err := scache.Set("key", "value", time.Hour)
value, found, err := scache.Get("key")
```

## 数据类型详解

### String（字符串）

字符串是最简单的数据类型，用于存储文本或二进制数据。

```go
// 设置字符串
err := scache.Set("greeting", "Hello, World!", time.Hour)
if err != nil {
    panic(err)
}

// 获取字符串
value, found, err := scache.Get("greeting")
if err != nil {
    panic(err)
}
if found {
    fmt.Printf("问候语: %s\n", value.(string))
}

// 检查类型
keyType, err := scache.Type("greeting")
if err != nil {
    panic(err)
}
fmt.Printf("数据类型: %s\n", keyType) // 输出: string
```

### List（列表）

列表用于存储有序的元素集合，类似于数组。

```go
// 添加元素到列表左侧
length, err := scache.LPush("numbers", 1, time.Hour)
if err != nil {
    panic(err)
}
fmt.Printf("列表长度: %d\n", length) // 输出: 1

// 继续添加元素
length, err = scache.LPush("numbers", 2, time.Hour)
fmt.Printf("列表长度: %d\n", length) // 输出: 2

// 从列表右侧弹出元素
value, err := scache.RPop("numbers")
if err != nil {
    panic(err)
}
fmt.Printf("弹出的元素: %v\n", value) // 输出: 1

// 再次弹出
value, err = scache.RPop("numbers")
fmt.Printf("弹出的元素: %v\n", value) // 输出: 2

// 尝试从空列表弹出
value, err = scache.RPop("numbers")
if err != nil {
    panic(err)
}
fmt.Printf("空列表弹出: %v\n", value) // 输出: <nil>
```

### Hash（哈希）

哈希用于存储键值对集合，类似于字典或映射。

```go
// 设置哈希字段
success, err := scache.HSet("user:1", "name", "Alice", time.Hour)
if err != nil {
    panic(err)
}
fmt.Printf("设置字段结果: %v\n", success) // 输出: true

// 设置更多字段
success, err = scache.HSet("user:1", "age", 30, time.Hour)
success, err = scache.HSet("user:1", "email", "alice@example.com", time.Hour)

// 获取哈希字段
name, err := scache.HGet("user:1", "name")
if err != nil {
    panic(err)
}
fmt.Printf("用户名: %v\n", name) // 输出: Alice

age, err := scache.HGet("user:1", "age")
fmt.Printf("年龄: %v\n", age) // 输出: 30

// 获取不存在的字段
city, err := scache.HGet("user:1", "city")
if err != nil {
    panic(err)
}
fmt.Printf("城市: %v\n", city) // 输出: <nil>
```

## 命令系统

SCache 采用命令模式，易于扩展和测试。

### 基本命令使用

```go
engine := scache.NewEngine()
executor := scache.NewExecutor(engine)

// SET 命令
result, err := executor.Execute("SET", "key", "value", time.Minute)
if err != nil {
    panic(err)
}

// GET 命令
result, err = executor.Execute("GET", "key")
if err != nil {
    panic(err)
}
fmt.Println("GET 结果:", result)

// DEL 命令
result, err = executor.Execute("DEL", "key")
if err != nil {
    panic(err)
}
fmt.Println("删除结果:", result.(bool))

// EXISTS 命令
result, err = executor.Execute("EXISTS", "key")
fmt.Println("存在性:", result.(bool))
```

### 批量操作

```go
// 批量设置多个键
keys := []string{"key1", "key2", "key3"}
values := []string{"value1", "value2", "value3"}

for i, key := range keys {
    _, err := executor.Execute("SET", key, values[i], time.Minute)
    if err != nil {
        panic(err)
    }
}

// 批量获取
for _, key := range keys {
    result, err := executor.Execute("GET", key)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s: %v\n", key, result)
}
```

## 配置管理

### 使用预定义配置

```go
import "scache/config"

// 小型配置 - 适合内存较小的环境
smallEngine := scache.NewEngine(config.SmallConfig...)

// 中等配置 - 适合一般应用
mediumEngine := scache.NewEngine(config.MediumConfig...)

// 大型配置 - 适合高负载应用
largeEngine := scache.NewEngine(config.LargeConfig...)
```

### 自定义配置

```go
// 创建自定义配置
engine := scache.NewEngine(
    config.WithMaxSize(5000),              // 最大5000个键
    config.WithDefaultExpiration(time.Hour), // 默认1小时过期
    config.WithMemoryThreshold(0.75),       // 75%内存阈值
    config.WithBackgroundCleanup(time.Minute*3), // 3分钟清理间隔
)

executor := scache.NewExecutor(engine)
```

### 配置参数说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| MaxSize | int | 0 | 最大缓存数量，0表示无限制 |
| MemoryThreshold | float64 | 0.8 | 内存检查阈值（0.0-1.0） |
| DefaultExpiration | time.Duration | 0 | 默认过期时间，0表示永不过期 |
| BackgroundCleanupInterval | time.Duration | 5分钟 | 后台清理间隔，0表示禁用 |

## 高级用法

### 过期时间管理

```go
// 设置不同过期时间的缓存
scache.Set("short_term", "很快过期", time.Second*30)    // 30秒
scache.Set("medium_term", "中等过期", time.Hour*2)     // 2小时
scache.Set("long_term", "长期有效", 0)                // 永不过期

// 动态设置过期时间
scache.Set("dynamic", "动态过期", time.Minute)       // 初始1分钟
success, err := scache.Expire("dynamic", time.Hour)  // 延长到1小时
fmt.Printf("设置过期时间结果: %v\n", success)

// 检查剩余生存时间
ttl, err := scache.TTL("dynamic")
switch ttl {
case -2:
    fmt.Println("键不存在")
case -1:
    fmt.Println("永不过期")
default:
    fmt.Printf("剩余时间: %d秒\n", ttl)
}
```

### 类型检查和转换

```go
// 设置不同类型的数据
scache.Set("text", "纯文本", time.Hour)
scache.LPush("mylist", "列表项", time.Hour)
scache.HSet("myhash", "字段", "值", time.Hour)

// 检查类型
types := []string{"text", "mylist", "myhash", "nonexistent"}
for _, key := range types {
    keyType, err := scache.Type(key)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s 的类型: %s\n", key, keyType)
}
```

### 统计信息监控

```go
// 执行一些操作
for i := 0; i < 100; i++ {
    scache.Set(fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i), time.Minute)
    scache.Get(fmt.Sprintf("key_%d", i))
    scache.Get(fmt.Sprintf("missing_%d", i)) // 这些会miss
}

// 获取统计信息
stats := scache.Stats()
statsMap := stats.(map[string]interface{})

fmt.Println("=== 缓存统计 ===")
fmt.Printf("命中次数: %d\n", statsMap["hits"])
fmt.Printf("未命中次数: %d\n", statsMap["misses"])
fmt.Printf("设置次数: %d\n", statsMap["sets"])
fmt.Printf("当前键数量: %d\n", statsMap["keys"])
fmt.Printf("命中率: %.2f%%\n", statsMap["hit_rate"].(float64)*100)
```

## 性能优化

### 批量操作优化

```go
// ❌ 低效的方式 - 多次调用
for i := 0; i < 1000; i++ {
    scache.Set(fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i), time.Minute)
}

// ✅ 高效的方式 - 使用执行器批量处理
engine := scache.NewEngine()
executor := scache.NewExecutor(engine)

for i := 0; i < 1000; i++ {
    executor.Execute("SET", fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i), time.Minute)
}
```

### 合理设置过期时间

```go
// ❌ 过短的过期时间 - 导致频繁的缓存失效
scache.Set("config", "配置数据", time.Second*10)

// ✅ 合理的过期时间 - 根据数据更新频率设置
scache.Set("config", "配置数据", time.Hour)    // 配置数据变化较少
scache.Set("session", "会话数据", time.Minute*30) // 会话数据
scache.Set("cache", "临时数据", time.Minute*5)    // 临时数据
```

### 内存管理

```go
// 设置合适的内存阈值
engine := scache.NewEngine(
    config.WithMaxSize(10000),           // 限制键数量
    config.WithMemoryThreshold(0.8),     // 80%内存使用率时开始清理
    config.WithBackgroundCleanup(time.Minute*2), // 定期清理
)

// 监控内存使用
stats := engine.Stats().(map[string]interface{})
fmt.Printf("当前键数量: %d\n", stats["keys"])
fmt.Printf("内存使用: %d bytes\n", stats["memory"])
```

## 常见问题

### Q1: 如何选择合适的过期时间？

**A:** 根据数据的特性和更新频率来决定：

- 配置数据：几小时到几天
- 用户会话：30分钟到几小时
- API响应：几分钟到几小时
- 临时计算结果：几秒到几分钟

```go
// 配置数据 - 变化频率低
scache.Set("app_config", configData, time.Hour*6)

// 用户数据 - 中等变化频率
scache.Set("user_profile", userData, time.Minute*30)

// API缓存 - 变化频率较高
scache.Set("api_response", responseData, time.Minute*5)
```

### Q2: 如何处理缓存雪崩？

**A:** 使用随机过期时间：

```go
import "math/rand"

func setWithRandomTTL(key string, value interface{}, baseTTL time.Duration) {
    // 添加随机偏移，避免同时过期
    randomOffset := time.Duration(rand.Intn(300)) * time.Second // 0-5分钟随机偏移
    finalTTL := baseTTL + randomOffset

    scache.Set(key, value, finalTTL)
}

// 使用示例
setWithRandomTTL("hot_data", "value", time.Hour)
```

### Q3: 如何监控缓存性能？

**A:** 定期检查统计信息：

```go
func monitorCache() {
    ticker := time.NewTicker(time.Minute * 5)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            stats := scache.Stats().(map[string]interface{})
            hitRate := stats["hit_rate"].(float64)

            if hitRate < 0.8 {
                log.Printf("警告: 命中率过低 %.2f%%", hitRate*100)
            }

            keys := stats["keys"].(int)
            log.Printf("缓存状态: %d 个键, 命中率 %.2f%%", keys, hitRate*100)
        }
    }
}
```

### Q4: 内存使用过高怎么办？

**A:** 调整配置参数：

```go
// 减少最大容量
engine := scache.NewEngine(
    config.WithMaxSize(5000),              // 降低最大键数量
    config.WithMemoryThreshold(0.7),       // 降低内存阈值
    config.WithBackgroundCleanup(time.Minute), // 增加清理频率
)

// 或者手动清理
if keys := scache.Keys(); len(keys) > 8000 {
    // 清理一些键
    for i := 0; i < 1000; i++ {
        scache.Delete(keys[i])
    }
}
```

## 实战案例

### 案例1: Web应用用户缓存

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "scache"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 模拟数据库
var database = map[int]User{
    1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
    2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func getUserFromDB(id int) (User, error) {
    // 模拟数据库查询延迟
    time.Sleep(time.Millisecond * 100)
    user, exists := database[id]
    if !exists {
        return User{}, fmt.Errorf("user not found")
    }
    return user, nil
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    userID := 1 // 简化处理

    // 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:%d", userID)
    cachedUser, found, err := scache.Get(cacheKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var user User
    if found {
        // 缓存命中
        user = cachedUser.(User)
        fmt.Printf("缓存命中: 用户 %s\n", user.Name)
    } else {
        // 缓存未命中，查询数据库
        user, err = getUserFromDB(userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        // 存入缓存，30分钟过期
        err = scache.Set(cacheKey, user, time.Minute*30)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Printf("数据库查询: 用户 %s\n", user.Name)
    }

    // 返回JSON响应
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func main() {
    // 配置缓存
    engine := scache.NewEngine(
        config.WithMaxSize(10000),
        config.WithDefaultExpiration(time.Minute*30),
    )
    executor := scache.NewExecutor(engine)

    // 启动缓存监控
    go monitorCache()

    http.HandleFunc("/user", getUserHandler)
    fmt.Println("服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}
```

### 案例2: API响应缓存

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "scache"
    "scache/config"
)

// WeatherAPI 天气API客户端
type WeatherAPI struct {
    cacheEngine interfaces.StorageEngine
    executor    *scache.Executor
}

func NewWeatherAPI() *WeatherAPI {
    engine := scache.NewEngine(
        config.WithMaxSize(1000),
        config.WithDefaultExpiration(time.Minute*10), // 天气数据10分钟过期
    )
    return &WeatherAPI{
        cacheEngine: engine,
        executor:    scache.NewExecutor(engine),
    }
}

func (w *WeatherAPI) GetWeather(city string) (string, error) {
    cacheKey := fmt.Sprintf("weather:%s", city)

    // 尝试从缓存获取
    result, err := w.executor.Execute("GET", cacheKey)
    if err != nil {
        return "", err
    }

    if result != nil {
        fmt.Printf("缓存命中: %s 天气数据\n", city)
        return result.(string), nil
    }

    fmt.Printf("缓存未命中: 请求 %s 天气数据\n", city)

    // 模拟API调用
    weather, err := w.callWeatherAPI(city)
    if err != nil {
        return "", err
    }

    // 存入缓存
    _, err = w.executor.Execute("SET", cacheKey, weather, time.Minute*10)
    if err != nil {
        return "", err
    }

    return weather, nil
}

func (w *WeatherAPI) callWeatherAPI(city string) (string, error) {
    // 模拟网络延迟
    time.Sleep(time.Millisecond * 200)
    return fmt.Sprintf("%s: 晴天, 25°C", city), nil
}

func main() {
    weatherAPI := NewWeatherAPI()

    cities := []string{"北京", "上海", "广州", "深圳"}

    // 第一次请求 - 会调用API
    fmt.Println("=== 第一次请求 ===")
    for _, city := range cities {
        weather, err := weatherAPI.GetWeather(city)
        if err != nil {
            panic(err)
        }
        fmt.Println(weather)
    }

    // 第二次请求 - 从缓存获取
    fmt.Println("\n=== 第二次请求 ===")
    for _, city := range cities {
        weather, err := weatherAPI.GetWeather(city)
        if err != nil {
            panic(err)
        }
        fmt.Println(weather)
    }

    // 显示统计信息
    stats := weatherAPI.cacheEngine.Stats().(map[string]interface{})
    fmt.Printf("\n=== 缓存统计 ===\n")
    fmt.Printf("命中次数: %d\n", stats["hits"])
    fmt.Printf("未命中次数: %d\n", stats["misses"])
    fmt.Printf("命中率: %.2f%%\n", stats["hit_rate"].(float64)*100)
}
```

### 案例3: 分布式锁实现

```go
package main

import (
    "fmt"
    "sync"
    "time"

    "scache"
    "scache/config"
)

// DistributedLock 分布式锁
type DistributedLock struct {
    cacheEngine interfaces.StorageEngine
    executor    *scache.Executor
}

func NewDistributedLock() *DistributedLock {
    engine := scache.NewEngine(
        config.WithMaxSize(1000),
        config.WithDefaultExpiration(time.Second*30), // 锁30秒自动过期
    )
    return &DistributedLock{
        cacheEngine: engine,
        executor:    scache.NewExecutor(engine),
    }
}

// TryLock 尝试获取锁
func (d *DistributedLock) TryLock(key string, ttl time.Duration) bool {
    lockKey := fmt.Sprintf("lock:%s", key)

    // 尝试设置锁
    result, err := d.executor.Execute("SET", lockKey, "locked", ttl)
    if err != nil {
        return false
    }

    return result == nil // SET成功表示获取到锁
}

// Release 释放锁
func (d *DistributedLock) Release(key string) bool {
    lockKey := fmt.Sprintf("lock:%s", key)

    result, err := d.executor.Execute("DEL", lockKey)
    if err != nil {
        return false
    }

    return result.(bool)
}

// IsLocked 检查锁状态
func (d *DistributedLock) IsLocked(key string) bool {
    lockKey := fmt.Sprintf("lock:%s", key)

    result, err := d.executor.Execute("EXISTS", lockKey)
    if err != nil {
        return false
    }

    return result.(bool)
}

func main() {
    lock := NewDistributedLock()
    var wg sync.WaitGroup

    // 模拟5个协程竞争锁
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            resourceKey := "shared_resource"

            // 尝试获取锁
            if lock.TryLock(resourceKey, time.Second*5) {
                fmt.Printf("协程 %d: 获取到锁\n", id)

                // 模拟临界区操作
                time.Sleep(time.Second * 2)
                fmt.Printf("协程 %d: 完成操作\n", id)

                // 释放锁
                lock.Release(resourceKey)
                fmt.Printf("协程 %d: 释放锁\n", id)
            } else {
                fmt.Printf("协程 %d: 获取锁失败\n", id)
            }
        }(i)
    }

    wg.Wait()
    fmt.Println("所有协程完成")
}
```

这个教程涵盖了 SCache 的主要特性和用法，从基础概念到高级实战案例。通过这些示例，你应该能够：

1. 理解 SCache 的基本架构和使用方法
2. 掌握不同数据类型的操作技巧
3. 学会配置和优化缓存性能
4. 解决常见的缓存问题
5. 在实际项目中应用缓存策略

更多详细信息请参考 [API文档](API.md) 和 [项目README](../README.md)。