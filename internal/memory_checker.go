package internal

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

var (
	// 内存监控状态
	memoryCheckEnabled int32 = 1 // 默认启用内存检查
)

// MemoryInfo 内存信息
type MemoryInfo struct {
	Alloc      uint64 // 已分配的堆内存 (字节)
	TotalAlloc uint64 // 累计分配的内存 (字节)
	Sys        uint64 // 从系统获取的内存 (字节)
	NumGC      uint32 // GC运行次数
}

// GetMemoryInfo 获取当前内存信息
func GetMemoryInfo() MemoryInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryInfo{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}
}

// IsMemoryThresholdExceeded 检查内存使用是否超过阈值
func IsMemoryThresholdExceeded(threshold float64) bool {
	if threshold <= 0 || threshold > 1 {
		return false
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 计算内存使用率（使用已分配内存 vs 系统内存）
	// 注意：这里使用堆内存作为参考，实际使用中可以根据需要调整
	memoryUsage := float64(m.Alloc) / float64(m.Sys)

	return memoryUsage > threshold
}

// CheckMemoryAvailability 检查内存可用性
func CheckMemoryAvailability(threshold float64) error {
	if !IsMemoryCheckEnabled() {
		return nil
	}

	if threshold <= 0 || threshold > 1 {
		return nil // 无效阈值，跳过检查
	}

	if IsMemoryThresholdExceeded(threshold) {
		memInfo := GetMemoryInfo()
		return fmt.Errorf("insufficient memory: current usage %d bytes exceeds threshold %.2f",
			memInfo.Alloc, threshold)
	}

	return nil
}

// EnableMemoryCheck 启用内存检查
func EnableMemoryCheck() {
	atomic.StoreInt32(&memoryCheckEnabled, 1)
}

// DisableMemoryCheck 禁用内存检查
func DisableMemoryCheck() {
	atomic.StoreInt32(&memoryCheckEnabled, 0)
}

// IsMemoryCheckEnabled 检查内存检查是否启用
func IsMemoryCheckEnabled() bool {
	return atomic.LoadInt32(&memoryCheckEnabled) == 1
}

// GetMemoryUsagePercentage 获取内存使用百分比
func GetMemoryUsagePercentage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if m.Sys == 0 {
		return 0
	}

	return float64(m.Alloc) / float64(m.Sys) * 100
}