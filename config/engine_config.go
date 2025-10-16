package config

import (
	"time"

	"scache/storage"
)

// EngineOption 引擎配置选项
type EngineOption func(*storage.EngineConfig)

// WithMaxSize 设置最大缓存数量
func WithMaxSize(size int) EngineOption {
	return func(config *storage.EngineConfig) {
		config.MaxSize = size
	}
}

// WithMemoryThreshold 设置内存阈值
func WithMemoryThreshold(threshold float64) EngineOption {
	return func(config *storage.EngineConfig) {
		if threshold >= 0 && threshold <= 1 {
			config.MemoryThreshold = threshold
		}
	}
}

// WithDefaultExpiration 设置默认过期时间
func WithDefaultExpiration(ttl time.Duration) EngineOption {
	return func(config *storage.EngineConfig) {
		config.DefaultExpiration = ttl
	}
}

// WithBackgroundCleanup 设置后台清理间隔
func WithBackgroundCleanup(interval time.Duration) EngineOption {
	return func(config *storage.EngineConfig) {
		config.BackgroundCleanupInterval = interval
	}
}

// 预定义配置

// DefaultConfig 默认配置
var DefaultConfig = []EngineOption{
	WithMaxSize(0),                         // 无限制
	WithMemoryThreshold(0.8),               // 80%
	WithDefaultExpiration(0),               // 永不过期
	WithBackgroundCleanup(5 * time.Minute), // 5分钟清理
}

// SmallConfig 小型配置（适用于内存较小的环境）
var SmallConfig = []EngineOption{
	WithMaxSize(1000),                      // 1000个键
	WithMemoryThreshold(0.7),               // 70%
	WithDefaultExpiration(time.Hour),       // 1小时过期
	WithBackgroundCleanup(2 * time.Minute), // 2分钟清理
}

// MediumConfig 中等配置（适用于一般应用）
var MediumConfig = []EngineOption{
	WithMaxSize(10000),                     // 10000个键
	WithMemoryThreshold(0.8),               // 80%
	WithDefaultExpiration(2 * time.Hour),   // 2小时过期
	WithBackgroundCleanup(5 * time.Minute), // 5分钟清理
}

// LargeConfig 大型配置（适用于高负载应用）
var LargeConfig = []EngineOption{
	WithMaxSize(100000),                     // 100000个键
	WithMemoryThreshold(0.85),               // 85%
	WithDefaultExpiration(6 * time.Hour),    // 6小时过期
	WithBackgroundCleanup(10 * time.Minute), // 10分钟清理
}
