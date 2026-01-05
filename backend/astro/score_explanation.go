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
		Dimension:   "综合运势",
		Score:       score,
		ScoreLevel:  scoreLevel,
		ScoreEmoji:  scoreEmoji,
		Explanation: ScoreExplanationDetail{
			Formula:     "五维度加权平均 → 标准化 → 最终分",
			BaseScore:   50,
			RawScore:    breakdown.OverallRaw,
			Description: "综合运势由事业、关系、健康、财务、灵性五个维度各占20%加权计算得出",
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
		return "月亮空亡期"
	case "custom":
		return "个人设定: " + f.Name
	default:
		return f.Name
	}
}

// getFactorCategory 获取因素类别
func getFactorCategory(factorType string) string {
	switch factorType {
	case "dignity":
		return "行星状态"
	case "retrograde":
		return "行星逆行"
	case "aspectPhase":
		return "行星相位"
	case "lunarPhase":
		return "月相"
	case "planetaryHour":
		return "行星时"
	case "voidOfCourse":
		return "月亮空亡"
	case "custom":
		return "个人因素"
	default:
		return "其他"
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
		return "增强"
	}
	return "减弱"
}

// getIntensity 获取强度
func getIntensity(value float64) string {
	absVal := abs(value)
	if absVal >= 0.5 {
		return "强"
	} else if absVal >= 0.2 {
		return "中"
	}
	return "弱"
}

// getFactorDescription 获取因素描述
func getFactorDescription(f FactorContribution) string {
	switch f.Type {
	case "dignity":
		if f.IsPositive {
			return "行星处于有利位置，能量得到提升"
		}
		return "行星处于不利位置，能量有所削弱"
	case "retrograde":
		return "行星逆行期间，相关领域需要回顾和反思"
	case "aspectPhase":
		if f.IsPositive {
			return "行星之间形成和谐角度，带来积极能量"
		}
		return "行星之间形成紧张角度，带来挑战"
	case "lunarPhase":
		return getLunarPhaseDescription(f.Name)
	case "planetaryHour":
		return "当前时段的行星能量影响"
	case "custom":
		return "个人设定的调整因素"
	default:
		return f.Description
	}
}

// getLunarPhaseDescription 获取月相描述
func getLunarPhaseDescription(name string) string {
	descriptions := map[string]string{
		"新月期": "新月时期，适合设立新目标和播种意图",
		"上弦月期": "月亮渐盈，能量上升，适合采取行动",
		"盈凸月期": "接近满月，准备收获的时期",
		"满月期": "满月高峰，情绪和能量达到顶点，适合展现成果",
		"亏凸月期": "满月后，适合分享和传播",
		"下弦月期": "能量下降，适合释放和放手",
		"残月期": "月亮即将隐没，适合休息和反思",
	}
	if desc, ok := descriptions[name]; ok {
		return desc
	}
	return "月相影响当日能量和情绪"
}

// getAstroExplanation 获取占星学解释
func getAstroExplanation(f FactorContribution) string {
	switch f.Type {
	case "dignity":
		return "根据托勒密尊贵度系统，行星在特定星座的能量表达有强弱之分"
	case "retrograde":
		return "从地球视角观察，行星呈现逆向移动，象征内省和重新评估"
	case "aspectPhase":
		return "行星之间的角度关系决定了能量的互动方式"
	case "lunarPhase":
		return "月亮周期影响情绪、身体节律和日常事务"
	case "planetaryHour":
		return "古典占星的行星时系统，每个时段由不同行星主管"
	default:
		return ""
	}
}

// getTimeLevelLabel 获取时间级别标签
func getTimeLevelLabel(level string) string {
	labels := map[string]string{
		"yearly":  "长期",
		"monthly": "月度",
		"weekly":  "本周",
		"daily":   "今日",
		"hourly":  "当前小时",
	}
	if l, ok := labels[level]; ok {
		return l
	}
	return level
}

// getValidPeriod 获取有效期描述
func getValidPeriod(level string) string {
	periods := map[string]string{
		"yearly":  "持续全年",
		"monthly": "持续整月",
		"weekly":  "本周有效",
		"daily":   "今日有效",
		"hourly":  "当前小时有效",
	}
	if p, ok := periods[level]; ok {
		return p
	}
	return "持续中"
}

// getScoreLevel 获取分数等级
func getScoreLevel(score float64) (string, string) {
	if score >= 85 {
		return "极佳", "🌟"
	} else if score >= 70 {
		return "良好", "✨"
	} else if score >= 55 {
		return "平稳", "💫"
	} else if score >= 40 {
		return "需注意", "🌙"
	}
	return "挑战", "⚡"
}

// getGranularityLabel 获取粒度标签
func getGranularityLabel(g string) string {
	labels := map[string]string{
		"hour":  "小时",
		"day":   "日",
		"week":  "周",
		"month": "月",
		"year":  "年",
	}
	if l, ok := labels[g]; ok {
		return l
	}
	return g
}

// getDimensionLabel 获取维度标签
func getDimensionLabel(d string) string {
	labels := map[string]string{
		"career":       "事业运",
		"relationship": "关系运",
		"health":       "健康运",
		"finance":      "财务运",
		"spiritual":    "灵性运",
		"overall":      "综合运势",
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
		return fmt.Sprintf("%d年%d月%d日 %02d:00", t.Year(), t.Month(), t.Day(), t.Hour())
	case "day":
		return fmt.Sprintf("%d年%d月%d日", t.Year(), t.Month(), t.Day())
	case "week":
		return fmt.Sprintf("%d年第%d周", t.Year(), getWeekNumber(t))
	case "month":
		return fmt.Sprintf("%d年%d月", t.Year(), t.Month())
	case "year":
		return fmt.Sprintf("%d年", t.Year())
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
		Formula:      "基础分 + 行星相位影响 + 其他天文因素 → 原始分 → 标准化 → 最终分",
		BaseScore:    dim.BaseScore,
		AspectEffect: round2(dim.AspectScore),
		FactorEffect: round2(dim.FactorScore),
		RawScore:     round2(dim.RawScore),
		Description: fmt.Sprintf(
			"基础分%.0f，行星相位带来%+.1f的影响，其他天文因素带来%+.1f的影响，"+
				"原始分%.1f经过标准化后得到最终分%.1f",
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
		if f.IsPositive && strongestPositive == "" && f.Intensity == "强" {
			strongestPositive = f.Name
		}
		if !f.IsPositive && strongestNegative == "" && f.Intensity == "强" {
			strongestNegative = f.Name
		}
	}
	
	summary := fmt.Sprintf("您的%s当前状态为「%s」(%.0f分)。", dimLabel, level, score)
	
	if strongestPositive != "" {
		summary += fmt.Sprintf("「%s」为您带来积极能量。", strongestPositive)
	}
	if strongestNegative != "" {
		summary += fmt.Sprintf("需注意「%s」带来的挑战。", strongestNegative)
	}
	
	return summary
}

// buildSuggestions 构建建议
func buildSuggestions(dimension string, score float64, factors []AstronomicalFactor) []string {
	var suggestions []string
	
	// 根据分数给出基础建议
	if score >= 80 {
		suggestions = append(suggestions, "运势极佳，可以大胆推进重要事项")
	} else if score >= 60 {
		suggestions = append(suggestions, "运势良好，保持积极心态继续努力")
	} else if score >= 40 {
		suggestions = append(suggestions, "运势平稳，适合稳扎稳打")
	} else {
		suggestions = append(suggestions, "运势有挑战，建议谨慎行事，等待时机")
	}
	
	// 根据因素给出具体建议
	for _, f := range factors {
		if f.Type == "retrograde" && !f.IsPositive {
			suggestions = append(suggestions, "有行星逆行，重要决定建议三思后行")
			break
		}
	}
	
	// 根据维度给出特定建议
	switch dimension {
	case "health":
		if score < 60 {
			suggestions = append(suggestions, "注意休息，避免过度劳累")
		}
	case "finance":
		if score < 60 {
			suggestions = append(suggestions, "财务方面保持谨慎，避免冲动消费")
		}
	case "relationship":
		if score >= 70 {
			suggestions = append(suggestions, "适合社交和增进感情")
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
		"您的综合运势为「%s」(%.0f分)。当前有%d个积极因素和%d个需注意的因素在影响您。",
		level, score, positiveCount, negativeCount)
}

// buildOverallSuggestions 构建综合建议
func buildOverallSuggestions(score float64) []string {
	if score >= 80 {
		return []string{
			"整体运势极佳，是推进重要事项的好时机",
			"保持积极心态，把握当下机遇",
		}
	} else if score >= 60 {
		return []string{
			"运势良好，稳步推进各项事务",
			"关注高分维度，发挥优势领域",
		}
	} else if score >= 40 {
		return []string{
			"运势平稳，适合处理日常事务",
			"避免做出重大决定，等待更好时机",
		}
	}
	return []string{
		"当前有较多挑战因素，建议谨慎行事",
		"专注于必要事务，保持耐心",
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

