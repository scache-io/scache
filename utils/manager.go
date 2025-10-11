package utils

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/scache/constants"
	"github.com/scache/types"
)

// ConstantsManager 常量管理器
type ConstantsManager struct{}

// NewConstantsManager 创建常量管理器
func NewConstantsManager() *ConstantsManager {
	return &ConstantsManager{}
}

// ListAllConstants 列出所有常量
func (cm *ConstantsManager) ListAllConstants() map[string]interface{} {
	return map[string]interface{}{
		"缓存配置常量": map[string]interface{}{
			"DefaultMaxSize":         constants.DefaultMaxSize,
			"DefaultShards":          constants.DefaultShards,
			"DefaultCleanupInterval": constants.DefaultCleanupInterval,
			"DefaultTTL":             constants.DefaultTTL,
		},
		"缓存策略常量": map[string]interface{}{
			"LRUStrategy":  constants.LRUStrategy,
			"LFUStrategy":  constants.LFUStrategy,
			"FIFOStrategy": constants.FIFOStrategy,
		},
		"全局缓存常量": map[string]interface{}{
			"DefaultCacheName": constants.DefaultCacheName,
			"ManagerTimeout":   constants.ManagerTimeout,
		},
		"性能相关常量": map[string]interface{}{
			"MaxKeyLength": constants.MaxKeyLength,
			"MinKeyLength": constants.MinKeyLength,
			"MaxValueSize": constants.MaxValueSize,
		},
		"错误消息常量": map[string]interface{}{
			"ErrCacheNotFound":      constants.ErrCacheNotFound,
			"ErrInvalidCacheName":   constants.ErrInvalidCacheName,
			"ErrCacheAlreadyExists": constants.ErrCacheAlreadyExists,
			"ErrInvalidStrategy":    constants.ErrInvalidStrategy,
			"ErrKeyNotFound":        constants.ErrKeyNotFound,
			"ErrKeyTooLong":         constants.ErrKeyTooLong,
			"ErrKeyEmpty":           constants.ErrKeyEmpty,
			"ErrValueTooLarge":      constants.ErrValueTooLarge,
			"ErrCacheClosed":        constants.ErrCacheClosed,
		},
		"日志相关常量": map[string]interface{}{
			"LogPrefixCache":   constants.LogPrefixCache,
			"LogPrefixManager": constants.LogPrefixManager,
			"LogPrefixGlobal":  constants.LogPrefixGlobal,
		},
		"统计相关常量": map[string]interface{}{
			"StatsUpdateInterval": constants.StatsUpdateInterval,
			"HitRateThreshold":    constants.HitRateThreshold,
		},
		"序列化常量": map[string]interface{}{
			"JSONEncoding": constants.JSONEncoding,
			"GobEncoding":  constants.GobEncoding,
		},
	}
}

// PrintConstants 打印所有常量（按类别）
func (cm *ConstantsManager) PrintConstants() {
	constants := cm.ListAllConstants()

	fmt.Println("=== SCache 常量总览 ===")
	for category, values := range constants {
		fmt.Printf("\n%s:\n", category)
		if m, ok := values.(map[string]interface{}); ok {
			// 按键名排序
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, key := range keys {
				fmt.Printf("  %s: %v\n", key, m[key])
			}
		}
	}
	fmt.Println()
}

// TypesManager 类型管理器
type TypesManager struct{}

// NewTypesManager 创建类型管理器
func NewTypesManager() *TypesManager {
	return &TypesManager{}
}

// ListAllTypes 列出所有重要类型
func (tm *TypesManager) ListAllTypes() map[string]string {
	return map[string]string{
		"CacheItem":           reflect.TypeOf(types.CacheItem{}).String(),
		"CacheStats":          reflect.TypeOf(types.CacheStats{}).String(),
		"CacheConfig":         reflect.TypeOf(types.CacheConfig{}).String(),
		"CacheShard":          reflect.TypeOf(types.CacheShard{}).String(),
		"ShardStats":          reflect.TypeOf(types.ShardStats{}).String(),
		"ManagerStats":        reflect.TypeOf(types.ManagerStats{}).String(),
		"ManagerConfig":       reflect.TypeOf(types.ManagerConfig{}).String(),
		"PolicyStats":         reflect.TypeOf(types.PolicyStats{}).String(),
		"LRUEntry":            reflect.TypeOf(types.LRUEntry{}).String(),
		"LFUEntry":            reflect.TypeOf(types.LFUEntry{}).String(),
		"FIFOEntry":           reflect.TypeOf(types.FIFOEntry{}).String(),
		"SerializationConfig": reflect.TypeOf(types.SerializationConfig{}).String(),
		"SerializedData":      reflect.TypeOf(types.SerializedData{}).String(),
		"ErrorInfo":           reflect.TypeOf(types.ErrorInfo{}).String(),
		"ErrorStats":          reflect.TypeOf(types.ErrorStats{}).String(),
		"Metrics":             reflect.TypeOf(types.Metrics{}).String(),
		"CacheMetrics":        reflect.TypeOf(types.CacheMetrics{}).String(),
		"SystemMetrics":       reflect.TypeOf(types.SystemMetrics{}).String(),
		"CacheEvent":          reflect.TypeOf(types.CacheEvent{}).String(),
		"EventConfig":         reflect.TypeOf(types.EventConfig{}).String(),
		"HealthStatus":        reflect.TypeOf(types.HealthStatus{}).String(),
		"HealthCheck":         reflect.TypeOf(types.HealthCheck{}).String(),
		"HealthSummary":       reflect.TypeOf(types.HealthSummary{}).String(),
		"ValidationResult":    reflect.TypeOf(types.ValidationResult{}).String(),
		"ValidationError":     reflect.TypeOf(types.ValidationError{}).String(),
		"ValidationWarning":   reflect.TypeOf(types.ValidationWarning{}).String(),
	}
}

// PrintTypes 打印所有类型信息
func (tm *TypesManager) PrintTypes() {
	types := tm.ListAllTypes()

	fmt.Println("=== SCache 类型总览 ===")

	// 按类别分组
	categories := map[string][]string{
		"缓存核心类型": {"CacheItem", "CacheStats", "CacheConfig", "CacheShard", "ShardStats"},
		"管理器类型":  {"ManagerStats", "ManagerConfig"},
		"策略类型":   {"PolicyStats", "LRUEntry", "LFUEntry", "FIFOEntry"},
		"序列化类型":  {"SerializationConfig", "SerializedData"},
		"错误处理类型": {"ErrorInfo", "ErrorStats"},
		"监控指标类型": {"Metrics", "CacheMetrics", "SystemMetrics"},
		"事件通知类型": {"CacheEvent", "EventConfig"},
		"健康检查类型": {"HealthStatus", "HealthCheck", "HealthSummary"},
		"配置验证类型": {"ValidationResult", "ValidationError", "ValidationWarning"},
	}

	for category, typeNames := range categories {
		fmt.Printf("\n%s:\n", category)
		for _, typeName := range typeNames {
			if typeDef, exists := types[typeName]; exists {
				fmt.Printf("  %s: %s\n", typeName, typeDef)
			}
		}
	}
	fmt.Println()
}

// ConfigValidator 配置验证器
type ConfigValidator struct{}

// NewConfigValidator 创建配置验证器
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// ValidateConfig 验证配置的完整性
func (cv *ConfigValidator) ValidateConfig(config interface{}) *types.ValidationResult {
	result := &types.ValidationResult{
		Valid:    true,
		Errors:   []types.ValidationError{},
		Warnings: []types.ValidationWarning{},
	}

	// 使用反射检查配置
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		result.Errors = append(result.Errors, types.ValidationError{
			Field:   "config",
			Message: "配置必须是结构体或结构体指针",
			Value:   fmt.Sprintf("%T", config),
		})
		result.Valid = false
		return result
	}

	t := v.Type()

	// 检查必需字段
	requiredFields := map[string]string{
		"MaxSize":              "int",
		"Shards":               "int",
		"DefaultTTL":           "time.Duration",
		"CleanupInterval":      "time.Duration",
		"EvictionPolicy":       "string",
		"EnableStatistics":     "bool",
		"EnableLazyExpiration": "bool",
		"EnableMetrics":        "bool",
	}

	for fieldName, expectedType := range requiredFields {
		field, found := t.FieldByName(fieldName)
		if !found {
			result.Errors = append(result.Errors, types.ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("缺少必需字段: %s", fieldName),
				Value:   "N/A",
			})
			result.Valid = false
			continue
		}

		actualType := field.Type.String()
		if !strings.Contains(actualType, expectedType) {
			result.Warnings = append(result.Warnings, types.ValidationWarning{
				Field:   fieldName,
				Message: fmt.Sprintf("字段类型不匹配，期望: %s, 实际: %s", expectedType, actualType),
				Value:   actualType,
			})
		}
	}

	return result
}

// PrintValidationResult 打印验证结果
func (cv *ConfigValidator) PrintValidationResult(result *types.ValidationResult) {
	fmt.Println("=== 配置验证结果 ===")
	fmt.Printf("验证状态: %s\n", map[bool]string{true: "✓ 通过", false: "✗ 失败"}[result.Valid])

	if len(result.Errors) > 0 {
		fmt.Printf("\n错误 (%d):\n", len(result.Errors))
		for i, err := range result.Errors {
			fmt.Printf("  %d. %s: %s\n", i+1, err.Field, err.Message)
			if err.Value != "" {
				fmt.Printf("     值: %v\n", err.Value)
			}
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("\n警告 (%d):\n", len(result.Warnings))
		for i, warn := range result.Warnings {
			fmt.Printf("  %d. %s: %s\n", i+1, warn.Field, warn.Message)
			if warn.Value != "" {
				fmt.Printf("     值: %v\n", warn.Value)
			}
		}
	}

	fmt.Println()
}

// ProjectInfo 项目信息
type ProjectInfo struct{}

// NewProjectInfo 创建项目信息
func NewProjectInfo() *ProjectInfo {
	return &ProjectInfo{}
}

// GetProjectStructure 获取项目结构信息
func (pi *ProjectInfo) GetProjectStructure() map[string][]string {
	return map[string][]string{
		"核心包": {
			"cache/      - 缓存核心实现",
			"interfaces/ - 接口定义",
			"constants/  - 常量定义",
			"types/      - 类型定义",
			"globals/    - 全局变量管理",
		},
		"策略包": {
			"policies/lru/  - LRU淘汰策略",
			"policies/lfu/  - LFU淘汰策略",
			"policies/fifo/ - FIFO淘汰策略",
		},
		"工具包": {
			"utils/       - 工具函数",
			"examples/    - 使用示例",
		},
		"入口文件": {
			"scache.go    - 主入口，API重新导出",
		},
	}
}

// PrintProjectStructure 打印项目结构
func (pi *ProjectInfo) PrintProjectStructure() {
	structure := pi.GetProjectStructure()

	fmt.Println("=== SCache 项目结构 ===")
	for category, items := range structure {
		fmt.Printf("\n%s:\n", category)
		for _, item := range items {
			fmt.Printf("  %s\n", item)
		}
	}
	fmt.Println()
}

// GetStatistics 获取项目统计信息
func (pi *ProjectInfo) GetStatistics() map[string]interface{} {
	return map[string]interface{}{
		"总包数":   8,
		"核心接口数": 6,
		"核心类型数": 25,
		"常量数":   30,
		"策略数":   3,
		"示例数":   4,
		"工具数":   4,
		"特性": []string{
			"统一常量管理",
			"类型定义组织",
			"全局变量管理",
			"配置验证",
			"错误处理",
			"性能监控",
			"健康检查",
			"事件通知",
		},
	}
}

// PrintStatistics 打印项目统计信息
func (pi *ProjectInfo) PrintStatistics() {
	stats := pi.GetStatistics()

	fmt.Println("=== SCache 项目统计 ===")
	for key, value := range stats {
		switch v := value.(type) {
		case []string:
			fmt.Printf("%s: %d\n", key, len(v))
			fmt.Println("详细列表:")
			for i, item := range v {
				fmt.Printf("  %d. %s\n", i+1, item)
			}
		default:
			fmt.Printf("%s: %v\n", key, v)
		}
	}
	fmt.Println()
}
