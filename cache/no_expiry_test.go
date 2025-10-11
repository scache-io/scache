package cache

import (
	"testing"
	"time"
)

// safeClose 辅助函数，用于在测试中安全关闭缓存
func safeClose(c Cache) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

func TestMemoryCache_DefaultConfigNoExpiry(t *testing.T) {
	cache := NewMemoryCache() // 使用默认配置
	defer safeClose(cache)

	// 设置缓存项
	err := cache.Set("permanent_key", "permanent_value")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 立即验证存在
	if value, exists := cache.Get("permanent_key"); !exists {
		t.Error("Expected key to exist immediately")
	} else if value != "permanent_value" {
		t.Errorf("Expected 'permanent_value', got '%v'", value)
	}

	// 等待一段时间，确保不会自动过期
	time.Sleep(100 * time.Millisecond)

	// 再次检查应该仍然存在
	if value, exists := cache.Get("permanent_key"); !exists {
		t.Error("Expected key to still exist (default config should not expire)")
	} else if value != "permanent_value" {
		t.Errorf("Expected 'permanent_value', got '%v'", value)
	}

	// 检查缓存大小，确保项目没有被清理
	if cache.Size() != 1 {
		t.Errorf("Expected cache size 1, got %d", cache.Size())
	}
}

func TestMemoryCache_WithDefaultTTLZero(t *testing.T) {
	cache := NewMemoryCache(WithDefaultTTL(0)) // 显式设置 TTL 为 0
	defer safeClose(cache)

	err := cache.Set("key", "value")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 等待一段时间
	time.Sleep(50 * time.Millisecond)

	// 应该仍然存在
	if _, exists := cache.Get("key"); !exists {
		t.Error("Expected key to exist when TTL is 0 (no expiry)")
	}
}

func TestMemoryCache_WithDefaultTTLNonZero(t *testing.T) {
	cache := NewMemoryCache(WithDefaultTTL(50 * time.Millisecond))
	defer safeClose(cache)

	err := cache.Set("key", "value")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 立即获取应该存在
	if _, exists := cache.Get("key"); !exists {
		t.Error("Expected key to exist immediately")
	}

	// 等待过期
	time.Sleep(60 * time.Millisecond)

	// 应该已经过期
	if _, exists := cache.Get("key"); exists {
		t.Error("Expected key to be expired")
	}
}

func TestMemoryCache_MixedExpiry(t *testing.T) {
	cache := NewMemoryCache(WithDefaultTTL(0)) // 默认永不过期
	defer safeClose(cache)

	// 设置永不过期的项
	cache.Set("permanent", "value1")

	// 设置有过期时间的项
	err := cache.SetWithTTL("temporary", "value2", 30*time.Millisecond)
	if err != nil {
		t.Fatalf("SetWithTTL failed: %v", err)
	}

	// 立即检查都存在
	if !cache.Exists("permanent") {
		t.Error("Expected permanent key to exist")
	}
	if !cache.Exists("temporary") {
		t.Error("Expected temporary key to exist")
	}

	// 等待临时项过期
	time.Sleep(40 * time.Millisecond)

	// 永久项应该存在
	if !cache.Exists("permanent") {
		t.Error("Expected permanent key to still exist")
	}

	// 临时项应该过期
	if cache.Exists("temporary") {
		t.Error("Expected temporary key to be expired")
	}
}
