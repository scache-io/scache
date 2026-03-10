package utils

import "time"

// ParseTTL 解析可选的TTLParameter
// 统一处理可选TTLParameter的逻辑，避免代码重复
func ParseTTL(ttl []time.Duration) time.Duration {
	if len(ttl) > 0 {
		return ttl[0]
	}
	return 0
}

// CalculateRemainingTTL 计算剩余生存时间
// 统一TTL计算逻辑
func CalculateRemainingTTL(expiresAt time.Time) (time.Duration, bool) {
	if expiresAt.IsZero() {
		return -1, true // 永不过期
	}

	remaining := time.Until(expiresAt)
	if remaining <= 0 {
		return 0, true // 已过期
	}

	return remaining, true
}
