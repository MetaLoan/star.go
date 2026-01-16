package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestMetrics 请求指标
type RequestMetrics struct {
	Path           string    `json:"path"`
	Method         string    `json:"method"`
	StatusCode     int       `json:"statusCode"`
	Duration       int64     `json:"duration"` // 毫秒
	Timestamp      time.Time `json:"timestamp"`
	ClientIP       string    `json:"clientIP"`
	UserAgent      string    `json:"userAgent"`
	ResponseSize   int       `json:"responseSize"`
	RequestSize    int       `json:"requestSize"`
}

// APIStats API统计信息
type APIStats struct {
	Path            string  `json:"path"`
	Method          string  `json:"method"`
	TotalRequests   int64   `json:"totalRequests"`
	SuccessRequests int64   `json:"successRequests"`
	ErrorRequests   int64   `json:"errorRequests"`
	AvgDuration     float64 `json:"avgDuration"` // 毫秒
	MinDuration     int64   `json:"minDuration"`
	MaxDuration     int64   `json:"maxDuration"`
	TotalDuration   int64   `json:"totalDuration"`
	LastAccess      time.Time `json:"lastAccess"`
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
	mu              sync.RWMutex
	recentRequests  []RequestMetrics       // 最近的请求（保留最新1000条）
	apiStats        map[string]*APIStats   // 按API路径统计
	startTime       time.Time
	totalRequests   int64
	activeRequests  int64
}

var (
	collector *MetricsCollector
	once      sync.Once
)

// GetCollector 获取单例收集器
func GetCollector() *MetricsCollector {
	once.Do(func() {
		collector = &MetricsCollector{
			recentRequests: make([]RequestMetrics, 0, 1000),
			apiStats:      make(map[string]*APIStats),
			startTime:     time.Now(),
		}
	})
	return collector
}

// MetricsMiddleware Gin中间件
func MetricsMiddleware() gin.HandlerFunc {
	collector := GetCollector()
	
	return func(c *gin.Context) {
		start := time.Now()
		
		// 增加活跃请求计数
		collector.mu.Lock()
		collector.activeRequests++
		collector.mu.Unlock()
		
		// 获取请求大小
		requestSize := c.Request.ContentLength
		
		// 处理请求
		c.Next()
		
		// 计算耗时
		duration := time.Since(start).Milliseconds()
		
		// 获取响应信息
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()
		
		// 记录指标
		metrics := RequestMetrics{
			Path:         c.Request.URL.Path,
			Method:       c.Request.Method,
			StatusCode:   statusCode,
			Duration:     duration,
			Timestamp:    time.Now(),
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			ResponseSize: responseSize,
			RequestSize:  int(requestSize),
		}
		
		collector.RecordRequest(metrics)
		
		// 减少活跃请求计数
		collector.mu.Lock()
		collector.activeRequests--
		collector.mu.Unlock()
	}
}

// RecordRequest 记录请求
func (mc *MetricsCollector) RecordRequest(metrics RequestMetrics) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	// 记录到最近请求列表
	mc.recentRequests = append(mc.recentRequests, metrics)
	if len(mc.recentRequests) > 1000 {
		mc.recentRequests = mc.recentRequests[1:]
	}
	
	// 更新总计数
	mc.totalRequests++
	
	// 更新API统计
	key := metrics.Method + " " + metrics.Path
	stats, exists := mc.apiStats[key]
	if !exists {
		stats = &APIStats{
			Path:        metrics.Path,
			Method:      metrics.Method,
			MinDuration: metrics.Duration,
			MaxDuration: metrics.Duration,
		}
		mc.apiStats[key] = stats
	}
	
	stats.TotalRequests++
	stats.TotalDuration += metrics.Duration
	stats.AvgDuration = float64(stats.TotalDuration) / float64(stats.TotalRequests)
	stats.LastAccess = metrics.Timestamp
	
	if metrics.StatusCode >= 200 && metrics.StatusCode < 400 {
		stats.SuccessRequests++
	} else {
		stats.ErrorRequests++
	}
	
	if metrics.Duration < stats.MinDuration {
		stats.MinDuration = metrics.Duration
	}
	if metrics.Duration > stats.MaxDuration {
		stats.MaxDuration = metrics.Duration
	}
}

// GetRecentRequests 获取最近的请求（限制数量）
func (mc *MetricsCollector) GetRecentRequests(limit int) []RequestMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	total := len(mc.recentRequests)
	if limit > total {
		limit = total
	}
	
	result := make([]RequestMetrics, limit)
	copy(result, mc.recentRequests[total-limit:])
	
	// 反转数组，最新的在前
	for i := 0; i < len(result)/2; i++ {
		j := len(result) - 1 - i
		result[i], result[j] = result[j], result[i]
	}
	
	return result
}

// GetAPIStats 获取API统计
func (mc *MetricsCollector) GetAPIStats() map[string]*APIStats {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	// 创建副本
	result := make(map[string]*APIStats)
	for k, v := range mc.apiStats {
		statsCopy := *v
		result[k] = &statsCopy
	}
	return result
}

// GetSummary 获取总体统计摘要
func (mc *MetricsCollector) GetSummary() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	uptime := time.Since(mc.startTime)
	
	// 计算总体成功率
	var totalSuccess, totalError int64
	for _, stats := range mc.apiStats {
		totalSuccess += stats.SuccessRequests
		totalError += stats.ErrorRequests
	}
	
	successRate := float64(0)
	if mc.totalRequests > 0 {
		successRate = float64(totalSuccess) / float64(mc.totalRequests) * 100
	}
	
	// 最近1分钟的请求数
	oneMinuteAgo := time.Now().Add(-time.Minute)
	recentCount := 0
	for i := len(mc.recentRequests) - 1; i >= 0; i-- {
		if mc.recentRequests[i].Timestamp.After(oneMinuteAgo) {
			recentCount++
		} else {
			break
		}
	}
	
	return map[string]interface{}{
		"startTime":        mc.startTime,
		"uptime":           uptime.String(),
		"uptimeSeconds":    int64(uptime.Seconds()),
		"totalRequests":    mc.totalRequests,
		"activeRequests":   mc.activeRequests,
		"successRequests":  totalSuccess,
		"errorRequests":    totalError,
		"successRate":      successRate,
		"totalAPIs":        len(mc.apiStats),
		"requestsLastMin":  recentCount,
		"avgRequestsPerMin": float64(mc.totalRequests) / (uptime.Minutes() + 0.001),
	}
}

// GetRealTimeStats 获取实时统计（最近N秒）
func (mc *MetricsCollector) GetRealTimeStats(seconds int) map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	cutoff := time.Now().Add(-time.Duration(seconds) * time.Second)
	
	var count int
	var totalDuration int64
	var statusCodes = make(map[int]int)
	var paths = make(map[string]int)
	
	for i := len(mc.recentRequests) - 1; i >= 0; i-- {
		req := mc.recentRequests[i]
		if req.Timestamp.After(cutoff) {
			count++
			totalDuration += req.Duration
			statusCodes[req.StatusCode]++
			paths[req.Path]++
		} else {
			break
		}
	}
	
	avgDuration := float64(0)
	if count > 0 {
		avgDuration = float64(totalDuration) / float64(count)
	}
	
	return map[string]interface{}{
		"timeWindow":    seconds,
		"requestCount":  count,
		"avgDuration":   avgDuration,
		"statusCodes":   statusCodes,
		"topPaths":      paths,
	}
}

// Reset 重置所有统计（慎用）
func (mc *MetricsCollector) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.recentRequests = make([]RequestMetrics, 0, 1000)
	mc.apiStats = make(map[string]*APIStats)
	mc.startTime = time.Now()
	mc.totalRequests = 0
	mc.activeRequests = 0
}
