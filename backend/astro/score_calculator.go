package astro

import (
	"math"
	"star/models"
	"time"
)

// ==================== 分数计算核心模块 ====================
// 实现：因子 → 维度分数 → 综合分数 的完整流程

// ScoreResult 分数计算结果
type ScoreResult struct {
	// 综合分数（由维度加权得出）
	Overall float64 `json:"overall"`

	// 五维度分数
	Dimensions models.DimensionScoresV2 `json:"dimensions"`

	// 本命盘基础分
	BaseScores models.NatalBaseScores `json:"baseScores"`

	// 因子计算结果
	Factors *models.FactorResult `json:"factors"`

	// 时间戳
	Timestamp time.Time `json:"timestamp"`
}

// CalculateScoresV2 计算某时刻的完整分数（新版）
func CalculateScoresV2(chart *models.NatalChart, date time.Time) *ScoreResult {
	// 1. 计算本命盘基础分
	baseScores := CalculateNatalBaseScores(chart)

	// 2. 获取行运位置
	transitPositions := GetTransitPositions(date)

	// 3. 计算所有影响因子（新版）
	factors := CalculateInfluenceFactorsV2(chart, date, transitPositions)

	// 4. 根据因子计算维度分数
	dimensions := calculateDimensionScoresFromFactors(baseScores, factors)

	// 5. 根据维度分数计算综合分数
	overall := calculateOverallFromDimensions(dimensions)

	return &ScoreResult{
		Overall:    overall,
		Dimensions: dimensions,
		BaseScores: baseScores,
		Factors:    factors,
		Timestamp:  date,
	}
}

// calculateDimensionScoresFromFactors 根据因子计算维度分数
func calculateDimensionScoresFromFactors(baseScores models.NatalBaseScores, factors *models.FactorResult) models.DimensionScoresV2 {
	// 从基础分开始
	scores := models.DimensionScoresV2{
		Career:       baseScores.Career,
		Relationship: baseScores.Relationship,
		Health:       baseScores.Health,
		Finance:      baseScores.Finance,
		Spiritual:    baseScores.Spiritual,
	}

	// 累加因子调整
	if factors != nil {
		scores.Career += factors.DimensionAdjustments.Career
		scores.Relationship += factors.DimensionAdjustments.Relationship
		scores.Health += factors.DimensionAdjustments.Health
		scores.Finance += factors.DimensionAdjustments.Finance
		scores.Spiritual += factors.DimensionAdjustments.Spiritual
	}

	// 标准化到 0-100
	scores.Career = NormalizeScoreV2(scores.Career)
	scores.Relationship = NormalizeScoreV2(scores.Relationship)
	scores.Health = NormalizeScoreV2(scores.Health)
	scores.Finance = NormalizeScoreV2(scores.Finance)
	scores.Spiritual = NormalizeScoreV2(scores.Spiritual)

	return scores
}

// calculateOverallFromDimensions 根据维度分数计算综合分数
func calculateOverallFromDimensions(dimensions models.DimensionScoresV2) float64 {
	weights := DefaultDimensionWeights

	overall := dimensions.Career*weights.Career +
		dimensions.Relationship*weights.Relationship +
		dimensions.Health*weights.Health +
		dimensions.Finance*weights.Finance +
		dimensions.Spiritual*weights.Spiritual

	return math.Round(overall*10000) / 10000 // 保留4位小数
}

// ==================== 标准化函数 ====================

// NormalizeScoreV2 改进的分数标准化函数
// 使用 tanh 压缩，确保输出在 0-100 范围内
func NormalizeScoreV2(raw float64) float64 {
	// 原始分数范围约为 0-100（基础50 ± 调整值）
	// 但极端情况可能超出范围

	// 偏移到以 50 为中心
	centered := raw - 50

	// 缩放因子：控制曲线陡峭程度
	// scale=30 时，±45 的原始偏移 → 约 10-90 的输出
	scale := 30.0

	// tanh 压缩
	compressed := math.Tanh(centered / scale)

	// 映射到 0-100
	normalized := 50 + 50*compressed

	// 确保边界并保留4位小数
	result := math.Max(0, math.Min(100, normalized))
	return math.Round(result*10000) / 10000
}

// ==================== 影响因子计算（新版） ====================

// CalculateInfluenceFactorsV2 计算影响因子（新版）
func CalculateInfluenceFactorsV2(chart *models.NatalChart, date time.Time, transitPositions []models.PlanetPosition) *models.FactorResult {
	weights := DefaultFactorWeights
	var factors []models.InfluenceFactor

	// 1. 尊贵度因子
	dignityFactors := calculateDignityFactorsV2(transitPositions, weights.Dignity)
	factors = append(factors, dignityFactors...)

	// 2. 逆行因子
	retrogradeFactors := calculateRetrogradeFactorsV2(transitPositions, weights.Retrograde, date)
	factors = append(factors, retrogradeFactors...)

	// 3. 相位因子
	aspectFactors := calculateAspectFactorsV2(chart, transitPositions, weights.AspectPhase, date)
	factors = append(factors, aspectFactors...)

	// 4. 月相因子
	lunarPhaseFactors := calculateLunarPhaseFactorsV2(transitPositions, weights.LunarPhase, date)
	factors = append(factors, lunarPhaseFactors...)

	// 5. 行星时因子
	planetaryHourFactors := calculatePlanetaryHourFactorsV2(chart, date, weights.PlanetaryHour)
	factors = append(factors, planetaryHourFactors...)

	// 6. 年主星因子
	profectionFactors := calculateProfectionLordFactorsV2(chart, date, transitPositions, weights.ProfectionLord)
	factors = append(factors, profectionFactors...)

	// 7. 月亮空亡因子
	jd := DateToJulianDay(date)
	vocFactors := calculateVoidOfCourseFactorsV2(jd, weights.VoidOfCourse, date)
	factors = append(factors, vocFactors...)

	// 构建结果
	return buildFactorResult(factors, date)
}

// buildFactorResult 构建因子计算结果
func buildFactorResult(factors []models.InfluenceFactor, date time.Time) *models.FactorResult {
	result := &models.FactorResult{
		Factors:              factors,
		YearlyFactors:        []models.InfluenceFactor{},
		MonthlyFactors:       []models.InfluenceFactor{},
		WeeklyFactors:        []models.InfluenceFactor{},
		DailyFactors:         []models.InfluenceFactor{},
		HourlyFactors:        []models.InfluenceFactor{},
		PositiveFactors:      []models.InfluenceFactor{},
		NegativeFactors:      []models.InfluenceFactor{},
		DimensionAdjustments: models.DimensionScoresV2{},
	}

	for i := range factors {
		factor := &factors[i]

		// 生成ID（基于类型和名称，确保去重）
		// 确保 Name 存在再生成 ID
		if factor.Name != "" {
			factor.ID = generateFactorID(factor.Type, factor.Name)
		} else {
			// 如果 Name 为空，使用备用方案
			factor.ID = string(factor.Type) + "_unnamed_" + itoa(i)
		}

		// 计算当前强度（根据生命周期）
		factor.CurrentStrength = CalculateFactorStrength(factor.Lifecycle, date)

		// 计算最终调整值
		factor.Adjustment = factor.BaseValue * factor.Weight * factor.CurrentStrength

		// 分配到各维度
		dimAdj := calculateDimensionAdjustment(factor)
		result.DimensionAdjustments.Career += dimAdj.Career
		result.DimensionAdjustments.Relationship += dimAdj.Relationship
		result.DimensionAdjustments.Health += dimAdj.Health
		result.DimensionAdjustments.Finance += dimAdj.Finance
		result.DimensionAdjustments.Spiritual += dimAdj.Spiritual

		// 累计总调整
		result.TotalAdjustment += factor.Adjustment

		// 按时间级别分组
		switch factor.TimeLevel {
		case models.TimeLevelYearly:
			result.YearlyFactors = append(result.YearlyFactors, *factor)
		case models.TimeLevelMonthly:
			result.MonthlyFactors = append(result.MonthlyFactors, *factor)
		case models.TimeLevelWeekly:
			result.WeeklyFactors = append(result.WeeklyFactors, *factor)
		case models.TimeLevelDaily:
			result.DailyFactors = append(result.DailyFactors, *factor)
		case models.TimeLevelHourly:
			result.HourlyFactors = append(result.HourlyFactors, *factor)
		}

		// 按正负分组
		if factor.IsPositive {
			result.PositiveFactors = append(result.PositiveFactors, *factor)
		} else {
			result.NegativeFactors = append(result.NegativeFactors, *factor)
		}
	}

	return result
}

// calculateDimensionAdjustment 计算因子对各维度的调整
func calculateDimensionAdjustment(factor *models.InfluenceFactor) models.DimensionScoresV2 {
	return models.DimensionScoresV2{
		Career:       factor.Adjustment * factor.DimensionImpact.Career,
		Relationship: factor.Adjustment * factor.DimensionImpact.Relationship,
		Health:       factor.Adjustment * factor.DimensionImpact.Health,
		Finance:      factor.Adjustment * factor.DimensionImpact.Finance,
		Spiritual:    factor.Adjustment * factor.DimensionImpact.Spiritual,
	}
}

// generateFactorID 生成因子ID（基于类型和名称，确保唯一性）
func generateFactorID(factorType models.InfluenceFactorType, name string) string {
	// 使用类型+名称作为唯一标识，确保去重
	return string(factorType) + "_" + name
}

// itoa 简单的整数转字符串
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

// ==================== 各类因子计算（新版） ====================

// calculateDignityFactorsV2 计算尊贵度因子（新版）
func calculateDignityFactorsV2(transitPositions []models.PlanetPosition, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	for _, p := range transitPositions {
		dignity := GetDignity(p.ID, p.Sign)
		score := GetDignityScore(dignity)

		if score == 0 {
			continue
		}

		var name, description, reason string
		var isPositive bool

		planetInfo := GetPlanetInfo(p.ID)
		zodiacInfo := GetZodiacInfo(p.Sign)

		switch dignity {
		case models.DignityDomicile:
			name = planetInfo.Name + " in Domicile"
			description = planetInfo.Name + " in " + zodiacInfo.Name + " - domicile position, powerful energy"
			reason = "Planet in its ruling sign, most natural and powerful expression"
			isPositive = true
		case models.DignityExaltation:
			name = planetInfo.Name + " Exalted"
			description = planetInfo.Name + " in " + zodiacInfo.Name + " - exalted position, enhanced energy"
			reason = "Planet in the sign that elevates its energy"
			isPositive = true
		case models.DignityDetriment:
			name = planetInfo.Name + " in Detriment"
			description = planetInfo.Name + " in " + zodiacInfo.Name + " - detriment position, restricted energy"
			reason = "Planet in opposite sign, expression is hindered"
			isPositive = false
		case models.DignityFall:
			name = planetInfo.Name + " in Fall"
			description = planetInfo.Name + " in " + zodiacInfo.Name + " - fall position, requires extra effort"
			reason = "Planet in opposite of exaltation, weakest energy"
			isPositive = false
		}

		factors = append(factors, models.InfluenceFactor{
			Type:            models.FactorDignity,
			Name:            name,
			Description:     description,
			TimeLevel:       models.TimeLevelMonthly,
			BaseValue:       score,
			Weight:          weight,
			DimensionImpact: GetPlanetDimensionImpact(p.ID),
			SourcePlanet:    p.ID,
			IsPositive:      isPositive,
			AstroReason:     reason,
		})
	}

	return factors
}

// calculateRetrogradeFactorsV2 计算逆行因子（新版）
func calculateRetrogradeFactorsV2(transitPositions []models.PlanetPosition, weight float64, date time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	for _, p := range transitPositions {
		if !p.Retrograde {
			continue
		}

		planetInfo := GetPlanetInfo(p.ID)
		value := -2.0
		reason := "During retrograde, energy turns inward, related matters need review"

		// 根据行星类型调整
		switch p.ID {
		case models.Mercury:
			value = -3.0
			reason = "Mercury rules communication, retrograde may cause miscommunication, be cautious with contracts"
		case models.Venus:
			value = -2.5
			reason = "Venus rules relationships and values, retrograde requires reassessing emotional and financial decisions"
		case models.Mars:
			value = -2.5
			reason = "Mars rules action, retrograde may hinder progress, not ideal for new initiatives"
		case models.Saturn, models.Jupiter:
			value = -1.5
			reason = "Outer planet retrograde lasts longer, impact is moderate but persistent"
		case models.Uranus, models.Neptune, models.Pluto:
			value = -1.0
			reason = "Trans-Saturnian planet retrograde nearly half year, normal phenomenon, milder impact"
		}

		// 创建生命周期
		durationDays := GetRetrogradeDuration(p.ID)
		lifecycle := CreateLifecycle(date.AddDate(0, 0, -int(durationDays/2)), durationDays*24)

		factors = append(factors, models.InfluenceFactor{
			Type:            models.FactorRetrograde,
			Name:            planetInfo.Name + " Retrograde",
			Description:     planetInfo.Name + " is retrograde, related areas may need review and adjustment",
			TimeLevel:       models.TimeLevelWeekly,
			Lifecycle:       lifecycle,
			BaseValue:       value,
			Weight:          weight,
			DimensionImpact: GetPlanetDimensionImpact(p.ID),
			SourcePlanet:    p.ID,
			IsPositive:      false,
			AstroReason:     reason,
		})
	}

	return factors
}

// calculateAspectFactorsV2 计算相位因子（新版）
func calculateAspectFactorsV2(chart *models.NatalChart, transitPositions []models.PlanetPosition, weight float64, date time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)

	for _, asp := range aspects {
		if asp.Strength < 0.5 {
			continue
		}

		aspectDef := GetAspectDefinition(asp.AspectType)
		if aspectDef == nil {
			continue
		}

		// 基础值
		baseValue := asp.Weight
		if aspectDef.Nature == "tense" {
			baseValue = -baseValue * 0.7
		}

		// 入相/离相调整
		if !asp.Applying {
			baseValue *= 0.8
		}

		// 容许度强度
		baseValue *= asp.Strength

		// 创建生命周期（基于容许度）
		lifecycle := CalculateAspectLifecycle(asp.Orb, date, asp.Applying)

		transitInfo := GetPlanetInfo(asp.Planet1)
		natalInfo := GetPlanetInfo(asp.Planet2)

		// 合并两颗行星的维度影响
		transitImpact := GetPlanetDimensionImpact(asp.Planet1)
		natalImpact := GetPlanetDimensionImpact(asp.Planet2)
		combinedImpact := models.DimensionImpact{
			Career:       (transitImpact.Career + natalImpact.Career) / 2,
			Relationship: (transitImpact.Relationship + natalImpact.Relationship) / 2,
			Health:       (transitImpact.Health + natalImpact.Health) / 2,
			Finance:      (transitImpact.Finance + natalImpact.Finance) / 2,
			Spiritual:    (transitImpact.Spiritual + natalImpact.Spiritual) / 2,
		}

		// 判断正负：harmonious=正, tense=负, neutral=根据baseValue判断
		isPositive := aspectDef.Nature == "harmonious"
		if aspectDef.Nature == "neutral" {
			// 合相：根据最终 baseValue 正负判断
			// 吉星（木星、金星）的合相倾向正面，凶星倾向负面
			isPositive = baseValue > 0
		}

		factors = append(factors, models.InfluenceFactor{
			Type:            models.FactorAspectPhase,
			Name:            transitInfo.Name + " " + aspectDef.Name + " " + natalInfo.Name,
			Description:     asp.Interpretation,
			TimeLevel:       models.TimeLevelDaily,
			Lifecycle:       lifecycle,
			BaseValue:       baseValue,
			Weight:          weight,
			DimensionImpact: combinedImpact,
			SourcePlanet:    asp.Planet1,
			IsPositive:      isPositive,
			AstroReason:     "Transit " + transitInfo.Name + " forms " + aspectDef.Name + " with natal " + natalInfo.Name,
		})
	}

	return factors
}

// calculateLunarPhaseFactorsV2 计算月相因子（新版）
func calculateLunarPhaseFactorsV2(transitPositions []models.PlanetPosition, weight float64, date time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	var sunLon, moonLon float64
	for _, p := range transitPositions {
		if p.ID == models.Sun {
			sunLon = p.Longitude
		}
		if p.ID == models.Moon {
			moonLon = p.Longitude
		}
	}

	angle := NormalizeAngle(moonLon - sunLon)
	phaseInfo := GetLunarPhase(angle)

	phaseValues := map[string]float64{
		"new":          1.0,
		"crescent":     1.5,
		"firstQuarter": 0.0,
		"gibbous":      1.0,
		"full":         2.0,
		"disseminating": 0.5,
		"lastQuarter":  -1.0,
		"balsamic":     -0.5,
	}

	value := phaseValues[phaseInfo.Phase]

	// 月相周期约3.5天
	lifecycle := CreateLifecycle(date.AddDate(0, 0, -1), 3.5*24)

	// 月相主要影响情感和健康
	moonImpact := GetPlanetDimensionImpact(models.Moon)

	factors = append(factors, models.InfluenceFactor{
		Type:            models.FactorLunarPhase,
		Name:            phaseInfo.Name,
		Description:     "Current lunar phase: " + phaseInfo.Name + ", " + phaseInfo.Keywords[0],
		TimeLevel:       models.TimeLevelDaily,
		Lifecycle:       lifecycle,
		BaseValue:       value,
		Weight:          weight,
		DimensionImpact: moonImpact,
		SourcePlanet:    models.Moon,
		IsPositive:      value > 0,
		AstroReason:     "Lunar phase reflects the relative position of Sun and Moon, affecting mood and energy rhythms",
	})

	return factors
}

// calculatePlanetaryHourFactorsV2 计算行星时因子（新版）
func calculatePlanetaryHourFactorsV2(chart *models.NatalChart, date time.Time, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	lat, lon := 0.0, 0.0
	if chart.BirthData.Latitude != 0 {
		lat = chart.BirthData.Latitude
		lon = chart.BirthData.Longitude
	}

	hourInfo := CalculatePlanetaryHourEnhanced(date, lat, lon)
	hourRulerInfo := GetPlanetInfo(hourInfo.Ruler)
	dayRulerInfo := GetPlanetInfo(hourInfo.DayRuler)

	value := hourInfo.Influence

	// 与命主星相关加成
	if hourInfo.Ruler == chart.ChartRuler {
		value += 2.0
	}
	if hourInfo.DayRuler == chart.ChartRuler {
		value += 1.0
	}

	// 行星时持续约1-1.5小时
	lifecycle := CreateLifecycle(date.Add(-30*time.Minute), 1.5)

	factors = append(factors, models.InfluenceFactor{
		Type:            models.FactorPlanetaryHour,
		Name:            dayRulerInfo.Name + " Day " + hourRulerInfo.Name + " Hour",
		Description:     "Today is " + dayRulerInfo.Name + " day, planetary hour #" + itoa(hourInfo.PlanetaryHour) + " ruled by " + hourRulerInfo.Name,
		TimeLevel:       models.TimeLevelHourly,
		Lifecycle:       lifecycle,
		BaseValue:       value,
		Weight:          weight,
		DimensionImpact: GetPlanetDimensionImpact(hourInfo.Ruler),
		SourcePlanet:    hourInfo.Ruler,
		IsPositive:      value > 0,
		AstroReason:     "Planetary hours originate from ancient Babylon, each hour ruled by a different planet affecting that period's energy",
	})

	return factors
}

// calculateProfectionLordFactorsV2 计算年主星因子（新版）
func calculateProfectionLordFactorsV2(chart *models.NatalChart, date time.Time, transitPositions []models.PlanetPosition, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	profection := GetProfectionForDate(chart, date)
	if profection == nil {
		return factors
	}

	var lordTransit *models.PlanetPosition
	for _, p := range transitPositions {
		if p.ID == profection.LordOfYear {
			lordTransit = &p
			break
		}
	}

	if lordTransit == nil {
		return factors
	}

	lordInfo := GetPlanetInfo(profection.LordOfYear)
	value := 0.0

	// 尊贵度影响
	dignity := GetDignity(lordTransit.ID, lordTransit.Sign)
	value += GetDignityScore(dignity)

	// 逆行影响
	if lordTransit.Retrograde {
		value -= 1.0
	}

	// 年主星影响全年
	lifecycle := CreateLifecycle(date.AddDate(0, -6, 0), 365*24)

	// 年主星影响对应宫位的维度
	houseDimension := GetDimensionForHouseV2(profection.House)
	impact := GetPlanetDimensionImpact(profection.LordOfYear)

	// 增强对应宫位维度的影响
	switch houseDimension {
	case "career":
		impact.Career += 0.2
	case "relationship":
		impact.Relationship += 0.2
	case "health":
		impact.Health += 0.2
	case "finance":
		impact.Finance += 0.2
	case "spiritual":
		impact.Spiritual += 0.2
	}

	// 归一化
	total := impact.Career + impact.Relationship + impact.Health + impact.Finance + impact.Spiritual
	if total > 0 {
		impact.Career /= total
		impact.Relationship /= total
		impact.Health /= total
		impact.Finance /= total
		impact.Spiritual /= total
	}

	factors = append(factors, models.InfluenceFactor{
		Type:            models.FactorProfectionLord,
		Name:            "Annual Lord " + lordInfo.Name + " Status",
		Description:     "This year's annual lord is " + lordInfo.Name + ", current state affects annual theme: " + profection.HouseTheme,
		TimeLevel:       models.TimeLevelYearly,
		Lifecycle:       lifecycle,
		BaseValue:       value,
		Weight:          weight,
		DimensionImpact: impact,
		SourcePlanet:    profection.LordOfYear,
		IsPositive:      value > 0,
		AstroReason:     "Annual Profections is a classical astrology technique, activating different houses each year",
	})

	return factors
}

// calculateVoidOfCourseFactorsV2 计算月亮空亡因子（新版）
func calculateVoidOfCourseFactorsV2(jd float64, weight float64, date time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	vocInfo := CalculateVoidOfCourse(jd, nil)

	if vocInfo.IsVoid {
		// 创建生命周期
		lifecycle := CreateLifecycle(date, vocInfo.Duration)

		factors = append(factors, models.InfluenceFactor{
			Type:            models.FactorVoidOfCourse,
			Name:            "Moon Void of Course",
			Description:     "Moon is void of course, lasting " + formatDuration(vocInfo.Duration) + " before entering " + vocInfo.NextSign + ". Not ideal for starting new matters",
			TimeLevel:       models.TimeLevelHourly,
			Lifecycle:       lifecycle,
			BaseValue:       vocInfo.Influence,
			Weight:          weight,
			DimensionImpact: GetPlanetDimensionImpact(models.Moon),
			SourcePlanet:    models.Moon,
			IsPositive:      false,
			AstroReason:     "Void of course means Moon won't form any major aspects before leaving current sign, traditionally believed that things started now may not develop as expected",
		})
	}

	return factors
}

// formatDuration 格式化持续时间
func formatDuration(hours float64) string {
	if hours < 1 {
		return itoa(int(hours*60)) + " minutes"
	}
	h := int(hours)
	m := int((hours - float64(h)) * 60)
	if m > 0 {
		return itoa(h) + "h " + itoa(m) + "m"
	}
	return itoa(h) + " hours"
}

