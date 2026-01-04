package astro

import (
	"star/models"
	"time"
)

// CalculateProgressions 计算推运
func CalculateProgressions(chart *models.NatalChart, targetDateStr string) *models.ProgressedChart {
	targetDate, err := time.Parse("2006-01-02", targetDateStr)
	if err != nil {
		targetDate = time.Now()
	}

	birthDate := chart.BirthData.ToTime()

	// 计算从出生到目标日期的年数
	duration := targetDate.Sub(birthDate)
	yearsFromBirth := duration.Hours() / (365.25 * 24)

	// 推运日期（1天 = 1年）
	daysProgressed := yearsFromBirth
	progressedDate := birthDate.AddDate(0, 0, int(daysProgressed))

	// 计算推运行星位置
	progressedJd := DateToJulianDay(progressedDate)
	progressedPositions := GetAllPlanetPositions(progressedJd)

	// 创建推运行星列表
	progressedPlanets := make([]models.ProgressedPlanet, len(progressedPositions))
	for i, progPos := range progressedPositions {
		natalPos := GetPlanetFromChart(chart, progPos.ID)
		var natalLon float64
		var signChanged, houseChanged bool

		if natalPos != nil {
			natalLon = natalPos.Longitude
			signChanged = natalPos.Sign != progPos.Sign
			houseChanged = natalPos.House != progPos.House
		}

		progressedPlanets[i] = models.ProgressedPlanet{
			ID:                  progPos.ID,
			Name:                progPos.Name,
			Symbol:              progPos.Symbol,
			NatalLongitude:      natalLon,
			ProgressedLongitude: progPos.Longitude,
			Movement:            progPos.Longitude - natalLon,
			Sign:                progPos.Sign,
			SignName:            progPos.SignName,
			SignChanged:         signChanged,
			House:               progPos.House,
			HouseChanged:        houseChanged,
		}
	}

	// 计算推运ASC和MC
	progressedHouses, progressedAsc, progressedMc := CalculateHouses(
		progressedJd,
		chart.BirthData.Latitude,
		chart.BirthData.Longitude,
	)
	_ = progressedHouses

	// 计算推运相位
	aspects := CalculateAspects(progressedPositions)

	// 计算推运月相
	var sunLon, moonLon float64
	for _, p := range progressedPositions {
		if p.ID == models.Sun {
			sunLon = p.Longitude
		}
		if p.ID == models.Moon {
			moonLon = p.Longitude
		}
	}
	lunarPhase := calculateProgressedLunarPhase(sunLon, moonLon)

	return &models.ProgressedChart{
		TargetDate:          targetDate,
		ProgressedDate:      progressedDate,
		DaysProgressed:      daysProgressed,
		YearsFromBirth:      yearsFromBirth,
		Planets:             progressedPlanets,
		ProgressedAscendant: progressedAsc,
		ProgressedMidheaven: progressedMc,
		Aspects:             aspects,
		LunarPhase:          lunarPhase,
	}
}

// calculateProgressedLunarPhase 计算推运月相
func calculateProgressedLunarPhase(sunLon, moonLon float64) models.LunarPhaseInfo {
	angle := NormalizeAngle(moonLon - sunLon)
	phaseInfo := GetLunarPhase(angle)

	// 添加描述
	descriptions := map[string]string{
		"new":          "新开始、播种期，适合设定意图和启动新项目",
		"crescent":     "行动期，需要突破阻力，建立基础",
		"firstQuarter": "危机调整期，需要做出重要决定",
		"gibbous":      "完善准备期，分析和优化",
		"full":         "高峰收获期，事情明朗化，情感充沛",
		"disseminating": "分享传播期，将收获与他人分享",
		"lastQuarter":  "反思释放期，放下不再需要的",
		"balsamic":     "整合休息期，为新周期做准备",
	}

	phaseInfo.Description = descriptions[phaseInfo.Phase]

	return phaseInfo
}

// GetProgressedLunarPhaseForAge 获取指定年龄的推运月相
func GetProgressedLunarPhaseForAge(chart *models.NatalChart, age int) models.LunarPhaseInfo {
	birthDate := chart.BirthData.ToTime()
	targetDate := birthDate.AddDate(age, 0, 0)

	progressedChart := CalculateProgressions(chart, targetDate.Format("2006-01-02"))
	return progressedChart.LunarPhase
}

// CalculateProgressedLunarCycle 计算推运月亮周期
func CalculateProgressedLunarCycle(chart *models.NatalChart) []struct {
	Age   int
	Phase models.LunarPhaseInfo
} {
	var cycle []struct {
		Age   int
		Phase models.LunarPhaseInfo
	}

	// 推运月相周期约27.3年
	for age := 0; age <= 80; age++ {
		phase := GetProgressedLunarPhaseForAge(chart, age)
		cycle = append(cycle, struct {
			Age   int
			Phase models.LunarPhaseInfo
		}{age, phase})
	}

	return cycle
}

