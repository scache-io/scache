package scache

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"scache/config"
	"scache/storage"
	"scache/types"
)

// BenchmarkStorageSet 存储引擎SET性能基准测试
func BenchmarkStorageSet(b *testing.B) {
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

// BenchmarkStorageGet 存储引擎GET性能基准测试
func BenchmarkStorageGet(b *testing.B) {
	engine := storage.NewStorageEngine(nil)

	// 预填充数据
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)
		obj := types.NewStringObject(value, time.Hour)
		engine.Set(key, obj)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench_key_%d", i%10000)
			engine.Get(key)
			i++
		}
	})
}

// BenchmarkStorageSetGet 存储引擎SET+GET组合性能基准测试
func BenchmarkStorageSetGet(b *testing.B) {
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

// BenchmarkExecutorSet 执行器SET性能基准测试
func BenchmarkExecutorSet(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

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

// BenchmarkExecutorGet 执行器GET性能基准测试
func BenchmarkExecutorGet(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

	// 预填充数据
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("exec_key_%d", i)
		value := fmt.Sprintf("exec_value_%d", i)
		executor.Execute("SET", key, value, time.Hour)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("exec_key_%d", i%10000)
			executor.Execute("GET", key)
			i++
		}
	})
}

// BenchmarkConvenienceAPI 便捷API性能基准测试
func BenchmarkConvenienceAPI(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("conv_key_%d", i)
			value := fmt.Sprintf("conv_value_%d", i)

			// SET
			Set(key, value, time.Hour)

			// GET
			Get(key)
			i++
		}
	})
}

// BenchmarkHashOperations 哈希操作性能基准测试
func BenchmarkHashOperations(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("hash_%d", i%1000)
			field := fmt.Sprintf("field_%d", i%100)
			value := fmt.Sprintf("value_%d", i)

			// HSET
			executor.Execute("HSET", key, field, value, time.Hour)

			// HGET
			executor.Execute("HGET", key, field)
			i++
		}
	})
}

// BenchmarkListOperations 列表操作性能基准测试
func BenchmarkListOperations(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

	// 预填充列表
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("list_%d", i)
		for j := 0; j < 10; j++ {
			executor.Execute("LPUSH", key, fmt.Sprintf("item_%d", j), time.Hour)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("list_%d", i%100)
			value := fmt.Sprintf("new_item_%d", i)

			// LPUSH
			executor.Execute("LPUSH", key, value, time.Hour)

			// RPOP
			executor.Execute("RPOP", key)
			i++
		}
	})
}

// BenchmarkTTLOperations TTL操作性能基准测试
func BenchmarkTTLOperations(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

	// 预填充数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("ttl_key_%d", i)
		executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Hour)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("ttl_key_%d", i%1000)

			// TTL
			executor.Execute("TTL", key)

			// EXPIRE
			executor.Execute("EXPIRE", key, time.Minute*30)
			i++
		}
	})
}

// BenchmarkMixedOperations 混合操作性能基准测试
func BenchmarkMixedOperations(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 8 {
			case 0:
				// SET
				key := fmt.Sprintf("mixed_key_%d", i)
				executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Hour)
			case 1:
				// GET
				key := fmt.Sprintf("mixed_key_%d", i%1000)
				executor.Execute("GET", key)
			case 2:
				// HSET
				key := fmt.Sprintf("mixed_hash_%d", i%100)
				executor.Execute("HSET", key, "field", fmt.Sprintf("value_%d", i), time.Hour)
			case 3:
				// HGET
				key := fmt.Sprintf("mixed_hash_%d", i%100)
				executor.Execute("HGET", key, "field")
			case 4:
				// LPUSH
				key := fmt.Sprintf("mixed_list_%d", i%100)
				executor.Execute("LPUSH", key, fmt.Sprintf("item_%d", i), time.Hour)
			case 5:
				// RPOP
				key := fmt.Sprintf("mixed_list_%d", i%100)
				executor.Execute("RPOP", key)
			case 6:
				// EXISTS
				key := fmt.Sprintf("mixed_key_%d", i%1000)
				executor.Execute("EXISTS", key)
			case 7:
				// TYPE
				key := fmt.Sprintf("mixed_key_%d", i%1000)
				executor.Execute("TYPE", key)
			}
			i++
		}
	})
}

// BenchmarkConcurrentAccess 并发访问性能基准测试
func BenchmarkConcurrentAccess(b *testing.B) {
	engine := storage.NewStorageEngine(&storage.EngineConfig{
		MaxSize: 10000,
	})
	executor := NewExecutor(engine)

	// 预填充一些数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("concurrent_key_%d", i)
		executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Hour)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("concurrent_key_%d", i%1000)
			executor.Execute("GET", key)
			i++
		}
	})
}

// BenchmarkMemoryUsage 内存使用基准测试
func BenchmarkMemoryUsage(b *testing.B) {
	engine := storage.NewStorageEngine(&storage.EngineConfig{
		MaxSize: 10000,
	})

	// 记录初始内存
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("mem_key_%d", i)
		value := make([]byte, 1024) // 1KB value
		obj := types.NewStringObject(string(value), time.Hour)
		engine.Set(key, obj)
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	memoryUsed := m2.Alloc - m1.Alloc
	avgMemoryPerItem := int64(memoryUsed) / int64(b.N)

	b.ReportMetric(float64(avgMemoryPerItem), "bytes/item")
	b.Logf("Memory used: %d bytes for %d items, avg: %d bytes/item",
		memoryUsed, b.N, avgMemoryPerItem)
}

// BenchmarkDifferentConfigs 不同配置性能对比
func BenchmarkDifferentConfigs(b *testing.B) {
	configs := []struct {
		name string
		opts []config.EngineOption
	}{
		{"Small", config.SmallConfig},
		{"Medium", config.MediumConfig},
		{"Large", config.LargeConfig},
	}

	for _, tc := range configs {
		b.Run(tc.name, func(b *testing.B) {
			config := &storage.EngineConfig{}
			for _, opt := range tc.opts {
				opt(config)
			}
			engine := storage.NewStorageEngine(config)
			executor := NewExecutor(engine)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := fmt.Sprintf("config_key_%d", i)
					value := fmt.Sprintf("config_value_%d", i)
					executor.Execute("SET", key, value, time.Minute)
					executor.Execute("GET", key)
					i++
				}
			})
		})
	}
}

// BenchmarkScalability 可扩展性基准测试
func BenchmarkScalability(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			engine := storage.NewStorageEngine(&storage.EngineConfig{
				MaxSize: size,
			})
			executor := NewExecutor(engine)

			// 预填充到指定大小
			for i := 0; i < size; i++ {
				key := fmt.Sprintf("scale_key_%d", i)
				executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Hour)
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := fmt.Sprintf("scale_key_%d", i%size)
					executor.Execute("GET", key)
					i++
				}
			})
		})
	}
}

// BenchmarkHighLoad 高负载基准测试
func BenchmarkHighLoad(b *testing.B) {
	engine := storage.NewStorageEngine(&storage.EngineConfig{
		MaxSize: 100000,
	})
	executor := NewExecutor(engine)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 模拟高负载场景
			switch i % 10 {
			case 0, 1, 2:
				// 30% 写操作
				key := fmt.Sprintf("load_key_%d", i)
				executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Minute)
			case 3, 4, 5, 6:
				// 40% 读操作
				key := fmt.Sprintf("load_key_%d", i%10000)
				executor.Execute("GET", key)
			case 7:
				// 10% 删除操作
				key := fmt.Sprintf("load_key_%d", i%5000)
				executor.Execute("DEL", key)
			case 8:
				// 10% 存在性检查
				key := fmt.Sprintf("load_key_%d", i%10000)
				executor.Execute("EXISTS", key)
			case 9:
				// 10% 类型检查
				key := fmt.Sprintf("load_key_%d", i%10000)
				executor.Execute("TYPE", key)
			}
			i++
		}
	})
}

// BenchmarkLatency 延迟基准测试
func BenchmarkLatency(b *testing.B) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

	// 预填充数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("latency_key_%d", i)
		executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Hour)
	}

	var totalLatency int64
	var operations int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("latency_key_%d", i%1000)

			start := time.Now()
			executor.Execute("GET", key)
			latency := time.Since(start).Nanoseconds()

			totalLatency += latency
			operations++
			i++
		}
	})

	avgLatency := totalLatency / operations
	b.ReportMetric(float64(avgLatency), "ns/op")
}

// TestPerformanceRegression 性能回归测试
func TestPerformanceRegression(t *testing.T) {
	engine := storage.NewStorageEngine(nil)
	executor := NewExecutor(engine)

	operations := 10000

	// SET 性能测试
	start := time.Now()
	for i := 0; i < operations; i++ {
		key := fmt.Sprintf("perf_key_%d", i)
		value := fmt.Sprintf("perf_value_%d", i)
		executor.Execute("SET", key, value, time.Hour)
	}
	setDuration := time.Since(start)
	setQPS := float64(operations) / setDuration.Seconds()

	// GET 性能测试
	start = time.Now()
	for i := 0; i < operations; i++ {
		key := fmt.Sprintf("perf_key_%d", i)
		executor.Execute("GET", key)
	}
	getDuration := time.Since(start)
	getQPS := float64(operations) / getDuration.Seconds()

	t.Logf("Performance Test Results:")
	t.Logf("  SET: %d ops in %v (%.2f QPS)", operations, setDuration, setQPS)
	t.Logf("  GET: %d ops in %v (%.2f QPS)", operations, getDuration, getQPS)

	// 性能基准断言
	if setQPS < 100000 {
		t.Errorf("SET performance regression: %.2f QPS < 100000 QPS", setQPS)
	}
	if getQPS < 200000 {
		t.Errorf("GET performance regression: %.2f QPS < 200000 QPS", getQPS)
	}
}

// BenchmarkConcurrencyScaling 并发扩展性测试
func BenchmarkConcurrencyScaling(b *testing.B) {
	concurrencies := []int{1, 2, 4, 8, 16, 32}

	for _, concurrency := range concurrencies {
		b.Run(fmt.Sprintf("Goroutines_%d", concurrency), func(b *testing.B) {
			engine := storage.NewStorageEngine(nil)
			executor := NewExecutor(engine)

			// 预填充数据
			for i := 0; i < 10000; i++ {
				key := fmt.Sprintf("scale_key_%d", i)
				executor.Execute("SET", key, fmt.Sprintf("value_%d", i), time.Hour)
			}

			b.ResetTimer()
			b.SetParallelism(concurrency)
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := fmt.Sprintf("scale_key_%d", i%10000)
					executor.Execute("GET", key)
					i++
				}
			})
		})
	}
}
