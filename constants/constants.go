package constants

import "time"

// 默认配置常量
const (
	DefaultMaxSize         = 0                // 默认最大缓存大小，0表示无限制
	DefaultExpiration      = 0                // 默认过期时间，0表示永不过期
	DefaultCleanupInterval = 10 * time.Minute // 默认清理间隔：10分钟
	DefaultInitialCapacity = 16               // 默认初始容量
	DefaultStatsEnabled    = true             // 默认启用统计功能
)

// LRU策略默认配置
const (
	DefaultLRUCapacity = 100 // LRU策略的默认容量
)

// 时间常量，方便用户使用
const (
	Second = time.Second    // 秒
	Minute = time.Minute    // 分钟
	Hour   = time.Hour      // 小时
	Day    = 24 * time.Hour // 天
	Week   = 7 * Day        // 周
)

// 缓存条目状态常量
const (
	StatusActive  = "active"  // 活跃状态
	StatusExpired = "expired" // 已过期
	StatusEvicted = "evicted" // 已被淘汰
)

// 错误消息常量
const (
	ErrKeyEmpty        = "cache key cannot be empty"       // 缓存键不能为空
	ErrValueNil        = "cache value cannot be nil"       // 缓存值不能为nil
	ErrCapacityInvalid = "cache capacity must be positive" // 缓存容量必须为正数
	ErrTTLInvalid      = "TTL must be non-negative"        // TTL必须为非负数
)
