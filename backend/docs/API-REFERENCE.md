# Star 后端 API 参考文档

本文档详细介绍了 Star 占星计算平台后端的 API 接口。

## 基础信息

- **Base URL**: `http://localhost:8080` (本地开发)
- **Content-Type**: `application/json`
- **时区说明**: API 支持通过 `timezone` 参数指定时区（如北京时间为 8），返回的时间戳通常带有时区偏移。

---

## 健康检查

### 服务状态
- **URL**: `/health`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "status": "ok",
    "service": "Star API (Go)",
    "version": "1.0.0",
    "dataSource": "Swiss Ephemeris",
    "features": ["natal-chart", "daily-forecast", "weekly-forecast", "life-trend", "profections", "transits", "progressions", "influence-factors", "user-management", "agent-api"]
  }
  ```

---

## 核心数据结构

### BirthData (出生数据)
用于几乎所有计算接口的基础请求数据。

```json
{
  "name": "张三",
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

### DimensionScores (五维度分数)
所有预测/时间序列接口返回的维度数据结构：

```json
{
  "career": 70.5,       // 事业
  "relationship": 65.3, // 关系
  "health": 72.1,       // 健康
  "finance": 68.4,      // 财务
  "spiritual": 64.7     // 灵性
}
```

**⚠️ 字段名映射**：
| 后端字段 | 中文含义 | 说明 |
|---------|---------|------|
| `career` | 事业 | 工作、职业发展 |
| `relationship` | 关系 | 人际、感情、合作 |
| `health` | 健康 | 身体、精力状态 |
| `finance` | 财务 | 财富、收入 |
| `spiritual` | 灵性 | 内在成长、直觉 |

---

## 计算 API (`/api/calc`)

### 1. 计算本命盘
- **URL**: `/api/calc/chart`
- **Method**: `POST`
- **Request**: `BirthData`
- **Response**:
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
      { "house": 1, "cusp": 120.5, "sign": "leo", "signName": "狮子座" }
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
    "elementBalance": { "fire": 0.3, "earth": 0.2, "air": 0.35, "water": 0.15 },
    "modalityBalance": { "cardinal": 0.4, "fixed": 0.3, "mutable": 0.3 },
    "dominantPlanets": ["sun", "mars"],
    "chartRuler": "sun"
  }
  ```

### 2. 每日预测
- **URL**: `/api/calc/daily`
- **Method**: `POST`
- **Request**: 
  ```json
  {
    "birthData": { ... },
    "date": "2026-01-06",
    "targetDate": "2026-01-06T12:00:00+08:00",
    "withFactors": true
  }
  ```
- **Response**:
  ```json
  {
    "date": "2026-01-06T00:00:00+08:00",
    "dayOfWeek": "Tuesday",
    "overallScore": 72.5,
    "overallTheme": "适合沟通与学习",
    "dimensions": {
      "career": 70.5,
      "relationship": 65.3,
      "health": 72.1,
      "finance": 68.4,
      "spiritual": 64.7
    },
    "moonPhase": {
      "phase": "waxing_crescent",
      "name": "蛾眉月",
      "angle": 45.2,
      "illumination": 0.25
    },
    "moonSign": {
      "sign": "gemini",
      "name": "双子座",
      "keywords": ["沟通", "学习", "灵活"]
    },
    "hourlyBreakdown": [
      { "hour": 0, "score": 68.5, "planetaryHour": "saturn", "bestFor": ["冥想", "反思"] },
      { "hour": 1, "score": 70.2, "planetaryHour": "jupiter", "bestFor": ["学习", "规划"] }
    ],
    "topFactors": [
      {
        "id": "transit_jupiter_sun",
        "type": "transit",
        "name": "木星合太阳",
        "description": "扩张与机遇的能量",
        "adjustment": 5.2,
        "isPositive": true
      }
    ]
  }
  ```

### 3. 每周预测
- **URL**: `/api/calc/weekly`
- **Method**: `POST`
- **Request**: 
  ```json
  {
    "birthData": { ... },
    "date": "2026-01-06",
    "withFactors": true
  }
  ```
- **Response**:
  ```json
  {
    "startDate": "2026-01-06T00:00:00+08:00",
    "endDate": "2026-01-12T23:59:59+08:00",
    "overallScore": 68.5,
    "overallTheme": "沟通与合作的一周",
    "dimensions": {
      "career": 72.0,
      "relationship": 65.5,
      "health": 70.0,
      "finance": 66.0,
      "spiritual": 62.5
    },
    "dailySummaries": [
      {
        "date": "2026-01-06",
        "dayOfWeek": "Monday",
        "overallScore": 72.5,
        "moonSign": "gemini",
        "keyTheme": "适合沟通"
      }
    ],
    "keyDates": [
      { "date": "2026-01-08", "event": "水星入摩羯", "significance": "思维变得务实" }
    ],
    "bestDaysFor": {
      "work": ["2026-01-07", "2026-01-10"],
      "relationship": ["2026-01-09"],
      "health": ["2026-01-06"]
    },
    "weeklyTransits": [
      {
        "planet": "venus",
        "aspect": "trine",
        "natalPlanet": "mars",
        "peak": "2026-01-08",
        "theme": "创造力与热情"
      }
    ]
  }
  ```

### 4. 人生趋势 (80年)
- **URL**: `/api/calc/life-trend`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "startYear": 0,
    "endYear": 80,
    "resolution": "yearly"
  }
  ```
- **Response**:
  ```json
  {
    "type": "lifeTrend",
    "birthDate": "1990-06-15T00:00:00+08:00",
    "points": [
      {
        "date": "1990-06-15T00:00:00+08:00",
        "year": 1990,
        "age": 0,
        "overallScore": 65.0,
        "harmonious": 55.0,
        "challenge": 30.0,
        "transformation": 20.0,
        "dimensions": {
          "career": 60.0,
          "relationship": 70.0,
          "health": 75.0,
          "finance": 55.0,
          "spiritual": 50.0
        },
        "dominantPlanet": "moon",
        "profection": { "house": 1, "theme": "自我", "lordOfYear": "mars" },
        "isMajorTransit": false,
        "lunarPhaseName": "新月",
        "lunarPhaseAngle": 0
      }
    ],
    "summary": {
      "overallTrend": "上升",
      "peakYears": [28, 35, 42],
      "challengeYears": [29, 44],
      "transformationYears": [30, 51]
    },
    "cycles": {
      "saturnCycles": [
        { "age": 29, "year": 2019, "description": "土星回归" },
        { "age": 58, "year": 2048, "description": "第二次土星回归" }
      ]
    }
  }
  ```

### 5. 统一时间序列 ⭐️ (图表核心接口)
- **URL**: `/api/calc/time-series`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "start": "2026-01-01T00:00:00+08:00",
    "end": "2026-01-31T23:59:59+08:00",
    "granularity": "day"
  }
  ```
- **Response**:
  ```json
  {
    "granularity": "day",
    "points": [
      {
        "time": "2026-01-01T00:00:00+08:00",
        "label": "01-01",
        "granularity": "day",
        "raw": {
          "value": 75.5,
          "factors": {
            "aspects": 2.5,
            "dignity": 1.0,
            "retrograde": -0.5
          }
        },
        "display": 68.2,
        "dimensions": {
          "career": 70.5,
          "relationship": 65.3,
          "health": 72.1,
          "finance": 68.4,
          "spiritual": 64.7
        },
        "volatility": 2.3
      }
    ],
    "stats": {
      "mean": 68.5,
      "stdDev": 5.2,
      "min": 55.0,
      "max": 85.0,
      "trend": "up",
      "volatility": 7.6
    }
  }
  ```

**⚠️ 前端解析要点**：
- 数据在 `points` 数组中（不是 `data`）
- 综合分数使用 `display` 字段
- 五维度在 `dimensions` 对象中：
  - `dimensions.career` (事业)
  - `dimensions.relationship` (关系)
  - `dimensions.health` (健康)
  - `dimensions.finance` (财务)
  - `dimensions.spiritual` (灵性)
- **没有 `scores`、`overall`、`love`、`wealth`、`learning`、`social` 等字段**

### 6. 年限法 (Profection)
- **URL**: `/api/calc/profection`
- **Method**: `POST`
- **Request**: `{ "birthData": { ... }, "age": 35 }`
- **Response**:
  ```json
  {
    "year": 2025,
    "age": 35,
    "house": 12,
    "houseName": "第十二宫",
    "houseTheme": "内省与灵性",
    "houseKeywords": ["隐退", "潜意识", "疗愈"],
    "sign": "cancer",
    "signName": "巨蟹座",
    "lordOfYear": "moon",
    "lordName": "月亮",
    "lordSymbol": "☽",
    "lordNatalHouse": 4,
    "lordNatalSign": "libra",
    "description": "今年聚焦内在成长与情感疗愈"
  }
  ```

### 7. 年限法地图 (Profection Map)
- **URL**: `/api/calc/profection-map`
- **Method**: `POST`
- **Request**: `{ "birthData": { ... } }`
- **Response**:
  ```json
  {
    "profections": [ /* 0-84岁的年限法数据 */ ],
    "currentYear": { /* 当前年限法 */ },
    "upcomingYears": [ /* 未来5年 */ ],
    "cycleAnalysis": {
      "currentCycleNumber": 3,
      "yearsIntoCurrentCycle": 11
    }
  }
  ```

### 8. 行运计算 (Transits)
- **URL**: `/api/calc/transits`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "startDate": "2026-01-01",
    "endDate": "2026-12-31"
  }
  ```
- **Response**:
  ```json
  {
    "startDate": "2026-01-01",
    "endDate": "2026-12-31",
    "events": [
      {
        "date": "2026-03-15T00:00:00+08:00",
        "transitPlanet": "saturn",
        "natalPlanet": "sun",
        "aspect": { "aspectType": "square", "orb": 0.5 },
        "phase": "exact",
        "intensity": 0.95,
        "duration": { "start": "2026-02-01", "peak": "2026-03-15", "end": "2026-04-30" },
        "interpretation": { "theme": "责任与挑战", "keywords": ["结构", "限制"], "advice": "耐心面对" }
      }
    ],
    "overallScore": 65.0,
    "dominantThemes": ["结构重建", "责任"]
  }
  ```

### 9. 推运计算 (Progressions)
- **URL**: `/api/calc/progressions`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "targetDate": "2026-01-06"
  }
  ```
- **Response**:
  ```json
  {
    "targetDate": "2026-01-06T00:00:00+08:00",
    "progressedDate": "1990-08-20T00:00:00+08:00",
    "daysProgressed": 35.5,
    "yearsFromBirth": 35.5,
    "planets": [
      {
        "id": "sun",
        "name": "太阳",
        "natalLongitude": 84.5,
        "progressedLongitude": 119.8,
        "movement": 35.3,
        "sign": "leo",
        "signChanged": true,
        "house": 11,
        "houseChanged": true
      }
    ],
    "progressedAscendant": 155.2,
    "progressedMidheaven": 65.8,
    "lunarPhase": {
      "phase": "full_moon",
      "name": "满月",
      "angle": 180,
      "description": "收获与完成的阶段"
    }
  }
  ```

### 10. 月亮空亡 (Void of Course)
- **URL**: `/api/calc/void-of-course`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "date": "2026-01-06T12:00:00",
    "latitude": 39.9042,
    "longitude": 116.4074
  }
  ```
- **Response**:
  ```json
  {
    "isVoid": true,
    "startTime": "2026-01-06T10:30:00+08:00",
    "endTime": "2026-01-06T14:45:00+08:00",
    "duration": 4.25,
    "nextSign": "cancer",
    "lastAspect": "trine to jupiter",
    "influence": -0.3
  }
  ```

### 11. 行星时 (Planetary Hour)
- **URL**: `/api/calc/planetary-hour`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "date": "2026-01-06T14:00:00",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "fullDay": false
  }
  ```
- **Response (单个行星时)**:
  ```json
  {
    "planetaryHour": 7,
    "ruler": "venus",
    "planetName": "金星",
    "planetSymbol": "♀",
    "dayRuler": "tuesday",
    "influence": 0.8,
    "bestFor": ["艺术", "社交", "美容"]
  }
  ```
- **Response (fullDay=true)**:
  ```json
  {
    "date": "2026-01-06",
    "hours": [
      { "planetaryHour": 1, "ruler": "mars", "planetName": "火星", "startTime": "06:45", "endTime": "07:45" },
      { "planetaryHour": 2, "ruler": "sun", "planetName": "太阳", "startTime": "07:45", "endTime": "08:45" }
    ]
  }
  ```

### 12. 分值组成详情 (单粒度，调试用)
- **URL**: `/api/calc/score-breakdown`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "queryTime": "2026-01-06T12:00:00+08:00",
    "granularity": "hour",
    "userId": "optional"
  }
  ```
- **Response**:
  ```json
  {
    "queryTime": "2026-01-06T12:00:00+08:00",
    "granularity": "hour",
    "overallScore": 72.5,
    "overallRaw": 85.2,
    "dimensions": [
      {
        "dimension": "career",
        "baseScore": 50,
        "aspectScore": 5.2,
        "factorScore": 15.3,
        "rawScore": 70.5,
        "finalScore": 70.5,
        "factors": [ /* 该维度相关因子 */ ]
      }
    ],
    "factorsByLevel": {
      "hourly": [ /* 小时级因子 */ ],
      "daily": [ /* 日级因子 */ ],
      "weekly": [ /* 周级因子 */ ]
    },
    "meta": {
      "dataSource": "Swiss Ephemeris",
      "visibleLevels": ["hourly", "daily", "weekly", "monthly", "yearly"],
      "totalFactorCount": 15,
      "positiveFactors": 10,
      "negativeFactors": 5
    }
  }
  ```

### 13. 分值组成详情 (多粒度，调试用)
- **URL**: `/api/calc/score-breakdown-all`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "queryTime": "2026-01-06T12:00:00+08:00",
    "userId": "optional"
  }
  ```
- **Response**:
  ```json
  {
    "queryTime": "2026-01-06T12:00:00+08:00",
    "breakdown": {
      "hour": { /* ScoreBreakdown */ },
      "day": { /* ScoreBreakdown */ },
      "week": { /* ScoreBreakdown */ },
      "month": { /* ScoreBreakdown */ },
      "year": { /* ScoreBreakdown */ }
    }
  }
  ```

### 14. 时间范围内活跃因子
- **URL**: `/api/calc/active-factors`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "queryTime": "2026-01-06",
    "granularity": "week",
    "infect": "core",
    "userId": "optional"
  }
  ```
- **Response**:
  ```json
  {
    "granularity": "week",
    "rangeStart": "2026-01-06T00:00:00+08:00",
    "rangeEnd": "2026-01-12T23:59:59+08:00",
    "infect": "core",
    "factors": [
      {
        "id": "transit_jupiter_sun",
        "name": "木星合太阳",
        "type": "transit",
        "timeLevel": "monthly",
        "baseValue": 5.0,
        "weight": 1.2,
        "description": "扩张与机遇",
        "isPositive": true,
        "startTime": "2026-01-01T00:00:00+08:00",
        "peakTime": "2026-01-08T12:00:00+08:00",
        "endTime": "2026-01-15T00:00:00+08:00",
        "maxStrength": 0.95
      }
    ],
    "totalCount": 12,
    "positiveCount": 8,
    "negativeCount": 4
  }
  ```

### 15. 分数解释 (面向C端用户)
- **URL**: `/api/calc/score-explain`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "birthData": { ... },
    "queryTime": "2026-01-06T12:00:00+08:00",
    "granularity": "day",
    "dimension": "career",
    "userId": "optional"
  }
  ```
- **Response**:
  ```json
  {
    "score": 72.5,
    "dimension": "career",
    "summary": "今日事业运势良好",
    "mainInfluences": [
      {
        "name": "木星三分太阳",
        "effect": "正面",
        "description": "带来扩张机遇，适合推进项目"
      },
      {
        "name": "水星逆行",
        "effect": "需注意",
        "description": "沟通可能出现误解，建议多确认"
      }
    ],
    "suggestions": [
      "适合进行重要会议",
      "避免签署重要合同",
      "下午3-5点能量最佳"
    ]
  }
  ```

---

## 用户管理 API (`/api/users`)

### 1. 获取所有用户
- **URL**: `/api/users`
- **Method**: `GET`
- **Response**: 
  ```json
  { "users": [ { "id": "...", "name": "...", "birthData": {...} } ] }
  ```

### 2. 创建用户
- **URL**: `/api/users`
- **Method**: `POST`
- **Request**: `{ "name": "Jack", "birthData": { ... } }`
- **Response**: 
  ```json
  { "id": "uuid", "name": "Jack", "birthData": {...}, "createdAt": "..." }
  ```

### 3. 获取单个用户
- **URL**: `/api/users/:id`
- **Method**: `GET`
- **Response**: 用户详情（含 BirthData 和 NatalChart）

### 4. 更新用户
- **URL**: `/api/users/:id`
- **Method**: `PUT`
- **Request**: `{ "name": "Jack Updated", "birthData": { ... } }`

### 5. 删除用户
- **URL**: `/api/users/:id`
- **Method**: `DELETE`
- **Response**: `204 No Content`

### 6. 获取用户预测
- **URL**: `/api/users/:id/forecast`
- **Method**: `GET`
- **Query Params**: `type=daily` 或 `type=weekly`
- **Response**: 同 `/api/calc/daily` 或 `/api/calc/weekly`

### 7. 获取用户快照 (Snapshot)
- **URL**: `/api/users/:id/snapshot`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "user": { "id": "...", "name": "...", "birthData": {...} },
    "currentDate": "2026-01-06T12:00:00+08:00",
    "age": 35,
    "dailyForecast": { /* DailyForecast */ },
    "profection": { /* AnnualProfection */ },
    "activeTransits": [ /* TransitEvent[] */ ],
    "progressedChart": { /* ProgressedChart */ }
  }
  ```

---

## 智能体接口 (`/api/agent`)

### 1. 获取全局上下文
- **URL**: `/api/agent/context`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "currentDate": "2026-01-06T12:00:00+08:00",
    "globalTransits": {
      "sunSign": "capricorn",
      "moonSign": "gemini",
      "moonPhase": "waxing_crescent",
      "retrogradePlanets": ["mercury"]
    },
    "users": [
      { "id": "...", "name": "...", "currentState": { "todayScore": 72 } }
    ]
  }
  ```

### 2. 智能体查询
- **URL**: `/api/agent/query`
- **Method**: `POST`
- **Request**: `{ "userId": "...", "query": "我今天运气怎么样？" }`
- **Response**:
  ```json
  {
    "response": "今天整体运势不错，综合评分72分...",
    "data": { "score": 72, "dimension": {...} }
  }
  ```

---

## 运营与配置 API (`/api/admin`)

### 1. 因子权重管理
- **GET**: `/api/admin/factor-weights` - 获取当前因子权重配置
- **PUT**: `/api/admin/factor-weights` - 更新因子权重
- **Request/Response**:
  ```json
  {
    "dignity": 1.0,
    "retrograde": 1.0,
    "aspectPhase": 0.8,
    "aspectOrb": 0.5,
    "outerPlanet": 1.2,
    "profectionLord": 1.0,
    "lunarPhase": 0.7,
    "planetaryHour": 0.3,
    "voidOfCourse": 0.5,
    "personal": 1.0,
    "custom": 1.0
  }
  ```

### 2. 维度权重管理
- **GET**: `/api/admin/dimension-weights` - 获取当前维度权重配置
- **PUT**: `/api/admin/dimension-weights` - 更新维度权重
- **Request/Response**:
  ```json
  {
    "career": 0.25,
    "relationship": 0.20,
    "health": 0.20,
    "finance": 0.20,
    "spiritual": 0.15
  }
  ```
- **Note**: 所有权重之和必须为 1.0。

### 3. 抖动配置管理
- **GET**: `/api/admin/jitter-config`
- **PUT**: `/api/admin/jitter-config`
- **Request/Response**:
  ```json
  {
    "enabled": true,
    "magnitude": 0.5,
    "seed": 0
  }
  ```

### 4. 添加自定义因子
- **URL**: `/api/admin/custom-factors`
- **Method**: `POST`
- **Request**:
  ```json
  {
    "userId": "default",
    "definition": "AddScore=(2*healthScore,2.5,202501171230)"
  }
  ```
- **Response**:
  ```json
  {
    "message": "自定义因子已添加",
    "factor": {
      "id": "custom_1704556800",
      "type": "custom",
      "operation": "AddScore",
      "value": 2.0,
      "targetDimension": "health",
      "duration": 2.5,
      "startTime": "2025-01-17T12:30:00+08:00"
    }
  }
  ```
- **公式格式**: `Operation=(value,duration,startTime)`
  - **操作**: `AddScore` | `SubScore` | `MulScore` | `SetScore`
  - **值**: 数值（可乘以维度如 `2*healthScore`）
  - **持续时长**: 小时数
  - **开始时间**: `YYYYMMDDHHmm` 格式

### 5. 获取用户自定义因子
- **URL**: `/api/admin/custom-factors/:userId`
- **Method**: `GET`
- **Response**: `{ "userId": "...", "count": 2, "factors": [...] }`

### 6. 清除用户自定义因子
- **URL**: `/api/admin/custom-factors/:userId`
- **Method**: `DELETE`
- **Response**: `{ "message": "自定义因子已清除", "userId": "..." }`

---

## 错误响应

所有接口在参数错误或服务异常时返回统一格式：

```json
{
  "error": "错误描述信息"
}
```

常见 HTTP 状态码：
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务内部错误

---

## 字段名速查表

| 用途 | 正确字段名 | ❌ 错误示例 |
|------|-----------|------------|
| 综合分数 | `display` 或 `overallScore` | `overall`, `score`, `value` |
| 时间序列数据 | `points` | `data`, `items` |
| 五维度对象 | `dimensions` | `scores`, `dims` |
| 事业维度 | `dimensions.career` | `scores.career`, `career` |
| 关系维度 | `dimensions.relationship` | `scores.love`, `love` |
| 健康维度 | `dimensions.health` | `scores.health` |
| 财务维度 | `dimensions.finance` | `scores.wealth`, `wealth` |
| 灵性维度 | `dimensions.spiritual` | `scores.spiritual` |
