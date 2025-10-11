package interfaces

import "time"

// Cache 定义缓存接口
type Cache interface {
	// Set 设置缓存项
	Set(key string, value interface{}, ttl time.Duration) error

	// Get 获取缓存项
	Get(key string) (interface{}, bool)

	// Delete 删除缓存项
	Delete(key string) bool

	// Exists 检查缓存项是否存在
	Exists(key string) bool

	// Flush 清空所有缓存项
	Flush()

	// Size 获取缓存项数量
	Size() int

	// Stats 获取缓存统计信息
	Stats() CacheStats
}

// CacheStats 缓存统计信息
type CacheStats struct {
	// 命中次数
	Hits int64

	// 未命中次数
	Misses int64

	// 设置次数
	Sets int64

	// 删除次数
	Deletes int64

	// 当前缓存项数量
	Size int

	// 最大容量
	MaxSize int

	// 命中率
	HitRate float64
}

// EvictionPolicy 淘汰策略接口
type EvictionPolicy interface {
	// Access 当访问 key 时调用
	Access(key string)

	// Set 当设置新 key 时调用
	Set(key string)

	// Delete 当删除 key 时调用
	Delete(key string)

	// Evict 获取需要淘汰的 key
	Evict() string

	// Size 获取当前策略状态
	Size() int

	// Clear 清空所有数据
	Clear()

	// Contains 检查 key 是否存在
	Contains(key string) bool

	// Keys 获取所有 key（按最近使用时间排序）
	Keys() []string

	// UpdateCapacity 更新容量限制
	UpdateCapacity(newCapacity int)
}
