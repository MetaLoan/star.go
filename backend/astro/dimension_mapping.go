package astro

import "star/models"

// ==================== 行星-维度影响矩阵 ====================
// 基于占星学理论：行星的自然象征意义决定其对各生活领域的影响比例
// 参考文献：Liz Greene, Stephen Arroyo 等现代心理占星学著作
// 每行总和 = 1.0

// PlanetDimensionMapping 行星-维度影响映射表
// 优化策略：每个行星最多2个主要维度，权重集中，确保标签清晰
// 单维度行星：金星、土星、天王星、海王星
// 双维度行星：太阳、月亮、水星、火星、木星、冥王星、北交点、凯龙
var PlanetDimensionMapping = map[models.PlanetID]models.DimensionImpact{
	// 太阳：领导力、生命力 → 事业+健康
	models.Sun: {
		Career:       0.6, // 领导、成就、职业表现
		Health:       0.4, // 生命活力、体能状态
		Relationship: 0,
		Finance:      0,
		Spiritual:    0,
	},

	// 月亮：情感、身心 → 关系+健康  
	models.Moon: {
		Relationship: 0.6, // 情感需求、亲密关系
		Health:       0.4, // 情绪健康、身体节律
		Career:       0,
		Finance:      0,
		Spiritual:    0,
	},

	// 水星：沟通、商业 → 事业+财运
	models.Mercury: {
		Career:       0.5, // 工作沟通、职业技能
		Finance:      0.5, // 商业交易、财务分析
		Relationship: 0,
		Health:       0,
		Spiritual:    0,
	},

	// 金星：爱与价值 → 关系（单维度）
	models.Venus: {
		Relationship: 0.8, // 爱情、美感、和谐关系
		Finance:      0.2, // 次要：价值观、享受消费
		Career:       0,
		Health:       0,
		Spiritual:    0,
	},

	// 火星：行动、能量 → 健康+事业
	models.Mars: {
		Health:       0.6, // 体力、身体能量
		Career:       0.4, // 行动力、竞争力
		Relationship: 0,
		Finance:      0,
		Spiritual:    0,
	},

	// 木星：扩张、机遇 → 灵性+财运
	models.Jupiter: {
		Spiritual:    0.6, // 信仰、智慧、精神成长
		Finance:      0.4, // 财富扩张、投资机会
		Career:       0,
		Relationship: 0,
		Health:       0,
	},

	// 土星：责任、成就 → 事业（单维度）
	models.Saturn: {
		Career:       0.8, // 职业责任、长期成就
		Health:       0.2, // 次要：慢性健康、骨骼
		Relationship: 0,
		Finance:      0,
		Spiritual:    0,
	},

	// 天王星：觉醒、创新 → 灵性（单维度）
	models.Uranus: {
		Spiritual:    0.8, // 意识觉醒、突破创新
		Career:       0.2, // 次要：科技创新
		Relationship: 0,
		Health:       0,
		Finance:      0,
	},

	// 海王星：灵性、理想 → 灵性（单维度）
	models.Neptune: {
		Spiritual:    0.8, // 灵性直觉、精神追求
		Relationship: 0.2, // 次要：理想化关系
		Career:       0,
		Health:       0,
		Finance:      0,
	},

	// 冥王星：转化、资源 → 灵性+财运
	models.Pluto: {
		Spiritual:    0.6, // 深层转化、重生
		Finance:      0.4, // 共享资源、深层财富
		Career:       0,
		Relationship: 0,
		Health:       0,
	},

	// 北交点：命运、成长 → 灵性+关系
	models.NorthNode: {
		Spiritual:    0.6, // 灵魂使命、精神方向
		Relationship: 0.4, // 业力关系、命运连接
		Career:       0,
		Health:       0,
		Finance:      0,
	},

	// 凯龙：疗愈、伤痛 → 灵性+健康
	models.Chiron: {
		Spiritual:    0.6, // 心灵疗愈、精神成长
		Health:       0.4, // 身心伤痛、康复
		Career:       0,
		Relationship: 0,
		Finance:      0,
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
// 时间级别说明：
// - Yearly（年级别）：影响持续数月至数年，如推运技术、外行星过境
// - Monthly（月级别）：影响持续数周至数月，如行星换座、逆行
// - Weekly（周级别）：影响持续数天至数周，如快速相位
// - Daily（日级别）：影响持续数小时至数天，如精确相位、月相
// - Hourly（小时级别）：影响持续数分钟至数小时，如行星时、月空
var FactorTimeLevelMapping = map[models.InfluenceFactorType]models.FactorTimeLevel{
	// ===== 基础因子 =====
	models.FactorDignity:        models.TimeLevelMonthly, // 行星在星座停留时间（数周至数月）
	models.FactorRetrograde:     models.TimeLevelWeekly,  // 逆行周期（数周）
	models.FactorAspectPhase:    models.TimeLevelDaily,   // 相位（根据行星速度，数天至数周）
	models.FactorAspectOrb:      models.TimeLevelDaily,   // 相位容许度
	models.FactorOuterPlanet:    models.TimeLevelYearly,  // 外行星过境（数月至数年）
	models.FactorProfectionLord: models.TimeLevelYearly,  // 年主星（整年）
	models.FactorLunarPhase:     models.TimeLevelDaily,   // 月相周期（约29.5天，但精确影响数天）
	models.FactorPlanetaryHour:  models.TimeLevelHourly,  // 行星时（约1-2小时）
	models.FactorVoidOfCourse:   models.TimeLevelHourly,  // 月亮空亡（数小时）
	models.FactorPersonal:       models.TimeLevelYearly,  // 个人因子如太阳回归（年度级）
	models.FactorCustom:         models.TimeLevelDaily,   // 自定义因子（默认日级，可配置）

	// ===== 日月食与交点 =====
	models.FactorEclipse:   models.TimeLevelMonthly, // 日月食影响期（前后2-4周）
	models.FactorLunarNode: models.TimeLevelWeekly,  // 月交点过境（数天至数周）

	// ===== 行星状态 =====
	models.FactorCombustion: models.TimeLevelDaily,  // 燃烧（数天，取决于行星速度）
	models.FactorStation:    models.TimeLevelDaily,  // 停滞（数天）
	models.FactorReception:  models.TimeLevelMonthly, // 互容（取决于行星在星座时间）

	// ===== 恒星与特殊点 =====
	models.FactorFixedStar:  models.TimeLevelDaily,  // 恒星合相（容许度小，持续数天）
	models.FactorArabicPart: models.TimeLevelDaily,  // 阿拉伯点（日级变化）
	models.FactorMidpoint:   models.TimeLevelDaily,  // 中点（日级）
	models.FactorAntiscion:  models.TimeLevelDaily,  // 反生点（日级）

	// ===== 界限与分度 =====
	models.FactorTerm:  models.TimeLevelWeekly, // 界限（行星在一个界限约数天至数周）
	models.FactorDecan: models.TimeLevelWeekly, // 十度面（行星在一个面约数天至数周）

	// ===== 推运技术 =====
	models.FactorSolarArc: models.TimeLevelYearly, // 太阳弧推进（年度级，影响约±6个月）
	models.FactorPrimary:  models.TimeLevelYearly, // 主限推进（年度级）
	models.FactorFirdaria: models.TimeLevelYearly, // 法达（多年级）
	models.FactorZodiacal: models.TimeLevelYearly, // 黄道释放（年度级）
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

