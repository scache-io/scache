package cache

import (
	"github.com/scache/interfaces"
)

// 重新导出接口类型
type (
	// Cache 缓存接口
	Cache = interfaces.Cache
	// CacheStats 缓存统计信息
	CacheStats = interfaces.CacheStats
	// CacheItem 缓存项
	CacheItem = interfaces.CacheItem
	// EvictionPolicy 淘汰策略接口
	EvictionPolicy = interfaces.EvictionPolicy
	// Serializer 序列化接口
	Serializer = interfaces.Serializer
)
