package internal

import (
	"fmt"
	"reflect"
)

// ValidateCacheKey 验证缓存键是否有效
func ValidateCacheKey(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	return nil
}

// ValidatePointerArgument 验证参数是否为指针类型
func ValidatePointerArgument(dest interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("invalid argument: got %T", dest)
	}
	return nil
}

// ValidateCapacity 验证容量参数是否有效
func ValidateCapacity(capacity int) error {
	if capacity < 0 {
		return fmt.Errorf("invalid argument: capacity must be non-negative")
	}
	return nil
}

// ValidateMemoryThreshold 验证内存阈值是否有效
func ValidateMemoryThreshold(threshold float64) error {
	if threshold < 0 || threshold > 1 {
		return fmt.Errorf("invalid argument: memory threshold must be between 0 and 1")
	}
	return nil
}

// ValidateStructName 验证结构体名称是否有效
func ValidateStructName(name string) error {
	if name == "" {
		return fmt.Errorf("invalid argument: struct name cannot be empty")
	}
	return nil
}