//go:build swe
// +build swe

package astro

import (
	"math"
	"star/models"
	"time"

	"github.com/mshafiee/swephgo"
)

// ==================== 精确相位时间搜索 ====================
// 使用二分法（Bisection）算法精确计算相位发生时间

// 黄经缓存
type longitudeCacheKey struct {
	planet models.PlanetID
	jd     float64
}

var longitudeCache = make(map[longitudeCacheKey]float64)
var longitudeCacheMaxSize = 500

// GetPlanetLongitudeAt 获取行星在指定儒略日的黄经
func GetPlanetLongitudeAt(planet models.PlanetID, jd float64) float64 {
	if !sweInitialized {
		InitSwissEphemeris("")
	}

	// 检查缓存（精确到6位小数）
	roundedJD := math.Floor(jd*1000000+0.5) / 1000000
	cacheKey := longitudeCacheKey{planet: planet, jd: roundedJD}
	if cached, ok := longitudeCache[cacheKey]; ok {
		return cached
	}

	sweBody, ok := sweBodyMap[planet]
	if !ok {
		return 0
	}

	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed
	xx := make([]float64, 6)
	serr := make([]byte, 256)

	ret := swephgo.Calc(jd, sweBody, flag, xx, serr)
	if ret < 0 {
		return 0
	}

	longitude := NormalizeAngle(xx[0])

	// 保存到缓存
	if len(longitudeCache) < longitudeCacheMaxSize {
		longitudeCache[cacheKey] = longitude
	}

	return longitude
}

// normalizeAngleDiff 将角度差归一化到 -180° 到 +180° 范围
func normalizeAngleDiff(diff float64) float64 {
	diff = math.Mod(diff, 360)
	if diff > 180 {
		diff -= 360
	} else if diff < -180 {
		diff += 360
	}
	return diff
}

// aspectFunction 计算两颗行星的角度差与目标相位角的偏差
// 返回值：当返回值为0时，表示精确形成相位
func aspectFunction(planet1, planet2 models.PlanetID, targetAngle, jd float64) float64 {
	lon1 := GetPlanetLongitudeAt(planet1, jd)
	lon2 := GetPlanetLongitudeAt(planet2, jd)
	
	// 计算角度差（0-360范围）
	diff := math.Mod(lon1-lon2+360, 360)
	
	// 对于相位，我们需要检查两种情况：
	// diff ≈ targetAngle 或 diff ≈ 360 - targetAngle
	// 返回最小的偏差
	
	deviation1 := diff - targetAngle
	deviation2 := diff - (360 - targetAngle)
	
	// 归一化偏差到 -180 到 +180
	deviation1 = normalizeAngleDiff(deviation1)
	deviation2 = normalizeAngleDiff(deviation2)
	
	// 返回绝对值更小的那个偏差
	if math.Abs(deviation1) < math.Abs(deviation2) {
		return deviation1
	}
	return deviation2
}

// FindExactAspectTime 使用二分法查找精确的相位时间
// 返回相位精确发生的儒略日时间
func FindExactAspectTime(planet1, planet2 models.PlanetID, targetAngle float64, startJd, endJd float64) (float64, bool) {
	const maxIterations = 50
	const tolerance = 1e-6 // 约 0.086 秒的精度

	// 检查是否存在符号变化（即是否有根）
	fStart := aspectFunction(planet1, planet2, targetAngle, startJd)
	fEnd := aspectFunction(planet1, planet2, targetAngle, endJd)

	// 如果两端符号相同，尝试在中间找到符号变化
	if fStart*fEnd > 0 {
		// 尝试更细的搜索
		step := (endJd - startJd) / 10
		found := false
		for t := startJd; t < endJd; t += step {
			fT := aspectFunction(planet1, planet2, targetAngle, t)
			fNext := aspectFunction(planet1, planet2, targetAngle, t+step)
			if fT*fNext <= 0 {
				startJd = t
				endJd = t + step
				fStart = fT
				fEnd = fNext
				found = true
				break
			}
		}
		if !found {
			return 0, false
		}
	}

	// 二分法搜索
	low, high := startJd, endJd
	for i := 0; i < maxIterations; i++ {
		mid := (low + high) / 2
		fMid := aspectFunction(planet1, planet2, targetAngle, mid)

		if math.Abs(fMid) < tolerance {
			return mid, true
		}

		if fStart*fMid <= 0 {
			high = mid
			fEnd = fMid
		} else {
			low = mid
			fStart = fMid
		}

		if high-low < tolerance {
			return (low + high) / 2, true
		}
	}

	return (low + high) / 2, true
}

// FindExactAspectTimeFromQuery 从查询时间出发，查找最近的相位精确时间
// 根据 isApplying 决定搜索方向
func FindExactAspectTimeFromQuery(planet1, planet2 models.PlanetID, targetAngle float64, orb float64, queryTime time.Time, isApplying bool) time.Time {
	// 将查询时间转换为儒略日
	queryJd := TimeToJulianDay(queryTime)

	// 估算搜索范围（基于行星速度）
	speed := GetPlanetSpeed(planet1)
	if speed < 0.01 {
		speed = 0.01
	}
	
	// 搜索范围：orb / speed 天
	searchDays := orb / speed
	if searchDays < 1 {
		searchDays = 1
	}
	if searchDays > 30 {
		searchDays = 30 // 限制最大搜索范围
	}

	var startJd, endJd float64
	if isApplying {
		// 入相：精确时间在未来
		startJd = queryJd
		endJd = queryJd + searchDays
	} else {
		// 离相：精确时间在过去
		startJd = queryJd - searchDays
		endJd = queryJd
	}

	exactJd, found := FindExactAspectTime(planet1, planet2, targetAngle, startJd, endJd)
	if !found {
		// 如果没找到，返回估算时间
		if isApplying {
			return queryTime.Add(time.Duration(orb/speed*24) * time.Hour)
		}
		return queryTime.Add(-time.Duration(orb/speed*24) * time.Hour)
	}

	return JulianDayToTime(exactJd)
}

// TimeToJulianDay 将 time.Time 转换为儒略日
func TimeToJulianDay(t time.Time) float64 {
	// 转换为 UTC
	t = t.UTC()
	
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	hour := float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600 + float64(t.Nanosecond())/3600e9
	
	// 使用 Swiss Ephemeris 的儒略日计算
	jd := swephgo.Julday(year, month, day, hour, swephgo.SeGregCal)
	return jd
}

// JulianDayToTime 将儒略日转换为 time.Time
func JulianDayToTime(jd float64) time.Time {
	// 使用 Swiss Ephemeris 逆向计算
	year := make([]int, 1)
	month := make([]int, 1)
	day := make([]int, 1)
	hour := make([]float64, 1)
	
	swephgo.Revjul(jd, swephgo.SeGregCal, year, month, day, hour)
	
	// 转换小时为时分秒
	h := int(hour[0])
	remainingMinutes := (hour[0] - float64(h)) * 60
	m := int(remainingMinutes)
	remainingSeconds := (remainingMinutes - float64(m)) * 60
	s := int(remainingSeconds)
	ns := int((remainingSeconds - float64(s)) * 1e9)
	
	return time.Date(year[0], time.Month(month[0]), day[0], h, m, s, ns, time.UTC)
}
