package models

import "time"

// ==================== 基础类型 ====================

// PlanetID 行星标识符
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

// ZodiacID 星座标识符
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

// AspectType 相位类型
type AspectType string

const (
	Conjunction AspectType = "conjunction" // 合相 0°
	Sextile     AspectType = "sextile"     // 六分相 60°
	Square      AspectType = "square"      // 四分相 90°
	Trine       AspectType = "trine"       // 三分相 120°
	Opposition  AspectType = "opposition"  // 对分相 180°
)

// TimeGranularity 时间粒度
type TimeGranularity string

const (
	GranularityHour  TimeGranularity = "hour"
	GranularityDay   TimeGranularity = "day"
	GranularityWeek  TimeGranularity = "week"
	GranularityMonth TimeGranularity = "month"
	GranularityYear  TimeGranularity = "year"
)

// Dignity 尊贵度
type Dignity string

const (
	DignityDomicile   Dignity = "domicile"   // 入庙
	DignityExaltation Dignity = "exaltation" // 旺相
	DignityDetriment  Dignity = "detriment"  // 落陷
	DignityFall       Dignity = "fall"       // 失势
	DignityPeregrine  Dignity = "peregrine"  // 游离
)

// ==================== 核心数据结构 ====================

// BirthData 出生数据
type BirthData struct {
	Name      string  `json:"name"`
	Year      int     `json:"year"`
	Month     int     `json:"month"`
	Day       int     `json:"day"`
	Hour      int     `json:"hour"`
	Minute    int     `json:"minute"`
	Second    int     `json:"second"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  float64 `json:"timezone"` // 时区偏移（小时），支持半时区如 5.5
}

// ToTime 将出生数据转换为 time.Time
func (b BirthData) ToTime() time.Time {
	loc := time.FixedZone("Birth", int(b.Timezone*3600))
	return time.Date(b.Year, time.Month(b.Month), b.Day, b.Hour, b.Minute, b.Second, 0, loc)
}

// PlanetPosition 行星位置
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

// HouseCusp 宫位
type HouseCusp struct {
	House    int      `json:"house"`
	Cusp     float64  `json:"cusp"`
	Sign     ZodiacID `json:"sign"`
	SignName string   `json:"signName"`
}

// AspectData 相位数据
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

// ==================== 星盘结构 ====================

// NatalChart 本命盘
type NatalChart struct {
	BirthData       BirthData          `json:"birthData"`
	Planets         []PlanetPosition   `json:"planets"`
	Houses          []HouseCusp        `json:"houses"`
	Ascendant       float64            `json:"ascendant"`
	Midheaven       float64            `json:"midheaven"`
	Aspects         []AspectData       `json:"aspects"`
	Patterns        []string           `json:"patterns"`
	ElementBalance  map[string]float64 `json:"elementBalance"`
	ModalityBalance map[string]float64 `json:"modalityBalance"`
	DominantPlanets []PlanetID         `json:"dominantPlanets"`
	ChartRuler      PlanetID           `json:"chartRuler"`
}

// ProgressedChart 推运盘
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

// ProgressedPlanet 推运行星
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

// LunarPhaseInfo 月相信息
type LunarPhaseInfo struct {
	Phase       string   `json:"phase"`
	Name        string   `json:"name"`
	Angle       float64  `json:"angle"`
	Description string   `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

// ==================== 预测结构 ====================

// DailyDimensions 每日维度分数
type DailyDimensions struct {
	Career       float64 `json:"career"`
	Relationship float64 `json:"relationship"`
	Health       float64 `json:"health"`
	Finance      float64 `json:"finance"`
	Spiritual    float64 `json:"spiritual"`
}

// MoonPhase 月相
type MoonPhase struct {
	Phase        string     `json:"phase"`
	Name         string     `json:"name"`
	Angle        float64    `json:"angle"`
	Illumination float64    `json:"illumination"`
	VoidStart    *time.Time `json:"voidStart,omitempty"`
	VoidEnd      *time.Time `json:"voidEnd,omitempty"`
}

// MoonSignInfo 月亮星座信息
type MoonSignInfo struct {
	Sign     ZodiacID `json:"sign"`
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
}

// HourlyForecast 小时预测
type HourlyForecast struct {
	Hour          int      `json:"hour"`
	Score         float64  `json:"score"`
	PlanetaryHour PlanetID `json:"planetaryHour"`
	BestFor       []string `json:"bestFor"`
}

// DailyForecast 每日预测
type DailyForecast struct {
	Date            time.Time        `json:"date"`
	DayOfWeek       string           `json:"dayOfWeek"`
	OverallScore    float64          `json:"overallScore"`
	OverallTheme    string           `json:"overallTheme"`
	Dimensions      DailyDimensions  `json:"dimensions"`
	MoonPhase       MoonPhase        `json:"moonPhase"`
	MoonSign        MoonSignInfo     `json:"moonSign"`
	HourlyBreakdown []HourlyForecast `json:"hourlyBreakdown"`
	ActiveAspects   []AspectData     `json:"activeAspects"`
	Factors         *FactorResult    `json:"factors,omitempty"`
	TopFactors      []InfluenceFactor `json:"topFactors,omitempty"`
}

// DailySummary 每日摘要
type DailySummary struct {
	Date         string   `json:"date"`
	DayOfWeek    string   `json:"dayOfWeek"`
	OverallScore float64  `json:"overallScore"`
	MoonSign     ZodiacID `json:"moonSign"`
	KeyTheme     string   `json:"keyTheme"`
}

// KeyDate 关键日期
type KeyDate struct {
	Date         string `json:"date"`
	Event        string `json:"event"`
	Significance string `json:"significance"`
}

// WeeklyTransit 周度行运
type WeeklyTransit struct {
	Planet      PlanetID   `json:"planet"`
	Aspect      AspectType `json:"aspect"`
	NatalPlanet PlanetID   `json:"natalPlanet"`
	Peak        string     `json:"peak"`
	Theme       string     `json:"theme"`
}

// WeeklyForecast 每周预测
type WeeklyForecast struct {
	StartDate       time.Time           `json:"startDate"`
	EndDate         time.Time           `json:"endDate"`
	OverallScore    float64             `json:"overallScore"`
	OverallTheme    string              `json:"overallTheme"`
	Dimensions      DailyDimensions     `json:"dimensions"`
	DailySummaries  []DailySummary      `json:"dailySummaries"`
	KeyDates        []KeyDate           `json:"keyDates"`
	BestDaysFor     map[string][]string `json:"bestDaysFor"`
	WeeklyTransits  []WeeklyTransit     `json:"weeklyTransits"`
	WeeklyFactors   *FactorResult       `json:"weeklyFactors,omitempty"`
	DimensionTrends map[string]string   `json:"dimensionTrends"`
}

// ==================== 年限法结构 ====================

// AnnualProfection 年度年限法
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

// CycleAnalysis 周期分析
type CycleAnalysis struct {
	FirstCycle            map[string]interface{} `json:"firstCycle"`
	SecondCycle           map[string]interface{} `json:"secondCycle"`
	ThirdCycle            map[string]interface{} `json:"thirdCycle"`
	CurrentCycleNumber    int                    `json:"currentCycleNumber"`
	YearsIntoCurrentCycle int                    `json:"yearsIntoCurrentCycle"`
}

// LifeProfectionMap 人生年限法地图
type LifeProfectionMap struct {
	Profections   []AnnualProfection `json:"profections"`
	CurrentYear   *AnnualProfection  `json:"currentYear"`
	UpcomingYears []AnnualProfection `json:"upcomingYears"`
	CycleAnalysis CycleAnalysis      `json:"cycleAnalysis"`
}

// ==================== 人生趋势结构 ====================

// ProfectionSummary 年限法摘要
type ProfectionSummary struct {
	House      int      `json:"house"`
	Theme      string   `json:"theme"`
	LordOfYear PlanetID `json:"lordOfYear"`
}

// LifeTrendPoint 人生趋势点
type LifeTrendPoint struct {
	Date             time.Time         `json:"date"`
	Year             int               `json:"year"`
	Age              int               `json:"age"`
	OverallScore     float64           `json:"overallScore"`
	Harmonious       float64           `json:"harmonious"`
	Challenge        float64           `json:"challenge"`
	Transformation   float64           `json:"transformation"`
	Dimensions       DailyDimensions   `json:"dimensions"`
	DominantPlanet   PlanetID          `json:"dominantPlanet"`
	Profection       ProfectionSummary `json:"profection"`
	IsMajorTransit   bool              `json:"isMajorTransit"`
	MajorTransitName string            `json:"majorTransitName,omitempty"`
	LunarPhaseName   string            `json:"lunarPhaseName"`
	LunarPhaseAngle  float64           `json:"lunarPhaseAngle"`
}

// LifeTrendSummary 人生趋势摘要
type LifeTrendSummary struct {
	OverallTrend        string `json:"overallTrend"`
	PeakYears           []int  `json:"peakYears"`
	ChallengeYears      []int  `json:"challengeYears"`
	TransformationYears []int  `json:"transformationYears"`
}

// SaturnCycle 土星周期
type SaturnCycle struct {
	Age         int    `json:"age"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}

// LifeCycles 人生周期
type LifeCycles struct {
	SaturnCycles     []SaturnCycle              `json:"saturnCycles"`
	JupiterCycles    []map[string]interface{}   `json:"jupiterCycles"`
	ProfectionCycles []map[string]interface{}   `json:"profectionCycles"`
}

// LifeTrendData 人生趋势数据
type LifeTrendData struct {
	Type      string           `json:"type"`
	BirthDate time.Time        `json:"birthDate"`
	Points    []LifeTrendPoint `json:"points"`
	Summary   LifeTrendSummary `json:"summary"`
	Cycles    LifeCycles       `json:"cycles"`
}

// ==================== 统一时间序列结构 ====================

// ScoreFactors 分数因子
type ScoreFactors struct {
	Aspects       float64 `json:"aspects"`
	Dignity       float64 `json:"dignity"`
	Retrograde    float64 `json:"retrograde"`
	PlanetaryHour float64 `json:"planetaryHour,omitempty"`
	MoonVoid      float64 `json:"moonVoid,omitempty"`
	MoonSign      float64 `json:"moonSign,omitempty"`
	MoonPhase     float64 `json:"moonPhase,omitempty"`
	WeeklyAspects float64 `json:"weeklyAspects,omitempty"`
	OuterTransits float64 `json:"outerTransits,omitempty"`
	MajorTransits float64 `json:"majorTransits,omitempty"`
	Profection    float64 `json:"profection,omitempty"`
	Progression   float64 `json:"progression,omitempty"`
}

// RawScore 原始分数
type RawScore struct {
	Value   float64      `json:"value"`
	Factors ScoreFactors `json:"factors"`
}

// TimeSeriesPoint 时间序列点
type TimeSeriesPoint struct {
	Time        time.Time          `json:"time"`
	Label       string             `json:"label"`
	Granularity TimeGranularity    `json:"granularity"`
	Raw         RawScore           `json:"raw"`
	Display     float64            `json:"display"`
	Dimensions  map[string]float64 `json:"dimensions"`
	Volatility  float64            `json:"volatility,omitempty"`
}

// TimeSeriesStats 时间序列统计
type TimeSeriesStats struct {
	Mean       float64 `json:"mean"`
	StdDev     float64 `json:"stdDev"`
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Trend      string  `json:"trend"`
	Volatility float64 `json:"volatility"`
}

// TimeSeries 时间序列
type TimeSeries struct {
	Granularity TimeGranularity   `json:"granularity"`
	Points      []TimeSeriesPoint `json:"points"`
	Stats       TimeSeriesStats   `json:"stats"`
}

// ==================== 影响因子结构 ====================

// InfluenceFactorType 影响因子类型
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

// InfluenceFactor 影响因子
type InfluenceFactor struct {
	ID          string              `json:"id"`
	Type        InfluenceFactorType `json:"type"`
	Name        string              `json:"name"`
	Value       float64             `json:"value"`
	Weight      float64             `json:"weight"`
	Adjustment  float64             `json:"adjustment"`
	Dimension   string              `json:"dimension,omitempty"`
	Description string              `json:"description"`
	Reason      string              `json:"reason,omitempty"`
	IsPositive  bool                `json:"isPositive"`
}

// FactorWeights 因子权重配置（可运营调整）
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

// FactorResult 因子计算结果
type FactorResult struct {
	Factors              []InfluenceFactor  `json:"factors"`
	TotalAdjustment      float64            `json:"totalAdjustment"`
	PositiveFactors      []InfluenceFactor  `json:"positiveFactors"`
	NegativeFactors      []InfluenceFactor  `json:"negativeFactors"`
	DimensionAdjustments map[string]float64 `json:"dimensionAdjustments"`
}

// ==================== 行运结构 ====================

// TransitDuration 行运持续时间
type TransitDuration struct {
	Start string `json:"start"`
	Peak  string `json:"peak"`
	End   string `json:"end"`
}

// TransitInterpretation 行运解读
type TransitInterpretation struct {
	Theme    string   `json:"theme"`
	Keywords []string `json:"keywords"`
	Advice   string   `json:"advice"`
}

// TransitEvent 行运事件
type TransitEvent struct {
	Date           time.Time             `json:"date"`
	TransitPlanet  PlanetID              `json:"transitPlanet"`
	NatalPlanet    PlanetID              `json:"natalPlanet"`
	Aspect         AspectData            `json:"aspect"`
	Phase          string                `json:"phase"`
	Intensity      float64               `json:"intensity"`
	Duration       TransitDuration       `json:"duration"`
	Interpretation TransitInterpretation `json:"interpretation"`
}

// MajorTransit 重大行运
type MajorTransit struct {
	Name         string    `json:"name"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	Significance string    `json:"significance"`
}

// TransitScore 行运分数
type TransitScore struct {
	Total      float64 `json:"total"`
	Harmonious float64 `json:"harmonious"`
	Tense      float64 `json:"tense"`
}

// TransitResult 行运结果
type TransitResult struct {
	StartDate      string         `json:"startDate"`
	EndDate        string         `json:"endDate"`
	Events         []TransitEvent `json:"events"`
	OverallScore   float64        `json:"overallScore"`
	DominantThemes []string       `json:"dominantThemes"`
}

// ==================== 用户结构 ====================

// UserSettings 用户设置
type UserSettings struct {
	FactorWeights  *FactorWeights         `json:"factorWeights,omitempty"`
	DisplayOptions map[string]interface{} `json:"displayOptions,omitempty"`
}

// User 用户
type User struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	BirthData  BirthData    `json:"birthData"`
	NatalChart *NatalChart  `json:"natalChart,omitempty"`
	Settings   UserSettings `json:"settings,omitempty"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
}

// UserSnapshot 用户快照
type UserSnapshot struct {
	User            *User            `json:"user"`
	CurrentDate     time.Time        `json:"currentDate"`
	Age             int              `json:"age"`
	DailyForecast   *DailyForecast   `json:"dailyForecast"`
	Profection      *AnnualProfection `json:"profection"`
	ActiveTransits  []TransitEvent   `json:"activeTransits"`
	ProgressedChart *ProgressedChart `json:"progressedChart"`
}

// AgentContext 智能体上下文
type AgentContext struct {
	Users          []AgentUserContext `json:"users"`
	CurrentDate    time.Time          `json:"currentDate"`
	GlobalTransits GlobalTransits     `json:"globalTransits"`
}

// AgentUserContext 智能体用户上下文
type AgentUserContext struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	CurrentState map[string]interface{} `json:"currentState"`
}

// GlobalTransits 全局行运
type GlobalTransits struct {
	SunSign           ZodiacID   `json:"sunSign"`
	MoonSign          ZodiacID   `json:"moonSign"`
	MoonPhase         string     `json:"moonPhase"`
	RetrogradePlanets []PlanetID `json:"retrogradePlanets"`
}

// AgentQueryResponse 智能体查询响应
type AgentQueryResponse struct {
	Response string                 `json:"response"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

