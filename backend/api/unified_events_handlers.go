package api

import (
	"fmt"
	"net/http"
	"sort"
	"star/astro"
	"star/i18n"
	"star/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== 统一事件接口 ====================
// 合并 daily-events 和 total-factors 功能
// 返回时间范围内的所有天体事件及其对分数的影响

// UnifiedEventRequest 统一事件请求
type UnifiedEventRequest struct {
	BirthData           models.BirthData `json:"birthData" binding:"required"`
	Birth               models.BirthData `json:"birth"`               // 别名
	StartTime           string           `json:"startTime"`           // ISO 8601 格式
	EndTime             string           `json:"endTime"`             // ISO 8601 格式
	Date                string           `json:"date"`                // 单日查询（简化格式）
	Timezone            int              `json:"timezone"`            // 时区偏移（小时）
	Granularity         string           `json:"granularity"`         // hour, day, week, month, year（过滤因子用）
	Language            string           `json:"language"`            // zh, en, ru (default: en)
	IncludeMinorAspects bool             `json:"includeMinorAspects"` // 是否包含次要相位
	IncludeFactors      bool             `json:"includeFactors"`      // 是否包含因子影响数据（默认true）
	IncludeTransitHouse bool             `json:"includeTransitHouse"` // 是否包含行运过宫
	IncludeProgressions bool             `json:"includeProgressions"` // 是否包含次限/三限推进
}

// UnifiedEvent 统一事件（包含事件和因子信息）
type UnifiedEvent struct {
	// 事件基本信息
	Time        time.Time `json:"time"`
	Type        string    `json:"type"` // aspect, sign_change, lunar_phase, planetary_hour_change, retrograde, dignity, voc, transit_house, secondary_progression, tertiary_progression
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Theme       string    `json:"theme"`
	Advice      string    `json:"advice"`
	IsPositive  bool      `json:"isPositive"`
	Intensity   string    `json:"intensity"` // high, medium, low

	// 事件详情（可选）
	Planet1 models.PlanetID `json:"planet1,omitempty"`
	Planet2 models.PlanetID `json:"planet2,omitempty"`
	Aspect  string          `json:"aspect,omitempty"`
	Sign    string          `json:"sign,omitempty"`
	Degree  float64         `json:"degree,omitempty"`
	House   int             `json:"house,omitempty"` // For transit house events

	// User trust fields
	IsExactToday   bool   `json:"isExactToday"`             // Whether event is exact today
	InfluencePhase string `json:"influencePhase,omitempty"` // approaching, active, fading
	
	// Duration info (for long-term events)
	StartDate    string  `json:"startDate,omitempty"`    // When influence started
	EndDate      string  `json:"endDate,omitempty"`      // When influence ends
	DurationDays float64 `json:"durationDays,omitempty"` // Duration in days
	DurationText string  `json:"durationText,omitempty"` // Formatted duration text
	
	// Localized fields
	EmotionalTitle        string   `json:"emotionalTitle,omitempty"`        // Emotional/psychological title
	DetailedInterpretation string   `json:"detailedInterpretation,omitempty"` // Detailed paragraph interpretation
	DimensionLabels       []string `json:"dimensionLabels,omitempty"`       // Dimension labels with arrows
	
	// Factor influence data
	Factor *FactorImpact `json:"factor,omitempty"`
}

// FactorImpact 因子影响数据
type FactorImpact struct {
	FactorType      string                  `json:"factorType"`
	TimeLevel       string                  `json:"timeLevel"`       // hourly, daily, weekly, monthly, yearly
	BaseValue       float64                 `json:"baseValue"`
	Weight          float64                 `json:"weight"`
	Strength        float64                 `json:"strength"`        // 当前强度 (0-1)
	DimensionImpact models.DimensionImpact  `json:"dimensionImpact"` // 对各维度的影响
	Lifecycle       *FactorLifecycleInfo    `json:"lifecycle,omitempty"`
}

// FactorLifecycleInfo 因子生命周期信息
type FactorLifecycleInfo struct {
	StartTime time.Time `json:"startTime"`
	PeakTime  time.Time `json:"peakTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  float64   `json:"durationHours"`
	Phase     string    `json:"phase"` // applying, exact, separating
}

// UnifiedEventResponse 统一事件响应
type UnifiedEventResponse struct {
	StartTime   string         `json:"startTime"`
	EndTime     string         `json:"endTime"`
	Timezone    int            `json:"timezone"`
	Events      []UnifiedEvent `json:"events"`
	EventCount  int            `json:"eventCount"`
	MajorEvents []UnifiedEvent `json:"majorEvents"` // 高强度事件
	Summary     string         `json:"summary"`
	DayTheme    string         `json:"dayTheme"`

	// 因子汇总
	FactorSummary *FactorSummary `json:"factorSummary,omitempty"`

	// 按时间级别分组的事件
	EventsByLevel map[string][]UnifiedEvent `json:"eventsByLevel,omitempty"`
}

// FactorSummary 因子汇总
type FactorSummary struct {
	TotalFactors    int     `json:"totalFactors"`
	PositiveFactors int     `json:"positiveFactors"`
	NegativeFactors int     `json:"negativeFactors"`
	NetInfluence    float64 `json:"netInfluence"`
	DominantFactor  string  `json:"dominantFactor,omitempty"`
}

// GetUnifiedEvents 统一事件查询
// POST /api/calc/unified-events
// 合并 daily-events 和 total-factors，返回时间范围内的所有天体事件及其对分数的影响
func GetUnifiedEvents(c *gin.Context) {
	var req UnifiedEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	// 支持 birth 别名
	if req.BirthData.Year == 0 && req.Birth.Year != 0 {
		req.BirthData = req.Birth
	}

	// 设置默认值
	if req.Granularity == "" {
		req.Granularity = "day"
	}
	if req.Language == "" {
		req.Language = "en"
	}
	
	// Create translator
	translator := i18n.New(req.Language)

	// 解析时间范围
	var startTime, endTime time.Time
	var err error

	if req.Date != "" {
		// 单日查询模式
		startTime, err = parseDate(req.Date, req.Timezone)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}
		endTime = startTime.Add(24 * time.Hour)
	} else if req.StartTime != "" && req.EndTime != "" {
		// 时间范围查询模式
		startTime, err = parseDateTime(req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startTime format: " + err.Error()})
			return
		}
		endTime, err = parseDateTime(req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endTime format: " + err.Error()})
			return
		}
	} else {
		// 默认今天
		location := time.FixedZone("Custom", req.Timezone*3600)
		now := time.Now().In(location)
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
		endTime = startTime.Add(24 * time.Hour)
	}

	// 计算本命盘
	chart := astro.CalculateNatalChart(req.BirthData)

	// 获取共享星象数据
	sharedData := astro.CalculateDailyAstroData(chart, startTime)

	// 获取每日事件
	dailyEvents := astro.CalculateDailyEvents(chart, startTime, req.IncludeMinorAspects)

	// 获取全因子数据
	var factorResponse *astro.TotalFactorsResponse
	factorResponse_val := astro.GetTotalFactors(chart, startTime.Add(12*time.Hour), req.Granularity, "")
	factorResponse = &factorResponse_val

	// 合并事件和因子
	unifiedEvents := mergeEventsAndFactors(dailyEvents, factorResponse, sharedData)
	
	// Add transit house events (planet through natal houses)
	if req.IncludeTransitHouse || req.Granularity == "day" || req.Granularity == "week" {
		transitHouseEvents := astro.GetActiveTransitHouseEvents(chart, startTime.Add(12*time.Hour), "medium")
		for _, e := range transitHouseEvents {
			unifiedEvents = append(unifiedEvents, convertTransitHouseToUnifiedEvent(e, translator, startTime))
		}
	}
	
	// Add progression events based on granularity
	if req.IncludeProgressions || req.Granularity == "month" || req.Granularity == "year" {
		// Secondary progressions (for yearly view)
		if req.Granularity == "year" || req.IncludeProgressions {
			spEvents := astro.GetSecondaryProgressionEvents(chart, startTime.Add(12*time.Hour))
			for _, e := range spEvents {
				unifiedEvents = append(unifiedEvents, convertProgressionToUnifiedEvent(e, translator, startTime))
			}
		}
		
		// Tertiary progressions (for monthly view)
		if req.Granularity == "month" || req.IncludeProgressions {
			tpEvents := astro.GetTertiaryProgressionEvents(chart, startTime.Add(12*time.Hour))
			for _, e := range tpEvents {
				unifiedEvents = append(unifiedEvents, convertProgressionToUnifiedEvent(e, translator, startTime))
			}
		}
	}
	
	// Apply translations and emotional titles to all events
	for i := range unifiedEvents {
		applyTranslationsToEvent(&unifiedEvents[i], translator, startTime)
	}

	// 筛选主要事件
	majorEvents := []UnifiedEvent{}
	for _, event := range unifiedEvents {
		if event.Intensity == "high" {
			majorEvents = append(majorEvents, event)
		}
	}

	// 生成主题和总结
	dayTheme := generateUnifiedDayTheme(unifiedEvents)
	summary := generateUnifiedSummary(unifiedEvents, majorEvents)

	// 因子汇总
	var factorSummary *FactorSummary
	if factorResponse != nil {
		factorSummary = &FactorSummary{
			TotalFactors:    factorResponse.Overall.PositiveCount + factorResponse.Overall.NegativeCount,
			PositiveFactors: factorResponse.Overall.PositiveCount,
			NegativeFactors: factorResponse.Overall.NegativeCount,
			NetInfluence:    factorResponse.Overall.NetAdjustment,
			DominantFactor:  getDominantFactorName(factorResponse),
		}
	}

	// 按时间级别分组事件
	eventsByLevel := groupEventsByTimeLevel(unifiedEvents)

	response := UnifiedEventResponse{
		StartTime:     startTime.Format(time.RFC3339),
		EndTime:       endTime.Format(time.RFC3339),
		Timezone:      req.Timezone,
		Events:        unifiedEvents,
		EventCount:    len(unifiedEvents),
		MajorEvents:   majorEvents,
		Summary:       summary,
		DayTheme:      dayTheme,
		FactorSummary: factorSummary,
		EventsByLevel: eventsByLevel,
	}

	c.JSON(http.StatusOK, response)
}

// groupEventsByTimeLevel 按时间级别分组事件
func groupEventsByTimeLevel(events []UnifiedEvent) map[string][]UnifiedEvent {
	grouped := map[string][]UnifiedEvent{
		"yearly":  {},
		"monthly": {},
		"weekly":  {},
		"daily":   {},
		"hourly":  {},
		"unknown": {}, // 没有因子数据的事件
	}

	for _, e := range events {
		if e.Factor == nil || e.Factor.TimeLevel == "" {
			grouped["unknown"] = append(grouped["unknown"], e)
		} else {
			level := e.Factor.TimeLevel
			if _, ok := grouped[level]; ok {
				grouped[level] = append(grouped[level], e)
			} else {
				grouped["unknown"] = append(grouped["unknown"], e)
			}
		}
	}

	// 删除空的分组
	for k, v := range grouped {
		if len(v) == 0 {
			delete(grouped, k)
		}
	}

	return grouped
}

// mergeEventsAndFactors 合并事件和因子数据
// 策略：以 Factor 系统为主（因为这才是影响分数的），补充 Daily Events 的精确时间
func mergeEventsAndFactors(events []astro.DailyEvent, factorResponse *astro.TotalFactorsResponse, sharedData *astro.DailyAstroData) []UnifiedEvent {
	unifiedEvents := []UnifiedEvent{}
	usedFactorIDs := make(map[string]bool) // 记录已处理的因子

	// 收集所有因子（从各级别合并）
	allFactors := []astro.TotalFactorDetail{}
	if factorResponse != nil {
		for _, factors := range factorResponse.FactorsByLevel {
			allFactors = append(allFactors, factors...)
		}
	}

	// 获取今天的日期范围（用于判断 isExactToday）
	today := time.Now()
	startOfToday := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfToday := startOfToday.Add(24 * time.Hour)

	// 步骤1：处理 Daily Events（有精确时间的事件）
	for _, e := range events {
		ue := UnifiedEvent{
			Time:         e.Time,
			Type:         e.Type,
			Title:        e.Title,
			Description:  e.Description,
			Theme:        e.Theme,
			Advice:       e.Advice,
			IsPositive:   e.IsPositive,
			Intensity:    e.Intensity,
			Planet1:      e.Planet1,
			Planet2:      e.Planet2,
			Aspect:       e.Aspect,
			Sign:         e.Sign,
			Degree:       e.Degree,
			IsExactToday: true, // Daily Events 都是今日精确的
		}

		// 尝试匹配因子数据
		factor, factorID := findMatchingFactorDetailWithID(e, allFactors)
		if factor != nil {
			ue.Factor = factor
			ue.InfluencePhase = getInfluencePhase(factor.Lifecycle, e.Time)
			usedFactorIDs[factorID] = true
		} else {
			// 没有匹配到因子时，根据事件类型和涉及的行星分配精确的 timeLevel 和 DimensionImpact
			factorType := eventTypeToFactorType(e.Type)
			timeLevel := getSmartTimeLevel(e.Type, e.Planet1, e.Planet2)
			dimensionImpact := models.DimensionImpact{}
			
			// 为相位事件计算维度影响
			if e.Type == "aspect" && e.Planet1 != "" && e.Planet2 != "" {
				dimensionImpact = getAspectDimensionImpactAPI(e.Planet1, e.Planet2)
			}
			
			ue.Factor = &FactorImpact{
				FactorType:      factorType,
				TimeLevel:       timeLevel,
				DimensionImpact: dimensionImpact,
			}
			ue.InfluencePhase = "active"
		}

		unifiedEvents = append(unifiedEvents, ue)
	}

	// 步骤2：添加 Factor 系统中有但 Daily Events 中没有的因子
	// 这些是"正在影响但非今日精确"的因子
	for _, f := range allFactors {
		factorID := f.ID
		if usedFactorIDs[factorID] {
			continue // 已经在步骤1中处理过
		}

		ue := convertFactorDetailToUnifiedEvent(f)
		
		// 判断是否今日精确
		if f.PeakTime != "" {
			if pt, err := time.Parse(time.RFC3339, f.PeakTime); err == nil {
				ue.IsExactToday = pt.After(startOfToday) && pt.Before(endOfToday)
			}
		}
		
		// 设置影响阶段
		ue.InfluencePhase = getInfluencePhaseFromStrings(f.StartTime, f.PeakTime, f.EndTime)

		unifiedEvents = append(unifiedEvents, ue)
	}

	// 按时间排序
	sort.Slice(unifiedEvents, func(i, j int) bool {
		return unifiedEvents[i].Time.Before(unifiedEvents[j].Time)
	})

	return unifiedEvents
}

// findMatchingFactorDetailWithID 查找匹配的因子并返回ID
func findMatchingFactorDetailWithID(event astro.DailyEvent, factors []astro.TotalFactorDetail) (*FactorImpact, string) {
	for _, f := range factors {
		if matchEventToFactorDetail(event, f) {
			return &FactorImpact{
				FactorType: f.Type,
				TimeLevel:  f.TimeLevel,
				BaseValue:  f.BaseValue,
				Weight:     f.Weight,
				Strength:   f.Strength,
				DimensionImpact: models.DimensionImpact{
					Career:       f.DimensionImpact.Career,
					Relationship: f.DimensionImpact.Relationship,
					Health:       f.DimensionImpact.Health,
					Finance:      f.DimensionImpact.Finance,
					Spiritual:    f.DimensionImpact.Spiritual,
				},
				Lifecycle: parseLifecycleFromFactor(f),
			}, f.ID
		}
	}
	return nil, ""
}

// getInfluencePhase 根据生命周期确定影响阶段
func getInfluencePhase(lifecycle *FactorLifecycleInfo, eventTime time.Time) string {
	if lifecycle == nil {
		return "active"
	}
	
	now := time.Now()
	if now.Before(lifecycle.PeakTime) {
		return "approaching" // 即将到来
	} else if now.After(lifecycle.EndTime) {
		return "fading" // 逐渐消退
	}
	return "active" // 正在影响
}

// getInfluencePhaseFromStrings 从字符串时间判断影响阶段
func getInfluencePhaseFromStrings(startTime, peakTime, endTime string) string {
	now := time.Now()
	
	if peakTime != "" {
		if pt, err := time.Parse(time.RFC3339, peakTime); err == nil {
			if now.Before(pt) {
				return "approaching"
			}
		}
	}
	
	if endTime != "" {
		if et, err := time.Parse(time.RFC3339, endTime); err == nil {
			if now.After(et) {
				return "fading"
			}
		}
	}
	
	return "active"
}

// matchEventToFactorDetail 匹配事件和因子
func matchEventToFactorDetail(event astro.DailyEvent, factor astro.TotalFactorDetail) bool {
	switch event.Type {
	case "aspect":
		if factor.Type == "aspectPhase" {
			// 检查行星是否匹配
			return factor.SourcePlanet == string(event.Planet1)
		}
	case "lunar_phase":
		return factor.Type == "lunarPhase"
	case "planetary_hour_change":
		return factor.Type == "planetaryHour"
	case "sign_change":
		return factor.Type == "dignity" && factor.SourcePlanet == string(event.Planet1)
	}
	return false
}

// isFactorDetailAlreadyInEvents 检查因子是否已经作为事件存在
func isFactorDetailAlreadyInEvents(factor astro.TotalFactorDetail, events []astro.DailyEvent) bool {
	for _, e := range events {
		if matchEventToFactorDetail(e, factor) {
			return true
		}
	}
	return false
}

// convertFactorDetailToUnifiedEvent 将因子转换为统一事件
func convertFactorDetailToUnifiedEvent(f astro.TotalFactorDetail) UnifiedEvent {
	eventType := f.Type
	intensity := "medium"
	if f.Strength > 0.7 {
		intensity = "high"
	} else if f.Strength < 0.3 {
		intensity = "low"
	}

	peakTime := time.Now()
	if f.PeakTime != "" {
		if pt, err := time.Parse(time.RFC3339, f.PeakTime); err == nil {
			peakTime = pt
		}
	}

	return UnifiedEvent{
		Time:        peakTime,
		Type:        eventType,
		Title:       f.Name,
		Description: f.Description,
		Theme:       getFactorThemeByType(f.Type),
		Advice:      getFactorAdviceByPositive(f.IsPositive),
		IsPositive:  f.IsPositive,
		Intensity:   intensity,
		Planet1:     models.PlanetID(f.SourcePlanet),
		Factor: &FactorImpact{
			FactorType: f.Type,
			TimeLevel:  f.TimeLevel,
			BaseValue:  f.BaseValue,
			Weight:     f.Weight,
			Strength:   f.Strength,
			DimensionImpact: models.DimensionImpact{
				Career:       f.DimensionImpact.Career,
				Relationship: f.DimensionImpact.Relationship,
				Health:       f.DimensionImpact.Health,
				Finance:      f.DimensionImpact.Finance,
				Spiritual:    f.DimensionImpact.Spiritual,
			},
			Lifecycle: parseLifecycleFromFactor(f),
		},
	}
}

// parseLifecycleFromFactor 从因子解析生命周期数据
func parseLifecycleFromFactor(f astro.TotalFactorDetail) *FactorLifecycleInfo {
	if f.StartTime == "" && f.EndTime == "" {
		return nil
	}

	lc := &FactorLifecycleInfo{}

	if f.StartTime != "" {
		if t, err := time.Parse(time.RFC3339, f.StartTime); err == nil {
			lc.StartTime = t
		}
	}
	if f.PeakTime != "" {
		if t, err := time.Parse(time.RFC3339, f.PeakTime); err == nil {
			lc.PeakTime = t
		}
	}
	if f.EndTime != "" {
		if t, err := time.Parse(time.RFC3339, f.EndTime); err == nil {
			lc.EndTime = t
		}
	}

	lc.Duration = f.RemainingDays * 24 // 转为小时

	// 判断阶段
	now := time.Now()
	if now.Before(lc.PeakTime) {
		lc.Phase = "applying"
	} else if now.After(lc.PeakTime) {
		lc.Phase = "separating"
	} else {
		lc.Phase = "exact"
	}

	return lc
}

// getFactorThemeByType 根据因子类型获取主题
func getFactorThemeByType(factorType string) string {
	themes := map[string]string{
		"dignity":        "Planetary strength and expression",
		"retrograde":     "Review, reflection, revisiting",
		"voidOfCourse":   "Pause period, avoid new starts",
		"profectionLord": "Annual themes and focus",
		"aspectPhase":    "Planetary interaction and energy",
		"lunarPhase":     "Emotional and intuitive cycles",
		"planetaryHour":  "Current planetary influence",
	}
	if theme, ok := themes[factorType]; ok {
		return theme
	}
	return "Energy influence"
}

// getFactorAdviceByPositive 根据正负获取建议
func getFactorAdviceByPositive(isPositive bool) string {
	if isPositive {
		return "Favorable period for related activities"
	}
	return "Be mindful and proceed with awareness"
}

// getDominantFactorName 获取主导因子名称
func getDominantFactorName(result *astro.TotalFactorsResponse) string {
	if result == nil {
		return ""
	}

	// 从所有正负因子中找最强的
	var dominantName string
	maxImpact := 0.0

	for _, f := range result.Overall.PositiveFactors {
		impact := f.Adjustment
		if impact < 0 {
			impact = -impact
		}
		if impact > maxImpact {
			maxImpact = impact
			dominantName = f.Name
		}
	}
	for _, f := range result.Overall.NegativeFactors {
		impact := f.Adjustment
		if impact < 0 {
			impact = -impact
		}
		if impact > maxImpact {
			maxImpact = impact
			dominantName = f.Name
		}
	}

	return dominantName
}

// generateUnifiedDayTheme 生成统一日主题
func generateUnifiedDayTheme(events []UnifiedEvent) string {
	// 统计事件类型
	aspectCount := 0
	positiveCount := 0
	highIntensity := 0

	for _, e := range events {
		if e.Type == "aspect" {
			aspectCount++
		}
		if e.IsPositive {
			positiveCount++
		}
		if e.Intensity == "high" {
			highIntensity++
		}
	}

	if highIntensity > 3 {
		return "Dynamic energy day with major cosmic influences"
	} else if positiveCount > len(events)/2 {
		return "Harmonious flow with supportive cosmic energies"
	} else if aspectCount > 5 {
		return "Active day with multiple planetary connections"
	}
	return "Balanced cosmic influences"
}

// generateUnifiedSummary 生成统一总结
func generateUnifiedSummary(events []UnifiedEvent, majorEvents []UnifiedEvent) string {
	if len(majorEvents) == 0 {
		return "A relatively calm day with subtle cosmic influences."
	}

	summary := ""
	if len(majorEvents) == 1 {
		summary = "Today features " + majorEvents[0].Title + ". "
	} else {
		summary = "Today features multiple significant events. "
	}

	positiveCount := 0
	for _, e := range majorEvents {
		if e.IsPositive {
			positiveCount++
		}
	}

	if positiveCount > len(majorEvents)/2 {
		summary += "Overall energy is supportive and favorable."
	} else {
		summary += "Navigate challenges with awareness and flexibility."
	}

	return summary
}

// parseDate 解析日期
func parseDate(dateStr string, timezone int) (time.Time, error) {
	location := time.FixedZone("Custom", timezone*3600)

	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"20060102",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, dateStr, location); err == nil {
			return t, nil
		}
	}

	return time.Time{}, nil
}

// parseDateTime 解析日期时间
func parseDateTime(dateTimeStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateTimeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, nil
}

// eventTypeToFactorType 将事件类型映射到因子类型
func eventTypeToFactorType(eventType string) string {
	mapping := map[string]string{
		"aspect":                "aspectPhase",
		"sign_change":           "dignity",
		"lunar_phase":           "lunarPhase",
		"planetary_hour_change": "planetaryHour",
	}
	if ft, ok := mapping[eventType]; ok {
		return ft
	}
	return eventType
}

// getDefaultTimeLevelForEventType 根据事件类型获取默认的时间级别
func getDefaultTimeLevelForEventType(eventType string) string {
	// 基于 FactorTimeLevelMapping 的逻辑
	mapping := map[string]string{
		"aspect":                 "daily",   // 相位默认日级别（实际由 getAspectTimeLevel 细分）
		"sign_change":            "monthly", // 换座默认月级别
		"lunar_phase":            "daily",   // 月相默认日级别
		"planetary_hour_change":  "hourly",  // 行星时默认小时级别
		"transit_house":          "weekly",  // 行运过宫默认周级别（实际由行星细分）
		"secondary_progression":  "yearly",  // 次限默认年级别
		"tertiary_progression":   "monthly", // 三限默认月级别
		"retrograde":             "monthly", // 逆行默认月级别（实际由行星细分）
		"dignity":                "monthly", // 尊贵度默认月级别
		"voidOfCourse":           "hourly",  // 月空默认小时级别
		"profectionLord":         "yearly",  // 年主星默认年级别
	}
	if level, ok := mapping[eventType]; ok {
		return level
	}
	return "daily" // 默认日级别
}

// ==================== 粒度分级系统 ====================
// 根据天体运行速度和影响持续时间，精确分配每个事件的时间级别

// getAspectTimeLevel 根据相位涉及的行星获取时间级别
// 规则：取行运行星（planet1）的速度决定
func getAspectTimeLevel(transitPlanet, natalPlanet models.PlanetID) string {
	return getPlanetTimeLevel(transitPlanet)
}

// getPlanetTimeLevel 根据行星速度获取时间级别
func getPlanetTimeLevel(planet models.PlanetID) string {
	switch planet {
	case models.Moon:
		return "hourly" // 月亮：约2.5天过一宫，相位持续数小时
	case models.Sun, models.Mercury, models.Venus:
		return "daily" // 内行星：相位持续1-3天
	case models.Mars:
		return "weekly" // 火星：相位持续1-2周
	case models.Jupiter, models.Saturn:
		return "monthly" // 木土：相位持续数周到数月
	case models.Uranus, models.Neptune, models.Pluto:
		return "yearly" // 外行星：相位持续数月到数年
	case models.NorthNode, models.Chiron:
		return "monthly" // 交点和凯龙：较慢移动
	default:
		return "daily"
	}
}

// getRetrogradeTimeLevel 根据逆行行星获取时间级别
func getRetrogradeTimeLevel(planet models.PlanetID) string {
	switch planet {
	case models.Mercury:
		return "weekly" // 水星逆行：约3周
	case models.Venus, models.Mars:
		return "monthly" // 金火逆行：约6周
	case models.Jupiter, models.Saturn, models.Uranus, models.Neptune, models.Pluto:
		return "yearly" // 外行星逆行：数月
	default:
		return "monthly"
	}
}

// getDignityTimeLevel 根据行星获取尊贵度时间级别
func getDignityTimeLevel(planet models.PlanetID) string {
	switch planet {
	case models.Moon:
		return "hourly" // 月亮换座：约2.5天
	case models.Sun:
		return "monthly" // 太阳换座：约30天
	case models.Mercury, models.Venus:
		return "weekly" // 水金换座：数天到数周
	case models.Mars:
		return "monthly" // 火星换座：约2个月
	case models.Jupiter:
		return "yearly" // 木星换座：约1年
	case models.Saturn, models.Uranus, models.Neptune, models.Pluto:
		return "yearly" // 土星及外行星：数年
	default:
		return "monthly"
	}
}

// getSignChangeTimeLevel 根据换座行星获取时间级别
func getSignChangeTimeLevel(planet models.PlanetID) string {
	// 换座事件的时间级别取决于行星在该星座停留多久
	switch planet {
	case models.Moon:
		return "hourly" // 月亮：约2.5天/星座
	case models.Sun:
		return "monthly" // 太阳：约30天/星座
	case models.Mercury, models.Venus:
		return "weekly" // 水金：约3-4周/星座（逆行时更长）
	case models.Mars:
		return "monthly" // 火星：约2个月/星座
	case models.Jupiter:
		return "yearly" // 木星：约1年/星座
	case models.Saturn:
		return "yearly" // 土星：约2.5年/星座
	case models.Uranus, models.Neptune, models.Pluto:
		return "yearly" // 外行星：7-20年/星座
	default:
		return "monthly"
	}
}

// getSmartTimeLevel 智能分配时间级别
// 根据事件类型和涉及的行星，精确分配粒度级别
func getSmartTimeLevel(eventType string, planet1, planet2 models.PlanetID) string {
	switch eventType {
	case "aspect":
		// 相位事件：根据行运行星（planet1）速度决定
		if planet1 != "" {
			return getAspectTimeLevel(planet1, planet2)
		}
		return "daily"
	
	case "sign_change":
		// 换座事件：根据换座行星决定
		if planet1 != "" {
			return getSignChangeTimeLevel(planet1)
		}
		return "monthly"
	
	case "retrograde":
		// 逆行事件：根据逆行行星决定
		if planet1 != "" {
			return getRetrogradeTimeLevel(planet1)
		}
		return "monthly"
	
	case "dignity":
		// 尊贵度事件：根据行星决定
		if planet1 != "" {
			return getDignityTimeLevel(planet1)
		}
		return "monthly"
	
	case "transit_house":
		// 行运过宫：根据行运行星决定
		if planet1 != "" {
			return getPlanetTimeLevel(planet1)
		}
		return "weekly"
	
	case "lunar_phase":
		return "daily" // 月相：固定日级别
	
	case "planetary_hour_change":
		return "hourly" // 行星时：固定小时级别
	
	case "voidOfCourse":
		return "hourly" // 月空：固定小时级别
	
	case "secondary_progression":
		return "yearly" // 次限推进：固定年级别
	
	case "tertiary_progression":
		return "monthly" // 三限推进：固定月级别
	
	case "profectionLord":
		return "yearly" // 年主星：固定年级别
	
	default:
		return getDefaultTimeLevelForEventType(eventType)
	}
}

// convertTransitHouseToUnifiedEvent converts a transit house event to unified event
func convertTransitHouseToUnifiedEvent(e astro.TransitHouseEvent, translator *i18n.Translator, now time.Time) UnifiedEvent {
	// Determine time level based on planet
	timeLevel := "weekly"
	switch e.Planet {
	case models.Moon:
		timeLevel = "hourly"
	case models.Sun, models.Mercury, models.Venus:
		timeLevel = "daily"
	case models.Mars:
		timeLevel = "weekly"
	case models.Jupiter, models.Saturn:
		timeLevel = "monthly"
	case models.Uranus, models.Neptune, models.Pluto:
		timeLevel = "yearly"
	}
	
	return UnifiedEvent{
		Time:           e.EntryTime,
		Type:           "transit_house",
		Title:          e.Title,
		Description:    e.Description,
		Theme:          e.Theme,
		Advice:         e.Advice,
		IsPositive:     e.IsPositive,
		Intensity:      e.Intensity,
		Planet1:        e.Planet,
		House:          e.House,
		IsExactToday:   false, // Transit houses are ongoing
		InfluencePhase: "active",
		StartDate:      e.EntryTime.Format("2006-01-02"),
		EndDate:        e.ExitTime.Format("2006-01-02"),
		DurationDays:   e.DurationDays,
		Factor: &FactorImpact{
			FactorType: "transitHouse",
			TimeLevel:  timeLevel,
			BaseValue:  0.5,
			Weight:     0.6,
			Strength:   0.8,
			DimensionImpact: getHouseDimensionImpactAPI(e.House),
		},
	}
}

// convertProgressionToUnifiedEvent converts a progression aspect event to unified event
func convertProgressionToUnifiedEvent(e astro.ProgressionAspectEvent, translator *i18n.Translator, now time.Time) UnifiedEvent {
	eventType := "secondary_progression"
	timeLevel := "yearly"
	if e.ProgressionType == astro.TertiaryProgression {
		eventType = "tertiary_progression"
		timeLevel = "monthly"
	}
	
	// Calculate orb-based strength
	strength := 1.0 - (e.Orb / 1.5)
	if strength < 0.3 {
		strength = 0.3
	}
	
	return UnifiedEvent{
		Time:           e.ExactDate,
		Type:           eventType,
		Title:          e.Title,
		Description:    e.Description,
		Theme:          e.Theme,
		Advice:         e.Advice,
		IsPositive:     e.IsPositive,
		Intensity:      e.Intensity,
		Planet1:        e.ProgressedPlanet,
		Planet2:        e.NatalPlanet,
		Aspect:         e.AspectType,
		Degree:         e.AspectAngle,
		IsExactToday:   e.IsExact,
		InfluencePhase: getProgressionInfluencePhase(e),
		StartDate:      e.StartDate.Format("2006-01-02"),
		EndDate:        e.EndDate.Format("2006-01-02"),
		DurationDays:   e.EndDate.Sub(e.StartDate).Hours() / 24,
		Factor: &FactorImpact{
			FactorType: eventType,
			TimeLevel:  timeLevel,
			BaseValue:  getProgressionBaseValue(e.IsPositive),
			Weight:     0.8,
			Strength:   strength,
			DimensionImpact: getProgressionDimensionImpactAPI(e.ProgressedPlanet, e.NatalPlanet),
		},
	}
}

// getProgressionInfluencePhase determines the influence phase of a progression
func getProgressionInfluencePhase(e astro.ProgressionAspectEvent) string {
	if e.IsApplying {
		return "approaching"
	}
	if e.IsExact {
		return "active"
	}
	return "fading"
}

// getProgressionBaseValue returns base value based on aspect harmony
func getProgressionBaseValue(isPositive bool) float64 {
	if isPositive {
		return 0.7
	}
	return 0.3
}

// getHouseDimensionImpactAPI returns dimension impact for house
func getHouseDimensionImpactAPI(house int) models.DimensionImpact {
	// Designed for max 2 dimension labels - matches transit_houses.go
	// Single dimension houses: 2, 4, 7, 9, 10, 12
	// Dual dimension houses: 1, 3, 5, 6, 8, 11
	impacts := map[int]models.DimensionImpact{
		1:  {Health: 0.6, Career: 0.4, Relationship: 0, Finance: 0, Spiritual: 0},        // 健康+事业
		2:  {Finance: 0.8, Career: 0, Relationship: 0, Health: 0, Spiritual: 0},          // 财运
		3:  {Career: 0.5, Relationship: 0.5, Health: 0, Finance: 0, Spiritual: 0},        // 事业+关系
		4:  {Relationship: 0.8, Career: 0, Health: 0, Finance: 0, Spiritual: 0},          // 关系
		5:  {Relationship: 0.6, Spiritual: 0.4, Career: 0, Health: 0, Finance: 0},        // 关系+灵性
		6:  {Health: 0.6, Career: 0.4, Relationship: 0, Finance: 0, Spiritual: 0},        // 健康+事业
		7:  {Relationship: 0.9, Career: 0, Health: 0, Finance: 0, Spiritual: 0},          // 关系
		8:  {Spiritual: 0.6, Finance: 0.4, Career: 0, Relationship: 0, Health: 0},        // 灵性+财运
		9:  {Spiritual: 0.8, Career: 0, Relationship: 0, Health: 0, Finance: 0},          // 灵性
		10: {Career: 0.9, Relationship: 0, Health: 0, Finance: 0, Spiritual: 0},          // 事业
		11: {Relationship: 0.6, Spiritual: 0.4, Career: 0, Health: 0, Finance: 0},        // 关系+灵性
		12: {Spiritual: 0.8, Career: 0, Relationship: 0, Health: 0, Finance: 0},          // 灵性
	}
	
	if impact, ok := impacts[house]; ok {
		return impact
	}
	return models.DimensionImpact{Career: 0.5, Relationship: 0, Health: 0, Finance: 0, Spiritual: 0.5}
}

// getAspectDimensionImpactAPI returns dimension impact for aspect events
func getAspectDimensionImpactAPI(planet1, planet2 models.PlanetID) models.DimensionImpact {
	// 取两颗行星的平均影响
	impact1 := astro.GetPlanetDimensionImpact(planet1)
	impact2 := astro.GetPlanetDimensionImpact(planet2)
	
	return models.DimensionImpact{
		Career:       (impact1.Career + impact2.Career) / 2,
		Relationship: (impact1.Relationship + impact2.Relationship) / 2,
		Health:       (impact1.Health + impact2.Health) / 2,
		Finance:      (impact1.Finance + impact2.Finance) / 2,
		Spiritual:    (impact1.Spiritual + impact2.Spiritual) / 2,
	}
}

// getProgressionDimensionImpactAPI returns dimension impact for progression
func getProgressionDimensionImpactAPI(progPlanet, natalPlanet models.PlanetID) models.DimensionImpact {
	progImpact := astro.GetPlanetDimensionImpact(progPlanet)
	natalImpact := astro.GetPlanetDimensionImpact(natalPlanet)
	
	return models.DimensionImpact{
		Career:       (progImpact.Career + natalImpact.Career) / 2,
		Relationship: (progImpact.Relationship + natalImpact.Relationship) / 2,
		Health:       (progImpact.Health + natalImpact.Health) / 2,
		Finance:      (progImpact.Finance + natalImpact.Finance) / 2,
		Spiritual:    (progImpact.Spiritual + natalImpact.Spiritual) / 2,
	}
}

// applyTranslationsToEvent applies translations, emotional titles, and formatting to an event
func applyTranslationsToEvent(event *UnifiedEvent, translator *i18n.Translator, now time.Time) {
	// Translate title components
	if event.Planet1 != "" {
		p1Name := translator.GetPlanetName(event.Planet1)
		p2Name := ""
		if event.Planet2 != "" {
			p2Name = translator.GetPlanetName(event.Planet2)
		}
		aspectName := translator.GetAspectName(event.Aspect)
		
		// Reconstruct title with translated names
		if p2Name != "" && aspectName != "" {
			event.Title = fmt.Sprintf("%s %s %s", p1Name, aspectName, p2Name)
		}
	}
	
	// Add emotional title and detailed interpretation
	houseStr := ""
	if event.House > 0 {
		houseStr = fmt.Sprintf("%d", event.House)
	}
	event.EmotionalTitle = translator.GetEmotionalTitle(
		event.Type,
		event.Planet1,
		event.Planet2,
		event.Aspect,
		houseStr,
		event.IsPositive,
	)
	event.DetailedInterpretation = translator.GetDetailedInterpretation(
		event.Type,
		event.Planet1,
		event.Planet2,
		event.Aspect,
		houseStr,
		event.IsPositive,
	)
	
	// Format duration text
	if event.StartDate != "" && event.EndDate != "" {
		startTime, _ := time.Parse("2006-01-02", event.StartDate)
		endTime, _ := time.Parse("2006-01-02", event.EndDate)
		
		daysAgo := int(now.Sub(startTime).Hours() / 24)
		daysUntil := int(endTime.Sub(now).Hours() / 24)
		
		// Convert to months if duration is long
		if daysAgo > 60 || daysUntil > 60 {
			monthsAgo := daysAgo / 30
			monthsUntil := daysUntil / 30
			event.DurationText = translator.FormatDuration(monthsAgo, monthsUntil)
		} else {
			event.DurationText = translator.FormatDurationDays(daysAgo, daysUntil)
		}
	}
	
	// Add dimension labels
	if event.Factor != nil {
		event.DimensionLabels = getDimensionLabels(event.Factor.DimensionImpact, translator)
	}
}

// getDimensionLabels returns dimension labels with arrows based on impact values
func getDimensionLabels(impact models.DimensionImpact, translator *i18n.Translator) []string {
	type dimValue struct {
		name  string
		value float64
	}
	
	dimensions := []dimValue{
		{"career", impact.Career},
		{"relationship", impact.Relationship},
		{"health", impact.Health},
		{"finance", impact.Finance},
		{"spiritual", impact.Spiritual},
	}
	
	// Filter dimensions with significant impact (abs value > 0.3)
	var significant []dimValue
	for _, d := range dimensions {
		if d.value > 0.3 || d.value < -0.3 {
			significant = append(significant, d)
		}
	}
	
	// Sort by absolute value (descending)
	for i := 0; i < len(significant); i++ {
		for j := i + 1; j < len(significant); j++ {
			absI := significant[i].value
			if absI < 0 {
				absI = -absI
			}
			absJ := significant[j].value
			if absJ < 0 {
				absJ = -absJ
			}
			if absJ > absI {
				significant[i], significant[j] = significant[j], significant[i]
			}
		}
	}
	
	// Return top 2 dimensions
	labels := []string{}
	maxDims := 2
	if len(significant) < maxDims {
		maxDims = len(significant)
	}
	
	for i := 0; i < maxDims; i++ {
		label := translator.GetDimensionLabel(significant[i].name, significant[i].value)
		labels = append(labels, label)
	}
	
	return labels
}
