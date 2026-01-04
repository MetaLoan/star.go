package astro

import (
	"math"
	"math/rand"
	"star/models"
)

// ==================== 视觉抖动系统 ====================
// 仅用于显示层，不影响实际计算
// 使用确定性随机数，相同输入产生相同抖动

// JitterConfig 抖动配置
type JitterConfig struct {
	Enabled   bool    `json:"enabled"`   // 是否启用抖动
	Magnitude float64 `json:"magnitude"` // 抖动幅度（默认0.5）
	Seed      int64   `json:"seed"`      // 随机种子（0表示使用时间戳）
}

// DefaultJitterConfig 默认抖动配置
var DefaultJitterConfig = JitterConfig{
	Enabled:   true,
	Magnitude: 0.5,
	Seed:      0,
}

// ApplyJitter 对单个分数应用抖动（仅显示用）
// seed 用于确定性随机，相同seed产生相同抖动
func ApplyJitter(score float64, seed int64) float64 {
	if !DefaultJitterConfig.Enabled {
		return score
	}

	// 使用确定性随机数
	r := rand.New(rand.NewSource(seed))

	// 抖动范围：±magnitude
	jitter := (r.Float64() - 0.5) * 2 * DefaultJitterConfig.Magnitude

	// 应用抖动
	jittered := score + jitter

	// 确保在 0-100 范围内
	jittered = math.Max(0, math.Min(100, jittered))

	// 保留4位小数
	return math.Round(jittered*10000) / 10000
}

// ApplyJitterToResult 对整个分数结果应用抖动
func ApplyJitterToResult(result *ScoreResult, baseSeed int64) *ScoreResult {
	if result == nil || !DefaultJitterConfig.Enabled {
		return result
	}

	// 创建新的结果，避免修改原始数据
	jittered := &ScoreResult{
		Overall:    ApplyJitter(result.Overall, baseSeed),
		Timestamp:  result.Timestamp,
		BaseScores: result.BaseScores,
		Factors:    result.Factors,
	}

	// 对各维度应用不同的seed
	jittered.Dimensions = models.DimensionScoresV2{
		Career:       ApplyJitter(result.Dimensions.Career, baseSeed+1),
		Relationship: ApplyJitter(result.Dimensions.Relationship, baseSeed+2),
		Health:       ApplyJitter(result.Dimensions.Health, baseSeed+3),
		Finance:      ApplyJitter(result.Dimensions.Finance, baseSeed+4),
		Spiritual:    ApplyJitter(result.Dimensions.Spiritual, baseSeed+5),
	}

	return jittered
}

// ApplyJitterToAggregated 对聚合分数应用抖动
func ApplyJitterToAggregated(score *AggregatedScore, baseSeed int64) *AggregatedScore {
	if score == nil || !DefaultJitterConfig.Enabled {
		return score
	}

	// 创建新的结果
	jittered := &AggregatedScore{
		Overall:     ApplyJitter(score.Overall, baseSeed),
		Features:    score.Features,
		StartTime:   score.StartTime,
		EndTime:     score.EndTime,
		Granularity: score.Granularity,
	}

	// 对各维度应用不同的seed
	jittered.Dimensions = models.DimensionScoresV2{
		Career:       ApplyJitter(score.Dimensions.Career, baseSeed+1),
		Relationship: ApplyJitter(score.Dimensions.Relationship, baseSeed+2),
		Health:       ApplyJitter(score.Dimensions.Health, baseSeed+3),
		Finance:      ApplyJitter(score.Dimensions.Finance, baseSeed+4),
		Spiritual:    ApplyJitter(score.Dimensions.Spiritual, baseSeed+5),
	}

	// 递归处理子分数
	if len(score.SubScores) > 0 {
		jittered.SubScores = make([]AggregatedScore, len(score.SubScores))
		for i, sub := range score.SubScores {
			subJittered := ApplyJitterToAggregated(&sub, baseSeed+int64(i)*10)
			jittered.SubScores[i] = *subJittered
		}
	}

	return jittered
}

// GenerateSeedFromTime 从时间生成种子
// 确保相同时间产生相同抖动
func GenerateSeedFromTime(year, month, day, hour int) int64 {
	return int64(year*10000000 + month*100000 + day*1000 + hour)
}

// SmoothTransition 平滑过渡函数
// 用于两个时间点之间的分数平滑插值
func SmoothTransition(score1, score2 float64, progress float64) float64 {
	// 使用缓动函数（ease-in-out）
	// smoothstep: 3*t^2 - 2*t^3
	t := progress
	smoothT := t * t * (3 - 2*t)

	// 线性插值
	result := score1 + (score2-score1)*smoothT

	return math.Round(result*10000) / 10000
}

// UpdateJitterConfig 更新抖动配置
func UpdateJitterConfig(config JitterConfig) {
	DefaultJitterConfig = config
}

