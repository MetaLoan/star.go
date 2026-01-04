package astro

import (
	"fmt"
	"star/models"
	"time"
)

// CalculateInfluenceFactors 计算影响因子
// 该系统设计为参数化配置，便于运营后期调整
func CalculateInfluenceFactors(chart *models.NatalChart, date time.Time, transitPositions []models.PlanetPosition) *models.FactorResult {
	weights := DefaultFactorWeights
	var factors []models.InfluenceFactor

	// 1. 尊贵度因子
	dignityFactors := calculateDignityFactors(chart, transitPositions, weights.Dignity)
	factors = append(factors, dignityFactors...)

	// 2. 逆行因子
	retrogradeFactors := calculateRetrogradeFactors(transitPositions, weights.Retrograde)
	factors = append(factors, retrogradeFactors...)

	// 3. 相位因子
	aspectFactors := calculateAspectFactors(chart, transitPositions, weights.AspectPhase, weights.AspectOrb)
	factors = append(factors, aspectFactors...)

	// 4. 月相因子
	lunarPhaseFactors := calculateLunarPhaseFactors(transitPositions, weights.LunarPhase)
	factors = append(factors, lunarPhaseFactors...)

	// 5. 行星时因子（增强版）
	planetaryHourFactors := calculatePlanetaryHourFactorsEnhanced(chart, date, weights.PlanetaryHour)
	factors = append(factors, planetaryHourFactors...)

	// 6. 年主星因子
	profectionFactors := calculateProfectionLordFactors(chart, date, transitPositions, weights.ProfectionLord)
	factors = append(factors, profectionFactors...)

	// 7. 月亮空亡因子
	jd := DateToJulianDay(date)
	vocFactors := calculateVoidOfCourseFactors(jd, weights.LunarPhase)
	factors = append(factors, vocFactors...)

	// 分类正负因子
	var positiveFactors, negativeFactors []models.InfluenceFactor
	var totalAdjustment float64
	dimensionAdjustments := make(map[string]float64)

	for i := range factors {
		factors[i].ID = fmt.Sprintf("%s_%d", factors[i].Type, i)
		totalAdjustment += factors[i].Adjustment

		if factors[i].IsPositive {
			positiveFactors = append(positiveFactors, factors[i])
		} else {
			negativeFactors = append(negativeFactors, factors[i])
		}

		// 累加维度调整
		if factors[i].Dimension != "" {
			dimensionAdjustments[factors[i].Dimension] += factors[i].Adjustment
		}
	}

	return &models.FactorResult{
		Factors:              factors,
		TotalAdjustment:      totalAdjustment,
		PositiveFactors:      positiveFactors,
		NegativeFactors:      negativeFactors,
		DimensionAdjustments: dimensionAdjustments,
	}
}

// calculateDignityFactors 计算尊贵度因子
func calculateDignityFactors(_ *models.NatalChart, transitPositions []models.PlanetPosition, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	for _, p := range transitPositions {
		dignity := GetDignity(p.ID, p.Sign)
		score := GetDignityScore(dignity)

		if score != 0 {
			var name, description string
			var isPositive bool

			planetInfo := GetPlanetInfo(p.ID)
			zodiacInfo := GetZodiacInfo(p.Sign)

			switch dignity {
			case models.DignityDomicile:
				name = fmt.Sprintf("%s入庙", planetInfo.Name)
				description = fmt.Sprintf("%s在%s入庙，能量强大", planetInfo.Name, zodiacInfo.Name)
				isPositive = true
			case models.DignityExaltation:
				name = fmt.Sprintf("%s旺相", planetInfo.Name)
				description = fmt.Sprintf("%s在%s旺相，能量提升", planetInfo.Name, zodiacInfo.Name)
				isPositive = true
			case models.DignityDetriment:
				name = fmt.Sprintf("%s落陷", planetInfo.Name)
				description = fmt.Sprintf("%s在%s落陷，能量受限", planetInfo.Name, zodiacInfo.Name)
				isPositive = false
			case models.DignityFall:
				name = fmt.Sprintf("%s失势", planetInfo.Name)
				description = fmt.Sprintf("%s在%s失势，需要额外努力", planetInfo.Name, zodiacInfo.Name)
				isPositive = false
			}

			factors = append(factors, models.InfluenceFactor{
				Type:        models.FactorDignity,
				Name:        name,
				Value:       score,
				Weight:      weight,
				Adjustment:  score * weight,
				Description: description,
				IsPositive:  isPositive,
			})
		}
	}

	return factors
}

// calculateRetrogradeFactors 计算逆行因子
func calculateRetrogradeFactors(transitPositions []models.PlanetPosition, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	for _, p := range transitPositions {
		if p.Retrograde {
			planetInfo := GetPlanetInfo(p.ID)

			// 逆行影响值
			value := -2.0

			// 外行星逆行影响较大
			if p.ID == models.Saturn || p.ID == models.Uranus ||
				p.ID == models.Neptune || p.ID == models.Pluto {
				value = -1.5 // 外行星逆行常态化，影响稍小
			}

			// 水星逆行特殊处理
			if p.ID == models.Mercury {
				value = -3.0 // 水星逆行影响较大
			}

			factors = append(factors, models.InfluenceFactor{
				Type:        models.FactorRetrograde,
				Name:        fmt.Sprintf("%s逆行", planetInfo.Name),
				Value:       value,
				Weight:      weight,
				Adjustment:  value * weight,
				Description: fmt.Sprintf("%s正在逆行，相关领域可能需要回顾和调整", planetInfo.Name),
				IsPositive:  false,
			})
		}
	}

	return factors
}

// calculateAspectFactors 计算相位因子
func calculateAspectFactors(chart *models.NatalChart, transitPositions []models.PlanetPosition, phaseWeight, _ float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)

	for _, asp := range aspects {
		if asp.Strength < 0.5 {
			continue // 跳过弱相位
		}

		aspectDef := GetAspectDefinition(asp.AspectType)
		if aspectDef == nil {
			continue
		}

		// 相位阶段因子（入相vs离相）
		phaseValue := 1.0
		if !asp.Applying {
			phaseValue = 0.8 // 离相相位影响减弱
		}

		// 容许度因子（越精确越强）
		orbValue := asp.Strength

		// 综合调整值
		baseValue := asp.Weight
		if aspectDef.Nature == "tense" {
			baseValue = -baseValue * 0.7
		}

		adjustment := baseValue * phaseValue * orbValue * phaseWeight

		transitInfo := GetPlanetInfo(asp.Planet1)
		natalInfo := GetPlanetInfo(asp.Planet2)

		factors = append(factors, models.InfluenceFactor{
			Type:        models.FactorAspectPhase,
			Name:        fmt.Sprintf("%s%s%s", transitInfo.Name, aspectDef.Name, natalInfo.Name),
			Value:       baseValue,
			Weight:      phaseWeight,
			Adjustment:  adjustment,
			Description: asp.Interpretation,
			IsPositive:  aspectDef.Nature == "harmonious",
		})
	}

	return factors
}

// calculateLunarPhaseFactors 计算月相因子
func calculateLunarPhaseFactors(transitPositions []models.PlanetPosition, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	// 找出太阳和月亮
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

	// 月相影响值
	phaseValues := map[string]float64{
		"new":          1.0,  // 新月 - 新开始
		"crescent":     1.5,  // 上升期
		"firstQuarter": 0.0,  // 危机调整
		"gibbous":      1.0,  // 准备期
		"full":         2.0,  // 满月 - 高峰
		"disseminating": 0.5, // 分享期
		"lastQuarter":  -1.0, // 释放期
		"balsamic":     -0.5, // 休息期
	}

	value := phaseValues[phaseInfo.Phase]

	factors = append(factors, models.InfluenceFactor{
		Type:        models.FactorLunarPhase,
		Name:        phaseInfo.Name,
		Value:       value,
		Weight:      weight,
		Adjustment:  value * weight,
		Description: fmt.Sprintf("当前月相为%s，%s", phaseInfo.Name, phaseInfo.Keywords[0]),
		IsPositive:  value > 0,
	})

	return factors
}

// calculatePlanetaryHourFactorsEnhanced 计算行星时因子（增强版）
func calculatePlanetaryHourFactorsEnhanced(chart *models.NatalChart, date time.Time, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	// 使用出生地点计算行星时
	lat, lon := 0.0, 0.0
	if chart.BirthData.Latitude != 0 {
		lat = chart.BirthData.Latitude
		lon = chart.BirthData.Longitude
	}

	// 获取当前行星时信息
	hourInfo := CalculatePlanetaryHourEnhanced(date, lat, lon)
	hourRulerInfo := GetPlanetInfo(hourInfo.Ruler)
	dayRulerInfo := GetPlanetInfo(hourInfo.DayRuler)

	// 基础行星时影响
	value := hourInfo.Influence

	// 检查是否与命主星相关
	if hourInfo.Ruler == chart.ChartRuler {
		value += 2.0
	}
	if hourInfo.DayRuler == chart.ChartRuler {
		value += 1.0
	}

	// 行星时因子
	factors = append(factors, models.InfluenceFactor{
		Type:        models.FactorPlanetaryHour,
		Name:        fmt.Sprintf("%s日 %s时", dayRulerInfo.Name, hourRulerInfo.Name),
		Value:       value,
		Weight:      weight,
		Adjustment:  value * weight,
		Description: fmt.Sprintf("今天是%s日，当前是第%d个行星时，由%s主管", dayRulerInfo.Name, hourInfo.PlanetaryHour, hourRulerInfo.Name),
		IsPositive:  value > 0,
	})

	return factors
}

// calculateVoidOfCourseFactors 计算月亮空亡因子
func calculateVoidOfCourseFactors(jd float64, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	vocInfo := CalculateVoidOfCourse(jd, nil)

	if vocInfo.IsVoid {
		// 月亮空亡期间，影响值为负
		factors = append(factors, models.InfluenceFactor{
			Type:        models.FactorLunarPhase,
			Name:        "月亮空亡",
			Value:       vocInfo.Influence,
			Weight:      weight,
			Adjustment:  vocInfo.Influence * weight,
			Description: fmt.Sprintf("月亮空亡中，持续%.1f小时后进入%s。此时不宜开始新事务", vocInfo.Duration, vocInfo.NextSign),
			IsPositive:  false,
		})
	}

	return factors
}

// calculateProfectionLordFactors 计算年主星因子
func calculateProfectionLordFactors(chart *models.NatalChart, date time.Time, transitPositions []models.PlanetPosition, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	// 获取当前年限法
	profection := GetProfectionForDate(chart, date)
	if profection == nil {
		return factors
	}

	// 找出行运中的年主星
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

	// 年主星状态影响
	lordInfo := GetPlanetInfo(profection.LordOfYear)
	value := 0.0

	// 尊贵度影响
	dignity := GetDignity(lordTransit.ID, lordTransit.Sign)
	value += GetDignityScore(dignity)

	// 逆行影响
	if lordTransit.Retrograde {
		value -= 1.0
	}

	factors = append(factors, models.InfluenceFactor{
		Type:        models.FactorProfectionLord,
		Name:        fmt.Sprintf("年主星%s状态", lordInfo.Name),
		Value:       value,
		Weight:      weight,
		Adjustment:  value * weight,
		Description: fmt.Sprintf("今年的年主星是%s，当前状态影响年度主题", lordInfo.Name),
		IsPositive:  value > 0,
		Dimension:   getDimensionForHouse(profection.House),
	})

	return factors
}

// getDimensionForHouse 根据宫位获取对应维度
func getDimensionForHouse(house int) string {
	switch house {
	case 1, 6:
		return "health"
	case 2, 8:
		return "finance"
	case 3, 9:
		return "spiritual"
	case 4, 10:
		return "career"
	case 5, 7, 11:
		return "relationship"
	default:
		return ""
	}
}

// UpdateFactorWeights 更新因子权重（供运营调整）
func UpdateFactorWeights(newWeights models.FactorWeights) {
	DefaultFactorWeights = newWeights
}

// GetCurrentFactorWeights 获取当前因子权重配置
func GetCurrentFactorWeights() models.FactorWeights {
	return DefaultFactorWeights
}

