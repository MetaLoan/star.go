# 数据结构文档

## 1. 基础类型

### 1.1 行星 ID (PlanetID)

```go
type PlanetID string

const (
    Sun       PlanetID = "sun"
    Moon      PlanetID = "moon"
    Mercury   PlanetID = "mercury"
    Venus     PlanetID = "venus"
    Mars      PlanetID = "mars"
    Jupiter   PlanetID = "jupiter"
    Saturn    PlanetID = "saturn"
    Uranus    PlanetID = "uranus"
    Neptune   PlanetID = "neptune"
    Pluto     PlanetID = "pluto"
    NorthNode PlanetID = "northNode"
    Chiron    PlanetID = "chiron"
)
```

### 1.2 星座 ID (ZodiacID)

```go
type ZodiacID string

const (
    Aries       ZodiacID = "aries"
    Taurus      ZodiacID = "taurus"
    Gemini      ZodiacID = "gemini"
    Cancer      ZodiacID = "cancer"
    Leo         ZodiacID = "leo"
    Virgo       ZodiacID = "virgo"
    Libra       ZodiacID = "libra"
    Scorpio     ZodiacID = "scorpio"
    Sagittarius ZodiacID = "sagittarius"
    Capricorn   ZodiacID = "capricorn"
    Aquarius    ZodiacID = "aquarius"
    Pisces      ZodiacID = "pisces"
)
```

### 1.3 相位类型 (AspectType)

```go
type AspectType string

const (
    Conjunction AspectType = "conjunction" // 合相 0°
    Sextile     AspectType = "sextile"     // 六分相 60°
    Square      AspectType = "square"      // 四分相 90°
    Trine       AspectType = "trine"       // 三分相 120°
    Opposition  AspectType = "opposition"  // 对分相 180°
)
```

### 1.4 时间粒度 (TimeGranularity)

```go
type TimeGranularity string

const (
    GranularityHour  TimeGranularity = "hour"
    GranularityDay   TimeGranularity = "day"
    GranularityWeek  TimeGranularity = "week"
    GranularityMonth TimeGranularity = "month"
    GranularityYear  TimeGranularity = "year"
)
```

### 1.5 尊贵度 (Dignity)

```go
type Dignity string

const (
    DignityDomicile   Dignity = "domicile"   // 入庙
    DignityExaltation Dignity = "exaltation" // 旺相
    DignityDetriment  Dignity = "detriment"  // 落陷
    DignityFall       Dignity = "fall"       // 失势
    DignityPeregrine  Dignity = "peregrine"  // 游离
)
```

---

## 2. 核心数据结构

### 2.1 出生数据 (BirthData)

```go
type BirthData struct {
    Name      string    `json:"name"`
    Date      time.Time `json:"date"`
    Latitude  float64   `json:"latitude"`
    Longitude float64   `json:"longitude"`
    Timezone  string    `json:"timezone"`
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 用户姓名（可选） |
| Date | time.Time | 出生日期时间（UTC） |
| Latitude | float64 | 纬度 (-90 ~ 90) |
| Longitude | float64 | 经度 (-180 ~ 180) |
| Timezone | string | 时区标识符 (IANA) |

### 2.2 行星位置 (PlanetPosition)

```go
type PlanetPosition struct {
    ID           PlanetID `json:"id"`
    Name         string   `json:"name"`
    Symbol       string   `json:"symbol"`
    Longitude    float64  `json:"longitude"`
    Latitude     float64  `json:"latitude"`
    Sign         ZodiacID `json:"sign"`
    SignName     string   `json:"signName"`
    SignSymbol   string   `json:"signSymbol"`
    SignDegree   float64  `json:"signDegree"`
    Retrograde   bool     `json:"retrograde"`
    House        int      `json:"house"`
    DignityScore float64  `json:"dignityScore"`
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | PlanetID | 行星标识符 |
| Name | string | 中文名称 |
| Symbol | string | Unicode 符号 |
| Longitude | float64 | 黄道经度 (0-360) |
| Latitude | float64 | 黄道纬度 |
| Sign | ZodiacID | 所在星座 |
| SignName | string | 星座中文名 |
| SignSymbol | string | 星座符号 |
| SignDegree | float64 | 在星座内度数 (0-30) |
| Retrograde | bool | 是否逆行 |
| House | int | 所在宫位 (1-12) |
| DignityScore | float64 | 尊贵度分数 |

### 2.3 宫位 (HouseCusp)

```go
type HouseCusp struct {
    House      int      `json:"house"`
    Cusp       float64  `json:"cusp"`
    Sign       ZodiacID `json:"sign"`
    SignName   string   `json:"signName"`
}
```

### 2.4 相位数据 (AspectData)

```go
type AspectData struct {
    Planet1        PlanetID   `json:"planet1"`
    Planet2        PlanetID   `json:"planet2"`
    AspectType     AspectType `json:"aspectType"`
    ExactAngle     float64    `json:"exactAngle"`
    ActualAngle    float64    `json:"actualAngle"`
    Orb            float64    `json:"orb"`
    Applying       bool       `json:"applying"`
    Strength       float64    `json:"strength"`
    Weight         float64    `json:"weight"`
    Interpretation string     `json:"interpretation,omitempty"`
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Planet1 | PlanetID | 第一个行星 |
| Planet2 | PlanetID | 第二个行星 |
| AspectType | AspectType | 相位类型 |
| ExactAngle | float64 | 精确相位角度 |
| ActualAngle | float64 | 实际角度 |
| Orb | float64 | 容许度 |
| Applying | bool | 是否入相（接近中） |
| Strength | float64 | 强度 (0-1) |
| Weight | float64 | 权重（考虑行星重要性） |
| Interpretation | string | 解读文本 |

---

## 3. 星盘结构

### 3.1 本命盘 (NatalChart)

```go
type NatalChart struct {
    BirthData       BirthData         `json:"birthData"`
    Planets         []PlanetPosition  `json:"planets"`
    Houses          []HouseCusp       `json:"houses"`
    Ascendant       float64           `json:"ascendant"`
    Midheaven       float64           `json:"midheaven"`
    Aspects         []AspectData      `json:"aspects"`
    Patterns        []string          `json:"patterns"`
    ElementBalance  map[string]float64 `json:"elementBalance"`
    ModalityBalance map[string]float64 `json:"modalityBalance"`
    DominantPlanets []PlanetID        `json:"dominantPlanets"`
    ChartRuler      PlanetID          `json:"chartRuler"`
}
```

### 3.2 推运盘 (ProgressedChart)

```go
type ProgressedChart struct {
    TargetDate          time.Time          `json:"targetDate"`
    ProgressedDate      time.Time          `json:"progressedDate"`
    DaysProgressed      float64            `json:"daysProgressed"`
    YearsFromBirth      float64            `json:"yearsFromBirth"`
    Planets             []ProgressedPlanet `json:"planets"`
    ProgressedAscendant float64            `json:"progressedAscendant"`
    ProgressedMidheaven float64            `json:"progressedMidheaven"`
    Aspects             []AspectData       `json:"aspects"`
    LunarPhase          LunarPhaseInfo     `json:"lunarPhase"`
}
```

### 3.3 推运行星 (ProgressedPlanet)

```go
type ProgressedPlanet struct {
    ID                  PlanetID `json:"id"`
    Name                string   `json:"name"`
    Symbol              string   `json:"symbol"`
    NatalLongitude      float64  `json:"natalLongitude"`
    ProgressedLongitude float64  `json:"progressedLongitude"`
    Movement            float64  `json:"movement"`
    Sign                ZodiacID `json:"sign"`
    SignName            string   `json:"signName"`
    SignChanged         bool     `json:"signChanged"`
    House               int      `json:"house"`
    HouseChanged        bool     `json:"houseChanged"`
}
```

---

## 4. 预测结构

### 4.1 每日预测 (DailyForecast)

```go
type DailyForecast struct {
    Date           time.Time          `json:"date"`
    DayOfWeek      string             `json:"dayOfWeek"`
    OverallScore   float64            `json:"overallScore"`
    OverallTheme   string             `json:"overallTheme"`
    Dimensions     DailyDimensions    `json:"dimensions"`
    MoonPhase      MoonPhase          `json:"moonPhase"`
    MoonSign       MoonSignInfo       `json:"moonSign"`
    HourlyBreakdown []HourlyForecast  `json:"hourlyBreakdown"`
    ActiveAspects  []AspectData       `json:"activeAspects"`
}
```

### 4.2 每日维度 (DailyDimensions)

```go
type DailyDimensions struct {
    Career       float64 `json:"career"`
    Relationship float64 `json:"relationship"`
    Health       float64 `json:"health"`
    Finance      float64 `json:"finance"`
    Spiritual    float64 `json:"spiritual"`
}
```

### 4.3 月相 (MoonPhase)

```go
type MoonPhase struct {
    Phase        string     `json:"phase"`
    Name         string     `json:"name"`
    Angle        float64    `json:"angle"`
    Illumination float64    `json:"illumination"`
    VoidStart    *time.Time `json:"voidStart,omitempty"`
    VoidEnd      *time.Time `json:"voidEnd,omitempty"`
}
```

### 4.4 小时预测 (HourlyForecast)

```go
type HourlyForecast struct {
    Hour          int      `json:"hour"`
    Score         float64  `json:"score"`
    PlanetaryHour PlanetID `json:"planetaryHour"`
    BestFor       []string `json:"bestFor"`
}
```

### 4.5 每周预测 (WeeklyForecast)

```go
type WeeklyForecast struct {
    StartDate      time.Time       `json:"startDate"`
    EndDate        time.Time       `json:"endDate"`
    OverallScore   float64         `json:"overallScore"`
    OverallTheme   string          `json:"overallTheme"`
    Dimensions     DailyDimensions `json:"dimensions"`
    DailySummaries []DailySummary  `json:"dailySummaries"`
    KeyDates       []KeyDate       `json:"keyDates"`
    BestDaysFor    map[string][]string `json:"bestDaysFor"`
    WeeklyTransits []WeeklyTransit `json:"weeklyTransits"`
}
```

---

## 5. 年限法结构

### 5.1 年度年限法 (AnnualProfection)

```go
type AnnualProfection struct {
    Year           int      `json:"year"`
    Age            int      `json:"age"`
    House          int      `json:"house"`
    HouseName      string   `json:"houseName"`
    HouseTheme     string   `json:"houseTheme"`
    HouseKeywords  []string `json:"houseKeywords"`
    Sign           ZodiacID `json:"sign"`
    SignName       string   `json:"signName"`
    LordOfYear     PlanetID `json:"lordOfYear"`
    LordName       string   `json:"lordName"`
    LordSymbol     string   `json:"lordSymbol"`
    LordNatalHouse int      `json:"lordNatalHouse"`
    LordNatalSign  ZodiacID `json:"lordNatalSign"`
    Description    string   `json:"description"`
}
```

### 5.2 人生年限法地图 (LifeProfectionMap)

```go
type LifeProfectionMap struct {
    Profections    []AnnualProfection `json:"profections"`
    CurrentYear    *AnnualProfection  `json:"currentYear"`
    UpcomingYears  []AnnualProfection `json:"upcomingYears"`
    CycleAnalysis  CycleAnalysis      `json:"cycleAnalysis"`
}
```

---

## 6. 人生趋势结构

### 6.1 人生趋势点 (LifeTrendPoint)

```go
type LifeTrendPoint struct {
    Date             time.Time       `json:"date"`
    Year             int             `json:"year"`
    Age              int             `json:"age"`
    OverallScore     float64         `json:"overallScore"`
    Harmonious       float64         `json:"harmonious"`
    Challenge        float64         `json:"challenge"`
    Transformation   float64         `json:"transformation"`
    Dimensions       LifeDimensions  `json:"dimensions"`
    DominantPlanet   PlanetID        `json:"dominantPlanet"`
    Profection       ProfectionSummary `json:"profection"`
    IsMajorTransit   bool            `json:"isMajorTransit"`
    MajorTransitName string          `json:"majorTransitName,omitempty"`
    LunarPhaseName   string          `json:"lunarPhaseName"`
    LunarPhaseAngle  float64         `json:"lunarPhaseAngle"`
}
```

### 6.2 人生维度 (LifeDimensions)

```go
type LifeDimensions struct {
    Career       float64 `json:"career"`
    Relationship float64 `json:"relationship"`
    Health       float64 `json:"health"`
    Finance      float64 `json:"finance"`
    Spiritual    float64 `json:"spiritual"`
}
```

### 6.3 人生趋势数据 (LifeTrendData)

```go
type LifeTrendData struct {
    BirthDate time.Time        `json:"birthDate"`
    Points    []LifeTrendPoint `json:"points"`
    Summary   LifeTrendSummary `json:"summary"`
    Cycles    LifeCycles       `json:"cycles"`
}
```

---

## 7. 统一时间序列结构

### 7.1 原始分数 (RawScore)

```go
type RawScore struct {
    Value   float64      `json:"value"`
    Factors ScoreFactors `json:"factors"`
}
```

### 7.2 分数因子 (ScoreFactors)

```go
type ScoreFactors struct {
    // 通用因子
    Aspects      float64 `json:"aspects"`
    Dignity      float64 `json:"dignity"`
    Retrograde   float64 `json:"retrograde"`
    
    // 小时级因子
    PlanetaryHour float64 `json:"planetaryHour,omitempty"`
    MoonVoid      float64 `json:"moonVoid,omitempty"`
    
    // 日级因子
    MoonSign     float64 `json:"moonSign,omitempty"`
    MoonPhase    float64 `json:"moonPhase,omitempty"`
    
    // 周级因子
    WeeklyAspects float64 `json:"weeklyAspects,omitempty"`
    
    // 月级因子
    OuterTransits float64 `json:"outerTransits,omitempty"`
    
    // 年级因子
    MajorTransits float64 `json:"majorTransits,omitempty"`
    Profection    float64 `json:"profection,omitempty"`
    Progression   float64 `json:"progression,omitempty"`
}
```

### 7.3 时间序列点 (TimeSeriesPoint)

```go
type TimeSeriesPoint struct {
    Time        time.Time          `json:"time"`
    Granularity TimeGranularity    `json:"granularity"`
    Raw         RawScore           `json:"raw"`
    Display     float64            `json:"display"`
    Dimensions  map[string]float64 `json:"dimensions"`
    Volatility  float64            `json:"volatility,omitempty"`
    Label       string             `json:"label,omitempty"`
}
```

### 7.4 时间序列统计 (TimeSeriesStats)

```go
type TimeSeriesStats struct {
    Mean       float64 `json:"mean"`
    StdDev     float64 `json:"stdDev"`
    Min        float64 `json:"min"`
    Max        float64 `json:"max"`
    Trend      string  `json:"trend"` // up, down, flat
    Volatility float64 `json:"volatility"`
}
```

### 7.5 时间序列 (TimeSeries)

```go
type TimeSeries struct {
    Granularity TimeGranularity   `json:"granularity"`
    Points      []TimeSeriesPoint `json:"points"`
    Stats       TimeSeriesStats   `json:"stats"`
}
```

---

## 8. 影响因子结构

### 8.1 影响因子类型 (InfluenceFactorType)

```go
type InfluenceFactorType string

const (
    FactorDignity        InfluenceFactorType = "dignity"
    FactorRetrograde     InfluenceFactorType = "retrograde"
    FactorAspectPhase    InfluenceFactorType = "aspectPhase"
    FactorAspectOrb      InfluenceFactorType = "aspectOrb"
    FactorOuterPlanet    InfluenceFactorType = "outerPlanet"
    FactorProfectionLord InfluenceFactorType = "profectionLord"
    FactorLunarPhase     InfluenceFactorType = "lunarPhase"
    FactorPlanetaryHour  InfluenceFactorType = "planetaryHour"
    FactorPersonal       InfluenceFactorType = "personal"
    FactorCustom         InfluenceFactorType = "custom"
)
```

### 8.2 影响因子 (InfluenceFactor)

```go
type InfluenceFactor struct {
    Type        InfluenceFactorType `json:"type"`
    Name        string              `json:"name"`
    Value       float64             `json:"value"`
    Weight      float64             `json:"weight"`
    Description string              `json:"description"`
    Source      string              `json:"source"`
    IsPositive  bool                `json:"isPositive"`
}
```

### 8.3 因子上下文 (FactorContext)

```go
type FactorContext struct {
    NatalChart       *NatalChart
    Date             time.Time
    TransitPositions []PlanetPosition
    ProgressedChart  *ProgressedChart
    Profection       *AnnualProfection
    MoonPhase        MoonPhase
    PlanetaryDay     PlanetID
    PlanetaryHour    PlanetID
}
```

### 8.4 因子权重 (FactorWeights)

```go
type FactorWeights struct {
    Dignity        float64 `json:"dignity"`
    Retrograde     float64 `json:"retrograde"`
    AspectPhase    float64 `json:"aspectPhase"`
    AspectOrb      float64 `json:"aspectOrb"`
    OuterPlanet    float64 `json:"outerPlanet"`
    ProfectionLord float64 `json:"profectionLord"`
    LunarPhase     float64 `json:"lunarPhase"`
    PlanetaryHour  float64 `json:"planetaryHour"`
    Personal       float64 `json:"personal"`
    Custom         float64 `json:"custom"`
}
```

### 8.5 因子结果 (FactorResult)

```go
type FactorResult struct {
    Factors              []InfluenceFactor  `json:"factors"`
    TotalAdjustment      float64            `json:"totalAdjustment"`
    PositiveFactors      []InfluenceFactor  `json:"positiveFactors"`
    NegativeFactors      []InfluenceFactor  `json:"negativeFactors"`
    DimensionAdjustments map[string]float64 `json:"dimensionAdjustments"`
}
```

---

## 9. 行运结构

### 9.1 行运事件 (TransitEvent)

```go
type TransitEvent struct {
    Date           time.Time            `json:"date"`
    TransitPlanet  PlanetID             `json:"transitPlanet"`
    NatalPlanet    PlanetID             `json:"natalPlanet"`
    Aspect         AspectData           `json:"aspect"`
    Phase          string               `json:"phase"` // approaching, exact, separating
    Intensity      float64              `json:"intensity"`
    Duration       TransitDuration      `json:"duration"`
    Interpretation TransitInterpretation `json:"interpretation"`
}
```

### 9.2 重大行运 (MajorTransit)

```go
type MajorTransit struct {
    Name         string    `json:"name"`
    Date         time.Time `json:"date"`
    Description  string    `json:"description"`
    Significance string    `json:"significance"` // high, medium, low
}
```

### 9.3 每日行运 (DailyTransit)

```go
type DailyTransit struct {
    Date      time.Time    `json:"date"`
    Aspects   []AspectData `json:"aspects"`
    Score     TransitScore `json:"score"`
    KeyTheme  string       `json:"keyTheme"`
    Intensity float64      `json:"intensity"`
}
```

### 9.4 行运分数 (TransitScore)

```go
type TransitScore struct {
    Total      float64 `json:"total"`
    Harmonious float64 `json:"harmonious"`
    Tense      float64 `json:"tense"`
}
```

---

## 10. 用户结构

### 10.1 用户 (User)

```go
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    BirthData BirthData `json:"birthData"`
    Settings  UserSettings `json:"settings,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

### 10.2 用户设置 (UserSettings)

```go
type UserSettings struct {
    FactorWeights  *FactorWeights `json:"factorWeights,omitempty"`
    DisplayOptions DisplayOptions `json:"displayOptions,omitempty"`
}
```

---

## 11. 常量数据

### 11.1 行星信息

```go
var Planets = []PlanetInfo{
    {Sun, "太阳", "☉", "#ffd700"},
    {Moon, "月亮", "☽", "#c0c0c0"},
    {Mercury, "水星", "☿", "#b5651d"},
    {Venus, "金星", "♀", "#ff69b4"},
    {Mars, "火星", "♂", "#dc143c"},
    {Jupiter, "木星", "♃", "#daa520"},
    {Saturn, "土星", "♄", "#8b7355"},
    {Uranus, "天王星", "♅", "#40e0d0"},
    {Neptune, "海王星", "♆", "#4169e1"},
    {Pluto, "冥王星", "♇", "#800080"},
    {NorthNode, "北交点", "☊", "#9932cc"},
    {Chiron, "凯龙星", "⚷", "#228b22"},
}
```

### 11.2 星座信息

```go
var ZodiacSigns = []ZodiacInfo{
    {Aries, "白羊座", "♈", "fire", "cardinal", Mars},
    {Taurus, "金牛座", "♉", "earth", "fixed", Venus},
    {Gemini, "双子座", "♊", "air", "mutable", Mercury},
    {Cancer, "巨蟹座", "♋", "water", "cardinal", Moon},
    {Leo, "狮子座", "♌", "fire", "fixed", Sun},
    {Virgo, "处女座", "♍", "earth", "mutable", Mercury},
    {Libra, "天秤座", "♎", "air", "cardinal", Venus},
    {Scorpio, "天蝎座", "♏", "water", "fixed", Pluto},
    {Sagittarius, "射手座", "♐", "fire", "mutable", Jupiter},
    {Capricorn, "摩羯座", "♑", "earth", "cardinal", Saturn},
    {Aquarius, "水瓶座", "♒", "air", "fixed", Uranus},
    {Pisces, "双鱼座", "♓", "water", "mutable", Neptune},
}
```

### 11.3 相位定义

```go
var AspectDefinitions = []AspectDefinition{
    {Conjunction, 0, 10, "合相"},
    {Sextile, 60, 6, "六分相"},
    {Square, 90, 8, "四分相"},
    {Trine, 120, 8, "三分相"},
    {Opposition, 180, 10, "对分相"},
}
```

### 11.4 行星权重

```go
var PlanetWeights = map[PlanetID]float64{
    Sun:       10,
    Moon:      10,
    Mercury:   4,
    Venus:     5,
    Mars:      6,
    Jupiter:   7,
    Saturn:    8,
    Uranus:    5,
    Neptune:   4,
    Pluto:     6,
    NorthNode: 3,
    Chiron:    3,
}
```

### 11.5 宫位主题

```go
var HouseDataList = []HouseData{
    {1, "命宫", "自我与身份", []string{"身份", "外表", "人格面具", "生命力"}},
    {2, "财帛宫", "资源与价值", []string{"金钱", "物质", "自我价值", "天赋"}},
    {3, "兄弟宫", "沟通与学习", []string{"沟通", "学习", "兄弟姐妹", "短途旅行"}},
    {4, "田宅宫", "家庭与根基", []string{"家庭", "根源", "安全感", "房产"}},
    {5, "子女宫", "创造与表达", []string{"创造力", "恋爱", "子女", "娱乐"}},
    {6, "奴仆宫", "服务与健康", []string{"工作", "健康", "日常习惯", "服务"}},
    {7, "夫妻宫", "关系与合作", []string{"婚姻", "合作", "公开敌人", "合约"}},
    {8, "疾厄宫", "转化与共享", []string{"死亡重生", "共享资源", "心理深度", "性"}},
    {9, "迁移宫", "探索与智慧", []string{"高等教育", "远行", "哲学", "宗教"}},
    {10, "官禄宫", "事业与成就", []string{"事业", "公众形象", "成就", "权威"}},
    {11, "福德宫", "愿景与社群", []string{"朋友", "团体", "理想", "希望"}},
    {12, "玄秘宫", "内省与超越", []string{"潜意识", "灵性", "隐退", "业力"}},
}
```

