# Star 后端算法说明

本文档面向技术同学，说明当前后端的底层计算链路、数据结构、评分模型、事件生成、时间聚合和缓存策略。

如果你不是占星背景，可以先记住一句话：

> 这套系统不是在“算命”，而是在把一组天文状态，翻译成可解释的时间评分和事件提示。

## 1. 总体架构

后端不是单一“打分函数”，而是一条完整的数据流：

1. 由出生数据生成本命盘 `NatalChart`
2. 在查询时刻计算行运位置 `TransitPositions`
3. 根据本命盘和行运位置生成影响因子 `InfluenceFactor`
4. 将因子映射到五个维度
5. 将五维原始分标准化为 0-100
6. 生成事件、指导、趋势、分数分解和时间聚合结果

主链路主要分成三层：

- **天文层**：负责行星、宫位、相位、月相、空亡、行星时等计算
- **评分层**：负责基础分、影响因子、维度分、综合分
- **编排层**：负责按小时/日/周/月/年聚合、缓存、预热、i18n 和 API 响应

### 1.1 大白话理解

- `本命盘` = 这个人的长期底子，像“出厂配置”
- `行运` = 今天、这周、这个月外部环境发生了什么
- `影响因子` = 某颗星此刻对某个生活领域的加减分理由
- `五维分数` = 把复杂天象压缩成事业、关系、健康、财务、灵性五个方向
- `时间槽` = 把某一小时/一天/一周的结果装进一个统一容器里

### 1.2 核心入口速查

下面这些是系统里最重要的入口函数。你可以把它们理解成“主按钮”。

| 入口函数 | 负责什么 | 输出给谁看 | 典型用途 |
|------|------|------|------|
| `CalculateNatalChart` | 生成本命盘 | 后续所有计算 | 先把出生数据变成一张星盘 |
| `CalculateNatalBaseScores` | 算本命底分 | 评分系统 | 看这个人长期底子怎么样 |
| `CalculateScoresV2` | 算某时刻的完整五维分 | 主接口、分数解释、人生趋势 | 单点评分的标准版本 |
| `CalculateScoresV2Lite` | 算轻量版五维分 | 趋势图、批量采样 | 快速算很多时间点 |
| `CalculateInfluenceFactorsV2` | 生成当前活跃因子 | 分数、分解、解释 | 找出“为什么会变动” |
| `CalculateScoreBreakdown` | 拆分分数来源 | 分数拆解接口 | 让用户看到每分怎么来的 |
| `CalculateDailyAstroData` | 统一计算一天的星象数据 | 日事件、因子系统 | 避免重复计算 |
| `CalculateDailyEvents` | 生成当天事件列表 | 日历/提醒/事件面板 | 直接给前端事件流 |
| `GetActiveFactorsInRange` | 查某时间范围内活跃因子 | 因子查询/排查 | 看某时刻有哪些因子还在生效 |
| `CalculateAnnualProfection` | 算年限法年度主题 | 年度趋势、人生趋势 | 看今年主线是什么 |
| `CalculateProgressions` | 算次限推运盘 | 推运分析 | 看长期成长和阶段变化 |
| `CalculateLifeTrend` | 算人生趋势点 | 人生时间轴 | 做长期轨迹图 |
| `CalculateTransits` | 算某段时间行运 | 行运分析 | 看一段日期里有哪些大事 |
| `CalculatePlanetaryHourEnhanced` | 算行星时 | 小时级提示 | 做短时段节奏建议 |
| `CalculateVoidOfCourse` | 算月亮空亡 | 小时级提示 | 提醒不要贸然开始新事 |
| `CalculateUnifiedHourlyScore` | 统一小时评分 | 老的统一分体系 | 小时级时间序列基础 |

大白话理解：

- `CalculateNatalChart` 负责“把人出生那一刻的天空拍下来”
- `CalculateScoresV2` 负责“把今天这刻的天空翻译成五维分”
- `CalculateInfluenceFactorsV2` 负责“告诉你这分数为什么这么变”
- `CalculateScoreBreakdown` 负责“把分数拆给你看”
- `CalculateLifeTrend` 负责“把很多年的变化串成一条时间线”
- `GetActiveFactorsInRange` 负责“告诉你某段时间里哪些因子正在生效”

---

## 2. 数据源与精度策略

当前系统以 Swiss Ephemeris 作为唯一权威数据源。启动时会先验证可用性，验证失败直接退出。

核心入口见：

- `backend/main.go`
- `backend/astro/ephemeris.go`
- `backend/astro/unified_ephemeris.go`

策略上分两类：

- **精确链路**：用于单点计算、事件精确时刻搜索、相位生命周期生成
- **轻量链路**：用于趋势曲线、批量采样、聚合计算，跳过精确时刻搜索，性能更高

这意味着同一个事件在不同接口里可能有不同的“计算深度”，但底层来源一致。

---

## 3. 核心数据模型

### 3.1 本命盘

`NatalChart` 由以下部分组成：

- 行星位置 `Planets`
- 宫位 `Houses`
- 上升点 `Ascendant`
- 中天 `Midheaven`
- 相位 `Aspects`
- 格局 `Patterns`
- 元素平衡 `ElementBalance`
- 模式平衡 `ModalityBalance`
- 主导行星 `DominantPlanets`
- 命主星 `ChartRuler`

生成入口：

- `backend/astro/natal_chart.go`

### 3.2 时间槽

`TimeSlot` 是统一时间单元，承载：

- 当前时间段
- 五维分数
- 事件列表
- 与上一周期的 `Delta`
- `Guidance`
- 子时间槽 `SubSlots`

相关定义：

- `backend/core/timeslot.go`
- `backend/models/types.go`

### 3.3 影响因子

`InfluenceFactor` 是评分系统的关键中间层，字段包含：

- `Type`
- `Name`
- `Lifecycle`
- `BaseValue`
- `Weight`
- `CurrentStrength`
- `Adjustment`
- `DimensionImpact`
- `SourcePlanet`
- `IsPositive`

这套模型把占星事件统一成可加权、可归类、可衰减的数值对象。

### 3.4 字段总览

下面这些字段是前后端最常见的“业务字段”，先看这一节，后面读接口和算法会轻松很多。

#### BirthData

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `name` | 用户名字 | 这次计算对应谁 |
| `year` | 出生年 | 出生的年份 |
| `month` | 出生月 | 出生的月份 |
| `day` | 出生日 | 出生的日期 |
| `hour` | 出生小时 | 出生时的钟点 |
| `minute` | 出生分钟 | 出生时的分钟 |
| `second` | 出生秒 | 更精细的出生时间 |
| `latitude` | 纬度 | 出生地南北位置，影响宫位 |
| `longitude` | 经度 | 出生地东西位置，影响宫位 |
| `timezone` | 时区偏移 | 出生地相对 UTC 的时间差 |

#### NatalChart

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `birthData` | 出生信息 | 这张盘是基于谁的出生时间算的 |
| `planets` | 行星位置 | 每颗星现在落在哪、状态如何 |
| `houses` | 宫位 | 12 个生活场景的分区 |
| `ascendant` | 上升点 | 出生时东方地平线上升起的点 |
| `midheaven` | 中天 | 事业/公众面向的顶部位置 |
| `aspects` | 本命相位 | 行星之间的配合或冲突关系 |
| `patterns` | 图形格局 | 多颗星组成的结构模式 |
| `elementBalance` | 元素平衡 | 火土风水哪类能量更强 |
| `modalityBalance` | 模式平衡 | 本位/固定/变动哪类更强 |
| `dominantPlanets` | 主导行星 | 影响最明显的星 |
| `chartRuler` | 命主星 | 上升星座的守护星 |

#### TimeSlot

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `userID` | 用户标识 | 这是谁的时间槽 |
| `startTime` | 开始时间 | 这个时间段从什么时候开始 |
| `endTime` | 结束时间 | 这个时间段什么时候结束 |
| `granularity` | 粒度 | 是小时、日、周、月还是年 |
| `scores` | 五维分数 | 这段时间的综合评分 |
| `events` | 事件列表 | 这段时间里有哪些天象事件 |
| `delta` | 变化量 | 跟上一周期比变好还是变差 |
| `guidance` | 指导建议 | 这段时间适合做什么、避开什么 |
| `subSlots` | 子时间槽 | 用来画曲线的下一级采样点 |

#### DimensionScores

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `overall` | 综合分 | 总体顺不顺 |
| `career` | 事业 | 工作、职业、成就 |
| `relationship` | 关系 | 感情、人际、合作 |
| `health` | 健康 | 身体、精力、恢复 |
| `finance` | 财务 | 收入、支出、资源 |
| `spiritual` | 灵性 | 内在成长、直觉、意义感 |

#### AstroEvent

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `EventID` | 事件 ID | 这条事件的唯一编号，用于去重 |
| `Type` | 类型 | 是相位、空亡、逆行还是别的事件 |
| `Title` | 标题 | 给前端看的简短标题 |
| `IsPositive` | 正负 | 这条事件是利好还是压力 |
| `Intensity` | 强度 | 这条事件有多强 |
| `TimeLevel` | 时间级别 | 这条事件是小时级还是年级 |
| `PrimaryPlanet` | 主星体 | 主要影响谁 |
| `SecondaryPlanet` | 次星体 | 相位事件里被影响的另一颗星 |
| `Aspect` | 相位类型 | 合相、三分、四分等 |
| `Impact` | 绝对影响 | 这条事件对五维的基础影响 |
| `ImpactDelta` | 相对变化 | 跟上一周期比增加还是减少 |
| `StartTime` | 开始时间 | 这条事件从什么时候开始有效 |
| `EndTime` | 结束时间 | 这条事件什么时候结束 |
| `ExactTime` | 精确时间 | 最强、最准确的那一刻 |
| `Interpretation` | 解读 | 解释这条事件代表什么 |
| `Advice` | 建议 | 建议怎么用这段时间 |

#### InfluenceFactor

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `ID` | 因子 ID | 这条原因的唯一编号 |
| `Type` | 因子类型 | 是逆行、相位、空亡还是别的 |
| `Name` | 因子名称 | 例如“Mercury Retrograde” |
| `Description` | 因子描述 | 更完整的说明 |
| `TimeLevel` | 时间级别 | 它影响多长时间 |
| `Lifecycle` | 生命周期 | 从开始到峰值再到结束 |
| `BaseValue` | 基础值 | 这个因子天然有多强 |
| `Weight` | 权重 | 系统给它多大重要性 |
| `CurrentStrength` | 当前强度 | 现在这个时刻它有多强 |
| `Adjustment` | 调整值 | 最终参与打分的增减量 |
| `DimensionImpact` | 维度影响 | 它对五维各自怎么影响 |
| `SourcePlanet` | 来源行星 | 这条因子主要来自哪颗星 |
| `IsPositive` | 正负 | 这条因子整体偏正还是偏负 |
| `AstroReason` | 占星理由 | 为什么占星上认为它有影响 |
| `RemainingDays` | 剩余天数 | 还会持续多久 |

#### FactorResult

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Factors` | 所有活跃因子 | 当前真的在起作用的原因列表 |
| `YearlyFactors` | 年级因子 | 一年级影响 |
| `MonthlyFactors` | 月级因子 | 几周到几个月的影响 |
| `WeeklyFactors` | 周级因子 | 几天到几周的影响 |
| `DailyFactors` | 日级因子 | 一天到几天的影响 |
| `HourlyFactors` | 小时级因子 | 一两个小时的影响 |
| `PositiveFactors` | 正向因子 | 偏利好的原因 |
| `NegativeFactors` | 负向因子 | 偏压力的原因 |
| `DimensionAdjustments` | 维度调整总和 | 所有因子合起来给五维加了多少 |
| `TotalAdjustment` | 总调整值 | 所有因子综合后的净变化 |

#### ScoreDelta

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Overall` | 综合变化 | 总体比上一周期升了还是降了 |
| `Dimensions` | 维度变化 | 事业、关系等分别怎么变 |
| `Reason` | 主要原因 | 这次变化主要是谁造成的 |

#### Guidance

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Summary` | 总结 | 一句话概括这段时间 |
| `Dos` | 宜做 | 适合做的事情 |
| `Donts` | 忌做 | 不太建议做的事情 |
| `Focus` | 重点维度 | 这段时间最值得关注哪一维 |

#### ScoreBreakdownResponse

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `QueryTime` | 查询时间 | 这次拆解对应哪个时刻 |
| `Granularity` | 粒度 | 是小时、日、周、月还是年 |
| `OverallScore` | 综合分 | 最终展示给用户的总分 |
| `OverallRaw` | 原始综合分 | 还没标准化前的综合值 |
| `Dimensions` | 维度分解 | 五个维度分别怎么来的 |
| `FactorsByLevel` | 按级别分组因子 | 年/月/周/日/小时因子各有哪些 |
| `Meta` | 元信息 | 数据来源、可见级别、因子数量等 |

#### DimensionBreakdown

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Dimension` | 维度名 | 事业/关系/健康/财务/灵性 |
| `BaseScore` | 基础分 | 这个维度的起始分 |
| `AspectScore` | 相位贡献 | 行星相位给这个维度加了多少 |
| `FactorScore` | 因子贡献 | 其他因子给这个维度加了多少 |
| `RawScore` | 原始分 | 加完以后没标准化的分数 |
| `FinalScore` | 最终分 | 标准化后的可展示分数 |
| `Factors` | 因子列表 | 具体是哪些原因造成的 |

#### FactorContribution

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `ID` | 因子 ID | 这条因子的唯一编号 |
| `Name` | 名称 | 因子名，比如“水星逆行” |
| `Type` | 类型 | 这个原因属于哪类 |
| `TimeLevel` | 时间级别 | 它是年级、月级还是小时级 |
| `BaseValue` | 基础值 | 因子本身的原始力度 |
| `Strength` | 强度 | 当前这个时刻有多强 |
| `Adjustment` | 调整值 | 真正参与加减分的数 |
| `Weight` | 权重 | 系统对它的重视程度 |
| `Dimension` | 维度 | 主要影响哪一维 |
| `Description` | 描述 | 这个因子怎么理解 |
| `IsPositive` | 正负 | 偏正向还是偏负向 |

#### DailyForecast

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Date` | 日期 | 这一天的预测 |
| `DayOfWeek` | 星期 | 周几 |
| `OverallScore` | 总分 | 这一天整体顺不顺 |
| `OverallTheme` | 总主题 | 这一天的主旋律 |
| `Dimensions` | 五维分数 | 事业、关系、健康、财务、灵性分别如何 |
| `MoonPhase` | 月相 | 当天月亮处于什么阶段 |
| `MoonSign` | 月亮星座 | 月亮当天落在哪个星座 |
| `HourlyBreakdown` | 小时拆分 | 当天每个小时的分数概况 |
| `ActiveAspects` | 活跃相位 | 当天起作用的相位 |
| `Factors` | 因子结果 | 当天所有因子的完整集合 |
| `TopFactors` | 重点因子 | 对当天影响最大的几个因子 |

#### WeeklyForecast

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `StartDate` | 周起始 | 这一周从哪天开始 |
| `EndDate` | 周结束 | 这一周到哪天结束 |
| `OverallScore` | 总分 | 这一周整体顺不顺 |
| `OverallTheme` | 总主题 | 这一周主旋律 |
| `Dimensions` | 五维分数 | 这周五个方向分别怎样 |
| `DailySummaries` | 每日摘要 | 一周里每天的简短摘要 |
| `KeyDates` | 关键日期 | 这一周里最值得注意的日子 |
| `BestDaysFor` | 最佳日映射 | 哪些事最适合在哪天做 |
| `WeeklyTransits` | 周行运 | 一周内的重要行运主题 |
| `WeeklyFactors` | 周因子集合 | 这周起作用的因子汇总 |
| `DimensionTrends` | 维度趋势 | 各维度是上升还是下降 |

#### AnnualProfection / LifeProfectionMap

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Year` | 年份 | 这是哪一年 |
| `Age` | 年龄 | 这个年限法对应多少岁 |
| `House` | 宫位 | 今年轮到哪个宫位当主角 |
| `HouseName` | 宫位名 | 这个宫位叫什么 |
| `HouseTheme` | 宫位主题 | 这个宫位代表什么主题 |
| `HouseKeywords` | 关键词 | 这个年份的几个关键词 |
| `Sign` | 星座 | 这一年对应什么星座能量 |
| `SignName` | 星座名 | 星座名字 |
| `LordOfYear` | 年主星 | 这一年的主导行星 |
| `LordName` | 年主星名 | 主导行星的名字 |
| `LordSymbol` | 年主星符号 | 主导行星符号 |
| `LordNatalHouse` | 本命宫位 | 这颗年主星在本命盘里落哪宫 |
| `LordNatalSign` | 本命星座 | 这颗年主星在本命盘里落哪座 |
| `Description` | 描述 | 这一年的总体解释 |
| `Profections` | 全部年限法 | 0 到 80 岁的完整列表 |
| `CurrentYear` | 当前年 | 现在这一岁的年限法结果 |
| `UpcomingYears` | 未来几年 | 接下来几年的年限法 |
| `CycleAnalysis` | 周期分析 | 年限法周期阶段说明 |

#### LifeTrendData / LifeTrendPoint

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Type` | 类型 | 年度、季度还是月度趋势 |
| `BirthDate` | 出生日期 | 这条人生趋势属于谁 |
| `Points` | 趋势点 | 时间轴上的每个点 |
| `Summary` | 摘要 | 整体趋势结论 |
| `Cycles` | 周期信息 | 土星周期、木星周期、年限法周期 |
| `Date` | 日期 | 这个趋势点对应哪一天/哪月 |
| `Year` | 年份 | 这个点属于哪一年 |
| `Age` | 年龄 | 当时几岁 |
| `OverallScore` | 总分 | 这个人生节点整体顺不顺 |
| `Harmonious` | 和谐分 | 顺风、轻松的部分 |
| `Challenge` | 挑战分 | 压力、阻力的部分 |
| `Transformation` | 转化分 | 变化、蜕变的部分 |
| `Dimensions` | 五维分数 | 事业、关系等五维状态 |
| `DominantPlanet` | 主导行星 | 这个节点最强的行星 |
| `Profection` | 年限法摘要 | 当年的宫位和主题 |
| `IsMajorTransit` | 是否重大行运 | 有没有特别大的行运事件 |
| `MajorTransitName` | 重大行运名 | 如果有，事件叫什么 |
| `LunarPhaseName` | 月相名 | 这时的月相是什么 |
| `LunarPhaseAngle` | 月相角度 | 太阳和月亮的角度差 |

#### ProgressedChart

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `TargetDate` | 目标日期 | 你想看的推运日期 |
| `ProgressedDate` | 推运日期 | 折算后的推运盘时间 |
| `DaysProgressed` | 推进天数 | 推了多少天 |
| `YearsFromBirth` | 出生后年数 | 从出生到目标日期过去了多少年 |
| `Planets` | 推运行星 | 推运后的行星位置 |
| `ProgressedAscendant` | 推运上升点 | 推运盘的上升点 |
| `ProgressedMidheaven` | 推运中天 | 推运盘的中天 |
| `Aspects` | 推运相位 | 推运盘行星之间的关系 |
| `LunarPhase` | 推运月相 | 推运盘里的月相状态 |

#### TransitEvent / TransitResult

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Date` | 日期 | 这条行运事件在哪天 |
| `TransitPlanet` | 行运星 | 当前在动的星 |
| `NatalPlanet` | 本命星 | 被影响的本命星 |
| `Aspect` | 相位 | 两星之间是什么关系 |
| `Phase` | 阶段 | exact/active 等状态 |
| `Intensity` | 强度 | 这条行运有多强 |
| `Duration` | 持续时间 | 这条行运从什么时候到什么时候 |
| `Interpretation` | 解读 | 行运的主题和建议 |
| `StartDate` | 开始日期 | 行运范围起点 |
| `EndDate` | 结束日期 | 行运范围终点 |
| `Events` | 事件列表 | 这段时间内有哪些行运 |
| `OverallScore` | 总分 | 整段行运的总体强弱 |
| `DominantThemes` | 主导主题 | 这段行运最常见的主题 |

#### TimeSeries

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Granularity` | 粒度 | 这个时间序列按小时还是按天 |
| `Points` | 数据点 | 一条曲线上的每个点 |
| `Stats` | 统计信息 | 平均值、波动、最大最小值等 |
| `Time` | 时间 | 这个点在哪一刻 |
| `Label` | 标签 | 前端显示名 |
| `Raw` | 原始分 | 没经过展示处理的分数 |
| `Display` | 展示分 | 前端直接画出来的分数 |
| `Dimensions` | 维度值 | 每个维度在这个点上的值 |
| `Volatility` | 波动率 | 这个点附近波动大不大 |

#### ActiveFactorsResponse / ActiveFactorInfo

这个结构专门回答“某个时间点或时间范围里，哪些影响因子还在生效”。

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `Granularity` | 查询粒度 | 你是在查年、月、周还是日 |
| `RangeStart` | 范围开始 | 这次查询从哪一刻开始 |
| `RangeEnd` | 范围结束 | 这次查询到哪一刻结束 |
| `Infect` | 过滤模式 | `all` 看全部，`core` 只看当前粒度可见的 |
| `Factors` | 活跃因子列表 | 在这个范围里真正起作用的因子 |
| `TotalCount` | 总数 | 一共有多少条活跃因子 |
| `PositiveCount` | 正向数 | 有多少条偏利好 |
| `NegativeCount` | 负向数 | 有多少条偏压力 |

`ActiveFactorInfo` 里每条因子还会带这些字段：

| 字段 | 含义 | 大白话 |
|------|------|--------|
| `ID` | 因子 ID | 这条因子的唯一编号 |
| `Name` | 因子名 | 比如“Mercury Retrograde” |
| `Type` | 类型 | 逆行、相位、空亡等 |
| `TimeLevel` | 时间级别 | 年级、月级、周级、日级、小时级 |
| `BaseValue` | 基础值 | 因子原始力度 |
| `Weight` | 权重 | 系统给它的重视程度 |
| `Description` | 描述 | 这条因子怎么理解 |
| `IsPositive` | 正负 | 偏正还是偏负 |
| `Effect` | 效果标签 | `positive` 或 `negative` |
| `StartTime` | 开始时间 | 这条因子什么时候开始生效 |
| `EndTime` | 结束时间 | 这条因子什么时候结束 |
| `PeakTime` | 峰值时间 | 这条因子什么时候最强 |
| `MaxStrength` | 最大强度 | 在你查的范围内它最强到什么程度 |

大白话：

- 如果你问“某个时间点有哪些影响因子”，这就是最直接的答案
- 如果你问“某个时间点为什么分数这么高/这么低”，先看这里，再看 `ScoreBreakdown`
- 如果你问“这条因子还会持续多久”，看 `StartTime`、`PeakTime`、`EndTime`、`MaxStrength`

---

## 4. 本命盘计算

本命盘入口是：

```go
CalculateNatalChart(birthData)
```

计算步骤：

1. 将出生时间转为儒略日
2. 计算行星黄经、逆行状态、星座落点
3. 计算宫位、上升点、中天
4. 为行星分配宫位
5. 计算本命相位
6. 检测格局
7. 计算元素/模式平衡
8. 找出主导行星与命主星

核心文件：

- `backend/astro/natal_chart.go`
- `backend/astro/houses.go`
- `backend/astro/aspects.go`
- `backend/astro/ephemeris.go`

---

## 5. 本命基础分模型

基础分不是查询时刻的分数，而是本命盘的长期底座。

入口：

- `CalculateNatalBaseScores(chart)`

你可以把它理解成“先看这个人的天赋底分”，再去看当天外部环境怎么推高或拉低。

计算公式可以概括为：

```text
基础分 = 50
        + 行星落宫贡献
        + 宫主星状态贡献
        + 相位格局贡献
        -> clamp 到 35~65
```

### 5.1 行星落宫贡献

对每颗行星：

1. 找到其所在宫位
2. 获取尊贵度
3. 获取行星自然权重
4. 获取该行星对五维的映射
5. 如果落入和该维度相关的宫位，额外加成

大白话解释：

- 太阳更像“我是谁、我要成为什么样的人”
- 月亮更像“我怎么感受、怎么建立安全感”
- 水星更像“我怎么说、怎么想、怎么交易”
- 金星更像“我怎么爱、怎么享受、怎么看待价值”
- 火星更像“我怎么行动、怎么争取”
- 木星更像“机会、扩张、好运”
- 土星更像“责任、压力、长期建设”

所以某颗星落在哪个宫位，就像“这股能力主要用在什么人生场景里”。

宫位相关维度映射在：

- `backend/astro/natal_base_scores.go`

### 5.2 宫主星贡献

对关键宫位 `1/2/6/7/8/9/10/11/12`：

1. 取宫头星座
2. 找该星座守护星
3. 看守护星尊贵度
4. 看守护星所在宫位力量
5. 若逆行，扣分

大白话解释：

这部分是在看“负责这个领域的管理员状态好不好”。

- 管理员状态强，这个领域通常更稳
- 管理员逆行或失势，这个领域更容易卡顿、拖延、反复

### 5.3 相位格局贡献

只做轻量格局增强：

- 和谐相位多于紧张相位时，给全维度小幅加成
- 目前大三角、大十字等格局还保留 TODO

大白话解释：

相位格局可以理解成“这盘棋是顺风局还是压力局”。

- 顺风局：很多星彼此配合，底分更容易上去
- 压力局：冲突相位多，底分不会特别高，但也可能代表强驱动力

---

## 6. 查询时刻评分主链路

当前核心评分入口：

- `CalculateScoresV2(chart, date)`
- `CalculateScoresV2WithPositions(chart, date, transitPositions)`
- `CalculateScoresV2Lite(chart, date)`

### 6.1 标准链路

1. 计算本命基础分
2. 计算影响因子
3. 将因子累加到五维
4. 对每个维度做标准化
5. 用维度权重求综合分

公式简化为：

```text
维度分 = Normalize(基础分 + 因子调整)
综合分 = Σ(维度分 × 维度权重)
```

这就是“评价”是怎么来的：

1. 先给每个维度一个初始分
2. 再把当前时刻的天象加减进去
3. 最后做一个压缩，避免分数无限变大或变小

所以最终看到的不是“绝对命运值”，而是“当前时段相对更顺还是更逆”。

### 6.2 轻量链路

`CalculateScoresV2Lite` 会跳过精确相位时间搜索，适合：

- 趋势图
- 批量采样
- 聚合计算

这是性能优化分支，不是功能简化版。

适用场景：

- 用户滑动时间轴时的连续刷新
- 趋势图大量采样
- 批量预计算
- 缓存预热

### 6.3 标准化函数

`NormalizeScoreV2(raw)` 使用 `tanh` 压缩：

```text
centered = raw - 50
compressed = tanh(centered / 30)
normalized = 50 + 50 * compressed
```

最后裁剪到 `0~100`，并保留 4 位小数。

这意味着分数会向 50 收敛，但极端高低值仍能保留方向性。

大白话：

- 50 分附近 = 中性
- 高于 50 = 偏顺
- 低于 50 = 偏紧张
- 但不会因为某个事件特别强就直接冲到 100 或跌到 0

---

## 7. 影响因子系统

影响因子入口：

- `CalculateInfluenceFactorsV2`
- `CalculateInfluenceFactorsLite`

系统会把不同类型的天象统一成 `InfluenceFactor`，然后汇总到 `FactorResult`。

大白话：

`InfluenceFactor` 就是一条“原因记录”。

例如：

- “水星逆行了，所以沟通和合同可能受影响”
- “月亮空亡，所以不适合立刻拍板”
- “木星与本命太阳成好相位，所以更容易出现机会”

### 7.1 当前启用的主要因子

当前 V2 主链路使用以下因子：

- 尊贵度 `dignity`
- 逆行 `retrograde`
- 相位 `aspectPhase`
- 月相 `lunarPhase`
- 行星时 `planetaryHour`
- 年主星 `profectionLord`
- 月亮空亡 `voidOfCourse`

这些因子不是全部占星知识，而是当前后端真正参与评分的子集。

### 7.2 因子生命周期

每个因子都有生命周期：

- `StartTime`
- `PeakTime`
- `EndTime`
- `Duration`

强度不是常量，而是随生命周期变化。

这表示“某个天象不是在整个时间段里都一样强”。

比如：

- 进入前，影响还没开始
- 接近精确时刻，影响最强
- 离开后，影响逐渐消散

### 7.3 强度曲线

`CalculateFactorStrength` 的策略：

- 生命周期外：强度 `0`
- 无生命周期：强度 `1`
- 前半段：正弦上升
- 后半段：指数衰减

即：

```text
前半段 strength = sin(pi * progress)
后半段 strength = exp(-3 * t)
```

这模拟了“入相慢、出相快”的占星节律。

如果不用占星术语，说白了就是：

- 事情临近发生时，影响会逐渐升高
- 真正到点以后，影响达到峰值
- 过去之后，影响会慢慢退场

### 7.4 因子调整值

```text
Adjustment = BaseValue × Weight × CurrentStrength
```

然后再根据维度影响矩阵分配到五维。

这里的三个乘数分别表示：

- `BaseValue`：这个事件先天有多强
- `Weight`：系统给它的重要程度
- `CurrentStrength`：当前时刻离峰值有多近

所以同样是“水星逆行”，在高峰期和尾声期影响是不一样的。

### 7.5 维度分配

V2 的维度分配不是简单平均，而是通过有符号影响矩阵完成：

- `FactorTypeDefaultImpactsV2`
- `RetrogradeImpactsByPlanetV2`

每个因子会先得到一个基准调整值，再乘以各维度的符号和比例。

大白话：

同一个事件，可能对不同生活领域产生不同方向的影响。

例如水星逆行：

- 事业：沟通混乱，偏负
- 关系：误会增加，偏负
- 灵性：适合复盘内省，偏正

所以系统不会给所有维度一个统一加减，而是“分领域处理”。

---

## 8. 各类因子的算法

### 8.1 尊贵度因子

入口：

- `calculateDignityFactorsV2`

逻辑：

1. 遍历当前行运行星
2. 计算尊贵度
3. 尊贵度为 `domicile / exaltation / detriment / fall` 时生成因子
4. 根据行星平均速度估算生命周期
5. 时间级别设为月度

大白话：

这是在看某颗星“待在自己舒服的地盘，还是别扭的地盘”。

- 舒服地盘 = 更容易正常发挥
- 别扭地盘 = 能量受限，需要更多代价才能表达出来

### 8.2 逆行因子

入口：

- `calculateRetrogradeFactorsV2`

逻辑：

1. 找出当前逆行行星
2. 不同行星给不同基础值
3. 生命周期按逆行周期估算
4. 时间级别按行星类型决定

大白话：

逆行不是真的“倒着走”，而是从地球视角看这颗星的能量更像“回头、复盘、延迟、反省”。

系统里最典型的理解是：

- 水星逆行：沟通、文书、合同容易反复
- 金星逆行：关系、审美、价值判断容易重估
- 火星逆行：行动力下降，但适合修整

### 8.3 相位因子

入口：

- `calculateAspectFactorsV2`
- `calculateAspectFactorsLite`
- `calculateAspectFactorsWithSharedData`

逻辑：

1. 计算行运对本命相位
2. 只保留强度足够高的相位
3. 依据和谐/紧张相位设定正负
4. 使用精确或估算的生命周期
5. 合并两颗行星的维度影响

大白话：

相位就是“两颗星之间的配合关系”。

- 和谐相位：像队友配合顺
- 紧张相位：像拉扯、冲突、压力测试
- 合相：两颗星靠得很近，能量叠加，正负要看星本身性质

#### 精确相位时间

精确相位时间搜索在：

- `backend/astro/aspect_search.go`

算法是二分法搜索：

1. 先在时间区间内找符号变化
2. 再用二分法逼近
3. 精度可到 `1e-6` 儒略日

这是高精度链路，代价高于估算链路。

为什么要做精确时间：

- 用于事件卡点
- 用于生命周期峰值
- 用于当天“最强时刻”的排序

### 8.4 月相因子

入口：

- `calculateLunarPhaseFactorsV2`
- `calculateLunarPhaseFactorsWithSharedData`

逻辑：

1. 计算太阳和月亮的黄经差
2. 映射到月相分类
3. 用月相表给出基础值
4. 生命周期固定为约 3.5 天

大白话：

月相就是“月亮和太阳的相对位置”。

它通常更像情绪节律、节奏切换、启动/收尾状态，不是长期结构性力量。

### 8.5 行星时因子

入口：

- `calculatePlanetaryHourFactorsV2`
- `calculatePlanetaryHourFactorsWithSharedData`

逻辑：

1. 使用迦勒底序列 `Saturn -> Jupiter -> Mars -> Sun -> Venus -> Mercury -> Moon`
2. 日主星由星期映射决定
3. 按日夜分段估算 24 个行星时
4. 当前小时的主宰星决定因子强度

实现上日出/日落目前是简化常量 `6:00 / 18:00`，不是实测天文日出日落。

大白话：

这更像“今天这个小时由哪颗星做主”。

适合用来做短时间内的节奏建议，比如：

- 适合开会还是整理
- 适合沟通还是独处
- 适合行动还是复盘

### 8.6 年主星因子

入口：

- `calculateProfectionLordFactorsV2`
- `CalculateAnnualProfection`

逻辑：

1. 年限法按 `age % 12 + 1` 映射到宫位
2. 找出该宫位对应星座及守护星
3. 看年主星当前状态
4. 用年主星状态影响全年主题

大白话：

这像“今年的年度主题”。

它不是说一年里每一天都一样，而是给全年一个主旋律。

### 8.7 月亮空亡因子

入口：

- `CalculateVoidOfCourse`
- `calculateVoidOfCourseFactorsV2`

逻辑：

1. 计算月亮当前星座剩余度数
2. 在月亮离座前向前搜索是否还会形成主要相位
3. 若不会，判定为空亡
4. 影响值随空亡持续时间增长而增强

大白话：

月亮空亡就是“事情还没真正接上下一步”。

常见理解是：

- 不适合立刻拍板
- 不适合开始特别重要的新事
- 更适合观察、缓冲、等待

### 8.8 查询某个时间点时，怎么同时看因子持续时间、相关事件和五维趋势

当技术同学问“某个时间点到底发生了什么”，实际要看三层结果：

1. 这个时间点有哪些因子正在生效
2. 这些因子对应哪些事件
3. 这些因子把五维趋势往哪边推

#### 第一步：先看因子是否还在有效期

每个 `InfluenceFactor` 都带有 `Lifecycle`：

- `StartTime`：什么时候开始生效
- `PeakTime`：什么时候最强
- `EndTime`：什么时候结束
- `Duration`：持续多久
- `RemainingDays`：从当前时刻到结束还剩多久

判断逻辑很简单：

- 如果当前时间早于 `StartTime`，说明还没开始
- 如果当前时间晚于 `EndTime`，说明已经结束
- 如果在中间，说明这个因子正在起作用

大白话：

- `StartTime` 前 = 还没上场
- `StartTime ~ PeakTime` = 越来越强
- `PeakTime` = 最强点
- `PeakTime ~ EndTime` = 逐渐退场
- `EndTime` 后 = 影响结束

#### 第二步：再看这个时间点关联了哪些事件

事件不是凭空来的，而是从同一套星象数据里抽出来的。

常见来源：

- 相位事件：`CalculateAspectEvents` / `CalculateDailyEvents`
- 月相事件：`calculateLunarPhaseEvents`
- 逆行事件：`calculateRetrogradeEvents`
- 月亮空亡：`CalculateVoidOfCourse`
- 行星时事件：`CalculatePlanetaryHourEnhanced`

事件对象里最重要的几个字段是：

- `Type`：事件类型
- `Title`：给前端看的标题
- `Intensity`：强度
- `StartTime` / `EndTime`：事件时间范围
- `ExactTime`：最准确的卡点
- `Impact`：这条事件对五维的绝对影响
- `ImpactDelta`：和上一周期相比的变化

大白话：

- `Impact` 说的是“这件事本身有多大影响”
- `ImpactDelta` 说的是“和上一段时间比，是变强还是变弱”

#### 第三步：最后看五维趋势怎么被推出来

五维趋势不是单靠一个事件，而是多个因子一起作用：

```text
五维趋势 = 当前基础分 + 所有活跃因子的调整 + 相位贡献
```

其中：

- `FactorResult.DimensionAdjustments` 是因子对五维的总影响
- `ScoreDelta.Dimensions` 是和上一周期对比后的变化
- `SubSlots` 或 `TimeSeries` 是把这个变化铺成曲线

如果想看“趋势影响”而不是“单点得分”，通常要看这三份东西：

1. `DimensionAdjustments`
2. `ScoreDelta`
3. `SubSlots` / `TimeSeries`

#### 怎么理解“趋势影响”

举个例子，某时刻可能出现：

- 火星和金星形成和谐相位
- 月亮处于情绪更活跃的阶段
- 水星逆行还在尾声

那系统不会只给一个总分，而是会分别推到：

- 事业：可能更适合沟通推进
- 关系：更容易有互动、合作、表达
- 健康：如果火星能量太强，也可能提示精力消耗
- 财务：需要看是否有冲动决策
- 灵性：如果逆行和月相叠加，反思感会更强

#### 实际查询时的推荐顺序

如果你要分析一个时间点，建议按这个顺序看：

1. 先看 `TimeSlot.scores`，知道整体好坏
2. 再看 `TimeSlot.events`，知道发生了什么
3. 再看 `FactorResult.Factors` 或 `ScoreBreakdownResponse`，知道为什么
4. 再看 `delta`，知道和上一周期比怎么变
5. 最后看 `subSlots` 或 `TimeSeries`，知道趋势是怎么走的

#### 适合的接口/入口

- 只想看当前这个时间点：`CalculateScoresV2`
- 想看为什么这么打分：`CalculateScoreBreakdown`
- 想看当天具体事件：`CalculateDailyEvents`
- 想看因子还会持续多久：看 `InfluenceFactor.Lifecycle` 和 `RemainingDays`
- 想看一段时间的趋势：`CalculateScoresV2Lite` 或 `CalculateTimeSeries`

### 8.9 如果要查“某个时间点对应哪些影响因子（天体事件）”，应该怎么查

这里先把名词说清楚：

- **影响因子**：在这个时间点仍然有效、会对五维分数产生作用的天体状态
- **天体事件**：前端能直接展示的一条条事件，比如合相、对冲、月相变化、逆行、空亡、行星时切换
- **查询某个时间点**：不是只看一个瞬间，而是要看这个瞬间落在哪些因子的有效区间里

如果技术同学问“这个时间为什么是这个分”，不要只找一个事件，而是要按下面顺序查：

1. 先找这个时间点**活跃的因子**
2. 再找这个时间点**对应的事件**
3. 再看这些因子**怎么改变五维**
4. 最后把单点结果和趋势曲线对起来

#### 1）先查这个时间点有哪些因子正在生效

核心入口是：

- `GetActiveFactorsInRange`

它做的事不是单纯返回一个数，而是返回“这段时间内哪些因子是活着的”。

常见输出字段可以这样理解：

- `Factors`：活跃因子列表
- `RangeStart` / `RangeEnd`：这次查询覆盖的时间范围
- `FactorCount`：因子数量
- `VisibleFactorCount`：前端应该显示出来的数量
- `MaxStrength`：这段范围里最强的因子强度

每个因子对象里最关键的是这些：

- `StartTime`：什么时候开始起作用
- `PeakTime`：什么时候最强
- `EndTime`：什么时候结束
- `Duration`：这个因子总共持续多久
- `RemainingDays`：从当前时刻到结束还剩多久
- `Strength`：当前时刻的强度
- `DimensionAdjustments`：它分别把五维往哪边推

大白话：

- `StartTime` 前面，说明还没开始
- `StartTime ~ PeakTime`，说明越来越强
- `PeakTime`，说明到顶了
- `PeakTime ~ EndTime`，说明慢慢退场
- `EndTime` 后面，说明已经不影响了

#### 2）再查这个时间点对应哪些事件

核心入口是：

- `CalculateDailyEvents`

它返回的是“前端可展示的事件流”，比单纯的因子更像业务语言。

常见事件来源包括：

- 相位事件：两颗星之间形成相位
- 月相事件：新月、满月、上弦、下弦等
- 逆行事件：行星速度方向变化
- 空亡事件：月亮暂时没有接上下一步
- 行星时事件：当前小时由哪颗星主导

这一步适合回答：

- 这个时间点发生了什么
- 这个时间点为什么会有提示
- 前端弹窗、日历、提醒卡片应该展示什么

#### 3）再查这些因子怎么影响五维

核心入口是：

- `CalculateScoreBreakdown`
- `CalculateScoresV2`

它们回答的是“为什么这个时间点分数变了”。

关键字段可以这样看：

- `ScoreBreakdownResponse`：总分为什么是现在这个样子
- `FactorContribution`：某个因子具体贡献了多少
- `DimensionBreakdown`：某个因子对事业、关系、健康、财务、灵性的分别影响
- `DimensionAdjustments`：把影响汇总到五维后的结果

这里最重要的一点是：

- **事件** 是“发生了什么”
- **因子** 是“什么力量还在起作用”
- **分数拆解** 是“这股力量怎么落到五维上”

#### 4）如果要看趋势，不要只看单点

单个时间点只能说明“这一刻”，但技术上常常还要回答：

- 这个影响是刚开始还是快结束了
- 下一小时会更强还是更弱
- 这条事件会不会持续一整天

这时要结合：

- `CalculateScoresV2Lite`
- `CalculateTimeSeries`
- `SubSlots`
- `TimeSeries`

它们的作用是把单点变成曲线，方便看上升、下降、拐点和峰值。

#### 推荐的排查顺序

如果是排查某个时间点，推荐直接按下面顺序看：

1. `GetActiveFactorsInRange`：这个时间点有哪些因子还在生效
2. `CalculateDailyEvents`：这个时间点对应哪些可展示事件
3. `CalculateScoreBreakdown`：每个因子具体推了哪一维
4. `CalculateScoresV2`：最终这个点的综合结果
5. `CalculateScoresV2Lite` / `CalculateTimeSeries`：前后趋势怎么走

#### 一句话总结

如果要查“某个时间点对应哪些影响因子（天体事件）”，不要只找一个接口。

正确做法是：

- 用 `GetActiveFactorsInRange` 找活跃因子
- 用 `CalculateDailyEvents` 找事件
- 用 `CalculateScoreBreakdown` 找原因
- 用 `CalculateScoresV2` 和趋势接口看结果

这样才能从“这个时间点有什么”一路追到“为什么会这样”。

---

## 9. 共享星象数据层

共享数据层见：

- `backend/astro/shared_astro_data.go`

它的目标是避免以下重复计算：

- `DailyEvents`
- `Factor` 系统

共享缓存 key：

- 日期 `YYYY-MM-DD`
- 本命盘 ID

缓存上限：

- 100 天

计算内容包括：

- 相位事件
- 月相事件
- 换座事件
- 行星时数据
- 当前月相
- 当前行星时

这个设计把“精确事件搜索”从多个调用点收敛到一个地方。

---

## 10. 事件系统

事件对象统一为 `AstroEvent`。

事件类型主要有：

- `aspect`
- `transit_house`
- `progression`
- `planetary_hour`
- `void_of_course`
- `retrograde`
- `lunar_phase`
- `sign_change`
- `dignity`

### 10.1 事件生成顺序

`Calculator.calculateEventsWithPositions` 的顺序是：

1. 相位事件
2. 行星时事件
3. 逆行事件
4. 月空事件
5. 月相事件
6. 换座事件
7. 尊贵度事件
8. 排序
9. i18n 翻译

### 10.2 事件去重

事件 ID 设计保证同类事件跨粒度尽量稳定：

- 相位生命周期事件用生命周期时间生成稳定 ID
- 逆行、行星时等按日期或小时生成 ID

聚合时会根据 ID 合并，强度更高者优先保留。

---

## 11. 时间聚合策略

聚合器入口：

- `backend/core/aggregator.go`

目标是把小时级结果压缩成日/周/月/年，而不是直接全量重算。

### 11.1 日级

采样点：

- 0, 4, 8, 12, 16, 20 点

### 11.2 周级

采样点：

- 周一、周三、周六中午

### 11.3 月级

采样点：

- 1, 8, 15, 22, 29 号中午

### 11.4 年级

采样点：

- 2、5、8、11 月 15 号中午
- 再映射到季度起始月

### 11.5 聚合规则

所有聚合都基于采样点的均值：

```text
父级分数 = 子级分数均值
```

为什么要采样而不是全量算：

- 全量算更精确，但成本高
- 采样算足够稳定，速度更快
- 对趋势图和大粒度页面更实用

事件则采用：

- 过滤粒度
- 合并区间重叠
- 按 `EventID` 或事件组去重

### 11.6 前一周期 Delta

每个日/周/月/年 slot 都会对比前一周期，生成：

- `ScoreDelta`
- `ImpactDelta`

`ImpactDelta` 是按事件与周期的重叠比例缩放后计算的。

---

## 12. 趋势计算

趋势接口：

- `/api/v2/astro/trend`

趋势点使用轻量分数链路：

- 小时趋势：一天 24 个点
- 日趋势：当月每天一个点
- 周趋势：当月每周一个点
- 月趋势：当年每月一个点
- 年趋势：每年一个点

趋势接口不生成事件，只输出分数曲线和摘要，目的是降低计算成本。

适合用在：

- 曲线图
- 周/月年对比
- “最近是变好还是变差”这种浏览场景

趋势汇总规则：

- `max`
- `min`
- `trend` = `upward / downward / stable`

---

## 13. 人生趋势与推运

### 13.1 推运

入口：

- `CalculateProgressions`

规则：

- 1 天 = 1 年
- 目标日期相对出生日期折算成年数
- 按这个年数推进本命盘

输出包含：

- 推运行星
- 推运 ASC/MC
- 推运相位
- 推运月相

### 13.2 人生趋势

入口：

- `CalculateLifeTrend`

它把多个维度合在一起：

- `CalculateScoresV2`
- 行运相位强弱
- 年限法
- 推运月相
- 重大行运

这个模块更像“人生时间轴分析器”，不是单点评分器。

适合用在：

- 人生阶段回顾
- 年龄节点分析
- 长周期转折点提示
- 个人成长轨迹图

---

## 14. 自定义因子

自定义因子格式：

```text
AddScore=(2*healthScore,2.5,202501171230)
```

含义：

- 操作类型：`AddScore / SubScore / MulScore / SetScore`
- 维度：`career / relationship / health / finance / spiritual / overall`
- 值：数值
- 持续时间：小时
- 开始时间：`YYYYMMDDHHmm`

实现：

- `backend/astro/custom_factors.go`

这些因子会被转换成标准影响因子后参与评分，也能直接作用于 `DimensionScoresV2`。

---

## 15. 缓存与预热

统一 API：

- `POST /api/v2/astro`

会先查全局缓存，再计算，再写回缓存。

缓存 key 由以下部分组成：

- userID
- granularity
- queryTime 归一化
- language

大粒度查询还会异步预热相邻周期：

- 周：前后各 1 周
- 月：前后各 1 月
- 年：前后各 1 年

这能明显降低翻页和连续浏览的延迟。

---

## 16. i18n 与展示层

i18n 不参与核心数值计算，只做展示映射：

- 行星名
- 相位名
- 维度名
- 情感化标题
- 解释文本

相关实现：

- `backend/i18n/*`

这层不会改变分数，只改变输出语言和文本风格。

---

## 17. 另一个评分分支：UnifiedScore

项目里还存在一套 `UnifiedScore`/`CalculateUnifiedHourlyScore` 体系：

- `backend/astro/unified_score.go`

它同样基于行运相位 + 因子，但输出更偏“统一时间序列”，并会加入视觉抖动。

当前这套体系更像兼容/扩展分支，不是 `POST /api/v2/astro` 的主返回模型。

---

## 18. 设计边界与已知近似

为了让技术同学正确理解系统，这里明确列出近似点：

1. 行星时使用固定 `6:00 / 18:00` 近似日出日落，不是实测天文时刻
2. 年/月/周聚合使用采样均值，不是全时段连续积分
3. 尊贵度、维度影响、权重矩阵都属于可调的启发式模型
4. 相位生命周期有精确搜索链路，但趋势和聚合常走轻量链路
5. 月空亡、逆行等时长在部分场景是估算值

因此这套系统的目标是：

- 一致
- 可解释
- 可扩展
- 性能可控

它更适合以下场景：

- 个人时间管理建议
- 事件节奏提示
- 内容运营中的“今日/本周/本月主题”
- 用户侧的长期趋势看板
- 解释型产品，而不是黑盒判定

不适合直接当作：

- 医疗诊断
- 投资建议
- 法律结论
- 唯一决策依据

这些输出更适合做“参考提示”，不是硬性结论。

而不是绝对天文精度优先。

---

## 19. 这个算法可以用在哪些地方

如果从产品角度看，当前这套算法最适合这些场景：

### 19.1 用户端

- 今日运势卡片
- 本周/本月趋势图
- 关键时间点提醒
- 重要事件提醒
- 个人成长回顾

### 19.2 内容端

- 每日主题文案
- 每周专题内容
- 情绪化标题生成
- 行动建议卡片

### 19.3 运营端

- 用户分层
- 时间窗口推送
- 缓存预热
- 热点日期预测

### 19.4 数据端

- 时间序列分析
- 事件贡献分解
- 因子命中率统计
- 维度波动追踪

### 19.5 技术端

- 可解释规则引擎
- 评分模型原型
- 多时间粒度聚合框架
- 可插拔因子系统

---

## 20. 推荐阅读顺序

如果要快速吃透代码，建议按这个顺序看：

1. `backend/astro/natal_chart.go`
2. `backend/astro/score_calculator.go`
3. `backend/astro/factor_lifecycle.go`
4. `backend/astro/score_breakdown.go`
5. `backend/core/timeslot.go`
6. `backend/core/calculator.go`
7. `backend/core/aggregator.go`
8. `backend/api/v2/astro_handler.go`
9. `backend/api/v2/trend_handler.go`
