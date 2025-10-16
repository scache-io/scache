# SCache API 文档

## 目录

- [核心接口](#核心接口)
- [数据类型](#数据类型)
- [命令系统](#命令系统)
- [配置管理](#配置管理)
- [存储引擎](#存储引擎)
- [统计信息](#统计信息)

## 核心接口

### StorageEngine

存储引擎是 SCache 的核心接口，提供了所有缓存操作的基础功能。

```go
type StorageEngine interface {
    // 基础操作
    Set(key string, obj DataObject) error
    Get(key string) (DataObject, bool)
    Delete(key string) bool
    Exists(key string) bool
    Keys() []string
    Flush() error
    Size() int

    // 类型检查
    Type(key string) (DataType, bool)

    // 过期管理
    Expire(key string, ttl time.Duration) bool
    TTL(key string) (time.Duration, bool)

    // 统计信息
    Stats() interface{}
}
```

### DataObject

所有缓存对象的基础接口。

```go
type DataObject interface {
    Type() DataType
    ExpiresAt() time.Time
    IsExpired() bool
    Size() int
}
```

## 数据类型

### StringObject

字符串类型对象，用于存储简单的字符串值。

```go
type StringObject interface {
    DataObject
    Value() string
    Set(value string)
}

// 创建字符串对象
obj := types.NewStringObject("hello", time.Hour)

// 获取值
value := obj.Value()

// 设置新值
obj.Set("world")
```

### ListObject

列表类型对象，用于存储有序的元素集合。

```go
type ListObject interface {
    DataObject
    Values() []interface{}
    Push(value interface{})
    Pop() (interface{}, bool)
    Index(index int) (interface{}, bool)
    Range(start, end int) []interface{}
    Len() int
}

// 创建列表对象
values := []interface{}{"item1", "item2", "item3"}
obj := types.NewListObject(values, time.Hour)

// 添加元素
obj.Push("item4")

// 弹出元素
value, exists := obj.Pop()

// 获取指定索引元素
item, exists := obj.Index(1)

// 获取范围元素
rangeItems := obj.Range(0, 2)
```

### HashObject

哈希类型对象，用于存储键值对集合。

```go
type HashObject interface {
    DataObject
    Fields() map[string]interface{}
    Get(field string) (interface{}, bool)
    Set(field string, value interface{})
    Delete(field string) bool
    Len() int
}

// 创建哈希对象
fields := map[string]interface{}{
    "name": "Alice",
    "age": 30,
}
obj := types.NewHashObject(fields, time.Hour)

// 设置字段
obj.Set("email", "alice@example.com")

// 获取字段
name, exists := obj.Get("name")

// 删除字段
deleted := obj.Delete("age")
```

## 命令系统

### Command

命令接口，所有命令都需要实现这个接口。

```go
type Command interface {
    Name() string
    Execute(ctx *Context) error
    Validate(args []interface{}) error
}
```

### Context

命令执行上下文，包含执行所需的所有信息。

```go
type Context struct {
    Storage StorageEngine
    Args    []interface{}
    Result  interface{}
    Error   error
}
```

### 内置命令

#### SetCommand

设置字符串值。

```go
// 命令格式: SET key value [ttl]
cmd := commands.NewSetCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mykey", "myvalue", time.Hour},
}
err := cmd.Execute(ctx)
```

#### GetCommand

获取字符串值。

```go
// 命令格式: GET key
cmd := commands.NewGetCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mykey"},
}
err := cmd.Execute(ctx)
value := ctx.Result
```

#### DeleteCommand

删除键。

```go
// 命令格式: DEL key
cmd := commands.NewDeleteCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mykey"},
}
err := cmd.Execute(ctx)
deleted := ctx.Result.(bool)
```

#### HSetCommand

设置哈希字段。

```go
// 命令格式: HSET key field value [ttl]
cmd := commands.NewHSetCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"user:1", "name", "Alice", time.Hour},
}
err := cmd.Execute(ctx)
```

#### HGetCommand

获取哈希字段。

```go
// 命令格式: HGET key field
cmd := commands.NewHGetCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"user:1", "name"},
}
err := cmd.Execute(ctx)
value := ctx.Result
```

#### LPushCommand

列表左侧推入元素。

```go
// 命令格式: LPUSH key value [ttl]
cmd := commands.NewLPushCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mylist", "item1", time.Hour},
}
err := cmd.Execute(ctx)
length := ctx.Result.(int)
```

#### RPopCommand

列表右侧弹出元素。

```go
// 命令格式: RPOP key
cmd := commands.NewRPopCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mylist"},
}
err := cmd.Execute(ctx)
value := ctx.Result
```

#### ExpireCommand

设置键的过期时间。

```go
// 命令格式: EXPIRE key ttl
cmd := commands.NewExpireCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mykey", time.Minute * 10},
}
err := cmd.Execute(ctx)
success := ctx.Result.(bool)
```

#### TTLCommand

获取键的剩余生存时间。

```go
// 命令格式: TTL key
cmd := commands.NewTTLCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mykey"},
}
err := cmd.Execute(ctx)
ttl := ctx.Result.(int) // 返回秒数，-1表示永不过期，-2表示不存在
```

#### TypeCommand

获取键的数据类型。

```go
// 命令格式: TYPE key
cmd := commands.NewTypeCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{"mykey"},
}
err := cmd.Execute(ctx)
dataType := ctx.Result.(string) // "string", "list", "hash", "none"
```

#### StatsCommand

获取缓存统计信息。

```go
// 命令格式: STATS
cmd := commands.NewStatsCommand()
ctx := &Context{
    Storage: engine,
    Args: []interface{}{},
}
err := cmd.Execute(ctx)
stats := ctx.Result.(map[string]interface{})
```

## 配置管理

### EngineConfig

存储引擎配置结构。

```go
type EngineConfig struct {
    MaxSize                   int           // 最大缓存数量
    MemoryThreshold           float64       // 内存阈值
    DefaultExpiration         time.Duration // 默认过期时间
    BackgroundCleanupInterval time.Duration // 后台清理间隔
}
```

### 配置选项

```go
type EngineOption func(*EngineConfig)

// 设置最大缓存数量
func WithMaxSize(size int) EngineOption

// 设置内存阈值
func WithMemoryThreshold(threshold float64) EngineOption

// 设置默认过期时间
func WithDefaultExpiration(ttl time.Duration) EngineOption

// 设置后台清理间隔
func WithBackgroundCleanup(interval time.Duration) EngineOption
```

### 预定义配置

```go
// 默认配置
var DefaultConfig = []EngineOption{
    WithMaxSize(0),                        // 无限制
    WithMemoryThreshold(0.8),              // 80%
    WithDefaultExpiration(0),              // 永不过期
    WithBackgroundCleanup(5 * time.Minute), // 5分钟清理
}

// 小型配置
var SmallConfig = []EngineOption{
    WithMaxSize(1000),                     // 1000个键
    WithMemoryThreshold(0.7),              // 70%
    WithDefaultExpiration(time.Hour),      // 1小时过期
    WithBackgroundCleanup(2 * time.Minute), // 2分钟清理
}

// 中等配置
var MediumConfig = []EngineOption{
    WithMaxSize(10000),                    // 10000个键
    WithMemoryThreshold(0.8),              // 80%
    WithDefaultExpiration(2 * time.Hour),  // 2小时过期
    WithBackgroundCleanup(5 * time.Minute), // 5分钟清理
}

// 大型配置
var LargeConfig = []EngineOption{
    WithMaxSize(100000),                   // 100000个键
    WithMemoryThreshold(0.85),             // 85%
    WithDefaultExpiration(6 * time.Hour),  // 6小时过期
    WithBackgroundCleanup(10 * time.Minute), // 10分钟清理
}
```

## 存储引擎

### NewStorageEngine

创建新的存储引擎实例。

```go
func NewStorageEngine(config *EngineConfig) StorageEngine

// 使用默认配置
engine := storage.NewStorageEngine(nil)

// 使用自定义配置
config := &storage.EngineConfig{
    MaxSize:                   10000,
    MemoryThreshold:           0.8,
    DefaultExpiration:         time.Hour,
    BackgroundCleanupInterval: 5 * time.Minute,
}
engine := storage.NewStorageEngine(config)
```

### 基本操作示例

```go
// 设置字符串对象
strObj := types.NewStringObject("hello", time.Hour)
err := engine.Set("greeting", strObj)

// 获取对象
obj, exists := engine.Get("greeting")
if exists {
    if strObj, ok := obj.(*types.StringObject); ok {
        value := strObj.Value()
        fmt.Println(value) // "hello"
    }
}

// 检查键是否存在
exists = engine.Exists("greeting")

// 获取键类型
dataType, exists := engine.Type("greeting")
fmt.Println(dataType) // "string"

// 设置过期时间
success := engine.Expire("greeting", time.Minute*30)

// 获取剩余生存时间
ttl, exists := engine.TTL("greeting")
fmt.Println(ttl) // 剩余秒数

// 删除键
deleted := engine.Delete("greeting")

// 获取所有键
keys := engine.Keys()

// 获取缓存大小
size := engine.Size()

// 清空所有缓存
err = engine.Flush()
```

## 统计信息

### 统计数据结构

```go
stats := engine.Stats()
statsMap := stats.(map[string]interface{})

// 统计字段包括：
// - hits:        int64  // 命中次数
// - misses:      int64  // 未命中次数
// - sets:        int64  // 设置次数
// - deletes:     int64  // 删除次数
// - evictions:   int64  // 淘汰次数
// - expirations: int64  // 过期次数
// - memory:      int64  // 内存使用量（字节）
// - keys:        int    // 当前键数量
// - hit_rate:    float64 // 命中率
```

### 使用示例

```go
stats := engine.Stats()
statsMap := stats.(map[string]interface{})

hits := statsMap["hits"].(int64)
misses := statsMap["misses"].(int64)
hitRate := statsMap["hit_rate"].(float64)
currentKeys := statsMap["keys"].(int)

fmt.Printf("命中率: %.2f%%\n", hitRate*100)
fmt.Printf("当前键数量: %d\n", currentKeys)
fmt.Printf("命中次数: %d\n", hits)
fmt.Printf("未命中次数: %d\n", misses)
```

## Executor 命令执行器

### NewExecutor

创建命令执行器。

```go
func NewExecutor(engine StorageEngine) *Executor

engine := storage.NewStorageEngine(nil)
executor := scache.NewExecutor(engine)
```

### Execute

执行命令。

```go
func (e *Executor) Execute(commandName string, args ...interface{}) (interface{}, error)

// 执行 SET 命令
result, err := executor.Execute("SET", "key", "value", time.Hour)

// 执行 GET 命令
result, err = executor.Execute("GET", "key")

// 执行 HSET 命令
result, err = executor.Execute("HSET", "user:1", "name", "Alice", time.Hour)
```

### RegisterCommand

注册自定义命令。

```go
func (e *Executor) RegisterCommand(cmd Command)

// 注册自定义命令
type MyCommand struct {
    commands.BaseCommand
}

func (c *MyCommand) Execute(ctx *interfaces.Context) error {
    ctx.Result = "Hello from custom command!"
    return nil
}

func (c *MyCommand) Name() string {
    return "MYCOMMAND"
}

executor.RegisterCommand(&MyCommand{})
```

## 便捷 API

SCache 提供了全局便捷函数，简化常用操作。

```go
// 字符串操作
err := scache.Set("key", "value", time.Hour)
value, found, err := scache.Get("key")

// 列表操作
length, err := scache.LPush("mylist", "item", time.Hour)
value, err := scache.RPop("mylist")

// 哈希操作
success, err := scache.HSet("myhash", "field", "value", time.Hour)
value, err = scache.HGet("myhash", "field")

// 通用操作
deleted, err := scache.Delete("key")
exists, err := scache.Exists("key")
keyType, err := scache.Type("key")
success, err := scache.Expire("key", time.Minute*30)
ttl, err := scache.TTL("key")

// 统计信息
stats := scache.Stats()
commands := scache.ListCommands()
```

## 错误处理

### 常见错误类型

```go
var (
    ErrUnknownCommand    = errors.New("unknown command")
    ErrKeyEmpty          = errors.New("key cannot be empty")
    ErrInvalidArgument   = errors.New("invalid argument")
    ErrTypeMismatch      = errors.New("type mismatch")
    ErrKeyNotFound       = errors.New("key not found")
    ErrFieldNotFound     = errors.New("field not found")
    ErrIndexOutOfRange   = errors.New("index out of range")
    ErrListEmpty         = errors.New("list is empty")
)
```

### 错误处理示例

```go
result, err := executor.Execute("SET", "", "value")
if err != nil {
    if err == scache.ErrKeyEmpty {
        fmt.Println("键不能为空")
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}

value, err := scache.Get("nonexistent")
if err != nil {
    fmt.Printf("获取错误: %v\n", err)
} else if !found {
    fmt.Println("键不存在")
} else {
    fmt.Printf("值: %v\n", value)
}
```

## 扩展开发

### 自定义数据类型

1. 定义新的数据类型接口

```go
type MyObject interface {
    DataObject
    MyMethod() string
}
```

2. 实现数据类型

```go
type myObject struct {
    *BaseObject
    data string
}

func NewMyObject(data string, ttl time.Duration) *myObject {
    return &myObject{
        BaseObject: NewBaseObject("mytype", ttl),
        data:       data,
    }
}

func (m *myObject) MyMethod() string {
    return m.data
}
```

3. 创建相关命令

```go
type MySetCommand struct {
    BaseCommand
}

func (c *MySetCommand) Execute(ctx *Context) error {
    // 实现命令逻辑
    return nil
}

func (c *MySetCommand) Name() string {
    return "MYSET"
}
```

4. 注册命令

```go
executor.RegisterCommand(NewMySetCommand())
```

### 自定义淘汰策略

实现 `EvictionPolicy` 接口：

```go
type MyPolicy struct {
    // 策略内部状态
}

func (p *MyPolicy) Access(key string) {
    // 访问时的处理逻辑
}

func (p *MyPolicy) Set(key string) {
    // 设置时的处理逻辑
}

func (p *MyPolicy) Delete(key string) {
    // 删除时的处理逻辑
}

func (p *MyPolicy) Evict() string {
    // 返回需要淘汰的键
    return ""
}

// 实现其他必需方法...
```