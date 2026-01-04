/**
 * API Client
 * 前端 API 调用层
 * 
 * 所有前端数据获取都通过这个客户端，实现前后端解耦
 */

// API 基础 URL，支持独立部署时配置
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '';

/**
 * 通用请求封装
 */
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });
  
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(error.error || `API Error: ${response.status}`);
  }
  
  return response.json();
}

// ============================================================
// 类型定义
// ============================================================

export interface BirthDataInput {
  name: string;
  date: string;  // ISO 格式
  latitude: number;
  longitude: number;
  timezone: string;
}

export interface NatalChartResponse {
  birthData: BirthDataInput;
  planets: Array<{
    id: string;
    name: string;
    symbol: string;
    longitude: number;
    latitude?: number;
    sign: string;
    signName: string;
    signSymbol?: string;
    signDegree: number;
    house: number;
    retrograde: boolean;
    dignityScore?: number;
  }>;
  houses: Array<{
    house: number;
    cusp: number;
    sign: string;
    signDegree?: number;
  }>;
  ascendant: number;
  midheaven: number;
  patterns: any[];
  elementBalance: Record<string, number>;
  modalityBalance: Record<string, number>;
  dominantPlanets: string[];
  chartRuler: string;
  aspects: Array<{
    planet1: string;
    planet2: string;
    aspectType: string;
    orb: number;
    applying: boolean;
    exactAngle?: number;
    actualAngle?: number;
    strength?: number;
    weight?: number;
  }>;
}

export interface DailyForecastResponse {
  type: 'daily';
  processed: boolean;
  date: string;
  overallScore: number;
  rawScore?: number;
  dimensions: {
    career: number;
    relationship: number;
    health: number;
    finance: number;
    spiritual: number;
  };
  factors?: {
    totalAdjustment: number;
    summary: string;
    appliedFactors: Array<{
      id: string;
      name: string;
      type: string;
      adjustment: number;
      dimension: string;
      reason: string;
    }>;
  };
  moonPhase: {
    phase: string;
    illumination: number;
  };
  moonSign: string;
  planetaryDay: string;
  theme: string;
}

export interface WeeklyForecastResponse {
  type: 'weekly';
  processed: boolean;
  startDate: string;
  endDate: string;
  overallScore: number;
  weeklyTheme: string;
  weeklyInsight: string;
  weeklyFactors?: {
    positive: Array<{ name: string; adjustment: number; reason: string }>;
    negative: Array<{ name: string; adjustment: number; reason: string }>;
  };
  dimensionTrends?: {
    career: number[];
    relationship: number[];
    health: number[];
    finance: number[];
    spiritual: number[];
  };
  days: Array<{
    date: string;
    overallScore: number;
    rawScore?: number;
    dimensions?: Record<string, number>;
  }>;
}

export interface UserSnapshotResponse {
  userId: string;
  timestamp: string;
  natal: {
    sunSign: string;
    moonSign: string;
    risingSign: string;
    dominantPlanets: string[];
    chartRuler: string;
  };
  current: {
    age: number;
    profectionHouse: number;
    profectionTheme: string;
    lordOfYear: string;
  };
  todayScore: number;
  todayTheme: string;
  activeTransits: Array<{
    transitPlanet: string;
    natalPlanet: string;
    aspect: string;
    interpretation: string;
  }>;
}

export interface LifeTrendResponse {
  type: 'yearly';
  birthDate: string;
  points: Array<{
    date: string;
    year: number;
    age: number;
    overallScore: number;
    dimensions: Record<string, number>;
    isMajorTransit: boolean;
    majorTransitName?: string;
  }>;
  summary: {
    overallTrend: string;
    peakYears: number[];
    challengeYears: number[];
  };
}

// ============================================================
// API 方法
// ============================================================

/**
 * 用户相关 API
 */
export const userApi = {
  /**
   * 创建用户（计算本命盘）
   */
  create: (birthData: BirthDataInput) =>
    request<{ id: string; natalChart: NatalChartResponse }>('/api/users', {
      method: 'POST',
      body: JSON.stringify(birthData),
    }),

  /**
   * 获取用户信息
   */
  get: (userId: string) =>
    request<{ user: any; natalChart: NatalChartResponse }>(`/api/users/${userId}`),

  /**
   * 获取用户实时快照
   */
  getSnapshot: (userId: string) =>
    request<UserSnapshotResponse>(`/api/users/${userId}/snapshot`),

  /**
   * 删除用户
   */
  delete: (userId: string) =>
    request<{ success: boolean }>(`/api/users/${userId}`, {
      method: 'DELETE',
    }),
};

/**
 * 预测相关 API
 */
export const forecastApi = {
  /**
   * 获取每日预测
   */
  getDaily: (userId: string, date?: string, withFactors = true) => {
    const params = new URLSearchParams({
      type: 'daily',
      withFactors: String(withFactors),
    });
    if (date) params.set('date', date);
    
    return request<DailyForecastResponse>(`/api/users/${userId}/forecast?${params}`);
  },

  /**
   * 获取每周预测
   */
  getWeekly: (userId: string, date?: string, withFactors = true) => {
    const params = new URLSearchParams({
      type: 'weekly',
      withFactors: String(withFactors),
    });
    if (date) params.set('date', date);
    
    return request<WeeklyForecastResponse>(`/api/users/${userId}/forecast?${params}`);
  },

  /**
   * 获取人生趋势
   */
  getLifeTrend: (userId: string, startYear?: number, endYear?: number) => {
    const params = new URLSearchParams({ type: 'yearly' });
    if (startYear) params.set('startYear', String(startYear));
    if (endYear) params.set('endYear', String(endYear));
    
    return request<LifeTrendResponse>(`/api/users/${userId}/forecast?${params}`);
  },

  /**
   * 获取年限法
   */
  getProfection: (userId: string, age?: number) => {
    const params = new URLSearchParams({ type: 'profection' });
    if (age !== undefined) params.set('age', String(age));
    
    return request<any>(`/api/users/${userId}/forecast?${params}`);
  },
};

/**
 * 计算相关 API（无需用户）
 */
export const calcApi = {
  /**
   * 直接计算本命盘（不创建用户）
   */
  calculateChart: (birthData: BirthDataInput) =>
    request<NatalChartResponse>('/api/calc/chart', {
      method: 'POST',
      body: JSON.stringify(birthData),
    }),

  /**
   * 直接计算每日预测（不创建用户）
   */
  calculateDaily: (birthData: BirthDataInput, date?: string, withFactors = true) =>
    request<DailyForecastResponse>('/api/calc/daily', {
      method: 'POST',
      body: JSON.stringify({ birthData, date, withFactors }),
    }),

  /**
   * 直接计算每周预测
   */
  calculateWeekly: (birthData: BirthDataInput, date?: string, withFactors = true) =>
    request<WeeklyForecastResponse>('/api/calc/weekly', {
      method: 'POST',
      body: JSON.stringify({ birthData, date, withFactors }),
    }),

  /**
   * 直接计算人生趋势
   */
  calculateLifeTrend: (birthData: BirthDataInput, startYear?: number, endYear?: number) =>
    request<LifeTrendResponse>('/api/calc/life-trend', {
      method: 'POST',
      body: JSON.stringify({ birthData, startYear, endYear }),
    }),
};

/**
 * 智能体相关 API
 */
export const agentApi = {
  /**
   * 获取用户上下文
   */
  getContext: (userId: string, query?: string) =>
    request<{ context: string; snapshot: UserSnapshotResponse }>('/api/agent/context', {
      method: 'POST',
      body: JSON.stringify({ userId, query }),
    }),

  /**
   * 智能体查询
   */
  query: (userId: string, query: string) =>
    request<{ content: string; suggestions: string[] }>('/api/agent/query', {
      method: 'POST',
      body: JSON.stringify({ userId, query }),
    }),
};

// 默认导出所有 API
export default {
  user: userApi,
  forecast: forecastApi,
  calc: calcApi,
  agent: agentApi,
};

