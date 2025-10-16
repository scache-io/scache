package main

import (
	"fmt"
	"log"
	"time"

	"github.com/scache-io/scache"
	"github.com/scache-io/scache/config"
)

func main() {
	fmt.Println("=== 简单缓存使用示例 ===")

	// 演示1: 局部缓存使用
	fmt.Println("\n1. 局部缓存使用示例")
	demonstrateLocalCache()

	// 演示2: 全局缓存使用
	fmt.Println("\n2. 全局缓存使用示例")
	demonstrateGlobalCache()

	// 演示3: 不同数据类型使用
	fmt.Println("\n3. 不同数据类型使用示例")
	demonstrateDataTypes()
}

func demonstrateLocalCache() {
	// 创建一个局部缓存实例，配置最大100个键，默认5分钟过期
	localCache := scache.New(
		config.WithMaxSize(100),
		config.WithDefaultExpiration(5*time.Minute),
	)

	// 设置和获取字符串
	err := localCache.SetString("username", "alice", 10*time.Minute)
	if err != nil {
		log.Printf("设置字符串失败: %v", err)
		return
	}

	if username, found := localCache.GetString("username"); found {
		fmt.Printf("用户名: %s\n", username)
	}

	// 设置和获取列表
	err = localCache.SetList("user_roles", []interface{}{"admin", "user", "moderator"}, 0)
	if err != nil {
		log.Printf("设置列表失败: %v", err)
		return
	}

	if roles, found := localCache.GetList("user_roles"); found {
		fmt.Printf("用户角色: %v\n", roles)
	}

	// 设置和获取哈希
	userProfile := map[string]interface{}{
		"name":  "Alice",
		"age":   25,
		"email": "alice@example.com",
	}
	err = localCache.SetHash("user_profile", userProfile, 0)
	if err != nil {
		log.Printf("设置哈希失败: %v", err)
		return
	}

	if profile, found := localCache.GetHash("user_profile"); found {
		fmt.Printf("用户档案: 姓名=%s, 年龄=%v, 邮箱=%s\n",
			profile["name"], profile["age"], profile["email"])
	}

	// 显示缓存统计
	stats := localCache.Stats()
	if statsMap, ok := stats.(map[string]interface{}); ok {
		fmt.Printf("缓存统计: 大小=%.0f, 命中率=%.2f%%\n",
			statsMap["keys"], statsMap["hit_rate"].(float64)*100)
	}
}

func demonstrateGlobalCache() {
	// 初始化全局缓存，使用中型配置
	scache.InitGlobalCache(config.MediumConfig...)

	// 全局缓存操作非常简单，直接调用函数即可
	scache.SetString("app_name", "MyApplication", 0)
	scache.SetString("version", "1.0.0", 0)
	scache.SetString("author", "Developer", 0)

	// 获取全局缓存值
	if appName, found := scache.GetString("app_name"); found {
		fmt.Printf("应用名称: %s\n", appName)
	}

	if version, found := scache.GetString("version"); found {
		fmt.Printf("版本: %s\n", version)
	}

	// 设置应用配置
	appConfig := map[string]interface{}{
		"debug":      true,
		"port":       8080,
		"max_conn":   1000,
		"timeout":    30,
		"enable_ssl": false,
	}
	scache.SetHash("app_config", appConfig, 0)

	if config, found := scache.GetHash("app_config"); found {
		fmt.Printf("应用配置: 端口=%v, 调试模式=%v\n",
			config["port"], config["debug"])
	}

	// 全局缓存统计
	globalStats := scache.Stats()
	if statsMap, ok := globalStats.(map[string]interface{}); ok {
		fmt.Printf("全局缓存统计: 大小=%.0f, 命中次数=%.0f\n",
			statsMap["keys"], statsMap["hits"])
	}
}

func demonstrateDataTypes() {
	// 创建一个专门用于演示数据类型的缓存
	typeCache := scache.New()

	// 字符串类型
	typeCache.SetString("message", "Hello, World!", 0)
	if message, found := typeCache.GetString("message"); found {
		fmt.Printf("字符串: %s\n", message)
	}

	// 列表类型
	todoList := []interface{}{
		"完成项目文档",
		"代码审查",
		"部署到生产环境",
		"性能优化",
	}
	typeCache.SetList("todos", todoList, 0)
	if todos, found := typeCache.GetList("todos"); found {
		fmt.Printf("待办事项: %v\n", todos)
	}

	// 哈希类型
	productInfo := map[string]interface{}{
		"id":          1001,
		"name":        "智能手机",
		"price":       2999.00,
		"stock":       50,
		"category":    "电子产品",
		"description": "高性能智能手机",
	}
	typeCache.SetHash("product_1001", productInfo, 0)
	if product, found := typeCache.GetHash("product_1001"); found {
		fmt.Printf("产品信息: 名称=%s, 价格=%v, 库存=%v\n",
			product["name"], product["price"], product["stock"])
	}

	// 演示过期操作
	fmt.Println("\n过期操作演示:")
	typeCache.SetString("temp_key", "这个值会过期", time.Second*2)
	fmt.Println("设置了2秒过期的临时键")

	// 检查TTL
	if ttl, exists := typeCache.TTL("temp_key"); exists {
		fmt.Printf("剩余时间: %v\n", ttl)
	}

	// 等待过期
	time.Sleep(time.Second * 3)
	if _, exists := typeCache.GetString("temp_key"); !exists {
		fmt.Println("临时键已过期")
	}

	// 演示其他操作
	fmt.Println("\n其他操作演示:")
	fmt.Printf("缓存大小: %d\n", typeCache.Size())
	fmt.Printf("所有键: %v\n", typeCache.Keys())
	fmt.Printf("键 'message' 是否存在: %t\n", typeCache.Exists("message"))

	// 删除一个键
	if typeCache.Delete("message") {
		fmt.Println("成功删除 'message' 键")
	}

	// 清空缓存
	typeCache.Flush()
	fmt.Printf("清空后缓存大小: %d\n", typeCache.Size())
}
