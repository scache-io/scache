# SCache - Go 结构体缓存代码生成工具

SCache 是一个智能的 Go 结构体缓存代码生成工具，自动扫描项目中的结构体并生成对应的缓存操作方法。

## 核心特性

- **智能代码生成** - 自动扫描Go结构体，生成懒汉式单例缓存方法
- **TTL过期机制** - 支持灵活的缓存过期时间设置
- **LRU淘汰策略** - 智能的缓存淘汰机制，支持容量限制
- **多种数据类型** - 支持String、List、Hash、Struct等数据类型
- **线程安全** - 内置锁机制，支持并发访问
- **高性能** - 基于内存存储，读写性能优异

## 安装与快速开始

### 安装

```bash
go install github.com/scache-io/scache/cmd/scache@latest

# 或者从源码安装
git clone https://github.com/scache-io/scache.git
cd scache
go install ./cmd/scache
```

### 生成缓存代码

```bash
# 生成所有结构体的缓存代码
scache gen

# 指定目录生成
scache gen -dir ./models

# 只生成指定结构体
scache gen -structs User,Product
```

### 生成的代码特点

为每个结构体生成懒汉式单例缓存管理器：

- **Get{Struct}Scache()** - 获取懒汉式单例实例（推荐）
- **Store/Load** - 基础缓存操作方法
- **MustStore/MustLoad** - 带panic的错误处理版本
- **WithKey()** - 链式调用支持格式化key
- **StorePtr/LoadPtr** - 指针类型的缓存操作

## 缓存功能

### TTL过期机制 & LRU淘汰策略

```go
// 配置缓存
cache := scache.New(
config.WithMaxSize(10000), // LRU容量限制
config.WithBackgroundCleanup(5*time.Minute), // 自动清理
)

// 设置不同TTL
cache.Store("user:1001", user, time.Hour) // 1小时过期
cache.Store("config", config, 24*time.Hour) // 24小时过期
cache.Store("temp", data, time.Minute) // 1分钟过期
```

### 数据类型支持

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
config.WithMaxSize(50000), // 最多50000个键
config.WithMemoryThreshold(0.9), // 内存阈值90%
config.WithDefaultExpiration(24*time.Hour), // 默认过期时间24小时
config.WithBackgroundCleanup(15*time.Minute), // 后台清理间隔15分钟
)
```

## 实践案例与最佳实践

### 用户会话管理

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

### 数据库查询缓存

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
var article Article
err := scache.Load("article:"+articleID, &article)
if err == nil {
return &article, nil // 缓存命中
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
```

### API限流

```go
func AllowRequest(clientID string) bool {
key := fmt.Sprintf("rate_limit:%s", clientID)

count, exists := scache.GetString(key)
if !exists {
scache.SetString(key, "1", time.Minute)
return true
}

currentCount := 0
fmt.Sscanf(count, "%d", &currentCount)

if currentCount >= 100 { // 每分钟100次
return false
}

scache.SetString(key, fmt.Sprintf("%d", currentCount+1), time.Minute)
return true
}
```

### 最佳实践建议

#### 缓存Key设计规范

```go
user:1001         // 用户信息
user:1001:profile // 用户资料
article:1001   // 文章内容
session:abc123 // 会话信息
rate_limit:client_001 // 限流计数
```

#### TTL过期策略

```go
time.Minute // 验证码等实时数据
time.Hour   // 频繁更新的数据
24*time.Hour // 用户信息、配置
0            // 永不过期数据
```

## 开发指南

### 构建

```bash
git clone https://github.com/scache-io/scache.git
cd scache
go mod tidy
go test ./...
go build -o scache ./cmd/scache
```

### 贡献

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

MIT License - 查看 [LICENSE](LICENSE) 文件了解详情。
