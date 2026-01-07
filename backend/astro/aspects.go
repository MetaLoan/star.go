package astro

import (
	"fmt"
	"math"
	"star/models"
)

// CalculateAspects 计算行星之间的相位
func CalculateAspects(planets []models.PlanetPosition) []models.AspectData {
	var aspects []models.AspectData

	for i := 0; i < len(planets); i++ {
		for j := i + 1; j < len(planets); j++ {
			p1, p2 := planets[i], planets[j]

			// 计算角距
			diff := math.Abs(p1.Longitude - p2.Longitude)
			if diff > 180 {
				diff = 360 - diff
			}

			// 检查每个相位类型
			for _, def := range AspectDefinitions {
				orb := math.Abs(diff - def.Angle)
				if orb <= def.Orb {
					// 计算强度 (0-1)
					strength := 1.0 - orb/def.Orb

					// 计算权重
					p1Weight := PlanetWeights[p1.ID]
					p2Weight := PlanetWeights[p2.ID]
					weight := strength * def.Weight * (p1Weight + p2Weight) / 20.0

					// 判断是否入相
					applying := isApplying(p1, p2, def.Angle)

					// 生成解读
					interpretation := generateAspectInterpretation(p1, p2, def)

					aspects = append(aspects, models.AspectData{
						Planet1:        p1.ID,
						Planet2:        p2.ID,
						AspectType:     def.Type,
						ExactAngle:     def.Angle,
						ActualAngle:    diff,
						Orb:            orb,
						Applying:       applying,
						Strength:       strength,
						Weight:         weight,
						Interpretation: interpretation,
					})
				}
			}
		}
	}

	return aspects
}

// CalculateTransitToNatalAspects 计算行运与本命的相位
func CalculateTransitToNatalAspects(transitPositions, natalPositions []models.PlanetPosition) []models.AspectData {
	var aspects []models.AspectData

	// 行运容许度收紧到 80%
	orbFactor := 0.8

	for _, transit := range transitPositions {
		for _, natal := range natalPositions {
			// 计算角距
			diff := math.Abs(transit.Longitude - natal.Longitude)
			if diff > 180 {
				diff = 360 - diff
			}

			for _, def := range AspectDefinitions {
				adjustedOrb := def.Orb * orbFactor
				orb := math.Abs(diff - def.Angle)

				if orb <= adjustedOrb {
					strength := 1.0 - orb/adjustedOrb

					p1Weight := PlanetWeights[transit.ID]
					p2Weight := PlanetWeights[natal.ID]
					weight := strength * def.Weight * (p1Weight + p2Weight) / 20.0

				interpretation := fmt.Sprintf("Transit %s forms %s with natal %s",
					transit.Name, def.Name, natal.Name)

					aspects = append(aspects, models.AspectData{
						Planet1:        transit.ID,
						Planet2:        natal.ID,
						AspectType:     def.Type,
						ExactAngle:     def.Angle,
						ActualAngle:    diff,
						Orb:            orb,
						Applying:       true, // 简化处理
						Strength:       strength,
						Weight:         weight,
						Interpretation: interpretation,
					})
				}
			}
		}
	}

	return aspects
}

// isApplying 判断相位是否入相（接近中）
func isApplying(p1, p2 models.PlanetPosition, _ float64) bool {
	// 简化判断：根据逆行状态
	// 如果两颗行星都不逆行，且p1速度快，则入相
	// 这里使用简化逻辑
	if p1.Retrograde || p2.Retrograde {
		return false
	}
	return true
}

// generateAspectInterpretation 生成相位解读
func generateAspectInterpretation(p1, p2 models.PlanetPosition, def AspectDefinition) string {
	p1Info := GetPlanetInfo(p1.ID)
	p2Info := GetPlanetInfo(p2.ID)

	if p1Info == nil || p2Info == nil {
		return ""
	}

	return fmt.Sprintf("%s forms %s with %s", p1Info.Name, def.Name, p2Info.Name)
}

// DetectPatterns 检测图形相位
func DetectPatterns(aspects []models.AspectData, planets []models.PlanetPosition) []string {
	var patterns []string

	// 检测大三角 (三个三分相)
	if hasGrandTrine(aspects) {
		patterns = append(patterns, "Grand Trine")
	}

	// 检测T三角 (两个四分相和一个对分相)
	if hasTSquare(aspects) {
		patterns = append(patterns, "T-Square")
	}

	// 检测大十字 (四个四分相和两个对分相)
	if hasGrandCross(aspects) {
		patterns = append(patterns, "Grand Cross")
	}

	// 检测风筝
	if hasKite(aspects) {
		patterns = append(patterns, "Kite")
	}

	return patterns
}

// hasGrandTrine 检测大三角
func hasGrandTrine(aspects []models.AspectData) bool {
	trines := filterAspectsByType(aspects, models.Trine)
	if len(trines) < 3 {
		return false
	}

	// 简化检测：检查是否有三颗行星互成三分相
	planetTrines := make(map[models.PlanetID]int)
	for _, a := range trines {
		planetTrines[a.Planet1]++
		planetTrines[a.Planet2]++
	}

	count := 0
	for _, v := range planetTrines {
		if v >= 2 {
			count++
		}
	}
	return count >= 3
}

// hasTSquare 检测T三角
func hasTSquare(aspects []models.AspectData) bool {
	squares := filterAspectsByType(aspects, models.Square)
	oppositions := filterAspectsByType(aspects, models.Opposition)

	if len(squares) < 2 || len(oppositions) < 1 {
		return false
	}

	// 检查是否有一颗行星同时与对分相的两颗行星成四分相
	for _, opp := range oppositions {
		for _, sq1 := range squares {
			for _, sq2 := range squares {
				if sq1.Planet1 == sq2.Planet1 && sq1.Planet1 != opp.Planet1 && sq1.Planet1 != opp.Planet2 {
					if (sq1.Planet2 == opp.Planet1 || sq1.Planet2 == opp.Planet2) &&
						(sq2.Planet2 == opp.Planet1 || sq2.Planet2 == opp.Planet2) {
						return true
					}
				}
			}
		}
	}
	return false
}

// hasGrandCross 检测大十字
func hasGrandCross(aspects []models.AspectData) bool {
	squares := filterAspectsByType(aspects, models.Square)
	oppositions := filterAspectsByType(aspects, models.Opposition)

	return len(squares) >= 4 && len(oppositions) >= 2
}

// hasKite 检测风筝
func hasKite(aspects []models.AspectData) bool {
	// 风筝 = 大三角 + 一颗行星与三角顶点对分相并与另两点六分相
	if !hasGrandTrine(aspects) {
		return false
	}

	oppositions := filterAspectsByType(aspects, models.Opposition)
	sextiles := filterAspectsByType(aspects, models.Sextile)

	return len(oppositions) >= 1 && len(sextiles) >= 2
}

// filterAspectsByType 按类型过滤相位
func filterAspectsByType(aspects []models.AspectData, aspectType models.AspectType) []models.AspectData {
	var result []models.AspectData
	for _, a := range aspects {
		if a.AspectType == aspectType {
			result = append(result, a)
		}
	}
	return result
}

// CalculateTransitScore 计算行运分数
func CalculateTransitScore(aspects []models.AspectData) models.TransitScore {
	var harmonious, tense float64

	for _, a := range aspects {
		switch a.AspectType {
		case models.Trine, models.Sextile:
			harmonious += a.Weight
		case models.Square, models.Opposition:
			tense += a.Weight
		case models.Conjunction:
			// 合相根据行星性质判断
			if isBenefic(a.Planet1) || isBenefic(a.Planet2) {
				harmonious += a.Weight * 0.5
			} else {
				tense += a.Weight * 0.3
			}
		}
	}

	return models.TransitScore{
		Total:      harmonious - tense,
		Harmonious: harmonious,
		Tense:      tense,
	}
}

// isBenefic 判断是否为吉星
func isBenefic(planet models.PlanetID) bool {
	switch planet {
	case models.Venus, models.Jupiter:
		return true
	case models.Sun, models.Moon, models.Mercury:
		return true // 中性偏吉
	default:
		return false
	}
}

