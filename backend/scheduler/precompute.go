package scheduler

import (
	"log"
	"star/astro"
	"star/cache"
	"star/core"
	"star/models"
	"sync"
	"time"
)

// ==================== 预计算调度器 ====================
// 在后台预计算热门数据，提升查询性能

// PrecomputeScheduler 预计算调度器
type PrecomputeScheduler struct {
	mu       sync.Mutex
	running  bool
	stopChan chan struct{}
	users    map[string]*models.NatalChart // 已注册用户的本命盘
}

// NewPrecomputeScheduler 创建调度器
func NewPrecomputeScheduler() *PrecomputeScheduler {
	return &PrecomputeScheduler{
		users:    make(map[string]*models.NatalChart),
		stopChan: make(chan struct{}),
	}
}

// Start 启动调度器
func (ps *PrecomputeScheduler) Start() {
	ps.mu.Lock()
	if ps.running {
		ps.mu.Unlock()
		return
	}
	ps.running = true
	ps.mu.Unlock()

	log.Println("🚀 Starting precompute scheduler...")

	// 启动后台任务
	go ps.runScheduler()
}

// Stop 停止调度器
func (ps *PrecomputeScheduler) Stop() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.running {
		return
	}

	ps.running = false
	close(ps.stopChan)
	log.Println("🛑 Precompute scheduler stopped")
}

// RegisterUser 注册用户（用于预热）
func (ps *PrecomputeScheduler) RegisterUser(userID string, chart *models.NatalChart) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.users[userID] = chart
}

// UnregisterUser 取消注册用户
func (ps *PrecomputeScheduler) UnregisterUser(userID string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.users, userID)
}

// runScheduler 运行调度器
func (ps *PrecomputeScheduler) runScheduler() {
	// 每小时执行一次预热
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	// 启动时立即执行一次
	ps.precomputeAll()

	for {
		select {
		case <-ticker.C:
			ps.precomputeAll()
		case <-ps.stopChan:
			return
		}
	}
}

// precomputeAll 预计算所有用户的数据
func (ps *PrecomputeScheduler) precomputeAll() {
	ps.mu.Lock()
	users := make(map[string]*models.NatalChart)
	for k, v := range ps.users {
		users[k] = v
	}
	ps.mu.Unlock()

	log.Printf("🔄 Precomputing data for %d users...", len(users))

	for userID, chart := range users {
		ps.precomputeUser(userID, chart)
	}

	log.Println("✅ Precompute completed")
}

// precomputeUser 预计算单个用户的数据
func (ps *PrecomputeScheduler) precomputeUser(userID string, chart *models.NatalChart) {
	globalCache := cache.GetGlobalCache()
	calculator := core.NewCalculator(chart, "en")
	aggregator := core.NewAggregator(calculator)

	now := time.Now()

	// 预热未来 7 天的日级数据
	for day := 0; day < 7; day++ {
		targetDate := now.AddDate(0, 0, day)
		cacheKey := cache.GenerateCacheKey(userID, core.GranularityDay, targetDate)

		// 检查是否已缓存
		if _, ok := globalCache.Get(cacheKey); ok {
			continue
		}

		// 计算并缓存
		slot := aggregator.AggregateDay(targetDate)
		globalCache.Set(cacheKey, slot, cache.DefaultTTL(core.GranularityDay))
	}

	// 预热当前小时
	cacheKey := cache.GenerateCacheKey(userID, core.GranularityHour, now)
	if _, ok := globalCache.Get(cacheKey); !ok {
		slot := calculator.CalculateHour(now)
		deltaCalc := core.NewDeltaCalculator(calculator)
		deltaCalc.ApplyDeltaToSlot(slot, core.GranularityHour, now)
		globalCache.Set(cacheKey, slot, cache.DefaultTTL(core.GranularityHour))
	}
}

// PrecomputeOnDemand 按需预热（用户首次查询时调用）
func (ps *PrecomputeScheduler) PrecomputeOnDemand(birth models.BirthData) {
	chart := astro.CalculateNatalChart(birth)
	if chart == nil {
		return
	}

	userID := generateUserID(birth)
	globalCache := cache.GetGlobalCache()
	calculator := core.NewCalculator(chart, "en")
	aggregator := core.NewAggregator(calculator)

	now := time.Now()

	// 预热未来 24 小时的小时级数据
	go func() {
		for hour := 0; hour < 24; hour++ {
			targetTime := now.Add(time.Duration(hour) * time.Hour)
			cacheKey := cache.GenerateCacheKey(userID, core.GranularityHour, targetTime)

			if _, ok := globalCache.Get(cacheKey); ok {
				continue
			}

			slot := calculator.CalculateHour(targetTime)
			deltaCalc := core.NewDeltaCalculator(calculator)
			deltaCalc.ApplyDeltaToSlot(slot, core.GranularityHour, targetTime)
			globalCache.Set(cacheKey, slot, cache.DefaultTTL(core.GranularityHour))
		}

		// 预热未来 7 天的日级数据
		for day := 0; day < 7; day++ {
			targetDate := now.AddDate(0, 0, day)
			cacheKey := cache.GenerateCacheKey(userID, core.GranularityDay, targetDate)

			if _, ok := globalCache.Get(cacheKey); ok {
				continue
			}

			slot := aggregator.AggregateDay(targetDate)
			globalCache.Set(cacheKey, slot, cache.DefaultTTL(core.GranularityDay))
		}
	}()
}

// generateUserID 生成用户 ID
func generateUserID(birth models.BirthData) string {
	return formatInt(birth.Year) + formatInt(birth.Month) + formatInt(birth.Day) + "_" +
		formatInt(birth.Hour) + formatInt(birth.Minute) + "_" +
		formatFloat(birth.Latitude) + "_" + formatFloat(birth.Longitude)
}

func formatInt(i int) string {
	if i < 10 {
		return "0" + string(rune('0'+i))
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}

func formatFloat(f float64) string {
	if f < 0 {
		return "n" + formatPositiveFloat(-f)
	}
	return formatPositiveFloat(f)
}

func formatPositiveFloat(f float64) string {
	intPart := int(f)
	decPart := int((f - float64(intPart)) * 100)
	return formatInt(intPart) + "p" + formatInt(decPart)
}

// ==================== 全局调度器实例 ====================

var (
	globalScheduler     *PrecomputeScheduler
	globalSchedulerOnce sync.Once
)

// GetGlobalScheduler 获取全局调度器
func GetGlobalScheduler() *PrecomputeScheduler {
	globalSchedulerOnce.Do(func() {
		globalScheduler = NewPrecomputeScheduler()
	})
	return globalScheduler
}
