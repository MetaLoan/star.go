package cache

import (
	"star/core"
	"sync"
	"time"
)

// ==================== L1 内存缓存实现 ====================

// cacheEntry 缓存条目
type cacheEntry struct {
	value     *core.TimeSlot
	expiresAt time.Time
}

// genericCacheEntry 通用缓存条目
type genericCacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// MemoryCache 内存缓存（L1）
type MemoryCache struct {
	data        map[string]*cacheEntry
	genericData map[string]*genericCacheEntry // 通用缓存（用于趋势数据等）
	mu          sync.RWMutex
	maxSize     int // 最大条目数
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(maxSize int) *MemoryCache {
	if maxSize <= 0 {
		maxSize = 10000 // 默认 10000 条
	}
	mc := &MemoryCache{
		data:        make(map[string]*cacheEntry),
		genericData: make(map[string]*genericCacheEntry),
		maxSize:     maxSize,
	}
	// 启动后台清理过期条目
	go mc.startCleanup()
	return mc
}

// Get 获取缓存
func (mc *MemoryCache) Get(key string) (*core.TimeSlot, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, ok := mc.data[key]
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.value, true
}

// Set 设置缓存
func (mc *MemoryCache) Set(key string, value *core.TimeSlot, ttl time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// 如果达到最大容量，先清理一些旧数据
	if len(mc.data) >= mc.maxSize {
		mc.evictOldest()
	}

	mc.data[key] = &cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Delete 删除缓存
func (mc *MemoryCache) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	delete(mc.data, key)
}

// Clear 清空缓存
func (mc *MemoryCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data = make(map[string]*cacheEntry)
	mc.genericData = make(map[string]*genericCacheEntry)
}

// GetGeneric 获取通用缓存
func (mc *MemoryCache) GetGeneric(key string) (interface{}, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, ok := mc.genericData[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.value, true
}

// SetGeneric 设置通用缓存
func (mc *MemoryCache) SetGeneric(key string, value interface{}, ttl time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.genericData[key] = &genericCacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Size 获取缓存大小
func (mc *MemoryCache) Size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return len(mc.data)
}

// evictOldest 清理最旧的 10% 条目
func (mc *MemoryCache) evictOldest() {
	now := time.Now()
	toDelete := make([]string, 0)

	// 先清理过期的
	for key, entry := range mc.data {
		if now.After(entry.expiresAt) {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		delete(mc.data, key)
	}

	// 如果还是太多，删除 10%
	if len(mc.data) >= mc.maxSize {
		count := 0
		target := mc.maxSize / 10
		for key := range mc.data {
			if count >= target {
				break
			}
			delete(mc.data, key)
			count++
		}
	}
}

// startCleanup 启动后台清理
func (mc *MemoryCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mc.cleanupExpired()
	}
}

// cleanupExpired 清理过期条目
func (mc *MemoryCache) cleanupExpired() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	for key, entry := range mc.data {
		if now.After(entry.expiresAt) {
			delete(mc.data, key)
		}
	}
	for key, entry := range mc.genericData {
		if now.After(entry.expiresAt) {
			delete(mc.genericData, key)
		}
	}
}

// ==================== 全局缓存实例 ====================

var (
	globalCache     *MultiLevelCache
	globalCacheOnce sync.Once
)

// GetGlobalCache 获取全局缓存实例
func GetGlobalCache() *MultiLevelCache {
	globalCacheOnce.Do(func() {
		l1 := NewMemoryCache(10000)
		// L2 暂时不实现，后续可以添加 Redis 或文件缓存
		globalCache = NewMultiLevelCache(l1, nil)
	})
	return globalCache
}

// SetGlobalCache 设置全局缓存实例（用于测试）
func SetGlobalCache(cache *MultiLevelCache) {
	globalCache = cache
}
