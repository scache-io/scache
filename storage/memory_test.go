package storage

import (
	"fmt"
	"testing"

	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/constants"
	"github.com/scache-io/scache/internal"
	"github.com/scache-io/scache/types"
)

func TestMemoryManagement(t *testing.T) {
	// 创建一个禁用自动清理的配置（默认行为）
	testConfig := &config.EngineConfig{
		MaxSize:                   constants.TestCapacity,   // 小容量
		MemoryThreshold:           0.1, // 低内存阈值，容易触发
		DefaultExpiration:         0,
		BackgroundCleanupInterval: 0, // 禁用自动清理
	}

	engine := NewStorageEngine(testConfig)

	// 启用内存检查
	internal.EnableMemoryCheck()

	// 添加一些数据，直到达到容量限制
	for i := 0; i < testConfig.MaxSize; i++ {
		key := fmt.Sprintf("key_%d", i)
		obj := types.NewStringObject("value", 0)
		err := engine.Set(key, obj)
		if err != nil {
			t.Fatalf("Expected no error for first 10 items, got: %v", err)
		}
	}

	// 尝试添加第11个项目（应该因为容量限制失败）
	obj := types.NewStringObject("overflow", 0)
	err := engine.Set("overflow", obj)
	if err == nil {
		t.Error("Expected error when exceeding max size without cleanup")
	}

	// 验证错误消息包含容量信息
	if err != nil && len(err.Error()) == 0 {
		t.Error("Expected error message to contain capacity information")
	}

	// 清理
	_ = engine.Flush()
}

func TestMemoryThresholdCheck(t *testing.T) {
	// 创建一个极低内存阈值的配置
	testConfig := &config.EngineConfig{
		MaxSize:                   100, // 大容量
		MemoryThreshold:           constants.TestMemoryThreshold, // 极低阈值，容易触发
		DefaultExpiration:         0,
		BackgroundCleanupInterval: 0, // 禁用自动清理
	}

	engine := NewStorageEngine(testConfig)

	// 启用内存检查
	internal.EnableMemoryCheck()

	// 创建一个较大的对象
	largeValue := string(make([]byte, 1024*1024)) // 1MB字符串
	obj := types.NewStringObject(largeValue, 0)

	// 这应该因为内存阈值而失败
	err := engine.Set("large", obj)
	if err == nil {
		t.Error("Expected error when memory threshold is exceeded")
	}

	// 验证错误消息包含内存信息
	if err != nil {
		expectedError := "memory limit exceeded"
		errStr := err.Error()
		if len(errStr) == 0 || !(contains(errStr, expectedError) || contains(errStr, "insufficient memory")) {
			t.Errorf("Expected memory-related error, got: %v", err)
		}
	}

	// 清理
	_ = engine.Flush()
}

func TestMemoryCheckToggle(t *testing.T) {
	// 测试内存检查开关
	internal.DisableMemoryCheck()
	if internal.IsMemoryCheckEnabled() {
		t.Error("Expected memory check to be disabled")
	}

	internal.EnableMemoryCheck()
	if !internal.IsMemoryCheckEnabled() {
		t.Error("Expected memory check to be enabled")
	}
}

func TestBackgroundCleanupDisabled(t *testing.T) {
	// 检查默认引擎配置是否禁用了自动清理
	defaultConfig := config.DefaultEngineConfig()
	if defaultConfig.BackgroundCleanupInterval != constants.DefaultCleanupInterval {
		t.Errorf("Expected BackgroundCleanupInterval to be %v (disabled), got: %v", constants.DefaultCleanupInterval, defaultConfig.BackgroundCleanupInterval)
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
			 s[len(s)-len(substr):] == substr ||
			 indexOf(s, substr) >= 0)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}