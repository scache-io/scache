package fifo

import (
	"sync"

	"github.com/scache/interfaces"
)

// FIFOPolicy FIFO (First In First Out) 淘汰策略
type FIFOPolicy struct {
	capacity int
	queue    []string
	items    map[string]int
	lock     sync.RWMutex
}

// NewFIFOPolicy 创建新的 FIFO 策略
func NewFIFOPolicy(capacity int) *FIFOPolicy {
	return &FIFOPolicy{
		capacity: capacity,
		queue:    make([]string, 0, capacity),
		items:    make(map[string]int),
	}
}

// OnAccess 当缓存项被访问时调用（FIFO 不需要处理访问事件）
func (f *FIFOPolicy) OnAccess(key string) {
	// FIFO 策略不关心访问次数，只需要关心插入顺序
}

// OnAdd 当新缓存项被添加时调用
func (f *FIFOPolicy) OnAdd(key string) {
	f.lock.Lock()
	defer f.lock.Unlock()

	// 如果键已存在，不需要做任何操作
	if _, exists := f.items[key]; exists {
		return
	}

	// 添加到队列末尾
	f.queue = append(f.queue, key)
	f.items[key] = len(f.queue) - 1

	// 如果超出容量，移除最旧的元素
	if len(f.queue) > f.capacity {
		f.removeOldest()
	}
}

// OnRemove 当缓存项被移除时调用
func (f *FIFOPolicy) OnRemove(key string) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if index, exists := f.items[key]; exists {
		// 从队列中移除
		f.queue = append(f.queue[:index], f.queue[index+1:]...)

		// 更新后面元素的索引
		for i := index; i < len(f.queue); i++ {
			itemKey := f.queue[i]
			f.items[itemKey] = i
		}

		delete(f.items, key)
	}
}

// ShouldEvict 判断是否需要淘汰缓存项
func (f *FIFOPolicy) ShouldEvict() (string, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if len(f.queue) >= f.capacity && len(f.queue) > 0 {
		return f.queue[0], true
	}
	return "", false
}

// SetMaxSize 设置最大容量
func (f *FIFOPolicy) SetMaxSize(size int) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.capacity = size

	// 如果当前大小超过新容量，移除多余的元素
	for len(f.queue) > f.capacity {
		f.removeOldest()
	}
}

// removeOldest 移除最旧的元素（内部方法，调用时需要持有锁）
func (f *FIFOPolicy) removeOldest() {
	if len(f.queue) > 0 {
		oldestKey := f.queue[0]
		f.queue = f.queue[1:]
		delete(f.items, oldestKey)
	}
}

// Len 返回当前缓存项数量
func (f *FIFOPolicy) Len() int {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return len(f.queue)
}

// Keys 返回所有键，按插入顺序排序
func (f *FIFOPolicy) Keys() []string {
	f.lock.RLock()
	defer f.lock.RUnlock()

	keys := make([]string, len(f.queue))
	copy(keys, f.queue)
	return keys
}

// Contains 检查是否包含指定键
func (f *FIFOPolicy) Contains(key string) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	_, exists := f.items[key]
	return exists
}

// Clear 清空所有缓存项
func (f *FIFOPolicy) Clear() {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.queue = make([]string, 0, f.capacity)
	f.items = make(map[string]int)
}

// GetPosition 获取指定键在队列中的位置
func (f *FIFOPolicy) GetPosition(key string) int {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if index, exists := f.items[key]; exists {
		return index
	}
	return -1
}

// 确保 FIFOPolicy 实现了 EvictionPolicy 接口
var _ interfaces.EvictionPolicy = (*FIFOPolicy)(nil)
