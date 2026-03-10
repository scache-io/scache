package api

import (
	"sync"
	"time"

	"github.com/scache-io/scache/cache"
	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/constants"
)

// LocalCache Local cache wrapper的别名，方便外部使用
type LocalCache = cache.LocalCache

// New 创建新的Local cache instance
func New(engineConfig *config.EngineConfig) *LocalCache {
	return cache.NewLocalCache(engineConfig)
}

// 全局默认实例
var (
	globalCache *LocalCache
	globalOnce  sync.Once
)

// GetGlobalCache 获取全局缓存实例（线程安全）
func GetGlobalCache() *LocalCache {
	globalOnce.Do(func() {
		// 创建中等配置
		mediumConfig := &config.EngineConfig{
			MaxSize:                   constants.MediumCapacity,
			MemoryThreshold:           constants.MediumMemoryThreshold,
			DefaultExpiration:         constants.TwoHours,
			BackgroundCleanupInterval: constants.TenMinutes,
		}
		globalCache = New(mediumConfig)
	})
	return globalCache
}

// InitGlobalCache Initialize全局缓存（可配置）
func InitGlobalCache(engineConfig *config.EngineConfig) {
	globalOnce.Do(func() {
		globalCache = New(engineConfig)
	})
}

// SetString 全局Set string value
func SetString(key, value string, ttl ...time.Duration) error {
	return GetGlobalCache().SetString(key, value, ttl...)
}

// GetString 全局Get string value
func GetString(key string) (string, bool) {
	return GetGlobalCache().GetString(key)
}

// SetList 全局Set list value
func SetList(key string, values []interface{}, ttl ...time.Duration) error {
	return GetGlobalCache().SetList(key, values, ttl...)
}

// GetList 全局Get list value
func GetList(key string) ([]interface{}, bool) {
	return GetGlobalCache().GetList(key)
}

// SetHash 全局Set hash value
func SetHash(key string, fields map[string]interface{}, ttl ...time.Duration) error {
	return GetGlobalCache().SetHash(key, fields, ttl...)
}

// GetHash 全局Get hash value
func GetHash(key string) (map[string]interface{}, bool) {
	return GetGlobalCache().GetHash(key)
}

// Store 全局Store struct值（JSON序列化，支持指针和非指针Type）
func Store(key string, obj interface{}, ttl ...time.Duration) error {
	return GetGlobalCache().Store(key, obj, ttl...)
}

// Load 全局Load struct值（JSON反序列化，要求指针Parameter）
func Load(key string, dest interface{}) error {
	return GetGlobalCache().Load(key, dest)
}

// Delete 全局Delete key
func Delete(key string) bool {
	return GetGlobalCache().Delete(key)
}

// Exists 全局Check if key exists
func Exists(key string) bool {
	return GetGlobalCache().Exists(key)
}

// Keys 全局Get all keys
func Keys() []string {
	return GetGlobalCache().Keys()
}

// Flush 全局清空所有数据
func Flush() error {
	return GetGlobalCache().Flush()
}

// Size 全局Get cache size
func Size() int {
	return GetGlobalCache().Size()
}

// Expire 全局Set expiration time
func Expire(key string, ttl time.Duration) bool {
	return GetGlobalCache().Expire(key, ttl)
}

// TTL 全局获取剩余生存时间
func TTL(key string) (time.Duration, bool) {
	return GetGlobalCache().TTL(key)
}

// Stats 全局Get statistics
func Stats() interface{} {
	return GetGlobalCache().Stats()
}
