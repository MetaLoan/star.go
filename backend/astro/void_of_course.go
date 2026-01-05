package astro

import (
	"math"
	"star/models"
)

// ==================== 月亮空亡计算 ====================

// VoidOfCourseInfo 月亮空亡信息
type VoidOfCourseInfo struct {
	IsVoid       bool    `json:"isVoid"`
	StartTime    string  `json:"startTime,omitempty"`
	EndTime      string  `json:"endTime,omitempty"`
	Duration     float64 `json:"duration"` // 小时
	NextSign     string  `json:"nextSign,omitempty"`
	LastAspect   string  `json:"lastAspect,omitempty"`
	Influence    float64 `json:"influence"` // -15 到 0
}

// CalculateVoidOfCourse 计算月亮空亡
// 月亮空亡：月亮在当前星座内不再形成主要相位，直到进入下一个星座
func CalculateVoidOfCourse(jd float64, natalPositions []models.PlanetPosition) VoidOfCourseInfo {
	moonPos := CalculatePlanetPositionUnified(models.Moon, jd)
	currentSign := moonPos.Sign

	// 获取当前行运行星位置 - 使用 Swiss Ephemeris
	transitPositions := GetPlanetPositionsUnified(jd)

	// 月亮在当前星座内剩余的度数
	degreesLeftInSign := 30 - moonPos.SignDegree

	// 月亮每天移动约 13.2 度
	moonDailyMotion := 13.2
	hoursToNextSign := (degreesLeftInSign / moonDailyMotion) * 24

	// 检查月亮在离开当前星座前是否会形成任何主要相位
	hasUpcomingAspect := false
	var lastAspectDesc string

	// 搜索接下来的时间内是否有相位
	searchSteps := 48 // 搜索 48 小时
	for step := 1; step <= searchSteps; step++ {
		futureJd := jd + float64(step)/24.0 // 每小时一步
		futureMoonPos := CalculatePlanetPositionUnified(models.Moon, futureJd)

		// 如果月亮已经换座，停止搜索
		if futureMoonPos.Sign != currentSign {
			break
		}

		// 检查是否有主要相位形成
		for _, transit := range transitPositions {
			if transit.ID == models.Moon {
				continue
			}

			// 计算角距
			diff := math.Abs(futureMoonPos.Longitude - transit.Longitude)
			if diff > 180 {
				diff = 360 - diff
			}

			// 检查主要相位（容许度 1°）
			for _, angle := range []float64{0, 60, 90, 120, 180} {
				if math.Abs(diff-angle) < 1.0 {
					hasUpcomingAspect = true
					aspectName := getVocAspectName(angle)
					planetInfo := GetPlanetInfo(transit.ID)
					lastAspectDesc = "月亮" + aspectName + planetInfo.Name
					break
				}
			}
			if hasUpcomingAspect {
				break
			}
		}
		if hasUpcomingAspect {
			break
		}
	}

	isVoid := !hasUpcomingAspect && hoursToNextSign > 0.5 // 至少 30 分钟

	// 计算影响值
	var influence float64
	if isVoid {
		// 空亡时间越长，影响越大（最大 -15）
		influence = -math.Min(15, hoursToNextSign*0.5)
	}

	// 获取下一个星座
	nextSignIdx := (int(moonPos.Longitude/30) + 1) % 12
	nextSign := ZodiacSigns[nextSignIdx].Name

	return VoidOfCourseInfo{
		IsVoid:      isVoid,
		Duration:    hoursToNextSign,
		NextSign:    nextSign,
		LastAspect:  lastAspectDesc,
		Influence:   influence,
	}
}

// IsVoidOfCourseMoon 简化版月亮空亡检测
func IsVoidOfCourseMoon(jd float64) bool {
	info := CalculateVoidOfCourse(jd, nil)
	return info.IsVoid
}

// getVocAspectName 获取相位中文名
func getVocAspectName(angle float64) string {
	switch int(angle) {
	case 0:
		return "合"
	case 60:
		return "六分"
	case 90:
		return "刑"
	case 120:
		return "拱"
	case 180:
		return "冲"
	default:
		return ""
	}
}

