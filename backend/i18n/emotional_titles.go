package i18n

import "star/models"

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

// Chinese emotional titles
func getChineseAspectTitle(key string, isPositive bool) string {
	titles := map[string]string{
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
		"jupiter_house_11": "贵人现身还请抓住",
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
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	return "过宫影响"
}

// English emotional titles
func getEnglishAspectTitle(key string, isPositive bool) string {
	titles := map[string]string{
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
		
		"moon_conjunction_venus": "Intimate Moments",
		"moon_trine_venus":       "Emotional Fulfillment",
		"moon_square_venus":      "Emotional Fluctuation",
		
		"venus_conjunction_mars": "Burning Desire",
		"venus_trine_mars":       "Passion & Charm",
		"venus_sextile_mars":     "Romantic Spark",
		
		"venus_conjunction_jupiter": "Social Charisma Blooms",
		"venus_trine_jupiter":       "Joy & Happiness",
		
		"mars_trine_jupiter": "Energetic Flow",
		"mars_square_saturn": "Resilient Endurance",
		
		"mercury_trine_jupiter": "Communication & Activity",
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
	
	if isPositive {
		return "Positive Development"
	}
	return "Growth Challenge"
}

func getEnglishTransitHouseTitle(key string) string {
	titles := map[string]string{
		"sun_house_1":  "Self-Awakening",
		"sun_house_10": "Career Spotlight",
		"moon_house_1": "Listen to Inner Voice",
		"venus_house_11": "Social Charisma Blooms",
		"mars_house_11": "Bold & Brave",
		"jupiter_house_11": "Benefactors Appear",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	return "House Transit"
}

// Russian emotional titles
func getRussianAspectTitle(key string, isPositive bool) string {
	titles := map[string]string{
		"sun_conjunction_moon": "Эмоциональная интеграция",
		"sun_trine_moon":       "Внутренняя гармония",
		"sun_conjunction_venus": "Сияющее очарование",
		"sun_trine_mars":       "Неудержимый воин",
		"moon_conjunction_venus": "Интимные моменты",
		"venus_conjunction_mars": "Пылающее желание",
		"venus_conjunction_jupiter": "Расцвет социальной харизмы",
		"mercury_trine_jupiter": "Общение и активность",
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
	
	if isPositive {
		return "Позитивное развитие"
	}
	return "Вызов роста"
}

func getRussianTransitHouseTitle(key string) string {
	titles := map[string]string{
		"sun_house_1":  "Пробуждение личности",
		"sun_house_10": "Карьерный прожектор",
		"moon_house_1": "Слушай внутренний голос",
		"venus_house_11": "Расцвет социальной харизмы",
		"jupiter_house_11": "Появление покровителей",
	}
	
	if title, ok := titles[key]; ok {
		return title
	}
	
	return "Транзит по дому"
}
