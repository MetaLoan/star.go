package core

import (
	"sort"
	"star/astro"
	"star/i18n"
	"star/models"
	"strings"
	"time"
)

// ==================== 统一计算器 ====================
// 合并 Event 和 Factor 的计算逻辑
// 产出：TimeSlot（不含 impactDelta，由 delta.go 动态计算）

// Calculator 统一计算器
type Calculator struct {
	chart    *models.NatalChart
	language string
}

// NewCalculator 创建计算器
func NewCalculator(chart *models.NatalChart, language string) *Calculator {
	if language == "" {
		language = "en"
	}
	return &Calculator{
		chart:    chart,
		language: language,
	}
}

// CalculateHour 计算单个小时的 TimeSlot
// 这是最基础的计算单元
func (c *Calculator) CalculateHour(t time.Time) *TimeSlot {
	// 归一化到小时开始
	hourStart := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	hourEnd := hourStart.Add(time.Hour)

	// 创建时间槽
	slot := NewTimeSlot(c.getUserID(), hourStart, hourEnd, GranularityHour)

	// 1. 获取行运位置（只获取一次，复用）
	transitPositions := astro.GetTransitPositions(hourStart)

	// 2. 计算分数（传入已获取的位置）
	scoreResult := astro.CalculateScoresV2WithPositions(c.chart, hourStart, transitPositions)
	slot.Scores = FromModelsDimensionScores(scoreResult.Dimensions)
	slot.Scores.Overall = scoreResult.Overall

	// 3. 获取天象数据并转换为事件（复用位置）
	events := c.calculateEventsWithPositions(hourStart, hourEnd, transitPositions)
	slot.Events = events

	// 4. 生成指导（基于事件和分数）
	slot.Guidance = c.generateGuidance(slot)

	return slot
}

// CalculateHourScoreOnly 只计算分数，不计算事件（用于大粒度聚合和趋势）
// 使用轻量版算法，跳过精确相位时间搜索，性能提升约 100 倍
func (c *Calculator) CalculateHourScoreOnly(t time.Time) *TimeSlot {
	hourStart := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	hourEnd := hourStart.Add(time.Hour)

	slot := NewTimeSlot(c.getUserID(), hourStart, hourEnd, GranularityHour)

	// 使用轻量版分数计算（跳过精确相位搜索）
	scoreResult := astro.CalculateScoresV2Lite(c.chart, hourStart)
	slot.Scores = FromModelsDimensionScores(scoreResult.Dimensions)
	slot.Scores.Overall = scoreResult.Overall

	return slot
}

// calculateEvents 计算时间范围内的所有事件
func (c *Calculator) calculateEvents(start, end time.Time) []AstroEvent {
	transitPositions := astro.GetTransitPositions(start)
	return c.calculateEventsWithPositions(start, end, transitPositions)
}

// calculateEventsWithPositions 使用已有的行运位置计算事件
func (c *Calculator) calculateEventsWithPositions(start, end time.Time, transitPositions []models.PlanetPosition) []AstroEvent {
	events := make([]AstroEvent, 0)

	// 1. 计算相位事件
	aspectEvents := c.calculateAspectEvents(start, transitPositions)
	events = append(events, aspectEvents...)

	// 2. 计算行星时事件
	planetaryHourEvent := c.calculatePlanetaryHourEvent(start)
	if planetaryHourEvent != nil {
		events = append(events, *planetaryHourEvent)
	}

	// 3. 计算逆行事件
	retrogradeEvents := c.calculateRetrogradeEvents(transitPositions, start)
	events = append(events, retrogradeEvents...)

	// 4. 计算月空事件
	vocEvent := c.calculateVoidOfCourseEvent(start)
	if vocEvent != nil {
		events = append(events, *vocEvent)
	}

	// 5. 计算月相事件
	lunarPhaseEvents := c.calculateLunarPhaseEvents(start, end)
	events = append(events, lunarPhaseEvents...)

	// 6. 计算换座事件
	signChangeEvents := c.calculateSignChangeEvents(start, end)
	events = append(events, signChangeEvents...)

	// 7. 计算尊贵度事件
	dignityEvents := c.calculateDignityEvents(transitPositions, start)
	events = append(events, dignityEvents...)

	// 8. 按时间排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].ExactTime.Before(events[j].ExactTime)
	})

	// 9. 应用 i18n
	translator := i18n.New(c.language)
	for i := range events {
		c.applyI18n(&events[i], translator)
	}

	return events
}

// calculateAspectEvents 计算相位事件
func (c *Calculator) calculateAspectEvents(t time.Time, transitPositions []models.PlanetPosition) []AstroEvent {
	events := make([]AstroEvent, 0)

	// 使用现有的相位计算逻辑
	factors := astro.CalculateInfluenceFactorsV2(c.chart, t, transitPositions)
	if factors == nil {
		return events
	}

	// 遍历所有因子，转换为事件
	allFactors := append(factors.DailyFactors, factors.HourlyFactors...)
	allFactors = append(allFactors, factors.WeeklyFactors...)
	allFactors = append(allFactors, factors.MonthlyFactors...)
	allFactors = append(allFactors, factors.YearlyFactors...)

	for _, f := range allFactors {
		if f.Type == models.FactorAspectPhase {
			event := c.factorToAspectEvent(&f, t)
			if event != nil {
				events = append(events, *event)
			}
		}
	}

	return events
}

// factorToAspectEvent 将因子转换为相位事件
func (c *Calculator) factorToAspectEvent(f *models.InfluenceFactor, t time.Time) *AstroEvent {
	// 确定时间层级
	timeLevel := c.getTimeLevel(f.TimeLevel)

	// 从 Name 解析行星和相位信息
	// Name 格式如: "Mars Trine Venus" 或 "Mercury Conjunction Sun"
	primaryPlanet, secondaryPlanet, aspect := parseFactorName(f.Name)
	if primaryPlanet == "" {
		primaryPlanet = string(f.SourcePlanet)
	}

	// 计算生命周期时间
	var startTime, endTime, exactTime time.Time
	if f.Lifecycle != nil {
		startTime = f.Lifecycle.StartTime
		endTime = f.Lifecycle.EndTime
		exactTime = f.Lifecycle.PeakTime
	} else {
		startTime = t.Add(-12 * time.Hour)
		endTime = t.Add(12 * time.Hour)
		exactTime = t
	}

	// 生成事件 ID（生命周期稳定 + 兼容无生命周期）
	var eventID string
	if f.Lifecycle != nil {
		eventID = GenerateAspectLifecycleEventID(models.PlanetID(primaryPlanet), models.PlanetID(secondaryPlanet), aspect, startTime, endTime)
	} else {
		eventID = GenerateEventID(EventTypeAspect, models.PlanetID(primaryPlanet), models.PlanetID(secondaryPlanet), aspect, exactTime)
	}

	return &AstroEvent{
		EventID:         eventID,
		Type:            EventTypeAspect,
		Title:           f.Description,
		IsPositive:      f.IsPositive,
		Intensity:       normalizeIntensity(f.CurrentStrength),
		TimeLevel:       timeLevel,
		PrimaryPlanet:   primaryPlanet,
		SecondaryPlanet: secondaryPlanet,
		Aspect:          aspect,
		Impact:          c.factorToImpact(f),
		StartTime:       startTime,
		EndTime:         endTime,
		ExactTime:       exactTime,
	}
}

// parseFactorName 解析因子名称获取行星和相位信息
// 格式: "Planet1 AspectType Planet2" 如 "Mars Trine Venus" 或 "North Node Opposition Venus"
func parseFactorName(name string) (primary, secondary, aspect string) {
	parts := splitBySpace(name)
	if len(parts) < 3 {
		if len(parts) >= 1 {
			primary = toLowerPlanetID(parts[0])
		}
		return
	}

	// 处理多词行星名称（如 "North Node", "South Node"）
	// 已知相位词列表
	aspectWords := map[string]bool{
		"conjunction": true, "trine": true, "square": true,
		"opposition": true, "sextile": true, "quincunx": true,
		"semi-sextile": true, "semi-square": true, "sesquiquadrate": true,
	}

	// 找到相位词的位置
	aspectIdx := -1
	for i, part := range parts {
		if aspectWords[toLower(part)] {
			aspectIdx = i
			break
		}
	}

	if aspectIdx == -1 {
		// 没找到相位词，使用原有逻辑
		primary = toLowerPlanetID(parts[0])
		aspect = toLowerAspect(parts[1])
		if len(parts) > 2 {
			secondary = toLowerPlanetID(parts[2])
		}
		return
	}

	// 相位词之前的部分是第一个行星
	primaryParts := parts[:aspectIdx]
	primary = toLowerPlanetID(joinWithSpace(primaryParts))

	// 相位词
	aspect = toLowerAspect(parts[aspectIdx])

	// 相位词之后的部分是第二个行星
	if aspectIdx+1 < len(parts) {
		secondaryParts := parts[aspectIdx+1:]
		secondary = toLowerPlanetID(joinWithSpace(secondaryParts))
	}

	return
}

// joinWithSpace 用空格连接字符串数组
func joinWithSpace(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}
	return result
}

// splitBySpace 按空格分割字符串
func splitBySpace(s string) []string {
	var parts []string
	var current string
	for _, r := range s {
		if r == ' ' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// toLowerPlanetID 转换行星名称为小写ID
func toLowerPlanetID(name string) string {
	lower := toLower(name)
	planetMap := map[string]string{
		"sun": "sun", "moon": "moon", "mercury": "mercury",
		"venus": "venus", "mars": "mars", "jupiter": "jupiter",
		"saturn": "saturn", "uranus": "uranus", "neptune": "neptune",
		"pluto": "pluto", "chiron": "chiron",
		"north node": "northNode", "northnode": "northNode", "north": "northNode",
		"south node": "southNode", "southnode": "southNode", "south": "southNode",
		"ascendant": "ascendant", "midheaven": "midheaven",
	}
	if id, ok := planetMap[lower]; ok {
		return id
	}
	return lower
}

// toLowerAspect 转换相位名称为小写
func toLowerAspect(name string) string {
	lower := toLower(name)
	aspectMap := map[string]string{
		"conjunction": "conjunction", "trine": "trine",
		"square": "square", "opposition": "opposition",
		"sextile": "sextile", "quincunx": "quincunx",
	}
	if asp, ok := aspectMap[lower]; ok {
		return asp
	}
	return lower
}

// toLower 转换为小写
func toLower(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

// calculatePlanetaryHourEvent 计算行星时事件
func (c *Calculator) calculatePlanetaryHourEvent(t time.Time) *AstroEvent {
	// 使用现有的行星时计算
	hourInfo := astro.CalculatePlanetaryHourEnhanced(t, c.chart.BirthData.Latitude, c.chart.BirthData.Longitude)

	// 计算行星时的时间范围（简化：每小时）
	hourStart := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	hourEnd := hourStart.Add(time.Hour)

	eventID := GenerateEventID(EventTypePlanetaryHour, hourInfo.Ruler, "", "", t)

	// 根据 Influence 确定是否正面
	isPositive := hourInfo.Influence >= 0

	return &AstroEvent{
		EventID:       eventID,
		Type:          EventTypePlanetaryHour,
		Title:         string(hourInfo.Ruler) + " Hour",
		IsPositive:    isPositive,
		Intensity:     normalizeIntensity(hourInfo.Influence / 10), // 归一化
		TimeLevel:     TimeLevelHourly,
		PrimaryPlanet: string(hourInfo.Ruler),
		Impact:        c.getPlanetaryHourImpact(hourInfo.Ruler),
		StartTime:     hourStart,
		EndTime:       hourEnd,
		ExactTime:     hourStart,
	}
}

// calculateRetrogradeEvents 计算逆行事件
func (c *Calculator) calculateRetrogradeEvents(transitPositions []models.PlanetPosition, t time.Time) []AstroEvent {
	events := make([]AstroEvent, 0)

	for _, p := range transitPositions {
		if p.Retrograde {
			eventID := GenerateEventID(EventTypeRetrograde, p.ID, "", "", t)

			// 估计逆行周期（实际应从星历表获取）
			duration := c.getRetrogradeDuration(p.ID)
			startTime := t.Add(-duration / 2)
			endTime := t.Add(duration / 2)

			events = append(events, AstroEvent{
				EventID:       eventID,
				Type:          EventTypeRetrograde,
				Title:         string(p.ID) + " Retrograde",
				IsPositive:    false,
				Intensity:     0.7,
				TimeLevel:     c.getRetrogradeTimeLevel(p.ID),
				PrimaryPlanet: string(p.ID),
				Impact:        c.getRetrogradeImpact(p.ID),
				StartTime:     startTime,
				EndTime:       endTime,
				ExactTime:     t,
			})
		}
	}

	return events
}

// calculateVoidOfCourseEvent 计算月空事件
func (c *Calculator) calculateVoidOfCourseEvent(t time.Time) *AstroEvent {
	// 使用现有的月空计算 - 需要儒略日和本命盘位置
	jd := astro.TimeToJulianDay(t)
	vocInfo := astro.CalculateVoidOfCourse(jd, c.chart.Planets)
	if !vocInfo.IsVoid {
		return nil
	}

	eventID := GenerateEventID(EventTypeVoidOfCourse, models.Moon, "", "", t)

	// 解析开始和结束时间（VoidOfCourseInfo 返回字符串格式）
	// 简化处理：使用当前时间和预估结束时间
	hourStart := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	hourEnd := hourStart.Add(time.Duration(vocInfo.Duration) * time.Hour)

	return &AstroEvent{
		EventID:       eventID,
		Type:          EventTypeVoidOfCourse,
		Title:         "Moon Void of Course",
		IsPositive:    false,
		Intensity:     normalizeIntensity(-vocInfo.Influence / 15), // 归一化
		TimeLevel:     TimeLevelHourly,
		PrimaryPlanet: string(models.Moon),
		Impact: DimensionImpact{
			Career:       -2,
			Relationship: -1,
			Health:       0,
			Finance:      -2,
			Spiritual:    1,
		},
		StartTime: hourStart,
		EndTime:   hourEnd,
		ExactTime: t,
	}
}

// calculateLunarPhaseEvents 计算月相事件
func (c *Calculator) calculateLunarPhaseEvents(start, end time.Time) []AstroEvent {
	events := make([]AstroEvent, 0)

	// 获取当天的共享数据
	sharedData := astro.CalculateDailyAstroData(c.chart, start)
	if sharedData == nil {
		return events
	}

	// 遍历月相事件，仅保留落在 [start, end] 内的
	for _, lp := range sharedData.LunarPhases {
		if lp.ExactTime.Before(start) || !lp.ExactTime.Before(end) {
			continue
		}

		eventID := GenerateEventID(EventTypeLunarPhase, models.Moon, "", lp.Phase, lp.ExactTime)

		events = append(events, AstroEvent{
			EventID:       eventID,
			Type:          EventTypeLunarPhase,
			Title:         lp.PhaseName,
			IsPositive:    c.isLunarPhasePositive(lp.Phase),
			Intensity:     0.8, // 月相通常有较高影响
			TimeLevel:     TimeLevelDaily,
			PrimaryPlanet: string(models.Moon),
			Aspect:        lp.Phase, // i18n 会使用此字段 (new, first_quarter, full, last_quarter)
			Impact:        c.getLunarPhaseImpact(lp.Phase),
			StartTime:     lp.ExactTime.Add(-6 * time.Hour),
			EndTime:       lp.ExactTime.Add(6 * time.Hour),
			ExactTime:     lp.ExactTime,
		})
	}

	return events
}

// calculateSignChangeEvents 计算换座事件
func (c *Calculator) calculateSignChangeEvents(start, end time.Time) []AstroEvent {
	events := make([]AstroEvent, 0)

	// 获取当天的共享数据
	sharedData := astro.CalculateDailyAstroData(c.chart, start)
	if sharedData == nil {
		return events
	}

	// 遍历换座事件，仅保留落在 [start, end] 内的
	for _, sc := range sharedData.SignChanges {
		if sc.ExactTime.Before(start) || !sc.ExactTime.Before(end) {
			continue
		}

		newSign := strings.ToLower(sc.NewSign)
		eventID := GenerateEventID(EventTypeSignChange, sc.Planet, "", newSign, sc.ExactTime)

		// 计算行星离开新星座的时间作为 endTime
		endTime := c.calculateSignChangeEndTime(sc.Planet, sc.ExactTime)

		events = append(events, AstroEvent{
			EventID:       eventID,
			Type:          EventTypeSignChange,
			Title:         string(sc.Planet) + " enters " + sc.NewSign,
			IsPositive:    c.isSignChangePositive(sc.Planet, newSign),
			Intensity:     c.getSignChangeIntensity(sc.Planet),
			TimeLevel:     c.getSignChangeTimeLevel(sc.Planet),
			PrimaryPlanet: string(sc.Planet),
			Aspect:        newSign, // i18n 会使用此字段作为新星座
			Impact:        c.getSignChangeImpact(sc.Planet, newSign),
			StartTime:     sc.ExactTime,
			EndTime:       endTime,
			ExactTime:     sc.ExactTime,
		})
	}

	return events
}

// calculateSignChangeEndTime 计算换座事件的结束时间（行星离开新星座的时间）
func (c *Calculator) calculateSignChangeEndTime(planetID models.PlanetID, entryTime time.Time) time.Time {
	// 从进入时间开始，查找下一次换座
	nextChange := astro.FindNextSignChange(planetID, entryTime.Add(time.Hour), 0) // +1小时避免找到同一个换座点
	if nextChange != nil {
		return *nextChange
	}
	// 回退：使用估算值
	return entryTime.Add(c.getEstimatedSignDuration(planetID))
}

// calculateDignityEvents 计算尊贵度事件
func (c *Calculator) calculateDignityEvents(transitPositions []models.PlanetPosition, t time.Time) []AstroEvent {
	events := make([]AstroEvent, 0)

	for _, p := range transitPositions {
		dignity := astro.GetDignity(p.ID, p.Sign)

		// 仅处理非游离（peregrine）的尊贵度
		var dignityType string
		switch dignity {
		case models.DignityDomicile:
			dignityType = "domicile"
		case models.DignityExaltation:
			dignityType = "exaltation"
		case models.DignityDetriment:
			dignityType = "detriment"
		case models.DignityFall:
			dignityType = "fall"
		default:
			continue // 跳过 peregrine
		}

		// 计算精确的开始和结束时间
		startTime, endTime := c.calculateDignityTimeRange(p.ID, t)

		eventID := GenerateEventID(EventTypeDignity, p.ID, "", dignityType, startTime)

		events = append(events, AstroEvent{
			EventID:       eventID,
			Type:          EventTypeDignity,
			Title:         string(p.ID) + " in " + dignityType,
			IsPositive:    dignity == models.DignityDomicile || dignity == models.DignityExaltation,
			Intensity:     c.getDignityIntensity(dignity),
			TimeLevel:     c.getDignityTimeLevel(p.ID),
			PrimaryPlanet: string(p.ID),
			Aspect:        dignityType, // i18n 会使用此字段 (domicile, exaltation, detriment, fall)
			Impact:        c.getDignityImpact(dignity),
			StartTime:     startTime,
			EndTime:       endTime,
			ExactTime:     t, // 查询时间作为"当前时刻"
		})
	}

	return events
}

// calculateDignityTimeRange 计算行星在当前星座的精确时间范围
func (c *Calculator) calculateDignityTimeRange(planetID models.PlanetID, t time.Time) (startTime, endTime time.Time) {
	// 计算进入当前星座的时间
	prevChange := astro.FindPrevSignChange(planetID, t, 0)
	if prevChange != nil {
		startTime = *prevChange
	} else {
		// 回退：使用估算值
		startTime = t.Add(-c.getEstimatedSignDuration(planetID) / 2)
	}

	// 计算离开当前星座的时间
	nextChange := astro.FindNextSignChange(planetID, t, 0)
	if nextChange != nil {
		endTime = *nextChange
	} else {
		// 回退：使用估算值
		endTime = t.Add(c.getEstimatedSignDuration(planetID) / 2)
	}

	return startTime, endTime
}

// getEstimatedSignDuration 获取行星在星座停留的估算时间
func (c *Calculator) getEstimatedSignDuration(planetID models.PlanetID) time.Duration {
	switch planetID {
	case models.Moon:
		return 2*24*time.Hour + 12*time.Hour // ~2.5 天
	case models.Sun:
		return 30 * 24 * time.Hour // ~30 天
	case models.Mercury:
		return 21 * 24 * time.Hour // ~3 周（可能逆行延长）
	case models.Venus:
		return 25 * 24 * time.Hour // ~25 天
	case models.Mars:
		return 45 * 24 * time.Hour // ~45 天
	case models.Jupiter:
		return 365 * 24 * time.Hour // ~1 年
	case models.Saturn:
		return 2 * 365 * 24 * time.Hour // ~2.5 年
	default:
		return 7 * 365 * 24 * time.Hour // 外行星：多年
	}
}

// ==================== 月相辅助方法 ====================

// isLunarPhasePositive 判断月相是否正面
func (c *Calculator) isLunarPhasePositive(phase string) bool {
	switch phase {
	case "new", "full":
		return true
	default:
		return true // 月相通常中性偏正
	}
}

// getLunarPhaseImpact 获取月相的影响
func (c *Calculator) getLunarPhaseImpact(phase string) DimensionImpact {
	switch phase {
	case "new":
		return DimensionImpact{Career: 2, Relationship: 1, Health: 1, Finance: 1, Spiritual: 3}
	case "first_quarter":
		return DimensionImpact{Career: 2, Relationship: 0, Health: 1, Finance: 1, Spiritual: 1}
	case "full":
		return DimensionImpact{Career: 1, Relationship: 3, Health: 1, Finance: 1, Spiritual: 2}
	case "last_quarter":
		return DimensionImpact{Career: 0, Relationship: 1, Health: 1, Finance: 0, Spiritual: 2}
	default:
		return DimensionImpact{Career: 1, Relationship: 1, Health: 1, Finance: 1, Spiritual: 1}
	}
}

// ==================== 换座辅助方法 ====================

// isSignChangePositive 判断换座是否正面
func (c *Calculator) isSignChangePositive(planet models.PlanetID, sign string) bool {
	// 入庙或旺的星座为正面
	dignityMap := map[models.PlanetID][]string{
		models.Sun:     {"aries", "leo"},
		models.Moon:    {"cancer", "taurus"},
		models.Mercury: {"gemini", "virgo"},
		models.Venus:   {"taurus", "libra", "pisces"},
		models.Mars:    {"aries", "scorpio", "capricorn"},
		models.Jupiter: {"sagittarius", "pisces", "cancer"},
		models.Saturn:  {"capricorn", "aquarius", "libra"},
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

// getSignChangeIntensity 获取换座强度
func (c *Calculator) getSignChangeIntensity(planet models.PlanetID) float64 {
	switch planet {
	case models.Sun, models.Moon:
		return 0.8
	case models.Mercury, models.Venus, models.Mars:
		return 0.6
	default:
		return 0.5
	}
}

// getSignChangeTimeLevel 获取换座时间层级
func (c *Calculator) getSignChangeTimeLevel(planet models.PlanetID) string {
	switch planet {
	case models.Moon:
		return TimeLevelDaily
	case models.Sun, models.Mercury, models.Venus:
		return TimeLevelMonthly
	default:
		return TimeLevelMonthly
	}
}

// getSignChangeImpact 获取换座的影响
func (c *Calculator) getSignChangeImpact(planet models.PlanetID, sign string) DimensionImpact {
	// 基于行星和星座特性的简化影响
	baseImpact := c.getPlanetaryHourImpact(planet)
	if c.isSignChangePositive(planet, sign) {
		// 正面换座增强影响
		return DimensionImpact{
			Career:       baseImpact.Career * 1.5,
			Relationship: baseImpact.Relationship * 1.5,
			Health:       baseImpact.Health * 1.5,
			Finance:      baseImpact.Finance * 1.5,
			Spiritual:    baseImpact.Spiritual * 1.5,
		}
	}
	return baseImpact
}

// ==================== 尊贵度辅助方法 ====================

// getDignityIntensity 获取尊贵度强度
func (c *Calculator) getDignityIntensity(dignity models.Dignity) float64 {
	switch dignity {
	case models.DignityDomicile:
		return 0.9
	case models.DignityExaltation:
		return 0.85
	case models.DignityDetriment:
		return 0.7
	case models.DignityFall:
		return 0.75
	default:
		return 0.5
	}
}

// getDignityTimeLevel 获取尊贵度时间层级
func (c *Calculator) getDignityTimeLevel(planet models.PlanetID) string {
	switch planet {
	case models.Moon:
		return TimeLevelDaily
	case models.Sun, models.Mercury, models.Venus, models.Mars:
		return TimeLevelMonthly
	default:
		return TimeLevelYearly
	}
}

// getDignityImpact 获取尊贵度的影响
func (c *Calculator) getDignityImpact(dignity models.Dignity) DimensionImpact {
	switch dignity {
	case models.DignityDomicile:
		return DimensionImpact{Career: 3, Relationship: 2, Health: 2, Finance: 2, Spiritual: 2}
	case models.DignityExaltation:
		return DimensionImpact{Career: 2, Relationship: 2, Health: 2, Finance: 2, Spiritual: 3}
	case models.DignityDetriment:
		return DimensionImpact{Career: -2, Relationship: -1, Health: -1, Finance: -1, Spiritual: 0}
	case models.DignityFall:
		return DimensionImpact{Career: -1, Relationship: -2, Health: -1, Finance: -1, Spiritual: 0}
	default:
		return DimensionImpact{}
	}
}

// ==================== 辅助方法 ====================

// getUserID 获取用户 ID（基于本命盘生成）
func (c *Calculator) getUserID() string {
	birthTime := c.chart.BirthData.ToTime()
	return birthTime.Format("20060102_1504") + "_" + formatCoord(c.chart.BirthData.Latitude, c.chart.BirthData.Longitude)
}

// formatCoord 格式化坐标
func formatCoord(lat, lng float64) string {
	return formatFloat(lat) + "_" + formatFloat(lng)
}

// formatFloat 格式化浮点数
func formatFloat(f float64) string {
	if f < 0 {
		return "n" + formatPositiveFloat(-f)
	}
	return formatPositiveFloat(f)
}

// formatPositiveFloat 格式化正浮点数
func formatPositiveFloat(f float64) string {
	intPart := int(f)
	decPart := int((f - float64(intPart)) * 100)
	return itoa(intPart) + "p" + itoa(decPart)
}

// itoa 整数转字符串
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}

// getTimeLevel 获取时间层级字符串
func (c *Calculator) getTimeLevel(level models.FactorTimeLevel) string {
	switch level {
	case models.TimeLevelHourly:
		return TimeLevelHourly
	case models.TimeLevelDaily:
		return TimeLevelDaily
	case models.TimeLevelWeekly:
		return TimeLevelWeekly
	case models.TimeLevelMonthly:
		return TimeLevelMonthly
	case models.TimeLevelYearly:
		return TimeLevelYearly
	default:
		return TimeLevelDaily
	}
}

// normalizeIntensity 归一化强度到 0-1
func normalizeIntensity(strength float64) float64 {
	if strength < 0 {
		strength = -strength
	}
	if strength > 1 {
		strength = 1
	}
	return strength
}

// factorToImpact 将因子转换为影响
func (c *Calculator) factorToImpact(f *models.InfluenceFactor) DimensionImpact {
	return DimensionImpact{
		Career:       f.DimensionImpact.Career,
		Relationship: f.DimensionImpact.Relationship,
		Health:       f.DimensionImpact.Health,
		Finance:      f.DimensionImpact.Finance,
		Spiritual:    f.DimensionImpact.Spiritual,
	}
}

// getPlanetaryHourImpact 获取行星时的影响
func (c *Calculator) getPlanetaryHourImpact(ruler models.PlanetID) DimensionImpact {
	switch ruler {
	case models.Sun:
		return DimensionImpact{Career: 2, Relationship: 0, Health: 1, Finance: 1, Spiritual: 1}
	case models.Moon:
		return DimensionImpact{Career: 0, Relationship: 2, Health: 1, Finance: 0, Spiritual: 2}
	case models.Mercury:
		return DimensionImpact{Career: 2, Relationship: 1, Health: 0, Finance: 1, Spiritual: 1}
	case models.Venus:
		return DimensionImpact{Career: 1, Relationship: 3, Health: 1, Finance: 2, Spiritual: 1}
	case models.Mars:
		return DimensionImpact{Career: 2, Relationship: -1, Health: 2, Finance: 1, Spiritual: 0}
	case models.Jupiter:
		return DimensionImpact{Career: 2, Relationship: 1, Health: 1, Finance: 3, Spiritual: 2}
	case models.Saturn:
		return DimensionImpact{Career: 1, Relationship: -1, Health: 0, Finance: 1, Spiritual: 1}
	default:
		return DimensionImpact{}
	}
}

// getRetrogradeImpact 获取逆行的影响
func (c *Calculator) getRetrogradeImpact(planet models.PlanetID) DimensionImpact {
	switch planet {
	case models.Mercury:
		return DimensionImpact{Career: -3, Relationship: -1, Health: 0, Finance: -2, Spiritual: 0}
	case models.Venus:
		return DimensionImpact{Career: -1, Relationship: -3, Health: 0, Finance: -2, Spiritual: 0}
	case models.Mars:
		return DimensionImpact{Career: -2, Relationship: -1, Health: -1, Finance: -1, Spiritual: 0}
	case models.Jupiter:
		return DimensionImpact{Career: -1, Relationship: 0, Health: 0, Finance: -1, Spiritual: 1}
	case models.Saturn:
		return DimensionImpact{Career: -1, Relationship: 0, Health: 0, Finance: 0, Spiritual: 1}
	case models.Uranus:
		return DimensionImpact{Career: -1, Relationship: -1, Health: 0, Finance: -1, Spiritual: 2}
	case models.Neptune:
		return DimensionImpact{Career: -1, Relationship: -1, Health: -1, Finance: -1, Spiritual: 2}
	case models.Pluto:
		return DimensionImpact{Career: -1, Relationship: -1, Health: -1, Finance: -1, Spiritual: 3}
	case models.NorthNode:
		return DimensionImpact{Career: 0, Relationship: 1, Health: 0, Finance: 0, Spiritual: 2}
	case models.Chiron:
		return DimensionImpact{Career: 0, Relationship: 0, Health: 1, Finance: 0, Spiritual: 2}
	default:
		return DimensionImpact{Career: 0, Relationship: 0, Health: 0, Finance: 0, Spiritual: 0}
	}
}

// getRetrogradeDuration 获取逆行持续时间
func (c *Calculator) getRetrogradeDuration(planet models.PlanetID) time.Duration {
	switch planet {
	case models.Mercury:
		return 21 * 24 * time.Hour // 约 3 周
	case models.Venus:
		return 40 * 24 * time.Hour // 约 40 天
	case models.Mars:
		return 70 * 24 * time.Hour // 约 70 天
	case models.Jupiter:
		return 120 * 24 * time.Hour // 约 4 个月
	case models.Saturn:
		return 140 * 24 * time.Hour // 约 4.5 个月
	default:
		return 30 * 24 * time.Hour
	}
}

// getRetrogradeTimeLevel 获取逆行的时间层级
func (c *Calculator) getRetrogradeTimeLevel(planet models.PlanetID) string {
	switch planet {
	case models.Mercury:
		return TimeLevelWeekly
	case models.Venus, models.Mars:
		return TimeLevelMonthly
	default:
		return TimeLevelMonthly
	}
}

// applyI18n 应用国际化
func (c *Calculator) applyI18n(event *AstroEvent, translator *i18n.Translator) {
	// 翻译行星名称
	event.PrimaryPlanetName = translator.GetPlanetName(models.PlanetID(event.PrimaryPlanet))
	if event.SecondaryPlanet != "" {
		event.SecondaryPlanetName = translator.GetPlanetName(models.PlanetID(event.SecondaryPlanet))
	}

	// 翻译相位/尊贵度名称（根据事件类型决定）
	if event.Aspect != "" {
		switch event.Type {
		case EventTypeDignity:
			// dignity 事件的 Aspect 字段存储 dignity 类型
			event.AspectName = translator.T("dignity." + event.Aspect)
		case EventTypeLunarPhase:
			// 月相事件的 Aspect 字段存储月相类型
			event.AspectName = translator.T("lunar_phase." + event.Aspect)
		default:
			// 其他事件使用相位翻译
			event.AspectName = translator.GetAspectName(event.Aspect)
		}
	}

	// 获取详细解读
	event.Interpretation = translator.GetDetailedInterpretation(
		event.Type,
		models.PlanetID(event.PrimaryPlanet),
		models.PlanetID(event.SecondaryPlanet),
		event.Aspect,
		"", // house - 相位事件通常无宫位
		event.IsPositive,
	)

	// 生成标题
	event.Title = c.generateEventTitle(event, translator)
}

// generateEventTitle 生成事件标题
func (c *Calculator) generateEventTitle(event *AstroEvent, translator *i18n.Translator) string {
	switch event.Type {
	case EventTypeAspect:
		return event.PrimaryPlanetName + " " + event.AspectName + " " + event.SecondaryPlanetName
	case EventTypePlanetaryHour:
		return event.PrimaryPlanetName + translator.T("planetary_hour_suffix")
	case EventTypeRetrograde:
		return event.PrimaryPlanetName + translator.T("retrograde_suffix")
	case EventTypeVoidOfCourse:
		return translator.T("void_of_course")
	case EventTypeLunarPhase:
		// 使用 i18n 的情感标题
		title := translator.GetEmotionalTitle(
			event.Type,
			models.PlanetID(event.PrimaryPlanet),
			"",
			event.Aspect, // phase key
			"",
			event.IsPositive,
		)
		if title != "" {
			return title
		}
		return event.Title
	case EventTypeSignChange:
		// 使用 i18n 的情感标题
		title := translator.GetEmotionalTitle(
			event.Type,
			models.PlanetID(event.PrimaryPlanet),
			"",
			event.Aspect, // new sign
			"",
			event.IsPositive,
		)
		if title != "" {
			return title
		}
		return event.PrimaryPlanetName + " " + translator.T("sign_change_suffix")
	case EventTypeDignity:
		// 使用 i18n 的情感标题
		title := translator.GetEmotionalTitle(
			event.Type,
			models.PlanetID(event.PrimaryPlanet),
			"",
			event.Aspect, // dignity type
			"",
			event.IsPositive,
		)
		if title != "" {
			return title
		}
		return event.PrimaryPlanetName + " " + event.Aspect
	default:
		return event.Title
	}
}

// generateGuidance 生成指导
func (c *Calculator) generateGuidance(slot *TimeSlot) *Guidance {
	if len(slot.Events) == 0 {
		return nil
	}

	// 找出最强的正面和负面事件
	var bestPositive, worstNegative *AstroEvent
	for i := range slot.Events {
		e := &slot.Events[i]
		if e.IsPositive && (bestPositive == nil || e.Intensity > bestPositive.Intensity) {
			bestPositive = e
		}
		if !e.IsPositive && (worstNegative == nil || e.Intensity > worstNegative.Intensity) {
			worstNegative = e
		}
	}

	guidance := &Guidance{
		Dos:   make([]string, 0),
		Donts: make([]string, 0),
	}

	// 生成摘要
	if bestPositive != nil && worstNegative != nil {
		guidance.Summary = bestPositive.Title + " brings opportunities, but watch out for " + worstNegative.Title
	} else if bestPositive != nil {
		guidance.Summary = "Good time for " + bestPositive.Title
	} else if worstNegative != nil {
		guidance.Summary = "Be cautious of " + worstNegative.Title
	}

	// 确定重点关注维度
	guidance.Focus = c.determineFocusDimension(slot)

	return guidance
}

// determineFocusDimension 确定重点关注维度
func (c *Calculator) determineFocusDimension(slot *TimeSlot) string {
	scores := slot.Scores
	dims := map[string]float64{
		"career":       scores.Career,
		"relationship": scores.Relationship,
		"health":       scores.Health,
		"finance":      scores.Finance,
		"spiritual":    scores.Spiritual,
	}

	// 找出最高分的维度
	maxDim := "career"
	maxScore := scores.Career
	for dim, score := range dims {
		if score > maxScore {
			maxScore = score
			maxDim = dim
		}
	}
	return maxDim
}
