package i18n

// getTranslations returns translations for a given language
func getTranslations(lang Language) map[string]string {
	switch lang {
	case Chinese:
		return chineseTranslations
	case Russian:
		return russianTranslations
	default:
		return englishTranslations
	}
}

var englishTranslations = map[string]string{
	// Planets
	"planet.sun":       "Sun",
	"planet.moon":      "Moon",
	"planet.mercury":   "Mercury",
	"planet.venus":     "Venus",
	"planet.mars":      "Mars",
	"planet.jupiter":   "Jupiter",
	"planet.saturn":    "Saturn",
	"planet.uranus":    "Uranus",
	"planet.neptune":   "Neptune",
	"planet.pluto":     "Pluto",
	"planet.northNode": "North Node",
	"planet.southNode": "South Node",
	"planet.chiron":    "Chiron",
	// Planet aliases
	"planet.north": "North Node",
	"planet.south": "South Node",

	// Event title suffixes
	"planetary_hour_suffix": " Hour",
	"retrograde_suffix":     " Retrograde",
	"void_of_course":        "Moon Void of Course",
	"sign_change_suffix":    " Sign Change",

	// Major Aspects
	"aspect.conjunction": "conjunction",
	"aspect.sextile":     "sextile",
	"aspect.square":      "square",
	"aspect.trine":       "trine",
	"aspect.opposition":  "opposition",
	// Minor Aspects
	"aspect.semi-sextile":   "semi-sextile",
	"aspect.semi-square":    "semi-square",
	"aspect.sesquiquadrate": "sesquiquadrate",
	"aspect.quincunx":       "quincunx",
	// Node aspects
	"aspect.node": "Node Axis",

	// Lunar Phases
	"lunar_phase.new":           "New Moon",
	"lunar_phase.first_quarter": "First Quarter",
	"lunar_phase.full":          "Full Moon",
	"lunar_phase.last_quarter":  "Last Quarter",
	"lunar_phase.waxing":        "Waxing",
	"lunar_phase.waning":        "Waning",

	// Dignities
	"dignity.domicile":   "Domicile",
	"dignity.exaltation": "Exaltation",
	"dignity.detriment":  "Detriment",
	"dignity.fall":       "Fall",

	// Aspect Phases
	"phase.applying":   "Applying",
	"phase.exact":      "Exact",
	"phase.separating": "Separating",

	// Transit States
	"transit.entering": "Entering",
	"transit.active":   "Active",
	"transit.leaving":  "Leaving",

	// Signs
	"sign.aries":       "Aries",
	"sign.taurus":      "Taurus",
	"sign.gemini":      "Gemini",
	"sign.cancer":      "Cancer",
	"sign.leo":         "Leo",
	"sign.virgo":       "Virgo",
	"sign.libra":       "Libra",
	"sign.scorpio":     "Scorpio",
	"sign.sagittarius": "Sagittarius",
	"sign.capricorn":   "Capricorn",
	"sign.aquarius":    "Aquarius",
	"sign.pisces":      "Pisces",

	// Dimensions
	"dimension.career":       "Career",
	"dimension.relationship": "Relationship",
	"dimension.health":       "Health",
	"dimension.finance":      "Finance",
	"dimension.spiritual":    "Spiritual",

	// Houses
	"house.1":  "1st House - Self",
	"house.2":  "2nd House - Resources",
	"house.3":  "3rd House - Communication",
	"house.4":  "4th House - Home",
	"house.5":  "5th House - Creativity",
	"house.6":  "6th House - Health",
	"house.7":  "7th House - Partnership",
	"house.8":  "8th House - Transformation",
	"house.9":  "9th House - Expansion",
	"house.10": "10th House - Career",
	"house.11": "11th House - Community",
	"house.12": "12th House - Spirituality",

	// Event types
	"event_type.aspect":                "Aspect",
	"event_type.sign_change":           "Sign Change",
	"event_type.lunar_phase":           "Lunar Phase",
	"event_type.planetary_hour_change": "Planetary Hour",
	"event_type.transit_house":         "Transit House",
	"event_type.secondary_progression": "Secondary Progression",
	"event_type.tertiary_progression":  "Tertiary Progression",
	"event_type.retrograde":            "Retrograde",
	"event_type.dignity":               "Dignity",
	"event_type.voidOfCourse":          "Void of Course",

	// Intensity
	"intensity.high":   "High",
	"intensity.medium": "Medium",
	"intensity.low":    "Low",

	// Time levels
	"time_level.yearly":  "Yearly",
	"time_level.monthly": "Monthly",
	"time_level.weekly":  "Weekly",
	"time_level.daily":   "Daily",
	"time_level.hourly":  "Hourly",
}

var chineseTranslations = map[string]string{
	// Planets
	"planet.sun":       "太阳",
	"planet.moon":      "月亮",
	"planet.mercury":   "水星",
	"planet.venus":     "金星",
	"planet.mars":      "火星",
	"planet.jupiter":   "木星",
	"planet.saturn":    "土星",
	"planet.uranus":    "天王星",
	"planet.neptune":   "海王星",
	"planet.pluto":     "冥王星",
	"planet.northNode": "北交点",
	"planet.southNode": "南交点",
	"planet.chiron":    "凯龙星",
	// Planet aliases
	"planet.north": "北交点",
	"planet.south": "南交点",

	// Event title suffixes
	"planetary_hour_suffix": "时",
	"retrograde_suffix":     "逆行",
	"void_of_course":        "月亮空亡",
	"sign_change_suffix":    "换座",

	// Major Aspects
	"aspect.conjunction": "合相",
	"aspect.sextile":     "六分相",
	"aspect.square":      "刑相",
	"aspect.trine":       "拱相",
	"aspect.opposition":  "对分相",
	// Minor Aspects
	"aspect.semi-sextile":   "半六分相",
	"aspect.semi-square":    "半刑相",
	"aspect.sesquiquadrate": "倍半刑相",
	"aspect.quincunx":       "梅花相",
	// Node aspects
	"aspect.node": "交点轴",

	// Lunar Phases
	"lunar_phase.new":           "新月",
	"lunar_phase.first_quarter": "上弦月",
	"lunar_phase.full":          "满月",
	"lunar_phase.last_quarter":  "下弦月",
	"lunar_phase.waxing":        "盈月",
	"lunar_phase.waning":        "亏月",

	// Dignities
	"dignity.domicile":   "入庙",
	"dignity.exaltation": "旺相",
	"dignity.detriment":  "落陷",
	"dignity.fall":       "失势",

	// Aspect Phases
	"phase.applying":   "入相",
	"phase.exact":      "精确",
	"phase.separating": "离相",

	// Transit States
	"transit.entering": "进入",
	"transit.active":   "活跃",
	"transit.leaving":  "离开",

	// Signs
	"sign.aries":       "白羊座",
	"sign.taurus":      "金牛座",
	"sign.gemini":      "双子座",
	"sign.cancer":      "巨蟹座",
	"sign.leo":         "狮子座",
	"sign.virgo":       "处女座",
	"sign.libra":       "天秤座",
	"sign.scorpio":     "天蝎座",
	"sign.sagittarius": "射手座",
	"sign.capricorn":   "摩羯座",
	"sign.aquarius":    "水瓶座",
	"sign.pisces":      "双鱼座",

	// Dimensions
	"dimension.career":       "事业",
	"dimension.relationship": "爱情",
	"dimension.health":       "健康",
	"dimension.finance":      "财运",
	"dimension.spiritual":    "灵性",

	// Houses
	"house.1":  "第1宫 - 自我",
	"house.2":  "第2宫 - 财富",
	"house.3":  "第3宫 - 沟通",
	"house.4":  "第4宫 - 家庭",
	"house.5":  "第5宫 - 创造",
	"house.6":  "第6宫 - 健康",
	"house.7":  "第7宫 - 伴侣",
	"house.8":  "第8宫 - 转化",
	"house.9":  "第9宫 - 扩展",
	"house.10": "第10宫 - 事业",
	"house.11": "第11宫 - 社群",
	"house.12": "第12宫 - 灵性",

	// Event types
	"event_type.aspect":                "相位",
	"event_type.sign_change":           "换座",
	"event_type.lunar_phase":           "月相",
	"event_type.planetary_hour_change": "行星时",
	"event_type.transit_house":         "行运过宫",
	"event_type.secondary_progression": "次限推进",
	"event_type.tertiary_progression":  "三限推进",
	"event_type.retrograde":            "逆行",
	"event_type.dignity":               "尊贵度",
	"event_type.voidOfCourse":          "月亮空亡",

	// Intensity
	"intensity.high":   "高",
	"intensity.medium": "中",
	"intensity.low":    "低",

	// Time levels
	"time_level.yearly":  "年度",
	"time_level.monthly": "月度",
	"time_level.weekly":  "周度",
	"time_level.daily":   "日度",
	"time_level.hourly":  "小时",
}

var russianTranslations = map[string]string{
	// Planets
	"planet.sun":       "Солнце",
	"planet.moon":      "Луна",
	"planet.mercury":   "Меркурий",
	"planet.venus":     "Венера",
	"planet.mars":      "Марс",
	"planet.jupiter":   "Юпитер",
	"planet.saturn":    "Сатурн",
	"planet.uranus":    "Уран",
	"planet.neptune":   "Нептун",
	"planet.pluto":     "Плутон",
	"planet.northNode": "Северный Узел",
	"planet.southNode": "Южный Узел",
	"planet.chiron":    "Хирон",
	// Planet aliases
	"planet.north": "Северный Узел",
	"planet.south": "Южный Узел",

	// Event title suffixes
	"planetary_hour_suffix": " Час",
	"retrograde_suffix":     " Ретроград",
	"void_of_course":        "Луна без курса",
	"sign_change_suffix":    " Смена знака",

	// Major Aspects
	"aspect.conjunction": "соединение",
	"aspect.sextile":     "секстиль",
	"aspect.square":      "квадрат",
	"aspect.trine":       "трин",
	"aspect.opposition":  "оппозиция",
	// Minor Aspects
	"aspect.semi-sextile":   "полусекстиль",
	"aspect.semi-square":    "полуквадрат",
	"aspect.sesquiquadrate": "полутораквадрат",
	"aspect.quincunx":       "квиконс",
	// Node aspects
	"aspect.node": "Ось Узлов",

	// Lunar Phases
	"lunar_phase.new":           "Новолуние",
	"lunar_phase.first_quarter": "Первая четверть",
	"lunar_phase.full":          "Полнолуние",
	"lunar_phase.last_quarter":  "Последняя четверть",
	"lunar_phase.waxing":        "Растущая",
	"lunar_phase.waning":        "Убывающая",

	// Dignities
	"dignity.domicile":   "Обитель",
	"dignity.exaltation": "Экзальтация",
	"dignity.detriment":  "Изгнание",
	"dignity.fall":       "Падение",

	// Aspect Phases
	"phase.applying":   "Сходящийся",
	"phase.exact":      "Точный",
	"phase.separating": "Расходящийся",

	// Transit States
	"transit.entering": "Входит",
	"transit.active":   "Активный",
	"transit.leaving":  "Выходит",

	// Signs
	"sign.aries":       "Овен",
	"sign.taurus":      "Телец",
	"sign.gemini":      "Близнецы",
	"sign.cancer":      "Рак",
	"sign.leo":         "Лев",
	"sign.virgo":       "Дева",
	"sign.libra":       "Весы",
	"sign.scorpio":     "Скорпион",
	"sign.sagittarius": "Стрелец",
	"sign.capricorn":   "Козерог",
	"sign.aquarius":    "Водолей",
	"sign.pisces":      "Рыбы",

	// Dimensions
	"dimension.career":       "Карьера",
	"dimension.relationship": "Отношения",
	"dimension.health":       "Здоровье",
	"dimension.finance":      "Финансы",
	"dimension.spiritual":    "Духовность",

	// Houses
	"house.1":  "1-й дом - Личность",
	"house.2":  "2-й дом - Ресурсы",
	"house.3":  "3-й дом - Общение",
	"house.4":  "4-й дом - Дом",
	"house.5":  "5-й дом - Творчество",
	"house.6":  "6-й дом - Здоровье",
	"house.7":  "7-й дом - Партнерство",
	"house.8":  "8-й дом - Трансформация",
	"house.9":  "9-й дом - Расширение",
	"house.10": "10-й дом - Карьера",
	"house.11": "11-й дом - Сообщество",
	"house.12": "12-й дом - Духовность",

	// Event types
	"event_type.aspect":                "Аспект",
	"event_type.sign_change":           "Смена знака",
	"event_type.lunar_phase":           "Лунная фаза",
	"event_type.planetary_hour_change": "Планетный час",
	"event_type.transit_house":         "Транзит по дому",
	"event_type.secondary_progression": "Вторичная прогрессия",
	"event_type.tertiary_progression":  "Третичная прогрессия",
	"event_type.retrograde":            "Ретроградность",
	"event_type.dignity":               "Достоинство",
	"event_type.voidOfCourse":          "Луна без курса",

	// Intensity
	"intensity.high":   "Высокий",
	"intensity.medium": "Средний",
	"intensity.low":    "Низкий",

	// Time levels
	"time_level.yearly":  "Годовой",
	"time_level.monthly": "Месячный",
	"time_level.weekly":  "Недельный",
	"time_level.daily":   "Дневной",
	"time_level.hourly":  "Часовой",
}
