package main

import (
	"fmt"
	"log"
	"time"

	"github.com/scache"
)

func main() {
	fmt.Println("=== SCache 全局缓存示例 ===")
	fmt.Println()

	// 1. 注册不同类型的全局缓存
	fmt.Println("1. 注册全局缓存")

	// 注册用户缓存（LRU策略）
	if err := scache.RegisterLRU("users", 1000); err != nil {
		log.Fatal("注册用户缓存失败:", err)
	}
	fmt.Println("✓ 注册 LRU 用户缓存 (users)")

	// 注册会话缓存（LFU策略，带过期时间）
	if err := scache.RegisterLFU("sessions", 500); err != nil {
		log.Fatal("注册会话缓存失败:", err)
	}
	fmt.Println("✓ 注册 LFU 会话缓存 (sessions)")

	// 注册产品缓存（FIFO策略）
	if err := scache.RegisterFIFO("products", 2000); err != nil {
		log.Fatal("注册产品缓存失败:", err)
	}
	fmt.Println("✓ 注册 FIFO 产品缓存 (products)")

	// 2. 使用不同的缓存
	fmt.Println("\n2. 使用不同类型的缓存")

	// 使用用户缓存
	usersCache, err := scache.Get("users")
	if err != nil {
		log.Fatal("获取用户缓存失败:", err)
	}

	if err := usersCache.Set("user:1", "Alice"); err != nil {
		log.Printf("设置用户缓存失败: %v", err)
	}
	fmt.Println("✓ 用户缓存: user:1 = Alice")

	// 使用会话缓存
	sessionsCache, err := scache.Get("sessions")
	if err != nil {
		log.Fatal("获取会话缓存失败:", err)
	}

	if err := sessionsCache.SetWithTTL("session:abc123", "user:1", 2*time.Second); err != nil {
		log.Printf("设置会话缓存失败: %v", err)
	}
	fmt.Println("✓ 会话缓存: session:abc123 = user:1 (2秒过期)")

	// 使用产品缓存
	productsCache, err := scache.Get("products")
	if err != nil {
		log.Fatal("获取产品缓存失败:", err)
	}

	if err := productsCache.Set("product:1001", "iPhone 15"); err != nil {
		log.Printf("设置产品缓存失败: %v", err)
	}
	fmt.Println("✓ 产品缓存: product:1001 = iPhone 15")

	// 3. 使用默认缓存
	fmt.Println("\n3. 使用默认全局缓存")

	if err := scache.Set("app:version", "1.0.0"); err != nil {
		log.Printf("设置默认缓存失败: %v", err)
	}

	if value, exists := scache.GetFromDefault("app:version"); exists {
		fmt.Printf("✓ 默认缓存: app:version = %v\n", value)
	}

	// 4. 查看所有注册的缓存
	fmt.Println("\n4. 查看所有注册的缓存")
	caches := scache.List()
	fmt.Printf("已注册的缓存: %v\n", caches)

	// 5. 查看缓存统计
	fmt.Println("\n5. 缓存统计信息")
	stats := scache.Stats()
	for name, stat := range stats {
		fmt.Printf("%s 缓存:\n", name)
		fmt.Printf("  大小: %d/%d\n", stat.Size, stat.MaxSize)
		fmt.Printf("  命中率: %.2f%%\n", stat.HitRate*100)
		fmt.Printf("  命中次数: %d\n", stat.Hits)
		fmt.Printf("  未命中次数: %d\n", stat.Misses)
	}

	// 6. 测试过期机制
	fmt.Println("\n6. 测试过期机制")
	fmt.Println("等待会话缓存过期...")
	time.Sleep(3 * time.Second)

	if value, exists := sessionsCache.Get("session:abc123"); !exists {
		fmt.Println("✓ 会话缓存已过期")
	} else {
		fmt.Printf("✗ 会话缓存未过期，值: %v\n", value)
	}

	// 7. 缓存列表和总大小
	fmt.Println("\n7. 全局缓存总览")
	fmt.Printf("缓存总数: %d\n", len(caches))
	fmt.Printf("缓存总大小: %d 项\n", scache.Size())

	// 8. 检查缓存是否存在
	fmt.Println("\n8. 缓存存在性检查")
	fmt.Printf("users 缓存存在: %v\n", scache.Exists("users"))
	fmt.Printf("nonexistent 缓存存在: %v\n", scache.Exists("nonexistent"))

	// 9. GetOrDefault 示例
	fmt.Println("\n9. GetOrDefault 示例")
	tempCache := scache.GetOrDefault("temp")
	if err := tempCache.Set("temp-key", "temp-value"); err != nil {
		log.Printf("设置临时缓存失败: %v", err)
	}
	fmt.Println("✓ 使用 GetOrDefault 创建临时缓存")

	// 10. 清理演示
	fmt.Println("\n10. 清理演示")

	// 清空默认缓存
	if err := scache.ClearDefault(); err != nil {
		log.Printf("清空默认缓存失败: %v", err)
	} else {
		fmt.Println("✓ 清空默认缓存")
	}

	// 移除临时缓存
	if err := scache.Remove("temp"); err != nil {
		log.Printf("移除临时缓存失败: %v", err)
	} else {
		fmt.Println("✓ 移除临时缓存")
	}

	// 11. 最终统计
	fmt.Println("\n11. 最终统计")
	finalStats := scache.Stats()
	fmt.Printf("剩余缓存数量: %d\n", len(finalStats))
	fmt.Printf("剩余缓存总大小: %d\n", scache.Size())

	// 12. 清理所有缓存
	fmt.Println("\n12. 清理所有缓存")
	if err := scache.Close(); err != nil {
		log.Printf("关闭所有缓存失败: %v", err)
	} else {
		fmt.Println("✓ 关闭所有缓存")
	}

	fmt.Println("\n=== 全局缓存示例完成 ===")
}
