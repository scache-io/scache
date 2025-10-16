package main

import (
	"fmt"
	"log"
	"time"

	"github.com/scache-io/scache"
	"github.com/scache-io/scache/config"
)

func main() {
	fmt.Println("=== Redis-like Cache Demo ===")

	// 创建一个自定义配置的引擎
	engine := scache.NewEngine(
		config.WithMaxSize(1000),
		config.WithDefaultExpiration(time.Hour),
		config.WithBackgroundCleanup(time.Minute),
	)

	// 创建命令执行器
	executor := scache.NewExecutor(engine)

	// === String 操作 ===
	fmt.Println("\n--- String Operations ---")

	// SET 操作
	result, err := executor.Execute("SET", "name", "Alice", time.Minute*30)
	if err != nil {
		log.Printf("SET error: %v", err)
	} else {
		fmt.Printf("SET result: %v\n", result)
	}

	// GET 操作
	result, err = executor.Execute("GET", "name")
	if err != nil {
		log.Printf("GET error: %v", err)
	} else {
		fmt.Printf("GET result: %v\n", result)
	}

	// === List 操作 ===
	fmt.Println("\n--- List Operations ---")

	// LPUSH 操作
	result, err = executor.Execute("LPUSH", "fruits", "apple", time.Hour)
	if err != nil {
		log.Printf("LPUSH error: %v", err)
	} else {
		fmt.Printf("LPUSH result: %v\n", result)
	}

	result, err = executor.Execute("LPUSH", "fruits", "banana", time.Hour)
	if err != nil {
		log.Printf("LPUSH error: %v", err)
	} else {
		fmt.Printf("LPUSH result: %v\n", result)
	}

	// RPOP 操作
	result, err = executor.Execute("RPOP", "fruits")
	if err != nil {
		log.Printf("RPOP error: %v", err)
	} else {
		fmt.Printf("RPOP result: %v\n", result)
	}

	// === Hash 操作 ===
	fmt.Println("\n--- Hash Operations ---")

	// HSET 操作
	result, err = executor.Execute("HSET", "user:1", "name", "Bob", time.Hour)
	if err != nil {
		log.Printf("HSET error: %v", err)
	} else {
		fmt.Printf("HSET result: %v\n", result)
	}

	result, err = executor.Execute("HSET", "user:1", "age", 30, time.Hour)
	if err != nil {
		log.Printf("HSET error: %v", err)
	} else {
		fmt.Printf("HSET result: %v\n", result)
	}

	// HGET 操作
	result, err = executor.Execute("HGET", "user:1", "name")
	if err != nil {
		log.Printf("HGET error: %v", err)
	} else {
		fmt.Printf("HGET result: %v\n", result)
	}

	result, err = executor.Execute("HGET", "user:1", "age")
	if err != nil {
		log.Printf("HGET error: %v", err)
	} else {
		fmt.Printf("HGET result: %v\n", result)
	}

	// === 通用操作 ===
	fmt.Println("\n--- General Operations ---")

	// TYPE 操作
	result, err = executor.Execute("TYPE", "name")
	if err != nil {
		log.Printf("TYPE error: %v", err)
	} else {
		fmt.Printf("TYPE name: %v\n", result)
	}

	result, err = executor.Execute("TYPE", "fruits")
	if err != nil {
		log.Printf("TYPE error: %v", err)
	} else {
		fmt.Printf("TYPE fruits: %v\n", result)
	}

	result, err = executor.Execute("TYPE", "user:1")
	if err != nil {
		log.Printf("TYPE error: %v", err)
	} else {
		fmt.Printf("TYPE user:1: %v\n", result)
	}

	// EXISTS 操作
	result, err = executor.Execute("EXISTS", "name")
	if err != nil {
		log.Printf("EXISTS error: %v", err)
	} else {
		fmt.Printf("EXISTS name: %v\n", result)
	}

	result, err = executor.Execute("EXISTS", "nonexistent")
	if err != nil {
		log.Printf("EXISTS error: %v", err)
	} else {
		fmt.Printf("EXISTS nonexistent: %v\n", result)
	}

	// TTL 操作
	result, err = executor.Execute("TTL", "name")
	if err != nil {
		log.Printf("TTL error: %v", err)
	} else {
		fmt.Printf("TTL name: %v seconds\n", result)
	}

	// EXPIRE 操作
	result, err = executor.Execute("EXPIRE", "name", time.Minute*5)
	if err != nil {
		log.Printf("EXPIRE error: %v", err)
	} else {
		fmt.Printf("EXPIRE name: %v\n", result)
	}

	// === 统计信息 ===
	fmt.Println("\n--- Statistics ---")
	result, err = executor.Execute("STATS")
	if err != nil {
		log.Printf("STATS error: %v", err)
	} else {
		fmt.Printf("Statistics: %+v\n", result)
	}

	// === 列出所有命令 ===
	fmt.Println("\n--- Available Commands ---")
	commands := executor.ListCommands()
	fmt.Printf("Available commands: %v\n", commands)

	// === 使用便捷API ===
	fmt.Println("\n--- Using Convenience API ---")

	// 使用全局便捷函数
	err = scache.Set("global_key", "global_value", time.Minute*10)
	if err != nil {
		log.Printf("Global SET error: %v", err)
	} else {
		fmt.Println("Global SET successful")
	}

	value, found, err := scache.Get("global_key")
	if err != nil {
		log.Printf("Global GET error: %v", err)
	} else if found {
		fmt.Printf("Global GET result: %v\n", value)
	} else {
		fmt.Println("Global GET: key not found")
	}

	// 最终统计
	fmt.Println("\n--- Final Statistics ---")
	finalStats := scache.Stats()
	fmt.Printf("Final stats: %+v\n", finalStats)

	// 关闭执行器
	executor.Close()
	fmt.Println("\nDemo completed!")
}
