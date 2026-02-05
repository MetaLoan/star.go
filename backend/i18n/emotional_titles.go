package i18n

import (
	"star/models"
	"strings"
)

// GetEmotionalTitle returns an emotional/psychological title for an event
func (t *Translator) GetEmotionalTitle(eventType string, planet1, planet2 models.PlanetID, aspect, house string, isPositive bool) string {
	// For transit house events
	if eventType == "transit_house" {
		return t.getTransitHouseTitle(planet1, house)
	}
	
	// For progression events
	if eventType == "secondary_progression" || eventType == "tertiary_progression" {
		return t.getProgressionTitle(planet1, planet2, aspect, isPositive)
	}
	
	// For aspect events
	if eventType == "aspect" {
		return t.getAspectTitle(planet1, planet2, aspect, isPositive)
	}
	
	// For retrograde events
	if eventType == "retrograde" {
		return t.getRetrogradeTitle(planet1)
	}
	
	// For annual lord (profection) events
	if eventType == "profectionLord" {
		return t.getProfectionLordTitle(planet1)
	}
	
	// For planetary hour (hourly); API may send "planetary_hour_change" or "planetaryHour"
	if eventType == "planetaryHour" || eventType == "planetary_hour_change" {
		return t.getPlanetaryHourTitle(planet1)
	}
	
	// For Moon void of course (hourly)
	if eventType == "voidOfCourse" {
		return t.getVoidOfCourseTitle()
	}
	
	// For lunar phase events
	if eventType == "lunar_phase" {
		return t.getLunarPhaseTitle(aspect)
	}
	
	// For sign change events
	if eventType == "sign_change" {
		return t.getSignChangeTitle(planet1, aspect)
	}
	
	// For dignity events
	if eventType == "dignity" {
		return t.getDignityTitle(planet1, aspect)
	}
	
	return ""
}

// getAspectTitle returns emotional title for aspects
func (t *Translator) getAspectTitle(p1, p2 models.PlanetID, aspect string, isPositive bool) string {
	key := string(p1) + "_" + aspect + "_" + string(p2)
	
	switch t.lang {
	case Chinese:
		return getChineseAspectTitle(key, isPositive)
	case Russian:
		return getRussianAspectTitle(key, isPositive)
	default:
		return getEnglishAspectTitle(key, isPositive)
	}
}

// getProgressionTitle returns emotional title for progressions
func (t *Translator) getProgressionTitle(p1, p2 models.PlanetID, aspect string, isPositive bool) string {
	key := string(p1) + "_" + aspect + "_" + string(p2)
	
	switch t.lang {
	case Chinese:
		return getChineseProgressionTitle(key, isPositive)
	case Russian:
		return getRussianProgressionTitle(key, isPositive)
	default:
		return getEnglishProgressionTitle(key, isPositive)
	}
}

// getTransitHouseTitle returns emotional title for transit houses
func (t *Translator) getTransitHouseTitle(planet models.PlanetID, house string) string {
	key := string(planet) + "_house_" + house
	
	switch t.lang {
	case Chinese:
		return getChineseTransitHouseTitle(key)
	case Russian:
		return getRussianTransitHouseTitle(key)
	default:
		return getEnglishTransitHouseTitle(key)
	}
}

// getRetrogradeTitle returns emotional title for retrograde events (planet = retrograde planet)
func (t *Translator) getRetrogradeTitle(planet models.PlanetID) string {
	switch t.lang {
	case Chinese:
		return getChineseRetrogradeTitle(planet)
	case Russian:
		return getRussianRetrogradeTitle(planet)
	default:
		return getEnglishRetrogradeTitle(planet)
	}
}

// getProfectionLordTitle returns emotional title for annual lord events (planet = lord planet)
func (t *Translator) getProfectionLordTitle(planet models.PlanetID) string {
	switch t.lang {
	case Chinese:
		return getChineseProfectionLordTitle(planet)
	case Russian:
		return getRussianProfectionLordTitle(planet)
	default:
		return getEnglishProfectionLordTitle(planet)
	}
}

func (t *Translator) getPlanetaryHourTitle(planet models.PlanetID) string {
	switch t.lang {
	case Chinese:
		return getChinesePlanetaryHourTitle(planet)
	case Russian:
		return getRussianPlanetaryHourTitle(planet)
	default:
		return getEnglishPlanetaryHourTitle(planet)
	}
}

func (t *Translator) getVoidOfCourseTitle() string {
	switch t.lang {
	case Chinese:
		return getChineseVoidOfCourseTitle()
	case Russian:
		return getRussianVoidOfCourseTitle()
	default:
		return getEnglishVoidOfCourseTitle()
	}
}

func (t *Translator) getLunarPhaseTitle(phase string) string {
	switch t.lang {
	case Chinese:
		return getChineseLunarPhaseTitle(phase)
	case Russian:
		return getRussianLunarPhaseTitle(phase)
	default:
		return getEnglishLunarPhaseTitle(phase)
	}
}

func (t *Translator) getSignChangeTitle(planet models.PlanetID, newSign string) string {
	switch t.lang {
	case Chinese:
		return getChineseSignChangeTitle(planet, newSign)
	case Russian:
		return getRussianSignChangeTitle(planet, newSign)
	default:
		return getEnglishSignChangeTitle(planet, newSign)
	}
}

func (t *Translator) getDignityTitle(planet models.PlanetID, dignityType string) string {
	switch t.lang {
	case Chinese:
		return getChineseDignityTitle(planet, dignityType)
	case Russian:
		return getRussianDignityTitle(planet, dignityType)
	default:
		return getEnglishDignityTitle(planet, dignityType)
	}
}

// Chinese emotional titles
func getChineseAspectTitle(key string, isPositive bool) string {
	titles := map[string]string{
		// Returns / self-conjunctions
		"sun_conjunction_sun":       "太阳回归·新岁启程",
		"moon_conjunction_moon":    "情绪与内在需求唤起",
		"mercury_conjunction_mercury": "思维与沟通强化",
		"venus_conjunction_venus": "金星回归·爱与美",
		"mars_conjunction_mars":   "行动力与斗志高峰",
		"jupiter_conjunction_jupiter": "木星回归·扩展之年",
		"saturn_conjunction_saturn": "土星回归·成熟之期",
		// Sun aspects
		"sun_conjunction_moon": "情感与意志的融合",
		"sun_trine_moon":       "内心和谐时光",
		"sun_sextile_moon":     "情绪舒畅",
		"sun_square_moon":      "内心矛盾",
		"sun_opposition_moon":  "情感拉扯",
		
		"sun_conjunction_venus": "魅力绽放",
		"sun_trine_venus":       "充满爱的时刻",
		"sun_sextile_venus":     "美好事物涌现",
		
		"sun_conjunction_mars": "勇敢行动",
		"sun_trine_mars":       "所向披靡的勇士",
		"sun_sextile_mars":     "充满活力",
		"sun_square_mars":      "冲动易怒",
		
		"sun_conjunction_jupiter": "幸运降临",
		"sun_trine_jupiter":       "扩展与成长",
		"sun_sextile_jupiter":     "机遇涌现",
		"sun_square_jupiter":      "过度乐观",
		"sun_opposition_jupiter":  "理想与现实拉扯",
		"sun_conjunction_uranus": "突然的自我觉醒",
		"sun_trine_uranus":        "创意与突破",
		"sun_conjunction_neptune": "灵感与慈悲",
		"sun_conjunction_pluto":   "深层的自我转化",
		"sun_conjunction_saturn": "面对责任",
		"sun_square_saturn":      "压力与考验",
		"sun_trine_saturn":       "稳健前行",
		
		// Moon aspects
		"moon_conjunction_venus": "亲密时光",
		"moon_trine_venus":       "情感满足",
		"moon_sextile_venus":     "温柔体贴",
		"moon_square_venus":      "情感波动",
		
		"moon_conjunction_mars": "情绪激动",
		"moon_square_mars":      "易怒烦躁",
		"moon_trine_mars":       "情感勇气",
		
		"moon_conjunction_saturn": "情感沉重",
		"moon_square_saturn":      "孤独压抑",
		"moon_trine_saturn":       "情感成熟",
		"moon_conjunction_jupiter": "情感扩展",
		"moon_trine_jupiter":      "内心充满希望",
		"moon_sextile_jupiter":    "希望与扩展的美好期",
		"moon_square_jupiter":     "过度乐观",
		"moon_conjunction_uranus": "情绪突变",
		"moon_conjunction_neptune": "感受力极强",
		"moon_conjunction_pluto":  "深层情绪触动",
		"moon_conjunction_northNode": "情感与命运方向对齐",
		"moon_conjunction_chiron": "旧伤或脆弱触及",
		// Venus aspects
		"venus_conjunction_mars": "欲望灼烧",
		"venus_trine_mars":       "激情与魅力",
		"venus_sextile_mars":     "爱情火花",
		"venus_square_mars":      "情欲纠葛",
		
		"venus_conjunction_jupiter": "社交魅力大绽放",
		"venus_trine_jupiter":       "快乐与幸福",
		"venus_sextile_jupiter":     "美好机遇",
		
		// Mars aspects
		"mars_conjunction_jupiter": "积极勇敢",
		"mars_trine_jupiter":       "充满能量",
		"mars_sextile_jupiter":     "行动力强",
		"mars_square_jupiter":      "过度冒险",
		
		"mars_conjunction_saturn": "行动受阻",
		"mars_square_saturn":      "坚韧耐力",
		"mars_trine_saturn":       "有纪律的行动",
		
		// Jupiter aspects
		"jupiter_conjunction_saturn": "稳健扩张",
		"jupiter_square_saturn":      "限制与机遇",
		"jupiter_trine_saturn":       "智慧成长",
		
		// Mercury aspects
		"mercury_conjunction_venus": "优雅表达",
		"mercury_trine_venus":       "沟通顺畅",
		"mercury_sextile_venus":     "愉快交流",
		
		"mercury_conjunction_mars": "言辞激烈",
		"mercury_square_mars":      "争论不休",
		"mercury_trine_mars":       "思维敏捷",
		
		"mercury_conjunction_jupiter": "思维开阔",
		"mercury_trine_jupiter":       "交流与忙碌",
		"mercury_sextile_jupiter":     "好消息传来",
		"mercury_conjunction_saturn": "严肃思考",
		"mercury_square_saturn":       "表达受阻",
		"mercury_conjunction_uranus": "突发灵感",
		"venus_conjunction_saturn":   "关系或价值考验",
		"venus_square_saturn":        "爱情考验",
		"venus_conjunction_uranus":  "突然的心动",
		"venus_conjunction_neptune": "浪漫与灵性",
		"venus_conjunction_pluto":   "深刻吸引",
		"mars_opposition_jupiter":   "冲与收的拉扯",
		"mars_conjunction_pluto":    "强烈意志",
		"jupiter_opposition_saturn":  "要更多与要更稳",
		"jupiter_conjunction_uranus": "突变机遇",
		"jupiter_conjunction_pluto": "深刻扩张",
		"sun_conjunction_northNode": "自我与命运方向对齐",
		"sun_conjunction_chiron":   "伤口与自我被看见",
		"venus_conjunction_northNode": "爱与价值往对的方向",
		"mars_conjunction_northNode": "行动与命运一致",
		"jupiter_conjunction_northNode": "扩张在对的方向",
		"saturn_conjunction_northNode": "责任在对的地方",
		"venus_conjunction_chiron":  "关系或价值中的伤",
		"mars_conjunction_chiron":   "行动或愤怒下的伤",
		"jupiter_conjunction_chiron": "成长通过伤与智慧",
		"saturn_conjunction_chiron": "责任与旧伤相遇",
		
		// Minor aspects
		"sun_semi-sextile_moon": "微妙的和谐",
		"mercury_semi-square_mars": "言语摩擦",
		"venus_sesquiquadrate_jupiter": "社交过度",
		"mars_quincunx_saturn": "行动受阻",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	// Default titles based on positivity
	if isPositive {
		return "和谐能量"
	}
	return "挑战时刻"
}

// getChineseProgressionTitle returns emotional title for progressions in Chinese.
func getChineseProgressionTitle(key string, isPositive bool) string {
	titles := map[string]string{
		"moon_trine_venus":  "亲密时光",
		"moon_square_venus": "情感调整",
		"moon_trine_mars":   "情感勇气",
		"moon_square_mars":  "情绪波动",
		"moon_trine_jupiter": "情感扩展",
		"moon_square_saturn": "情感成熟",
		
		"sun_sextile_sun":      "自我认同",
		"sun_trine_sun":        "生命力旺盛",
		"sun_square_sun":       "自我挑战",
		"sun_opposition_saturn": "责任与压力",
		"sun_trine_jupiter":     "自信成长",
		
		"mercury_sextile_venus":  "优雅沟通",
		"mercury_trine_jupiter":  "思维活跃",
		"mercury_square_saturn":  "深度思考",
		
		"venus_conjunction_sun":  "魅力绽放",
		"venus_trine_jupiter":    "快乐社交",
		"venus_square_saturn":    "爱情考验",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}

	// Fallback to aspect title with progression context
	aspectTitle := getChineseAspectTitle(key, isPositive)
	if aspectTitle != "和谐能量" && aspectTitle != "挑战时刻" {
		return aspectTitle + " (长期趋势)"
	}
	
	if isPositive {
		return "积极发展"
	}
	return "成长挑战"
}

func getChineseTransitHouseTitle(key string) string {
	titles := map[string]string{
		"sun_house_1":  "自我觉醒",
		"sun_house_2":  "财富聚焦",
		"sun_house_3":  "表达沟通",
		"sun_house_4":  "家庭温暖",
		"sun_house_5":  "创造喜悦",
		"sun_house_6":  "健康关注",
		"sun_house_7":  "伴侣互动",
		"sun_house_8":  "深度转化",
		"sun_house_9":  "视野开阔",
		"sun_house_10": "事业高光",
		"sun_house_11": "社交活跃",
		"sun_house_12": "内在探索",
		
		"moon_house_1": "听从内心的声音",
		"moon_house_2": "情感安全感",
		"moon_house_3": "情绪表达",
		"moon_house_4": "家的归属",
		"moon_house_5": "情感喜悦",
		"moon_house_6": "身心调理",
		"moon_house_7": "情感需求",
		"moon_house_8": "情绪深潜",
		"moon_house_9": "情感探索",
		"moon_house_10": "公众形象",
		"moon_house_11": "友谊温暖",
		"moon_house_12": "潜意识涌现",
		
		"venus_house_1":  "魅力提升",
		"venus_house_2":  "享受物质",
		"venus_house_3":  "愉快交流",
		"venus_house_4":  "家居美化",
		"venus_house_5":  "浪漫爱情",
		"venus_house_6":  "工作和谐",
		"venus_house_7":  "关系甜蜜",
		"venus_house_8":  "亲密连接",
		"venus_house_9":  "文化欣赏",
		"venus_house_10": "职场人缘",
		"venus_house_11": "社交魅力大绽放",
		"venus_house_12": "精神之爱",
		
		"mars_house_1":  "行动力爆发",
		"mars_house_2":  "赚钱冲劲",
		"mars_house_3":  "言辞有力",
		"mars_house_4":  "家务行动",
		"mars_house_5":  "竞技精神",
		"mars_house_6":  "工作拼搏",
		"mars_house_7":  "关系张力",
		"mars_house_8":  "欲望强烈",
		"mars_house_9":  "冒险旅行",
		"mars_house_10": "事业野心",
		"mars_house_11": "积极勇敢",
		"mars_house_12": "隐藏行动",
		
		"jupiter_house_1":  "自信扩张",
		"jupiter_house_2":  "财富与机遇",
		"jupiter_house_3":  "学习机会",
		"jupiter_house_4":  "家庭幸福",
		"jupiter_house_5":  "创作丰盛",
		"jupiter_house_6":  "健康改善",
		"jupiter_house_7":  "关系增长",
		"jupiter_house_8":  "转化的祝福",
		"jupiter_house_9":  "智慧启蒙",
		"jupiter_house_10": "事业机遇",
		"jupiter_house_11": "贵人现身，请抓住机会",
		"jupiter_house_12": "灵性成长",
		
		"saturn_house_1":  "自律建设",
		"saturn_house_2":  "财务责任",
		"saturn_house_3":  "严肃沟通",
		"saturn_house_4":  "家庭责任",
		"saturn_house_5":  "创作纪律",
		"saturn_house_6":  "健康警示",
		"saturn_house_7":  "关系考验",
		"saturn_house_8":  "深度整顿",
		"saturn_house_9":  "哲学思考",
		"saturn_house_10": "事业建构",
		"saturn_house_11": "社交筛选",
		"saturn_house_12": "业力清理",

		"mercury_house_1":  "表达与思考",
		"mercury_house_2":  "财务与价值思考",
		"mercury_house_3":  "沟通与学习",
		"mercury_house_4":  "家庭沟通",
		"mercury_house_5":  "创意表达",
		"mercury_house_6":  "工作与健康沟通",
		"mercury_house_7":  "伴侣沟通",
		"mercury_house_8":  "深度交流",
		"mercury_house_9":  "智慧与远见",
		"mercury_house_10": "事业表达",
		"mercury_house_11": "社交与理想沟通",
		"mercury_house_12": "内在与灵感",

		"uranus_house_1":  "自我革新",
		"uranus_house_2":  "财务变革",
		"uranus_house_3":  "沟通创新",
		"uranus_house_4":  "家庭变革",
		"uranus_house_5":  "创意突破",
		"uranus_house_6":  "工作与健康革新",
		"uranus_house_7":  "关系突变",
		"uranus_house_8":  "深度与资源变革",
		"uranus_house_9":  "信念与视野突破",
		"uranus_house_10": "事业革新",
		"uranus_house_11": "社群与理想变革",
		"uranus_house_12": "内在觉醒",

		"neptune_house_1":  "灵性自我",
		"neptune_house_2":  "财务与价值直觉",
		"neptune_house_3":  "灵感沟通",
		"neptune_house_4":  "家庭与归属感",
		"neptune_house_5":  "创意与灵感",
		"neptune_house_6": "身心与灵性调理",
		"neptune_house_7":  "关系与共情",
		"neptune_house_8":  "深度共情与转化",
		"neptune_house_9":  "信念与灵性探索",
		"neptune_house_10": "事业与理想化",
		"neptune_house_11": "社群与愿景",
		"neptune_house_12": "灵性沉淀",

		"pluto_house_1":  "自我转化",
		"pluto_house_2":  "财富与价值深度",
		"pluto_house_3":  "沟通深化",
		"pluto_house_4":  "家庭与根源转化",
		"pluto_house_5":  "创意与欲望深度",
		"pluto_house_6": "工作与健康转化",
		"pluto_house_7":  "关系深度转化",
		"pluto_house_8":  "深度与重生",
		"pluto_house_9":  "信念与哲学转化",
		"pluto_house_10": "事业与权力转化",
		"pluto_house_11": "社群与集体转化",
		"pluto_house_12": "业力与潜意识",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	return "过宫影响"
}

// English emotional titles
func getEnglishAspectTitle(key string, isPositive bool) string {
	titles := map[string]string{
		// Returns / self-conjunctions
		"sun_conjunction_sun":       "Solar Return · Fresh Start",
		"moon_conjunction_moon":    "Emotions & Inner Needs",
		"mercury_conjunction_mercury": "Mind & Communication",
		"venus_conjunction_venus": "Venus Return · Love & Beauty",
		"mars_conjunction_mars":   "Drive & Spirit Peak",
		"jupiter_conjunction_jupiter": "Jupiter Return · Year of Expansion",
		"saturn_conjunction_saturn": "Saturn Return · Maturity",
		"sun_conjunction_moon": "Emotional Integration",
		"sun_trine_moon":       "Inner Harmony",
		"sun_sextile_moon":     "Emotional Flow",
		"sun_square_moon":      "Inner Conflict",
		
		"sun_conjunction_venus": "Radiant Charm",
		"sun_trine_venus":       "Loving Moments",
		"sun_sextile_venus":     "Beauty Emerges",
		
		"sun_conjunction_mars": "Bold Action",
		"sun_trine_mars":       "Unstoppable Warrior",
		"sun_sextile_mars":     "Energetic Flow",
		"sun_square_mars":      "Impulsive Tension",
		
		"sun_conjunction_jupiter": "Fortune Arrives",
		"sun_trine_jupiter":       "Growth & Expansion",
		"sun_sextile_jupiter":     "Opportunities Emerge",
		"sun_square_jupiter":      "Over-Optimism",
		"sun_opposition_jupiter":  "Ideal vs. Reality",
		"sun_conjunction_uranus":  "Sudden Self-Awakening",
		"sun_trine_uranus":        "Creativity & Breakthrough",
		"sun_conjunction_neptune": "Inspiration & Compassion",
		"sun_conjunction_pluto":  "Deep Self-Transformation",
		"moon_conjunction_venus": "Intimate Moments",
		"moon_trine_venus":       "Emotional Fulfillment",
		"moon_square_venus":      "Emotional Fluctuation",
		"moon_conjunction_jupiter": "Emotional Expansion",
		"moon_trine_jupiter":     "Inner Hope",
		"moon_sextile_jupiter":   "Hope & Expansion",
		"moon_square_jupiter":    "Over-Optimism",
		"moon_conjunction_uranus": "Sudden Mood Shifts",
		"moon_conjunction_neptune": "Heightened Sensitivity",
		"moon_conjunction_pluto": "Deep Emotion Stirred",
		"moon_conjunction_northNode": "Emotion Aligns with Destiny",
		"moon_conjunction_chiron": "Old Wounds Touched",
		"venus_conjunction_mars": "Burning Desire",
		"venus_trine_mars":       "Passion & Charm",
		"venus_sextile_mars":     "Romantic Spark",
		
		"venus_conjunction_jupiter": "Social Charisma Blooms",
		"venus_trine_jupiter":       "Joy & Happiness",
		"venus_conjunction_saturn": "Love or Values Tested",
		"venus_square_saturn":      "Love Tested",
		"venus_conjunction_uranus": "Sudden Attraction",
		"venus_conjunction_neptune": "Romance & Spirituality",
		"venus_conjunction_pluto": "Deep Attraction",
		"mars_trine_jupiter": "Energetic Flow",
		"mars_square_saturn": "Resilient Endurance",
		"mars_opposition_jupiter": "Push vs. Hold Back",
		"mars_conjunction_pluto": "Strong Will",
		"mercury_trine_jupiter": "Communication & Activity",
		"mercury_conjunction_saturn": "Serious Thought",
		"mercury_square_saturn":   "Expression Blocked",
		"mercury_conjunction_uranus": "Sudden Inspiration",
		"jupiter_conjunction_saturn": "Steady Expansion",
		"jupiter_opposition_saturn": "More vs. Stability",
		"jupiter_conjunction_uranus": "Sudden Opportunity",
		"jupiter_conjunction_pluto": "Deep Expansion",
		"sun_conjunction_northNode": "Self Aligns with Destiny",
		"sun_conjunction_chiron":  "Wound & Self Seen",
		"venus_conjunction_northNode": "Love Toward Right Direction",
		"mars_conjunction_northNode": "Action & Destiny Align",
		"jupiter_conjunction_northNode": "Expansion in Right Direction",
		"saturn_conjunction_northNode": "Responsibility in Right Place",
		"venus_conjunction_chiron": "Wound in Relationship or Values",
		"mars_conjunction_chiron": "Wound Under Action or Anger",
		"jupiter_conjunction_chiron": "Growth Through Wound & Wisdom",
		"saturn_conjunction_chiron": "Duty Meets Old Wound",

		// Minor aspects
		"sun_semi-sextile_moon": "Subtle Harmony",
		"mercury_semi-square_mars": "Verbal Friction",
		"venus_sesquiquadrate_jupiter": "Social Excess",
		"mars_quincunx_saturn": "Action Incoordination",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	if isPositive {
		return "Harmonious Energy"
	}
	return "Challenging Moment"
}

func getEnglishProgressionTitle(key string, isPositive bool) string {
	titles := map[string]string{
		"moon_trine_venus":  "Intimate Moments",
		"sun_sextile_sun":   "Self-Recognition",
		"venus_conjunction_sun": "Radiant Charm",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}

	// Fallback to aspect title with progression context
	aspectTitle := getEnglishAspectTitle(key, isPositive)
	if aspectTitle != "Harmonious Energy" && aspectTitle != "Challenging Moment" {
		return aspectTitle + " (Long-term)"
	}
	
	if isPositive {
		return "Positive Development"
	}
	return "Growth Challenge"
}

func getEnglishTransitHouseTitle(key string) string {
	titles := map[string]string{
		"sun_house_1":  "Self-Awakening",
		"sun_house_2":  "Wealth Focus",
		"sun_house_3":  "Expression & Communication",
		"sun_house_4":  "Family Warmth",
		"sun_house_5":  "Creation & Joy",
		"sun_house_6":  "Health Focus",
		"sun_house_7":  "Partnership",
		"sun_house_8":  "Depth & Transformation",
		"sun_house_9":  "Wider Horizons",
		"sun_house_10": "Career Spotlight",
		"sun_house_11": "Social Activity",
		"sun_house_12": "Inner Exploration",
		"moon_house_1":  "Listen to Inner Voice",
		"moon_house_2":  "Emotional Security",
		"moon_house_3":  "Emotional Expression",
		"moon_house_4":  "Sense of Belonging",
		"moon_house_5":  "Emotional Joy",
		"moon_house_6":  "Body & Mind Care",
		"moon_house_7":  "Emotional Needs",
		"moon_house_8":  "Emotional Depth",
		"moon_house_9":  "Emotional Exploration",
		"moon_house_10": "Public Image",
		"moon_house_11": "Friendship Warmth",
		"moon_house_12": "Subconscious Surfacing",
		"venus_house_1":  "Charm Boost",
		"venus_house_2":  "Enjoying Matter",
		"venus_house_3":  "Pleasant Exchange",
		"venus_house_4":  "Home Beauty",
		"venus_house_5":  "Romance",
		"venus_house_6":  "Work Harmony",
		"venus_house_7":  "Sweet Relations",
		"venus_house_8":  "Intimate Connection",
		"venus_house_9":  "Cultural Appreciation",
		"venus_house_10": "Workplace Charm",
		"venus_house_11": "Social Charisma Blooms",
		"venus_house_12": "Spiritual Love",
		"mars_house_1":  "Drive Surge",
		"mars_house_2":  "Earning Drive",
		"mars_house_3":  "Forceful Words",
		"mars_house_4":  "Home Action",
		"mars_house_5":  "Competitive Spirit",
		"mars_house_6":  "Work Push",
		"mars_house_7":  "Relationship Tension",
		"mars_house_8":  "Strong Desire",
		"mars_house_9":  "Adventure & Travel",
		"mars_house_10": "Career Ambition",
		"mars_house_11": "Bold & Brave",
		"mars_house_12": "Hidden Action",
		"jupiter_house_1":  "Confidence Expansion",
		"jupiter_house_2":  "Wealth & Opportunity",
		"jupiter_house_3":  "Learning Opportunity",
		"jupiter_house_4":  "Family Joy",
		"jupiter_house_5":  "Creative Abundance",
		"jupiter_house_6":  "Health Improvement",
		"jupiter_house_7":  "Relationship Growth",
		"jupiter_house_8":  "Blessing of Transformation",
		"jupiter_house_9":  "Wisdom Enlightenment",
		"jupiter_house_10": "Career Opportunity",
		"jupiter_house_11": "Benefactors Appear",
		"jupiter_house_12": "Spiritual Growth",
		"saturn_house_1":  "Discipline & Structure",
		"saturn_house_2":  "Financial Responsibility",
		"saturn_house_3":  "Serious Communication",
		"saturn_house_4":  "Family Responsibility",
		"saturn_house_5":  "Creative Discipline",
		"saturn_house_6":  "Health Caution",
		"saturn_house_7":  "Relationship Test",
		"saturn_house_8":  "Depth & Order",
		"saturn_house_9":  "Philosophical Thought",
		"saturn_house_10": "Career Structure",
		"saturn_house_11": "Social Filtering",
		"saturn_house_12": "Karma Clearing",

		"mercury_house_1":  "Expression & Thought",
		"mercury_house_2":  "Values & Money Thought",
		"mercury_house_3":  "Communication & Learning",
		"mercury_house_4":  "Family Communication",
		"mercury_house_5":  "Creative Expression",
		"mercury_house_6":  "Work & Health Communication",
		"mercury_house_7":  "Partner Communication",
		"mercury_house_8":  "Deep Exchange",
		"mercury_house_9":  "Wisdom & Vision",
		"mercury_house_10": "Career Expression",
		"mercury_house_11": "Social & Ideal Communication",
		"mercury_house_12": "Inner & Inspiration",

		"uranus_house_1":  "Self Revolution",
		"uranus_house_2":  "Financial Change",
		"uranus_house_3":  "Communication Innovation",
		"uranus_house_4":  "Family Revolution",
		"uranus_house_5":  "Creative Breakthrough",
		"uranus_house_6":  "Work & Health Revolution",
		"uranus_house_7":  "Relationship Sudden Shift",
		"uranus_house_8":  "Depth & Resource Change",
		"uranus_house_9":  "Belief & Vision Breakthrough",
		"uranus_house_10": "Career Revolution",
		"uranus_house_11": "Community & Ideal Change",
		"uranus_house_12": "Inner Awakening",

		"neptune_house_1":  "Spiritual Self",
		"neptune_house_2":  "Values & Intuition",
		"neptune_house_3":  "Inspired Communication",
		"neptune_house_4":  "Home & Belonging",
		"neptune_house_5":  "Creativity & Inspiration",
		"neptune_house_6":  "Body, Mind & Spirit",
		"neptune_house_7":  "Relationship & Empathy",
		"neptune_house_8":  "Deep Empathy & Transformation",
		"neptune_house_9":  "Faith & Spiritual Quest",
		"neptune_house_10": "Career & Idealization",
		"neptune_house_11": "Community & Vision",
		"neptune_house_12": "Spiritual Settling",

		"pluto_house_1":  "Self Transformation",
		"pluto_house_2":  "Wealth & Value Depth",
		"pluto_house_3":  "Communication Depth",
		"pluto_house_4":  "Family & Roots Transformation",
		"pluto_house_5":  "Creativity & Desire Depth",
		"pluto_house_6":  "Work & Health Transformation",
		"pluto_house_7":  "Relationship Deep Transformation",
		"pluto_house_8":  "Depth & Rebirth",
		"pluto_house_9":  "Belief & Philosophy Transformation",
		"pluto_house_10": "Career & Power Transformation",
		"pluto_house_11": "Community & Collective Transformation",
		"pluto_house_12": "Karma & Unconscious",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	return "House Transit"
}

// Russian emotional titles
func getRussianAspectTitle(key string, isPositive bool) string {
	titles := map[string]string{
		// Returns / self-conjunctions
		"sun_conjunction_sun":       "Солнечное возвращение · Новый цикл",
		"moon_conjunction_moon":    "Эмоции и внутренние потребности",
		"mercury_conjunction_mercury": "Ум и общение",
		"venus_conjunction_venus": "Возвращение Венеры · Любовь и красота",
		"mars_conjunction_mars":   "Энергия и боевой дух на пике",
		"jupiter_conjunction_jupiter": "Возвращение Юпитера · Расширение",
		"saturn_conjunction_saturn": "Возвращение Сатурна · Зрелость",
		"sun_conjunction_moon": "Эмоциональная интеграция",
		"sun_trine_moon":       "Внутренняя гармония",
		"sun_conjunction_venus": "Сияющее очарование",
		"sun_trine_mars":       "Неудержимый воин",
		"sun_conjunction_jupiter": "Удача приходит",
		"sun_trine_jupiter":    "Рост и расширение",
		"sun_sextile_jupiter":  "Возможности появляются",
		"sun_conjunction_uranus": "Внезапное пробуждение себя",
		"sun_conjunction_neptune": "Вдохновение и сострадание",
		"sun_conjunction_pluto": "Глубинное преображение себя",
		"moon_conjunction_venus": "Интимные моменты",
		"moon_trine_venus":     "Эмоциональное удовлетворение",
		"moon_conjunction_jupiter": "Эмоциональное расширение",
		"moon_trine_jupiter":   "Внутренняя надежда",
		"moon_sextile_jupiter": "Надежда и расширение",
		"moon_conjunction_uranus": "Резкие смены настроения",
		"moon_conjunction_neptune": "Высокая чувствительность",
		"moon_conjunction_pluto": "Глубокие эмоции затронуты",
		"moon_conjunction_northNode": "Эмоция совпадает с судьбой",
		"moon_conjunction_chiron": "Старые раны затронуты",
		"venus_conjunction_mars": "Пылающее желание",
		"venus_conjunction_jupiter": "Расцвет социальной харизмы",
		"venus_conjunction_saturn": "Любовь или ценности под испытанием",
		"venus_conjunction_uranus": "Внезапное влечение",
		"venus_conjunction_neptune": "Романтика и духовность",
		"venus_conjunction_pluto": "Глубокое влечение",
		"mars_trine_jupiter": "Полны энергии",
		"mars_square_saturn": "Упругая выносливость",
		"mars_conjunction_pluto": "Сильная воля",
		"mercury_trine_jupiter": "Общение и активность",
		"jupiter_conjunction_saturn": "Устойчивое расширение",
		"jupiter_conjunction_uranus": "Внезапная возможность",
		"jupiter_conjunction_pluto": "Глубокое расширение",
		"sun_conjunction_northNode": "Самость совпадает с судьбой",
		"sun_conjunction_chiron": "Рана и самость видимы",
		"venus_conjunction_northNode": "Любовь в верном направлении",
		"mars_conjunction_northNode": "Действие и судьба совпадают",
		"jupiter_conjunction_northNode": "Расширение в верном направлении",
		"saturn_conjunction_northNode": "Ответственность в верном месте",
		"venus_conjunction_chiron": "Рана в отношениях или ценностях",
		"mars_conjunction_chiron": "Рана под действием или гневом",
		"jupiter_conjunction_chiron": "Рост через рану и мудрость",
		"saturn_conjunction_chiron": "Долг встречает старую рану",

		// Minor aspects
		"sun_semi-sextile_moon": "Тонкая гармония",
		"mercury_semi-square_mars": "Словесные трения",
		"venus_sesquiquadrate_jupiter": "Социальный избыток",
		"mars_quincunx_saturn": "Несогласованность действий",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	if isPositive {
		return "Гармоничная энергия"
	}
	return "Сложный момент"
}

func getRussianProgressionTitle(key string, isPositive bool) string {
	titles := map[string]string{
		"moon_trine_venus": "Интимные моменты",
		"sun_sextile_sun":  "Самопознание",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}

	// Fallback to aspect title with progression context
	aspectTitle := getRussianAspectTitle(key, isPositive)
	if aspectTitle != "Гармоничная энергия" && aspectTitle != "Сложный момент" {
		return aspectTitle + " (Длительный)"
	}
	
	if isPositive {
		return "Позитивное развитие"
	}
	return "Вызов роста"
}

func getRussianTransitHouseTitle(key string) string {
	titles := map[string]string{
		"sun_house_1":  "Пробуждение личности",
		"sun_house_2":  "Фокус на благосостоянии",
		"sun_house_3":  "Самовыражение и общение",
		"sun_house_4":  "Семейное тепло",
		"sun_house_5":  "Творчество и радость",
		"sun_house_6":  "Фокус на здоровье",
		"sun_house_7":  "Партнёрство",
		"sun_house_8":  "Глубина и преображение",
		"sun_house_9":  "Широкие горизонты",
		"sun_house_10": "Карьерный прожектор",
		"sun_house_11": "Социальная активность",
		"sun_house_12": "Внутреннее исследование",
		"moon_house_1":  "Слушай внутренний голос",
		"moon_house_2":  "Эмоциональная безопасность",
		"moon_house_3":  "Эмоциональное выражение",
		"moon_house_4":  "Чувство принадлежности",
		"moon_house_5":  "Эмоциональная радость",
		"moon_house_6":  "Забота о теле и душе",
		"moon_house_7":  "Эмоциональные потребности",
		"moon_house_8":  "Эмоциональная глубина",
		"moon_house_9":  "Эмоциональное исследование",
		"moon_house_10": "Публичный образ",
		"moon_house_11": "Тепло дружбы",
		"moon_house_12": "Подсознание на поверхности",
		"venus_house_1":  "Усиление обаяния",
		"venus_house_2":  "Наслаждение материальным",
		"venus_house_3":  "Приятный обмен",
		"venus_house_4":  "Красота дома",
		"venus_house_5":  "Романтика",
		"venus_house_6":  "Гармония в работе",
		"venus_house_7":  "Сладкие отношения",
		"venus_house_8":  "Интимная связь",
		"venus_house_9":  "Культурное восприятие",
		"venus_house_10": "Обаяние на работе",
		"venus_house_11": "Расцвет социальной харизмы",
		"venus_house_12": "Духовная любовь",
		"mars_house_1":  "Всплеск энергии",
		"mars_house_2":  "Драйв к заработку",
		"mars_house_3":  "Напористые слова",
		"mars_house_4":  "Действия дома",
		"mars_house_5":  "Соревновательный дух",
		"mars_house_6":  "Рабочий напор",
		"mars_house_7":  "Напряжение в отношениях",
		"mars_house_8":  "Сильное желание",
		"mars_house_9":  "Приключения и путешествия",
		"mars_house_10": "Карьерные амбиции",
		"mars_house_11": "Смелость и отвага",
		"mars_house_12": "Скрытое действие",
		"jupiter_house_1":  "Расширение уверенности",
		"jupiter_house_2":  "Богатство и возможности",
		"jupiter_house_3":  "Возможность учиться",
		"jupiter_house_4":  "Семейная радость",
		"jupiter_house_5":  "Творческое изобилие",
		"jupiter_house_6":  "Улучшение здоровья",
		"jupiter_house_7":  "Рост отношений",
		"jupiter_house_8":  "Благословение трансформации",
		"jupiter_house_9":  "Просветление мудрости",
		"jupiter_house_10": "Карьерная возможность",
		"jupiter_house_11": "Появление покровителей",
		"jupiter_house_12": "Духовный рост",
		"saturn_house_1":  "Дисциплина и структура",
		"saturn_house_2":  "Финансовая ответственность",
		"saturn_house_3":  "Серьёзное общение",
		"saturn_house_4":  "Семейная ответственность",
		"saturn_house_5":  "Творческая дисциплина",
		"saturn_house_6":  "Осторожность со здоровьем",
		"saturn_house_7":  "Испытание отношений",
		"saturn_house_8":  "Глубина и порядок",
		"saturn_house_9":  "Философская мысль",
		"saturn_house_10": "Карьерная структура",
		"saturn_house_11": "Социальная фильтрация",
		"saturn_house_12": "Очищение кармы",

		"mercury_house_1":  "Выражение и мысль",
		"mercury_house_2":  "Ценности и мысли о деньгах",
		"mercury_house_3":  "Общение и обучение",
		"mercury_house_4":  "Семейное общение",
		"mercury_house_5":  "Творческое выражение",
		"mercury_house_6":  "Общение о работе и здоровье",
		"mercury_house_7":  "Общение с партнёром",
		"mercury_house_8":  "Глубокий обмен",
		"mercury_house_9":  "Мудрость и видение",
		"mercury_house_10": "Карьерное выражение",
		"mercury_house_11": "Социальное и идеальное общение",
		"mercury_house_12": "Внутреннее и вдохновение",

		"uranus_house_1":  "Революция себя",
		"uranus_house_2":  "Финансовые перемены",
		"uranus_house_3":  "Инновации в общении",
		"uranus_house_4":  "Семейная революция",
		"uranus_house_5":  "Творческий прорыв",
		"uranus_house_6":  "Революция в работе и здоровье",
		"uranus_house_7":  "Внезапный сдвиг в отношениях",
		"uranus_house_8":  "Глубина и перемены в ресурсах",
		"uranus_house_9":  "Прорыв в убеждениях и видении",
		"uranus_house_10": "Карьерная революция",
		"uranus_house_11": "Перемены в сообществе и идеалах",
		"uranus_house_12": "Внутреннее пробуждение",

		"neptune_house_1":  "Духовное я",
		"neptune_house_2":  "Ценности и интуиция",
		"neptune_house_3":  "Вдохновлённое общение",
		"neptune_house_4":  "Дом и принадлежность",
		"neptune_house_5":  "Творчество и вдохновение",
		"neptune_house_6":  "Тело, ум и дух",
		"neptune_house_7":  "Отношения и эмпатия",
		"neptune_house_8":  "Глубокая эмпатия и трансформация",
		"neptune_house_9":  "Вера и духовный поиск",
		"neptune_house_10": "Карьера и идеализация",
		"neptune_house_11": "Сообщество и видение",
		"neptune_house_12": "Духовное оседание",

		"pluto_house_1":  "Трансформация себя",
		"pluto_house_2":  "Глубина богатства и ценностей",
		"pluto_house_3":  "Глубина общения",
		"pluto_house_4":  "Трансформация семьи и корней",
		"pluto_house_5":  "Глубина творчества и желания",
		"pluto_house_6":  "Трансформация работы и здоровья",
		"pluto_house_7":  "Глубокая трансформация отношений",
		"pluto_house_8":  "Глубина и возрождение",
		"pluto_house_9":  "Трансформация убеждений и философии",
		"pluto_house_10": "Трансформация карьеры и власти",
		"pluto_house_11": "Трансформация сообщества и коллектива",
		"pluto_house_12": "Карма и бессознательное",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	return "Транзит по дому"
}

// Retrograde emotional titles (by planet)
func getChineseRetrogradeTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"mercury": "水星逆行·复盘与修正",
		"venus":   "金星逆行·价值与关系回顾",
		"mars":    "火星逆行·策略与收束",
		"jupiter": "木星逆行·夯实与内化",
		"saturn":  "土星逆行·责任与成熟",
		"uranus":  "天王星逆行·内在变革",
		"neptune": "海王星逆行·边界与落地",
		"pluto":   "冥王星逆行·转化与放下",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "逆行期·回顾与调整"
}

func getEnglishRetrogradeTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"mercury": "Mercury Retrograde · Review & Revise",
		"venus":   "Venus Retrograde · Values & Relationships",
		"mars":    "Mars Retrograde · Strategy & Restraint",
		"jupiter": "Jupiter Retrograde · Consolidate & Integrate",
		"saturn":  "Saturn Retrograde · Duty & Maturity",
		"uranus":  "Uranus Retrograde · Inner Change",
		"neptune": "Neptune Retrograde · Boundaries & Grounding",
		"pluto":   "Pluto Retrograde · Transform & Release",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "Retrograde · Review & Adjust"
}

func getRussianRetrogradeTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"mercury": "Ретроградный Меркурий · Обзор и правки",
		"venus":   "Ретроградная Венера · Ценности и отношения",
		"mars":    "Ретроградный Марс · Стратегия и сдержанность",
		"jupiter": "Ретроградный Юпитер · Укрепление и интеграция",
		"saturn":  "Ретроградный Сатурн · Долг и зрелость",
		"uranus":  "Ретроградный Уран · Внутренние перемены",
		"neptune": "Ретроградный Нептун · Границы и заземление",
		"pluto":   "Ретроградный Плутон · Трансформация и отпускание",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "Ретроградность · Обзор и корректировка"
}

// Profection Lord (Annual Lord) emotional titles (by planet)
func getChineseProfectionLordTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"sun":     "年主星太阳·自我与活力",
		"moon":    "年主星月亮·情绪与家庭",
		"mercury": "年主星水星·沟通与学习",
		"venus":   "年主星金星·爱与价值",
		"mars":    "年主星火星·行动与斗志",
		"jupiter": "年主星木星·扩展与机遇",
		"saturn":  "年主星土星·责任与结构",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "年主星·年度主题"
}

func getEnglishProfectionLordTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"sun":     "Annual Lord Sun · Self & Vitality",
		"moon":    "Annual Lord Moon · Emotions & Home",
		"mercury": "Annual Lord Mercury · Communication & Learning",
		"venus":   "Annual Lord Venus · Love & Value",
		"mars":    "Annual Lord Mars · Action & Drive",
		"jupiter": "Annual Lord Jupiter · Expansion & Opportunity",
		"saturn":  "Annual Lord Saturn · Responsibility & Structure",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "Annual Lord · Yearly Theme"
}

func getRussianProfectionLordTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"sun":     "Годовой управитель Солнце · Самость и сила",
		"moon":    "Годовой управитель Луна · Эмоции и дом",
		"mercury": "Годовой управитель Меркурий · Общение и учёба",
		"venus":   "Годовой управитель Венера · Любовь и ценность",
		"mars":    "Годовой управитель Марс · Действие и драйв",
		"jupiter": "Годовой управитель Юпитер · Расширение и возможности",
		"saturn":  "Годовой управитель Сатурн · Ответственность и структура",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "Годовой управитель · Тема года"
}

// Planetary Hour emotional titles (hourly)
func getChinesePlanetaryHourTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"sun":     "太阳时·自我与决断",
		"moon":    "月亮时·情绪与直觉",
		"mercury": "水星时·沟通与文书",
		"venus":   "金星时·关系与金钱",
		"mars":    "火星时·行动与竞争",
		"jupiter": "木星时·扩展与机遇",
		"saturn":  "土星时·责任与收尾",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "行星时·当下能量"
}

func getEnglishPlanetaryHourTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"sun":     "Sun Hour · Self & Decision",
		"moon":    "Moon Hour · Emotion & Intuition",
		"mercury": "Mercury Hour · Communication & Paperwork",
		"venus":   "Venus Hour · Relationship & Money",
		"mars":    "Mars Hour · Action & Competition",
		"jupiter": "Jupiter Hour · Expansion & Opportunity",
		"saturn":  "Saturn Hour · Duty & Wrap-up",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "Planetary Hour · Current Energy"
}

func getRussianPlanetaryHourTitle(planet models.PlanetID) string {
	titles := map[string]string{
		"sun":     "Час Солнца · Самость и решение",
		"moon":    "Час Луны · Эмоция и интуиция",
		"mercury": "Час Меркурия · Общение и документы",
		"venus":   "Час Венеры · Отношения и деньги",
		"mars":    "Час Марса · Действие и соревнование",
		"jupiter": "Час Юпитера · Расширение и возможности",
		"saturn":  "Час Сатурна · Долг и завершение",
	}
	if s, ok := titles[string(planet)]; ok {
		return s
	}
	return "Планетарный час · Текущая энергия"
}

// Void of Course (Moon) emotional titles (hourly)
func getChineseVoidOfCourseTitle() string {
	return "月亮空亡·宜收尾不宜新启"
}

func getEnglishVoidOfCourseTitle() string {
	return "Moon Void of Course · Wrap Up, Don't Start New"
}

func getRussianVoidOfCourseTitle() string {
	return "Луна без курса · Завершайте, не начинайте нового"
}

// Lunar phase emotional titles
func getChineseLunarPhaseTitle(phase string) string {
	if phase == "first_quarter" {
		phase = "firstQuarter"
	} else if phase == "last_quarter" {
		phase = "lastQuarter"
	}
	titles := map[string]string{
		"new": "新月·设定意图", "crescent": "月牙·突破行动", "firstQuarter": "上弦月·考验与选择",
		"gibbous": "渐盈月·打磨准备", "full": "满月·收获与释放", "disseminating": "渐亏月·分享智慧",
		"lastQuarter": "下弦月·放下调整", "balsamic": "残月·整合与休息",
	}
	if s, ok := titles[phase]; ok {
		return s
	}
	return "月相·能量周期"
}

func getEnglishLunarPhaseTitle(phase string) string {
	if phase == "first_quarter" {
		phase = "firstQuarter"
	} else if phase == "last_quarter" {
		phase = "lastQuarter"
	}
	titles := map[string]string{
		"new": "New Moon · Set Intentions", "crescent": "Crescent · Break Through", "firstQuarter": "First Quarter · Test & Choose",
		"gibbous": "Gibbous · Refine & Prepare", "full": "Full Moon · Harvest & Release", "disseminating": "Disseminating · Share Wisdom",
		"lastQuarter": "Last Quarter · Let Go & Adjust", "balsamic": "Balsamic · Integrate & Rest",
	}
	if s, ok := titles[phase]; ok {
		return s
	}
	return "Lunar Phase · Energy Cycle"
}

func getRussianLunarPhaseTitle(phase string) string {
	if phase == "first_quarter" {
		phase = "firstQuarter"
	} else if phase == "last_quarter" {
		phase = "lastQuarter"
	}
	titles := map[string]string{
		"new": "Новолуние · Намерения", "crescent": "Молодая луна · Прорыв", "firstQuarter": "Первая четверть · Испытание",
		"gibbous": "Растущая · Отделка", "full": "Полнолуние · Урожай и освобождение", "disseminating": "Убывающая · Делиться мудростью",
		"lastQuarter": "Последняя четверть · Отпускание", "balsamic": "Бальзамическая · Интеграция",
	}
	if s, ok := titles[phase]; ok {
		return s
	}
	return "Лунная фаза · Энергетический цикл"
}

// Sign change emotional titles (short; full text in detailed interpretation)
func getChineseSignChangeTitle(planet models.PlanetID, newSign string) string {
	p := string(planet)
	names := map[string]string{"sun": "太阳", "moon": "月亮", "mercury": "水星", "venus": "金星", "mars": "火星", "jupiter": "木星", "saturn": "土星"}
	signs := map[string]string{"aries": "白羊", "taurus": "金牛", "gemini": "双子", "cancer": "巨蟹", "leo": "狮子", "virgo": "处女", "libra": "天秤", "scorpio": "天蝎", "sagittarius": "射手", "capricorn": "摩羯", "aquarius": "水瓶", "pisces": "双鱼"}
	n, sn := names[p], signs[strings.ToLower(newSign)]
	if n == "" {
		n = p
	}
	if sn == "" {
		sn = newSign
	}
	return n + "进入" + sn + "座"
}

func getEnglishSignChangeTitle(planet models.PlanetID, newSign string) string {
	names := map[string]string{"sun": "Sun", "moon": "Moon", "mercury": "Mercury", "venus": "Venus", "mars": "Mars", "jupiter": "Jupiter", "saturn": "Saturn"}
	n := names[string(planet)]
	if n == "" {
		n = string(planet)
	}
	return n + " Enters " + newSign
}

func getRussianSignChangeTitle(planet models.PlanetID, newSign string) string {
	names := map[string]string{"sun": "Солнце", "moon": "Луна", "mercury": "Меркурий", "venus": "Венера", "mars": "Марс", "jupiter": "Юпитер", "saturn": "Сатурн"}
	n := names[string(planet)]
	if n == "" {
		n = string(planet)
	}
	return n + " входит в " + newSign
}

// Dignity emotional titles
func getChineseDignityTitle(planet models.PlanetID, dignityType string) string {
	names := map[string]string{"sun": "太阳", "moon": "月亮", "mercury": "水星", "venus": "金星", "mars": "火星", "jupiter": "木星", "saturn": "土星"}
	d := map[string]string{"domicile": "入庙", "exaltation": "旺相", "detriment": "落陷", "fall": "失势"}
	n, dn := names[string(planet)], d[dignityType]
	if n == "" {
		n = string(planet)
	}
	if dn == "" {
		dn = dignityType
	}
	return n + "·" + dn
}

func getEnglishDignityTitle(planet models.PlanetID, dignityType string) string {
	names := map[string]string{"sun": "Sun", "moon": "Moon", "mercury": "Mercury", "venus": "Venus", "mars": "Mars", "jupiter": "Jupiter", "saturn": "Saturn"}
	n := names[string(planet)]
	if n == "" {
		n = string(planet)
	}
	return n + " in " + dignityType
}

func getRussianDignityTitle(planet models.PlanetID, dignityType string) string {
	names := map[string]string{"sun": "Солнце", "moon": "Луна", "mercury": "Меркурий", "venus": "Венера", "mars": "Марс", "jupiter": "Юпитер", "saturn": "Сатурн"}
	d := map[string]string{"domicile": "в доме", "exaltation": "в экзальтации", "detriment": "в изгнании", "fall": "в падении"}
	n, dn := names[string(planet)], d[dignityType]
	if n == "" {
		n = string(planet)
	}
	if dn == "" {
		dn = dignityType
	}
	return n + " " + dn
}
