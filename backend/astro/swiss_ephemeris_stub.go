//go:build !swe
// +build !swe

package astro

import (
	"star/models"
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

