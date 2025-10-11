package types

import (
	"sync"
	"time"
)

// ===== 缓存核心结构体 =====

// CacheItem 缓存项结构体
type CacheItem struct {
	Key         string      `json:"key"`          // 缓存键
	Value       interface{} `json:"value"`        // 缓存值
	ExpiresAt   time.Time   `json:"expires_at"`   // 过期时间，零值表示永不过期
	CreatedAt   time.Time   `json:"created_at"`   // 创建时间
	LastAccess  time.Time   `json:"last_access"`  // 最后访问时间
	AccessCount int64       `json:"access_count"` // 访问次数
	Size        int64       `json:"size"`         // 缓存项大小（字节）
}

// IsExpired 检查缓存项是否过期
func (item *CacheItem) IsExpired() bool {
	return !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt)
}

// UpdateAccess 更新访问信息
func (item *CacheItem) UpdateAccess() {
	item.LastAccess = time.Now()
	item.AccessCount++
}

// CacheStats 缓存统计信息结构体
type CacheStats struct {
	Hits       int64     `json:"hits"`        // 命中次数
	Misses     int64     `json:"misses"`      // 未命中次数
	HitRate    float64   `json:"hit_rate"`    // 命中率
	Size       int       `json:"size"`        // 当前缓存项数量
	MaxSize    int       `json:"max_size"`    // 最大缓存项数量
	Evictions  int64     `json:"evictions"`   // 淘汰次数
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	LastAccess time.Time `json:"last_access"` // 最后访问时间
}

// CacheConfig 缓存配置结构体
type CacheConfig struct {
	MaxSize              int           `json:"max_size"`               // 最大缓存项数量
	Shards               int           `json:"shards"`                 // 分片数量
	DefaultTTL           time.Duration `json:"default_ttl"`            // 默认TTL
	CleanupInterval      time.Duration `json:"cleanup_interval"`       // 清理间隔
	EvictionPolicy       string        `json:"eviction_policy"`        // 淘汰策略
	EnableStatistics     bool          `json:"enable_statistics"`      // 是否启用统计
	EnableLazyExpiration bool          `json:"enable_lazy_expiration"` // 是否启用延迟过期
	EnableMetrics        bool          `json:"enable_metrics"`         // 是否启用指标
	Serializer           string        `json:"serializer"`             // 序列化器类型
}

// ===== 缓存分片结构体 =====

// CacheShard 缓存分片结构体
type CacheShard struct {
	Items  map[string]*CacheItem `json:"items"`  // 缓存项映射
	Lock   sync.RWMutex          `json:"-"`      // 读写锁（不序列化）
	Policy interface{}           `json:"policy"` // 淘汰策略
	Stats  ShardStats            `json:"stats"`  // 分片统计
}

// ShardStats 分片统计信息
type ShardStats struct {
	Hits      int64     `json:"hits"`      // 分片命中次数
	Misses    int64     `json:"misses"`    // 分片未命中次数
	Evictions int64     `json:"evictions"` // 分片淘汰次数
	Size      int       `json:"size"`      // 分片当前大小
	LastSync  time.Time `json:"last_sync"` // 最后同步时间
}

// ===== 全局管理结构体 =====

// ManagerStats 管理器统计信息
type ManagerStats struct {
	TotalCaches    int                    `json:"total_caches"`     // 总缓存数量
	TotalItems     int64                  `json:"total_items"`      // 总缓存项数量
	TotalHits      int64                  `json:"total_hits"`       // 总命中次数
	TotalMisses    int64                  `json:"total_misses"`     // 总未命中次数
	OverallHitRate float64                `json:"overall_hit_rate"` // 总体命中率
	CacheStats     map[string]interface{} `json:"cache_stats"`      // 各缓存统计信息
	StartTime      time.Time              `json:"start_time"`       // 启动时间
	LastUpdateTime time.Time              `json:"last_update_time"` // 最后更新时间
}

// ManagerConfig 管理器配置
type ManagerConfig struct {
	MaxCaches           int           `json:"max_caches"`            // 最大缓存数量
	DefaultCacheConfig  CacheConfig   `json:"default_cache_config"`  // 默认缓存配置
	EnableGlobalStats   bool          `json:"enable_global_stats"`   // 是否启用全局统计
	StatsUpdateInterval time.Duration `json:"stats_update_interval"` // 统计更新间隔
}

// ===== 淘汰策略结构体 =====

// PolicyStats 策略统计信息
type PolicyStats struct {
	PolicyType string    `json:"policy_type"`  // 策略类型
	Operations int64     `json:"operations"`   // 操作次数
	Evictions  int64     `json:"evictions"`    // 淘汰次数
	LastOpTime time.Time `json:"last_op_time"` // 最后操作时间
}

// LRUEntry LRU策略的链表节点
type LRUEntry struct {
	Key   string      `json:"key"`   // 缓存键
	Value interface{} `json:"value"` // 缓存值
}

// LFUEntry LFU策略的堆节点
type LFUEntry struct {
	Key       string `json:"key"`       // 缓存键
	Frequency int    `json:"frequency"` // 访问频率
	Index     int    `json:"index"`     // 堆索引
}

// FIFOEntry FIFO策略的队列节点
type FIFOEntry struct {
	Key     string    `json:"key"`      // 缓存键
	AddTime time.Time `json:"add_time"` // 添加时间
}

// ===== 序列化相关结构体 =====

// SerializationConfig 序列化配置
type SerializationConfig struct {
	Encoding    string                 `json:"encoding"`    // 编码方式
	Compression bool                   `json:"compression"` // 是否压缩
	Options     map[string]interface{} `json:"options"`     // 序列化选项
}

// SerializedData 序列化数据结构
type SerializedData struct {
	Version   string                 `json:"version"`   // 版本号
	Timestamp time.Time              `json:"timestamp"` // 时间戳
	Data      map[string]interface{} `json:"data"`      // 实际数据
	Metadata  map[string]interface{} `json:"metadata"`  // 元数据
}

// ===== 错误处理结构体 =====

// ErrorInfo 错误信息结构体
type ErrorInfo struct {
	Code      int                    `json:"code"`      // 错误代码
	Message   string                 `json:"message"`   // 错误消息
	Details   map[string]interface{} `json:"details"`   // 错误详情
	Timestamp time.Time              `json:"timestamp"` // 错误时间
	Stack     string                 `json:"stack"`     // 错误堆栈
}

// ErrorStats 错误统计信息
type ErrorStats struct {
	TotalErrors   int64            `json:"total_errors"`    // 总错误数
	ErrorsByType  map[string]int64 `json:"errors_by_type"`  // 按类型统计的错误
	ErrorsByCode  map[int]int64    `json:"errors_by_code"`  // 按代码统计的错误
	LastError     *ErrorInfo       `json:"last_error"`      // 最后一个错误
	ErrorRate     float64          `json:"error_rate"`      // 错误率
	LastResetTime time.Time        `json:"last_reset_time"` // 最后重置时间
}

// ===== 指标监控结构体 =====

// Metrics 指标数据结构
type Metrics struct {
	CacheMetrics   *CacheMetrics    `json:"cache_metrics"`   // 缓存指标
	SystemMetrics  *SystemMetrics   `json:"system_metrics"`  // 系统指标
	CustomMetrics  map[string]int64 `json:"custom_metrics"`  // 自定义指标
	CollectionTime time.Time        `json:"collection_time"` // 收集时间
}

// CacheMetrics 缓存指标
type CacheMetrics struct {
	OperationRate   float64 `json:"operation_rate"`   // 操作速率
	AverageLatency  float64 `json:"average_latency"`  // 平均延迟
	P95Latency      float64 `json:"p95_latency"`      // P95延迟
	P99Latency      float64 `json:"p99_latency"`      // P99延迟
	Throughput      float64 `json:"throughput"`       // 吞吐量
	MemoryUsage     int64   `json:"memory_usage"`     // 内存使用
	CompressionRate float64 `json:"compression_rate"` // 压缩率
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPUUsage       float64 `json:"cpu_usage"`       // CPU使用率
	MemoryUsage    float64 `json:"memory_usage"`    // 内存使用率
	GoroutineCount int     `json:"goroutine_count"` // 协程数量
	GCCount        int64   `json:"gc_count"`        // GC次数
	PauseTime      int64   `json:"pause_time"`      // GC暂停时间
}

// ===== 事件通知结构体 =====

// CacheEvent 缓存事件结构体
type CacheEvent struct {
	Type      string                 `json:"type"`       // 事件类型
	CacheName string                 `json:"cache_name"` // 缓存名称
	Key       string                 `json:"key"`        // 缓存键
	Value     interface{}            `json:"value"`      // 缓存值
	Timestamp time.Time              `json:"timestamp"`  // 事件时间
	Metadata  map[string]interface{} `json:"metadata"`   // 事件元数据
}

// EventConfig 事件配置
type EventConfig struct {
	Enabled    bool              `json:"enabled"`     // 是否启用事件
	EventTypes []string          `json:"event_types"` // 启用的事件类型
	BufferSize int               `json:"buffer_size"` // 事件缓冲区大小
	Handlers   map[string]string `json:"handlers"`    // 事件处理器
	Filters    map[string]string `json:"filters"`     // 事件过滤器
}

// ===== 健康检查结构体 =====

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string                  `json:"status"`    // 健康状态
	Timestamp time.Time               `json:"timestamp"` // 检查时间
	Checks    map[string]*HealthCheck `json:"checks"`    // 各项检查结果
	Summary   *HealthSummary          `json:"summary"`   // 健康摘要
}

// HealthCheck 单项健康检查
type HealthCheck struct {
	Name         string                 `json:"name"`          // 检查名称
	Status       string                 `json:"status"`        // 检查状态
	Message      string                 `json:"message"`       // 检查消息
	LastCheck    time.Time              `json:"last_check"`    // 最后检查时间
	ResponseTime int64                  `json:"response_time"` // 响应时间（毫秒）
	Details      map[string]interface{} `json:"details"`       // 检查详情
}

// HealthSummary 健康摘要
type HealthSummary struct {
	TotalChecks  int     `json:"total_checks"`  // 总检查数
	PassedChecks int     `json:"passed_checks"` // 通过检查数
	FailedChecks int     `json:"failed_checks"` // 失败检查数
	OverallScore float64 `json:"overall_score"` // 总体评分
}

// ===== 配置验证结构体 =====

// ValidationResult 配置验证结果
type ValidationResult struct {
	Valid    bool                `json:"valid"`    // 是否有效
	Errors   []ValidationError   `json:"errors"`   // 验证错误
	Warnings []ValidationWarning `json:"warnings"` // 验证警告
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string      `json:"field"`   // 字段名
	Message string      `json:"message"` // 错误消息
	Value   interface{} `json:"value"`   // 错误值
}

// ValidationWarning 验证警告
type ValidationWarning struct {
	Field   string      `json:"field"`   // 字段名
	Message string      `json:"message"` // 警告消息
	Value   interface{} `json:"value"`   // 警告值
}
