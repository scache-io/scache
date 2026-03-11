package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/scache-io/scache"
	"github.com/scache-io/scache/config"
)

// ==================== GC Optimization Benchmarks ====================
// These benchmarks are designed to measure the effectiveness of the
// object pooling and GC optimization improvements

// BenchmarkGCStringOperations benchmarks string operations with GC optimization
func BenchmarkGCStringOperations(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		cache.SetString(key, fmt.Sprintf("value-%d", i), time.Minute)
		if i%10 == 0 {
			cache.Delete(key)
		}
	}
}

// BenchmarkGCListOperations benchmarks list operations with GC optimization
func BenchmarkGCListOperations(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())
	list := []interface{}{"item1", "item2", "item3", "item4", "item5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("list-%d", i%1000)
		cache.SetList(key, list, time.Minute)
		if i%10 == 0 {
			cache.Delete(key)
		}
	}
}

// BenchmarkGCHashOperations benchmarks hash operations with GC optimization
func BenchmarkGCHashOperations(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())
	hash := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("hash-%d", i%1000)
		cache.SetHash(key, hash, time.Minute)
		if i%10 == 0 {
			cache.Delete(key)
		}
	}
}

// BenchmarkGCStructOperations benchmarks struct operations with GC optimization
func BenchmarkGCStructOperations(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	type User struct {
		ID   int
		Name string
		Age  int
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("user-%d", i%1000)
		user := User{ID: i, Name: "User", Age: 25}
		cache.Store(key, user, time.Minute)
		if i%10 == 0 {
			cache.Delete(key)
		}
	}
}

// BenchmarkGCHighDeleteRate benchmarks with high delete rate to test pool effectiveness
func BenchmarkGCHighDeleteRate(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%100)
		cache.SetString(key, "value", time.Minute)
		// Delete 50% of the time to trigger pool returns
		if i%2 == 0 {
			cache.Delete(key)
		}
	}
}

// BenchmarkGCConcurrentDelete benchmarks concurrent operations with deletes
func BenchmarkGCConcurrentDelete(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%1000)
			if i%3 == 0 {
				cache.Delete(key)
			} else {
				cache.GetString(key)
			}
			i++
		}
	})
}

// BenchmarkGCExpireOperations benchmarks TTL expiration with pool reuse
func BenchmarkGCExpireOperations(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		cache.SetString(key, "value", 10*time.Millisecond)
		// Some keys will expire and trigger pool returns
		if i%100 == 0 {
			time.Sleep(5 * time.Millisecond)
		}
	}
}

// BenchmarkGCFlush benchmarks flush operations with pool returns
func BenchmarkGCFlush(b *testing.B) {
	cfg := &config.EngineConfig{
		MaxSize: 1000,
	}
	cache := scache.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Fill cache
		for j := 0; j < 100; j++ {
			key := fmt.Sprintf("key-%d", j)
			cache.SetString(key, "value", time.Minute)
		}
		// Flush to trigger pool returns
		cache.Flush()
	}
}

// BenchmarkGCEviction benchmarks LRU eviction with pool returns
func BenchmarkGCEviction(b *testing.B) {
	cfg := &config.EngineConfig{
		MaxSize: 1000,
	}
	cache := scache.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
		// Evictions will trigger pool returns
	}
}

// BenchmarkGCMixedDataTypes benchmarks mixed data type operations
func BenchmarkGCMixedDataTypes(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())
	list := []interface{}{"item1", "item2", "item3"}
	hash := map[string]interface{}{"f1": "v1", "f2": "v2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			key := fmt.Sprintf("str-%d", i%1000)
			cache.SetString(key, "value", time.Minute)
		case 1:
			key := fmt.Sprintf("list-%d", i%1000)
			cache.SetList(key, list, time.Minute)
		case 2:
			key := fmt.Sprintf("hash-%d", i%1000)
			cache.SetHash(key, hash, time.Minute)
		case 3:
			key := fmt.Sprintf("user-%d", i%1000)
			cache.Store(key, i, time.Minute)
		}
		if i%20 == 0 {
			key := fmt.Sprintf("key-%d", i%1000)
			cache.Delete(key)
		}
	}
}

// BenchmarkGCStats benchmarks stats collection overhead
func BenchmarkGCStats(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Stats()
	}
}

// BenchmarkGCPreallocatedMap benchmarks pre-allocated map capacity
func BenchmarkGCPreallocatedMap(b *testing.B) {
	// Small cache to test pre-allocation effectiveness
	cfg := &config.EngineConfig{
		MaxSize: 100,
	}
	cache := scache.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%100)
		cache.SetString(key, "value", time.Minute)
	}
}

// BenchmarkGCStressTest stress test with high GC pressure
func BenchmarkGCStressTest(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(10)
		for w := 0; w < 10; w++ {
			go func(workerID, iter int) {
				defer wg.Done()
				key := fmt.Sprintf("key-%d-%d", workerID, iter%100)
				cache.SetString(key, "value", time.Minute)
				cache.GetString(key)
				if iter%5 == 0 {
					cache.Delete(key)
				}
			}(w, i)
		}
		wg.Wait()
	}
}

// BenchmarkGCMemoryUsage measures memory allocation efficiency
func BenchmarkGCMemoryUsage(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		value := fmt.Sprintf("value-%d", i)
		cache.SetString(key, value, time.Minute)
		if i%100 == 0 {
			cache.Delete(key)
		}
	}
	b.ReportAllocs()
}

// BenchmarkGCConcurrentReadWrite benchmarks concurrent read/write with pool operations
func BenchmarkGCConcurrentReadWrite(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				key := fmt.Sprintf("key-%d", i%5000)
				cache.GetString(key)
			} else {
				key := fmt.Sprintf("key-%d", i)
				cache.SetString(key, "value", time.Minute)
			}
			i++
		}
	})
}

// BenchmarkGCLongRunning long-running benchmark to test GC over time
func BenchmarkGCLongRunning(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		cache.SetString(key, "value", time.Minute)

		// Periodic cleanup
		if i%1000 == 0 {
			for j := 0; j < 100; j++ {
				deleteKey := fmt.Sprintf("key-%d", (i+j)%10000)
				cache.Delete(deleteKey)
			}
		}
	}
}
