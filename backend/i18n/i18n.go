package i18n

import "star/models"

// Language supported languages
type Language string

const (
	Chinese Language = "zh"
	English Language = "en"
	Russian Language = "ru"
)

// Translator provides translation functions
type Translator struct {
	lang Language
}

// New creates a new translator
func New(lang string) *Translator {
	l := Language(lang)
	if l != Chinese && l != English && l != Russian {
		l = English // Default to English
	}
	return &Translator{lang: l}
}

// T translates a key
func (t *Translator) T(key string, params ...interface{}) string {
	translations := getTranslations(t.lang)
	if val, ok := translations[key]; ok {
		// Simple parameter substitution if needed
		return val
	}
	return key
}

// GetPlanetName returns translated planet name
func (t *Translator) GetPlanetName(planet models.PlanetID) string {
	key := "planet." + string(planet)
	return t.T(key)
}

// GetAspectName returns translated aspect name
func (t *Translator) GetAspectName(aspect string) string {
	key := "aspect." + aspect
	return t.T(key)
}

// GetSignName returns translated sign name
func (t *Translator) GetSignName(sign string) string {
	key := "sign." + sign
	return t.T(key)
}

// GetHouseName returns translated house name
func (t *Translator) GetHouseName(house int) string {
	key := "house." + string(rune('0'+house))
	return t.T(key)
}

// GetDimensionName returns translated dimension name
func (t *Translator) GetDimensionName(dimension string) string {
	key := "dimension." + dimension
	return t.T(key)
}

// GetEventTypeName returns translated event type name
func (t *Translator) GetEventTypeName(eventType string) string {
	key := "event_type." + eventType
	return t.T(key)
}

// GetIntensityLabel returns translated intensity label
func (t *Translator) GetIntensityLabel(intensity string) string {
	key := "intensity." + intensity
	return t.T(key)
}

// GetTimeLevelName returns translated time level name
func (t *Translator) GetTimeLevelName(level string) string {
	key := "time_level." + level
	return t.T(key)
}

// FormatDuration formats duration text like "Started 4 months ago, ends in 2 months"
func (t *Translator) FormatDuration(startedMonthsAgo, endsInMonths int) string {
	switch t.lang {
	case Chinese:
		if endsInMonths > 0 {
			return formatString("已开始%d个月，%d个月后结束", startedMonthsAgo, endsInMonths)
		} else if endsInMonths == 0 {
			return formatString("已开始%d个月，今天结束", startedMonthsAgo)
		} else {
			return formatString("已开始%d个月，已结束", startedMonthsAgo)
		}
	case Russian:
		if endsInMonths > 0 {
			return formatString("Началось %d мес. назад, закончится через %d мес.", startedMonthsAgo, endsInMonths)
		} else if endsInMonths == 0 {
			return formatString("Началось %d мес. назад, заканчивается сегодня", startedMonthsAgo)
		} else {
			return formatString("Началось %d мес. назад, уже закончилось", startedMonthsAgo)
		}
	default: // English
		if endsInMonths > 0 {
			return formatString("Started %d months ago, ends in %d months", startedMonthsAgo, endsInMonths)
		} else if endsInMonths == 0 {
			return formatString("Started %d months ago, ends today", startedMonthsAgo)
		} else {
			return formatString("Started %d months ago, already ended", startedMonthsAgo)
		}
	}
}

// FormatDurationDays formats duration in days
func (t *Translator) FormatDurationDays(startedDaysAgo, endsInDays int) string {
	switch t.lang {
	case Chinese:
		if endsInDays > 0 {
			return formatString("已开始%d天，%d天后结束", startedDaysAgo, endsInDays)
		} else if endsInDays == 0 {
			return formatString("已开始%d天，今天结束", startedDaysAgo)
		} else {
			return formatString("已开始%d天，已结束", startedDaysAgo)
		}
	case Russian:
		if endsInDays > 0 {
			return formatString("Началось %d дн. назад, закончится через %d дн.", startedDaysAgo, endsInDays)
		} else if endsInDays == 0 {
			return formatString("Началось %d дн. назад, заканчивается сегодня", startedDaysAgo)
		} else {
			return formatString("Началось %d дн. назад, уже закончилось", startedDaysAgo)
		}
	default: // English
		if endsInDays > 0 {
			return formatString("Started %d days ago, ends in %d days", startedDaysAgo, endsInDays)
		} else if endsInDays == 0 {
			return formatString("Started %d days ago, ends today", startedDaysAgo)
		} else {
			return formatString("Started %d days ago, already ended", startedDaysAgo)
		}
	}
}

// GetDimensionLabel returns dimension label with arrow
func (t *Translator) GetDimensionLabel(dimension string, value float64) string {
	arrow := "→"
	if value > 0.6 {
		arrow = "↑"
	} else if value < 0.4 {
		arrow = "↓"
	}
	
	dimName := t.GetDimensionName(dimension)
	return dimName + " " + arrow
}

func formatString(format string, args ...interface{}) string {
	// Simple sprintf-like formatting
	result := format
	for _, arg := range args {
		// Replace first %d with the argument
		// This is a simplified version - in production use fmt.Sprintf
		switch v := arg.(type) {
		case int:
			result = replaceFirst(result, "%d", intToString(v))
		case string:
			result = replaceFirst(result, "%s", v)
		}
	}
	return result
}

func replaceFirst(s, old, new string) string {
	for i := 0; i < len(s)-len(old)+1; i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	
	negative := n < 0
	if negative {
		n = -n
	}
	
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	
	return string(digits)
}
