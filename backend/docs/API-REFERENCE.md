# Star API 参考文档

本文档描述 Star 占星计算平台后端的所有 API 接口。

算法说明见：[ALGORITHM.md](./ALGORITHM.md)

---

## 基础信息

- **Base URL**: `http://localhost:8888`
- **Content-Type**: `application/json`
- **版本**: 2.0.0

---

## 接口概览

| 路径 | 方法 | 功能 |
|------|------|------|
| `/` | GET | 根路径信息 |
| `/health` | GET | 健康检查 |
| `/api/calc/chart` | POST | 基础星盘数据查询 |
| `/api/v2/astro` | POST | 五维运势统一接口（核心） |
| `/api/v2/astro/trend` | POST | 趋势图数据接口（独立） |
| `/api/monitor/dashboard` | GET | 监控仪表板页面 |
| `/api/monitor/summary` | GET | 监控概览统计 |
| `/api/monitor/stats` | GET | API 统计详情 |
| `/api/monitor/recent` | GET | 最近请求日志 |
| `/api/monitor/realtime` | GET | 实时统计 |
| `/api/monitor/reset` | POST | 重置统计 |

## 接口职责速查

| 接口 | 内部核心入口 | 主要功能 | 适合场景 |
|------|------|------|------|
| `/health` | 无 | 服务状态检查 | 启动探活、监控 |
| `/api/calc/chart` | `CalculateNatalChart` | 生成本命盘 | 展示出生盘基础信息 |
| `/api/v2/astro` | `CalculateScoresV2` / `CalculateScoreBreakdown` / `CalculateDailyEvents` | 返回分数、事件、指导、delta | 主页面、详情页、卡片页 |
| `/api/v2/astro/trend` | `CalculateScoresV2Lite` | 返回趋势曲线 | 折线图、趋势面板 |
| `/api/monitor/*` | 监控模块 | 统计请求和性能 | 运维、排障、观察系统健康 |

大白话：

- `chart` 接口只负责“出生盘长什么样”
- `astro` 接口负责“某个时刻整体怎么样”
- `trend` 接口负责“这一段时间是往上还是往下”
- `monitor` 接口负责“系统自己运行得好不好”

---

## 健康检查

### GET /

返回服务根路径信息，适合快速查看服务状态和已暴露入口。

**响应示例**：

```json
{
  "service": "Star API (Go)",
  "version": "2.0.0",
  "status": "running",
  "dataSource": "Swiss Ephemeris (High Precision)",
  "endpoints": {
    "health": "GET /health",
    "chart": "POST /api/calc/chart",
    "astro": "POST /api/v2/astro",
    "monitor_dashboard": "GET /api/monitor/dashboard",
    "monitor_summary": "GET /api/monitor/summary"
  },
  "docs": "See /docs/API-REFERENCE.md for detailed API documentation"
}
```

**字段说明**：

- `service`：服务名称
- `version`：当前版本号
- `status`：当前服务状态
- `dataSource`：当前使用的数据源
- `endpoints`：已注册的核心入口
- `docs`：文档提示信息

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
| `name` | string | - | 出生者名称或昵称，主要用于展示和回显 |
| `year` | int | ✓ | 出生年份 |
| `month` | int | ✓ | 出生月份（1-12） |
| `day` | int | ✓ | 出生日期（1-31） |
| `hour` | int | ✓ | 出生小时（0-23） |
| `minute` | int | ✓ | 出生分钟（0-59） |
| `second` | int | - | 出生秒数（默认 0） |
| `latitude` | float | ✓ | 出生地纬度（-90 到 90） |
| `longitude` | float | ✓ | 出生地经度（-180 到 180） |
| `timezone` | float | ✓ | 时区偏移（如北京为 8，支持 5.5 这种半时区） |

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
  "birthData": {
    "name": "",
    "year": 1990,
    "month": 5,
    "day": 15,
    "hour": 10,
    "minute": 30,
    "second": 0,
    "latitude": 39.9,
    "longitude": 116.4,
    "timezone": 8
  },
  "planets": [
    {
      "id": "sun",
      "name": "Sun",
      "symbol": "☉",
      "longitude": 54.3,
      "latitude": 0,
      "sign": "taurus",
      "signName": "Taurus",
      "signSymbol": "♉",
      "signDegree": 24.3,
      "retrograde": false,
      "house": 7,
      "dignityScore": 0
    }
  ],
  "houses": [
    {
      "house": 1,
      "cusp": 225.1,
      "sign": "scorpio",
      "signName": "Scorpio"
    }
  ],
  "ascendant": 225.1,
  "midheaven": 144.5,
  "aspects": [
    {
      "planet1": "sun",
      "planet2": "moon",
      "aspectType": "trine",
      "exactAngle": 120,
      "actualAngle": 117.9,
      "orb": 2.1,
      "applying": true,
      "strength": 0.74,
      "weight": 4.46,
      "interpretation": "Sun forms Trine with Moon"
    }
  ],
  "elementBalance": {
    "fire": 7.04,
    "earth": 61.97,
    "air": 4.23,
    "water": 26.76
  },
  "modalityBalance": {
    "cardinal": 54.93,
    "fixed": 36.62,
    "mutable": 8.45
  },
  "patterns": ["Grand Trine", "T-Square"],
  "dominantPlanets": ["moon", "saturn", "sun", "jupiter", "uranus"],
  "chartRuler": "pluto"
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
| `patterns` | array | 识别出的整体结构模式 |
| `dominantPlanets` | array | 主导行星 |
| `chartRuler` | string | 命盘主星 |

**字段解释**：

- `birthData`：本次计算使用的出生信息，会把请求里的出生参数原样回显出来
- `planets`：每颗行星当前落在哪个星座、哪个宫位、是否逆行
- `houses`：12 个生活领域的分区，类似“人生地图的 12 个房间”
- `ascendant`：上升点，代表出生那一刻东方地平线升起的位置
- `midheaven`：中天，通常和事业、社会角色、公众表现有关
- `aspects`：行星彼此之间的角度关系，决定配合还是冲突
- `elementBalance`：火、土、风、水四种元素的相对占比，当前实现通常是百分比数值
- `modalityBalance`：本位、固定、变动三种模式的相对占比，当前实现通常是百分比数值
- `patterns`：识别出来的结构模式，比如大三角、T 三角、T 字架
- `dominantPlanets`：对这张盘最有存在感的几颗星
- `chartRuler`：上升星座的守护星，类似整张盘的“总控星”

**行星与相位的附加字段**：

- `planets[].name`：行星英文名，当前实现直接返回数据源原名
- `planets[].signName`：星座英文名
- `planets[].dignityScore`：这颗星在当前星座里的尊贵度分值，越高说明状态越顺
- `aspects[].weight`：这条相位在整体盘里的权重，越大说明越重要
- `aspects[].interpretation`：相位的简短文字解释，通常用于调试或文案展示

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

**字段说明**：

- `birth`：这次计算对应谁的出生信息，是所有后续计算的底座
- `time`：你想看“这个时刻附近”的运势，不是出生时间
- `granularity`：决定输出是小时级、日级、周级、月级还是年级
- `language`：只影响文案语言，不影响分数本身

**返回的 `slot` 字段**：

- `userId`：这次计算对应的用户标识，通常由出生数据拼出来
- `startTime` / `endTime`：这个时间槽的起止范围
- `granularity`：这个结果属于哪个时间粒度
- `scores`：五维评分结果，`overall` 是总分，其余是分维度分
- `events`：这个时间槽里命中的天象事件
- `delta`：和上一周期相比的变化量
- `guidance`：给用户看的行动建议
- `subSlots`：用于画折线图的子时间点

**返回的 `meta` 字段**：

- `cached`：是否直接命中了缓存
- `computeTime`：后端实际计算耗时
- `eventCount`：返回了多少条事件

**补充说明**：

- 这个接口返回的是“已经整理好的结果”，不是完整的内部因子对象
- 如果你想看某个时间点的单个因子持续时间、峰值和剩余天数，要结合算法文档里的生命周期说明
- 如果你想看“为什么这个时间点变成这样”，优先看 `slot.events`、`slot.delta`，再去看分数拆解逻辑
- 如果你要查“某个时间点有哪些影响因子仍然活跃”，当前代码里对应的是内部方法 `GetActiveFactorsInRange`
- 这个方法更偏排查和分析，不是目前对外暴露的独立 HTTP 接口；前端要展示这类内容时，通常由 `astro` 结果加上算法层分析一起拼出来

**请求示例**：

```bash
curl -X POST http://localhost:8888/api/v2/astro \
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
| `meta.cacheAge` | string | 预留字段，当前 live 响应里通常不返回 |
| `meta.computeTime` | string | 计算耗时 |

**响应样例速读**：

- `slot.userId`：这条结果属于哪个用户
- `slot.startTime` / `slot.endTime`：当前查询覆盖的时间段
- `slot.scores.overall`：总分，越高表示这段时间整体越顺
- `slot.scores.career`：事业维度分数
- `slot.scores.relationship`：关系维度分数
- `slot.scores.health`：健康维度分数
- `slot.scores.finance`：财务维度分数
- `slot.scores.spiritual`：灵性维度分数
- `slot.events[0].type`：这条事件属于哪类天象
- `slot.events[0].title`：前端展示给用户看的短标题
- `slot.events[0].isPositive`：这条事件是利好还是压力
- `slot.events[0].intensity`：事件强度，越大越明显
- `slot.events[0].impact`：这条事件对五维的“绝对影响”
- `slot.events[0].impactDelta`：跟上一周期相比，这条事件变化了多少
- `slot.events[0].exactTime`：最强的卡点时间
- `slot.delta.reason`：为什么这段时间比上一段更好或更差
- `slot.guidance.summary`：一句话总结
- `slot.guidance.dos`：建议做什么
- `slot.guidance.donts`：建议少做什么
- `slot.guidance.focus`：当前最值得关注的维度
- `slot.subSlots`：子粒度采样点，用来画折线
- `meta.cached`：是否直接用缓存返回
- `meta.computeTime`：后端计算花了多久
- `meta.eventCount`：最终返回了多少事件
- `meta.cacheAge`：预留字段，当前 live 响应通常没有这个值

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
| `year` | 4 个季度分数 |

---

## 趋势图接口（独立）

### POST /api/v2/astro/trend

一次性返回指定粒度下的所有趋势数据点，用于绘制趋势曲线图。

**设计目的**：
- 前端无需多次请求，一次获取整条曲线数据
- 与主接口解耦，可并行请求，加速首屏加载
- 只返回分数，不包含事件/指导等重数据

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `birth` | object | ✓ | 出生信息（BirthData） |
| `start_time` | string | ✓ | 基准时间，ISO 8601 格式 |
| `granularity` | string | - | 粒度：`hour`/`day`/`week`/`month`/`year`，默认 `day` |
| `language` | string | - | 语言：`zh`/`en`/`ru`，默认 `en` |

**字段说明**：

- `birth`：同上，表示这个趋势属于谁
- `start_time`：趋势曲线从哪个时间点开始算
- `granularity`：决定趋势点的密度
- `language`：只影响标签和文案

**返回的 `points[]` 字段**：

- `time`：这个趋势点对应的实际时间
- `label`：前端展示标签，比如 `1月`、`W1`、`08:00`
- `scores`：这个时间点的五维分数

**返回的 `summary` 字段**：

- `max`：整条趋势曲线里的最高总分
- `min`：整条趋势曲线里的最低总分
- `trend`：整体趋势方向，`upward`、`downward` 或 `stable`

**返回的 `meta` 字段**：

- `cached`：是否命中缓存
- `computeTime`：这条趋势曲线算了多久

**补充说明**：

- 趋势接口只给“走势”，不直接给完整因子生命周期
- 如果你想知道某个趋势点为什么高或低，要回到对应时间点看 `slot.events` 和分数拆解
- 趋势本质上是很多点连起来的结果，不是单个事件决定的

**请求示例**：

```bash
curl -X POST http://localhost:8888/api/v2/astro/trend \
  -H "Content-Type: application/json" \
  -d '{
    "birth": {
      "year": 1990, "month": 5, "day": 15,
      "hour": 10, "minute": 30,
      "latitude": 39.9, "longitude": 116.4,
      "timezone": 8
    },
    "start_time": "2026-02-06",
    "granularity": "month",
    "language": "zh"
  }'
```

**响应结构**：

```json
{
  "granularity": "month",
  "points": [
    {
      "time": "2026-01-15T12:00:00+08:00",
      "label": "1月",
      "scores": {
        "overall": 72.5,
        "career": 70.2,
        "relationship": 65.8,
        "health": 75.1,
        "finance": 68.4,
        "spiritual": 80.3
      }
    },
    {
      "time": "2026-02-15T12:00:00+08:00",
      "label": "2月",
      "scores": { ... }
    }
  ],
  "summary": {
    "max": 78.5,
    "min": 65.2,
    "trend": "upward"
  },
  "meta": {
    "cached": false,
    "computeTime": "5.8s",
    "eventCount": 0
  }
}
```

**响应字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `granularity` | string | 返回数据的粒度 |
| `points` | array | 趋势数据点数组 |
| `points[].time` | string | 数据点时间 |
| `points[].label` | string | 显示标签（如 "Jan"、"1月"、"W1"） |
| `points[].scores` | object | 五维分数 |
| `summary.max` | float | 最高综合分 |
| `summary.min` | float | 最低综合分 |
| `summary.trend` | string | 趋势方向：`upward`/`downward`/`stable` |
| `meta.cached` | bool | 是否缓存命中 |
| `meta.computeTime` | string | 计算耗时 |
| `meta.eventCount` | int | 当前趋势接口复用的计数字段，通常为 `0` |

**响应样例速读**：

- `granularity`：这条趋势曲线是按小时、天、周、月还是年生成
- `points`：整条曲线上的所有点
- `points[0].time`：这个点对应的实际时间
- `points[0].label`：给前端显示的短标签
- `points[0].scores.overall`：这个点的总分
- `points[0].scores.career`：这个点的事业分
- `points[0].scores.relationship`：这个点的关系分
- `summary.max`：曲线里最高的总分
- `summary.min`：曲线里最低的总分
- `summary.trend`：整体走势是向上、向下还是横盘
- `meta.cached`：是否命中趋势缓存
- `meta.computeTime`：生成这条曲线耗时多久
- `meta.eventCount`：趋势接口目前会带上这个字段，但它不承载事件数量语义，通常是 `0`

**各粒度返回点数**：

| 粒度 | 数据点数 | 说明 |
|------|----------|------|
| `hour` | 24 | 当天 00:00-23:00 每小时 |
| `day` | 28-31 | 当月每天中午 |
| `week` | 4-5 | 当月每周一（及月初） |
| `month` | 12 | 当年每月 15 号 |
| `year` | 5 | 前后各 2 年（共 5 年） |

**Label 格式**：

| 粒度 | 英文 | 中文 |
|------|------|------|
| `hour` | "00:00", "01:00"... | 同左 |
| `day` | "1", "2"... | 同左 |
| `week` | "W1", "W2"... | 同左 |
| `month` | "Jan", "Feb"... | "1月", "2月"... |
| `year` | "2024", "2025"... | 同左 |

**前端使用建议**：

```javascript
// 并行请求主数据和趋势数据
const [mainData, trendData] = await Promise.all([
  fetch('/api/v2/astro', { method: 'POST', body: JSON.stringify({
    birth: birthData,
    time: queryTime,
    granularity: 'year',
    language: 'zh'
  })}),
  fetch('/api/v2/astro/trend', { method: 'POST', body: JSON.stringify({
    birth: birthData,
    start_time: queryTime,
    granularity: 'year',
    language: 'zh'
  })})
]);

// 主数据先渲染（分数、事件、指导）
// 趋势图异步绘制
```

---

## 监控接口

### GET /api/monitor/dashboard

返回监控仪表板 HTML 页面。

**访问方式**：浏览器打开 `http://localhost:8888/api/monitor/dashboard`

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
| `startTime` | string | 服务启动时间 |
| `uptime` | string | 运行时长（人类可读） |
| `uptimeSeconds` | int | 运行时长（秒） |
| `totalRequests` | int | 总请求数 |
| `activeRequests` | int | 正在处理的请求数 |
| `successRequests` | int | 成功请求数 |
| `errorRequests` | int | 失败请求数 |
| `successRate` | float | 成功率（百分比） |
| `totalAPIs` | int | 已统计的接口数量 |
| `requestsLastMin` | int | 最近 1 分钟请求数 |
| `avgRequestsPerMin` | float | 平均每分钟请求数 |

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
| `path` | string | 请求路径 |
| `method` | string | 请求方法 |
| `totalRequests` | int | 总请求数 |
| `successRequests` | int | 成功请求数 |
| `errorRequests` | int | 失败请求数 |
| `avgDuration` | float | 平均响应时间（毫秒） |
| `minDuration` | int | 最快响应时间（毫秒） |
| `maxDuration` | int | 最慢响应时间（毫秒） |
| `totalDuration` | int | 总耗时（毫秒） |
| `lastAccess` | string | 最近一次访问时间 |

---

### GET /api/monitor/recent

获取最近的请求记录。

**查询参数**：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `limit` | int | 50 | 返回记录数（最多 1000） |

**说明**：

- `limit` 默认是 `50`
- 当前实现会限制为正整数，超大值由采集层控制

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

**说明**：

- `seconds` 默认是 `60`
- live 统计会按这个窗口重新聚合请求

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
    'http://localhost:8888/api/v2/astro',
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
  const response = await fetch('http://localhost:8888/api/v2/astro', {
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
curl http://localhost:8888/health

# 计算本命盘
curl -X POST http://localhost:8888/api/calc/chart \
  -H "Content-Type: application/json" \
  -d '{
    "year": 1990, "month": 5, "day": 15,
    "hour": 10, "minute": 30,
    "latitude": 39.9, "longitude": 116.4,
    "timezone": 8
  }'

# 查询五维运势（日粒度）
curl -X POST http://localhost:8888/api/v2/astro \
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
curl -X POST http://localhost:8888/api/v2/astro \
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
curl http://localhost:8888/api/monitor/summary
```

---

## 性能说明

### 主接口 `/api/v2/astro`

| 粒度 | 首次请求 | 缓存命中 |
|------|----------|----------|
| hour | < 200ms | < 1ms |
| day | < 500ms | < 1ms |
| week | 5-7s | < 1ms |
| month | 2-3s | < 1ms |
| year | 7-14s | < 1ms |

### 趋势接口 `/api/v2/astro/trend`

| 粒度 | 数据点 | 首次请求 | 缓存命中 |
|------|--------|----------|----------|
| hour | 24 | ~350ms | < 1ms |
| day | 28-31 | ~14s | < 1ms |
| week | 4-5 | ~50ms | < 1ms |
| month | 12 | ~6s | < 1ms |
| year | 5 | ~2s | < 1ms |

**缓存机制**：
- V2 接口内置多级缓存
- 相同查询参数会命中缓存，响应时间 < 1ms
- 缓存 TTL 根据粒度自动设置（hour 1h, day 6h, week/month 12h, year 24h）
- 首次请求后自动异步预热相邻时间段

---

## 版本历史

### 2.0.2 (2026-02-06)

- 新增趋势图接口 `/api/v2/astro/trend`：一次性返回所有趋势数据点
- Year 粒度性能优化：采样点从 12 个月减少到 4 个季度
- 添加异步缓存预热机制：首次请求后自动预热相邻时间段
- 前端可并行请求主数据和趋势数据，加速首屏加载

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
