# 架构重构待办事项

## Event/Factor 架构问题

### 问题描述

当前系统中，同一个天象事件被两套系统分别处理：

| 系统 | 文件 | 用途 |
|------|------|------|
| DailyEvent | `astro/daily_events.go` | 日历视图：展示精确时间点发生的事件 |
| Factor | `astro/score_calculator.go` | 评分计算：计算影响强度和生命周期 |

### 具体问题

1. **重复计算**：同一个相位被两个系统分别计算
2. **数据不一致风险**：Event 和 Factor 可能对同一事件有不同的判断
3. **覆盖不一致**：
   - `voidOfCourse` 只在 Factor 系统中，没有对应的 DailyEvent
   - `planetary_hour_change` 是 DailyEvent 类型，Factor 中命名为 `planetaryHour`
4. **命名不统一**：相同概念使用不同的 eventType 名称

### unified-events 的临时方案

`/api/calc/unified-events` 接口试图合并两者：
1. 从 DailyEvent 获取精确时间事件
2. 从 Factor 获取影响数据
3. 匹配后合并，未匹配的 Factor 单独添加为事件

但这只是**表面统一**，底层仍是两套计算系统。

### 建议的重构方向

1. **统一事件模型**：定义一个统一的 `AstroEvent` 结构
2. **单一数据源**：所有天象计算从同一个源头产生
3. **视图分离**：日历视图和评分视图从同一数据派生不同展示
4. **命名规范**：统一 eventType 命名（建议使用 snake_case）

### 受影响的文件

- `backend/astro/daily_events.go`
- `backend/astro/score_calculator.go`
- `backend/api/unified_events_handlers.go`
- `backend/i18n/detailed_interpretations.go`
- `backend/i18n/emotional_titles.go`

### 优先级

中 - 当前系统可用，但长期维护成本较高

---

*创建日期：2026-02-05*
