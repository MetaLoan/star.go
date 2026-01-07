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
	{models.Sun, "Sun", "☉", "#ffd700", 10},
	{models.Moon, "Moon", "☽", "#c0c0c0", 10},
	{models.Mercury, "Mercury", "☿", "#b5651d", 4},
	{models.Venus, "Venus", "♀", "#ff69b4", 5},
	{models.Mars, "Mars", "♂", "#dc143c", 6},
	{models.Jupiter, "Jupiter", "♃", "#daa520", 7},
	{models.Saturn, "Saturn", "♄", "#8b7355", 8},
	{models.Uranus, "Uranus", "♅", "#40e0d0", 5},
	{models.Neptune, "Neptune", "♆", "#4169e1", 4},
	{models.Pluto, "Pluto", "♇", "#800080", 6},
	{models.NorthNode, "North Node", "☊", "#9932cc", 3},
	{models.Chiron, "Chiron", "⚷", "#228b22", 3},
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
	{models.Aries, "Aries", "♈", "fire", "cardinal", models.Mars},
	{models.Taurus, "Taurus", "♉", "earth", "fixed", models.Venus},
	{models.Gemini, "Gemini", "♊", "air", "mutable", models.Mercury},
	{models.Cancer, "Cancer", "♋", "water", "cardinal", models.Moon},
	{models.Leo, "Leo", "♌", "fire", "fixed", models.Sun},
	{models.Virgo, "Virgo", "♍", "earth", "mutable", models.Mercury},
	{models.Libra, "Libra", "♎", "air", "cardinal", models.Venus},
	{models.Scorpio, "Scorpio", "♏", "water", "fixed", models.Pluto},
	{models.Sagittarius, "Sagittarius", "♐", "fire", "mutable", models.Jupiter},
	{models.Capricorn, "Capricorn", "♑", "earth", "cardinal", models.Saturn},
	{models.Aquarius, "Aquarius", "♒", "air", "fixed", models.Uranus},
	{models.Pisces, "Pisces", "♓", "water", "mutable", models.Neptune},
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
	{models.Conjunction, 0, 10, "Conjunction", 10, "neutral"},
	{models.Sextile, 60, 6, "Sextile", 3, "harmonious"},
	{models.Square, 90, 8, "Square", 7, "tense"},
	{models.Trine, 120, 8, "Trine", 6, "harmonious"},
	{models.Opposition, 180, 10, "Opposition", 8, "tense"},
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
	{1, "1st House", "Self & Identity", []string{"identity", "appearance", "persona", "vitality"}},
	{2, "2nd House", "Resources & Values", []string{"money", "possessions", "self-worth", "talents"}},
	{3, "3rd House", "Communication & Learning", []string{"communication", "learning", "siblings", "short trips"}},
	{4, "4th House", "Home & Foundation", []string{"family", "roots", "security", "property"}},
	{5, "5th House", "Creativity & Expression", []string{"creativity", "romance", "children", "entertainment"}},
	{6, "6th House", "Service & Health", []string{"work", "health", "daily routines", "service"}},
	{7, "7th House", "Relationships & Partnership", []string{"marriage", "partnerships", "open enemies", "contracts"}},
	{8, "8th House", "Transformation & Shared Resources", []string{"death/rebirth", "shared resources", "psychology", "intimacy"}},
	{9, "9th House", "Exploration & Wisdom", []string{"higher education", "travel", "philosophy", "religion"}},
	{10, "10th House", "Career & Achievement", []string{"career", "public image", "achievement", "authority"}},
	{11, "11th House", "Vision & Community", []string{"friends", "groups", "ideals", "hopes"}},
	{12, "12th House", "Introspection & Transcendence", []string{"subconscious", "spirituality", "seclusion", "karma"}},
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
	{0, 45, "new", "New Moon", []string{"new beginnings", "seeding", "intentions"}},
	{45, 90, "crescent", "Crescent Moon", []string{"action", "building", "breakthrough"}},
	{90, 135, "firstQuarter", "First Quarter", []string{"crisis", "adjustment", "decision"}},
	{135, 180, "gibbous", "Gibbous Moon", []string{"refinement", "analysis", "preparation"}},
	{180, 225, "full", "Full Moon", []string{"peak", "harvest", "awakening"}},
	{225, 270, "disseminating", "Disseminating Moon", []string{"sharing", "spreading", "teaching"}},
	{270, 315, "lastQuarter", "Last Quarter", []string{"reflection", "release", "adjustment"}},
	{315, 360, "balsamic", "Balsamic Moon", []string{"integration", "rest", "preparation"}},
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
		Name:     "New Moon",
		Angle:    normalizedAngle,
		Keywords: []string{"new beginnings", "seeding", "intentions"},
	}
}

// ==================== 月亮星座关键词 ====================

// MoonSignKeywords 月亮星座关键词
var MoonSignKeywords = map[models.ZodiacID][]string{
	models.Aries:       {"action", "impulse", "directness"},
	models.Taurus:      {"stability", "enjoyment", "persistence"},
	models.Gemini:      {"communication", "learning", "adaptability"},
	models.Cancer:      {"emotion", "family", "sensitivity"},
	models.Leo:         {"expression", "creativity", "passion"},
	models.Virgo:       {"analysis", "practicality", "detail"},
	models.Libra:       {"harmony", "sociability", "balance"},
	models.Scorpio:     {"depth", "transformation", "intensity"},
	models.Sagittarius: {"exploration", "optimism", "freedom"},
	models.Capricorn:   {"responsibility", "achievement", "pragmatism"},
	models.Aquarius:    {"innovation", "independence", "humanitarianism"},
	models.Pisces:      {"spirituality", "dreams", "compassion"},
}

