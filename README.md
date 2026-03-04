# SCache - Go 结构体缓存代码生成工具

SCache 是一个智能的 Go 结构体缓存代码生成工具，自动扫描项目中的结构体并生成对应的缓存操作方法。支持泛型版本（推荐）和传统版本。

## 🚀 核心特性

- **智能代码生成** - 自动扫描Go结构体，生成懒汉式单例缓存方法
- **泛型支持** - 基于Go 1.18+泛型，代码更简洁、类型更安全（推荐）
- **传统版本** - 兼容旧版Go，完整的缓存功能
- **TTL过期机制** - 支持灵活的缓存过期时间设置
- **LRU淘汰策略** - 智能的缓存淘汰机制，支持容量限制
- **多种数据类型** - 支持String、List、Hash、Struct等数据类型
- **线程安全** - 内置锁机制，支持并发访问
- **高性能** - 基于内存存储，读写性能优异

## 📦 安装

### 通过 Go 安装（推荐）

```bash
go install github.com/scache-io/scache/cmd/scache@latest
```

### 从源码安装

```bash
git clone https://github.com/scache-io/scache.git
cd scache
go install ./cmd/scache
```

## 🎯 快速开始

### 1. 生成泛型版本代码（推荐，Go 1.18+）

```bash
# 生成泛型版本缓存代码
scache gen --generic

# 指定目录生成
scache gen --generic -dir ./models

# 只生成指定结构体
scache gen --generic -structs User,Product
```

### 2. 生成传统版本代码（兼容旧版Go）

```bash
# 生成传统版本缓存代码（默认）
scache gen

# 指定目录生成
scache gen -dir ./models

# 只生成指定结构体
scache gen -structs User,Product
```

### 3. 生成代码示例

假设你有以下结构体：

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

**泛型版本生成**：

```go
// 类型特定的缓存实例定义
var (
    defaultUserScache *Scache[User]
    defaultUserScacheOnce sync.Once

    defaultProductScache *Scache[Product]
    defaultProductScacheOnce sync.Once
)

// 泛型缓存管理器
type Scache[T any] struct {
    cache *scache.LocalCache
}

// 便捷函数
func GetUserScache() *Scache[User] { /* ... */ }
func GetProductScache() *Scache[Product] { /* ... */ }
```

**传统版本生成**：

```go
// 每个结构体独立的缓存管理器
type UserScache struct { /* ... */ }
type ProductScache struct { /* ... */ }

// 独立的便捷函数
func GetUserScache() *UserScache { /* ... */ }
func GetProductScache() *ProductScache { /* ... */ }
```

## 🔧 使用方法

### 泛型版本使用（推荐）

```go
package main

import (
    "fmt"
    "time"
    "yourproject/cache" // 生成的缓存包
)

func main() {
    // 获取用户缓存实例
    userCache := cache.GetUserScache()

    // 存储用户数据
    user := User{ID: 1, Name: "张三", Age: 25}
    err := userCache.Store("user:1", user, time.Hour)
    if err != nil {
        panic(err)
    }

    // 读取用户数据
    loadedUser, err := userCache.Load("user:1")
    if err != nil {
        panic(err)
    }
    fmt.Printf("加载的用户: %+v\n", loadedUser)

    // 获取产品缓存实例
    productCache := cache.GetProductScache()

    // 存储产品数据
    product := Product{ID: "p001", Name: "iPhone", Price: 999.99}
    err = productCache.Store("product:p001", product, 2*time.Hour)
    if err != nil {
        panic(err)
    }

    // 检查是否存在
    if productCache.Exists("product:p001") {
        fmt.Println("产品缓存存在")
    }

    // 获取剩余生存时间
    if ttl, exists := productCache.GetTTL("product:p001"); exists {
        fmt.Printf("产品缓存剩余时间: %v\n", ttl)
    }
}
```

### 传统版本使用

```go
package main

import (
    "fmt"
    "time"
    "yourproject/cache" // 生成的缓存包
)

func main() {
    // 获取用户缓存实例
    userCache := cache.GetUserScache()

    // 存储用户数据
    user := User{ID: 1, Name: "张三", Age: 25}
    err := userCache.Store("user:1", user, time.Hour)
    if err != nil {
        panic(err)
    }

    // 读取用户数据
    loadedUser, err := userCache.Load("user:1")
    if err != nil {
        panic(err)
    }
    fmt.Printf("加载的用户: %+v\n", loadedUser)
}
```

## 📦 导入方式

### 推荐方式（v2.0+）

```go
// 导入 pkg 包
import scache "github.com/scache-io/scache/pkg"

// 使用
cache := scache.New(config.MediumConfig...)
scache.SetString("key", "value", time.Hour)
```

### 向后兼容方式（保持原有 API）

```go
// 导入根目录包（自动重新导出 pkg 内容）
import "github.com/scache-io/scache"

// 使用方式不变
cache := scache.New(config.MediumConfig...)
scache.SetString("key", "value", time.Hour)
```

### API 迁移指南

**从旧版本迁移到 v2.0**：

1. 更新导入路径（可选，旧方式仍然支持）：
   ```go
   // 旧方式（仍然有效）
   import "github.com/scache-io/scache"

   // 新方式（推荐）
   import scache "github.com/scache-io/scache/pkg"
   ```

2. API 完全兼容，无需修改代码

3. 功能增强：
   - 更清晰的代码结构
   - 更好的类型安全
   - 更高的可维护性

## 🎛️ 命令行选项

### 基本选项

```bash
# 查看帮助
scache gen --help

# 基本用法
scache gen [flags]

# 标志说明
  -d, --dir string          项目目录路径 (默认 ".")
  -p, --package string      包名（默认为目录名）
  -e, --exclude string      排除的目录，用逗号分隔 (默认 "vendor,node_modules,.git")
  -s, --structs string      指定结构体名称，用逗号分隔（默认生成所有）
  --generic                 使用泛型版本（支持Go 1.18+）
```

### 使用示例

```bash
# 生成当前目录所有结构体的泛型版本
scache gen --generic

# 生成指定目录的传统版本
scache gen -dir ./models -package myapp

# 只生成指定结构体的泛型版本
scache gen --generic -structs User,Product,Order

# 排除特定目录
scache gen --generic -exclude "vendor,test,docs"
```

## 🏗️ 生成代码结构

### 泛型版本结构

```
生成的代码文件结构：
├── 类型特定的缓存实例定义（最顶部）
│   ├── defaultXXXScache 变量声明
│   └── defaultXXXScacheOnce 同步锁
├── Scache[T] 泛型缓存管理器
├── 构造函数
│   └── NewScache[T]()
├── 类型特定的便捷函数
│   ├── GetXXXScache() 单例获取
│   └── NewXXXScache() 便捷构造
└── 核心方法（按功能分组）
    ├── 存储读取操作 (Store, Load)
    ├── 键管理操作 (Delete, Exists)
    ├── 过期时间管理 (SetTTL, GetTTL)
    └── 缓存管理操作 (Clear, Size, Keys)
```

### 传统版本结构

```
生成的代码文件结构：
├── 结构体定义和单例变量
├── XXXScache 结构体缓存管理器
├── 构造函数 NewXXXScache()
├── 核心方法
│   ├── Store/Load (基础操作)
│   ├── Delete/Clear/Size/Keys/Exists (管理操作)
│   └── SetTTL/GetTTL (过期时间管理)
```

## 🆚 泛型版本 vs 传统版本

| 特性 | 泛型版本 | 传统版本 |
|------|----------|----------|
| **代码量** | ⭐ 极少 (40行核心) | ⭐⭐ 较多 (96行/结构体) |
| **类型安全** | ⭐⭐⭐ 编译时检查 | ⭐⭐ 运行时检查 |
| **性能** | ⭐⭐⭐ 优秀 | ⭐⭐⭐ 优秀 |
| **内存占用** | ⭐⭐⭐ 极低 | ⭐⭐ 较高 |
| **Go版本要求** | Go 1.18+ | Go 1.10+ |
| **API一致性** | ⭐⭐⭐ 统一接口 | ⭐⭐ 独立接口 |
| **维护性** | ⭐⭐⭐ 极佳 | ⭐⭐ 良好 |

**推荐选择**：
- ✅ **新项目** → 使用泛型版本
- ✅ **Go 1.18+** → 使用泛型版本
- ⚠️ **旧项目兼容** → 使用传统版本
- ⚠️ **Go < 1.18** → 使用传统版本

## 📁 项目结构

```
scache/
├── cmd/scache/              # 命令行工具
│   ├── main.go
│   ├── generator/           # 代码生成器
│   └── go.mod
│
├── pkg/                     # 公共 API 层
│   ├── api/                 # API 实现
│   │   ├── global.go        # 全局缓存 API
│   │   └── local.go         # 本地缓存 API
│   └── scache.go            # 包入口
│
├── cache/                   # 缓存核心
│   ├── cache.go
│   └── cache_test.go
│
├── storage/                 # 存储引擎
│   ├── engine.go
│   └── *_test.go
│
├── types/                   # 数据类型
│   ├── data_objects.go
│   └── *_test.go
│
├── policies/                # 淘汰策略
│   └── lru/
│       ├── lru.go
│       └── lru_test.go
│
├── config/                  # 配置
│   └── engine_config.go
│
├── constants/               # 常量
│   └── cache_constants.go
│
├── errors/                  # 错误定义
│   └── errors.go
│
├── utils/                   # 工具函数
│   ├── ttl_helper.go
│   ├── type_helper.go
│   └── validation.go
│
├── internal/                # 内部包
│   └── memory_checker.go   # 内存监控
│
├── interfaces/              # 接口定义
│   └── interfaces.go
│
├── scache_export.go         # 向后兼容导出
└── scache_test.go          # 根目录测试
```

## 🎯 缓存功能

### TTL过期机制 & LRU淘汰策略

```go
// 创建自定义配置的缓存
cache := scache.New(
    config.WithMaxSize(10000), // LRU容量限制
    config.WithBackgroundCleanup(5*time.Minute), // 自动清理
)

// 设置不同TTL
cache.Store("user:1001", user, time.Hour) // 1小时过期
cache.Store("config", config, 24*time.Hour) // 24小时过期
cache.Store("temp", data, time.Minute) // 1分钟过期
cache.Store("permanent", data) // 永不过期
```

### 数据类型支持

#### 局部缓存 (LocalCache)

```go
package main

import (
    "fmt"
    "time"

    scache "github.com/scache-io/scache/pkg"  // 推荐方式
    // 或
    // "github.com/scache-io/scache"  // 向后兼容
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

    scache "github.com/scache-io/scache/pkg"  // 推荐方式
    // 或
    // "github.com/scache-io/scache"  // 向后兼容
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

### 缓存配置

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

## 🎨 实践案例与最佳实践

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
user:1001           // 用户信息
user:1001:profile   // 用户资料
article:1001        // 文章内容
session:abc123      // 会话信息
rate_limit:client_001 // 限流计数
```

#### TTL过期策略

```go
time.Minute    // 验证码等实时数据
time.Hour      // 频繁更新的数据
24*time.Hour   // 用户信息、配置
0              // 永不过期数据
```

#### 迁移指南

**从传统版本迁移到泛型版本**：

1. 确保Go版本 >= 1.18
2. 重新生成代码：`scache gen --generic`
3. 更新导入：`cache.GetUserScache()` → `cache.GetUserScache()`
4. API保持兼容，无需修改业务逻辑

## 🔧 开发指南

### 构建

```bash
git clone https://github.com/scache-io/scache.git
cd scache
go mod tidy
go test ./...
go build -o scache ./cmd/scache
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./cache
go test ./storage
go test ./types

# 运行基准测试
go test -bench=. ./cache
```

### 贡献

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📋 更新日志

### v2.0.0 (最新)
- ✨ 新增泛型版本支持（Go 1.18+）
- 🎯 大幅简化模板代码，减少70%重复代码
- 🔄 优化文件生成逻辑，支持覆盖写入
- 📚 改进文档和使用示例
- 🗂️ 重构代码结构，更清晰的组织方式
  - 新增 `pkg/` 公共 API 层
  - 新增 `errors/` 独立错误包
  - 新增 `utils/` 工具函数包
  - 优化 `internal/` 包的使用规范
  - 保持完全向后兼容

### v1.x.x
- 🎉 初始版本发布
- ✨ 基础缓存功能
- ✨ TTL和LRU支持
- ✨ 多数据类型支持

## 📄 许可证

MIT License - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🤝 贡献者

感谢所有为 SCache 项目做出贡献的开发者！

---

**SCache** - 让 Go 缓存开发更简单！ 🚀