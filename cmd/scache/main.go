package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/scache-io/scache/cmd/scache/generator"
	"github.com/spf13/cobra"
)

const (
	appName = "scache"
	version = "0.0.1"
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

	// 检测并自动安装 scache 包
	if err := ensureScachePackage(dir); err != nil {
		return fmt.Errorf("安装scache包失败: %w", err)
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

// ensureScachePackage 检测并自动安装 scache 包
func ensureScachePackage(dir string) error {
	// 查找项目根目录的 go.mod 文件
	projectRoot, err := findProjectRoot(dir)
	if err != nil {
		// 如果找不到 go.mod，在当前目录初始化一个
		fmt.Println("未找到 go.mod 文件，正在初始化...")
		if err := initGoMod(dir); err != nil {
			return fmt.Errorf("初始化 go.mod 失败: %w", err)
		}
		projectRoot = dir
	}

	// 检测 scache 包是否已安装
	if !isScachePackageInstalled(projectRoot) {
		fmt.Println("正在安装 scache 包...")
		if err := installScachePackage(projectRoot); err != nil {
			return fmt.Errorf("安装 scache 包失败: %w", err)
		}
		fmt.Println("scache 包安装成功")
	} else {
		fmt.Println("scache 包已存在")
	}

	return nil
}

// initGoMod go.mod 文件初始化
func initGoMod(dir string) error {
	// 获取最后一个目录名作为模块名
	dirName := filepath.Base(dir)
	if dirName == "." || dirName == "/" {
		dirName = "project"
	}
	moduleName := "example.com/" + dirName

	// 在传入的目录下执行 go mod init
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行 go mod init 失败: %v, output: %s", err, string(output))
	}
	return nil
}

// findProjectRoot 查找项目的根目录（包含go.mod文件的目录）
func findProjectRoot(dir string) (string, error) {
	currentDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("获取绝对路径失败: %w", err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		// 移动到父目录
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// 已经到达根目录
			return "", fmt.Errorf("未找到 go.mod 文件，请先在项目根目录初始化 go.mod")
		}
		currentDir = parent
	}
}

// isScachePackageInstalled 检测 scache 包是否已安装
func isScachePackageInstalled(dir string) bool {
	cmd := exec.Command("go", "list", "-m", "all")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	// 检查输出中是否包含 scache 包
	return strings.Contains(string(output), "github.com/scache-io/scache")
}

// installScachePackage 安装 scache 包
func installScachePackage(dir string) error {
	// 运行 go get 安装包
	cmd := exec.Command("go", "get", "github.com/scache-io/scache@latest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 go get 失败: %v, output: %s", err, string(output))
	}
	fmt.Printf("go get 输出: %s\n", string(output))

	// 运行 go mod tidy 确保依赖完整
	tidyCmd := exec.Command("go", "mod", "tidy")
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行 go mod tidy 失败: %v, output: %s", err, string(output))
	}
	fmt.Printf("go mod tidy 输出: %s\n", string(output))

	return nil
}
