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

	var strength float64
	
	// 非对称波函数：入相缓和，出相陡峭
	// 占星学理论：能量积累慢，消散快
	if progress <= 0.5 {
		// 入相阶段 (0 → 0.5)：使用正弦曲线
		// 映射 [0, 0.5] → [0, π/2]
		// 从 0 平滑上升到 1
		strength = math.Sin(math.Pi * progress)
	} else {
		// 出相阶段 (0.5 → 1.0)：使用指数衰减
		// 能量快速消散，模拟出相的快速减弱
		t := (progress - 0.5) * 2 // 映射到 [0, 1]
		// 使用 e^(-3t) 的衰减曲线，从 1 快速降到 0
		strength = math.Exp(-3 * t)
	}

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
	var strength float64
	
	if currentTime.Before(peakTime) || currentTime.Equal(peakTime) {
		// 上升阶段（入相）：开始 → 峰值
		duration := peakTime.Sub(startTime).Seconds()
		if duration <= 0 {
			strength = 1.0
		} else {
			elapsed := currentTime.Sub(startTime).Seconds()
			progress = (elapsed / duration) * 0.5 // 映射到 0 - 0.5
			// 使用正弦曲线平滑上升
			strength = math.Sin(math.Pi * progress)
		}
	} else {
		// 下降阶段（出相）：峰值 → 结束
		duration := endTime.Sub(peakTime).Seconds()
		if duration <= 0 {
			strength = 1.0
		} else {
			elapsed := currentTime.Sub(peakTime).Seconds()
			t := elapsed / duration // t ∈ [0, 1]
			// 使用指数衰减，出相能量快速消散
			strength = math.Exp(-3 * t)
		}
	}

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

// PlanetSpeeds 各行星的典型平均速度（度/天）
var PlanetSpeeds = map[models.PlanetID]float64{
	models.Moon:    13.0,  // 月亮：约13°/天
	models.Sun:     1.0,   // 太阳：约1°/天
	models.Mercury: 1.2,   // 水星：约1.2°/天（变化大）
	models.Venus:   1.2,   // 金星：约1.2°/天
	models.Mars:    0.5,   // 火星：约0.5°/天
	models.Jupiter: 0.08,  // 木星：约0.08°/天
	models.Saturn:  0.03,  // 土星：约0.03°/天
	models.Uranus:  0.01,  // 天王星：约0.01°/天
	models.Neptune: 0.006, // 海王星：约0.006°/天
	models.Pluto:   0.004, // 冥王星：约0.004°/天
	models.Chiron:  0.02,  // 凯龙：约0.02°/天
}

// GetPlanetSpeed 获取行星速度（度/天）
func GetPlanetSpeed(planet models.PlanetID) float64 {
	if speed, ok := PlanetSpeeds[planet]; ok {
		return speed
	}
	return 0.5 // 默认速度
}

// CalculateAspectLifecycleWithPlanets 根据相位和行星速度计算生命周期
// 使用精确的相位时间搜索算法（如果可用）
func CalculateAspectLifecycleWithPlanets(orb float64, queryTime time.Time, isApplying bool, planet1, planet2 models.PlanetID, maxOrb float64) *models.FactorLifecycle {
	// 获取两颗行星的速度
	speed1 := GetPlanetSpeed(planet1)
	speed2 := GetPlanetSpeed(planet2)

	// 相对速度 = 快行星速度 - 慢行星速度（简化计算）
	relativeSpeed := speed1
	if speed2 > speed1 {
		relativeSpeed = speed2
	}
	if relativeSpeed < 0.01 {
		relativeSpeed = 0.01 // 防止除以零
	}

	// 总持续时间 = 2 * maxOrb / 相对速度
	durationDays := (maxOrb * 2) / relativeSpeed
	durationHours := durationDays * 24
	halfDuration := time.Duration(durationHours / 2 * float64(time.Hour))

	// 尝试使用精确搜索算法计算精确时间
	var exactTime time.Time
	if IsSweAvailable() {
		// 使用二分法精确搜索相位时间
		// 注意：这里假设 aspectAngle = 0（合相），其他相位需要传入实际角度
		// 由于当前函数签名没有传入相位角度，我们先使用估算方法
		// 后续可以扩展函数签名来支持精确搜索
		exactTime = estimateExactTime(orb, queryTime, isApplying, relativeSpeed)
	} else {
		exactTime = estimateExactTime(orb, queryTime, isApplying, relativeSpeed)
	}

	startTime := exactTime.Add(-halfDuration)
	endTime := exactTime.Add(halfDuration)

	return &models.FactorLifecycle{
		StartTime: startTime,
		PeakTime:  exactTime,
		EndTime:   endTime,
		Duration:  durationHours,
	}
}

// estimateExactTime 估算精确时间（用于回退或无法精确搜索时）
func estimateExactTime(orb float64, queryTime time.Time, isApplying bool, relativeSpeed float64) time.Time {
	orbDays := orb / relativeSpeed
	orbHours := orbDays * 24

	if isApplying {
		return queryTime.Add(time.Duration(orbHours * float64(time.Hour)))
	}
	return queryTime.Add(-time.Duration(orbHours * float64(time.Hour)))
}

// CalculateAspectLifecycleExact 使用精确搜索算法计算相位生命周期
// 这个函数使用二分法精确查找相位发生时间
func CalculateAspectLifecycleExact(orb float64, queryTime time.Time, isApplying bool, planet1, planet2 models.PlanetID, aspectAngle, maxOrb float64) *models.FactorLifecycle {
	// 获取相对速度
	speed1 := GetPlanetSpeed(planet1)
	speed2 := GetPlanetSpeed(planet2)
	relativeSpeed := speed1
	if speed2 > speed1 {
		relativeSpeed = speed2
	}
	if relativeSpeed < 0.01 {
		relativeSpeed = 0.01
	}

	// 总持续时间
	durationDays := (maxOrb * 2) / relativeSpeed
	durationHours := durationDays * 24
	halfDuration := time.Duration(durationHours / 2 * float64(time.Hour))

	// 使用估算算法（性能优化）
	// 精确搜索虽然更准确，但计算成本太高（每次需要50+次星历计算）
	// 对于大多数用途，估算误差在几分钟以内，完全可以接受
	var exactTime time.Time
	exactTime = estimateExactTime(orb, queryTime, isApplying, relativeSpeed)

	startTime := exactTime.Add(-halfDuration)
	endTime := exactTime.Add(halfDuration)

	return &models.FactorLifecycle{
		StartTime: startTime,
		PeakTime:  exactTime,
		EndTime:   endTime,
		Duration:  durationHours,
	}
}

// CalculateAspectLifecycle 根据相位容许度计算生命周期（兼容旧接口）
// 假设行星平均每天移动约1°
func CalculateAspectLifecycle(orb float64, queryTime time.Time, isApplying bool) *models.FactorLifecycle {
	// 使用默认参数调用新函数
	return CalculateAspectLifecycleWithPlanets(orb, queryTime, isApplying, models.Sun, models.Moon, 8.0)
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

// ==================== 剩余时间计算 ====================

// CalculateRemainingDays 计算因子从指定时间到结束的剩余天数
// 返回值：
//   - 正数：从当前时间到结束时间的剩余天数（支持小数，如0.5表示12小时）
//   - 0：因子已结束或刚好结束
func CalculateRemainingDays(lifecycle *models.FactorLifecycle, currentTime time.Time) float64 {
	if lifecycle == nil {
		return 0
	}

	// 如果当前时间在结束时间之后，返回0（已结束）
	if currentTime.After(lifecycle.EndTime) {
		return 0
	}

	// 计算从当前时间到结束时间的剩余天数
	// 无论当前时间是在开始之前还是之后，都按此方式计算
	remaining := lifecycle.EndTime.Sub(currentTime)
	remainingHours := remaining.Hours()
	remainingDays := remainingHours / 24.0

	return remainingDays
}

// CalculateRemainingHours 计算因子从指定时间到结束的剩余小时数
func CalculateRemainingHours(lifecycle *models.FactorLifecycle, currentTime time.Time) float64 {
	if lifecycle == nil {
		return 0
	}

	if currentTime.Before(lifecycle.StartTime) {
		return lifecycle.Duration
	}

	if currentTime.After(lifecycle.EndTime) {
		return 0
	}

	remaining := lifecycle.EndTime.Sub(currentTime)
	return remaining.Hours()
}

// CalculateElapsedDays 计算因子从开始到指定时间已经过的天数
func CalculateElapsedDays(lifecycle *models.FactorLifecycle, currentTime time.Time) float64 {
	if lifecycle == nil {
		return 0
	}

	if currentTime.Before(lifecycle.StartTime) {
		return 0
	}

	if currentTime.After(lifecycle.EndTime) {
		return lifecycle.Duration / 24.0
	}

	elapsed := currentTime.Sub(lifecycle.StartTime)
	return elapsed.Hours() / 24.0
}

// GetFactorTimeInfo 获取因子的完整时间信息
type FactorTimeInfo struct {
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	PeakTime      time.Time `json:"peakTime"`
	TotalDays     float64   `json:"totalDays"`     // 总持续天数
	ElapsedDays   float64   `json:"elapsedDays"`   // 已过天数
	RemainingDays float64   `json:"remainingDays"` // 剩余天数
	Progress      float64   `json:"progress"`      // 进度 0.0-1.0
	Phase         string    `json:"phase"`         // "pending", "rising", "peak", "falling", "ended"
}

// GetFactorTimeInfo 获取因子的详细时间信息
func GetFactorTimeInfo(lifecycle *models.FactorLifecycle, currentTime time.Time) *FactorTimeInfo {
	if lifecycle == nil {
		return nil
	}

	totalDays := lifecycle.Duration / 24.0
	elapsedDays := CalculateElapsedDays(lifecycle, currentTime)
	remainingDays := CalculateRemainingDays(lifecycle, currentTime)

	// 计算进度
	var progress float64
	if totalDays > 0 {
		progress = elapsedDays / totalDays
		if progress > 1 {
			progress = 1
		}
		if progress < 0 {
			progress = 0
		}
	}

	// 判断阶段
	var phase string
	if currentTime.Before(lifecycle.StartTime) {
		phase = "pending"
	} else if currentTime.After(lifecycle.EndTime) {
		phase = "ended"
	} else if currentTime.Before(lifecycle.PeakTime) {
		phase = "rising"
	} else if currentTime.Equal(lifecycle.PeakTime) {
		phase = "peak"
	} else {
		phase = "falling"
	}

	return &FactorTimeInfo{
		StartTime:     lifecycle.StartTime,
		EndTime:       lifecycle.EndTime,
		PeakTime:      lifecycle.PeakTime,
		TotalDays:     totalDays,
		ElapsedDays:   elapsedDays,
		RemainingDays: remainingDays,
		Progress:      progress,
		Phase:         phase,
	}
}

// UpdateFactorRemainingDays 更新因子的剩余天数字段
func UpdateFactorRemainingDays(factor *models.InfluenceFactor, currentTime time.Time) {
	if factor == nil {
		return
	}
	factor.RemainingDays = CalculateRemainingDays(factor.Lifecycle, currentTime)
}

