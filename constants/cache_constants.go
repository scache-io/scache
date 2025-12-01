package constants

import "time"

// 默认配置常量
const (
	DefaultMaxSize         = 0    // 默认最大缓存大小，0表示无限制
	DefaultExpiration      = 0    // 默认过期时间，0表示永不过期
	DefaultCleanupInterval = 0    // 默认清理间隔，0表示不执行清理
	DefaultInitialCapacity = 16   // 默认初始容量
	DefaultStatsEnabled    = true // 默认启用统计功能
)

// DefaultLRUCapacity LRU策略默认配置
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

// 容量常量
const (
	SmallCapacity  = 1000   // 小型缓存容量
	MediumCapacity = 10000  // 中型缓存容量
	LargeCapacity  = 100000 // 大型缓存容量
	TestCapacity   = 10     // 测试用容量
)

// 内存阈值常量
const (
	SmallMemoryThreshold   = 0.7    // 小型配置内存阈值 70%
	MediumMemoryThreshold  = 0.8    // 中型配置内存阈值 80%
	LargeMemoryThreshold   = 0.85   // 大型配置内存阈值 85%
	DefaultMemoryThreshold = 0.8    // 默认内存阈值 80%
	MinMemoryThreshold     = 0.0    // 最小内存阈值
	MaxMemoryThreshold     = 1.0    // 最大内存阈值
	TestMemoryThreshold    = 0.0001 // 测试用极低内存阈值
)

// 常用时间间隔常量
const (
	TwoMinutes    = 2 * Minute  // 2分钟
	TenMinutes    = 10 * Minute // 10分钟
	ThirtyMinutes = 30 * Minute // 30分钟
	TwoHours      = 2 * Hour    // 2小时
	SixHours      = 6 * Hour    // 6小时
	ThirtySeconds = 30 * Second // 30秒
	OneHour       = 1 * Hour    // 1小时
)
