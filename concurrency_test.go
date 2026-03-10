package scache

import (
	"sync"
	"testing"
	"time"

	"github.com/scache-io/scache/config"
)

// TestConcurrentExpiredDeletion 测试并发删除过期键的竞态条件
func TestConcurrentExpiredDeletion(t *testing.T) {
	Flush()

	// 设置一个立即过期的键
	SetString("expire_key", "value", time.Millisecond)

	// 等待过期
	time.Sleep(10 * time.Millisecond)

	// 并发读取过期键，验证不会出现竞态条件
	var wg sync.WaitGroup
	errors := make(chan error, 100)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				GetString("expire_key")
				Exists("expire_key")
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}
}

// TestMaxSizeZeroDisablesEviction 测试 MaxSize=0 时禁用淘汰
func TestMaxSizeZeroDisablesEviction(t *testing.T) {
	Flush()

	// 创建 MaxSize=0 的局部缓存（无限制）
	testConfig := &config.EngineConfig{
		MaxSize:                   0, // 无限制
		MemoryThreshold:           0.9,
		DefaultExpiration:         0,
		BackgroundCleanupInterval: time.Minute,
	}
	localCache := New(testConfig)

	// 插入超过原默认容量(100)的数据
	for i := 0; i < 200; i++ {
		err := localCache.SetString(string(rune(i)), "value", time.Hour)
		if err != nil {
			t.Errorf("Expected no error with MaxSize=0, got: %v", err)
		}
	}

	// 验证所有数据都存在
	if localCache.Size() != 200 {
		t.Errorf("Expected 200 items with MaxSize=0, got %d", localCache.Size())
	}
}
