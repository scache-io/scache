package scache

import (
	"fmt"
	"testing"
	"time"
)

func TestGlobalCacheBasicOperations(t *testing.T) {
	// 测试全局SET操作
	err := Set("test_key", "test_value", time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 测试全局GET操作
	value, found, err := Get("test_key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !found {
		t.Fatal("Key should exist")
	}
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %v", value)
	}

	// 测试全局EXISTS操作
	exists, err := Exists("test_key")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Key should exist")
	}

	// 测试全局DELETE操作
	deleted, err := Delete("test_key")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !deleted {
		t.Error("Delete should succeed")
	}

	// 验证删除后不存在
	_, exists, err = Get("test_key")
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if exists {
		t.Error("Key should not exist after delete")
	}
}

func TestGlobalCacheListOperations(t *testing.T) {
	// 测试全局LPUSH操作
	length, err := LPush("test_list", "item1", time.Minute)
	if err != nil {
		t.Fatalf("LPush failed: %v", err)
	}
	if length != 1 {
		t.Errorf("Expected length 1, got %d", length)
	}

	// 测试全局RPOP操作
	value, err := RPop("test_list")
	if err != nil {
		t.Fatalf("RPop failed: %v", err)
	}
	if value != "item1" {
		t.Errorf("Expected 'item1', got %v", value)
	}

	// 测试空列表弹出
	value, err = RPop("test_list")
	if err != nil {
		t.Fatalf("RPop on empty list failed: %v", err)
	}
	if value != nil {
		t.Error("RPop on empty list should return nil")
	}
}

func TestGlobalCacheHashOperations(t *testing.T) {
	// 测试全局HSET操作
	result, err := HSet("test_hash", "field1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("HSet failed: %v", err)
	}
	if result != 1 {
		t.Errorf("HSet should return 1 for new field, got %d", result)
	}

	// 测试全局HGET操作
	value, err := HGet("test_hash", "field1")
	if err != nil {
		t.Fatalf("HGet failed: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	// 测试获取不存在的字段
	value, err = HGet("test_hash", "nonexistent")
	if err != nil {
		t.Fatalf("HGet nonexistent field failed: %v", err)
	}
	if value != nil {
		t.Error("HGet nonexistent field should return nil")
	}
}

func TestGlobalCacheTTL(t *testing.T) {
	// 测试TTL操作
	ttl, err := TTL("nonexistent_key")
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}
	if ttl != -2 {
		t.Errorf("Expected TTL -2 for nonexistent key, got %d", ttl)
	}

	// 设置一个键（永不过期）
	err = Set("ttl_key", "ttl_value")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 检查TTL（永不过期）
	ttl, err = TTL("ttl_key")
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}
	if ttl != -1 {
		t.Errorf("Expected TTL -1 for non-expiring key, got %d", ttl)
	}

	// 设置过期时间
	success, err := Expire("ttl_key", time.Minute*30)
	if err != nil {
		t.Fatalf("Expire failed: %v", err)
	}
	if !success {
		t.Error("Expire should succeed")
	}

	// 检查新的TTL（30分钟过期）
	ttl, err = TTL("ttl_key")
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}
	if ttl <= 0 || ttl > 1800 { // 30分钟 = 1800秒
		t.Errorf("Expected TTL between 1 and 1800, got %d", ttl)
	}
}

func TestGlobalCacheType(t *testing.T) {
	// 测试字符串类型
	err := Set("string_key", "string_value", time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	keyType, err := Type("string_key")
	if err != nil {
		t.Fatalf("Type failed: %v", err)
	}
	if keyType != "string" {
		t.Errorf("Expected type 'string', got '%s'", keyType)
	}

	// 测试不存在的键类型
	keyType, err = Type("nonexistent_key")
	if err != nil {
		t.Fatalf("Type failed: %v", err)
	}
	if keyType != "none" {
		t.Errorf("Expected type 'none' for nonexistent key, got '%s'", keyType)
	}
}

func TestGlobalCacheStats(t *testing.T) {
	// 执行一些操作来生成统计信息
	Set("stats_key1", "value1", time.Minute)
	Set("stats_key2", "value2", time.Minute)
	Get("stats_key1")
	Get("nonexistent")
	Delete("stats_key1")

	// 获取统计信息
	stats := Stats()
	if stats == nil {
		t.Fatal("Stats should not be nil")
	}

	// 验证统计结构
	statsMap, ok := stats.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", stats)
	}

	// 检查基本统计字段
	expectedFields := []string{"hits", "misses", "sets", "deletes", "keys"}
	for _, field := range expectedFields {
		if _, exists := statsMap[field]; !exists {
			t.Errorf("Missing stats field: %s", field)
		}
	}

	t.Logf("Global cache stats: %+v", stats)
}

func TestGlobalCacheCommands(t *testing.T) {
	// 测试直接执行命令
	result, err := Execute("SET", "cmd_key", "cmd_value", time.Minute)
	if err != nil {
		t.Fatalf("Execute SET failed: %v", err)
	}
	if result != nil {
		t.Errorf("SET should return nil, got %v", result)
	}

	result, err = Execute("GET", "cmd_key")
	if err != nil {
		t.Fatalf("Execute GET failed: %v", err)
	}
	if result != "cmd_value" {
		t.Errorf("GET should return 'cmd_value', got %v", result)
	}

	result, err = Execute("EXISTS", "cmd_key")
	if err != nil {
		t.Fatalf("Execute EXISTS failed: %v", err)
	}
	if !result.(bool) {
		t.Error("EXISTS should return true")
	}

	result, err = Execute("TYPE", "cmd_key")
	if err != nil {
		t.Fatalf("Execute TYPE failed: %v", err)
	}
	if result.(string) != "string" {
		t.Errorf("TYPE should return 'string', got %v", result)
	}
}

func TestGlobalCacheConcurrentAccess(t *testing.T) {
	// 并发测试全局缓存的线程安全性
	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// 每个协程执行多个操作
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("concurrent_key_%d_%d", id, j)
				value := fmt.Sprintf("concurrent_value_%d_%d", id, j)

				// SET
				err := Set(key, value, time.Minute)
				if err != nil {
					t.Errorf("Goroutine %d SET failed: %v", id, err)
					return
				}

				// GET
				_, found, err := Get(key)
				if err != nil {
					t.Errorf("Goroutine %d GET failed: %v", id, err)
					return
				}
				if !found {
					t.Errorf("Goroutine %d key should exist", id)
					return
				}

				// DELETE
				_, err = Delete(key)
				if err != nil {
					t.Errorf("Goroutine %d DELETE failed: %v", id, err)
					return
				}
			}
		}(i)
	}

	// 等待所有协程完成
	for i := 0; i < concurrency; i++ {
		<-done
	}

	// 验证最终统计信息
	stats := Stats().(map[string]interface{})
	if sets, ok := stats["sets"].(int64); ok && sets > 0 {
		t.Logf("Concurrent test completed successfully, sets: %d", sets)
	}
}

func TestGlobalCachePersistence(t *testing.T) {
	// 测试全局缓存在不同操作间的持久性
	// 设置数据
	err := Set("persistent_key", "persistent_value", time.Hour)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 第一次获取
	value1, found1, err := Get("persistent_key")
	if err != nil {
		t.Fatalf("First Get failed: %v", err)
	}
	if !found1 || value1 != "persistent_value" {
		t.Fatal("First GET should return persistent_value")
	}

	// 修改数据
	err = Set("persistent_key", "updated_value", time.Hour)
	if err != nil {
		t.Fatalf("Update Set failed: %v", err)
	}

	// 第二次获取
	value2, found2, err := Get("persistent_key")
	if err != nil {
		t.Fatalf("Second Get failed: %v", err)
	}
	if !found2 || value2 != "updated_value" {
		t.Fatal("Second GET should return updated_value")
	}

	// 验证统计信息更新
	stats := Stats().(map[string]interface{})
	if sets, ok := stats["sets"].(int64); ok && sets >= 2 {
		t.Logf("Global cache persistence verified, sets: %d", sets)
	}
}
