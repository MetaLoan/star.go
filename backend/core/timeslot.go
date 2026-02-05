package core

import (
	"star/models"
	"time"
)

// ==================== 统一数据模型 ====================
// TimeSlot 是系统的核心数据单元
// 一个用户的一个时间周期 = 一个 TimeSlot

// TimeSlot 时间槽（核心数据单元）
type TimeSlot struct {
	// 时间标识
	UserID      string    `json:"userId"`      // 用户/本命盘ID
	StartTime   time.Time `json:"startTime"`   // 开始时间
	EndTime     time.Time `json:"endTime"`     // 结束时间
	Granularity string    `json:"granularity"` // hour/day/week/month/year

	// 五维分数（当前周期的聚合分数）
	Scores DimensionScores `json:"scores"`

	// 天体事件列表（所有事件，前端自行筛选）
	Events []AstroEvent `json:"events"`

	// 与上一时间槽的变化
	Delta *ScoreDelta `json:"delta,omitempty"`

	// 综合指导
	Guidance *Guidance `json:"guidance,omitempty"`

	// 子时间槽（用于绘制分数曲线）
	// - 查询 day 时：包含 24 个 hour 的简化数据
	// - 查询 week 时：包含 7 个 day 的简化数据
	// - 查询 month 时：包含 ~30 个 day 的简化数据
	// - 查询 hour 时：为空（最小粒度）
	SubSlots []SubSlot `json:"subSlots,omitempty"`
}

// SubSlot 子时间槽（简化版，用于曲线绘制）
type SubSlot struct {
	StartTime  time.Time       `json:"startTime"`
	Scores     DimensionScores `json:"scores"`     // 五维分数
	EventCount int             `json:"eventCount"` // 该时段事件数量
}

// DimensionScores 五维分数
type DimensionScores struct {
	Overall      float64 `json:"overall"`      // 综合分（由五维加权计算）
	Career       float64 `json:"career"`       // 事业
	Relationship float64 `json:"relationship"` // 关系
	Health       float64 `json:"health"`       // 健康
	Finance      float64 `json:"finance"`      // 财富
	Spiritual    float64 `json:"spiritual"`    // 心灵
}

// DimensionImpact 五维影响值
type DimensionImpact struct {
	Career       float64 `json:"career"`
	Relationship float64 `json:"relationship"`
	Health       float64 `json:"health"`
	Finance      float64 `json:"finance"`
	Spiritual    float64 `json:"spiritual"`
}

// ScoreDelta 分数变化
type ScoreDelta struct {
	Overall    float64         `json:"overall"`    // 综合变化（正=上升，负=下降，0=平稳）
	Dimensions DimensionScores `json:"dimensions"` // 各维度变化
	Reason     string          `json:"reason"`     // 主要变化原因
}

// Guidance 综合指导
type Guidance struct {
	Summary string   `json:"summary"` // 概要
	Dos     []string `json:"dos"`     // 宜做
	Donts   []string `json:"donts"`   // 忌做
	Focus   string   `json:"focus"`   // 重点关注维度
}

// ==================== 天体事件 ====================

// AstroEvent 统一天体事件
type AstroEvent struct {
	// 唯一标识（用于前端去重、跨粒度追踪同一事件）
	EventID string `json:"eventId"` // 如 "mercury_retrograde_20260125", "moon_trine_venus_20260205_1423"

	// 基础信息
	Type       string  `json:"type"`       // aspect/transit_house/progression/planetary_hour/void_of_course/retrograde
	Title      string  `json:"title"`      // 事件标题
	IsPositive bool    `json:"isPositive"` // 正面/负面
	Intensity  float64 `json:"intensity"`  // 强度 0-1
	TimeLevel  string  `json:"timeLevel"`  // hourly/daily/weekly/monthly/yearly（事件的时间层级）

	// 星体信息（便于前端按星体筛选）
	PrimaryPlanet       string `json:"primaryPlanet"`                 // 主影响星体 ID (sun/moon/mars...)
	PrimaryPlanetName   string `json:"primaryPlanetName"`             // 主影响星体名称（本地化）
	SecondaryPlanet     string `json:"secondaryPlanet,omitempty"`     // 次要星体 ID（相位事件时）
	SecondaryPlanetName string `json:"secondaryPlanetName,omitempty"` // 次要星体名称

	// 相位详情（仅 type=aspect 时有值）
	Aspect     string  `json:"aspect,omitempty"`     // conjunction/trine/square/opposition/sextile
	AspectName string  `json:"aspectName,omitempty"` // 相位名称（本地化）
	Orb        float64 `json:"orb,omitempty"`        // 容许度

	// 对五维的影响
	// 重要：同一事件在不同粒度下 ImpactDelta 值不同！
	Impact      DimensionImpact `json:"impact"`      // 该事件对各维度的绝对影响值
	ImpactDelta DimensionImpact `json:"impactDelta"` // 相对上一周期的变化值（正=上升，负=下降，0=平稳）

	// 事件时间范围
	StartTime time.Time `json:"startTime"` // 事件/影响开始时间
	EndTime   time.Time `json:"endTime"`   // 事件/影响结束时间
	ExactTime time.Time `json:"exactTime"` // 精确时刻（相位精确成相的时间）

	// 解读文本
	Interpretation string `json:"interpretation"` // 详细解读
	Advice         string `json:"advice"`         // 行动建议
}

// ==================== 常量定义 ====================

// 时间粒度
const (
	GranularityHour  = "hour"
	GranularityDay   = "day"
	GranularityWeek  = "week"
	GranularityMonth = "month"
	GranularityYear  = "year"
)

// 事件类型
const (
	EventTypeAspect        = "aspect"
	EventTypeTransitHouse  = "transit_house"
	EventTypeProgression   = "progression"
	EventTypePlanetaryHour = "planetary_hour"
	EventTypeVoidOfCourse  = "void_of_course"
	EventTypeRetrograde    = "retrograde"
	EventTypeLunarPhase    = "lunar_phase"
	EventTypeSignChange    = "sign_change"
	EventTypeDignity       = "dignity"
)


// 时间层级
const (
	TimeLevelHourly  = "hourly"
	TimeLevelDaily   = "daily"
	TimeLevelWeekly  = "weekly"
	TimeLevelMonthly = "monthly"
	TimeLevelYearly  = "yearly"
)

// ==================== 辅助方法 ====================

// GenerateEventID 生成事件唯一ID
func GenerateEventID(eventType string, planet1 models.PlanetID, planet2 models.PlanetID, aspect string, exactTime time.Time) string {
	if eventType == EventTypeAspect && planet2 != "" {
		return eventType + "_" + string(planet1) + "_" + aspect + "_" + string(planet2) + "_" + exactTime.Format("20060102_1504")
	}
	if eventType == EventTypeRetrograde {
		return eventType + "_" + string(planet1) + "_" + exactTime.Format("20060102")
	}
	if eventType == EventTypePlanetaryHour {
		return eventType + "_" + string(planet1) + "_" + exactTime.Format("20060102_15")
	}
	return eventType + "_" + string(planet1) + "_" + exactTime.Format("20060102_1504")
}

// NewTimeSlot 创建新的时间槽
func NewTimeSlot(userID string, startTime, endTime time.Time, granularity string) *TimeSlot {
	return &TimeSlot{
		UserID:      userID,
		StartTime:   startTime,
		EndTime:     endTime,
		Granularity: granularity,
		Events:      make([]AstroEvent, 0),
		SubSlots:    make([]SubSlot, 0),
	}
}

// AddEvent 添加事件到时间槽
func (ts *TimeSlot) AddEvent(event AstroEvent) {
	ts.Events = append(ts.Events, event)
}

// SetScores 设置分数
func (ts *TimeSlot) SetScores(scores DimensionScores) {
	ts.Scores = scores
}

// SetDelta 设置变化
func (ts *TimeSlot) SetDelta(delta *ScoreDelta) {
	ts.Delta = delta
}

// SetGuidance 设置指导
func (ts *TimeSlot) SetGuidance(guidance *Guidance) {
	ts.Guidance = guidance
}

// AddSubSlot 添加子时间槽
func (ts *TimeSlot) AddSubSlot(subSlot SubSlot) {
	ts.SubSlots = append(ts.SubSlots, subSlot)
}

// GetEventCount 获取事件数量
func (ts *TimeSlot) GetEventCount() int {
	return len(ts.Events)
}

// ==================== 类型转换 ====================

// FromModelsDimensionScores 从 models.DimensionScoresV2 转换
// 注意：DimensionScoresV2 不包含 Overall，需要单独设置
func FromModelsDimensionScores(scores models.DimensionScoresV2) DimensionScores {
	return DimensionScores{
		Overall:      0, // 需要单独设置（来自 ScoreResult.Overall）
		Career:       scores.Career,
		Relationship: scores.Relationship,
		Health:       scores.Health,
		Finance:      scores.Finance,
		Spiritual:    scores.Spiritual,
	}
}

// ToModelsDimensionScores 转换为 models.DimensionScoresV2
// 注意：Overall 会丢失，因为 DimensionScoresV2 不包含该字段
func (ds DimensionScores) ToModelsDimensionScores() models.DimensionScoresV2 {
	return models.DimensionScoresV2{
		Career:       ds.Career,
		Relationship: ds.Relationship,
		Health:       ds.Health,
		Finance:      ds.Finance,
		Spiritual:    ds.Spiritual,
	}
}

// Add 加法
func (di DimensionImpact) Add(other DimensionImpact) DimensionImpact {
	return DimensionImpact{
		Career:       di.Career + other.Career,
		Relationship: di.Relationship + other.Relationship,
		Health:       di.Health + other.Health,
		Finance:      di.Finance + other.Finance,
		Spiritual:    di.Spiritual + other.Spiritual,
	}
}

// Subtract 减法
func (di DimensionImpact) Subtract(other DimensionImpact) DimensionImpact {
	return DimensionImpact{
		Career:       di.Career - other.Career,
		Relationship: di.Relationship - other.Relationship,
		Health:       di.Health - other.Health,
		Finance:      di.Finance - other.Finance,
		Spiritual:    di.Spiritual - other.Spiritual,
	}
}

// Scale 缩放
func (di DimensionImpact) Scale(factor float64) DimensionImpact {
	return DimensionImpact{
		Career:       di.Career * factor,
		Relationship: di.Relationship * factor,
		Health:       di.Health * factor,
		Finance:      di.Finance * factor,
		Spiritual:    di.Spiritual * factor,
	}
}
