package astro

import (
	"math"
	"star/models"
	"time"
)

// CalculateHouses 计算宫位 (Placidus 分宫制)
func CalculateHouses(jd float64, latitude, longitude float64) ([]models.HouseCusp, float64, float64) {
	// 计算本地恒星时
	lst := calculateLocalSiderealTime(jd, longitude)

	// 计算上升点和中天
	ascendant := calculateAscendant(lst, latitude)
	midheaven := calculateMidheaven(lst)

	// 计算12宫位
	houses := make([]models.HouseCusp, 12)
	houses[0] = createHouseCusp(1, ascendant)
	houses[9] = createHouseCusp(10, midheaven)

	// 计算其他宫位 (Placidus)
	houses[1] = createHouseCusp(2, calculatePlacidusHouse(2, ascendant, midheaven, latitude))
	houses[2] = createHouseCusp(3, calculatePlacidusHouse(3, ascendant, midheaven, latitude))
	houses[3] = createHouseCusp(4, NormalizeAngle(midheaven+180)) // IC
	houses[4] = createHouseCusp(5, calculatePlacidusHouse(5, ascendant, midheaven, latitude))
	houses[5] = createHouseCusp(6, calculatePlacidusHouse(6, ascendant, midheaven, latitude))
	houses[6] = createHouseCusp(7, NormalizeAngle(ascendant+180)) // DC
	houses[7] = createHouseCusp(8, calculatePlacidusHouse(8, ascendant, midheaven, latitude))
	houses[8] = createHouseCusp(9, calculatePlacidusHouse(9, ascendant, midheaven, latitude))
	houses[10] = createHouseCusp(11, calculatePlacidusHouse(11, ascendant, midheaven, latitude))
	houses[11] = createHouseCusp(12, calculatePlacidusHouse(12, ascendant, midheaven, latitude))

	return houses, ascendant, midheaven
}

// calculateLocalSiderealTime 计算本地恒星时
func calculateLocalSiderealTime(jd float64, longitude float64) float64 {
	T := (jd - J2000) / 36525.0

	// 格林威治平恒星时
	GMST := 280.46061837 + 360.98564736629*(jd-J2000) +
		0.000387933*T*T - T*T*T/38710000

	// 本地恒星时
	LST := NormalizeAngle(GMST + longitude)

	return LST
}

// calculateAscendant 计算上升点
func calculateAscendant(lst float64, latitude float64) float64 {
	RAMC := lst * DEG_TO_RAD
	phi := latitude * DEG_TO_RAD
	eps := OBLIQUITY * DEG_TO_RAD

	// ASC = atan2(cos(RAMC), -sin(RAMC)*cos(ε) - tan(φ)*sin(ε))
	y := math.Cos(RAMC)
	x := -math.Sin(RAMC)*math.Cos(eps) - math.Tan(phi)*math.Sin(eps)

	asc := math.Atan2(y, x) * RAD_TO_DEG

	return NormalizeAngle(asc)
}

// calculateMidheaven 计算中天
func calculateMidheaven(lst float64) float64 {
	RAMC := lst * DEG_TO_RAD
	eps := OBLIQUITY * DEG_TO_RAD

	// MC = atan2(sin(RAMC), cos(RAMC)*cos(ε))
	mc := math.Atan2(math.Sin(RAMC), math.Cos(RAMC)*math.Cos(eps)) * RAD_TO_DEG

	return NormalizeAngle(mc)
}

// calculatePlacidusHouse 计算 Placidus 宫位
func calculatePlacidusHouse(house int, asc, mc float64, _ float64) float64 {
	// 简化实现：使用插值法
	// 完整实现需要迭代求解

	switch house {
	case 2:
		return NormalizeAngle(asc + 30)
	case 3:
		return NormalizeAngle(asc + 60)
	case 5:
		return NormalizeAngle(mc + 210)
	case 6:
		return NormalizeAngle(mc + 240)
	case 8:
		return NormalizeAngle(asc + 210)
	case 9:
		return NormalizeAngle(asc + 240)
	case 11:
		return NormalizeAngle(mc + 30)
	case 12:
		return NormalizeAngle(mc + 60)
	default:
		return 0
	}
}

// createHouseCusp 创建宫位对象
func createHouseCusp(house int, cusp float64) models.HouseCusp {
	zodiac := GetZodiacByLongitude(cusp)
	return models.HouseCusp{
		House:    house,
		Cusp:     cusp,
		Sign:     zodiac.ID,
		SignName: zodiac.Name,
	}
}

// GetPlanetHouse 获取行星所在宫位
func GetPlanetHouse(planetLongitude float64, houses []models.HouseCusp) int {
	for i := 0; i < 12; i++ {
		nextHouse := (i + 1) % 12
		cusp1 := houses[i].Cusp
		cusp2 := houses[nextHouse].Cusp

		// 处理跨越0度的情况
		if cusp2 < cusp1 {
			cusp2 += 360
		}

		testLon := planetLongitude
		if testLon < cusp1 {
			testLon += 360
		}

		if testLon >= cusp1 && testLon < cusp2 {
			return houses[i].House
		}
	}
	return 1 // 默认返回1宫
}

// AssignHousesToPlanets 为行星分配宫位
func AssignHousesToPlanets(planets []models.PlanetPosition, houses []models.HouseCusp) []models.PlanetPosition {
	result := make([]models.PlanetPosition, len(planets))
	copy(result, planets)

	for i := range result {
		result[i].House = GetPlanetHouse(result[i].Longitude, houses)
	}

	return result
}

// CalculateElementBalance 计算四元素平衡
func CalculateElementBalance(planets []models.PlanetPosition) map[string]float64 {
	elements := map[string]float64{
		"fire":  0,
		"earth": 0,
		"air":   0,
		"water": 0,
	}

	totalWeight := 0.0

	for _, p := range planets {
		weight := PlanetWeights[p.ID]
		zodiac := GetZodiacInfo(p.Sign)
		if zodiac != nil {
			elements[zodiac.Element] += weight
			totalWeight += weight
		}
	}

	// 转换为百分比
	if totalWeight > 0 {
		for k := range elements {
			elements[k] = (elements[k] / totalWeight) * 100
		}
	}

	return elements
}

// CalculateModalityBalance 计算三模式平衡
func CalculateModalityBalance(planets []models.PlanetPosition) map[string]float64 {
	modalities := map[string]float64{
		"cardinal": 0,
		"fixed":    0,
		"mutable":  0,
	}

	totalWeight := 0.0

	for _, p := range planets {
		weight := PlanetWeights[p.ID]
		zodiac := GetZodiacInfo(p.Sign)
		if zodiac != nil {
			modalities[zodiac.Modality] += weight
			totalWeight += weight
		}
	}

	// 转换为百分比
	if totalWeight > 0 {
		for k := range modalities {
			modalities[k] = (modalities[k] / totalWeight) * 100
		}
	}

	return modalities
}

// FindDominantPlanets 找出主导行星
func FindDominantPlanets(planets []models.PlanetPosition, aspects []models.AspectData) []models.PlanetID {
	// 计算每颗行星的综合权重
	scores := make(map[models.PlanetID]float64)

	for _, p := range planets {
		// 基础权重
		scores[p.ID] = PlanetWeights[p.ID]

		// 尊贵度加成
		scores[p.ID] += p.DignityScore
	}

	// 相位加成
	for _, a := range aspects {
		scores[a.Planet1] += a.Weight * 0.5
		scores[a.Planet2] += a.Weight * 0.5
	}

	// 找出得分最高的行星
	var dominant []models.PlanetID
	threshold := 10.0 // 阈值

	for id, score := range scores {
		if score >= threshold {
			dominant = append(dominant, id)
		}
	}

	// 如果没有超过阈值，返回得分最高的两颗
	if len(dominant) == 0 {
		var maxID1, maxID2 models.PlanetID
		var maxScore1, maxScore2 float64

		for id, score := range scores {
			if score > maxScore1 {
				maxScore2 = maxScore1
				maxID2 = maxID1
				maxScore1 = score
				maxID1 = id
			} else if score > maxScore2 {
				maxScore2 = score
				maxID2 = id
			}
		}

		dominant = []models.PlanetID{maxID1, maxID2}
	}

	return dominant
}

// GetChartRuler 获取命主星
func GetChartRuler(ascendant float64) models.PlanetID {
	zodiac := GetZodiacByLongitude(ascendant)
	if zodiac != nil {
		return zodiac.Ruler
	}
	return models.Sun
}

// CalculateMoonPhase 计算月相
func CalculateMoonPhase(sunLon, moonLon float64) models.MoonPhase {
	angle := NormalizeAngle(moonLon - sunLon)
	phaseInfo := GetLunarPhase(angle)

	// 计算照明度 (0-1)
	illumination := (1 - math.Cos(angle*DEG_TO_RAD)) / 2

	return models.MoonPhase{
		Phase:        phaseInfo.Phase,
		Name:         phaseInfo.Name,
		Angle:        angle,
		Illumination: illumination,
	}
}

// CalculatePlanetaryDay 计算行星日
func CalculatePlanetaryDay(date time.Time) models.PlanetID {
	// 星期日=太阳, 星期一=月亮, 星期二=火星...
	dayOrder := []models.PlanetID{
		models.Sun, models.Moon, models.Mars, models.Mercury,
		models.Jupiter, models.Venus, models.Saturn,
	}
	return dayOrder[date.Weekday()]
}

// CalculatePlanetaryHour 计算行星时
func CalculatePlanetaryHour(date time.Time, hour int) models.PlanetID {
	// 行星时序列
	hourOrder := []models.PlanetID{
		models.Saturn, models.Jupiter, models.Mars, models.Sun,
		models.Venus, models.Mercury, models.Moon,
	}

	// 计算从星期日午夜开始的小时数
	dayOfWeek := int(date.Weekday())
	totalHours := dayOfWeek*24 + hour

	// 每颗行星按照迦勒底顺序循环
	index := totalHours % 7
	return hourOrder[index]
}

