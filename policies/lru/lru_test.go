package lru

import (
	"sync"
	"testing"
)

func TestNewLRUPolicy(t *testing.T) {
	policy := NewLRUPolicy(10)
	if policy == nil {
		t.Fatal("NewLRUPolicy() returned nil")
	}
}

func TestLRUPolicy_Access(t *testing.T) {
	policy := NewLRUPolicy(3)

	// 访问新键
	policy.Access("key1")
	if !policy.Contains("key1") {
		t.Error("key1 should exist after access")
	}

	if policy.Size() != 1 {
		t.Errorf("Expected size 1, got %d", policy.Size())
	}

	// 访问相同键
	policy.Access("key1")
	if policy.Size() != 1 {
		t.Error("Accessing existing key should not change size")
	}

	// 访问多个键
	policy.Access("key2")
	policy.Access("key3")
	if policy.Size() != 3 {
		t.Errorf("Expected size 3, got %d", policy.Size())
	}

	// 访问新键应淘汰最旧的
	policy.Access("key4")
	if policy.Size() != 3 {
		t.Error("Size should remain at capacity")
	}

	if policy.Contains("key1") {
		t.Error("key1 should be evicted")
	}

	if !policy.Contains("key4") {
		t.Error("key4 should exist")
	}
}

func TestLRUPolicy_Set(t *testing.T) {
	policy := NewLRUPolicy(2)

	policy.Set("key1")
	policy.Set("key2")

	if policy.Size() != 2 {
		t.Errorf("Expected size 2, got %d", policy.Size())
	}

	policy.Set("key3") // 应该淘汰一个

	if policy.Size() != 2 {
		t.Error("Size should remain at capacity")
	}
}

func TestLRUPolicy_Delete(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.Access("key1")
	policy.Access("key2")

	if !policy.Contains("key1") {
		t.Error("key1 should exist")
	}

	policy.Delete("key1")

	if policy.Contains("key1") {
		t.Error("key1 should be deleted")
	}

	if policy.Size() != 1 {
		t.Errorf("Expected size 1 after deletion, got %d", policy.Size())
	}

	// 删除不存在的键
	policy.Delete("nonexistent")
	if policy.Size() != 1 {
		t.Error("Deleting non-existent key should not change size")
	}
}

func TestLRUPolicy_Evict(t *testing.T) {
	policy := NewLRUPolicy(3)

	// 添加一些键
	policy.Access("key1")
	policy.Access("key2")
	policy.Access("key3")

	// 重新访问 key1 使其成为最新
	policy.Access("key1")

	// 淘汰应该移除最旧的键 (key2)
	evicted := policy.Evict()
	if evicted != "key2" {
		t.Errorf("Expected to evict 'key2', got '%s'", evicted)
	}

	if policy.Contains("key2") {
		t.Error("key2 should be evicted")
	}

	if !policy.Contains("key1") && !policy.Contains("key3") {
		t.Error("key1 and key3 should still exist")
	}
}

func TestLRUPolicy_Keys(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.Access("key1")
	policy.Access("key2")
	policy.Access("key3")

	// 重新访问 key1
	policy.Access("key1")

	keys := policy.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// 最新的键应该在前面
	if keys[0] != "key1" {
		t.Errorf("Expected first key to be 'key1', got '%s'", keys[0])
	}

	// 最旧的键应该在后面
	if keys[2] != "key2" {
		t.Errorf("Expected last key to be 'key2', got '%s'", keys[2])
	}
}

func TestLRUPolicy_Contains(t *testing.T) {
	policy := NewLRUPolicy(3)

	if policy.Contains("key1") {
		t.Error("Empty policy should not contain any keys")
	}

	policy.Access("key1")
	if !policy.Contains("key1") {
		t.Error("Policy should contain key1 after access")
	}
}

func TestLRUPolicy_UpdateCapacity(t *testing.T) {
	policy := NewLRUPolicy(2)

	policy.Access("key1")
	policy.Access("key2")

	// 增加容量
	policy.UpdateCapacity(3)
	policy.Access("key3")

	if policy.Size() != 3 {
		t.Errorf("Expected size 3 after capacity increase, got %d", policy.Size())
	}

	// 减少容量
	policy.UpdateCapacity(1)

	if policy.Size() != 1 {
		t.Errorf("Expected size 1 after capacity decrease, got %d", policy.Size())
	}
}

func TestLRUPolicy_Clear(t *testing.T) {
	policy := NewLRUPolicy(3)

	policy.Access("key1")
	policy.Access("key2")
	policy.Access("key3")

	if policy.Size() != 3 {
		t.Errorf("Expected size 3, got %d", policy.Size())
	}

	policy.Clear()

	if policy.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", policy.Size())
	}

	if policy.Contains("key1") || policy.Contains("key2") || policy.Contains("key3") {
		t.Error("No keys should exist after clear")
	}
}

func TestLRUPolicy_ConcurrentAccess(t *testing.T) {
	policy := NewLRUPolicy(10)
	var wg sync.WaitGroup

	// 并发访问
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			policy.Access(string(rune(i)))
		}(i)
	}

	wg.Wait()

	if policy.Size() != 10 {
		t.Errorf("Expected size 10 (capacity), got %d", policy.Size())
	}
}

func BenchmarkLRUPolicy_Access(b *testing.B) {
	policy := NewLRUPolicy(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		policy.Access("key")
	}
}

func BenchmarkLRUPolicy_Contains(b *testing.B) {
	policy := NewLRUPolicy(1000)
	policy.Access("key")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		policy.Contains("key")
	}
}
