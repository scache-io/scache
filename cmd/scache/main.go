package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/scache-io/scache/cmd/scache/generator"
	"github.com/spf13/cobra"
)

//go:embed version.info
var versionInfo string

//go:embed help/root.txt
var rootHelp string

//go:embed help/gen.txt
var genHelp string

//go:embed help/version.txt
var versionHelp string

func getVersion() string {
	return strings.TrimSpace(versionInfo)
}

func main() {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(newGenCmd())
	rootCmd.AddCommand(newVersionCmd())

	// 默认执行 gen
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "gen")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "scache",
		Short:   "Go 结构体缓存代码生成工具",
		Long:    rootHelp,
		Version: getVersion(),
	}
}

func newGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "生成结构体缓存代码",
		Long:  genHelp,
		RunE:  runGen,
	}

	cmd.Flags().StringP("dir", "d", ".", "项目目录路径")
	cmd.Flags().StringP("package", "p", "", "包名（默认为目录名）")
	cmd.Flags().StringP("exclude", "e", "vendor,node_modules,.git", "排除的目录")
	cmd.Flags().StringP("structs", "s", "", "指定结构体名称，用逗号分隔")
	cmd.Flags().BoolP("generic", "g", false, "使用泛型版本（Go 1.18+）")

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Long:  versionHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("SCache version %s\n", getVersion())
			fmt.Println("GitHub: https://github.com/scache-io/scache")
			fmt.Println("Docs: https://github.com/scache-io/scache/blob/main/README.md")
		},
	}
}

func runGen(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("dir")
	pkgName, _ := cmd.Flags().GetString("package")
	excludes, _ := cmd.Flags().GetString("exclude")
	structs, _ := cmd.Flags().GetString("structs")
	useGeneric, _ := cmd.Flags().GetBool("generic")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", dir)
	}

	packageName := pkgName
	if packageName == "" {
		packageName = filepath.Base(dir)
		if strings.Contains(packageName, "-") {
			parts := strings.Split(packageName, "/")
			packageName = parts[len(parts)-1]
		}
	}

	excludeDirs := strings.Split(excludes, ",")
	for i, d := range excludeDirs {
		excludeDirs[i] = strings.TrimSpace(d)
	}

	var targetStructs []string
	if structs != "" {
		targetStructs = strings.Split(structs, ",")
		for i, s := range targetStructs {
			targetStructs[i] = strings.TrimSpace(s)
		}
	}

	config := &generator.Config{
		Dir:           dir,
		Package:       packageName,
		ExcludeDirs:   excludeDirs,
		TargetStructs: targetStructs,
		SplitPackages: false,
		UseGeneric:    useGeneric,
	}

	if err := ensureScachePackage(dir); err != nil {
		return fmt.Errorf("安装 scache 包失败: %w", err)
	}

	if err := generator.Generate(config); err != nil {
		return fmt.Errorf("生成失败: %w", err)
	}

	// 生成代码后执行 go mod tidy，确保依赖完整
	fmt.Println("正在整理依赖...")
	projectRoot, err := findProjectRoot(dir)
	if err == nil {
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectRoot
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			fmt.Printf("警告: go mod tidy 失败: %v\n", err)
		} else if len(output) > 0 {
			fmt.Printf("依赖已更新\n")
		}
	}

	printSuccess(config, packageName, dir, targetStructs)
	return nil
}

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
	fmt.Printf("  // 使用缓存实例\n")
	if len(targetStructs) > 0 {
		fmt.Printf("  cache := %s.Get%sScache()\n", packageName, targetStructs[0])
		fmt.Printf("  cache.Store(\"key\", %s{}, time.Hour)\n", targetStructs[0])
	} else {
		fmt.Printf("  cache := %s.GetExampleScache()\n", packageName)
		fmt.Printf("  cache.Store(\"key\", Example{}, time.Hour)\n")
	}
}

func ensureScachePackage(dir string) error {
	projectRoot, err := findProjectRoot(dir)
	if err != nil {
		fmt.Println("未找到 go.mod 文件，正在初始化...")
		if err := initGoMod(dir); err != nil {
			return fmt.Errorf("初始化 go.mod 失败: %w", err)
		}
		projectRoot = dir
	}

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

func initGoMod(dir string) error {
	dirName := filepath.Base(dir)
	if dirName == "." || dirName == "/" {
		dirName = "project"
	}
	moduleName := "example.com/" + dirName

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行 go mod init 失败: %v, output: %s", err, string(output))
	}
	return nil
}

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

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("未找到 go.mod 文件")
		}
		currentDir = parent
	}
}

func isScachePackageInstalled(dir string) bool {
	cmd := exec.Command("go", "list", "-m", "all")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "github.com/scache-io/scache")
}

func installScachePackage(dir string) error {
	// 使用当前版本而不是 @latest，确保版本一致
	version := getVersion()
	pkgPath := fmt.Sprintf("github.com/scache-io/scache@v%s", version)

	cmd := exec.Command("go", "get", pkgPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果指定版本失败，尝试使用 @latest
		fmt.Printf("警告: 安装 v%s 失败，尝试使用最新版本\n", version)
		cmd = exec.Command("go", "get", "github.com/scache-io/scache@latest")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("执行 go get 失败: %v, output: %s", err, string(output))
		}
	}
	fmt.Printf("scache 包安装成功\n")

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = dir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行 go mod tidy 失败: %v, output: %s", err, string(output))
	}

	return nil
}
