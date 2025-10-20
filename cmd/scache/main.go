package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scache-io/scache/generator"
	"github.com/spf13/cobra"
)

const (
	appName = "scache"
	version = "1.0.0"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     appName,
		Short:   "Go 结构体缓存代码生成工具",
		Long:    ``,
		Version: version,
	}

	// 添加 gen 子命令
	var genCmd = &cobra.Command{
		Use:   "gen",
		Short: "生成结构体模板代码",
		Long:  ``,
		RunE:  runGen,
	}

	// gen 命令参数
	var (
		dir      string
		pkgName  string
		excludes string
		structs  string
	)

	genCmd.Flags().StringVarP(&dir, "dir", "d", ".", "项目目录路径")
	genCmd.Flags().StringVarP(&pkgName, "package", "p", "", "包名（默认为目录名）")
	genCmd.Flags().StringVarP(&excludes, "exclude", "e", "vendor,node_modules,.git", "排除的目录，用逗号分隔")
	genCmd.Flags().StringVarP(&structs, "structs", "s", "", "指定结构体名称，用逗号分隔（默认生成所有）")

	// 设置 gen 命令为默认命令
	rootCmd.AddCommand(genCmd)

	// 如果没有参数，默认执行 gen
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "gen")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// runGen 执行生成命令
func runGen(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("dir")
	pkgName, _ := cmd.Flags().GetString("package")
	excludes, _ := cmd.Flags().GetString("exclude")
	structs, _ := cmd.Flags().GetString("structs")

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", dir)
	}

	// 确定包名
	packageName := pkgName
	if packageName == "" {
		packageName = filepath.Base(dir)
		// 如果目录名包含go mod路径，提取最后部分
		if strings.Contains(packageName, "-") {
			parts := strings.Split(packageName, "/")
			packageName = parts[len(parts)-1]
		}
	}

	// 解析排除目录
	excludeDirs := strings.Split(excludes, ",")
	for i, dir := range excludeDirs {
		excludeDirs[i] = strings.TrimSpace(dir)
	}

	// 解析指定的结构体
	var targetStructs []string
	if structs != "" {
		targetStructs = strings.Split(structs, ",")
		// 去除空白字符
		for i, s := range targetStructs {
			targetStructs[i] = strings.TrimSpace(s)
		}
	}

	// 创建生成器配置
	config := &generator.Config{
		Dir:           dir,
		Package:       packageName,
		ExcludeDirs:   excludeDirs,
		TargetStructs: targetStructs,
		SplitPackages: true, // 默认分包模式
	}

	// 执行代码生成
	if err := generator.Generate(config); err != nil {
		return fmt.Errorf("生成失败: %w", err)
	}

	printSuccess(config, packageName, dir, targetStructs)
	return nil
}

// printSuccess 打印成功信息
func printSuccess(config *generator.Config, packageName, dir string, targetStructs []string) {
	fmt.Printf("缓存代码已生成到: %s\n", dir)
	fmt.Printf("包名: %s\n", packageName)
	fmt.Printf("扫描目录: %s\n", dir)
	fmt.Printf("生成方式: 按包生成 _scache.go 文件\n")

	if len(targetStructs) > 0 {
		fmt.Printf("指定结构体: %v (%d个)\n", targetStructs, config.GeneratedCount)
	} else {
		fmt.Printf("生成所有结构体 (%d个)\n", config.GeneratedCount)
	}

	fmt.Printf("\n使用示例:\n")
	fmt.Printf("  import \"yourproject/%s\"\n", packageName)
	fmt.Printf("  \n")
	fmt.Printf("  // 使用默认缓存实例\n")
	if len(targetStructs) > 0 {
		fmt.Printf("  cache := %s.Get%sScache()\n", packageName, targetStructs[0])
		fmt.Printf("  cache.Store(\"key\", %s{}, time.Hour)\n", targetStructs[0])
	} else {
		fmt.Printf("  cache := %s.GetExampleScache()\n", packageName)
		fmt.Printf("  cache.Store(\"key\", Example{}, time.Hour)\n")
	}
}
