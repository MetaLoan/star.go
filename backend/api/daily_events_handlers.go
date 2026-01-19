package api

import (
	"fmt"
	"net/http"
	"star/astro"
	"star/models"
	"time"

	"github.com/gin-gonic/gin"
)

// DailyEventsRequest 每日星象事件请求
type DailyEventsRequest struct {
	BirthData          models.BirthData `json:"birthData" binding:"required"`
	Date               string           `json:"date"`               // 格式：2026-01-20 或 RFC3339
	Timezone           int              `json:"timezone"`           // 时区偏移（小时）
	IncludeMinorAspects bool            `json:"includeMinorAspects"` // 是否包含次要相位
}

// DailyEventsResponse 每日星象事件响应
type DailyEventsResponse struct {
	Date         string              `json:"date"`
	Timezone     int                 `json:"timezone"`
	Events       []astro.DailyEvent  `json:"events"`
	EventCount   int                 `json:"eventCount"`
	MajorEvents  []astro.DailyEvent  `json:"majorEvents"`  // 高强度事件
	Summary      string              `json:"summary"`
	DayTheme     string              `json:"dayTheme"`
}

// CalculateDailyEvents 计算每日星象事件
func CalculateDailyEvents(c *gin.Context) {
	var req DailyEventsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	// 解析日期
	var targetDate time.Time
	var err error

	if req.Date == "" {
		// 如果没有提供日期，使用今天
		targetDate = time.Now()
	} else {
		// 尝试多种日期格式
		formats := []string{
			"2006-01-02",
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
		}

		parsed := false
		for _, format := range formats {
			targetDate, err = time.Parse(format, req.Date)
			if err == nil {
				parsed = true
				break
			}
		}

		if !parsed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, please use 2026-01-20 or RFC3339 format"})
			return
		}
	}

	// 应用时区
	if req.Timezone != 0 {
		location := time.FixedZone("Custom", req.Timezone*3600)
		targetDate = time.Date(
			targetDate.Year(), targetDate.Month(), targetDate.Day(),
			0, 0, 0, 0, location,
		)
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 计算每日事件（精确模式）
	events := astro.CalculateDailyEvents(chart, targetDate, req.IncludeMinorAspects)

	// 筛选主要事件（高强度）
	majorEvents := []astro.DailyEvent{}
	for _, event := range events {
		if event.Intensity == "high" {
			majorEvents = append(majorEvents, event)
		}
	}

	// 生成每日主题和总结
	dayTheme := generateDayTheme(events)
	summary := generateDaySummary(events, majorEvents)

	response := DailyEventsResponse{
		Date:        targetDate.Format("2006-01-02"),
		Timezone:    req.Timezone,
		Events:      events,
		EventCount:  len(events),
		MajorEvents: majorEvents,
		Summary:     summary,
		DayTheme:    dayTheme,
	}

	c.JSON(http.StatusOK, response)
}

// generateDayTheme 生成每日主题
func generateDayTheme(events []astro.DailyEvent) string {
	if len(events) == 0 {
		return "A peaceful and harmonious day"
	}

	// Count positive and negative events
	positiveCount := 0
	negativeCount := 0
	highIntensityCount := 0

	for _, event := range events {
		if event.IsPositive {
			positiveCount++
		} else {
			negativeCount++
		}
		if event.Intensity == "high" {
			highIntensityCount++
		}
	}

	// Generate theme
	if highIntensityCount > 2 {
		return "A day of intense energy and importance"
	}
	if positiveCount > negativeCount*2 {
		return "Full of opportunities and harmony"
	}
	if negativeCount > positiveCount*2 {
		return "Need to handle challenges carefully"
	}
	return "Balanced development, opportunities and challenges coexist"
}

// generateDaySummary 生成每日总结
func generateDaySummary(events []astro.DailyEvent, majorEvents []astro.DailyEvent) string {
	if len(events) == 0 {
		return "No significant astrological events today, a peaceful day."
	}

	summary := ""
	
	if len(majorEvents) > 0 {
		summary += fmt.Sprintf("Today has %d major astrological events. ", len(majorEvents))
	}

	// Count event types
	aspectCount := 0
	signChangeCount := 0
	lunarPhaseCount := 0

	for _, event := range events {
		switch event.Type {
		case "aspect":
			aspectCount++
		case "sign_change":
			signChangeCount++
		case "lunar_phase":
			lunarPhaseCount++
		}
	}

	if signChangeCount > 0 {
		summary += "Planetary sign changes indicate energy shifts. "
	}
	if lunarPhaseCount > 0 {
		summary += "Lunar phase changes bring emotional fluctuations. "
	}
	if aspectCount > 3 {
		summary += "Active aspects indicate frequent interactions. "
	}

	if summary == "" {
		summary = "Normal astrological activity, suitable for following plans."
	}

	return summary
}

// GetDailyEventsSimple 简化版每日事件（只需要日期，不需要出生信息）
// 适用于查看当天的普遍星象
func GetDailyEventsSimple(c *gin.Context) {
	dateStr := c.Query("date")
	timezoneStr := c.DefaultQuery("timezone", "8") // 默认东八区

	var targetDate time.Time
	var err error

	if dateStr == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	}

	// 解析时区
	timezone := 8
	if timezoneStr != "" {
		if tz, err := time.ParseDuration(timezoneStr + "h"); err == nil {
			timezone = int(tz.Hours())
		}
	}

	// 使用默认出生信息（只看行运）
	defaultBirthData := models.BirthData{
		Year:      2000,
		Month:     1,
		Day:       1,
		Hour:      12,
		Minute:    0,
		Latitude:  39.9042,
		Longitude: 116.4074,
		Timezone:  float64(timezone),
	}

	chart := astro.CalculateNatalChart(defaultBirthData)
	events := astro.CalculateDailyEvents(chart, targetDate, false)

	// 只返回非相位事件（星座变化、月相等普遍事件）
	universalEvents := []astro.DailyEvent{}
	for _, event := range events {
		if event.Type != "aspect" {
			universalEvents = append(universalEvents, event)
		}
	}

	response := gin.H{
		"date":       targetDate.Format("2006-01-02"),
		"timezone":   timezone,
		"events":     universalEvents,
		"eventCount": len(universalEvents),
	}

	c.JSON(http.StatusOK, response)
}
