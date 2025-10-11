package lfu

import (
	"container/heap"
	"sync"

	"github.com/scache/interfaces"
)

// LFUPolicy LFU (Least Frequently Used) 淘汰策略
type LFUPolicy struct {
	capacity int
	items    map[string]*lfuItem
	heap     *lfuHeap
	lock     sync.RWMutex
}

// lfuItem LFU 缓存项
type lfuItem struct {
	key       string
	frequency int
	index     int
}

// lfuHeap LFU 堆实现
type lfuHeap []*lfuItem

func (h lfuHeap) Len() int           { return len(h) }
func (h lfuHeap) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h lfuHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *lfuHeap) Push(x interface{}) {
	item := x.(*lfuItem)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h *lfuHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1
	*h = old[0 : n-1]
	return item
}

// NewLFUPolicy 创建新的 LFU 策略
func NewLFUPolicy(capacity int) *LFUPolicy {
	h := &lfuHeap{}
	heap.Init(h)

	return &LFUPolicy{
		capacity: capacity,
		items:    make(map[string]*lfuItem),
		heap:     h,
	}
}

// OnAccess 当缓存项被访问时调用
func (l *LFUPolicy) OnAccess(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if item, exists := l.items[key]; exists {
		item.frequency++
		heap.Fix(l.heap, item.index)
	}
}

// OnAdd 当新缓存项被添加时调用
func (l *LFUPolicy) OnAdd(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if _, exists := l.items[key]; exists {
		return
	}

	item := &lfuItem{
		key:       key,
		frequency: 1,
	}
	l.items[key] = item
	heap.Push(l.heap, item)

	// 如果超出容量，移除频率最低的元素
	if l.heap.Len() > l.capacity {
		l.removeLowestFrequency()
	}
}

// OnRemove 当缓存项被移除时调用
func (l *LFUPolicy) OnRemove(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if item, exists := l.items[key]; exists {
		heap.Remove(l.heap, item.index)
		delete(l.items, key)
	}
}

// ShouldEvict 判断是否需要淘汰缓存项
func (l *LFUPolicy) ShouldEvict() (string, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if l.heap.Len() >= l.capacity {
		if l.heap.Len() > 0 {
			item := (*l.heap)[0]
			return item.key, true
		}
	}
	return "", false
}

// SetMaxSize 设置最大容量
func (l *LFUPolicy) SetMaxSize(size int) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.capacity = size

	// 如果当前大小超过新容量，移除多余的元素
	for l.heap.Len() > l.capacity {
		l.removeLowestFrequency()
	}
}

// removeLowestFrequency 移除频率最低的元素（内部方法，调用时需要持有锁）
func (l *LFUPolicy) removeLowestFrequency() {
	if l.heap.Len() > 0 {
		item := heap.Pop(l.heap).(*lfuItem)
		delete(l.items, item.key)
	}
}

// Len 返回当前缓存项数量
func (l *LFUPolicy) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.heap.Len()
}

// Keys 返回所有键，按频率从低到高排序
func (l *LFUPolicy) Keys() []string {
	l.lock.RLock()
	defer l.lock.RUnlock()

	keys := make([]string, 0, l.heap.Len())
	for _, item := range *l.heap {
		keys = append(keys, item.key)
	}
	return keys
}

// Contains 检查是否包含指定键
func (l *LFUPolicy) Contains(key string) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	_, exists := l.items[key]
	return exists
}

// Clear 清空所有缓存项
func (l *LFUPolicy) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.items = make(map[string]*lfuItem)
	l.heap = &lfuHeap{}
	heap.Init(l.heap)
}

// GetFrequency 获取指定键的访问频率
func (l *LFUPolicy) GetFrequency(key string) int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if item, exists := l.items[key]; exists {
		return item.frequency
	}
	return 0
}

// 确保 LFUPolicy 实现了 EvictionPolicy 接口
var _ interfaces.EvictionPolicy = (*LFUPolicy)(nil)
