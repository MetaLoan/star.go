package astro

import (
	"star/models"
	"time"
)

/*
分值组成查询模块

根据开发备忘录的可见性规范：
- 年度级别影响因子：在年/月/周/日/小时中可见
- 月级别影响因子：在月/周/日/小时中可见
- 周级别影响因子：在周/日/小时中可见
- 日级别影响因子：在日/小时中可见
- 小时级别影响因子：仅在小时中可见
*/

// ==================== 数据结构 ====================

// ScoreBreakdownRequest 分值组成查询请求
type ScoreBreakdownRequest struct {
	BirthData   models.BirthData `json:"birthData"`
	QueryTime   string           `json:"queryTime"`   // ISO 8601 格式
	Granularity string           `json:"granularity"` // hour, day, week, month, year
	UserID      string           `json:"userId,omitempty"`
}

// FactorContribution 因子贡献详情
type FactorContribution struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Type        string             `json:"type"`      // dignity, retrograde, aspectPhase, etc.
	TimeLevel   string             `json:"timeLevel"` // yearly, monthly, weekly, daily, hourly
	BaseValue   float64            `json:"baseValue"` // 基础值
	Strength    float64            `json:"strength"`  // 当前强度 (0-1)
	Adjustment  float64            `json:"adjustment"`// 实际调整值 = BaseValue * Strength * Weight
	Weight      float64            `json:"weight"`    // 权重
	Dimension   string             `json:"dimension,omitempty"` // 主要影响维度
	Description string             `json:"description"`
	IsPositive  bool               `json:"isPositive"`
}

// DimensionBreakdown 维度分值分解
type DimensionBreakdown struct {
	Dimension   string               `json:"dimension"`   // career, relationship, health, finance, spiritual
	BaseScore   float64              `json:"baseScore"`   // 基础分（50）
	AspectScore float64              `json:"aspectScore"` // 相位贡献
	FactorScore float64              `json:"factorScore"` // 因子贡献
	RawScore    float64              `json:"rawScore"`    // 原始分 = Base + Aspect + Factor
	FinalScore  float64              `json:"finalScore"`  // 最终分（标准化后）
	Factors     []FactorContribution `json:"factors"`     // 贡献的因子列表
}

// ScoreBreakdownResponse 分值组成查询响应
type ScoreBreakdownResponse struct {
	QueryTime      string               `json:"queryTime"`
	Granularity    string               `json:"granularity"`
	
	// 综合分
	OverallScore   float64              `json:"overallScore"`
	OverallRaw     float64              `json:"overallRaw"`
	
	// 五维度分解
	Dimensions     []DimensionBreakdown `json:"dimensions"`
	
	// 所有可见因子（按时间级别分组）
	FactorsByLevel map[string][]FactorContribution `json:"factorsByLevel"`
	
	// 元信息
	Meta           ScoreBreakdownMeta   `json:"meta"`
}

// ScoreBreakdownMeta 元信息
type ScoreBreakdownMeta struct {
	DataSource        string   `json:"dataSource"`        // Swiss Ephemeris
	VisibleLevels     []string `json:"visibleLevels"`     // 当前粒度可见的时间级别
	TotalFactorCount  int      `json:"totalFactorCount"`  // 总因子数
	PositiveFactors   int      `json:"positiveFactors"`   // 正向因子数
	NegativeFactors   int      `json:"negativeFactors"`   // 负向因子数
}

// ==================== 核心计算 ====================

// GetVisibleTimeLevels 根据查询粒度获取可见的时间级别
func GetVisibleTimeLevels(granularity string) []models.FactorTimeLevel {
	switch granularity {
	case "hour":
		// 小时级：所有级别可见
		return []models.FactorTimeLevel{
			models.TimeLevelYearly,
			models.TimeLevelMonthly,
			models.TimeLevelWeekly,
			models.TimeLevelDaily,
			models.TimeLevelHourly,
		}
	case "day":
		// 日级：年/月/周/日可见，不含小时
		return []models.FactorTimeLevel{
			models.TimeLevelYearly,
			models.TimeLevelMonthly,
			models.TimeLevelWeekly,
			models.TimeLevelDaily,
		}
	case "week":
		// 周级：年/月/周可见
		return []models.FactorTimeLevel{
			models.TimeLevelYearly,
			models.TimeLevelMonthly,
			models.TimeLevelWeekly,
		}
	case "month":
		// 月级：年/月可见
		return []models.FactorTimeLevel{
			models.TimeLevelYearly,
			models.TimeLevelMonthly,
		}
	case "year":
		// 年级：仅年可见
		return []models.FactorTimeLevel{
			models.TimeLevelYearly,
		}
	default:
		return []models.FactorTimeLevel{models.TimeLevelDaily}
	}
}

// IsFactorVisible 判断因子在当前粒度是否可见
func IsFactorVisible(factorLevel models.FactorTimeLevel, queryGranularity string) bool {
	visibleLevels := GetVisibleTimeLevels(queryGranularity)
	return isTimeLevelVisible(factorLevel, visibleLevels)
}

// isTimeLevelVisible 检查时间级别是否在可见列表中
func isTimeLevelVisible(factorLevel models.FactorTimeLevel, visibleLevels []models.FactorTimeLevel) bool {
	for _, level := range visibleLevels {
		if level == factorLevel {
			return true
		}
	}
	return false
}

// CalculateScoreBreakdown 计算分值组成详情
func CalculateScoreBreakdown(chart *models.NatalChart, t time.Time, granularity string, userID string) ScoreBreakdownResponse {
	// 1. 获取行运位置
	transitPositions := GetTransitPositions(t)
	
	// 2. 计算相位
	aspects := CalculateTransitToNatalAspects(transitPositions, chart.Planets)
	
	// 3. 计算每个维度的相位贡献
	dimensionAspectScores := calculateDimensionAspectScoresDetailed(aspects)
	
	// 4. 获取所有影响因子
	factorResult := CalculateInfluenceFactors(chart, t, transitPositions)
	
	// 5. 过滤可见因子
	visibleLevels := GetVisibleTimeLevels(granularity)
	visibleLevelStrings := make([]string, len(visibleLevels))
	for i, l := range visibleLevels {
		visibleLevelStrings[i] = string(l)
	}
	
	// 6. 分类因子
	factorsByLevel := make(map[string][]FactorContribution)
	var allVisibleFactors []FactorContribution
	positiveCount := 0
	negativeCount := 0
	
	if factorResult != nil {
		for _, f := range factorResult.Factors {
			// 检查因子是否可见
			if !IsFactorVisible(f.TimeLevel, granularity) {
				continue
			}
			
			fc := FactorContribution{
				ID:          f.ID,
				Name:        f.Name,
				Type:        string(f.Type),
				TimeLevel:   string(f.TimeLevel),
				BaseValue:   f.BaseValue,
				Strength:    f.CurrentStrength,
				Adjustment:  f.Adjustment,
				Weight:      f.Weight,
				Dimension:   string(f.SourcePlanet),
				Description: f.Description,
				IsPositive:  f.IsPositive,
			}
			
			levelStr := string(f.TimeLevel)
			factorsByLevel[levelStr] = append(factorsByLevel[levelStr], fc)
			allVisibleFactors = append(allVisibleFactors, fc)
			
			if f.IsPositive {
				positiveCount++
			} else {
				negativeCount++
			}
		}
	}
	
	// 7. 添加自定义因子
	if userID != "" {
		customFactors := GetActiveCustomFactors(userID, t)
		for _, cf := range customFactors {
			lifecycle := &models.FactorLifecycle{
				StartTime: cf.StartTime,
				PeakTime:  cf.StartTime.Add(time.Duration(cf.Duration/2) * time.Hour),
				EndTime:   cf.EndTime,
				Duration:  cf.Duration,
			}
			strength := CalculateFactorStrength(lifecycle, t)
			
			fc := FactorContribution{
				ID:          cf.ID,
				Name:        cf.Description,
				Type:        "custom",
				TimeLevel:   "hourly", // 自定义因子默认为小时级
				BaseValue:   cf.Value,
				Strength:    strength,
				Adjustment:  cf.Value * strength,
				Weight:      1.0,
				Dimension:   cf.Dimension,
				Description: cf.Description,
				IsPositive:  cf.Value > 0,
			}
			
			// 自定义因子在所有粒度都可见
			factorsByLevel["custom"] = append(factorsByLevel["custom"], fc)
			allVisibleFactors = append(allVisibleFactors, fc)
			
			if cf.Value > 0 {
				positiveCount++
			} else {
				negativeCount++
			}
		}
	}
	
	// 8. 计算每个维度的详细分解
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}
	var dimensionBreakdowns []DimensionBreakdown
	
	for _, dim := range dimensions {
		baseScore := 50.0 // 基础分
		aspectScore := dimensionAspectScores[dim] - baseScore // 相位贡献（减去基础分）
		
		// 计算因子对该维度的贡献
		factorScore := 0.0
		var dimFactors []FactorContribution
		
		for _, fc := range allVisibleFactors {
			// 计算该因子对该维度的贡献
			contribution := calculateFactorDimensionContribution(fc, dim)
			if contribution != 0 {
				factorScore += contribution
				// 创建带维度贡献的因子副本
				fcCopy := fc
				fcCopy.Adjustment = contribution
				dimFactors = append(dimFactors, fcCopy)
			}
		}
		
		// 计算原始分和最终分
		rawScore := baseScore + aspectScore + factorScore
		finalScore := NormalizeScoreV2(rawScore)
		
		dimensionBreakdowns = append(dimensionBreakdowns, DimensionBreakdown{
			Dimension:   dim,
			BaseScore:   baseScore,
			AspectScore: aspectScore,
			FactorScore: factorScore,
			RawScore:    rawScore,
			FinalScore:  finalScore,
			Factors:     dimFactors,
		})
	}
	
	// 9. 计算综合分
	var totalRaw float64
	for _, db := range dimensionBreakdowns {
		totalRaw += db.RawScore * 0.2 // 五维度等权重
	}
	overallScore := NormalizeScoreV2(totalRaw)
	
	// 10. 构建响应
	return ScoreBreakdownResponse{
		QueryTime:    t.Format(time.RFC3339),
		Granularity:  granularity,
		OverallScore: overallScore,
		OverallRaw:   totalRaw,
		Dimensions:   dimensionBreakdowns,
		FactorsByLevel: factorsByLevel,
		Meta: ScoreBreakdownMeta{
			DataSource:       "Swiss Ephemeris",
			VisibleLevels:    visibleLevelStrings,
			TotalFactorCount: len(allVisibleFactors),
			PositiveFactors:  positiveCount,
			NegativeFactors:  negativeCount,
		},
	}
}

// calculateDimensionAspectScoresDetailed 计算维度相位分数（详细版）
func calculateDimensionAspectScoresDetailed(aspects []models.AspectData) map[string]float64 {
	dimensions := []string{"career", "relationship", "health", "finance", "spiritual"}
	scores := make(map[string]float64)
	
	// 初始化基础分
	for _, d := range dimensions {
		scores[d] = 50 // 基础分50
	}
	
	// 遍历每个相位
	for _, aspect := range aspects {
		transitPlanet := aspect.Planet1
		natalPlanet := aspect.Planet2
		
		// 计算相位分值
		aspectValue := getUnifiedAspectValue(string(aspect.AspectType), aspect.Orb)
		
		// 根据行星-维度权重分配分值
		for _, d := range dimensions {
			transitWeight := 0.5
			if w, ok := PlanetDimensionWeight[transitPlanet]; ok {
				if dw, ok := w[d]; ok {
					transitWeight = dw
				}
			}
			
			natalWeight := 0.5
			if w, ok := PlanetDimensionWeight[natalPlanet]; ok {
				if dw, ok := w[d]; ok {
					natalWeight = dw
				}
			}
			
			combinedWeight := (transitWeight + natalWeight) / 2
			scores[d] += aspectValue * combinedWeight
		}
	}
	
	return scores
}

// calculateFactorDimensionContribution 计算因子对特定维度的贡献
func calculateFactorDimensionContribution(fc FactorContribution, dimension string) float64 {
	// 如果因子指定了特定维度
	if fc.Dimension == dimension {
		return fc.Adjustment
	}
	
	// 如果是 overall 或未指定，按行星权重分配
	if fc.Dimension == "overall" || fc.Dimension == "" {
		// 平均分配到五个维度
		return fc.Adjustment / 5
	}
	
	// 根据来源行星分配
	planetID := models.PlanetID(fc.Dimension)
	if weights, ok := PlanetDimensionWeight[planetID]; ok {
		if w, ok := weights[dimension]; ok {
			totalWeight := 0.0
			for _, wt := range weights {
				totalWeight += wt
			}
			return fc.Adjustment * (w / totalWeight)
		}
	}
	
	return 0
}

// ==================== 多粒度查询 ====================

// GetMultiGranularityBreakdown 获取多粒度分值组成
func GetMultiGranularityBreakdown(chart *models.NatalChart, t time.Time, userID string) map[string]ScoreBreakdownResponse {
	granularities := []string{"hour", "day", "month", "year"}
	result := make(map[string]ScoreBreakdownResponse)
	
	for _, g := range granularities {
		result[g] = CalculateScoreBreakdown(chart, t, g, userID)
	}
	
	return result
}

// ==================== 时间范围内活跃因子查询 ====================

// ActiveFactorInfo 活跃因子信息
type ActiveFactorInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	TimeLevel   string  `json:"timeLevel"`
	BaseValue   float64 `json:"baseValue"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
	IsPositive  bool    `json:"isPositive"`
	Effect      string  `json:"effect"` // "positive" 或 "negative"
	// 生命周期信息
	StartTime   string  `json:"startTime,omitempty"`
	EndTime     string  `json:"endTime,omitempty"`
	PeakTime    string  `json:"peakTime,omitempty"`
	// 在查询范围内的最大强度
	MaxStrength float64 `json:"maxStrength"`
}

// ActiveFactorsResponse 活跃因子查询响应
type ActiveFactorsResponse struct {
	Granularity   string             `json:"granularity"`   // year/month/week/day
	RangeStart    string             `json:"rangeStart"`    // 查询范围开始
	RangeEnd      string             `json:"rangeEnd"`      // 查询范围结束
	Infect        string             `json:"infect"`        // all/core（是否过滤）
	Factors       []ActiveFactorInfo `json:"factors"`       // 活跃因子列表
	TotalCount    int                `json:"totalCount"`    // 总数
	PositiveCount int                `json:"positiveCount"` // 正向因子数
	NegativeCount int                `json:"negativeCount"` // 负向因子数
}

// GetActiveFactorsInRange 获取时间范围内活跃的所有因子
// granularity: year/month/week/day
// t: 范围内的任意时间点
// infect: "all" 返回所有因子, "core" 按可见性规则过滤
func GetActiveFactorsInRange(chart *models.NatalChart, t time.Time, granularity string, infect string, userID string) ActiveFactorsResponse {
	// 1. 根据粒度计算时间范围
	rangeStart, rangeEnd := calculateTimeRange(t, granularity)
	
	// 2. 获取当前粒度可见的时间级别（用于 core 过滤）
	visibleLevels := GetVisibleTimeLevels(granularity)
	shouldFilter := infect == "core"
	
	// 3. 采样时间点（根据粒度决定采样频率）
	samplePoints := generateSamplePoints(rangeStart, rangeEnd, granularity)
	
	// 4. 收集所有采样点的因子
	factorMap := make(map[string]*ActiveFactorInfo)
	
	for _, sampleTime := range samplePoints {
		transitPositions := GetTransitPositions(sampleTime)
		factorResult := CalculateInfluenceFactors(chart, sampleTime, transitPositions)
		
		if factorResult == nil {
			continue
		}
		
		for _, f := range factorResult.Factors {
			// core 模式：按可见性规则过滤
			if shouldFilter && !isTimeLevelVisible(f.TimeLevel, visibleLevels) {
				continue
			}
			
			key := f.ID
			
			// 计算当前强度
			strength := 1.0
			if f.Lifecycle != nil {
				strength = CalculateFactorStrength(f.Lifecycle, sampleTime)
			}
			
			if existing, ok := factorMap[key]; ok {
				// 更新最大强度
				if strength > existing.MaxStrength {
					existing.MaxStrength = strength
				}
			} else {
				// 新因子
				effect := "negative"
				if f.IsPositive {
					effect = "positive"
				}
				info := &ActiveFactorInfo{
					ID:          f.ID,
					Name:        f.Name,
					Type:        string(f.Type),
					TimeLevel:   string(f.TimeLevel),
					BaseValue:   f.BaseValue,
					Weight:      f.Weight,
					Description: f.Description,
					IsPositive:  f.IsPositive,
					Effect:      effect,
					MaxStrength: strength,
				}
				
				// 添加生命周期信息
				if f.Lifecycle != nil {
					info.StartTime = f.Lifecycle.StartTime.Format(time.RFC3339)
					info.EndTime = f.Lifecycle.EndTime.Format(time.RFC3339)
					info.PeakTime = f.Lifecycle.PeakTime.Format(time.RFC3339)
				}
				
				factorMap[key] = info
			}
		}
	}
	
	// 4. 添加自定义因子
	if userID != "" {
		customFactors := GetActiveCustomFactorsInRange(userID, rangeStart, rangeEnd)
		for _, cf := range customFactors {
			key := "custom_" + cf.ID
			if _, ok := factorMap[key]; !ok {
				isPositive := cf.Value > 0
				cfEffect := "negative"
				if isPositive {
					cfEffect = "positive"
				}
				factorMap[key] = &ActiveFactorInfo{
					ID:          cf.ID,
					Name:        cf.Description,
					Type:        "custom",
					TimeLevel:   "custom",
					BaseValue:   cf.Value,
					Weight:      1.0,
					Description: cf.Description,
					IsPositive:  isPositive,
					Effect:      cfEffect,
					StartTime:   cf.StartTime.Format(time.RFC3339),
					EndTime:     cf.EndTime.Format(time.RFC3339),
					MaxStrength: 1.0,
				}
			}
		}
	}
	
	// 5. 转换为切片并统计
	factors := make([]ActiveFactorInfo, 0, len(factorMap))
	positiveCount := 0
	negativeCount := 0
	
	for _, f := range factorMap {
		factors = append(factors, *f)
		if f.IsPositive {
			positiveCount++
		} else {
			negativeCount++
		}
	}
	
	return ActiveFactorsResponse{
		Granularity:   granularity,
		RangeStart:    rangeStart.Format(time.RFC3339),
		RangeEnd:      rangeEnd.Format(time.RFC3339),
		Infect:        infect,
		Factors:       factors,
		TotalCount:    len(factors),
		PositiveCount: positiveCount,
		NegativeCount: negativeCount,
	}
}

// calculateTimeRange 根据粒度和时间点计算时间范围
func calculateTimeRange(t time.Time, granularity string) (time.Time, time.Time) {
	switch granularity {
	case "year":
		// 整年
		start := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
		end := time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
		return start, end
		
	case "month":
		// 整月
		start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		end := start.AddDate(0, 1, 0)
		return start, end
		
	case "week":
		// 整周（周一到周日）
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
		end := start.AddDate(0, 0, 7)
		return start, end
		
	case "day":
		// 整天
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end := start.AddDate(0, 0, 1)
		return start, end
		
	case "hour":
		// 整小时
		start := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		end := start.Add(time.Hour)
		return start, end
		
	default:
		// 默认为天
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end := start.AddDate(0, 0, 1)
		return start, end
	}
}

// generateSamplePoints 生成采样时间点
func generateSamplePoints(start, end time.Time, granularity string) []time.Time {
	var points []time.Time
	
	switch granularity {
	case "year":
		// 年粒度：每月采样一次
		current := start
		for current.Before(end) {
			points = append(points, current)
			current = current.AddDate(0, 1, 0)
		}
		
	case "month":
		// 月粒度：每天采样一次
		current := start
		for current.Before(end) {
			points = append(points, current)
			current = current.AddDate(0, 0, 1)
		}
		
	case "week":
		// 周粒度：每天采样两次（早晚）
		current := start
		for current.Before(end) {
			points = append(points, current)
			points = append(points, current.Add(12*time.Hour))
			current = current.AddDate(0, 0, 1)
		}
		
	case "day":
		// 日粒度：每2小时采样一次
		current := start
		for current.Before(end) {
			points = append(points, current)
			current = current.Add(2 * time.Hour)
		}
		
	case "hour":
		// 小时粒度：每10分钟采样一次
		current := start
		for current.Before(end) {
			points = append(points, current)
			current = current.Add(10 * time.Minute)
		}
		
	default:
		points = append(points, start)
	}
	
	return points
}

// GetActiveCustomFactorsInRange 获取时间范围内活跃的自定义因子
func GetActiveCustomFactorsInRange(userID string, start, end time.Time) []CustomFactorDefinition {
	factors := GetAllCustomFactors(userID)
	var active []CustomFactorDefinition
	
	for _, f := range factors {
		// 检查因子的生命周期是否与查询范围有交集
		if !f.EndTime.Before(start) && !f.StartTime.After(end) {
			active = append(active, f)
		}
	}
	
	return active
}

