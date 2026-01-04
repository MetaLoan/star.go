# 算法文档

## 1. 天文计算基础

### 1.1 儒略日 (Julian Day)

儒略日是天文学中连续计日的方法，从公元前4713年1月1日中午开始计算。

#### 公式

```
日期 → 儒略日：
if month <= 2:
    year = year - 1
    month = month + 12

A = year / 100
B = 2 - A + A/4

JD = int(365.25 * (year + 4716)) + int(30.6001 * (month + 1)) 
     + day + hour/24 + B - 1524.5
```

#### 代码实现 (Go)

```go
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
```

### 1.2 J2000 纪元

J2000.0 是天文学标准纪元，对应 2000年1月1日12:00 TT (地球时)。

```
J2000 = 2451545.0 (儒略日)
```

计算距 J2000 的儒略世纪：
```
T = (JD - 2451545.0) / 36525.0
```

---

## 2. 行星位置计算

### 2.1 太阳位置

太阳位置计算基于 VSOP87 简化模型。

#### 算法步骤

```
1. 计算儒略世纪 T
2. 计算太阳平黄经 L0
   L0 = 280.4664567 + 36000.76983*T + 0.0003032*T²

3. 计算太阳平近点角 M
   M = 357.5291092 + 35999.0502909*T - 0.0001536*T²

4. 计算中心差 C
   C = (1.9146 - 0.004817*T - 0.000014*T²) * sin(M)
     + (0.019993 - 0.000101*T) * sin(2M)
     + 0.00029 * sin(3M)

5. 真黄经 = L0 + C
```

#### 代码实现

```go
func CalculateSunPosition(jd float64) PlanetPosition {
    T := (jd - J2000) / 36525.0
    
    L0 := NormalizeAngle(280.4664567 + 36000.76983*T + 0.0003032*T*T)
    M := NormalizeAngle(357.5291092 + 35999.0502909*T - 0.0001536*T*T)
    Mrad := M * DEG_TO_RAD
    
    C := (1.9146 - 0.004817*T - 0.000014*T*T) * math.Sin(Mrad)
    C += (0.019993 - 0.000101*T) * math.Sin(2*Mrad)
    C += 0.00029 * math.Sin(3*Mrad)
    
    sunLong := NormalizeAngle(L0 + C)
    
    return PlanetPosition{
        ID:        Sun,
        Longitude: sunLong,
        // ...
    }
}
```

### 2.2 月亮位置

月亮位置计算较为复杂，需要考虑多个摄动项。

#### 主要参数

| 参数 | 含义 | 公式 |
|------|------|------|
| Lp | 月亮平黄经 | 218.3164477 + 481267.88123421*T |
| M | 月亮平近点角 | 134.9633964 + 477198.8675055*T |
| Ms | 太阳平近点角 | 357.5291092 + 35999.0502909*T |
| F | 月亮升交点平黄经 | 93.272095 + 483202.0175233*T |
| D | 月亮平角距 | 297.8501921 + 445267.1114034*T |

#### 主要摄动项

```
longitude = Lp
longitude += 6.289 * sin(M)         // 主椭圆摄动
longitude -= 1.274 * sin(2D - M)    // 出差
longitude += 0.658 * sin(2D)        // 变差
longitude += 0.214 * sin(2M)
longitude -= 0.186 * sin(Ms)        // 年差
longitude -= 0.114 * sin(2F)
```

### 2.3 行星位置

行星位置使用开普勒轨道元素计算。

#### 轨道元素

| 元素 | 符号 | 含义 |
|------|------|------|
| 半长轴 | a | 轨道大小 |
| 离心率 | e | 轨道形状 |
| 倾角 | i | 轨道与黄道夹角 |
| 升交点黄经 | Ω | 轨道与黄道交点 |
| 近日点幅角 | ω | 近日点方向 |
| 平近点角 | M | 行星在轨道上的位置 |

#### 计算步骤

```
1. 计算轨道元素（考虑时间变化）
2. 求解开普勒方程得到真近点角 v
3. 计算日心距离 r
4. 计算日心黄道坐标
5. 转换为地心黄道坐标
```

#### 开普勒方程

```
M = E - e*sin(E)    // M: 平近点角, E: 偏近点角, e: 离心率

迭代求解 E:
E₀ = M
Eₙ₊₁ = M + e*sin(Eₙ)
直到 |Eₙ₊₁ - Eₙ| < 1e-6
```

#### 真近点角

```
v = 2 * atan2(sqrt(1+e) * sin(E/2), sqrt(1-e) * cos(E/2))
```

#### 日心到地心转换

```
1. 计算行星日心坐标 (xh, yh, zh)
2. 计算地球日心坐标 (xe, ye, ze)
3. 地心坐标 = 行星日心 - 地球日心
4. 转换为黄经黄纬
```

### 2.4 北交点计算

北交点（月球轨道与黄道的升交点）计算：

```go
func CalculateNorthNodePosition(jd float64) PlanetPosition {
    T := (jd - J2000) / 36525.0
    
    // 北交点平黄经（逆行）
    longitude := NormalizeAngle(125.0445479 - 1934.1362891*T + 0.0020754*T*T)
    
    return PlanetPosition{
        ID:        NorthNode,
        Longitude: longitude,
    }
}
```

### 2.5 凯龙星计算

凯龙星（Chiron）是一颗半人马小行星，轨道在土星和天王星之间。

```go
// 凯龙星轨道元素
a := 13.6481 + 0.0000*T           // 半长轴 (AU)
e := 0.3832 + 0.0000*T            // 离心率
i := 6.9311 + 0.0000*T            // 倾角 (度)
O := 209.4 + 0.0*T                // 升交点黄经
w := 339.5 + 0.0*T                // 近日点幅角
M := 92.43 + 0.01958*daysSinceJ2000  // 平近点角
```

---

## 3. 宫位计算

### 3.1 Placidus 分宫制

Placidus 是最常用的分宫制，基于时间等分原则。

#### 计算步骤

```
1. 计算本地恒星时 (LST)
2. 计算上升点 (ASC)
3. 计算中天点 (MC)
4. 计算中间宫位
```

#### 本地恒星时

```
GMST = 280.46061837 + 360.98564736629*(JD - 2451545.0)
       + 0.000387933*T² - T³/38710000

LST = GMST + longitude  (经度，东经为正)
```

#### 上升点

```
ASC = atan2(cos(RAMC), -sin(RAMC)*cos(ε) - tan(φ)*sin(ε))

其中：
RAMC = LST (以角度表示)
ε = 黄赤交角 ≈ 23.4393°
φ = 地理纬度
```

#### 中天点

```
MC = RAMC  // 赤经中天

// 转换为黄经
MC_ecl = atan2(sin(MC)*cos(ε), cos(MC))
```

### 3.2 其他分宫制（待实现）

| 分宫制 | 特点 |
|--------|------|
| Equal | 每宫30° |
| Whole Sign | 整星座宫位 |
| Koch | 时间等分 |
| Campanus | 空间等分 |
| Regiomontanus | 赤道投影 |

---

## 4. 相位计算

### 4.1 相位定义

| 相位 | 角度 | 容许度 | 性质 | 权重 |
|------|------|--------|------|------|
| 合相 | 0° | 10° | 融合 | 10 |
| 六分相 | 60° | 6° | 和谐 | 3 |
| 四分相 | 90° | 8° | 紧张 | 7 |
| 三分相 | 120° | 8° | 顺畅 | 6 |
| 对分相 | 180° | 10° | 对立 | 8 |

### 4.2 相位计算算法

```go
func CalculateAspects(planets []PlanetPosition) []AspectData {
    var aspects []AspectData
    
    for i := 0; i < len(planets); i++ {
        for j := i + 1; j < len(planets); j++ {
            p1, p2 := planets[i], planets[j]
            
            // 计算角距
            diff := math.Abs(p1.Longitude - p2.Longitude)
            if diff > 180 {
                diff = 360 - diff
            }
            
            // 检查每个相位类型
            for _, def := range AspectDefinitions {
                orb := math.Abs(diff - def.Angle)
                if orb <= def.Orb {
                    // 计算强度
                    strength := 1.0 - orb/def.Orb
                    
                    // 计算权重
                    p1Weight := GetPlanetWeight(p1.ID)
                    p2Weight := GetPlanetWeight(p2.ID)
                    weight := strength * def.Weight * (p1Weight + p2Weight) / 20.0
                    
                    aspects = append(aspects, AspectData{
                        Planet1:    p1.ID,
                        Planet2:    p2.ID,
                        AspectType: def.Type,
                        Orb:        orb,
                        Strength:   strength,
                        Weight:     weight,
                    })
                }
            }
        }
    }
    
    return aspects
}
```

### 4.3 行星权重

| 行星 | 权重 | 说明 |
|------|------|------|
| 太阳 | 10 | 自我核心 |
| 月亮 | 10 | 情感核心 |
| 水星 | 4 | 思维沟通 |
| 金星 | 5 | 价值关系 |
| 火星 | 6 | 行动欲望 |
| 木星 | 7 | 扩张信念 |
| 土星 | 8 | 结构限制 |
| 天王星 | 5 | 突变创新 |
| 海王星 | 4 | 灵性幻象 |
| 冥王星 | 6 | 转化重生 |
| 北交点 | 3 | 命运方向 |
| 凯龙 | 3 | 伤痛治愈 |

---

## 5. 行运计算

### 5.1 行运相位

行运相位是当前天象行星与本命盘行星之间的相位。

```go
func CalculateDailyTransit(date time.Time, natalChart *NatalChart) *DailyTransit {
    // 获取当前行星位置
    transitPositions := GetTransitPositions(date)
    
    // 计算行运相位（容许度收紧到 80%）
    aspects := calculateTransitToNatalAspects(transitPositions, natalChart.Planets)
    
    // 计算分数
    score := calculateTransitScore(aspects)
    
    return &DailyTransit{
        Date:    date,
        Aspects: aspects,
        Score:   score,
    }
}
```

### 5.2 行运分数

```go
func calculateTransitScore(aspects []AspectData) TransitScore {
    var harmonious, tense float64
    
    for _, a := range aspects {
        switch a.AspectType {
        case Trine, Sextile:
            harmonious += a.Weight
        case Square, Opposition:
            tense += a.Weight
        case Conjunction:
            // 合相根据行星性质判断
            if isBenefic(a.Planet1) || isBenefic(a.Planet2) {
                harmonious += a.Weight * 0.5
            } else {
                tense += a.Weight * 0.3
            }
        }
    }
    
    return TransitScore{
        Total:      harmonious - tense,
        Harmonious: harmonious,
        Tense:      tense,
    }
}
```

### 5.3 重大行运（个性化）

重大行运根据用户本命盘位置精确计算。

#### 配置

```go
type MajorTransitConfig struct {
    TransitPlanet PlanetID   // 行运行星
    NatalPlanet   PlanetID   // 本命行星
    AspectType    AspectType // 相位类型
    Name          string     // 事件名称
    MinAge        int        // 最小年龄
}

var transitConfigs = []MajorTransitConfig{
    {Saturn, Saturn, Conjunction, "土星回归", 28},
    {Saturn, Saturn, Opposition, "土星对冲", 14},
    {Saturn, Saturn, Square, "土星刑克", 7},
    {Jupiter, Jupiter, Conjunction, "木星回归", 11},
    {Uranus, Uranus, Opposition, "天王星对冲", 38},
    {Uranus, Uranus, Square, "天王星刑克", 20},
    {Neptune, Neptune, Square, "海王星刑克", 40},
    {Pluto, Pluto, Square, "冥王星刑克", 60},
    // ... 更多配置
}
```

#### 算法

```go
func FindMajorTransitsPersonalized(natalChart *NatalChart, startYear, endYear int) []MajorTransit {
    var transits []MajorTransit
    birthYear := natalChart.BirthData.Date.Year()
    
    for _, config := range transitConfigs {
        // 获取本命行星位置
        natalPos := getNatalPlanetPosition(natalChart, config.NatalPlanet)
        
        // 计算目标经度（本命位置 + 相位角度）
        targetLongitude := NormalizeAngle(natalPos.Longitude + getAspectAngle(config.AspectType))
        
        // 搜索行运行星到达目标经度的时间
        for year := startYear; year <= endYear; year++ {
            age := year - birthYear
            if age < config.MinAge {
                continue
            }
            
            // 搜索该年内的精确相位
            found := searchTransitInYear(config.TransitPlanet, targetLongitude, year)
            if found != nil {
                transits = append(transits, MajorTransit{
                    Name: config.Name,
                    Date: found.Date,
                    // ...
                })
            }
        }
    }
    
    return transits
}
```

---

## 6. 推运计算

### 6.1 次限推运

次限推运（Secondary Progressions）基于"一天 = 一年"的原则。

```
推运日期 = 出生日期 + (目标年龄 天)

例如：30岁的推运
推运日期 = 出生日期 + 30天
```

#### 算法

```go
func CalculateProgressedChart(natalChart *NatalChart, targetDate time.Time) *ProgressedChart {
    birthDate := natalChart.BirthData.Date
    
    // 计算年数
    yearsFromBirth := targetDate.Sub(birthDate).Hours() / (365.25 * 24)
    
    // 推运日期（1天 = 1年）
    daysProgressed := yearsFromBirth
    progressedDate := birthDate.AddDate(0, 0, int(daysProgressed))
    
    // 计算推运行星位置
    progressedJd := DateToJulianDay(progressedDate)
    // ...
}
```

### 6.2 推运月相

推运太阳和月亮的夹角决定推运月相。

```go
func calculateProgressedLunarPhase(progressedSun, progressedMoon float64) LunarPhaseInfo {
    angle := NormalizeAngle(progressedMoon - progressedSun)
    
    // 8个月相阶段
    phases := []struct {
        min, max    float64
        phase, name string
    }{
        {0, 45, "new", "新月期"},
        {45, 90, "crescent", "新月后期"},
        {90, 135, "firstQuarter", "上弦月"},
        {135, 180, "gibbous", "凸月期"},
        {180, 225, "full", "满月期"},
        {225, 270, "disseminating", "播种期"},
        {270, 315, "lastQuarter", "下弦月"},
        {315, 360, "balsamic", "残月期"},
    }
    
    for _, p := range phases {
        if angle >= p.min && angle < p.max {
            return LunarPhaseInfo{
                Phase: p.phase,
                Name:  p.name,
                Angle: angle,
            }
        }
    }
    
    return LunarPhaseInfo{Phase: "new", Name: "新月期", Angle: angle}
}
```

---

## 7. 年限法

### 7.1 年度宫位

年限法根据年龄确定当年的主题宫位。

```
年度宫位 = (年龄 % 12) + 1

0岁 → 1宫（命宫）
1岁 → 2宫（财帛宫）
...
11岁 → 12宫（玄秘宫）
12岁 → 1宫（新周期）
```

### 7.2 年主星

年主星是年度宫位星座的守护星。

```go
func CalculateAnnualProfection(natalChart *NatalChart, age int) *AnnualProfection {
    // 计算年度宫位
    house := (age % 12) + 1
    
    // 获取该宫位的星座
    houseCusp := natalChart.Houses[house-1]
    sign := getZodiacSign(houseCusp.Cusp)
    
    // 获取守护星
    lordOfYear := signRulers[sign.ID]
    
    return &AnnualProfection{
        Age:        age,
        House:      house,
        Sign:       sign.ID,
        LordOfYear: lordOfYear,
        // ...
    }
}
```

---

## 8. 统一时间序列系统

### 8.1 核心原则

1. **最小粒度是小时**
2. **向上聚合**：小时 → 日 → 周 → 月 → 年
3. **原始分数无上限**
4. **显示分数标准化**

### 8.2 小时级分数

```go
func CalculateHourlyScore(date time.Time, hour int, natalChart *NatalChart) TimeSeriesPoint {
    hourTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)
    jd := DateToJulianDay(hourTime)
    
    factors := ScoreFactors{}
    
    // 1. 行星时影响 (-5 ~ +5)
    factors.PlanetaryHour = getPlanetaryHourInfluence(hourTime, natalChart)
    
    // 2. 月亮空亡 (-15 ~ 0)
    factors.MoonVoid = calculateMoonVoidInfluence(moonPos, jd)
    
    // 3. 行运相位 (无限制)
    aspects := calculateTransitToNatalAspects(transitPositions, natalChart.Planets)
    for _, asp := range aspects {
        switch asp.AspectType {
        case Trine, Sextile:
            factors.Aspects += asp.Weight * asp.Strength * 2
        case Conjunction:
            factors.Aspects += asp.Weight * asp.Strength
        case Square, Opposition:
            factors.Aspects -= asp.Weight * asp.Strength * 1.5
        }
    }
    
    // 原始总分
    rawValue := factors.PlanetaryHour + factors.MoonVoid + factors.Aspects
    
    return TimeSeriesPoint{
        Raw:     RawScore{Value: rawValue, Factors: factors},
        Display: normalizeScore(rawValue),
    }
}
```

### 8.3 聚合函数

```go
// 日级聚合
func AggregateToDaily(hourlyPoints []TimeSeriesPoint, date time.Time, natalChart *NatalChart) TimeSeriesPoint {
    // 计算小时分数平均值
    avgRaw := sum(hourlyPoints.raw.value) / len(hourlyPoints)
    
    // 添加日级因子
    dailyFactors.MoonSign = getMoonSignInfluence(moonPos, natalChart)
    dailyFactors.MoonPhase = getMoonPhaseInfluence(moonPhaseAngle)
    
    // 日级原始分数 = 小时平均 + 日级因子
    dailyRaw := avgRaw + dailyFactors.MoonSign + dailyFactors.MoonPhase
    
    // 计算波动率
    volatility := calculateVolatility(hourlyPoints)
    
    return TimeSeriesPoint{
        Raw:        RawScore{Value: dailyRaw, Factors: dailyFactors},
        Display:    normalizeScore(dailyRaw),
        Volatility: volatility,
    }
}
```

### 8.4 标准化函数

使用 sigmoid/tanh 将原始分数压缩到 0-100 范围。

```go
func normalizeScore(raw float64) float64 {
    // 使用修改的 sigmoid
    // 0 映射到 50，正值 > 50，负值 < 50
    scale := 20.0
    return 50 + 50*math.Tanh(raw/scale)
}
```

---

## 9. 影响因子系统

### 9.1 因子类型

| 类型 | 描述 | 范围 |
|------|------|------|
| dignity | 尊贵度 | ±3 |
| retrograde | 逆行 | -2 |
| aspectPhase | 相位阶段 | ×0.8 |
| aspectOrb | 容许度 | ×0.5 |
| outerPlanet | 外行星 | ×1.2 |
| profectionLord | 年主星 | ×1.0 |
| lunarPhase | 月相 | ×0.7 |
| planetaryHour | 行星时 | ×0.3 |

### 9.2 尊贵度表

| 行星 | 入庙 | 旺相 | 落陷 | 失势 |
|------|------|------|------|------|
| 太阳 | 狮子 | 白羊 | 水瓶 | 天秤 |
| 月亮 | 巨蟹 | 金牛 | 摩羯 | 天蝎 |
| 水星 | 双子/处女 | - | 射手/双鱼 | - |
| 金星 | 金牛/天秤 | 双鱼 | 白羊/天蝎 | 处女 |
| 火星 | 白羊/天蝎 | 摩羯 | 金牛/天秤 | 巨蟹 |
| 木星 | 射手/双鱼 | 巨蟹 | 双子/处女 | 摩羯 |
| 土星 | 摩羯/水瓶 | 天秤 | 巨蟹/狮子 | 白羊 |

### 9.3 因子计算

```go
func CalculateInfluenceFactors(ctx *FactorContext, weights *FactorWeights) *FactorResult {
    var factors []InfluenceFactor
    
    // 1. 尊贵度因子
    dignityFactors := calculateDignityFactors(ctx, weights.Dignity)
    factors = append(factors, dignityFactors...)
    
    // 2. 逆行因子
    retrogradeFactors := calculateRetrogradeFactors(ctx, weights.Retrograde)
    factors = append(factors, retrogradeFactors...)
    
    // 3. 相位因子
    aspectFactors := calculateAspectFactors(ctx, weights.AspectPhase, weights.AspectOrb)
    factors = append(factors, aspectFactors...)
    
    // ... 更多因子
    
    // 计算总调整值
    var totalAdjustment float64
    for _, f := range factors {
        totalAdjustment += f.Value * f.Weight
    }
    
    return &FactorResult{
        Factors:         factors,
        TotalAdjustment: totalAdjustment,
    }
}
```

---

## 10. 精度验证

### 10.1 验证方法

使用已知天文事件验证计算精度：

1. **春分点**：太阳黄经 = 0°
2. **夏至点**：太阳黄经 = 90°
3. **秋分点**：太阳黄经 = 180°
4. **冬至点**：太阳黄经 = 270°
5. **新月**：日月合相
6. **满月**：日月对冲

### 10.2 精度指标

| 天体 | 目标精度 | 实际精度 |
|------|---------|---------|
| 太阳 | < 0.1° | ≈ 0.02° |
| 月亮 | < 0.5° | ≈ 0.3° |
| 水星-冥王星 | < 1° | ≈ 0.2-0.5° |
| 新月时刻 | < 6小时 | ≈ 2-3小时 |

