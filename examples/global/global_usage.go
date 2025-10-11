package main

import (
	"fmt"
	"time"

	"scache/cache"
)

func main() {
	fmt.Println("=== 全局缓存使用示例 ===")

	// 1. 使用全局缓存（无需实例化）
	fmt.Println("\n1. 使用全局缓存设置和获取")

	err := cache.Set("global:key1", "全局值1", 0)
	if err != nil {
		fmt.Printf("设置失败: %v\n", err)
		return
	}
	fmt.Println("设置 global:key1 = 全局值1")

	if value, found := cache.Get("global:key1"); found {
		fmt.Printf("获取 global:key1 = %v\n", value)
	}

	// 2. 全局缓存统计
	fmt.Println("\n2. 全局缓存统计")
	stats := cache.Stats()
	fmt.Printf("命中次数: %d\n", stats.Hits)
	fmt.Printf("未命中次数: %d\n", stats.Misses)
	fmt.Printf("命中率: %.2f%%\n", stats.HitRate*100)
	fmt.Printf("当前大小: %d\n", stats.Size)

	// 3. 全局缓存配置（首次使用前）
	fmt.Println("\n3. 配置全局缓存示例")
	fmt.Println("(注意: 全局缓存一旦初始化就无法重新配置)")

	// 4. 全局缓存的所有键
	fmt.Println("\n4. 获取全局缓存所有键")

	// 添加更多数据
	cache.Set("global:key2", "全局值2", time.Hour)
	cache.Set("global:key3", "全局值3", time.Hour)

	keys := cache.Keys()
	fmt.Printf("全局缓存中的键: %v\n", keys)

	// 5. 带过期时间的获取
	fmt.Println("\n5. 带过期时间的获取")
	cache.Set("global:temp", "临时数据", time.Minute*5)

	value, expiration, found := cache.GetWithExpiration("global:temp")
	if found {
		fmt.Printf("值: %v, 过期时间: %v\n", value, expiration.Format("2006-01-02 15:04:05"))
	}

	// 6. 全局缓存的存在性检查
	fmt.Println("\n6. 全局缓存存在性检查")

	if cache.Exists("global:key1") {
		fmt.Println("global:key1 存在")
	}

	if !cache.Exists("global:nonexistent") {
		fmt.Println("global:nonexistent 不存在")
	}

	// 7. 全局缓存大小
	fmt.Println("\n7. 全局缓存大小")
	fmt.Printf("当前缓存大小: %d\n", cache.Size())

	// 8. 全局缓存删除
	fmt.Println("\n8. 全局缓存删除")

	if cache.Delete("global:key2") {
		fmt.Println("成功删除 global:key2")
	}

	fmt.Printf("删除后缓存大小: %d\n", cache.Size())

	// 9. 全局缓存清空
	fmt.Println("\n9. 清空全局缓存")
	cache.Flush()
	fmt.Printf("清空后缓存大小: %d\n", cache.Size())

	fmt.Println("\n=== 全局缓存使用完成 ===")
	fmt.Println("提示: 全局缓存在程序整个生命周期中保持单例状态")
}
