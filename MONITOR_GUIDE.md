# Star API 监控系统使用指南

**版本**：v1.0  
**更新日期**：2026-01-16

---

## 功能概述

监控系统提供实时的API调用统计和性能分析，帮助你：
- ✅ 实时查看API调用情况
- ✅ 监控响应时间和性能指标
- ✅ 分析错误率和成功率
- ✅ 追踪客户端IP和请求详情

---

## 快速开始

### 访问监控仪表板

**URL**：`http://localhost:8080/api/monitor/dashboard`

打开浏览器访问上述地址，即可看到一个美观的实时监控仪表板。

**特性**：
- 🔄 每3秒自动刷新
- 📊 实时数据展示
- 🎨 现代化UI设计
- 📱 响应式布局

### 仪表板功能

#### 1. 核心指标卡片

| 指标 | 说明 |
|------|------|
| **总请求数** | 自服务启动以来的所有请求 |
| **活跃请求** | 当前正在处理的请求数 |
| **成功率** | 2xx/3xx状态码的请求占比 |
| **最近1分钟** | 最近60秒内的请求数 |

#### 2. API端点统计表

展示每个API端点的详细统计：
- 请求方法和路径
- 总请求数、成功数、失败数
- 平均响应时间
- 最慢请求时间
- 最后访问时间

#### 3. 实时监控

最近30秒的实时统计：
- 请求数量
- 平均响应时间
- 状态码分布

#### 4. 最近请求日志

最新50条请求的详细记录：
- 请求时间
- 方法和路径
- 状态码
- 响应时间
- 客户端IP

---

## API接口

除了Web仪表板，还提供以下API接口供程序化访问：

### 1. 概览统计

**端点**：`GET /api/monitor/summary`

**响应示例**：
```json
{
  "startTime": "2026-01-16T18:00:00+08:00",
  "uptime": "2h15m30s",
  "uptimeSeconds": 8130,
  "totalRequests": 1523,
  "activeRequests": 3,
  "successRequests": 1498,
  "errorRequests": 25,
  "successRate": 98.36,
  "totalAPIs": 12,
  "requestsLastMin": 45,
  "avgRequestsPerMin": 11.2
}
```

**字段说明**：
- `startTime`: 服务启动时间
- `uptime`: 运行时长（人类可读）
- `uptimeSeconds`: 运行时长（秒）
- `totalRequests`: 总请求数
- `activeRequests`: 当前活跃请求数
- `successRequests`: 成功请求数（2xx/3xx）
- `errorRequests`: 错误请求数（4xx/5xx）
- `successRate`: 成功率（百分比）
- `totalAPIs`: API端点总数
- `requestsLastMin`: 最近1分钟请求数
- `avgRequestsPerMin`: 平均每分钟请求数

---

### 2. API统计

**端点**：`GET /api/monitor/stats`

**响应示例**：
```json
{
  "POST /api/calc/daily": {
    "path": "/api/calc/daily",
    "method": "POST",
    "totalRequests": 456,
    "successRequests": 450,
    "errorRequests": 6,
    "avgDuration": 892.5,
    "minDuration": 654,
    "maxDuration": 1523,
    "totalDuration": 407100,
    "lastAccess": "2026-01-16T18:15:30+08:00"
  },
  "POST /api/calc/time-series": {
    "path": "/api/calc/time-series",
    "method": "POST",
    "totalRequests": 234,
    "successRequests": 234,
    "errorRequests": 0,
    "avgDuration": 723.8,
    "minDuration": 512,
    "maxDuration": 1102,
    "totalDuration": 169369,
    "lastAccess": "2026-01-16T18:15:25+08:00"
  }
}
```

**字段说明**：
- `totalRequests`: 该API的总请求数
- `successRequests`: 成功请求数
- `errorRequests`: 失败请求数
- `avgDuration`: 平均响应时间（毫秒）
- `minDuration`: 最快响应时间（毫秒）
- `maxDuration`: 最慢响应时间（毫秒）
- `lastAccess`: 最后访问时间

---

### 3. 最近请求

**端点**：`GET /api/monitor/recent?limit=50`

**参数**：
- `limit`（可选）：返回记录数，默认50，最多保留1000条

**响应示例**：
```json
[
  {
    "path": "/api/calc/daily",
    "method": "POST",
    "statusCode": 200,
    "duration": 892,
    "timestamp": "2026-01-16T18:15:30.123456+08:00",
    "clientIP": "192.168.1.100",
    "userAgent": "Mozilla/5.0...",
    "responseSize": 2456,
    "requestSize": 512
  },
  {
    "path": "/api/calc/time-series",
    "method": "POST",
    "statusCode": 200,
    "duration": 723,
    "timestamp": "2026-01-16T18:15:25.789012+08:00",
    "clientIP": "192.168.1.100",
    "userAgent": "Mozilla/5.0...",
    "responseSize": 8734,
    "requestSize": 589
  }
]
```

---

### 4. 实时统计

**端点**：`GET /api/monitor/realtime?seconds=30`

**参数**：
- `seconds`（可选）：时间窗口，默认60秒

**响应示例**：
```json
{
  "timeWindow": 30,
  "requestCount": 45,
  "avgDuration": 856.3,
  "statusCodes": {
    "200": 42,
    "400": 2,
    "500": 1
  },
  "topPaths": {
    "/api/calc/daily": 23,
    "/api/calc/time-series": 12,
    "/health": 10
  }
}
```

**字段说明**：
- `timeWindow`: 统计的时间窗口（秒）
- `requestCount`: 该时间窗口内的请求数
- `avgDuration`: 平均响应时间（毫秒）
- `statusCodes`: 状态码分布
- `topPaths`: 热门请求路径

---

### 5. 重置统计

**端点**：`POST /api/monitor/reset`

**说明**：重置所有监控统计数据（慎用！）

**响应示例**：
```json
{
  "message": "监控统计已重置",
  "time": "2026-01-16T18:20:00+08:00"
}
```

---

## 使用场景

### 场景1：性能问题排查

```bash
# 1. 查看整体概况
curl http://localhost:8080/api/monitor/summary

# 2. 找出最慢的API
curl http://localhost:8080/api/monitor/stats | jq 'to_entries | sort_by(.value.maxDuration) | reverse | .[0:5]'

# 3. 查看最近的慢请求
curl http://localhost:8080/api/monitor/recent?limit=100 | jq '[.[] | select(.duration > 1000)]'
```

### 场景2：监控实时流量

```bash
# 实时查看最近30秒的请求统计
watch -n 3 'curl -s http://localhost:8080/api/monitor/realtime?seconds=30 | jq .'
```

### 场景3：分析错误情况

```bash
# 查看所有错误的API统计
curl http://localhost:8080/api/monitor/stats | jq 'to_entries | map(select(.value.errorRequests > 0))'

# 查看最近的错误请求
curl http://localhost:8080/api/monitor/recent?limit=100 | jq '[.[] | select(.statusCode >= 400)]'
```

### 场景4：客户端IP分析

```bash
# 统计最活跃的客户端IP
curl http://localhost:8080/api/monitor/recent?limit=1000 | jq 'group_by(.clientIP) | map({ip: .[0].clientIP, count: length}) | sort_by(.count) | reverse'
```

---

## 监控最佳实践

### 1. 定期检查

建议每天查看监控仪表板，关注：
- ✅ 成功率是否正常（> 95%）
- ✅ 平均响应时间是否稳定
- ✅ 是否有异常的错误峰值

### 2. 设置告警阈值

虽然当前系统未内置告警，但你可以通过定时脚本监控关键指标：

```bash
#!/bin/bash
# monitor_alert.sh - 简单的监控告警脚本

SUCCESS_RATE=$(curl -s http://localhost:8080/api/monitor/summary | jq '.successRate')
AVG_REQUESTS=$(curl -s http://localhost:8080/api/monitor/summary | jq '.avgRequestsPerMin')

if (( $(echo "$SUCCESS_RATE < 90" | bc -l) )); then
    echo "警告: 成功率低于90% (当前: $SUCCESS_RATE%)"
    # 发送告警邮件或通知
fi

if (( $(echo "$AVG_REQUESTS > 1000" | bc -l) )); then
    echo "警告: 请求量过高 (当前: $AVG_REQUESTS req/min)"
    # 发送告警邮件或通知
fi
```

### 3. 数据导出

监控数据保存在内存中，服务重启后会清空。如需持久化：

```bash
# 定期导出统计数据
crontab -e
# 添加：每小时导出一次
0 * * * * curl -s http://localhost:8080/api/monitor/stats > /var/log/star-api-stats-$(date +\%Y\%m\%d-\%H).json
```

---

## 性能影响

### 监控开销

监控中间件对性能的影响：
- CPU：< 1%
- 内存：约10-20MB（保存1000条最近请求）
- 响应时间：< 0.1ms

### 内存管理

- 最近请求：最多保留1000条（FIFO淘汰）
- API统计：按API路径聚合，无数量限制
- 自动清理：服务重启时清空所有数据

---

## 常见问题

### Q1: 监控数据能持久化吗？

A: 当前版本的监控数据存储在内存中，服务重启后会清空。如需持久化：
- 方案1：定期通过API导出数据到文件/数据库
- 方案2：集成Prometheus等专业监控系统（未来版本）

### Q2: 能否添加自定义指标？

A: 当前版本支持基础HTTP指标。如需自定义业务指标，可以：
1. 在`middleware/metrics.go`中扩展`RequestMetrics`结构
2. 在中间件中添加自定义字段的收集逻辑
3. 在仪表板中展示新字段

### Q3: 如何限制监控API的访问？

A: 建议添加认证中间件：

```go
// 在routes.go中
monitor := api.Group("/monitor")
monitor.Use(AuthMiddleware()) // 添加认证
{
    monitor.GET("/dashboard", GetMonitorDashboard)
    // ...
}
```

### Q4: 能否按时间范围查询历史数据？

A: 当前版本仅支持最近N条记录。如需按时间范围查询，建议：
- 定期导出数据到时序数据库（InfluxDB、TimescaleDB等）
- 或集成Prometheus + Grafana

---

## 未来规划

### v1.1 (计划中)

- [ ] 告警功能（邮件/Webhook）
- [ ] 数据持久化选项
- [ ] 按时间范围查询
- [ ] 更多可视化图表

### v2.0 (规划中)

- [ ] Prometheus metrics导出
- [ ] Grafana仪表板模板
- [ ] 分布式追踪集成
- [ ] 业务指标自定义

---

## 相关链接

- [API文档](./backend/docs/API-REFERENCE.md)
- [性能优化文档](./PERFORMANCE_OPTIMIZATION_2026-01-16.md)

---

## 支持

如有问题或建议，请联系开发团队。

**最后更新**：2026-01-16
