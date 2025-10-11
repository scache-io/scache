package lfu

import (
	"fmt"
	"testing"

	"github.com/scache/interfaces"
)

func TestLFUPolicy_BasicOperations(t *testing.T) {
	policy := NewLFUPolicy(3)

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
	policy.OnAdd("key4") // Should trigger eviction of lowest frequency item

	if policy.Len() != 3 {
		t.Errorf("Expected length to remain 3 after eviction, got %d", policy.Len())
	}

	// Check that one of the original keys was evicted (the one with lowest frequency)
	evictedKeys := 0
	for _, key := range []string{"key1", "key2", "key3"} {
		if !policy.Contains(key) {
			evictedKeys++
		}
	}

	if evictedKeys != 1 {
		t.Errorf("Expected exactly 1 key to be evicted, got %d", evictedKeys)
	}

	if !policy.Contains("key4") {
		t.Error("Expected policy to contain key4")
	}
}

func TestLFUPolicy_FrequencyTracking(t *testing.T) {
	policy := NewLFUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	// Access key1 multiple times to increase its frequency
	for i := 0; i < 5; i++ {
		policy.OnAccess("key1")
	}

	// Access key2 a few times
	for i := 0; i < 2; i++ {
		policy.OnAccess("key2")
	}

	// key3 has frequency 1, key2 has frequency 3, key1 has frequency 6

	// Adding a new key should evict key3 (lowest frequency)
	policy.OnAdd("key4")

	if policy.Contains("key3") {
		t.Error("Expected key3 to be evicted (lowest frequency)")
	}

	if !policy.Contains("key1") {
		t.Error("Expected key1 to remain (highest frequency)")
	}

	if !policy.Contains("key2") {
		t.Error("Expected key2 to remain (middle frequency)")
	}
}

func TestLFUPolicy_GetFrequency(t *testing.T) {
	policy := NewLFUPolicy(3)

	policy.OnAdd("key1")

	if freq := policy.GetFrequency("key1"); freq != 1 {
		t.Errorf("Expected initial frequency 1, got %d", freq)
	}

	// Access the key multiple times
	for i := 0; i < 5; i++ {
		policy.OnAccess("key1")
	}

	if freq := policy.GetFrequency("key1"); freq != 6 {
		t.Errorf("Expected frequency 6 after 5 accesses, got %d", freq)
	}

	if freq := policy.GetFrequency("nonexistent"); freq != 0 {
		t.Errorf("Expected frequency 0 for nonexistent key, got %d", freq)
	}
}

func TestLFUPolicy_OnRemove(t *testing.T) {
	policy := NewLFUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	// Access key1 to increase its frequency
	policy.OnAccess("key1")
	policy.OnAccess("key1")

	policy.OnRemove("key2")

	if policy.Contains("key2") {
		t.Error("Expected key2 to be removed")
	}

	if policy.Len() != 2 {
		t.Errorf("Expected length 2 after removal, got %d", policy.Len())
	}

	// Verify frequencies are maintained
	if freq := policy.GetFrequency("key1"); freq != 3 {
		t.Errorf("Expected key1 frequency to be 3, got %d", freq)
	}
}

func TestLFUPolicy_SetMaxSize(t *testing.T) {
	policy := NewLFUPolicy(5)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	// Give different frequencies
	policy.OnAccess("key2")
	policy.OnAccess("key2")

	policy.SetMaxSize(2) // Should trigger evictions

	if policy.Len() != 2 {
		t.Errorf("Expected length 2 after setting max size to 2, got %d", policy.Len())
	}

	// Should keep the items with highest frequencies (key2 and key1 or key3)
	if !policy.Contains("key2") {
		t.Error("Expected key2 to remain (highest frequency)")
	}
}

func TestLFUPolicy_Clear(t *testing.T) {
	policy := NewLFUPolicy(3)

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

func TestLFUPolicy_Keys(t *testing.T) {
	policy := NewLFUPolicy(3)

	policy.OnAdd("key1")
	policy.OnAdd("key2")
	policy.OnAdd("key3")

	keys := policy.Keys()

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// In LFU, keys are returned in frequency order (lowest to highest)
	keySet := make(map[string]bool)
	for _, key := range keys {
		if keySet[key] {
			t.Errorf("Found duplicate key: %s", key)
		}
		keySet[key] = true
	}

	// All keys should be present
	for _, key := range []string{"key1", "key2", "key3"} {
		if !keySet[key] {
			t.Errorf("Expected key %s to be present in keys", key)
		}
	}
}

func TestLFUPolicy_ConcurrentAccess(t *testing.T) {
	policy := NewLFUPolicy(100)
	done := make(chan bool, 10)

	// Start multiple goroutines adding and accessing items
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				policy.OnAdd(key)

				// Access some items to increase frequency
				if j%2 == 0 {
					policy.OnAccess(key)
				}
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
		t.Errorf("Expected 100 items after concurrent operations, got %d", policy.Len())
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

func TestLFUPolicy_EdgeCases(t *testing.T) {
	t.Run("Zero capacity", func(t *testing.T) {
		policy := NewLFUPolicy(0)

		policy.OnAdd("key1")

		if policy.Len() != 0 {
			t.Error("Zero capacity policy should not store items")
		}

		if policy.Contains("key1") {
			t.Error("Zero capacity policy should not contain items")
		}
	})

	t.Run("Single item", func(t *testing.T) {
		policy := NewLFUPolicy(1)

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

	t.Run("Accessing non-existent item", func(t *testing.T) {
		policy := NewLFUPolicy(3)

		// Should not panic
		policy.OnAccess("nonexistent")

		if policy.Len() != 0 {
			t.Error("Accessing non-existent key should not affect policy")
		}
	})
}

func TestLFUPolicy_FrequencyTieBreaking(t *testing.T) {
	policy := NewLFUPolicy(3) // Use capacity 3 instead of 2

	policy.OnAdd("key1")
	policy.OnAdd("key2")

	// Give both items same frequency
	policy.OnAccess("key1")
	policy.OnAccess("key2")

	// Add a third item - should not trigger eviction yet
	policy.OnAdd("key3")

	if policy.Len() != 3 {
		t.Errorf("Expected 3 items, got %d", policy.Len())
	}

	// All three keys should be present
	if !policy.Contains("key1") || !policy.Contains("key2") || !policy.Contains("key3") {
		t.Error("All keys should be present")
	}

	// Now add a fourth item to trigger eviction
	policy.OnAdd("key4")

	if policy.Len() != 3 {
		t.Errorf("Expected 3 items after eviction, got %d", policy.Len())
	}

	// Should contain key4 and two of the original three
	if !policy.Contains("key4") {
		t.Error("Expected key4 to be present")
	}

	// Count present items
	presentKeys := 0
	for _, key := range []string{"key1", "key2", "key3"} {
		if policy.Contains(key) {
			presentKeys++
		}
	}

	if presentKeys != 2 {
		t.Errorf("Expected exactly 2 of the original keys to be present, got %d", presentKeys)
	}
}

// Ensure LFUPolicy implements EvictionPolicy interface
var _ interfaces.EvictionPolicy = (*LFUPolicy)(nil)