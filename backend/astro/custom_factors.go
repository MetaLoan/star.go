package astro

import (
	"fmt"
	"regexp"
	"star/models"
	"strconv"
	"strings"
	"time"
)

// ==================== 自定义因子系统 ====================
// 格式：AddScore=(2*healthScore,2.5,202501171230)
// 含义：从202501171230开始，健康值×2，维持2.5小时

// CustomFactorDefinition 自定义因子定义
type CustomFactorDefinition struct {
	ID          string    `json:"id"`
	Operation   string    `json:"operation"`   // AddScore, SubScore, MulScore, SetScore
	Dimension   string    `json:"dimension"`   // career, relationship, health, finance, spiritual, overall
	Value       float64   `json:"value"`       // 操作值
	Duration    float64   `json:"duration"`    // 持续时长（小时）
	StartTime   time.Time `json:"startTime"`   // 开始时间
	EndTime     time.Time `json:"endTime"`     // 结束时间（自动计算）
	Description string    `json:"description"` // 描述
}

// CustomFactorStore 自定义因子存储（内存）
var CustomFactorStore = make(map[string][]CustomFactorDefinition) // userID -> factors

// ParseCustomFactor 解析自定义因子字符串
// 格式：AddScore=(2*healthScore,2.5,202501171230)
func ParseCustomFactor(input string) (*CustomFactorDefinition, error) {
	// 正则匹配：Operation=(value*dimension,duration,startTime)
	// 或：Operation=(value,duration,startTime) 表示对 overall
	re := regexp.MustCompile(`^(\w+)=\(([\d.]+)\*?(\w+)?,\s*([\d.]+),\s*(\d{12})\)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return nil, fmt.Errorf("无效的自定义因子格式: %s", input)
	}

	operation := matches[1]
	valueStr := matches[2]
	dimension := matches[3]
	durationStr := matches[4]
	timeStr := matches[5]

	// 解析值
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的值: %s", valueStr)
	}

	// 解析持续时长
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的持续时长: %s", durationStr)
	}

	// 解析时间 YYYYMMDDHHmm
	startTime, err := parseCustomTime(timeStr)
	if err != nil {
		return nil, err
	}

	// 默认维度为 overall
	if dimension == "" {
		dimension = "overall"
	} else {
		dimension = normalizeDimensionName(dimension)
	}

	// 验证操作类型
	validOps := map[string]bool{
		"AddScore": true,
		"SubScore": true,
		"MulScore": true,
		"SetScore": true,
	}
	if !validOps[operation] {
		return nil, fmt.Errorf("无效的操作类型: %s (支持: AddScore, SubScore, MulScore, SetScore)", operation)
	}

	endTime := startTime.Add(time.Duration(duration * float64(time.Hour)))

	return &CustomFactorDefinition{
		ID:        fmt.Sprintf("custom_%d", startTime.Unix()),
		Operation: operation,
		Dimension: dimension,
		Value:     value,
		Duration:  duration,
		StartTime: startTime,
		EndTime:   endTime,
		Description: fmt.Sprintf("%s %s %.2f, 持续%.1f小时",
			operation, dimension, value, duration),
	}, nil
}

// parseCustomTime 解析自定义时间格式 YYYYMMDDHHmm
func parseCustomTime(timeStr string) (time.Time, error) {
	if len(timeStr) != 12 {
		return time.Time{}, fmt.Errorf("时间格式错误，应为12位: YYYYMMDDHHmm")
	}

	year, _ := strconv.Atoi(timeStr[0:4])
	month, _ := strconv.Atoi(timeStr[4:6])
	day, _ := strconv.Atoi(timeStr[6:8])
	hour, _ := strconv.Atoi(timeStr[8:10])
	minute, _ := strconv.Atoi(timeStr[10:12])

	// 验证范围
	if month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("月份无效: %d", month)
	}
	if day < 1 || day > 31 {
		return time.Time{}, fmt.Errorf("日期无效: %d", day)
	}
	if hour < 0 || hour > 23 {
		return time.Time{}, fmt.Errorf("小时无效: %d", hour)
	}
	if minute < 0 || minute > 59 {
		return time.Time{}, fmt.Errorf("分钟无效: %d", minute)
	}

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC), nil
}

// normalizeDimensionName 标准化维度名称
func normalizeDimensionName(name string) string {
	name = strings.ToLower(name)
	switch name {
	case "career", "careerscore":
		return "career"
	case "relationship", "relationshipscore":
		return "relationship"
	case "health", "healthscore":
		return "health"
	case "finance", "financescore":
		return "finance"
	case "spiritual", "spiritualscore":
		return "spiritual"
	case "overall", "overallscore":
		return "overall"
	default:
		return "overall"
	}
}

// AddCustomFactor 添加自定义因子
func AddCustomFactor(userID string, factorDef string) (*CustomFactorDefinition, error) {
	factor, err := ParseCustomFactor(factorDef)
	if err != nil {
		return nil, err
	}

	CustomFactorStore[userID] = append(CustomFactorStore[userID], *factor)
	return factor, nil
}

// GetActiveCustomFactors 获取某时刻的活跃自定义因子
func GetActiveCustomFactors(userID string, t time.Time) []CustomFactorDefinition {
	var active []CustomFactorDefinition

	factors := CustomFactorStore[userID]
	for _, f := range factors {
		if t.After(f.StartTime) && t.Before(f.EndTime) {
			active = append(active, f)
		}
	}

	return active
}

// ApplyCustomFactors 应用自定义因子到分数
func ApplyCustomFactors(userID string, scores *models.DimensionScoresV2, overall *float64, t time.Time) {
	factors := GetActiveCustomFactors(userID, t)

	for _, f := range factors {
		// 计算当前强度（正弦曲线）
		lifecycle := &models.FactorLifecycle{
			StartTime: f.StartTime,
			PeakTime:  f.StartTime.Add(time.Duration(f.Duration/2) * time.Hour),
			EndTime:   f.EndTime,
			Duration:  f.Duration,
		}
		strength := CalculateFactorStrength(lifecycle, t)

		// 根据强度调整值
		adjustedValue := f.Value * strength

		// 应用到对应维度
		switch f.Dimension {
		case "career":
			applyOperation(&scores.Career, f.Operation, adjustedValue)
		case "relationship":
			applyOperation(&scores.Relationship, f.Operation, adjustedValue)
		case "health":
			applyOperation(&scores.Health, f.Operation, adjustedValue)
		case "finance":
			applyOperation(&scores.Finance, f.Operation, adjustedValue)
		case "spiritual":
			applyOperation(&scores.Spiritual, f.Operation, adjustedValue)
		case "overall":
			applyOperation(overall, f.Operation, adjustedValue)
		}
	}
}

// applyOperation 应用操作到分数
func applyOperation(score *float64, operation string, value float64) {
	switch operation {
	case "AddScore":
		*score += value
	case "SubScore":
		*score -= value
	case "MulScore":
		*score *= value
	case "SetScore":
		*score = value
	}

	// 确保分数在0-100范围内
	if *score < 0 {
		*score = 0
	}
	if *score > 100 {
		*score = 100
	}
}

// ConvertCustomFactorToInfluenceFactor 将自定义因子转换为标准影响因子
func ConvertCustomFactorToInfluenceFactor(cfd CustomFactorDefinition, t time.Time) models.InfluenceFactor {
	// 计算生命周期
	lifecycle := &models.FactorLifecycle{
		StartTime: cfd.StartTime,
		PeakTime:  cfd.StartTime.Add(time.Duration(cfd.Duration/2) * time.Hour),
		EndTime:   cfd.EndTime,
		Duration:  cfd.Duration,
	}

	// 计算当前强度
	strength := CalculateFactorStrength(lifecycle, t)

	// 确定维度影响
	var impact models.DimensionImpact
	switch cfd.Dimension {
	case "career":
		impact = models.DimensionImpact{Career: 1.0}
	case "relationship":
		impact = models.DimensionImpact{Relationship: 1.0}
	case "health":
		impact = models.DimensionImpact{Health: 1.0}
	case "finance":
		impact = models.DimensionImpact{Finance: 1.0}
	case "spiritual":
		impact = models.DimensionImpact{Spiritual: 1.0}
	default:
		// overall: 平均分配
		impact = models.DimensionImpact{
			Career:       0.2,
			Relationship: 0.2,
			Health:       0.2,
			Finance:      0.2,
			Spiritual:    0.2,
		}
	}

	return models.InfluenceFactor{
		ID:              cfd.ID,
		Type:            models.FactorCustom,
		Name:            cfd.Description,
		Description:     cfd.Description,
		TimeLevel:       models.TimeLevelHourly,
		Lifecycle:       lifecycle,
		BaseValue:       cfd.Value,
		Weight:          DefaultFactorWeights.Custom,
		CurrentStrength: strength,
		Adjustment:      cfd.Value * strength * DefaultFactorWeights.Custom,
		DimensionImpact: impact,
		IsPositive:      cfd.Value > 0,
		AstroReason:     "自定义因子：用户手动设置",
	}
}

// ClearCustomFactors 清除用户的自定义因子
func ClearCustomFactors(userID string) {
	delete(CustomFactorStore, userID)
}

// GetAllCustomFactors 获取用户的所有自定义因子
func GetAllCustomFactors(userID string) []CustomFactorDefinition {
	return CustomFactorStore[userID]
}

