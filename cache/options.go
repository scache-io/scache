package cache

import (
	"time"

	"scache/types"
)

// 导出配置选项函数

// WithMaxSize 设置最大容量
func WithMaxSize(size int) types.CacheOption {
	return types.WithMaxSize(size)
}

// WithDefaultExpiration 设置默认过期时间
func WithDefaultExpiration(d time.Duration) types.CacheOption {
	return types.WithDefaultExpiration(d)
}

// WithCleanupInterval 设置清理间隔
func WithCleanupInterval(d time.Duration) types.CacheOption {
	return types.WithCleanupInterval(d)
}

// WithStats 设置是否启用统计
func WithStats(enable bool) types.CacheOption {
	return types.WithStats(enable)
}

// WithInitialCapacity 设置初始容量
func WithInitialCapacity(capacity int) types.CacheOption {
	return types.WithInitialCapacity(capacity)
}
