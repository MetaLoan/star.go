# API 完整参考文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **编码**: UTF-8

---

## 1. 健康检查

### GET /health

检查服务状态。

**响应示例**：
```json
{
  "status": "ok",
  "service": "Star API (Go)",
  "version": "1.0.0",
  "features": [
    "natal-chart",
    "daily-forecast",
    "weekly-forecast",
    "life-trend",
    "profections",
    "transits",
    "progressions",
    "influence-factors",
    "user-management",
    "agent-api"
  ]
}
```

---

## 2. 计算接口

### 2.1 POST /api/calc/chart

计算本命盘。

**请求**：
```json
{
  "name": "用户名",
  "date": "1990-05-15T10:30:00Z",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "timezone": "Asia/Shanghai"
}
```

**请求字段**：
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 用户姓名 |
| date | string | 是 | ISO 8601 格式出生时间 |
| latitude | float | 是 | 出生地纬度 (-90 ~ 90) |
| longitude | float | 是 | 出生地经度 (-180 ~ 180) |
| timezone | string | 是 | 时区标识符 |

**响应**：
```json
{
  "birthData": {
    "name": "用户名",
    "date": "1990-05-15T10:30:00Z",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "timezone": "Asia/Shanghai"
  },
  "planets": [
    {
      "id": "sun",
      "name": "太阳",
      "symbol": "☉",
      "longitude": 54.321,
      "latitude": 0,
      "sign": "taurus",
      "signName": "金牛座",
      "signSymbol": "♉",
      "signDegree": 24.321,
      "retrograde": false,
      "house": 10,
      "dignityScore": 0
    }
    // ... 其他行星
  ],
  "houses": [
    {
      "house": 1,
      "longitude": 120.5,
      "sign": "leo",
      "signName": "狮子座"
    }
    // ... 其他宫位
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
      "strength": 0.81,
      "weight": 4.86,
      "interpretation": "太阳与月亮形成三分相"
    }
    // ... 其他相位
  ],
  "patterns": ["大三角"],
  "elementBalance": {
    "fire": 25.5,
    "earth": 35.2,
    "air": 20.1,
    "water": 19.2
  },
  "modalityBalance": {
    "cardinal": 30.0,
    "fixed": 45.0,
    "mutable": 25.0
  },
  "dominantPlanets": ["venus", "saturn"],
  "chartRuler": "sun"
}
```

### 2.2 POST /api/calc/daily

计算每日预测。

**请求**：
```json
{
  "birthData": {
    "name": "用户名",
    "date": "1990-05-15T10:30:00Z",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "timezone": "Asia/Shanghai"
  },
  "date": "2026-01-04",
  "withFactors": true
}
```

**请求字段**：
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| birthData | object | 是 | 出生数据 |
| date | string | 否 | 预测日期，默认今天 |
| withFactors | bool | 否 | 是否包含影响因子，默认 true |

**响应**：
```json
{
  "date": "2026-01-04T00:00:00Z",
  "dayOfWeek": "Saturday",
  "overallScore": 72,
  "overallTheme": "沟通与学习",
  "dimensions": {
    "career": 68,
    "relationship": 75,
    "health": 70,
    "finance": 65,
    "spiritual": 78
  },
  "moonPhase": {
    "phase": "waxingCrescent",
    "name": "盈凸月",
    "angle": 95.5,
    "illumination": 0.65,
    "voidStart": "2026-01-04T14:30:00Z",
    "voidEnd": "2026-01-04T18:45:00Z"
  },
  "moonSign": {
    "sign": "gemini",
    "name": "双子座",
    "keywords": ["沟通", "学习", "多变"]
  },
  "hourlyBreakdown": [
    {
      "hour": 0,
      "score": 65,
      "planetaryHour": "saturn",
      "bestFor": ["冥想", "计划"]
    }
    // ... 24小时
  ],
  "activeAspects": [
    {
      "planet1": "moon",
      "planet2": "venus",
      "aspectType": "trine",
      "orb": 2.1,
      "interpretation": "情感和谐，适合社交"
    }
  ],
  "factors": {
    "factors": [...],
    "totalAdjustment": 5.2,
    "positiveFactors": [...],
    "negativeFactors": [...]
  },
  "topFactors": [
    {
      "id": "dignity_exaltation_0",
      "type": "dignity",
      "name": "金星旺相",
      "adjustment": 3,
      "dimension": "relationship",
      "reason": "金星在双鱼座旺相"
    }
  ]
}
```

### 2.3 POST /api/calc/weekly

计算每周预测。

**请求**：
```json
{
  "birthData": {...},
  "date": "2026-01-04",
  "withFactors": true
}
```

**响应**：
```json
{
  "startDate": "2026-01-04",
  "endDate": "2026-01-10",
  "overallScore": 68,
  "overallTheme": "关系与合作",
  "dimensions": {
    "career": 70,
    "relationship": 72,
    "health": 65,
    "finance": 68,
    "spiritual": 66
  },
  "dailySummaries": [
    {
      "date": "2026-01-04",
      "dayOfWeek": "Saturday",
      "overallScore": 72,
      "moonSign": "gemini",
      "keyTheme": "沟通"
    }
    // ... 7天
  ],
  "keyDates": [
    {
      "date": "2026-01-06",
      "event": "木星三分太阳精确",
      "significance": "high"
    }
  ],
  "bestDaysFor": {
    "career": ["2026-01-05", "2026-01-08"],
    "relationship": ["2026-01-04", "2026-01-07"],
    "health": ["2026-01-06"]
  },
  "weeklyTransits": [
    {
      "planet": "mars",
      "aspect": "square",
      "natalPlanet": "saturn",
      "peak": "2026-01-07",
      "theme": "行动受阻"
    }
  ],
  "weeklyFactors": [...],
  "dimensionTrends": {
    "career": "stable",
    "relationship": "rising",
    "health": "declining"
  }
}
```

### 2.4 POST /api/calc/life-trend

计算人生趋势。

**请求**：
```json
{
  "birthData": {...},
  "startYear": 1990,
  "endYear": 2070,
  "resolution": "yearly"
}
```

**请求字段**：
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| birthData | object | 是 | 出生数据 |
| startYear | int | 否 | 开始年份，默认出生年 |
| endYear | int | 否 | 结束年份，默认出生年+80 |
| resolution | string | 否 | 分辨率：yearly/quarterly/monthly |

**响应**：
```json
{
  "type": "yearly",
  "birthDate": "1990-05-15T10:30:00Z",
  "points": [
    {
      "date": "1990-01-01T00:00:00Z",
      "year": 1990,
      "age": 0,
      "overallScore": 55,
      "harmonious": 60,
      "challenge": 40,
      "transformation": 30,
      "dimensions": {
        "career": 50,
        "relationship": 55,
        "health": 60,
        "finance": 52,
        "spiritual": 58
      },
      "dominantPlanet": "moon",
      "profection": {
        "house": 1,
        "theme": "自我与身份",
        "lordOfYear": "sun"
      },
      "isMajorTransit": false,
      "majorTransitName": "",
      "lunarPhaseName": "新月期",
      "lunarPhaseAngle": 15.5
    },
    {
      "year": 2019,
      "age": 29,
      "isMajorTransit": true,
      "majorTransitName": "土星回归",
      // ...
    }
    // ... 更多年份
  ],
  "summary": {
    "overallTrend": "fluctuating",
    "peakYears": [1995, 2007, 2019, 2031],
    "challengeYears": [2000, 2014, 2029],
    "transformationYears": [2008, 2020, 2038]
  },
  "cycles": {
    "saturnCycles": [
      {"age": 29, "year": 2019, "description": "第一次土星回归"},
      {"age": 58, "year": 2048, "description": "第二次土星回归"}
    ],
    "jupiterCycles": [
      {"age": 12, "year": 2002, "description": "木星回归"},
      {"age": 24, "year": 2014, "description": "木星回归"},
      // ...
    ],
    "profectionCycles": [
      {"startAge": 0, "endAge": 11, "theme": "第一轮周期"},
      {"startAge": 12, "endAge": 23, "theme": "第二轮周期"},
      {"startAge": 24, "endAge": 35, "theme": "第三轮周期"}
    ]
  }
}
```

### 2.5 POST /api/calc/time-series

生成统一时间序列。

**请求**：
```json
{
  "birthData": {...},
  "start": "2026-01-01T00:00:00Z",
  "end": "2026-01-31T00:00:00Z",
  "granularity": "day"
}
```

**请求字段**：
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| birthData | object | 是 | 出生数据 |
| start | string | 否 | 开始时间 |
| end | string | 否 | 结束时间 |
| granularity | string | 否 | 粒度：hour/day/week/month/year |

**响应**：
```json
{
  "granularity": "day",
  "points": [
    {
      "time": "2026-01-01T00:00:00Z",
      "label": "01-01",
      "granularity": "day",
      "raw": {
        "value": 45.5,
        "factors": {
          "aspects": 38.2,
          "dignity": 3,
          "retrograde": -2,
          "moonSign": 4,
          "moonPhase": 2.3
        }
      },
      "display": 72.5,
      "dimensions": {
        "career": 68,
        "relationship": 75,
        "health": 70,
        "finance": 65,
        "spiritual": 78
      },
      "volatility": 5.2
    }
    // ... 更多数据点
  ],
  "stats": {
    "mean": 50.5,
    "stdDev": 15.2,
    "min": 22.3,
    "max": 85.6,
    "trend": "up",
    "volatility": 8.5
  }
}
```

### 2.6 POST /api/calc/profection

计算年限法。

**请求**：
```json
{
  "birthData": {...},
  "age": 35
}
```

**响应**：
```json
{
  "year": 2025,
  "age": 35,
  "house": 12,
  "houseName": "玄秘宫",
  "houseTheme": "内省与超越",
  "houseKeywords": ["潜意识", "灵性", "隐退", "业力"],
  "sign": "cancer",
  "signName": "巨蟹座",
  "lordOfYear": "moon",
  "lordName": "月亮",
  "lordSymbol": "☽",
  "lordNatalHouse": 4,
  "lordNatalSign": "libra",
  "description": "今年是12宫年（玄秘宫），主题是内省与超越..."
}
```

### 2.7 POST /api/calc/profection-map

计算完整年限法地图。

**请求**：
```json
{
  "birthData": {...}
}
```

**响应**：
```json
{
  "profections": [
    // 0-80岁的完整列表
  ],
  "currentYear": {...},
  "upcomingYears": [
    // 未来5年
  ],
  "cycleAnalysis": {
    "firstCycle": {"startYear": 0, "endYear": 11, "theme": "探索与学习"},
    "secondCycle": {"startYear": 12, "endYear": 23, "theme": "建立自我"},
    "thirdCycle": {"startYear": 24, "endYear": 35, "theme": "成就与责任"},
    "currentCycleNumber": 3,
    "yearsIntoCurrentCycle": 11
  }
}
```

### 2.8 POST /api/calc/transits

计算行运。

**请求**：
```json
{
  "birthData": {...},
  "startDate": "2026-01-01",
  "endDate": "2026-12-31"
}
```

**响应**：
```json
{
  "startDate": "2026-01-01",
  "endDate": "2026-12-31",
  "events": [
    {
      "date": "2026-03-15T00:00:00Z",
      "transitPlanet": "jupiter",
      "natalPlanet": "sun",
      "aspect": {...},
      "phase": "exact",
      "intensity": 85,
      "duration": {
        "start": "2026-03-01",
        "peak": "2026-03-15",
        "end": "2026-03-30"
      },
      "interpretation": {
        "theme": "扩张与机遇",
        "keywords": ["成长", "乐观", "机会"],
        "advice": "适合开始新项目"
      }
    }
  ],
  "overallScore": 65,
  "dominantThemes": ["扩张", "沟通", "转变"]
}
```

### 2.9 POST /api/calc/progressions

计算推运。

**请求**：
```json
{
  "birthData": {...},
  "targetDate": "2026-06-15"
}
```

**响应**：
```json
{
  "targetDate": "2026-06-15T00:00:00Z",
  "progressedDate": "1990-06-20T00:00:00Z",
  "daysProgressed": 36.08,
  "yearsFromBirth": 36.08,
  "planets": [
    {
      "id": "sun",
      "name": "太阳",
      "symbol": "☉",
      "natalLongitude": 54.3,
      "progressedLongitude": 90.5,
      "movement": 36.2,
      "sign": "cancer",
      "signName": "巨蟹座",
      "signChanged": true,
      "house": 11,
      "houseChanged": true
    }
  ],
  "progressedAscendant": 156.8,
  "progressedMidheaven": 66.5,
  "aspects": [...],
  "lunarPhase": {
    "phase": "gibbous",
    "name": "凸月期",
    "angle": 155.5,
    "description": "完善与准备阶段",
    "keywords": ["完善", "分析", "准备"]
  }
}
```

---

## 3. 用户管理

### 3.1 GET /api/users

获取所有用户。

**响应**：
```json
{
  "users": [
    {
      "id": "user_001",
      "name": "张三",
      "birthData": {...},
      "createdAt": "2026-01-01T00:00:00Z",
      "updatedAt": "2026-01-04T12:00:00Z"
    }
  ]
}
```

### 3.2 POST /api/users

创建用户。

**请求**：
```json
{
  "name": "张三",
  "birthData": {
    "date": "1990-05-15T10:30:00Z",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "timezone": "Asia/Shanghai"
  }
}
```

**响应**：
```json
{
  "id": "user_001",
  "name": "张三",
  "birthData": {...},
  "createdAt": "2026-01-04T12:00:00Z"
}
```

### 3.3 GET /api/users/:id

获取单个用户。

### 3.4 PUT /api/users/:id

更新用户。

### 3.5 DELETE /api/users/:id

删除用户。

### 3.6 GET /api/users/:id/forecast

获取用户预测。

**查询参数**：
| 参数 | 类型 | 说明 |
|------|------|------|
| type | string | daily/weekly/monthly |

### 3.7 GET /api/users/:id/snapshot

获取用户当前状态快照。

**响应**：
```json
{
  "user": {...},
  "currentDate": "2026-01-04T12:00:00Z",
  "age": 35,
  "dailyForecast": {...},
  "profection": {...},
  "activetransits": [...],
  "progressedChart": {...}
}
```

---

## 4. 智能体接口

### 4.1 GET /api/agent/context

获取所有用户的占星上下文。

**响应**：
```json
{
  "users": [
    {
      "id": "user_001",
      "name": "张三",
      "currentState": {
        "sunSign": "金牛座",
        "moonSign": "天秤座",
        "ascendant": "狮子座",
        "currentMoonSign": "双子座",
        "currentMoonPhase": "上弦月",
        "profectionHouse": 12,
        "profectionTheme": "内省与超越",
        "lordOfYear": "月亮",
        "majorTransitsActive": ["木星三分太阳"]
      }
    }
  ],
  "currentDate": "2026-01-04T12:00:00Z",
  "globalTransits": {
    "sunSign": "摩羯座",
    "moonSign": "双子座",
    "moonPhase": "上弦月",
    "retrogradeplanets": ["mercury"]
  }
}
```

### 4.2 POST /api/agent/query

自然语言查询。

**请求**：
```json
{
  "userId": "user_001",
  "query": "今天的运势如何？"
}
```

**响应**：
```json
{
  "response": "今天整体运势良好，综合分数72分...",
  "data": {
    "dailyForecast": {...},
    "relevantAspects": [...]
  }
}
```

---

## 5. 错误处理

### 错误格式

```json
{
  "error": "错误描述",
  "code": "ERROR_CODE",
  "details": {}
}
```

### 错误代码

| HTTP 状态码 | 错误代码 | 说明 |
|-------------|---------|------|
| 400 | BAD_REQUEST | 请求参数错误 |
| 404 | NOT_FOUND | 资源不存在 |
| 500 | INTERNAL_ERROR | 服务器内部错误 |

---

## 6. 类型定义

### BirthData

```typescript
interface BirthData {
  name?: string;
  date: string;      // ISO 8601
  latitude: number;  // -90 ~ 90
  longitude: number; // -180 ~ 180
  timezone: string;  // IANA 时区
}
```

### PlanetPosition

```typescript
interface PlanetPosition {
  id: PlanetID;
  name: string;
  symbol: string;
  longitude: number;
  latitude: number;
  sign: ZodiacID;
  signName: string;
  signSymbol: string;
  signDegree: number;
  retrograde: boolean;
  house: number;
  dignityScore: number;
}
```

### AspectData

```typescript
interface AspectData {
  planet1: PlanetID;
  planet2: PlanetID;
  aspectType: AspectType;
  exactAngle: number;
  actualAngle: number;
  orb: number;
  applying: boolean;
  strength: number;
  weight: number;
  interpretation?: string;
}
```

### TimeSeriesPoint

```typescript
interface TimeSeriesPoint {
  time: string;
  label: string;
  granularity: 'hour' | 'day' | 'week' | 'month' | 'year';
  raw: {
    value: number;
    factors: ScoreFactors;
  };
  display: number;
  dimensions: {
    career: number;
    relationship: number;
    health: number;
    finance: number;
    spiritual: number;
  };
  volatility: number;
}
```

