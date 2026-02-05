package core

import (
	"star/models"
	"testing"
	"time"
)

func TestDeltaCalculator_CalculateImpactDelta(t *testing.T) {
	// 准备模拟数据
	birth := models.BirthData{
		Year: 1990, Month: 5, Day: 15, Hour: 10, Minute: 30,
		Latitude: 39.9, Longitude: 116.4, Timezone: 8,
	}
	chart := &models.NatalChart{BirthData: birth}
	calc := NewCalculator(chart, "en")
	dc := NewDeltaCalculator(calc)

	queryTime := time.Date(2026, 1, 16, 14, 0, 0, 0, time.UTC)

	// 模拟一个事件，覆盖当前周期和上一周期
	// 当前周期 (Day): 2026-01-16 00:00 - 2026-01-17 00:00
	// 上一周期 (Day): 2026-01-15 00:00 - 2026-01-16 00:00
	event := &AstroEvent{
		EventID: "test_event",
		Type:    EventTypeAspect,
		Impact: DimensionImpact{
			Career: -5,
		},
		// 事件从 1月15日 12:00 开始，到 1月17日 12:00 结束
		// 在 1月15日 (上一周期) 覆盖 12 小时 (50%)
		// 在 1月16日 (当前周期) 覆盖 24 小时 (100%)
		StartTime: time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 1, 17, 12, 0, 0, 0, time.UTC),
	}

	// 1. 测试小时粒度
	// 当前小时: 14:00 - 15:00 (100% 覆盖)
	// 上一小时: 13:00 - 14:00 (100% 覆盖)
	// delta 应为 0
	deltaHour := dc.CalculateImpactDelta(event, GranularityHour, queryTime)
	if deltaHour.Career != 0 {
		t.Errorf("Expected hourly delta 0, got %f", deltaHour.Career)
	}

	// 2. 测试日粒度
	// 当前周期 (1月16日): 100% 覆盖 -> impact = -5
	// 上一周期 (1月15日): 50% 覆盖  -> impact = -2.5
	// delta = -5 - (-2.5) = -2.5
	deltaDay := dc.CalculateImpactDelta(event, GranularityDay, queryTime)
	expectedDelta := -2.5
	if deltaDay.Career != expectedDelta {
		t.Errorf("Expected daily delta %f, got %f", expectedDelta, deltaDay.Career)
	}
}

func TestGetPreviousPeriod(t *testing.T) {
	now := time.Date(2026, 1, 16, 14, 30, 0, 0, time.UTC)

	// Hour
	start, end := GetPreviousPeriod(now, GranularityHour)
	expectedStart := time.Date(2026, 1, 16, 13, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2026, 1, 16, 14, 0, 0, 0, time.UTC)
	if !start.Equal(expectedStart) || !end.Equal(expectedEnd) {
		t.Errorf("Hour: expected %v-%v, got %v-%v", expectedStart, expectedEnd, start, end)
	}

	// Day
	start, end = GetPreviousPeriod(now, GranularityDay)
	expectedStart = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedEnd = time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)
	if !start.Equal(expectedStart) || !end.Equal(expectedEnd) {
		t.Errorf("Day: expected %v-%v, got %v-%v", expectedStart, expectedEnd, start, end)
	}
}
