package tests

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/scache-io/scache/cache"
)

// BenchmarkTPSRead 测试读操作 TPS
func BenchmarkTPSRead(b *testing.B) {
	c := cache.NewLocalCache(nil)
	
	// 预热数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		c.SetString(key, "value")
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key:%d", i%1000)
			c.GetString(key)
			i++
		}
	})
}

// BenchmarkTPSWrite 测试写操作 TPS
func BenchmarkTPSWrite(b *testing.B) {
	c := cache.NewLocalCache(nil)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key:%d", i%1000)
			c.SetString(key, "value")
			i++
		}
	})
}

// BenchmarkTPSMixed 测试混合读写 TPS (80% read, 20% write)
func BenchmarkTPSMixed(b *testing.B) {
	c := cache.NewLocalCache(nil)
	
	// 预热数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		c.SetString(key, "value")
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key:%d", i%1000)
			if i%5 == 0 {
				c.SetString(key, "newvalue")
			} else {
				c.GetString(key)
			}
			i++
		}
	})
}

// TestTPSMeasurement 测试 TPS 测量
func TestTPSMeasurement(t *testing.T) {
	c := cache.NewLocalCache(nil)
	
	// 预热数据
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key:%d", i)
		c.SetString(key, fmt.Sprintf("value:%d", i))
	}
	
	// 测试参数
	workers := 10
	duration := 3 * time.Second
	
	// 读 TPS 测试
	readOps := int64(0)
	var wg sync.WaitGroup
	start := time.Now()
	
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			i := 0
			for time.Since(start) < duration {
				key := fmt.Sprintf("key:%d", (workerID*1000+i)%10000)
				c.GetString(key)
				atomic.AddInt64(&readOps, 1)
				i++
			}
		}(w)
	}
	wg.Wait()
	
	readTPS := float64(readOps) / duration.Seconds()
	t.Logf("Read TPS: %.0f ops/sec (%d workers)", readTPS, workers)
	
	// 写 TPS 测试
	writeOps := int64(0)
	start = time.Now()
	
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			i := 0
			for time.Since(start) < duration {
				key := fmt.Sprintf("key:%d", (workerID*1000+i)%10000)
				c.SetString(key, fmt.Sprintf("newvalue:%d", i))
				atomic.AddInt64(&writeOps, 1)
				i++
			}
		}(w)
	}
	wg.Wait()
	
	writeTPS := float64(writeOps) / duration.Seconds()
	t.Logf("Write TPS: %.0f ops/sec (%d workers)", writeTPS, workers)
	
	// 混合读写 TPS 测试 (80% read, 20% write)
	mixedOps := int64(0)
	start = time.Now()
	
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			i := 0
			for time.Since(start) < duration {
				key := fmt.Sprintf("key:%d", (workerID*1000+i)%10000)
				if i%5 == 0 {
					c.SetString(key, fmt.Sprintf("newvalue:%d", i))
				} else {
					c.GetString(key)
				}
				atomic.AddInt64(&mixedOps, 1)
				i++
			}
		}(w)
	}
	wg.Wait()
	
	mixedTPS := float64(mixedOps) / duration.Seconds()
	t.Logf("Mixed TPS (80%% read, 20%% write): %.0f ops/sec (%d workers)", mixedTPS, workers)
	
	// 单线程 TPS 测试
	singleOps := 0
	start = time.Now()
	for time.Since(start) < duration {
		key := fmt.Sprintf("key:%d", singleOps%10000)
		c.GetString(key)
		singleOps++
	}
	
	singleTPS := float64(singleOps) / duration.Seconds()
	t.Logf("Single-thread Read TPS: %.0f ops/sec", singleTPS)
}

// TestTPSWithDifferentWorkers 测试不同并发数下的 TPS
func TestTPSWithDifferentWorkers(t *testing.T) {
	c := cache.NewLocalCache(nil)
	
	// 预热数据
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key:%d", i)
		c.SetString(key, fmt.Sprintf("value:%d", i))
	}
	
	duration := 2 * time.Second
	workerCounts := []int{1, 2, 4, 8, 16, 32}
	
	t.Log("\n=== Read TPS with Different Workers ===")
	for _, workers := range workerCounts {
		ops := int64(0)
		var wg sync.WaitGroup
		start := time.Now()
		
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				i := 0
				for time.Since(start) < duration {
					key := fmt.Sprintf("key:%d", (workerID*1000+i)%10000)
					c.GetString(key)
					atomic.AddInt64(&ops, 1)
					i++
				}
			}(w)
		}
		wg.Wait()
		
		tps := float64(ops) / duration.Seconds()
		t.Logf("%2d workers: %12.0f ops/sec (%.0f ops/sec per worker)", workers, tps, tps/float64(workers))
	}
	
	t.Log("\n=== Write TPS with Different Workers ===")
	for _, workers := range workerCounts {
		ops := int64(0)
		var wg sync.WaitGroup
		start := time.Now()
		
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				i := 0
				for time.Since(start) < duration {
					key := fmt.Sprintf("key:%d", (workerID*1000+i)%10000)
					c.SetString(key, fmt.Sprintf("newvalue:%d", i))
					atomic.AddInt64(&ops, 1)
					i++
				}
			}(w)
		}
		wg.Wait()
		
		tps := float64(ops) / duration.Seconds()
		t.Logf("%2d workers: %12.0f ops/sec (%.0f ops/sec per worker)", workers, tps, tps/float64(workers))
	}
}
