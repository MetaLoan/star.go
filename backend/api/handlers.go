package api

import (
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
		"status":  "ok",
		"service": "Star API (Go)",
		"version": "1.0.0",
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
		WithFactors bool             `json:"withFactors"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析日期
	date := time.Now()
	if req.Date != "" {
		if parsed, err := time.Parse("2006-01-02", req.Date); err == nil {
			date = parsed
		}
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	forecast := astro.CalculateDailyForecast(chart, date, req.WithFactors)
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
		Granularity string           `json:"granularity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart := astro.CalculateNatalChart(req.BirthData)
	series := astro.CalculateTimeSeries(chart, req.Start, req.End, req.Granularity)
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
			"error": "维度权重之和必须为1.0",
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

