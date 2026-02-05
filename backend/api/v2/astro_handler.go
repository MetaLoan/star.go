package v2

import (
	"net/http"
	"star/astro"
	"star/cache"
	"star/core"
	"star/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== 统一 API 接口 ====================
// POST /api/v2/astro
// 单一接口返回所有数据：分数、事件、影响、指导

// AstroRequest 请求参数
type AstroRequest struct {
	Birth       models.BirthData `json:"birth" binding:"required"`
	Time        string           `json:"time" binding:"required"` // ISO 8601 格式
	Granularity string           `json:"granularity"`             // hour/day/week/month/year，默认 day
	Language    string           `json:"language"`                // zh/en/ru，默认 en
}

// AstroResponse 响应结构
type AstroResponse struct {
	Slot *core.TimeSlot `json:"slot"`
	Meta *ResponseMeta  `json:"meta"`
}

// ResponseMeta 响应元数据
type ResponseMeta struct {
	Cached      bool   `json:"cached"`
	CacheAge    string `json:"cacheAge,omitempty"`
	ComputeTime string `json:"computeTime"`
	EventCount  int    `json:"eventCount"`
}

// HandleAstro 处理统一 API 请求
func HandleAstro(c *gin.Context) {
	startTime := time.Now()

	// 解析请求
	var req AstroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Granularity == "" {
		req.Granularity = core.GranularityDay
	}
	if req.Language == "" {
		req.Language = "en"
	}

	// 解析查询时间
	queryTime, err := parseTime(req.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid time format",
			"message": err.Error(),
		})
		return
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.Birth)
	if chart == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate natal chart",
		})
		return
	}

	// 生成缓存 key
	userID := generateUserID(req.Birth)
	cacheKey := cache.GenerateCacheKey(userID, req.Granularity, queryTime)

	// 尝试从缓存获取
	globalCache := cache.GetGlobalCache()
	if slot, ok := globalCache.Get(cacheKey); ok {
		// 缓存命中
		computeTime := time.Since(startTime)
		c.JSON(http.StatusOK, AstroResponse{
			Slot: slot,
			Meta: &ResponseMeta{
				Cached:      true,
				ComputeTime: computeTime.String(),
				EventCount:  len(slot.Events),
			},
		})
		return
	}

	// 缓存未命中，计算数据
	slot := calculateSlot(chart, queryTime, req.Granularity, req.Language)

	// 写入缓存
	ttl := cache.DefaultTTL(req.Granularity)
	globalCache.Set(cacheKey, slot, ttl)

	// 返回响应
	computeTime := time.Since(startTime)
	c.JSON(http.StatusOK, AstroResponse{
		Slot: slot,
		Meta: &ResponseMeta{
			Cached:      false,
			ComputeTime: computeTime.String(),
			EventCount:  len(slot.Events),
		},
	})
}

// calculateSlot 计算时间槽
func calculateSlot(chart *models.NatalChart, queryTime time.Time, granularity, language string) *core.TimeSlot {
	calculator := core.NewCalculator(chart, language)
	aggregator := core.NewAggregator(calculator)

	var slot *core.TimeSlot

	switch granularity {
	case core.GranularityHour:
		slot = calculator.CalculateHour(queryTime)
		// 为小时粒度也计算 impactDelta
		deltaCalc := core.NewDeltaCalculator(calculator)
		deltaCalc.ApplyDeltaToSlot(slot, granularity, queryTime)
	case core.GranularityDay:
		slot = aggregator.AggregateDay(queryTime)
	case core.GranularityWeek:
		slot = aggregator.AggregateWeek(queryTime)
	case core.GranularityMonth:
		slot = aggregator.AggregateMonth(queryTime)
	case core.GranularityYear:
		slot = aggregator.AggregateYear(queryTime)
	default:
		slot = aggregator.AggregateDay(queryTime)
	}

	return slot
}

// parseTime 解析时间字符串
func parseTime(timeStr string) (time.Time, error) {
	// 尝试多种格式
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, &TimeParseError{Input: timeStr}
}

// TimeParseError 时间解析错误
type TimeParseError struct {
	Input string
}

func (e *TimeParseError) Error() string {
	return "unable to parse time: " + e.Input
}

// generateUserID 生成用户 ID
func generateUserID(birth models.BirthData) string {
	return formatInt(birth.Year) + formatInt(birth.Month) + formatInt(birth.Day) + "_" +
		formatInt(birth.Hour) + formatInt(birth.Minute) + "_" +
		formatFloat(birth.Latitude) + "_" + formatFloat(birth.Longitude)
}

// formatInt 格式化整数
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

// formatFloat 格式化浮点数
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

// ==================== 路由注册 ====================

// RegisterRoutes 注册 v2 路由
func RegisterRoutes(router *gin.RouterGroup) {
	v2 := router.Group("/v2")
	{
		v2.POST("/astro", HandleAstro)
	}
}
