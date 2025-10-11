package main

import (
	"fmt"
	"strings"

	"github.com/scache/utils"
)

func main() {
	fmt.Println("=== SCache 常量、变量、结构体管理展示 ===")
	fmt.Println()

	// 1. 项目结构信息
	fmt.Println("1. 项目结构信息")
	fmt.Println(strings.Repeat("=", 50))

	projectInfo := utils.NewProjectInfo()
	projectInfo.PrintProjectStructure()

	// 2. 项目统计信息
	fmt.Println("2. 项目统计信息")
	fmt.Println(strings.Repeat("=", 50))

	projectInfo.PrintStatistics()

	// 3. 常量管理展示
	fmt.Println("3. 常量管理展示")
	fmt.Println(strings.Repeat("=", 50))

	constantsManager := utils.NewConstantsManager()
	constantsManager.PrintConstants()

	// 4. 类型管理展示
	fmt.Println("4. 类型管理展示")
	fmt.Println(strings.Repeat("=", 50))

	typesManager := utils.NewTypesManager()
	typesManager.PrintTypes()

	// 5. 配置验证展示
	fmt.Println("5. 配置验证展示")
	fmt.Println(strings.Repeat("=", 50))

	validator := utils.NewConfigValidator()

	// 测试有效配置
	validConfig := struct {
		MaxSize              int
		Shards               int
		DefaultTTL           interface{} // 使用interface{}来模拟time.Duration
		CleanupInterval      interface{}
		EvictionPolicy       string
		EnableStatistics     bool
		EnableLazyExpiration bool
		EnableMetrics        bool
	}{
		MaxSize:              10000,
		Shards:               16,
		DefaultTTL:           "0s",
		CleanupInterval:      "10m",
		EvictionPolicy:       "lru",
		EnableStatistics:     true,
		EnableLazyExpiration: true,
		EnableMetrics:        false,
	}

	fmt.Println("验证有效配置:")
	result := validator.ValidateConfig(validConfig)
	validator.PrintValidationResult(result)

	// 测试无效配置
	invalidConfig := struct {
		Name string // 错误的字段名
	}{
		Name: "test",
	}

	fmt.Println("验证无效配置:")
	result = validator.ValidateConfig(invalidConfig)
	validator.PrintValidationResult(result)

	// 6. 管理功能总结
	fmt.Println("6. 管理功能总结")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("✓ 统一常量管理 - 所有配置常量集中管理")
	fmt.Println("✓ 类型定义组织 - 结构体按功能分类管理")
	fmt.Println("✓ 全局变量控制 - 线程安全的全局状态管理")
	fmt.Println("✓ 配置验证 - 自动检查配置完整性和正确性")
	fmt.Println("✓ 项目信息统计 - 提供详细的项目结构信息")
	fmt.Println("✓ 错误处理统一 - 标准化的错误类型和消息")
	fmt.Println("✓ 工具函数丰富 - 提供便捷的管理和查询功能")

	fmt.Println("\n=== 管理展示完成 ===")
}
