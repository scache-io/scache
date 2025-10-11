package constants

import "time"

// 默认配置常量
const (
	// DefaultMaxSize 默认最大容量
	DefaultMaxSize = 0 // 0 表示无限制

	// DefaultExpiration 默认过期时间
	DefaultExpiration = 0 // 0 表示永不过期

	// DefaultCleanupInterval 默认清理间隔
	DefaultCleanupInterval = 10 * time.Minute

	// DefaultInitialCapacity 默认初始容量
	DefaultInitialCapacity = 16

	// DefaultStatsEnabled 默认是否启用统计
	DefaultStatsEnabled = true
)

// LRU策略常量
const (
	// DefaultLRUCapacity 默认LRU容量
	DefaultLRUCapacity = 100
)

// 时间相关常量
const (
	// Second 秒
	Second = time.Second

	// Minute 分钟
	Minute = time.Minute

	// Hour 小时
	Hour = time.Hour

	// Day 天
	Day = 24 * time.Hour

	// Week 周
	Week = 7 * Day
)

// 缓存状态常量
const (
	// StatusActive 缓存活跃状态
	StatusActive = "active"

	// StatusExpired 缓存过期状态
	StatusExpired = "expired"

	// StatusEvicted 缓存被淘汰状态
	StatusEvicted = "evicted"
)

// 错误消息常量
const (
	// ErrKeyEmpty 键为空的错误
	ErrKeyEmpty = "cache key cannot be empty"

	// ErrValueNil 值为nil的错误
	ErrValueNil = "cache value cannot be nil"

	// ErrCapacityInvalid 容量无效的错误
	ErrCapacityInvalid = "cache capacity must be positive"

	// ErrTTLInvalid TTL无效的错误
	ErrTTLInvalid = "TTL must be non-negative"
)
