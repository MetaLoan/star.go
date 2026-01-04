/**
 * Star 前端类型定义
 * 与后端 API 数据结构完全对应
 */

// ==================== 枚举类型 ====================

// 后端使用小写
export type PlanetID = 
  | 'sun' | 'moon' | 'mercury' | 'venus' | 'mars' 
  | 'jupiter' | 'saturn' | 'uranus' | 'neptune' | 'pluto' 
  | 'northNode' | 'chiron';

export type ZodiacID = 
  | 'aries' | 'taurus' | 'gemini' | 'cancer' | 'leo' | 'virgo'
  | 'libra' | 'scorpio' | 'sagittarius' | 'capricorn' | 'aquarius' | 'pisces';

export type AspectType = 
  | 'conjunction' | 'sextile' | 'square' | 'trine' | 'opposition';

export type TimeGranularity = 'hour' | 'day' | 'week' | 'month' | 'year';

export type Dignity = 'Domicile' | 'Exaltation' | 'Detriment' | 'Fall' | 'Peregrine';

export type Dimension = 'career' | 'relationship' | 'health' | 'finance' | 'spiritual';

// ==================== 基础数据结构 ====================

export interface BirthData {
  year: number;
  month: number;
  day: number;
  hour: number;
  minute: number;
  second?: number;
  latitude: number;
  longitude: number;
  timezone: number;
}

export interface PlanetPosition {
  id: PlanetID;
  name: string;
  symbol: string;
  longitude: number;
  latitude: number;
  sign: ZodiacID;
  signName: string;
  signSymbol: string;
  signDegree: number;
  retrograde: boolean;
  house: number;
  dignityScore: number;
}

export interface HouseCusp {
  house: number;
  cusp: number;
  sign: ZodiacID;
  signName: string;
}

export interface AspectData {
  planet1: PlanetID;
  planet2: PlanetID;
  aspectType: AspectType;
  exactAngle: number;
  actualAngle: number;
  orb: number;
  applying: boolean;
  strength: number;
  weight: number;
  interpretation: string;
}

// ==================== 星盘数据 ====================

export interface NatalChart {
  birthData: BirthData;
  planets: PlanetPosition[];
  houses: HouseCusp[];
  aspects: AspectData[];
  ascendant: number;
  midheaven: number;
  dominantPlanets: PlanetID[];
  chartRuler: PlanetID;
  elementBalance: Record<string, number>;
  modalityBalance: Record<string, number>;
  patterns?: unknown;
}

export interface ProgressedChart {
  progressedDate: string;
  planets: PlanetPosition[];
  lunarPhase: number;
  lunarPhaseName: string;
}

// ==================== 维度分数 ====================

export interface DimensionScores {
  career: number;
  relationship: number;
  health: number;
  finance: number;
  spiritual: number;
}

// ==================== 预测数据 ====================

export interface HourlyForecast {
  hour: number;
  score: number;
  planetaryHour: PlanetID;
  bestFor: string[];
}

export interface MoonPhaseInfo {
  phase: string;
  name: string;
  angle: number;
  illumination: number;
}

export interface MoonSignInfo {
  sign: ZodiacID;
  name: string;
  keywords: string[];
}

export interface DailyForecast {
  date: string;
  dayOfWeek?: string;
  overallScore: number;
  overallTheme?: string;
  dimensions?: DimensionScores;
  moonPhase?: MoonPhaseInfo;
  moonSign?: MoonSignInfo;
  hourlyBreakdown?: HourlyForecast[];
  activeAspects?: AspectData[];
}

export interface DailySummary {
  date: string;
  dayOfWeek: string;
  overallScore: number;
  moonSign?: string;
  moonPhase?: string;
  theme?: string;
  keyTheme?: string;
  dimensions?: DimensionScores;
}

export interface KeyDate {
  date: string;
  event?: string;
  reason?: string;
  significance?: string;
}

export interface WeeklyTransit {
  planet: PlanetID;
  aspect: AspectType;
  natalPlanet: PlanetID;
  peak: string;
  theme: string;
}

export interface WeeklyForecast {
  startDate: string;
  endDate: string;
  overallScore: number;
  overallTheme: string;
  dimensions: DimensionScores;
  dailySummaries: DailySummary[];
  keyDates: KeyDate[];
  bestDaysFor: Record<string, string[]>;
  weeklyTransits: WeeklyTransit[];
  dimensionTrends: Record<string, string>;
  // 兼容字段
  weekStart?: string;
  weekEnd?: string;
  dailyForecasts?: DailyForecast[];
  weeklyAverage?: number;
  weeklyThemes?: string[];
  bestDays?: string[];
  challengeDays?: string[];
}

export interface TransitEvent {
  transitPlanet: PlanetID;
  natalPlanet: PlanetID;
  aspect: AspectType;
  exactTime: string;
  influence: number;
  description: string;
}

// ==================== 月亮空亡数据 ====================

export interface VoidOfCourseInfo {
  isVoid: boolean;
  startTime?: string;
  endTime?: string;
  duration: number;
  nextSign: string;
  lastAspect: string;
  influence: number;
}

// ==================== 行星时数据 ====================

export interface PlanetaryHourInfo {
  planetaryHour: number;
  ruler?: PlanetID;
  planetName: string;
  planetSymbol: string;
  dayRuler?: PlanetID;
  influence: number;
  bestFor: string[];
}

export interface PlanetaryHoursResponse {
  date: string;
  hours: PlanetaryHourInfo[];
}

// ==================== 年限法数据 ====================

export interface AnnualProfection {
  age: number;
  year: number;
  house: number;
  houseName: string;
  houseTheme: string;
  houseKeywords: string[];
  sign: ZodiacID;
  signName: string;
  lordOfYear: PlanetID;
  lordName: string;
  lordSymbol: string;
  lordNatalHouse: number;
  lordNatalSign: ZodiacID;
  description: string;
}

export interface ProfectionMap {
  profections: AnnualProfection[];
  currentYear: AnnualProfection;
  upcomingYears: AnnualProfection[];
  cycleAnalysis: {
    currentCycle: number;
    cycleTheme: string;
    yearsInCycle: number;
  };
}

// ==================== 生命趋势数据 ====================

export interface LifeTrendProfection {
  house: number;
  theme: string;
  lordOfYear: PlanetID;
}

export interface LifeTrendPoint {
  date: string;
  year: number;
  age: number;
  overallScore: number;
  harmonious: number;
  challenge: number;
  transformation: number;
  dimensions: DimensionScores;
  dominantPlanet: PlanetID;
  profection: LifeTrendProfection;
  isMajorTransit: boolean;
  majorTransits?: string[];
  lunarPhaseName: string;
  lunarPhaseAngle: number;
}

export interface LifeTrendSummary {
  overallTrend: string;
  peakYears: number[];
  challengeYears: number[];
  transformationYears: number[];
}

export interface LifeTrendCycles {
  saturn: { ages: number[]; description: string };
  jupiter: { ages: number[]; description: string };
  uranus: { ages: number[]; description: string };
  chiron: { ages: number[]; description: string };
}

export interface LifeTrend {
  type: string;
  birthDate: string;
  points: LifeTrendPoint[];
  summary: LifeTrendSummary;
  cycles?: LifeTrendCycles;
}

// ==================== 时间序列数据 ====================

export interface TimeSeriesPoint {
  timestamp: string;
  granularity: TimeGranularity;
  rawScore: number;
  normalizedScore: number;
  dimensions: DimensionScores;
  factors: InfluenceFactor[];
}

// ==================== 影响因子 ====================

export interface InfluenceFactor {
  id: string;
  type: string;
  name: string;
  value: number;
  weight: number;
  adjustment: number;
  description: string;
  isPositive: boolean;
  dimension?: string;
}

export interface InfluenceFactorGroup {
  category: string;
  factors: InfluenceFactor[];
  totalWeight: number;
}

// ==================== 用户数据 ====================

export interface User {
  id: string;
  name: string;
  birthData: BirthData;
  natalChart?: NatalChart;
  settings: UserSettings;
  createdAt: string;
  updatedAt: string;
}

export interface UserSettings {
  factorWeights: FactorWeights;
  displayOptions: DisplayOptions;
}

export interface FactorWeights {
  transit: number;
  progression: number;
  profection: number;
  lunar: number;
}

export interface DisplayOptions {
  showRetrograde: boolean;
  showAspects: boolean;
  showHouses: boolean;
  colorScheme: 'cosmic' | 'classic' | 'minimal';
}

// ==================== API 响应类型 ====================

export interface ApiResponse<T> {
  data: T;
  error?: string;
  timestamp: string;
}

export interface HealthResponse {
  status: string;
  service: string;
  version: string;
  features: string[];
}

// ==================== 组件 Props 类型 ====================

export interface ChartDisplayProps {
  chart: NatalChart;
  size?: number;
  showAspects?: boolean;
  highlightPlanet?: PlanetID;
  onPlanetClick?: (planet: PlanetID) => void;
}

export interface TimelineProps {
  data: TimeSeriesPoint[];
  dimension?: Dimension | 'overall';
  height?: number;
  onPointClick?: (point: TimeSeriesPoint) => void;
}

export interface ForecastCardProps {
  forecast: DailyForecast;
  isToday?: boolean;
  onClick?: () => void;
}
