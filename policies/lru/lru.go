package lru

import (
	"container/list"
	"sync"

	"github.com/scache/interfaces"
)

// LRUPolicy LRU (Least Recently Used) 淘汰策略
type LRUPolicy struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	lock     sync.RWMutex
}

// entry LRU 缓存条目
type entry struct {
	key   string
	value interface{}
}

// NewLRUPolicy 创建新的 LRU 策略
func NewLRUPolicy(capacity int) *LRUPolicy {
	return &LRUPolicy{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// OnAccess 当缓存项被访问时调用
func (l *LRUPolicy) OnAccess(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if elem, exists := l.cache[key]; exists {
		l.list.MoveToFront(elem)
	}
}

// OnAdd 当新缓存项被添加时调用
func (l *LRUPolicy) OnAdd(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// 如果键已存在，更新它
	if elem, exists := l.cache[key]; exists {
		l.list.MoveToFront(elem)
		return
	}

	// 添加新元素到前面
	elem := l.list.PushFront(&entry{key: key})
	l.cache[key] = elem

	// 如果超出容量，移除最旧的元素
	if l.list.Len() > l.capacity {
		l.removeOldest()
	}
}

// OnRemove 当缓存项被移除时调用
func (l *LRUPolicy) OnRemove(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if elem, exists := l.cache[key]; exists {
		l.list.Remove(elem)
		delete(l.cache, key)
	}
}

// ShouldEvict 判断是否需要淘汰缓存项，返回要淘汰的键
func (l *LRUPolicy) ShouldEvict() (string, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if l.list.Len() >= l.capacity {
		if elem := l.list.Back(); elem != nil {
			entry := elem.Value.(*entry)
			return entry.key, true
		}
	}
	return "", false
}

// SetMaxSize 设置最大容量
func (l *LRUPolicy) SetMaxSize(size int) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.capacity = size

	// 如果当前大小超过新容量，移除多余的元素
	for l.list.Len() > l.capacity {
		l.removeOldest()
	}
}

// removeOldest 移除最旧的元素（内部方法，调用时需要持有锁）
func (l *LRUPolicy) removeOldest() {
	if elem := l.list.Back(); elem != nil {
		entry := elem.Value.(*entry)
		l.list.Remove(elem)
		delete(l.cache, entry.key)
	}
}

// Len 返回当前缓存项数量
func (l *LRUPolicy) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.list.Len()
}

// Keys 返回所有键，从新到旧排序
func (l *LRUPolicy) Keys() []string {
	l.lock.RLock()
	defer l.lock.RUnlock()

	keys := make([]string, 0, l.list.Len())
	for elem := l.list.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*entry)
		keys = append(keys, entry.key)
	}
	return keys
}

// Contains 检查是否包含指定键
func (l *LRUPolicy) Contains(key string) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	_, exists := l.cache[key]
	return exists
}

// Clear 清空所有缓存项
func (l *LRUPolicy) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.cache = make(map[string]*list.Element)
	l.list = list.New()
}

// 确保 LRUPolicy 实现了 EvictionPolicy 接口
var _ interfaces.EvictionPolicy = (*LRUPolicy)(nil)
