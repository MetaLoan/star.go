package v2

import (
	"fmt"
	"net/http"
	"star/astro"
	"star/cache"
	"star/core"
	"star/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== 趋势图 API 接口 ====================
// POST /api/v2/astro/trend
// 一次性返回指定粒度下的所有趋势数据点

// TrendRequest 趋势请求参数
type TrendRequest struct {
	Birth       models.BirthData `json:"birth" binding:"required"`
	StartTime   string           `json:"start_time" binding:"required"` // ISO 8601 格式
	Granularity string           `json:"granularity"`                   // hour/day/week/month/year，默认 day
	Language    string           `json:"language"`                      // zh/en/ru，默认 en
}

// TrendResponse 趋势响应
type TrendResponse struct {
	Granularity string        `json:"granularity"`
	Points      []TrendPoint  `json:"points"`
	Summary     *TrendSummary `json:"summary"`
	Meta        *ResponseMeta `json:"meta"`
}

// TrendPoint 单个趋势点
type TrendPoint struct {
	Time   time.Time            `json:"time"`
	Label  string               `json:"label"`
	Scores core.DimensionScores `json:"scores"`
}

// TrendSummary 趋势汇总
type TrendSummary struct {
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
	Trend string  `json:"trend"` // upward/downward/stable
}

// HandleTrend 处理趋势图请求
func HandleTrend(c *gin.Context) {
	startTime := time.Now()

	// 解析请求
	var req TrendRequest
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
	queryTime, err := parseTime(req.StartTime)
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
	cacheKey := generateTrendCacheKey(userID, req.Granularity, queryTime, req.Language)

	// 尝试从缓存获取
	globalCache := cache.GetGlobalCache()
	if cached, ok := globalCache.GetTrend(cacheKey); ok {
		if cachedResponse, ok := cached.(*TrendResponse); ok {
			computeTime := time.Since(startTime)
			cachedResponse.Meta = &ResponseMeta{
				Cached:      true,
				ComputeTime: computeTime.String(),
			}
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	// 计算趋势数据
	calculator := core.NewCalculator(chart, req.Language)
	points := calculateTrendPoints(calculator, queryTime, req.Granularity, req.Language)

	// 计算汇总
	summary := calculateTrendSummary(points)

	// 构建响应
	response := &TrendResponse{
		Granularity: req.Granularity,
		Points:      points,
		Summary:     summary,
	}

	// 写入缓存
	ttl := cache.DefaultTTL(req.Granularity)
	globalCache.SetTrend(cacheKey, response, ttl)

	// 返回响应
	computeTime := time.Since(startTime)
	response.Meta = &ResponseMeta{
		Cached:      false,
		ComputeTime: computeTime.String(),
	}
	c.JSON(http.StatusOK, response)
}

// generateTrendCacheKey 生成趋势缓存 key
func generateTrendCacheKey(userID, granularity string, t time.Time, language string) string {
	var normalizedTime string
	switch granularity {
	case core.GranularityHour:
		// 按天归一化
		normalizedTime = t.Format("2006-01-02")
	case core.GranularityDay:
		// 按月归一化
		normalizedTime = t.Format("2006-01")
	case core.GranularityWeek:
		// 按月归一化
		normalizedTime = t.Format("2006-01")
	case core.GranularityMonth:
		// 按年归一化
		normalizedTime = t.Format("2006")
	case core.GranularityYear:
		// 按年归一化（中心年份）
		normalizedTime = fmt.Sprintf("%d", t.Year())
	default:
		normalizedTime = t.Format("2006-01-02")
	}
	return fmt.Sprintf("trend:%s:%s:%s:%s", userID, granularity, normalizedTime, language)
}

// calculateTrendPoints 计算趋势数据点
func calculateTrendPoints(calculator *core.Calculator, queryTime time.Time, granularity, language string) []TrendPoint {
	var points []TrendPoint

	switch granularity {
	case core.GranularityHour:
		points = calculateHourTrendPoints(calculator, queryTime, language)
	case core.GranularityDay:
		points = calculateDayTrendPoints(calculator, queryTime, language)
	case core.GranularityWeek:
		points = calculateWeekTrendPoints(calculator, queryTime, language)
	case core.GranularityMonth:
		points = calculateMonthTrendPoints(calculator, queryTime, language)
	case core.GranularityYear:
		points = calculateYearTrendPoints(calculator, queryTime, language)
	default:
		points = calculateDayTrendPoints(calculator, queryTime, language)
	}

	return points
}

// calculateHourTrendPoints 计算小时趋势点（当天 24 个小时）
func calculateHourTrendPoints(calculator *core.Calculator, queryTime time.Time, language string) []TrendPoint {
	dayStart := time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
	points := make([]TrendPoint, 24)

	for hour := 0; hour < 24; hour++ {
		t := dayStart.Add(time.Duration(hour) * time.Hour)
		slot := calculator.CalculateHour(t)
		points[hour] = TrendPoint{
			Time:   t,
			Label:  fmt.Sprintf("%02d:00", hour),
			Scores: slot.Scores,
		}
	}

	return points
}

// calculateDayTrendPoints 计算日趋势点（当月所有天）
func calculateDayTrendPoints(calculator *core.Calculator, queryTime time.Time, language string) []TrendPoint {
	monthStart := time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	daysInMonth := int(monthEnd.Sub(monthStart).Hours() / 24)

	points := make([]TrendPoint, daysInMonth)

	for day := 0; day < daysInMonth; day++ {
		// 每天中午采样
		t := time.Date(queryTime.Year(), queryTime.Month(), day+1, 12, 0, 0, 0, queryTime.Location())
		slot := calculator.CalculateHour(t)
		points[day] = TrendPoint{
			Time:   t,
			Label:  fmt.Sprintf("%d", day+1),
			Scores: slot.Scores,
		}
	}

	return points
}

// calculateWeekTrendPoints 计算周趋势点（当月 4-5 个周）
func calculateWeekTrendPoints(calculator *core.Calculator, queryTime time.Time, language string) []TrendPoint {
	monthStart := time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	var points []TrendPoint
	weekNum := 1

	// 找到第一个周一
	current := monthStart
	for current.Weekday() != time.Monday {
		current = current.AddDate(0, 0, 1)
	}

	// 如果第一个周一在月初之后，先添加月初那周
	if current.After(monthStart) {
		t := time.Date(monthStart.Year(), monthStart.Month(), monthStart.Day(), 12, 0, 0, 0, monthStart.Location())
		slot := calculator.CalculateHour(t)
		points = append(points, TrendPoint{
			Time:   t,
			Label:  fmt.Sprintf("W%d", weekNum),
			Scores: slot.Scores,
		})
		weekNum++
	}

	// 遍历每个周一
	for current.Before(monthEnd) {
		t := time.Date(current.Year(), current.Month(), current.Day(), 12, 0, 0, 0, current.Location())
		slot := calculator.CalculateHour(t)
		points = append(points, TrendPoint{
			Time:   t,
			Label:  fmt.Sprintf("W%d", weekNum),
			Scores: slot.Scores,
		})
		weekNum++
		current = current.AddDate(0, 0, 7)
	}

	return points
}

// calculateMonthTrendPoints 计算月趋势点（当年 12 个月）
func calculateMonthTrendPoints(calculator *core.Calculator, queryTime time.Time, language string) []TrendPoint {
	points := make([]TrendPoint, 12)
	monthLabelsEn := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	monthLabelsZh := []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}

	for month := 0; month < 12; month++ {
		// 每月 15 号中午采样
		t := time.Date(queryTime.Year(), time.Month(month+1), 15, 12, 0, 0, 0, queryTime.Location())
		slot := calculator.CalculateHour(t)

		label := monthLabelsEn[month]
		if language == "zh" {
			label = monthLabelsZh[month]
		}

		points[month] = TrendPoint{
			Time:   t,
			Label:  label,
			Scores: slot.Scores,
		}
	}

	return points
}

// calculateYearTrendPoints 计算年趋势点（前后各 2 年，共 5 个点）
func calculateYearTrendPoints(calculator *core.Calculator, queryTime time.Time, language string) []TrendPoint {
	centerYear := queryTime.Year()
	points := make([]TrendPoint, 5)

	for i := 0; i < 5; i++ {
		year := centerYear - 2 + i
		// 每年 6 月 15 号中午采样
		t := time.Date(year, 6, 15, 12, 0, 0, 0, queryTime.Location())
		slot := calculator.CalculateHour(t)
		points[i] = TrendPoint{
			Time:   t,
			Label:  fmt.Sprintf("%d", year),
			Scores: slot.Scores,
		}
	}

	return points
}

// calculateTrendSummary 计算趋势汇总
func calculateTrendSummary(points []TrendPoint) *TrendSummary {
	if len(points) == 0 {
		return nil
	}

	maxScore := points[0].Scores.Overall
	minScore := points[0].Scores.Overall

	for _, p := range points {
		if p.Scores.Overall > maxScore {
			maxScore = p.Scores.Overall
		}
		if p.Scores.Overall < minScore {
			minScore = p.Scores.Overall
		}
	}

	// 判断趋势
	trend := "stable"
	if len(points) >= 2 {
		firstHalf := 0.0
		secondHalf := 0.0
		mid := len(points) / 2

		for i := 0; i < mid; i++ {
			firstHalf += points[i].Scores.Overall
		}
		for i := mid; i < len(points); i++ {
			secondHalf += points[i].Scores.Overall
		}

		firstHalf /= float64(mid)
		secondHalf /= float64(len(points) - mid)

		diff := secondHalf - firstHalf
		if diff > 3 {
			trend = "upward"
		} else if diff < -3 {
			trend = "downward"
		}
	}

	return &TrendSummary{
		Max:   maxScore,
		Min:   minScore,
		Trend: trend,
	}
}
