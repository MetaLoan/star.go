package astro

import "star/models"

// ==================== 行星-维度影响矩阵 ====================
// 基于占星学理论：行星的自然象征意义决定其对各生活领域的影响比例
// 参考文献：Liz Greene, Stephen Arroyo 等现代心理占星学著作
// 每行总和 = 1.0

// PlanetDimensionMapping 行星-维度影响映射表
var PlanetDimensionMapping = map[models.PlanetID]models.DimensionImpact{
	// 太阳：自我实现、生命力核心
	// 事业(10宫)0.35, 健康(1宫活力)0.25, 灵性(自我意识)0.15
	models.Sun: {
		Career:       0.35,
		Relationship: 0.15,
		Health:       0.25,
		Finance:      0.10,
		Spiritual:    0.15,
	},

	// 月亮：情感需求、身体节律、潜意识
	// 关系(情感)0.30, 健康(身体节律)0.30, 灵性(潜意识)0.20
	models.Moon: {
		Career:       0.10,
		Relationship: 0.30,
		Health:       0.30,
		Finance:      0.10,
		Spiritual:    0.20,
	},

	// 水星：思维沟通、商业交易、学习
	// 事业(工作技能)0.30, 财务(商业)0.25, 关系(沟通)0.25
	models.Mercury: {
		Career:       0.30,
		Relationship: 0.25,
		Health:       0.10,
		Finance:      0.25,
		Spiritual:    0.10,
	},

	// 金星：爱与美、价值观、享乐
	// 关系(爱情)0.45, 财务(价值/金钱)0.25
	models.Venus: {
		Career:       0.10,
		Relationship: 0.45,
		Health:       0.10,
		Finance:      0.25,
		Spiritual:    0.10,
	},

	// 火星：行动力、身体能量、竞争
	// 健康(体力)0.40, 事业(行动力)0.30
	models.Mars: {
		Career:       0.30,
		Relationship: 0.10,
		Health:       0.40,
		Finance:      0.10,
		Spiritual:    0.10,
	},

	// 木星：扩张、机遇、信仰、财富
	// 财务(扩张)0.30, 事业(机遇)0.25, 灵性(信仰)0.20
	models.Jupiter: {
		Career:       0.25,
		Relationship: 0.15,
		Health:       0.10,
		Finance:      0.30,
		Spiritual:    0.20,
	},

	// 土星：责任、结构、限制、时间
	// 事业(责任/成就)0.40, 健康(骨骼/慢性)0.20, 财务(保守)0.20
	models.Saturn: {
		Career:       0.40,
		Relationship: 0.10,
		Health:       0.20,
		Finance:      0.20,
		Spiritual:    0.10,
	},

	// 天王星：突破、觉醒、创新、科技
	// 灵性(觉醒)0.50, 事业(创新)0.20
	models.Uranus: {
		Career:       0.20,
		Relationship: 0.10,
		Health:       0.10,
		Finance:      0.10,
		Spiritual:    0.50,
	},

	// 海王星：灵感、幻想、灵性、迷惑
	// 灵性(灵感)0.50, 关系(理想化)0.20
	models.Neptune: {
		Career:       0.10,
		Relationship: 0.20,
		Health:       0.10,
		Finance:      0.10,
		Spiritual:    0.50,
	},

	// 冥王星：转化、权力、深层资源
	// 财务(深层资源)0.25, 灵性(转化)0.20, 事业(权力)0.20, 健康(再生)0.20
	models.Pluto: {
		Career:       0.20,
		Relationship: 0.15,
		Health:       0.20,
		Finance:      0.25,
		Spiritual:    0.20,
	},

	// 北交点：命运方向、灵魂成长
	// 灵性(灵魂使命)0.45, 关系(业力关系)0.20
	models.NorthNode: {
		Career:       0.15,
		Relationship: 0.20,
		Health:       0.10,
		Finance:      0.10,
		Spiritual:    0.45,
	},

	// 凯龙：伤痛与疗愈
	// 灵性(疗愈)0.35, 健康(身心伤痛)0.30, 关系(疗愈关系)0.20
	models.Chiron: {
		Career:       0.10,
		Relationship: 0.20,
		Health:       0.30,
		Finance:      0.05,
		Spiritual:    0.35,
	},
}

// GetPlanetDimensionImpact 获取行星的维度影响分配
func GetPlanetDimensionImpact(planet models.PlanetID) models.DimensionImpact {
	if impact, ok := PlanetDimensionMapping[planet]; ok {
		return impact
	}
	// 默认平均分配
	return models.DimensionImpact{
		Career:       0.20,
		Relationship: 0.20,
		Health:       0.20,
		Finance:      0.20,
		Spiritual:    0.20,
	}
}

// ==================== 默认维度权重 ====================

// DefaultDimensionWeights 默认维度权重
var DefaultDimensionWeights = models.DimensionWeights{
	Career:       0.25, // 事业权重最高（10宫天顶）
	Relationship: 0.20, // 关系（7宫对宫）
	Health:       0.20, // 健康（1+6宫）
	Finance:      0.20, // 财务（2+8宫）
	Spiritual:    0.15, // 灵性（9+12宫）
}

// ==================== 因子时间级别映射 ====================

// FactorTimeLevelMapping 因子类型到时间级别的映射
var FactorTimeLevelMapping = map[models.InfluenceFactorType]models.FactorTimeLevel{
	models.FactorDignity:        models.TimeLevelMonthly, // 行星换座周期
	models.FactorRetrograde:     models.TimeLevelWeekly,  // 逆行周期数周
	models.FactorAspectPhase:    models.TimeLevelDaily,   // 相位变化较快
	models.FactorAspectOrb:      models.TimeLevelDaily,
	models.FactorOuterPlanet:    models.TimeLevelYearly,  // 外行星影响长期
	models.FactorProfectionLord: models.TimeLevelYearly,  // 年主星为年度级
	models.FactorLunarPhase:     models.TimeLevelDaily,   // 月相周期约29.5天
	models.FactorPlanetaryHour:  models.TimeLevelHourly,  // 行星时为小时级
	models.FactorVoidOfCourse:   models.TimeLevelHourly,  // 月亮空亡为小时级
	models.FactorPersonal:       models.TimeLevelDaily,   // 个人因子默认日级
	models.FactorCustom:         models.TimeLevelHourly,  // 自定义因子可配置
}

// GetFactorTimeLevel 获取因子的时间级别
func GetFactorTimeLevel(factorType models.InfluenceFactorType) models.FactorTimeLevel {
	if level, ok := FactorTimeLevelMapping[factorType]; ok {
		return level
	}
	return models.TimeLevelDaily
}

// IsFactorVisibleAtLevel 检查因子在某个视图级别是否可见
func IsFactorVisibleAtLevel(factorLevel, viewLevel models.FactorTimeLevel) bool {
	levelOrder := map[models.FactorTimeLevel]int{
		models.TimeLevelYearly:  1,
		models.TimeLevelMonthly: 2,
		models.TimeLevelWeekly:  3,
		models.TimeLevelDaily:   4,
		models.TimeLevelHourly:  5,
	}

	factorOrder := levelOrder[factorLevel]
	viewOrder := levelOrder[viewLevel]

	// 大级别因子在小级别视图中可见
	return factorOrder <= viewOrder
}

// ==================== 宫位-维度映射 ====================

// HouseDimensionMapping 宫位到主维度的映射
var HouseDimensionMapping = map[int]string{
	1:  "health",       // 命宫 - 身体/自我
	2:  "finance",      // 财帛宫 - 个人资源
	3:  "relationship", // 兄弟宫 - 沟通/短途
	4:  "career",       // 田宅宫 - 家庭根基（影响事业基础）
	5:  "relationship", // 子女宫 - 恋爱/创造
	6:  "health",       // 奴仆宫 - 健康/工作
	7:  "relationship", // 夫妻宫 - 婚姻/合作
	8:  "finance",      // 疾厄宫 - 共享资源/转化
	9:  "spiritual",    // 迁移宫 - 高等教育/信仰
	10: "career",       // 官禄宫 - 事业/社会地位
	11: "relationship", // 福德宫 - 社交/愿望
	12: "spiritual",    // 玄秘宫 - 灵性/潜意识
}

// GetDimensionForHouseV2 根据宫位获取主维度
func GetDimensionForHouseV2(house int) string {
	if dim, ok := HouseDimensionMapping[house]; ok {
		return dim
	}
	return "spiritual"
}

