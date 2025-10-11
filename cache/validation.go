package cache

import (
	"github.com/scache/constants"
)

// validateConfig 验证配置
func validateConfig(config *Config) {
	if config.MaxSize <= 0 {
		config.MaxSize = constants.DefaultMaxSize
	}
	if config.Shards <= 0 {
		config.Shards = constants.DefaultShards
	}
	if config.CleanupInterval <= 0 {
		config.CleanupInterval = constants.DefaultCleanupInterval
	}
	// 验证策略类型
	switch config.EvictionPolicy {
	case constants.LRUStrategy, constants.LFUStrategy, constants.FIFOStrategy:
		// 有效策略
	default:
		config.EvictionPolicy = constants.LRUStrategy
	}
}
