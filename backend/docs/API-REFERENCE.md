# Star API 参考文档

本文档描述 Star 占星计算平台后端的所有 API 接口。

---

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **版本**: 2.0.0

---

## 接口概览

| 路径 | 方法 | 功能 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/api/calc/chart` | POST | 基础星盘数据查询 |
| `/api/v2/astro` | POST | 五维运势统一接口（核心） |
| `/api/monitor/dashboard` | GET | 监控仪表板页面 |
| `/api/monitor/summary` | GET | 监控概览统计 |
| `/api/monitor/stats` | GET | API 统计详情 |
| `/api/monitor/recent` | GET | 最近请求日志 |
| `/api/monitor/realtime` | GET | 实时统计 |
| `/api/monitor/reset` | POST | 重置统计 |

---

## 健康检查

### GET /health

检查服务运行状态。

**响应示例**：

```json
{
  "status": "ok",
  "service": "Star API (Go)",
  "version": "2.0.0",
  "dataSource": "Swiss Ephemeris (High Precision)",
  "features": [
    "natal-chart",
    "v2-unified-api",
    "five-dimension-forecast"
  ]
}
```

---

## 核心数据结构

### BirthData（出生数据）

所有计算接口的基础请求数据。

```json
{
  "year": 1990,
  "month": 6,
  "day": 15,
  "hour": 12,
  "minute": 30,
  "second": 0,
  "latitude": 39.9042,
  "longitude": 116.4074,
  "timezone": 8
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `year` | int | ✓ | 出生年份 |
| `month` | int | ✓ | 出生月份（1-12） |
| `day` | int | ✓ | 出生日期（1-31） |
| `hour` | int | ✓ | 出生小时（0-23） |
| `minute` | int | ✓ | 出生分钟（0-59） |
| `second` | int | - | 出生秒数（默认 0） |
| `latitude` | float | ✓ | 出生地纬度（-90 到 90） |
| `longitude` | float | ✓ | 出生地经度（-180 到 180） |
| `timezone` | float | ✓ | 时区偏移（如北京为 8） |

### DimensionScores（五维度分数）

五维运势分数结构。

```json
{
  "overall": 72.5,
  "career": 70.5,
  "relationship": 65.3,
  "health": 72.1,
  "finance": 68.4,
  "spiritual": 64.7
}
```

| 维度 | 字段 | 说明 |
|------|------|------|
| 综合 | `overall` | 整体运势评分 |
| 事业 | `career` | 工作、职业发展 |
| 关系 | `relationship` | 人际、感情、合作 |
| 健康 | `health` | 身体、精力状态 |
| 财务 | `finance` | 财富、收入 |
| 灵性 | `spiritual` | 内在成长、直觉 |

---

## 基础星盘查询

### POST /api/calc/chart

计算本命盘数据。

**请求体**：

```json
{
  "year": 1990,
  "month": 6,
  "day": 15,
  "hour": 12,
  "minute": 30,
  "second": 0,
  "latitude": 39.9042,
  "longitude": 116.4074,
  "timezone": 8
}
```

**响应示例**：

```json
{
  "birthData": { ... },
  "planets": [
    {
      "id": "sun",
      "name": "太阳",
      "symbol": "☉",
      "longitude": 84.5,
      "latitude": 0,
      "sign": "gemini",
      "signName": "双子座",
      "signSymbol": "♊",
      "signDegree": 24.5,
      "retrograde": false,
      "house": 10,
      "dignityScore": 0
    }
  ],
  "houses": [
    {
      "house": 1,
      "cusp": 120.5,
      "sign": "leo",
      "signName": "狮子座"
    }
  ],
  "ascendant": 120.5,
  "midheaven": 30.2,
  "aspects": [
    {
      "planet1": "sun",
      "planet2": "moon",
      "aspectType": "trine",
      "exactAngle": 120,
      "actualAngle": 118.5,
      "orb": 1.5,
      "applying": true,
      "strength": 0.85
    }
  ],
  "elementBalance": {
    "fire": 0.3,
    "earth": 0.2,
    "air": 0.35,
    "water": 0.15
  },
  "modalityBalance": {
    "cardinal": 0.4,
    "fixed": 0.3,
    "mutable": 0.3
  },
  "dominantPlanets": ["sun", "mars"],
  "chartRuler": "sun"
}
```

**响应字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `planets` | array | 十大行星信息 |
| `houses` | array | 十二宫位信息 |
| `ascendant` | float | 上升点黄经度数 |
| `midheaven` | float | 中天黄经度数 |
| `aspects` | array | 行星相位列表 |
| `elementBalance` | object | 元素平衡（火/土/风/水） |
| `modalityBalance` | object | 模式平衡（本位/固定/变动） |
| `dominantPlanets` | array | 主导行星 |
| `chartRuler` | string | 命盘主星 |

---

## 五维运势统一接口（核心）

### POST /api/v2/astro

单一接口返回所有数据：五维分数、天体事件、影响因子、曲线数据。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `birth` | object | ✓ | 出生信息（BirthData） |
| `time` | string | ✓ | 查询时间，ISO 8601 格式 |
| `granularity` | string | - | 粒度：`hour`/`day`/`week`/`month`/`year`，默认 `day` |
| `language` | string | - | 语言：`zh`/`en`/`ru`，默认 `en` |

**请求示例**：

```bash
curl -X POST http://localhost:8080/api/v2/astro \
  -H "Content-Type: application/json" \
  -d '{
    "birth": {
      "year": 1990, "month": 5, "day": 15,
      "hour": 10, "minute": 30,
      "latitude": 39.9, "longitude": 116.4,
      "timezone": 8
    },
    "time": "2026-01-16T14:00:00+08:00",
    "granularity": "day",
    "language": "zh"
  }'
```

**响应结构**：

```json
{
  "slot": {
    "userId": "19900515_1030_39p90_116p40",
    "startTime": "2026-01-16T00:00:00+08:00",
    "endTime": "2026-01-17T00:00:00+08:00",
    "granularity": "day",
    "scores": {
      "overall": 72.5,
      "career": 70.2,
      "relationship": 65.8,
      "health": 75.1,
      "finance": 68.4,
      "spiritual": 80.3
    },
    "events": [
      {
        "eventId": "aspect_mars_trine_venus_20260116",
        "type": "aspect",
        "title": "火星 拱相 金星",
        "isPositive": true,
        "intensity": 0.85,
        "primaryPlanet": "mars",
        "primaryPlanetName": "火星",
        "impact": {
          "career": 2.5,
          "relationship": 4.0,
          "health": 1.0,
          "finance": 2.0,
          "spiritual": 1.5
        },
        "impactDelta": {
          "career": 0.5,
          "relationship": 1.2,
          "health": 0.0,
          "finance": 0.3,
          "spiritual": 0.2
        },
        "startTime": "2026-01-14T10:00:00Z",
        "endTime": "2026-01-18T10:00:00Z",
        "exactTime": "2026-01-16T08:30:00Z",
        "interpretation": "详细解读文本...",
        "advice": "行动建议..."
      }
    ],
    "delta": {
      "overall": 2.3,
      "dimensions": {
        "career": 1.5,
        "relationship": 2.0,
        "health": 0.5,
        "finance": 1.0,
        "spiritual": 0.8
      },
      "reason": "火星拱金星增强了社交能量"
    },
    "guidance": {
      "summary": "今日整体运势上扬...",
      "dos": ["主动沟通", "处理财务"],
      "donts": ["过度劳累"],
      "focus": "relationship"
    },
    "subSlots": [
      {
        "startTime": "2026-01-16T00:00:00+08:00",
        "scores": {
          "overall": 70.1,
          "career": 68.5,
          "relationship": 64.0,
          "health": 73.0,
          "finance": 67.0,
          "spiritual": 78.0
        },
        "eventCount": 5
      }
    ]
  },
  "meta": {
    "cached": false,
    "computeTime": "450ms",
    "eventCount": 12
  }
}
```

**响应字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `slot.scores` | object | 当前时间段的五维分数 |
| `slot.events` | array | 天体事件列表 |
| `slot.delta` | object | 与上一周期的变化 |
| `slot.guidance` | object | 综合指导建议 |
| `slot.subSlots` | array | 子周期曲线数据 |
| `meta.cached` | bool | 是否缓存命中 |
| `meta.computeTime` | string | 计算耗时 |

**事件类型**：

| type | 说明 |
|------|------|
| `aspect` | 行运相位 |
| `lunar_phase` | 月相事件 |
| `sign_change` | 行星换座 |
| `dignity` | 尊贵度状态 |
| `retrograde` | 逆行状态 |
| `transit_house` | 行星过宫 |
| `void_of_course` | 月亮空亡 |
| `planetary_hour` | 行星时 |

**事件过滤规则**：

- **返回与时间槽有交集的所有事件**
- 事件只需与查询的时间范围有交集即可返回
- `startTime` 可能早于时间槽开始，`endTime` 可能晚于时间槽结束
- 客户端可根据 `startTime`/`endTime` 自行决定展示逻辑

**粒度与子周期**：

| 查询粒度 | subSlots 内容 |
|---------|--------------|
| `hour` | 无 |
| `day` | 24 个小时分数 |
| `week` | 7 天分数 |
| `month` | 30 天分数 |
| `year` | 12 个月分数 |

---

## 监控接口

### GET /api/monitor/dashboard

返回监控仪表板 HTML 页面。

**访问方式**：浏览器打开 `http://localhost:8080/api/monitor/dashboard`

**功能**：
- 核心指标卡片（总请求数、活跃请求、成功率）
- API 端点统计表格
- 实时监控数据
- 最近请求日志
- 每 3 秒自动刷新

---

### GET /api/monitor/summary

获取监控概览统计。

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
  "totalAPIs": 8,
  "requestsLastMin": 45,
  "avgRequestsPerMin": 11.2
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `uptime` | string | 运行时长（人类可读） |
| `totalRequests` | int | 总请求数 |
| `activeRequests` | int | 正在处理的请求数 |
| `successRate` | float | 成功率（百分比） |
| `requestsLastMin` | int | 最近 1 分钟请求数 |

---

### GET /api/monitor/stats

获取每个 API 端点的详细统计。

**响应示例**：

```json
{
  "POST /api/v2/astro": {
    "path": "/api/v2/astro",
    "method": "POST",
    "totalRequests": 456,
    "successRequests": 450,
    "errorRequests": 6,
    "avgDuration": 892.5,
    "minDuration": 654,
    "maxDuration": 1523,
    "lastAccess": "2026-01-16T18:15:30+08:00"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `totalRequests` | int | 总请求数 |
| `avgDuration` | float | 平均响应时间（毫秒） |
| `maxDuration` | int | 最慢响应时间（毫秒） |

---

### GET /api/monitor/recent

获取最近的请求记录。

**查询参数**：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `limit` | int | 50 | 返回记录数（最多 1000） |

**响应示例**：

```json
[
  {
    "path": "/api/v2/astro",
    "method": "POST",
    "statusCode": 200,
    "duration": 892,
    "timestamp": "2026-01-16T18:15:30.123456+08:00",
    "clientIP": "192.168.1.100",
    "userAgent": "Mozilla/5.0...",
    "responseSize": 2456,
    "requestSize": 512
  }
]
```

---

### GET /api/monitor/realtime

获取实时统计数据。

**查询参数**：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `seconds` | int | 60 | 时间窗口（秒） |

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
    "/api/v2/astro": 23,
    "/health": 10
  }
}
```

---

### POST /api/monitor/reset

重置所有监控统计数据。

**响应示例**：

```json
{
  "message": "监控统计已重置",
  "time": "2026-01-16T18:30:00+08:00"
}
```

---

## 错误响应

所有接口在出错时返回统一格式：

```json
{
  "error": "错误类型",
  "message": "详细错误描述"
}
```

**HTTP 状态码**：

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务内部错误 |

---

## 使用示例

### Python 示例

```python
import requests

# 查询五维运势
response = requests.post(
    'http://localhost:8080/api/v2/astro',
    json={
        'birth': {
            'year': 1990, 'month': 5, 'day': 15,
            'hour': 10, 'minute': 30,
            'latitude': 39.9, 'longitude': 116.4,
            'timezone': 8
        },
        'time': '2026-01-16T14:00:00+08:00',
        'granularity': 'day',
        'language': 'zh'
    }
)

data = response.json()
slot = data['slot']

print(f"综合运势: {slot['scores']['overall']}")
print(f"事业运势: {slot['scores']['career']}")
print(f"关系运势: {slot['scores']['relationship']}")

for event in slot['events'][:3]:
    print(f"事件: {event['title']} - {'正面' if event['isPositive'] else '负面'}")
```

### JavaScript 示例

```javascript
async function getAstroData(birthData, queryTime) {
  const response = await fetch('http://localhost:8080/api/v2/astro', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      birth: birthData,
      time: queryTime,
      granularity: 'day',
      language: 'zh'
    })
  });
  
  const data = await response.json();
  return data.slot;
}

// 使用示例
const birthData = {
  year: 1990, month: 5, day: 15,
  hour: 10, minute: 30,
  latitude: 39.9, longitude: 116.4,
  timezone: 8
};

getAstroData(birthData, '2026-01-16T14:00:00+08:00')
  .then(slot => {
    console.log('五维分数:', slot.scores);
    console.log('事件数量:', slot.events.length);
    console.log('指导建议:', slot.guidance.summary);
  });
```

### cURL 示例

```bash
# 健康检查
curl http://localhost:8080/health

# 计算本命盘
curl -X POST http://localhost:8080/api/calc/chart \
  -H "Content-Type: application/json" \
  -d '{
    "year": 1990, "month": 5, "day": 15,
    "hour": 10, "minute": 30,
    "latitude": 39.9, "longitude": 116.4,
    "timezone": 8
  }'

# 查询五维运势（日粒度）
curl -X POST http://localhost:8080/api/v2/astro \
  -H "Content-Type: application/json" \
  -d '{
    "birth": {
      "year": 1990, "month": 5, "day": 15,
      "hour": 10, "minute": 30,
      "latitude": 39.9, "longitude": 116.4,
      "timezone": 8
    },
    "time": "2026-01-16T14:00:00+08:00",
    "granularity": "day",
    "language": "zh"
  }'

# 查询周运势曲线
curl -X POST http://localhost:8080/api/v2/astro \
  -H "Content-Type: application/json" \
  -d '{
    "birth": {
      "year": 1990, "month": 5, "day": 15,
      "hour": 10, "minute": 30,
      "latitude": 39.9, "longitude": 116.4,
      "timezone": 8
    },
    "time": "2026-01-16T14:00:00+08:00",
    "granularity": "week",
    "language": "zh"
  }'

# 查看监控统计
curl http://localhost:8080/api/monitor/summary
```

---

## 性能说明

| 接口 | 典型响应时间 |
|------|-------------|
| `/health` | < 10ms |
| `/api/calc/chart` | < 100ms |
| `/api/v2/astro` (hour) | < 200ms |
| `/api/v2/astro` (day) | < 500ms |
| `/api/v2/astro` (week) | < 1s |
| `/api/v2/astro` (month) | < 2s |
| `/api/v2/astro` (year) | < 5s |

**缓存机制**：
- V2 接口内置多级缓存
- 相同查询参数会命中缓存，响应时间 < 50ms
- 缓存 TTL 根据粒度自动设置

---

## 版本历史

### 2.0.1 (2026-02-05)

- 移除 `phase` 字段：客户端根据 startTime/endTime 自行判断状态
- 优化事件过滤：返回与时间槽有交集的所有事件
- 修复 Dignity/SignChange 事件时间精度：使用真实的星座进入/离开时间
- 修复并发安全问题（positionCache, longitudeCache, customFactorStore）
- 扩展 i18n 翻译：添加南交点、次要相位、月相、尊贵度等翻译
- 修复 "North Node" 等多词行星名称的解析问题

### 2.0.0 (2026-01-16)

- 精简 API 接口，移除冗余接口
- 保留核心功能：星盘查询、五维运势、监控
- V2 统一接口作为主要数据接口
