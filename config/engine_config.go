package config

import (
	"time"

	"github.com/scache-io/scache/constants"
)

// EngineConfig 存储引擎配置
type EngineConfig struct {
	MaxSize                   int           // 最大缓存数量
	MemoryThreshold           float64       // 内存阈值
	DefaultExpiration         time.Duration // 默认过期时间
	BackgroundCleanupInterval time.Duration // 后台清理间隔
}

// DefaultEngineConfig 默认引擎配置
func DefaultEngineConfig() *EngineConfig {
	return &EngineConfig{
		MaxSize:                   constants.DefaultMaxSize,         // 无限制
		MemoryThreshold:           constants.DefaultMemoryThreshold, // 80%
		DefaultExpiration:         constants.DefaultExpiration,      // 永不过期
		BackgroundCleanupInterval: constants.DefaultCleanupInterval, // 禁用自动清理
	}
}
