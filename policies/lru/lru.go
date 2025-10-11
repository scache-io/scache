package lru

import (
	"container/list"
	"sync"

	"scache/constants"
	"scache/interfaces"
)

// lruPolicy LRU 淘汰策略实现
type lruPolicy struct {
	capacity int                      // 容量限制
	cache    map[string]*list.Element // 键到链表节点的映射
	list     *list.List               // 双向链表
	mu       sync.RWMutex             // 读写锁
}

// lruNode LRU 链表节点
type lruNode struct {
	key string
}

// NewLRUPolicy 创建 LRU 策略
func NewLRUPolicy(capacity int) interfaces.EvictionPolicy {
	if capacity <= constants.DefaultExpiration {
		capacity = constants.DefaultLRUCapacity // 默认容量
	}

	return &lruPolicy{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Access 当访问 key 时调用
func (l *lruPolicy) Access(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, exists := l.cache[key]; exists {
		// 如果 key 存在，移动到链表头部
		l.list.MoveToFront(elem)
	} else {
		// 如果 key 不存在，添加到链表头部
		elem := l.list.PushFront(&lruNode{key: key})
		l.cache[key] = elem

		// 检查容量限制
		if l.list.Len() > l.capacity {
			l.evictInternal()
		}
	}
}

// Set 当设置新 key 时调用
func (l *lruPolicy) Set(key string) {
	l.Access(key) // LRU 策略中，Set 和 Access 处理相同
}

// Delete 当删除 key 时调用
func (l *lruPolicy) Delete(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, exists := l.cache[key]; exists {
		l.list.Remove(elem)
		delete(l.cache, key)
	}
}

// Evict 获取需要淘汰的 key
func (l *lruPolicy) Evict() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.evictInternal()
}

// evictInternal 内部淘汰方法（调用时需要持有锁）
func (l *lruPolicy) evictInternal() string {
	if l.list.Len() == constants.DefaultExpiration {
		return ""
	}

	// 获取链表尾部节点（最少使用的）
	elem := l.list.Back()
	if elem != nil {
		node := elem.Value.(*lruNode)
		l.list.Remove(elem)
		delete(l.cache, node.key)
		return node.key
	}

	return ""
}

// Size 获取当前策略状态
func (l *lruPolicy) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.list.Len()
}

// Keys 获取所有 key（按最近使用时间排序，最新的在前）
func (l *lruPolicy) Keys() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	keys := make([]string, 0, l.list.Len())
	for elem := l.list.Front(); elem != nil; elem = elem.Next() {
		node := elem.Value.(*lruNode)
		keys = append(keys, node.key)
	}

	return keys
}

// Contains 检查 key 是否存在
func (l *lruPolicy) Contains(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	_, exists := l.cache[key]
	return exists
}

// UpdateCapacity 更新容量限制
func (l *lruPolicy) UpdateCapacity(newCapacity int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if newCapacity <= constants.DefaultExpiration {
		return
	}

	l.capacity = newCapacity

	// 如果当前数量超过新容量，淘汰多余的项
	for l.list.Len() > l.capacity {
		l.evictInternal()
	}
}

// Clear 清空所有数据
func (l *lruPolicy) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = make(map[string]*list.Element)
	l.list.Init()
}
