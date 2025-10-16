package types

import (
	"sync"
	"time"

	"github.com/scache-io/scache/interfaces"
)

// BaseObject 基础对象实现
type BaseObject struct {
	dataType  interfaces.DataType
	expiresAt time.Time
	created   time.Time
	accessed  time.Time
	mu        sync.RWMutex
}

// NewBaseObject 创建基础对象
func NewBaseObject(dataType interfaces.DataType, ttl time.Duration) *BaseObject {
	now := time.Now()
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = now.Add(ttl)
	}

	return &BaseObject{
		dataType:  dataType,
		expiresAt: expiresAt,
		created:   now,
		accessed:  now,
	}
}

// Type 返回数据类型
func (o *BaseObject) Type() interfaces.DataType {
	return o.dataType
}

// ExpiresAt 返回过期时间
func (o *BaseObject) ExpiresAt() time.Time {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.expiresAt
}

// IsExpired 检查是否过期
func (o *BaseObject) IsExpired() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(o.expiresAt)
}

// UpdateAccess 更新访问时间
func (o *BaseObject) UpdateAccess() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.accessed = time.Now()
}

// CreatedAt 返回创建时间
func (o *BaseObject) CreatedAt() time.Time {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.created
}

// StringObject 字符串对象实现
type StringObject struct {
	*BaseObject
	value string
	mu    sync.RWMutex
}

// NewStringObject 创建字符串对象
func NewStringObject(value string, ttl time.Duration) *StringObject {
	return &StringObject{
		BaseObject: NewBaseObject(interfaces.DataTypeString, ttl),
		value:      value,
	}
}

// Value 返回字符串值
func (s *StringObject) Value() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.UpdateAccess()
	return s.value
}

// Set 设置字符串值
func (s *StringObject) Set(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = value
	s.UpdateAccess()
}

// Size 返回对象大小（字节）
func (s *StringObject) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.value)
}

// ListObject 列表对象实现
type ListObject struct {
	*BaseObject
	values []interface{}
	mu     sync.RWMutex
}

// NewListObject 创建列表对象
func NewListObject(values []interface{}, ttl time.Duration) *ListObject {
	return &ListObject{
		BaseObject: NewBaseObject(interfaces.DataTypeList, ttl),
		values:     values,
	}
}

// Values 返回所有值
func (l *ListObject) Values() []interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.UpdateAccess()

	// 返回副本避免外部修改
	result := make([]interface{}, len(l.values))
	copy(result, l.values)
	return result
}

// Push 在列表末尾添加元素
func (l *ListObject) Push(value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.values = append(l.values, value)
	l.UpdateAccess()
}

// Pop 从列表末尾移除元素
func (l *ListObject) Pop() (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.values) == 0 {
		return nil, false
	}

	index := len(l.values) - 1
	value := l.values[index]
	l.values = l.values[:index]
	l.UpdateAccess()
	return value, true
}

// Index 返回指定索引的元素
func (l *ListObject) Index(index int) (interface{}, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if index < 0 || index >= len(l.values) {
		return nil, false
	}

	l.UpdateAccess()
	return l.values[index], true
}

// Range 返回指定范围的元素
func (l *ListObject) Range(start, end int) []interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if start < 0 {
		start = 0
	}
	if end >= len(l.values) {
		end = len(l.values) - 1
	}
	if start > end {
		return nil
	}

	l.UpdateAccess()

	result := make([]interface{}, end-start+1)
	copy(result, l.values[start:end+1])
	return result
}

// Len 返回列表长度
func (l *ListObject) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.UpdateAccess()
	return len(l.values)
}

// Size 返回对象大小
func (l *ListObject) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.values) * 8 // 估算每个元素8字节
}

// HashObject 哈希对象实现
type HashObject struct {
	*BaseObject
	fields map[string]interface{}
	mu     sync.RWMutex
}

// NewHashObject 创建哈希对象
func NewHashObject(fields map[string]interface{}, ttl time.Duration) *HashObject {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	return &HashObject{
		BaseObject: NewBaseObject(interfaces.DataTypeHash, ttl),
		fields:     fields,
	}
}

// Fields 返回所有字段
func (h *HashObject) Fields() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	h.UpdateAccess()

	// 返回副本避免外部修改
	result := make(map[string]interface{}, len(h.fields))
	for k, v := range h.fields {
		result[k] = v
	}
	return result
}

// Get 获取字段值
func (h *HashObject) Get(field string) (interface{}, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	value, exists := h.fields[field]
	if exists {
		h.UpdateAccess()
	}
	return value, exists
}

// Set 设置字段值
func (h *HashObject) Set(field string, value interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[field] = value
	h.UpdateAccess()
}

// Delete 删除字段
func (h *HashObject) Delete(field string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.fields[field]; exists {
		delete(h.fields, field)
		h.UpdateAccess()
		return true
	}
	return false
}

// Len 返回字段数量
func (h *HashObject) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	h.UpdateAccess()
	return len(h.fields)
}

// Size 返回对象大小
func (h *HashObject) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	size := 0
	for k := range h.fields {
		size += len(k) + 8 // 键长度 + 值估算
	}
	return size
}
