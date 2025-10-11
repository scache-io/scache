package lru

import (
	"testing"
)

func TestLRUPolicy_OnAccess(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	// 访问 key1，应该将其移到前面
	policy.OnAccess("key1")

	// 添加 key4，应该淘汰 key2（最旧的）
	key, shouldEvict := policy.ShouldEvict()
	if shouldEvict && key == "key2" {
		policy.OnAdd("key4")
	}

	if policy.Contains("key2") {
		t.Error("Expected key2 to be evicted")
	}

	if !policy.Contains("key1") {
		t.Error("Expected key1 to still exist")
	}
}

func TestLRUPolicy_OnAdd(t *testing.T) {
	policy := NewLRUPolicy(2)

	policy.OnAdd("key1")
	if policy.Len() != 1 {
		t.Errorf("Expected length 1, got %d", policy.Len())
	}

	policy.OnAdd("key2")
	if policy.Len() != 2 {
		t.Errorf("Expected length 2, got %d", policy.Len())
	}

	// 添加第三个，应该淘汰第一个
	policy.OnAdd("key3")
	if policy.Len() != 2 {
		t.Errorf("Expected length 2, got %d", policy.Len())
	}

	if policy.Contains("key1") {
		t.Error("Expected key1 to be evicted")
	}
}

func TestLRUPolicy_OnRemove(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	policy.OnRemove("key2")

	if policy.Contains("key2") {
		t.Error("Expected key2 to be removed")
	}

	if policy.Len() != 2 {
		t.Errorf("Expected length 2, got %d", policy.Len())
	}
}

func TestLRUPolicy_ShouldEvict(t *testing.T) {
	policy := NewLRUPolicy(2)

	// 空策略不应该需要淘汰
	key, shouldEvict := policy.ShouldEvict()
	if shouldEvict {
		t.Error("Empty policy should not need eviction")
	}

	policy.OnAdd("key1")
	policy.OnAdd("key2")

	// 达到容量后应该需要淘汰
	key, shouldEvict = policy.ShouldEvict()
	if !shouldEvict {
		t.Error("Full policy should need eviction")
	}
	if key != "key1" {
		t.Errorf("Expected key1 to be evicted, got %s", key)
	}
}

func TestLRUPolicy_SetMaxSize(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	// 减少容量
	policy.SetMaxSize(2)

	if policy.Len() != 2 {
		t.Errorf("Expected length 2 after resizing, got %d", policy.Len())
	}

	// 应该保留最新的两个元素
	keys := policy.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestLRUPolicy_Keys(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	// 访问 key1，改变顺序
	policy.OnAccess("key1")

	keys := policy.Keys()
	expected := []string{"key1", "key3", "key2"} // 从新到旧

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	for i, expectedKey := range expected {
		if i >= len(keys) || keys[i] != expectedKey {
			t.Errorf("Expected key %s at position %d, got %s", expectedKey, i, keys[i])
		}
	}
}

func TestLRUPolicy_Clear(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")

	policy.Clear()

	if policy.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", policy.Len())
	}

	if policy.Contains("key1") || policy.Contains("key2") {
		t.Error("Expected no keys to exist after clear")
	}
}

func TestLRUPolicy_Contains(t *testing.T) {
	policy := NewLRUPolicy(3)

	if policy.Contains("nonexistent") {
		t.Error("Expected nonexistent key to not exist")
	}

	policy.OnAdd("key1")

	if !policy.Contains("key1") {
		t.Error("Expected key1 to exist")
	}
}

// 基准测试
func BenchmarkLRUPolicy_OnAccess(b *testing.B) {
	policy := NewLRUPolicy(1000)

	// 预填充
	for i := 0; i < 1000; i++ {
		policy.OnAdd("key" + string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		policy.OnAccess("key500")
	}
}

func BenchmarkLRUPolicy_OnAdd(b *testing.B) {
	policy := NewLRUPolicy(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		policy.OnAdd("key" + string(rune(i)))
	}
}

func BenchmarkLRUPolicy_ConcurrentOperations(b *testing.B) {
	policy := NewLRUPolicy(1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "key" + string(rune(i%1000))
			policy.OnAdd(key)
			policy.OnAccess(key)
			i++
		}
	})
}
