package interfaces

import (
	"time"

	"github.com/scache/types"
)

// 重新导出核心类型
type (
	// CacheItem 缓存项
	CacheItem = types.CacheItem
	// CacheStats 缓存统计信息
	CacheStats = types.CacheStats
	// CacheConfig 缓存配置
	CacheConfig = types.CacheConfig
	// CacheShard 缓存分片
	CacheShard = types.CacheShard
	// ManagerStats 管理器统计信息
	ManagerStats = types.ManagerStats
	// ErrorInfo 错误信息
	ErrorInfo = types.ErrorInfo
	// ValidationResult 验证结果
	ValidationResult = types.ValidationResult
	// HealthStatus 健康状态
	HealthStatus = types.HealthStatus
)

// Cache 缓存接口
type Cache interface {
	// 基本操作
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) bool
	Exists(key string) bool
	Clear() error
	Close() error

	// 批量操作
	SetBatch(items map[string]interface{}) error
	GetBatch(keys []string) map[string]interface{}
	DeleteBatch(keys []string) map[string]bool

	// 状态查询
	Size() int
	Keys() []string
	Stats() CacheStats
}

// EvictionPolicy 淘汰策略接口
type EvictionPolicy interface {
	// 基本操作
	OnAccess(key string)
	OnAdd(key string)
	OnRemove(key string)
	ShouldEvict() (string, bool)
	SetMaxSize(size int)
	Clear()

	// 状态查询
	Len() int
	Keys() []string
	Contains(key string) bool
}

// Manager 缓存管理器接口
type Manager interface {
	// 缓存管理
	Register(name string, cache Cache) error
	Get(name string) (Cache, error)
	Remove(name string) error
	List() []string
	Clear() error
	Close() error

	// 状态查询
	Stats() map[string]CacheStats
	Size() int
	Exists(name string) bool
}

// Serializer 序列化接口
type Serializer interface {
	// 基本序列化
	Serialize(value interface{}) ([]byte, error)
	Deserialize(data []byte, target interface{}) error

	// 类型支持
	SupportType(value interface{}) bool
	SupportedContentTypes() []string

	// 配置
	Configure(config map[string]interface{}) error
	GetConfig() map[string]interface{}
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event *types.CacheEvent) error
	ShouldHandle(eventType string) bool
	GetName() string
}

// Validator 配置验证器接口
type Validator interface {
	Validate(config interface{}) ValidationResult
	GetSupportedTypes() []string
}

// HealthChecker 健康检查器接口
type HealthChecker interface {
	Check() *types.HealthCheck
	GetName() string
	IsEnabled() bool
}

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	Collect() *types.Metrics
	Start() error
	Stop() error
	IsRunning() bool
}
