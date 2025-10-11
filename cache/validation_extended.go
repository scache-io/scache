package cache

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/scache/constants"
)

// ValidationErrors 多个验证错误的集合
type ValidationErrors []error

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}

	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// validateKey 验证缓存键的有效性
func validateKey(key string) error {
	if key == "" {
		return errors.New(constants.ErrKeyEmpty)
	}

	if len(key) > constants.MaxKeyLength {
		return fmt.Errorf("%s: max %d characters, got %d",
			constants.ErrKeyTooLong, constants.MaxKeyLength, len(key))
	}

	if len(key) < constants.MinKeyLength {
		return fmt.Errorf("key too short: min %d characters, got %d",
			constants.MinKeyLength, len(key))
	}

	// 检查键是否包含无效字符
	for _, r := range key {
		if r < 32 || r > 126 { // 控制字符
			return fmt.Errorf("key contains invalid character: %q", r)
		}
	}

	return nil
}

// validateValue 验证缓存值的有效性
func validateValue(value interface{}) error {
	if value == nil {
		return errors.New("value cannot be nil")
	}

	// 检查值是否可以被序列化
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return fmt.Errorf("value of type %s cannot be cached", val.Kind())
	}

	return nil
}

// validateTTL 验证TTL的有效性
func validateTTL(ttl time.Duration) error {
	if ttl < 0 {
		return fmt.Errorf("TTL cannot be negative: %v", ttl)
	}

	// 设置一个合理的上限（比如30天）
	maxTTL := 30 * 24 * time.Hour
	if ttl > maxTTL {
		return fmt.Errorf("TTL too large: max %v, got %v", maxTTL, ttl)
	}

	return nil
}

// validateMaxSize 验证最大缓存大小的有效性
func validateMaxSize(size int) error {
	if size <= 0 {
		return fmt.Errorf("max size must be positive, got %d", size)
	}

	// 设置一个合理的上限（比如1亿个项目）
	maxSize := 100000000
	if size > maxSize {
		return fmt.Errorf("max size too large: max %d, got %d", maxSize, size)
	}

	return nil
}

// validateShards 验证分片数量的有效性
func validateShards(shards int) error {
	if shards <= 0 {
		return fmt.Errorf("shards must be positive, got %d", shards)
	}

	// 分片数量应该是2的幂次，这样哈希分布更均匀
	if shards&(shards-1) != 0 {
		return fmt.Errorf("shards should be a power of 2 for optimal performance, got %d", shards)
	}

	// 设置合理的上限
	maxShards := 1024
	if shards > maxShards {
		return fmt.Errorf("too many shards: max %d, got %d", maxShards, shards)
	}

	return nil
}

// validateEvictionPolicy 验证淘汰策略的有效性
func validateEvictionPolicy(policy string) error {
	if policy == "" {
		return errors.New("eviction policy cannot be empty")
	}

	validPolicies := []string{
		constants.LRUStrategy,
		constants.LFUStrategy,
		constants.FIFOStrategy,
	}

	for _, valid := range validPolicies {
		if policy == valid {
			return nil
		}
	}

	return fmt.Errorf("%s: valid options are %v",
		constants.ErrInvalidStrategy, validPolicies)
}

// validateSerializer 验证序列化器的有效性
func validateSerializer(serializer string) error {
	if serializer == "" {
		return nil // 空值表示使用默认序列化器
	}

	validSerializers := []string{
		constants.JSONEncoding,
		constants.GobEncoding,
	}

	for _, valid := range validSerializers {
		if serializer == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid serializer: %q, valid options are %v",
		serializer, validSerializers)
}

// validateCleanupInterval 验证清理间隔的有效性
func validateCleanupInterval(interval time.Duration) error {
	if interval <= 0 {
		return fmt.Errorf("cleanup interval must be positive, got %v", interval)
	}

	// 设置合理的上限（比如1小时）
	maxInterval := time.Hour
	if interval > maxInterval {
		return fmt.Errorf("cleanup interval too large: max %v, got %v", maxInterval, interval)
	}

	return nil
}

// ComprehensiveValidator 全面的配置验证器
type ComprehensiveValidator struct {
	enableStrictMode bool
}

// NewComprehensiveValidator 创建新的全面验证器
func NewComprehensiveValidator(enableStrictMode bool) *ComprehensiveValidator {
	return &ComprehensiveValidator{
		enableStrictMode: enableStrictMode,
	}
}

// ValidateConfig 全面验证配置
func (v *ComprehensiveValidator) ValidateConfig(config *Config) error {
	var errors ValidationErrors

	// 验证各个字段
	if err := validateMaxSize(config.MaxSize); err != nil {
		errors = append(errors, fmt.Errorf("MaxSize: %w", err))
	}

	if err := validateShards(config.Shards); err != nil {
		errors = append(errors, fmt.Errorf("Shards: %w", err))
	}

	if err := validateTTL(config.DefaultTTL); err != nil {
		errors = append(errors, fmt.Errorf("DefaultTTL: %w", err))
	}

	if err := validateCleanupInterval(config.CleanupInterval); err != nil {
		errors = append(errors, fmt.Errorf("CleanupInterval: %w", err))
	}

	if err := validateEvictionPolicy(config.EvictionPolicy); err != nil {
		errors = append(errors, fmt.Errorf("EvictionPolicy: %w", err))
	}

	if err := validateSerializer(config.Serializer); err != nil {
		errors = append(errors, fmt.Errorf("Serializer: %w", err))
	}

	// 验证配置之间的关系
	if config.MaxSize < config.Shards {
		errors = append(errors, fmt.Errorf("MaxSize (%d) should be >= Shards (%d)",
			config.MaxSize, config.Shards))
	}

	// 严格模式下的额外检查
	if v.enableStrictMode {
		if config.MaxSize < 100 {
			errors = append(errors, fmt.Errorf("MaxSize (%d) is too small for production use", config.MaxSize))
		}

		if config.Shards > config.MaxSize/10 {
			errors = append(errors, fmt.Errorf("too many shards (%d) for cache size (%d)",
				config.Shards, config.MaxSize))
		}

		if config.CleanupInterval > 30*time.Minute {
			errors = append(errors, fmt.Errorf("cleanup interval (%v) is too long", config.CleanupInterval))
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// ValidateOperation 验证缓存操作参数
func (v *ComprehensiveValidator) ValidateOperation(operation string, key string, value interface{}) error {
	var errors ValidationErrors

	if operation == "" {
		errors = append(errors, fmt.Errorf("operation cannot be empty"))
	}

	if err := validateKey(key); err != nil {
		errors = append(errors, fmt.Errorf("key validation failed: %w", err))
	}

	if value != nil {
		if err := validateValue(value); err != nil {
			errors = append(errors, fmt.Errorf("value validation failed: %w", err))
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// ValidateBatchOperation 验证批量操作参数
func (v *ComprehensiveValidator) ValidateBatchOperation(items map[string]interface{}) error {
	var errors ValidationErrors

	if len(items) == 0 {
		errors = append(errors, fmt.Errorf("batch operation requires at least one item"))
	}

	if len(items) > 10000 {
		errors = append(errors, fmt.Errorf("batch operation too large: max 10000 items, got %d", len(items)))
	}
	validatedCount := 0

	for key, value := range items {
		if err := validateKey(key); err != nil {
			errors = append(errors, fmt.Errorf("key %q: %w", key, err))
			continue
		}

		if err := validateValue(value); err != nil {
			errors = append(errors, fmt.Errorf("value for key %q: %w", key, err))
			continue
		}

		validatedCount++
	}

	if validatedCount == 0 {
		errors = append(errors, fmt.Errorf("no valid items in batch operation"))
	}

	if len(errors) > 0 && v.enableStrictMode {
		return errors
	}

	return nil
}

// GetValidationSummary 获取验证摘要
func (v *ComprehensiveValidator) GetValidationSummary(err error) map[string]interface{} {
	summary := map[string]interface{}{
		"valid": err == nil,
	}

		if err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			summary["error_count"] = len(ve)
			summary["errors"] = ve.Error()
		} else {
			summary["error"] = err.Error()
		}
	}

	return summary
}