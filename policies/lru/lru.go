package lru

import (
	"container/list"
	"sync"

	"github.com/scache-io/scache/constants"
	"github.com/scache-io/scache/interfaces"
)

// 本包实现了LRU（Least Recently Used）缓存淘汰策略

// lruPolicy LRU淘汰策略的实现结构体
type lruPolicy struct {
	capacity int                      // 缓存容量
	cache    map[string]*list.Element // 键到链表节点的映射，用于O(1)查找
	list     *list.List               // 双向链表，头部为最近使用，尾部为最久未使用
	mu       sync.RWMutex             // 读写锁，保护并发访问
}

// lruNode 链表中存储的节点数据
type lruNode struct {
	key string // 缓存键
}

// NewLRUPolicy 创建一个新的LRU淘汰策略实例
// capacity: 缓存容量，如果小于等于0则使用默认值
func NewLRUPolicy(capacity int) interfaces.EvictionPolicy {
	if capacity <= 0 {
		capacity = constants.DefaultLRUCapacity // 使用默认容量
	}

	return &lruPolicy{
		capacity: capacity,
		cache:    make(map[string]*list.Element), // 初始化映射表
		list:     list.New(),                     // 初始化双向链表
	}
}

// Access 访问指定键，将其标记为最近使用
// 如果键不存在，则添加到缓存；如果超过容量，则淘汰最久未使用的条目
func (l *lruPolicy) Access(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, exists := l.cache[key]; exists {
		l.list.MoveToFront(elem) // 移动到链表头部，标记为最近使用
	} else {
		elem := l.list.PushFront(&lruNode{key: key}) // 添加新节点到头部
		l.cache[key] = elem                          // 建立映射关系

		if l.list.Len() > l.capacity {
			l.evictInternal() // 超过容量时淘汰最久未使用的条目
		}
	}
}

// Set 设置指定键的值，等同于Access操作
func (l *lruPolicy) Set(key string) {
	l.Access(key)
}

// Delete 从缓存中删除指定键的条目
func (l *lruPolicy) Delete(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, exists := l.cache[key]; exists {
		l.list.Remove(elem)  // 从链表中移除节点
		delete(l.cache, key) // 从映射表中删除键
	}
}

// Evict 淘汰最久未使用的缓存条目，返回被淘汰的键
func (l *lruPolicy) Evict() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.evictInternal()
}

// evictInternal 内部淘汰方法，必须在持有锁的情况下调用
func (l *lruPolicy) evictInternal() string {
	if l.list.Len() == 0 {
		return "" // 空缓存，无需淘汰
	}

	elem := l.list.Back() // 获取链表尾部元素（最久未使用）
	if elem != nil {
		node := elem.Value.(*lruNode)
		l.list.Remove(elem)       // 从链表中移除
		delete(l.cache, node.key) // 从映射表中删除
		return node.key           // 返回被淘汰的键
	}

	return ""
}

// Size 返回当前缓存中的条目数量
func (l *lruPolicy) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.list.Len()
}

// Keys 返回缓存中所有键的列表，按最近使用顺序排列
func (l *lruPolicy) Keys() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	keys := make([]string, 0, l.list.Len()) // 预分配切片容量
	for elem := l.list.Front(); elem != nil; elem = elem.Next() {
		node := elem.Value.(*lruNode)
		keys = append(keys, node.key) // 从头部开始遍历，添加到结果中
	}

	return keys
}

// Contains 检查指定键是否存在于缓存中
func (l *lruPolicy) Contains(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	_, exists := l.cache[key]
	return exists
}

// UpdateCapacity 更新缓存容量，如果新容量小于当前条目数，则淘汰多余的条目
func (l *lruPolicy) UpdateCapacity(newCapacity int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if newCapacity <= 0 {
		return // 无效容量，忽略更新
	}

	l.capacity = newCapacity

	// 如果当前条目数超过新容量，持续淘汰直到符合容量限制
	for l.list.Len() > l.capacity {
		l.evictInternal()
	}
}

// Clear 清空缓存中的所有条目
func (l *lruPolicy) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = make(map[string]*list.Element) // 重新创建映射表
	l.list.Init()                            // 重置链表
}
