/**
 * useAstroData Hook
 * 统一管理占星数据状态
 */

import { useState, useCallback, useMemo } from 'react';
import type {
  BirthData,
  NatalChart,
  DailyForecast,
  WeeklyForecast,
  LifeTrend,
  TimeSeriesPoint,
  AnnualProfection,
  ProfectionMap,
  InfluenceFactor,
} from '../types';
import { apiClient } from '../api/client';

interface AstroDataState {
  birthData: BirthData | null;
  natalChart: NatalChart | null;
  dailyForecast: DailyForecast | null;
  weeklyForecast: WeeklyForecast | null;
  lifeTrend: LifeTrend | null;
  timeSeries: TimeSeriesPoint[];
  profection: AnnualProfection | null;
  profectionMap: ProfectionMap | null;
  influenceFactors: InfluenceFactor[];
  loading: boolean;
  error: string | null;
}

interface UseAstroDataReturn extends AstroDataState {
  setBirthData: (data: BirthData) => Promise<void>;
  refreshDaily: () => Promise<void>;
  refreshWeekly: () => Promise<void>;
  loadLifeTrend: (startAge?: number, endAge?: number) => Promise<void>;
  loadTimeSeries: (start: string, end: string, granularity: 'hour' | 'day' | 'week') => Promise<void>;
  loadProfectionMap: (startAge?: number, endAge?: number) => Promise<void>;
  currentAge: number;
  clearError: () => void;
  isReady: boolean;
}

export function useAstroData(): UseAstroDataReturn {
  const [state, setState] = useState<AstroDataState>({
    birthData: null,
    natalChart: null,
    dailyForecast: null,
    weeklyForecast: null,
    lifeTrend: null,
    timeSeries: [],
    profection: null,
    profectionMap: null,
    influenceFactors: [],
    loading: false,
    error: null,
  });

  const setLoading = useCallback((loading: boolean) => {
    setState(prev => ({ ...prev, loading }));
  }, []);

  const setError = useCallback((error: string | null) => {
    setState(prev => ({ ...prev, error, loading: false }));
  }, []);

  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }));
  }, []);

  // 设置出生数据并计算本命盘
  const setBirthData = useCallback(async (data: BirthData) => {
    setLoading(true);
    try {
      // 计算本命盘
      const natalChart = await apiClient.calculateChart(data);
      
      // 计算今日预测
      const today = new Date().toISOString().split('T')[0];
      const dailyForecast = await apiClient.getDailyForecast(data, today);
      
      // 计算当前年限法
      const birthDate = new Date(data.year, data.month - 1, data.day);
      const currentAge = Math.floor(
        (Date.now() - birthDate.getTime()) / (365.25 * 24 * 60 * 60 * 1000)
      );
      const profection = await apiClient.getProfection(data, currentAge);

      setState(prev => ({
        ...prev,
        birthData: data,
        natalChart,
        dailyForecast,
        profection,
        loading: false,
        error: null,
      }));
    } catch (error) {
      setError(error instanceof Error ? error.message : '计算失败');
    }
  }, [setLoading, setError]);

  // 刷新每日预测
  const refreshDaily = useCallback(async () => {
    if (!state.birthData) return;
    
    setLoading(true);
    try {
      const today = new Date().toISOString().split('T')[0];
      const dailyForecast = await apiClient.getDailyForecast(state.birthData, today);
      setState(prev => ({ ...prev, dailyForecast, loading: false }));
    } catch (error) {
      setError(error instanceof Error ? error.message : '刷新失败');
    }
  }, [state.birthData, setLoading, setError]);

  // 刷新每周预测
  const refreshWeekly = useCallback(async () => {
    if (!state.birthData) return;
    
    setLoading(true);
    try {
      const today = new Date().toISOString().split('T')[0];
      const weeklyForecast = await apiClient.getWeeklyForecast(state.birthData, today);
      setState(prev => ({ ...prev, weeklyForecast, loading: false }));
    } catch (error) {
      setError(error instanceof Error ? error.message : '刷新失败');
    }
  }, [state.birthData, setLoading, setError]);

  // 加载生命趋势
  const loadLifeTrend = useCallback(async (startAge = 0, endAge = 80) => {
    if (!state.birthData) return;
    
    setLoading(true);
    try {
      const lifeTrend = await apiClient.getLifeTrend(state.birthData, startAge, endAge);
      setState(prev => ({ ...prev, lifeTrend, loading: false }));
    } catch (error) {
      setError(error instanceof Error ? error.message : '加载失败');
    }
  }, [state.birthData, setLoading, setError]);

  // 加载时间序列
  const loadTimeSeries = useCallback(async (
    start: string,
    end: string,
    granularity: 'hour' | 'day' | 'week'
  ) => {
    if (!state.birthData) return;
    
    setLoading(true);
    try {
      const timeSeries = await apiClient.getTimeSeries(
        state.birthData,
        start,
        end,
        granularity
      );
      setState(prev => ({ ...prev, timeSeries, loading: false }));
    } catch (error) {
      setError(error instanceof Error ? error.message : '加载失败');
    }
  }, [state.birthData, setLoading, setError]);

  // 加载年限法地图
  const loadProfectionMap = useCallback(async (startAge = 0, endAge = 80) => {
    if (!state.birthData) return;
    
    setLoading(true);
    try {
      const profectionMap = await apiClient.getProfectionMap(state.birthData, startAge, endAge);
      setState(prev => ({ ...prev, profectionMap, loading: false }));
    } catch (error) {
      setError(error instanceof Error ? error.message : '加载失败');
    }
  }, [state.birthData, setLoading, setError]);

  // 计算当前年龄
  const currentAge = useMemo(() => {
    if (!state.birthData) return 0;
    const birthDate = new Date(state.birthData.year, state.birthData.month - 1, state.birthData.day);
    return Math.floor((Date.now() - birthDate.getTime()) / (365.25 * 24 * 60 * 60 * 1000));
  }, [state.birthData]);

  // 是否准备就绪（有本命盘数据）
  const isReady = useMemo(() => {
    return state.birthData !== null && state.natalChart !== null;
  }, [state.birthData, state.natalChart]);

  return {
    ...state,
    setBirthData,
    refreshDaily,
    refreshWeekly,
    loadLifeTrend,
    loadTimeSeries,
    loadProfectionMap,
    currentAge,
    clearError,
    isReady,
  };
}

export default useAstroData;

