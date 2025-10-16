package types

import (
	"testing"
	"time"
)

func TestStringObject(t *testing.T) {
	// 创建字符串对象
	strObj := NewStringObject("hello", time.Minute)

	// 测试类型
	if strObj.Type() != "string" {
		t.Errorf("Expected 'string', got '%s'", strObj.Type())
	}

	// 测试值
	if strObj.Value() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", strObj.Value())
	}

	// 测试设置新值
	strObj.Set("world")
	if strObj.Value() != "world" {
		t.Errorf("Expected 'world', got '%s'", strObj.Value())
	}

	// 测试大小
	if strObj.Size() != 5 {
		t.Errorf("Expected size 5, got %d", strObj.Size())
	}

	// 测试过期
	expiredObj := NewStringObject("expire", time.Millisecond*10)
	if expiredObj.IsExpired() {
		t.Error("Object should not be expired initially")
	}

	time.Sleep(time.Millisecond * 20)
	if !expiredObj.IsExpired() {
		t.Error("Object should be expired")
	}
}

func TestListObject(t *testing.T) {
	values := []interface{}{"item1", "item2", "item3"}
	listObj := NewListObject(values, time.Minute)

	// 测试类型
	if listObj.Type() != "list" {
		t.Errorf("Expected 'list', got '%s'", listObj.Type())
	}

	// 测试初始值
	initialValues := listObj.Values()
	if len(initialValues) != 3 {
		t.Errorf("Expected 3 items, got %d", len(initialValues))
	}

	// 测试Push
	listObj.Push("item4")
	newValues := listObj.Values()
	if len(newValues) != 4 {
		t.Errorf("Expected 4 items after push, got %d", len(newValues))
	}

	// 测试Pop
	popped, exists := listObj.Pop()
	if !exists {
		t.Error("Pop should succeed")
	}
	if popped != "item4" {
		t.Errorf("Expected 'item4', got '%v'", popped)
	}

	// 测试Index
	item, exists := listObj.Index(1)
	if !exists {
		t.Error("Index should exist")
	}
	if item != "item2" {
		t.Errorf("Expected 'item2', got '%v'", item)
	}

	// 测试Range
	rangeItems := listObj.Range(0, 1)
	if len(rangeItems) != 2 {
		t.Errorf("Expected 2 items in range, got %d", len(rangeItems))
	}

	// 测试Len
	if listObj.Len() != 3 {
		t.Errorf("Expected length 3, got %d", listObj.Len())
	}

	// 测试空列表Pop
	emptyList := NewListObject([]interface{}{}, time.Minute)
	_, exists = emptyList.Pop()
	if exists {
		t.Error("Pop from empty list should fail")
	}
}

func TestHashObject(t *testing.T) {
	fields := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
	}
	hashObj := NewHashObject(fields, time.Minute)

	// 测试类型
	if hashObj.Type() != "hash" {
		t.Errorf("Expected 'hash', got '%s'", hashObj.Type())
	}

	// 测试初始字段
	initialFields := hashObj.Fields()
	if len(initialFields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(initialFields))
	}

	// 测试Get
	name, exists := hashObj.Get("name")
	if !exists {
		t.Error("Field 'name' should exist")
	}
	if name != "Alice" {
		t.Errorf("Expected 'Alice', got '%v'", name)
	}

	// 测试Set
	hashObj.Set("city", "Beijing")
	city, exists := hashObj.Get("city")
	if !exists {
		t.Error("Field 'city' should exist after set")
	}
	if city != "Beijing" {
		t.Errorf("Expected 'Beijing', got '%v'", city)
	}

	// 测试Delete
	deleted := hashObj.Delete("age")
	if !deleted {
		t.Error("Delete should succeed")
	}

	_, exists = hashObj.Get("age")
	if exists {
		t.Error("Field 'age' should not exist after delete")
	}

	// 测试Len
	if hashObj.Len() != 3 {
		t.Errorf("Expected length 3, got %d", hashObj.Len())
	}

	// 测试Size
	size := hashObj.Size()
	if size <= 0 {
		t.Error("Size should be positive")
	}
}

func TestBaseObject(t *testing.T) {
	// 测试永不过期的对象
	neverExpire := NewBaseObject("string", 0)
	if neverExpire.IsExpired() {
		t.Error("Object with zero TTL should never expire")
	}

	// 测试会过期的对象
	willExpire := NewBaseObject("string", time.Millisecond*10)
	if willExpire.IsExpired() {
		t.Error("Object should not be expired initially")
	}

	time.Sleep(time.Millisecond * 20)
	if !willExpire.IsExpired() {
		t.Error("Object should be expired")
	}

	// 测试UpdateAccess
	before := willExpire.CreatedAt()
	willExpire.UpdateAccess()
	after := willExpire.CreatedAt()
	if before != after {
		t.Error("CreatedAt should not change when updating access")
	}
}
