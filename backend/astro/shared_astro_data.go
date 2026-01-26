package astro

import (
	"star/models"
	"sync"
	"time"
)

// ==================== 共享星象数据层 ====================
// 统一管理星象事件计算，避免 Daily Events 和 Factor 系统重复计算
// 所有精确的星象事件都从这里获取

// ==================== 共享数据结构 ====================

// AspectEventData 相位事件数据（共享）
type AspectEventData struct {
	ExactTime    time.Time       // 精确时间
	ExactJD      float64         // 儒略日
	TransitPlanet models.PlanetID // 行运行星
	NatalPlanet  models.PlanetID // 本命行星
	AspectAngle  float64         // 相位角度
	AspectName   string          // 相位名称 (conjunction, sextile, square, trine, opposition)
	Orb          float64         // 容许度
	IsPositive   bool            // 是否正面相位
}

// LunarPhaseEventData 月相事件数据（共享）
type LunarPhaseEventData struct {
	ExactTime time.Time // 精确时间
	ExactJD   float64   // 儒略日
	Phase     string    // 月相类型 (new, first_quarter, full, last_quarter)
	PhaseName string    // 月相名称
	Angle     float64   // 日月角度
}

// SignChangeEventData 换座事件数据（共享）
type SignChangeEventData struct {
	ExactTime time.Time       // 精确时间
	Planet    models.PlanetID // 行星
	NewSign   string          // 新星座
}

// PlanetaryHourData 行星时数据（共享）
type PlanetaryHourData struct {
	StartTime     time.Time       // 开始时间
	EndTime       time.Time       // 结束时间
	Ruler         models.PlanetID // 主宰行星
	DayRuler      models.PlanetID // 日主星
	PlanetaryHour int             // 第几个行星时
	Influence     float64         // 影响值
}

// DailyAstroData 某一天的所有星象数据（共享）
type DailyAstroData struct {
	Date           time.Time             // 日期
	AspectEvents   []AspectEventData     // 相位事件
	LunarPhases    []LunarPhaseEventData // 月相事件
	SignChanges    []SignChangeEventData // 换座事件
	PlanetaryHours []PlanetaryHourData   // 行星时数据
	
	// 当前时刻的状态（用于分数计算）
	CurrentLunarPhase    *models.LunarPhaseInfo // 当前月相
	CurrentPlanetaryHour *PlanetaryHourInfo     // 当前行星时
}

// ==================== 缓存机制 ====================

type dailyAstroCacheKey struct {
	date     string // YYYY-MM-DD 格式
	chartID  string // 本命盘 ID（影响相位计算）
}

var (
	dailyAstroCache     = make(map[dailyAstroCacheKey]*DailyAstroData)
	dailyAstroCacheMu   sync.RWMutex
	dailyAstroCacheSize = 100 // 最多缓存 100 天的数据
)

// getChartID 生成本命盘的唯一标识
func getChartID(chart *models.NatalChart) string {
	if chart == nil {
		return "universal" // 通用查询（不考虑本命盘）
	}
	bd := chart.BirthData
	return bd.ToTime().Format("20060102150405")
}

// getCacheKey 生成缓存键
func getCacheKey(date time.Time, chart *models.NatalChart) dailyAstroCacheKey {
	return dailyAstroCacheKey{
		date:    date.Format("2006-01-02"),
		chartID: getChartID(chart),
	}
}

// ==================== 核心计算函数 ====================

// CalculateDailyAstroData 计算指定日期的所有星象数据（带缓存）
// 这是统一入口，Daily Events 和 Factor 系统都应该调用这个函数
func CalculateDailyAstroData(chart *models.NatalChart, date time.Time) *DailyAstroData {
	// 检查缓存
	cacheKey := getCacheKey(date, chart)
	
	dailyAstroCacheMu.RLock()
	if cached, ok := dailyAstroCache[cacheKey]; ok {
		dailyAstroCacheMu.RUnlock()
		return cached
	}
	dailyAstroCacheMu.RUnlock()
	
	// 计算数据
	data := calculateDailyAstroDataInternal(chart, date)
	
	// 存入缓存
	dailyAstroCacheMu.Lock()
	if len(dailyAstroCache) >= dailyAstroCacheSize {
		// 清理一半的缓存（简单策略）
		count := 0
		for k := range dailyAstroCache {
			delete(dailyAstroCache, k)
			count++
			if count >= dailyAstroCacheSize/2 {
				break
			}
		}
	}
	dailyAstroCache[cacheKey] = data
	dailyAstroCacheMu.Unlock()
	
	return data
}

// calculateDailyAstroDataInternal 内部计算函数（无缓存）
func calculateDailyAstroDataInternal(chart *models.NatalChart, date time.Time) *DailyAstroData {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	data := &DailyAstroData{
		Date:           date,
		AspectEvents:   []AspectEventData{},
		LunarPhases:    []LunarPhaseEventData{},
		SignChanges:    []SignChangeEventData{},
		PlanetaryHours: []PlanetaryHourData{},
	}
	
	// 1. 计算相位事件（需要本命盘）
	if chart != nil {
		data.AspectEvents = calculateSharedAspectEvents(chart, startOfDay, endOfDay)
	}
	
	// 2. 计算月相事件（通用）
	data.LunarPhases = calculateSharedLunarPhaseEvents(startOfDay, endOfDay)
	
	// 3. 计算换座事件（通用）
	data.SignChanges = calculateSharedSignChangeEvents(startOfDay, endOfDay)
	
	// 4. 计算行星时数据（通用）
	data.PlanetaryHours = calculateSharedPlanetaryHours(startOfDay, endOfDay)
	
	// 5. 获取当前时刻的状态
	data.CurrentLunarPhase = getCurrentLunarPhaseInfo(date)
	data.CurrentPlanetaryHour = getCurrentPlanetaryHourInfo(date, chart)
	
	return data
}

// ==================== 相位事件计算 ====================

// SharedAspectInfo 共享相位信息（用于避免与 constants.go 的冲突）
type SharedAspectInfo struct {
	Angle      float64
	Name       string // English name
	Orb        float64
	IsPositive bool
}

// 主要相位定义（共享版本）
var SharedMajorAspects = []SharedAspectInfo{
	{0, "conjunction", 8.0, true},
	{60, "sextile", 6.0, true},
	{90, "square", 8.0, false},
	{120, "trine", 8.0, true},
	{180, "opposition", 8.0, false},
}

// 次要相位定义（共享版本）
var SharedMinorAspects = []SharedAspectInfo{
	{30, "semi-sextile", 2.0, true},
	{45, "semi-square", 2.0, false},
	{135, "sesquiquadrate", 2.0, false},
	{150, "quincunx", 2.0, false},
}

// calculateSharedAspectEvents 计算相位事件（共享）
func calculateSharedAspectEvents(chart *models.NatalChart, startTime, endTime time.Time) []AspectEventData {
	events := []AspectEventData{}
	
	// 行运行星
	transitPlanets := []models.PlanetID{
		models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars,
		models.Jupiter, models.Saturn,
	}
	
	searchStart := TimeToJulianDay(startTime)
	searchEnd := TimeToJulianDay(endTime)
	
	// 对每个行运行星和本命行星的组合
	for _, transitPlanet := range transitPlanets {
		for _, natalPlanet := range chart.Planets {
			for _, aspect := range SharedMajorAspects {
				exactJD, found := FindExactAspectTime(
					transitPlanet,
					natalPlanet.ID,
					aspect.Angle,
					searchStart,
					searchEnd,
				)
				
				if found {
					events = append(events, AspectEventData{
						ExactTime:     JulianDayToTime(exactJD),
						ExactJD:       exactJD,
						TransitPlanet: transitPlanet,
						NatalPlanet:   natalPlanet.ID,
						AspectAngle:   aspect.Angle,
						AspectName:    aspect.Name,
						Orb:           aspect.Orb,
						IsPositive:    aspect.IsPositive,
					})
				}
			}
		}
	}
	
	return events
}

// ==================== 月相事件计算 ====================

// calculateSharedLunarPhaseEvents 计算月相事件（共享）
func calculateSharedLunarPhaseEvents(startTime, endTime time.Time) []LunarPhaseEventData {
	events := []LunarPhaseEventData{}
	
	phaseDefinitions := []struct {
		Angle float64
		Phase string
		Name  string
	}{
		{0, "new", "New Moon"},
		{90, "first_quarter", "First Quarter"},
		{180, "full", "Full Moon"},
		{270, "last_quarter", "Last Quarter"},
	}
	
	for _, phase := range phaseDefinitions {
		exactTime := findLunarPhaseTime(phase.Angle, startTime, endTime)
		if exactTime != nil {
			events = append(events, LunarPhaseEventData{
				ExactTime: *exactTime,
				ExactJD:   TimeToJulianDay(*exactTime),
				Phase:     phase.Phase,
				PhaseName: phase.Name,
				Angle:     phase.Angle,
			})
		}
	}
	
	return events
}

// getCurrentLunarPhaseInfo 获取当前月相信息
func getCurrentLunarPhaseInfo(date time.Time) *models.LunarPhaseInfo {
	jd := TimeToJulianDay(date)
	sunPos := CalculatePlanetPositionSwe(models.Sun, jd)
	moonPos := CalculatePlanetPositionSwe(models.Moon, jd)
	angle := NormalizeAngle(moonPos.Longitude - sunPos.Longitude)
	phaseInfo := GetLunarPhase(angle)
	return &phaseInfo
}

// ==================== 换座事件计算 ====================

// calculateSharedSignChangeEvents 计算换座事件（共享）
func calculateSharedSignChangeEvents(startTime, endTime time.Time) []SignChangeEventData {
	events := []SignChangeEventData{}
	
	planetsToCheck := []models.PlanetID{
		models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars,
	}
	
	for _, planetID := range planetsToCheck {
		exactTime := findSignChangeTime(planetID, startTime, endTime)
		if exactTime != nil {
			pos := CalculatePlanetPositionSwe(planetID, TimeToJulianDay(*exactTime))
			zodiacInfo := GetZodiacByLongitude(pos.Longitude)
			events = append(events, SignChangeEventData{
				ExactTime: *exactTime,
				Planet:    planetID,
				NewSign:   string(zodiacInfo.ID),
			})
		}
	}
	
	return events
}

// ==================== 行星时计算 ====================

// calculateSharedPlanetaryHours 计算行星时数据（共享）
func calculateSharedPlanetaryHours(startTime, endTime time.Time) []PlanetaryHourData {
	hours := []PlanetaryHourData{}
	
	// 简化计算：每2小时一个行星时
	// 完整计算需要考虑日出日落时间
	currentTime := startTime
	hourIndex := 0
	
	for currentTime.Before(endTime) {
		nextTime := currentTime.Add(2 * time.Hour)
		if nextTime.After(endTime) {
			nextTime = endTime
		}
		
		hourInfo := CalculatePlanetaryHourEnhanced(currentTime.Add(time.Hour), 0, 0)
		
		hours = append(hours, PlanetaryHourData{
			StartTime:     currentTime,
			EndTime:       nextTime,
			Ruler:         hourInfo.Ruler,
			DayRuler:      hourInfo.DayRuler,
			PlanetaryHour: hourIndex + 1,
			Influence:     hourInfo.Influence,
		})
		
		currentTime = nextTime
		hourIndex++
	}
	
	return hours
}

// getCurrentPlanetaryHourInfo 获取当前行星时信息
func getCurrentPlanetaryHourInfo(date time.Time, chart *models.NatalChart) *PlanetaryHourInfo {
	lat, lon := 0.0, 0.0
	if chart != nil && chart.BirthData.Latitude != 0 {
		lat = chart.BirthData.Latitude
		lon = chart.BirthData.Longitude
	}
	
	hourInfo := CalculatePlanetaryHourEnhanced(date, lat, lon)
	return &hourInfo
}

// ==================== 工具函数 ====================

// GetSharedAspectByAngle 根据角度获取共享相位定义
func GetSharedAspectByAngle(angle float64) *SharedAspectInfo {
	for i := range SharedMajorAspects {
		if SharedMajorAspects[i].Angle == angle {
			return &SharedMajorAspects[i]
		}
	}
	for i := range SharedMinorAspects {
		if SharedMinorAspects[i].Angle == angle {
			return &SharedMinorAspects[i]
		}
	}
	return nil
}

// FindAspectEventAt 在共享数据中查找指定时间附近的相位事件
func FindAspectEventAt(data *DailyAstroData, transitPlanet, natalPlanet models.PlanetID, aspectAngle float64) *AspectEventData {
	for i := range data.AspectEvents {
		e := &data.AspectEvents[i]
		if e.TransitPlanet == transitPlanet && e.NatalPlanet == natalPlanet && e.AspectAngle == aspectAngle {
			return e
		}
	}
	return nil
}

// ClearDailyAstroCache 清除缓存（用于测试或内存管理）
func ClearDailyAstroCache() {
	dailyAstroCacheMu.Lock()
	dailyAstroCache = make(map[dailyAstroCacheKey]*DailyAstroData)
	dailyAstroCacheMu.Unlock()
}
