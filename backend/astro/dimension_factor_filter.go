package astro

import "star/models"

/*
维度因子过滤器
根据因子的天文学属性（行星、宫位、类型）决定是否影响某个维度
目标：让不同维度受到不同的因子集合影响，打破"共享因子"导致的同步波动
*/

// ShouldFactorAffectDimension 判断因子是否应该影响某个维度
func ShouldFactorAffectDimension(factor *models.InfluenceFactor, dimension string) bool {
	// 1. 基于因子类型的过滤规则
	switch factor.Type {
	case models.FactorPlanetaryHour:
		// 行星时：主要影响当下行动，对长期维度（灵性）影响较小
		return dimension != "spiritual"
		
	case models.FactorLunarPhase:
		// 月相：主要影响情绪和灵性（仅这两个）
		return dimension == "relationship" || dimension == "spiritual"
		
	case models.FactorVoidOfCourse:
		// 月空亡：主要影响决策和行动，灵性维度反而有利
		return dimension == "career" || dimension == "finance" || dimension == "spiritual"
		
	case models.FactorRetrograde:
		// 逆行：根据行星类型决定
		return shouldRetrogradeAffectDimension(factor.SourcePlanet, dimension)
		
	case models.FactorAspectPhase:
		// 相位：根据涉及的行星决定
		return shouldAspectAffectDimension(factor.SourcePlanet, dimension)
		
	case models.FactorDignity:
		// 尊贵度：影响所有维度（但权重不同）
		return true
		
	case models.FactorProfectionLord:
		// 年主星：影响所有维度
		return true
		
	default:
		return true
	}
}

// shouldRetrogradeAffectDimension 判断逆行是否影响某个维度
func shouldRetrogradeAffectDimension(planet models.PlanetID, dimension string) bool {
	switch planet {
	case models.Mercury:
		// 水星逆行：沟通、合同、学习
		return dimension == "career" || dimension == "relationship"
		
	case models.Venus:
		// 金星逆行：关系、金钱、价值观
		return dimension == "relationship" || dimension == "finance"
		
	case models.Mars:
		// 火星逆行：行动、冲突、健康
		return dimension == "career" || dimension == "health"
		
	case models.Jupiter:
		// 木星逆行：信仰、财富、成长（仅这两个）
		return dimension == "finance" || dimension == "spiritual"
		
	case models.Saturn:
		// 土星逆行：责任、结构、健康（仅这三个）
		return dimension == "career" || dimension == "health" || dimension == "spiritual"
		
	case models.Uranus:
		// 天王星逆行：创新、突破（仅灵性）
		return dimension == "spiritual"
		
	case models.NorthNode:
		// 北交点逆行：命运、业力（仅灵性）
		return dimension == "spiritual"
		
	default:
		// 其他外行星逆行：主要影响灵性
		return dimension == "spiritual"
	}
}

// shouldAspectAffectDimension 判断相位是否影响某个维度
func shouldAspectAffectDimension(planet models.PlanetID, dimension string) bool {
	switch planet {
	case models.Sun, models.Mars, models.Saturn:
		// 阳性行星：事业、行动、责任
		return dimension == "career" || dimension == "health"
		
	case models.Moon, models.Venus:
		// 阴性行星：情感、关系
		return dimension == "relationship" || dimension == "health" || dimension == "spiritual"
		
	case models.Jupiter:
		// 木星：财富、机会
		return dimension == "career" || dimension == "finance"
		
	case models.Mercury:
		// 水星：沟通、学习
		return dimension == "career" || dimension == "relationship"
		
	default:
		// 外行星：深层转化
		return dimension == "spiritual" || dimension == "health"
	}
}

// ApplyDimensionFilter 应用过滤规则到因子计算
// 如果因子不应该影响某个维度，返回 0 作为影响权重
func ApplyDimensionFilter(factor *models.InfluenceFactor, dimension string, originalImpact float64) float64 {
	if !ShouldFactorAffectDimension(factor, dimension) {
		return 0.0 // 完全不影响该维度
	}
	return originalImpact // 保持原始影响
}
