package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/scache/cache"
)

const (
	// 并发协程数量
	numGoroutines = 100
	// 每个协程的操作数量
	operationsPerGoroutine = 1000
	// 键空间大小
	keySpace = 1000
)

// Statistics 用于收集操作统计
type Statistics struct {
	totalSets int64
	totalGets int64
	hits      int64
	misses    int64
	errors    int64
	startTime time.Time
	endTime   time.Time
}

// Worker 表示一个工作协程
type Worker struct {
	id    int
	cache cache.Cache
	stats *Statistics
	wg    *sync.WaitGroup
	rand  *rand.Rand
}

// NewWorker 创建新的工作协程
func NewWorker(id int, c cache.Cache, stats *Statistics, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:    id,
		cache: c,
		stats: stats,
		wg:    wg,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano() + int64(id))),
	}
}

// Run 执行工作协程
func (w *Worker) Run() {
	defer w.wg.Done()

	for i := 0; i < operationsPerGoroutine; i++ {
		// 随机选择操作类型
		operation := w.rand.Intn(100)

		if operation < 60 {
			// 60% 概率执行 GET 操作
			w.doGet()
		} else if operation < 90 {
			// 30% 概率执行 SET 操作
			w.doSet()
		} else {
			// 10% 概率执行 DELETE 操作
			w.doDelete()
		}

		// 随机休息一小段时间，模拟真实场景
		if w.rand.Intn(100) < 10 {
			time.Sleep(time.Duration(w.rand.Intn(5)) * time.Microsecond)
		}
	}
}

// doGet 执行获取操作
func (w *Worker) doGet() {
	key := fmt.Sprintf("key:%d", w.rand.Intn(keySpace))

	_, exists := w.cache.Get(key)
	atomic.AddInt64(&w.stats.totalGets, 1)

	if exists {
		atomic.AddInt64(&w.stats.hits, 1)
	} else {
		atomic.AddInt64(&w.stats.misses, 1)
	}
}

// doSet 执行设置操作
func (w *Worker) doSet() {
	key := fmt.Sprintf("key:%d", w.rand.Intn(keySpace))
	value := fmt.Sprintf("value:%d:%d", w.id, w.rand.Intn(10000))

	err := w.cache.Set(key, value)
	if err != nil {
		atomic.AddInt64(&w.stats.errors, 1)
		log.Printf("Worker %d: SET error: %v", w.id, err)
	}
	atomic.AddInt64(&w.stats.totalSets, 1)
}

// doDelete 执行删除操作
func (w *Worker) doDelete() {
	key := fmt.Sprintf("key:%d", w.rand.Intn(keySpace))

	w.cache.Delete(key)
	// 注意：Delete 操作不会产生错误，即使键不存在
}

// BenchmarkCache 执行缓存基准测试
func BenchmarkCache(cacheType string, c cache.Cache) *Statistics {
	fmt.Printf("\n=== %s 缓存并发测试 ===\n", cacheType)

	stats := &Statistics{
		startTime: time.Now(),
	}

	var wg sync.WaitGroup

	// 创建并启动工作协程
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		worker := NewWorker(i, c, stats, &wg)
		go worker.Run()
	}

	// 等待所有协程完成
	wg.Wait()

	stats.endTime = time.Now()

	// 打印测试结果
	PrintBenchmarkResults(cacheType, stats, c)

	return stats
}

// PrintBenchmarkResults 打印基准测试结果
func PrintBenchmarkResults(cacheType string, stats *Statistics, c cache.Cache) {
	duration := stats.endTime.Sub(stats.startTime)

	fmt.Printf("测试类型: %s\n", cacheType)
	fmt.Printf("测试时长: %v\n", duration)
	fmt.Printf("并发协程数: %d\n", numGoroutines)
	fmt.Printf("每个协程操作数: %d\n", operationsPerGoroutine)
	fmt.Printf("总操作数: %d\n", stats.totalSets+stats.totalGets)
	fmt.Printf("SET 操作: %d\n", stats.totalSets)
	fmt.Printf("GET 操作: %d\n", stats.totalGets)
	fmt.Printf("缓存命中: %d\n", stats.hits)
	fmt.Printf("缓存未命中: %d\n", stats.misses)
	fmt.Printf("错误数: %d\n", stats.errors)

	totalOps := stats.totalSets + stats.totalGets
	if totalOps > 0 {
		hitRate := float64(stats.hits) / float64(stats.totalGets) * 100
		opsPerSec := float64(totalOps) / duration.Seconds()

		fmt.Printf("命中率: %.2f%%\n", hitRate)
		fmt.Printf("操作/秒: %.0f\n", opsPerSec)
		fmt.Printf("平均延迟: %v\n", duration/time.Duration(totalOps))
	}

	// 打印缓存内部统计
	cacheStats := c.Stats()
	fmt.Printf("缓存大小: %d/%d\n", cacheStats.Size, cacheStats.MaxSize)
	fmt.Printf("内部命中率: %.2f%%\n", cacheStats.HitRate*100)
}

// StressTest 压力测试
func StressTest() {
	fmt.Println("\n=== 压力测试 - 短时间内大量操作 ===")

	// 创建小容量缓存，更容易触发淘汰策略
	c := cache.NewLRU(100, // 小容量
		cache.WithShards(32), // 更多分片
		cache.WithStatistics(true),
	)
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("Cache close error: %v", err)
		}
	}()

	stats := &Statistics{
		startTime: time.Now(),
	}

	var wg sync.WaitGroup

	// 启动大量协程，每个协程执行少量操作
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("stress:%d", rand.Intn(50)) // 小键空间，增加冲突
				value := fmt.Sprintf("value:%d:%d", id, j)

				if err := c.Set(key, value); err != nil {
					log.Printf("Stress test Set error: %v", err)
				}
				c.Get(key)

				if rand.Intn(10) == 0 {
					c.Delete(key)
				}
			}
		}(i)
	}

	wg.Wait()
	stats.endTime = time.Now()

	duration := stats.endTime.Sub(stats.startTime)
	fmt.Printf("压力测试完成，耗时: %v\n", duration)
	fmt.Printf("最终缓存大小: %d\n", c.Size())
	fmt.Printf("命中率: %.2f%%\n", c.Stats().HitRate*100)
}

// ComparePolicies 比较不同淘汰策略的性能
func ComparePolicies() {
	fmt.Println("\n=== 淘汰策略性能比较 ===")

	policies := []struct {
		name   string
		create func() cache.Cache
	}{
		{
			"LRU",
			func() cache.Cache {
				return cache.NewLRU(1000, cache.WithStatistics(true))
			},
		},
		{
			"LFU",
			func() cache.Cache {
				return cache.NewLFU(1000, cache.WithStatistics(true))
			},
		},
		{
			"FIFO",
			func() cache.Cache {
				return cache.NewFIFO(1000, cache.WithStatistics(true))
			},
		},
	}

	results := make(map[string]*Statistics)

	for _, policy := range policies {
		c := policy.create()
		stats := BenchmarkCache(policy.name, c)
		results[policy.name] = stats
		if err := c.Close(); err != nil {
			log.Printf("Cache close error: %v", err)
		}

		// 在测试之间休息一下
		time.Sleep(1 * time.Second)
	}

	// 比较结果
	fmt.Println("\n=== 性能比较总结 ===")
	for name, stats := range results {
		duration := stats.endTime.Sub(stats.startTime)
		totalOps := stats.totalSets + stats.totalGets
		opsPerSec := float64(totalOps) / duration.Seconds()
		hitRate := float64(stats.hits) / float64(stats.totalGets) * 100

		fmt.Printf("%-6s: %8.0f ops/s, 命中率: %5.1f%%\n",
			name, opsPerSec, hitRate)
	}
}

func main() {
	fmt.Println("=== SCache 高并发性能测试 ===")
	fmt.Printf("测试配置:\n")
	fmt.Printf("  并发协程数: %d\n", numGoroutines)
	fmt.Printf("  每个协程操作数: %d\n", operationsPerGoroutine)
	fmt.Printf("  总操作数: %d\n", numGoroutines*operationsPerGoroutine)
	fmt.Printf("  键空间大小: %d\n", keySpace)

	// 测试不同类型的缓存
	caches := []struct {
		name   string
		create func() cache.Cache
	}{
		{
			"默认配置",
			func() cache.Cache {
				return cache.New(cache.WithStatistics(true))
			},
		},
		{
			"高性能配置",
			func() cache.Cache {
				return cache.New(
					cache.WithShards(64),     // 更多分片
					cache.WithMaxSize(10000), // 更大容量
					cache.WithStatistics(true),
				)
			},
		},
	}

	for _, cacheConfig := range caches {
		c := cacheConfig.create()
		BenchmarkCache(cacheConfig.name, c)
		if err := c.Close(); err != nil {
			log.Printf("Cache close error: %v", err)
		}

		// 休息一下，让系统恢复
		time.Sleep(2 * time.Second)
	}

	// 比较不同淘汰策略
	ComparePolicies()

	// 压力测试
	StressTest()

	fmt.Println("\n=== 测试完成 ===")
}
