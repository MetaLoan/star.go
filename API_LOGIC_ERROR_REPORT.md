# API 逻辑错误报告

## 问题概述
前端在不同日期查询同一个天体因子时，返回的 `remainingDays` 数值存在逻辑矛盾，怀疑后端计算 `remainingDays` 时使用了服务器当前时间，而非请求参数中的 `targetDate`。

---

## 前端调用方式

### API 端点
```
POST https://nathalie-clothlike-urgently.ngrok-free.dev/api/calc/daily
```

### 请求参数结构
```typescript
{
  birthData: {
    name: "Emma",
    year: 1990,
    month: 1,
    day: 9,
    hour: 10,
    minute: 30,
    second: 0,
    latitude: 39.9042,
    longitude: 116.4074,
    timezone: 8
  },
  date: "2026-01-10",           // YYYY-MM-DD 格式
  targetDate: "2026-01-10T12:00:00+08:00",  // ISO 8601 格式
  withFactors: true
}
```

### 前端代码实现
```typescript
// 文件位置：freyav1/pages/DailyInsight.tsx (545-547行)

const date = formatDateSimple(selectedTime);           // "2026-01-10"
const targetDate = formatDateToISO(selectedTime, 8);   // "2026-01-10T12:00:00+08:00"
const forecast = await getDailyForecast(
  DEFAULT_BIRTH_DATA, 
  date, 
  targetDate, 
  true
);
```

---

## 实际测试案例

### 案例 1：查询 2026-01-10 的数据

#### 请求参数
```json
{
  "birthData": {
    "name": "Emma",
    "year": 1990,
    "month": 1,
    "day": 9,
    "hour": 10,
    "minute": 30,
    "second": 0,
    "latitude": 39.9042,
    "longitude": 116.4074,
    "timezone": 8
  },
  "date": "2026-01-10",
  "targetDate": "2026-01-10T12:00:00+08:00",
  "withFactors": true
}
```

#### 预期返回（基于截图）
```json
{
  "date": "2026-01-10",
  "overallScore": 81,
  "topFactors": [
    {
      "name": "Mars Conjunction Sun",
      "adjustment": 5.3,
      "isPositive": true,
      "remainingDays": 0.708,  // 约17小时
      "lifecycle": {
        "startTime": "2026-01-09T00:00:00Z",
        "endTime": "2026-01-11T05:00:00Z"  // 假设结束时间是1月11日凌晨5点
      }
    }
  ]
}
```

#### 前端显示结果
- **日期**：January 10, 2026
- **因子名称**：Mars Conjunction Sun +5.3
- **剩余时间**：Ends in 17 hours ✅ **正确**

---

### 案例 2：查询 2026-01-11 的数据

#### 请求参数
```json
{
  "birthData": {
    "name": "Emma",
    "year": 1990,
    "month": 1,
    "day": 9,
    "hour": 10,
    "minute": 30,
    "second": 0,
    "latitude": 39.9042,
    "longitude": 116.4074,
    "timezone": 8
  },
  "date": "2026-01-11",
  "targetDate": "2026-01-11T12:00:00+08:00",  // ⚠️ 注意：查询时间是1月11日中午
  "withFactors": true
}
```

#### 实际返回（有问题）
```json
{
  "date": "2026-01-11",
  "overallScore": 81,
  "topFactors": [
    {
      "name": "Mars Conjunction Sun",
      "adjustment": 4.3,
      "isPositive": true,
      "remainingDays": 1.x,  // 大于1天（显示为 "Ends tomorrow"）
      "lifecycle": {
        "startTime": "2026-01-09T00:00:00Z",
        "endTime": "2026-01-11T05:00:00Z"  // 结束时间应该还是1月11日凌晨5点
      }
    }
  ]
}
```

#### 前端显示结果
- **日期**：January 11, 2026
- **因子名称**：Mars Conjunction Sun +4.3
- **剩余时间**：Ends tomorrow ❌ **错误**

---

## 逻辑错误分析

### 错误点 1：remainingDays 计算错误

**预期逻辑**：
```
remainingDays = (因子结束时间 - 查询目标时间) / 24小时

案例1 (2026-01-10 12:00):
remainingDays = (2026-01-11 05:00 - 2026-01-10 12:00) / 24h
              = 17小时 / 24h
              = 0.708天
              ≈ 显示 "Ends in 17 hours" ✅

案例2 (2026-01-11 12:00):
remainingDays = (2026-01-11 05:00 - 2026-01-11 12:00) / 24h
              = -7小时 / 24h
              = -0.29天
              ≈ 显示 "Ended" 或不显示此因子 ✅
```

**实际返回**（怀疑）：
```
后端可能使用了：
remainingDays = (因子结束时间 - 服务器当前实时时间) / 24小时

假设服务器当前实时时间是 2026-01-10 08:00:
remainingDays = (2026-01-11 05:00 - 2026-01-10 08:00) / 24h
              = 21小时 / 24h
              = 0.875天

但前端查询的是 2026-01-11 12:00，此时因子已结束7小时！
```

### 错误点 2：已结束因子仍然显示

如果因子在 `2026-01-11 05:00` 结束，那么查询 `2026-01-11 12:00` 时：
- **预期行为**：
  - `remainingDays` 应该为负数或0
  - 前端判断为已结束，不显示此因子
  - 或显示 "Ended"

- **实际行为**：
  - `remainingDays` 仍然大于1
  - 显示 "Ends tomorrow"
  - 因子依然出现在列表中

---

## 后端需要修复的代码位置

### 文件：`astro/factor_lifecycle.go`

#### 错误代码（推测）
```go
func CalculateRemainingDays(factor InfluenceFactor) float64 {
    if factor.Lifecycle == nil {
        return 0
    }
    
    // ❌ 错误：使用服务器当前时间
    now := time.Now()
    remaining := factor.Lifecycle.EndTime.Sub(now).Hours() / 24
    
    return remaining
}
```

#### 正确代码
```go
func CalculateRemainingDays(factor InfluenceFactor, targetDate time.Time) float64 {
    if factor.Lifecycle == nil {
        return 0
    }
    
    // ✅ 正确：使用请求参数中的 targetDate
    remaining := factor.Lifecycle.EndTime.Sub(targetDate).Hours() / 24
    
    return remaining
}
```

### 文件：`astro/score_calculator.go`

#### 需要修改的地方
在填充 `RemainingDays` 字段时，应该传入 `targetDate` 参数：

```go
// ❌ 错误
factor.RemainingDays = CalculateRemainingDays(factor)

// ✅ 正确
factor.RemainingDays = CalculateRemainingDays(factor, targetDate)
```

---

## 验证方法

### 测试步骤
1. 查询 **2026-01-10** 的数据，记录某个因子的 `remainingDays`（例如：0.708天）
2. 查询 **2026-01-11** 的数据，检查同一因子的 `remainingDays`
3. **预期**：应该减少约1天（例如：-0.292天或0，表示已结束）
4. **实际**：如果数值增加或保持不变，则证明使用了服务器实时时间

### 正确的数据示例

假设因子在 `2026-01-11 05:00` 结束：

| 查询日期 | targetDate | remainingDays | 显示文本 | 是否正确 |
|---------|-----------|---------------|---------|---------|
| 2026-01-10 | 2026-01-10T12:00 | 0.708 | Ends in 17 hours | ✅ |
| 2026-01-11 | 2026-01-11T00:00 | 0.208 | Ends in 5 hours | ✅ |
| 2026-01-11 | 2026-01-11T06:00 | -0.042 | Ended | ✅ |
| 2026-01-11 | 2026-01-11T12:00 | -0.292 | Ended | ✅ |

### 错误的数据示例（当前问题）

| 查询日期 | targetDate | remainingDays | 显示文本 | 问题 |
|---------|-----------|---------------|---------|------|
| 2026-01-10 | 2026-01-10T12:00 | 0.708 | Ends in 17 hours | ✅ |
| 2026-01-11 | 2026-01-11T12:00 | 1.x | Ends tomorrow | ❌ 应该已结束 |

---

## 前端代码位置

### 接口定义
**文件**：`freyav1/services/starApiService.ts`
**行数**：101-114

```typescript
topFactors?: Array<{
  id: string;
  type: string;
  name: string;
  description: string;
  adjustment: number;
  isPositive: boolean;
  remainingDays?: number;  // 后端返回的剩余天数
  lifecycle?: {
    startTime: string;
    peakTime: string;
    endTime: string;
    duration: number;
  };
}>;
```

### 显示逻辑
**文件**：`freyav1/pages/DailyInsight.tsx`
**行数**：1352-1390

```typescript
// 显示剩余天数（支持小数）
if (factor.remainingDays !== undefined && factor.remainingDays > 0) {
  const days = Math.floor(factor.remainingDays);
  const hours = Math.round((factor.remainingDays - days) * 24);
  
  if (days === 0 && hours < 24) {
    // 小于1天，显示小时
    return `Ends in ${hours} hours`;
  } else if (days === 1) {
    return 'Ends tomorrow';
  } else {
    return `Ends in ${days} days`;
  }
} else if (factor.remainingDays === 0 || factor.remainingDays < 0) {
  return 'Ended';
}
```

---

## 总结

### 问题根因
后端在计算 `remainingDays` 时，使用了 **服务器当前实时时间**，而非 API 请求中传入的 **targetDate** 参数。

### 影响范围
- 用户查询历史日期时，显示的剩余天数不准确
- 用户查询未来日期时，显示的剩余天数也不准确
- 已结束的因子在查询历史日期时仍然显示

### 修复建议
1. 修改 `CalculateRemainingDays()` 函数，接收 `targetDate` 参数
2. 修改 `score_calculator.go`，在计算时传入请求的 `targetDate`
3. 确保所有时间计算都使用请求参数中的时间，而非 `time.Now()`

### 验证方法
修复后，查询不同日期时，`remainingDays` 应该呈现递减趋势，直到为0或负数（表示已结束）。
