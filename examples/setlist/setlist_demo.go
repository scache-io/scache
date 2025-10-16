package main

import (
	"fmt"
	"time"

	"scache/cache"
)

func main() {
	fmt.Println("=== SetList 批量设置示例 ===")

	// 1. 创建缓存实例
	c := cache.NewCache(
		cache.WithMaxSize(1000),
		cache.WithDefaultExpiration(time.Minute*10),
	)

	// 2. 准备各种类型的测试数据
	fmt.Println("\n2. 准备测试数据")

	// 数字数组
	numbers := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 字符串数组
	strings := []interface{}{"apple", "banana", "orange", "grape", "watermelon"}

	// 混合类型数组
	mixed := []interface{}{1, "hello", true, 3.14, "world", false}

	// 用户对象数组
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	user1 := User{ID: 1, Name: "Alice", Age: 25}
	user2 := User{ID: 2, Name: "Bob", Age: 30}
	user3 := User{ID: 3, Name: "Charlie", Age: 35}

	users := []interface{}{user1, user2, user3}

	// 3. 批量设置测试
	fmt.Println("\n3. 批量设置数字数组")
	err1 := c.SetList("numbers", numbers, 0) // 使用默认TTL
	if err1 != nil {
		fmt.Printf("设置失败: %v\n", err1)
	} else {
		fmt.Println("设置成功")
	}

	fmt.Println("\n4. 批量设置字符串数组")
	err2 := c.SetList("strings", strings, time.Minute*5) // 设置5分钟过期
	if err2 != nil {
		fmt.Printf("设置失败: %v\n", err2)
	} else {
		fmt.Println("设置成功")
	}

	fmt.Println("\n5. 批量设置混合类型数组")
	err3 := c.SetList("mixed", mixed, time.Hour*2) // 设置2小时过期
	if err3 != nil {
		fmt.Printf("设置失败: %v\n", err3)
	} else {
		fmt.Println("设置成功")
	}

	fmt.Println("\n6. 批量设置用户对象数组")
	err4 := c.SetList("users", users, time.Minute*15)
	if err4 != nil {
		fmt.Printf("设置失败: %v\n", err4)
	} else {
		fmt.Println("设置成功")
	}

	// 4. 验证批量设置的数据
	fmt.Println("\n7. 验证批量设置的数据")

	// 验证数字数组
	if value, found := c.Get("numbers"); found {
		fmt.Printf("✅ numbers = %v (类型: %T)\n", value, value)
		if nums, ok := value.([]interface{}); ok {
			fmt.Printf("   数组元素: %v\n", nums)
			fmt.Printf("   数组长度: %d\n", len(nums))
		}
	}

	// 验证字符串数组
	if value, found := c.Get("strings"); found {
		fmt.Printf("✅ strings = %v (类型: %T)\n", value, value)
		if strs, ok := value.([]interface{}); ok {
			fmt.Printf("   字符串列表: ")
			for i, s := range strs {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%v", s)
			}
			fmt.Println()
		}
	}

	// 验证混合类型数组
	if value, found := c.Get("mixed"); found {
		fmt.Printf("✅ mixed = %v (类型: %T)\n", value, value)
		if items, ok := value.([]interface{}); ok {
			fmt.Printf("   混合元素: ")
			for i, item := range items {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%v(%T)", item, item)
			}
			fmt.Println()
		}
	}

	// 验证用户对象数组
	if value, found := c.Get("users"); found {
		fmt.Printf("✅ users = %v (类型: %T)\n", value, value)
		if userList, ok := value.([]interface{}); ok {
			fmt.Printf("   用户列表:\n")
			for i, user := range userList {
				if u, ok := user.(User); ok {
					fmt.Printf("     [%d] ID=%d, Name=%s, Age=%d\n", i+1, u.ID, u.Name, u.Age)
				}
			}
		}
	}

	// 5. 显示最终统计
	fmt.Println("\n8. 最终缓存统计")
	stats := c.Stats()
	fmt.Printf("总设置次数: %d\n", stats.Sets)
	fmt.Printf("当前缓存大小: %d\n", stats.Size)
	fmt.Printf("命中率: %.2f%%\n", stats.HitRate*100)

	// 6. 测试错误情况
	fmt.Println("\n9. 测试错误情况")
	// 空键测试
	err5 := c.SetList("", []interface{}{"test"}, 0)
	if err5 != nil {
		fmt.Printf("空键设置失败（符合预期）: %v\n", err5)
	} else {
		fmt.Println("空键设置意外成功")
	}

	// 空值测试
	err6 := c.SetList("empty_key", nil, 0)
	if err6 != nil {
		fmt.Printf("空值设置失败: %v\n", err6)
	} else {
		fmt.Println("空值设置成功")
	}

	// 7. 测试分页查询
	fmt.Println("\n10. 测试分页查询")

	// 添加更多数据用于分页测试
	for i := 0; i < 25; i++ {
		key := fmt.Sprintf("page_item_%d", i)
		value := []interface{}{i, fmt.Sprintf("item_%d", i), i%2 == 0}
		c.Set(key, value, time.Minute*30)
	}

	// 分页查询
	page1 := c.KeysPage(1, 10)
	fmt.Printf("第1页 (每页10条): 总数=%d, 当前页=%d, 是否有下一页=%t\n",
		page1.Total, page1.Page, page1.HasNext)
	fmt.Printf("   键列表: %v\n", page1.Keys)

	page2 := c.KeysPage(2, 10)
	fmt.Printf("第2页 (每页10条): 总数=%d, 当前页=%d, 是否有下一页=%t\n",
		page2.Total, page2.Page, page2.HasNext)
	fmt.Printf("   键列表: %v\n", page2.Keys)

	page3 := c.KeysPage(3, 10)
	fmt.Printf("第3页 (每页10条): 总数=%d, 当前页=%d, 是否有下一页=%t\n",
		page3.Total, page3.Page, page3.HasNext)
	fmt.Printf("   键列表: %v\n", page3.Keys)

	fmt.Println("\n=== SetList 演示完成 ===")
}
