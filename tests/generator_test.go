package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scache-io/scache/cmd/scache/generator"
)

// ==================== 生成器直接测试 ====================

func TestGeneratorGeneric(t *testing.T) {
	testdataDir := getTestdataDir(t)
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	// 清理可能存在的旧文件
	os.Remove(outputFile)

	cfg := &generator.Config{
		Dir:           testdataDir,
		Package:       "models",
		ExcludeDirs:   []string{"vendor", "node_modules", ".git"},
		TargetStructs: nil, // 生成所有
		UseGeneric:    true,
	}

	err := generator.Generate(cfg)
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("生成的文件不存在: %s", outputFile)
	}

	// 读取生成的代码
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证泛型特征
	if !strings.Contains(contentStr, "[T any]") {
		t.Error("泛型代码应包含 [T any]")
	}
	if !strings.Contains(contentStr, "type Scache[T any] struct") {
		t.Error("泛型代码应包含 Scache[T any] 类型定义")
	}

	// 验证包含所有结构体
	expectedStructs := []string{"User", "Product", "Order"}
	for _, s := range expectedStructs {
		if !strings.Contains(contentStr, "Get"+s+"Scache") {
			t.Errorf("生成的代码应包含 Get%sScache", s)
		}
	}

	// 验证代码格式
	expectedElements := []string{
		"package models",
		"github.com/scache-io/scache",
		"func (s *Scache[T]) Store",
		"func (s *Scache[T]) Load",
		"func (s *Scache[T]) Delete",
		"time.Duration",
	}
	for _, elem := range expectedElements {
		if !strings.Contains(contentStr, elem) {
			t.Errorf("生成的代码应包含 '%s'", elem)
		}
	}

	// 清理
	os.Remove(outputFile)
}

func TestGeneratorClassic(t *testing.T) {
	testdataDir := getTestdataDir(t)
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	os.Remove(outputFile)

	cfg := &generator.Config{
		Dir:           testdataDir,
		Package:       "models",
		ExcludeDirs:   []string{"vendor", "node_modules", ".git"},
		TargetStructs: nil,
		UseGeneric:    false, // 传统版本
	}

	err := generator.Generate(cfg)
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证非泛型特征
	if strings.Contains(contentStr, "[T any]") {
		t.Error("传统版本代码不应包含泛型语法 [T any]")
	}
	if strings.Contains(contentStr, "type Scache[T any]") {
		t.Error("传统版本代码不应包含泛型类型定义")
	}

	// 验证包含所有结构体
	expectedStructs := []string{"User", "Product", "Order"}
	for _, s := range expectedStructs {
		if !strings.Contains(contentStr, "type "+s+"Scache struct") {
			t.Errorf("传统版本应包含 %sScache 类型定义", s)
		}
		if !strings.Contains(contentStr, "Get"+s+"Scache") {
			t.Errorf("生成的代码应包含 Get%sScache", s)
		}
	}

	os.Remove(outputFile)
}

func TestGeneratorSpecificStruct(t *testing.T) {
	testdataDir := getTestdataDir(t)
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	os.Remove(outputFile)

	cfg := &generator.Config{
		Dir:           testdataDir,
		Package:       "models",
		TargetStructs: []string{"User"}, // 只生成 User
		UseGeneric:    true,
	}

	err := generator.Generate(cfg)
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}

	content, err := os.ReadFile(outputFile)
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

	os.Remove(outputFile)
}

func TestGeneratorMultipleStructs(t *testing.T) {
	testdataDir := getTestdataDir(t)
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	os.Remove(outputFile)

	cfg := &generator.Config{
		Dir:           testdataDir,
		Package:       "models",
		TargetStructs: []string{"User", "Product"},
		UseGeneric:    true,
	}

	err := generator.Generate(cfg)
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}

	content, err := os.ReadFile(outputFile)
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

	os.Remove(outputFile)
}

// ==================== 生成代码验证测试 ====================

func TestGeneratedCodeValidation(t *testing.T) {
	testdataDir := getTestdataDir(t)
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	os.Remove(outputFile)

	// 生成泛型代码
	cfg := &generator.Config{
		Dir:        testdataDir,
		Package:    "models",
		UseGeneric: true,
	}

	err := generator.Generate(cfg)
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}

	// 验证代码格式
	fmtCmd := exec.Command("go", "fmt", outputFile)
	if output, err := fmtCmd.CombinedOutput(); err != nil {
		t.Errorf("生成的代码格式不正确: %v\n输出: %s", err, string(output))
	}

	// 读取生成的代码，验证基本结构
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证必要的导入
	if !strings.Contains(contentStr, `"github.com/scache-io/scache"`) {
		t.Error("生成的代码应导入 scache 包")
	}
	if !strings.Contains(contentStr, `"github.com/scache-io/scache/config"`) {
		t.Error("生成的代码应导入 config 包")
	}
	if !strings.Contains(contentStr, `"time"`) {
		t.Error("生成的代码应导入 time 包")
	}

	// 验证基本代码结构
	if !strings.Contains(contentStr, "func NewScache[T any]") {
		t.Error("生成的代码应包含 NewScache 函数")
	}
	if !strings.Contains(contentStr, "func (s *Scache[T]) Store") {
		t.Error("生成的代码应包含 Store 方法")
	}
	if !strings.Contains(contentStr, "func (s *Scache[T]) Load") {
		t.Error("生成的代码应包含 Load 方法")
	}

	os.Remove(outputFile)
}

func TestGeneratedCodeClassicValidation(t *testing.T) {
	testdataDir := getTestdataDir(t)
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	os.Remove(outputFile)

	// 生成传统代码
	cfg := &generator.Config{
		Dir:        testdataDir,
		Package:    "models",
		UseGeneric: false,
	}

	err := generator.Generate(cfg)
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}

	// 验证代码格式
	fmtCmd := exec.Command("go", "fmt", outputFile)
	if output, err := fmtCmd.CombinedOutput(); err != nil {
		t.Errorf("生成的代码格式不正确: %v\n输出: %s", err, string(output))
	}

	// 读取生成的代码，验证基本结构
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证必要的导入
	if !strings.Contains(contentStr, `"github.com/scache-io/scache"`) {
		t.Error("生成的代码应导入 scache 包")
	}
	if !strings.Contains(contentStr, `"github.com/scache-io/scache/config"`) {
		t.Error("生成的代码应导入 config 包")
	}
	if !strings.Contains(contentStr, `"time"`) {
		t.Error("生成的代码应导入 time 包")
	}

	// 验证传统版本的结构体特定方法
	if !strings.Contains(contentStr, "type UserScache struct") {
		t.Error("传统版本应包含 UserScache 类型定义")
	}
	if !strings.Contains(contentStr, "func (s *UserScache) Store") {
		t.Error("传统版本应包含 UserScache 的 Store 方法")
	}

	os.Remove(outputFile)
}

// ==================== 生成代码功能测试 ====================

func TestGeneratedCodeUsage(t *testing.T) {
	// 这个测试验证生成的代码API是否合理
	// 由于无法在测试中实际使用生成的代码，我们只验证生成的代码结构

	testdataDir := getTestdataDir(t)
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	os.Remove(outputFile)

	cfg := &generator.Config{
		Dir:           testdataDir,
		Package:       "models",
		TargetStructs: []string{"User"},
		UseGeneric:    true,
	}

	err := generator.Generate(cfg)
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证生成的代码包含必要的API
	requiredAPIs := []string{
		"func GetUserScache()",
		"func NewUserScache(cfg *config.EngineConfig)",
		"func (s *Scache[T]) Store(key string, obj T, ttl ...time.Duration) error",
		"func (s *Scache[T]) Load(key string) (T, error)",
		"func (s *Scache[T]) Delete(key string) bool",
		"func (s *Scache[T]) Exists(key string) bool",
		"func (s *Scache[T]) SetTTL(key string, ttl time.Duration) bool",
		"func (s *Scache[T]) GetTTL(key string) (time.Duration, bool)",
		"func (s *Scache[T]) Clear() error",
		"func (s *Scache[T]) Size() int",
		"func (s *Scache[T]) Keys() []string",
		"func (s *Scache[T]) Stats() interface{}",
	}

	for _, api := range requiredAPIs {
		if !strings.Contains(contentStr, api) {
			t.Errorf("生成的代码应包含 API: %s", api)
		}
	}

	os.Remove(outputFile)
}

// ==================== 边界条件测试 ====================

func TestGeneratorEmptyStructs(t *testing.T) {
	// 创建一个临时目录，没有结构体
	tempDir := t.TempDir()

	// 创建一个空的 go 文件
	goFile := filepath.Join(tempDir, "empty.go")
	if err := os.WriteFile(goFile, []byte("package empty\n\n// no structs\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &generator.Config{
		Dir:        tempDir,
		Package:    "empty",
		UseGeneric: true,
	}

	err := generator.Generate(cfg)
	if err == nil {
		t.Error("没有结构体时应该返回错误")
	}

	if !strings.Contains(err.Error(), "未发现任何结构体") {
		t.Errorf("错误信息应包含'未发现任何结构体': %v", err)
	}
}

func TestGeneratorNonexistentStruct(t *testing.T) {
	testdataDir := getTestdataDir(t)

	cfg := &generator.Config{
		Dir:           testdataDir,
		Package:       "models",
		TargetStructs: []string{"NonexistentStruct"},
		UseGeneric:    true,
	}

	err := generator.Generate(cfg)
	if err == nil {
		t.Error("指定不存在的结构体时应该返回错误")
	}

	if !strings.Contains(err.Error(), "未找到指定的结构体") {
		t.Errorf("错误信息应包含'未找到指定的结构体': %v", err)
	}
}

func TestGeneratorInvalidDir(t *testing.T) {
	cfg := &generator.Config{
		Dir:        "/nonexistent/path",
		Package:    "test",
		UseGeneric: true,
	}

	err := generator.Generate(cfg)
	if err == nil {
		t.Error("不存在的目录应该返回错误")
	}
}

// ==================== 性能测试 ====================

func BenchmarkGeneratorGeneric(b *testing.B) {
	testdataDir, _ := filepath.Abs("testdata")
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	cfg := &generator.Config{
		Dir:        testdataDir,
		Package:    "models",
		UseGeneric: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.Remove(outputFile)
		generator.Generate(cfg)
	}
	os.Remove(outputFile)
}

func BenchmarkGeneratorClassic(b *testing.B) {
	testdataDir, _ := filepath.Abs("testdata")
	outputFile := filepath.Join(testdataDir, "models_scache.go")

	cfg := &generator.Config{
		Dir:        testdataDir,
		Package:    "models",
		UseGeneric: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.Remove(outputFile)
		generator.Generate(cfg)
	}
	os.Remove(outputFile)
}

// ==================== 辅助函数 ====================
// getTestdataDir 定义在 cmd_test.go 中
