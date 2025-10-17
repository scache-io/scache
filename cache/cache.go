package cache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/storage"
	"github.com/scache-io/scache/types"
)

// NewEngine 创建新的存储引擎实例
func NewEngine(opts ...config.EngineOption) interfaces.StorageEngine {
	// 创建配置
	engineConfig := &storage.EngineConfig{}
	for _, opt := range opts {
		opt(engineConfig)
	}
	return storage.NewStorageEngine(engineConfig)
}

// LocalCache 局部缓存封装
type LocalCache struct {
	engine interfaces.StorageEngine
}

// NewLocalCache 创建局部缓存实例
func NewLocalCache(opts ...config.EngineOption) *LocalCache {
	return &LocalCache{
		engine: NewEngine(opts...),
	}
}

// SetString 设置字符串值
func (c *LocalCache) SetString(key, value string, ttl ...time.Duration) error {
	var expiration time.Duration
	if len(ttl) > 0 {
		expiration = ttl[0]
	}
	obj := types.NewStringObject(value, expiration)
	return c.engine.Set(key, obj)
}

// GetString 获取字符串值
func (c *LocalCache) GetString(key string) (string, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return "", false
	}

	if strObj, ok := obj.(*types.StringObject); ok {
		return strObj.Value(), true
	}
	return "", false
}

// SetList 设置列表值
func (c *LocalCache) SetList(key string, values []interface{}, ttl ...time.Duration) error {
	var expiration time.Duration
	if len(ttl) > 0 {
		expiration = ttl[0]
	}
	obj := types.NewListObject(values, expiration)
	return c.engine.Set(key, obj)
}

// GetList 获取列表值
func (c *LocalCache) GetList(key string) ([]interface{}, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return nil, false
	}

	if listObj, ok := obj.(*types.ListObject); ok {
		return listObj.Values(), true
	}
	return nil, false
}

// SetHash 设置哈希值
func (c *LocalCache) SetHash(key string, fields map[string]interface{}, ttl ...time.Duration) error {
	var expiration time.Duration
	if len(ttl) > 0 {
		expiration = ttl[0]
	}
	obj := types.NewHashObject(fields, expiration)
	return c.engine.Set(key, obj)
}

// GetHash 获取哈希值
func (c *LocalCache) GetHash(key string) (map[string]interface{}, bool) {
	obj, exists := c.engine.Get(key)
	if !exists {
		return nil, false
	}

	if hashObj, ok := obj.(*types.HashObject); ok {
		return hashObj.Fields(), true
	}
	return nil, false
}

// Store 存储结构体值（JSON序列化，支持指针和非指针类型）
func (c *LocalCache) Store(key string, obj interface{}, ttl ...time.Duration) error {
	var expiration time.Duration
	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	stringObj := types.NewStringObject(string(jsonBytes), expiration)
	return c.engine.Set(key, stringObj)
}

// Load 加载结构体值（JSON反序列化，要求指针参数）
func (c *LocalCache) Load(key string, dest interface{}) error {
	// 检查参数是否为指针类型
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("Load requires a pointer argument, got %T", dest)
	}

	obj, exists := c.engine.Get(key)
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	stringObj, ok := obj.(*types.StringObject)
	if !ok {
		return fmt.Errorf("value is not a struct object")
	}

	return json.Unmarshal([]byte(stringObj.Value()), dest)
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
