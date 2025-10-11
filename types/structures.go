package types

import (
	"scache/constants"
	"sync"
	"time"
)

// CacheItem 缓存项结构
type CacheItem struct {
	Key       string      // 缓存键
	Value     interface{} // 缓存值
	ExpiresAt time.Time   // 过期时间
	CreatedAt time.Time   // 创建时间
	AccessAt  time.Time   // 最后访问时间
	Hits      int64       // 访问次数
}

// IsExpired 检查缓存项是否过期
func (item *CacheItem) IsExpired() bool {
	if item.ExpiresAt.IsZero() {
		return false // 零值表示永不过期
	}
	return time.Now().After(item.ExpiresAt)
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// 默认过期时间，0 表示永不过期
	DefaultExpiration time.Duration

	// 清理间隔
	CleanupInterval time.Duration

	// 最大容量，0 表示无限制
	MaxSize int

	// 是否启用统计信息
	EnableStats bool

	// 初始容量
	InitialCapacity int
}

// DefaultCacheConfig 返回默认配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultExpiration: constants.DefaultExpiration,      // 永不过期
		CleanupInterval:   constants.DefaultCleanupInterval, // 10分钟清理一次
		MaxSize:           constants.DefaultMaxSize,         // 无容量限制
		EnableStats:       true,                             // 启用统计
		InitialCapacity:   constants.DefaultInitialCapacity, // 初始容量16
	}
}

// CacheOption 配置选项函数
type CacheOption func(*CacheConfig)

// WithDefaultExpiration 设置默认过期时间
func WithDefaultExpiration(d time.Duration) CacheOption {
	return func(c *CacheConfig) {
		c.DefaultExpiration = d
	}
}

// WithCleanupInterval 设置清理间隔
func WithCleanupInterval(d time.Duration) CacheOption {
	return func(c *CacheConfig) {
		c.CleanupInterval = d
	}
}

// WithMaxSize 设置最大容量
func WithMaxSize(size int) CacheOption {
	return func(c *CacheConfig) {
		c.MaxSize = size
	}
}

// WithStats 设置是否启用统计
func WithStats(enable bool) CacheOption {
	return func(c *CacheConfig) {
		c.EnableStats = enable
	}
}

// WithInitialCapacity 设置初始容量
func WithInitialCapacity(capacity int) CacheOption {
	return func(c *CacheConfig) {
		c.InitialCapacity = capacity
	}
}

// CacheStats 缓存统计信息（线程安全版本）
type CacheStats struct {
	mu      sync.RWMutex
	hits    int64
	misses  int64
	sets    int64
	deletes int64
}

// Hit 增加命中计数
func (s *CacheStats) Hit() {
	s.mu.Lock()
	s.hits++
	s.mu.Unlock()
}

// Miss 增加未命中计数
func (s *CacheStats) Miss() {
	s.mu.Lock()
	s.misses++
	s.mu.Unlock()
}

// Set 增加设置计数
func (s *CacheStats) Set() {
	s.mu.Lock()
	s.sets++
	s.mu.Unlock()
}

// Delete 增加删除计数
func (s *CacheStats) Delete() {
	s.mu.Lock()
	s.deletes++
	s.mu.Unlock()
}

// GetStats 获取统计信息快照
func (s *CacheStats) GetStats() (hits, misses, sets, deletes int64) {
	s.mu.RLock()
	hits = s.hits
	misses = s.misses
	sets = s.sets
	deletes = s.deletes
	s.mu.RUnlock()
	return
}

// HitRate 计算命中率
func (s *CacheStats) HitRate() float64 {
	hits, misses, _, _ := s.GetStats()
	total := hits + misses
	if total == constants.DefaultExpiration {
		return constants.DefaultExpiration
	}
	return float64(hits) / float64(total)
}

// Reset 重置统计信息
func (s *CacheStats) Reset() {
	s.mu.Lock()
	s.hits = constants.DefaultExpiration
	s.misses = constants.DefaultExpiration
	s.sets = constants.DefaultExpiration
	s.deletes = constants.DefaultExpiration
	s.mu.Unlock()
}
