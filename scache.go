package scache

import (
	"sync"
	"time"

	"github.com/scache-io/scache/cache"
	"github.com/scache-io/scache/config"
)

// LocalCache 局部缓存封装的别名，方便外部使用
type LocalCache = cache.LocalCache

// New 创建新的局部缓存实例
func New(opts ...config.EngineOption) *LocalCache {
	return cache.NewLocalCache(opts...)
}

// 全局默认实例
var (
	globalCache *LocalCache
	globalOnce  sync.Once
)

// GetGlobalCache 获取全局缓存实例（线程安全）
func GetGlobalCache() *LocalCache {
	globalOnce.Do(func() {
		// 初始化默认实例，使用中等配置
		globalCache = New(config.MediumConfig...)
	})
	return globalCache
}

// InitGlobalCache 初始化全局缓存（可配置）
func InitGlobalCache(opts ...config.EngineOption) {
	globalOnce.Do(func() {
		globalCache = New(opts...)
	})
}

// SetString 全局设置字符串值
func SetString(key, value string, ttl ...time.Duration) error {
	return GetGlobalCache().SetString(key, value, ttl...)
}

// GetString 全局获取字符串值
func GetString(key string) (string, bool) {
	return GetGlobalCache().GetString(key)
}

// SetList 全局设置列表值
func SetList(key string, values []interface{}, ttl ...time.Duration) error {
	return GetGlobalCache().SetList(key, values, ttl...)
}

// GetList 全局获取列表值
func GetList(key string) ([]interface{}, bool) {
	return GetGlobalCache().GetList(key)
}

// SetHash 全局设置哈希值
func SetHash(key string, fields map[string]interface{}, ttl ...time.Duration) error {
	return GetGlobalCache().SetHash(key, fields, ttl...)
}

// GetHash 全局获取哈希值
func GetHash(key string) (map[string]interface{}, bool) {
	return GetGlobalCache().GetHash(key)
}

// Delete 全局删除键
func Delete(key string) bool {
	return GetGlobalCache().Delete(key)
}

// Exists 全局检查键是否存在
func Exists(key string) bool {
	return GetGlobalCache().Exists(key)
}

// Keys 全局获取所有键
func Keys() []string {
	return GetGlobalCache().Keys()
}

// Flush 全局清空所有数据
func Flush() error {
	return GetGlobalCache().Flush()
}

// Size 全局获取缓存大小
func Size() int {
	return GetGlobalCache().Size()
}

// Expire 全局设置过期时间
func Expire(key string, ttl time.Duration) bool {
	return GetGlobalCache().Expire(key, ttl)
}

// TTL 全局获取剩余生存时间
func TTL(key string) (time.Duration, bool) {
	return GetGlobalCache().TTL(key)
}

// Stats 全局获取统计信息
func Stats() interface{} {
	return GetGlobalCache().Stats()
}

// SetStruct 全局设置结构体值（JSON序列化）
func SetStruct(key string, obj interface{}, ttl ...time.Duration) error {
	return GetGlobalCache().SetStruct(key, obj, ttl...)
}

// GetStruct 全局获取结构体值（JSON反序列化）
func GetStruct(key string, dest interface{}) error {
	return GetGlobalCache().GetStruct(key, dest)
}