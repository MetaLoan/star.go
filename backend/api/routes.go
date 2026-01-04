package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回路由器
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

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
		}
	}

	return router
}

