package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/policies/lru"
	"github.com/scache-io/scache/types"
)

// StorageEngine 存储引擎实现
type StorageEngine struct {
	mu        sync.RWMutex
	data      map[string]interfaces.DataObject
	policy    interfaces.EvictionPolicy
	config    *EngineConfig
	stats     *EngineStats
	stopChan  chan struct{}
	bgCleanup chan struct{}
}

// EngineConfig 引擎配置
type EngineConfig struct {
	MaxSize                   int           // 最大缓存数量
	MemoryThreshold           float64       // 内存阈值
	DefaultExpiration         time.Duration // 默认过期时间
	BackgroundCleanupInterval time.Duration // 后台清理间隔
}

// DefaultEngineConfig 默认引擎配置
func DefaultEngineConfig() *EngineConfig {
	return &EngineConfig{
		MaxSize:                   0,               // 无限制
		MemoryThreshold:           0.8,             // 80%
		DefaultExpiration:         0,               // 永不过期
		BackgroundCleanupInterval: 5 * time.Minute, // 5分钟
	}
}

// EngineStats 引擎统计
type EngineStats struct {
	mu          sync.RWMutex
	hits        int64
	misses      int64
	sets        int64
	deletes     int64
	evictions   int64
	expirations int64
	memoryUsage int64 // 字节
}

// NewStorageEngine 创建新的存储引擎
func NewStorageEngine(config *EngineConfig) interfaces.StorageEngine {
	if config == nil {
		config = DefaultEngineConfig()
	}

	engine := &StorageEngine{
		data:      make(map[string]interfaces.DataObject),
		policy:    lru.NewLRUPolicy(config.MaxSize),
		config:    config,
		stats:     &EngineStats{},
		stopChan:  make(chan struct{}),
		bgCleanup: make(chan struct{}),
	}

	// 启动后台清理
	if config.BackgroundCleanupInterval > 0 {
		engine.startBackgroundCleanup()
	}

	return engine
}

// Set 存储对象
func (e *StorageEngine) Set(key string, obj interfaces.DataObject) error {
	if key == "" {
		return errors.New("cache key cannot be empty")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// 检查是否需要淘汰（在设置新键之前）
	if e.config.MaxSize > 0 && len(e.data) >= e.config.MaxSize && e.data[key] == nil {
		e.evictOne()
	}

	e.data[key] = obj
	e.policy.Set(key)
	e.stats.recordSet()

	return nil
}

// Get 获取对象
func (e *StorageEngine) Get(key string) (interfaces.DataObject, bool) {
	if key == "" {
		return nil, false
	}

	e.mu.RLock()
	obj, exists := e.data[key]
	e.mu.RUnlock()

	if !exists {
		e.stats.recordMiss()
		return nil, false
	}

	// 检查过期
	if obj.IsExpired() {
		go e.Delete(key)
		e.stats.recordMiss()
		e.stats.recordExpiration()
		return nil, false
	}

	e.policy.Access(key)
	e.stats.recordHit()
	return obj, true
}

// Delete 删除对象
func (e *StorageEngine) Delete(key string) bool {
	if key == "" {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.data[key]; exists {
		delete(e.data, key)
		e.policy.Delete(key)
		e.stats.recordDelete()
		return true
	}

	return false
}

// Exists 检查键是否存在
func (e *StorageEngine) Exists(key string) bool {
	if key == "" {
		return false
	}

	e.mu.RLock()
	obj, exists := e.data[key]
	e.mu.RUnlock()

	if !exists {
		return false
	}

	if obj.IsExpired() {
		go e.Delete(key)
		return false
	}

	return true
}

// Keys 获取所有键
func (e *StorageEngine) Keys() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	keys := make([]string, 0, len(e.data))
	for key := range e.data {
		keys = append(keys, key)
	}
	return keys
}

// Flush 清空所有数据
func (e *StorageEngine) Flush() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.data = make(map[string]interfaces.DataObject)
	e.policy.Clear()
	e.stats.reset()
	return nil
}

// Size 返回当前大小
func (e *StorageEngine) Size() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.data)
}

// Type 获取键的类型
func (e *StorageEngine) Type(key string) (interfaces.DataType, bool) {
	e.mu.RLock()
	obj, exists := e.data[key]
	e.mu.RUnlock()

	if !exists {
		return "", false
	}

	if obj.IsExpired() {
		go e.Delete(key)
		return "", false
	}

	return obj.Type(), true
}

// Expire 设置过期时间
func (e *StorageEngine) Expire(key string, ttl time.Duration) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	obj, exists := e.data[key]
	if !exists {
		return false
	}

	// 创建新的对象以更新过期时间
	switch t := obj.(type) {
	case *types.StringObject:
		newObj := types.NewStringObject(t.Value(), ttl)
		e.data[key] = newObj
		return true
	case *types.ListObject:
		newObj := types.NewListObject(t.Values(), ttl)
		e.data[key] = newObj
		return true
	case *types.HashObject:
		newObj := types.NewHashObject(t.Fields(), ttl)
		e.data[key] = newObj
		return true
	}

	return false
}

// TTL 获取剩余生存时间
func (e *StorageEngine) TTL(key string) (time.Duration, bool) {
	e.mu.RLock()
	obj, exists := e.data[key]
	e.mu.RUnlock()

	if !exists {
		return -1, false
	}

	if obj.ExpiresAt().IsZero() {
		return -1, true // 永不过期
	}

	remaining := time.Until(obj.ExpiresAt())
	if remaining <= 0 {
		return 0, true // 已过期
	}

	return remaining, true
}

// Stats 获取统计信息
func (e *StorageEngine) Stats() interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return map[string]interface{}{
		"hits":        e.stats.hits,
		"misses":      e.stats.misses,
		"sets":        e.stats.sets,
		"deletes":     e.stats.deletes,
		"evictions":   e.stats.evictions,
		"expirations": e.stats.expirations,
		"memory":      e.stats.memoryUsage,
		"keys":        len(e.data),
		"hit_rate":    e.stats.hitRate(),
	}
}

// evictOne 淘汰一个键
func (e *StorageEngine) evictOne() {
	if key := e.policy.Evict(); key != "" {
		delete(e.data, key)
		e.stats.recordEviction()
	}
}

// startBackgroundCleanup 启动后台清理
func (e *StorageEngine) startBackgroundCleanup() {
	go func() {
		ticker := time.NewTicker(e.config.BackgroundCleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.cleanupExpired()
			case <-e.stopChan:
				return
			case <-e.bgCleanup:
				return
			}
		}
	}()
}

// cleanupExpired 清理过期项目
func (e *StorageEngine) cleanupExpired() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for key, obj := range e.data {
		if obj.IsExpired() {
			delete(e.data, key)
			e.policy.Delete(key)
			e.stats.recordExpiration()
		}
	}
}

// GetConfig 获取引擎配置
func (e *StorageEngine) GetConfig() *EngineConfig {
	return e.config
}

// Close 关闭引擎
func (e *StorageEngine) Close() {
	close(e.stopChan)
}

// EngineStats 方法实现

func (s *EngineStats) recordHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hits++
}

func (s *EngineStats) recordMiss() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.misses++
}

func (s *EngineStats) recordSet() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sets++
}

func (s *EngineStats) recordDelete() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deletes++
}

func (s *EngineStats) recordEviction() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evictions++
}

func (s *EngineStats) recordExpiration() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expirations++
}

func (s *EngineStats) hitRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := s.hits + s.misses
	if total == 0 {
		return 0
	}
	return float64(s.hits) / float64(total)
}

func (s *EngineStats) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hits = 0
	s.misses = 0
	s.sets = 0
	s.deletes = 0
	s.evictions = 0
	s.expirations = 0
}
