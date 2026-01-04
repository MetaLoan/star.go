package astro

import (
	"math"
	"star/models"
	"time"
)

// CalculateDailyForecast 计算每日预测
func CalculateDailyForecast(chart *models.NatalChart, date time.Time, withFactors bool) *models.DailyForecast {
	jd := DateToJulianDay(date)

	// 获取当前行星位置
	transitPositions := GetTransitPositions(date)

	// 计算行运相位
	activeAspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)

	// 获取月亮信息
	var moonPos models.PlanetPosition
	var sunPos models.PlanetPosition
	for _, p := range transitPositions {
		if p.ID == models.Moon {
			moonPos = p
		}
		if p.ID == models.Sun {
			sunPos = p
		}
	}

	// 计算月相
	moonPhase := CalculateMoonPhase(sunPos.Longitude, moonPos.Longitude)

	// 月亮星座信息
	moonSign := models.MoonSignInfo{
		Sign:     moonPos.Sign,
		Name:     moonPos.SignName,
		Keywords: MoonSignKeywords[moonPos.Sign],
	}

	// 计算行运分数
	transitScore := CalculateTransitScore(activeAspects)

	// 计算影响因子
	var factors *models.FactorResult
	var topFactors []models.InfluenceFactor
	if withFactors {
		factors = CalculateInfluenceFactors(chart, date, transitPositions)
		topFactors = getTopFactors(factors, 5)
	}

	// 计算基础分数
	baseScore := 50 + transitScore.Total
	if factors != nil {
		baseScore += factors.TotalAdjustment
	}

	// 标准化到 0-100
	overallScore := normalizeScore(baseScore)

	// 计算维度分数
	dimensions := calculateDailyDimensions(chart, transitPositions, activeAspects, factors)

	// 计算小时预测
	hourlyBreakdown := calculateHourlyBreakdown(chart, date, transitPositions)

	// 生成主题
	overallTheme := generateDailyTheme(moonSign, activeAspects)

	// 获取星期
	dayOfWeek := getDayOfWeekChinese(date.Weekday())

	_ = jd

	return &models.DailyForecast{
		Date:            date,
		DayOfWeek:       dayOfWeek,
		OverallScore:    overallScore,
		OverallTheme:    overallTheme,
		Dimensions:      dimensions,
		MoonPhase:       moonPhase,
		MoonSign:        moonSign,
		HourlyBreakdown: hourlyBreakdown,
		ActiveAspects:   activeAspects,
		Factors:         factors,
		TopFactors:      topFactors,
	}
}

// CalculateWeeklyForecast 计算每周预测
func CalculateWeeklyForecast(chart *models.NatalChart, startDate time.Time, withFactors bool) *models.WeeklyForecast {
	// 确保从周一开始
	weekday := int(startDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startDate = startDate.AddDate(0, 0, -(weekday - 1))
	endDate := startDate.AddDate(0, 0, 6)

	// 计算每日摘要
	dailySummaries := make([]models.DailySummary, 7)
	var totalScore float64
	var totalDimensions models.DailyDimensions

	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		dailyForecast := CalculateDailyForecast(chart, date, false)

		dailySummaries[i] = models.DailySummary{
			Date:         date.Format("2006-01-02"),
			DayOfWeek:    dailyForecast.DayOfWeek,
			OverallScore: dailyForecast.OverallScore,
			MoonSign:     dailyForecast.MoonSign.Sign,
			KeyTheme:     dailyForecast.OverallTheme,
		}

		totalScore += dailyForecast.OverallScore
		totalDimensions.Career += dailyForecast.Dimensions.Career
		totalDimensions.Relationship += dailyForecast.Dimensions.Relationship
		totalDimensions.Health += dailyForecast.Dimensions.Health
		totalDimensions.Finance += dailyForecast.Dimensions.Finance
		totalDimensions.Spiritual += dailyForecast.Dimensions.Spiritual
	}

	// 计算平均值
	overallScore := totalScore / 7
	dimensions := models.DailyDimensions{
		Career:       totalDimensions.Career / 7,
		Relationship: totalDimensions.Relationship / 7,
		Health:       totalDimensions.Health / 7,
		Finance:      totalDimensions.Finance / 7,
		Spiritual:    totalDimensions.Spiritual / 7,
	}

	// 找出关键日期和最佳日期
	keyDates := findKeyDates(dailySummaries)
	bestDaysFor := findBestDaysFor(chart, startDate)

	// 获取周度行运
	weeklyTransits := getWeeklyTransits(chart, startDate, endDate)

	// 计算维度趋势
	dimensionTrends := calculateDimensionTrends(dailySummaries)

	// 生成主题
	overallTheme := generateWeeklyTheme(weeklyTransits)

	// 周度因子
	var weeklyFactors *models.FactorResult
	if withFactors {
		weeklyFactors = CalculateInfluenceFactors(chart, startDate, GetTransitPositions(startDate))
	}

	return &models.WeeklyForecast{
		StartDate:       startDate,
		EndDate:         endDate,
		OverallScore:    overallScore,
		OverallTheme:    overallTheme,
		Dimensions:      dimensions,
		DailySummaries:  dailySummaries,
		KeyDates:        keyDates,
		BestDaysFor:     bestDaysFor,
		WeeklyTransits:  weeklyTransits,
		WeeklyFactors:   weeklyFactors,
		DimensionTrends: dimensionTrends,
	}
}

// normalizeScore 将原始分数标准化到 0-100
func normalizeScore(raw float64) float64 {
	// 使用 tanh 压缩
	scale := 20.0
	normalized := 50 + 50*math.Tanh(raw/scale)
	return math.Max(0, math.Min(100, normalized))
}

// calculateDailyDimensions 计算每日维度分数
func calculateDailyDimensions(_ *models.NatalChart, _ []models.PlanetPosition, aspects []models.AspectData, factors *models.FactorResult) models.DailyDimensions {
	base := 50.0

	// 根据相位调整各维度
	career := base
	relationship := base
	health := base
	finance := base
	spiritual := base

	for _, asp := range aspects {
		weight := asp.Weight * asp.Strength

		// 事业相关：土星、木星、太阳、MC
		if asp.Planet1 == models.Saturn || asp.Planet2 == models.Saturn ||
			asp.Planet1 == models.Jupiter || asp.Planet2 == models.Jupiter {
			switch asp.AspectType {
			case models.Trine, models.Sextile:
				career += weight
			case models.Square, models.Opposition:
				career -= weight * 0.5
			}
		}

		// 关系相关：金星、月亮、7宫主
		if asp.Planet1 == models.Venus || asp.Planet2 == models.Venus ||
			asp.Planet1 == models.Moon || asp.Planet2 == models.Moon {
			switch asp.AspectType {
			case models.Trine, models.Sextile:
				relationship += weight
			case models.Square, models.Opposition:
				relationship -= weight * 0.3
			}
		}

		// 健康相关：火星、太阳、6宫主
		if asp.Planet1 == models.Mars || asp.Planet2 == models.Mars ||
			asp.Planet1 == models.Sun || asp.Planet2 == models.Sun {
			switch asp.AspectType {
			case models.Trine, models.Sextile:
				health += weight * 0.5
			case models.Square, models.Opposition:
				health -= weight * 0.3
			}
		}

		// 财务相关：木星、金星、2/8宫主
		if asp.Planet1 == models.Jupiter || asp.Planet2 == models.Jupiter ||
			asp.Planet1 == models.Venus || asp.Planet2 == models.Venus {
			switch asp.AspectType {
			case models.Trine, models.Sextile:
				finance += weight
			case models.Square, models.Opposition:
				finance -= weight * 0.4
			}
		}

		// 灵性相关：海王星、冥王星、12宫主
		if asp.Planet1 == models.Neptune || asp.Planet2 == models.Neptune ||
			asp.Planet1 == models.Pluto || asp.Planet2 == models.Pluto {
			switch asp.AspectType {
			case models.Trine, models.Sextile:
				spiritual += weight
			case models.Conjunction:
				spiritual += weight * 0.5
			}
		}
	}

	// 应用因子调整
	if factors != nil {
		for dim, adj := range factors.DimensionAdjustments {
			switch dim {
			case "career":
				career += adj
			case "relationship":
				relationship += adj
			case "health":
				health += adj
			case "finance":
				finance += adj
			case "spiritual":
				spiritual += adj
			}
		}
	}

	return models.DailyDimensions{
		Career:       normalizeScore(career - 50),
		Relationship: normalizeScore(relationship - 50),
		Health:       normalizeScore(health - 50),
		Finance:      normalizeScore(finance - 50),
		Spiritual:    normalizeScore(spiritual - 50),
	}
}

// calculateHourlyBreakdown 计算小时预测
func calculateHourlyBreakdown(chart *models.NatalChart, date time.Time, _ []models.PlanetPosition) []models.HourlyForecast {
	breakdown := make([]models.HourlyForecast, 24)

	for hour := 0; hour < 24; hour++ {
		planetaryHour := CalculatePlanetaryHour(date, hour)
		score := calculateHourlyScore(chart, date, hour, planetaryHour)
		bestFor := getBestActivitiesForHour(planetaryHour)

		breakdown[hour] = models.HourlyForecast{
			Hour:          hour,
			Score:         score,
			PlanetaryHour: planetaryHour,
			BestFor:       bestFor,
		}
	}

	return breakdown
}

// calculateHourlyScore 计算小时分数
func calculateHourlyScore(_ *models.NatalChart, _ time.Time, _ int, planetaryHour models.PlanetID) float64 {
	base := 50.0

	// 根据行星时调整
	switch planetaryHour {
	case models.Venus, models.Jupiter:
		base += 10 // 吉星时
	case models.Saturn, models.Mars:
		base -= 5 // 凶星时
	case models.Sun, models.Moon:
		base += 5 // 发光体时
	}

	return normalizeScore(base - 50)
}

// getBestActivitiesForHour 获取该小时最佳活动
func getBestActivitiesForHour(planetaryHour models.PlanetID) []string {
	activities := map[models.PlanetID][]string{
		models.Saturn:  {"计划", "冥想", "独处", "整理"},
		models.Jupiter: {"学习", "旅行", "社交", "决策"},
		models.Mars:    {"运动", "竞争", "启动项目", "主动行动"},
		models.Sun:     {"领导", "创意", "展示自我", "重要会议"},
		models.Venus:   {"艺术", "社交", "约会", "美容"},
		models.Mercury: {"沟通", "写作", "商务", "学习"},
		models.Moon:    {"家庭", "休息", "情感交流", "直觉决策"},
	}

	if acts, ok := activities[planetaryHour]; ok {
		return acts
	}
	return []string{"一般活动"}
}

// generateDailyTheme 生成每日主题
func generateDailyTheme(moonSign models.MoonSignInfo, _ []models.AspectData) string {
	// 基于月亮星座生成主题
	themes := map[models.ZodiacID]string{
		models.Aries:       "行动与开始",
		models.Taurus:      "稳定与享受",
		models.Gemini:      "沟通与学习",
		models.Cancer:      "情感与家庭",
		models.Leo:         "表达与创造",
		models.Virgo:       "分析与服务",
		models.Libra:       "关系与平衡",
		models.Scorpio:     "转化与深度",
		models.Sagittarius: "探索与扩张",
		models.Capricorn:   "成就与责任",
		models.Aquarius:    "创新与自由",
		models.Pisces:      "灵性与直觉",
	}

	if theme, ok := themes[moonSign.Sign]; ok {
		return theme
	}
	return "一般日子"
}

// generateWeeklyTheme 生成每周主题
func generateWeeklyTheme(transits []models.WeeklyTransit) string {
	if len(transits) == 0 {
		return "平稳的一周"
	}

	// 根据最重要的行运生成主题
	themes := []string{}
	for _, t := range transits {
		themes = append(themes, t.Theme)
	}

	if len(themes) > 0 {
		return themes[0]
	}
	return "多元发展的一周"
}

// getDayOfWeekChinese 获取中文星期
func getDayOfWeekChinese(day time.Weekday) string {
	days := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
	return days[day]
}

// findKeyDates 找出关键日期
func findKeyDates(summaries []models.DailySummary) []models.KeyDate {
	var keyDates []models.KeyDate

	for _, s := range summaries {
		if s.OverallScore >= 75 {
			keyDates = append(keyDates, models.KeyDate{
				Date:         s.Date,
				Event:        "高能量日",
				Significance: "high",
			})
		} else if s.OverallScore <= 35 {
			keyDates = append(keyDates, models.KeyDate{
				Date:         s.Date,
				Event:        "挑战日",
				Significance: "medium",
			})
		}
	}

	return keyDates
}

// findBestDaysFor 找出各活动最佳日期
func findBestDaysFor(chart *models.NatalChart, startDate time.Time) map[string][]string {
	bestDays := map[string][]string{
		"career":       {},
		"relationship": {},
		"health":       {},
	}

	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		forecast := CalculateDailyForecast(chart, date, false)
		dateStr := date.Format("2006-01-02")

		if forecast.Dimensions.Career >= 65 {
			bestDays["career"] = append(bestDays["career"], dateStr)
		}
		if forecast.Dimensions.Relationship >= 65 {
			bestDays["relationship"] = append(bestDays["relationship"], dateStr)
		}
		if forecast.Dimensions.Health >= 65 {
			bestDays["health"] = append(bestDays["health"], dateStr)
		}
	}

	return bestDays
}

// getWeeklyTransits 获取周度行运
func getWeeklyTransits(chart *models.NatalChart, startDate, _ time.Time) []models.WeeklyTransit {
	var transits []models.WeeklyTransit

	// 简化实现：检查一周内的重要行运
	midWeek := startDate.AddDate(0, 0, 3)
	transitPositions := GetTransitPositions(midWeek)
	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)

	for _, asp := range aspects {
		if asp.Strength > 0.7 {
			transits = append(transits, models.WeeklyTransit{
				Planet:      asp.Planet1,
				Aspect:      asp.AspectType,
				NatalPlanet: asp.Planet2,
				Peak:        midWeek.Format("2006-01-02"),
				Theme:       asp.Interpretation,
			})
		}
	}

	return transits
}

// calculateDimensionTrends 计算维度趋势
func calculateDimensionTrends(_ []models.DailySummary) map[string]string {
	// 简化实现：比较周初和周末
	return map[string]string{
		"career":       "stable",
		"relationship": "stable",
		"health":       "stable",
	}
}

// getTopFactors 获取最重要的因子
func getTopFactors(factors *models.FactorResult, n int) []models.InfluenceFactor {
	if factors == nil || len(factors.Factors) == 0 {
		return nil
	}

	// 按绝对值排序
	sorted := make([]models.InfluenceFactor, len(factors.Factors))
	copy(sorted, factors.Factors)

	// 简单冒泡排序
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if math.Abs(sorted[j].Adjustment) < math.Abs(sorted[j+1].Adjustment) {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	if len(sorted) > n {
		return sorted[:n]
	}
	return sorted
}

