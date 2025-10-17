# SCache 代码生成工具

## 快速开始

SCache 是一个智能的 Go 结构体缓存代码生成工具，自动扫描项目中的结构体并生成对应的缓存操作方法。

## 安装

```bash
go install github.com/scache-io/scache/cmd/scache@latest
```

## 使用方法

### 基本用法

```bash
# 在当前目录生成所有结构体的缓存代码
scache gen

# 指定目录生成
scache gen -dir ./models

# 只生成指定结构体
scache gen -structs User,Product
```

## 生成结果

工具默认按结构体分包生成，生成目录结构如下：

```
project/
├── scache/                 # 主包
│   └── cache.go           # 工具函数
├── scache/user/            # User结构体包
│   └── cache.go           # User相关方法
├── scache/product/         # Product结构体包
│   └── cache.go           # Product相关方法
└── ...
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

## 特性

- ✅ **智能扫描** - 自动发现所有Go结构体
- ✅ **分包生成** - 每个结构体独立包，便于维护
- ✅ **类型安全** - 支持指针和值类型
- ✅ **错误处理** - 完整的错误处理机制
- ✅ **高性能** - 基于JSON序列化
- ✅ **可配置** - 灵活的结构体选择

## 最佳实践

1. **大型项目** - 利用分包模式保持代码整洁
2. **增量开发** - 使用 `-structs` 只生成需要的结构体
3. **团队协作** - 分包模式减少代码冲突

## 许可证

MIT License