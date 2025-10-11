package cache

import (
	"fmt"
	"sync"
)

// 全局缓存管理器
type Manager struct {
	caches map[string]Cache
	mutex  sync.RWMutex
}

var (
	globalManager *Manager
	once          sync.Once
)

// GetGlobalManager 获取全局缓存管理器实例（单例模式）
func GetGlobalManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			caches: make(map[string]Cache),
		}
	})
	return globalManager
}

// Register 注册一个命名缓存
func (m *Manager) Register(name string, c Cache) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.caches[name]; exists {
		return fmt.Errorf("cache '%s' already registered", name)
	}

	m.caches[name] = c
	return nil
}

// Get 获取已注册的缓存
func (m *Manager) Get(name string) (Cache, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	c, exists := m.caches[name]
	if !exists {
		return nil, fmt.Errorf("cache '%s' not found", name)
	}
	return c, nil
}

// GetOrDefault 获取缓存，如果不存在则使用默认配置创建
func (m *Manager) GetOrDefault(name string, opts ...Option) Cache {
	if c, err := m.Get(name); err == nil {
		return c
	}

	c := New(opts...)
	m.Register(name, c)
	return c
}

// Remove 移除已注册的缓存
func (m *Manager) Remove(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if c, exists := m.caches[name]; exists {
		c.Close()
		delete(m.caches, name)
		return nil
	}
	return fmt.Errorf("cache '%s' not found", name)
}

// List 列出所有已注册的缓存名称
func (m *Manager) List() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.caches))
	for name := range m.caches {
		names = append(names, name)
	}
	return names
}

// Clear 清空所有缓存
func (m *Manager) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for name, c := range m.caches {
		if err := c.Clear(); err != nil {
			return fmt.Errorf("failed to clear cache '%s': %w", name, err)
		}
	}
	return nil
}

// Close 关闭所有缓存并清理管理器
func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var errors []error
	for name, c := range m.caches {
		if err := c.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close cache '%s': %w", name, err))
		}
	}

	m.caches = make(map[string]Cache)

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while closing caches: %v", errors)
	}
	return nil
}

// Stats 获取所有缓存的统计信息
func (m *Manager) Stats() map[string]CacheStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string]CacheStats)
	for name, c := range m.caches {
		stats[name] = c.Stats()
	}
	return stats
}

// Size 获取所有缓存的总大小
func (m *Manager) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := 0
	for _, c := range m.caches {
		total += c.Size()
	}
	return total
}

// Exists 检查缓存是否已注册
func (m *Manager) Exists(name string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.caches[name]
	return exists
}
