package scache

import (
	"testing"
	"time"
)

func TestGlobalCacheStructOperations(t *testing.T) {
	// 清空全局缓存
	Flush()

	// 测试全局缓存存储和读取结构体
	type User struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Age       int       `json:"age"`
		CreatedAt time.Time `json:"created_at"`
	}

	// 存储结构体
	user := User{
		ID:        1,
		Name:      "测试用户",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: time.Now(),
	}
	err := Store("user:1", user, time.Hour)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// 读取结构体
	var loadedUser User
	err = Load("user:1", &loadedUser)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// 验证字段
	if loadedUser.ID != user.ID {
		t.Errorf("ID mismatch: expected %d, got %d", user.ID, loadedUser.ID)
	}
	if loadedUser.Name != user.Name {
		t.Errorf("Name mismatch: expected %s, got %s", user.Name, loadedUser.Name)
	}
	if loadedUser.Email != user.Email {
		t.Errorf("Email mismatch: expected %s, got %s", user.Email, loadedUser.Email)
	}
	if loadedUser.Age != user.Age {
		t.Errorf("Age mismatch: expected %d, got %d", user.Age, loadedUser.Age)
	}
}

func TestGlobalCacheMultipleOperations(t *testing.T) {
	// 清空全局缓存
	Flush()

	// 测试多种数据类型
	// 字符串
	SetString("str:key", "value", time.Hour)
	val, exists := GetString("str:key")
	if !exists || val != "value" {
		t.Error("String operation failed")
	}

	// 列表
	SetList("list:key", []interface{}{1, 2, 3}, time.Hour)
	list, exists := GetList("list:key")
	if !exists || len(list) != 3 {
		t.Error("List operation failed")
	}

	// 哈希
	SetHash("hash:key", map[string]interface{}{"a": 1, "b": 2}, time.Hour)
	hash, exists := GetHash("hash:key")
	if !exists || len(hash) != 2 {
		t.Error("Hash operation failed")
	}

	// 验证缓存大小
	if Size() != 3 {
		t.Errorf("Expected cache size 3, got %d", Size())
	}
}

func TestGlobalCacheTTL(t *testing.T) {
	// 清空全局缓存
	Flush()

	// 设置带 TTL 的数据
	SetString("ttl:key", "value", time.Second)

	// 检查 TTL
	ttl, exists := TTL("ttl:key")
	if !exists {
		t.Error("Key should exist")
	}
	if ttl <= 0 || ttl > time.Second {
		t.Errorf("Unexpected TTL: %v", ttl)
	}

	// 修改 TTL
	Expire("ttl:key", 2*time.Minute)

	// 检查修改后的 TTL
	ttl, exists = TTL("ttl:key")
	if !exists {
		t.Error("Key should still exist")
	}
	if ttl < 1*time.Minute || ttl > 3*time.Minute {
		t.Errorf("Unexpected modified TTL: %v", ttl)
	}
}
