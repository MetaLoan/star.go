/**
 * API Hooks
 * React Hooks 封装 API 调用
 */

'use client';

import { useState, useCallback } from 'react';
import { calcApi, type BirthDataInput, type DailyForecastResponse, type WeeklyForecastResponse } from './client';

// ============================================================
// 类型定义
// ============================================================

export interface UseApiState<T> {
  data: T | null;
  loading: boolean;
  error: Error | null;
}

export interface ChartData {
  birthData: BirthDataInput;
  planets: any[];
  houses: any[];
  ascendant: number;
  midheaven: number;
  aspects: any[];
  patterns: any[];
  elementBalance: Record<string, number>;
  modalityBalance: Record<string, number>;
  dominantPlanets: string[];
  chartRuler: string;
}

// ============================================================
// Hooks
// ============================================================

/**
 * 本命盘计算 Hook
 */
export function useChartCalculation() {
  const [state, setState] = useState<UseApiState<ChartData>>({
    data: null,
    loading: false,
    error: null,
  });

  const calculate = useCallback(async (birthData: BirthDataInput) => {
    setState({ data: null, loading: true, error: null });
    
    try {
      const result = await calcApi.calculateChart(birthData);
      setState({ data: result as ChartData, loading: false, error: null });
      return result;
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Unknown error');
      setState({ data: null, loading: false, error: err });
      throw err;
    }
  }, []);

  const reset = useCallback(() => {
    setState({ data: null, loading: false, error: null });
  }, []);

  return {
    ...state,
    calculate,
    reset,
  };
}

/**
 * 每日预测 Hook
 */
export function useDailyForecast() {
  const [state, setState] = useState<UseApiState<DailyForecastResponse>>({
    data: null,
    loading: false,
    error: null,
  });

  const calculate = useCallback(async (
    birthData: BirthDataInput,
    date?: string,
    withFactors = true
  ) => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    
    try {
      const result = await calcApi.calculateDaily(birthData, date, withFactors);
      setState({ data: result, loading: false, error: null });
      return result;
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Unknown error');
      setState(prev => ({ ...prev, loading: false, error: err }));
      throw err;
    }
  }, []);

  return { ...state, calculate };
}

/**
 * 每周预测 Hook
 */
export function useWeeklyForecast() {
  const [state, setState] = useState<UseApiState<WeeklyForecastResponse>>({
    data: null,
    loading: false,
    error: null,
  });

  const calculate = useCallback(async (
    birthData: BirthDataInput,
    date?: string,
    withFactors = true
  ) => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    
    try {
      const result = await calcApi.calculateWeekly(birthData, date, withFactors);
      setState({ data: result, loading: false, error: null });
      return result;
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Unknown error');
      setState(prev => ({ ...prev, loading: false, error: err }));
      throw err;
    }
  }, []);

  return { ...state, calculate };
}

/**
 * 人生趋势 Hook
 */
export function useLifeTrend() {
  const [state, setState] = useState<UseApiState<any>>({
    data: null,
    loading: false,
    error: null,
  });

  const calculate = useCallback(async (
    birthData: BirthDataInput,
    startYear?: number,
    endYear?: number
  ) => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    
    try {
      const result = await calcApi.calculateLifeTrend(birthData, startYear, endYear);
      setState({ data: result, loading: false, error: null });
      return result;
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Unknown error');
      setState(prev => ({ ...prev, loading: false, error: err }));
      throw err;
    }
  }, []);

  return { ...state, calculate };
}

/**
 * 综合数据 Hook（一次性获取所有需要的数据）
 */
export function useAstroData() {
  const [birthData, setBirthData] = useState<BirthDataInput | null>(null);
  const [chart, setChart] = useState<ChartData | null>(null);
  const [daily, setDaily] = useState<DailyForecastResponse | null>(null);
  const [weekly, setWeekly] = useState<WeeklyForecastResponse | null>(null);
  const [lifeTrend, setLifeTrend] = useState<any | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const initialize = useCallback(async (data: BirthDataInput) => {
    setLoading(true);
    setError(null);
    setBirthData(data);
    
    try {
      // 并行请求
      const [chartResult, dailyResult, weeklyResult] = await Promise.all([
        calcApi.calculateChart(data),
        calcApi.calculateDaily(data, undefined, true),
        calcApi.calculateWeekly(data, undefined, true),
      ]);
      
      setChart(chartResult as ChartData);
      setDaily(dailyResult);
      setWeekly(weeklyResult);
      
      // 人生趋势较大，延迟加载
      const birthYear = new Date(data.date).getFullYear();
      const trendResult = await calcApi.calculateLifeTrend(data, birthYear, birthYear + 80);
      setLifeTrend(trendResult);
      
      setLoading(false);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Unknown error');
      setError(error);
      setLoading(false);
    }
  }, []);

  const refreshDaily = useCallback(async (date?: string) => {
    if (!birthData) return;
    try {
      const result = await calcApi.calculateDaily(birthData, date, true);
      setDaily(result);
    } catch (err) {
      console.error('Failed to refresh daily:', err);
    }
  }, [birthData]);

  const refreshWeekly = useCallback(async (date?: string) => {
    if (!birthData) return;
    try {
      const result = await calcApi.calculateWeekly(birthData, date, true);
      setWeekly(result);
    } catch (err) {
      console.error('Failed to refresh weekly:', err);
    }
  }, [birthData]);

  const reset = useCallback(() => {
    setBirthData(null);
    setChart(null);
    setDaily(null);
    setWeekly(null);
    setLifeTrend(null);
    setError(null);
  }, []);

  return {
    birthData,
    chart,
    daily,
    weekly,
    lifeTrend,
    loading,
    error,
    initialize,
    refreshDaily,
    refreshWeekly,
    reset,
  };
}

