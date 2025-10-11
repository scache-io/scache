package cache

import (
	"sync"
	"time"

	"scache/interfaces"
	"scache/types"
)

var (
	// globalCache 全局缓存实例
	globalCache interfaces.Cache

	// globalOnce 确保全局缓存只初始化一次
	globalOnce sync.Once

	// globalMu 保护全局缓存的访问
	globalMu sync.RWMutex
)

// initGlobalCache 初始化全局缓存
func initGlobalCache() {
	globalCache = NewCache(
		types.WithDefaultExpiration(time.Hour),
		types.WithCleanupInterval(time.Minute*10),
		types.WithMaxSize(10000),
		types.WithStats(true),
		types.WithInitialCapacity(128),
	)
}

// GetGlobalCache 获取全局缓存实例
func GetGlobalCache() interfaces.Cache {
	globalOnce.Do(initGlobalCache)
	return globalCache
}

// Set 设置全局缓存项
func Set(key string, value interface{}, ttl ...time.Duration) error {
	var expiration time.Duration
	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	globalMu.Lock()
	defer globalMu.Unlock()

	return GetGlobalCache().Set(key, value, expiration)
}

// Get 获取全局缓存项
func Get(key string) (interface{}, bool) {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return GetGlobalCache().Get(key)
}

// Delete 删除全局缓存项
func Delete(key string) bool {
	globalMu.Lock()
	defer globalMu.Unlock()

	return GetGlobalCache().Delete(key)
}

// Exists 检查全局缓存项是否存在
func Exists(key string) bool {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return GetGlobalCache().Exists(key)
}

// Flush 清空全局缓存
func Flush() {
	globalMu.Lock()
	defer globalMu.Unlock()

	GetGlobalCache().Flush()
}

// Size 获取全局缓存大小
func Size() int {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return GetGlobalCache().Size()
}

// Stats 获取全局缓存统计
func Stats() interfaces.CacheStats {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return GetGlobalCache().Stats()
}

// Keys 获取全局缓存所有键
func Keys() []string {
	globalMu.RLock()
	defer globalMu.RUnlock()

	if cache, ok := GetGlobalCache().(*MemoryCache); ok {
		return cache.Keys()
	}
	return nil
}

// GetWithExpiration 获取全局缓存项和过期时间
func GetWithExpiration(key string) (interface{}, time.Time, bool) {
	globalMu.RLock()
	defer globalMu.RUnlock()

	if cache, ok := GetGlobalCache().(*MemoryCache); ok {
		return cache.GetWithExpiration(key)
	}
	return nil, time.Time{}, false
}

// ConfigureGlobalCache 配置全局缓存（在首次使用前调用）
func ConfigureGlobalCache(opts ...types.CacheOption) {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalCache == nil {
		globalCache = NewCache(opts...)
	}
}

// CloseGlobalCache 关闭全局缓存
func CloseGlobalCache() {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalCache != nil {
		if cache, ok := globalCache.(*MemoryCache); ok {
			cache.Close()
		}
	}
}
