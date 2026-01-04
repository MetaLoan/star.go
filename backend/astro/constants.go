package astro

import "star/models"

// ==================== 天文常量 ====================

const (
	// J2000 纪元儒略日
	J2000 = 2451545.0

	// 角度转换
	DEG_TO_RAD = 0.017453292519943295 // π/180
	RAD_TO_DEG = 57.29577951308232    // 180/π

	// 黄赤交角 (度)
	OBLIQUITY = 23.4393
)

// ==================== 行星信息 ====================

// PlanetInfo 行星信息
type PlanetInfo struct {
	ID     models.PlanetID
	Name   string
	Symbol string
	Color  string
	Weight float64
}

// Planets 行星列表
var Planets = []PlanetInfo{
	{models.Sun, "太阳", "☉", "#ffd700", 10},
	{models.Moon, "月亮", "☽", "#c0c0c0", 10},
	{models.Mercury, "水星", "☿", "#b5651d", 4},
	{models.Venus, "金星", "♀", "#ff69b4", 5},
	{models.Mars, "火星", "♂", "#dc143c", 6},
	{models.Jupiter, "木星", "♃", "#daa520", 7},
	{models.Saturn, "土星", "♄", "#8b7355", 8},
	{models.Uranus, "天王星", "♅", "#40e0d0", 5},
	{models.Neptune, "海王星", "♆", "#4169e1", 4},
	{models.Pluto, "冥王星", "♇", "#800080", 6},
	{models.NorthNode, "北交点", "☊", "#9932cc", 3},
	{models.Chiron, "凯龙星", "⚷", "#228b22", 3},
}

// PlanetWeights 行星权重
var PlanetWeights = map[models.PlanetID]float64{
	models.Sun:       10,
	models.Moon:      10,
	models.Mercury:   4,
	models.Venus:     5,
	models.Mars:      6,
	models.Jupiter:   7,
	models.Saturn:    8,
	models.Uranus:    5,
	models.Neptune:   4,
	models.Pluto:     6,
	models.NorthNode: 3,
	models.Chiron:    3,
}

// GetPlanetInfo 获取行星信息
func GetPlanetInfo(id models.PlanetID) *PlanetInfo {
	for _, p := range Planets {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

// ==================== 星座信息 ====================

// ZodiacInfo 星座信息
type ZodiacInfo struct {
	ID       models.ZodiacID
	Name     string
	Symbol   string
	Element  string // fire, earth, air, water
	Modality string // cardinal, fixed, mutable
	Ruler    models.PlanetID
}

// ZodiacSigns 星座列表
var ZodiacSigns = []ZodiacInfo{
	{models.Aries, "白羊座", "♈", "fire", "cardinal", models.Mars},
	{models.Taurus, "金牛座", "♉", "earth", "fixed", models.Venus},
	{models.Gemini, "双子座", "♊", "air", "mutable", models.Mercury},
	{models.Cancer, "巨蟹座", "♋", "water", "cardinal", models.Moon},
	{models.Leo, "狮子座", "♌", "fire", "fixed", models.Sun},
	{models.Virgo, "处女座", "♍", "earth", "mutable", models.Mercury},
	{models.Libra, "天秤座", "♎", "air", "cardinal", models.Venus},
	{models.Scorpio, "天蝎座", "♏", "water", "fixed", models.Pluto},
	{models.Sagittarius, "射手座", "♐", "fire", "mutable", models.Jupiter},
	{models.Capricorn, "摩羯座", "♑", "earth", "cardinal", models.Saturn},
	{models.Aquarius, "水瓶座", "♒", "air", "fixed", models.Uranus},
	{models.Pisces, "双鱼座", "♓", "water", "mutable", models.Neptune},
}

// GetZodiacInfo 获取星座信息
func GetZodiacInfo(id models.ZodiacID) *ZodiacInfo {
	for _, z := range ZodiacSigns {
		if z.ID == id {
			return &z
		}
	}
	return nil
}

// GetZodiacByLongitude 根据黄经获取星座
func GetZodiacByLongitude(longitude float64) *ZodiacInfo {
	// 确保黄经在0-360范围内
	lon := NormalizeAngle(longitude)
	index := int(lon / 30)
	if index >= 12 {
		index = 0
	}
	return &ZodiacSigns[index]
}

// ==================== 相位定义 ====================

// AspectDefinition 相位定义
type AspectDefinition struct {
	Type     models.AspectType
	Angle    float64
	Orb      float64
	Name     string
	Weight   float64
	Nature   string // harmonious, tense, neutral
}

// AspectDefinitions 相位定义列表
var AspectDefinitions = []AspectDefinition{
	{models.Conjunction, 0, 10, "合相", 10, "neutral"},
	{models.Sextile, 60, 6, "六分相", 3, "harmonious"},
	{models.Square, 90, 8, "四分相", 7, "tense"},
	{models.Trine, 120, 8, "三分相", 6, "harmonious"},
	{models.Opposition, 180, 10, "对分相", 8, "tense"},
}

// GetAspectDefinition 获取相位定义
func GetAspectDefinition(t models.AspectType) *AspectDefinition {
	for _, a := range AspectDefinitions {
		if a.Type == t {
			return &a
		}
	}
	return nil
}

// ==================== 宫位信息 ====================

// HouseInfo 宫位信息
type HouseInfo struct {
	House    int
	Name     string
	Theme    string
	Keywords []string
}

// HouseDataList 宫位数据列表
var HouseDataList = []HouseInfo{
	{1, "命宫", "自我与身份", []string{"身份", "外表", "人格面具", "生命力"}},
	{2, "财帛宫", "资源与价值", []string{"金钱", "物质", "自我价值", "天赋"}},
	{3, "兄弟宫", "沟通与学习", []string{"沟通", "学习", "兄弟姐妹", "短途旅行"}},
	{4, "田宅宫", "家庭与根基", []string{"家庭", "根源", "安全感", "房产"}},
	{5, "子女宫", "创造与表达", []string{"创造力", "恋爱", "子女", "娱乐"}},
	{6, "奴仆宫", "服务与健康", []string{"工作", "健康", "日常习惯", "服务"}},
	{7, "夫妻宫", "关系与合作", []string{"婚姻", "合作", "公开敌人", "合约"}},
	{8, "疾厄宫", "转化与共享", []string{"死亡重生", "共享资源", "心理深度", "性"}},
	{9, "迁移宫", "探索与智慧", []string{"高等教育", "远行", "哲学", "宗教"}},
	{10, "官禄宫", "事业与成就", []string{"事业", "公众形象", "成就", "权威"}},
	{11, "福德宫", "愿景与社群", []string{"朋友", "团体", "理想", "希望"}},
	{12, "玄秘宫", "内省与超越", []string{"潜意识", "灵性", "隐退", "业力"}},
}

// GetHouseInfo 获取宫位信息
func GetHouseInfo(house int) *HouseInfo {
	if house < 1 || house > 12 {
		return nil
	}
	return &HouseDataList[house-1]
}

// ==================== 尊贵度表 ====================

// DignityTable 尊贵度表
// key: PlanetID, value: map[dignity][]ZodiacID
var DignityTable = map[models.PlanetID]map[models.Dignity][]models.ZodiacID{
	models.Sun: {
		models.DignityDomicile:   {models.Leo},
		models.DignityExaltation: {models.Aries},
		models.DignityDetriment:  {models.Aquarius},
		models.DignityFall:       {models.Libra},
	},
	models.Moon: {
		models.DignityDomicile:   {models.Cancer},
		models.DignityExaltation: {models.Taurus},
		models.DignityDetriment:  {models.Capricorn},
		models.DignityFall:       {models.Scorpio},
	},
	models.Mercury: {
		models.DignityDomicile:  {models.Gemini, models.Virgo},
		models.DignityDetriment: {models.Sagittarius, models.Pisces},
	},
	models.Venus: {
		models.DignityDomicile:   {models.Taurus, models.Libra},
		models.DignityExaltation: {models.Pisces},
		models.DignityDetriment:  {models.Aries, models.Scorpio},
		models.DignityFall:       {models.Virgo},
	},
	models.Mars: {
		models.DignityDomicile:   {models.Aries, models.Scorpio},
		models.DignityExaltation: {models.Capricorn},
		models.DignityDetriment:  {models.Taurus, models.Libra},
		models.DignityFall:       {models.Cancer},
	},
	models.Jupiter: {
		models.DignityDomicile:   {models.Sagittarius, models.Pisces},
		models.DignityExaltation: {models.Cancer},
		models.DignityDetriment:  {models.Gemini, models.Virgo},
		models.DignityFall:       {models.Capricorn},
	},
	models.Saturn: {
		models.DignityDomicile:   {models.Capricorn, models.Aquarius},
		models.DignityExaltation: {models.Libra},
		models.DignityDetriment:  {models.Cancer, models.Leo},
		models.DignityFall:       {models.Aries},
	},
}

// GetDignity 获取行星在星座中的尊贵度
func GetDignity(planet models.PlanetID, sign models.ZodiacID) models.Dignity {
	if table, ok := DignityTable[planet]; ok {
		for dignity, signs := range table {
			for _, s := range signs {
				if s == sign {
					return dignity
				}
			}
		}
	}
	return models.DignityPeregrine
}

// GetDignityScore 获取尊贵度分数
func GetDignityScore(dignity models.Dignity) float64 {
	switch dignity {
	case models.DignityDomicile:
		return 3.0
	case models.DignityExaltation:
		return 2.0
	case models.DignityDetriment:
		return -2.0
	case models.DignityFall:
		return -3.0
	default:
		return 0.0
	}
}

// ==================== 默认因子权重 ====================

// DefaultFactorWeights 默认因子权重（可通过配置调整）
var DefaultFactorWeights = models.FactorWeights{
	Dignity:        1.0,
	Retrograde:     1.0,
	AspectPhase:    0.8,
	AspectOrb:      0.5,
	OuterPlanet:    1.2,
	ProfectionLord: 1.0,
	LunarPhase:     0.7,
	PlanetaryHour:  0.3,
	VoidOfCourse:   0.8,
	Personal:       1.0,
	Custom:         1.0,
}

// ==================== 月相名称 ====================

// LunarPhases 月相定义
var LunarPhases = []struct {
	Min      float64
	Max      float64
	Phase    string
	Name     string
	Keywords []string
}{
	{0, 45, "new", "新月期", []string{"新开始", "播种", "意图"}},
	{45, 90, "crescent", "新月后期", []string{"行动", "建设", "突破"}},
	{90, 135, "firstQuarter", "上弦月", []string{"危机", "调整", "决断"}},
	{135, 180, "gibbous", "凸月期", []string{"完善", "分析", "准备"}},
	{180, 225, "full", "满月期", []string{"高峰", "收获", "觉醒"}},
	{225, 270, "disseminating", "播种期", []string{"分享", "传播", "教导"}},
	{270, 315, "lastQuarter", "下弦月", []string{"反思", "释放", "调整"}},
	{315, 360, "balsamic", "残月期", []string{"整合", "休息", "准备"}},
}

// GetLunarPhase 根据角度获取月相信息
func GetLunarPhase(angle float64) models.LunarPhaseInfo {
	normalizedAngle := NormalizeAngle(angle)
	for _, lp := range LunarPhases {
		if normalizedAngle >= lp.Min && normalizedAngle < lp.Max {
			return models.LunarPhaseInfo{
				Phase:    lp.Phase,
				Name:     lp.Name,
				Angle:    normalizedAngle,
				Keywords: lp.Keywords,
			}
		}
	}
	// 默认返回新月
	return models.LunarPhaseInfo{
		Phase:    "new",
		Name:     "新月期",
		Angle:    normalizedAngle,
		Keywords: []string{"新开始", "播种", "意图"},
	}
}

// ==================== 月亮星座关键词 ====================

// MoonSignKeywords 月亮星座关键词
var MoonSignKeywords = map[models.ZodiacID][]string{
	models.Aries:       {"行动", "冲动", "直接"},
	models.Taurus:      {"稳定", "享受", "固执"},
	models.Gemini:      {"沟通", "学习", "多变"},
	models.Cancer:      {"情感", "家庭", "敏感"},
	models.Leo:         {"表达", "创造", "热情"},
	models.Virgo:       {"分析", "实用", "细节"},
	models.Libra:       {"和谐", "社交", "平衡"},
	models.Scorpio:     {"深度", "转化", "强烈"},
	models.Sagittarius: {"探索", "乐观", "自由"},
	models.Capricorn:   {"责任", "成就", "务实"},
	models.Aquarius:    {"创新", "独立", "人道"},
	models.Pisces:      {"灵性", "梦想", "同情"},
}

