package tests

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/scache-io/scache"
	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/storage"
	"github.com/scache-io/scache/types"
)

// BenchmarkSetOperation SET操作性能测试
func BenchmarkSetOperation(b *testing.B) {
	engine := storage.NewStorageEngine(nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench_key_%d", i)
			value := fmt.Sprintf("bench_value_%d", i)
			obj := types.NewStringObject(value, time.Hour)
			engine.Set(key, obj)
			i++
		}
	})
}

// BenchmarkGetOperation GET操作性能测试
func BenchmarkGetOperation(b *testing.B) {
	engine := storage.NewStorageEngine(nil)

	// 预填充数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)
		obj := types.NewStringObject(value, time.Hour)
		engine.Set(key, obj)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench_key_%d", i%1000)
			engine.Get(key)
			i++
		}
	})
}

// BenchmarkSetGetOperation SET+GET组合操作性能测试
func BenchmarkSetGetOperation(b *testing.B) {
	engine := storage.NewStorageEngine(nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench_key_%d", i)
			value := fmt.Sprintf("bench_value_%d", i)

			// SET
			obj := types.NewStringObject(value, time.Hour)
			engine.Set(key, obj)

			// GET
			engine.Get(key)
			i++
		}
	})
}

// BenchmarkConcurrentOperations 并发操作性能测试
func BenchmarkConcurrentOperations(b *testing.B) {
	engine := storage.NewStorageEngine(nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%4 == 0 {
				// SET 操作
				key := fmt.Sprintf("bench_key_%d", i)
				value := fmt.Sprintf("bench_value_%d", i)
				obj := types.NewStringObject(value, time.Hour)
				engine.Set(key, obj)
			} else if i%4 == 1 {
				// GET 操作
				key := fmt.Sprintf("bench_key_%d", i%100)
				engine.Get(key)
			} else if i%4 == 2 {
				// EXISTS 操作
				key := fmt.Sprintf("bench_key_%d", i%100)
				engine.Exists(key)
			} else {
				// DELETE 操作
				key := fmt.Sprintf("bench_key_%d", i%50)
				engine.Delete(key)
			}
			i++
		}
	})
}

// BenchmarkHashOperations 哈希操作性能测试
func BenchmarkHashOperations(b *testing.B) {
	engine := storage.NewStorageEngine(nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("hash_%d", i%100)
			field := fmt.Sprintf("field_%d", i%10)
			value := fmt.Sprintf("value_%d", i)

			// HSET
			obj := types.NewHashObject(map[string]interface{}{
				field: value,
			}, time.Hour)
			engine.Set(key, obj)

			// HGET
			if retrievedObj, exists := engine.Get(key); exists {
				if hashObj, ok := retrievedObj.(*types.HashObject); ok {
					hashObj.Get(field)
				}
			}
			i++
		}
	})
}

// BenchmarkListOperations 列表操作性能测试
func BenchmarkListOperations(b *testing.B) {
	engine := storage.NewStorageEngine(nil)

	// 预填充一些列表
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("list_%d", i)
		values := make([]interface{}, 10)
		for j := 0; j < 10; j++ {
			values[j] = fmt.Sprintf("item_%d", j)
		}
		obj := types.NewListObject(values, time.Hour)
		engine.Set(key, obj)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("list_%d", i%100)
			value := fmt.Sprintf("new_item_%d", i)

			if retrievedObj, exists := engine.Get(key); exists {
				if listObj, ok := retrievedObj.(*types.ListObject); ok {
					// Push
					listObj.Push(value)

					// Pop
					listObj.Pop()

					// Index
					listObj.Index(0)
				}
			}
			i++
		}
	})
}

// BenchmarkExecutorSet 执行器SET操作性能测试
func BenchmarkExecutorSet(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("exec_key_%d", i)
			value := fmt.Sprintf("exec_value_%d", i)
			executor.Execute("SET", key, value, time.Hour)
			i++
		}
	})
}

// BenchmarkExecutorGet 执行器GET操作性能测试
func BenchmarkExecutorGet(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	// 预填充数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("exec_key_%d", i)
		value := fmt.Sprintf("exec_value_%d", i)
		executor.Execute("SET", key, value, time.Hour)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("exec_key_%d", i%1000)
			executor.Execute("GET", key)
			i++
		}
	})
}

// BenchmarkConvenienceAPI 便捷API性能测试
func BenchmarkConvenienceAPI(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("conv_key_%d", i)
			value := fmt.Sprintf("conv_value_%d", i)

			// SET
			scache.Set(key, value, time.Hour)

			// GET
			scache.Get(key)
			i++
		}
	})
}

// TestMemoryUsage 内存使用测试
func TestMemoryUsage(t *testing.T) {
	engine := storage.NewStorageEngine(&storage.EngineConfig{
		MaxSize: 10000,
	})

	// 记录初始内存
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// 添加大量数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("mem_test_key_%d", i)
		value := make([]byte, 1024) // 1KB per value
		obj := types.NewStringObject(string(value), time.Hour)
		engine.Set(key, obj)
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	memoryUsed := m2.Alloc - m1.Alloc
	expectedMemory := int64(1000 * 1024) // 约1MB

	// 允许一定的内存开销（最多3倍）
	if int64(memoryUsed) > expectedMemory*3 {
		t.Errorf("Memory usage too high: used %d bytes, expected around %d bytes", memoryUsed, expectedMemory)
	}

	t.Logf("Memory used for 1000 items: %d bytes", memoryUsed)
}

// TestHighConcurrency 高并发测试
func TestHighConcurrency(t *testing.T) {
	engine := storage.NewStorageEngine(&storage.EngineConfig{
		MaxSize: 10000,
	})
	executor := scache.NewExecutor(engine)

	concurrency := 100
	operations := 1000
	var wg sync.WaitGroup
	errors := make(chan error, concurrency)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("worker_%d_key_%d", workerID, j)
				value := fmt.Sprintf("worker_%d_value_%d", workerID, j)

				// SET
				_, err := executor.Execute("SET", key, value, time.Minute)
				if err != nil {
					errors <- fmt.Errorf("SET failed: %w", err)
					return
				}

				// GET
				_, err = executor.Execute("GET", key)
				if err != nil {
					errors <- fmt.Errorf("GET failed: %w", err)
					return
				}

				// EXISTS
				_, err = executor.Execute("EXISTS", key)
				if err != nil {
					errors <- fmt.Errorf("EXISTS failed: %w", err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)
	totalOps := concurrency * operations * 3 // 每个循环3个操作
	qps := float64(totalOps) / duration.Seconds()

	// 检查错误
	for err := range errors {
		t.Error(err)
	}

	t.Logf("High concurrency test completed:")
	t.Logf("  Concurrency: %d", concurrency)
	t.Logf("  Operations per worker: %d", operations)
	t.Logf("  Total operations: %d", totalOps)
	t.Logf("  Duration: %v", duration)
	t.Logf("  QPS: %.2f", qps)

	// 验证最终状态
	stats := engine.Stats().(map[string]interface{})
	if hits, ok := stats["hits"].(int64); ok {
		t.Logf("  Hits: %d", hits)
	}
	if misses, ok := stats["misses"].(int64); ok {
		t.Logf("  Misses: %d", misses)
	}
	if keys, ok := stats["keys"].(int); ok {
		t.Logf("  Keys: %d", keys)
	}
}

// TestThroughputThroughput 吞吐量测试
func TestThroughput(t *testing.T) {
	configs := []struct {
		name string
		opts []config.EngineOption
	}{
		{"Small", config.SmallConfig},
		{"Medium", config.MediumConfig},
		{"Large", config.LargeConfig},
	}

	for _, tc := range configs {
		t.Run(tc.name, func(t *testing.T) {
			engine := storage.NewStorageEngine(nil)
			for _, opt := range tc.opts {
				opt(engine.(*storage.StorageEngine).GetConfig())
			}

			executor := scache.NewExecutor(engine)

			operations := 10000
			start := time.Now()

			// SET operations
			for i := 0; i < operations; i++ {
				key := fmt.Sprintf("throughput_key_%d", i)
				value := fmt.Sprintf("throughput_value_%d", i)
				_, err := executor.Execute("SET", key, value, time.Minute)
				if err != nil {
					t.Fatalf("SET failed: %v", err)
				}
			}

			setDuration := time.Since(start)
			setQPS := float64(operations) / setDuration.Seconds()

			// GET operations
			start = time.Now()
			for i := 0; i < operations; i++ {
				key := fmt.Sprintf("throughput_key_%d", i)
				_, err := executor.Execute("GET", key)
				if err != nil {
					t.Fatalf("GET failed: %v", err)
				}
			}

			getDuration := time.Since(start)
			getQPS := float64(operations) / getDuration.Seconds()

			t.Logf("Throughput test - %s config:", tc.name)
			t.Logf("  SET: %d ops in %v (%.2f QPS)", operations, setDuration, setQPS)
			t.Logf("  GET: %d ops in %v (%.2f QPS)", operations, getDuration, getQPS)

			// 性能断言
			if setQPS < 100000 {
				t.Errorf("SET QPS too low: %.2f (expected >= 100000)", setQPS)
			}
			if getQPS < 200000 {
				t.Errorf("GET QPS too low: %.2f (expected >= 200000)", getQPS)
			}
		})
	}
}

// TestLatency 延迟测试
func TestLatency(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := scache.NewExecutor(engine)

	operations := 10000
	setLatencies := make([]time.Duration, operations)
	getLatencies := make([]time.Duration, operations)

	// SET latency test
	for i := 0; i < operations; i++ {
		key := fmt.Sprintf("latency_key_%d", i)
		value := fmt.Sprintf("latency_value_%d", i)

		start := time.Now()
		_, err := executor.Execute("SET", key, value, time.Minute)
		setLatencies[i] = time.Since(start)

		if err != nil {
			t.Fatalf("SET failed: %v", err)
		}
	}

	// GET latency test
	for i := 0; i < operations; i++ {
		key := fmt.Sprintf("latency_key_%d", i)

		start := time.Now()
		_, err := executor.Execute("GET", key)
		getLatencies[i] = time.Since(start)

		if err != nil {
			t.Fatalf("GET failed: %v", err)
		}
	}

	// 计算统计数据
	calculateLatencyStats := func(latencies []time.Duration, op string) {
		var sum time.Duration
		max := latencies[0]
		min := latencies[0]

		for _, lat := range latencies {
			sum += lat
			if lat > max {
				max = lat
			}
			if lat < min {
				min = lat
			}
		}

		avg := sum / time.Duration(len(latencies))

		// 计算 P99
		sorted := make([]time.Duration, len(latencies))
		copy(sorted, latencies)

		// 简单排序（冒泡排序，仅用于测试）
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i] > sorted[j] {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}

		p99Index := int(float64(len(sorted)) * 0.99)
		p99 := sorted[p99Index]

		t.Logf("%s latency stats:", op)
		t.Logf("  Average: %v", avg)
		t.Logf("  Min: %v", min)
		t.Logf("  Max: %v", max)
		t.Logf("  P99: %v", p99)

		// 延迟断言
		if avg > time.Microsecond*100 {
			t.Errorf("%s average latency too high: %v (expected <= 100μs)", op, avg)
		}
		if p99 > time.Microsecond*500 {
			t.Errorf("%s P99 latency too high: %v (expected <= 500μs)", op, p99)
		}
	}

	calculateLatencyStats(setLatencies, "SET")
	calculateLatencyStats(getLatencies, "GET")
}
