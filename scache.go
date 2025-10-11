// SCache - 高性能 Go 缓存框架
// 提供传统的实例化缓存和全局缓存两种使用方式
package scache

import (
	"github.com/scache/cache"
	"github.com/scache/constants"
	"github.com/scache/types"
)

// 重新导出核心接口和类型
type (
	// Cache 缓存接口
	Cache = cache.Cache
	// Config 配置结构
	Config = cache.Config
	// CacheStats 缓存统计信息
	CacheStats = cache.CacheStats
	// CacheItem 缓存项
	CacheItem = types.CacheItem
	// EvictionPolicy 淘汰策略接口
	EvictionPolicy = cache.EvictionPolicy
	// Serializer 序列化接口
	Serializer = cache.Serializer
	// EvictionPolicyType 淘汰策略类型
	EvictionPolicyType = cache.EvictionPolicyType
	// Manager 缓存管理器
	Manager = cache.Manager
	// CacheConfig 缓存配置详细结构
	CacheConfig = types.CacheConfig
	// ManagerStats 管理器统计信息
	ManagerStats = types.ManagerStats
	// ErrorInfo 错误信息
	ErrorInfo = types.ErrorInfo
	// ValidationResult 验证结果
	ValidationResult = types.ValidationResult
	// HealthStatus 健康状态
	HealthStatus = types.HealthStatus
)

// 重新导出配置选项函数
var (
	WithMaxSize         = cache.WithMaxSize
	WithDefaultTTL      = cache.WithDefaultTTL
	WithEvictionPolicy  = cache.WithEvictionPolicy
	WithShards          = cache.WithShards
	WithCleanupInterval = cache.WithCleanupInterval
	WithStatistics      = cache.WithStatistics
	WithSerializer      = cache.WithSerializer
	WithLazyExpiration  = cache.WithLazyExpiration
	WithMetrics         = cache.WithMetrics
)

// 传统缓存实例化方式

// New 创建一个新的缓存实例
func New(opts ...cache.Option) Cache {
	return cache.New(opts...)
}

// NewLRU 创建一个使用 LRU 策略的缓存
func NewLRU(maxSize int, opts ...cache.Option) Cache {
	allOpts := append([]cache.Option{cache.WithMaxSize(maxSize), cache.WithEvictionPolicy(constants.LRUStrategy)}, opts...)
	return cache.New(allOpts...)
}

// NewLFU 创建一个使用 LFU 策略的缓存
func NewLFU(maxSize int, opts ...cache.Option) Cache {
	allOpts := append([]cache.Option{cache.WithMaxSize(maxSize), cache.WithEvictionPolicy(constants.LFUStrategy)}, opts...)
	return cache.New(allOpts...)
}

// NewFIFO 创建一个使用 FIFO 策略的缓存
func NewFIFO(maxSize int, opts ...cache.Option) Cache {
	allOpts := append([]cache.Option{cache.WithMaxSize(maxSize), cache.WithEvictionPolicy(constants.FIFOStrategy)}, opts...)
	return cache.New(allOpts...)
}

// NewWithConfig 使用配置创建缓存
func NewWithConfig(config *Config) Cache {
	opts := []cache.Option{
		cache.WithMaxSize(config.MaxSize),
		cache.WithDefaultTTL(config.DefaultTTL),
		cache.WithEvictionPolicy(config.EvictionPolicy),
		cache.WithShards(config.Shards),
		cache.WithCleanupInterval(config.CleanupInterval),
		cache.WithStatistics(config.EnableStatistics),
		cache.WithSerializer(config.Serializer),
		cache.WithLazyExpiration(config.EnableLazyExpiration),
		cache.WithMetrics(config.EnableMetrics),
	}
	return cache.New(opts...)
}

// 全局缓存方式

// GetGlobalManager 获取全局缓存管理器
var GetGlobalManager = cache.GetGlobalManager

// 全局便捷函数
var (
	// Register 注册全局缓存
	Register = cache.Register
	// RegisterLRU 注册 LRU 缓存
	RegisterLRU = cache.RegisterLRU
	// RegisterLFU 注册 LFU 缓存
	RegisterLFU = cache.RegisterLFU
	// RegisterFIFO 注册 FIFO 缓存
	RegisterFIFO = cache.RegisterFIFO
	// Get 获取全局缓存
	Get = cache.Get
	// GetOrDefault 获取缓存，如果不存在则创建默认缓存
	GetOrDefault = cache.GetOrDefault
	// Remove 移除全局缓存
	Remove = cache.Remove
	// List 列出所有全局缓存名称
	List = cache.List
	// Clear 清空所有全局缓存
	Clear = cache.Clear
	// Close 关闭所有全局缓存
	Close = cache.Close
	// Stats 获取所有全局缓存的统计信息
	Stats = cache.Stats
	// Size 获取所有全局缓存的总大小
	Size = cache.Size
	// Exists 检查全局缓存是否存在
	Exists = cache.Exists
)

// 默认缓存的便捷操作
var (
	// Set 在默认缓存中设置键值
	Set = cache.Set
	// SetWithTTL 在默认缓存中设置带过期时间的键值
	SetWithTTL = cache.SetWithTTL
	// GetFromDefault 从默认缓存中获取值
	GetFromDefault = cache.GetFromDefault
	// DeleteFromDefault 从默认缓存中删除键
	DeleteFromDefault = cache.DeleteFromDefault
	// ExistsInDefault 检查默认缓存中是否存在键
	ExistsInDefault = cache.ExistsInDefault
	// ClearDefault 清空默认缓存
	ClearDefault = cache.ClearDefault
)

// 重新导出常量
var (
	// 缓存策略常量
	LRU  = cache.LRU
	LFU  = cache.LFU
	FIFO = cache.FIFO

	// 默认配置常量
	DefaultMaxSize         = constants.DefaultMaxSize
	DefaultShards          = constants.DefaultShards
	DefaultCleanupInterval = constants.DefaultCleanupInterval
	DefaultTTL             = constants.DefaultTTL

	// 错误消息常量
	ErrCacheNotFound      = constants.ErrCacheNotFound
	ErrInvalidCacheName   = constants.ErrInvalidCacheName
	ErrCacheAlreadyExists = constants.ErrCacheAlreadyExists
	ErrInvalidStrategy    = constants.ErrInvalidStrategy
	ErrKeyNotFound        = constants.ErrKeyNotFound
	ErrKeyTooLong         = constants.ErrKeyTooLong
	ErrKeyEmpty           = constants.ErrKeyEmpty
	ErrValueTooLarge      = constants.ErrValueTooLarge
	ErrCacheClosed        = constants.ErrCacheClosed
)
