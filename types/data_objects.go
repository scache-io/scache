package types

import (
	"sync"
	"time"

	"github.com/scache-io/scache/interfaces"
)

// BaseObject Base object implementation
type BaseObject struct {
	dataType  interfaces.DataType
	expiresAt time.Time
	created   time.Time
	accessed  time.Time
	mu        sync.RWMutex
}

// NewBaseObject Create base object
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

// Type 返回Data type
func (o *BaseObject) Type() interfaces.DataType {
	return o.dataType
}

// ExpiresAt Return expiration time
func (o *BaseObject) ExpiresAt() time.Time {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.expiresAt
}

// IsExpired Check if expired
func (o *BaseObject) IsExpired() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return isExpiredUnsafe(o.expiresAt)
}

// isExpiredUnsafe 内部过期检查Method（不加锁）
func isExpiredUnsafe(expiresAt time.Time) bool {
	if expiresAt.IsZero() {
		return false
	}
	return time.Now().After(expiresAt)
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

// StringObject String object实现
type StringObject struct {
	*BaseObject
	value string
	mu    sync.RWMutex
}

// NewStringObject 创建String object
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

// Set Set string value
func (s *StringObject) Set(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = value
	s.UpdateAccess()
}

// StructObject Struct object实现（复用StringObject，增加JSON支持）
type StructObject struct {
	*StringObject
}

// NewStructObject 创建Struct object
func NewStructObject(data string, ttl time.Duration) *StructObject {
	return &StructObject{
		StringObject: NewStringObject(data, ttl),
	}
}

// Data 返回JSON数据
func (s *StructObject) Data() string {
	return s.Value()
}

// Set 设置JSON数据
func (s *StructObject) Set(data string) {
	s.StringObject.Set(data)
}

// Size Return object size（字节）
func (s *StringObject) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.value)
}

// ListObject List object实现
type ListObject struct {
	*BaseObject
	values []interface{}
	mu     sync.RWMutex
}

// NewListObject 创建List object
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

// Size Return object size
func (l *ListObject) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.values) * 8 // 估算每个元素8字节
}

// HashObject Hash object实现
type HashObject struct {
	*BaseObject
	fields map[string]interface{}
	mu     sync.RWMutex
}

// NewHashObject 创建Hash object
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

// Size Return object size
func (h *HashObject) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	size := 0
	for k := range h.fields {
		size += len(k) + 8 // 键长度 + 值估算
	}
	return size
}
