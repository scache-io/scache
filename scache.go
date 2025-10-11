// Package scache 是一个高性能的 Go 语言缓存库
//
// 特性：
//   - 支持内存缓存，基于 Go map 和 RWLock 实现
//   - 支持 TTL (Time To Live) 过期时间
//   - 支持 LRU (Least Recently Used) 淘汰策略
//   - 支持统计信息（命中率、操作次数等）
//   - 提供实例化和全局单例两种使用方式
//   - 线程安全，支持高并发访问
//   - 支持定期清理过期项目
//
// 基本使用：
//
//	// 1. 实例化使用
//	cache := scache.NewCache(
//		scache.WithMaxSize(1000),
//		scache.WithDefaultExpiration(time.Hour),
//	)
//
//	cache.Set("key", "value", time.Minute*10)
//	value, found := cache.Get("key")
//
//	// 2. 全局单例使用
//	scache.Set("global_key", "global_value", time.Hour)
//	value, found := scache.Get("global_key")
//
// 配置选项：
//
//   - WithMaxSize: 设置最大容量
//   - WithDefaultExpiration: 设置默认过期时间
//   - WithCleanupInterval: 设置清理间隔
//   - WithStats: 启用/禁用统计信息
//   - WithInitialCapacity: 设置初始容量
//
// API 方法：
//
//	Set(key string, value interface{}, ttl time.Duration) error
//	Get(key string) (interface{}, bool)
//	Delete(key string) bool
//	Exists(key string) bool
//	Flush()
//	Size() int
//	Stats() CacheStats
package scache

// 导出主要类型和函数
import (
	"scache/cache"
	"scache/interfaces"
	"scache/policies/lru"
	"scache/types"
)

// 主要类型别名
type (
	// Cache 缓存接口
	Cache = interfaces.Cache

	// CacheStats 缓存统计信息
	CacheStats = interfaces.CacheStats

	// EvictionPolicy 淘汰策略接口
	EvictionPolicy = interfaces.EvictionPolicy

	// CacheItem 缓存项
	CacheItem = types.CacheItem

	// CacheConfig 缓存配置
	CacheConfig = types.CacheConfig

	// CacheOption 缓存配置选项
	CacheOption = types.CacheOption
)

// 主要构造函数
var (
	// NewCache 创建新的缓存实例
	NewCache = cache.NewCache

	// GetGlobalCache 获取全局缓存实例
	GetGlobalCache = cache.GetGlobalCache

	// NewLRUPolicy 创建 LRU 策略
	NewLRUPolicy = lru.NewLRUPolicy
)

// 全局缓存函数
var (
	// Set 设置全局缓存项
	Set = cache.Set

	// Get 获取全局缓存项
	Get = cache.Get

	// Delete 删除全局缓存项
	Delete = cache.Delete

	// Exists 检查全局缓存项是否存在
	Exists = cache.Exists

	// Flush 清空全局缓存
	Flush = cache.Flush

	// Size 获取全局缓存大小
	Size = cache.Size

	// Stats 获取全局缓存统计
	Stats = cache.Stats

	// Keys 获取全局缓存所有键
	Keys = cache.Keys

	// GetWithExpiration 获取全局缓存项和过期时间
	GetWithExpiration = cache.GetWithExpiration

	// ConfigureGlobalCache 配置全局缓存
	ConfigureGlobalCache = cache.ConfigureGlobalCache

	// CloseGlobalCache 关闭全局缓存
	CloseGlobalCache = cache.CloseGlobalCache
)

// 配置选项函数
var (
	// WithDefaultExpiration 设置默认过期时间
	WithDefaultExpiration = types.WithDefaultExpiration

	// WithCleanupInterval 设置清理间隔
	WithCleanupInterval = types.WithCleanupInterval

	// WithMaxSize 设置最大容量
	WithMaxSize = types.WithMaxSize

	// WithStats 设置是否启用统计
	WithStats = types.WithStats

	// WithInitialCapacity 设置初始容量
	WithInitialCapacity = types.WithInitialCapacity
)

// 默认配置函数
var (
	// DefaultCacheConfig 返回默认配置
	DefaultCacheConfig = types.DefaultCacheConfig
)
