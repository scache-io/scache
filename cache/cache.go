package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"scache/constants"
	"scache/interfaces"
	"scache/policies/lru"
	"scache/types"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	mu          sync.RWMutex                // 读写锁
	items       map[string]*types.CacheItem // 缓存数据
	config      *types.CacheConfig          // 配置
	stats       *types.CacheStats           // 统计信息
	policy      interfaces.EvictionPolicy   // 淘汰策略
	stopChan    chan struct{}               // 停止清理协程的通道
	cleanupOnce sync.Once                   // 确保清理协程只启动一次
}

// NewCache 创建新的缓存实例
func NewCache(opts ...types.CacheOption) interfaces.Cache {
	config := types.DefaultCacheConfig()

	// 应用配置选项
	for _, opt := range opts {
		opt(config)
	}

	cache := &MemoryCache{
		items:    make(map[string]*types.CacheItem, config.InitialCapacity),
		config:   config,
		stats:    &types.CacheStats{},
		policy:   lru.NewLRUPolicy(config.MaxSize),
		stopChan: make(chan struct{}),
	}

	// 启动清理协程 - 只有在有过期机制时才需要
	if config.CleanupInterval > constants.DefaultExpiration && (config.DefaultExpiration > constants.DefaultExpiration) {
		cache.startCleanup()
	}

	return cache
}

// Set 设置缓存项
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	if key == "" {
		return errors.New("cache key cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 计算过期时间
	var expiresAt time.Time
	if ttl > constants.DefaultExpiration {
		expiresAt = time.Now().Add(ttl)
	} else if c.config.DefaultExpiration > constants.DefaultExpiration {
		expiresAt = time.Now().Add(c.config.DefaultExpiration)
	}

	// 创建或更新缓存项
	now := time.Now()
	item := &types.CacheItem{
		Key:       key,
		Value:     value,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		AccessAt:  now,
		Hits:      constants.DefaultExpiration,
	}

	// 如果 key 已存在，更新访问次数
	if oldItem, exists := c.items[key]; exists {
		item.Hits = oldItem.Hits
	}

	c.items[key] = item

	// 更新淘汰策略
	c.policy.Set(key)

	// 检查容量限制
	if c.config.MaxSize > constants.DefaultExpiration && len(c.items) > c.config.MaxSize {
		c.evict()
	}

	// 更新统计
	if c.config.EnableStats {
		c.stats.Set()
	}

	return nil
}

// Get 获取缓存项
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		if c.config.EnableStats {
			c.stats.Miss()
		}
		return nil, false
	}

	// 检查是否过期
	if item.IsExpired() {
		// 异步删除过期项
		go c.Delete(key)
		if c.config.EnableStats {
			c.stats.Miss()
		}
		return nil, false
	}

	// 更新访问信息
	c.mu.Lock()
	item.AccessAt = time.Now()
	item.Hits++
	c.mu.Unlock()

	// 更新淘汰策略
	c.policy.Access(key)

	// 更新统计
	if c.config.EnableStats {
		c.stats.Hit()
	}

	return item.Value, true
}

// Delete 删除缓存项
func (c *MemoryCache) Delete(key string) bool {
	if key == "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; exists {
		delete(c.items, key)
		c.policy.Delete(key)

		if c.config.EnableStats {
			c.stats.Delete()
		}
		return true
	}

	return false
}

// Exists 检查缓存项是否存在
func (c *MemoryCache) Exists(key string) bool {
	if key == "" {
		return false
	}

	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		return false
	}

	// 检查是否过期
	if item.IsExpired() {
		// 异步删除过期项
		go c.Delete(key)
		return false
	}

	return true
}

// Flush 清空所有缓存项
func (c *MemoryCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*types.CacheItem, c.config.InitialCapacity)
	c.policy.Clear()

	// 重置统计信息
	if c.config.EnableStats {
		c.stats.Reset()
	}
}

// Size 获取缓存项数量
func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// Stats 获取缓存统计信息
func (c *MemoryCache) Stats() interfaces.CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hits, misses, sets, deletes := c.stats.GetStats()

	return interfaces.CacheStats{
		Hits:    hits,
		Misses:  misses,
		Sets:    sets,
		Deletes: deletes,
		Size:    len(c.items),
		MaxSize: c.config.MaxSize,
		HitRate: c.stats.HitRate(),
	}
}

// evict 淘汰一个缓存项
func (c *MemoryCache) evict() {
	if key := c.policy.Evict(); key != "" {
		delete(c.items, key)
	}
}

// startCleanup 启动过期清理协程
func (c *MemoryCache) startCleanup() {
	c.cleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(c.config.CleanupInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					c.cleanupExpired()
				case <-c.stopChan:
					return
				}
			}
		}()
	})
}

// cleanupExpired 清理过期的缓存项
func (c *MemoryCache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if !item.ExpiresAt.IsZero() && now.After(item.ExpiresAt) {
			delete(c.items, key)
			c.policy.Delete(key)
		}
	}
}

// Close 关闭缓存，停止清理协程
func (c *MemoryCache) Close() {
	close(c.stopChan)
}

// Keys 获取所有缓存键
func (c *MemoryCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

// GetWithExpiration 获取缓存项和过期时间
func (c *MemoryCache) GetWithExpiration(key string) (interface{}, time.Time, bool) {
	if key == "" {
		return nil, time.Time{}, false
	}

	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		if c.config.EnableStats {
			c.stats.Miss()
		}
		return nil, time.Time{}, false
	}

	// 检查是否过期
	if item.IsExpired() {
		// 异步删除过期项
		go c.Delete(key)
		if c.config.EnableStats {
			c.stats.Miss()
		}
		return nil, time.Time{}, false
	}

	// 更新访问信息
	c.mu.Lock()
	item.AccessAt = time.Now()
	item.Hits++
	c.mu.Unlock()

	// 更新淘汰策略
	c.policy.Access(key)

	// 更新统计
	if c.config.EnableStats {
		c.stats.Hit()
	}

	return item.Value, item.ExpiresAt, true
}

// String 返回缓存的字符串表示
func (c *MemoryCache) String() string {
	stats := c.Stats()
	return fmt.Sprintf("Cache{Size: %d, Hits: %d, Misses: %d, HitRate: %.2f%%}",
		stats.Size, stats.Hits, stats.Misses, stats.HitRate*100)
}
