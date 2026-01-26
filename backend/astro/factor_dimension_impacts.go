package astro

import "star/models"

/*
因子维度有符号影响映射系统

核心理念：
天体事件的影响在不同维度可以是正向或负向的。
例如：水星逆行对沟通（事业/关系）是负面的，但对反思（灵性）是正面的。

DimensionImpact 现在包含有符号的影响值（-1.0 到 +1.0）：
- 正值：该维度受到正向影响
- 负值：该维度受到负向影响
- 0值：该维度不受影响或影响中性

这样同一个因子可以在不同维度产生不同方向的影响。
*/

// SignedDimensionImpact 有符号的维度影响
// 值域：-1.0 到 +1.0
type SignedDimensionImpact struct {
	Career       float64 // 事业：工作、职业发展、成就
	Relationship float64 // 关系：人际、感情、合作
	Health       float64 // 健康：身体、精力、活力
	Finance      float64 // 财务：金钱、收入、资源
	Spiritual    float64 // 灵性：内在成长、直觉、意义
}

// ==================== 基础因子类型的维度影响 ====================

// FactorTypeDefaultImpacts 因子类型的默认维度影响
// 注意：这是默认模板，实际使用时还会根据具体行星、相位等调整
var FactorTypeDefaultImpacts = map[models.InfluenceFactorType]SignedDimensionImpact{

	// ===== 逆行因子 =====
	// 逆行：回溯、延迟、反思、内省
	models.FactorRetrograde: {
		Career:       -0.30, // 负：计划延迟、沟通错乱、项目推迟
		Relationship: -0.20, // 负：沟通不畅、误解增加
		Health:       -0.10, // 轻微负：注意力下降
		Finance:      -0.15, // 负：交易延迟、合同问题
		Spiritual:    +0.25, // 正：适合反思、复盘、内省
	},

	// ===== 月空亡 =====
	// 月空亡：能量空白、事务停滞、不宜决策
	models.FactorVoidOfCourse: {
		Career:       -0.30, // 负：决策失误、计划难推进
		Relationship: -0.10, // 轻微负：沟通效率低
		Health:       -0.05, // 微负：精力不足
		Finance:      -0.35, // 负：不宜投资、签约
		Spiritual:    +0.20, // 正：适合冥想、放空
	},

	// ===== 日月食 =====
	// 日月食：转折点、剧变、释放旧模式
	models.FactorEclipse: {
		Career:       -0.20, // 负：可能有突然变动
		Relationship: -0.15, // 负：关系考验、揭示隐藏问题
		Health:       -0.25, // 负：能量不稳定
		Finance:      -0.10, // 负：财务波动
		Spiritual:    +0.30, // 正：觉醒、转化的机会
	},

	// ===== 燃烧因子 =====
	// 燃烧：行星能量被太阳压制、能见度下降
	models.FactorCombustion: {
		Career:       -0.25, // 负：才能难以展现
		Relationship: -0.20, // 负：表达受阻
		Health:       -0.15, // 负：活力下降
		Finance:      -0.15, // 负：资源流失
		Spiritual:    -0.05, // 轻微负：直觉模糊
	},

	// ===== 停滞因子 =====
	// 停滞：行星能量极度强化、影响深刻
	models.FactorStation: {
		Career:       +0.20, // 正：关键转折、深化主题
		Relationship: +0.10, // 正：关系议题凸显
		Health:       -0.10, // 负：能量停滞可能不适
		Finance:      +0.15, // 正：财务事宜需重视
		Spiritual:    +0.25, // 正：深刻内省的时机
	},

	// ===== 互容因子 =====
	// 互容：行星互相支持、能量交换顺畅
	models.FactorReception: {
		Career:       +0.30, // 正：合作顺利、支持增加
		Relationship: +0.35, // 正：关系和谐、互相理解
		Health:       +0.20, // 正：身心协调
		Finance:      +0.25, // 正：资源流通顺畅
		Spiritual:    +0.20, // 正：内外一致
	},

	// ===== 月相因子 =====
	// 月相：会根据具体相位（新月/满月等）动态调整
	models.FactorLunarPhase: {
		Career:       0.00,  // 中性：根据月相类型调整
		Relationship: +0.10, // 轻微正：情绪共鸣
		Health:       +0.05, // 轻微正：生理节律
		Finance:      0.00,  // 中性
		Spiritual:    +0.20, // 正：潮汐、节律感知
	},

	// ===== 行星时 =====
	// 行星时：短期能量波动
	models.FactorPlanetaryHour: {
		Career:       +0.15, // 轻微正：择时而动
		Relationship: +0.10, // 轻微正
		Health:       +0.10, // 轻微正
		Finance:      +0.15, // 轻微正
		Spiritual:    +0.10, // 轻微正
	},

	// ===== 月交点 =====
	// 月交点：命运方向、业力课题
	models.FactorLunarNode: {
		Career:       +0.15, // 正：命运推动
		Relationship: +0.20, // 正：业力关系显现
		Health:       0.00,  // 中性
		Finance:      +0.10, // 轻微正
		Spiritual:    +0.35, // 正：灵魂使命
	},

	// ===== 恒星因子 =====
	// 恒星：根据恒星性质会有很大差异
	models.FactorFixedStar: {
		Career:       0.00,  // 中性：需根据具体恒星调整
		Relationship: 0.00,  // 中性
		Health:       0.00,  // 中性
		Finance:      0.00,  // 中性
		Spiritual:    +0.20, // 轻微正：连接宇宙能量
	},

	// ===== 阿拉伯点 =====
	// 阿拉伯点：福点、精神点等
	models.FactorArabicPart: {
		Career:       +0.15, // 轻微正
		Relationship: +0.15, // 轻微正
		Health:       +0.15, // 轻微正
		Finance:      +0.20, // 正：福点相关
		Spiritual:    +0.25, // 正：精神点相关
	},

	// ===== 界限与十度面 =====
	models.FactorTerm: {
		Career:       +0.10, // 轻微正：尊贵度提升
		Relationship: +0.05, // 微正
		Health:       +0.05, // 微正
		Finance:      +0.05, // 微正
		Spiritual:    +0.05, // 微正
	},

	models.FactorDecan: {
		Career:       +0.08, // 轻微正：面尊贵
		Relationship: +0.05, // 微正
		Health:       +0.05, // 微正
		Finance:      +0.05, // 微正
		Spiritual:    +0.05, // 微正
	},

	// ===== 太阳弧推进 =====
	models.FactorSolarArc: {
		Career:       +0.25, // 正：重要人生转折
		Relationship: +0.20, // 正：关系进展
		Health:       +0.15, // 正：生命阶段转换
		Finance:      +0.20, // 正：财务转机
		Spiritual:    +0.20, // 正：成长契机
	},

	// ===== 法达时间主星 =====
	models.FactorFirdaria: {
		Career:       +0.25, // 正：主导时期
		Relationship: +0.15, // 正
		Health:       +0.15, // 正
		Finance:      +0.20, // 正
		Spiritual:    +0.15, // 正
	},

	// ===== 黄道释放 =====
	models.FactorZodiacal: {
		Career:       +0.20, // 正：释放期
		Relationship: +0.15, // 正
		Health:       +0.10, // 轻微正
		Finance:      +0.20, // 正
		Spiritual:    +0.15, // 正
	},

	// ===== 年主星（小限法）=====
	models.FactorProfectionLord: {
		Career:       +0.20, // 正：年度主题
		Relationship: +0.15, // 正
		Health:       +0.15, // 正
		Finance:      +0.15, // 正
		Spiritual:    +0.15, // 正
	},

	// ===== 相位容许度 =====
	// 相位容许度：相位精确度加成
	models.FactorAspectOrb: {
		Career:       +0.05, // 轻微正：精确相位影响强
		Relationship: +0.05, // 轻微正
		Health:       +0.05, // 轻微正
		Finance:      +0.05, // 轻微正
		Spiritual:    +0.05, // 轻微正
	},

	// ===== 外行星过境 =====
	// 外行星过境：天王星、海王星、冥王星的长期影响
	models.FactorOuterPlanet: {
		Career:       +0.20, // 正：重大转变
		Relationship: +0.15, // 正：关系深化或重组
		Health:       -0.10, // 负：需要适应变化
		Finance:      +0.15, // 正：资源重组
		Spiritual:    +0.30, // 正：深刻觉醒与转化
	},

	// ===== 个人因子 =====
	// 个人因子：太阳回归、月亮回归等个人周期
	models.FactorPersonal: {
		Career:       +0.20, // 正：个人周期重启
		Relationship: +0.15, // 正
		Health:       +0.15, // 正：生命力更新
		Finance:      +0.15, // 正
		Spiritual:    +0.20, // 正：自我认知深化
	},

	// ===== 中点技术 =====
	// 中点：两颗行星的中点被激活
	models.FactorMidpoint: {
		Career:       +0.10, // 轻微正：融合能量
		Relationship: +0.15, // 正：关系融合点
		Health:       +0.05, // 微正
		Finance:      +0.10, // 轻微正
		Spiritual:    +0.10, // 轻微正：整合能量
	},

	// ===== 反生点 =====
	// 反生点：镜像点技术，隐藏的联系
	models.FactorAntiscion: {
		Career:       +0.08, // 微正：隐藏机会
		Relationship: +0.12, // 轻微正：隐藏连接
		Health:       +0.05, // 微正
		Finance:      +0.08, // 微正：隐藏资源
		Spiritual:    +0.15, // 正：揭示隐藏联系
	},

	// ===== 主限推进 =====
	// 主限推进：基于地球自转的推运技术
	models.FactorPrimary: {
		Career:       +0.25, // 正：重要时间节点
		Relationship: +0.20, // 正：关系里程碑
		Health:       +0.15, // 正：生命周期转换
		Finance:      +0.20, // 正：财务转折
		Spiritual:    +0.25, // 正：灵性成长阶段
	},

	// ===== 自定义因子 =====
	models.FactorCustom: {
		Career:       0.00, // 由用户指定
		Relationship: 0.00,
		Health:       0.00,
		Finance:      0.00,
		Spiritual:    0.00,
	},
}

// ==================== 行星特定的逆行影响 ====================

// RetrogradeImpactsByPlanet 不同行星逆行的维度影响
var RetrogradeImpactsByPlanet = map[models.PlanetID]SignedDimensionImpact{
	// 水星逆行：沟通、交通、电子设备
	models.Mercury: {
		Career:       -0.40, // 强负：沟通障碍、合同延误
		Relationship: -0.35, // 强负：误解增多
		Health:       -0.10, // 轻微负：神经紧张
		Finance:      -0.25, // 负：交易延迟
		Spiritual:    +0.30, // 正：反思、重新思考
	},

	// 金星逆行：爱情、美感、价值观
	models.Venus: {
		Career:       -0.15, // 负：创意受阻
		Relationship: -0.40, // 强负：感情考验、旧情复燃
		Health:       -0.05, // 微负
		Finance:      -0.25, // 负：消费欲望混乱
		Spiritual:    +0.25, // 正：重新审视价值观
	},

	// 火星逆行：行动力、竞争、性能量
	models.Mars: {
		Career:       -0.35, // 强负：行动受阻、竞争力下降
		Relationship: -0.30, // 负：冲突增多、性张力
		Health:       -0.25, // 负：精力不足、发炎
		Finance:      -0.20, // 负：投资冲动受阻
		Spiritual:    +0.20, // 正：学习控制冲动
	},

	// 木星逆行：扩张、机遇、信念
	models.Jupiter: {
		Career:       -0.25, // 负：机会延迟、扩张受限
		Relationship: -0.10, // 轻微负
		Health:       -0.10, // 轻微负：过度乐观
		Finance:      -0.30, // 负：投资时机需谨慎
		Spiritual:    +0.35, // 正：内在信仰深化
	},

	// 土星逆行：责任、结构、限制
	models.Saturn: {
		Career:       -0.20, // 负：责任加重、进展缓慢
		Relationship: -0.15, // 负：承诺议题
		Health:       -0.20, // 负：慢性问题凸显
		Finance:      -0.15, // 负：保守、限制
		Spiritual:    +0.40, // 强正：深刻内省、业力清理
	},

	// 天王星逆行：突破、自由、创新
	models.Uranus: {
		Career:       -0.15, // 负：创新受阻
		Relationship: -0.20, // 负：独立性冲突
		Health:       -0.10, // 轻微负：神经系统
		Finance:      -0.10, // 轻微负：投资波动
		Spiritual:    +0.35, // 正：内在觉醒
	},

	// 海王星逆行：幻想、灵感、迷惑
	models.Neptune: {
		Career:       -0.10, // 轻微负：理想受质疑
		Relationship: -0.15, // 负：幻想破灭
		Health:       -0.15, // 负：免疫力、中毒
		Finance:      -0.20, // 负：财务迷惑澄清
		Spiritual:    +0.40, // 强正：灵性修炼深化
	},

	// 冥王星逆行：转化、权力、深层资源
	models.Pluto: {
		Career:       -0.15, // 负：权力斗争内化
		Relationship: -0.20, // 负：控制欲、深层问题
		Health:       -0.20, // 负：深层疗愈需要
		Finance:      -0.15, // 负：隐藏财务问题浮现
		Spiritual:    +0.45, // 强正：深刻转化、重生
	},
}

// ==================== 相位类型的维度影响修正 ====================

// AspectTypeImpactModifiers 相位类型对维度影响的修正系数
// 这会与基础相位影响值相乘
var AspectTypeImpactModifiers = map[string]map[string]float64{
	// 合相：能量融合，强化
	"conjunction": {
		"career":       1.0,
		"relationship": 1.0,
		"health":       1.0,
		"finance":      1.0,
		"spiritual":    1.0,
	},

	// 六分相：和谐机会
	"sextile": {
		"career":       1.0,
		"relationship": 1.2, // 关系更和谐
		"health":       1.0,
		"finance":      1.0,
		"spiritual":    1.1,
	},

	// 四分相：紧张挑战
	"square": {
		"career":       1.0,
		"relationship": 0.7, // 关系压力
		"health":       0.8, // 健康压力
		"finance":      0.9,
		"spiritual":    1.2, // 挑战促进成长
	},

	// 三分相：和谐流动
	"trine": {
		"career":       1.1,
		"relationship": 1.2, // 关系顺畅
		"health":       1.1,
		"finance":      1.1,
		"spiritual":    1.0,
	},

	// 冲相：对立平衡
	"opposition": {
		"career":       0.9,
		"relationship": 0.6, // 关系对立
		"health":       0.9,
		"finance":      0.9,
		"spiritual":    1.1, // 需要平衡与整合
	},
}

// ==================== 尊贵度的维度影响 ====================

// DignityImpactsByType 尊贵度类型的维度影响
var DignityImpactsByType = map[string]SignedDimensionImpact{
	// 入庙（Domicile）：行星在自己守护的星座
	"domicile": {
		Career:       +0.30,
		Relationship: +0.25,
		Health:       +0.25,
		Finance:      +0.25,
		Spiritual:    +0.20,
	},

	// 旺相（Exaltation）：行星能量最强
	"exaltation": {
		Career:       +0.35,
		Relationship: +0.30,
		Health:       +0.30,
		Finance:      +0.30,
		Spiritual:    +0.25,
	},

	// 失势（Detriment）：行星在对宫星座
	"detriment": {
		Career:       -0.25,
		Relationship: -0.30,
		Health:       -0.20,
		Finance:      -0.25,
		Spiritual:    -0.15,
	},

	// 落陷（Fall）：行星能量最弱
	"fall": {
		Career:       -0.30,
		Relationship: -0.35,
		Health:       -0.25,
		Finance:      -0.30,
		Spiritual:    -0.20,
	},
}

// ==================== 辅助函数 ====================

// GetFactorDimensionImpact 获取因子的维度影响
// 这是新的接口，返回有符号的维度影响
// 现在使用V2强化版，引入维度对立性
func GetFactorDimensionImpact(factor *models.InfluenceFactor) SignedDimensionImpact {
	// 使用V2强化版本
	impact := GetFactorDimensionImpactV2(factor)
	
	// 应用维度过滤器：让不同维度受到不同因子集合影响
	impact.Career = ApplyDimensionFilter(factor, "career", impact.Career)
	impact.Relationship = ApplyDimensionFilter(factor, "relationship", impact.Relationship)
	impact.Health = ApplyDimensionFilter(factor, "health", impact.Health)
	impact.Finance = ApplyDimensionFilter(factor, "finance", impact.Finance)
	impact.Spiritual = ApplyDimensionFilter(factor, "spiritual", impact.Spiritual)
	
	return impact
}

// extractDignityType 从因子名称提取尊贵度类型
func extractDignityType(name string) string {
	// 简单的字符串匹配
	if contains(name, "Exalted") {
		return "exaltation"
	}
	if contains(name, "Domicile") || contains(name, "Home") {
		return "domicile"
	}
	if contains(name, "Detriment") {
		return "detriment"
	}
	if contains(name, "Fall") {
		return "fall"
	}
	return "domicile" // 默认
}

// extractAspectType 从因子名称提取相位类型
func extractAspectType(name string) string {
	if contains(name, "Conjunction") {
		return "conjunction"
	}
	if contains(name, "Sextile") {
		return "sextile"
	}
	if contains(name, "Square") {
		return "square"
	}
	if contains(name, "Trine") {
		return "trine"
	}
	if contains(name, "Opposition") {
		return "opposition"
	}
	return "conjunction" // 默认
}

// adjustImpactByPlanetNature 根据行星特性调整影响
func adjustImpactByPlanetNature(impact SignedDimensionImpact, planet models.PlanetID) SignedDimensionImpact {
	// 根据行星的自然属性，放大某些维度的影响
	switch planet {
	case models.Mercury:
		// 水星：沟通、思维
		impact.Career *= 1.2
		impact.Relationship *= 1.1

	case models.Venus:
		// 金星：爱、美、和谐
		impact.Relationship *= 1.3
		impact.Finance *= 1.1

	case models.Mars:
		// 火星：行动、竞争
		impact.Career *= 1.2
		impact.Health *= 1.1

	case models.Jupiter:
		// 木星：扩张、幸运
		impact.Career *= 1.1
		impact.Finance *= 1.2
		impact.Spiritual *= 1.1

	case models.Saturn:
		// 土星：责任、限制
		impact.Career *= 1.2
		impact.Health *= 1.1

	case models.Uranus:
		// 天王星：创新、突破
		impact.Career *= 1.1
		impact.Spiritual *= 1.2

	case models.Neptune:
		// 海王星：灵感、幻想
		impact.Spiritual *= 1.3
		impact.Relationship *= 1.1

	case models.Pluto:
		// 冥王星：转化、权力
		impact.Finance *= 1.2
		impact.Spiritual *= 1.2
	}

	return impact
}

// contains 简单的字符串包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findInString(s, substr))))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ConvertToModelsDimensionImpact 将SignedDimensionImpact转换为models.DimensionImpact
// 用于向后兼容，将有符号影响转换为无符号权重
func ConvertToModelsDimensionImpact(signed SignedDimensionImpact) models.DimensionImpact {
	// 取绝对值并归一化为总和1.0的权重
	abs := func(v float64) float64 {
		if v < 0 {
			return -v
		}
		return v
	}

	total := abs(signed.Career) + abs(signed.Relationship) +
		abs(signed.Health) + abs(signed.Finance) + abs(signed.Spiritual)

	if total == 0 {
		// 如果全为0，返回平均分配
		return models.DimensionImpact{
			Career: 0.2, Relationship: 0.2, Health: 0.2,
			Finance: 0.2, Spiritual: 0.2,
		}
	}

	return models.DimensionImpact{
		Career:       abs(signed.Career) / total,
		Relationship: abs(signed.Relationship) / total,
		Health:       abs(signed.Health) / total,
		Finance:      abs(signed.Finance) / total,
		Spiritual:    abs(signed.Spiritual) / total,
	}
}
