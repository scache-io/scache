package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/internal"
	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/storage"
	"github.com/scache-io/scache/types"
)

// NewEngine 创建新的存储引擎实例
func NewEngine(engineConfig *config.EngineConfig) interfaces.StorageEngine {
	return storage.NewStorageEngine(engineConfig)
}

// LocalCache 局部缓存封装
type LocalCache struct {
	engine interfaces.StorageEngine
}

// NewLocalCache 创建局部缓存实例
func NewLocalCache(engineConfig *config.EngineConfig) *LocalCache {
	return &LocalCache{
		engine: NewEngine(engineConfig),
	}
}

// SetString 设置字符串值
func (c *LocalCache) SetString(key, value string, ttl ...time.Duration) error {
	obj := types.NewStringObject(value, internal.ParseTTL(ttl))
	return c.engine.Set(key, obj)
}

// GetString 获取字符串值
func (c *LocalCache) GetString(key string) (string, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return "", false
	}

	return internal.ExtractStringValue(obj)
}

// SetList 设置列表值
func (c *LocalCache) SetList(key string, values []interface{}, ttl ...time.Duration) error {
	obj := types.NewListObject(values, internal.ParseTTL(ttl))
	return c.engine.Set(key, obj)
}

// GetList 获取列表值
func (c *LocalCache) GetList(key string) ([]interface{}, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return nil, false
	}

	return internal.ExtractListValue(obj)
}

// SetHash 设置哈希值
func (c *LocalCache) SetHash(key string, fields map[string]interface{}, ttl ...time.Duration) error {
	obj := types.NewHashObject(fields, internal.ParseTTL(ttl))
	return c.engine.Set(key, obj)
}

// GetHash 获取哈希值
func (c *LocalCache) GetHash(key string) (map[string]interface{}, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return nil, false
	}

	return internal.ExtractHashValue(obj)
}

// Store 存储结构体值（JSON序列化，支持指针和非指针类型）
func (c *LocalCache) Store(key string, obj interface{}, ttl ...time.Duration) error {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	stringObj := types.NewStringObject(string(jsonBytes), internal.ParseTTL(ttl))
	return c.engine.Set(key, stringObj)
}

// Load 加载结构体值（JSON反序列化，要求指针参数）
func (c *LocalCache) Load(key string, dest interface{}) error {
	// 验证参数
	if err := internal.ValidatePointerArgument(dest); err != nil {
		return err
	}

	obj, exists := c.engine.Get(key)
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	jsonData, ok := internal.ExtractStructValue(obj)
	if !ok {
		return fmt.Errorf("type mismatch")
	}

	return json.Unmarshal([]byte(jsonData), dest)
}

// Delete 删除键
func (c *LocalCache) Delete(key string) bool {
	return c.engine.Delete(key)
}

// Exists 检查键是否存在
func (c *LocalCache) Exists(key string) bool {
	return c.engine.Exists(key)
}

// Keys 获取所有键
func (c *LocalCache) Keys() []string {
	return c.engine.Keys()
}

// Flush 清空所有数据
func (c *LocalCache) Flush() error {
	return c.engine.Flush()
}

// Size 获取缓存大小
func (c *LocalCache) Size() int {
	return c.engine.Size()
}

// Expire 设置过期时间
func (c *LocalCache) Expire(key string, ttl time.Duration) bool {
	return c.engine.Expire(key, ttl)
}

// TTL 获取剩余生存时间
func (c *LocalCache) TTL(key string) (time.Duration, bool) {
	return c.engine.TTL(key)
}

// Stats 获取统计信息
func (c *LocalCache) Stats() interface{} {
	return c.engine.Stats()
}

// GetEngine 获取底层引擎（用于高级操作）
func (c *LocalCache) GetEngine() interfaces.StorageEngine {
	return c.engine
}
