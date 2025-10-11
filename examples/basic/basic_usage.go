package main

import (
	"fmt"
	"log"
	"time"

	"scache/cache"
)

func main() {
	fmt.Println("=== 基本缓存使用示例 ===")

	// 1. 创建缓存实例
	fmt.Println("\n1. 创建缓存实例")
	c := cache.NewCache(
		cache.WithMaxSize(100),
		cache.WithDefaultExpiration(time.Minute*5),
		cache.WithCleanupInterval(time.Minute),
		cache.WithStats(true),
	)
	fmt.Printf("缓存创建成功，大小: %d\n", c.Size())

	// 2. 基本操作
	fmt.Println("\n2. 基本操作演示")

	// 设置缓存项
	err := c.Set("user:1001", "张三", time.Minute*10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("设置 user:1001 = 张三 (10分钟过期)")

	// 获取缓存项
	if value, found := c.Get("user:1001"); found {
		fmt.Printf("获取 user:1001 = %v\n", value)
	} else {
		fmt.Println("未找到 user:1001")
	}

	// 检查是否存在
	if c.Exists("user:1001") {
		fmt.Println("user:1001 存在于缓存中")
	}

	// 3. 缓存统计
	fmt.Println("\n3. 缓存统计信息")
	stats := c.Stats()
	fmt.Printf("命中次数: %d\n", stats.Hits)
	fmt.Printf("未命中次数: %d\n", stats.Misses)
	fmt.Printf("命中率: %.2f%%\n", stats.HitRate*100)
	fmt.Printf("当前大小: %d\n", stats.Size)
	fmt.Printf("最大容量: %d\n", stats.MaxSize)

	// 4. 批量操作
	fmt.Println("\n4. 批量操作")
	users := map[string]string{
		"user:1002": "李四",
		"user:1003": "王五",
		"user:1004": "赵六",
	}

	for key, value := range users {
		c.Set(key, value, 0) // 使用默认过期时间
		fmt.Printf("设置 %s = %s\n", key, value)
	}

	fmt.Printf("缓存当前大小: %d\n", c.Size())

	// 5. 获取不存在的键
	fmt.Println("\n5. 获取不存在的键")
	if value, found := c.Get("user:9999"); found {
		fmt.Printf("获取 user:9999 = %v\n", value)
	} else {
		fmt.Println("user:9999 不存在于缓存中")
	}

	// 6. 删除操作
	fmt.Println("\n6. 删除操作")
	if c.Delete("user:1002") {
		fmt.Println("成功删除 user:1002")
	}

	if !c.Exists("user:1002") {
		fmt.Println("user:1002 已被删除")
	}

	// 7. 过期测试
	fmt.Println("\n7. 过期测试")
	c.Set("temp_data", "这是临时数据", time.Second*2) // 2秒过期
	fmt.Println("设置临时数据，2秒后过期")

	time.Sleep(time.Second * 3) // 等待过期

	if value, found := c.Get("temp_data"); found {
		fmt.Printf("临时数据仍然存在: %v\n", value)
	} else {
		fmt.Println("临时数据已过期")
	}

	// 8. 最终统计
	fmt.Println("\n8. 最终统计信息")
	finalStats := c.Stats()
	fmt.Printf("总命中次数: %d\n", finalStats.Hits)
	fmt.Printf("总未命中次数: %d\n", finalStats.Misses)
	fmt.Printf("总设置次数: %d\n", finalStats.Sets)
	fmt.Printf("总删除次数: %d\n", finalStats.Deletes)
	fmt.Printf("最终命中率: %.2f%%\n", finalStats.HitRate*100)
	fmt.Printf("最终缓存大小: %d\n", finalStats.Size)

	// 9. 清空缓存
	fmt.Println("\n9. 清空缓存")
	c.Flush()
	fmt.Printf("清空后缓存大小: %d\n", c.Size())
}
