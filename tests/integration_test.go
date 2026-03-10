package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/scache-io/scache"
	"github.com/scache-io/scache/config"
)

// ==================== Basic operation tests ====================

func TestStringOperations(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	// Set & Get
	err := cache.SetString("key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("SetString failed: %v", err)
	}

	value, found := cache.GetString("key1")
	if !found {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %s", value)
	}

	// Delete
	if !cache.Delete("key1") {
		t.Error("Delete should return true")
	}
	if cache.Exists("key1") {
		t.Error("Key should not exist after deletion")
	}
}

func TestListOperations(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	err := cache.SetList("list1", []interface{}{"item1", "item2"}, time.Minute)
	if err != nil {
		t.Fatalf("SetList failed: %v", err)
	}

	items, found := cache.GetList("list1")
	if !found {
		t.Error("Expected list to exist")
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestHashOperations(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	data := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}
	err := cache.SetHash("hash1", data, time.Minute)
	if err != nil {
		t.Fatalf("SetHash failed: %v", err)
	}

	hash, found := cache.GetHash("hash1")
	if !found {
		t.Error("Expected hash to exist")
	}
	if hash["name"] != "Alice" {
		t.Errorf("Expected 'Alice', got %v", hash["name"])
	}
}

func TestStructOperations(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	user := User{ID: 1, Name: "Test"}
	err := cache.Store("user:1", user, time.Hour)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	var loadedUser User
	err = cache.Load("user:1", &loadedUser)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loadedUser.Name != user.Name {
		t.Errorf("Expected %s, got %s", user.Name, loadedUser.Name)
	}
}

// ==================== TTL 测试 ====================

func TestTTL(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	// 设置带 TTL 的数据
	cache.SetString("ttl_key", "value", 100*time.Millisecond)

	// 立即读取应该存在
	if _, found := cache.GetString("ttl_key"); !found {
		t.Error("Key should exist immediately after set")
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 过期后应该不存在
	if _, found := cache.GetString("ttl_key"); found {
		t.Error("Key should not exist after TTL expiration")
	}
}

func TestExpire(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	cache.SetString("expire_key", "value", time.Hour)

	// 修改 TTL
	cache.Expire("expire_key", 100*time.Millisecond)

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	if _, found := cache.GetString("expire_key"); found {
		t.Error("Key should be expired after Expire()")
	}
}

// ==================== 并发测试 ====================

func TestConcurrentAccess(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// 并发写入
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := string(rune(n*10 + j))
				if err := cache.SetString(key, "value", time.Hour); err != nil {
					errors <- err
				}
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := string(rune(n*10 + j))
				cache.GetString(key)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}
}

func TestConcurrentExpiredDeletion(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	// 设置一个立即过期的键
	cache.SetString("expire_key", "value", time.Millisecond)

	// 等待过期
	time.Sleep(10 * time.Millisecond)

	// 并发读取过期键，验证不会出现竞态条件
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				cache.GetString("expire_key")
				cache.Exists("expire_key")
			}
		}()
	}

	wg.Wait()
}

// ==================== 容量与淘汰测试 ====================

func TestMaxSizeLimit(t *testing.T) {
	cfg := &config.EngineConfig{
		MaxSize:                   10,
		MemoryThreshold:           0.9,
		DefaultExpiration:         0,
		BackgroundCleanupInterval: time.Minute,
	}
	cache := scache.New(cfg)

	// 插入超过容量的数据
	for i := 0; i < 15; i++ {
		cache.SetString(string(rune(i)), "value", time.Hour)
	}

	// 容量应该被限制
	if cache.Size() > 10 {
		t.Errorf("Cache size should not exceed MaxSize, got %d", cache.Size())
	}
}

func TestMaxSizeZeroDisablesEviction(t *testing.T) {
	cfg := &config.EngineConfig{
		MaxSize:                   0, // 无限制
		MemoryThreshold:           0.9,
		DefaultExpiration:         0,
		BackgroundCleanupInterval: time.Minute,
	}
	cache := scache.New(cfg)

	// 插入超过原默认容量(100)的数据
	for i := 0; i < 200; i++ {
		err := cache.SetString(string(rune(i)), "value", time.Hour)
		if err != nil {
			t.Errorf("Expected no error with MaxSize=0, got: %v", err)
		}
	}

	// 验证所有数据都存在
	if cache.Size() != 200 {
		t.Errorf("Expected 200 items with MaxSize=0, got %d", cache.Size())
	}
}

// ==================== 全局缓存测试 ====================

func TestGlobalCache(t *testing.T) {
	scache.Flush()

	// 测试全局字符串操作
	err := scache.SetString("global_key", "value", time.Minute)
	if err != nil {
		t.Fatalf("Global SetString failed: %v", err)
	}

	value, found := scache.GetString("global_key")
	if !found || value != "value" {
		t.Errorf("Global GetString failed")
	}

	// 测试管理操作
	if !scache.Exists("global_key") {
		t.Error("Global Exists failed")
	}

	if scache.Size() != 1 {
		t.Errorf("Global Size failed, expected 1, got %d", scache.Size())
	}

	scache.Flush()
	if scache.Size() != 0 {
		t.Error("Global Flush failed")
	}
}

// ==================== 统计测试 ====================

func TestStats(t *testing.T) {
	cache := scache.New(config.DefaultEngineConfig())

	cache.SetString("stats_key", "value", time.Minute)
	cache.GetString("stats_key") // 命中
	cache.GetString("not_exist") // 未命中

	stats := cache.Stats().(map[string]interface{})

	if stats["hits"].(int64) != 1 {
		t.Errorf("Expected 1 hit, got %d", stats["hits"])
	}
	if stats["misses"].(int64) != 1 {
		t.Errorf("Expected 1 miss, got %d", stats["misses"])
	}
}
