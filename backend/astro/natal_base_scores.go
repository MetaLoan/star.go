package astro

import (
	"star/models"
)

// ==================== 本命盘基础分计算 ====================
// 基于占星学理论：本命盘决定个人在各维度的先天禀赋
// 计算因素：
// 1. 行星落入宫位的贡献
// 2. 宫主星状态
// 3. 相位格局加成

// CalculateNatalBaseScores 计算本命盘的五维度基础分
func CalculateNatalBaseScores(chart *models.NatalChart) models.NatalBaseScores {
	// 基准分为50
	baseScores := models.NatalBaseScores{
		Career:       50.0,
		Relationship: 50.0,
		Health:       50.0,
		Finance:      50.0,
		Spiritual:    50.0,
	}

	// 1. 计算行星落入宫位的贡献
	planetContributions := calculatePlanetInHouseContributions(chart)
	addContributions(&baseScores, planetContributions)

	// 2. 计算宫主星状态贡献
	rulerContributions := calculateHouseRulerContributions(chart)
	addContributions(&baseScores, rulerContributions)

	// 3. 计算相位格局加成
	aspectContributions := calculateAspectPatternContributions(chart)
	addContributions(&baseScores, aspectContributions)

	// 4. 确保基础分在合理范围内 (35-65)
	clampBaseScores(&baseScores)

	return baseScores
}

// calculatePlanetInHouseContributions 计算行星落入宫位的贡献
func calculatePlanetInHouseContributions(chart *models.NatalChart) models.NatalBaseScores {
	contributions := models.NatalBaseScores{}

	for _, planet := range chart.Planets {
		// 获取行星所在宫位
		house := getPlanetHouse(planet.Longitude, chart.Houses)
		if house == 0 {
			continue
		}

		// 获取行星的尊贵度
		dignity := GetDignity(planet.ID, planet.Sign)
		dignityMultiplier := getDignityMultiplier(dignity)

		// 获取行星的维度影响分配
		impact := GetPlanetDimensionImpact(planet.ID)

		// 基础贡献值：行星的自然影响力
		baseValue := getPlanetNaturalWeight(planet.ID)

		// 宫位是否为该维度的关联宫位（额外加成）
		houseDimension := GetDimensionForHouseV2(house)
		houseBonus := 1.0
		if houseDimension != "" {
			houseBonus = 1.3 // 落入相关宫位有30%加成
		}

		// 计算对各维度的贡献
		contribution := baseValue * dignityMultiplier * houseBonus

		contributions.Career += contribution * impact.Career
		contributions.Relationship += contribution * impact.Relationship
		contributions.Health += contribution * impact.Health
		contributions.Finance += contribution * impact.Finance
		contributions.Spiritual += contribution * impact.Spiritual
	}

	return contributions
}

// calculateHouseRulerContributions 计算宫主星状态贡献
func calculateHouseRulerContributions(chart *models.NatalChart) models.NatalBaseScores {
	contributions := models.NatalBaseScores{}

	// 关键宫位及其对应维度
	keyHouses := map[int]string{
		1:  "health",       // 1宫 → 健康
		2:  "finance",      // 2宫 → 财务
		6:  "health",       // 6宫 → 健康
		7:  "relationship", // 7宫 → 关系
		8:  "finance",      // 8宫 → 财务
		9:  "spiritual",    // 9宫 → 灵性
		10: "career",       // 10宫 → 事业
		11: "relationship", // 11宫 → 关系
		12: "spiritual",    // 12宫 → 灵性
	}

	for houseNum, dimension := range keyHouses {
		// 获取宫头星座
		if houseNum > len(chart.Houses) {
			continue
		}
		houseCusp := chart.Houses[houseNum-1]

		// 获取该星座的守护星
		ruler := GetZodiacRuler(houseCusp.Sign)
		if ruler == "" {
			continue
		}

		// 找到守护星在本命盘的位置
		var rulerPlanet *models.PlanetPosition
		for i := range chart.Planets {
			if chart.Planets[i].ID == ruler {
				rulerPlanet = &chart.Planets[i]
				break
			}
		}

		if rulerPlanet == nil {
			continue
		}

		// 计算宫主星状态
		dignity := GetDignity(rulerPlanet.ID, rulerPlanet.Sign)
		dignityScore := GetDignityScore(dignity)

		// 宫主星落入的宫位是否有力
		rulerHouse := getPlanetHouse(rulerPlanet.Longitude, chart.Houses)
		houseStrength := getHouseStrength(rulerHouse)

		// 逆行减分
		retroPenalty := 0.0
		if rulerPlanet.Retrograde {
			retroPenalty = -1.0
		}

		contribution := dignityScore + houseStrength + retroPenalty

		// 应用到对应维度
		switch dimension {
		case "career":
			contributions.Career += contribution
		case "relationship":
			contributions.Relationship += contribution
		case "health":
			contributions.Health += contribution
		case "finance":
			contributions.Finance += contribution
		case "spiritual":
			contributions.Spiritual += contribution
		}
	}

	return contributions
}

// calculateAspectPatternContributions 计算相位格局加成
func calculateAspectPatternContributions(chart *models.NatalChart) models.NatalBaseScores {
	contributions := models.NatalBaseScores{}

	// 统计和谐相位和紧张相位数量
	harmoniousCount := 0
	tenseCount := 0

	for _, asp := range chart.Aspects {
		switch asp.AspectType {
		case models.Trine, models.Sextile:
			harmoniousCount++
		case models.Square, models.Opposition:
			tenseCount++
		}
	}

	// 和谐相位多：整体加分
	if harmoniousCount > tenseCount {
		bonus := float64(harmoniousCount-tenseCount) * 0.5
		contributions.Career += bonus
		contributions.Relationship += bonus
		contributions.Health += bonus
		contributions.Finance += bonus
		contributions.Spiritual += bonus
	}

	// 检查特殊格局
	// TODO: 实现大三角、大十字等格局检测

	return contributions
}

// ==================== 辅助函数 ====================

// getPlanetHouse 获取行星所在宫位
func getPlanetHouse(longitude float64, houses []models.HouseCusp) int {
	if len(houses) < 12 {
		return 0
	}

	for i := 0; i < 12; i++ {
		currentCusp := houses[i].Cusp
		nextCusp := houses[(i+1)%12].Cusp

		// 处理跨越0度的情况
		if nextCusp < currentCusp {
			if longitude >= currentCusp || longitude < nextCusp {
				return i + 1
			}
		} else {
			if longitude >= currentCusp && longitude < nextCusp {
				return i + 1
			}
		}
	}

	return 1 // 默认第一宫
}

// getDignityMultiplier 获取尊贵度乘数
func getDignityMultiplier(dignity models.Dignity) float64 {
	switch dignity {
	case models.DignityDomicile:
		return 1.5 // 入庙
	case models.DignityExaltation:
		return 1.3 // 旺相
	case models.DignityDetriment:
		return 0.7 // 落陷
	case models.DignityFall:
		return 0.5 // 失势
	default:
		return 1.0 // 普通
	}
}

// getPlanetNaturalWeight 获取行星的自然影响权重
func getPlanetNaturalWeight(planet models.PlanetID) float64 {
	weights := map[models.PlanetID]float64{
		models.Sun:       3.0, // 发光体，最重要
		models.Moon:      3.0, // 发光体，最重要
		models.Mercury:   2.0,
		models.Venus:     2.0,
		models.Mars:      2.0,
		models.Jupiter:   2.5, // 大吉星
		models.Saturn:    2.5, // 重要限制星
		models.Uranus:    1.5, // 外行星
		models.Neptune:   1.5,
		models.Pluto:     1.5,
		models.NorthNode: 1.0,
		models.Chiron:    1.0,
	}

	if w, ok := weights[planet]; ok {
		return w
	}
	return 1.0
}

// getHouseStrength 获取宫位力量
func getHouseStrength(house int) float64 {
	// 始宫最强，续宫次之，果宫最弱
	switch house {
	case 1, 4, 7, 10:
		return 2.0 // 始宫（Angular）
	case 2, 5, 8, 11:
		return 1.0 // 续宫（Succedent）
	case 3, 6, 9, 12:
		return 0.0 // 果宫（Cadent）
	default:
		return 0.0
	}
}

// addContributions 将贡献加到基础分
func addContributions(base *models.NatalBaseScores, contrib models.NatalBaseScores) {
	base.Career += contrib.Career
	base.Relationship += contrib.Relationship
	base.Health += contrib.Health
	base.Finance += contrib.Finance
	base.Spiritual += contrib.Spiritual
}

// clampBaseScores 限制基础分范围
func clampBaseScores(scores *models.NatalBaseScores) {
	clamp := func(v float64) float64 {
		if v < 35 {
			return 35
		}
		if v > 65 {
			return 65
		}
		return v
	}

	scores.Career = clamp(scores.Career)
	scores.Relationship = clamp(scores.Relationship)
	scores.Health = clamp(scores.Health)
	scores.Finance = clamp(scores.Finance)
	scores.Spiritual = clamp(scores.Spiritual)
}

// GetZodiacRuler 获取星座的守护星
func GetZodiacRuler(sign models.ZodiacID) models.PlanetID {
	rulers := map[models.ZodiacID]models.PlanetID{
		models.Aries:       models.Mars,
		models.Taurus:      models.Venus,
		models.Gemini:      models.Mercury,
		models.Cancer:      models.Moon,
		models.Leo:         models.Sun,
		models.Virgo:       models.Mercury,
		models.Libra:       models.Venus,
		models.Scorpio:     models.Pluto, // 现代守护星
		models.Sagittarius: models.Jupiter,
		models.Capricorn:   models.Saturn,
		models.Aquarius:    models.Uranus, // 现代守护星
		models.Pisces:      models.Neptune, // 现代守护星
	}

	if ruler, ok := rulers[sign]; ok {
		return ruler
	}
	return ""
}

