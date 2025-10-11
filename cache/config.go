package cache

import (
	"time"

	"github.com/scache/constants"
	"github.com/scache/types"
)

// Config 缓存配置（重新导出类型）
type Config = types.CacheConfig

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		MaxSize:              constants.DefaultMaxSize,
		DefaultTTL:           constants.DefaultTTL,
		CleanupInterval:      constants.DefaultCleanupInterval,
		EvictionPolicy:       constants.LRUStrategy,
		Shards:               constants.DefaultShards,
		EnableStatistics:     true,
		Serializer:           constants.JSONEncoding,
		EnableLazyExpiration: true,
		EnableMetrics:        false,
	}
}

// ValidateConfig 验证配置（外部函数）
func ValidateConfig(c *Config) error {
	if c.MaxSize <= 0 {
		c.MaxSize = constants.DefaultMaxSize
	}
	if c.Shards <= 0 {
		c.Shards = constants.DefaultShards
	}
	if c.CleanupInterval <= 0 {
		c.CleanupInterval = constants.DefaultCleanupInterval
	}
	// 验证策略类型
	switch c.EvictionPolicy {
	case constants.LRUStrategy, constants.LFUStrategy, constants.FIFOStrategy:
		// 有效策略
	default:
		c.EvictionPolicy = constants.LRUStrategy
	}
	return nil
}

// Option 配置选项函数
type Option func(*Config)

// WithMaxSize 设置最大缓存大小
func WithMaxSize(size int) Option {
	return func(c *Config) {
		c.MaxSize = size
	}
}

// WithDefaultTTL 设置默认TTL
func WithDefaultTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.DefaultTTL = ttl
	}
}

// WithEvictionPolicy 设置淘汰策略
func WithEvictionPolicy(policy string) Option {
	return func(c *Config) {
		c.EvictionPolicy = policy
	}
}

// WithShards 设置分片数量
func WithShards(shards int) Option {
	return func(c *Config) {
		c.Shards = shards
	}
}

// WithCleanupInterval 设置清理间隔
func WithCleanupInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.CleanupInterval = interval
	}
}

// WithStatistics 启用统计
func WithStatistics(enable bool) Option {
	return func(c *Config) {
		c.EnableStatistics = enable
	}
}

// WithSerializer 设置序列化器
func WithSerializer(serializer string) Option {
	return func(c *Config) {
		c.Serializer = serializer
	}
}

// WithLazyExpiration 启用懒过期
func WithLazyExpiration(enable bool) Option {
	return func(c *Config) {
		c.EnableLazyExpiration = enable
	}
}

// WithMetrics 启用指标
func WithMetrics(enable bool) Option {
	return func(c *Config) {
		c.EnableMetrics = enable
	}
}
