package astro

import (
	"star/models"
	"time"
)

// ==================== 行星时计算（增强版） ====================

// PlanetaryHourInfo 行星时信息
type PlanetaryHourInfo struct {
	PlanetaryHour int             `json:"planetaryHour"` // 1-24
	Ruler         models.PlanetID `json:"ruler"`
	PlanetName    string          `json:"planetName"`
	PlanetSymbol  string          `json:"planetSymbol"`
	DayRuler      models.PlanetID `json:"dayRuler"`
	Influence     float64         `json:"influence"`
	BestFor       []string        `json:"bestFor"`
}

// 迦勒底行星顺序（行星时循环顺序）
var chaldeanOrder = []models.PlanetID{
	models.Saturn, models.Jupiter, models.Mars, models.Sun,
	models.Venus, models.Mercury, models.Moon,
}

// 日主星映射（周几对应哪颗行星）
var dayRulers = map[time.Weekday]models.PlanetID{
	time.Sunday:    models.Sun,
	time.Monday:    models.Moon,
	time.Tuesday:   models.Mars,
	time.Wednesday: models.Mercury,
	time.Thursday:  models.Jupiter,
	time.Friday:    models.Venus,
	time.Saturday:  models.Saturn,
}

// 行星时适合的活动
var planetaryHourActivities = map[models.PlanetID][]string{
	models.Sun:     {"leadership", "creativity", "self-expression", "important meetings"},
	models.Moon:    {"emotions", "family", "intuition", "rest"},
	models.Mercury: {"communication", "writing", "business", "learning"},
	models.Venus:   {"art", "socializing", "dating", "beauty"},
	models.Mars:    {"exercise", "competition", "initiating projects", "taking action"},
	models.Jupiter: {"learning", "travel", "networking", "decision-making"},
	models.Saturn:  {"planning", "meditation", "solitude", "organizing"},
}

// CalculatePlanetaryHourEnhanced 计算行星时（增强版）
func CalculatePlanetaryHourEnhanced(t time.Time, lat, lon float64) PlanetaryHourInfo {
	// 获取日主星
	dayRuler := dayRulers[t.Weekday()]

	// 计算日出日落时间（简化：使用固定的 6:00 和 18:00）
	// 完整实现需要根据纬度和日期计算实际日出日落时间
	sunriseHour := 6.0
	sunsetHour := 18.0

	// 日间和夜间小时数
	dayLength := sunsetHour - sunriseHour   // 12 小时
	nightLength := 24 - dayLength           // 12 小时

	// 当前时间（小时）
	currentHour := float64(t.Hour()) + float64(t.Minute())/60.0

	// 判断是日间还是夜间，计算行星时编号
	var planetaryHourNum int
	var hourLength float64

	if currentHour >= sunriseHour && currentHour < sunsetHour {
		// 日间
		hourLength = dayLength / 12.0
		hoursFromSunrise := currentHour - sunriseHour
		planetaryHourNum = int(hoursFromSunrise/hourLength) + 1
	} else {
		// 夜间
		hourLength = nightLength / 12.0
		var hoursFromSunset float64
		if currentHour >= sunsetHour {
			hoursFromSunset = currentHour - sunsetHour
		} else {
			hoursFromSunset = (24 - sunsetHour) + currentHour
		}
		planetaryHourNum = int(hoursFromSunset/hourLength) + 13
	}

	// 限制在 1-24 范围
	if planetaryHourNum > 24 {
		planetaryHourNum = 24
	}
	if planetaryHourNum < 1 {
		planetaryHourNum = 1
	}

	// 找到日主星在迦勒底顺序中的位置
	dayRulerIdx := 0
	for i, p := range chaldeanOrder {
		if p == dayRuler {
			dayRulerIdx = i
			break
		}
	}

	// 计算当前行星时的主管行星
	// 行星时从日主星开始，按迦勒底顺序循环
	rulerIdx := (dayRulerIdx + planetaryHourNum - 1) % 7
	ruler := chaldeanOrder[rulerIdx]

	// 获取行星信息
	planetInfo := GetPlanetInfo(ruler)

	// 计算影响值（基于行星吉凶性质）
	influence := getPlanetaryHourInfluence(ruler, dayRuler)

	return PlanetaryHourInfo{
		PlanetaryHour: planetaryHourNum,
		Ruler:         ruler,
		PlanetName:    planetInfo.Name,
		PlanetSymbol:  planetInfo.Symbol,
		DayRuler:      dayRuler,
		Influence:     influence,
		BestFor:       planetaryHourActivities[ruler],
	}
}

// getPlanetaryHourInfluence 获取行星时影响值
func getPlanetaryHourInfluence(ruler, dayRuler models.PlanetID) float64 {
	// 基础影响值
	baseInfluence := map[models.PlanetID]float64{
		models.Sun:     4.0,
		models.Moon:    3.0,
		models.Mercury: 2.0,
		models.Venus:   4.0,
		models.Mars:    -1.0,
		models.Jupiter: 5.0,
		models.Saturn:  -2.0,
	}

	influence := baseInfluence[ruler]

	// 如果行星时主管与日主星相同，加成
	if ruler == dayRuler {
		influence += 2.0
	}

	return influence
}

// GetPlanetaryHoursForDate 获取指定日期所有行星时
func GetPlanetaryHoursForDate(date time.Time, lat, lon float64) []PlanetaryHourInfo {
	hours := make([]PlanetaryHourInfo, 24)
	for h := 0; h < 24; h++ {
		t := time.Date(date.Year(), date.Month(), date.Day(), h, 0, 0, 0, date.Location())
		hours[h] = CalculatePlanetaryHourEnhanced(t, lat, lon)
	}
	return hours
}

