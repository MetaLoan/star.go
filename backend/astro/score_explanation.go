package astro

import (
	"fmt"
	"sort"
	"star/models"
	"time"
)

/*
分值解释模块 - 面向C端用户
用通俗易懂的语言解释分数是如何计算的
*/

// ==================== C端友好的数据结构 ====================

// ScoreExplanationResponse C端用户友好的分值解释响应
type ScoreExplanationResponse struct {
	// 基本信息
	QueryTime   string `json:"queryTime"`
	TimeLabel   string `json:"timeLabel"`   // "2026年1月5日 14:30"
	Granularity string `json:"granularity"` // "小时" | "日" | "周" | "月" | "年"
	Dimension   string `json:"dimension"`   // "健康" | "事业" | "关系" | "财务" | "灵性" | "综合"

	// 分数
	Score       float64 `json:"score"`       // 最终分数
	ScoreLevel  string  `json:"scoreLevel"`  // "极佳" | "良好" | "平稳" | "需注意" | "挑战"
	ScoreEmoji  string  `json:"scoreEmoji"`  // 🌟 | ✨ | 💫 | 🌙 | ⚡

	// 分数解释
	Explanation ScoreExplanationDetail `json:"explanation"`

	// 天文现象影响
	AstronomicalFactors []AstronomicalFactor `json:"astronomicalFactors"`

	// 总结
	Summary     string   `json:"summary"`     // 一句话总结
	Suggestions []string `json:"suggestions"` // 建议
}

// ScoreExplanationDetail 分数解释详情
type ScoreExplanationDetail struct {
	Formula      string  `json:"formula"`      // "基础分 + 相位影响 + 其他因素 = 原始分 → 标准化 → 最终分"
	BaseScore    float64 `json:"baseScore"`    // 基础分
	AspectEffect float64 `json:"aspectEffect"` // 相位影响
	FactorEffect float64 `json:"factorEffect"` // 其他因素影响
	RawScore     float64 `json:"rawScore"`     // 原始分
	Description  string  `json:"description"`  // 解释文字
}

// AstronomicalFactor 天文现象因素
type AstronomicalFactor struct {
	// 展示信息
	Name        string  `json:"name"`        // "火星与月亮形成和谐相位"
	Category    string  `json:"category"`    // "行星相位" | "行星状态" | "月相" | "行星时" | "逆行"
	Icon        string  `json:"icon"`        // ♂️ | 🌙 | etc
	Effect      string  `json:"effect"`      // "增强" | "减弱"
	EffectValue float64 `json:"effectValue"` // +0.52 或 -0.30
	Intensity   string  `json:"intensity"`   // "强" | "中" | "弱"
	IsPositive  bool    `json:"isPositive"`  // 是否正向影响
	Type        string  `json:"type"`        // 因素类型
	
	// 详细解释
	Description      string `json:"description"`      // 对用户友好的解释
	AstroExplanation string `json:"astroExplanation"` // 占星学解释
	
	// 时效性
	TimeLevel   string `json:"timeLevel"`   // "长期" | "月度" | "本周" | "今日" | "当前小时"
	ValidPeriod string `json:"validPeriod"` // "持续整月" | "本周有效" | "今日有效" | "当前小时有效"
}

// ==================== 核心计算函数 ====================

// GetScoreExplanation 获取分数解释（面向C端用户）
func GetScoreExplanation(chart *models.NatalChart, t time.Time, granularity string, dimension string, userID string) ScoreExplanationResponse {
	// 1. 计算详细分值组成
	breakdown := CalculateScoreBreakdown(chart, t, granularity, userID)
	
	// 2. 找到目标维度
	var dimBreakdown *DimensionBreakdown
	for _, db := range breakdown.Dimensions {
		if db.Dimension == dimension {
			dimBreakdown = &db
			break
		}
	}
	
	// 如果是综合分
	if dimension == "overall" || dimBreakdown == nil {
		return buildOverallExplanation(breakdown, t, granularity)
	}
	
	// 3. 构建用户友好的响应
	return buildDimensionExplanation(*dimBreakdown, breakdown, t, granularity, dimension)
}

// buildDimensionExplanation 构建维度解释
func buildDimensionExplanation(dim DimensionBreakdown, breakdown ScoreBreakdownResponse, t time.Time, granularity string, dimension string) ScoreExplanationResponse {
	score := dim.FinalScore
	
	// 分数等级
	scoreLevel, scoreEmoji := getScoreLevel(score)
	
	// 转换天文因素
	astroFactors := convertToAstronomicalFactors(dim.Factors)
	
	// 排序：按影响值绝对值降序
	sort.Slice(astroFactors, func(i, j int) bool {
		return abs(astroFactors[i].EffectValue) > abs(astroFactors[j].EffectValue)
	})
	
	// 生成解释和建议
	explanation := buildExplanationDetail(dim)
	summary := buildSummary(dimension, score, astroFactors)
	suggestions := buildSuggestions(dimension, score, astroFactors)
	
	return ScoreExplanationResponse{
		QueryTime:           t.Format(time.RFC3339),
		TimeLabel:           formatTimeLabel(t, granularity),
		Granularity:         getGranularityLabel(granularity),
		Dimension:           getDimensionLabel(dimension),
		Score:               score,
		ScoreLevel:          scoreLevel,
		ScoreEmoji:          scoreEmoji,
		Explanation:         explanation,
		AstronomicalFactors: astroFactors,
		Summary:             summary,
		Suggestions:         suggestions,
	}
}

// buildOverallExplanation 构建综合分解释
func buildOverallExplanation(breakdown ScoreBreakdownResponse, t time.Time, granularity string) ScoreExplanationResponse {
	score := breakdown.OverallScore
	scoreLevel, scoreEmoji := getScoreLevel(score)
	
	// 收集所有维度的因素
	var allFactors []FactorContribution
	for _, dim := range breakdown.Dimensions {
		allFactors = append(allFactors, dim.Factors...)
	}
	
	// 去重并转换
	astroFactors := convertToAstronomicalFactors(allFactors)
	
	// 排序
	sort.Slice(astroFactors, func(i, j int) bool {
		return abs(astroFactors[i].EffectValue) > abs(astroFactors[j].EffectValue)
	})
	
	// 只保留前15个最重要的
	if len(astroFactors) > 15 {
		astroFactors = astroFactors[:15]
	}
	
	return ScoreExplanationResponse{
		QueryTime:   t.Format(time.RFC3339),
		TimeLabel:   formatTimeLabel(t, granularity),
		Granularity: getGranularityLabel(granularity),
		Dimension:   "Overall Fortune",
		Score:       score,
		ScoreLevel:  scoreLevel,
		ScoreEmoji:  scoreEmoji,
		Explanation: ScoreExplanationDetail{
			Formula:     "Five Dimension Weighted Average → Normalization → Final Score",
			BaseScore:   50,
			RawScore:    breakdown.OverallRaw,
			Description: "Overall fortune is calculated from career, relationship, health, finance, and spiritual dimensions, each weighted at 20%",
		},
		AstronomicalFactors: astroFactors,
		Summary:             buildOverallSummary(score, astroFactors),
		Suggestions:         buildOverallSuggestions(score),
	}
}

// ==================== 辅助函数 ====================

// convertToAstronomicalFactors 转换为天文因素
func convertToAstronomicalFactors(factors []FactorContribution) []AstronomicalFactor {
	// 去重
	seen := make(map[string]bool)
	var result []AstronomicalFactor
	
	for _, f := range factors {
		if seen[f.ID] {
			continue
		}
		seen[f.ID] = true
		
		af := AstronomicalFactor{
			Name:             getFactorFriendlyName(f),
			Category:         getFactorCategory(f.Type),
			Icon:             getFactorIcon(f),
			Effect:           getEffectLabel(f.IsPositive),
			EffectValue:      round2(f.Adjustment),
			Intensity:        getIntensity(f.Adjustment),
			IsPositive:       f.IsPositive,
			Type:             f.Type,
			Description:      getFactorDescription(f),
			AstroExplanation: getAstroExplanation(f),
			TimeLevel:        getTimeLevelLabel(f.TimeLevel),
			ValidPeriod:      getValidPeriod(f.TimeLevel),
		}
		result = append(result, af)
	}
	
	return result
}

// getFactorFriendlyName 获取因素的用户友好名称
func getFactorFriendlyName(f FactorContribution) string {
	// 根据类型和名称生成友好名称
	switch f.Type {
	case "dignity":
		return f.Name // 已经是友好名称如 "火星旺相"
	case "retrograde":
		return f.Name // 如 "木星逆行"
	case "aspectPhase":
		return f.Name // 如 "太阳六分相月亮"
	case "lunarPhase":
		return f.Name // 如 "满月期"
	case "planetaryHour":
		return f.Name // 如 "月亮日 土星时"
	case "voidOfCourse":
		return "Moon Void of Course"
	case "custom":
		return "Personal Setting: " + f.Name
	default:
		return f.Name
	}
}

// getFactorCategory 获取因素类别
func getFactorCategory(factorType string) string {
	switch factorType {
	case "dignity":
		return "Planetary Dignity"
	case "retrograde":
		return "Planetary Retrograde"
	case "aspectPhase":
		return "Planetary Aspect"
	case "lunarPhase":
		return "Lunar Phase"
	case "planetaryHour":
		return "Planetary Hour"
	case "voidOfCourse":
		return "Moon Void of Course"
	case "custom":
		return "Personal Factor"
	default:
		return "Other"
	}
}

// getFactorIcon 获取因素图标
func getFactorIcon(f FactorContribution) string {
	switch f.Type {
	case "dignity":
		return getPlanetEmoji(f.Dimension)
	case "retrograde":
		return "℞"
	case "aspectPhase":
		return "✦"
	case "lunarPhase":
		return "🌙"
	case "planetaryHour":
		return "⏰"
	case "voidOfCourse":
		return "🌑"
	case "custom":
		return "⚙️"
	default:
		return "✧"
	}
}

// getPlanetEmoji 获取行星符号
func getPlanetEmoji(planet string) string {
	emojis := map[string]string{
		"sun": "☀️", "moon": "🌙", "mercury": "☿️", "venus": "♀️",
		"mars": "♂️", "jupiter": "♃", "saturn": "♄", "uranus": "⛢",
		"neptune": "♆", "pluto": "⯓", "chiron": "⚷", "northNode": "☊",
	}
	if e, ok := emojis[planet]; ok {
		return e
	}
	return "✧"
}

// getEffectLabel 获取影响标签
func getEffectLabel(isPositive bool) string {
	if isPositive {
		return "Enhance"
	}
	return "Weaken"
}

// getIntensity 获取强度
func getIntensity(value float64) string {
	absVal := abs(value)
	if absVal >= 0.5 {
		return "Strong"
	} else if absVal >= 0.2 {
		return "Medium"
	}
	return "Weak"
}

// getFactorDescription 获取因素描述
func getFactorDescription(f FactorContribution) string {
	switch f.Type {
	case "dignity":
		if f.IsPositive {
			return "Planet in favorable position, energy enhanced"
		}
		return "Planet in unfavorable position, energy weakened"
	case "retrograde":
		return "During planetary retrograde, related areas need review and reflection"
	case "aspectPhase":
		if f.IsPositive {
			return "Planets form harmonious angle, bringing positive energy"
		}
		return "Planets form tense angle, bringing challenges"
	case "lunarPhase":
		return getLunarPhaseDescription(f.Name)
	case "planetaryHour":
		return "Current planetary hour energy influence"
	case "custom":
		return "Personal adjustment factor"
	default:
		return f.Description
	}
}

// getLunarPhaseDescription 获取月相描述
func getLunarPhaseDescription(name string) string {
	descriptions := map[string]string{
		"New Moon":           "New Moon period, ideal for setting new goals and planting intentions",
		"Crescent Moon":      "Moon waxing, energy rising, good time for action",
		"Gibbous Moon":       "Approaching full moon, time to prepare for harvest",
		"Full Moon":          "Full Moon peak, emotions and energy at their height, ideal for showcasing results",
		"Disseminating Moon": "After full moon, good for sharing and spreading",
		"Last Quarter":       "Energy decreasing, time for release and letting go",
		"Balsamic Moon":      "Moon about to disappear, time for rest and reflection",
	}
	if desc, ok := descriptions[name]; ok {
		return desc
	}
	return "Lunar phase influences daily energy and emotions"
}

// getAstroExplanation 获取占星学解释
func getAstroExplanation(f FactorContribution) string {
	switch f.Type {
	case "dignity":
		return "According to Ptolemy's dignity system, planets have varying energy expression in different signs"
	case "retrograde":
		return "From Earth's perspective, planet appears to move backward, symbolizing introspection and reassessment"
	case "aspectPhase":
		return "Angular relationships between planets determine how energies interact"
	case "lunarPhase":
		return "Lunar cycle influences mood, body rhythms, and daily affairs"
	case "planetaryHour":
		return "Classical astrology's planetary hour system, each period ruled by a different planet"
	default:
		return ""
	}
}

// getTimeLevelLabel 获取时间级别标签
func getTimeLevelLabel(level string) string {
	labels := map[string]string{
		"yearly":  "Long-term",
		"monthly": "Monthly",
		"weekly":  "This Week",
		"daily":   "Today",
		"hourly":  "Current Hour",
	}
	if l, ok := labels[level]; ok {
		return l
	}
	return level
}

// getValidPeriod 获取有效期描述
func getValidPeriod(level string) string {
	periods := map[string]string{
		"yearly":  "Lasts all year",
		"monthly": "Lasts all month",
		"weekly":  "Valid this week",
		"daily":   "Valid today",
		"hourly":  "Valid this hour",
	}
	if p, ok := periods[level]; ok {
		return p
	}
	return "Ongoing"
}

// getScoreLevel 获取分数等级
func getScoreLevel(score float64) (string, string) {
	if score >= 85 {
		return "Excellent", "🌟"
	} else if score >= 70 {
		return "Good", "✨"
	} else if score >= 55 {
		return "Stable", "💫"
	} else if score >= 40 {
		return "Caution", "🌙"
	}
	return "Challenge", "⚡"
}

// getGranularityLabel 获取粒度标签
func getGranularityLabel(g string) string {
	labels := map[string]string{
		"hour":  "Hourly",
		"day":   "Daily",
		"week":  "Weekly",
		"month": "Monthly",
		"year":  "Yearly",
	}
	if l, ok := labels[g]; ok {
		return l
	}
	return g
}

// getDimensionLabel 获取维度标签
func getDimensionLabel(d string) string {
	labels := map[string]string{
		"career":       "Career Fortune",
		"relationship": "Relationship Fortune",
		"health":       "Health Fortune",
		"finance":      "Finance Fortune",
		"spiritual":    "Spiritual Fortune",
		"overall":      "Overall Fortune",
	}
	if l, ok := labels[d]; ok {
		return l
	}
	return d
}

// formatTimeLabel 格式化时间标签
func formatTimeLabel(t time.Time, granularity string) string {
	switch granularity {
	case "hour":
		return fmt.Sprintf("%s %d, %d %02d:00", t.Month().String()[:3], t.Day(), t.Year(), t.Hour())
	case "day":
		return fmt.Sprintf("%s %d, %d", t.Month().String()[:3], t.Day(), t.Year())
	case "week":
		return fmt.Sprintf("%d Week %d", t.Year(), getWeekNumber(t))
	case "month":
		return fmt.Sprintf("%s %d", t.Month().String(), t.Year())
	case "year":
		return fmt.Sprintf("%d", t.Year())
	default:
		return t.Format("2006-01-02 15:04")
	}
}

// getWeekNumber 获取周数
func getWeekNumber(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}

// buildExplanationDetail 构建解释详情
func buildExplanationDetail(dim DimensionBreakdown) ScoreExplanationDetail {
	return ScoreExplanationDetail{
		Formula:      "Base Score + Planetary Aspects + Other Factors → Raw Score → Normalization → Final Score",
		BaseScore:    dim.BaseScore,
		AspectEffect: round2(dim.AspectScore),
		FactorEffect: round2(dim.FactorScore),
		RawScore:     round2(dim.RawScore),
		Description: fmt.Sprintf(
			"Base score %.0f, planetary aspects contribute %+.1f, other factors contribute %+.1f, "+
				"raw score %.1f normalized to final score %.1f",
			dim.BaseScore, dim.AspectScore, dim.FactorScore, dim.RawScore, dim.FinalScore),
	}
}

// buildSummary 构建总结
func buildSummary(dimension string, score float64, factors []AstronomicalFactor) string {
	dimLabel := getDimensionLabel(dimension)
	level, _ := getScoreLevel(score)
	
	// 找出最强的正负因素
	var strongestPositive, strongestNegative string
	for _, f := range factors {
		if f.IsPositive && strongestPositive == "" && f.Intensity == "Strong" {
			strongestPositive = f.Name
		}
		if !f.IsPositive && strongestNegative == "" && f.Intensity == "Strong" {
			strongestNegative = f.Name
		}
	}
	
	summary := fmt.Sprintf("Your %s status is '%s' (%.0f points). ", dimLabel, level, score)
	
	if strongestPositive != "" {
		summary += fmt.Sprintf("'%s' brings you positive energy. ", strongestPositive)
	}
	if strongestNegative != "" {
		summary += fmt.Sprintf("Pay attention to challenges from '%s'. ", strongestNegative)
	}
	
	return summary
}

// buildSuggestions 构建建议
func buildSuggestions(dimension string, score float64, factors []AstronomicalFactor) []string {
	var suggestions []string
	
	// 根据分数给出基础建议
	if score >= 80 {
		suggestions = append(suggestions, "Excellent fortune, good time to push forward on important matters")
	} else if score >= 60 {
		suggestions = append(suggestions, "Good fortune, maintain positive attitude and keep going")
	} else if score >= 40 {
		suggestions = append(suggestions, "Stable fortune, suitable for steady progress")
	} else {
		suggestions = append(suggestions, "Challenging period, proceed with caution and wait for better timing")
	}
	
	// 根据因素给出具体建议
	for _, f := range factors {
		if f.Type == "retrograde" && !f.IsPositive {
			suggestions = append(suggestions, "Planetary retrograde present, think twice before major decisions")
			break
		}
	}
	
	// 根据维度给出特定建议
	switch dimension {
	case "health":
		if score < 60 {
			suggestions = append(suggestions, "Get adequate rest, avoid overexertion")
		}
	case "finance":
		if score < 60 {
			suggestions = append(suggestions, "Be cautious with finances, avoid impulsive spending")
		}
	case "relationship":
		if score >= 70 {
			suggestions = append(suggestions, "Good time for socializing and deepening relationships")
		}
	}
	
	return suggestions
}

// buildOverallSummary 构建综合总结
func buildOverallSummary(score float64, factors []AstronomicalFactor) string {
	level, _ := getScoreLevel(score)
	
	positiveCount := 0
	negativeCount := 0
	for _, f := range factors {
		if f.IsPositive {
			positiveCount++
		} else {
			negativeCount++
		}
	}
	
	return fmt.Sprintf(
		"Your overall fortune is '%s' (%.0f points). Currently %d positive factors and %d factors requiring attention are influencing you.",
		level, score, positiveCount, negativeCount)
}

// buildOverallSuggestions 构建综合建议
func buildOverallSuggestions(score float64) []string {
	if score >= 80 {
		return []string{
			"Excellent overall fortune, great time to push forward on important matters",
			"Maintain positive attitude and seize current opportunities",
		}
	} else if score >= 60 {
		return []string{
			"Good fortune, steadily advance various affairs",
			"Focus on high-scoring dimensions, leverage your strengths",
		}
	} else if score >= 40 {
		return []string{
			"Stable fortune, suitable for handling daily affairs",
			"Avoid major decisions, wait for better timing",
		}
	}
	return []string{
		"Multiple challenging factors present, proceed with caution",
		"Focus on essential matters, maintain patience",
	}
}

// round2 保留两位小数
func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

// abs 绝对值
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

