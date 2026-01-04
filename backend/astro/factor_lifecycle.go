package astro

import (
	"math"
	"star/models"
	"time"
)

// ==================== 正弦曲线生命周期 ====================
// 基于占星学"入相-精确-离相"过程
// 使用正弦曲线模拟因子从开始到结束的强度变化

// CalculateFactorStrength 计算因子在某时刻的强度（正弦曲线）
// 返回值范围：0.0 - 1.0
//
// 曲线形态：
//
//	1.0 ┤     ╭───╮
//	    │    ╱     ╲
//	0.5 ┤   ╱       ╲
//	    │  ╱         ╲
//	0.0 ┼─╱───────────╲──
//	    开始   峰值   结束
func CalculateFactorStrength(lifecycle *models.FactorLifecycle, currentTime time.Time) float64 {
	if lifecycle == nil {
		return 1.0 // 无生命周期信息时返回满强度
	}

	startTime := lifecycle.StartTime
	endTime := lifecycle.EndTime

	// 检查是否在生效期间
	if currentTime.Before(startTime) || currentTime.After(endTime) {
		return 0.0
	}

	// 计算进度 (0.0 - 1.0)
	duration := endTime.Sub(startTime).Seconds()
	if duration <= 0 {
		return 1.0
	}

	elapsed := currentTime.Sub(startTime).Seconds()
	progress := elapsed / duration

	// 正弦曲线：sin(π × progress)
	// progress=0 时 strength=0
	// progress=0.5 时 strength=1 (峰值)
	// progress=1 时 strength=0
	strength := math.Sin(math.Pi * progress)

	return math.Max(0, math.Min(1, strength))
}

// CalculateFactorStrengthWithPeak 带峰值时间的强度计算
// 支持非对称曲线（峰值不一定在中间）
func CalculateFactorStrengthWithPeak(lifecycle *models.FactorLifecycle, currentTime time.Time) float64 {
	if lifecycle == nil {
		return 1.0
	}

	startTime := lifecycle.StartTime
	peakTime := lifecycle.PeakTime
	endTime := lifecycle.EndTime

	// 检查是否在生效期间
	if currentTime.Before(startTime) || currentTime.After(endTime) {
		return 0.0
	}

	var progress float64

	if currentTime.Before(peakTime) || currentTime.Equal(peakTime) {
		// 上升阶段：开始 → 峰值
		duration := peakTime.Sub(startTime).Seconds()
		if duration <= 0 {
			progress = 0.5
		} else {
			elapsed := currentTime.Sub(startTime).Seconds()
			progress = (elapsed / duration) * 0.5 // 映射到 0 - 0.5
		}
	} else {
		// 下降阶段：峰值 → 结束
		duration := endTime.Sub(peakTime).Seconds()
		if duration <= 0 {
			progress = 0.5
		} else {
			elapsed := currentTime.Sub(peakTime).Seconds()
			progress = 0.5 + (elapsed/duration)*0.5 // 映射到 0.5 - 1.0
		}
	}

	// 正弦曲线
	strength := math.Sin(math.Pi * progress)
	return math.Max(0, math.Min(1, strength))
}

// ==================== 因子持续时间配置 ====================

// FactorDurations 各类因子的典型持续时间（小时）
// 基于天文周期
var FactorDurations = map[models.InfluenceFactorType]float64{
	// 年度级
	models.FactorProfectionLord: 365 * 24, // 年主星：1年
	models.FactorOuterPlanet:    180 * 24, // 外行星相位：约6个月

	// 月度级
	models.FactorDignity: 30 * 24, // 行星换座：约30天（太阳周期）

	// 周度级
	models.FactorRetrograde: 21 * 24, // 水星逆行：约21天

	// 日度级
	models.FactorAspectPhase: 3 * 24,  // 相位影响：约3天（1°容许度/天）
	models.FactorAspectOrb:   3 * 24,  // 相位容许度
	models.FactorLunarPhase:  3.5 * 24, // 月相阶段：约3.5天

	// 小时级
	models.FactorPlanetaryHour: 1.5,   // 行星时：约1-1.5小时
	models.FactorVoidOfCourse:  12,    // 月亮空亡：数小时到半天
	models.FactorPersonal:      24,    // 个人因子：默认1天
	models.FactorCustom:        1,     // 自定义：默认1小时
}

// GetFactorDuration 获取因子的默认持续时间（小时）
func GetFactorDuration(factorType models.InfluenceFactorType) float64 {
	if duration, ok := FactorDurations[factorType]; ok {
		return duration
	}
	return 24 // 默认1天
}

// CreateLifecycle 创建因子生命周期
func CreateLifecycle(startTime time.Time, durationHours float64) *models.FactorLifecycle {
	endTime := startTime.Add(time.Duration(durationHours * float64(time.Hour)))
	peakTime := startTime.Add(time.Duration(durationHours / 2 * float64(time.Hour)))

	return &models.FactorLifecycle{
		StartTime: startTime,
		PeakTime:  peakTime,
		EndTime:   endTime,
		Duration:  durationHours,
	}
}

// CreateLifecycleWithPeak 创建带自定义峰值时间的生命周期
func CreateLifecycleWithPeak(startTime, peakTime, endTime time.Time) *models.FactorLifecycle {
	duration := endTime.Sub(startTime).Hours()
	return &models.FactorLifecycle{
		StartTime: startTime,
		PeakTime:  peakTime,
		EndTime:   endTime,
		Duration:  duration,
	}
}

// ==================== 相位生命周期计算 ====================

// CalculateAspectLifecycle 根据相位容许度计算生命周期
// 假设行星平均每天移动约1°
func CalculateAspectLifecycle(orb float64, exactTime time.Time, isApplying bool) *models.FactorLifecycle {
	// 容许度决定持续天数（每天约1°）
	durationDays := orb * 2 // 入相+离相
	durationHours := durationDays * 24

	// 计算开始和结束时间
	halfDuration := time.Duration(durationHours / 2 * float64(time.Hour))

	startTime := exactTime.Add(-halfDuration)
	endTime := exactTime.Add(halfDuration)

	// 如果正在入相，峰值在后面
	// 如果正在离相，峰值已过
	var peakTime time.Time
	if isApplying {
		peakTime = exactTime // 精确时刻即将到来
	} else {
		peakTime = exactTime // 精确时刻已过
	}

	return &models.FactorLifecycle{
		StartTime: startTime,
		PeakTime:  peakTime,
		EndTime:   endTime,
		Duration:  durationHours,
	}
}

// ==================== 逆行生命周期 ====================

// RetrogradeDurations 各行星逆行典型持续天数
var RetrogradeDurations = map[models.PlanetID]float64{
	models.Mercury: 21,  // 水星逆行：约21天
	models.Venus:   42,  // 金星逆行：约42天
	models.Mars:    72,  // 火星逆行：约72天
	models.Jupiter: 120, // 木星逆行：约4个月
	models.Saturn:  140, // 土星逆行：约4.5个月
	models.Uranus:  155, // 天王星逆行：约5个月
	models.Neptune: 160, // 海王星逆行：约5个月
	models.Pluto:   165, // 冥王星逆行：约5.5个月
}

// GetRetrogradeDuration 获取行星逆行持续时间（天）
func GetRetrogradeDuration(planet models.PlanetID) float64 {
	if duration, ok := RetrogradeDurations[planet]; ok {
		return duration
	}
	return 21 // 默认21天
}

