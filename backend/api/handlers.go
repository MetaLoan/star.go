package api

import (
	"net/http"
	"star/astro"
	"star/models"

	"github.com/gin-gonic/gin"
)

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
