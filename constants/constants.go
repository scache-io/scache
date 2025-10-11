package constants

import "time"

// 缓存默认配置常量
const (
	// DefaultMaxSize 默认最大缓存项数量
	DefaultMaxSize = 10000
	// DefaultShards 默认分片数量
	DefaultShards = 16
	// DefaultCleanupInterval 默认清理间隔
	DefaultCleanupInterval = 10 * time.Minute
	// DefaultTTL 默认TTL（0表示永不过期）
	DefaultTTL = 0
)

// 缓存策略常量
const (
	// LRUStrategy LRU淘汰策略
	LRUStrategy = "lru"
	// LFUStrategy LFU淘汰策略
	LFUStrategy = "lfu"
	// FIFOStrategy FIFO淘汰策略
	FIFOStrategy = "fifo"
)

// 全局缓存常量
const (
	// DefaultCacheName 默认缓存名称
	DefaultCacheName = "default"
	// ManagerTimeout 管理器操作超时时间
	ManagerTimeout = 30 * time.Second
)

// 性能相关常量
const (
	// MaxKeyLength 最大键长度
	MaxKeyLength = 256
	// MinKeyLength 最小键长度
	MinKeyLength = 1
	// MaxValueSize 最大值大小（字节）
	MaxValueSize = 10 * 1024 * 1024 // 10MB
)

// 错误消息常量
const (
	// ErrCacheNotFound 缓存未找到错误
	ErrCacheNotFound = "cache not found"
	// ErrInvalidCacheName 无效缓存名称错误
	ErrInvalidCacheName = "invalid cache name"
	// ErrCacheAlreadyExists 缓存已存在错误
	ErrCacheAlreadyExists = "cache already exists"
	// ErrInvalidStrategy 无效策略错误
	ErrInvalidStrategy = "invalid eviction strategy"
	// ErrKeyNotFound 键未找到错误
	ErrKeyNotFound = "key not found"
	// ErrKeyTooLong 键过长错误
	ErrKeyTooLong = "key too long"
	// ErrKeyEmpty 键为空错误
	ErrKeyEmpty = "key empty"
	// ErrValueTooLarge 值过大错误
	ErrValueTooLarge = "value too large"
	// ErrCacheClosed 缓存已关闭错误
	ErrCacheClosed = "cache is closed"
)

// 日志相关常量
const (
	// LogPrefixCache 缓存日志前缀
	LogPrefixCache = "[SCache]"
	// LogPrefixManager 管理器日志前缀
	LogPrefixManager = "[SCache-Manager]"
	// LogPrefixGlobal 全局缓存日志前缀
	LogPrefixGlobal = "[SCache-Global]"
)

// 统计相关常量
const (
	// StatsUpdateInterval 统计信息更新间隔
	StatsUpdateInterval = time.Second
	// HitRateThreshold 命中率阈值
	HitRateThreshold = 0.8
)

// 序列化常量
const (
	// JSONEncoding JSON编码
	JSONEncoding = "json"
	// GobEncoding Gob编码
	GobEncoding = "gob"
)
