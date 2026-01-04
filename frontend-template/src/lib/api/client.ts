/**
 * API Client
 * 前端 API 调用层
 */

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '';

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
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

// 类型定义
export interface BirthDataInput {
  name: string;
  date: string;
  latitude: number;
  longitude: number;
  timezone: string;
}

export interface NatalChartResponse {
  birthData: BirthDataInput;
  planets: any[];
  houses: any[];
  ascendant: number;
  midheaven: number;
  aspects: any[];
  dominantPlanets: string[];
  chartRuler: string;
  elementBalance: Record<string, number>;
  modalityBalance: Record<string, number>;
}

export interface DailyForecastResponse {
  type: 'daily';
  processed: boolean;
  date: string;
  overallScore: number;
  rawScore?: number;
  dimensions: Record<string, number>;
  factors?: {
    totalAdjustment: number;
    summary: string;
    appliedFactors: any[];
  };
  theme: string;
}

export interface WeeklyForecastResponse {
  type: 'weekly';
  processed: boolean;
  startDate: string;
  endDate: string;
  overallScore: number;
  weeklyTheme: string;
  days: any[];
}

// API 方法
export const calcApi = {
  calculateChart: (birthData: BirthDataInput) =>
    request<NatalChartResponse>('/api/calc/chart', {
      method: 'POST',
      body: JSON.stringify(birthData),
    }),

  calculateDaily: (birthData: BirthDataInput, date?: string, withFactors = true) =>
    request<DailyForecastResponse>('/api/calc/daily', {
      method: 'POST',
      body: JSON.stringify({ birthData, date, withFactors }),
    }),

  calculateWeekly: (birthData: BirthDataInput, date?: string, withFactors = true) =>
    request<WeeklyForecastResponse>('/api/calc/weekly', {
      method: 'POST',
      body: JSON.stringify({ birthData, date, withFactors }),
    }),

  calculateLifeTrend: (birthData: BirthDataInput, startYear?: number, endYear?: number) =>
    request<any>('/api/calc/life-trend', {
      method: 'POST',
      body: JSON.stringify({ birthData, startYear, endYear }),
    }),
};

export default { calc: calcApi };

