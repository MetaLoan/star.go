package astro

import (
	"star/models"
	"time"
)

// ==================== 旧版兼容层 ====================
// CalculateInfluenceFactors 保持向后兼容，内部调用新版实现

// CalculateInfluenceFactors 计算影响因子（兼容旧版API）
// Deprecated: 请使用 CalculateInfluenceFactorsV2
func CalculateInfluenceFactors(chart *models.NatalChart, date time.Time, transitPositions []models.PlanetPosition) *models.FactorResult {
	return CalculateInfluenceFactorsV2(chart, date, transitPositions)
}

// UpdateFactorWeights 更新因子权重（供运营调整）
func UpdateFactorWeights(newWeights models.FactorWeights) {
	DefaultFactorWeights = newWeights
}

// GetCurrentFactorWeights 获取当前因子权重配置
func GetCurrentFactorWeights() models.FactorWeights {
	return DefaultFactorWeights
}

// UpdateDimensionWeights 更新维度权重（供运营调整）
func UpdateDimensionWeights(newWeights models.DimensionWeights) {
	DefaultDimensionWeights = newWeights
}

// GetCurrentDimensionWeights 获取当前维度权重配置
func GetCurrentDimensionWeights() models.DimensionWeights {
	return DefaultDimensionWeights
}
