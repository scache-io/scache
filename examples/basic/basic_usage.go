package main

import (
	"fmt"
	"log"
	"time"

	"github.com/scache-io/scache/cache"
	"github.com/scache-io/scache/config"
)

func main() {
	fmt.Println("=== 基本缓存使用示例 ===")

	// 1. 局部缓存实例使用
	fmt.Println("\n1. 局部缓存实例使用")
	c := cache.NewLocalCache(
		config.WithMaxSize(100),
		config.WithDefaultExpiration(time.Minute*5),
	)
	fmt.Printf("缓存创建成功，大小: %d\n", c.Size())

	// 2. 局部缓存基本操作
	fmt.Println("\n2. 局部缓存基本操作演示")

	// 设置字符串
	err := c.SetString("user:1001", "张三", time.Minute*10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("设置 user:1001 = 张三 (10分钟过期)")

	// 获取字符串
	if value, found := c.GetString("user:1001"); found {
		fmt.Printf("获取 user:1001 = %v\n", value)
	} else {
		fmt.Println("未找到 user:1001")
	}

	// 检查是否存在
	if c.Exists("user:1001") {
		fmt.Println("user:1001 存在于缓存中")
	}

	// 3. 局部缓存列表操作
	fmt.Println("\n3. 列表操作演示")
	err = c.SetList("permissions:1001", []interface{}{"read", "write", "delete"}, 0)
	if err != nil {
		log.Fatal(err)
	}

	if permissions, found := c.GetList("permissions:1001"); found {
		fmt.Printf("用户权限: %v\n", permissions)
	}

	// 4. 局部缓存哈希操作
	fmt.Println("\n4. 哈希操作演示")
	userProfile := map[string]interface{}{
		"name":  "张三",
		"age":   30,
		"email": "zhangsan@example.com",
	}
	err = c.SetHash("profile:1001", userProfile, 0)
	if err != nil {
		log.Fatal(err)
	}

	if profile, found := c.GetHash("profile:1001"); found {
		fmt.Printf("用户档案: %+v\n", profile)
	}

	// 5. 缓存统计
	fmt.Println("\n5. 缓存统计信息")
	stats := c.Stats()
	if statsMap, ok := stats.(map[string]interface{}); ok {
		fmt.Printf("命中次数: %.0f\n", statsMap["hits"])
		fmt.Printf("未命中次数: %.0f\n", statsMap["misses"])
		fmt.Printf("命中率: %.2f%%\n", statsMap["hit_rate"].(float64)*100)
		fmt.Printf("当前大小: %.0f\n", statsMap["keys"])
	}

	// 6. 全局缓存使用演示
	fmt.Println("\n6. 全局缓存使用演示")

	// 使用全局缓存
	cache.SetString("global:user:2001", "李四", time.Minute*15)
	cache.SetString("global:user:2002", "王五", time.Minute*15)

	if name, found := cache.GetString("global:user:2001"); found {
		fmt.Printf("全局缓存获取: %v\n", name)
	}

	// 设置全局哈希
	globalProfile := map[string]interface{}{
		"department": "技术部",
		"position":   "高级工程师",
		"salary":     15000,
	}
	cache.SetHash("global:profile:2001", globalProfile, 0)

	if profile, found := cache.GetHash("global:profile:2001"); found {
		fmt.Printf("全局用户档案: %+v\n", profile)
	}

	// 7. 全局缓存统计
	fmt.Println("\n7. 全局缓存统计")
	globalStats := cache.Stats()
	if globalStatsMap, ok := globalStats.(map[string]interface{}); ok {
		fmt.Printf("全局缓存命中次数: %.0f\n", globalStatsMap["hits"])
		fmt.Printf("全局缓存大小: %.0f\n", globalStatsMap["keys"])
	}

	// 8. 过期测试
	fmt.Println("\n8. 过期测试")
	c.SetString("temp_data", "这是临时数据", time.Second*2) // 2秒过期
	fmt.Println("设置临时数据，2秒后过期")

	time.Sleep(time.Second * 3) // 等待过期

	if value, found := c.GetString("temp_data"); found {
		fmt.Printf("临时数据仍然存在: %v\n", value)
	} else {
		fmt.Println("临时数据已过期")
	}

	// 9. TTL操作
	fmt.Println("\n9. TTL操作演示")
	c.SetString("ttl_test", "TTL测试数据", 0)

	// 设置5秒TTL
	success := c.Expire("ttl_test", time.Second*5)
	if success {
		if ttl, exists := c.TTL("ttl_test"); exists {
			fmt.Printf("TTL剩余时间: %v\n", ttl)
		}
	}

	// 10. 批量操作
	fmt.Println("\n10. 批量操作")
	users := map[string]string{
		"user:1002": "李四",
		"user:1003": "王五",
		"user:1004": "赵六",
	}

	for key, value := range users {
		c.SetString(key, value, 0) // 使用默认过期时间
		fmt.Printf("设置 %s = %s\n", key, value)
	}

	fmt.Printf("缓存当前大小: %d\n", c.Size())

	// 11. 获取所有键
	fmt.Println("\n11. 获取所有键")
	keys := c.Keys()
	fmt.Printf("所有键: %v\n", keys)

	globalKeys := cache.Keys()
	fmt.Printf("全局缓存所有键: %v\n", globalKeys)

	// 12. 清空缓存
	fmt.Println("\n12. 清空缓存")
	c.Flush()
	fmt.Printf("清空后局部缓存大小: %d\n", c.Size())

	cache.Flush()
	fmt.Printf("清空后全局缓存大小: %d\n", cache.Size())
}
