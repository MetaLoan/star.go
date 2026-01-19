package astro

import (
	"math"
	"sort"
	"star/models"
	"time"
)

// DailyEvent 每日星象事件
type DailyEvent struct {
	Time        time.Time              `json:"time"`
	Type        string                 `json:"type"` // aspect, sign_change, lunar_phase, planetary_hour_change
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Theme       string                 `json:"theme"`
	Advice      string                 `json:"advice"`
	Planet1     models.PlanetID        `json:"planet1,omitempty"`
	Planet2     models.PlanetID        `json:"planet2,omitempty"`
	Aspect      string                 `json:"aspect,omitempty"`
	Sign        string                 `json:"sign,omitempty"`
	Degree      float64                `json:"degree,omitempty"`
	IsPositive  bool                   `json:"isPositive"`
	Intensity   string                 `json:"intensity"` // high, medium, low
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CalculateDailyEvents 计算指定日期的所有星象事件（精确到分钟）
func CalculateDailyEvents(chart *models.NatalChart, date time.Time, includeMinorAspects bool) []DailyEvent {
	events := []DailyEvent{}

	// 设置时间范围：从当天00:00到次日00:00
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 1. 查找行星换座事件
	signChangeEvents := findSignChangeEvents(startOfDay, endOfDay)
	events = append(events, signChangeEvents...)

	// 2. 查找重要相位事件
	aspectEvents := findAspectEvents(chart, startOfDay, endOfDay, includeMinorAspects)
	events = append(events, aspectEvents...)

	// 3. 查找月相事件
	lunarPhaseEvents := findLunarPhaseEvents(startOfDay, endOfDay)
	events = append(events, lunarPhaseEvents...)

	// 4. 查找行星时变化（每2小时一次）
	planetaryHourEvents := findPlanetaryHourEvents(startOfDay, endOfDay)
	events = append(events, planetaryHourEvents...)

	// 按时间排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	return events
}

// findSignChangeEvents 查找行星换座事件
func findSignChangeEvents(startTime, endTime time.Time) []DailyEvent {
	events := []DailyEvent{}

	// 要检查的行星（快行星更可能换座）
	planetsToCheck := []models.PlanetID{
		models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars,
	}

	for _, planetID := range planetsToCheck {
		// 检查这个行星是否在这一天换座
		exactTime := findSignChangeTime(planetID, startTime, endTime)
		if exactTime != nil {
			// 获取新星座
			pos := CalculatePlanetPositionSwe(planetID, TimeToJulianDay(*exactTime))
			zodiacInfo := GetZodiacByLongitude(pos.Longitude)
			newSign := string(zodiacInfo.ID)

			event := DailyEvent{
				Time:        *exactTime,
				Type:        "sign_change",
				Title:       getPlanetName(planetID) + " enters " + getSignName(newSign),
				Description: getPlanetName(planetID) + " enters " + getSignName(newSign),
				Theme:       getSignChangeTheme(planetID, newSign),
				Advice:      getSignChangeAdvice(planetID, newSign),
				Planet1:     planetID,
				Sign:        newSign,
				Degree:      0.0,
				IsPositive:  isSignChangePositive(planetID, newSign),
				Intensity:   getSignChangeIntensity(planetID),
			}
			events = append(events, event)
		}
	}

	return events
}

// findSignChangeTime 使用二分法精确查找行星换座时间
func findSignChangeTime(planetID models.PlanetID, startTime, endTime time.Time) *time.Time {
	startJD := TimeToJulianDay(startTime)
	endJD := TimeToJulianDay(endTime)

	// 获取起始和结束位置
	startPos := CalculatePlanetPositionSwe(planetID, startJD)
	endPos := CalculatePlanetPositionSwe(planetID, endJD)

	startZodiac := GetZodiacByLongitude(startPos.Longitude)
	endZodiac := GetZodiacByLongitude(endPos.Longitude)

	// 如果星座没变化，返回nil
	if startZodiac.ID == endZodiac.ID {
		return nil
	}

	// 使用二分法精确查找换座时间（精确到1分钟）
	maxIterations := 20 // 2^20分钟 > 2年，足够了
	currentStart := startJD
	currentEnd := endJD
	startSign := startZodiac.ID

	for i := 0; i < maxIterations; i++ {
		midJD := (currentStart + currentEnd) / 2.0
		midPos := CalculatePlanetPositionSwe(planetID, midJD)
		midZodiac := GetZodiacByLongitude(midPos.Longitude)

		if midZodiac.ID == startSign {
			currentStart = midJD
		} else {
			currentEnd = midJD
		}

		// 如果精度足够（< 1分钟）
		if (currentEnd-currentStart)*1440 < 1.0 { // 1440分钟/天
			break
		}
	}

	exactTime := JulianDayToTime((currentStart + currentEnd) / 2.0)
	return &exactTime
}

// findAspectEvents 查找相位事件
func findAspectEvents(chart *models.NatalChart, startTime, endTime time.Time, includeMinor bool) []DailyEvent {
	events := []DailyEvent{}

	// 行运行星
	transitPlanets := []models.PlanetID{
		models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars,
		models.Jupiter, models.Saturn,
	}

	// 主要相位
	majorAspects := []struct {
		Angle float64
		Name  string
		Orb   float64
	}{
		{0, "合相", 8.0},
		{60, "六合", 6.0},
		{90, "刑相", 8.0},
		{120, "三合", 8.0},
		{180, "对分", 8.0},
	}

	// 次要相位
	minorAspects := []struct {
		Angle float64
		Name  string
		Orb   float64
	}{
		{30, "半六合", 2.0},
		{45, "半刑", 2.0},
		{135, "补八分", 2.0},
		{150, "梅花", 2.0},
	}

	aspects := majorAspects
	if includeMinor {
		aspects = append(aspects, minorAspects...)
	}

	// 对每个行运行星和本命行星的组合
	for _, transitPlanet := range transitPlanets {
		for _, natalPlanet := range chart.Planets {
			for _, aspect := range aspects {
				// 使用精确搜索找到相位发生的准确时间
				// 搜索时间窗口：前后各12小时
				searchStart := TimeToJulianDay(startTime)
				searchEnd := TimeToJulianDay(endTime)
				
				exactJD, found := FindExactAspectTime(
					transitPlanet,
					natalPlanet.ID,
					aspect.Angle,
					searchStart,
					searchEnd,
				)

				if found {
					exactTime := JulianDayToTime(exactJD)
					event := DailyEvent{
						Time:        exactTime,
						Type:        "aspect",
						Title:       getPlanetName(transitPlanet) + " " + getAspectNameEN(aspect.Name) + " " + getPlanetName(natalPlanet.ID),
						Description: "Transiting " + getPlanetName(transitPlanet) + " forms " + getAspectNameEN(aspect.Name) + " with natal " + getPlanetName(natalPlanet.ID),
						Theme:       getAspectTheme(transitPlanet, natalPlanet.ID, aspect.Name),
						Advice:      getAspectAdvice(transitPlanet, natalPlanet.ID, aspect.Name),
						Planet1:     transitPlanet,
						Planet2:     natalPlanet.ID,
						Aspect:      getAspectNameEN(aspect.Name),
						Degree:      aspect.Angle,
						IsPositive:  isAspectPositive(aspect.Name),
						Intensity:   getAspectIntensity(transitPlanet, natalPlanet.ID, aspect.Name),
					}
					events = append(events, event)
				}
			}
		}
	}

	return events
}

// findLunarPhaseEvents 查找月相事件（新月、满月等）
func findLunarPhaseEvents(startTime, endTime time.Time) []DailyEvent {
	events := []DailyEvent{}

	// 检查关键月相角度
	phaseAngles := []struct {
		Angle float64
		Name  string
	}{
		{0, "新月"},
		{90, "上弦月"},
		{180, "满月"},
		{270, "下弦月"},
	}

	for _, phase := range phaseAngles {
		exactTime := findLunarPhaseTime(phase.Angle, startTime, endTime)
		if exactTime != nil {
			event := DailyEvent{
				Time:        *exactTime,
				Type:        "lunar_phase",
				Title:       getLunarPhaseNameEN(phase.Name),
				Description: getLunarPhaseNameEN(phase.Name) + " occurs",
				Theme:       getLunarPhaseTheme(phase.Name),
				Advice:      getLunarPhaseAdvice(phase.Name),
				IsPositive:  true,
				Intensity:   "high",
			}
			events = append(events, event)
		}
	}

	return events
}

// findLunarPhaseTime 精确查找月相时间
func findLunarPhaseTime(targetAngle float64, startTime, endTime time.Time) *time.Time {
	startJD := TimeToJulianDay(startTime)
	endJD := TimeToJulianDay(endTime)

	// 计算起始和结束的月相角度
	startSun := CalculatePlanetPositionSwe(models.Sun, startJD)
	startMoon := CalculatePlanetPositionSwe(models.Moon, startJD)
	startAngle := normalizeDegrees(startMoon.Longitude - startSun.Longitude)

	endSun := CalculatePlanetPositionSwe(models.Sun, endJD)
	endMoon := CalculatePlanetPositionSwe(models.Moon, endJD)
	endAngle := normalizeDegrees(endMoon.Longitude - endSun.Longitude)

	// 检查是否跨越目标角度
	if !crossesAngle(startAngle, endAngle, targetAngle) {
		return nil
	}

	// 二分法精确查找
	maxIterations := 20
	currentStart := startJD
	currentEnd := endJD

	for i := 0; i < maxIterations; i++ {
		midJD := (currentStart + currentEnd) / 2.0
		midSun := CalculatePlanetPositionSwe(models.Sun, midJD)
		midMoon := CalculatePlanetPositionSwe(models.Moon, midJD)
		midAngle := normalizeDegrees(midMoon.Longitude - midSun.Longitude)

		if crossesAngle(startAngle, midAngle, targetAngle) {
			currentEnd = midJD
			endAngle = midAngle
		} else {
			currentStart = midJD
			startAngle = midAngle
		}

		if (currentEnd-currentStart)*1440 < 1.0 {
			break
		}
	}

	exactTime := JulianDayToTime((currentStart + currentEnd) / 2.0)
	return &exactTime
}

// findPlanetaryHourEvents 查找行星时变化
func findPlanetaryHourEvents(startTime, endTime time.Time) []DailyEvent {
	events := []DailyEvent{}

	// 行星时序列
	planetaryHours := []models.PlanetID{
		models.Saturn, models.Jupiter, models.Mars, models.Sun,
		models.Venus, models.Mercury, models.Moon,
	}

	// 从startTime开始，每2小时一个行星时
	current := startTime
	hourIndex := 0
	for current.Before(endTime) {
		rulingPlanet := planetaryHours[hourIndex%7]
		event := DailyEvent{
			Time:        current,
			Type:        "planetary_hour_change",
			Title:       getPlanetName(rulingPlanet) + " Hour",
			Description: getPlanetName(rulingPlanet) + " rules this hour",
			Theme:       getPlanetaryHourTheme(rulingPlanet),
			Advice:      getPlanetaryHourAdvice(rulingPlanet),
			Planet1:     rulingPlanet,
			IsPositive:  true,
			Intensity:   "low",
		}
		events = append(events, event)
		current = current.Add(2 * time.Hour)
		hourIndex++
	}

	return events
}

// 辅助函数

func crossesAngle(start, end, target float64) bool {
	// 处理角度跨越0度的情况
	if start > end {
		return target >= start || target <= end
	}
	return target >= start && target <= end
}

func normalizeDegrees(deg float64) float64 {
	result := math.Mod(deg, 360.0)
	if result < 0 {
		result += 360.0
	}
	return result
}

func isAspectPositive(aspectName string) bool {
	positiveAspects := map[string]bool{
		"三合": true,
		"六合": true,
	}
	return positiveAspects[aspectName]
}

func getAspectIntensity(planet1, planet2 models.PlanetID, aspectName string) string {
	// 内行星相位通常影响较弱
	if (planet1 == models.Moon || planet1 == models.Mercury) &&
		(planet2 == models.Moon || planet2 == models.Mercury) {
		return "low"
	}

	// 外行星相位影响较强
	outerPlanets := map[models.PlanetID]bool{
		models.Jupiter: true,
		models.Saturn:  true,
		models.Uranus:  true,
		models.Neptune: true,
		models.Pluto:   true,
	}

	if outerPlanets[planet1] || outerPlanets[planet2] {
		return "high"
	}

	return "medium"
}

func getSignChangeIntensity(planet models.PlanetID) string {
	intensityMap := map[models.PlanetID]string{
		models.Sun:     "high",
		models.Moon:    "medium",
		models.Mercury: "medium",
		models.Venus:   "medium",
		models.Mars:    "medium",
	}
	if intensity, ok := intensityMap[planet]; ok {
		return intensity
	}
	return "low"
}

func isSignChangePositive(planet models.PlanetID, sign string) bool {
	// 简化判断：入庙旺的星座为正面
	dignityMap := map[models.PlanetID][]string{
		models.Sun:     {"aries", "leo"},
		models.Moon:    {"cancer", "taurus"},
		models.Mercury: {"gemini", "virgo"},
		models.Venus:   {"taurus", "libra"},
		models.Mars:    {"aries", "scorpio"},
	}

	if signs, ok := dignityMap[planet]; ok {
		for _, s := range signs {
			if s == sign {
				return true
			}
		}
	}
	return false
}

// 主题和建议函数（简化版，可以后续扩展）

func getSignChangeTheme(planet models.PlanetID, sign string) string {
	themes := map[string]map[string]string{
		string(models.Sun): {
			"aquarius": "Humanitarianism, social awareness, independent innovation",
			"pisces":   "Dreams, intuition, compassion",
			"aries":    "Initiative, action, courage",
		},
		string(models.Moon): {
			"cancer":    "Emotional richness, family warmth",
			"scorpio":   "Deep emotions, transformative power",
			"capricorn": "Emotional stability, responsibility",
		},
	}

	if planetThemes, ok := themes[string(planet)]; ok {
		if theme, ok := planetThemes[sign]; ok {
			return theme
		}
	}
	return "Energy transition period"
}

func getSignChangeAdvice(planet models.PlanetID, sign string) string {
	return "Adapt to new energy, adjust your approach"
}

func getAspectTheme(planet1, planet2 models.PlanetID, aspect string) string {
	if aspect == "三合" || aspect == "六合" {
		return "Harmonious energy, opportunities emerge"
	}
	if aspect == "刑相" || aspect == "对分" {
		return "Challenges and growth opportunities"
	}
	return "Energy interaction"
}

func getAspectAdvice(planet1, planet2 models.PlanetID, aspect string) string {
	adviceMap := map[string]string{
		"三合": "Seize opportunities, go with the flow",
		"六合": "Take initiative, create opportunities",
		"刑相": "Face challenges, break through limitations",
		"对分": "Seek balance, integrate opposites",
		"合相": "Focus energy, concentrate on goals",
	}
	if advice, ok := adviceMap[aspect]; ok {
		return advice
	}
	return "Be aware of energy changes"
}

func getLunarPhaseTheme(phaseName string) string {
	themes := map[string]string{
		"新月":  "New beginnings, planting intentions",
		"上弦月": "Increased action, overcoming obstacles",
		"满月":  "Manifestation, emotional peak",
		"下弦月": "Reflection, release and letting go",
	}
	if theme, ok := themes[phaseName]; ok {
		return theme
	}
	return "Lunar phase transition"
}

func getLunarPhaseAdvice(phaseName string) string {
	advice := map[string]string{
		"新月":  "Set goals, start new plans",
		"上弦月": "Take action, advance projects",
		"满月":  "Celebrate achievements, express gratitude",
		"下弦月": "Organize, reflect, clear old things",
	}
	if adv, ok := advice[phaseName]; ok {
		return adv
	}
	return "Follow the lunar rhythm"
}

func getPlanetaryHourTheme(planet models.PlanetID) string {
	themes := map[models.PlanetID]string{
		models.Sun:     "Vitality, leadership, creativity",
		models.Moon:    "Emotions, intuition, nurturing",
		models.Mercury: "Communication, learning, commerce",
		models.Venus:   "Love, art, enjoyment",
		models.Mars:    "Action, competition, courage",
		models.Jupiter: "Expansion, learning, travel",
		models.Saturn:  "Focus, discipline, building",
	}
	if theme, ok := themes[planet]; ok {
		return theme
	}
	return "Planetary energy"
}

func getPlanetaryHourAdvice(planet models.PlanetID) string {
	advice := map[models.PlanetID]string{
		models.Sun:     "Handle important matters, show leadership",
		models.Moon:    "Attend to emotional needs, rest and relax",
		models.Mercury: "Communicate, learn new things",
		models.Venus:   "Socialize, artistic creation",
		models.Mars:    "Physical exercise, decisive action",
		models.Jupiter: "Learn and expand, stay optimistic",
		models.Saturn:  "Focus on work, handle serious matters",
	}
	if adv, ok := advice[planet]; ok {
		return adv
	}
	return "Follow planetary energy"
}

func getPlanetName(planet models.PlanetID) string {
	names := map[models.PlanetID]string{
		models.Sun:     "Sun",
		models.Moon:    "Moon",
		models.Mercury: "Mercury",
		models.Venus:   "Venus",
		models.Mars:    "Mars",
		models.Jupiter: "Jupiter",
		models.Saturn:  "Saturn",
		models.Uranus:  "Uranus",
		models.Neptune: "Neptune",
		models.Pluto:   "Pluto",
	}
	if name, ok := names[planet]; ok {
		return name
	}
	return string(planet)
}

func getSignName(sign string) string {
	names := map[string]string{
		"aries":       "Aries",
		"taurus":      "Taurus",
		"gemini":      "Gemini",
		"cancer":      "Cancer",
		"leo":         "Leo",
		"virgo":       "Virgo",
		"libra":       "Libra",
		"scorpio":     "Scorpio",
		"sagittarius": "Sagittarius",
		"capricorn":   "Capricorn",
		"aquarius":    "Aquarius",
		"pisces":      "Pisces",
	}
	if name, ok := names[sign]; ok {
		return name
	}
	return sign
}

// getAspectNameEN 获取相位英文名称
func getAspectNameEN(aspectNameCN string) string {
	aspectMap := map[string]string{
		"合相":   "Conjunction",
		"六合":   "Sextile",
		"刑相":   "Square",
		"三合":   "Trine",
		"对分":   "Opposition",
		"半六合":  "Semi-Sextile",
		"半刑":   "Semi-Square",
		"补八分":  "Sesquiquadrate",
		"梅花":   "Quincunx",
	}
	if name, ok := aspectMap[aspectNameCN]; ok {
		return name
	}
	return aspectNameCN
}

// getLunarPhaseNameEN 获取月相英文名称
func getLunarPhaseNameEN(phaseNameCN string) string {
	phaseMap := map[string]string{
		"新月":  "New Moon",
		"上弦月": "First Quarter",
		"满月":  "Full Moon",
		"下弦月": "Last Quarter",
	}
	if name, ok := phaseMap[phaseNameCN]; ok {
		return name
	}
	return phaseNameCN
}
