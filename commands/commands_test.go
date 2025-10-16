package commands

import (
	"testing"
	"time"

	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/storage"
	"github.com/scache-io/scache/types"
)

func createTestEngine() interfaces.StorageEngine {
	config := &storage.EngineConfig{
		MaxSize:           100,
		DefaultExpiration: 0,
	}
	return storage.NewStorageEngine(config)
}

func TestSetCommand(t *testing.T) {
	engine := createTestEngine()
	cmd := NewSetCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"key1", "value1", time.Minute},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// 验证结果
	obj, exists := engine.Get("key1")
	if !exists {
		t.Fatal("Key not found")
	}

	strObj, ok := obj.(*types.StringObject)
	if !ok {
		t.Fatal("Object is not StringObject")
	}

	if strObj.Value() != "value1" {
		t.Errorf("Expected 'value1', got '%s'", strObj.Value())
	}
}

func TestGetCommand(t *testing.T) {
	engine := createTestEngine()

	// 先设置一个值
	obj := types.NewStringObject("test_value", time.Minute)
	engine.Set("test_key", obj)

	cmd := NewGetCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"test_key"},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", ctx.Result)
	}
}

func TestLPushCommand(t *testing.T) {
	engine := createTestEngine()
	cmd := NewLPushCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"list_key", "item1", time.Minute},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != 1 {
		t.Errorf("Expected length 1, got %v", ctx.Result)
	}
}

func TestHSetCommand(t *testing.T) {
	engine := createTestEngine()
	cmd := NewHSetCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"hash_key", "field1", "value1", time.Minute},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != 1 {
		t.Errorf("Expected result 1, got %v", ctx.Result)
	}

	// 验证哈希对象
	obj, exists := engine.Get("hash_key")
	if !exists {
		t.Fatal("Hash key not found")
	}

	hashObj, ok := obj.(*types.HashObject)
	if !ok {
		t.Fatal("Object is not HashObject")
	}

	value, exists := hashObj.Get("field1")
	if !exists {
		t.Fatal("Field not found")
	}

	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

func TestDeleteCommand(t *testing.T) {
	engine := createTestEngine()

	// 先设置一个值
	obj := types.NewStringObject("delete_me", time.Minute)
	engine.Set("delete_key", obj)

	cmd := NewDeleteCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"delete_key"},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != true {
		t.Error("Expected true, got false")
	}

	// 验证删除
	_, exists := engine.Get("delete_key")
	if exists {
		t.Error("Key still exists after deletion")
	}
}

func TestExistsCommand(t *testing.T) {
	engine := createTestEngine()

	// 先设置一个值
	obj := types.NewStringObject("exists_test", time.Minute)
	engine.Set("exists_key", obj)

	cmd := NewExistsCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"exists_key"},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != true {
		t.Error("Expected true, got false")
	}

	// 测试不存在的键
	ctx.Args = []interface{}{"nonexistent_key"}
	err = cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != false {
		t.Error("Expected false, got true")
	}
}

func TestTypeCommand(t *testing.T) {
	engine := createTestEngine()

	// 测试字符串类型
	strObj := types.NewStringObject("test", time.Minute)
	engine.Set("str_key", strObj)

	cmd := NewTypeCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"str_key"},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != "string" {
		t.Errorf("Expected 'string', got '%v'", ctx.Result)
	}

	// 测试不存在的键
	ctx.Args = []interface{}{"nonexistent_key"}
	err = cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != "none" {
		t.Errorf("Expected 'none', got '%v'", ctx.Result)
	}
}

func TestExpireCommand(t *testing.T) {
	engine := createTestEngine()

	// 先设置一个值（永不过期）
	obj := types.NewStringObject("expire_test", 0)
	engine.Set("expire_key", obj)

	cmd := NewExpireCommand()
	ctx := &interfaces.Context{
		Storage: engine,
		Args:    []interface{}{"expire_key", time.Minute},
	}

	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if ctx.Result != true {
		t.Error("Expected true, got false")
	}

	// 验证设置了过期时间
	ttl, exists := engine.TTL("expire_key")
	if !exists {
		t.Fatal("Key not found")
	}

	if ttl <= 0 {
		t.Error("TTL should be positive")
	}
}

func TestCommandRegistry(t *testing.T) {
	registry := NewCommandRegistry()

	// 测试注册命令
	testCmd := NewSetCommand()
	registry.Register(testCmd)

	// 测试获取命令
	cmd, exists := registry.Get("SET")
	if !exists {
		t.Error("Command not found")
	}

	if cmd.Name() != "SET" {
		t.Errorf("Expected 'SET', got '%s'", cmd.Name())
	}

	// 测试大小写不敏感
	cmd, exists = registry.Get("set")
	if !exists {
		t.Error("Command not found (case insensitive)")
	}

	// 测试列出命令
	commands := registry.List()
	if len(commands) == 0 {
		t.Error("No commands found")
	}
}
