package cache

import (
	"testing"
	"time"
)

// mustClose 辅助函数，用于在测试中安全关闭缓存
func mustClose(c Cache) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

func TestMemoryCache_SetAndGet(t *testing.T) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	// 测试设置和获取
	err := cache.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	value, exists := cache.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

func TestMemoryCache_SetWithTTL(t *testing.T) {
	cache := NewMemoryCache(WithDefaultTTL(100 * time.Millisecond))
	defer mustClose(cache)

	// 测试带TTL的设置
	err := cache.SetWithTTL("key1", "value1", 50*time.Millisecond)
	if err != nil {
		t.Errorf("SetWithTTL failed: %v", err)
	}

	// 立即获取应该存在
	_, exists := cache.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist immediately after setting")
	}

	// 等待过期
	time.Sleep(60 * time.Millisecond)
	_, exists = cache.Get("key1")
	if exists {
		t.Error("Expected key1 to be expired")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	cache.Set("key1", "value1")

	// 测试删除
	deleted := cache.Delete("key1")
	if !deleted {
		t.Error("Expected key1 to be deleted")
	}

	// 验证删除后不存在
	_, exists := cache.Get("key1")
	if exists {
		t.Error("Expected key1 to not exist after deletion")
	}

	// 删除不存在的键
	deleted = cache.Delete("nonexistent")
	if deleted {
		t.Error("Expected deletion of nonexistent key to return false")
	}
}

func TestMemoryCache_Exists(t *testing.T) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	cache.Set("key1", "value1")

	// 测试存在性检查
	if !cache.Exists("key1") {
		t.Error("Expected key1 to exist")
	}

	if cache.Exists("nonexistent") {
		t.Error("Expected nonexistent key to not exist")
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// 测试清空
	err := cache.Clear()
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}

	if cache.Size() != 0 {
		t.Error("Expected cache to be empty after clear")
	}
}

func TestMemoryCache_BatchOperations(t *testing.T) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	// 测试批量设置
	items := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	err := cache.SetBatch(items)
	if err != nil {
		t.Errorf("SetBatch failed: %v", err)
	}

	// 测试批量获取
	keys := []string{"key1", "key2", "key3", "nonexistent"}
	results := cache.GetBatch(keys)
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// 测试批量删除
	deleteResults := cache.DeleteBatch([]string{"key1", "key3"})
	if len(deleteResults) != 2 {
		t.Errorf("Expected 2 delete results, got %d", len(deleteResults))
	}
	if !deleteResults["key1"] || !deleteResults["key3"] {
		t.Error("Expected both keys to be deleted")
	}
}

func TestMemoryCache_SizeAndKeys(t *testing.T) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}

	keys := cache.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestMemoryCache_Stats(t *testing.T) {
	cache := NewMemoryCache(WithStatistics(true))
	defer mustClose(cache)

	// 初始统计
	stats := cache.Stats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Error("Expected initial stats to be zero")
	}

	cache.Set("key1", "value1")

	// 命中
	cache.Get("key1")
	stats = cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}

	// 未命中
	cache.Get("nonexistent")
	stats = cache.Stats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}

	expectedHitRate := float64(1) / float64(2)
	if stats.HitRate != expectedHitRate {
		t.Errorf("Expected hit rate %f, got %f", expectedHitRate, stats.HitRate)
	}
}

func TestMemoryCache_ConcurrentAccess(t *testing.T) {
	cache := NewMemoryCache(WithShards(16))
	defer mustClose(cache)

	done := make(chan bool, 2)

	// 并发写入
	go func() {
		for i := 0; i < 1000; i++ {
			cache.Set("key"+string(rune(i)), "value"+string(rune(i)))
		}
		done <- true
	}()

	// 并发读取
	go func() {
		for i := 0; i < 1000; i++ {
			cache.Get("key" + string(rune(i)))
		}
		done <- true
	}()

	// 等待完成
	<-done
	<-done

	if cache.Size() == 0 {
		t.Error("Expected cache to have items after concurrent operations")
	}
}

func TestMemoryCache_Eviction(t *testing.T) {
	// 使用单个分片来确保淘汰策略正常工作
	cache := NewMemoryCache(WithMaxSize(2), WithShards(1), WithEvictionPolicy("lru"))
	defer mustClose(cache)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3") // 应该淘汰 key1

	if cache.Exists("key1") {
		t.Error("Expected key1 to be evicted")
	}

	if !cache.Exists("key3") {
		t.Error("Expected key3 to exist")
	}
}

// 基准测试
func BenchmarkMemoryCache_Set(b *testing.B) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value")
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	cache.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkMemoryCache_SetWithTTL(b *testing.B) {
	cache := NewMemoryCache()
	defer mustClose(cache)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.SetWithTTL("key", "value", time.Minute)
	}
}

func BenchmarkMemoryCache_ConcurrentOperations(b *testing.B) {
	cache := NewMemoryCache(WithShards(16))
	defer mustClose(cache)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "key" + string(rune(i%1000))
			cache.Set(key, "value")
			cache.Get(key)
			i++
		}
	})
}
