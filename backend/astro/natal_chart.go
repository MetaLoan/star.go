package astro

import (
	"star/models"
)

// CalculateNatalChart 计算本命盘
func CalculateNatalChart(birthData models.BirthData) *models.NatalChart {
	// 计算儒略日
	jd := DateToJulianDay(birthData.ToTime())

	// 计算行星位置
	planets := GetAllPlanetPositions(jd)

	// 计算宫位
	houses, ascendant, midheaven := CalculateHouses(jd, birthData.Latitude, birthData.Longitude)

	// 为行星分配宫位
	planets = AssignHousesToPlanets(planets, houses)

	// 计算相位
	aspects := CalculateAspects(planets)

	// 检测图形相位
	patterns := DetectPatterns(aspects, planets)

	// 计算元素和模式平衡
	elementBalance := CalculateElementBalance(planets)
	modalityBalance := CalculateModalityBalance(planets)

	// 找出主导行星
	dominantPlanets := FindDominantPlanets(planets, aspects)

	// 确定命主星
	chartRuler := GetChartRuler(ascendant)

	return &models.NatalChart{
		BirthData:       birthData,
		Planets:         planets,
		Houses:          houses,
		Ascendant:       ascendant,
		Midheaven:       midheaven,
		Aspects:         aspects,
		Patterns:        patterns,
		ElementBalance:  elementBalance,
		ModalityBalance: modalityBalance,
		DominantPlanets: dominantPlanets,
		ChartRuler:      chartRuler,
	}
}

// GetPlanetFromChart 从星盘中获取指定行星
func GetPlanetFromChart(chart *models.NatalChart, planetID models.PlanetID) *models.PlanetPosition {
	for _, p := range chart.Planets {
		if p.ID == planetID {
			return &p
		}
	}
	return nil
}

