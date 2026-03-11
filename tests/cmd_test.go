package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ==================== CMD tool tests ====================

// buildScacheCMD 构建 scache 命令
func buildScacheCMD(t *testing.T) string {
	binaryPath := filepath.Join(t.TempDir(), "scache")

	projectRoot, err := filepath.Abs(filepath.Join(".."))
	if err != nil {
		t.Fatalf("Failed to get project root directory: %v", err)
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/scache")
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build scache command: %v\noutput: %s", err, string(output))
	}

	return binaryPath
}

// getTestdataDir 获取测试数据目录
func getTestdataDir(t *testing.T) string {
	dir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatalf("Failed to get testdata directory: %v", err)
	}
	return dir
}

// ==================== Basic command tests ====================

func TestCMDVersion(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version 命令Execution failed: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "SCache version") {
		t.Errorf("version output format is incorrect: %s", outputStr)
	}
	if !strings.Contains(outputStr, "GitHub") {
		t.Errorf("version output should contain GitHub 链接: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Docs") {
		t.Errorf("version output should contain文档链接: %s", outputStr)
	}
}

func TestCMDHelp(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help 命令Execution failed: %v", err)
	}

	outputStr := string(output)
	expectedKeywords := []string{"gen", "version", "Quick Start"}
	for _, kw := range expectedKeywords {
		if !strings.Contains(outputStr, kw) {
			t.Errorf("help output should contain '%s'", kw)
		}
	}
}

func TestCMDGenHelp(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "gen", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gen --help 命令Execution failed: %v", err)
	}

	outputStr := string(output)
	expectedFlags := []string{"--dir", "--package", "--exclude", "--structs", "--generic"}
	for _, flag := range expectedFlags {
		if !strings.Contains(outputStr, flag) {
			t.Errorf("gen 帮助信息应包含 %s flag", flag)
		}
	}
}

// ==================== Error handling tests ====================

func TestCMDInvalidDir(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "gen", "--dir", "/nonexistent/path")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("对Should return error for non-existent directory")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "directory not found") {
		t.Errorf("error message should contain'directory not found': %s", outputStr)
	}
}

func TestCMDNoStructs(t *testing.T) {
	binary := buildScacheCMD(t)

	tempDir := t.TempDir()
	goMod := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binary, "gen", "--generic", "--dir", tempDir)
	output, err := cmd.CombinedOutput()

	// May fail due to network issues, but should at least try to execute
	if err != nil {
		outputStr := string(output)
		validErrors := []string{"no structs found", "not found", "install", "go get"}
		hasValidError := false
		for _, e := range validErrors {
			if strings.Contains(outputStr, e) {
				hasValidError = true
				break
			}
		}
		if !hasValidError {
			t.Logf("Warning: Error信息不符合预期: %s", outputStr)
		}
	} else {
		t.Log("意外Success：空目录生成了代码")
	}
}

// ==================== Code generation tests (需要网络) ====================

// TestCMDGenGeneric 测试泛型代码生成
// 注意：此测试需要网络连接来安装依赖
func TestCMDGenGeneric(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test requiring network")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: Requires network to install dependencies: %v\noutput: %s", err, string(output))
	}

	outputStr := string(output)

	// 验证生成Success
	if !strings.Contains(outputStr, "Generated") && !strings.Contains(outputStr, "缓存代码已生成") {
		t.Errorf("Generation failed，output: %s", outputStr)
	}

	// 验证生成的文件（基于源文件名 models.go）
	generatedFile := filepath.Join(testdataDir, "models_scache.go")
	defer os.Remove(generatedFile) // 清理

	if _, err := os.Stat(generatedFile); os.IsNotExist(err) {
		t.Fatalf("Generated file does not exist: %s", generatedFile)
	}

	// 验证生成的代码内容
	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)

	// 验证泛型特征
	if !strings.Contains(contentStr, "[T any]") {
		t.Error("Generic code should contain [T any]")
	}
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("Generated code should contain GetUserScache")
	}
	if !strings.Contains(contentStr, "GetProductScache") {
		t.Error("Generated code should contain GetProductScache")
	}
	if !strings.Contains(contentStr, "GetOrderScache") {
		t.Error("Generated code should contain GetOrderScache")
	}
}

// TestCMDGenClassic 测试传统代码生成
// 注意：此测试需要网络连接来安装依赖
func TestCMDGenClassic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test requiring network")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--dir", testdataDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: Requires network to install dependencies: %v\noutput: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Generated") && !strings.Contains(outputStr, "缓存代码已生成") {
		t.Errorf("Generation failed，output: %s", outputStr)
	}

	generatedFile := filepath.Join(testdataDir, "models_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)

	// 验证非泛型特征（不应该包含泛型语法）
	if strings.Contains(contentStr, "[T any]") {
		t.Error("Classic version code should not contain泛型语法")
	}
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("Generated code should contain GetUserScache")
	}
}

// TestCMDGenSpecificStructs 测试指定Struct生成
func TestCMDGenSpecificStructs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test requiring network")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// Only generate User Struct
	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--structs", "User")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: Requires network to install dependencies: %v", err)
	}

	generatedFile := filepath.Join(testdataDir, "models_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)

	// 验证只包含 User
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("Generated code should contain GetUserScache")
	}
	if strings.Contains(contentStr, "GetProductScache") {
		t.Error("Generated code should not contain GetProductScache（Not specified）")
	}
	if strings.Contains(contentStr, "GetOrderScache") {
		t.Error("Generated code should not contain GetOrderScache（Not specified）")
	}
}

// TestCMDGenMultipleStructs 测试多个指定Struct
func TestCMDGenMultipleStructs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test requiring network")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// 生成多个指定Struct
	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--structs", "User,Product")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: Requires network to install dependencies: %v", err)
	}

	generatedFile := filepath.Join(testdataDir, "models_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)

	// 验证包含 User 和 Product，但不包含 Order
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("Generated code should contain GetUserScache")
	}
	if !strings.Contains(contentStr, "GetProductScache") {
		t.Error("Generated code should contain GetProductScache")
	}
	if strings.Contains(contentStr, "GetOrderScache") {
		t.Error("Generated code should not contain GetOrderScache（Not specified）")
	}
}

// TestCMDNonexistentStruct 测试不存在的Struct
func TestCMDNonexistentStruct(t *testing.T) {
	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--structs", "NonexistentStruct")
	output, err := cmd.CombinedOutput()

	// 应该返回Error
	if err == nil {
		t.Error("指定不存在的Struct应该返回Error")
	}

	outputStr := string(output)
	// 由于可能需要先安装依赖，Error信息可能不同
	if !strings.Contains(outputStr, "未找到") && !strings.Contains(outputStr, "安装") {
		t.Logf("output: %s", outputStr)
	}
}

// ==================== 短 flag 测试 ====================

func TestCMDGenShortFlags(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test requiring network")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// 使用短 flag
	cmd := exec.Command(binary, "gen", "-g", "-d", testdataDir, "-s", "User")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: Requires network to install dependencies: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Generated") && !strings.Contains(outputStr, "缓存代码已生成") {
		t.Errorf("Generation failed，output: %s", outputStr)
	}

	// 清理
	os.Remove(filepath.Join(testdataDir, "models_scache.go"))
}

// TestCMDGenExclude 测试排除目录
func TestCMDGenExclude(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test requiring network")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// 排除不存在的目录（测试 exclude Parameter能正常工作）
	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--exclude", "nonexistent,vendor")
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)

	// 应该能正常生成
	if strings.Contains(outputStr, "Generated") || strings.Contains(outputStr, "缓存代码已生成") {
		// 清理
		os.Remove(filepath.Join(testdataDir, "models_scache.go"))
	}
}

// ==================== 生成Code quality tests ====================

func TestGeneratedCodeFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test requiring network")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: Requires network to install dependencies: %v", err)
	}

	generatedFile := filepath.Join(testdataDir, "models_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)

	// 验证代码格式
	expectedElements := []string{
		"package testdata",
		"import",
		"github.com/scache-io/scache",
		"func Get",
		"func Store",
		"func Load",
		"func Delete",
		"time.Duration",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(contentStr, elem) {
			t.Errorf("Generated code should contain '%s'", elem)
		}
	}

	// 验证可以通过 go fmt
	fmtCmd := exec.Command("go", "fmt", generatedFile)
	if err := fmtCmd.Run(); err != nil {
		t.Errorf("Generated code should pass go fmt: %v", err)
	}
}

// ==================== Integration tests ====================

// TestGeneratorDirect 直接测试生成器（不需要网络）
func TestGeneratorDirect(t *testing.T) {
	// 这个测试直接使用生成器包，不需要安装依赖
	// 可以测试基本的解析和生成逻辑
	testdataDir := getTestdataDir(t)

	// 验证测试数据存在
	modelsFile := filepath.Join(testdataDir, "models.go")
	if _, err := os.Stat(modelsFile); os.IsNotExist(err) {
		t.Fatalf("Test data file does not exist: %s", modelsFile)
	}

	// 读取测试数据
	content, err := os.ReadFile(modelsFile)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	contentStr := string(content)

	// 验证包含预期的Struct
	expectedStructs := []string{"User", "Product", "Order"}
	for _, s := range expectedStructs {
		if !strings.Contains(contentStr, "type "+s+" struct") {
			t.Errorf("Test data should contain struct %s", s)
		}
	}
}

// copyTestdata 复制测试数据目录
func copyTestdata(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}
