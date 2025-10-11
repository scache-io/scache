package types

import (
	"testing"
	"time"
)

func TestCacheItem_IsExpired(t *testing.T) {
	// 测试永不过期的项
	item := &CacheItem{}
	if item.IsExpired() {
		t.Error("Zero expiration time should not be expired")
	}

	// 测试未过期的项
	item = &CacheItem{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if item.IsExpired() {
		t.Error("Future expiration should not be expired")
	}

	// 测试已过期的项
	item = &CacheItem{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	if !item.IsExpired() {
		t.Error("Past expiration should be expired")
	}
}

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	if config.DefaultExpiration != 0 {
		t.Errorf("Expected default expiration 0, got %v", config.DefaultExpiration)
	}

	if config.CleanupInterval != time.Minute*10 {
		t.Errorf("Expected cleanup interval 10m, got %v", config.CleanupInterval)
	}

	if config.MaxSize != 0 {
		t.Errorf("Expected max size 0, got %d", config.MaxSize)
	}

	if !config.EnableStats {
		t.Error("Expected stats to be enabled by default")
	}

	if config.InitialCapacity != 16 {
		t.Errorf("Expected initial capacity 16, got %d", config.InitialCapacity)
	}
}

func TestCacheOptions(t *testing.T) {
	config := DefaultCacheConfig()

	// 测试 WithDefaultExpiration
	option := WithDefaultExpiration(time.Hour)
	option(config)
	if config.DefaultExpiration != time.Hour {
		t.Errorf("WithDefaultExpiration failed, got %v", config.DefaultExpiration)
	}

	// 重置配置
	config = DefaultCacheConfig()

	// 测试 WithCleanupInterval
	option = WithCleanupInterval(time.Minute * 5)
	option(config)
	if config.CleanupInterval != time.Minute*5 {
		t.Errorf("WithCleanupInterval failed, got %v", config.CleanupInterval)
	}

	// 重置配置
	config = DefaultCacheConfig()

	// 测试 WithMaxSize
	option = WithMaxSize(1000)
	option(config)
	if config.MaxSize != 1000 {
		t.Errorf("WithMaxSize failed, got %d", config.MaxSize)
	}

	// 重置配置
	config = DefaultCacheConfig()

	// 测试 WithStats
	option = WithStats(false)
	option(config)
	if config.EnableStats {
		t.Error("WithStats(false) failed")
	}

	// 重置配置
	config = DefaultCacheConfig()

	// 测试 WithInitialCapacity
	option = WithInitialCapacity(32)
	option(config)
	if config.InitialCapacity != 32 {
		t.Errorf("WithInitialCapacity failed, got %d", config.InitialCapacity)
	}
}

func TestCacheStats_Hit(t *testing.T) {
	stats := &CacheStats{}

	stats.Hit()
	hits, misses, sets, deletes := stats.GetStats()

	if hits != 1 {
		t.Errorf("Expected 1 hit, got %d", hits)
	}
	if misses != 0 || sets != 0 || deletes != 0 {
		t.Error("Other stats should remain zero")
	}
}

func TestCacheStats_Miss(t *testing.T) {
	stats := &CacheStats{}

	stats.Miss()
	hits, misses, sets, deletes := stats.GetStats()

	if misses != 1 {
		t.Errorf("Expected 1 miss, got %d", misses)
	}
	if hits != 0 || sets != 0 || deletes != 0 {
		t.Error("Other stats should remain zero")
	}
}

func TestCacheStats_Set(t *testing.T) {
	stats := &CacheStats{}

	stats.Set()
	hits, misses, sets, deletes := stats.GetStats()

	if sets != 1 {
		t.Errorf("Expected 1 set, got %d", sets)
	}
	if hits != 0 || misses != 0 || deletes != 0 {
		t.Error("Other stats should remain zero")
	}
}

func TestCacheStats_Delete(t *testing.T) {
	stats := &CacheStats{}

	stats.Delete()
	hits, misses, sets, deletes := stats.GetStats()

	if deletes != 1 {
		t.Errorf("Expected 1 delete, got %d", deletes)
	}
	if hits != 0 || misses != 0 || sets != 0 {
		t.Error("Other stats should remain zero")
	}
}

func TestCacheStats_HitRate(t *testing.T) {
	stats := &CacheStats{}

	// 测试零命中率
	if stats.HitRate() != 0 {
		t.Error("Initial hit rate should be 0")
	}

	// 测试一些命中和未命中
	stats.Hit()
	stats.Hit()
	stats.Miss()

	expected := 2.0 / 3.0
	if stats.HitRate() != expected {
		t.Errorf("Expected hit rate %f, got %f", expected, stats.HitRate())
	}

	// 测试全部命中
	stats.Hit()
	stats.Hit()

	expected = 4.0 / 5.0
	if stats.HitRate() != expected {
		t.Errorf("Expected hit rate %f, got %f", expected, stats.HitRate())
	}
}

func TestCacheStats_Reset(t *testing.T) {
	stats := &CacheStats{}

	// 添加一些统计数据
	stats.Hit()
	stats.Miss()
	stats.Set()
	stats.Delete()

	// 重置
	stats.Reset()

	// 验证重置后的统计
	hits, misses, sets, deletes := stats.GetStats()
	if hits != 0 || misses != 0 || sets != 0 || deletes != 0 {
		t.Error("Stats should be zero after reset")
	}

	if stats.HitRate() != 0 {
		t.Error("Hit rate should be 0 after reset")
	}
}

func TestCacheStats_ConcurrentAccess(t *testing.T) {
	stats := &CacheStats{}
	done := make(chan bool, 4)

	// 并发增加统计
	go func() {
		for i := 0; i < 1000; i++ {
			stats.Hit()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 500; i++ {
			stats.Miss()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 300; i++ {
			stats.Set()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 200; i++ {
			stats.Delete()
		}
		done <- true
	}()

	// 等待所有协程完成
	for i := 0; i < 4; i++ {
		<-done
	}

	// 验证最终统计
	hits, misses, sets, deletes := stats.GetStats()
	if hits != 1000 {
		t.Errorf("Expected 1000 hits, got %d", hits)
	}
	if misses != 500 {
		t.Errorf("Expected 500 misses, got %d", misses)
	}
	if sets != 300 {
		t.Errorf("Expected 300 sets, got %d", sets)
	}
	if deletes != 200 {
		t.Errorf("Expected 200 deletes, got %d", deletes)
	}
}
