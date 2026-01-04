package astro

import (
	"math"
	"star/models"
	"time"
)

// CalculateTimeSeries 计算统一时间序列
func CalculateTimeSeries(chart *models.NatalChart, startStr, endStr, granularity string) *models.TimeSeries {
	// 解析时间范围
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		start = time.Now()
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		end = start.AddDate(0, 1, 0) // 默认一个月
	}

	// 验证粒度
	var g models.TimeGranularity
	switch granularity {
	case "hour":
		g = models.GranularityHour
	case "day":
		g = models.GranularityDay
	case "week":
		g = models.GranularityWeek
	case "month":
		g = models.GranularityMonth
	case "year":
		g = models.GranularityYear
	default:
		g = models.GranularityDay
	}

	// 生成数据点
	points := generateTimeSeriesPoints(chart, start, end, g)

	// 计算统计信息
	stats := calculateTimeSeriesStats(points)

	return &models.TimeSeries{
		Granularity: g,
		Points:      points,
		Stats:       stats,
	}
}

// generateTimeSeriesPoints 生成时间序列数据点
func generateTimeSeriesPoints(chart *models.NatalChart, start, end time.Time, granularity models.TimeGranularity) []models.TimeSeriesPoint {
	var points []models.TimeSeriesPoint
	current := start

	for current.Before(end) || current.Equal(end) {
		point := calculateTimeSeriesPoint(chart, current, granularity)
		points = append(points, point)

		// 移动到下一个时间点
		switch granularity {
		case models.GranularityHour:
			current = current.Add(time.Hour)
		case models.GranularityDay:
			current = current.AddDate(0, 0, 1)
		case models.GranularityWeek:
			current = current.AddDate(0, 0, 7)
		case models.GranularityMonth:
			current = current.AddDate(0, 1, 0)
		case models.GranularityYear:
			current = current.AddDate(1, 0, 0)
		}
	}

	return points
}

// calculateTimeSeriesPoint 计算单个时间序列点
func calculateTimeSeriesPoint(chart *models.NatalChart, t time.Time, granularity models.TimeGranularity) models.TimeSeriesPoint {
	// 获取行运位置
	transitPositions := GetTransitPositions(t)

	// 计算行运相位
	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)
	transitScore := CalculateTransitScore(aspects)

	// 计算因子
	factors := CalculateInfluenceFactors(chart, t, transitPositions)

	// 构建分数因子
	// 从因子列表中提取尊贵度和逆行的总调整
	var dignityTotal, retrogradeTotal float64
	for _, f := range factors.Factors {
		switch f.Type {
		case models.FactorDignity:
			dignityTotal += f.Adjustment
		case models.FactorRetrograde:
			retrogradeTotal += f.Adjustment
		}
	}
	scoreFactors := models.ScoreFactors{
		Aspects:    transitScore.Total,
		Dignity:    dignityTotal,
		Retrograde: retrogradeTotal,
	}

	// 根据粒度添加特定因子
	switch granularity {
	case models.GranularityHour:
		planetaryHour := CalculatePlanetaryHour(t, t.Hour())
		scoreFactors.PlanetaryHour = getPlanetaryHourValue(planetaryHour)
	case models.GranularityDay, models.GranularityWeek:
		scoreFactors.MoonSign = getMoonSignInfluence(transitPositions, chart)
		scoreFactors.MoonPhase = getMoonPhaseInfluence(transitPositions)
	}

	// 计算原始分数
	rawValue := transitScore.Total + factors.TotalAdjustment

	// 标准化显示分数
	displayScore := normalizeScore(rawValue)

	// 计算维度分数
	dimensions := map[string]float64{
		"career":       normalizeScore(rawValue * 0.9),
		"relationship": normalizeScore(rawValue * 1.1),
		"health":       normalizeScore(rawValue * 0.85),
		"finance":      normalizeScore(rawValue * 0.95),
		"spiritual":    normalizeScore(rawValue * 1.05),
	}

	// 生成标签
	label := generateTimeLabel(t, granularity)

	return models.TimeSeriesPoint{
		Time:        t,
		Label:       label,
		Granularity: granularity,
		Raw: models.RawScore{
			Value:   rawValue,
			Factors: scoreFactors,
		},
		Display:    displayScore,
		Dimensions: dimensions,
	}
}

// getPlanetaryHourValue 获取行星时影响值
func getPlanetaryHourValue(planet models.PlanetID) float64 {
	values := map[models.PlanetID]float64{
		models.Sun:     5,
		models.Moon:    3,
		models.Mercury: 2,
		models.Venus:   5,
		models.Mars:    -2,
		models.Jupiter: 5,
		models.Saturn:  -3,
	}
	if v, ok := values[planet]; ok {
		return v
	}
	return 0
}

// getMoonSignInfluence 获取月亮星座影响
func getMoonSignInfluence(transitPositions []models.PlanetPosition, chart *models.NatalChart) float64 {
	for _, p := range transitPositions {
		if p.ID == models.Moon {
			// 检查是否与本命盘有吉相位
			natalMoon := GetPlanetFromChart(chart, models.Moon)
			if natalMoon != nil {
				if p.Sign == natalMoon.Sign {
					return 5 // 回归本命月亮星座
				}
				// 检查元素相合
				transitZodiac := GetZodiacInfo(p.Sign)
				natalZodiac := GetZodiacInfo(natalMoon.Sign)
				if transitZodiac != nil && natalZodiac != nil {
					if transitZodiac.Element == natalZodiac.Element {
						return 3
					}
				}
			}
			return 0
		}
	}
	return 0
}

// getMoonPhaseInfluence 获取月相影响
func getMoonPhaseInfluence(transitPositions []models.PlanetPosition) float64 {
	var sunLon, moonLon float64
	for _, p := range transitPositions {
		if p.ID == models.Sun {
			sunLon = p.Longitude
		}
		if p.ID == models.Moon {
			moonLon = p.Longitude
		}
	}

	angle := NormalizeAngle(moonLon - sunLon)
	phaseInfo := GetLunarPhase(angle)

	// 月相影响值
	values := map[string]float64{
		"new":          2,
		"crescent":     3,
		"firstQuarter": 0,
		"gibbous":      2,
		"full":         5,
		"disseminating": 1,
		"lastQuarter":  -2,
		"balsamic":     -1,
	}

	if v, ok := values[phaseInfo.Phase]; ok {
		return v
	}
	return 0
}

// generateTimeLabel 生成时间标签
func generateTimeLabel(t time.Time, granularity models.TimeGranularity) string {
	switch granularity {
	case models.GranularityHour:
		return t.Format("15:00")
	case models.GranularityDay:
		return t.Format("01-02")
	case models.GranularityWeek:
		return t.Format("01-02") + "周"
	case models.GranularityMonth:
		return t.Format("2006-01")
	case models.GranularityYear:
		return t.Format("2006")
	default:
		return t.Format("2006-01-02")
	}
}

// calculateTimeSeriesStats 计算时间序列统计
func calculateTimeSeriesStats(points []models.TimeSeriesPoint) models.TimeSeriesStats {
	if len(points) == 0 {
		return models.TimeSeriesStats{}
	}

	// 计算均值
	var sum float64
	min := points[0].Display
	max := points[0].Display

	for _, p := range points {
		sum += p.Display
		if p.Display < min {
			min = p.Display
		}
		if p.Display > max {
			max = p.Display
		}
	}

	mean := sum / float64(len(points))

	// 计算标准差
	var sumSquares float64
	for _, p := range points {
		diff := p.Display - mean
		sumSquares += diff * diff
	}
	stdDev := math.Sqrt(sumSquares / float64(len(points)))

	// 计算趋势
	trend := "flat"
	if len(points) >= 2 {
		firstHalf := 0.0
		secondHalf := 0.0
		mid := len(points) / 2

		for i, p := range points {
			if i < mid {
				firstHalf += p.Display
			} else {
				secondHalf += p.Display
			}
		}

		firstAvg := firstHalf / float64(mid)
		secondAvg := secondHalf / float64(len(points)-mid)

		if secondAvg > firstAvg*1.05 {
			trend = "up"
		} else if firstAvg > secondAvg*1.05 {
			trend = "down"
		}
	}

	// 计算波动率 (类似 ATR)
	volatility := stdDev / mean * 100

	return models.TimeSeriesStats{
		Mean:       mean,
		StdDev:     stdDev,
		Min:        min,
		Max:        max,
		Trend:      trend,
		Volatility: volatility,
	}
}

// CalculateHourlyTimeSeries 计算小时级时间序列
func CalculateHourlyTimeSeries(chart *models.NatalChart, date time.Time) []models.TimeSeriesPoint {
	var points []models.TimeSeriesPoint

	for hour := 0; hour < 24; hour++ {
		t := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
		point := calculateTimeSeriesPoint(chart, t, models.GranularityHour)
		points = append(points, point)
	}

	return points
}

// AggregateToDaily 将小时级数据聚合为日级
func AggregateToDaily(hourlyPoints []models.TimeSeriesPoint) models.TimeSeriesPoint {
	if len(hourlyPoints) == 0 {
		return models.TimeSeriesPoint{}
	}

	var sumRaw, sumDisplay float64
	dimensions := make(map[string]float64)

	for _, p := range hourlyPoints {
		sumRaw += p.Raw.Value
		sumDisplay += p.Display
		for dim, val := range p.Dimensions {
			dimensions[dim] += val
		}
	}

	n := float64(len(hourlyPoints))
	avgRaw := sumRaw / n
	avgDisplay := sumDisplay / n

	for dim := range dimensions {
		dimensions[dim] /= n
	}

	// 计算波动率
	var sumSquares float64
	for _, p := range hourlyPoints {
		diff := p.Display - avgDisplay
		sumSquares += diff * diff
	}
	volatility := math.Sqrt(sumSquares / n)

	return models.TimeSeriesPoint{
		Time:        hourlyPoints[0].Time,
		Label:       hourlyPoints[0].Time.Format("01-02"),
		Granularity: models.GranularityDay,
		Raw: models.RawScore{
			Value: avgRaw,
		},
		Display:    avgDisplay,
		Dimensions: dimensions,
		Volatility: volatility,
	}
}

