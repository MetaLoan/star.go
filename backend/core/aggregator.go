package core

import (
	"sort"
	"time"
)

// ==================== 时间聚合器 ====================
// 将小时级数据聚合为日/周/月/年级数据

// Aggregator 时间聚合器
type Aggregator struct {
	calculator      *Calculator
	deltaCalculator *DeltaCalculator
}

// NewAggregator 创建聚合器
func NewAggregator(calculator *Calculator) *Aggregator {
	return &Aggregator{
		calculator:      calculator,
		deltaCalculator: NewDeltaCalculator(calculator),
	}
}

// AggregateDay 聚合为日级数据
func (a *Aggregator) AggregateDay(t time.Time) *TimeSlot {
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	dayEnd := dayStart.AddDate(0, 0, 1)

	// 计算 24 个小时的数据
	hourSlots := make([]*TimeSlot, 24)
	for hour := 0; hour < 24; hour++ {
		hourTime := dayStart.Add(time.Duration(hour) * time.Hour)
		hourSlots[hour] = a.calculator.CalculateHour(hourTime)
	}

	// 创建日级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), dayStart, dayEnd, GranularityDay)

	// 聚合分数
	slot.Scores = a.aggregateScores(hourSlots)

	// 设置子槽的粒度，以便后续推断父槽粒度
	for i := range hourSlots {
		if hourSlots[i] != nil {
			hourSlots[i].Granularity = GranularityHour
		}
	}

	// 合并事件（去重并按粒度过滤）
	slot.Events = a.mergeAndDeduplicateEvents(hourSlots)

	// 生成子时间槽（用于曲线）
	for _, hs := range hourSlots {
		slot.SubSlots = append(slot.SubSlots, SubSlot{
			StartTime:  hs.StartTime,
			Scores:     hs.Scores,
			EventCount: len(hs.Events),
		})
	}

	// 计算事件的 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityDay, t)

	// 计算 slot 级别的 delta（与前一天对比）
	prevDay := dayStart.AddDate(0, 0, -1)
	prevSlot := a.calculateDaySlotWithoutDelta(prevDay)
	slot.Delta = a.deltaCalculator.CalculateSlotDelta(slot, prevSlot)

	// 生成指导
	slot.Guidance = a.generateDayGuidance(slot)

	return slot
}

// AggregateWeek 聚合为周级数据
// 性能优化：使用 3 天采样（周一、周三、周六）
func (a *Aggregator) AggregateWeek(t time.Time) *TimeSlot {
	// 获取本周一
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := t.AddDate(0, 0, -(weekday - 1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, t.Location())
	weekEnd := weekStart.AddDate(0, 0, 7)

	// 采样 3 个关键天（周一、周三、周六）
	sampleDays := []int{0, 2, 5} // 0=周一, 2=周三, 5=周六
	sampleSlots := make([]*TimeSlot, len(sampleDays))
	for i, day := range sampleDays {
		noonTime := weekStart.AddDate(0, 0, day).Add(12 * time.Hour)
		sampleSlots[i] = a.calculator.CalculateHour(noonTime)
		sampleSlots[i].StartTime = weekStart.AddDate(0, 0, day)
		sampleSlots[i].Granularity = GranularityDay
	}

	// 创建周级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), weekStart, weekEnd, GranularityWeek)

	// 聚合分数（基于采样点）
	slot.Scores = a.aggregateScoresFromSlots(sampleSlots)

	// 独立获取事件（使用已有的采样点，不额外计算）
	slot.Events = a.getEventsFromSlots(sampleSlots, weekStart, weekEnd, GranularityWeek)

	// 生成子时间槽（3 个采样点）
	for _, ss := range sampleSlots {
		slot.SubSlots = append(slot.SubSlots, SubSlot{
			StartTime:  ss.StartTime,
			Scores:     ss.Scores,
			EventCount: 0,
		})
	}

	// 计算事件的 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityWeek, t)

	// 计算前一周 delta（使用简化采样）
	prevWeekStart := weekStart.AddDate(0, 0, -7)
	prevSlot := a.calculateWeekSlotWithoutDelta(prevWeekStart)
	slot.Delta = a.deltaCalculator.CalculateSlotDelta(slot, prevSlot)

	// 生成指导
	slot.Guidance = a.generateWeekGuidance(slot)

	return slot
}

// AggregateMonth 聚合为月级数据
// 性能优化：每周采样一次（720次 -> 5次）
func (a *Aggregator) AggregateMonth(t time.Time) *TimeSlot {
	monthStart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	// 每周采样一次（第 1、8、15、22、29 天，最多 5 个点）
	sampleDays := []int{1, 8, 15, 22}
	daysInMonth := int(monthEnd.Sub(monthStart).Hours() / 24)
	if daysInMonth >= 29 {
		sampleDays = append(sampleDays, 29)
	}

	sampleSlots := make([]*TimeSlot, len(sampleDays))
	for i, day := range sampleDays {
		noonTime := time.Date(monthStart.Year(), monthStart.Month(), day, 12, 0, 0, 0, monthStart.Location())
		sampleSlots[i] = a.calculator.CalculateHour(noonTime)
		sampleSlots[i].StartTime = time.Date(monthStart.Year(), monthStart.Month(), day, 0, 0, 0, 0, monthStart.Location())
		sampleSlots[i].Granularity = GranularityDay
	}

	// 创建月级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), monthStart, monthEnd, GranularityMonth)

	// 聚合分数（基于采样点）
	slot.Scores = a.aggregateScoresFromSlots(sampleSlots)

	// 独立获取事件
	slot.Events = a.getEventsForTimeRange(monthStart, monthEnd, GranularityMonth)

	// 生成子时间槽（每周一个点）
	for _, ss := range sampleSlots {
		slot.SubSlots = append(slot.SubSlots, SubSlot{
			StartTime:  ss.StartTime,
			Scores:     ss.Scores,
			EventCount: 0,
		})
	}

	// 计算事件的 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityMonth, t)

	// 计算前一月 delta
	prevMonthStart := monthStart.AddDate(0, -1, 0)
	prevSlot := a.calculateMonthSlotWithoutDelta(prevMonthStart)
	slot.Delta = a.deltaCalculator.CalculateSlotDelta(slot, prevSlot)

	// 生成指导
	slot.Guidance = a.generateMonthGuidance(slot)

	return slot
}

// AggregateYear 聚合为年级数据
// 性能优化：使用季度采样（12次 -> 4次）
func (a *Aggregator) AggregateYear(t time.Time) *TimeSlot {
	yearStart := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	yearEnd := yearStart.AddDate(1, 0, 0)

	// 每季度采样中间月份（2月、5月、8月、11月）的15号中午
	sampleMonths := []int{2, 5, 8, 11}
	sampleSlots := make([]*TimeSlot, len(sampleMonths))
	for i, month := range sampleMonths {
		sampleTime := time.Date(t.Year(), time.Month(month), 15, 12, 0, 0, 0, t.Location())
		sampleSlots[i] = a.calculator.CalculateHour(sampleTime)
		// 设置为季度起始月份（1月、4月、7月、10月）
		quarterStartMonth := ((month-1)/3)*3 + 1
		sampleSlots[i].StartTime = time.Date(t.Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, t.Location())
		sampleSlots[i].Granularity = GranularityMonth
	}

	// 创建年级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), yearStart, yearEnd, GranularityYear)

	// 聚合分数（基于采样点）
	slot.Scores = a.aggregateScoresFromSlots(sampleSlots)

	// 独立获取事件
	slot.Events = a.getEventsForTimeRange(yearStart, yearEnd, GranularityYear)

	// 生成子时间槽（每季度一个点）
	for _, ss := range sampleSlots {
		slot.SubSlots = append(slot.SubSlots, SubSlot{
			StartTime:  ss.StartTime,
			Scores:     ss.Scores,
			EventCount: 0,
		})
	}

	// 计算事件的 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityYear, t)

	// 计算前一年 delta
	prevYearStart := yearStart.AddDate(-1, 0, 0)
	prevSlot := a.calculateYearSlotWithoutDelta(prevYearStart)
	slot.Delta = a.deltaCalculator.CalculateSlotDelta(slot, prevSlot)

	// 生成指导
	slot.Guidance = a.generateYearGuidance(slot)

	return slot
}

// ==================== 分数聚合 ====================

// aggregateScores 从小时槽聚合分数
func (a *Aggregator) aggregateScores(hourSlots []*TimeSlot) DimensionScores {
	if len(hourSlots) == 0 {
		return DimensionScores{}
	}

	var sum DimensionScores
	count := 0

	for _, slot := range hourSlots {
		if slot != nil {
			sum.Overall += slot.Scores.Overall
			sum.Career += slot.Scores.Career
			sum.Relationship += slot.Scores.Relationship
			sum.Health += slot.Scores.Health
			sum.Finance += slot.Scores.Finance
			sum.Spiritual += slot.Scores.Spiritual
			count++
		}
	}

	if count == 0 {
		return DimensionScores{}
	}

	return DimensionScores{
		Overall:      sum.Overall / float64(count),
		Career:       sum.Career / float64(count),
		Relationship: sum.Relationship / float64(count),
		Health:       sum.Health / float64(count),
		Finance:      sum.Finance / float64(count),
		Spiritual:    sum.Spiritual / float64(count),
	}
}

// aggregateScoresFromSlots 从时间槽切片聚合分数
func (a *Aggregator) aggregateScoresFromSlots(slots []*TimeSlot) DimensionScores {
	if len(slots) == 0 {
		return DimensionScores{}
	}

	var sum DimensionScores
	count := 0

	for _, slot := range slots {
		if slot != nil {
			sum.Overall += slot.Scores.Overall
			sum.Career += slot.Scores.Career
			sum.Relationship += slot.Scores.Relationship
			sum.Health += slot.Scores.Health
			sum.Finance += slot.Scores.Finance
			sum.Spiritual += slot.Scores.Spiritual
			count++
		}
	}

	if count == 0 {
		return DimensionScores{}
	}

	return DimensionScores{
		Overall:      sum.Overall / float64(count),
		Career:       sum.Career / float64(count),
		Relationship: sum.Relationship / float64(count),
		Health:       sum.Health / float64(count),
		Finance:      sum.Finance / float64(count),
		Spiritual:    sum.Spiritual / float64(count),
	}
}

// ==================== 事件合并与去重 ====================

// mergeAndDeduplicateEvents 合并并去重事件
func (a *Aggregator) mergeAndDeduplicateEvents(hourSlots []*TimeSlot) []AstroEvent {
	eventMap := make(map[string]AstroEvent)
	queryGranularity := GranularityDay // 这个方法只在 AggregateDay 中被调用

	for _, slot := range hourSlots {
		if slot == nil {
			continue
		}
		for _, event := range slot.Events {
			// 粒度过滤：只包含符合当前查询粒度权重的事件
			if !a.shouldIncludeEvent(event, queryGranularity) {
				continue
			}

			// 使用 eventID 去重
			if existing, ok := eventMap[event.EventID]; ok {
				// 如果事件已存在，保留强度更高的
				if event.Intensity > existing.Intensity {
					eventMap[event.EventID] = event
				}
			} else {
				eventMap[event.EventID] = event
			}
		}
	}

	// 转换为切片
	events := make([]AstroEvent, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, event)
	}

	// 按时间排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].ExactTime.Before(events[j].ExactTime)
	})

	return events
}

// mergeEventsFromSlots 从时间槽切片合并事件
func (a *Aggregator) mergeEventsFromSlots(slots []*TimeSlot) []AstroEvent {
	eventMap := make(map[string]AstroEvent)
	var queryGranularity string
	if len(slots) > 0 && slots[0] != nil {
		// 根据子槽的粒度推断父槽粒度
		subGran := slots[0].Granularity
		switch subGran {
		case GranularityHour:
			queryGranularity = GranularityDay
		case GranularityDay:
			if len(slots) == 7 {
				queryGranularity = GranularityWeek
			} else {
				queryGranularity = GranularityMonth
			}
		case GranularityMonth:
			queryGranularity = GranularityYear
		default:
			queryGranularity = GranularityDay // 安全回退
		}
	}

	for _, slot := range slots {
		if slot == nil {
			continue
		}
		for _, event := range slot.Events {
			// 粒度过滤
			if !a.shouldIncludeEvent(event, queryGranularity) {
				continue
			}

			if existing, ok := eventMap[event.EventID]; ok {
				if event.Intensity > existing.Intensity {
					eventMap[event.EventID] = event
				}
			} else {
				eventMap[event.EventID] = event
			}
		}
	}

	events := make([]AstroEvent, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].ExactTime.Before(events[j].ExactTime)
	})

	return events
}

// getEventsForTimeRange 获取时间范围内的事件（性能优化版）
// 只计算几个采样点的事件并合并去重，而非遍历每小时
func (a *Aggregator) getEventsForTimeRange(start, end time.Time, granularity string) []AstroEvent {
	eventMap := make(map[string]AstroEvent)

	// 计算采样点：起点、中点、终点
	midPoint := start.Add(end.Sub(start) / 2)
	samplePoints := []time.Time{start, midPoint, end.Add(-time.Hour)}

	for _, t := range samplePoints {
		slot := a.calculator.CalculateHour(t)
		if slot == nil {
			continue
		}
		for _, event := range slot.Events {
			// 粒度过滤
			if !a.shouldIncludeEvent(event, granularity) {
				continue
			}
			// 检查事件是否与时间范围有交集
			if event.EndTime.Before(start) || event.StartTime.After(end) {
				continue
			}
			if existing, ok := eventMap[event.EventID]; ok {
				if event.Intensity > existing.Intensity {
					eventMap[event.EventID] = event
				}
			} else {
				eventMap[event.EventID] = event
			}
		}
	}

	events := make([]AstroEvent, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].ExactTime.Before(events[j].ExactTime)
	})

	return events
}

// getEventsFromSlots 从已计算的 slots 中提取事件（避免重复计算）
func (a *Aggregator) getEventsFromSlots(slots []*TimeSlot, start, end time.Time, granularity string) []AstroEvent {
	eventMap := make(map[string]AstroEvent)

	for _, slot := range slots {
		if slot == nil {
			continue
		}
		for _, event := range slot.Events {
			// 粒度过滤
			if !a.shouldIncludeEvent(event, granularity) {
				continue
			}
			// 检查事件是否与时间范围有交集
			if event.EndTime.Before(start) || event.StartTime.After(end) {
				continue
			}
			if existing, ok := eventMap[event.EventID]; ok {
				if event.Intensity > existing.Intensity {
					eventMap[event.EventID] = event
				}
			} else {
				eventMap[event.EventID] = event
			}
		}
	}

	events := make([]AstroEvent, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].ExactTime.Before(events[j].ExactTime)
	})

	return events
}

// shouldIncludeEvent 判断事件是否应包含在指定查询粒度中
// 规则：事件的 timeLevel 必须 >= 查询粒度
// 例如：day 粒度只显示 daily、weekly、monthly、yearly 事件，不显示 hourly
func (a *Aggregator) shouldIncludeEvent(event AstroEvent, queryGranularity string) bool {
	// 定义 timeLevel 权重（越高=生命周期越长）
	levelWeight := map[string]int{
		TimeLevelHourly:  1,
		TimeLevelDaily:   2,
		TimeLevelWeekly:  3,
		TimeLevelMonthly: 4,
		TimeLevelYearly:  5,
	}

	// 定义查询粒度对应的最小 timeLevel 门槛
	minThreshold := map[string]int{
		GranularityHour:  1, // 小时视图：显示所有事件 (hourly+)
		GranularityDay:   2, // 天视图：只显示 daily 及以上 (daily, weekly, monthly, yearly)
		GranularityWeek:  3, // 周视图：只显示 weekly 及以上 (weekly, monthly, yearly)
		GranularityMonth: 4, // 月视图：只显示 monthly 及以上 (monthly, yearly)
		GranularityYear:  5, // 年视图：只显示 yearly
	}

	// 获取事件的 timeLevel 权重
	weight, ok1 := levelWeight[event.TimeLevel]
	if !ok1 {
		// 未知 timeLevel，默认不包含
		return false
	}

	// 获取查询粒度的门槛
	threshold, ok2 := minThreshold[queryGranularity]
	if !ok2 {
		// 未知粒度，默认使用 day 门槛
		threshold = 2
	}

	// 事件的 timeLevel 权重必须 >= 查询粒度门槛
	return weight >= threshold
}

// ==================== 指导生成 ====================

func (a *Aggregator) generateDayGuidance(slot *TimeSlot) *Guidance {
	return a.generateGuidance(slot, "day")
}

func (a *Aggregator) generateWeekGuidance(slot *TimeSlot) *Guidance {
	return a.generateGuidance(slot, "week")
}

func (a *Aggregator) generateMonthGuidance(slot *TimeSlot) *Guidance {
	return a.generateGuidance(slot, "month")
}

func (a *Aggregator) generateYearGuidance(slot *TimeSlot) *Guidance {
	return a.generateGuidance(slot, "year")
}

func (a *Aggregator) generateGuidance(slot *TimeSlot, period string) *Guidance {
	guidance := &Guidance{
		Dos:   make([]string, 0),
		Donts: make([]string, 0),
	}

	// 找出最重要的正面和负面事件
	var topPositive, topNegative *AstroEvent
	for i := range slot.Events {
		e := &slot.Events[i]
		if e.IsPositive && (topPositive == nil || e.Intensity > topPositive.Intensity) {
			topPositive = e
		}
		if !e.IsPositive && (topNegative == nil || e.Intensity > topNegative.Intensity) {
			topNegative = e
		}
	}

	// 生成摘要
	if topPositive != nil && topNegative != nil {
		guidance.Summary = "This " + period + " features " + topPositive.Title + ", but be mindful of " + topNegative.Title
	} else if topPositive != nil {
		guidance.Summary = "This " + period + " is favorable for activities related to " + topPositive.Title
	} else if topNegative != nil {
		guidance.Summary = "This " + period + " requires caution regarding " + topNegative.Title
	} else {
		guidance.Summary = "A relatively stable " + period + " without major events"
	}

	// 确定重点关注维度
	guidance.Focus = a.findHighestDimension(slot.Scores)

	return guidance
}

// findHighestDimension 找出最高分的维度
func (a *Aggregator) findHighestDimension(scores DimensionScores) string {
	dims := map[string]float64{
		"career":       scores.Career,
		"relationship": scores.Relationship,
		"health":       scores.Health,
		"finance":      scores.Finance,
		"spiritual":    scores.Spiritual,
	}

	maxDim := "career"
	maxScore := scores.Career
	for dim, score := range dims {
		if score > maxScore {
			maxScore = score
			maxDim = dim
		}
	}
	return maxDim
}

// ==================== 辅助聚合方法（不含 delta，用于对比前一周期） ====================

// calculateDaySlotWithoutDelta 计算日级 slot（不含 delta，避免递归）
// 性能优化：只采样 4 个点（每 6 小时一个）用于 delta 对比
func (a *Aggregator) calculateDaySlotWithoutDelta(t time.Time) *TimeSlot {
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	dayEnd := dayStart.AddDate(0, 0, 1)

	// 性能优化：只采样 4 个点（0, 6, 12, 18 点）
	sampleHours := []int{0, 6, 12, 18}
	sampleSlots := make([]*TimeSlot, len(sampleHours))
	for i, hour := range sampleHours {
		hourTime := dayStart.Add(time.Duration(hour) * time.Hour)
		sampleSlots[i] = a.calculator.CalculateHour(hourTime)
	}

	// 创建日级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), dayStart, dayEnd, GranularityDay)

	// 聚合分数
	slot.Scores = a.aggregateScores(sampleSlots)

	return slot
}

// calculateWeekSlotWithoutDelta 计算周级 slot（不含 delta，避免递归）
// 性能优化：使用每天中午采样
func (a *Aggregator) calculateWeekSlotWithoutDelta(weekStart time.Time) *TimeSlot {
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	weekEnd := weekStart.AddDate(0, 0, 7)

	// 只采样 2 个点（周二、周五）用于 delta 计算
	sampleDays := []int{1, 4} // 周二、周五
	sampleSlots := make([]*TimeSlot, len(sampleDays))
	for i, day := range sampleDays {
		noonTime := weekStart.AddDate(0, 0, day).Add(12 * time.Hour)
		sampleSlots[i] = a.calculator.CalculateHour(noonTime)
	}

	slot := NewTimeSlot(a.calculator.getUserID(), weekStart, weekEnd, GranularityWeek)
	slot.Scores = a.aggregateScores(sampleSlots)

	return slot
}

// calculateMonthSlotWithoutDelta 计算月级 slot（不含 delta，避免递归）
// 性能优化：只采样 4 个点
func (a *Aggregator) calculateMonthSlotWithoutDelta(monthStart time.Time) *TimeSlot {
	monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	sampleDays := []int{1, 8, 15, 22}
	sampleSlots := make([]*TimeSlot, len(sampleDays))
	for i, day := range sampleDays {
		sampleTime := time.Date(monthStart.Year(), monthStart.Month(), day, 12, 0, 0, 0, monthStart.Location())
		sampleSlots[i] = a.calculator.CalculateHour(sampleTime)
	}

	slot := NewTimeSlot(a.calculator.getUserID(), monthStart, monthEnd, GranularityMonth)
	slot.Scores = a.aggregateScores(sampleSlots)

	return slot
}

// calculateYearSlotWithoutDelta 计算年级 slot（不含 delta，避免递归）
// 性能优化：使用季度采样（4次）
func (a *Aggregator) calculateYearSlotWithoutDelta(yearStart time.Time) *TimeSlot {
	yearStart = time.Date(yearStart.Year(), 1, 1, 0, 0, 0, 0, yearStart.Location())
	yearEnd := yearStart.AddDate(1, 0, 0)

	// 每季度采样中间月份（2月、5月、8月、11月）
	sampleMonths := []int{2, 5, 8, 11}
	sampleSlots := make([]*TimeSlot, len(sampleMonths))
	for i, month := range sampleMonths {
		sampleTime := time.Date(yearStart.Year(), time.Month(month), 15, 12, 0, 0, 0, yearStart.Location())
		sampleSlots[i] = a.calculator.CalculateHour(sampleTime)
	}

	slot := NewTimeSlot(a.calculator.getUserID(), yearStart, yearEnd, GranularityYear)
	slot.Scores = a.aggregateScores(sampleSlots)

	return slot
}
