package astro

import (
	"star/models"
	"time"
)

// CalculateLifeTrend 计算人生趋势
func CalculateLifeTrend(chart *models.NatalChart, startYear, endYear int, resolution string) *models.LifeTrendData {
	birthYear := chart.BirthData.ToTime().Year()

	// 设置默认值
	if startYear == 0 {
		startYear = birthYear
	}
	if endYear == 0 {
		endYear = birthYear + 80
	}
	if resolution == "" {
		resolution = "yearly"
	}

	var points []models.LifeTrendPoint

	// 根据分辨率生成数据点
	switch resolution {
	case "yearly":
		points = generateYearlyTrend(chart, startYear, endYear)
	case "quarterly":
		points = generateQuarterlyTrend(chart, startYear, endYear)
	case "monthly":
		points = generateMonthlyTrend(chart, startYear, endYear)
	default:
		points = generateYearlyTrend(chart, startYear, endYear)
	}

	// 生成摘要
	summary := generateLifeTrendSummary(points)

	// 生成周期信息
	cycles := generateLifeCycles(chart, startYear, endYear)

	return &models.LifeTrendData{
		Type:      resolution,
		BirthDate: chart.BirthData.ToTime(),
		Points:    points,
		Summary:   summary,
		Cycles:    cycles,
	}
}

// generateYearlyTrend 生成年度趋势
func generateYearlyTrend(chart *models.NatalChart, startYear, endYear int) []models.LifeTrendPoint {
	var points []models.LifeTrendPoint
	birthYear := chart.BirthData.ToTime().Year()

	for year := startYear; year <= endYear; year++ {
		age := year - birthYear
		if age < 0 {
			continue
		}

		date := time.Date(year, 6, 15, 12, 0, 0, 0, time.UTC) // 使用年中点
		point := calculateLifeTrendPoint(chart, date, year, age)
		points = append(points, point)
	}

	return points
}

// generateQuarterlyTrend 生成季度趋势
func generateQuarterlyTrend(chart *models.NatalChart, startYear, endYear int) []models.LifeTrendPoint {
	var points []models.LifeTrendPoint
	birthYear := chart.BirthData.ToTime().Year()

	for year := startYear; year <= endYear; year++ {
		age := year - birthYear
		if age < 0 {
			continue
		}

		for quarter := 1; quarter <= 4; quarter++ {
			month := (quarter-1)*3 + 2 // 每季度中间月
			date := time.Date(year, time.Month(month), 15, 12, 0, 0, 0, time.UTC)
			point := calculateLifeTrendPoint(chart, date, year, age)
			points = append(points, point)
		}
	}

	return points
}

// generateMonthlyTrend 生成月度趋势
func generateMonthlyTrend(chart *models.NatalChart, startYear, endYear int) []models.LifeTrendPoint {
	var points []models.LifeTrendPoint
	birthYear := chart.BirthData.ToTime().Year()

	for year := startYear; year <= endYear; year++ {
		age := year - birthYear
		if age < 0 {
			continue
		}

		for month := 1; month <= 12; month++ {
			date := time.Date(year, time.Month(month), 15, 12, 0, 0, 0, time.UTC)
			point := calculateLifeTrendPoint(chart, date, year, age)
			points = append(points, point)
		}
	}

	return points
}

// calculateLifeTrendPoint 计算单个人生趋势点
func calculateLifeTrendPoint(chart *models.NatalChart, date time.Time, year, age int) models.LifeTrendPoint {
	// 获取行运位置
	transitPositions := GetTransitPositions(date)

	// 计算行运相位
	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)
	transitScore := CalculateTransitScore(aspects)

	// 计算年限法
	profection := CalculateAnnualProfection(chart, age)

	// 计算推运月相
	lunarPhase := GetProgressedLunarPhaseForAge(chart, age)

	// 检查是否有重大行运
	majorTransits := FindMajorTransits(chart, year, year)
	var isMajorTransit bool
	var majorTransitName string
	if len(majorTransits) > 0 {
		isMajorTransit = true
		majorTransitName = majorTransits[0].Name
	}

	// 计算各项分数
	harmonious := transitScore.Harmonious + 30 // 基础值
	challenge := transitScore.Tense
	transformation := 0.0

	// 转化分数 - 与外行星相关
	for _, asp := range aspects {
		if asp.Planet1 == models.Pluto || asp.Planet2 == models.Pluto ||
			asp.Planet1 == models.Uranus || asp.Planet2 == models.Uranus {
			transformation += asp.Weight
		}
	}

	// 整体分数
	overallScore := normalizeScore(transitScore.Total)

	// 维度分数
	dimensions := models.DailyDimensions{
		Career:       normalizeScore(transitScore.Total * 0.9),
		Relationship: normalizeScore(transitScore.Total * 1.1),
		Health:       normalizeScore(transitScore.Total * 0.85),
		Finance:      normalizeScore(transitScore.Total * 0.95),
		Spiritual:    normalizeScore(transitScore.Total + transformation*0.5),
	}

	// 确定主导行星
	dominantPlanet := findDominantTransitPlanet(aspects)

	return models.LifeTrendPoint{
		Date:             date,
		Year:             year,
		Age:              age,
		OverallScore:     overallScore,
		Harmonious:       harmonious,
		Challenge:        challenge,
		Transformation:   transformation,
		Dimensions:       dimensions,
		DominantPlanet:   dominantPlanet,
		Profection:       models.ProfectionSummary{
			House:      profection.House,
			Theme:      profection.HouseTheme,
			LordOfYear: profection.LordOfYear,
		},
		IsMajorTransit:   isMajorTransit,
		MajorTransitName: majorTransitName,
		LunarPhaseName:   lunarPhase.Name,
		LunarPhaseAngle:  lunarPhase.Angle,
	}
}

// findDominantTransitPlanet 找出主导行运行星
func findDominantTransitPlanet(aspects []models.AspectData) models.PlanetID {
	planetScores := make(map[models.PlanetID]float64)

	for _, asp := range aspects {
		planetScores[asp.Planet1] += asp.Weight
	}

	var dominant models.PlanetID
	var maxScore float64
	for p, score := range planetScores {
		if score > maxScore {
			maxScore = score
			dominant = p
		}
	}

	if dominant == "" {
		return models.Sun
	}
	return dominant
}

// generateLifeTrendSummary 生成人生趋势摘要
func generateLifeTrendSummary(points []models.LifeTrendPoint) models.LifeTrendSummary {
	var peakYears, challengeYears, transformationYears []int

	for _, p := range points {
		if p.OverallScore >= 70 {
			peakYears = append(peakYears, p.Year)
		}
		if p.Challenge >= 30 {
			challengeYears = append(challengeYears, p.Year)
		}
		if p.Transformation >= 20 {
			transformationYears = append(transformationYears, p.Year)
		}
	}

	// 判断整体趋势
	overallTrend := "fluctuating"
	if len(points) > 10 {
		firstHalf := 0.0
		secondHalf := 0.0
		mid := len(points) / 2
		for i, p := range points {
			if i < mid {
				firstHalf += p.OverallScore
			} else {
				secondHalf += p.OverallScore
			}
		}
		if secondHalf > firstHalf*1.1 {
			overallTrend = "rising"
		} else if firstHalf > secondHalf*1.1 {
			overallTrend = "declining"
		}
	}

	return models.LifeTrendSummary{
		OverallTrend:        overallTrend,
		PeakYears:           peakYears,
		ChallengeYears:      challengeYears,
		TransformationYears: transformationYears,
	}
}

// generateLifeCycles 生成人生周期信息
func generateLifeCycles(chart *models.NatalChart, _, _ int) models.LifeCycles {
	birthYear := chart.BirthData.ToTime().Year()

	// 土星周期
	saturnCycles := []models.SaturnCycle{
		{Age: 29, Year: birthYear + 29, Description: "第一次土星回归"},
		{Age: 58, Year: birthYear + 58, Description: "第二次土星回归"},
	}

	// 木星周期
	jupiterCycles := make([]map[string]interface{}, 0)
	for age := 12; age <= 72; age += 12 {
		jupiterCycles = append(jupiterCycles, map[string]interface{}{
			"age":         age,
			"year":        birthYear + age,
			"description": "木星回归",
		})
	}

	// 年限法周期
	profectionCycles := []map[string]interface{}{
		{"startAge": 0, "endAge": 11, "theme": "第一轮周期 - 探索与学习"},
		{"startAge": 12, "endAge": 23, "theme": "第二轮周期 - 建立自我"},
		{"startAge": 24, "endAge": 35, "theme": "第三轮周期 - 成就与责任"},
		{"startAge": 36, "endAge": 47, "theme": "第四轮周期 - 整合与深化"},
		{"startAge": 48, "endAge": 59, "theme": "第五轮周期 - 智慧与传承"},
		{"startAge": 60, "endAge": 71, "theme": "第六轮周期 - 超越与觉醒"},
	}

	return models.LifeCycles{
		SaturnCycles:     saturnCycles,
		JupiterCycles:    jupiterCycles,
		ProfectionCycles: profectionCycles,
	}
}

