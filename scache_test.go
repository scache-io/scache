package main

import (
	"testing"
	"time"
)

func TestGlobalStringOperations(t *testing.T) {
	// 清空全局缓存确保测试独立
	Flush()

	// 测试全局字符串操作
	err := SetString("key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("SetString failed: %v", err)
	}

	value, found := GetString("key1")
	if !found {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %s", value)
	}
}

func TestGlobalListOperations(t *testing.T) {
	// 清空全局缓存确保测试独立
	Flush()

	// 测试全局列表操作
	err := SetList("list1", []interface{}{"item1", "item2"}, time.Minute)
	if err != nil {
		t.Fatalf("SetList failed: %v", err)
	}

	items, found := GetList("list1")
	if !found {
		t.Error("Expected list to exist")
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
	if items[0] != "item1" {
		t.Errorf("Expected 'item1' as first item, got %v", items[0])
	}
}

func TestGlobalHashOperations(t *testing.T) {
	// 清空全局缓存确保测试独立
	Flush()

	// 测试全局哈希操作
	data := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}
	err := SetHash("hash1", data, time.Minute)
	if err != nil {
		t.Fatalf("SetHash failed: %v", err)
	}

	hash, found := GetHash("hash1")
	if !found {
		t.Error("Expected hash to exist")
	}
	if hash["name"] != "Alice" {
		t.Errorf("Expected 'Alice', got %v", hash["name"])
	}
	if hash["age"] != 25 {
		t.Errorf("Expected 25, got %v", hash["age"])
	}
}

func TestGlobalManagementOperations(t *testing.T) {
	// 清空全局缓存确保测试独立
	Flush()

	// 测试全局管理操作
	err := SetString("temp_key", "temp_value", time.Minute)
	if err != nil {
		t.Fatalf("SetString failed: %v", err)
	}

	// EXISTS
	if !Exists("temp_key") {
		t.Error("Expected key to exist")
	}

	// SIZE
	size := Size()
	if size != 1 {
		t.Errorf("Expected size 1, got %d", size)
	}

	// KEYS
	keys := Keys()
	if len(keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(keys))
	}
	if keys[0] != "temp_key" {
		t.Errorf("Expected 'temp_key', got %s", keys[0])
	}

	// DELETE
	if !Delete("temp_key") {
		t.Error("Delete should return true")
	}

	// 验证删除后不存在
	if Exists("temp_key") {
		t.Error("Key should not exist after deletion")
	}
}

func TestGlobalStats(t *testing.T) {
	// 清空全局缓存确保测试独立
	Flush()

	// 测试统计信息
	SetString("stats_key", "stats_value", time.Minute)
	GetString("stats_key") // 增加命中次数

	stats := Stats()
	if stats == nil {
		t.Error("Stats should not be nil")
	}
}

func TestLocalCacheOperations(t *testing.T) {
	// 清空全局缓存确保测试独立
	Flush()

	// 测试局部缓存
	localCache := New()

	err := localCache.SetString("local_key", "local_value", time.Minute)
	if err != nil {
		t.Fatalf("Local SetString failed: %v", err)
	}

	value, found := localCache.GetString("local_key")
	if !found {
		t.Error("Expected local key to exist")
	}
	if value != "local_value" {
		t.Errorf("Expected 'local_value', got %s", value)
	}

	// 确保全局缓存不受影响
	if _, found = GetString("local_key"); found {
		t.Error("Local cache should not affect global cache")
	}
}
