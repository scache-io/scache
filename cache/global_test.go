package cache

import (
	"sync"
	"testing"
	"time"
)

func TestGlobalCache(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 测试全局缓存初始化
	cache := GetGlobalCache()
	if cache == nil {
		t.Fatal("GetGlobalCache() returned nil")
	}

	// 测试第二次调用返回相同实例
	cache2 := GetGlobalCache()
	if cache != cache2 {
		t.Error("GetGlobalCache() should return same instance")
	}
}

func TestGlobalSet(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 测试设置
	err := Set("key1", "value1", 0)
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// 验证设置成功
	value, found := Get("key1")
	if !found {
		t.Error("Get() should find set key")
	}
	if value != "value1" {
		t.Errorf("Get() = %v, want %v", value, "value1")
	}
}

func TestGlobalGet(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 设置值
	Set("key1", "value1", 0)

	// 测试获取
	value, found := Get("key1")
	if !found {
		t.Error("Get() should find existing key")
	}
	if value != "value1" {
		t.Errorf("Get() = %v, want %v", value, "value1")
	}

	// 测试获取不存在的键
	_, found = Get("nonexistent")
	if found {
		t.Error("Get() should not find nonexistent key")
	}
}

func TestGlobalDelete(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 设置值
	Set("key1", "value1", 0)

	// 删除存在的键
	deleted := Delete("key1")
	if !deleted {
		t.Error("Delete() should return true for existing key")
	}

	// 验证已删除
	_, found := Get("key1")
	if found {
		t.Error("Item should not be found after deletion")
	}

	// 删除不存在的键
	deleted = Delete("nonexistent")
	if deleted {
		t.Error("Delete() should return false for nonexistent key")
	}
}

func TestGlobalExists(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 设置值
	Set("key1", "value1", 0)

	// 测试存在的键
	if !Exists("key1") {
		t.Error("Exists() should return true for existing key")
	}

	// 测试不存在的键
	if Exists("nonexistent") {
		t.Error("Exists() should return false for nonexistent key")
	}
}

func TestGlobalFlush(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 设置多个值
	Set("key1", "value1", 0)
	Set("key2", "value2", 0)

	if Size() != 2 {
		t.Errorf("Expected size 2, got %d", Size())
	}

	// 清空
	Flush()

	if Size() != 0 {
		t.Errorf("Expected size 0 after flush, got %d", Size())
	}
}

func TestGlobalSize(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	if Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", Size())
	}

	Set("key1", "value1", 0)
	Set("key2", "value2", 0)

	if Size() != 2 {
		t.Errorf("Expected size 2, got %d", Size())
	}
}

func TestGlobalStats(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 配置全局缓存启用统计
	ConfigureGlobalCache(WithStats(true))

	// 执行一些操作
	Set("key1", "value1", 0)
	Get("key1")
	Get("nonexistent")
	Delete("key1")

	stats := Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	if stats.Sets != 1 {
		t.Errorf("Expected 1 set, got %d", stats.Sets)
	}
	if stats.Deletes != 1 {
		t.Errorf("Expected 1 delete, got %d", stats.Deletes)
	}
}

func TestGlobalKeys(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 设置一些键
	Set("key1", "value1", 0)
	Set("key2", "value2", 0)
	Set("key3", "value3", 0)

	keys := Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// 检查键是否在列表中
	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}

	for _, expectedKey := range []string{"key1", "key2", "key3"} {
		if !keySet[expectedKey] {
			t.Errorf("Key '%s' not found in keys list", expectedKey)
		}
	}
}

func TestGlobalGetWithExpiration(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 设置带过期时间的值
	Set("key1", "value1", time.Hour)

	value, expiration, found := GetWithExpiration("key1")
	if !found {
		t.Error("GetWithExpiration() should find existing key")
	}
	if value != "value1" {
		t.Errorf("GetWithExpiration() = %v, want %v", value, "value1")
	}
	if expiration.IsZero() {
		t.Error("Expiration should not be zero")
	}
}

func TestGlobalConcurrentAccess(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Set(string(rune(i)), i, 0)
		}(i)
	}

	// 并发读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Get(string(rune(i)))
		}(i)
	}

	wg.Wait()

	if Size() != 100 {
		t.Errorf("Expected size 100, got %d", Size())
	}
}

func TestConfigureGlobalCache(t *testing.T) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}

	// 配置全局缓存
	ConfigureGlobalCache(
		WithMaxSize(100),
		WithDefaultExpiration(time.Hour),
		WithStats(false),
	)

	// 验证配置生效
	Set("key1", "value1", 0) // 应该使用默认过期时间

	value, found := Get("key1")
	if !found {
		t.Error("Get() should find existing key")
	}
	if value != "value1" {
		t.Errorf("Get() = %v, want %v", value, "value1")
	}
}

func BenchmarkGlobalSet(b *testing.B) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}
	ConfigureGlobalCache(WithStats(false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Set("key", "value", 0)
	}
}

func BenchmarkGlobalGet(b *testing.B) {
	// 重置全局缓存
	globalCache = nil
	globalOnce = sync.Once{}
	ConfigureGlobalCache(WithStats(false))

	Set("key", "value", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get("key")
	}
}
