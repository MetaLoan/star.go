//go:build swe
// +build swe

package astro

import (
	"math"
	"star/models"
	"time"

	"github.com/mshafiee/swephgo"
)

// ==================== Swiss Ephemeris 初始化 ====================

var sweInitialized bool
var sweAvailable = true // 标记 Swiss Ephemeris 是否可用

// InitSwissEphemeris 初始化 Swiss Ephemeris
// ephePath: 星历表文件路径，如果为空则使用内置 Moshier 算法
func InitSwissEphemeris(ephePath string) {
	if ephePath != "" {
		swephgo.SetEphePath([]byte(ephePath))
	}
	sweInitialized = true
}

// CloseSwissEphemeris 关闭 Swiss Ephemeris
func CloseSwissEphemeris() {
	swephgo.Close()
	sweInitialized = false
}

// IsSweAvailable 返回 Swiss Ephemeris 是否可用
func IsSweAvailable() bool {
	return sweAvailable
}

// ==================== 行星 ID 映射 ====================

// sweBodyMap Swiss Ephemeris 天体 ID 映射
var sweBodyMap = map[models.PlanetID]int{
	models.Sun:       swephgo.SeSun,
	models.Moon:      swephgo.SeMoon,
	models.Mercury:   swephgo.SeMercury,
	models.Venus:     swephgo.SeVenus,
	models.Mars:      swephgo.SeMars,
	models.Jupiter:   swephgo.SeJupiter,
	models.Saturn:    swephgo.SeSaturn,
	models.Uranus:    swephgo.SeUranus,
	models.Neptune:   swephgo.SeNeptune,
	models.Pluto:     swephgo.SePluto,
	models.NorthNode: swephgo.SeTrueNode,
	models.Chiron:    swephgo.SeChiron,
}

// ==================== 高精度行星位置计算 ====================

// CalculatePlanetPositionSwe 使用 Swiss Ephemeris 计算行星位置
func CalculatePlanetPositionSwe(planet models.PlanetID, jd float64) models.PlanetPosition {
	if !sweInitialized {
		InitSwissEphemeris("")
	}

	sweBody, ok := sweBodyMap[planet]
	if !ok {
		// 回退到内置算法
		return CalculatePlanetPosition(planet, jd)
	}

	// 计算标志：使用 Swiss Ephemeris + 速度
	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed

	// 准备输出缓冲区
	xx := make([]float64, 6)
	serr := make([]byte, 256)

	// 计算行星位置
	ret := swephgo.Calc(jd, sweBody, flag, xx, serr)
	if ret < 0 {
		// 如果失败，回退到内置算法
		return CalculatePlanetPosition(planet, jd)
	}

	longitude := NormalizeAngle(xx[0])
	latitude := xx[1]
	speed := xx[3]
	retrograde := speed < 0

	// 获取星座信息
	zodiac := GetZodiacByLongitude(longitude)
	signDegree := math.Mod(longitude, 30)

	// 计算尊贵度
	dignity := GetDignity(planet, zodiac.ID)
	dignityScore := GetDignityScore(dignity)

	planetInfo := GetPlanetInfo(planet)

	return models.PlanetPosition{
		ID:           planet,
		Name:         planetInfo.Name,
		Symbol:       planetInfo.Symbol,
		Longitude:    longitude,
		Latitude:     latitude,
		Sign:         zodiac.ID,
		SignName:     zodiac.Name,
		SignSymbol:   zodiac.Symbol,
		SignDegree:   signDegree,
		Retrograde:   retrograde,
		House:        0, // 需要额外计算
		DignityScore: dignityScore,
	}
}

// GetAllPlanetPositionsSwe 使用 Swiss Ephemeris 获取所有行星位置
func GetAllPlanetPositionsSwe(jd float64) []models.PlanetPosition {
	planets := []models.PlanetID{
		models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars,
		models.Jupiter, models.Saturn, models.Uranus, models.Neptune, models.Pluto,
		models.NorthNode, models.Chiron,
	}

	positions := make([]models.PlanetPosition, 0, len(planets))
	for _, p := range planets {
		positions = append(positions, CalculatePlanetPositionSwe(p, jd))
	}
	return positions
}

// ==================== 高精度宫位计算 ====================

// CalculateHousesSwe 使用 Swiss Ephemeris 计算宫位
func CalculateHousesSwe(jd float64, lat, lon float64) ([]models.HouseCusp, float64, float64) {
	if !sweInitialized {
		InitSwissEphemeris("")
	}

	// Placidus 分宫制 = 'P'
	cusps := make([]float64, 13)   // 13 个宫位尖端 (0 未使用)
	ascmc := make([]float64, 10)   // ASC, MC 等

	ret := swephgo.Houses(jd, lat, lon, 'P', cusps, ascmc)
	if ret < 0 {
		// 回退到内置算法
		return CalculateHouses(jd, lat, lon)
	}

	houses := make([]models.HouseCusp, 12)
	for i := 0; i < 12; i++ {
		cusp := NormalizeAngle(cusps[i+1]) // cusps[1] 是第一宫
		zodiac := GetZodiacByLongitude(cusp)
		houses[i] = models.HouseCusp{
			House:    i + 1,
			Cusp:     cusp,
			Sign:     zodiac.ID,
			SignName: zodiac.Name,
		}
	}

	asc := NormalizeAngle(ascmc[0]) // ASC
	mc := NormalizeAngle(ascmc[1])  // MC

	return houses, asc, mc
}
