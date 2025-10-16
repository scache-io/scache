//go:build integration
// +build integration

package tests

import (
	"fmt"
	"testing"
	"time"

	"scache"
	"scache/config"
	"scache/interfaces"
	"scache/storage"
	"scache/types"
)

// TestBasicIntegration 基础集成测试
func TestBasicIntegration(t *testing.T) {
	// 创建存储引擎
	engine := storage.NewStorageEngine(nil)

	// 创建执行器
	executor := scache.NewExecutor(engine)

	// 测试 SET 和 GET
	result, err := executor.Execute("SET", "test_key", "test_value", time.Minute)
	if err != nil {
		t.Fatalf("SET failed: %v", err)
	}
	if result != nil {
		t.Errorf("SET should return nil, got %v", result)
	}

	result, err = executor.Execute("GET", "test_key")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got %v", result)
	}

	// 测试 EXISTS
	result, err = executor.Execute("EXISTS", "test_key")
	if err != nil {
		t.Fatalf("EXISTS failed: %v", err)
	}
	if result != true {
		t.Errorf("Expected true, got %v", result)
	}

	// 测试 TYPE
	result, err = executor.Execute("TYPE", "test_key")
	if err != nil {
		t.Fatalf("TYPE failed: %v", err)
	}
	if result != "string" {
		t.Errorf("Expected 'string', got %v", result)
	}
}

// TestListIntegration 列表操作集成测试
func TestListIntegration(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	// 测试 LPUSH
	result, err := executor.Execute("LPUSH", "mylist", "item1", time.Minute)
	if err != nil {
		t.Fatalf("LPUSH failed: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected length 1, got %v", result)
	}

	result, err = executor.Execute("LPUSH", "mylist", "item2", time.Minute)
	if err != nil {
		t.Fatalf("LPUSH failed: %v", err)
	}
	if result != 2 {
		t.Errorf("Expected length 2, got %v", result)
	}

	// 测试 RPOP
	result, err = executor.Execute("RPOP", "mylist")
	if err != nil {
		t.Fatalf("RPOP failed: %v", err)
	}
	if result != "item1" {
		t.Errorf("Expected 'item1', got %v", result)
	}

	// 再次 RPOP
	result, err = executor.Execute("RPOP", "mylist")
	if err != nil {
		t.Fatalf("RPOP failed: %v", err)
	}
	if result != "item2" {
		t.Errorf("Expected 'item2', got %v", result)
	}

	// 测试空列表
	result, err = executor.Execute("RPOP", "mylist")
	if err != nil {
		t.Fatalf("RPOP failed: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil for empty list, got %v", result)
	}
}

// TestHashIntegration 哈希操作集成测试
func TestHashIntegration(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	// 测试 HSET
	result, err := executor.Execute("HSET", "myhash", "field1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("HSET failed: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1, got %v", result)
	}

	result, err = executor.Execute("HSET", "myhash", "field2", "value2", time.Minute)
	if err != nil {
		t.Fatalf("HSET failed: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1, got %v", result)
	}

	// 测试 HGET
	result, err = executor.Execute("HGET", "myhash", "field1")
	if err != nil {
		t.Fatalf("HGET failed: %v", err)
	}
	if result != "value1" {
		t.Errorf("Expected 'value1', got %v", result)
	}

	// 测试不存在的字段
	result, err = executor.Execute("HGET", "myhash", "nonexistent")
	if err != nil {
		t.Fatalf("HGET failed: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil for nonexistent field, got %v", result)
	}
}

// TestExpirationIntegration 过期时间集成测试
func TestExpirationIntegration(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	// 设置一个短期过期的键
	_, err := executor.Execute("SET", "expire_key", "expire_value", time.Millisecond*100)
	if err != nil {
		t.Fatalf("SET failed: %v", err)
	}

	// 立即获取应该成功
	result, err := executor.Execute("GET", "expire_key")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if result != "expire_value" {
		t.Errorf("Expected 'expire_value', got %v", result)
	}

	// 等待过期
	time.Sleep(time.Millisecond * 150)

	// 再次获取应该失败
	result, err = executor.Execute("GET", "expire_key")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil for expired key, got %v", result)
	}
}

// TestTTLIntegration TTL操作集成测试
func TestTTLIntegration(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	// 设置一个永不过期的键
	_, err := executor.Execute("SET", "forever_key", "forever_value", 0)
	if err != nil {
		t.Fatalf("SET failed: %v", err)
	}

	// 检查TTL
	result, err := executor.Execute("TTL", "forever_key")
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}
	if result != -1 {
		t.Errorf("Expected -1 for non-expiring key, got %v", result)
	}

	// 设置过期时间
	result, err = executor.Execute("EXPIRE", "forever_key", time.Second*30)
	if err != nil {
		t.Fatalf("EXPIRE failed: %v", err)
	}
	if result != true {
		t.Errorf("Expected true, got %v", result)
	}

	// 再次检查TTL
	result, err = executor.Execute("TTL", "forever_key")
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}
	ttl := result.(int)
	if ttl <= 0 || ttl > 30 {
		t.Errorf("Expected TTL between 1 and 30, got %d", ttl)
	}

	// 测试不存在的键
	result, err = executor.Execute("TTL", "nonexistent_key")
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}
	if result != -2 {
		t.Errorf("Expected -2 for nonexistent key, got %v", result)
	}
}

// TestStatsIntegration 统计信息集成测试
func TestStatsIntegration(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	// 执行一些操作
	_, err := executor.Execute("SET", "stat_key", "stat_value", time.Minute)
	if err != nil {
		t.Fatalf("SET failed: %v", err)
	}

	_, err = executor.Execute("GET", "stat_key")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	_, err = executor.Execute("GET", "nonexistent")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	_, err = executor.Execute("DEL", "stat_key")
	if err != nil {
		t.Fatalf("DEL failed: %v", err)
	}

	// 获取统计信息
	result, err := executor.Execute("STATS")
	if err != nil {
		t.Fatalf("STATS failed: %v", err)
	}

	stats, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// 验证统计字段
	expectedFields := []string{"hits", "misses", "sets", "deletes", "keys"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Missing stats field: %s", field)
		}
	}

	// 验证基本统计
	if sets, ok := stats["sets"].(int64); !ok || sets != 1 {
		t.Errorf("Expected sets=1, got %v", stats["sets"])
	}

	if hits, ok := stats["hits"].(int64); !ok || hits != 1 {
		t.Errorf("Expected hits=1, got %v", stats["hits"])
	}

	if misses, ok := stats["misses"].(int64); !ok || misses != 1 {
		t.Errorf("Expected misses=1, got %v", stats["misses"])
	}

	if deletes, ok := stats["deletes"].(int64); !ok || deletes != 1 {
		t.Errorf("Expected deletes=1, got %v", stats["deletes"])
	}
}

// TestConcurrentIntegration 并发集成测试
func TestConcurrentIntegration(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	// 并发写入
	concurrency := 100
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			key := fmt.Sprintf("concurrent_key_%d", id)
			value := fmt.Sprintf("concurrent_value_%d", id)

			_, err := executor.Execute("SET", key, value, time.Minute)
			if err != nil {
				t.Errorf("SET failed in goroutine %d: %v", id, err)
			}

			// 读取验证
			result, err := executor.Execute("GET", key)
			if err != nil {
				t.Errorf("GET failed in goroutine %d: %v", id, err)
			}
			if result != value {
				t.Errorf("Value mismatch in goroutine %d: expected %s, got %v", id, value, result)
			}

			done <- true
		}(i)
	}

	// 等待所有协程完成
	for i := 0; i < concurrency; i++ {
		<-done
	}

	// 验证所有键都存在
	keys, err := executor.Execute("KEYS")
	if err != nil {
		t.Fatalf("KEYS failed: %v", err)
	}

	keyList, ok := keys.([]string)
	if !ok {
		t.Fatalf("Expected []string, got %T", keys)
	}

	if len(keyList) != concurrency {
		t.Errorf("Expected %d keys, got %d", concurrency, len(keyList))
	}
}

// TestConfigurationIntegration 配置集成测试
func TestConfigurationIntegration(t *testing.T) {
	// 创建自定义配置的引擎
	config := &storage.EngineConfig{
		MaxSize:                   10,
		DefaultExpiration:         time.Minute,
		MemoryThreshold:           0.8,
		BackgroundCleanupInterval: time.Second,
	}

	engine := storage.NewStorageEngine(config)
	executor := scache.NewExecutor(engine)

	// 添加超过最大容量的键
	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("overflow_key_%d", i)
		_, err := executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Minute)
		if err != nil {
			t.Fatalf("SET failed for key %s: %v", key, err)
		}
	}

	// 验证引擎大小不超过最大值
	size := engine.Size()
	if size > 10 {
		t.Errorf("Engine size %d exceeds MaxSize 10", size)
	}

	// 验证统计信息中有淘汰记录
	result, err := executor.Execute("STATS")
	if err != nil {
		t.Fatalf("STATS failed: %v", err)
	}

	stats := result.(map[string]interface{})
	if evictions, ok := stats["evictions"].(int64); !ok || evictions <= 0 {
		t.Errorf("Expected evictions > 0, got %v", stats["evictions"])
	}
}

// TestConvenienceAPIIntegration 便捷API集成测试
func TestConvenienceAPIIntegration(t *testing.T) {
	// 测试便捷API函数
	err := scache.Set("conv_key", "conv_value", time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	value, found, err := scache.Get("conv_key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !found {
		t.Error("Key should exist")
	}
	if value != "conv_value" {
		t.Errorf("Expected 'conv_value', got %v", value)
	}

	// 测试列表操作
	length, err := scache.LPush("conv_list", "conv_item", time.Minute)
	if err != nil {
		t.Fatalf("LPush failed: %v", err)
	}
	if length != 1 {
		t.Errorf("Expected length 1, got %d", length)
	}

	// 测试哈希操作
	success, err := scache.HSet("conv_hash", "conv_field", "conv_field_value", time.Minute)
	if err != nil {
		t.Fatalf("HSet failed: %v", err)
	}
	if !success {
		t.Error("HSet should succeed")
	}

	fieldValue, err := scache.HGet("conv_hash", "conv_field")
	if err != nil {
		t.Fatalf("HGet failed: %v", err)
	}
	if fieldValue != "conv_field_value" {
		t.Errorf("Expected 'conv_field_value', got %v", fieldValue)
	}

	// 测试其他操作
	exists, err := scache.Exists("conv_key")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Key should exist")
	}

	keyType, err := scache.Type("conv_key")
	if err != nil {
		t.Fatalf("Type failed: %v", err)
	}
	if keyType != "string" {
		t.Errorf("Expected 'string', got %s", keyType)
	}

	// 测试统计信息
	stats := scache.Stats()
	if stats == nil {
		t.Error("Stats should not be nil")
	}
}
