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

// ANSI color codes
const (
	colorReset  = "\x1b[0m"
	colorGreen  = "\x1b[32m"
	colorRed    = "\x1b[31m"
	colorYellow = "\x1b[33m"
	colorCyan   = "\x1b[36m"
)

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

	// 禁用 Cobra 的自动错误输出，我们自己处理
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s✗%s %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "scache",
		Short:   "Go struct cache code generator",
		Long:    rootHelp,
		Version: getVersion(),
	}
}

func newGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate struct cache code",
		Long:  genHelp,
		RunE:  runGen,
	}

	cmd.Flags().StringP("dir", "d", ".", "Project directory")
	cmd.Flags().StringP("package", "p", "", "Package name (default: directory name)")
	cmd.Flags().StringP("exclude", "e", "vendor,node_modules,.git", "Exclude directories")
	cmd.Flags().StringP("structs", "s", "", "Specific structs (comma-separated)")
	cmd.Flags().BoolP("generic", "g", false, "Use generic version (Go 1.18+)")

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
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
		return fmt.Errorf("directory not found: %s", dir)
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
		return err
	}

	if err := generator.Generate(config); err != nil {
		return err
	}

	// Auto run go mod tidy
	projectRoot, err := findProjectRoot(dir)
	if err == nil {
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectRoot
		_ = tidyCmd.Run()
	}

	// Success output
	printSuccess(config, packageName, dir, targetStructs)
	return nil
}

func printSuccess(config *generator.Config, packageName, dir string, targetStructs []string) {
	fmt.Printf("%s✓%s Generated %d struct(s): %s\n", colorGreen, colorReset, config.GeneratedCount, dir)
}

func ensureScachePackage(dir string) error {
	projectRoot, err := findProjectRoot(dir)
	if err != nil {
		fmt.Printf("%s→%s Initializing go.mod...\n", colorCyan, colorReset)
		if err := initGoMod(dir); err != nil {
			return err
		}
		projectRoot = dir
	}

	if !isScachePackageInstalled(projectRoot) {
		fmt.Printf("%s→%s Installing scache...\n", colorCyan, colorReset)
		if err := installScachePackage(projectRoot); err != nil {
			return err
		}
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
		return fmt.Errorf("go mod init failed: %v, output: %s", err, string(output))
	}
	return nil
}

func findProjectRoot(dir string) (string, error) {
	currentDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("get absolute path failed: %w", err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("go.mod not found")
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
	version := getVersion()
	pkgPath := fmt.Sprintf("github.com/scache-io/scache@v%s", version)

	cmd := exec.Command("go", "get", pkgPath)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		cmd = exec.Command("go", "get", "github.com/scache-io/scache@latest")
		cmd.Dir = dir
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("install failed: %v, output: %s", err, string(output))
		}
	}

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = dir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v, output: %s", err, string(output))
	}

	return nil
}
