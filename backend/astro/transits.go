package astro

import (
	"star/models"
	"time"
)

// MajorTransitConfig 重大行运配置
type MajorTransitConfig struct {
	TransitPlanet models.PlanetID
	NatalPlanet   models.PlanetID
	AspectType    models.AspectType
	Name          string
	MinAge        int
	Significance  string
}

// MajorTransitConfigs 重大行运配置列表
var MajorTransitConfigs = []MajorTransitConfig{
	{models.Saturn, models.Saturn, models.Conjunction, "Saturn Return", 28, "high"},
	{models.Saturn, models.Saturn, models.Opposition, "Saturn Opposition", 14, "high"},
	{models.Saturn, models.Saturn, models.Square, "Saturn Square", 7, "medium"},
	{models.Jupiter, models.Jupiter, models.Conjunction, "Jupiter Return", 11, "medium"},
	{models.Uranus, models.Uranus, models.Opposition, "Uranus Opposition", 38, "high"},
	{models.Uranus, models.Uranus, models.Square, "Uranus Square", 20, "medium"},
	{models.Neptune, models.Neptune, models.Square, "Neptune Square", 40, "high"},
	{models.Pluto, models.Pluto, models.Square, "Pluto Square", 60, "high"},
	{models.NorthNode, models.NorthNode, models.Conjunction, "North Node Return", 18, "medium"},
	{models.Chiron, models.Chiron, models.Conjunction, "Chiron Return", 50, "high"},
}

// CalculateTransits 计算行运
func CalculateTransits(chart *models.NatalChart, startDateStr, endDateStr string) *models.TransitResult {
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	if startDate.IsZero() {
		startDate = time.Now()
	}
	if endDate.IsZero() {
		endDate = startDate.AddDate(1, 0, 0) // 默认一年
	}

	var events []models.TransitEvent
	var dominantThemes []string

	// 遍历日期范围，检测行运事件
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 7) { // 每周检查
		transitPositions := GetTransitPositions(date)
		aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)

		for _, aspect := range aspects {
			if aspect.Strength > 0.8 { // 只记录强相位
				event := models.TransitEvent{
					Date:          date,
					TransitPlanet: aspect.Planet1,
					NatalPlanet:   aspect.Planet2,
					Aspect:        aspect,
					Phase:         "exact",
					Intensity:     aspect.Strength * 100,
					Duration: models.TransitDuration{
						Start: date.AddDate(0, 0, -7).Format("2006-01-02"),
						Peak:  date.Format("2006-01-02"),
						End:   date.AddDate(0, 0, 7).Format("2006-01-02"),
					},
					Interpretation: generateTransitInterpretation(aspect),
				}
				events = append(events, event)
			}
		}
	}

	// 计算整体分数
	var totalScore float64
	if len(events) > 0 {
		for _, e := range events {
			totalScore += e.Intensity
		}
		totalScore /= float64(len(events))
	} else {
		totalScore = 50
	}

	// 提取主导主题
	themeMap := make(map[string]int)
	for _, e := range events {
		theme := e.Interpretation.Theme
		if theme != "" {
			themeMap[theme]++
		}
	}
	for theme := range themeMap {
		dominantThemes = append(dominantThemes, theme)
	}

	return &models.TransitResult{
		StartDate:      startDateStr,
		EndDate:        endDateStr,
		Events:         events,
		OverallScore:   totalScore,
		DominantThemes: dominantThemes,
	}
}

// generateTransitInterpretation 生成行运解读
func generateTransitInterpretation(aspect models.AspectData) models.TransitInterpretation {
	transitInfo := GetPlanetInfo(aspect.Planet1)
	natalInfo := GetPlanetInfo(aspect.Planet2)

	if transitInfo == nil || natalInfo == nil {
		return models.TransitInterpretation{}
	}

	// 根据相位类型生成主题
	var theme string
	var keywords []string
	var advice string

	switch aspect.AspectType {
	case models.Conjunction:
		theme = "Fusion & New Beginnings"
		keywords = []string{"unity", "initiation", "intensification"}
		advice = "Pay attention to the fusion and release of energy"
	case models.Trine:
		theme = "Flow & Opportunity"
		keywords = []string{"harmony", "opportunity", "flow"}
		advice = "Seize opportunities and go with the flow"
	case models.Sextile:
		theme = "Support & Cooperation"
		keywords = []string{"cooperation", "support", "small opportunities"}
		advice = "Actively seek cooperation opportunities"
	case models.Square:
		theme = "Challenge & Growth"
		keywords = []string{"tension", "pressure", "breakthrough"}
		advice = "Face challenges head-on, seek breakthrough"
	case models.Opposition:
		theme = "Balance & Awareness"
		keywords = []string{"polarity", "balance", "awakening"}
		advice = "Seek balance, integrate opposites"
	}

	return models.TransitInterpretation{
		Theme:    theme,
		Keywords: keywords,
		Advice:   advice,
	}
}

// FindMajorTransits 查找重大行运
func FindMajorTransits(chart *models.NatalChart, startYear, endYear int) []models.MajorTransit {
	var transits []models.MajorTransit
	birthYear := chart.BirthData.ToTime().Year()

	for _, config := range MajorTransitConfigs {
		// 获取本命行星位置
		natalPlanet := GetPlanetFromChart(chart, config.NatalPlanet)
		if natalPlanet == nil {
			continue
		}

		// 计算该行星回归/相位的大致年份
		for year := startYear; year <= endYear; year++ {
			age := year - birthYear
			if age < config.MinAge {
				continue
			}

			// 检查是否是该事件的典型年龄
			if isTypicalAgeForTransit(age, config) {
				date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
				transits = append(transits, models.MajorTransit{
					Name:         config.Name,
					Date:         date,
					Description:  generateMajorTransitDescription(config, age),
					Significance: config.Significance,
				})
			}
		}
	}

	return transits
}

// isTypicalAgeForTransit 判断是否是行运的典型年龄
func isTypicalAgeForTransit(age int, config MajorTransitConfig) bool {
	switch config.Name {
	case "Saturn Return":
		return age == 29 || age == 30 || age == 58 || age == 59 || age == 60
	case "Saturn Opposition":
		return age == 14 || age == 15 || age == 44 || age == 45
	case "Saturn Square":
		return age == 7 || age == 21 || age == 22 || age == 36 || age == 37 || age == 51 || age == 52
	case "Jupiter Return":
		return age%12 == 0 && age > 0
	case "Uranus Opposition":
		return age >= 40 && age <= 42
	case "Uranus Square":
		return age == 21 || age == 63
	case "Neptune Square":
		return age >= 40 && age <= 42
	case "Pluto Square":
		return age >= 60 && age <= 90 // varies by person
	case "North Node Return":
		return age == 18 || age == 19 || age == 37 || age == 38 || age == 56 || age == 57
	case "Chiron Return":
		return age >= 50 && age <= 51
	default:
		return false
	}
}

// generateMajorTransitDescription 生成重大行运描述
func generateMajorTransitDescription(config MajorTransitConfig, _ int) string {
	descriptions := map[string]string{
		"Saturn Return":     "Major life structure upgrade, taking on more responsibility, redefining life direction",
		"Saturn Opposition": "Critical period facing responsibility and challenges, balancing personal needs with external expectations",
		"Saturn Square":     "Growth pressure test, adjusting structural issues",
		"Jupiter Return":    "Beginning of expansion and opportunity cycle, ideal for starting new growth",
		"Uranus Opposition": "Midlife awakening, may bring breakthrough changes and new self-awareness",
		"Uranus Square":     "Period of breaking traditional constraints, pursuing personal uniqueness",
		"Neptune Square":    "Spiritual transformation, distinguishing ideals from reality",
		"Pluto Square":      "Deep transformation involving power, control, and rebirth themes",
		"North Node Return": "Destiny node, important adjustment period for life direction",
		"Chiron Return":     "Integration of wounds and healing, becoming a healer for others",
	}

	if desc, ok := descriptions[config.Name]; ok {
		return desc
	}
	return "Important life transition period"
}

// GetCurrentTransits 获取当前活跃的行运
func GetCurrentTransits(chart *models.NatalChart, date time.Time) []models.TransitEvent {
	transitPositions := GetTransitPositions(date)
	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)

	var events []models.TransitEvent
	for _, aspect := range aspects {
		if aspect.Strength > 0.5 { // 中等强度以上
			events = append(events, models.TransitEvent{
				Date:           date,
				TransitPlanet:  aspect.Planet1,
				NatalPlanet:    aspect.Planet2,
				Aspect:         aspect,
				Phase:          "active",
				Intensity:      aspect.Strength * 100,
				Interpretation: generateTransitInterpretation(aspect),
			})
		}
	}

	return events
}

