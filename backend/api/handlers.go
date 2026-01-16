package api

import (
	"fmt"
	"net/http"
	"star/astro"
	"star/models"
	"star/services"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"service":    "Star API (Go)",
		"version":    "1.0.0",
		"dataSource": astro.GetDataSource(), // 显示天文数据唯一来源
		"features": []string{
			"natal-chart",
			"daily-forecast",
			"weekly-forecast",
			"life-trend",
			"profections",
			"transits",
			"progressions",
			"influence-factors",
			"user-management",
			"agent-api",
		},
	})
}

// CalculateChart 计算本命盘
func CalculateChart(c *gin.Context) {
	var req models.BirthData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(req)
	c.JSON(http.StatusOK, chart)
}

// CalculateDailyForecast 计算每日预测
func CalculateDailyForecast(c *gin.Context) {
	var req struct {
		BirthData   models.BirthData `json:"birthData"`
		Date        string           `json:"date"`
		TargetDate  string           `json:"targetDate"`
		WithFactors bool             `json:"withFactors"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析日期 - 优先使用更精确的 targetDate
	// 优先级：targetDate (ISO 8601 带时间) > date (可能是纯日期)
	date := time.Now()

	// 首先尝试解析 targetDate（更精确的时间）
	if req.TargetDate != "" {
		if parsed, err := time.Parse(time.RFC3339, req.TargetDate); err == nil {
			date = parsed
		} else if parsed, err := time.Parse("2006-01-02T15:04:05-07:00", req.TargetDate); err == nil {
			date = parsed
		} else if parsed, err := time.Parse("2006-01-02T15:04:05", req.TargetDate); err == nil {
			date = parsed
		} else if parsed, err := time.Parse("2006-01-02", req.TargetDate); err == nil {
			date = parsed
		}
	} else if req.Date != "" {
		// 如果没有 targetDate，则使用 date
		if parsed, err := time.Parse(time.RFC3339, req.Date); err == nil {
			date = parsed
		} else if parsed, err := time.Parse("2006-01-02T15:04:05-07:00", req.Date); err == nil {
			date = parsed
		} else if parsed, err := time.Parse("2006-01-02", req.Date); err == nil {
			date = parsed
		}
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	// 默认启用 factors 计算，以确保与时间序列 API 一致
	forecast := astro.CalculateDailyForecast(chart, date, true)
	c.JSON(http.StatusOK, forecast)
}

// CalculateWeeklyForecast 计算每周预测
func CalculateWeeklyForecast(c *gin.Context) {
	var req struct {
		BirthData   models.BirthData `json:"birthData"`
		Date        string           `json:"date"`
		WithFactors bool             `json:"withFactors"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date := time.Now()
	if req.Date != "" {
		if parsed, err := time.Parse("2006-01-02", req.Date); err == nil {
			date = parsed
		}
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	forecast := astro.CalculateWeeklyForecast(chart, date, req.WithFactors)
	c.JSON(http.StatusOK, forecast)
}

// CalculateLifeTrend 计算人生趋势
func CalculateLifeTrend(c *gin.Context) {
	var req struct {
		BirthData  models.BirthData `json:"birthData"`
		StartYear  int              `json:"startYear"`
		EndYear    int              `json:"endYear"`
		Resolution string           `json:"resolution"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	trend := astro.CalculateLifeTrend(chart, req.StartYear, req.EndYear, req.Resolution)
	c.JSON(http.StatusOK, trend)
}

// CalculateTimeSeries 生成统一时间序列
func CalculateTimeSeries(c *gin.Context) {
	var req struct {
		BirthData   models.BirthData `json:"birthData"`
		Start       string           `json:"start"`
		End         string           `json:"end"`
		StartTime   string           `json:"startTime"`   // 兼容前端字段名
		EndTime     string           `json:"endTime"`     // 兼容前端字段名
		Granularity string           `json:"granularity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 兼容两种字段名（优先使用有值的字段）
	start := req.Start
	if req.StartTime != "" {
		start = req.StartTime
	}
	end := req.End
	if req.EndTime != "" {
		end = req.EndTime
	}

	// 调试日志 - 总是打印
	fmt.Printf("DEBUG Handler: start='%s', end='%s', startTime='%s', endTime='%s', final_start='%s', final_end='%s'\n", 
		req.Start, req.End, req.StartTime, req.EndTime, start, end)

	chart := astro.CalculateNatalChart(req.BirthData)
	series := astro.CalculateTimeSeries(chart, start, end, req.Granularity)
	c.JSON(http.StatusOK, series)
}

// CalculateProfection 计算年限法
func CalculateProfection(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		Age       int              `json:"age"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	profection := astro.CalculateAnnualProfection(chart, req.Age)
	c.JSON(http.StatusOK, profection)
}

// CalculateProfectionMap 计算完整年限法地图
func CalculateProfectionMap(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	profectionMap := astro.CalculateProfectionMap(chart)
	c.JSON(http.StatusOK, profectionMap)
}

// CalculateTransits 计算行运
func CalculateTransits(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		StartDate string           `json:"startDate"`
		EndDate   string           `json:"endDate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	transits := astro.CalculateTransits(chart, req.StartDate, req.EndDate)
	c.JSON(http.StatusOK, transits)
}

// CalculateProgressions 计算推运
func CalculateProgressions(c *gin.Context) {
	var req struct {
		BirthData  models.BirthData `json:"birthData"`
		TargetDate string           `json:"targetDate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	progressions := astro.CalculateProgressions(chart, req.TargetDate)
	c.JSON(http.StatusOK, progressions)
}

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {
	users := services.GetAllUsers()
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req struct {
		Name      string           `json:"name"`
		BirthData models.BirthData `json:"birthData"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := services.CreateUser(req.Name, req.BirthData)
	c.JSON(http.StatusCreated, user)
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := services.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name      string           `json:"name"`
		BirthData models.BirthData `json:"birthData"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UpdateUser(id, req.Name, req.BirthData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetUserForecast 获取用户预测
func GetUserForecast(c *gin.Context) {
	id := c.Param("id")
	forecastType := c.DefaultQuery("type", "daily")

	user, err := services.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(user.BirthData)
	var forecast interface{}
	switch forecastType {
	case "weekly":
		forecast = astro.CalculateWeeklyForecast(chart, time.Now(), true)
	default:
		forecast = astro.CalculateDailyForecast(chart, time.Now(), true)
	}

	c.JSON(http.StatusOK, forecast)
}

// GetUserSnapshot 获取用户当前状态快照
func GetUserSnapshot(c *gin.Context) {
	id := c.Param("id")

	user, err := services.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	snapshot := services.GetUserSnapshot(user)
	c.JSON(http.StatusOK, snapshot)
}

// GetAgentContext 获取智能体上下文
func GetAgentContext(c *gin.Context) {
	context := services.GetAgentContext()
	c.JSON(http.StatusOK, context)
}

// AgentQuery 智能体查询
func AgentQuery(c *gin.Context) {
	var req struct {
		UserID string `json:"userId"`
		Query  string `json:"query"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := services.HandleAgentQuery(req.UserID, req.Query)
	c.JSON(http.StatusOK, response)
}

// CalculateVoidOfCourse 计算月亮空亡
func CalculateVoidOfCourse(c *gin.Context) {
	var req struct {
		Date      string  `json:"date"`      // 可选，默认当前时间
		Latitude  float64 `json:"latitude"`  // 可选
		Longitude float64 `json:"longitude"` // 可选
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析日期
	var date time.Time
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05", req.Date)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", req.Date)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日期格式"})
				return
			}
		}
		date = parsed
	} else {
		date = time.Now()
	}

	jd := astro.DateToJulianDay(date)
	vocInfo := astro.CalculateVoidOfCourse(jd, nil)

	c.JSON(http.StatusOK, vocInfo)
}

// CalculatePlanetaryHour 计算行星时
func CalculatePlanetaryHour(c *gin.Context) {
	var req struct {
		Date      string  `json:"date"`      // 可选，默认当前时间
		Latitude  float64 `json:"latitude"`  // 可选，默认 0
		Longitude float64 `json:"longitude"` // 可选，默认 0
		FullDay   bool    `json:"fullDay"`   // 是否返回全天行星时
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析日期
	var date time.Time
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05", req.Date)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", req.Date)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日期格式"})
				return
			}
		}
		date = parsed
	} else {
		date = time.Now()
	}

	if req.FullDay {
		// 返回全天行星时
		hours := astro.GetPlanetaryHoursForDate(date, req.Latitude, req.Longitude)
		c.JSON(http.StatusOK, gin.H{
			"date":  date.Format("2006-01-02"),
			"hours": hours,
		})
	} else {
		// 返回当前行星时
		hourInfo := astro.CalculatePlanetaryHourEnhanced(date, req.Latitude, req.Longitude)
		c.JSON(http.StatusOK, hourInfo)
	}
}

// ==================== 运营配置 API ====================

// GetFactorWeights 获取因子权重配置
func GetFactorWeights(c *gin.Context) {
	weights := astro.GetCurrentFactorWeights()
	c.JSON(http.StatusOK, gin.H{
		"weights": weights,
		"description": gin.H{
			"dignity":        "尊贵度因子权重（入庙/旺相/落陷/失势）",
			"retrograde":     "逆行因子权重",
			"aspectPhase":    "相位阶段因子权重（入相/离相）",
			"aspectOrb":      "相位容许度因子权重",
			"outerPlanet":    "外行星因子权重",
			"profectionLord": "年主星因子权重",
			"lunarPhase":     "月相因子权重",
			"planetaryHour":  "行星时因子权重",
			"voidOfCourse":   "月亮空亡因子权重",
			"personal":       "个人因子权重",
			"custom":         "自定义因子权重",
		},
	})
}

// UpdateFactorWeights 更新因子权重配置
func UpdateFactorWeights(c *gin.Context) {
	var req models.FactorWeights
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	astro.UpdateFactorWeights(req)
	c.JSON(http.StatusOK, gin.H{
		"message": "因子权重已更新",
		"weights": req,
	})
}

// GetDimensionWeights 获取维度权重配置
func GetDimensionWeights(c *gin.Context) {
	weights := astro.GetCurrentDimensionWeights()
	c.JSON(http.StatusOK, gin.H{
		"weights": weights,
		"description": gin.H{
			"career":       "事业维度权重（默认0.25，对应10/6/1宫）",
			"relationship": "关系维度权重（默认0.20，对应7/5/11宫）",
			"health":       "健康维度权重（默认0.20，对应1/6/8宫）",
			"finance":      "财务维度权重（默认0.20，对应2/8宫）",
			"spiritual":    "灵性维度权重（默认0.15，对应9/12宫）",
		},
		"note": "所有权重之和应为1.0",
	})
}

// UpdateDimensionWeights 更新维度权重配置
func UpdateDimensionWeights(c *gin.Context) {
	var req models.DimensionWeights
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证权重总和
	total := req.Career + req.Relationship + req.Health + req.Finance + req.Spiritual
	if total < 0.99 || total > 1.01 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "维度权重之和必须为1.0",
			"current_total": total,
		})
		return
	}

	astro.UpdateDimensionWeights(req)
	c.JSON(http.StatusOK, gin.H{
		"message": "维度权重已更新",
		"weights": req,
	})
}

// GetJitterConfig 获取抖动配置
func GetJitterConfig(c *gin.Context) {
	config := astro.DefaultJitterConfig
	c.JSON(http.StatusOK, gin.H{
		"config": config,
		"description": gin.H{
			"enabled":   "是否启用视觉抖动（仅影响显示，不影响计算）",
			"magnitude": "抖动幅度（±范围，默认0.5）",
			"seed":      "随机种子（0表示使用时间戳）",
		},
	})
}

// UpdateJitterConfig 更新抖动配置
func UpdateJitterConfig(c *gin.Context) {
	var req astro.JitterConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	astro.UpdateJitterConfig(req)
	c.JSON(http.StatusOK, gin.H{
		"message": "抖动配置已更新",
		"config":  req,
	})
}

// ==================== 自定义因子 API ====================

// AddCustomFactor 添加自定义因子
func AddCustomFactor(c *gin.Context) {
	var req struct {
		UserID     string `json:"userId"`
		Definition string `json:"definition"` // 格式：AddScore=(2*healthScore,2.5,202501171230)
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID == "" {
		req.UserID = "default"
	}

	factor, err := astro.AddCustomFactor(req.UserID, req.Definition)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"example": "AddScore=(2*healthScore,2.5,202501171230)",
			"format": gin.H{
				"operation":  "AddScore | SubScore | MulScore | SetScore",
				"value":      "数值（可选乘以维度如 2*healthScore）",
				"duration":   "持续小时数",
				"startTime":  "开始时间 YYYYMMDDHHmm",
				"dimensions": "career | relationship | health | finance | spiritual | overall",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "自定义因子已添加",
		"factor":  factor,
	})
}

// GetCustomFactors 获取用户的自定义因子
func GetCustomFactors(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = "default"
	}

	factors := astro.GetAllCustomFactors(userID)
	c.JSON(http.StatusOK, gin.H{
		"userId":  userID,
		"count":   len(factors),
		"factors": factors,
	})
}

// ClearCustomFactors 清除用户的自定义因子
func ClearCustomFactors(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = "default"
	}

	astro.ClearCustomFactors(userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "自定义因子已清除",
		"userId":  userID,
	})
}

// ==================== 分值组成查询 API ====================

// GetScoreBreakdown 获取指定时间点的分值组成详情
// POST /api/calc/score-breakdown
// 请求体：
//
//	{
//	  "birthData": {...},
//	  "queryTime": "2026-01-05T12:00:00+08:00",
//	  "granularity": "hour" | "day" | "week" | "month" | "year",
//	  "userId": "optional"
//	}
func GetScoreBreakdown(c *gin.Context) {
	var req struct {
		BirthData   models.BirthData `json:"birthData"`
		QueryTime   string           `json:"queryTime"`
		Granularity string           `json:"granularity"`
		UserID      string           `json:"userId,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证粒度
	validGranularities := map[string]bool{
		"hour": true, "day": true, "week": true, "month": true, "year": true,
	}
	if !validGranularities[req.Granularity] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的粒度参数",
			"valid": []string{"hour", "day", "week", "month", "year"},
		})
		return
	}

	// 解析时间
	queryTime, err := time.Parse(time.RFC3339, req.QueryTime)
	if err != nil {
		// 尝试其他格式
		queryTime, err = time.Parse("2006-01-02T15:04:05", req.QueryTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "无效的时间格式",
				"example": "2026-01-05T12:00:00+08:00",
			})
			return
		}
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 计算分值组成
	breakdown := astro.CalculateScoreBreakdown(chart, queryTime, req.Granularity, req.UserID)

	c.JSON(http.StatusOK, breakdown)
}

// GetMultiGranularityBreakdown 获取多粒度分值组成
// POST /api/calc/score-breakdown-all
// 一次返回 hour/day/month/year 四个粒度的分值组成
func GetMultiGranularityBreakdown(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
		UserID    string           `json:"userId,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析时间
	queryTime, err := time.Parse(time.RFC3339, req.QueryTime)
	if err != nil {
		queryTime, err = time.Parse("2006-01-02T15:04:05", req.QueryTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "无效的时间格式",
				"example": "2026-01-05T12:00:00+08:00",
			})
			return
		}
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 获取多粒度分值组成
	result := astro.GetMultiGranularityBreakdown(chart, queryTime, req.UserID)

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"breakdown": result,
	})
}

// GetActiveFactorsInRange 获取时间范围内活跃的所有因子
// POST /api/calc/active-factors
// 查询指定时间范围内（年/月/周/日）活跃的所有影响因子
// infect 参数：
//   - "all": 返回所有因子（不过滤）
//   - "core": 按可见性规则过滤（年度级在年/月/周/日/小时可见，月级在月/周/日/小时可见...）
func GetActiveFactorsInRange(c *gin.Context) {
	var req struct {
		BirthData   models.BirthData `json:"birthData"`
		QueryTime   string           `json:"queryTime"`   // 范围内任意时间点
		Granularity string           `json:"granularity"` // year/month/week/day
		Infect      string           `json:"infect"`      // all/core，默认 all
		UserID      string           `json:"userId,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证粒度参数
	validGranularities := map[string]bool{"year": true, "month": true, "week": true, "day": true, "hour": true}
	if !validGranularities[req.Granularity] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的粒度参数",
			"valid":   []string{"year", "month", "week", "day", "hour"},
			"example": "week",
		})
		return
	}

	// 验证 infect 参数，默认为 all
	if req.Infect == "" {
		req.Infect = "all"
	}
	validInfect := map[string]bool{"all": true, "core": true}
	if !validInfect[req.Infect] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的 infect 参数",
			"valid":   []string{"all", "core"},
			"example": "core",
			"description": map[string]string{
				"all":  "返回所有因子（不过滤）",
				"core": "按可见性规则过滤（年度级在年/月/周/日/小时可见，月级在月/周/日/小时可见...）",
			},
		})
		return
	}

	// 解析时间
	queryTime, err := time.Parse(time.RFC3339, req.QueryTime)
	if err != nil {
		queryTime, err = time.Parse("2006-01-02T15:04:05", req.QueryTime)
		if err != nil {
			queryTime, err = time.Parse("2006-01-02", req.QueryTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "无效的时间格式",
					"example": "2026-01-05 或 2026-01-05T12:00:00+08:00",
				})
				return
			}
		}
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 获取时间范围内活跃的因子
	result := astro.GetActiveFactorsInRange(chart, queryTime, req.Granularity, req.Infect, req.UserID)

	c.JSON(http.StatusOK, result)
}

// GetTotalFactors 获取全因子数据
// POST /api/calc/total-factors
// 返回指定时间点的所有影响因子，按粒度过滤，包含正负影响和出相时间
func GetTotalFactors(c *gin.Context) {
	var req struct {
		BirthData   models.BirthData `json:"birthData"`
		QueryTime   string           `json:"queryTime"`   // ISO 8601 格式
		Granularity string           `json:"granularity"` // hour, day, week, month, year
		UserID      string           `json:"userId,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证粒度参数
	validGranularities := map[string]bool{
		"hour": true, "day": true, "week": true, "month": true, "year": true,
	}
	if req.Granularity == "" {
		req.Granularity = "day" // 默认为天级别
	}
	if !validGranularities[req.Granularity] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的粒度参数",
			"valid":   []string{"hour", "day", "week", "month", "year"},
			"example": "day",
		})
		return
	}

	// 解析时间
	queryTime, err := time.Parse(time.RFC3339, req.QueryTime)
	if err != nil {
		queryTime, err = time.Parse("2006-01-02T15:04:05", req.QueryTime)
		if err != nil {
			queryTime, err = time.Parse("2006-01-02", req.QueryTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "无效的时间格式",
					"example": "2026-01-15 或 2026-01-15T12:00:00+08:00",
				})
				return
			}
		}
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 获取全因子数据
	result := astro.GetTotalFactors(chart, queryTime, req.Granularity, req.UserID)

	c.JSON(http.StatusOK, result)
}

// ==================== C端用户友好接口 ====================

// GetScoreExplanation 获取分数解释（面向C端用户）
// POST /api/calc/score-explain
// 用通俗易懂的语言解释分数是如何计算的，受哪些天文现象影响
func GetScoreExplanation(c *gin.Context) {
	var req struct {
		BirthData   models.BirthData `json:"birthData"`
		QueryTime   string           `json:"queryTime"`
		Granularity string           `json:"granularity"` // hour, day, week, month, year
		Dimension   string           `json:"dimension"`   // career, relationship, health, finance, spiritual, overall
		UserID      string           `json:"userId,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 默认值
	if req.Granularity == "" {
		req.Granularity = "hour"
	}
	if req.Dimension == "" {
		req.Dimension = "overall"
	}

	// 验证粒度
	validGranularities := map[string]bool{
		"hour": true, "day": true, "week": true, "month": true, "year": true,
	}
	if !validGranularities[req.Granularity] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的粒度参数",
			"valid": []string{"hour", "day", "week", "month", "year"},
		})
		return
	}

	// 验证维度
	validDimensions := map[string]bool{
		"career": true, "relationship": true, "health": true,
		"finance": true, "spiritual": true, "overall": true,
	}
	if !validDimensions[req.Dimension] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的维度参数",
			"valid": []string{"career", "relationship", "health", "finance", "spiritual", "overall"},
		})
		return
	}

	// 解析时间
	queryTime, err := time.Parse(time.RFC3339, req.QueryTime)
	if err != nil {
		queryTime, err = time.Parse("2006-01-02T15:04:05", req.QueryTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "无效的时间格式",
				"example": "2026-01-05T12:00:00+08:00",
			})
			return
		}
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 获取分数解释
	explanation := astro.GetScoreExplanation(chart, queryTime, req.Granularity, req.Dimension, req.UserID)

	c.JSON(http.StatusOK, explanation)
}

// ==================== 高级因子查询接口 ====================

// GetAdvancedFactors 获取所有高级因子
func GetAdvancedFactors(c *gin.Context) {
	var req struct {
		BirthData  models.BirthData `json:"birthData"`
		QueryTime  string           `json:"queryTime"`
		FactorType string           `json:"factorType"` // 可选：指定查询的因子类型
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析时间
	queryTime := time.Now()
	if req.QueryTime != "" {
		if parsed, err := time.Parse(time.RFC3339, req.QueryTime); err == nil {
			queryTime = parsed
		} else if parsed, err := time.Parse("2006-01-02", req.QueryTime); err == nil {
			queryTime = parsed
		}
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 构建出生时间
	birthTime := time.Date(
		req.BirthData.Year, time.Month(req.BirthData.Month), req.BirthData.Day,
		req.BirthData.Hour, req.BirthData.Minute, req.BirthData.Second,
		0, time.FixedZone("", int(req.BirthData.Timezone*3600)),
	)

	// 计算所有高级因子
	allFactors := astro.CalculateAllAdvancedFactors(chart, birthTime, queryTime)

	// 如果指定了因子类型，进行过滤
	if req.FactorType != "" {
		var filtered []models.InfluenceFactor
		for _, f := range allFactors {
			if string(f.Type) == req.FactorType {
				filtered = append(filtered, f)
			}
		}
		allFactors = filtered
	}

	// 按类型分组
	grouped := make(map[string][]models.InfluenceFactor)
	for _, f := range allFactors {
		key := string(f.Type)
		grouped[key] = append(grouped[key], f)
	}

	c.JSON(http.StatusOK, gin.H{
		"queryTime":  queryTime.Format(time.RFC3339),
		"totalCount": len(allFactors),
		"factors":    allFactors,
		"grouped":    grouped,
		"summary":    getFactorTypeSummary(allFactors),
	})
}

// GetEclipseFactors 获取日月食因子
func GetEclipseFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)
	factors := astro.CalculateEclipseFactors(chart, queryTime)

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"count":     len(factors),
		"factors":   factors,
	})
}

// GetLunarNodeFactors 获取月交点因子
func GetLunarNodeFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)
	factors := astro.CalculateLunarNodeFactors(chart, queryTime)

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"count":     len(factors),
		"factors":   factors,
	})
}

// GetCombustionFactors 获取燃烧因子
func GetCombustionFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)
	factors := astro.CalculateCombustionFactors(chart, queryTime)

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"count":     len(factors),
		"factors":   factors,
	})
}

// GetStationFactors 获取停滞因子
func GetStationFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)
	factors := astro.CalculateStationFactors(chart, queryTime)

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"count":     len(factors),
		"factors":   factors,
	})
}

// GetReceptionFactors 获取互容因子
func GetReceptionFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)
	factors := astro.CalculateReceptionFactors(chart, queryTime)

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"count":     len(factors),
		"factors":   factors,
	})
}

// GetFixedStarFactors 获取恒星因子
func GetFixedStarFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)
	factors := astro.CalculateFixedStarFactors(chart, queryTime)

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"count":     len(factors),
		"factors":   factors,
		"starList":  astro.ImportantFixedStars, // 返回支持的恒星列表
	})
}

// GetArabicPartFactors 获取阿拉伯点因子
func GetArabicPartFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)
	factors := astro.CalculateArabicPartFactors(chart, queryTime)

	// 计算福点和精神点位置
	asc := chart.Ascendant
	sunPos := astro.GetPlanetFromChart(chart, models.Sun)
	moonPos := astro.GetPlanetFromChart(chart, models.Moon)
	var sun, moon float64
	if sunPos != nil {
		sun = sunPos.Longitude
	}
	if moonPos != nil {
		moon = moonPos.Longitude
	}
	var fortunePart, spiritPart float64
	if sun > 180 {
		fortunePart = float64(int(asc+sun-moon+360) % 360)
		spiritPart = float64(int(asc+moon-sun+360) % 360)
	} else {
		fortunePart = float64(int(asc+moon-sun+360) % 360)
		spiritPart = float64(int(asc+sun-moon+360) % 360)
	}

	c.JSON(http.StatusOK, gin.H{
		"queryTime": queryTime.Format(time.RFC3339),
		"count":     len(factors),
		"factors":   factors,
		"parts": gin.H{
			"fortune": gin.H{
				"longitude": fortunePart,
				"sign":      astro.GetSignFromLongitude(fortunePart),
			},
			"spirit": gin.H{
				"longitude": spiritPart,
				"sign":      astro.GetSignFromLongitude(spiritPart),
			},
		},
	})
}

// GetTermDecanFactors 获取界限和十度面因子
func GetTermDecanFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)

	termFactors := astro.CalculateTermFactors(chart, queryTime)
	decanFactors := astro.CalculateDecanFactors(chart, queryTime)

	c.JSON(http.StatusOK, gin.H{
		"queryTime":    queryTime.Format(time.RFC3339),
		"termCount":    len(termFactors),
		"decanCount":   len(decanFactors),
		"termFactors":  termFactors,
		"decanFactors": decanFactors,
	})
}

// GetSolarArcFactors 获取太阳弧推进因子
func GetSolarArcFactors(c *gin.Context) {
	var req struct {
		BirthData models.BirthData `json:"birthData"`
		QueryTime string           `json:"queryTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryTime := parseQueryTime(req.QueryTime)
	chart := astro.CalculateNatalChart(req.BirthData)

	// 构建出生时间
	birthTime := time.Date(
		req.BirthData.Year, time.Month(req.BirthData.Month), req.BirthData.Day,
		req.BirthData.Hour, req.BirthData.Minute, req.BirthData.Second,
		0, time.FixedZone("", int(req.BirthData.Timezone*3600)),
	)

	factors := astro.CalculateSolarArcFactors(chart, birthTime, queryTime)

	// 计算当前太阳弧度数
	age := queryTime.Sub(birthTime).Hours() / 24 / 365.25
	solarArc := age * 0.9856

	c.JSON(http.StatusOK, gin.H{
		"queryTime":      queryTime.Format(time.RFC3339),
		"count":          len(factors),
		"factors":        factors,
		"age":            age,
		"solarArcDegree": solarArc,
	})
}

// GetFactorTypes 获取所有支持的因子类型
func GetFactorTypes(c *gin.Context) {
	types := []gin.H{
		// 基础因子
		{"type": "dignity", "name": "尊贵度", "description": "行星入庙、旺、弱、陷状态"},
		{"type": "retrograde", "name": "逆行", "description": "行星逆行状态"},
		{"type": "aspectPhase", "name": "相位", "description": "行星间主要相位"},
		{"type": "aspectOrb", "name": "相位容许度", "description": "相位精确度"},
		{"type": "outerPlanet", "name": "外行星过境", "description": "外行星过境影响"},
		{"type": "profectionLord", "name": "年主星", "description": "小限法年主星"},
		{"type": "lunarPhase", "name": "月相", "description": "月相周期"},
		{"type": "planetaryHour", "name": "行星时", "description": "当前行星时"},
		{"type": "voidOfCourse", "name": "月空亡", "description": "月亮空亡期"},
		{"type": "personal", "name": "个人因子", "description": "太阳回归等个人技术"},
		{"type": "custom", "name": "自定义", "description": "用户自定义因子"},
		// 高级因子
		{"type": "eclipse", "name": "日月食", "description": "日食、月食影响"},
		{"type": "lunarNode", "name": "月交点", "description": "北交/南交点过境"},
		{"type": "combustion", "name": "燃烧", "description": "行星被太阳灼烧"},
		{"type": "station", "name": "停滞", "description": "行星逆行前后停滞"},
		{"type": "reception", "name": "互容", "description": "行星互容关系"},
		{"type": "fixedStar", "name": "恒星", "description": "重要恒星影响"},
		{"type": "arabicPart", "name": "阿拉伯点", "description": "福点、精神点等"},
		{"type": "midpoint", "name": "中点", "description": "中点技术"},
		{"type": "antiscion", "name": "反生点", "description": "反生点技术"},
		{"type": "term", "name": "界限", "description": "埃及界限"},
		{"type": "decan", "name": "十度面", "description": "迦勒底十度面"},
		{"type": "solarArc", "name": "太阳弧", "description": "太阳弧推进"},
		{"type": "primary", "name": "主限推进", "description": "主限方向"},
		{"type": "firdaria", "name": "法达", "description": "法达时间主星"},
		{"type": "zodiacal", "name": "黄道释放", "description": "黄道释放技术"},
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(types),
		"types": types,
	})
}

// ==================== 辅助函数 ====================

func parseQueryTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Now()
	}
	if parsed, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return parsed
	}
	if parsed, err := time.Parse("2006-01-02", timeStr); err == nil {
		return parsed
	}
	return time.Now()
}

func getFactorTypeSummary(factors []models.InfluenceFactor) gin.H {
	summary := make(gin.H)

	for _, f := range factors {
		key := string(f.Type)
		if _, exists := summary[key]; !exists {
			summary[key] = gin.H{
				"count":           0,
				"totalAdjustment": 0.0,
				"positiveCount":   0,
				"negativeCount":   0,
			}
		}

		s := summary[key].(gin.H)
		s["count"] = s["count"].(int) + 1
		s["totalAdjustment"] = s["totalAdjustment"].(float64) + f.Adjustment
		if f.IsPositive {
			s["positiveCount"] = s["positiveCount"].(int) + 1
		} else {
			s["negativeCount"] = s["negativeCount"].(int) + 1
		}
		summary[key] = s
	}

	return summary
}
