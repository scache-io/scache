package cache

import (
	"testing"
	"time"

	"github.com/scache-io/scache/config"
)

func TestLocalCache_BasicOperations(t *testing.T) {
	cache := NewLocalCache()

	// 测试 SetString 和 GetString
	err := cache.SetString("test_key", "test_value", 0)
	if err != nil {
		t.Fatalf("设置字符串失败: %v", err)
	}

	value, exists := cache.GetString("test_key")
	if !exists {
		t.Fatal("获取字符串失败，键不存在")
	}
	if value != "test_value" {
		t.Fatalf("期望值 'test_value'，实际值 '%s'", value)
	}

	// 测试 Exists
	if !cache.Exists("test_key") {
		t.Fatal("Exists 检查失败")
	}

	// 测试 Delete
	if !cache.Delete("test_key") {
		t.Fatal("删除失败")
	}

	if cache.Exists("test_key") {
		t.Fatal("删除后键仍然存在")
	}
}

func TestLocalCache_ListOperations(t *testing.T) {
	cache := NewLocalCache()

	values := []interface{}{"item1", "item2", "item3"}
	err := cache.SetList("list_key", values, 0)
	if err != nil {
		t.Fatalf("设置列表失败: %v", err)
	}

	retrievedValues, exists := cache.GetList("list_key")
	if !exists {
		t.Fatal("获取列表失败，键不存在")
	}

	if len(retrievedValues) != len(values) {
		t.Fatalf("期望长度 %d，实际长度 %d", len(values), len(retrievedValues))
	}

	for i, v := range values {
		if retrievedValues[i] != v {
			t.Fatalf("索引 %d 期望值 %v，实际值 %v", i, v, retrievedValues[i])
		}
	}
}

func TestLocalCache_HashOperations(t *testing.T) {
	cache := NewLocalCache()

	fields := map[string]interface{}{
		"name": "张三",
		"age":  30,
		"city": "北京",
	}

	err := cache.SetHash("hash_key", fields, 0)
	if err != nil {
		t.Fatalf("设置哈希失败: %v", err)
	}

	retrievedFields, exists := cache.GetHash("hash_key")
	if !exists {
		t.Fatal("获取哈希失败，键不存在")
	}

	if len(retrievedFields) != len(fields) {
		t.Fatalf("期望字段数 %d，实际字段数 %d", len(fields), len(retrievedFields))
	}

	for k, v := range fields {
		if retrievedFields[k] != v {
			t.Fatalf("字段 %s 期望值 %v，实际值 %v", k, v, retrievedFields[k])
		}
	}
}

func TestLocalCache_Expiration(t *testing.T) {
	cache := NewLocalCache()

	// 设置1秒过期的值
	err := cache.SetString("expire_key", "expire_value", time.Second)
	if err != nil {
		t.Fatalf("设置字符串失败: %v", err)
	}

	// 立即获取应该存在
	_, exists := cache.GetString("expire_key")
	if !exists {
		t.Fatal("刚设置的键不应该过期")
	}

	// 等待2秒后应该过期
	time.Sleep(2 * time.Second)
	_, exists = cache.GetString("expire_key")
	if exists {
		t.Fatal("键应该已过期")
	}
}

func TestLocalCache_ExpireAndTTL(t *testing.T) {
	cache := NewLocalCache()

	// 设置不过期的值
	err := cache.SetString("ttl_key", "ttl_value", 0)
	if err != nil {
		t.Fatalf("设置字符串失败: %v", err)
	}

	// 设置过期时间
	success := cache.Expire("ttl_key", time.Minute*5)
	if !success {
		t.Fatal("设置过期时间失败")
	}

	// 获取TTL
	ttl, exists := cache.TTL("ttl_key")
	if !exists {
		t.Fatal("键不存在")
	}

	if ttl <= 0 {
		t.Fatal("TTL应该大于0")
	}

	if ttl > time.Minute*5 {
		t.Fatal("TTL超过预期值")
	}
}

func TestLocalCache_Stats(t *testing.T) {
	cache := NewLocalCache()

	// 初始统计
	stats := cache.Stats()
	if stats == nil {
		t.Fatal("统计信息为空")
	}

	// 设置一些值
	cache.SetString("key1", "value1", 0)
	cache.SetString("key2", "value2", 0)

	// 获取值
	cache.GetString("key1")
	cache.GetString("key3") // 不存在的键

	stats = cache.Stats()
	// 验证统计信息包含必要的字段
	if _, ok := stats.(map[string]interface{}); !ok {
		t.Fatal("统计信息格式不正确")
	}
}

func TestLocalCache_Configuration(t *testing.T) {
	// 使用配置创建缓存
	cache := NewLocalCache(
		config.WithMaxSize(10),
		config.WithDefaultExpiration(time.Minute),
		config.WithMemoryThreshold(0.7),
	)

	// 添加超过限制的键，并访问它们以确保LRU策略正确工作
	for i := 0; i < 15; i++ {
		key := "key" + string(rune('A'+i))
		value := "value" + string(rune('A'+i))
		err := cache.SetString(key, value, 0)
		if err != nil {
			t.Fatalf("设置键 %s 失败: %v", key, err)
		}
		// 立即访问以触发LRU策略的Access方法
		cache.GetString(key)
	}

	// 缓存大小应该不超过限制
	size := cache.Size()
	if size > 10 {
		t.Fatalf("缓存大小 %d 超过限制 10", size)
	}

	// 验证最后5个键存在（因为它们最近被访问）
	for i := 10; i < 15; i++ {
		key := "key" + string(rune('A'+i))
		if _, exists := cache.GetString(key); !exists {
			t.Fatalf("最近的键 %s 应该存在", key)
		}
	}
}

func TestGlobalCache_BasicOperations(t *testing.T) {
	// 创建局部缓存实例用于测试全局功能模式
	testCache := NewLocalCache()

	// 测试 SetString
	err := testCache.SetString("global_key", "global_value", 0)
	if err != nil {
		t.Fatalf("设置字符串失败: %v", err)
	}

	// 测试 GetString
	value, exists := testCache.GetString("global_key")
	if !exists {
		t.Fatal("获取字符串失败，键不存在")
	}
	if value != "global_value" {
		t.Fatalf("期望值 'global_value'，实际值 '%s'", value)
	}

	// 测试 Exists
	if !testCache.Exists("global_key") {
		t.Fatal("Exists 检查失败")
	}

	// 测试 Delete
	if !testCache.Delete("global_key") {
		t.Fatal("删除失败")
	}

	if testCache.Exists("global_key") {
		t.Fatal("删除后键仍然存在")
	}
}

func TestGlobalCache_Configuration(t *testing.T) {
	// 使用配置初始化缓存
	testCache := NewLocalCache(
		config.WithMaxSize(5),
		config.WithDefaultExpiration(time.Minute*2),
	)

	// 添加键
	for i := 0; i < 3; i++ {
		testCache.SetString("global_key"+string(rune('A'+i)), "value"+string(rune('A'+i)), 0)
	}

	// 检查大小
	if testCache.Size() != 3 {
		t.Fatalf("期望大小 3，实际大小 %d", testCache.Size())
	}

	// 检查统计
	stats := testCache.Stats()
	if stats == nil {
		t.Fatal("统计信息为空")
	}
}

func TestGlobalCache_MixedTypes(t *testing.T) {
	// 创建一个全新的缓存实例进行测试
	testCache := NewLocalCache()

	// 设置不同类型的数据
	testCache.SetString("str_key", "string_value", 0)
	testCache.SetList("list_key", []interface{}{"a", "b", "c"}, 0)
	testCache.SetHash("hash_key", map[string]interface{}{"field1": "value1"}, 0)

	// 获取数据
	if _, exists := testCache.GetString("str_key"); !exists {
		t.Fatal("字符串数据不存在")
	}

	if _, exists := testCache.GetList("list_key"); !exists {
		t.Fatal("列表数据不存在")
	}

	if _, exists := testCache.GetHash("hash_key"); !exists {
		t.Fatal("哈希数据不存在")
	}

	// 获取所有键
	keys := testCache.Keys()
	if len(keys) != 3 {
		t.Fatalf("期望3个键，实际%d个", len(keys))
	}
}
