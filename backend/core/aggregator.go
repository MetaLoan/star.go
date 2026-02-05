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

	// 计算 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityDay, t)

	// 生成指导
	slot.Guidance = a.generateDayGuidance(slot)

	return slot
}

// AggregateWeek 聚合为周级数据
func (a *Aggregator) AggregateWeek(t time.Time) *TimeSlot {
	// 获取本周一
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := t.AddDate(0, 0, -(weekday - 1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, t.Location())
	weekEnd := weekStart.AddDate(0, 0, 7)

	// 计算 7 天的数据
	daySlots := make([]*TimeSlot, 7)
	for day := 0; day < 7; day++ {
		dayTime := weekStart.AddDate(0, 0, day)
		daySlots[day] = a.AggregateDay(dayTime)
	}

	// 创建周级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), weekStart, weekEnd, GranularityWeek)

	// 聚合分数
	slot.Scores = a.aggregateScoresFromSlots(daySlots)

	// 合并事件（去重）
	slot.Events = a.mergeEventsFromSlots(daySlots)

	// 生成子时间槽
	for _, ds := range daySlots {
		slot.SubSlots = append(slot.SubSlots, SubSlot{
			StartTime:  ds.StartTime,
			Scores:     ds.Scores,
			EventCount: len(ds.Events),
		})
	}

	// 计算 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityWeek, t)

	// 生成指导
	slot.Guidance = a.generateWeekGuidance(slot)

	return slot
}

// AggregateMonth 聚合为月级数据
func (a *Aggregator) AggregateMonth(t time.Time) *TimeSlot {
	monthStart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	daysInMonth := int(monthEnd.Sub(monthStart).Hours() / 24)

	// 计算每天的数据
	daySlots := make([]*TimeSlot, daysInMonth)
	for day := 0; day < daysInMonth; day++ {
		dayTime := monthStart.AddDate(0, 0, day)
		daySlots[day] = a.AggregateDay(dayTime)
	}

	// 创建月级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), monthStart, monthEnd, GranularityMonth)

	// 聚合分数
	slot.Scores = a.aggregateScoresFromSlots(daySlots)

	// 合并事件（去重）
	slot.Events = a.mergeEventsFromSlots(daySlots)

	// 生成子时间槽
	for _, ds := range daySlots {
		slot.SubSlots = append(slot.SubSlots, SubSlot{
			StartTime:  ds.StartTime,
			Scores:     ds.Scores,
			EventCount: len(ds.Events),
		})
	}

	// 计算 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityMonth, t)

	// 生成指导
	slot.Guidance = a.generateMonthGuidance(slot)

	return slot
}

// AggregateYear 聚合为年级数据
func (a *Aggregator) AggregateYear(t time.Time) *TimeSlot {
	yearStart := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	yearEnd := yearStart.AddDate(1, 0, 0)

	// 计算 12 个月的数据
	monthSlots := make([]*TimeSlot, 12)
	for month := 0; month < 12; month++ {
		monthTime := time.Date(t.Year(), time.Month(month+1), 15, 0, 0, 0, 0, t.Location())
		monthSlots[month] = a.AggregateMonth(monthTime)
	}

	// 创建年级时间槽
	slot := NewTimeSlot(a.calculator.getUserID(), yearStart, yearEnd, GranularityYear)

	// 聚合分数
	slot.Scores = a.aggregateScoresFromSlots(monthSlots)

	// 合并事件（去重）
	slot.Events = a.mergeEventsFromSlots(monthSlots)

	// 生成子时间槽（按月）
	for _, ms := range monthSlots {
		slot.SubSlots = append(slot.SubSlots, SubSlot{
			StartTime:  ms.StartTime,
			Scores:     ms.Scores,
			EventCount: len(ms.Events),
		})
	}

	// 计算 impactDelta
	a.deltaCalculator.ApplyDeltaToSlot(slot, GranularityYear, t)

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

// shouldIncludeEvent 判断事件是否应包含在指定查询粒度中
func (a *Aggregator) shouldIncludeEvent(event AstroEvent, queryGranularity string) bool {
	// 定义粒度等级权重
	levelWeight := map[string]int{
		TimeLevelHourly:  1,
		TimeLevelDaily:   2,
		TimeLevelWeekly:  3,
		TimeLevelMonthly: 4,
		TimeLevelYearly:  5,
	}

	// 定义查询粒度对应的最小显示门槛
	minThreshold := map[string]int{
		GranularityHour:  1, // 小时视图看所有
		GranularityDay:   2, // 天视图看 daily 及以上
		GranularityWeek:  3, // 周视图看 weekly 及以上
		GranularityMonth: 4, // 月视图看 monthly 及以上
		GranularityYear:  5, // 年视图只看 yearly
	}

	// 如果没有定义权重或门槛，默认包含（安全回退）
	weight, ok1 := levelWeight[event.TimeLevel]
	threshold, ok2 := minThreshold[queryGranularity]
	if !ok1 || !ok2 {
		return true
	}

	// 强制过滤逻辑：如果查询粒度大于 Hour，则过滤掉所有 Hourly 事件
	if queryGranularity != "" && queryGranularity != GranularityHour && event.TimeLevel == TimeLevelHourly {
		return false
	}

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
