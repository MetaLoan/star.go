/**
 * Server-side API utilities
 * 服务端 API 工具（用于 SSR）
 */

import { headers } from 'next/headers';

// 获取 API 基础 URL（服务端）
function getApiBaseUrl(): string {
  // 优先使用环境变量
  if (process.env.API_URL) {
    return process.env.API_URL;
  }
  
  // 在 Vercel 等平台上使用 VERCEL_URL
  if (process.env.VERCEL_URL) {
    return `https://${process.env.VERCEL_URL}`;
  }
  
  // 开发环境默认
  return 'http://localhost:3000';
}

/**
 * 服务端 fetch 封装
 */
async function serverFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const baseUrl = getApiBaseUrl();
  const url = `${baseUrl}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    // SSR 时禁用缓存或设置合适的缓存策略
    cache: 'no-store',
  });
  
  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`);
  }
  
  return response.json();
}

// ============================================================
// SSR 数据获取函数
// ============================================================

export interface BirthDataInput {
  name: string;
  date: string;
  latitude: number;
  longitude: number;
  timezone: string;
}

/**
 * SSR: 获取用户数据
 */
export async function fetchUserServer(userId: string) {
  return serverFetch<any>(`/api/users/${userId}`);
}

/**
 * SSR: 获取用户快照
 */
export async function fetchUserSnapshotServer(userId: string) {
  return serverFetch<any>(`/api/users/${userId}/snapshot`);
}

/**
 * SSR: 计算本命盘
 */
export async function calculateChartServer(birthData: BirthDataInput) {
  return serverFetch<any>('/api/calc/chart', {
    method: 'POST',
    body: JSON.stringify(birthData),
  });
}

/**
 * SSR: 计算每日预测
 */
export async function calculateDailyServer(
  birthData: BirthDataInput, 
  date?: string, 
  withFactors = true
) {
  return serverFetch<any>('/api/calc/daily', {
    method: 'POST',
    body: JSON.stringify({ birthData, date, withFactors }),
  });
}

/**
 * SSR: 计算人生趋势
 */
export async function calculateLifeTrendServer(
  birthData: BirthDataInput,
  startYear?: number,
  endYear?: number
) {
  return serverFetch<any>('/api/calc/life-trend', {
    method: 'POST',
    body: JSON.stringify({ birthData, startYear, endYear }),
  });
}

