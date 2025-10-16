package config

import (
	"testing"
	"time"

	"github.com/scache-io/scache/storage"
)

func TestEngineConfig(t *testing.T) {
	// 测试默认配置
	defaultConfig := storage.DefaultEngineConfig()

	if defaultConfig.MaxSize != 0 {
		t.Errorf("Expected MaxSize 0, got %d", defaultConfig.MaxSize)
	}

	if defaultConfig.MemoryThreshold != 0.8 {
		t.Errorf("Expected MemoryThreshold 0.8, got %f", defaultConfig.MemoryThreshold)
	}

	if defaultConfig.DefaultExpiration != 0 {
		t.Errorf("Expected DefaultExpiration 0, got %v", defaultConfig.DefaultExpiration)
	}

	if defaultConfig.BackgroundCleanupInterval != 5*time.Minute {
		t.Errorf("Expected BackgroundCleanupInterval 5m, got %v", defaultConfig.BackgroundCleanupInterval)
	}
}

func TestWithMaxSize(t *testing.T) {
	config := storage.DefaultEngineConfig()
	option := WithMaxSize(1000)
	option(config)

	if config.MaxSize != 1000 {
		t.Errorf("Expected MaxSize 1000, got %d", config.MaxSize)
	}
}

func TestWithMemoryThreshold(t *testing.T) {
	config := storage.DefaultEngineConfig()

	// 测试有效阈值
	option := WithMemoryThreshold(0.7)
	option(config)

	if config.MemoryThreshold != 0.7 {
		t.Errorf("Expected MemoryThreshold 0.7, got %f", config.MemoryThreshold)
	}

	// 测试无效阈值（应该被忽略）
	originalThreshold := config.MemoryThreshold
	option = WithMemoryThreshold(1.5) // 超过范围
	option(config)

	if config.MemoryThreshold != originalThreshold {
		t.Error("Invalid threshold should be ignored")
	}

	// 测试边界值
	option = WithMemoryThreshold(1.0)
	option(config)

	if config.MemoryThreshold != 1.0 {
		t.Error("Valid boundary value 1.0 should be accepted")
	}

	option = WithMemoryThreshold(0.0)
	option(config)

	if config.MemoryThreshold != 0.0 {
		t.Error("Valid boundary value 0.0 should be accepted")
	}
}

func TestWithDefaultExpiration(t *testing.T) {
	config := storage.DefaultEngineConfig()
	option := WithDefaultExpiration(time.Hour)
	option(config)

	if config.DefaultExpiration != time.Hour {
		t.Errorf("Expected DefaultExpiration 1h, got %v", config.DefaultExpiration)
	}
}

func TestWithBackgroundCleanup(t *testing.T) {
	config := storage.DefaultEngineConfig()
	option := WithBackgroundCleanup(time.Minute * 10)
	option(config)

	if config.BackgroundCleanupInterval != 10*time.Minute {
		t.Errorf("Expected BackgroundCleanupInterval 10m, got %v", config.BackgroundCleanupInterval)
	}
}

func TestPredefinedConfigs(t *testing.T) {
	// 测试小型配置
	smallConfig := &storage.EngineConfig{}
	for _, option := range SmallConfig {
		option(smallConfig)
	}

	if smallConfig.MaxSize != 1000 {
		t.Errorf("SmallConfig: Expected MaxSize 1000, got %d", smallConfig.MaxSize)
	}

	if smallConfig.MemoryThreshold != 0.7 {
		t.Errorf("SmallConfig: Expected MemoryThreshold 0.7, got %f", smallConfig.MemoryThreshold)
	}

	if smallConfig.DefaultExpiration != time.Hour {
		t.Errorf("SmallConfig: Expected DefaultExpiration 1h, got %v", smallConfig.DefaultExpiration)
	}

	if smallConfig.BackgroundCleanupInterval != 2*time.Minute {
		t.Errorf("SmallConfig: Expected BackgroundCleanupInterval 2m, got %v", smallConfig.BackgroundCleanupInterval)
	}

	// 测试中等配置
	mediumConfig := &storage.EngineConfig{}
	for _, option := range MediumConfig {
		option(mediumConfig)
	}

	if mediumConfig.MaxSize != 10000 {
		t.Errorf("MediumConfig: Expected MaxSize 10000, got %d", mediumConfig.MaxSize)
	}

	if mediumConfig.MemoryThreshold != 0.8 {
		t.Errorf("MediumConfig: Expected MemoryThreshold 0.8, got %f", mediumConfig.MemoryThreshold)
	}

	if mediumConfig.DefaultExpiration != 2*time.Hour {
		t.Errorf("MediumConfig: Expected DefaultExpiration 2h, got %v", mediumConfig.DefaultExpiration)
	}

	// 测试大型配置
	largeConfig := &storage.EngineConfig{}
	for _, option := range LargeConfig {
		option(largeConfig)
	}

	if largeConfig.MaxSize != 100000 {
		t.Errorf("LargeConfig: Expected MaxSize 100000, got %d", largeConfig.MaxSize)
	}

	if largeConfig.MemoryThreshold != 0.85 {
		t.Errorf("LargeConfig: Expected MemoryThreshold 0.85, got %f", largeConfig.MemoryThreshold)
	}

	if largeConfig.DefaultExpiration != 6*time.Hour {
		t.Errorf("LargeConfig: Expected DefaultExpiration 6h, got %v", largeConfig.DefaultExpiration)
	}

	if largeConfig.BackgroundCleanupInterval != 10*time.Minute {
		t.Errorf("LargeConfig: Expected BackgroundCleanupInterval 10m, got %v", largeConfig.BackgroundCleanupInterval)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := &storage.EngineConfig{}
	for _, option := range DefaultConfig {
		option(config)
	}

	if config.MaxSize != 0 {
		t.Errorf("DefaultConfig: Expected MaxSize 0, got %d", config.MaxSize)
	}

	if config.MemoryThreshold != 0.8 {
		t.Errorf("DefaultConfig: Expected MemoryThreshold 0.8, got %f", config.MemoryThreshold)
	}

	if config.DefaultExpiration != 0 {
		t.Errorf("DefaultConfig: Expected DefaultExpiration 0, got %v", config.DefaultExpiration)
	}

	if config.BackgroundCleanupInterval != 5*time.Minute {
		t.Errorf("DefaultConfig: Expected BackgroundCleanupInterval 5m, got %v", config.BackgroundCleanupInterval)
	}
}

func TestMultipleOptions(t *testing.T) {
	config := storage.DefaultEngineConfig()

	options := []EngineOption{
		WithMaxSize(5000),
		WithMemoryThreshold(0.75),
		WithDefaultExpiration(30 * time.Minute),
		WithBackgroundCleanup(time.Minute),
	}

	for _, option := range options {
		option(config)
	}

	if config.MaxSize != 5000 {
		t.Errorf("Expected MaxSize 5000, got %d", config.MaxSize)
	}

	if config.MemoryThreshold != 0.75 {
		t.Errorf("Expected MemoryThreshold 0.75, got %f", config.MemoryThreshold)
	}

	if config.DefaultExpiration != 30*time.Minute {
		t.Errorf("Expected DefaultExpiration 30m, got %v", config.DefaultExpiration)
	}

	if config.BackgroundCleanupInterval != time.Minute {
		t.Errorf("Expected BackgroundCleanupInterval 1m, got %v", config.BackgroundCleanupInterval)
	}
}

func TestConfigImmutability(t *testing.T) {
	// 确保不同的配置不会相互影响
	config1 := storage.DefaultEngineConfig()
	config2 := storage.DefaultEngineConfig()

	// 修改 config1
	WithMaxSize(1000)(config1)

	// config2 应该不受影响
	if config2.MaxSize != 0 {
		t.Error("Config2 should not be affected by changes to config1")
	}

	// 修改 config2
	WithMemoryThreshold(0.9)(config2)

	// config1 应该不受影响
	if config1.MemoryThreshold != 0.8 {
		t.Error("Config1 should not be affected by changes to config2")
	}
}

func BenchmarkWithMaxSize(b *testing.B) {
	config := storage.DefaultEngineConfig()
	option := WithMaxSize(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		option(config)
	}
}

func BenchmarkWithMemoryThreshold(b *testing.B) {
	config := storage.DefaultEngineConfig()
	option := WithMemoryThreshold(0.8)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		option(config)
	}
}

func BenchmarkDefaultEngineConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = storage.DefaultEngineConfig()
	}
}
