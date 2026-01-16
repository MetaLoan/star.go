package astro

import "star/models"

/*
因子维度影响系统 V2 - 强化版

核心设计原则：
1. 维度对立性：某个维度强了，其他维度可能变弱（零和博弈）
2. 维度专属性：某些因子只显著影响特定维度
3. 生活平衡理论：时间和精力是有限的，专注一个领域会牺牲其他领域

示例：
- 土星过境第10宫 → 事业大涨(+0.50)，但关系疏远(-0.30)、健康透支(-0.20)
- 金星-火星合相 → 关系热烈(+0.50)，但工作分心(-0.20)
- 水星逆行 → 灵性成长(+0.40)，但沟通障碍(-0.40)
*/

// ==================== 基础因子类型的维度影响（V2强化版）====================

// FactorTypeDefaultImpactsV2 因子类型的默认维度影响（强化版）
var FactorTypeDefaultImpactsV2 = map[models.InfluenceFactorType]SignedDimensionImpact{

	// ===== 逆行因子 =====
	// 原则：物质世界受阻 ↔ 灵性世界深化
	models.FactorRetrograde: {
		Career:       -0.35, // 负：外在行动受阻
		Relationship: -0.25, // 负：沟通误解
		Health:       -0.15, // 负：注意力下降
		Finance:      -0.20, // 负：交易延迟
		Spiritual:    +0.50, // 强正：被迫内省，灵性深化
	},

	// ===== 月空亡 =====
	// 原则：世俗事务停滞 ↔ 灵性空间开启
	models.FactorVoidOfCourse: {
		Career:       -0.40, // 强负：决策失误
		Relationship: -0.15, // 负：沟通空洞
		Health:       -0.10, // 轻微负
		Finance:      -0.45, // 强负：不宜投资
		Spiritual:    +0.35, // 正：适合冥想、放空
	},

	// ===== 日月食 =====
	// 原则：物质动荡 ↔ 灵性觉醒
	models.FactorEclipse: {
		Career:       -0.30, // 负：突然变动
		Relationship: -0.25, // 负：关系考验
		Health:       -0.35, // 负：能量失衡
		Finance:      -0.20, // 负：财务波动
		Spiritual:    +0.60, // 强正：重大觉醒、转化机会
	},

	// ===== 燃烧因子 =====
	// 原则：全面压制，无对立维度
	models.FactorCombustion: {
		Career:       -0.30, // 负：才能被遮蔽
		Relationship: -0.25, // 负：表达受阻
		Health:       -0.20, // 负：活力下降
		Finance:      -0.20, // 负：资源流失
		Spiritual:    -0.15, // 负：直觉模糊
	},

	// ===== 停滞因子 =====
	// 原则：能量凝聚，专注vs健康
	models.FactorStation: {
		Career:       +0.35, // 正：关键转折
		Relationship: +0.20, // 正：关系议题凸显
		Health:       -0.25, // 负：能量停滞不适
		Finance:      +0.25, // 正：财务凝聚
		Spiritual:    +0.40, // 正：深刻内省
	},

	// ===== 互容因子 =====
	// 原则：全面和谐（少数全正因子）
	models.FactorReception: {
		Career:       +0.35, // 正：合作顺利
		Relationship: +0.45, // 强正：关系和谐
		Health:       +0.25, // 正：身心协调
		Finance:      +0.30, // 正：资源流通
		Spiritual:    +0.25, // 正：内外一致
	},

	// ===== 月相因子 =====
	// 原则：情绪vs理性
	models.FactorLunarPhase: {
		Career:       -0.15, // 负：情绪影响决策
		Relationship: +0.30, // 正：情绪共鸣
		Health:       +0.10, // 轻微正：生理节律
		Finance:      -0.10, // 负：冲动消费
		Spiritual:    +0.35, // 正：潮汐、节律感知
	},

	// ===== 行星时 =====
	// 原则：择时而动（专属性，根据具体行星调整）
	models.FactorPlanetaryHour: {
		Career:       +0.20, // 正：择时而动
		Relationship: +0.10, // 轻微正
		Health:       +0.10, // 轻微正
		Finance:      +0.20, // 正：时机把握
		Spiritual:    +0.15, // 正：能量对齐
	},

	// ===== 月交点 =====
	// 原则：命运vs日常
	models.FactorLunarNode: {
		Career:       +0.20, // 正：命运推动
		Relationship: +0.30, // 正：业力关系
		Health:       -0.05, // 微负：忽视健康
		Finance:      +0.10, // 轻微正
		Spiritual:    +0.50, // 强正：灵魂使命
	},

	// ===== 恒星因子 =====
	// 原则：宇宙vs尘世（根据具体恒星调整）
	models.FactorFixedStar: {
		Career:       -0.10, // 负：超越世俗
		Relationship: 0.00,  // 中性
		Health:       +0.10, // 轻微正：宇宙能量
		Finance:      -0.05, // 微负：不看重物质
		Spiritual:    +0.40, // 正：连接宇宙
	},

	// ===== 阿拉伯点 =====
	// 原则：福点（物质）vs精神点（灵性）
	models.FactorArabicPart: {
		Career:       +0.20,
		Relationship: +0.20,
		Health:       +0.20,
		Finance:      +0.30, // 福点相关
		Spiritual:    +0.35, // 精神点相关
	},

	// ===== 界限与十度面 =====
	models.FactorTerm: {
		Career:       +0.15,
		Relationship: +0.08,
		Health:       +0.08,
		Finance:      +0.08,
		Spiritual:    +0.08,
	},

	models.FactorDecan: {
		Career:       +0.12,
		Relationship: +0.08,
		Health:       +0.08,
		Finance:      +0.08,
		Spiritual:    +0.08,
	},

	// ===== 推运技术 =====
	models.FactorSolarArc: {
		Career:       +0.35,
		Relationship: +0.30,
		Health:       +0.20,
		Finance:      +0.30,
		Spiritual:    +0.30,
	},

	models.FactorFirdaria: {
		Career:       +0.35,
		Relationship: +0.20,
		Health:       +0.20,
		Finance:      +0.30,
		Spiritual:    +0.20,
	},

	models.FactorZodiacal: {
		Career:       +0.30,
		Relationship: +0.20,
		Health:       +0.15,
		Finance:      +0.30,
		Spiritual:    +0.20,
	},

	models.FactorProfectionLord: {
		Career:       +0.30,
		Relationship: +0.20,
		Health:       +0.20,
		Finance:      +0.20,
		Spiritual:    +0.20,
	},

	// ===== 其他因子 =====
	models.FactorAspectOrb: {
		Career:       +0.08,
		Relationship: +0.08,
		Health:       +0.08,
		Finance:      +0.08,
		Spiritual:    +0.08,
	},

	models.FactorOuterPlanet: {
		Career:       +0.25,
		Relationship: -0.10, // 负：变革打乱关系
		Health:       -0.15, // 负：适应压力
		Finance:      +0.20,
		Spiritual:    +0.45, // 正：深刻转化
	},

	models.FactorPersonal: {
		Career:       +0.25,
		Relationship: +0.20,
		Health:       +0.25,
		Finance:      +0.20,
		Spiritual:    +0.30,
	},

	models.FactorMidpoint: {
		Career:       +0.15,
		Relationship: +0.25, // 正：关系融合
		Health:       +0.10,
		Finance:      +0.15,
		Spiritual:    +0.15,
	},

	models.FactorAntiscion: {
		Career:       +0.12,
		Relationship: +0.20, // 正：隐藏连接
		Health:       +0.08,
		Finance:      +0.12,
		Spiritual:    +0.25, // 正：揭示隐藏
	},

	models.FactorPrimary: {
		Career:       +0.35,
		Relationship: +0.30,
		Health:       +0.20,
		Finance:      +0.30,
		Spiritual:    +0.35,
	},

	models.FactorCustom: {
		Career:       0.00,
		Relationship: 0.00,
		Health:       0.00,
		Finance:      0.00,
		Spiritual:    0.00,
	},
}

// ==================== 行星特定的逆行影响（V2强化版）====================

// RetrogradeImpactsByPlanetV2 不同行星逆行的维度影响（强化版）
var RetrogradeImpactsByPlanetV2 = map[models.PlanetID]SignedDimensionImpact{
	// 水星逆行：沟通vs内省
	models.Mercury: {
		Career:       -0.50, // 强负：沟通障碍、合同延误
		Relationship: -0.45, // 强负：误解增多
		Health:       -0.15, // 负：神经紧张
		Finance:      -0.35, // 负：交易延迟
		Spiritual:    +0.60, // 强正：深度思考、复盘
	},

	// 金星逆行：关系vs自我价值
	models.Venus: {
		Career:       -0.20, // 负：创意受阻
		Relationship: -0.55, // 强负：感情考验、旧情复燃
		Health:       -0.10, // 轻微负
		Finance:      -0.35, // 负：消费混乱
		Spiritual:    +0.50, // 正：重新审视价值观
	},

	// 火星逆行：行动vs耐心
	models.Mars: {
		Career:       -0.45, // 强负：行动受阻
		Relationship: -0.40, // 负：冲突增多
		Health:       -0.35, // 负：精力不足、炎症
		Finance:      -0.30, // 负：投资冲动
		Spiritual:    +0.40, // 正：学习控制冲动
	},

	// 木星逆行：扩张vs内在信仰
	models.Jupiter: {
		Career:       -0.35, // 负：机会延迟
		Relationship: -0.15, // 轻微负
		Health:       -0.15, // 轻微负
		Finance:      -0.40, // 负：投资需谨慎
		Spiritual:    +0.55, // 强正：内在信仰深化
	},

	// 土星逆行：外在责任vs内在成熟
	models.Saturn: {
		Career:       -0.30, // 负：责任加重、进展缓慢
		Relationship: -0.25, // 负：承诺议题
		Health:       -0.30, // 负：慢性问题
		Finance:      -0.25, // 负：保守、限制
		Spiritual:    +0.65, // 强正：深刻内省、业力清理
	},

	// 天王星逆行：外在突破vs内在觉醒
	models.Uranus: {
		Career:       -0.25, // 负：创新受阻
		Relationship: -0.30, // 负：独立性冲突
		Health:       -0.15, // 负：神经系统
		Finance:      -0.20, // 负：投资波动
		Spiritual:    +0.55, // 强正：内在觉醒
	},

	// 海王星逆行：幻想vs真实灵性
	models.Neptune: {
		Career:       -0.15, // 负：理想受质疑
		Relationship: -0.25, // 负：幻想破灭
		Health:       -0.25, // 负：免疫力、迷惑
		Finance:      -0.30, // 负：财务迷惑
		Spiritual:    +0.60, // 强正：真实灵性修炼
	},

	// 冥王星逆行：外在权力vs内在转化
	models.Pluto: {
		Career:       -0.25, // 负：权力斗争内化
		Relationship: -0.30, // 负：控制欲、深层问题
		Health:       -0.30, // 负：深层疗愈需要
		Finance:      -0.25, // 负：隐藏财务问题
		Spiritual:    +0.70, // 强正：深刻转化、重生
	},
}

// ==================== 尊贵度的维度影响（V2强化版）====================

// DignityImpactsByTypeV2 尊贵度类型的维度影响（引入对立性）
var DignityImpactsByTypeV2 = map[string]SignedDimensionImpact{
	// 入庙（Domicile）：全面强化，但略有侧重
	"domicile": {
		Career:       +0.35,
		Relationship: +0.30,
		Health:       +0.30,
		Finance:      +0.30,
		Spiritual:    +0.25,
	},

	// 旺相（Exaltation）：极强，但可能过度（根据行星调整）
	"exaltation": {
		Career:       +0.45,
		Relationship: +0.35,
		Health:       +0.35,
		Finance:      +0.35,
		Spiritual:    +0.30,
	},

	// 失势（Detriment）：对立性强
	"detriment": {
		Career:       -0.35,
		Relationship: -0.40,
		Health:       -0.25,
		Finance:      -0.35,
		Spiritual:    -0.20,
	},

	// 落陷（Fall）：严重削弱
	"fall": {
		Career:       -0.40,
		Relationship: -0.45,
		Health:       -0.35,
		Finance:      -0.40,
		Spiritual:    -0.25,
	},
}

// ==================== 相位类型的维度影响修正（V2强化版）====================

// AspectTypeImpactModifiersV2 相位类型的维度影响修正系数（强化版）
var AspectTypeImpactModifiersV2 = map[string]map[string]float64{
	// 合相：能量融合
	"conjunction": {
		"career":       1.0,
		"relationship": 1.0,
		"health":       1.0,
		"finance":      1.0,
		"spiritual":    1.0,
	},

	// 六分相：和谐机会，关系加成
	"sextile": {
		"career":       1.0,
		"relationship": 1.4, // 关系显著加成
		"health":       1.0,
		"finance":      1.0,
		"spiritual":    1.2,
	},

	// 四分相：紧张挑战，物质vs灵性
	"square": {
		"career":       0.5, // 事业压力大
		"relationship": 0.4, // 关系紧张
		"health":       0.6, // 健康压力
		"finance":      0.7,
		"spiritual":    1.5, // 挑战促进灵性成长
	},

	// 三分相：和谐流动
	"trine": {
		"career":       1.2,
		"relationship": 1.3,
		"health":       1.1,
		"finance":      1.2,
		"spiritual":    1.0,
	},

	// 冲相：对立平衡，关系挑战vs灵性整合
	"opposition": {
		"career":       0.7,
		"relationship": 0.3, // 关系强烈对立
		"health":       0.8,
		"finance":      0.8,
		"spiritual":    1.4, // 需要整合与平衡
	},
}

// ==================== 行星特性的维度调整（V2强化版）====================

// adjustImpactByPlanetNatureV2 根据行星特性调整影响（引入对立性）
func adjustImpactByPlanetNatureV2(impact SignedDimensionImpact, planet models.PlanetID) SignedDimensionImpact {
	switch planet {
	case models.Sun:
		// 太阳：自我vs关系（零和博弈）
		impact.Career *= 1.3
		impact.Relationship *= 0.7 // 自我强→关系弱
		impact.Spiritual *= 1.2

	case models.Moon:
		// 月亮：情绪vs理性
		impact.Relationship *= 1.4
		impact.Health *= 1.2
		impact.Career *= 0.8 // 情绪强→理性弱

	case models.Mercury:
		// 水星：沟通vs行动
		impact.Career *= 1.3
		impact.Relationship *= 1.2

	case models.Venus:
		// 金星：关系vs事业（零和博弈）
		impact.Relationship *= 1.5
		impact.Finance *= 1.2
		impact.Career *= 0.7 // 享受→工作分心

	case models.Mars:
		// 火星：行动vs关系（零和博弈）
		impact.Career *= 1.4
		impact.Health *= 1.2
		impact.Relationship *= 0.6 // 竞争→关系紧张

	case models.Jupiter:
		// 木星：扩张（可能过度）
		impact.Career *= 1.2
		impact.Finance *= 1.3
		impact.Spiritual *= 1.2
		impact.Health *= 0.9 // 过度乐观→忽视健康

	case models.Saturn:
		// 土星：事业vs生活平衡（零和博弈）
		impact.Career *= 1.4
		impact.Relationship *= 0.6 // 责任→关系疏远
		impact.Health *= 0.7       // 压力→健康透支

	case models.Uranus:
		// 天王星：创新vs稳定（零和博弈）
		impact.Career *= 1.2
		impact.Spiritual *= 1.4
		impact.Relationship *= 0.7 // 独立→关系距离

	case models.Neptune:
		// 海王星：灵性vs现实（零和博弈）
		impact.Spiritual *= 1.5
		impact.Career *= 0.6       // 理想→现实模糊
		impact.Finance *= 0.7      // 幻想→财务迷惑
		impact.Relationship *= 1.2 // 浪漫

	case models.Pluto:
		// 冥王星：转化vs舒适（零和博弈）
		impact.Spiritual *= 1.5
		impact.Finance *= 1.3
		impact.Health *= 0.7 // 深层转化→健康压力
	}

	return impact
}

// GetFactorDimensionImpactV2 获取因子的维度影响（V2强化版）
func GetFactorDimensionImpactV2(factor *models.InfluenceFactor) SignedDimensionImpact {
	// 1. 获取因子类型的默认影响（V2版本）
	defaultImpact, exists := FactorTypeDefaultImpactsV2[factor.Type]
	if !exists {
		return SignedDimensionImpact{}
	}

	impact := defaultImpact

	// 2. 针对特定因子类型进行细化
	switch factor.Type {
	case models.FactorRetrograde:
		// 逆行：使用V2强化版
		if planetImpact, ok := RetrogradeImpactsByPlanetV2[factor.SourcePlanet]; ok {
			impact = planetImpact
		}

	case models.FactorDignity:
		// 尊贵度：使用V2强化版
		dignityType := extractDignityType(factor.Name)
		if dignityImpact, ok := DignityImpactsByTypeV2[dignityType]; ok {
			impact = dignityImpact
			// 根据行星进一步调整对立性
			impact = adjustDignityByPlanet(impact, factor.SourcePlanet)
		}

	case models.FactorAspectPhase:
		// 相位：使用V2强化版修正
		aspectType := extractAspectType(factor.Name)
		if modifiers, ok := AspectTypeImpactModifiersV2[aspectType]; ok {
			impact.Career *= modifiers["career"]
			impact.Relationship *= modifiers["relationship"]
			impact.Health *= modifiers["health"]
			impact.Finance *= modifiers["finance"]
			impact.Spiritual *= modifiers["spiritual"]
		}

		// 应用行星特性的对立性调整
		impact = adjustImpactByPlanetNatureV2(impact, factor.SourcePlanet)
	}

	return impact
}

// adjustDignityByPlanet 根据行星调整尊贵度的维度对立性
func adjustDignityByPlanet(impact SignedDimensionImpact, planet models.PlanetID) SignedDimensionImpact {
	switch planet {
	case models.Mars:
		// 火星旺相：健康强但关系弱
		if impact.Health > 0 {
			impact.Health *= 1.5
			impact.Relationship *= 0.5
		}

	case models.Venus:
		// 金星旺相：关系强但事业弱
		if impact.Relationship > 0 {
			impact.Relationship *= 1.4
			impact.Career *= 0.7
		}

	case models.Saturn:
		// 土星旺相：事业强但关系、健康弱
		if impact.Career > 0 {
			impact.Career *= 1.5
			impact.Relationship *= 0.6
			impact.Health *= 0.7
		}

	case models.Jupiter:
		// 木星旺相：财务强但可能健康忽视
		if impact.Finance > 0 {
			impact.Finance *= 1.4
			impact.Health *= 0.8
		}
	}

	return impact
}
