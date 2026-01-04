'use client';

import { useState, useCallback } from 'react';
import { calcApi, type BirthDataInput, type DailyForecastResponse, type WeeklyForecastResponse, type NatalChartResponse } from './client';

export interface UseApiState<T> {
  data: T | null;
  loading: boolean;
  error: Error | null;
}

export function useChartCalculation() {
  const [state, setState] = useState<UseApiState<NatalChartResponse>>({
    data: null,
    loading: false,
    error: null,
  });

  const calculate = useCallback(async (birthData: BirthDataInput) => {
    setState({ data: null, loading: true, error: null });
    
    try {
      const result = await calcApi.calculateChart(birthData);
      setState({ data: result, loading: false, error: null });
      return result;
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Unknown error');
      setState({ data: null, loading: false, error: err });
      throw err;
    }
  }, []);

  return { ...state, calculate };
}

export function useAstroData() {
  const [birthData, setBirthData] = useState<BirthDataInput | null>(null);
  const [chart, setChart] = useState<NatalChartResponse | null>(null);
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
      const [chartResult, dailyResult, weeklyResult] = await Promise.all([
        calcApi.calculateChart(data),
        calcApi.calculateDaily(data, undefined, true),
        calcApi.calculateWeekly(data, undefined, true),
      ]);
      
      setChart(chartResult);
      setDaily(dailyResult);
      setWeekly(weeklyResult);
      
      const birthYear = new Date(data.date).getFullYear();
      const trendResult = await calcApi.calculateLifeTrend(data, birthYear, birthYear + 80);
      setLifeTrend(trendResult);
      
      setLoading(false);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
      setLoading(false);
    }
  }, []);

  const refreshDaily = useCallback(async (date?: string) => {
    if (!birthData) return;
    const result = await calcApi.calculateDaily(birthData, date, true);
    setDaily(result);
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
    reset,
  };
}

