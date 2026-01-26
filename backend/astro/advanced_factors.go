// advanced_factors.go - 高级占星因子计算
// 实现日月食、交点、燃烧、停滞、恒星、阿拉伯点等高级技术

package astro

import (
	"fmt"
	"math"
	"time"

	"star/models"
)

// ==================== 日月食计算 ====================

// EclipseType 日月食类型
type EclipseType string

const (
	EclipseSolar        EclipseType = "solar"        // 日食
	EclipseLunar        EclipseType = "lunar"        // 月食
	EclipseSolarTotal   EclipseType = "solarTotal"   // 日全食
	EclipseSolarAnnular EclipseType = "solarAnnular" // 日环食
	EclipseLunarTotal   EclipseType = "lunarTotal"   // 月全食
	EclipseLunarPartial EclipseType = "lunarPartial" // 月偏食
)

// CalculateEclipseFactors 计算日月食因子
func CalculateEclipseFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	// 获取太阳和月亮位置
	sunPos := GetPlanetFromChart(chart, models.Sun)
	moonPos := GetPlanetFromChart(chart, models.Moon)
	nodePos := GetPlanetFromChart(chart, models.NorthNode)

	if sunPos == nil || moonPos == nil || nodePos == nil {
		return factors
	}

	sunLong := sunPos.Longitude
	moonLong := moonPos.Longitude
	nodeLong := nodePos.Longitude

	// 计算太阳-月亮角度差
	sunMoonDiff := math.Abs(sunLong - moonLong)
	if sunMoonDiff > 180 {
		sunMoonDiff = 360 - sunMoonDiff
	}

	// 检查是否接近交点（容许度18度内可能发生食相）
	sunNodeDiff := math.Abs(sunLong - nodeLong)
	if sunNodeDiff > 180 {
		sunNodeDiff = 360 - sunNodeDiff
	}
	sunSouthNodeDiff := math.Abs(sunLong - math.Mod(nodeLong+180, 360))
	if sunSouthNodeDiff > 180 {
		sunSouthNodeDiff = 360 - sunSouthNodeDiff
	}

	nearNode := sunNodeDiff < 18 || sunSouthNodeDiff < 18

	// 新月附近（日食可能）
	if sunMoonDiff < 12 && nearNode {
		strength := 1.0 - (sunMoonDiff / 12.0)
		lifecycle := &models.FactorLifecycle{
			StartTime: queryTime.AddDate(0, 0, -14),
			PeakTime:  queryTime,
			EndTime:   queryTime.AddDate(0, 0, 14),
			Duration:  672,
		}
		factor := models.InfluenceFactor{
			ID:          fmt.Sprintf("eclipse_solar_%s", queryTime.Format("200601")),
			Type:        models.FactorEclipse,
			Name:        "日食影响期",
			Description: fmt.Sprintf("太阳与月亮合相于%.1f°，接近月交点，日食能量活跃", sunLong),
			BaseValue:   strength * 3,
			Weight:      1.0,
			Adjustment:  strength * 3,
			IsPositive:  false,
			TimeLevel:   models.TimeLevelMonthly,
			Lifecycle:   lifecycle,
			DimensionImpact: models.DimensionImpact{
				Career: 0.3, Relationship: 0.2, Health: 0.1, Finance: 0.2, Spiritual: 0.2,
			},
		}
		factors = append(factors, factor)
	}

	// 满月附近（月食可能）
	if sunMoonDiff > 168 && sunMoonDiff < 192 && nearNode {
		strength := 1.0 - math.Abs(sunMoonDiff-180)/12.0
		lifecycle := &models.FactorLifecycle{
			StartTime: queryTime.AddDate(0, 0, -7),
			PeakTime:  queryTime,
			EndTime:   queryTime.AddDate(0, 0, 7),
			Duration:  336,
		}
		factor := models.InfluenceFactor{
			ID:          fmt.Sprintf("eclipse_lunar_%s", queryTime.Format("200601")),
			Type:        models.FactorEclipse,
			Name:        "月食影响期",
			Description: "太阳与月亮冲相，接近月交点，月食能量活跃",
			BaseValue:   strength * 2.5,
			Weight:      1.0,
			Adjustment:  strength * 2.5,
			IsPositive:  false,
			TimeLevel:   models.TimeLevelMonthly,
			Lifecycle:   lifecycle,
			DimensionImpact: models.DimensionImpact{
				Career: 0.2, Relationship: 0.3, Health: 0.15, Finance: 0.15, Spiritual: 0.2,
			},
		}
		factors = append(factors, factor)
	}

	return factors
}

// ==================== 月交点计算 ====================

// CalculateLunarNodeFactors 计算月交点因子
func CalculateLunarNodeFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	nodePos := GetPlanetFromChart(chart, models.NorthNode)
	if nodePos == nil {
		return factors
	}

	northNode := nodePos.Longitude
	southNode := math.Mod(northNode+180, 360)

	// 检查各行星与交点的相位
	planets := []models.PlanetID{models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars, models.Jupiter, models.Saturn}

	for _, planetID := range planets {
		planet := GetPlanetFromChart(chart, planetID)
		if planet == nil {
			continue
		}

		// 与北交点的角度差
		diffNN := math.Abs(planet.Longitude - northNode)
		if diffNN > 180 {
			diffNN = 360 - diffNN
		}

		// 与南交点的角度差
		diffSN := math.Abs(planet.Longitude - southNode)
		if diffSN > 180 {
			diffSN = 360 - diffSN
		}

		planetChinese := getPlanetChineseName(planetID)

		// 合相北交点（容许度8度）
		if diffNN < 8 {
			strength := 1.0 - (diffNN / 8.0)
			factor := models.InfluenceFactor{
				ID:           fmt.Sprintf("node_nn_%s_%s", planetID, queryTime.Format("20060102")),
				Type:         models.FactorLunarNode,
				Name:         fmt.Sprintf("%s合北交点", planetChinese),
				Description:  fmt.Sprintf("%s与北交点合相，命运方向指引，未来发展机遇", planetChinese),
				BaseValue:    strength * 2,
				Weight:       0.9,
				Adjustment:   strength * 1.8,
				IsPositive:   true,
				TimeLevel:    models.TimeLevelWeekly,
				SourcePlanet: planetID,
				DimensionImpact: models.DimensionImpact{
					Career: 0.3, Relationship: 0.2, Health: 0.1, Finance: 0.2, Spiritual: 0.2,
				},
			}
			factors = append(factors, factor)
		}

		// 合相南交点（容许度8度）
		if diffSN < 8 {
			strength := 1.0 - (diffSN / 8.0)
			factor := models.InfluenceFactor{
				ID:           fmt.Sprintf("node_sn_%s_%s", planetID, queryTime.Format("20060102")),
				Type:         models.FactorLunarNode,
				Name:         fmt.Sprintf("%s合南交点", planetChinese),
				Description:  fmt.Sprintf("%s与南交点合相，业力功课，需要放下过去模式", planetChinese),
				BaseValue:    strength * -1.5,
				Weight:       0.9,
				Adjustment:   strength * -1.35,
				IsPositive:   false,
				TimeLevel:    models.TimeLevelWeekly,
				SourcePlanet: planetID,
				DimensionImpact: models.DimensionImpact{
					Career: 0.2, Relationship: 0.25, Health: 0.15, Finance: 0.15, Spiritual: 0.25,
				},
			}
			factors = append(factors, factor)
		}
	}

	return factors
}

// ==================== 燃烧计算 ====================

// CalculateCombustionFactors 计算燃烧因子（行星被太阳灼烧）
func CalculateCombustionFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	sunPos := GetPlanetFromChart(chart, models.Sun)
	if sunPos == nil {
		return factors
	}
	sunLong := sunPos.Longitude

	// 燃烧容许度
	combustionOrb := 8.5
	underBeamsOrb := 17.0

	// 可被燃烧的行星
	planets := []models.PlanetID{models.Mercury, models.Venus, models.Mars, models.Jupiter, models.Saturn}

	for _, planetID := range planets {
		planet := GetPlanetFromChart(chart, planetID)
		if planet == nil {
			continue
		}

		diff := math.Abs(planet.Longitude - sunLong)
		if diff > 180 {
			diff = 360 - diff
		}

		planetChinese := getPlanetChineseName(planetID)

		if diff < combustionOrb {
			strength := 1.0 - (diff / combustionOrb)
			factor := models.InfluenceFactor{
				ID:              fmt.Sprintf("combustion_%s_%s", planetID, queryTime.Format("20060102")),
				Type:            models.FactorCombustion,
				Name:            fmt.Sprintf("%s燃烧", planetChinese),
				Description:     fmt.Sprintf("%s距太阳仅%.1f°，被太阳光芒遮蔽，能量受损", planetChinese, diff),
				BaseValue:       strength * -2.5,
				Weight:          1.0,
				Adjustment:      strength * -2.5,
				IsPositive:      false,
				TimeLevel:       models.TimeLevelDaily,
				SourcePlanet:    planetID,
				DimensionImpact: GetPlanetDimensionImpact(planetID),
			}
			factors = append(factors, factor)
		} else if diff < underBeamsOrb {
			strength := 1.0 - ((diff - combustionOrb) / (underBeamsOrb - combustionOrb))
			factor := models.InfluenceFactor{
				ID:              fmt.Sprintf("underbeams_%s_%s", planetID, queryTime.Format("20060102")),
				Type:            models.FactorCombustion,
				Name:            fmt.Sprintf("%s在光下", planetChinese),
				Description:     fmt.Sprintf("%s距太阳%.1f°，处于太阳光芒之下，能量略受影响", planetChinese, diff),
				BaseValue:       strength * -1.0,
				Weight:          0.7,
				Adjustment:      strength * -0.7,
				IsPositive:      false,
				TimeLevel:       models.TimeLevelDaily,
				SourcePlanet:    planetID,
				DimensionImpact: GetPlanetDimensionImpact(planetID),
			}
			factors = append(factors, factor)
		}
	}

	return factors
}

// ==================== 停滞计算 ====================

// CalculateStationFactors 计算停滞因子
// 注意：由于当前 PlanetPosition 不包含速度信息，此函数通过实时计算行星速度来检测停滞
func CalculateStationFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	planets := []models.PlanetID{models.Mercury, models.Venus, models.Mars, models.Jupiter, models.Saturn, models.Uranus, models.Neptune, models.Pluto}

	// 计算当前儒略日
	jd := TimeToJulianDay(queryTime)

	for _, planetID := range planets {
		planet := GetPlanetFromChart(chart, planetID)
		if planet == nil {
			continue
		}

		// 通过计算前后时刻的位置来估算速度
		long1 := GetPlanetLongitudeAt(planetID, jd-0.5) // 12小时前
		long2 := GetPlanetLongitudeAt(planetID, jd+0.5) // 12小时后

		// 计算速度（度/天）
		speed := long2 - long1
		// 处理跨越0度的情况
		if speed > 180 {
			speed -= 360
		} else if speed < -180 {
			speed += 360
		}

		absSpeed := math.Abs(speed)
		stationThreshold := getStationThreshold(planetID)

		if absSpeed < stationThreshold {
			stationType := "顺转逆"
			if planet.Retrograde {
				stationType = "逆转顺"
			}

			planetChinese := getPlanetChineseName(planetID)
			strength := 1.0 - (absSpeed / stationThreshold)
			factor := models.InfluenceFactor{
				ID:              fmt.Sprintf("station_%s_%s", planetID, queryTime.Format("20060102")),
				Type:            models.FactorStation,
				Name:            fmt.Sprintf("%s停滞（%s）", planetChinese, stationType),
				Description:     fmt.Sprintf("%s速度仅%.4f°/天，处于停滞状态，该行星能量极度强化", planetChinese, absSpeed),
				BaseValue:       strength * 2.0,
				Weight:          1.2,
				Adjustment:      strength * 2.4,
				IsPositive:      true,
				TimeLevel:       models.TimeLevelDaily,
				SourcePlanet:    planetID,
				DimensionImpact: GetPlanetDimensionImpact(planetID),
			}
			factors = append(factors, factor)
		}
	}

	return factors
}

// getStationThreshold 获取行星停滞速度阈值
func getStationThreshold(planet models.PlanetID) float64 {
	thresholds := map[models.PlanetID]float64{
		models.Mercury: 0.2,
		models.Venus:   0.1,
		models.Mars:    0.05,
		models.Jupiter: 0.02,
		models.Saturn:  0.01,
		models.Uranus:  0.005,
		models.Neptune: 0.003,
		models.Pluto:   0.002,
	}
	if t, ok := thresholds[planet]; ok {
		return t
	}
	return 0.01
}

// ==================== 互容计算 ====================

// CalculateReceptionFactors 计算互容因子
func CalculateReceptionFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	// 行星守护星座映射
	rulershipMap := map[models.PlanetID][]string{
		models.Sun:     {"Leo"},
		models.Moon:    {"Cancer"},
		models.Mercury: {"Gemini", "Virgo"},
		models.Venus:   {"Taurus", "Libra"},
		models.Mars:    {"Aries", "Scorpio"},
		models.Jupiter: {"Sagittarius", "Pisces"},
		models.Saturn:  {"Capricorn", "Aquarius"},
	}

	planets := []models.PlanetID{models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars, models.Jupiter, models.Saturn}

	for i, planet1ID := range planets {
		planet1 := GetPlanetFromChart(chart, planet1ID)
		if planet1 == nil {
			continue
		}

		for j := i + 1; j < len(planets); j++ {
			planet2ID := planets[j]
			planet2 := GetPlanetFromChart(chart, planet2ID)
			if planet2 == nil {
				continue
			}

			planet1Sign := getSignFromLongitudeInternal(planet1.Longitude)
			planet2Sign := getSignFromLongitudeInternal(planet2.Longitude)

			planet1InPlanet2Signs := containsSign(rulershipMap[planet2ID], planet1Sign)
			planet2InPlanet1Signs := containsSign(rulershipMap[planet1ID], planet2Sign)

			if planet1InPlanet2Signs && planet2InPlanet1Signs {
				planet1Chinese := getPlanetChineseName(planet1ID)
				planet2Chinese := getPlanetChineseName(planet2ID)
				factor := models.InfluenceFactor{
					ID:          fmt.Sprintf("reception_%s_%s_%s", planet1ID, planet2ID, queryTime.Format("20060102")),
					Type:        models.FactorReception,
					Name:        fmt.Sprintf("%s与%s互容", planet1Chinese, planet2Chinese),
					Description: fmt.Sprintf("%s在%s，%s在%s，形成互容关系，双方能量互相支持", planet1Chinese, planet2Sign, planet2Chinese, planet1Sign),
					BaseValue:   2.5,
					Weight:      1.0,
					Adjustment:  2.5,
					IsPositive:  true,
					TimeLevel:   models.TimeLevelMonthly, // 互容持续时间取决于行星在星座时间
					DimensionImpact: models.DimensionImpact{
						Career: 0.2, Relationship: 0.2, Health: 0.2, Finance: 0.2, Spiritual: 0.2,
					},
				}
				factors = append(factors, factor)
			}
		}
	}

	return factors
}

func containsSign(signs []string, sign string) bool {
	for _, s := range signs {
		if s == sign {
			return true
		}
	}
	return false
}

// ==================== 恒星计算 ====================

// FixedStarInfo 恒星信息
type FixedStarInfo struct {
	Name       string  `json:"name"`
	Chinese    string  `json:"chinese"`
	Longitude  float64 `json:"longitude"`
	Magnitude  float64 `json:"magnitude"`
	Nature     string  `json:"nature"`
	IsPositive bool    `json:"isPositive"`
	Keywords   string  `json:"keywords"`
}

// ImportantFixedStars 重要恒星数据
var ImportantFixedStars = []FixedStarInfo{
	{Name: "Aldebaran", Chinese: "毕宿五", Longitude: 69.47 + 10.0/60, Magnitude: 0.85, Nature: "Mars", IsPositive: true, Keywords: "荣耀、成功、勇气"},
	{Name: "Regulus", Chinese: "轩辕十四", Longitude: 149.50 + 50.0/60, Magnitude: 1.35, Nature: "Mars-Jupiter", IsPositive: true, Keywords: "权力、领导、成功"},
	{Name: "Antares", Chinese: "心宿二", Longitude: 249.46 + 12.0/60, Magnitude: 0.96, Nature: "Mars-Jupiter", IsPositive: false, Keywords: "战争、危险、执着"},
	{Name: "Fomalhaut", Chinese: "北落师门", Longitude: 333.50 + 30.0/60, Magnitude: 1.16, Nature: "Venus-Mercury", IsPositive: true, Keywords: "理想、名声、魔法"},
	{Name: "Algol", Chinese: "大陵五", Longitude: 56.10 + 10.0/60, Magnitude: 2.12, Nature: "Saturn-Jupiter", IsPositive: false, Keywords: "暴力、不幸、魔鬼"},
	{Name: "Spica", Chinese: "角宿一", Longitude: 203.50 + 20.0/60, Magnitude: 0.97, Nature: "Venus-Mars", IsPositive: true, Keywords: "才华、成功、财富"},
	{Name: "Vega", Chinese: "织女一", Longitude: 285.15 + 20.0/60, Magnitude: 0.03, Nature: "Venus-Mercury", IsPositive: true, Keywords: "艺术、魅力、变化"},
	{Name: "Sirius", Chinese: "天狼星", Longitude: 104.07, Magnitude: -1.46, Nature: "Jupiter-Mars", IsPositive: true, Keywords: "荣耀、野心、危险"},
}

// CalculateFixedStarFactors 计算恒星因子
func CalculateFixedStarFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	yearsFrom2000 := float64(queryTime.Year() - 2000)
	precessionCorrection := yearsFrom2000 * 50.3 / 3600

	// 检查点位
	type checkPoint struct {
		Name      string
		Longitude float64
	}
	checkPoints := []checkPoint{
		{"太阳", GetPlanetFromChart(chart, models.Sun).Longitude},
		{"月亮", GetPlanetFromChart(chart, models.Moon).Longitude},
		{"上升", chart.Ascendant},
		{"天顶", chart.Midheaven},
	}

	for _, star := range ImportantFixedStars {
		currentStarLong := star.Longitude + precessionCorrection

		for _, point := range checkPoints {
			diff := math.Abs(point.Longitude - currentStarLong)
			if diff > 180 {
				diff = 360 - diff
			}

			if diff < 1.5 {
				strength := 1.0 - (diff / 1.5)
				value := strength * 2.0
				if !star.IsPositive {
					value = -value
				}

				factor := models.InfluenceFactor{
					ID:          fmt.Sprintf("star_%s_%s_%s", star.Name, point.Name, queryTime.Format("20060102")),
					Type:        models.FactorFixedStar,
					Name:        fmt.Sprintf("%s合%s", point.Name, star.Chinese),
					Description: fmt.Sprintf("%s与恒星%s（%s）合相，距离%.2f°。%s", point.Name, star.Chinese, star.Name, diff, star.Keywords),
					BaseValue:   value,
					Weight:      1.0,
					Adjustment:  value,
					IsPositive:  star.IsPositive,
					TimeLevel:   models.TimeLevelDaily,
					DimensionImpact: models.DimensionImpact{
						Career: 0.25, Relationship: 0.2, Health: 0.15, Finance: 0.2, Spiritual: 0.2,
					},
				}
				factors = append(factors, factor)
			}
		}
	}

	return factors
}

// ==================== 阿拉伯点计算 ====================

// CalculateArabicPartFactors 计算阿拉伯点因子
func CalculateArabicPartFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	asc := chart.Ascendant
	sunPos := GetPlanetFromChart(chart, models.Sun)
	moonPos := GetPlanetFromChart(chart, models.Moon)
	if sunPos == nil || moonPos == nil {
		return factors
	}
	sun := sunPos.Longitude
	moon := moonPos.Longitude

	var fortunePart, spiritPart float64
	if sun > 180 {
		fortunePart = math.Mod(asc+sun-moon+360, 360)
		spiritPart = math.Mod(asc+moon-sun+360, 360)
	} else {
		fortunePart = math.Mod(asc+moon-sun+360, 360)
		spiritPart = math.Mod(asc+sun-moon+360, 360)
	}

	planets := []models.PlanetID{models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars, models.Jupiter, models.Saturn}
	for _, planetID := range planets {
		planet := GetPlanetFromChart(chart, planetID)
		if planet == nil {
			continue
		}
		planetChinese := getPlanetChineseName(planetID)

		diffFortune := math.Abs(planet.Longitude - fortunePart)
		if diffFortune > 180 {
			diffFortune = 360 - diffFortune
		}
		if diffFortune < 5 {
			strength := 1.0 - (diffFortune / 5.0)
			factor := models.InfluenceFactor{
				ID:           fmt.Sprintf("fortune_%s_%s", planetID, queryTime.Format("20060102")),
				Type:         models.FactorArabicPart,
				Name:         fmt.Sprintf("福点合%s", planetChinese),
				Description:  fmt.Sprintf("福点（%.1f°%s）与%s合相，物质层面的幸运与机遇", fortunePart, getSignFromLongitudeInternal(fortunePart), planetChinese),
				BaseValue:    strength * 2.0,
				Weight:       0.8,
				Adjustment:   strength * 1.6,
				IsPositive:   true,
				TimeLevel:    models.TimeLevelDaily,
				SourcePlanet: planetID,
				DimensionImpact: models.DimensionImpact{
					Career: 0.2, Relationship: 0.15, Health: 0.2, Finance: 0.3, Spiritual: 0.15,
				},
			}
			factors = append(factors, factor)
		}

		diffSpirit := math.Abs(planet.Longitude - spiritPart)
		if diffSpirit > 180 {
			diffSpirit = 360 - diffSpirit
		}
		if diffSpirit < 5 {
			strength := 1.0 - (diffSpirit / 5.0)
			factor := models.InfluenceFactor{
				ID:           fmt.Sprintf("spirit_%s_%s", planetID, queryTime.Format("20060102")),
				Type:         models.FactorArabicPart,
				Name:         fmt.Sprintf("精神点合%s", planetChinese),
				Description:  fmt.Sprintf("精神点（%.1f°%s）与%s合相，精神层面的成长与领悟", spiritPart, getSignFromLongitudeInternal(spiritPart), planetChinese),
				BaseValue:    strength * 1.8,
				Weight:       0.8,
				Adjustment:   strength * 1.44,
				IsPositive:   true,
				TimeLevel:    models.TimeLevelDaily,
				SourcePlanet: planetID,
				DimensionImpact: models.DimensionImpact{
					Career: 0.15, Relationship: 0.2, Health: 0.15, Finance: 0.15, Spiritual: 0.35,
				},
			}
			factors = append(factors, factor)
		}
	}

	return factors
}

// ==================== 界限与十度面 ====================

// TermRuler 界限主星
type TermRuler struct {
	Planet string
	Start  float64
	End    float64
}

// EgyptianTerms 埃及界限表
var EgyptianTerms = map[string][]TermRuler{
	"Aries":       {{Planet: "Jupiter", Start: 0, End: 6}, {Planet: "Venus", Start: 6, End: 12}, {Planet: "Mercury", Start: 12, End: 20}, {Planet: "Mars", Start: 20, End: 25}, {Planet: "Saturn", Start: 25, End: 30}},
	"Taurus":      {{Planet: "Venus", Start: 0, End: 8}, {Planet: "Mercury", Start: 8, End: 14}, {Planet: "Jupiter", Start: 14, End: 22}, {Planet: "Saturn", Start: 22, End: 27}, {Planet: "Mars", Start: 27, End: 30}},
	"Gemini":      {{Planet: "Mercury", Start: 0, End: 6}, {Planet: "Jupiter", Start: 6, End: 12}, {Planet: "Venus", Start: 12, End: 17}, {Planet: "Mars", Start: 17, End: 24}, {Planet: "Saturn", Start: 24, End: 30}},
	"Cancer":      {{Planet: "Mars", Start: 0, End: 7}, {Planet: "Venus", Start: 7, End: 13}, {Planet: "Mercury", Start: 13, End: 19}, {Planet: "Jupiter", Start: 19, End: 26}, {Planet: "Saturn", Start: 26, End: 30}},
	"Leo":         {{Planet: "Jupiter", Start: 0, End: 6}, {Planet: "Venus", Start: 6, End: 11}, {Planet: "Saturn", Start: 11, End: 18}, {Planet: "Mercury", Start: 18, End: 24}, {Planet: "Mars", Start: 24, End: 30}},
	"Virgo":       {{Planet: "Mercury", Start: 0, End: 7}, {Planet: "Venus", Start: 7, End: 17}, {Planet: "Jupiter", Start: 17, End: 21}, {Planet: "Mars", Start: 21, End: 28}, {Planet: "Saturn", Start: 28, End: 30}},
	"Libra":       {{Planet: "Saturn", Start: 0, End: 6}, {Planet: "Mercury", Start: 6, End: 14}, {Planet: "Jupiter", Start: 14, End: 21}, {Planet: "Venus", Start: 21, End: 28}, {Planet: "Mars", Start: 28, End: 30}},
	"Scorpio":     {{Planet: "Mars", Start: 0, End: 7}, {Planet: "Venus", Start: 7, End: 11}, {Planet: "Mercury", Start: 11, End: 19}, {Planet: "Jupiter", Start: 19, End: 24}, {Planet: "Saturn", Start: 24, End: 30}},
	"Sagittarius": {{Planet: "Jupiter", Start: 0, End: 12}, {Planet: "Venus", Start: 12, End: 17}, {Planet: "Mercury", Start: 17, End: 21}, {Planet: "Saturn", Start: 21, End: 26}, {Planet: "Mars", Start: 26, End: 30}},
	"Capricorn":   {{Planet: "Mercury", Start: 0, End: 7}, {Planet: "Jupiter", Start: 7, End: 14}, {Planet: "Venus", Start: 14, End: 22}, {Planet: "Saturn", Start: 22, End: 26}, {Planet: "Mars", Start: 26, End: 30}},
	"Aquarius":    {{Planet: "Mercury", Start: 0, End: 7}, {Planet: "Venus", Start: 7, End: 13}, {Planet: "Jupiter", Start: 13, End: 20}, {Planet: "Mars", Start: 20, End: 25}, {Planet: "Saturn", Start: 25, End: 30}},
	"Pisces":      {{Planet: "Venus", Start: 0, End: 12}, {Planet: "Jupiter", Start: 12, End: 16}, {Planet: "Mercury", Start: 16, End: 19}, {Planet: "Mars", Start: 19, End: 28}, {Planet: "Saturn", Start: 28, End: 30}},
}

// DecanRulers 十度面主星（迦勒底次序）
var DecanRulers = []string{"Mars", "Sun", "Venus", "Mercury", "Moon", "Saturn", "Jupiter"}

// CalculateTermFactors 计算界限因子
func CalculateTermFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	mainPlanets := []models.PlanetID{models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars}

	for _, planetID := range mainPlanets {
		planet := GetPlanetFromChart(chart, planetID)
		if planet == nil {
			continue
		}
		sign := getSignFromLongitudeInternal(planet.Longitude)
		degreeInSign := math.Mod(planet.Longitude, 30)

		terms, exists := EgyptianTerms[sign]
		if !exists {
			continue
		}

		planetChinese := getPlanetChineseName(planetID)
		planetNameStr := getPlanetEnglishName(planetID)

		for _, term := range terms {
			if degreeInSign >= term.Start && degreeInSign < term.End {
				if term.Planet == planetNameStr {
					factor := models.InfluenceFactor{
						ID:              fmt.Sprintf("term_%s_%s", planetID, queryTime.Format("20060102")),
						Type:            models.FactorTerm,
						Name:            fmt.Sprintf("%s在本界", planetChinese),
						Description:     fmt.Sprintf("%s在%s的%s界限内，获得界限尊贵", planetChinese, sign, term.Planet),
						BaseValue:       1.0,
						Weight:          0.5,
						Adjustment:      0.5,
						IsPositive:      true,
						TimeLevel:       models.TimeLevelWeekly, // 界限持续数天至数周
						SourcePlanet:    planetID,
						DimensionImpact: GetPlanetDimensionImpact(planetID),
					}
					factors = append(factors, factor)
				}
				break
			}
		}
	}

	return factors
}

// CalculateDecanFactors 计算十度面因子
func CalculateDecanFactors(chart *models.NatalChart, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	mainPlanets := []models.PlanetID{models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars}

	for _, planetID := range mainPlanets {
		planet := GetPlanetFromChart(chart, planetID)
		if planet == nil {
			continue
		}

		absoluteDegree := planet.Longitude
		decanIndex := int(absoluteDegree/10) % 36
		decanRuler := DecanRulers[decanIndex%7]

		planetNameStr := getPlanetEnglishName(planetID)
		if decanRuler == planetNameStr {
			planetChinese := getPlanetChineseName(planetID)
			factor := models.InfluenceFactor{
				ID:              fmt.Sprintf("decan_%s_%s", planetID, queryTime.Format("20060102")),
				Type:            models.FactorDecan,
				Name:            fmt.Sprintf("%s在本面", planetChinese),
				Description:     fmt.Sprintf("%s在自己主管的十度面内，获得面尊贵", planetChinese),
				BaseValue:       0.8,
				Weight:          0.4,
				Adjustment:      0.32,
				IsPositive:      true,
				TimeLevel:       models.TimeLevelWeekly, // 十度面持续数天至数周
				SourcePlanet:    planetID,
				DimensionImpact: GetPlanetDimensionImpact(planetID),
			}
			factors = append(factors, factor)
		}
	}

	return factors
}

// ==================== 推运技术 ====================

// CalculateSolarArcFactors 计算太阳弧推进因子
func CalculateSolarArcFactors(chart *models.NatalChart, birthTime, queryTime time.Time) []models.InfluenceFactor {
	var factors []models.InfluenceFactor

	age := queryTime.Sub(birthTime).Hours() / 24 / 365.25
	solarArc := age * 0.9856

	natalPoints := []struct {
		Name      string
		Longitude float64
	}{
		{"上升", chart.Ascendant},
		{"天顶", chart.Midheaven},
		{"太阳", GetPlanetFromChart(chart, models.Sun).Longitude},
		{"月亮", GetPlanetFromChart(chart, models.Moon).Longitude},
	}

	progressedPlanets := []models.PlanetID{models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars}

	for _, progPlanetID := range progressedPlanets {
		progPlanet := GetPlanetFromChart(chart, progPlanetID)
		if progPlanet == nil {
			continue
		}
		natalPos := progPlanet.Longitude
		progPos := math.Mod(natalPos+solarArc, 360)

		progPlanetChinese := getPlanetChineseName(progPlanetID)

		for _, natalPoint := range natalPoints {
			if progPlanetChinese == natalPoint.Name {
				continue
			}

			diff := math.Abs(progPos - natalPoint.Longitude)
			if diff > 180 {
				diff = 360 - diff
			}

			aspects := []struct {
				angle float64
				name  string
			}{
				{0, "合"},
				{90, "刑"},
				{180, "冲"},
				{120, "拱"},
				{60, "六合"},
			}

			for _, asp := range aspects {
				aspDiff := math.Abs(diff - asp.angle)
				if aspDiff < 1.0 {
					strength := 1.0 - aspDiff
					isPositive := asp.angle == 0 || asp.angle == 120 || asp.angle == 60

					lifecycle := &models.FactorLifecycle{
						StartTime: queryTime.AddDate(0, -6, 0),
						PeakTime:  queryTime,
						EndTime:   queryTime.AddDate(0, 6, 0),
						Duration:  8760,
					}

					factor := models.InfluenceFactor{
						ID:           fmt.Sprintf("solararc_%s_%s_%s_%s", progPlanetID, asp.name, natalPoint.Name, queryTime.Format("200601")),
						Type:         models.FactorSolarArc,
						Name:         fmt.Sprintf("太阳弧%s%s本命%s", progPlanetChinese, asp.name, natalPoint.Name),
						Description:  fmt.Sprintf("太阳弧推进%s（%.1f°）%s本命%s，重要人生转折点", progPlanetChinese, progPos, asp.name, natalPoint.Name),
						BaseValue:    strength * 3.0,
						Weight:       1.2,
						Adjustment:   strength * 3.6,
						IsPositive:   isPositive,
						TimeLevel:    models.TimeLevelYearly,
						Lifecycle:    lifecycle,
						SourcePlanet: progPlanetID,
						DimensionImpact: models.DimensionImpact{
							Career: 0.25, Relationship: 0.25, Health: 0.15, Finance: 0.2, Spiritual: 0.15,
						},
					}
					factors = append(factors, factor)
				}
			}
		}
	}

	return factors
}

// ==================== 辅助函数 ====================

// getPlanetChineseName 获取行星中文名
func getPlanetChineseName(planet models.PlanetID) string {
	names := map[models.PlanetID]string{
		models.Sun:       "太阳",
		models.Moon:      "月亮",
		models.Mercury:   "水星",
		models.Venus:     "金星",
		models.Mars:      "火星",
		models.Jupiter:   "木星",
		models.Saturn:    "土星",
		models.Uranus:    "天王星",
		models.Neptune:   "海王星",
		models.Pluto:     "冥王星",
		models.NorthNode: "北交点",
		models.Chiron:    "凯龙",
	}
	if name, ok := names[planet]; ok {
		return name
	}
	return string(planet)
}

// getPlanetEnglishName 获取行星英文名（用于界限/十度面匹配）
func getPlanetEnglishName(planet models.PlanetID) string {
	names := map[models.PlanetID]string{
		models.Sun:     "Sun",
		models.Moon:    "Moon",
		models.Mercury: "Mercury",
		models.Venus:   "Venus",
		models.Mars:    "Mars",
		models.Jupiter: "Jupiter",
		models.Saturn:  "Saturn",
	}
	if name, ok := names[planet]; ok {
		return name
	}
	return string(planet)
}

// getSignFromLongitudeInternal 从黄经获取星座（内部使用）
func getSignFromLongitudeInternal(longitude float64) string {
	signs := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	index := int(longitude/30) % 12
	return signs[index]
}

// GetSignFromLongitude 从黄经获取星座（导出用于API）
func GetSignFromLongitude(longitude float64) string {
	return getSignFromLongitudeInternal(longitude)
}

// ==================== 统一计算入口 ====================

// CalculateAllAdvancedFactors 计算所有高级因子
func CalculateAllAdvancedFactors(chart *models.NatalChart, birthTime, queryTime time.Time) []models.InfluenceFactor {
	var allFactors []models.InfluenceFactor

	allFactors = append(allFactors, CalculateEclipseFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateLunarNodeFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateCombustionFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateStationFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateReceptionFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateFixedStarFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateArabicPartFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateTermFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateDecanFactors(chart, queryTime)...)
	allFactors = append(allFactors, CalculateSolarArcFactors(chart, birthTime, queryTime)...)

	return allFactors
}
