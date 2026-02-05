package api

import (
	"star/api/v2"
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
		// ==================== V2 统一接口 ====================
		// POST /api/v2/astro - 单一接口返回所有数据（五维运势各粒度曲线趋势）
		v2.RegisterRoutes(api)

		// ==================== 基础查询接口 ====================
		calc := api.Group("/calc")
		{
			// 基础星盘数据查询
			calc.POST("/chart", CalculateChart)
		}

		// ==================== 监控接口 ====================
		monitor := api.Group("/monitor")
		{
			monitor.GET("/dashboard", GetMonitorDashboard)   // 监控仪表板页面
			monitor.GET("/summary", GetMonitorSummary)       // 概览统计
			monitor.GET("/stats", GetMonitorStats)           // API统计
			monitor.GET("/recent", GetMonitorRecentRequests) // 最近请求
			monitor.GET("/realtime", GetMonitorRealtime)     // 实时统计
			monitor.POST("/reset", ResetMonitorStats)        // 重置统计（管理员）
		}
	}

	return router
}
