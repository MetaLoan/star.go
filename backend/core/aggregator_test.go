package core

import (
	"star/models"
	"testing"
)

func TestAggregator_MergeAndDeduplicateEvents(t *testing.T) {
	birth := models.BirthData{
		Year: 1990, Month: 5, Day: 15, Hour: 10, Minute: 30,
		Latitude: 39.9, Longitude: 116.4, Timezone: 8,
	}
	chart := &models.NatalChart{BirthData: birth}
	calc := NewCalculator(chart, "en")
	agg := NewAggregator(calc)

	event1 := AstroEvent{EventID: "e1", Type: EventTypeAspect, Intensity: 0.5, TimeLevel: TimeLevelDaily}
	event2 := AstroEvent{EventID: "e2", Type: EventTypeAspect, Intensity: 0.8, TimeLevel: TimeLevelHourly}
	
	hour1 := &TimeSlot{Events: []AstroEvent{event1, event2}}
	
	// 模拟 AggregateDay (包含 24 个 hourSlots)
	hourSlots := make([]*TimeSlot, 24)
	for i := 0; i < 24; i++ {
		hourSlots[i] = hour1
	}

	merged := agg.mergeAndDeduplicateEvents(hourSlots)

	// 在 Day 粒度下，hourly 事件 (e2) 应该被过滤掉，只剩下 e1
	if len(merged) != 1 {
		t.Errorf("Expected 1 filtered event, got %d", len(merged))
	}

	if merged[0].EventID != "e1" {
		t.Errorf("Expected event e1, got %s", merged[0].EventID)
	}
}
