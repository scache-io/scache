package interfaces

import "time"

// DataType 数据类型枚举
type DataType string

const (
	DataTypeString DataType = "string"
	DataTypeList   DataType = "list"
	DataTypeHash   DataType = "hash"
)

// DataObject 通用数据对象接口
type DataObject interface {
	Type() DataType
	ExpiresAt() time.Time
	IsExpired() bool
	Size() int
}

// StringObject 字符串对象接口
type StringObject interface {
	DataObject
	Value() string
	Set(value string)
}

// ListObject 列表对象接口
type ListObject interface {
	DataObject
	Values() []interface{}
	Push(value interface{})
	Pop() (interface{}, bool)
	Index(index int) (interface{}, bool)
	Range(start, end int) []interface{}
	Len() int
}

// HashObject 哈希对象接口
type HashObject interface {
	DataObject
	Fields() map[string]interface{}
	Get(field string) (interface{}, bool)
	Set(field string, value interface{})
	Delete(field string) bool
	Len() int
}

// StorageEngine 存储引擎接口
type StorageEngine interface {
	// 基础操作
	Set(key string, obj DataObject) error
	Get(key string) (DataObject, bool)
	Delete(key string) bool
	Exists(key string) bool
	Keys() []string
	Flush() error
	Size() int

	// 类型检查
	Type(key string) (DataType, bool)

	// 过期管理
	Expire(key string, ttl time.Duration) bool
	TTL(key string) (time.Duration, bool)

	// 统计信息
	Stats() interface{}
}


// EvictionPolicy 淘汰策略接口
type EvictionPolicy interface {
	// Access 当访问 key 时调用
	Access(key string)

	// Set 当设置新 key 时调用
	Set(key string)

	// Delete 当删除 key 时调用
	Delete(key string)

	// Evict 获取需要淘汰的 key
	Evict() string

	// Size 获取当前策略状态
	Size() int

	// Clear 清空所有数据
	Clear()

	// Contains 检查 key 是否存在
	Contains(key string) bool

	// Keys 获取所有 key（按最近使用时间排序）
	Keys() []string

	// UpdateCapacity 更新容量限制
	UpdateCapacity(newCapacity int)
}

