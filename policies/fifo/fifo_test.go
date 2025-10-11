package fifo

import (
	"fmt"
	"testing"

	"github.com/scache/interfaces"
)

func TestFIFOPolicy_BasicOperations(t *testing.T) {
	policy := NewFIFOPolicy(3)

	// Test initial state
	if policy.Len() != 0 {
		t.Errorf("Expected initial length 0, got %d", policy.Len())
	}

	if key, should := policy.ShouldEvict(); should || key != "" {
		t.Errorf("Expected no eviction for empty policy, got should=%v, key=%s", should, key)
	}

	// Test Add operations
	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	if policy.Len() != 3 {
		t.Errorf("Expected length 3 after adding 3 items, got %d", policy.Len())
	}

	// Test Contains
	if !policy.Contains("key1") {
		t.Error("Expected policy to contain key1")
	}

	if policy.Contains("nonexistent") {
		t.Error("Expected policy to not contain nonexistent key")
	}

	// Test eviction trigger
	policy.OnAdd("key4") // Should trigger eviction of key1

	if policy.Len() != 3 {
		t.Errorf("Expected length to remain 3 after eviction, got %d", policy.Len())
	}

	if policy.Contains("key1") {
		t.Error("Expected key1 to be evicted")
	}

	if !policy.Contains("key4") {
		t.Error("Expected policy to contain key4")
	}
}

func TestFIFOPolicy_EvictionOrder(t *testing.T) {
	policy := NewFIFOPolicy(2)

	policy.OnAdd("first")
	policy.OnAdd("second")

	// Should evict "first" as it was added first
	if evicted, should := policy.ShouldEvict(); !should || evicted != "first" {
		t.Errorf("Expected to evict 'first', got should=%v, evicted=%s", should, evicted)
	}

	policy.OnAdd("third") // Triggers eviction

	keys := policy.Keys()
	expectedKeys := []string{"second", "third"}

	if len(keys) != len(expectedKeys) {
		t.Fatalf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}

	for i, key := range expectedKeys {
		if keys[i] != key {
			t.Errorf("Expected key[%d] to be '%s', got '%s'", i, key, keys[i])
		}
	}
}

func TestFIFOPolicy_OnAccess(t *testing.T) {
	policy := NewFIFOPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")

	// FIFO doesn't change order on access
	policy.OnAccess("key1")

	// Since capacity is 3 and we have 3 items, no eviction should happen yet
	if evicted, should := policy.ShouldEvict(); should || evicted != "" {
		t.Errorf("Expected no eviction yet, got should=%v, evicted=%s", should, evicted)
	}
}

func TestFIFOPolicy_OnRemove(t *testing.T) {
	policy := NewFIFOPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	policy.OnRemove("key2")

	if policy.Contains("key2") {
		t.Error("Expected key2 to be removed")
	}

	if policy.Len() != 2 {
		t.Errorf("Expected length 2 after removal, got %d", policy.Len())
	}

	// Should not evict anything yet since we have exactly 3 items
	if evicted, should := policy.ShouldEvict(); should || evicted != "" {
		t.Errorf("Expected no eviction yet, got should=%v, evicted=%s", should, evicted)
	}
}

func TestFIFOPolicy_SetMaxSize(t *testing.T) {
	policy := NewFIFOPolicy(5)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	policy.SetMaxSize(2) // Should trigger evictions

	if policy.Len() != 2 {
		t.Errorf("Expected length 2 after setting max size to 2, got %d", policy.Len())
	}

	// Should keep the last 2 items (key2, key3)
	if policy.Contains("key1") {
		t.Error("Expected key1 to be evicted when reducing max size")
	}
}

func TestFIFOPolicy_Clear(t *testing.T) {
	policy := NewFIFOPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	policy.Clear()

	if policy.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", policy.Len())
	}

	if key, should := policy.ShouldEvict(); should || key != "" {
		t.Errorf("Expected no eviction after clear, got should=%v, key=%s", should, key)
	}

	if policy.Contains("key1") || policy.Contains("key2") || policy.Contains("key3") {
		t.Error("Expected no keys to exist after clear")
	}
}

func TestFIFOPolicy_GetPosition(t *testing.T) {
	policy := NewFIFOPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	if pos := policy.GetPosition("key1"); pos != 0 {
		t.Errorf("Expected key1 at position 0, got %d", pos)
	}

	if pos := policy.GetPosition("key2"); pos != 1 {
		t.Errorf("Expected key2 at position 1, got %d", pos)
	}

	if pos := policy.GetPosition("key3"); pos != 2 {
		t.Errorf("Expected key3 at position 2, got %d", pos)
	}

	if pos := policy.GetPosition("nonexistent"); pos != -1 {
		t.Errorf("Expected -1 for nonexistent key, got %d", pos)
	}
}

func TestFIFOPolicy_ConcurrentAccess(t *testing.T) {
	policy := NewFIFOPolicy(100)
	done := make(chan bool, 10)

	// Start multiple goroutines adding items
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				policy.OnAdd(key)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all items were added
	if policy.Len() != 100 {
		t.Errorf("Expected 100 items after concurrent adds, got %d", policy.Len())
	}

	// Verify no duplicates
	keys := policy.Keys()
	keySet := make(map[string]bool)
	for _, key := range keys {
		if keySet[key] {
			t.Errorf("Found duplicate key: %s", key)
		}
		keySet[key] = true
	}
}

func TestFIFOPolicy_EdgeCases(t *testing.T) {
	t.Run("Zero capacity", func(t *testing.T) {
		policy := NewFIFOPolicy(0)

		policy.OnAdd("key1")

		if policy.Len() != 0 {
			t.Error("Zero capacity policy should not store items")
		}

		if policy.Contains("key1") {
			t.Error("Zero capacity policy should not contain items")
		}
	})

	t.Run("Single item", func(t *testing.T) {
		policy := NewFIFOPolicy(1)

		policy.OnAdd("only_key")

		if policy.Len() != 1 {
			t.Errorf("Expected 1 item, got %d", policy.Len())
		}

		policy.OnAdd("another_key")

		if policy.Contains("only_key") {
			t.Error("Expected first item to be evicted")
		}

		if !policy.Contains("another_key") {
			t.Error("Expected new item to be stored")
		}
	})
}

// Ensure FIFOPolicy implements EvictionPolicy interface
var _ interfaces.EvictionPolicy = (*FIFOPolicy)(nil)