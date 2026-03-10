package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/storage"
	"github.com/scache-io/scache/types"
	"github.com/scache-io/scache/utils"
)

// NewEngine Create new storage engine instance
func NewEngine(engineConfig *config.EngineConfig) interfaces.StorageEngine {
	return storage.NewStorageEngine(engineConfig)
}

// LocalCache Local cache wrapper
type LocalCache struct {
	engine interfaces.StorageEngine
}

// NewLocalCache Create local cache instance
func NewLocalCache(engineConfig *config.EngineConfig) *LocalCache {
	return &LocalCache{
		engine: NewEngine(engineConfig),
	}
}

// SetString Set string value
func (c *LocalCache) SetString(key, value string, ttl ...time.Duration) error {
	obj := types.NewStringObject(value, utils.ParseTTL(ttl))
	return c.engine.Set(key, obj)
}

// GetString Get string value
func (c *LocalCache) GetString(key string) (string, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return "", false
	}

	return utils.ExtractStringValue(obj)
}

// SetList Set list value
func (c *LocalCache) SetList(key string, values []interface{}, ttl ...time.Duration) error {
	obj := types.NewListObject(values, utils.ParseTTL(ttl))
	return c.engine.Set(key, obj)
}

// GetList Get list value
func (c *LocalCache) GetList(key string) ([]interface{}, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return nil, false
	}

	return utils.ExtractListValue(obj)
}

// SetHash Set hash value
func (c *LocalCache) SetHash(key string, fields map[string]interface{}, ttl ...time.Duration) error {
	obj := types.NewHashObject(fields, utils.ParseTTL(ttl))
	return c.engine.Set(key, obj)
}

// GetHash Get hash value
func (c *LocalCache) GetHash(key string) (map[string]interface{}, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return nil, false
	}

	return utils.ExtractHashValue(obj)
}

// Store Store struct值（JSON序列化，支持指针和非指针Type）
func (c *LocalCache) Store(key string, obj interface{}, ttl ...time.Duration) error {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	stringObj := types.NewStringObject(string(jsonBytes), utils.ParseTTL(ttl))
	return c.engine.Set(key, stringObj)
}

// Load Load struct值（JSON反序列化，要求指针Parameter）
func (c *LocalCache) Load(key string, dest interface{}) error {
	// 验证Parameter
	if err := utils.ValidatePointerArgument(dest); err != nil {
		return err
	}

	obj, exists := c.engine.Get(key)
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	jsonData, ok := utils.ExtractStructValue(obj)
	if !ok {
		return fmt.Errorf("type mismatch")
	}

	return json.Unmarshal([]byte(jsonData), dest)
}

// Delete Delete key
func (c *LocalCache) Delete(key string) bool {
	return c.engine.Delete(key)
}

// Exists Check if key exists
func (c *LocalCache) Exists(key string) bool {
	return c.engine.Exists(key)
}

// Keys Get all keys
func (c *LocalCache) Keys() []string {
	return c.engine.Keys()
}

// Flush 清空所有数据
func (c *LocalCache) Flush() error {
	return c.engine.Flush()
}

// Size Get cache size
func (c *LocalCache) Size() int {
	return c.engine.Size()
}

// Expire Set expiration time
func (c *LocalCache) Expire(key string, ttl time.Duration) bool {
	return c.engine.Expire(key, ttl)
}

// TTL 获取剩余生存时间
func (c *LocalCache) TTL(key string) (time.Duration, bool) {
	return c.engine.TTL(key)
}

// Stats Get statistics
func (c *LocalCache) Stats() interface{} {
	return c.engine.Stats()
}

// GetEngine 获取底层引擎（用于高级操作）
func (c *LocalCache) GetEngine() interfaces.StorageEngine {
	return c.engine
}
