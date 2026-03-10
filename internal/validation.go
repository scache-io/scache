package internal

import (
	"fmt"
	"reflect"
)

// ValidateCacheKey 验证Cache key是否有效
func ValidateCacheKey(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	return nil
}

// ValidatePointerArgument 验证Parameter是否为指针Type
func ValidatePointerArgument(dest interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("invalid argument: got %T", dest)
	}
	return nil
}

// ValidateCapacity 验证容量Parameter是否有效
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

// ValidateStructName 验证Struct name是否有效
func ValidateStructName(name string) error {
	if name == "" {
		return fmt.Errorf("invalid argument: struct name cannot be empty")
	}
	return nil
}
