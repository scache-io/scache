package cache

import (
	"sync"
	"testing"
	"time"

	"scache/constants"
)

func TestNewCache(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("NewCache() returned nil")
	}

	// 测试默认配置
	if c.Size() != 0 {
		t.Errorf("Expected cache size 0, got %d", c.Size())
	}

	stats := c.Stats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Errorf("Expected zero stats, got Hits: %d, Misses: %d", stats.Hits, stats.Misses)
	}
}

func TestCache_Set(t *testing.T) {
	c := NewCache(WithMaxSize(2))

	// 测试基本设置
	err := c.Set("key1", "value1", 0)
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// 测试空键
	err = c.Set("", "value", 0)
	if err == nil {
		t.Error("Set() with empty key should return error")
	}

	// 测试容量限制
	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0) // 应该淘汰一个

	if c.Size() > 2 {
		t.Errorf("Cache size exceeded limit: expected <= 2, got %d", c.Size())
	}
}

func TestCache_Get(t *testing.T) {
	c := NewCache()

	// 设置缓存项
	c.Set("key1", "value1", 0)

	// 测试获取存在的键
	value, found := c.Get("key1")
	if !found {
		t.Error("Get() should find existing key")
	}
	if value != "value1" {
		t.Errorf("Get() = %v, want %v", value, "value1")
	}

	// 测试获取不存在的键
	_, found = c.Get("nonexistent")
	if found {
		t.Error("Get() should not find nonexistent key")
	}

	// 测试空键
	_, found = c.Get("")
	if found {
		t.Error("Get() should not find empty key")
	}
}

func TestCache_Expiration(t *testing.T) {
	c := NewCache()

	// 设置带过期时间的缓存项
	c.Set("key1", "value1", time.Millisecond*100)

	// 立即获取应该成功
	_, found := c.Get("key1")
	if !found {
		t.Error("Item should be found immediately after set")
	}

	// 等待过期
	time.Sleep(time.Millisecond * 150)

	// 过期后获取应该失败
	_, found = c.Get("key1")
	if found {
		t.Error("Expired item should not be found")
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache()

	// 设置缓存项
	c.Set("key1", "value1", 0)

	// 删除存在的键
	deleted := c.Delete("key1")
	if !deleted {
		t.Error("Delete() should return true for existing key")
	}

	// 验证已删除
	_, found := c.Get("key1")
	if found {
		t.Error("Item should not be found after deletion")
	}

	// 删除不存在的键
	deleted = c.Delete("nonexistent")
	if deleted {
		t.Error("Delete() should return false for nonexistent key")
	}

	// 删除空键
	deleted = c.Delete("")
	if deleted {
		t.Error("Delete() should return false for empty key")
	}
}

func TestCache_Exists(t *testing.T) {
	c := NewCache()

	// 设置缓存项
	c.Set("key1", "value1", 0)

	// 测试存在的键
	if !c.Exists("key1") {
		t.Error("Exists() should return true for existing key")
	}

	// 测试不存在的键
	if c.Exists("nonexistent") {
		t.Error("Exists() should return false for nonexistent key")
	}

	// 测试空键
	if c.Exists("") {
		t.Error("Exists() should return false for empty key")
	}
}

func TestCache_Flush(t *testing.T) {
	c := NewCache()

	// 设置多个缓存项
	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0)

	if c.Size() != 3 {
		t.Errorf("Expected cache size 3, got %d", c.Size())
	}

	// 清空缓存
	c.Flush()

	if c.Size() != 0 {
		t.Errorf("Expected cache size 0 after flush, got %d", c.Size())
	}

	// 验证统计信息重置
	stats := c.Stats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Sets != 0 {
		t.Error("Stats should be reset after flush")
	}
}

func TestCache_Stats(t *testing.T) {
	c := NewCache(WithStats(true))

	// 初始统计
	stats := c.Stats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Sets != 0 {
		t.Error("Initial stats should be zero")
	}

	// 设置操作
	c.Set("key1", "value1", 0)
	stats = c.Stats()
	if stats.Sets != 1 {
		t.Errorf("Expected 1 set operation, got %d", stats.Sets)
	}

	// 命中操作
	c.Get("key1")
	stats = c.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}

	// 未命中操作
	c.Get("nonexistent")
	stats = c.Stats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}

	// 删除操作
	c.Delete("key1")
	stats = c.Stats()
	if stats.Deletes != 1 {
		t.Errorf("Expected 1 delete, got %d", stats.Deletes)
	}

	// 测试命中率
	expectedHitRate := float64(stats.Hits) / float64(stats.Hits+stats.Misses)
	if stats.HitRate != expectedHitRate {
		t.Errorf("Expected hit rate %f, got %f", expectedHitRate, stats.HitRate)
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(string(rune(i)), i, 0)
		}(i)
	}

	// 并发读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Get(string(rune(i)))
		}(i)
	}

	wg.Wait()

	// 验证数据一致性
	if c.Size() != 100 {
		t.Errorf("Expected cache size 100, got %d", c.Size())
	}

	stats := c.Stats()
	if stats.Sets != 100 {
		t.Errorf("Expected 100 sets, got %d", stats.Sets)
	}
}

func TestCache_Cleanup(t *testing.T) {
	c := NewCache(
		WithCleanupInterval(time.Millisecond*50),
		WithStats(false), // 关闭统计避免干扰
	)

	// 设置带过期时间的缓存项
	c.Set("key1", "value1", time.Millisecond*100)
	c.Set("key2", "value2", time.Millisecond*100)

	// 等待清理
	time.Sleep(time.Millisecond * 200)

	// 检查是否被清理
	if c.Exists("key1") || c.Exists("key2") {
		t.Error("Expired items should be cleaned up")
	}
}

func TestCache_GetWithExpiration(t *testing.T) {
	c := NewCache()
	if cache, ok := c.(*MemoryCache); ok {
		// 设置带过期时间的缓存项
		cache.Set("key1", "value1", time.Hour)

		value, expiration, found := cache.GetWithExpiration("key1")
		if !found {
			t.Error("Item should be found")
		}

		if value != "value1" {
			t.Errorf("Expected value 'value1', got %v", value)
		}

		if expiration.IsZero() {
			t.Error("Expiration should not be zero")
		}
	}
}

func TestCache_DefaultExpiration(t *testing.T) {
	c := NewCache(
		WithDefaultExpiration(time.Millisecond * 100),
	)

	// 设置不带TTL的缓存项（应使用默认过期时间）
	c.Set("key1", "value1", 0)

	// 立即获取应该成功
	_, found := c.Get("key1")
	if !found {
		t.Error("Item should be found immediately")
	}

	// 等待过期
	time.Sleep(time.Millisecond * 150)

	// 过期后获取应该失败
	_, found = c.Get("key1")
	if found {
		t.Error("Item with default expiration should expire")
	}
}

func BenchmarkCache_Set(b *testing.B) {
	c := NewCache(WithStats(false))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set("key", "value", 0)
	}
}

func BenchmarkCache_Get(b *testing.B) {
	c := NewCache(WithStats(false))
	c.Set("key", "value", 0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Get("key")
	}
}

func BenchmarkCache_Concurrent(b *testing.B) {
	c := NewCache(WithStats(false))
	c.Set("key", "value", 0)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set("key", "value", 0)
			c.Get("key")
		}
	})
}
