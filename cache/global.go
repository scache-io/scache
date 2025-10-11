package cache

import (
	"sync"
	"time"
)

var (
	// 全局便捷函数使用管理器的单例
	managerInstance = GetGlobalManager()
)

// Register 注册全局缓存
func Register(name string, c Cache) error {
	return managerInstance.Register(name, c)
}

// RegisterLRU 注册 LRU 缓存
func RegisterLRU(name string, maxSize int, opts ...Option) error {
	allOpts := append([]Option{WithMaxSize(maxSize), WithEvictionPolicy("lru")}, opts...)
	c := New(allOpts...)
	return managerInstance.Register(name, c)
}

// RegisterLFU 注册 LFU 缓存
func RegisterLFU(name string, maxSize int, opts ...Option) error {
	allOpts := append([]Option{WithMaxSize(maxSize), WithEvictionPolicy("lfu")}, opts...)
	c := New(allOpts...)
	return managerInstance.Register(name, c)
}

// RegisterFIFO 注册 FIFO 缓存
func RegisterFIFO(name string, maxSize int, opts ...Option) error {
	allOpts := append([]Option{WithMaxSize(maxSize), WithEvictionPolicy("fifo")}, opts...)
	c := New(allOpts...)
	return managerInstance.Register(name, c)
}

// Get 获取全局缓存
func Get(name string) (Cache, error) {
	return managerInstance.Get(name)
}

// GetOrDefault 获取缓存，如果不存在则创建默认缓存
func GetOrDefault(name string, opts ...Option) Cache {
	return managerInstance.GetOrDefault(name, opts...)
}

// Remove 移除全局缓存
func Remove(name string) error {
	return managerInstance.Remove(name)
}

// List 列出所有全局缓存名称
func List() []string {
	return managerInstance.List()
}

// Clear 清空所有全局缓存
func Clear() error {
	return managerInstance.Clear()
}

// Close 关闭所有全局缓存
func Close() error {
	return managerInstance.Close()
}

// Stats 获取所有全局缓存的统计信息
func Stats() map[string]CacheStats {
	return managerInstance.Stats()
}

// Size 获取所有全局缓存的总大小
func Size() int {
	return managerInstance.Size()
}

// Exists 检查全局缓存是否存在
func Exists(name string) bool {
	return managerInstance.Exists(name)
}

// 默认缓存的便捷操作
var (
	defaultCacheOnce sync.Once
	defaultCache     Cache
)

// getDefaultCache 获取默认缓存实例
func getDefaultCache() Cache {
	defaultCacheOnce.Do(func() {
		defaultCache = New()
		Register("default", defaultCache)
	})
	return defaultCache
}

// Set 在默认缓存中设置键值
func Set(key string, value interface{}) error {
	return getDefaultCache().Set(key, value)
}

// SetWithTTL 在默认缓存中设置带过期时间的键值
func SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	return getDefaultCache().SetWithTTL(key, value, ttl)
}

// GetFromDefault 从默认缓存中获取值
func GetFromDefault(key string) (interface{}, bool) {
	return getDefaultCache().Get(key)
}

// DeleteFromDefault 从默认缓存中删除键
func DeleteFromDefault(key string) bool {
	return getDefaultCache().Delete(key)
}

// ExistsInDefault 检查默认缓存中是否存在键
func ExistsInDefault(key string) bool {
	return getDefaultCache().Exists(key)
}

// ClearDefault 清空默认缓存
func ClearDefault() error {
	return getDefaultCache().Clear()
}
