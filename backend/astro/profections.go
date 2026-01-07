package astro

import (
	"fmt"
	"star/models"
	"time"
)

// CalculateAnnualProfection 计算年限法
func CalculateAnnualProfection(chart *models.NatalChart, age int) *models.AnnualProfection {
	// 计算年度宫位 (年龄 % 12 + 1)
	house := (age % 12) + 1

	// 获取宫位信息
	houseInfo := GetHouseInfo(house)
	if houseInfo == nil {
		houseInfo = &HouseInfo{House: house, Name: "Unknown House", Theme: "Unknown", Keywords: []string{}}
	}

	// 获取该宫位的星座
	var houseCusp float64
	for _, h := range chart.Houses {
		if h.House == house {
			houseCusp = h.Cusp
			break
		}
	}

	zodiac := GetZodiacByLongitude(houseCusp)
	if zodiac == nil {
		zodiac = &ZodiacSigns[0]
	}

	// 获取守护星（年主星）
	lordOfYear := zodiac.Ruler
	lordInfo := GetPlanetInfo(lordOfYear)

	// 获取年主星在本命盘中的位置
	var lordNatalHouse int
	var lordNatalSign models.ZodiacID
	lordPlanet := GetPlanetFromChart(chart, lordOfYear)
	if lordPlanet != nil {
		lordNatalHouse = lordPlanet.House
		lordNatalSign = lordPlanet.Sign
	}

	// 计算当前年份
	currentYear := chart.BirthData.ToTime().Year() + age

	// 生成描述
	description := generateProfectionDescription(house, houseInfo, zodiac, lordInfo)

	return &models.AnnualProfection{
		Year:           currentYear,
		Age:            age,
		House:          house,
		HouseName:      houseInfo.Name,
		HouseTheme:     houseInfo.Theme,
		HouseKeywords:  houseInfo.Keywords,
		Sign:           zodiac.ID,
		SignName:       zodiac.Name,
		LordOfYear:     lordOfYear,
		LordName:       lordInfo.Name,
		LordSymbol:     lordInfo.Symbol,
		LordNatalHouse: lordNatalHouse,
		LordNatalSign:  lordNatalSign,
		Description:    description,
	}
}

// generateProfectionDescription 生成年限法描述
func generateProfectionDescription(house int, houseInfo *HouseInfo, _ *ZodiacInfo, lordInfo *PlanetInfo) string {
	return fmt.Sprintf(
		"This is a %s year, theme: %s. Annual lord is %s (%s), "+
			"main energy focuses on %s related areas. Key areas: %s.",
		houseInfo.Name, houseInfo.Theme,
		lordInfo.Name, lordInfo.Symbol,
		houseInfo.Theme,
		joinKeywords(houseInfo.Keywords),
	)
}

// joinKeywords 连接关键词
func joinKeywords(keywords []string) string {
	if len(keywords) == 0 {
		return ""
	}
	result := keywords[0]
	for i := 1; i < len(keywords); i++ {
		result += ", " + keywords[i]
	}
	return result
}

// CalculateProfectionMap 计算完整年限法地图
func CalculateProfectionMap(chart *models.NatalChart) *models.LifeProfectionMap {
	// 计算0-80岁的年限法
	profections := make([]models.AnnualProfection, 81)
	for age := 0; age <= 80; age++ {
		profections[age] = *CalculateAnnualProfection(chart, age)
	}

	// 计算当前年龄
	now := time.Now()
	birthDate := chart.BirthData.ToTime()
	currentAge := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		currentAge--
	}
	if currentAge < 0 {
		currentAge = 0
	}
	if currentAge > 80 {
		currentAge = 80
	}

	// 获取当前年和未来5年
	var currentYear *models.AnnualProfection
	var upcomingYears []models.AnnualProfection

	if currentAge >= 0 && currentAge <= 80 {
		currentYear = &profections[currentAge]
	}

	for i := 1; i <= 5; i++ {
		age := currentAge + i
		if age <= 80 {
			upcomingYears = append(upcomingYears, profections[age])
		}
	}

	// 周期分析
	currentCycleNumber := currentAge/12 + 1
	yearsIntoCurrentCycle := currentAge % 12

	cycleAnalysis := models.CycleAnalysis{
		FirstCycle: map[string]interface{}{
			"startAge": 0,
			"endAge":   11,
			"theme":    "Exploration & Learning",
		},
		SecondCycle: map[string]interface{}{
			"startAge": 12,
			"endAge":   23,
			"theme":    "Building Self",
		},
		ThirdCycle: map[string]interface{}{
			"startAge": 24,
			"endAge":   35,
			"theme":    "Achievement & Responsibility",
		},
		CurrentCycleNumber:    currentCycleNumber,
		YearsIntoCurrentCycle: yearsIntoCurrentCycle,
	}

	return &models.LifeProfectionMap{
		Profections:   profections,
		CurrentYear:   currentYear,
		UpcomingYears: upcomingYears,
		CycleAnalysis: cycleAnalysis,
	}
}

// GetProfectionForDate 获取指定日期的年限法
func GetProfectionForDate(chart *models.NatalChart, date time.Time) *models.AnnualProfection {
	birthDate := chart.BirthData.ToTime()
	age := date.Year() - birthDate.Year()
	if date.YearDay() < birthDate.YearDay() {
		age--
	}
	if age < 0 {
		age = 0
	}
	return CalculateAnnualProfection(chart, age)
}

