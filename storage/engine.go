package storage

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/internal"
	"github.com/scache-io/scache/policies/lru"
	"github.com/scache-io/scache/types"
	"github.com/scache-io/scache/utils"
)

// StorageEngine Storage engine实现
type StorageEngine struct {
	mu        sync.RWMutex
	data      map[string]interfaces.DataObject
	policy    interfaces.EvictionPolicy
	config    *config.EngineConfig
	stats     *EngineStats
	stopChan  chan struct{}
	bgCleanup chan struct{}
}

// EngineStats 引擎统计
type EngineStats struct {
	mu            sync.RWMutex
	hits          int64
	misses        int64
	sets          int64
	deletes       int64
	evictions     int64
	expirations   int64
	memoryUsage   int64 // 字节
	gcCycles      int64 // GC cycles count
	poolHits      int64 // Object pool hits
	poolAllocs    int64 // Object pool allocations (new objects created)
	lastGCTime    time.Time
}

// NewStorageEngine 创建新的Storage engine
func NewStorageEngine(engineConfig *config.EngineConfig) interfaces.StorageEngine {
	if engineConfig == nil {
		engineConfig = config.DefaultEngineConfig()
	}

	// Pre-allocate map capacity based on MaxSize to reduce GC pressure
	initialCapacity := 64
	if engineConfig.MaxSize > 0 && engineConfig.MaxSize < 10000 {
		initialCapacity = engineConfig.MaxSize
	}

	engine := &StorageEngine{
		data:      make(map[string]interfaces.DataObject, initialCapacity),
		policy:    lru.NewLRUPolicy(engineConfig.MaxSize),
		config:    engineConfig,
		stats:     &EngineStats{},
		stopChan:  make(chan struct{}),
		bgCleanup: make(chan struct{}),
	}

	// 启动后台清理
	if engineConfig.BackgroundCleanupInterval > 0 {
		engine.startBackgroundCleanup()
	}

	return engine
}

// Set 存储对象
func (e *StorageEngine) Set(key string, obj interfaces.DataObject) error {
	// 验证Parameter
	if err := utils.ValidateCacheKey(key); err != nil {
		return err
	}

	// 检查内存可用性（仅在禁用自动清理时进行严格检查）
	if e.config.BackgroundCleanupInterval == 0 {
		if err := internal.CheckMemoryAvailability(e.config.MemoryThreshold); err != nil {
			return fmt.Errorf("memory limit exceeded: %w", err)
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// 检查是否需要淘汰（仅在配置了MaxSize时进行淘汰）
	if e.config.MaxSize > 0 && len(e.data) >= e.config.MaxSize && e.data[key] == nil {
		// 如果没有自动清理，则拒绝新数据
		if e.config.BackgroundCleanupInterval == 0 {
			return fmt.Errorf("storage capacity exceeded: max size %d reached", e.config.MaxSize)
		}
		e.evictOne()
	}

	// 再次检查内存（添加对象后的预估内存使用）
	if e.config.BackgroundCleanupInterval == 0 {
		// 预估新增对象的大小
		estimatedSize := int64(obj.Size())
		e.stats.updateMemoryUsage(estimatedSize)

		// 检查是否超过内存阈值
		if err := internal.CheckMemoryAvailability(e.config.MemoryThreshold); err != nil {
			// 回滚内存使用统计
			e.stats.updateMemoryUsage(-estimatedSize)
			return fmt.Errorf("insufficient memory for new object: %w", err)
		}
	}

	e.data[key] = obj
	e.policy.Set(key)
	e.stats.recordSet()

	return nil
}

// Get Get object
func (e *StorageEngine) Get(key string) (interfaces.DataObject, bool) {
	// 验证Parameter
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

	// Check expiration
	if obj.IsExpired() {
		e.deleteExpired(key)
		e.stats.recordMiss()
		e.stats.recordExpiration()
		return nil, false
	}

	e.policy.Access(key)
	e.stats.recordHit()
	return obj, true
}

// deleteExpired Synchronously delete expired key（避免竞态条件）
func (e *StorageEngine) deleteExpired(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if obj, exists := e.data[key]; exists && obj.IsExpired() {
		// Return object to pool before deletion
		e.returnObjectToPool(obj)
		delete(e.data, key)
		e.policy.Delete(key)
	}
}

// Delete Delete object
func (e *StorageEngine) Delete(key string) bool {
	// 验证Parameter
	if key == "" {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if obj, exists := e.data[key]; exists {
		// 更新内存使用统计（如果启用了内存管理）
		if e.config.BackgroundCleanupInterval == 0 {
			e.stats.updateMemoryUsage(-int64(obj.Size()))
		}

		// Return object to pool before deletion
		e.returnObjectToPool(obj)

		delete(e.data, key)
		e.policy.Delete(key)
		e.stats.recordDelete()
		return true
	}

	return false
}

// returnObjectToPool returns an object to the appropriate pool for reuse
func (e *StorageEngine) returnObjectToPool(obj interfaces.DataObject) {
	switch o := obj.(type) {
	case *types.StringObject:
		types.ReleaseStringObject(o)
		e.stats.recordPoolHit()
	case *types.ListObject:
		types.ReleaseListObject(o)
		e.stats.recordPoolHit()
	case *types.HashObject:
		types.ReleaseHashObject(o)
		e.stats.recordPoolHit()
	default:
		// Object type not supported for pooling
		e.stats.recordPoolAlloc()
	}
}

// Exists Check if key exists
func (e *StorageEngine) Exists(key string) bool {
	// 验证Parameter
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
		e.deleteExpired(key)
		return false
	}

	return true
}

// Keys Get all keys
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

	// Return all objects to pool before clearing
	for _, obj := range e.data {
		e.returnObjectToPool(obj)
	}

	e.data = make(map[string]interfaces.DataObject, len(e.data))
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

// Type Get key type
func (e *StorageEngine) Type(key string) (interfaces.DataType, bool) {
	e.mu.RLock()
	obj, exists := e.data[key]
	e.mu.RUnlock()

	if !exists {
		return "", false
	}

	if obj.IsExpired() {
		e.deleteExpired(key)
		return "", false
	}

	return obj.Type(), true
}

// Expire Set expiration time
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
	// 验证Parameter
	if key == "" {
		return -1, false
	}

	e.mu.RLock()
	obj, exists := e.data[key]
	e.mu.RUnlock()

	if !exists {
		return -1, false
	}

	return utils.CalculateRemainingTTL(obj.ExpiresAt())
}

// Stats Get statistics
func (e *StorageEngine) Stats() interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Get GC stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Update GC cycle count if it has changed
	if memStats.NumGC > 0 {
		e.stats.updateGCCycles(int64(memStats.NumGC))
	}

	return map[string]interface{}{
		"hits":         e.stats.hits,
		"misses":       e.stats.misses,
		"sets":         e.stats.sets,
		"deletes":      e.stats.deletes,
		"evictions":    e.stats.evictions,
		"expirations":  e.stats.expirations,
		"memory":       e.stats.memoryUsage,
		"keys":         len(e.data),
		"hit_rate":     e.stats.hitRate(),
		"gc_cycles":    e.stats.gcCycles,
		"pool_hits":    e.stats.poolHits,
		"pool_allocs":  e.stats.poolAllocs,
		"heap_alloc":   memStats.HeapAlloc,
		"heap_sys":     memStats.HeapSys,
		"num_gc":       memStats.NumGC,
		"gc_cpu_frac":  memStats.GCCPUFraction,
	}
}

// evictOne 淘汰一个键
func (e *StorageEngine) evictOne() {
	if key := e.policy.Evict(); key != "" {
		if obj, exists := e.data[key]; exists {
			// Return object to pool before eviction
			e.returnObjectToPool(obj)
		}
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
			// Return object to pool before deletion
			e.returnObjectToPool(obj)
			delete(e.data, key)
			e.policy.Delete(key)
			e.stats.recordExpiration()
		}
	}
}

// GetConfig 获取引擎配置
func (e *StorageEngine) GetConfig() *config.EngineConfig {
	return e.config
}

// Close 关闭引擎
func (e *StorageEngine) Close() {
	close(e.stopChan)
}

// EngineStats Method实现

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

func (s *EngineStats) recordPoolHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.poolHits++
}

func (s *EngineStats) recordPoolAlloc() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.poolAllocs++
}

func (s *EngineStats) updateGCCycles(cycles int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gcCycles = cycles
	s.lastGCTime = time.Now()
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
	s.gcCycles = 0
	s.poolHits = 0
	s.poolAllocs = 0
}

// updateMemoryUsage 更新内存使用统计
func (s *EngineStats) updateMemoryUsage(delta int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memoryUsage += delta
}
