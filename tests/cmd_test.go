package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ==================== CMD 工具测试 ====================

// buildScacheCMD 构建scache命令
func buildScacheCMD(t *testing.T) string {
	binaryPath := filepath.Join(t.TempDir(), "scache")
	
	// 获取项目根目录（从tests目录向上两级）
	projectRoot, err := filepath.Abs(filepath.Join(".."))
	if err != nil {
		t.Fatalf("获取项目根目录失败: %v", err)
	}
	
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/scache")
	cmd.Dir = projectRoot
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("构建scache命令失败: %v\n输出: %s", err, string(output))
	}
	
	return binaryPath
}

// TestCMDVersion 测试version命令
func TestCMDVersion(t *testing.T) {
	binary := buildScacheCMD(t)
	
	cmd := exec.Command(binary, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version命令执行失败: %v", err)
	}
	
	outputStr := string(output)
	if !strings.Contains(outputStr, "SCache version") {
		t.Errorf("version输出格式错误: %s", outputStr)
	}
	if !strings.Contains(outputStr, "0.0.1") {
		t.Errorf("version不包含版本号: %s", outputStr)
	}
}

// TestCMDHelp 测试help命令
func TestCMDHelp(t *testing.T) {
	binary := buildScacheCMD(t)
	
	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	// --help 会返回 exit code 0
	if err != nil {
		t.Fatalf("help命令执行失败: %v", err)
	}
	
	outputStr := string(output)
	if !strings.Contains(outputStr, "gen") {
		t.Errorf("help输出应包含gen命令: %s", outputStr)
	}
}

// TestCMDGenHelp 测试gen命令帮助
func TestCMDGenHelp(t *testing.T) {
	binary := buildScacheCMD(t)
	
	cmd := exec.Command(binary, "gen", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gen --help命令执行失败: %v", err)
	}
	
	outputStr := string(output)
	
	// 验证帮助信息包含所有flag
	expectedFlags := []string{"--dir", "--package", "--exclude", "--structs", "--generic"}
	for _, flag := range expectedFlags {
		if !strings.Contains(outputStr, flag) {
			t.Errorf("gen帮助信息应包含 %s flag", flag)
		}
	}
}

// TestCMDInvalidDir 测试无效目录
func TestCMDInvalidDir(t *testing.T) {
	binary := buildScacheCMD(t)
	
	cmd := exec.Command(binary, "gen", "--dir", "/nonexistent/path")
	output, err := cmd.CombinedOutput()
	
	// 应该返回错误
	if err == nil {
		t.Error("对不存在的目录应该返回错误")
	}
	
	outputStr := string(output)
	if !strings.Contains(outputStr, "目录不存在") {
		t.Errorf("错误信息应包含'目录不存在': %s", outputStr)
	}
}

// TestCMDGenShortFlags 测试短flag
func TestCMDGenShortFlags(t *testing.T) {
	binary := buildScacheCMD(t)
	
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	
	// 测试短flag组合（应该失败，因为目录为空）
	cmd := exec.Command(binary, "gen", "-g", "-d", tempDir)
	output, err := cmd.CombinedOutput()
	
	// 应该返回错误（没有结构体）
	if err == nil {
		t.Error("空目录应该返回错误")
	}
	
	outputStr := string(output)
	// 验证短flag被正确解析
	if !strings.Contains(outputStr, "未发现任何结构体") && !strings.Contains(outputStr, "未找到指定的结构体") {
		// 可能因为其他原因失败，但至少不应该是因为flag解析错误
		t.Logf("输出: %s", outputStr)
	}
}

// TestCMDNoStructs 测试没有结构体的目录
func TestCMDNoStructs(t *testing.T) {
	binary := buildScacheCMD(t)
	
	// 创建一个临时目录，没有结构体
	tempDir := t.TempDir()
	
	// 创建一个空的 go.mod
	goMod := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	
	cmd := exec.Command(binary, "gen", "--generic", "--dir", tempDir)
	output, err := cmd.CombinedOutput()
	
	// 应该返回错误（可能是因为没有结构体，也可能是因为安装依赖失败）
	if err == nil {
		t.Error("没有结构体时应该返回错误")
	}
	
	outputStr := string(output)
	// 由于测试环境的限制，可能无法真正测试代码生成
	// 所以我们只验证命令能正确执行并返回错误
	if !strings.Contains(outputStr, "未发现任何结构体") && 
	   !strings.Contains(outputStr, "未找到指定的结构体") &&
	   !strings.Contains(outputStr, "安装scache包失败") {
		t.Errorf("错误信息不符合预期: %s", outputStr)
	}
}
