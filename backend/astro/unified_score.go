package astro

import (
	"math"
	"star/models"
	"time"
)

/*
统一分数计算架构

层级聚合关系：
  小时分 → 日分(24小时平均) → 周分(7日平均) → 月分(~30日平均) → 年分(12月平均)

五维分数独立计算：
  每个维度基于行星-维度映射矩阵独立计算，不是简单的综合分乘以系数
*/

// =============================================================================
// 行星-维度影响矩阵
// =============================================================================

// PlanetDimensionWeight 定义每个行星对五个维度的影响权重
var PlanetDimensionWeight = map[models.PlanetID]map[string]float64{
	models.Sun: {
		"career": 0.8, "relationship": 0.3, "health": 0.5, "finance": 0.4, "spiritual": 0.6,
	},
	models.Moon: {
		"career": 0.2, "relationship": 0.8, "health": 0.6, "finance": 0.3, "spiritual": 0.7,
	},
	models.Mercury: {
		"career": 0.7, "relationship": 0.5, "health": 0.3, "finance": 0.6, "spiritual": 0.4,
	},
	models.Venus: {
		"career": 0.3, "relationship": 0.9, "health": 0.4, "finance": 0.7, "spiritual": 0.5,
	},
	models.Mars: {
		"career": 0.8, "relationship": 0.4, "health": 0.7, "finance": 0.5, "spiritual": 0.3,
	},
	models.Jupiter: {
		"career": 0.7, "relationship": 0.5, "health": 0.6, "finance": 0.9, "spiritual": 0.8,
	},
	models.Saturn: {
		"career": 0.9, "relationship": 0.3, "health": 0.5, "finance": 0.6, "spiritual": 0.4,
	},
	models.Uranus: {
		"career": 0.5, "relationship": 0.4, "health": 0.3, "finance": 0.4, "spiritual": 0.7,
	},
	models.Neptune: {
		"career": 0.2, "relationship": 0.6, "health": 0.4, "finance": 0.2, "spiritual": 0.9,
	},
	models.Pluto: {
		"career": 0.6, "relationship": 0.5, "health": 0.4, "finance": 0.5, "spiritual": 0.8,
	},
}

// =============================================================================
// 统一分数结果结构
// =============================================================================

// UnifiedScore 统一分数结果
type UnifiedScore struct {
	Overall    float64            `json:"overall"`
	Dimensions map[string]float64 `json:"dimensions"`
	RawValue   float64            `json:"rawValue"`
	Factors    []FactorDetail     `json:"factors,omitempty"`
}

// FactorDetail 因子详情
type FactorDetail struct {
	Name       string             `json:"name"`
	Type       string             `json:"type"`
	Value      float64            `json:"value"`
	Dimensions map[string]float64 `json:"dimensions"` // 该因子对各维度的贡献
}

// =============================================================================
// 核心计算函数
// =============================================================================

// CalculateUnifiedHourlyScore 计算小时级别分数（最细粒度，是所有计算的基础）
// 默认使用空用户ID（不应用自定义因子）
func CalculateUnifiedHourlyScore(chart *models.NatalChart, t time.Time) UnifiedScore {
	return CalculateUnifiedHourlyScoreWithUser(chart, t, "")
}

// CalculateUnifiedHourlyScoreWithUser 计算小时级别分数（支持自定义因子）
func CalculateUnifiedHourlyScoreWithUser(chart *models.NatalChart, t time.Time, userID string) UnifiedScore {
	// 1. 获取行运位置
	transitPositions := GetTransitPositions(t)

	// 2. 计算行运相位
	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)

	// 3. 计算每个维度的相位分数
	dimensionScores := calculateDimensionAspectScores(aspects)

	// 4. 计算影响因子并分配到各维度
	factors := CalculateInfluenceFactors(chart, t, transitPositions)
	factorDetails := distributeFactorsToDimensions(factors)

	// 5. 应用因子调整到各维度
	for _, fd := range factorDetails {
		for dim, contrib := range fd.Dimensions {
			dimensionScores[dim] += contrib
		}
	}

	// 6. 应用自定义因子（如果有用户ID）
	if userID != "" {
		// 转换维度分数为 DimensionScoresV2 格式
		dimScoresV2 := &models.DimensionScoresV2{
			Career:       dimensionScores["career"],
			Relationship: dimensionScores["relationship"],
			Health:       dimensionScores["health"],
			Finance:      dimensionScores["finance"],
			Spiritual:    dimensionScores["spiritual"],
		}
		
		// 应用自定义因子
		var overallForCustom float64 = 0 // 用于接收综合分的自定义因子
		ApplyCustomFactors(userID, dimScoresV2, &overallForCustom, t)
		
		// 将修改后的分数写回
		dimensionScores["career"] = dimScoresV2.Career
		dimensionScores["relationship"] = dimScoresV2.Relationship
		dimensionScores["health"] = dimScoresV2.Health
		dimensionScores["finance"] = dimScoresV2.Finance
		dimensionScores["spiritual"] = dimScoresV2.Spiritual
	}

	// 7. 标准化各维度分数
	normalizedDimensions := make(map[string]float64)
	for dim, score := range dimensionScores {
		normalizedDimensions[dim] = NormalizeScoreV2(score)
	}

	// 8. 计算综合分（五维度加权平均）
	overallRaw := (dimensionScores["career"]*0.2 +
		dimensionScores["relationship"]*0.2 +
		dimensionScores["health"]*0.2 +
		dimensionScores["finance"]*0.2 +
		dimensionScores["spiritual"]*0.2)

	overall := NormalizeScoreV2(overallRaw)

	// 9. 应用视觉抖动
	jitterSeed := GenerateSeedFromTime(t.Year(), int(t.Month()), t.Day(), t.Hour())
	overall = ApplyJitter(overall, jitterSeed)
	for dim := range normalizedDimensions {
		normalizedDimensions[dim] = ApplyJitter(normalizedDimensions[dim], jitterSeed+int64(dim[0]))
	}

	// 10. 添加活跃的自定义因子到返回结果
	if userID != "" {
		activeCustomFactors := GetActiveCustomFactors(userID, t)
		for _, cf := range activeCustomFactors {
			factorDetails = append(factorDetails, FactorDetail{
				Name:  cf.Description,
				Type:  "custom",
				Value: cf.Value,
				Dimensions: map[string]float64{
					cf.Dimension: cf.Value,
				},
			})
		}
	}

	return UnifiedScore{
		Overall:    overall,
		Dimensions: normalizedDimensions,
		RawValue:   overallRaw,
		Factors:    factorDetails,
	}
}

// CalculateDailyScore 计算日级别分数（24小时平均）
func CalculateDailyScore(chart *models.NatalChart, date time.Time) UnifiedScore {
	// 设置为当天 00:00
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// 计算24个小时的分数
	var totalOverall float64
	totalDimensions := make(map[string]float64)
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}
	for _, d := range dimensions {
		totalDimensions[d] = 0
	}

	for hour := 0; hour < 24; hour++ {
		t := startOfDay.Add(time.Duration(hour) * time.Hour)
		hourlyScore := CalculateUnifiedHourlyScore(chart, t)

		totalOverall += hourlyScore.Overall
		for _, d := range dimensions {
			totalDimensions[d] += hourlyScore.Dimensions[d]
		}
	}

	// 计算平均值
	avgOverall := totalOverall / 24
	avgDimensions := make(map[string]float64)
	for _, d := range dimensions {
		avgDimensions[d] = totalDimensions[d] / 24
	}

	return UnifiedScore{
		Overall:    math.Round(avgOverall*10000) / 10000,
		Dimensions: avgDimensions,
		RawValue:   avgOverall, // 日级别的 rawValue 就是平均值
	}
}

// CalculateWeeklyScore 计算周级别分数（7天平均）
func CalculateWeeklyScore(chart *models.NatalChart, startDate time.Time) UnifiedScore {
	// 确保从周一开始
	weekday := int(startDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := startDate.AddDate(0, 0, -(weekday - 1))
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())

	// 计算7天的分数
	var totalOverall float64
	totalDimensions := make(map[string]float64)
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}

	for day := 0; day < 7; day++ {
		date := monday.AddDate(0, 0, day)
		dailyScore := CalculateDailyScore(chart, date)

		totalOverall += dailyScore.Overall
		for _, d := range dimensions {
			totalDimensions[d] += dailyScore.Dimensions[d]
		}
	}

	// 计算平均值
	avgOverall := totalOverall / 7
	avgDimensions := make(map[string]float64)
	for _, d := range dimensions {
		avgDimensions[d] = totalDimensions[d] / 7
	}

	return UnifiedScore{
		Overall:    math.Round(avgOverall*10000) / 10000,
		Dimensions: avgDimensions,
		RawValue:   avgOverall,
	}
}

// CalculateMonthlyScore 计算月级别分数（当月所有天平均）
func CalculateMonthlyScore(chart *models.NatalChart, year int, month time.Month) UnifiedScore {
	// 获取当月天数
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	daysInMonth := lastDay.Day()

	// 计算每天的分数
	var totalOverall float64
	totalDimensions := make(map[string]float64)
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}

	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dailyScore := CalculateDailyScore(chart, date)

		totalOverall += dailyScore.Overall
		for _, d := range dimensions {
			totalDimensions[d] += dailyScore.Dimensions[d]
		}
	}

	// 计算平均值
	avgOverall := totalOverall / float64(daysInMonth)
	avgDimensions := make(map[string]float64)
	for _, d := range dimensions {
		avgDimensions[d] = totalDimensions[d] / float64(daysInMonth)
	}

	return UnifiedScore{
		Overall:    math.Round(avgOverall*10000) / 10000,
		Dimensions: avgDimensions,
		RawValue:   avgOverall,
	}
}

// CalculateYearlyScore 计算年级别分数（12个月平均）
func CalculateYearlyScore(chart *models.NatalChart, year int) UnifiedScore {
	var totalOverall float64
	totalDimensions := make(map[string]float64)
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}

	for month := time.January; month <= time.December; month++ {
		monthlyScore := CalculateMonthlyScore(chart, year, month)

		totalOverall += monthlyScore.Overall
		for _, d := range dimensions {
			totalDimensions[d] += monthlyScore.Dimensions[d]
		}
	}

	// 计算平均值
	avgOverall := totalOverall / 12
	avgDimensions := make(map[string]float64)
	for _, d := range dimensions {
		avgDimensions[d] = totalDimensions[d] / 12
	}

	return UnifiedScore{
		Overall:    math.Round(avgOverall*10000) / 10000,
		Dimensions: avgDimensions,
		RawValue:   avgOverall,
	}
}

// =============================================================================
// 辅助函数
// =============================================================================

// calculateDimensionAspectScores 计算每个维度的相位分数
func calculateDimensionAspectScores(aspects []models.AspectData) map[string]float64 {
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}
	scores := make(map[string]float64)

	// 初始化基础分
	for _, d := range dimensions {
		scores[d] = 50 // 基础分50
	}

	// 遍历每个相位
	for _, aspect := range aspects {
		// 获取行运行星和本命行星
		transitPlanet := aspect.Planet1
		natalPlanet := aspect.Planet2

		// 计算相位分值
		aspectValue := getUnifiedAspectValue(string(aspect.AspectType), aspect.Orb)

		// 根据行星-维度权重分配分值
		for _, d := range dimensions {
			// 行运行星权重
			transitWeight := 0.5
			if w, ok := PlanetDimensionWeight[transitPlanet]; ok {
				if dw, ok := w[d]; ok {
					transitWeight = dw
				}
			}

			// 本命行星权重
			natalWeight := 0.5
			if w, ok := PlanetDimensionWeight[natalPlanet]; ok {
				if dw, ok := w[d]; ok {
					natalWeight = dw
				}
			}

			// 综合权重
			combinedWeight := (transitWeight + natalWeight) / 2
			scores[d] += aspectValue * combinedWeight
		}
	}

	return scores
}

// getUnifiedAspectValue 获取相位的分值
func getUnifiedAspectValue(aspectType string, orb float64) float64 {
	// 基础分值
	baseValue := 0.0
	switch aspectType {
	case "conjunction":
		baseValue = 3.0 // 合相：中性偏强
	case "sextile":
		baseValue = 2.0 // 六分：和谐
	case "square":
		baseValue = -2.5 // 四分：紧张
	case "trine":
		baseValue = 3.0 // 三分：和谐
	case "opposition":
		baseValue = -2.0 // 对分：紧张
	}

	// 根据容许度衰减
	maxOrb := 10.0
	orbFactor := 1.0 - (orb / maxOrb)
	if orbFactor < 0 {
		orbFactor = 0
	}

	return baseValue * orbFactor
}

// distributeFactorsToDimensions 将因子分配到各维度
func distributeFactorsToDimensions(factors *models.FactorResult) []FactorDetail {
	if factors == nil {
		return nil
	}

	var details []FactorDetail
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}

	for _, f := range factors.Factors {
		detail := FactorDetail{
			Name:       f.Name,
			Type:       string(f.Type),
			Value:      f.Adjustment,
			Dimensions: make(map[string]float64),
		}

		// 根据因子来源行星分配到各维度
		sourcePlanet := f.SourcePlanet
		if sourcePlanet == "" {
			// 如果没有来源行星，平均分配
			for _, d := range dimensions {
				detail.Dimensions[d] = f.Adjustment / 5
			}
		} else {
			// 根据行星-维度权重分配
			planetID := models.PlanetID(sourcePlanet)
			if weights, ok := PlanetDimensionWeight[planetID]; ok {
				totalWeight := 0.0
				for _, w := range weights {
					totalWeight += w
				}
				for _, d := range dimensions {
					if w, ok := weights[d]; ok {
						detail.Dimensions[d] = f.Adjustment * (w / totalWeight)
					}
				}
			} else {
				// 平均分配
				for _, d := range dimensions {
					detail.Dimensions[d] = f.Adjustment / 5
				}
			}
		}

		details = append(details, detail)
	}

	return details
}

