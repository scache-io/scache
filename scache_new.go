// Package scache 是一个高性能的 Go 语言缓存库，类似 Redis 架构
//
// 特性：
//   - 支持多种数据类型：String, List, Hash
//   - 支持 TTL (Time To Live) 过期时间
//   - 支持 LRU (Least Recently Used) 淘汰策略
//   - 支持命令模式，易于扩展
//   - 线程安全，支持高并发访问
//   - 支持后台清理过期项目
//   - 提供统计信息（命中率、操作次数等）
//
// 基本使用：
//
//	// 创建存储引擎
//	engine := scache.NewEngine(
//		scache.WithMaxSize(1000),
//		scache.WithDefaultExpiration(time.Hour),
//	)
//
//	// 使用命令执行器
//	executor := scache.NewExecutor(engine)
//
//	// SET 操作
//	result, err := executor.Execute("SET", "key", "value", time.Minute*10)
//
//	// GET 操作
//	result, err = executor.Execute("GET", "key")
//
// 支持的命令：
//   - SET key value [ttl] - 设置字符串值
//   - GET key - 获取字符串值
//   - DEL key - 删除键
//   - EXISTS key - 检查键是否存在
//   - TYPE key - 获取键的类型
//   - EXPIRE key ttl - 设置过期时间
//   - TTL key - 获取剩余生存时间
//   - LPUSH key value [ttl] - 列表头部插入
//   - RPOP key - 列表尾部弹出
//   - HSET key field value [ttl] - 哈希字段设置
//   - HGET key field - 哈希字段获取
//   - STATS - 获取统计信息
package scache

import (
	"sync"
	"time"

	"github.com/scache-io/scache/commands"
	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/storage"
)

// Engine 存储引擎别名
type Engine = interfaces.StorageEngine

// NewEngine 创建新的存储引擎
func NewEngine(opts ...config.EngineOption) Engine {
	engineConfig := storage.DefaultEngineConfig()
	for _, opt := range opts {
		opt(engineConfig)
	}
	return storage.NewStorageEngine(engineConfig)
}

// Executor 命令执行器
type Executor struct {
	storage  Engine
	registry *commands.CommandRegistry
}

// NewExecutor 创建新的命令执行器
func NewExecutor(engine Engine) *Executor {
	return &Executor{
		storage:  engine,
		registry: commands.DefaultRegistry(),
	}
}

// Execute 执行命令
func (e *Executor) Execute(commandName string, args ...interface{}) (interface{}, error) {
	cmd, exists := e.registry.Get(commandName)
	if !exists {
		return nil, ErrUnknownCommand
	}

	ctx := &interfaces.Context{
		Storage: e.storage,
		Args:    args,
	}

	if err := cmd.Execute(ctx); err != nil {
		return nil, err
	}

	return ctx.Result, nil
}

// RegisterCommand 注册自定义命令
func (e *Executor) RegisterCommand(cmd interfaces.Command) {
	e.registry.Register(cmd)
}

// ListCommands 列出所有可用命令
func (e *Executor) ListCommands() []string {
	return e.registry.List()
}

// Stats 获取统计信息
func (e *Executor) Stats() interface{} {
	return e.storage.Stats()
}

// Close 关闭执行器
func (e *Executor) Close() {
	if closer, ok := e.storage.(interface{ Close() }); ok {
		closer.Close()
	}
}

// 全局默认实例
var (
	defaultExecutor *Executor
	defaultOnce     sync.Once
)

// GetGlobalExecutor 获取全局执行器实例（线程安全）
func GetGlobalExecutor() *Executor {
	defaultOnce.Do(func() {
		// 初始化默认实例，使用中等配置
		engine := NewEngine(config.MediumConfig...)
		defaultExecutor = NewExecutor(engine)
	})
	return defaultExecutor
}

// Execute 使用默认实例执行命令
func Execute(commandName string, args ...interface{}) (interface{}, error) {
	return GetGlobalExecutor().Execute(commandName, args...)
}

// Set 使用默认实例设置值
func Set(key string, value interface{}, ttl ...time.Duration) error {
	var args []interface{}
	args = append(args, key, value)
	if len(ttl) > 0 {
		args = append(args, ttl[0])
	}
	_, err := GetGlobalExecutor().Execute("SET", args...)
	return err
}

// Get 使用默认实例获取值
func Get(key string) (interface{}, bool, error) {
	result, err := GetGlobalExecutor().Execute("GET", key)
	if err != nil {
		return nil, false, err
	}
	if result == nil {
		return nil, false, nil
	}
	return result, true, nil
}

// Delete 使用默认实例删除键
func Delete(key string) (bool, error) {
	result, err := GetGlobalExecutor().Execute("DEL", key)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// Exists 使用默认实例检查键是否存在
func Exists(key string) (bool, error) {
	result, err := GetGlobalExecutor().Execute("EXISTS", key)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// Type 使用默认实例获取键类型
func Type(key string) (string, error) {
	result, err := GetGlobalExecutor().Execute("TYPE", key)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Expire 使用默认实例设置过期时间
func Expire(key string, ttl time.Duration) (bool, error) {
	result, err := GetGlobalExecutor().Execute("EXPIRE", key, ttl)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// TTL 使用默认实例获取剩余生存时间
func TTL(key string) (int, error) {
	result, err := GetGlobalExecutor().Execute("TTL", key)
	if err != nil {
		return -3, err
	}
	return result.(int), nil
}

// LPush 使用默认实例在列表头部插入
func LPush(key string, value interface{}, ttl ...time.Duration) (int, error) {
	var args []interface{}
	args = append(args, key, value)
	if len(ttl) > 0 {
		args = append(args, ttl[0])
	}
	result, err := GetGlobalExecutor().Execute("LPUSH", args...)
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}

// RPop 使用默认实例从列表尾部弹出
func RPop(key string) (interface{}, error) {
	return GetGlobalExecutor().Execute("RPOP", key)
}

// HSet 使用默认实例设置哈希字段
func HSet(key, field string, value interface{}, ttl ...time.Duration) (int, error) {
	var args []interface{}
	args = append(args, key, field, value)
	if len(ttl) > 0 {
		args = append(args, ttl[0])
	}
	result, err := GetGlobalExecutor().Execute("HSET", args...)
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}

// HGet 使用默认实例获取哈希字段
func HGet(key, field string) (interface{}, error) {
	return GetGlobalExecutor().Execute("HGET", key, field)
}

// Stats 使用默认实例获取统计信息
func Stats() interface{} {
	return GetGlobalExecutor().Stats()
}

// ListCommands 列出所有可用命令
func ListCommands() []string {
	return GetGlobalExecutor().ListCommands()
}

// RegisterCommand 注册自定义命令到默认实例
func RegisterCommand(cmd interfaces.Command) {
	GetGlobalExecutor().RegisterCommand(cmd)
}
