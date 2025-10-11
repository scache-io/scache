package constants

import (
	"testing"
	"time"
)

func TestConstants_ConfigValues(t *testing.T) {
	// Test basic configuration constants
	if DefaultMaxSize != 10000 {
		t.Errorf("Expected DefaultMaxSize = 10000, got %d", DefaultMaxSize)
	}

	if DefaultShards != 16 {
		t.Errorf("Expected DefaultShards = 16, got %d", DefaultShards)
	}

	if DefaultCleanupInterval != 10*time.Minute {
		t.Errorf("Expected DefaultCleanupInterval = 10m, got %v", DefaultCleanupInterval)
	}

	if DefaultTTL != 0 {
		t.Errorf("Expected DefaultTTL = 0, got %v", DefaultTTL)
	}
}

func TestConstants_StrategyConstants(t *testing.T) {
	// Test strategy constants
	if LRUStrategy != "lru" {
		t.Errorf("Expected LRUStrategy = 'lru', got '%s'", LRUStrategy)
	}

	if LFUStrategy != "lfu" {
		t.Errorf("Expected LFUStrategy = 'lfu', got '%s'", LFUStrategy)
	}

	if FIFOStrategy != "fifo" {
		t.Errorf("Expected FIFOStrategy = 'fifo', got '%s'", FIFOStrategy)
	}
}

func TestConstants_GlobalConstants(t *testing.T) {
	// Test global cache constants
	if DefaultCacheName != "default" {
		t.Errorf("Expected DefaultCacheName = 'default', got '%s'", DefaultCacheName)
	}

	if ManagerTimeout != 30*time.Second {
		t.Errorf("Expected ManagerTimeout = 30s, got %v", ManagerTimeout)
	}
}

func TestConstants_PerformanceConstants(t *testing.T) {
	// Test performance limits
	if MaxKeyLength != 256 {
		t.Errorf("Expected MaxKeyLength = 256, got %d", MaxKeyLength)
	}

	if MinKeyLength != 1 {
		t.Errorf("Expected MinKeyLength = 1, got %d", MinKeyLength)
	}

	expectedMaxValueSize := int64(10 * 1024 * 1024) // 10MB
	if MaxValueSize != expectedMaxValueSize {
		t.Errorf("Expected MaxValueSize = %d, got %d", expectedMaxValueSize, MaxValueSize)
	}
}

func TestConstants_ErrorMessages(t *testing.T) {
	// Test error message constants are not empty
	errorConstants := []struct {
		name  string
		value string
	}{
		{"ErrCacheNotFound", ErrCacheNotFound},
		{"ErrInvalidCacheName", ErrInvalidCacheName},
		{"ErrCacheAlreadyExists", ErrCacheAlreadyExists},
		{"ErrInvalidStrategy", ErrInvalidStrategy},
		{"ErrKeyNotFound", ErrKeyNotFound},
		{"ErrKeyTooLong", ErrKeyTooLong},
		{"ErrKeyEmpty", ErrKeyEmpty},
		{"ErrValueTooLarge", ErrValueTooLarge},
		{"ErrCacheClosed", ErrCacheClosed},
	}

	for _, tc := range errorConstants {
		t.Run(tc.name, func(t *testing.T) {
			if tc.value == "" {
				t.Errorf("Error constant %s should not be empty", tc.name)
			}
		})
	}
}

func TestConstants_LogPrefixes(t *testing.T) {
	// Test log prefixes
	if LogPrefixCache == "" {
		t.Error("LogPrefixCache should not be empty")
	}

	if LogPrefixManager == "" {
		t.Error("LogPrefixManager should not be empty")
	}

	if LogPrefixGlobal == "" {
		t.Error("LogPrefixGlobal should not be empty")
	}
}

func TestConstants_StatsConstants(t *testing.T) {
	// Test statistics constants
	if StatsUpdateInterval != time.Second {
		t.Errorf("Expected StatsUpdateInterval = 1s, got %v", StatsUpdateInterval)
	}

	if HitRateThreshold != 0.8 {
		t.Errorf("Expected HitRateThreshold = 0.8, got %f", HitRateThreshold)
	}
}

func TestConstants_SerializationConstants(t *testing.T) {
	// Test serialization constants
	if JSONEncoding != "json" {
		t.Errorf("Expected JSONEncoding = 'json', got '%s'", JSONEncoding)
	}

	if GobEncoding != "gob" {
		t.Errorf("Expected GobEncoding = 'gob', got '%s'", GobEncoding)
	}
}

func TestConstants_Validations(t *testing.T) {
	// Validate constant relationships
	if DefaultMaxSize <= 0 {
		t.Error("DefaultMaxSize should be positive")
	}

	if DefaultShards <= 0 {
		t.Error("DefaultShards should be positive")
	}

	if DefaultCleanupInterval <= 0 {
		t.Error("DefaultCleanupInterval should be positive")
	}

	if DefaultMaxSize < DefaultShards {
		t.Error("DefaultMaxSize should be >= DefaultShards")
	}

	// Validate strategy constants are lowercase
	if LRUStrategy != "lru" || LFUStrategy != "lfu" || FIFOStrategy != "fifo" {
		t.Error("Strategy constants should be lowercase")
	}

	// Validate performance constants are reasonable
	if MaxKeyLength <= 0 || MaxKeyLength > 1024 {
		t.Errorf("MaxKeyLength should be between 1 and 1024, got %d", MaxKeyLength)
	}

	if MinKeyLength <= 0 || MinKeyLength > MaxKeyLength {
		t.Errorf("MinKeyLength should be between 1 and MaxKeyLength, got %d", MinKeyLength)
	}

	if MaxValueSize <= 0 {
		t.Errorf("MaxValueSize should be positive, got %d", MaxValueSize)
	}
}

func TestConstants_StrategyUniqueness(t *testing.T) {
	// Test that strategy constants are unique
	strategies := []string{LRUStrategy, LFUStrategy, FIFOStrategy}
	seen := make(map[string]bool)

	for _, strategy := range strategies {
		if seen[strategy] {
			t.Errorf("Duplicate strategy found: %s", strategy)
		}
		seen[strategy] = true
	}
}

// Benchmark tests
func BenchmarkConstants_Access(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultMaxSize
		_ = LRUStrategy
		_ = ErrCacheNotFound
	}
}