package storage

import (
	"testing"
	"time"

	"scache/interfaces"
	"scache/types"
)

func TestNewStorageEngine(t *testing.T) {
	config := &EngineConfig{
		MaxSize:                   100,
		DefaultExpiration:         time.Hour,
		MemoryThreshold:           0.8,
		BackgroundCleanupInterval: time.Minute,
	}

	engine := NewStorageEngine(config)
	if engine == nil {
		t.Fatal("Engine is nil")
	}

	// 验证初始状态
	if engine.Size() != 0 {
		t.Errorf("Expected size 0, got %d", engine.Size())
	}
}

func TestStringOperations(t *testing.T) {
	engine := NewStorageEngine(nil)
	strObj := types.NewStringObject("hello", time.Minute)

	// Test Set
	err := engine.Set("key1", strObj)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	obj, exists := engine.Get("key1")
	if !exists {
		t.Fatal("Key not found")
	}

	retrievedStr, ok := obj.(*types.StringObject)
	if !ok {
		t.Fatal("Object is not StringObject")
	}

	if retrievedStr.Value() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", retrievedStr.Value())
	}

	// Test Type
	dataType, exists := engine.Type("key1")
	if !exists {
		t.Fatal("Type check failed")
	}

	if dataType != interfaces.DataTypeString {
		t.Errorf("Expected 'string', got '%s'", dataType)
	}
}

func TestListOperations(t *testing.T) {
	engine := NewStorageEngine(nil)
	values := []interface{}{"item1", "item2", "item3"}
	listObj := types.NewListObject(values, time.Minute)

	// Test Set
	err := engine.Set("list1", listObj)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	obj, exists := engine.Get("list1")
	if !exists {
		t.Fatal("Key not found")
	}

	retrievedList, ok := obj.(*types.ListObject)
	if !ok {
		t.Fatal("Object is not ListObject")
	}

	retrievedValues := retrievedList.Values()
	if len(retrievedValues) != 3 {
		t.Errorf("Expected 3 items, got %d", len(retrievedValues))
	}

	// Test Type
	dataType, exists := engine.Type("list1")
	if !exists {
		t.Fatal("Type check failed")
	}

	if dataType != interfaces.DataTypeList {
		t.Errorf("Expected 'list', got '%s'", dataType)
	}
}

func TestHashOperations(t *testing.T) {
	engine := NewStorageEngine(nil)
	fields := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
	}
	hashObj := types.NewHashObject(fields, time.Minute)

	// Test Set
	err := engine.Set("hash1", hashObj)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	obj, exists := engine.Get("hash1")
	if !exists {
		t.Fatal("Key not found")
	}

	retrievedHash, ok := obj.(*types.HashObject)
	if !ok {
		t.Fatal("Object is not HashObject")
	}

	retrievedFields := retrievedHash.Fields()
	if len(retrievedFields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(retrievedFields))
	}

	// Test Type
	dataType, exists := engine.Type("hash1")
	if !exists {
		t.Fatal("Type check failed")
	}

	if dataType != interfaces.DataTypeHash {
		t.Errorf("Expected 'hash', got '%s'", dataType)
	}
}

func TestExpiration(t *testing.T) {
	engine := NewStorageEngine(nil)

	// Test with expiration
	strObj := types.NewStringObject("expire_me", time.Millisecond*100)
	err := engine.Set("expire_key", strObj)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Should exist initially
	_, exists := engine.Get("expire_key")
	if !exists {
		t.Error("Key should exist initially")
	}

	// Wait for expiration
	time.Sleep(time.Millisecond * 150)

	// Should be expired
	_, exists = engine.Get("expire_key")
	if exists {
		t.Error("Key should be expired")
	}
}

func TestTTL(t *testing.T) {
	engine := NewStorageEngine(nil)

	// Test with TTL
	strObj := types.NewStringObject("ttl_test", time.Minute)
	err := engine.Set("ttl_key", strObj)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test TTL
	ttl, exists := engine.TTL("ttl_key")
	if !exists {
		t.Fatal("Key not found")
	}

	if ttl <= 0 {
		t.Error("TTL should be positive")
	}

	// Test Expire
	success := engine.Expire("ttl_key", time.Second*30)
	if !success {
		t.Error("Expire should succeed")
	}

	// Verify new TTL
	newTTL, exists := engine.TTL("ttl_key")
	if !exists {
		t.Fatal("Key not found")
	}

	if newTTL <= time.Second*25 || newTTL > time.Second*30 {
		t.Errorf("TTL should be around 30 seconds, got %v", newTTL)
	}
}

func TestDelete(t *testing.T) {
	engine := NewStorageEngine(nil)
	strObj := types.NewStringObject("delete_me", time.Minute)

	// Set key
	err := engine.Set("delete_key", strObj)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify exists
	_, exists := engine.Get("delete_key")
	if !exists {
		t.Fatal("Key should exist")
	}

	// Delete key
	deleted := engine.Delete("delete_key")
	if !deleted {
		t.Error("Delete should succeed")
	}

	// Verify deleted
	_, exists = engine.Get("delete_key")
	if exists {
		t.Error("Key should be deleted")
	}
}

func TestExists(t *testing.T) {
	engine := NewStorageEngine(nil)
	strObj := types.NewStringObject("exists_test", time.Minute)

	// Test non-existent key
	exists := engine.Exists("nonexistent")
	if exists {
		t.Error("Key should not exist")
	}

	// Set key
	err := engine.Set("exists_key", strObj)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test existing key
	exists = engine.Exists("exists_key")
	if !exists {
		t.Error("Key should exist")
	}
}

func TestFlush(t *testing.T) {
	engine := NewStorageEngine(nil)

	// Add multiple keys
	keys := []string{"key1", "key2", "key3"}
	for i, key := range keys {
		strObj := types.NewStringObject("value"+string(rune('1'+i)), time.Minute)
		err := engine.Set(key, strObj)
		if err != nil {
			t.Fatalf("Set failed for %s: %v", key, err)
		}
	}

	// Verify keys exist
	if engine.Size() != 3 {
		t.Errorf("Expected size 3, got %d", engine.Size())
	}

	// Flush
	err := engine.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// Verify all keys are gone
	if engine.Size() != 0 {
		t.Errorf("Expected size 0 after flush, got %d", engine.Size())
	}

	for _, key := range keys {
		_, exists := engine.Get(key)
		if exists {
			t.Errorf("Key %s should not exist after flush", key)
		}
	}
}

func TestKeys(t *testing.T) {
	engine := NewStorageEngine(nil)

	// Initially empty
	keys := engine.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(keys))
	}

	// Add keys
	testKeys := []string{"a", "b", "c"}
	for _, key := range testKeys {
		strObj := types.NewStringObject("value", time.Minute)
		err := engine.Set(key, strObj)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
	}

	// Get keys
	keys = engine.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Verify all test keys are present
	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}

	for _, testKey := range testKeys {
		if !keySet[testKey] {
			t.Errorf("Key %s not found in keys list", testKey)
		}
	}
}

func TestStats(t *testing.T) {
	engine := NewStorageEngine(nil)
	strObj := types.NewStringObject("stats_test", time.Minute)

	// Initial stats
	stats := engine.Stats()
	statsMap := stats.(map[string]interface{})
	if hits, ok := statsMap["hits"].(int64); ok && hits != 0 {
		t.Errorf("Expected 0 hits, got %d", hits)
	}

	// Perform operations
	engine.Set("stats_key", strObj) // set
	engine.Get("stats_key")         // hit
	engine.Get("nonexistent")       // miss
	engine.Delete("stats_key")      // delete

	// Check stats
	stats = engine.Stats()
	statsMap = stats.(map[string]interface{})

	if hits, ok := statsMap["hits"].(int64); ok && hits != 1 {
		t.Errorf("Expected 1 hit, got %d", hits)
	}

	if misses, ok := statsMap["misses"].(int64); ok && misses != 1 {
		t.Errorf("Expected 1 miss, got %d", misses)
	}

	if sets, ok := statsMap["sets"].(int64); ok && sets != 1 {
		t.Errorf("Expected 1 set, got %d", sets)
	}

	if deletes, ok := statsMap["deletes"].(int64); ok && deletes != 1 {
		t.Errorf("Expected 1 delete, got %d", deletes)
	}
}
