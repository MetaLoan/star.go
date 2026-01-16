package api

import (
	"net/http"
	"strconv"
	"star/middleware"

	"github.com/gin-gonic/gin"
)

// GetMonitorDashboard 返回监控仪表板HTML页面
func GetMonitorDashboard(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Star API 监控仪表板</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            color: #333;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        .header {
            background: white;
            padding: 20px 30px;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            margin-bottom: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .header h1 {
            color: #667eea;
            font-size: 28px;
        }
        .status-badge {
            background: #10b981;
            color: white;
            padding: 8px 16px;
            border-radius: 20px;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .status-dot {
            width: 8px;
            height: 8px;
            background: white;
            border-radius: 50%;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }
        .card {
            background: white;
            padding: 25px;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .card-title {
            font-size: 14px;
            color: #6b7280;
            margin-bottom: 10px;
            font-weight: 600;
            text-transform: uppercase;
        }
        .card-value {
            font-size: 36px;
            font-weight: 700;
            color: #1f2937;
            margin-bottom: 5px;
        }
        .card-subtitle {
            font-size: 13px;
            color: #9ca3af;
        }
        .card-large {
            grid-column: span 2;
        }
        .table-container {
            overflow-x: auto;
            margin-top: 15px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 14px;
        }
        th {
            background: #f3f4f6;
            padding: 12px;
            text-align: left;
            font-weight: 600;
            color: #374151;
            border-bottom: 2px solid #e5e7eb;
        }
        td {
            padding: 12px;
            border-bottom: 1px solid #e5e7eb;
        }
        tr:hover {
            background: #f9fafb;
        }
        .method-badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 600;
        }
        .method-GET { background: #dbeafe; color: #1e40af; }
        .method-POST { background: #d1fae5; color: #065f46; }
        .method-PUT { background: #fef3c7; color: #92400e; }
        .method-DELETE { background: #fee2e2; color: #991b1b; }
        .status-success { color: #10b981; font-weight: 600; }
        .status-error { color: #ef4444; font-weight: 600; }
        .refresh-info {
            text-align: center;
            color: white;
            margin-top: 20px;
            font-size: 14px;
        }
        .trend-up { color: #10b981; }
        .trend-down { color: #ef4444; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🌟 Star API 监控仪表板</h1>
            <div class="status-badge">
                <div class="status-dot"></div>
                <span>运行中</span>
            </div>
        </div>

        <div class="grid">
            <div class="card">
                <div class="card-title">总请求数</div>
                <div class="card-value" id="totalRequests">-</div>
                <div class="card-subtitle">自启动以来</div>
            </div>
            <div class="card">
                <div class="card-title">活跃请求</div>
                <div class="card-value" id="activeRequests">-</div>
                <div class="card-subtitle">正在处理中</div>
            </div>
            <div class="card">
                <div class="card-title">成功率</div>
                <div class="card-value" id="successRate">-</div>
                <div class="card-subtitle">所有请求</div>
            </div>
            <div class="card">
                <div class="card-title">最近1分钟</div>
                <div class="card-value" id="requestsLastMin">-</div>
                <div class="card-subtitle" id="avgPerMin">-</div>
            </div>
        </div>

        <div class="grid">
            <div class="card card-large">
                <div class="card-title">API 端点统计</div>
                <div class="table-container">
                    <table id="apiStatsTable">
                        <thead>
                            <tr>
                                <th>方法</th>
                                <th>路径</th>
                                <th>请求数</th>
                                <th>成功</th>
                                <th>失败</th>
                                <th>平均耗时</th>
                                <th>最慢</th>
                            </tr>
                        </thead>
                        <tbody id="apiStatsBody">
                            <tr><td colspan="7" style="text-align:center;color:#9ca3af;">加载中...</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div class="card">
                <div class="card-title">实时监控（最近30秒）</div>
                <div style="margin-top: 15px;">
                    <div style="margin-bottom: 15px;">
                        <div style="font-size: 13px; color: #6b7280; margin-bottom: 5px;">请求数</div>
                        <div style="font-size: 24px; font-weight: 600;" id="realtime30s">-</div>
                    </div>
                    <div style="margin-bottom: 15px;">
                        <div style="font-size: 13px; color: #6b7280; margin-bottom: 5px;">平均响应时间</div>
                        <div style="font-size: 24px; font-weight: 600;" id="avgDuration30s">-</div>
                    </div>
                    <div>
                        <div style="font-size: 13px; color: #6b7280; margin-bottom: 5px;">状态码分布</div>
                        <div id="statusCodes30s" style="font-size: 12px; color: #4b5563;">-</div>
                    </div>
                </div>
            </div>
        </div>

        <div class="card">
            <div class="card-title">最近请求 (最新50条)</div>
            <div class="table-container">
                <table>
                    <thead>
                        <tr>
                            <th>时间</th>
                            <th>方法</th>
                            <th>路径</th>
                            <th>状态码</th>
                            <th>耗时</th>
                            <th>客户端IP</th>
                        </tr>
                    </thead>
                    <tbody id="recentRequestsBody">
                        <tr><td colspan="6" style="text-align:center;color:#9ca3af;">加载中...</td></tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="refresh-info">
            ⏱️ 每3秒自动刷新 · 运行时间: <span id="uptime">-</span>
        </div>
    </div>

    <script>
        function formatDuration(ms) {
            if (ms < 1000) return ms + 'ms';
            return (ms / 1000).toFixed(2) + 's';
        }

        function formatTime(timestamp) {
            const date = new Date(timestamp);
            return date.toLocaleTimeString('zh-CN', { hour12: false });
        }

        function updateDashboard() {
            // 获取概览数据
            fetch('/api/monitor/summary')
                .then(res => res.json())
                .then(data => {
                    document.getElementById('totalRequests').textContent = data.totalRequests.toLocaleString();
                    document.getElementById('activeRequests').textContent = data.activeRequests;
                    document.getElementById('successRate').textContent = data.successRate.toFixed(1) + '%';
                    document.getElementById('requestsLastMin').textContent = data.requestsLastMin;
                    document.getElementById('avgPerMin').textContent = '平均 ' + data.avgRequestsPerMin.toFixed(1) + ' req/min';
                    document.getElementById('uptime').textContent = data.uptime;
                });

            // 获取API统计
            fetch('/api/monitor/stats')
                .then(res => res.json())
                .then(data => {
                    const tbody = document.getElementById('apiStatsBody');
                    const stats = Object.values(data).sort((a, b) => b.totalRequests - a.totalRequests);
                    
                    tbody.innerHTML = stats.map(stat => {
                        const methodClass = 'method-' + stat.method;
                        return` + " `" + `
                            <tr>
                                <td><span class="method-badge ${methodClass}">${stat.method}</span></td>
                                <td style="font-family: monospace; font-size: 12px;">${stat.path}</td>
                                <td>${stat.totalRequests.toLocaleString()}</td>
                                <td class="status-success">${stat.successRequests}</td>
                                <td class="status-error">${stat.errorRequests}</td>
                                <td>${stat.avgDuration.toFixed(0)}ms</td>
                                <td>${stat.maxDuration}ms</td>
                            </tr>
                        ` + "`" + `;
                    }).join('');
                });

            // 获取最近请求
            fetch('/api/monitor/recent?limit=50')
                .then(res => res.json())
                .then(data => {
                    const tbody = document.getElementById('recentRequestsBody');
                    tbody.innerHTML = data.map(req => {
                        const methodClass = 'method-' + req.method;
                        const statusClass = req.statusCode >= 200 && req.statusCode < 400 ? 'status-success' : 'status-error';
                        return` + " `" + `
                            <tr>
                                <td style="font-size: 12px;">${formatTime(req.timestamp)}</td>
                                <td><span class="method-badge ${methodClass}">${req.method}</span></td>
                                <td style="font-family: monospace; font-size: 12px;">${req.path}</td>
                                <td class="${statusClass}">${req.statusCode}</td>
                                <td>${req.duration}ms</td>
                                <td style="font-size: 12px;">${req.clientIP}</td>
                            </tr>
                        ` + "`" + `;
                    }).join('');
                });

            // 获取实时统计
            fetch('/api/monitor/realtime?seconds=30')
                .then(res => res.json())
                .then(data => {
                    document.getElementById('realtime30s').textContent = data.requestCount;
                    document.getElementById('avgDuration30s').textContent = data.avgDuration.toFixed(0) + 'ms';
                    
                    const statusCodesHtml = Object.entries(data.statusCodes || {})
                        .map(([code, count]) =>` + " `<span style=\"margin-right: 10px;\">${code}: ${count}</span>`" + `)
                        .join('');
                    document.getElementById('statusCodes30s').innerHTML = statusCodesHtml || '无数据';
                });
        }

        // 初始加载
        updateDashboard();
        
        // 每3秒刷新
        setInterval(updateDashboard, 3000);
    </script>
</body>
</html>`
	
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// GetMonitorSummary 获取监控摘要
func GetMonitorSummary(c *gin.Context) {
	collector := middleware.GetCollector()
	summary := collector.GetSummary()
	c.JSON(http.StatusOK, summary)
}

// GetMonitorStats 获取API统计
func GetMonitorStats(c *gin.Context) {
	collector := middleware.GetCollector()
	stats := collector.GetAPIStats()
	c.JSON(http.StatusOK, stats)
}

// GetMonitorRecentRequests 获取最近的请求
func GetMonitorRecentRequests(c *gin.Context) {
	collector := middleware.GetCollector()
	
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	
	requests := collector.GetRecentRequests(limit)
	c.JSON(http.StatusOK, requests)
}

// GetMonitorRealtime 获取实时统计
func GetMonitorRealtime(c *gin.Context) {
	collector := middleware.GetCollector()
	
	seconds := 60
	if secondsStr := c.Query("seconds"); secondsStr != "" {
		if s, err := strconv.Atoi(secondsStr); err == nil && s > 0 {
			seconds = s
		}
	}
	
	stats := collector.GetRealTimeStats(seconds)
	c.JSON(http.StatusOK, stats)
}

// ResetMonitorStats 重置统计（需要管理员权限）
func ResetMonitorStats(c *gin.Context) {
	collector := middleware.GetCollector()
	collector.Reset()
	c.JSON(http.StatusOK, gin.H{
		"message": "监控统计已重置",
		"time":    collector.GetSummary()["startTime"],
	})
}
