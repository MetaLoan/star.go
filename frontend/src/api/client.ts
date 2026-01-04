/**
 * Star API 客户端
 * 处理与后端的所有通信，使用 JSON 格式
 */

import type {
  BirthData,
  NatalChart,
  DailyForecast,
  WeeklyForecast,
  LifeTrend,
  TimeSeriesPoint,
  AnnualProfection,
  ProfectionMap,
  ProgressedChart,
  TransitEvent,
  User,
  HealthResponse,
  TimeGranularity,
  VoidOfCourseInfo,
  PlanetaryHourInfo,
  PlanetaryHoursResponse,
} from '../types';

// API 基础 URL - 可通过环境变量配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * 通用请求函数
 */
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error(`API request failed: ${endpoint}`, error);
    throw error;
  }
}

// ==================== 健康检查 ====================

export async function checkHealth(): Promise<HealthResponse> {
  return request<HealthResponse>('/health');
}

// ==================== 星盘计算 ====================

export async function calculateChart(birthData: BirthData): Promise<NatalChart> {
  return request<NatalChart>('/api/calc/chart', {
    method: 'POST',
    body: JSON.stringify(birthData),
  });
}

// ==================== 预测计算 ====================

export async function getDailyForecast(
  birthData: BirthData,
  date: string
): Promise<DailyForecast> {
  return request<DailyForecast>('/api/calc/daily', {
    method: 'POST',
    body: JSON.stringify({ birthData, date }),
  });
}

export async function getWeeklyForecast(
  birthData: BirthData,
  startDate: string
): Promise<WeeklyForecast> {
  return request<WeeklyForecast>('/api/calc/weekly', {
    method: 'POST',
    body: JSON.stringify({ birthData, startDate }),
  });
}

// ==================== 生命趋势 ====================

export async function getLifeTrend(
  birthData: BirthData,
  startAge: number,
  endAge: number
): Promise<LifeTrend> {
  return request<LifeTrend>('/api/calc/life-trend', {
    method: 'POST',
    body: JSON.stringify({ birthData, startAge, endAge }),
  });
}

// ==================== 时间序列 ====================

export async function getTimeSeries(
  birthData: BirthData,
  start: string,
  end: string,
  granularity: TimeGranularity
): Promise<TimeSeriesPoint[]> {
  return request<TimeSeriesPoint[]>('/api/calc/time-series', {
    method: 'POST',
    body: JSON.stringify({ birthData, start, end, granularity }),
  });
}

// ==================== 年限法 ====================

export async function getProfection(
  birthData: BirthData,
  age: number
): Promise<AnnualProfection> {
  return request<AnnualProfection>('/api/calc/profection', {
    method: 'POST',
    body: JSON.stringify({ birthData, age }),
  });
}

export async function getProfectionMap(
  birthData: BirthData,
  startAge: number,
  endAge: number
): Promise<ProfectionMap> {
  return request<ProfectionMap>('/api/calc/profection-map', {
    method: 'POST',
    body: JSON.stringify({ birthData, startAge, endAge }),
  });
}

// ==================== 行运 ====================

export async function getTransits(
  birthData: BirthData,
  date: string
): Promise<TransitEvent[]> {
  return request<TransitEvent[]>('/api/calc/transits', {
    method: 'POST',
    body: JSON.stringify({ birthData, date }),
  });
}

// ==================== 推运 ====================

export async function getProgressions(
  birthData: BirthData,
  date: string
): Promise<ProgressedChart> {
  return request<ProgressedChart>('/api/calc/progressions', {
    method: 'POST',
    body: JSON.stringify({ birthData, date }),
  });
}

// ==================== 月亮空亡 ====================

export async function getVoidOfCourse(
  date?: string,
  latitude?: number,
  longitude?: number
): Promise<VoidOfCourseInfo> {
  return request<VoidOfCourseInfo>('/api/calc/void-of-course', {
    method: 'POST',
    body: JSON.stringify({ date, latitude, longitude }),
  });
}

// ==================== 行星时 ====================

export async function getPlanetaryHour(
  date?: string,
  latitude?: number,
  longitude?: number
): Promise<PlanetaryHourInfo> {
  return request<PlanetaryHourInfo>('/api/calc/planetary-hour', {
    method: 'POST',
    body: JSON.stringify({ date, latitude, longitude, fullDay: false }),
  });
}

export async function getPlanetaryHoursForDay(
  date?: string,
  latitude?: number,
  longitude?: number
): Promise<PlanetaryHoursResponse> {
  return request<PlanetaryHoursResponse>('/api/calc/planetary-hour', {
    method: 'POST',
    body: JSON.stringify({ date, latitude, longitude, fullDay: true }),
  });
}

// ==================== 用户管理 ====================

export async function createUser(data: {
  name: string;
  birthData: BirthData;
}): Promise<User> {
  return request<User>('/api/users', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getUser(id: string): Promise<User> {
  return request<User>(`/api/users/${id}`);
}

export async function updateUser(
  id: string,
  data: Partial<User>
): Promise<User> {
  return request<User>(`/api/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteUser(id: string): Promise<void> {
  await request<void>(`/api/users/${id}`, {
    method: 'DELETE',
  });
}

export async function getUserForecast(
  userId: string,
  startDate: string,
  days: number
): Promise<DailyForecast[]> {
  return request<DailyForecast[]>(
    `/api/users/${userId}/forecast?start=${startDate}&days=${days}`
  );
}

// ==================== Agent API ====================

export async function getAgentContext(userId: string): Promise<{
  natalChart: NatalChart;
  currentTransits: TransitEvent[];
  todayForecast: DailyForecast;
  profection: AnnualProfection;
}> {
  return request(`/api/agent/context?userId=${userId}`);
}

export async function queryAgent(
  userId: string,
  query: string
): Promise<{ response: string; data?: unknown }> {
  return request('/api/agent/query', {
    method: 'POST',
    body: JSON.stringify({ userId, query }),
  });
}

// ==================== 导出默认客户端 ====================

export const apiClient = {
  checkHealth,
  calculateChart,
  getDailyForecast,
  getWeeklyForecast,
  getLifeTrend,
  getTimeSeries,
  getProfection,
  getProfectionMap,
  getTransits,
  getProgressions,
  getVoidOfCourse,
  getPlanetaryHour,
  getPlanetaryHoursForDay,
  createUser,
  getUser,
  updateUser,
  deleteUser,
  getUserForecast,
  getAgentContext,
  queryAgent,
};

export default apiClient;

