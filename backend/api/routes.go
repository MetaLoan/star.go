package api

import (
	"star/middleware"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回路由器
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Metrics 中间件 - 收集所有请求的监控数据
	router.Use(middleware.MetricsMiddleware())

	// CORS 配置 - 支持多客户端
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// 健康检查
	router.GET("/health", HealthCheck)

	// API 路由组
	api := router.Group("/api")
	{
		// 计算接口
		calc := api.Group("/calc")
		{
			calc.POST("/chart", CalculateChart)
			calc.POST("/daily", CalculateDailyForecast)
			calc.POST("/weekly", CalculateWeeklyForecast)
			calc.POST("/life-trend", CalculateLifeTrend)
			calc.POST("/time-series", CalculateTimeSeries)
			calc.POST("/profection", CalculateProfection)
			calc.POST("/profection-map", CalculateProfectionMap)
			calc.POST("/transits", CalculateTransits)
			calc.POST("/progressions", CalculateProgressions)
			calc.POST("/void-of-course", CalculateVoidOfCourse)
			calc.POST("/planetary-hour", CalculatePlanetaryHour)

			// 每日星象事件（精确版）⭐ 新增
			calc.POST("/daily-events", CalculateDailyEvents)      // 精确到分钟的每日星象
			calc.GET("/daily-events/simple", GetDailyEventsSimple) // 简化版（无需出生信息）

			// 分值组成查询（详细因子分解）
			calc.POST("/score-breakdown", GetScoreBreakdown)                // 单粒度（开发调试用）
			calc.POST("/score-breakdown-all", GetMultiGranularityBreakdown) // 多粒度（开发调试用）
			calc.POST("/active-factors", GetActiveFactorsInRange)           // 时间范围内活跃因子

			// C端用户友好接口
			calc.POST("/score-explain", GetScoreExplanation) // 分数解释（面向用户）

			// 全因子数据接口
			calc.POST("/total-factors", GetTotalFactors) // 全因子数据（按粒度过滤，包含出相时间）
		}

		// 高级因子查询接口
		factors := api.Group("/factors")
		{
			factors.GET("/types", GetFactorTypes)              // 获取所有支持的因子类型
			factors.POST("/all", GetAdvancedFactors)           // 获取所有高级因子
			factors.POST("/eclipse", GetEclipseFactors)        // 日月食因子
			factors.POST("/lunar-node", GetLunarNodeFactors)   // 月交点因子
			factors.POST("/combustion", GetCombustionFactors)  // 燃烧因子
			factors.POST("/station", GetStationFactors)        // 停滞因子
			factors.POST("/reception", GetReceptionFactors)    // 互容因子
			factors.POST("/fixed-star", GetFixedStarFactors)   // 恒星因子
			factors.POST("/arabic-part", GetArabicPartFactors) // 阿拉伯点因子
			factors.POST("/term-decan", GetTermDecanFactors)   // 界限和十度面因子
			factors.POST("/solar-arc", GetSolarArcFactors)     // 太阳弧推进因子
		}

		// 用户管理
		users := api.Group("/users")
		{
			users.GET("", GetUsers)
			users.POST("", CreateUser)
			users.GET("/:id", GetUser)
			users.PUT("/:id", UpdateUser)
			users.DELETE("/:id", DeleteUser)
			users.GET("/:id/forecast", GetUserForecast)
			users.GET("/:id/snapshot", GetUserSnapshot)
		}

		// 智能体接口
		agent := api.Group("/agent")
		{
			agent.GET("/context", GetAgentContext)
			agent.POST("/query", AgentQuery)
		}

		// 运营配置接口
		admin := api.Group("/admin")
		{
			// 因子权重配置
			admin.GET("/factor-weights", GetFactorWeights)
			admin.PUT("/factor-weights", UpdateFactorWeights)

			// 维度权重配置
			admin.GET("/dimension-weights", GetDimensionWeights)
			admin.PUT("/dimension-weights", UpdateDimensionWeights)

			// 抖动配置
			admin.GET("/jitter-config", GetJitterConfig)
			admin.PUT("/jitter-config", UpdateJitterConfig)

			// 自定义因子管理
			admin.POST("/custom-factors", AddCustomFactor)
			admin.GET("/custom-factors/:userId", GetCustomFactors)
			admin.DELETE("/custom-factors/:userId", ClearCustomFactors)
		}

		// 监控接口 ⭐ 新增
		monitor := api.Group("/monitor")
		{
			monitor.GET("/dashboard", GetMonitorDashboard)       // 监控仪表板页面
			monitor.GET("/summary", GetMonitorSummary)           // 概览统计
			monitor.GET("/stats", GetMonitorStats)               // API统计
			monitor.GET("/recent", GetMonitorRecentRequests)     // 最近请求
			monitor.GET("/realtime", GetMonitorRealtime)         // 实时统计
			monitor.POST("/reset", ResetMonitorStats)            // 重置统计（管理员）
		}
	}

	return router
}
