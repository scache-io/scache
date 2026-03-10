package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/scache-io/scache"
	"github.com/scache-io/scache/config"
)

// ==================== Basic Operations Benchmarks ====================

func BenchmarkStoreString(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}
}

func BenchmarkLoadString(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate cache
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		cache.GetString(key)
	}
}

func BenchmarkDeleteString(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
		cache.Delete(key)
	}
}

// ==================== Struct Operations Benchmarks ====================

type BenchmarkUser struct {
	ID   int
	Name string
	Age  int
}

func BenchmarkStoreStruct(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("user-%d", i)
		user := BenchmarkUser{ID: i, Name: "User", Age: 25}
		cache.Store(key, user, time.Minute)
	}
}

func BenchmarkLoadStruct(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate cache
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("user-%d", i)
		user := BenchmarkUser{ID: i, Name: "User", Age: 25}
		cache.Store(key, user, time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("user-%d", i%10000)
		var user BenchmarkUser
		cache.Load(key, &user)
	}
}

// ==================== Concurrent Operations Benchmarks ====================

func BenchmarkConcurrentStore(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i)
			cache.SetString(key, "value", time.Minute)
			i++
		}
	})
}

func BenchmarkConcurrentLoad(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate cache
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%10000)
			cache.GetString(key)
			i++
		}
	})
}

func BenchmarkConcurrentReadWrite(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				key := fmt.Sprintf("key-%d", i%1000)
				cache.GetString(key)
			} else {
				key := fmt.Sprintf("key-%d", i)
				cache.SetString(key, "value", time.Minute)
			}
			i++
		}
	})
}

// ==================== LRU Eviction Benchmarks ====================

func BenchmarkLRUEviction(b *testing.B) {
	cfg := &config.EngineConfig{
		MaxSize: 1000,
	}
	cache := scache.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}
}

func BenchmarkLRUEvictionWithLoad(b *testing.B) {
	cfg := &config.EngineConfig{
		MaxSize: 1000,
	}
	cache := scache.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)

		// Randomly load some keys
		if i%10 == 0 {
			cache.GetString(fmt.Sprintf("key-%d", i%500))
		}
	}
}

// ==================== TTL Expiration Benchmarks ====================

func BenchmarkTTLExpiration(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		// Very short TTL to trigger expiration
		cache.SetString(key, "value", 10*time.Millisecond)
	}
}

func BenchmarkTTLCheck(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate with various TTLs
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		ttl := time.Duration(i%10+1) * time.Second
		cache.SetString(key, "value", ttl)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		cache.GetString(key)
	}
}

// ==================== Data Type Benchmarks ====================

func BenchmarkStoreList(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	list := []interface{}{"item1", "item2", "item3", "item4", "item5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("list-%d", i)
		cache.SetList(key, list, time.Minute)
	}
}

func BenchmarkLoadList(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	list := []interface{}{"item1", "item2", "item3", "item4", "item5"}
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("list-%d", i)
		cache.SetList(key, list, time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("list-%d", i%10000)
		cache.GetList(key)
	}
}

func BenchmarkStoreHash(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	hash := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
		"field4": "value4",
		"field5": "value5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("hash-%d", i)
		cache.SetHash(key, hash, time.Minute)
	}
}

func BenchmarkLoadHash(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	hash := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
		"field4": "value4",
		"field5": "value5",
	}
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("hash-%d", i)
		cache.SetHash(key, hash, time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("hash-%d", i%10000)
		cache.GetHash(key)
	}
}

// ==================== Large Data Benchmarks ====================

func BenchmarkStoreLargeStruct(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	type LargeStruct struct {
		ID      int
		Name    string
		Email   string
		Address string
		Phone   string
		Data    []byte
	}

	largeData := make([]byte, 1024) // 1KB data

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("large-%d", i)
		obj := LargeStruct{
			ID:      i,
			Name:    "User Name",
			Email:   "user@example.com",
			Address: "123 Main St",
			Phone:   "555-1234",
			Data:    largeData,
		}
		cache.Store(key, obj, time.Minute)
	}
}

func BenchmarkLoadLargeStruct(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	type LargeStruct struct {
		ID      int
		Name    string
		Email   string
		Address string
		Phone   string
		Data    []byte
	}

	largeData := make([]byte, 1024) // 1KB data

	// Pre-populate
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("large-%d", i)
		obj := LargeStruct{
			ID:      i,
			Name:    "User Name",
			Email:   "user@example.com",
			Address: "123 Main St",
			Phone:   "555-1234",
			Data:    largeData,
		}
		cache.Store(key, obj, time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("large-%d", i%1000)
		var obj LargeStruct
		cache.Load(key, &obj)
	}
}

// ==================== Mixed Workload Benchmarks ====================

func BenchmarkMixedWorkload(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	// Pre-populate
	for i := 0; i < 5000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 10 {
		case 0, 1, 2, 3, 4: // 50% reads
			key := fmt.Sprintf("key-%d", i%5000)
			cache.GetString(key)
		case 5, 6: // 20% writes
			key := fmt.Sprintf("key-%d", i)
			cache.SetString(key, "value", time.Minute)
		case 7: // 10% deletes
			key := fmt.Sprintf("key-%d", i%5000)
			cache.Delete(key)
		case 8: // 10% exists checks
			key := fmt.Sprintf("key-%d", i%5000)
			cache.Exists(key)
		case 9: // 10% struct stores
			key := fmt.Sprintf("user-%d", i)
			user := BenchmarkUser{ID: i, Name: "User", Age: 25}
			cache.Store(key, user, time.Minute)
		}
	}
}

// ==================== Capacity Benchmarks ====================

func BenchmarkUnlimitedCapacity(b *testing.B) {
	cfg := &config.EngineConfig{
		MaxSize: 0, // Unlimited
	}
	cache := scache.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}
}

func BenchmarkLargeCapacity(b *testing.B) {
	cfg := &config.EngineConfig{
		MaxSize: 100000, // 100K capacity
	}
	cache := scache.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.SetString(key, "value", time.Minute)
	}
}

// ==================== Stress Tests ====================

func BenchmarkHighConcurrency(b *testing.B) {
	cache := scache.New(config.DefaultEngineConfig())

	var wg sync.WaitGroup
	workers := 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(workers)
		for w := 0; w < workers; w++ {
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					key := fmt.Sprintf("key-%d-%d", workerID, j)
					cache.SetString(key, "value", time.Minute)
					cache.GetString(key)
					cache.Delete(key)
				}
			}(w)
		}
		wg.Wait()
	}
}
