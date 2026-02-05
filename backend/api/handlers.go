package api

import (
	"net/http"
	"star/astro"
	"star/models"

	"github.com/gin-gonic/gin"
)

// RootHandler 根路径处理器
func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":    "Star API (Go)",
		"version":    "2.0.0",
		"status":     "running",
		"dataSource": astro.GetDataSource(),
		"endpoints": map[string]string{
			"health":           "GET /health",
			"chart":            "POST /api/calc/chart",
			"astro":            "POST /api/v2/astro",
			"monitor_dashboard": "GET /api/monitor/dashboard",
			"monitor_summary":  "GET /api/monitor/summary",
		},
		"docs": "See /docs/API-REFERENCE.md for detailed API documentation",
	})
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"service":    "Star API (Go)",
		"version":    "2.0.0",
		"dataSource": astro.GetDataSource(),
		"features": []string{
			"natal-chart",
			"v2-unified-api",
			"five-dimension-forecast",
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
