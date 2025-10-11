package cache

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/scache/interfaces"
	"github.com/scache/types"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	shards    []*cacheShard
	config    *Config
	stats     *interfaces.CacheStats
	statsLock sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// cacheShard 缓存分片，减少锁竞争
type cacheShard struct {
	items  map[string]*types.CacheItem
	lock   sync.RWMutex
	policy interfaces.EvictionPolicy
}

// NewMemoryCache 创建新的内存缓存
func NewMemoryCache(opts ...Option) *MemoryCache {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	validateConfig(config)

	ctx, cancel := context.WithCancel(context.Background())

	cache := &MemoryCache{
		config: config,
		stats: &CacheStats{
			MaxSize:   config.MaxSize,
			CreatedAt: time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化分片
	cache.shards = make([]*cacheShard, config.Shards)
	shardMaxSize := config.MaxSize / config.Shards
	for i := 0; i < config.Shards; i++ {
		cache.shards[i] = &cacheShard{
			items: make(map[string]*types.CacheItem),
		}
		// 初始化淘汰策略
		cache.shards[i].policy = newEvictionPolicy(config.EvictionPolicy, shardMaxSize)
	}

	// 启动清理协程
	go cache.cleanup()

	return cache
}

// getShard 获取键对应的分片
func (c *MemoryCache) getShard(key string) *cacheShard {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return c.shards[hash.Sum32()%uint32(c.config.Shards)]
}

// Set 设置缓存项
func (c *MemoryCache) Set(key string, value interface{}) error {
	if c.config.DefaultTTL == 0 {
		return c.SetWithTTL(key, value, 0)
	}
	return c.SetWithTTL(key, value, c.config.DefaultTTL)
}

// SetWithTTL 设置带TTL的缓存项
func (c *MemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	shard := c.getShard(key)

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	item := &types.CacheItem{
		Key:         key,
		Value:       value,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		AccessCount: 1,
		LastAccess:  time.Now(),
	}

	shard.lock.Lock()
	defer shard.lock.Unlock()

	shardMaxSize := c.config.MaxSize / c.config.Shards
	if len(shard.items) >= shardMaxSize {
		if evictKey, shouldEvict := shard.policy.ShouldEvict(); shouldEvict {
			delete(shard.items, evictKey)
			shard.policy.OnRemove(evictKey)
		}
	}

	if _, exists := shard.items[key]; exists {
		shard.policy.OnRemove(key)
	}

	shard.items[key] = item
	shard.policy.OnAdd(key)

	c.updateLastAccess()
	return nil
}

// Get 获取缓存项
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	shard := c.getShard(key)

	shard.lock.RLock()
	item, exists := shard.items[key]
	shard.lock.RUnlock()

	if !exists {
		c.recordMiss()
		return nil, false
	}

	// 检查是否过期
	if c.config.EnableLazyExpiration && item.IsExpired() {
		c.Delete(key)
		c.recordMiss()
		return nil, false
	}

	// 更新访问信息
	shard.lock.Lock()
	item.AccessCount++
	item.LastAccess = time.Now()
	shard.lock.Unlock()

	shard.policy.OnAccess(key)
	c.recordHit()
	c.updateLastAccess()

	return item.Value, true
}

// Delete 删除缓存项
func (c *MemoryCache) Delete(key string) bool {
	shard := c.getShard(key)

	shard.lock.Lock()
	defer shard.lock.Unlock()

	if _, exists := shard.items[key]; exists {
		delete(shard.items, key)
		shard.policy.OnRemove(key)
		return true
	}
	return false
}

// Exists 检查缓存项是否存在
func (c *MemoryCache) Exists(key string) bool {
	shard := c.getShard(key)

	shard.lock.RLock()
	item, exists := shard.items[key]
	shard.lock.RUnlock()

	if !exists {
		return false
	}

	// 检查是否过期
	if c.config.EnableLazyExpiration && item.IsExpired() {
		c.Delete(key)
		return false
	}

	return true
}

// Clear 清空所有缓存
func (c *MemoryCache) Clear() error {
	for _, shard := range c.shards {
		shard.lock.Lock()
		shard.items = make(map[string]*types.CacheItem)
		shard.lock.Unlock()
	}

	c.resetStats()
	return nil
}

// SetBatch 批量设置缓存项
func (c *MemoryCache) SetBatch(items map[string]interface{}) error {
	for key, value := range items {
		if err := c.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

// GetBatch 批量获取缓存项
func (c *MemoryCache) GetBatch(keys []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range keys {
		if value, exists := c.Get(key); exists {
			result[key] = value
		}
	}
	return result
}

// DeleteBatch 批量删除缓存项
func (c *MemoryCache) DeleteBatch(keys []string) map[string]bool {
	result := make(map[string]bool)
	for _, key := range keys {
		result[key] = c.Delete(key)
	}
	return result
}

// Size 获取缓存项总数
func (c *MemoryCache) Size() int {
	total := 0
	for _, shard := range c.shards {
		shard.lock.RLock()
		total += len(shard.items)
		shard.lock.RUnlock()
	}
	return total
}

// Keys 获取所有键
func (c *MemoryCache) Keys() []string {
	var keys []string
	for _, shard := range c.shards {
		shard.lock.RLock()
		for key := range shard.items {
			keys = append(keys, key)
		}
		shard.lock.RUnlock()
	}
	return keys
}

// Stats 获取缓存统计信息
func (c *MemoryCache) Stats() CacheStats {
	c.statsLock.RLock()
	defer c.statsLock.RUnlock()

	stats := *c.stats
	stats.Size = c.Size()

	hits := atomic.LoadInt64(&c.stats.Hits)
	misses := atomic.LoadInt64(&c.stats.Misses)
	total := hits + misses
	if total > 0 {
		stats.HitRate = float64(hits) / float64(total)
	}

	return stats
}

// Close 关闭缓存
func (c *MemoryCache) Close() error {
	c.cancel()
	return c.Clear()
}

// cleanup 定期清理过期项
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.cleanupExpired()
		}
	}
}

// cleanupExpired 清理过期的缓存项
func (c *MemoryCache) cleanupExpired() {
	now := time.Now()
	for _, shard := range c.shards {
		var expiredKeys []string

		shard.lock.RLock()
		for key, item := range shard.items {
			if !item.ExpiresAt.IsZero() && now.After(item.ExpiresAt) {
				expiredKeys = append(expiredKeys, key)
			}
		}
		shard.lock.RUnlock()

		for _, key := range expiredKeys {
			c.Delete(key)
		}
	}
}

// recordHit 记录命中
func (c *MemoryCache) recordHit() {
	if !c.config.EnableStatistics {
		return
	}
	atomic.AddInt64(&c.stats.Hits, 1)
}

// recordMiss 记录未命中
func (c *MemoryCache) recordMiss() {
	if !c.config.EnableStatistics {
		return
	}
	atomic.AddInt64(&c.stats.Misses, 1)
}

// updateLastAccess 更新最后访问时间
func (c *MemoryCache) updateLastAccess() {
	c.statsLock.Lock()
	c.stats.LastAccess = time.Now()
	c.statsLock.Unlock()
}

// resetStats 重置统计信息
func (c *MemoryCache) resetStats() {
	c.statsLock.Lock()
	c.stats.Hits = 0
	c.stats.Misses = 0
	c.stats.HitRate = 0
	c.statsLock.Unlock()
}

// Serialize 序列化缓存数据
func (c *MemoryCache) Serialize() ([]byte, error) {
	data := make(map[string]*types.CacheItem)
	for _, shard := range c.shards {
		shard.lock.RLock()
		for key, item := range shard.items {
			data[key] = item
		}
		shard.lock.RUnlock()
	}
	return json.Marshal(data)
}

// Deserialize 反序列化缓存数据
func (c *MemoryCache) Deserialize(data []byte) error {
	var items map[string]*types.CacheItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	for key, item := range items {
		c.Set(key, item.Value)
	}
	return nil
}
