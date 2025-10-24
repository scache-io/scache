package generator

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed cache.tpl
var cacheTemplateContent string

//go:embed cache_generic.tpl
var cacheGenericTemplateContent string

// Config 生成器配置
type Config struct {
	Dir            string   // 扫描目录
	Package        string   // 包名
	ExcludeDirs    []string // 排除的目录
	TargetStructs  []string // 指定的结构体名称
	SplitPackages  bool     // 是否按结构体分包
	GeneratedCount int      // 生成的结构体数量
	UseGeneric     bool     // 是否使用泛型版本
}

// StructInfo 结构体信息
type StructInfo struct {
	Name   string      // 结构体名称
	Fields []FieldInfo // 字段信息
	Pkg    string      // 包名
	Source string      // 源文件路径
}

// FieldInfo 字段信息
type FieldInfo struct {
	Name string // 字段名
	Type string // 字段类型
	Tag  string // 标签
}

// Generate 执行代码生成
func Generate(config *Config) error {
	// 扫描结构体
	structs, err := scanStructs(config)
	if err != nil {
		return fmt.Errorf("扫描结构体失败: %w", err)
	}

	if len(structs) == 0 {
		return fmt.Errorf("未发现任何结构体")
	}

	// 过滤指定的结构体
	if len(config.TargetStructs) > 0 {
		filtered := make([]StructInfo, 0)
		structMap := make(map[string]StructInfo)
		for _, s := range structs {
			structMap[s.Name] = s
		}
		for _, target := range config.TargetStructs {
			if s, exists := structMap[target]; exists {
				filtered = append(filtered, s)
			}
		}
		structs = filtered
	}

	if len(structs) == 0 {
		return fmt.Errorf("未找到指定的结构体")
	}

	// 记录生成的结构体数量
	config.GeneratedCount = len(structs)

	// 直接在同目录生成文件，不再分包
	return generateInPlace(config, structs)
}

// scanStructs 扫描目录中的所有结构体
func scanStructs(config *Config) ([]StructInfo, error) {
	var structs []StructInfo
	fset := token.NewFileSet()

	err := filepath.Walk(config.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			// 检查是否需要排除
			for _, exclude := range config.ExcludeDirs {
				if strings.Contains(path, exclude) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// 只处理.go文件
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// 跳过测试文件
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 解析Go文件
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // 忽略解析错误的文件
		}

		// 提取结构体
		fileStructs := extractStructs(file, path)
		structs = append(structs, fileStructs...)

		return nil
	})

	return structs, err
}

// extractStructs 从AST中提取结构体
func extractStructs(file *ast.File, sourcePath string) []StructInfo {
	var structs []StructInfo

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 提取字段信息
			var fields []FieldInfo
			if structType.Fields != nil {
				for _, field := range structType.Fields.List {
					fieldInfo := FieldInfo{}

					// 字段名
					if len(field.Names) > 0 {
						fieldInfo.Name = field.Names[0].Name
					} else {
						// 匿名字段
						fieldInfo.Name = ""
					}

					// 字段类型
					fieldInfo.Type = fieldTypeToString(field.Type)

					// 标签
					if field.Tag != nil {
						fieldInfo.Tag = strings.Trim(field.Tag.Value, "`")
					}

					fields = append(fields, fieldInfo)
				}
			}

			structs = append(structs, StructInfo{
				Name:   typeSpec.Name.Name,
				Fields: fields,
				Pkg:    file.Name.Name,
				Source: sourcePath,
			})
		}
	}

	return structs
}

// fieldTypeToString 将字段类型转换为字符串
func fieldTypeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", fieldTypeToString(t.X), t.Sel.Name)
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", fieldTypeToString(t.Elt))
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", fieldTypeToString(t.X))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", fieldTypeToString(t.Key), fieldTypeToString(t.Value))
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return fmt.Sprintf("...%s", fieldTypeToString(t.Elt))
	default:
		return "unknown"
	}
}

// generateInPlace 在原文件同目录下生成_scache.go文件
func generateInPlace(config *Config, structs []StructInfo) error {
	// 按包分组结构体
	packageGroups := make(map[string][]StructInfo)
	for _, structInfo := range structs {
		packageGroups[structInfo.Pkg] = append(packageGroups[structInfo.Pkg], structInfo)
	}

	// 为每个包生成_scache.go文件
	for pkgName, pkgStructs := range packageGroups {
		if err := generatePackageScache(config, pkgName, pkgStructs); err != nil {
			return err
		}
	}

	return nil
}

// generatePackageScache 为指定包生成_scache.go文件
func generatePackageScache(config *Config, pkgName string, structs []StructInfo) error {
	// 找到该包的第一个结构体所在目录
	targetDir := findPackageDirectory(pkgName, structs)
	if targetDir == "" {
		return fmt.Errorf("找不到包 %s 的目录", pkgName)
	}

	// 生成文件名
	filename := filepath.Join(targetDir, pkgName+"_scache.go")

	// 生成包代码
	content, err := generatePackageCode(pkgName, structs, config.UseGeneric)
	if err != nil {
		return fmt.Errorf("生成代码失败: %w", err)
	}

	// 写入文件
	return generatePackageFile(filename, content)
}

// findPackageDirectory 找到包的目录
func findPackageDirectory(pkgName string, structs []StructInfo) string {
	for _, structInfo := range structs {
		if structInfo.Pkg == pkgName {
			return filepath.Dir(structInfo.Source)
		}
	}
	return ""
}

// TemplateData 模板数据结构
type TemplateData struct {
	Package string
	Structs []StructInfo
}

// loadTemplate 加载模板文件
func loadTemplate(useGeneric bool) (*template.Template, error) {
	templateName := "cache"
	templateContent := cacheTemplateContent

	if useGeneric {
		templateName = "cache_generic"
		templateContent = cacheGenericTemplateContent
	}

	return template.New(templateName).Parse(templateContent)
}

// generatePackageCode 为指定包生成缓存代码
func generatePackageCode(pkgName string, structs []StructInfo, useGeneric bool) (string, error) {
	// 加载嵌入的模板
	tmpl, err := loadTemplate(useGeneric)
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %w", err)
	}

	data := TemplateData{
		Package: pkgName,
		Structs: structs,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}

// generatePackageFile 生成包文件
func generatePackageFile(filePath, content string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 检查文件是否已存在，如果存在则先删除
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("🗑️  删除现有缓存文件: %s\n", filePath)
		if removeErr := os.Remove(filePath); removeErr != nil {
			return fmt.Errorf("删除现有文件失败: %w", removeErr)
		}
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Printf("✅ 生成缓存文件: %s\n", filePath)
	return nil
}
