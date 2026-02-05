package cache

import (
	"star/core"
	"time"
)

// ==================== 缓存接口定义 ====================

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(key string) (*core.TimeSlot, bool)

	// Set 设置缓存
	Set(key string, value *core.TimeSlot, ttl time.Duration)

	// Delete 删除缓存
	Delete(key string)

	// Clear 清空缓存
	Clear()

	// Size 获取缓存大小
	Size() int
}

// ==================== 缓存 Key 生成 ====================

// GenerateCacheKey 生成缓存 Key
// 格式: astro:{userId}:{granularity}:{timeKey}:{language}
func GenerateCacheKey(userID, granularity string, t time.Time, language ...string) string {
	var timeKey string
	switch granularity {
	case core.GranularityHour:
		timeKey = t.Format("2006010215") // 2026020514
	case core.GranularityDay:
		timeKey = t.Format("20060102") // 20260205
	case core.GranularityWeek:
		year, week := t.ISOWeek()
		timeKey = formatWeekKey(year, week) // 2026W06
	case core.GranularityMonth:
		timeKey = t.Format("200601") // 202602
	case core.GranularityYear:
		timeKey = t.Format("2006") // 2026
	default:
		timeKey = t.Format("2006010215")
	}

	// 如果提供了语言参数，添加到缓存键
	lang := "en"
	if len(language) > 0 && language[0] != "" {
		lang = language[0]
	}
	return "astro:" + userID + ":" + granularity + ":" + timeKey + ":" + lang
}

// formatWeekKey 格式化周 Key
func formatWeekKey(year, week int) string {
	weekStr := "0"
	if week < 10 {
		weekStr = "0" + string(rune('0'+week))
	} else {
		weekStr = string(rune('0'+week/10)) + string(rune('0'+week%10))
	}
	return string(rune('0'+year/1000)) + string(rune('0'+(year/100)%10)) + string(rune('0'+(year/10)%10)) + string(rune('0'+year%10)) + "W" + weekStr
}

// ==================== 多级缓存管理器 ====================

// MultiLevelCache 多级缓存管理器
type MultiLevelCache struct {
	L1 Cache // 内存缓存
	L2 Cache // 文件/Redis 缓存
}

// NewMultiLevelCache 创建多级缓存
func NewMultiLevelCache(l1, l2 Cache) *MultiLevelCache {
	return &MultiLevelCache{
		L1: l1,
		L2: l2,
	}
}

// Get 获取缓存（先 L1，再 L2）
func (m *MultiLevelCache) Get(key string) (*core.TimeSlot, bool) {
	// 先查 L1
	if slot, ok := m.L1.Get(key); ok {
		return slot, true
	}

	// L1 miss，查 L2
	if m.L2 != nil {
		if slot, ok := m.L2.Get(key); ok {
			// 回填 L1
			m.L1.Set(key, slot, time.Hour)
			return slot, true
		}
	}

	return nil, false
}

// Set 设置缓存（同时写入 L1 和 L2）
func (m *MultiLevelCache) Set(key string, value *core.TimeSlot, ttl time.Duration) {
	m.L1.Set(key, value, ttl)
	if m.L2 != nil {
		m.L2.Set(key, value, ttl*24) // L2 TTL 更长
	}
}

// Delete 删除缓存
func (m *MultiLevelCache) Delete(key string) {
	m.L1.Delete(key)
	if m.L2 != nil {
		m.L2.Delete(key)
	}
}

// Clear 清空缓存
func (m *MultiLevelCache) Clear() {
	m.L1.Clear()
	if m.L2 != nil {
		m.L2.Clear()
	}
}

// Size 获取 L1 缓存大小
func (m *MultiLevelCache) Size() int {
	return m.L1.Size()
}

// ==================== 默认 TTL 配置 ====================

// DefaultTTL 获取默认 TTL
func DefaultTTL(granularity string) time.Duration {
	switch granularity {
	case core.GranularityHour:
		return time.Hour
	case core.GranularityDay:
		return 6 * time.Hour
	case core.GranularityWeek:
		return 12 * time.Hour
	case core.GranularityMonth:
		return 24 * time.Hour
	case core.GranularityYear:
		return 24 * time.Hour
	default:
		return time.Hour
	}
}
