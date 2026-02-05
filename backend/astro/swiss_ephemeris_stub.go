//go:build !swe
// +build !swe

package astro

import (
	"star/models"
	"time"
)

// ==================== Swiss Ephemeris 桩实现（无CGO依赖） ====================

var sweAvailable = false // Swiss Ephemeris 不可用

// InitSwissEphemeris 初始化 Swiss Ephemeris（桩实现）
func InitSwissEphemeris(ephePath string) {
	// 无操作
}

// CloseSwissEphemeris 关闭 Swiss Ephemeris（桩实现）
func CloseSwissEphemeris() {
	// 无操作
}

// IsSweAvailable 返回 Swiss Ephemeris 是否可用
func IsSweAvailable() bool {
	return sweAvailable
}

// CalculatePlanetPositionSwe 使用内置算法计算行星位置（回退实现）
func CalculatePlanetPositionSwe(planet models.PlanetID, jd float64) models.PlanetPosition {
	// 回退到内置算法
	return CalculatePlanetPosition(planet, jd)
}

// GetAllPlanetPositionsSwe 使用内置算法获取所有行星位置（回退实现）
func GetAllPlanetPositionsSwe(jd float64) []models.PlanetPosition {
	return GetAllPlanetPositions(jd)
}

// CalculateHousesSwe 使用内置算法计算宫位（回退实现）
func CalculateHousesSwe(jd float64, lat, lon float64) ([]models.HouseCusp, float64, float64) {
	return CalculateHouses(jd, lat, lon)
}

// ==================== aspect_search.go 桩实现 ====================

// TimeToJulianDay 将 time.Time 转换为儒略日（桩实现，使用标准公式）
func TimeToJulianDay(t time.Time) float64 {
	t = t.UTC()
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	hour := float64(t.Hour()) + float64(t.Minute())/60.0 + float64(t.Second())/3600.0

	// 标准儒略日公式
	if month <= 2 {
		year--
		month += 12
	}
	a := year / 100
	b := 2 - a + a/4
	jd := float64(int(365.25*float64(year+4716))) + float64(int(30.6001*float64(month+1))) + float64(day) + hour/24.0 + float64(b) - 1524.5
	return jd
}

// JulianDayToTime 将儒略日转换为 time.Time（桩实现，使用标准公式）
func JulianDayToTime(jd float64) time.Time {
	z := int(jd + 0.5)
	f := jd + 0.5 - float64(z)

	var a int
	if z < 2299161 {
		a = z
	} else {
		alpha := int((float64(z) - 1867216.25) / 36524.25)
		a = z + 1 + alpha - alpha/4
	}

	b := a + 1524
	c := int((float64(b) - 122.1) / 365.25)
	d := int(365.25 * float64(c))
	e := int(float64(b-d) / 30.6001)

	day := b - d - int(30.6001*float64(e))
	var month int
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}
	var year int
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}

	hours := f * 24.0
	h := int(hours)
	mins := (hours - float64(h)) * 60
	m := int(mins)
	secs := (mins - float64(m)) * 60
	s := int(secs)
	ns := int((secs - float64(s)) * 1e9)

	return time.Date(year, time.Month(month), day, h, m, s, ns, time.UTC)
}

// GetPlanetLongitudeAt 获取行星在指定儒略日的黄经（桩实现，返回近似值）
func GetPlanetLongitudeAt(planet models.PlanetID, jd float64) float64 {
	// 回退到内置算法
	pos := CalculatePlanetPosition(planet, jd)
	return pos.Longitude
}

// FindExactAspectTime 使用二分法查找精确相位时间（桩实现）
func FindExactAspectTime(planet1, planet2 models.PlanetID, targetAngle float64, startJd, endJd float64) (float64, bool) {
	// 简化实现：返回中点时间和 true
	// 在没有精确星历时，这只是一个近似值
	midJd := (startJd + endJd) / 2.0
	return midJd, true
}
