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
	{models.Saturn, models.Saturn, models.Conjunction, "土星回归", 28, "high"},
	{models.Saturn, models.Saturn, models.Opposition, "土星对冲", 14, "high"},
	{models.Saturn, models.Saturn, models.Square, "土星刑克", 7, "medium"},
	{models.Jupiter, models.Jupiter, models.Conjunction, "木星回归", 11, "medium"},
	{models.Uranus, models.Uranus, models.Opposition, "天王星对冲", 38, "high"},
	{models.Uranus, models.Uranus, models.Square, "天王星刑克", 20, "medium"},
	{models.Neptune, models.Neptune, models.Square, "海王星刑克", 40, "high"},
	{models.Pluto, models.Pluto, models.Square, "冥王星刑克", 60, "high"},
	{models.NorthNode, models.NorthNode, models.Conjunction, "北交点回归", 18, "medium"},
	{models.Chiron, models.Chiron, models.Conjunction, "凯龙回归", 50, "high"},
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
		theme = "融合与新开始"
		keywords = []string{"合一", "启动", "强化"}
		advice = "注意能量的融合与释放"
	case models.Trine:
		theme = "顺畅与机遇"
		keywords = []string{"和谐", "机会", "流动"}
		advice = "把握机遇，顺势而为"
	case models.Sextile:
		theme = "支持与合作"
		keywords = []string{"合作", "支持", "小机会"}
		advice = "积极寻求合作机会"
	case models.Square:
		theme = "挑战与成长"
		keywords = []string{"紧张", "压力", "突破"}
		advice = "直面挑战，寻求突破"
	case models.Opposition:
		theme = "平衡与觉察"
		keywords = []string{"对立", "平衡", "觉醒"}
		advice = "寻求平衡，整合对立面"
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
	case "土星回归":
		return age == 29 || age == 30 || age == 58 || age == 59 || age == 60
	case "土星对冲":
		return age == 14 || age == 15 || age == 44 || age == 45
	case "土星刑克":
		return age == 7 || age == 21 || age == 22 || age == 36 || age == 37 || age == 51 || age == 52
	case "木星回归":
		return age%12 == 0 && age > 0
	case "天王星对冲":
		return age >= 40 && age <= 42
	case "天王星刑克":
		return age == 21 || age == 63
	case "海王星刑克":
		return age >= 40 && age <= 42
	case "冥王星刑克":
		return age >= 60 && age <= 90 // 因人而异
	case "北交点回归":
		return age == 18 || age == 19 || age == 37 || age == 38 || age == 56 || age == 57
	case "凯龙回归":
		return age >= 50 && age <= 51
	default:
		return false
	}
}

// generateMajorTransitDescription 生成重大行运描述
func generateMajorTransitDescription(config MajorTransitConfig, _ int) string {
	descriptions := map[string]string{
		"土星回归":   "人生结构的重大升级期，需要承担更多责任，重新定义人生方向",
		"土星对冲":   "面对责任与挑战的关键时期，需要平衡个人需求与外界期望",
		"土星刑克":   "成长的压力测试，需要调整结构性问题",
		"木星回归":   "扩张与机遇的周期起点，适合开始新的成长",
		"天王星对冲":  "中年觉醒期，可能带来突破性的变化和新的自我认知",
		"天王星刑克":  "突破传统束缚的时期，追求个人独特性",
		"海王星刑克":  "灵性转化期，需要区分理想与现实",
		"冥王星刑克":  "深层转化期，涉及权力、控制和重生主题",
		"北交点回归":  "命运的节点，人生方向的重要调整期",
		"凯龙回归":   "伤痛与治愈的整合期，成为他人的疗愈者",
	}

	if desc, ok := descriptions[config.Name]; ok {
		return desc
	}
	return "重要的人生转折期"
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

