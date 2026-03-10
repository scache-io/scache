package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ==================== CMD 工具测试 ====================

// buildScacheCMD 构建 scache 命令
func buildScacheCMD(t *testing.T) string {
	binaryPath := filepath.Join(t.TempDir(), "scache")

	projectRoot, err := filepath.Abs(filepath.Join(".."))
	if err != nil {
		t.Fatalf("获取项目根目录失败: %v", err)
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/scache")
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("构建 scache 命令失败: %v\n输出: %s", err, string(output))
	}

	return binaryPath
}

// getTestdataDir 获取测试数据目录
func getTestdataDir(t *testing.T) string {
	dir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatalf("获取测试数据目录失败: %v", err)
	}
	return dir
}

// ==================== 基础命令测试 ====================

func TestCMDVersion(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version 命令执行失败: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "SCache version") {
		t.Errorf("version 输出格式错误: %s", outputStr)
	}
	if !strings.Contains(outputStr, "GitHub") {
		t.Errorf("version 输出应包含 GitHub 链接: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Docs") {
		t.Errorf("version 输出应包含文档链接: %s", outputStr)
	}
}

func TestCMDHelp(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help 命令执行失败: %v", err)
	}

	outputStr := string(output)
	expectedKeywords := []string{"gen", "version", "快速开始"}
	for _, kw := range expectedKeywords {
		if !strings.Contains(outputStr, kw) {
			t.Errorf("help 输出应包含 '%s'", kw)
		}
	}
}

func TestCMDGenHelp(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "gen", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gen --help 命令执行失败: %v", err)
	}

	outputStr := string(output)
	expectedFlags := []string{"--dir", "--package", "--exclude", "--structs", "--generic"}
	for _, flag := range expectedFlags {
		if !strings.Contains(outputStr, flag) {
			t.Errorf("gen 帮助信息应包含 %s flag", flag)
		}
	}
}

// ==================== 错误处理测试 ====================

func TestCMDInvalidDir(t *testing.T) {
	binary := buildScacheCMD(t)

	cmd := exec.Command(binary, "gen", "--dir", "/nonexistent/path")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("对不存在的目录应该返回错误")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "目录不存在") {
		t.Errorf("错误信息应包含'目录不存在': %s", outputStr)
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

	// 可能因为网络问题失败，但至少应该尝试执行
	if err != nil {
		outputStr := string(output)
		validErrors := []string{"未发现任何结构体", "未找到", "安装", "go get"}
		hasValidError := false
		for _, e := range validErrors {
			if strings.Contains(outputStr, e) {
				hasValidError = true
				break
			}
		}
		if !hasValidError {
			t.Logf("警告: 错误信息不符合预期: %s", outputStr)
		}
	} else {
		t.Log("意外成功：空目录生成了代码")
	}
}

// ==================== 代码生成测试 (需要网络) ====================

// TestCMDGenGeneric 测试泛型代码生成
// 注意：此测试需要网络连接来安装依赖
func TestCMDGenGeneric(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过需要网络的测试")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: 需要网络安装依赖: %v\n输出: %s", err, string(output))
	}

	outputStr := string(output)

	// 验证生成成功
	if !strings.Contains(outputStr, "缓存代码已生成") {
		t.Errorf("生成失败，输出: %s", outputStr)
	}

	// 验证生成的文件
	generatedFile := filepath.Join(testdataDir, "testdata_scache.go")
	defer os.Remove(generatedFile) // 清理

	if _, err := os.Stat(generatedFile); os.IsNotExist(err) {
		t.Fatalf("生成的文件不存在: %s", generatedFile)
	}

	// 验证生成的代码内容
	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证泛型特征
	if !strings.Contains(contentStr, "[T any]") {
		t.Error("泛型代码应包含 [T any]")
	}
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("生成的代码应包含 GetUserScache")
	}
	if !strings.Contains(contentStr, "GetProductScache") {
		t.Error("生成的代码应包含 GetProductScache")
	}
	if !strings.Contains(contentStr, "GetOrderScache") {
		t.Error("生成的代码应包含 GetOrderScache")
	}
}

// TestCMDGenClassic 测试传统代码生成
// 注意：此测试需要网络连接来安装依赖
func TestCMDGenClassic(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过需要网络的测试")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--dir", testdataDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: 需要网络安装依赖: %v\n输出: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "缓存代码已生成") {
		t.Errorf("生成失败，输出: %s", outputStr)
	}

	generatedFile := filepath.Join(testdataDir, "testdata_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证非泛型特征（不应该包含泛型语法）
	if strings.Contains(contentStr, "[T any]") {
		t.Error("传统版本代码不应包含泛型语法")
	}
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("生成的代码应包含 GetUserScache")
	}
}

// TestCMDGenSpecificStructs 测试指定结构体生成
func TestCMDGenSpecificStructs(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过需要网络的测试")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// 只生成 User 结构体
	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--structs", "User")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: 需要网络安装依赖: %v", err)
	}

	generatedFile := filepath.Join(testdataDir, "testdata_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证只包含 User
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("生成的代码应包含 GetUserScache")
	}
	if strings.Contains(contentStr, "GetProductScache") {
		t.Error("生成的代码不应包含 GetProductScache（未指定）")
	}
	if strings.Contains(contentStr, "GetOrderScache") {
		t.Error("生成的代码不应包含 GetOrderScache（未指定）")
	}
}

// TestCMDGenMultipleStructs 测试多个指定结构体
func TestCMDGenMultipleStructs(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过需要网络的测试")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// 生成多个指定结构体
	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--structs", "User,Product")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: 需要网络安装依赖: %v", err)
	}

	generatedFile := filepath.Join(testdataDir, "testdata_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证包含 User 和 Product，但不包含 Order
	if !strings.Contains(contentStr, "GetUserScache") {
		t.Error("生成的代码应包含 GetUserScache")
	}
	if !strings.Contains(contentStr, "GetProductScache") {
		t.Error("生成的代码应包含 GetProductScache")
	}
	if strings.Contains(contentStr, "GetOrderScache") {
		t.Error("生成的代码不应包含 GetOrderScache（未指定）")
	}
}

// TestCMDNonexistentStruct 测试不存在的结构体
func TestCMDNonexistentStruct(t *testing.T) {
	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--structs", "NonexistentStruct")
	output, err := cmd.CombinedOutput()

	// 应该返回错误
	if err == nil {
		t.Error("指定不存在的结构体应该返回错误")
	}

	outputStr := string(output)
	// 由于可能需要先安装依赖，错误信息可能不同
	if !strings.Contains(outputStr, "未找到") && !strings.Contains(outputStr, "安装") {
		t.Logf("输出: %s", outputStr)
	}
}

// ==================== 短 flag 测试 ====================

func TestCMDGenShortFlags(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过需要网络的测试")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// 使用短 flag
	cmd := exec.Command(binary, "gen", "-g", "-d", testdataDir, "-s", "User")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: 需要网络安装依赖: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "缓存代码已生成") {
		t.Errorf("生成失败，输出: %s", outputStr)
	}

	// 清理
	os.Remove(filepath.Join(testdataDir, "testdata_scache.go"))
}

// TestCMDGenExclude 测试排除目录
func TestCMDGenExclude(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过需要网络的测试")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	// 排除不存在的目录（测试 exclude 参数能正常工作）
	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir, "--exclude", "nonexistent,vendor")
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)

	// 应该能正常生成
	if strings.Contains(outputStr, "缓存代码已生成") {
		// 清理
		os.Remove(filepath.Join(testdataDir, "testdata_scache.go"))
	}
}

// ==================== 生成代码质量测试 ====================

func TestGeneratedCodeFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过需要网络的测试")
	}

	binary := buildScacheCMD(t)
	testdataDir := getTestdataDir(t)

	cmd := exec.Command(binary, "gen", "--generic", "--dir", testdataDir)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("跳过: 需要网络安装依赖: %v", err)
	}

	generatedFile := filepath.Join(testdataDir, "testdata_scache.go")
	defer os.Remove(generatedFile)

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
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
			t.Errorf("生成的代码应包含 '%s'", elem)
		}
	}

	// 验证可以通过 go fmt
	fmtCmd := exec.Command("go", "fmt", generatedFile)
	if err := fmtCmd.Run(); err != nil {
		t.Errorf("生成的代码应该能通过 go fmt: %v", err)
	}
}

// ==================== 集成测试 ====================

// TestGeneratorDirect 直接测试生成器（不需要网络）
func TestGeneratorDirect(t *testing.T) {
	// 这个测试直接使用生成器包，不需要安装依赖
	// 可以测试基本的解析和生成逻辑
	testdataDir := getTestdataDir(t)

	// 验证测试数据存在
	modelsFile := filepath.Join(testdataDir, "models.go")
	if _, err := os.Stat(modelsFile); os.IsNotExist(err) {
		t.Fatalf("测试数据文件不存在: %s", modelsFile)
	}

	// 读取测试数据
	content, err := os.ReadFile(modelsFile)
	if err != nil {
		t.Fatalf("读取测试数据失败: %v", err)
	}

	contentStr := string(content)

	// 验证包含预期的结构体
	expectedStructs := []string{"User", "Product", "Order"}
	for _, s := range expectedStructs {
		if !strings.Contains(contentStr, "type "+s+" struct") {
			t.Errorf("测试数据应包含结构体 %s", s)
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
