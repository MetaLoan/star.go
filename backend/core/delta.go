package core

import (
	"time"
)

// ==================== impactDelta 动态计算器 ====================
// 根据查询粒度计算事件相对于上一周期的影响变化

// DeltaCalculator impactDelta 计算器
type DeltaCalculator struct {
	calculator *Calculator
}

// NewDeltaCalculator 创建 delta 计算器
func NewDeltaCalculator(calculator *Calculator) *DeltaCalculator {
	return &DeltaCalculator{
		calculator: calculator,
	}
}

// CalculateImpactDelta 计算事件相对上一周期的影响变化
// 参数：
//   - event: 当前事件
//   - granularity: 查询粒度 (hour/day/week/month/year)
//   - queryTime: 查询时间点
func (dc *DeltaCalculator) CalculateImpactDelta(event *AstroEvent, granularity string, queryTime time.Time) DimensionImpact {
	// 1. 获取上一周期的时间范围
	prevStart, prevEnd := GetPreviousPeriod(queryTime, granularity)

	// 2. 获取该事件在上一周期的影响
	prevImpact := dc.getEventImpactInPeriod(event, prevStart, prevEnd)

	// 3. 获取当前周期的影响
	currStart, currEnd := GetCurrentPeriod(queryTime, granularity)
	currImpact := dc.getEventImpactInPeriod(event, currStart, currEnd)

	// 4. 计算差值
	return currImpact.Subtract(prevImpact)
}

// CalculateSlotDelta 计算整个时间槽相对上一周期的分数变化
func (dc *DeltaCalculator) CalculateSlotDelta(currentSlot *TimeSlot, prevSlot *TimeSlot) *ScoreDelta {
	if prevSlot == nil {
		return &ScoreDelta{
			Overall: currentSlot.Scores.Overall,
			Dimensions: DimensionScores{
				Career:       currentSlot.Scores.Career,
				Relationship: currentSlot.Scores.Relationship,
				Health:       currentSlot.Scores.Health,
				Finance:      currentSlot.Scores.Finance,
				Spiritual:    currentSlot.Scores.Spiritual,
			},
			Reason: "No previous period data",
		}
	}

	delta := &ScoreDelta{
		Overall: currentSlot.Scores.Overall - prevSlot.Scores.Overall,
		Dimensions: DimensionScores{
			Career:       currentSlot.Scores.Career - prevSlot.Scores.Career,
			Relationship: currentSlot.Scores.Relationship - prevSlot.Scores.Relationship,
			Health:       currentSlot.Scores.Health - prevSlot.Scores.Health,
			Finance:      currentSlot.Scores.Finance - prevSlot.Scores.Finance,
			Spiritual:    currentSlot.Scores.Spiritual - prevSlot.Scores.Spiritual,
		},
	}

	// 确定主要变化原因
	delta.Reason = dc.findMainChangeReason(currentSlot, prevSlot)

	return delta
}

// getEventImpactInPeriod 获取事件在指定时间段内的影响
func (dc *DeltaCalculator) getEventImpactInPeriod(event *AstroEvent, periodStart, periodEnd time.Time) DimensionImpact {
	eventStart := event.StartTime
	eventEnd := event.EndTime

	// 无重叠：该周期没有这个事件
	if eventEnd.Before(periodStart) || eventStart.After(periodEnd) {
		return DimensionImpact{} // 全零
	}

	// 计算重叠比例
	overlapStart := maxTime(eventStart, periodStart)
	overlapEnd := minTime(eventEnd, periodEnd)
	overlapDuration := overlapEnd.Sub(overlapStart)
	periodDuration := periodEnd.Sub(periodStart)

	if periodDuration <= 0 {
		return event.Impact
	}

	overlapRatio := float64(overlapDuration) / float64(periodDuration)

	// 按重叠比例缩放影响
	return event.Impact.Scale(overlapRatio)
}

// findMainChangeReason 找出主要变化原因
func (dc *DeltaCalculator) findMainChangeReason(current, prev *TimeSlot) string {
	// 找出变化最大的事件
	if len(current.Events) == 0 {
		return "No significant events"
	}

	// 简单策略：返回强度最高的事件
	var maxEvent *AstroEvent
	maxIntensity := 0.0
	for i := range current.Events {
		e := &current.Events[i]
		if e.Intensity > maxIntensity {
			maxIntensity = e.Intensity
			maxEvent = e
		}
	}

	if maxEvent != nil {
		return maxEvent.Title
	}
	return "Multiple factors"
}

// ApplyDeltaToSlot 为时间槽中的所有事件计算 impactDelta
func (dc *DeltaCalculator) ApplyDeltaToSlot(slot *TimeSlot, granularity string, queryTime time.Time) {
	for i := range slot.Events {
		slot.Events[i].ImpactDelta = dc.CalculateImpactDelta(&slot.Events[i], granularity, queryTime)
	}
}

// ==================== 时间周期计算 ====================

// GetPreviousPeriod 获取上一周期的时间范围
func GetPreviousPeriod(t time.Time, granularity string) (start, end time.Time) {
	switch granularity {
	case GranularityHour:
		end = t.Truncate(time.Hour)
		start = end.Add(-time.Hour)
	case GranularityDay:
		end = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		start = end.AddDate(0, 0, -1)
	case GranularityWeek:
		// 获取本周一，然后减7天
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		monday := t.AddDate(0, 0, -(weekday - 1))
		end = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
		start = end.AddDate(0, 0, -7)
	case GranularityMonth:
		end = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		start = end.AddDate(0, -1, 0)
	case GranularityYear:
		end = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
		start = end.AddDate(-1, 0, 0)
	default:
		end = t.Truncate(time.Hour)
		start = end.Add(-time.Hour)
	}
	return
}

// GetCurrentPeriod 获取当前周期的时间范围
func GetCurrentPeriod(t time.Time, granularity string) (start, end time.Time) {
	switch granularity {
	case GranularityHour:
		start = t.Truncate(time.Hour)
		end = start.Add(time.Hour)
	case GranularityDay:
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end = start.AddDate(0, 0, 1)
	case GranularityWeek:
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		monday := t.AddDate(0, 0, -(weekday - 1))
		start = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
		end = start.AddDate(0, 0, 7)
	case GranularityMonth:
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		end = start.AddDate(0, 1, 0)
	case GranularityYear:
		start = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
		end = start.AddDate(1, 0, 0)
	default:
		start = t.Truncate(time.Hour)
		end = start.Add(time.Hour)
	}
	return
}

// GetPeriodDuration 获取周期持续时间
func GetPeriodDuration(granularity string) time.Duration {
	switch granularity {
	case GranularityHour:
		return time.Hour
	case GranularityDay:
		return 24 * time.Hour
	case GranularityWeek:
		return 7 * 24 * time.Hour
	case GranularityMonth:
		return 30 * 24 * time.Hour // 近似值
	case GranularityYear:
		return 365 * 24 * time.Hour // 近似值
	default:
		return time.Hour
	}
}

// ==================== 工具函数 ====================

// maxTime 返回较大的时间
func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// minTime 返回较小的时间
func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
