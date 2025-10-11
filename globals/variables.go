package globals

import (
	"sync"
	"time"

	"github.com/scache/interfaces"
)

// 全局缓存管理器实例
var (
	// globalManager 全局缓存管理器单例
	globalManager Manager
	// managerOnce 确保管理器只初始化一次
	managerOnce sync.Once
	// managerMutex 管理器访问锁
	managerMutex sync.RWMutex
)

// 全局默认缓存实例
var (
	// defaultCacheInstance 默认缓存实例
	defaultCacheInstance interfaces.Cache
	// defaultCacheOnce 确保默认缓存只初始化一次
	defaultCacheOnce sync.Once
	// defaultCacheMutex 默认缓存访问锁
	defaultCacheMutex sync.RWMutex
)

// 全局配置变量
var (
	// defaultConfig 默认配置（可通过函数修改）
	defaultConfig = &Config{
		MaxSize:              10000,
		Shards:               16,
		DefaultTTL:           0,
		CleanupInterval:      10 * time.Minute,
		EvictionPolicy:       "lru",
		EnableStatistics:     false,
		EnableLazyExpiration: true,
		EnableMetrics:        false,
	}
	// configMutex 配置访问锁
	configMutex sync.RWMutex
)

// 全局统计变量
var (
	// globalStats 全局统计信息
	globalStats = &GlobalStats{
		TotalCaches:    0,
		TotalItems:     0,
		TotalHits:      0,
		TotalMisses:    0,
		StartTime:      time.Now(),
		LastUpdateTime: time.Now(),
	}
	// statsMutex 统计信息访问锁
	statsMutex sync.RWMutex
)

// 全局错误变量
var (
	// commonErrors 常见错误集合
	commonErrors = &ErrorSet{
		CacheNotFound:        &CacheError{Code: 1001, Message: "cache not found"},
		InvalidCacheName:     &CacheError{Code: 1002, Message: "invalid cache name"},
		CacheAlreadyExists:   &CacheError{Code: 1003, Message: "cache already exists"},
		InvalidStrategy:      &CacheError{Code: 1004, Message: "invalid eviction strategy"},
		KeyNotFound:          &CacheError{Code: 2001, Message: "key not found"},
		KeyTooLong:           &CacheError{Code: 2002, Message: "key too long"},
		KeyEmpty:             &CacheError{Code: 2003, Message: "key empty"},
		ValueTooLarge:        &CacheError{Code: 2004, Message: "value too large"},
		CacheClosed:          &CacheError{Code: 3001, Message: "cache is closed"},
		SerializationError:   &CacheError{Code: 4001, Message: "serialization error"},
		DeserializationError: &CacheError{Code: 4002, Message: "deserialization error"},
	}
)

// 全局注册表
var (
	// strategyRegistry 策略注册表
	strategyRegistry = make(map[string]StrategyFactory)
	// serializerRegistry 序列化器注册表
	serializerRegistry = make(map[string]SerializerFactory)
	// registryMutex 注册表访问锁
	registryMutex sync.RWMutex
)

// Manager 全局管理器类型（前向声明）
type Manager interface {
	Register(name string, cache interfaces.Cache) error
	Get(name string) (interfaces.Cache, error)
	Remove(name string) error
	List() []string
	Clear() error
	Close() error
	Stats() map[string]interfaces.CacheStats
	Size() int
	Exists(name string) bool
}

// Config 配置结构
type Config struct {
	MaxSize              int
	Shards               int
	DefaultTTL           time.Duration
	CleanupInterval      time.Duration
	EvictionPolicy       string
	EnableStatistics     bool
	EnableLazyExpiration bool
	EnableMetrics        bool
	Serializer           string
}

// GlobalStats 全局统计信息
type GlobalStats struct {
	TotalCaches    int
	TotalItems     int64
	TotalHits      int64
	TotalMisses    int64
	StartTime      time.Time
	LastUpdateTime time.Time
}

// CacheError 缓存错误类型
type CacheError struct {
	Code    int
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

// ErrorSet 错误集合
type ErrorSet struct {
	CacheNotFound        *CacheError
	InvalidCacheName     *CacheError
	CacheAlreadyExists   *CacheError
	InvalidStrategy      *CacheError
	KeyNotFound          *CacheError
	KeyTooLong           *CacheError
	KeyEmpty             *CacheError
	ValueTooLarge        *CacheError
	CacheClosed          *CacheError
	SerializationError   *CacheError
	DeserializationError *CacheError
}

// StrategyFactory 策略工厂函数类型
type StrategyFactory func(int) interfaces.EvictionPolicy

// SerializerFactory 序列化器工厂函数类型
type SerializerFactory func() interfaces.Serializer

// 全局访问函数

// GetGlobalManager 获取全局管理器实例
func GetGlobalManager() Manager {
	managerOnce.Do(func() {
		// 这里会在实际的实现中被注入
	})
	return globalManager
}

// SetGlobalManager 设置全局管理器实例（用于初始化）
func SetGlobalManager(manager Manager) {
	managerMutex.Lock()
	defer managerMutex.Unlock()
	globalManager = manager
}

// GetDefaultCacheInstance 获取默认缓存实例
func GetDefaultCacheInstance() interfaces.Cache {
	return defaultCacheInstance
}

// SetDefaultCacheInstance 设置默认缓存实例（用于初始化）
func SetDefaultCacheInstance(cache interfaces.Cache) {
	defaultCacheMutex.Lock()
	defer defaultCacheMutex.Unlock()
	defaultCacheInstance = cache
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return defaultConfig
}

// UpdateDefaultConfig 更新默认配置
func UpdateDefaultConfig(updateFn func(*Config)) {
	configMutex.Lock()
	defer configMutex.Unlock()
	updateFn(defaultConfig)
}

// GetGlobalStats 获取全局统计信息
func GetGlobalStats() *GlobalStats {
	statsMutex.RLock()
	defer statsMutex.RUnlock()
	// 返回副本以避免并发问题
	return &GlobalStats{
		TotalCaches:    globalStats.TotalCaches,
		TotalItems:     globalStats.TotalItems,
		TotalHits:      globalStats.TotalHits,
		TotalMisses:    globalStats.TotalMisses,
		StartTime:      globalStats.StartTime,
		LastUpdateTime: globalStats.LastUpdateTime,
	}
}

// UpdateGlobalStats 更新全局统计信息
func UpdateGlobalStats(updateFn func(*GlobalStats)) {
	statsMutex.Lock()
	defer statsMutex.Unlock()
	updateFn(globalStats)
	globalStats.LastUpdateTime = time.Now()
}

// GetCommonErrors 获取常见错误集合
func GetCommonErrors() *ErrorSet {
	return commonErrors
}

// RegisterStrategy 注册策略
func RegisterStrategy(name string, factory StrategyFactory) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	strategyRegistry[name] = factory
}

// GetStrategy 获取策略工厂
func GetStrategy(name string) (StrategyFactory, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	factory, exists := strategyRegistry[name]
	return factory, exists
}

// ListStrategies 列出所有注册的策略
func ListStrategies() []string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	strategies := make([]string, 0, len(strategyRegistry))
	for name := range strategyRegistry {
		strategies = append(strategies, name)
	}
	return strategies
}

// RegisterSerializer 注册序列化器
func RegisterSerializer(name string, factory SerializerFactory) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	serializerRegistry[name] = factory
}

// GetSerializer 获取序列化器工厂
func GetSerializer(name string) (SerializerFactory, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	factory, exists := serializerRegistry[name]
	return factory, exists
}

// ListSerializers 列出所有注册的序列化器
func ListSerializers() []string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	serializers := make([]string, 0, len(serializerRegistry))
	for name := range serializerRegistry {
		serializers = append(serializers, name)
	}
	return serializers
}
