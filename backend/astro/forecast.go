package astro

import (
	"math"
	"star/models"
	"time"
)

/*
每日/每周预测

使用统一分数计算架构：
- 日分数 = 24小时平均
- 周分数 = 7天平均
- 五维分数基于行星-维度映射独立计算
*/

// CalculateDailyForecast 计算每日预测
func CalculateDailyForecast(chart *models.NatalChart, date time.Time, withFactors bool) *models.DailyForecast {
	// 使用统一计算逻辑获取日分数（24小时平均）
	dailyScore := CalculateDailyScore(chart, date)

	// 获取当前行星位置（用于月亮信息和相位）
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

	// 计算影响因子
	var factors *models.FactorResult
	var topFactors []models.InfluenceFactor
	if withFactors {
		factors = CalculateInfluenceFactors(chart, date, transitPositions)
		topFactors = getTopFactors(factors, 5)
	}

	// 构建五维分数（从统一计算结果中获取）
	dimensions := models.DailyDimensions{
		Career:       dailyScore.Dimensions["career"],
		Relationship: dailyScore.Dimensions["relationship"],
		Health:       dailyScore.Dimensions["health"],
		Finance:      dailyScore.Dimensions["finance"],
		Spiritual:    dailyScore.Dimensions["spiritual"],
	}

	// 计算小时预测（使用统一计算逻辑）
	hourlyBreakdown := calculateUnifiedHourlyBreakdown(chart, date)

	// 生成主题
	overallTheme := generateDailyTheme(moonSign, activeAspects)

	// 获取星期
	dayOfWeek := getDayOfWeekChinese(date.Weekday())

	return &models.DailyForecast{
		Date:            date,
		DayOfWeek:       dayOfWeek,
		OverallScore:    dailyScore.Overall,
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

	// 使用统一计算逻辑获取周分数（7天平均）
	weeklyScore := CalculateWeeklyScore(chart, startDate)

	// 计算每日摘要
	dailySummaries := make([]models.DailySummary, 7)
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
	}

	// 构建五维分数（从统一计算结果中获取）
	dimensions := models.DailyDimensions{
		Career:       weeklyScore.Dimensions["career"],
		Relationship: weeklyScore.Dimensions["relationship"],
		Health:       weeklyScore.Dimensions["health"],
		Finance:      weeklyScore.Dimensions["finance"],
		Spiritual:    weeklyScore.Dimensions["spiritual"],
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
		OverallScore:    weeklyScore.Overall,
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

// calculateUnifiedHourlyBreakdown 使用统一计算逻辑计算小时预测
func calculateUnifiedHourlyBreakdown(chart *models.NatalChart, date time.Time) []models.HourlyForecast {
	breakdown := make([]models.HourlyForecast, 24)

	for hour := 0; hour < 24; hour++ {
		t := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
		hourlyScore := CalculateUnifiedHourlyScore(chart, t)

		planetaryHour := CalculatePlanetaryHour(date, hour)
		bestFor := getBestActivitiesForHour(planetaryHour)

		breakdown[hour] = models.HourlyForecast{
			Hour:          hour,
			Score:         hourlyScore.Overall,
			PlanetaryHour: planetaryHour,
			BestFor:       bestFor,
		}
	}

	return breakdown
}

// normalizeScore 将原始分数标准化到 0-100
// 确保输出在 0-100 范围内，保留4位小数
func normalizeScore(raw float64) float64 {
	// 使用 tanh 压缩
	scale := 20.0
	normalized := 50 + 50*math.Tanh(raw/scale)
	result := math.Max(0, math.Min(100, normalized))
	// 保留4位小数
	return math.Round(result*10000) / 10000
}

// generateDailyTheme 生成每日主题
func generateDailyTheme(moonSign models.MoonSignInfo, _ []models.AspectData) string {
	// 基于月亮星座生成主题
	themes := map[models.ZodiacID]string{
		models.Aries:       "Action & Initiative",
		models.Taurus:      "Stability & Enjoyment",
		models.Gemini:      "Communication & Learning",
		models.Cancer:      "Emotion & Family",
		models.Leo:         "Expression & Creativity",
		models.Virgo:       "Analysis & Service",
		models.Libra:       "Relationships & Balance",
		models.Scorpio:     "Transformation & Depth",
		models.Sagittarius: "Exploration & Expansion",
		models.Capricorn:   "Achievement & Responsibility",
		models.Aquarius:    "Innovation & Freedom",
		models.Pisces:      "Spirituality & Intuition",
	}

	if theme, ok := themes[moonSign.Sign]; ok {
		return theme
	}
	return "Normal Day"
}

// generateWeeklyTheme 生成每周主题
func generateWeeklyTheme(transits []models.WeeklyTransit) string {
	if len(transits) == 0 {
		return "Steady Week"
	}

	// 根据最重要的行运生成主题
	themes := []string{}
	for _, t := range transits {
		themes = append(themes, t.Theme)
	}

	if len(themes) > 0 {
		return themes[0]
	}
	return "Week of Diverse Development"
}

// getDayOfWeekChinese 获取星期名称
func getDayOfWeekChinese(day time.Weekday) string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	return days[day]
}

// findKeyDates 找出关键日期
func findKeyDates(summaries []models.DailySummary) []models.KeyDate {
	var keyDates []models.KeyDate

	for _, s := range summaries {
		if s.OverallScore >= 75 {
			keyDates = append(keyDates, models.KeyDate{
				Date:         s.Date,
				Event:        "High Energy Day",
				Significance: "high",
			})
		} else if s.OverallScore <= 35 {
			keyDates = append(keyDates, models.KeyDate{
				Date:         s.Date,
				Event:        "Challenge Day",
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

// getBestActivitiesForHour 获取该小时最佳活动
func getBestActivitiesForHour(planetaryHour models.PlanetID) []string {
	activities := map[models.PlanetID][]string{
		models.Saturn:  {"planning", "meditation", "solitude", "organizing"},
		models.Jupiter: {"learning", "travel", "networking", "decision-making"},
		models.Mars:    {"exercise", "competition", "initiating projects", "taking action"},
		models.Sun:     {"leadership", "creativity", "self-expression", "important meetings"},
		models.Venus:   {"art", "socializing", "dating", "beauty"},
		models.Mercury: {"communication", "writing", "business", "learning"},
		models.Moon:    {"family", "rest", "emotional connection", "intuitive decisions"},
	}

	if acts, ok := activities[planetaryHour]; ok {
		return acts
	}
	return []string{"general activities"}
}
