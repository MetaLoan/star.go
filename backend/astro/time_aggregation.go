package astro

import (
	"math"
	"star/models"
	"time"
)

// ==================== 时间层级聚合系统 ====================
// 实现：小时分 → 日分 → 周分 → 月分 → 年分 的聚合
// 采用特征提取法：Average, Max, Min, Trend, Stability

// AggregatedScore 聚合后的分数
type AggregatedScore struct {
	// 综合分数（聚合后）
	Overall float64 `json:"overall"`

	// 五维度分数（聚合后）
	Dimensions models.DimensionScoresV2 `json:"dimensions"`

	// 统计特征
	Features ScoreFeatures `json:"features"`

	// 时间范围
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`

	// 粒度
	Granularity string `json:"granularity"`

	// 子时间段分数（可选，用于下钻）
	SubScores []AggregatedScore `json:"subScores,omitempty"`
}

// ScoreFeatures 分数统计特征
type ScoreFeatures struct {
	Average   float64 `json:"average"`   // 平均值
	Max       float64 `json:"max"`       // 最高值
	Min       float64 `json:"min"`       // 最低值
	Trend     float64 `json:"trend"`     // 趋势斜率（正=上升，负=下降）
	Stability float64 `json:"stability"` // 稳定性（标准差的倒数，越高越稳定）
}

// ==================== 小时分计算 ====================

// CalculateHourlyScore 计算单个小时的分数
func CalculateHourlyScore(chart *models.NatalChart, hourTime time.Time) *ScoreResult {
	return CalculateScoresV2(chart, hourTime)
}

// ==================== 日分聚合 ====================

// AggregateDailyFromHourly 从小时分聚合成日分
// 采用特征提取法
func AggregateDailyFromHourly(chart *models.NatalChart, date time.Time) *AggregatedScore {
	// 生成24个小时的分数
	hourlyScores := make([]*ScoreResult, 24)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for hour := 0; hour < 24; hour++ {
		hourTime := startOfDay.Add(time.Duration(hour) * time.Hour)
		hourlyScores[hour] = CalculateHourlyScore(chart, hourTime)
	}

	// 提取综合分特征
	overallValues := make([]float64, 24)
	for i, s := range hourlyScores {
		overallValues[i] = s.Overall
	}
	overallFeatures := extractFeatures(overallValues)

	// 提取各维度特征
	careerValues := make([]float64, 24)
	relationshipValues := make([]float64, 24)
	healthValues := make([]float64, 24)
	financeValues := make([]float64, 24)
	spiritualValues := make([]float64, 24)

	for i, s := range hourlyScores {
		careerValues[i] = s.Dimensions.Career
		relationshipValues[i] = s.Dimensions.Relationship
		healthValues[i] = s.Dimensions.Health
		financeValues[i] = s.Dimensions.Finance
		spiritualValues[i] = s.Dimensions.Spiritual
	}

	// 聚合维度分数
	dimensions := models.DimensionScoresV2{
		Career:       aggregateWithFeatures(careerValues),
		Relationship: aggregateWithFeatures(relationshipValues),
		Health:       aggregateWithFeatures(healthValues),
		Finance:      aggregateWithFeatures(financeValues),
		Spiritual:    aggregateWithFeatures(spiritualValues),
	}

	// 聚合综合分
	overall := aggregateWithFeatures(overallValues)

	return &AggregatedScore{
		Overall:     overall,
		Dimensions:  dimensions,
		Features:    overallFeatures,
		StartTime:   startOfDay,
		EndTime:     startOfDay.Add(24 * time.Hour),
		Granularity: "daily",
	}
}

// ==================== 周分聚合 ====================

// AggregateWeeklyFromDaily 从日分聚合成周分
func AggregateWeeklyFromDaily(chart *models.NatalChart, weekStart time.Time) *AggregatedScore {
	// 生成7天的日分
	dailyScores := make([]*AggregatedScore, 7)

	for day := 0; day < 7; day++ {
		dayTime := weekStart.AddDate(0, 0, day)
		dailyScores[day] = AggregateDailyFromHourly(chart, dayTime)
	}

	// 提取特征
	overallValues := make([]float64, 7)
	for i, s := range dailyScores {
		overallValues[i] = s.Overall
	}
	overallFeatures := extractFeatures(overallValues)

	// 聚合各维度
	careerValues := make([]float64, 7)
	relationshipValues := make([]float64, 7)
	healthValues := make([]float64, 7)
	financeValues := make([]float64, 7)
	spiritualValues := make([]float64, 7)

	for i, s := range dailyScores {
		careerValues[i] = s.Dimensions.Career
		relationshipValues[i] = s.Dimensions.Relationship
		healthValues[i] = s.Dimensions.Health
		financeValues[i] = s.Dimensions.Finance
		spiritualValues[i] = s.Dimensions.Spiritual
	}

	dimensions := models.DimensionScoresV2{
		Career:       aggregateWithFeatures(careerValues),
		Relationship: aggregateWithFeatures(relationshipValues),
		Health:       aggregateWithFeatures(healthValues),
		Finance:      aggregateWithFeatures(financeValues),
		Spiritual:    aggregateWithFeatures(spiritualValues),
	}

	overall := aggregateWithFeatures(overallValues)

	// 转换 SubScores
	subScores := make([]AggregatedScore, 7)
	for i, s := range dailyScores {
		subScores[i] = *s
	}

	return &AggregatedScore{
		Overall:     overall,
		Dimensions:  dimensions,
		Features:    overallFeatures,
		StartTime:   weekStart,
		EndTime:     weekStart.AddDate(0, 0, 7),
		Granularity: "weekly",
		SubScores:   subScores,
	}
}

// ==================== 月分聚合 ====================

// AggregateMonthlyFromDaily 从日分聚合成月分
func AggregateMonthlyFromDaily(chart *models.NatalChart, year int, month time.Month) *AggregatedScore {
	// 获取该月的天数
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	daysInMonth := lastOfMonth.Day()

	// 生成每天的日分
	dailyScores := make([]*AggregatedScore, daysInMonth)
	for day := 0; day < daysInMonth; day++ {
		dayTime := firstOfMonth.AddDate(0, 0, day)
		dailyScores[day] = AggregateDailyFromHourly(chart, dayTime)
	}

	// 提取特征
	overallValues := make([]float64, daysInMonth)
	for i, s := range dailyScores {
		overallValues[i] = s.Overall
	}
	overallFeatures := extractFeatures(overallValues)

	// 聚合各维度
	careerValues := make([]float64, daysInMonth)
	relationshipValues := make([]float64, daysInMonth)
	healthValues := make([]float64, daysInMonth)
	financeValues := make([]float64, daysInMonth)
	spiritualValues := make([]float64, daysInMonth)

	for i, s := range dailyScores {
		careerValues[i] = s.Dimensions.Career
		relationshipValues[i] = s.Dimensions.Relationship
		healthValues[i] = s.Dimensions.Health
		financeValues[i] = s.Dimensions.Finance
		spiritualValues[i] = s.Dimensions.Spiritual
	}

	dimensions := models.DimensionScoresV2{
		Career:       aggregateWithFeatures(careerValues),
		Relationship: aggregateWithFeatures(relationshipValues),
		Health:       aggregateWithFeatures(healthValues),
		Finance:      aggregateWithFeatures(financeValues),
		Spiritual:    aggregateWithFeatures(spiritualValues),
	}

	overall := aggregateWithFeatures(overallValues)

	return &AggregatedScore{
		Overall:     overall,
		Dimensions:  dimensions,
		Features:    overallFeatures,
		StartTime:   firstOfMonth,
		EndTime:     lastOfMonth.AddDate(0, 0, 1),
		Granularity: "monthly",
	}
}

// ==================== 年分聚合 ====================

// AggregateYearlyFromMonthly 从月分聚合成年分
func AggregateYearlyFromMonthly(chart *models.NatalChart, year int) *AggregatedScore {
	// 生成12个月的月分
	monthlyScores := make([]*AggregatedScore, 12)

	for m := 0; m < 12; m++ {
		month := time.Month(m + 1)
		monthlyScores[m] = AggregateMonthlyFromDaily(chart, year, month)
	}

	// 提取特征
	overallValues := make([]float64, 12)
	for i, s := range monthlyScores {
		overallValues[i] = s.Overall
	}
	overallFeatures := extractFeatures(overallValues)

	// 聚合各维度
	careerValues := make([]float64, 12)
	relationshipValues := make([]float64, 12)
	healthValues := make([]float64, 12)
	financeValues := make([]float64, 12)
	spiritualValues := make([]float64, 12)

	for i, s := range monthlyScores {
		careerValues[i] = s.Dimensions.Career
		relationshipValues[i] = s.Dimensions.Relationship
		healthValues[i] = s.Dimensions.Health
		financeValues[i] = s.Dimensions.Finance
		spiritualValues[i] = s.Dimensions.Spiritual
	}

	dimensions := models.DimensionScoresV2{
		Career:       aggregateWithFeatures(careerValues),
		Relationship: aggregateWithFeatures(relationshipValues),
		Health:       aggregateWithFeatures(healthValues),
		Finance:      aggregateWithFeatures(financeValues),
		Spiritual:    aggregateWithFeatures(spiritualValues),
	}

	overall := aggregateWithFeatures(overallValues)

	// 转换 SubScores
	subScores := make([]AggregatedScore, 12)
	for i, s := range monthlyScores {
		subScores[i] = *s
	}

	return &AggregatedScore{
		Overall:     overall,
		Dimensions:  dimensions,
		Features:    overallFeatures,
		StartTime:   time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:     time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC),
		Granularity: "yearly",
		SubScores:   subScores,
	}
}

// ==================== 特征提取函数 ====================

// extractFeatures 从数值数组中提取统计特征
func extractFeatures(values []float64) ScoreFeatures {
	if len(values) == 0 {
		return ScoreFeatures{}
	}

	n := float64(len(values))

	// 平均值
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	avg := sum / n

	// 最大最小值
	maxVal := values[0]
	minVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
		if v < minVal {
			minVal = v
		}
	}

	// 趋势（线性回归斜率）
	// y = ax + b, 求 a
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	for i, v := range values {
		x := float64(i)
		sumX += x
		sumY += v
		sumXY += x * v
		sumX2 += x * x
	}
	denominator := n*sumX2 - sumX*sumX
	var trend float64
	if denominator != 0 {
		trend = (n*sumXY - sumX*sumY) / denominator
	}

	// 标准差
	sumSq := 0.0
	for _, v := range values {
		diff := v - avg
		sumSq += diff * diff
	}
	stdDev := math.Sqrt(sumSq / n)

	// 稳定性（标准差的倒数，上限100）
	var stability float64
	if stdDev > 0 {
		stability = math.Min(100, 10/stdDev)
	} else {
		stability = 100 // 完全稳定
	}

	return ScoreFeatures{
		Average:   math.Round(avg*10000) / 10000,
		Max:       math.Round(maxVal*10000) / 10000,
		Min:       math.Round(minVal*10000) / 10000,
		Trend:     math.Round(trend*10000) / 10000,
		Stability: math.Round(stability*10000) / 10000,
	}
}

// aggregateWithFeatures 使用特征加权聚合
// 公式：result = average * 0.6 + max * 0.15 + min * 0.10 + (average + trend*5) * 0.15
func aggregateWithFeatures(values []float64) float64 {
	features := extractFeatures(values)

	// 加权组合
	result := features.Average*0.6 +
		features.Max*0.15 +
		features.Min*0.10 +
		(features.Average+features.Trend*5)*0.15

	// 标准化并保留4位小数
	normalized := NormalizeScoreV2(result)
	return math.Round(normalized*10000) / 10000
}

// ==================== 快速查询接口 ====================

// GetScoreAtTime 获取某个时间点的分数（自动选择合适的粒度）
func GetScoreAtTime(chart *models.NatalChart, t time.Time, granularity string) *AggregatedScore {
	switch granularity {
	case "hourly":
		result := CalculateHourlyScore(chart, t)
		return &AggregatedScore{
			Overall:     result.Overall,
			Dimensions:  result.Dimensions,
			StartTime:   t,
			EndTime:     t.Add(time.Hour),
			Granularity: "hourly",
		}
	case "daily":
		return AggregateDailyFromHourly(chart, t)
	case "weekly":
		// 找到周一
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekStart := t.AddDate(0, 0, -(weekday - 1))
		return AggregateWeeklyFromDaily(chart, weekStart)
	case "monthly":
		return AggregateMonthlyFromDaily(chart, t.Year(), t.Month())
	case "yearly":
		return AggregateYearlyFromMonthly(chart, t.Year())
	default:
		return AggregateDailyFromHourly(chart, t)
	}
}

