package main

import (
	"fmt"
	"log"
	"time"

	"github.com/scache"
)

func main() {
	fmt.Println("=== SCache 传统方式 vs 全局缓存方式对比 ===")
	fmt.Println()

	// 传统方式 - 创建实例
	fmt.Println("1. 传统方式 - 创建缓存实例")
	fmt.Println("----------------------------------")

	// 创建不同类型的缓存实例
	lruCache := scache.NewLRU(3)
	lfuCache := scache.NewLFU(3)
	fifoCache := scache.NewFIFO(3)

	// 使用传统缓存
	if err := lruCache.Set("key1", "value1"); err != nil {
		log.Printf("LRU Set error: %v", err)
	}
	if err := lfuCache.Set("key2", "value2"); err != nil {
		log.Printf("LFU Set error: %v", err)
	}
	if err := fifoCache.Set("key3", "value3"); err != nil {
		log.Printf("FIFO Set error: %v", err)
	}

	fmt.Printf("LRU 缓存键: %v\n", lruCache.Keys())
	fmt.Printf("LFU 缓存键: %v\n", lfuCache.Keys())
	fmt.Printf("FIFO 缓存键: %v\n", fifoCache.Keys())

	// 关闭传统缓存
	if err := lruCache.Close(); err != nil {
		log.Printf("LRU cache close error: %v", err)
	}
	if err := lfuCache.Close(); err != nil {
		log.Printf("LFU cache close error: %v", err)
	}
	if err := fifoCache.Close(); err != nil {
		log.Printf("FIFO cache close error: %v", err)
	}

	// 全局缓存方式 - 注册机制
	fmt.Println("\n2. 全局缓存方式 - 注册机制")
	fmt.Println("----------------------------------")

	// 注册全局缓存
	if err := scache.RegisterLRU("global-lru", 3); err != nil {
		log.Fatal("注册 LRU 缓存失败:", err)
	}
	if err := scache.RegisterLFU("global-lfu", 3); err != nil {
		log.Fatal("注册 LFU 缓存失败:", err)
	}
	if err := scache.RegisterFIFO("global-fifo", 3); err != nil {
		log.Fatal("注册 FIFO 缓存失败:", err)
	}

	// 使用全局缓存
	globalLRU, err := scache.Get("global-lru")
	if err != nil {
		log.Fatal("获取 LRU 缓存失败:", err)
	}

	globalLFU, err := scache.Get("global-lfu")
	if err != nil {
		log.Fatal("获取 LFU 缓存失败:", err)
	}

	globalFIFO, err := scache.Get("global-fifo")
	if err != nil {
		log.Fatal("获取 FIFO 缓存失败:", err)
	}

	if err := globalLRU.Set("key1", "value1"); err != nil {
		log.Printf("Global LRU Set error: %v", err)
	}
	if err := globalLFU.Set("key2", "value2"); err != nil {
		log.Printf("Global LFU Set error: %v", err)
	}
	if err := globalFIFO.Set("key3", "value3"); err != nil {
		log.Printf("Global FIFO Set error: %v", err)
	}

	fmt.Printf("全局 LRU 缓存键: %v\n", globalLRU.Keys())
	fmt.Printf("全局 LFU 缓存键: %v\n", globalLFU.Keys())
	fmt.Printf("全局 FIFO 缓存键: %v\n", globalFIFO.Keys())

	// 默认缓存使用
	fmt.Println("\n3. 默认全局缓存使用")
	fmt.Println("----------------------------------")

	// 设置默认缓存值
	if err := scache.Set("app:name", "MyApp"); err != nil {
		log.Printf("设置默认缓存失败: %v", err)
	}

	if err := scache.SetWithTTL("app:session", "session123", 2*time.Second); err != nil {
		log.Printf("设置带过期时间的默认缓存失败: %v", err)
	}

	// 获取默认缓存值
	if value, exists := scache.GetFromDefault("app:name"); exists {
		fmt.Printf("应用名称: %v\n", value)
	}

	if value, exists := scache.GetFromDefault("app:session"); exists {
		fmt.Printf("会话ID: %v\n", value)
	}

	// 等待过期
	fmt.Println("等待会话过期...")
	time.Sleep(3 * time.Second)

	if value, exists := scache.GetFromDefault("app:session"); exists {
		fmt.Printf("会话ID (应该已过期): %v\n", value)
	} else {
		fmt.Println("✓ 会话已过期")
	}

	// 查看所有注册的缓存
	fmt.Println("\n4. 全局缓存管理")
	fmt.Println("----------------------------------")

	caches := scache.List()
	fmt.Printf("所有注册的缓存: %v\n", caches)

	stats := scache.Stats()
	for name, stat := range stats {
		fmt.Printf("%s: 大小=%d, 命中率=%.2f%%\n", name, stat.Size, stat.HitRate*100)
	}

	fmt.Printf("缓存总大小: %d\n", scache.Size())

	// GetOrDefault 演示
	fmt.Println("\n5. GetOrDefault 使用")
	fmt.Println("----------------------------------")

	// 第一次调用会创建新缓存
	tempCache1 := scache.GetOrDefault("temp-cache")
	if err := tempCache1.Set("temp", "value"); err != nil {
		log.Printf("设置临时缓存失败: %v", err)
	}
	fmt.Println("✓ 创建临时缓存并设置值")

	// 第二次调用返回相同的缓存
	tempCache2 := scache.GetOrDefault("temp-cache")
	if value, exists := tempCache2.Get("temp"); exists {
		fmt.Printf("✓ 获取到相同的临时缓存值: %v\n", value)
	}

	// 便捷的默认缓存操作
	fmt.Println("\n6. 便捷的默认缓存操作")
	fmt.Println("----------------------------------")

	// 批量设置
	keys := []string{"user:1", "user:2", "user:3"}
	values := []string{"Alice", "Bob", "Charlie"}

	for i, key := range keys {
		if err := scache.Set(key, values[i]); err != nil {
			log.Printf("设置用户 %s 失败: %v", key, err)
		}
	}

	// 批量获取
	for _, key := range keys {
		if value, exists := scache.GetFromDefault(key); exists {
			fmt.Printf("✓ %s = %v\n", key, value)
		}
	}

	// 检查存在性
	fmt.Printf("user:1 存在: %v\n", scache.ExistsInDefault("user:1"))
	fmt.Printf("nonexistent 存在: %v\n", scache.ExistsInDefault("nonexistent"))

	// 删除操作
	if scache.DeleteFromDefault("user:2") {
		fmt.Println("✓ 删除 user:2")
	}

	// 清理演示
	fmt.Println("\n7. 清理操作")
	fmt.Println("----------------------------------")

	// 清空默认缓存
	if err := scache.ClearDefault(); err != nil {
		log.Printf("清空默认缓存失败: %v", err)
	} else {
		fmt.Println("✓ 清空默认缓存")
	}

	// 移除特定缓存
	if err := scache.Remove("temp-cache"); err != nil {
		log.Printf("移除临时缓存失败: %v", err)
	} else {
		fmt.Println("✓ 移除临时缓存")
	}

	// 最终状态
	fmt.Println("\n8. 最终状态")
	fmt.Println("----------------------------------")

	finalCaches := scache.List()
	fmt.Printf("剩余缓存: %v\n", finalCaches)
	fmt.Printf("剩余缓存总大小: %d\n", scache.Size())

	// 关闭所有全局缓存
	if err := scache.Close(); err != nil {
		log.Printf("关闭所有缓存失败: %v", err)
	} else {
		fmt.Println("✓ 关闭所有全局缓存")
	}

	fmt.Println("\n=== 对比示例完成 ===")
}
