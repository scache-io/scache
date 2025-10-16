package commands

import (
	"strings"
	"sync"

	"scache/interfaces"
)

// CommandRegistry 命令注册表
type CommandRegistry struct {
	commands map[string]interfaces.Command
	mu       sync.RWMutex
}

// NewCommandRegistry 创建新的命令注册表
func NewCommandRegistry() *CommandRegistry {
	registry := &CommandRegistry{
		commands: make(map[string]interfaces.Command),
	}

	// 注册默认命令
	registry.RegisterDefaults()
	return registry
}

// Register 注册命令
func (r *CommandRegistry) Register(command interfaces.Command) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := strings.ToUpper(command.Name())
	r.commands[name] = command
}

// Get 获取命令
func (r *CommandRegistry) Get(name string) (interfaces.Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	name = strings.ToUpper(name)
	cmd, exists := r.commands[name]
	return cmd, exists
}

// List 列出所有命令
func (r *CommandRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}

// RegisterDefaults 注册默认命令
func (r *CommandRegistry) RegisterDefaults() {
	// 字符串命令
	r.Register(NewSetCommand())
	r.Register(NewGetCommand())

	// 列表命令
	r.Register(NewLPushCommand())
	r.Register(NewRPopCommand())

	// 哈希命令
	r.Register(NewHSetCommand())
	r.Register(NewHGetCommand())

	// 通用命令
	r.Register(NewDeleteCommand())
	r.Register(NewExistsCommand())
	r.Register(NewTypeCommand())
	r.Register(NewExpireCommand())
	r.Register(NewTTLCommand())
	r.Register(NewStatsCommand())
}

// 全局命令注册表
var (
	defaultRegistry *CommandRegistry
	registryOnce    sync.Once
)

// DefaultRegistry 获取默认命令注册表
func DefaultRegistry() *CommandRegistry {
	registryOnce.Do(func() {
		defaultRegistry = NewCommandRegistry()
	})
	return defaultRegistry
}

// Register 注册命令到默认注册表
func Register(command interfaces.Command) {
	DefaultRegistry().Register(command)
}

// Get 从默认注册表获取命令
func Get(name string) (interfaces.Command, bool) {
	return DefaultRegistry().Get(name)
}

// List 列出默认注册表的所有命令
func List() []string {
	return DefaultRegistry().List()
}
