package main

import (
	"fmt"
	"time"

	"scache"
)

func main() {
	// 使用全局缓存API
	err := scache.Set("test_key", "test_value", time.Minute)
	if err != nil {
		panic(err)
	}

	value, found, err := scache.Get("test_key")
	if err != nil {
		panic(err)
	}

	if found {
		fmt.Printf("全局缓存测试成功: %s\n", value)
	}

	// 获取统计信息
	stats := scache.Stats()
	fmt.Printf("缓存统计: %+v\n", stats)
}
