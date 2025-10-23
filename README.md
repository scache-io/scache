# SCache - Go 结构体缓存代码生成工具

## 项目概述

SCache 是一个智能的 Go 结构体缓存代码生成工具，自动扫描项目中的结构体并生成对应的缓存操作方法。项目采用模块化设计，代码结构清晰，便于维护和扩展。

SCache 提供两种使用方式：
1. **代码生成器** - 自动为你的结构体生成专用的缓存操作方法
2. **缓存库** - 提供完整的局部缓存和全局缓存功能，支持多种数据类型

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

## 核心功能

### 1. 代码生成器

自动扫描项目中的结构体，生成专用的缓存操作方法。为每个结构体生成以下方法：

- `Store{Struct}` - 存储结构体
- `Load{Struct}` - 加载结构体
- `MustStore{Struct}` - 存储（失败时panic）
- `MustLoad{Struct}` - 加载（失败时panic）
- `Store{Struct}WithKey` - 格式化key存储
- `Load{Struct}WithKey` - 格式化key加载

### 2. 缓存库

提供完整的缓存功能，支持多种数据类型：

#### 局部缓存 (LocalCache)

```go
package main

import (
    "fmt"
    "time"
    "github.com/scache-io/scache"
    "github.com/scache-io/scache/config"
)

func main() {
    // 创建局部缓存实例
    cache := scache.New(config.MediumConfig...)

    // 字符串操作
    err := cache.SetString("user:name", "张三", time.Hour)
    if err != nil {
        panic(err)
    }

    name, exists := cache.GetString("user:name")
    if exists {
        fmt.Printf("用户名: %s\n", name)
    }

    // 结构体操作
    type User struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    user := User{Name: "李四", Age: 30}
    err = cache.Store("user:1001", &user, 2*time.Hour)
    if err != nil {
        panic(err)
    }

    var loadedUser User
    err = cache.Load("user:1001", &loadedUser)
    if err != nil {
        panic(err)
    }
    fmt.Printf("加载的用户: %+v\n", loadedUser)

    // 列表操作
    tags := []interface{}{"Go", "缓存", "高性能"}
    err = cache.SetList("tags:go", tags, time.Hour)
    if err != nil {
        panic(err)
    }

    loadedTags, exists := cache.GetList("tags:go")
    if exists {
        fmt.Printf("标签: %v\n", loadedTags)
    }

    // 哈希操作
    profile := map[string]interface{}{
        "email": "user@example.com",
        "phone": "13800138000",
        "city":  "北京",
    }
    err = cache.SetHash("profile:1001", profile, time.Hour)
    if err != nil {
        panic(err)
    }

    loadedProfile, exists := cache.GetHash("profile:1001")
    if exists {
        fmt.Printf("用户资料: %v\n", loadedProfile)
    }
}
```

#### 全局缓存 (GlobalCache)

```go
package main

import (
    "fmt"
    "time"
    "github.com/scache-io/scache"
    "github.com/scache-io/scache/config"
)

func init() {
    // 初始化全局缓存（可选，不调用则使用默认配置）
    scache.InitGlobalCache(config.LargeConfig...)
}

func main() {
    // 直接使用全局缓存函数
    err := scache.SetString("global:counter", "42", time.Hour)
    if err != nil {
        panic(err)
    }

    counter, exists := scache.GetString("global:counter")
    if exists {
        fmt.Printf("计数器: %s\n", counter)
    }

    // 全局缓存操作都是线程安全的
    go func() {
        scache.SetString("concurrent:test", "goroutine 1", time.Minute)
    }()

    go func() {
        scache.SetString("concurrent:test", "goroutine 2", time.Minute)
    }()

    time.Sleep(100 * time.Millisecond)
    value, _ := scache.GetString("concurrent:test")
    fmt.Printf("并发测试结果: %s\n", value)
}
```

### 3. 缓存配置

SCache 提供多种预定义配置：

```go
// 小型配置（适用于内存较小的环境）
cache := scache.New(config.SmallConfig...)

// 中等配置（适用于一般应用，默认配置）
cache := scache.New(config.MediumConfig...)

// 大型配置（适用于高负载应用）
cache := scache.New(config.LargeConfig...)

// 自定义配置
cache := scache.New(
    config.WithMaxSize(50000),                       // 最多50000个键
    config.WithMemoryThreshold(0.9),                 // 内存阈值90%
    config.WithDefaultExpiration(24*time.Hour),      // 默认过期时间24小时
    config.WithBackgroundCleanup(15*time.Minute),    // 后台清理间隔15分钟
)
```

## 优秀使用案例

### 案例1：用户会话管理

```go
type Session struct {
    UserID    string    `json:"user_id"`
    Username  string    `json:"username"`
    LoginTime time.Time `json:"login_time"`
    LastSeen  time.Time `json:"last_seen"`
    Role      string    `json:"role"`
}

func CreateSession(userID, username string) error {
    session := Session{
        UserID:    userID,
        Username:  username,
        LoginTime: time.Now(),
        LastSeen:  time.Now(),
        Role:      "user",
    }

    // 使用全局缓存存储会话，24小时过期
    return scache.Store("session:"+userID, &session, 24*time.Hour)
}

func GetSession(userID string) (*Session, error) {
    var session Session
    err := scache.Load("session:"+userID, &session)
    if err != nil {
        return nil, err
    }

    // 更新最后访问时间
    session.LastSeen = time.Now()
    scache.Store("session:"+userID, &session, 24*time.Hour)

    return &session, nil
}
```

### 案例2：数据库查询结果缓存

```go
type Article struct {
    ID       string    `json:"id"`
    Title    string    `json:"title"`
    Content  string    `json:"content"`
    Author   string    `json:"author"`
    CreateAt time.Time `json:"create_at"`
    Views    int       `json:"views"`
}

func GetArticle(articleID string) (*Article, error) {
    // 先尝试从缓存获取
    var article Article
    err := scache.Load("article:"+articleID, &article)
    if err == nil {
        return &article, nil
    }

    // 缓存未命中，从数据库查询
    article, err = queryArticleFromDB(articleID)
    if err != nil {
        return nil, err
    }

    // 存入缓存，1小时过期
    scache.Store("article:"+articleID, &article, time.Hour)
    return &article, nil
}

func GetPopularArticles() ([]Article, error) {
    cacheKey := "articles:popular"

    var articles []Article
    err := scache.Load(cacheKey, &articles)
    if err == nil {
        return articles, nil
    }

    // 查询热门文章
    articles, err = queryPopularArticlesFromDB()
    if err != nil {
        return nil, err
    }

    // 缓存30分钟
    scache.Store(cacheKey, &articles, 30*time.Minute)
    return articles, nil
}
```

### 案例3：API限流和计数

```go
type RateLimiter struct {
    window time.Duration
    limit  int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
    return &RateLimiter{window: window, limit: limit}
}

func (r *RateLimiter) Allow(clientID string) bool {
    key := fmt.Sprintf("rate_limit:%s", clientID)

    // 获取当前计数
    count, exists := scache.GetString(key)
    if !exists {
        scache.SetString(key, "1", r.window)
        return true
    }

    currentCount := 0
    fmt.Sscanf(count, "%d", &currentCount)

    if currentCount >= r.limit {
        return false
    }

    // 增加计数
    scache.SetString(key, fmt.Sprintf("%d", currentCount+1), r.window)
    return true
}

// 使用示例
func APIMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    limiter := NewRateLimiter(time.Minute, 100) // 每分钟100次请求

    return func(w http.ResponseWriter, r *http.Request) {
        clientID := getClientID(r)

        if !limiter.Allow(clientID) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        handler(w, r)
    }
}
```

### 案例4：分布式锁实现

```go
type DistributedLock struct {
    key        string
    ttl        time.Duration
    localCache *scache.LocalCache
}

func NewDistributedLock(key string, ttl time.Duration) *DistributedLock {
    return &DistributedLock{
        key:        key,
        ttl:        ttl,
        localCache: scache.New(config.SmallConfig...),
    }
}

func (l *DistributedLock) TryLock() bool {
    // 检查锁是否已被获取
    if _, exists := l.localCache.GetString(l.key); exists {
        return false
    }

    // 尝试获取锁
    err := l.localCache.SetString(l.key, "locked", l.ttl)
    return err == nil
}

func (l *DistributedLock) Unlock() {
    l.localCache.Delete(l.key)
}

func (l *DistributedLock) IsLocked() bool {
    _, exists := l.localCache.GetString(l.key)
    return exists
}

// 使用示例
func ProcessTask(taskID string) {
    lock := NewDistributedLock("lock:task:"+taskID, 5*time.Minute)

    if !lock.TryLock() {
        fmt.Printf("任务 %s 正在被其他实例处理\n", taskID)
        return
    }

    defer lock.Unlock()

    // 处理任务
    fmt.Printf("开始处理任务 %s\n", taskID)
    time.Sleep(2 * time.Second)
    fmt.Printf("任务 %s 处理完成\n", taskID)
}
```

### 案例5：缓存预热和失效处理

```go
type CacheWarmer struct {
    cache *scache.LocalCache
}

func NewCacheWarmer() *CacheWarmer {
    return &CacheWarmer{
        cache: scache.New(config.LargeConfig...),
    }
}

func (w *CacheWarmer) WarmupUserData(userIDs []string) error {
    for _, userID := range userIDs {
        // 预加载用户基本信息
        user, err := getUserFromDB(userID)
        if err != nil {
            continue
        }
        w.cache.Store("user:"+userID, user, time.Hour)

        // 预加载用户权限
        permissions, err := getUserPermissionsFromDB(userID)
        if err != nil {
            continue
        }
        w.cache.SetHash("permissions:"+userID, permissions, time.Hour)

        // 预加载用户偏好设置
        preferences, err := getUserPreferencesFromDB(userID)
        if err != nil {
            continue
        }
        w.cache.SetHash("preferences:"+userID, preferences, 24*time.Hour)
    }
    return nil
}

func (w *CacheWarmer) InvalidateUserCache(userID string) {
    // 删除用户相关的所有缓存
    keys := []string{
        "user:" + userID,
        "permissions:" + userID,
        "preferences:" + userID,
        "user:profile:" + userID,
        "user:stats:" + userID,
    }

    for _, key := range keys {
        w.cache.Delete(key)
    }

    // 可选：重新预热热点数据
    go w.WarmupUserData([]string{userID})
}
```

## 数据类型支持

SCache 支持以下数据类型：

| 数据类型 | 说明 | 用途 |
|---------|------|------|
| **String** | 字符串类型 | 存储简单文本、JSON数据等 |
| **List** | 列表类型 | 存储数组、队列等 |
| **Hash** | 哈希类型 | 存储键值对、对象属性等 |
| **Struct** | 结构体类型 | 存储复杂对象，自动JSON序列化 |

## 缓存策略

SCache 内置 LRU (Least Recently Used) 缓存淘汰策略：

- **访问时更新** - 每次访问都会更新元素的访问时间
- **容量限制** - 可配置最大缓存项数量
- **自动淘汰** - 超出容量时自动淘汰最少使用的项
- **内存管理** - 支持内存阈值监控

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

## 性能特性

- ✅ **高性能** - 基于内存存储，读写性能优异
- ✅ **线程安全** - 支持并发访问，内置锁机制
- ✅ **智能扫描** - 自动发现所有Go结构体
- ✅ **模块化设计** - 清晰的目录结构和职责分离
- ✅ **类型安全** - 强类型支持，编译时检查
- ✅ **错误处理** - 完整的错误处理机制
- ✅ **可配置** - 灵活的配置选项和策略
- ✅ **易于测试** - 完整的单元测试覆盖
- ✅ **可扩展** - 支持自定义缓存策略和存储引擎

## 最佳实践

### 1. 缓存Key设计

```go
// 推荐的Key命名规范
user:1001                    // 用户信息
user:1001:profile           // 用户资料
user:1001:permissions       // 用户权限
article:1001                // 文章内容
article:popular             // 热门文章列表
cache:api:user:list         // API缓存
session:abc123              // 会话信息
rate_limit:client_001       // 限流计数
```

### 2. 过期时间策略

```go
// 根据数据特性设置不同的过期时间
time.Minute         // 实时性要求高的数据（如验证码）
time.Hour           // 频繁更新的数据（如计数器）
24 * time.Hour      // 用户信息、配置数据
7 * 24 * time.Hour  // 统计数据、报告
0                  // 永不过期的配置数据
```

### 3. 错误处理

```go
// 推荐的错误处理方式
func GetUser(id string) (*User, error) {
    var user User
    err := scache.Load("user:"+id, &user)
    if err != nil {
        // 记录缓存未命中日志
        log.Printf("Cache miss for user %s: %v", id, err)

        // 从数据库获取
        user, err = getUserFromDB(id)
        if err != nil {
            return nil, err
        }

        // 存入缓存
        scache.Store("user:"+id, &user, time.Hour)
    }
    return &user, nil
}
```

### 4. 缓存预热

```go
// 应用启动时预热关键数据
func warmupCache() {
    // 预热热门配置
    go func() {
        configs := getHotConfigs()
        for _, config := range configs {
            scache.Store("config:"+config.ID, &config, 24*time.Hour)
        }
    }()

    // 预热热门内容
    go func() {
        articles := getPopularArticles()
        scache.Store("articles:popular", &articles, time.Hour)
    }()
}
```

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

### 安装命令行工具

```bash
# 通过 go install 安装
go install github.com/scache-io/scache/cmd/scache@latest

# 或者从源码安装
git clone https://github.com/scache-io/scache.git
cd scache
go install ./cmd/scache
```
