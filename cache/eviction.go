package cache

import (
	"fmt"

	"github.com/scache/constants"
	"github.com/scache/interfaces"
	"github.com/scache/policies/fifo"
	"github.com/scache/policies/lfu"
	"github.com/scache/policies/lru"
)

// newEvictionPolicy 创建淘汰策略
func newEvictionPolicy(policyType string, maxSize int) interfaces.EvictionPolicy {
	switch policyType {
	case constants.LRUStrategy:
		return lru.NewLRUPolicy(maxSize)
	case constants.LFUStrategy:
		return lfu.NewLFUPolicy(maxSize)
	case constants.FIFOStrategy:
		return fifo.NewFIFOPolicy(maxSize)
	default:
		return lru.NewLRUPolicy(maxSize)
	}
}

// 重新导出常量
var (
	// LRU 淘汰策略
	LRU = constants.LRUStrategy
	// LFU 淘汰策略
	LFU = constants.LFUStrategy
	// FIFO 淘汰策略
	FIFO = constants.FIFOStrategy
)

// EvictionPolicyType 淘汰策略类型（兼容性）
type EvictionPolicyType string

// String 返回策略类型的字符串表示
func (e EvictionPolicyType) String() string {
	return string(e)
}

// IsValid 检查策略类型是否有效
func (e EvictionPolicyType) IsValid() bool {
	switch string(e) {
	case constants.LRUStrategy, constants.LFUStrategy, constants.FIFOStrategy:
		return true
	default:
		return false
	}
}

// ParseEvictionPolicy 解析淘汰策略类型
func ParseEvictionPolicy(s string) (EvictionPolicyType, error) {
	policy := EvictionPolicyType(s)
	if !policy.IsValid() {
		return "", fmt.Errorf(constants.ErrInvalidStrategy+": %s", s)
	}
	return policy, nil
}
