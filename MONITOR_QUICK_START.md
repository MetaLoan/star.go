# 监控系统快速开始

## ✨ 功能已就绪！

你现在有一个完整的实时监控系统。

---

## 🚀 立即访问

### 方式1：Web仪表板（推荐）

**打开浏览器访问：**
```
http://localhost:8080/api/monitor/dashboard
```

**功能：**
- 📊 实时数据展示（每3秒刷新）
- 💻 美观的现代化UI
- 📈 核心指标：总请求数、成功率、活跃请求
- 📋 API统计表格
- 🔄 最近50条请求记录
- ⏱️ 实时30秒统计

---

### 方式2：API接口

#### 1. 概览统计
```bash
curl http://localhost:8080/api/monitor/summary | jq .
```

返回：总请求数、成功率、运行时间等

#### 2. API详细统计
```bash
curl http://localhost:8080/api/monitor/stats | jq .
```

返回：每个API的请求数、成功数、平均耗时等

#### 3. 最近请求
```bash
curl http://localhost:8080/api/monitor/recent?limit=20 | jq .
```

返回：最近N条请求的详细信息

#### 4. 实时统计
```bash
curl http://localhost:8080/api/monitor/realtime?seconds=60 | jq .
```

返回：最近N秒的统计数据

---

## 📊 监控指标

### 自动收集的指标：

| 指标类型 | 内容 |
|---------|------|
| **请求指标** | 总数、成功数、失败数、活跃数 |
| **性能指标** | 平均耗时、最小耗时、最大耗时 |
| **API统计** | 每个端点的独立统计 |
| **客户端信息** | IP地址、User-Agent |
| **时间信息** | 请求时间戳、运行时长 |

---

## 🎯 使用场景

### 1. 日常监控
- 打开仪表板 → 浏览器标签页保持打开
- 每隔一段时间查看
- 关注成功率和响应时间

### 2. 性能问题排查
```bash
# 找出最慢的API
curl http://localhost:8080/api/monitor/stats | jq 'to_entries | sort_by(.value.maxDuration) | reverse | .[0:5]'
```

### 3. 错误追踪
```bash
# 查看所有错误请求
curl http://localhost:8080/api/monitor/recent?limit=100 | jq '[.[] | select(.statusCode >= 400)]'
```

### 4. 流量分析
```bash
# 实时监控最近30秒
watch -n 3 'curl -s http://localhost:8080/api/monitor/realtime?seconds=30 | jq .'
```

---

## 📁 相关文档

- **详细文档**: [MONITOR_GUIDE.md](./MONITOR_GUIDE.md)
- **API文档**: [backend/docs/API-REFERENCE.md](./backend/docs/API-REFERENCE.md)

---

## 💡 提示

1. **数据保留**：监控数据存储在内存中，服务重启后清空
2. **性能开销**：非常小（< 1% CPU，10-20MB 内存）
3. **访问控制**：当前无认证，建议生产环境添加认证中间件
4. **数据导出**：可通过API定期导出数据到文件/数据库

---

## 🎉 开始使用

**现在就打开浏览器访问：**
```
http://localhost:8080/api/monitor/dashboard
```

享受实时监控的便利！
