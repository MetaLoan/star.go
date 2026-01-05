package astro

import (
	"fmt"
	"star/models"
)

// ==================== 统一天文计算入口 ====================
// 所有天体数据的唯一来源是 Swiss Ephemeris
// 这些函数是整个系统的统一入口点

// GetPlanetPositionsUnified 统一获取所有行星位置
// 强制使用 Swiss Ephemeris 作为唯一数据源
func GetPlanetPositionsUnified(jd float64) []models.PlanetPosition {
	if !IsSweAvailable() {
		panic("Swiss Ephemeris is required but not available. Build with -tags swe")
	}
	return GetAllPlanetPositionsSwe(jd)
}

// CalculatePlanetPositionUnified 统一计算单个行星位置
// 强制使用 Swiss Ephemeris 作为唯一数据源
func CalculatePlanetPositionUnified(planet models.PlanetID, jd float64) models.PlanetPosition {
	if !IsSweAvailable() {
		panic("Swiss Ephemeris is required but not available. Build with -tags swe")
	}
	return CalculatePlanetPositionSwe(planet, jd)
}

// CalculateHousesUnified 统一计算宫位
// 强制使用 Swiss Ephemeris 作为唯一数据源
func CalculateHousesUnified(jd float64, lat, lon float64) ([]models.HouseCusp, float64, float64) {
	if !IsSweAvailable() {
		panic("Swiss Ephemeris is required but not available. Build with -tags swe")
	}
	return CalculateHousesSwe(jd, lat, lon)
}

// ValidateSwissEphemeris 验证 Swiss Ephemeris 是否可用
// 在应用启动时调用此函数确保数据源正确
func ValidateSwissEphemeris() error {
	if !IsSweAvailable() {
		return fmt.Errorf("Swiss Ephemeris is not available. The application requires Swiss Ephemeris as the sole data source. Please rebuild with: CGO_ENABLED=1 go build -tags swe")
	}
	
	// 初始化 Swiss Ephemeris
	InitSwissEphemeris("")
	
	return nil
}

// GetDataSource 返回当前数据源信息
func GetDataSource() string {
	if IsSweAvailable() {
		return "Swiss Ephemeris (High Precision)"
	}
	return "Built-in Algorithm (NOT RECOMMENDED)"
}

