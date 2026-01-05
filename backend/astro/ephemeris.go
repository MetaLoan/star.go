package astro

import (
	"math"
	"star/models"
	"time"
)

// ==================== 儒略日计算 ====================

// DateToJulianDay 将日期转换为儒略日
func DateToJulianDay(t time.Time) float64 {
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	hour := float64(t.Hour()) + float64(t.Minute())/60.0 + float64(t.Second())/3600.0

	if month <= 2 {
		year--
		month += 12
	}

	A := year / 100
	B := 2 - A + A/4

	jd := float64(int(365.25*float64(year+4716))) +
		float64(int(30.6001*float64(month+1))) +
		float64(day) + hour/24.0 + float64(B) - 1524.5

	return jd
}

// JulianDayToDate 将儒略日转换为日期
func JulianDayToDate(jd float64) time.Time {
	jd = jd + 0.5
	Z := int(jd)
	F := jd - float64(Z)

	var A int
	if Z < 2299161 {
		A = Z
	} else {
		alpha := int((float64(Z) - 1867216.25) / 36524.25)
		A = Z + 1 + alpha - alpha/4
	}

	B := A + 1524
	C := int((float64(B) - 122.1) / 365.25)
	D := int(365.25 * float64(C))
	E := int(float64(B-D) / 30.6001)

	day := B - D - int(30.6001*float64(E))

	var month int
	if E < 14 {
		month = E - 1
	} else {
		month = E - 13
	}

	var year int
	if month > 2 {
		year = C - 4716
	} else {
		year = C - 4715
	}

	hour := F * 24
	h := int(hour)
	minute := int((hour - float64(h)) * 60)
	second := int(((hour - float64(h)) * 60 - float64(minute)) * 60)

	return time.Date(year, time.Month(month), day, h, minute, second, 0, time.UTC)
}

// ==================== 角度工具函数 ====================

// NormalizeAngle 将角度归一化到0-360范围
func NormalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}
	return angle
}

// AngleDifference 计算两个角度之间的最短差值
func AngleDifference(a1, a2 float64) float64 {
	diff := math.Abs(a1 - a2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// ==================== 太阳位置计算 ====================

// CalculateSunPosition 计算太阳位置
func CalculateSunPosition(jd float64) models.PlanetPosition {
	T := (jd - J2000) / 36525.0

	// 太阳平黄经
	L0 := NormalizeAngle(280.4664567 + 36000.76983*T + 0.0003032*T*T)

	// 太阳平近点角
	M := NormalizeAngle(357.5291092 + 35999.0502909*T - 0.0001536*T*T)
	Mrad := M * DEG_TO_RAD

	// 中心差
	C := (1.9146 - 0.004817*T - 0.000014*T*T) * math.Sin(Mrad)
	C += (0.019993 - 0.000101*T) * math.Sin(2*Mrad)
	C += 0.00029 * math.Sin(3*Mrad)

	// 真黄经
	sunLong := NormalizeAngle(L0 + C)

	// 获取星座信息
	zodiac := GetZodiacByLongitude(sunLong)
	signDegree := math.Mod(sunLong, 30)

	// 计算尊贵度
	dignity := GetDignity(models.Sun, zodiac.ID)
	dignityScore := GetDignityScore(dignity)

	planetInfo := GetPlanetInfo(models.Sun)

	return models.PlanetPosition{
		ID:           models.Sun,
		Name:         planetInfo.Name,
		Symbol:       planetInfo.Symbol,
		Longitude:    sunLong,
		Latitude:     0,
		Sign:         zodiac.ID,
		SignName:     zodiac.Name,
		SignSymbol:   zodiac.Symbol,
		SignDegree:   signDegree,
		Retrograde:   false, // 太阳不逆行
		DignityScore: dignityScore,
	}
}

// ==================== 月亮位置计算 ====================

// CalculateMoonPosition 计算月亮位置
func CalculateMoonPosition(jd float64) models.PlanetPosition {
	T := (jd - J2000) / 36525.0

	// 月亮平黄经
	Lp := NormalizeAngle(218.3164477 + 481267.88123421*T -
		0.0015786*T*T + T*T*T/538841 - T*T*T*T/65194000)

	// 月亮平近点角
	M := NormalizeAngle(134.9633964 + 477198.8675055*T +
		0.0087414*T*T + T*T*T/69699 - T*T*T*T/14712000)

	// 太阳平近点角
	Ms := NormalizeAngle(357.5291092 + 35999.0502909*T -
		0.0001536*T*T + T*T*T/24490000)

	// 月亮升交点平黄经
	F := NormalizeAngle(93.272095 + 483202.0175233*T -
		0.0036539*T*T - T*T*T/3526000 + T*T*T*T/863310000)

	// 月日角距
	D := NormalizeAngle(297.8501921 + 445267.1114034*T -
		0.0018819*T*T + T*T*T/545868 - T*T*T*T/113065000)

	// 转换为弧度
	Mrad := M * DEG_TO_RAD
	Msrad := Ms * DEG_TO_RAD
	Frad := F * DEG_TO_RAD
	Drad := D * DEG_TO_RAD

	// 黄经摄动项
	longitude := Lp
	longitude += 6.289 * math.Sin(Mrad)                    // 主椭圆
	longitude -= 1.274 * math.Sin(2*Drad-Mrad)             // 出差
	longitude += 0.658 * math.Sin(2*Drad)                  // 变差
	longitude += 0.214 * math.Sin(2*Mrad)                  // 年差
	longitude -= 0.186 * math.Sin(Msrad)                   // 二均差
	longitude -= 0.114 * math.Sin(2*Frad)                  // 章动
	longitude += 0.059 * math.Sin(2*Drad-2*Mrad)
	longitude -= 0.057 * math.Sin(2*Drad-Msrad-Mrad)
	longitude += 0.053 * math.Sin(2*Drad+Mrad)
	longitude += 0.046 * math.Sin(2*Drad-Msrad)
	longitude += 0.041 * math.Sin(Mrad-Msrad)
	longitude -= 0.035 * math.Sin(Drad)
	longitude -= 0.031 * math.Sin(Mrad+Msrad)

	moonLong := NormalizeAngle(longitude)

	// 获取星座信息
	zodiac := GetZodiacByLongitude(moonLong)
	signDegree := math.Mod(moonLong, 30)

	// 计算尊贵度
	dignity := GetDignity(models.Moon, zodiac.ID)
	dignityScore := GetDignityScore(dignity)

	planetInfo := GetPlanetInfo(models.Moon)

	return models.PlanetPosition{
		ID:           models.Moon,
		Name:         planetInfo.Name,
		Symbol:       planetInfo.Symbol,
		Longitude:    moonLong,
		Latitude:     0,
		Sign:         zodiac.ID,
		SignName:     zodiac.Name,
		SignSymbol:   zodiac.Symbol,
		SignDegree:   signDegree,
		Retrograde:   false, // 月亮不逆行
		DignityScore: dignityScore,
	}
}

// ==================== 行星轨道元素 ====================

// OrbitalElements 轨道元素
type OrbitalElements struct {
	A  float64 // 半长轴 (AU)
	E  float64 // 离心率
	I  float64 // 倾角 (度)
	O  float64 // 升交点黄经 (度)
	W  float64 // 近日点幅角 (度)
	L  float64 // 平黄经 (度)
}

// 地球轨道元素结构 (单独定义，因为地球不是占星行星)
var earthOrbitalElements = struct {
	A0, AdT float64
	E0, EdT float64
	I0, IdT float64
	O0, OdT float64
	W0, WdT float64
	L0, LdT float64
}{1.00000261, 0.00000562, 0.01671123, -0.00004392, 0.00001531, -0.01294668, 0.0, 0.0, 102.93768193, 0.32327364, 100.46457166, 35999.37244981}

// 行星轨道元素 (J2000 参考，来自 JPL)
var planetOrbitalElements = map[models.PlanetID]struct {
	A0, AdT  float64
	E0, EdT  float64
	I0, IdT  float64
	O0, OdT  float64
	W0, WdT  float64
	L0, LdT  float64
}{
	models.Mercury: {0.38709927, 0.00000037, 0.20563593, 0.00001906, 7.00497902, -0.00594749, 48.33076593, -0.12534081, 77.45779628, 0.16047689, 252.25032350, 149472.67411175},
	models.Venus:   {0.72333566, 0.00000390, 0.00677672, -0.00004107, 3.39467605, -0.00078890, 76.67984255, -0.27769418, 131.60246718, 0.00268329, 181.97909950, 58517.81538729},
	models.Mars:    {1.52371034, 0.00001847, 0.09339410, 0.00007882, 1.84969142, -0.00813131, 49.55953891, -0.29257343, -23.94362959, 0.44441088, -4.55343205, 19140.30268499},
	models.Jupiter: {5.20288700, -0.00011607, 0.04838624, -0.00013253, 1.30439695, -0.00183714, 100.47390909, 0.20469106, 14.72847983, 0.21252668, 34.39644051, 3034.74612775},
	models.Saturn:  {9.53667594, -0.00125060, 0.05386179, -0.00050991, 2.48599187, 0.00193609, 113.66242448, -0.28867794, 92.59887831, -0.41897216, 49.95424423, 1222.49362201},
	models.Uranus:  {19.18916464, -0.00196176, 0.04725744, -0.00004397, 0.77263783, -0.00242939, 74.01692503, 0.04240589, 170.95427630, 0.40805281, 313.23810451, 428.48202785},
	models.Neptune: {30.06992276, 0.00026291, 0.00859048, 0.00005105, 1.77004347, 0.00035372, 131.78422574, -0.00508664, 44.96476227, -0.32241464, -55.12002969, 218.45945325},
	models.Pluto:   {39.48211675, -0.00031596, 0.24882730, 0.00005170, 17.14001206, 0.00004818, 110.30393684, -0.01183482, 224.06891629, -0.04062942, 238.92903833, 145.20780515},
}

// GetOrbitalElements 获取行星轨道元素
func GetOrbitalElements(planet models.PlanetID, T float64) OrbitalElements {
	elem := planetOrbitalElements[planet]
	return OrbitalElements{
		A: elem.A0 + elem.AdT*T,
		E: elem.E0 + elem.EdT*T,
		I: elem.I0 + elem.IdT*T,
		O: elem.O0 + elem.OdT*T,
		W: elem.W0 + elem.WdT*T,
		L: elem.L0 + elem.LdT*T,
	}
}

// ==================== 行星位置计算 ====================

// CalculatePlanetPosition 计算行星位置（简化算法）
func CalculatePlanetPosition(planet models.PlanetID, jd float64) models.PlanetPosition {
	// Debug print for 1990-06-15 (JD approx 2448057.5)
	if planet == models.Saturn && jd > 2448000 && jd < 2448100 {
		// fmt.Printf("DEBUG: Calculating Saturn for JD=%f\n", jd)
	}

	T := (jd - J2000) / 36525.0

	var longitude float64
	var retrograde bool

	switch planet {
	case models.Sun:
		return CalculateSunPosition(jd)
	case models.Moon:
		return CalculateMoonPosition(jd)
	case models.NorthNode:
		return CalculateNorthNodePosition(jd)
	case models.Chiron:
		return CalculateChironPosition(jd)
	default:
		longitude, retrograde = calculatePlanetLongitude(planet, T, jd)
	}

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
		Latitude:     0,
		Sign:         zodiac.ID,
		SignName:     zodiac.Name,
		SignSymbol:   zodiac.Symbol,
		SignDegree:   signDegree,
		Retrograde:   retrograde,
		DignityScore: dignityScore,
	}
}

// getEarthOrbitalElements 获取地球轨道元素
func getEarthOrbitalElements(T float64) OrbitalElements {
	e := earthOrbitalElements
	return OrbitalElements{
		A: e.A0 + e.AdT*T,
		E: e.E0 + e.EdT*T,
		I: e.I0 + e.IdT*T,
		O: e.O0 + e.OdT*T,
		W: e.W0 + e.WdT*T,
		L: e.L0 + e.LdT*T,
	}
}

// calculateHeliocentricPosition 计算日心坐标
func calculateHeliocentricPosition(elem OrbitalElements) (x, y, z, r float64) {
	// 平近点角
	M := NormalizeAngle(elem.L - elem.W)
	Mrad := M * DEG_TO_RAD

	// 求解开普勒方程
	E := solveKepler(Mrad, elem.E)

	// 计算真近点角
	v := 2 * math.Atan2(
		math.Sqrt(1+elem.E)*math.Sin(E/2),
		math.Sqrt(1-elem.E)*math.Cos(E/2),
	)

	// 日心距离
	r = elem.A * (1 - elem.E*math.Cos(E))

	// 近日点幅角 (弧度)
	// W 是近日点黄经 (Longitude of Perihelion, varpi)
	// O 是升交点黄经 (Longitude of Ascending Node, Omega)
	// w 是近日点幅角 (Argument of Perihelion, omega) = varpi - Omega
	w := elem.W - elem.O
	wRad := w * DEG_TO_RAD
	oRad := elem.O * DEG_TO_RAD
	iRad := elem.I * DEG_TO_RAD

	// 轨道平面坐标
	xOrb := r * math.Cos(v)
	yOrb := r * math.Sin(v)

	// 转换到黄道坐标系
	cosO := math.Cos(oRad)
	sinO := math.Sin(oRad)
	cosI := math.Cos(iRad)
	sinI := math.Sin(iRad)
	cosW := math.Cos(wRad)
	sinW := math.Sin(wRad)

	x = (cosO*cosW - sinO*sinW*cosI)*xOrb + (-cosO*sinW - sinO*cosW*cosI)*yOrb
	y = (sinO*cosW + cosO*sinW*cosI)*xOrb + (-sinO*sinW + cosO*cosW*cosI)*yOrb
	z = (sinW*sinI)*xOrb + (cosW*sinI)*yOrb

	return x, y, z, r
}

// calculatePlanetLongitude 计算行星黄经（改进版）
func calculatePlanetLongitude(planet models.PlanetID, T, jd float64) (float64, bool) {
	// 获取行星轨道元素
	elem := GetOrbitalElements(planet, T)

	// 计算行星日心坐标
	xP, yP, zP, rP := calculateHeliocentricPosition(elem)

	// 获取地球轨道元素
	earthElem := getEarthOrbitalElements(T)

	// 计算地球日心坐标
	xE, yE, zE, _ := calculateHeliocentricPosition(earthElem)

	// 计算地心坐标
	xGeo := xP - xE
	yGeo := yP - yE
	zGeo := zP - zE

	// 计算黄经
	geoLon := NormalizeAngle(math.Atan2(yGeo, xGeo) * RAD_TO_DEG)

	// 判断逆行 - 计算前后位置比较
	prevJd := jd - 1
	nextJd := jd + 1

	prevLon := calculateGeocentricLongitude(planet, (prevJd-J2000)/36525.0)
	nextLon := calculateGeocentricLongitude(planet, (nextJd-J2000)/36525.0)

	// 计算运动方向
	motion := nextLon - prevLon
	if motion < -180 {
		motion += 360
	}
	if motion > 180 {
		motion -= 360
	}
	retrograde := motion < 0

	// 忽略未使用变量警告
	_ = zP
	_ = zE
	_ = rP
	_ = zGeo

	return geoLon, retrograde
}

// calculateGeocentricLongitude 简化的地心黄经计算（用于逆行判断）
func calculateGeocentricLongitude(planet models.PlanetID, T float64) float64 {
	elem := GetOrbitalElements(planet, T)
	earthElem := getEarthOrbitalElements(T)

	xP, yP, _, _ := calculateHeliocentricPosition(elem)
	xE, yE, _, _ := calculateHeliocentricPosition(earthElem)

	xGeo := xP - xE
	yGeo := yP - yE

	return NormalizeAngle(math.Atan2(yGeo, xGeo) * RAD_TO_DEG)
}

// solveKepler 求解开普勒方程
func solveKepler(M, e float64) float64 {
	E := M
	for i := 0; i < 10; i++ {
		dE := (E - e*math.Sin(E) - M) / (1 - e*math.Cos(E))
		E = E - dE
		if math.Abs(dE) < 1e-9 {
			break
		}
	}
	return E
}

// ==================== 北交点计算 ====================

// CalculateNorthNodePosition 计算北交点位置
func CalculateNorthNodePosition(jd float64) models.PlanetPosition {
	T := (jd - J2000) / 36525.0

	// 北交点平黄经（逆行）
	longitude := NormalizeAngle(125.0445479 - 1934.1362891*T + 0.0020754*T*T)

	zodiac := GetZodiacByLongitude(longitude)
	signDegree := math.Mod(longitude, 30)
	planetInfo := GetPlanetInfo(models.NorthNode)

	return models.PlanetPosition{
		ID:         models.NorthNode,
		Name:       planetInfo.Name,
		Symbol:     planetInfo.Symbol,
		Longitude:  longitude,
		Latitude:   0,
		Sign:       zodiac.ID,
		SignName:   zodiac.Name,
		SignSymbol: zodiac.Symbol,
		SignDegree: signDegree,
		Retrograde: true, // 北交点总是逆行
	}
}

// ==================== 凯龙星计算 ====================

// CalculateChironPosition 计算凯龙星位置（简化）
func CalculateChironPosition(jd float64) models.PlanetPosition {
	T := (jd - J2000) / 36525.0
	daysSinceJ2000 := jd - J2000

	// 凯龙星轨道元素（简化）
	a := 13.6481 // 半长轴 (AU)
	e := 0.3832  // 离心率

	// 平近点角
	M := NormalizeAngle(92.43 + 0.01958*daysSinceJ2000)
	Mrad := M * DEG_TO_RAD

	// 求解开普勒方程
	E := solveKepler(Mrad, e)

	// 计算真近点角
	v := 2 * math.Atan2(
		math.Sqrt(1+e)*math.Sin(E/2),
		math.Sqrt(1-e)*math.Cos(E/2),
	)

	// 计算黄经
	w := 339.5 // 近日点幅角
	longitude := NormalizeAngle((v * RAD_TO_DEG) + w)

	// 简化地心转换
	earthLon := CalculateSunPosition(jd).Longitude + 180
	geoLon := longitude + (180/a)*math.Sin((longitude-earthLon)*DEG_TO_RAD)
	geoLon = NormalizeAngle(geoLon)

	zodiac := GetZodiacByLongitude(geoLon)
	signDegree := math.Mod(geoLon, 30)
	planetInfo := GetPlanetInfo(models.Chiron)

	// 判断逆行
	prevLon := NormalizeAngle(92.43 + 0.01958*(daysSinceJ2000-1))
	nextLon := NormalizeAngle(92.43 + 0.01958*(daysSinceJ2000+1))
	motion := nextLon - prevLon
	if motion < -180 {
		motion += 360
	}
	if motion > 180 {
		motion -= 360
	}
	retrograde := motion < 0

	_ = T // 避免未使用警告

	return models.PlanetPosition{
		ID:         models.Chiron,
		Name:       planetInfo.Name,
		Symbol:     planetInfo.Symbol,
		Longitude:  geoLon,
		Latitude:   0,
		Sign:       zodiac.ID,
		SignName:   zodiac.Name,
		SignSymbol: zodiac.Symbol,
		SignDegree: signDegree,
		Retrograde: retrograde,
	}
}

// ==================== 获取所有行星位置 ====================

// GetAllPlanetPositions 获取所有行星位置
func GetAllPlanetPositions(jd float64) []models.PlanetPosition {
	planets := []models.PlanetID{
		models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars,
		models.Jupiter, models.Saturn, models.Uranus, models.Neptune, models.Pluto,
		models.NorthNode, models.Chiron,
	}

	positions := make([]models.PlanetPosition, 0, len(planets))
	for _, p := range planets {
		positions = append(positions, CalculatePlanetPosition(p, jd))
	}
	return positions
}

// GetTransitPositions 获取当前行运行星位置
// 使用 Swiss Ephemeris 作为唯一数据源
func GetTransitPositions(date time.Time) []models.PlanetPosition {
	jd := DateToJulianDay(date)
	return GetPlanetPositionsUnified(jd)
}

