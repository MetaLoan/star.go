/**
 * Star 占星计算平台 - 主应用组件
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Tabs, Tab, Spinner, Button, Switch } from '@heroui/react';
import { useAstroData } from './hooks/useAstroData';
import { AstroChartContainer } from './components/chart/AstroChartContainer';
import { DimensionRadarChart5832XY } from './components/chart/DimensionRadarChart5832XY';
import { BirthDataForm2943KL } from './components/input/BirthDataForm2943KL';
import { ScoreCard5612XY, DimensionScoresCard5612XY } from './components/ui/ScoreCard5612XY';
import { GranularitySelector4721VW, type TimeGranularity } from './components/ui/GranularitySelector4721VW';
import { DailyForecastCard7821MN } from './components/forecast/DailyForecastCard7821MN';
import { LifeTimeline4529PQ } from './components/timeline/LifeTimeline4529PQ';
import { ProfectionWheel6183RS } from './components/timeline/ProfectionWheel6183RS';
import { InteractiveTrendChart9823EF } from './components/timeline/InteractiveTrendChart9823EF';
import { InfluenceFactorsPanel8274TU } from './components/factors/InfluenceFactorsPanel8274TU';
import { CustomFactorEditor9456DE } from './components/factors/CustomFactorEditor9456DE';
import { RealtimeDimensionDashboard7392WZ } from './components/ui/RealtimeDimensionDashboard7392WZ';
import { MultiGranularityScoreViewer8475QR } from './components/ui/MultiGranularityScoreViewer8475QR';
import { ScoreBreakdownPopup5932MN } from './components/ui/ScoreBreakdownPopup5932MN';
import MiniAppDemo from './miniapp/MiniAppDemo';
import type { PlanetID, BirthData, InfluenceFactor, ScoreBreakdownAllResponse, ActiveFactorsResponse } from './types';
import { PLANET_NAMES, PLANET_SYMBOLS, PLANET_COLORS, formatDegree } from './utils/astro';
import { apiClient } from './api/client';

// 模拟影响因子数据（后续从 API 获取）
const MOCK_INFLUENCE_FACTORS: InfluenceFactor[] = [
  { id: '1', type: 'dignity', name: '太阳入庙狮子', value: 3, weight: 1, adjustment: 3, description: '太阳在狮子座获得入庙尊贵', isPositive: true },
  { id: '2', type: 'dignity', name: '金星入旺双鱼', value: 2, weight: 1, adjustment: 2, description: '金星在双鱼座获得旺相尊贵', isPositive: true },
  { id: '3', type: 'retrograde', name: '水星逆行', value: -2, weight: 1, adjustment: -2, description: '水星逆行期间沟通需谨慎', isPositive: false },
  { id: '4', type: 'aspectPhase', name: '木星三分太阳', value: 1.5, weight: 0.8, adjustment: 1.2, description: '木星与太阳形成和谐相位', isPositive: true },
  { id: '5', type: 'aspectPhase', name: '土星四分月亮', value: -1.2, weight: 0.8, adjustment: -0.96, description: '土星与月亮形成紧张相位', isPositive: false },
  { id: '6', type: 'lunarPhase', name: '月亮上弦', value: 0.5, weight: 0.7, adjustment: 0.35, description: '月相处于上弦阶段，适合行动', isPositive: true },
  { id: '7', type: 'profectionLord', name: '年主星木星', value: 1.0, weight: 1, adjustment: 1, description: '今年由木星主管，带来扩张机遇', isPositive: true },
];

function App() {
  const {
    birthData,
    natalChart,
    dailyForecast,
    weeklyForecast,
    lifeTrend,
    timeSeries,
    profection,
    profectionMap,
    currentAge,
    loading,
    error,
    isReady,
    setBirthData,
    refreshWeekly,
    loadLifeTrend,
    loadTimeSeries,
    extendTimeSeries,
    loadProfectionMap,
    clearError,
  } = useAstroData();

  const [selectedTab, setSelectedTab] = useState('chart');
  const [showMiniApp, setShowMiniApp] = useState(false);
  const [highlightedPlanet, setHighlightedPlanet] = useState<PlanetID | null>(null);
  const [expandedForecast, setExpandedForecast] = useState<string | null>(null);
  const [showFactorEditor, setShowFactorEditor] = useState(false);
  
  // 新增：多粒度趋势图状态
  const [trendGranularity, setTrendGranularity] = useState<TimeGranularity>('daily');
  
  // 时间序列数据范围跟踪（用于动态加载更多数据）
  const [timeSeriesRange, setTimeSeriesRange] = useState<{
    start: Date;
    end: Date;
  } | null>(null);
  const [isLoadingMoreData, setIsLoadingMoreData] = useState(false);

  // 分数组成浮窗状态（点击趋势点触发）
  const [breakdownOpen, setBreakdownOpen] = useState(false);
  const [breakdownPosition, setBreakdownPosition] = useState({ x: 0, y: 0 }); // 浮窗位置
  const [breakdownQueryTime, setBreakdownQueryTime] = useState<string | null>(null);
  const [breakdownLoading, setBreakdownLoading] = useState(false);
  const [breakdownError, setBreakdownError] = useState<string | null>(null);
  const [breakdownData, setBreakdownData] = useState<ScoreBreakdownAllResponse | null>(null);
  const [activeFactorsData, setActiveFactorsData] = useState<ActiveFactorsResponse | null>(null);
  const [breakdownDimension, setBreakdownDimension] = useState<'overall' | 'career' | 'relationship' | 'health' | 'finance' | 'spiritual'>('overall');
  const [breakdownGranularity, setBreakdownGranularity] = useState<'hour' | 'day' | 'week' | 'month' | 'year'>('hour');
  const breakdownReqIdRef = useRef(0);
  
  // 格式化日期为本地 ISO 时间字符串（带时区偏移）
  const formatLocalISO = useCallback((date: Date, timezone: number = 8) => {
    const offsetHours = Math.floor(Math.abs(timezone));
    const offsetMins = Math.round((Math.abs(timezone) % 1) * 60);
    const sign = timezone >= 0 ? '+' : '-';
    const pad = (n: number) => n.toString().padStart(2, '0');
    
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}${sign}${pad(offsetHours)}:${pad(offsetMins)}`;
  }, []);

  // 新增：实时运势状态
  const [realtimeScore, setRealtimeScore] = useState<{
    score: number;
    dimensions: { career: number; relationship: number; health: number; finance: number; spiritual: number };
    time: string;
  } | null>(null);
  
  // 新增：自定义因子状态
  interface CustomFactor {
    id: string;
    operation: 'AddScore' | 'SubScore' | 'MulScore' | 'SetScore';
    value: number;
    dimension: 'career' | 'relationship' | 'health' | 'finance' | 'spiritual';
    duration: number;
    startTime: string;
    name?: string;
  }
  const [customFactors, setCustomFactors] = useState<CustomFactor[]>([]);
  
  // 添加自定义因子
  const handleAddCustomFactor = (factor: Omit<CustomFactor, 'id'>) => {
    const newFactor = { ...factor, id: Date.now().toString() };
    setCustomFactors([...customFactors, newFactor]);
    // TODO: 调用后端 API 保存
    console.log('添加自定义因子:', newFactor);
  };
  
  // 删除自定义因子
  const handleRemoveCustomFactor = (id: string) => {
    setCustomFactors(customFactors.filter(f => f.id !== id));
    // TODO: 调用后端 API 删除
    console.log('删除自定义因子:', id);
  };

  // 加载实时运势（每分钟刷新）
  useEffect(() => {
    if (!isReady || !birthData) return;
    
    const fetchRealtimeScore = async () => {
      try {
        // 获取当前 UTC 时间，然后转换为用户时区的本地时间
        const nowUtc = Date.now();
        const userTimezoneOffset = birthData.timezone * 60 * 60 * 1000; // 毫秒
        const userLocalTime = new Date(nowUtc + userTimezoneOffset + new Date().getTimezoneOffset() * 60 * 1000);
        
        const start = new Date(userLocalTime);
        start.setMinutes(0, 0, 0);
        const end = new Date(start);
        end.setHours(end.getHours() + 1);
        
        const startStr = formatLocalISO(start, birthData.timezone);
        const endStr = formatLocalISO(end, birthData.timezone);
        
        console.log('[实时运势] 请求时间范围:', startStr, '-', endStr);
        
        const response = await fetch('http://localhost:8080/api/calc/time-series', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            birthData: {
              year: birthData.year,
              month: birthData.month,
              day: birthData.day,
              hour: birthData.hour,
              minute: birthData.minute,
              latitude: birthData.latitude,
              longitude: birthData.longitude,
              timezone: birthData.timezone,
            },
            start: startStr,
            end: endStr,
            granularity: 'hour',
          }),
        });
        
        if (response.ok) {
          const data = await response.json();
          if (data.points && data.points.length > 0) {
            const point = data.points[0];
            setRealtimeScore({
              score: point.display,
              dimensions: point.dimensions || {
                career: 50, relationship: 50, health: 50, finance: 50, spiritual: 50,
              },
              time: userLocalTime.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
            });
          }
        }
      } catch (err) {
        console.error('获取实时运势失败:', err);
      }
    };
    
    // 立即获取一次
    fetchRealtimeScore();
    
    // 每分钟刷新
    const interval = setInterval(fetchRealtimeScore, 60000);
    return () => clearInterval(interval);
  }, [isReady, birthData]);

  // 加载周预测（当有出生数据时）
  useEffect(() => {
    if (isReady && !weeklyForecast) {
      refreshWeekly();
    }
  }, [isReady, weeklyForecast, refreshWeekly]);

  // 加载趋势数据（当切换到趋势 tab 时）
  useEffect(() => {
    if (isReady && selectedTab === 'trend') {
      if (!lifeTrend) {
        loadLifeTrend(0, 80);
      }
      if (!profectionMap) {
        loadProfectionMap(0, 80);
      }
    }
  }, [isReady, selectedTab, lifeTrend, profectionMap, loadLifeTrend, loadProfectionMap]);

  // 加载时间序列数据（当粒度变化时）
  useEffect(() => {
    if (isReady && selectedTab === 'trend') {
      // 年度视图使用 lifeTrend 数据（80年），不需要额外加载
      if (trendGranularity === 'yearly') {
        return;
      }
      
      const now = new Date();
      let startDate: Date;
      let endDate: Date;
      let granularity: 'hour' | 'day' | 'week' | 'month' | 'year';
      
      // 使用用户时区或默认 UTC+8
      const tz = birthData?.timezone ?? 8;
      
      switch (trendGranularity) {
        case 'hourly':
          // 显示过去24小时（每小时一个点）- 从整点开始
          endDate = new Date(now);
          endDate.setMinutes(0, 0, 0);
          startDate = new Date(endDate.getTime() - 24 * 60 * 60 * 1000);
          granularity = 'hour';
          break;
        case 'daily':
          // 显示最近45天（每天一个点，增加范围减少边界触发）
          endDate = new Date(now);
          endDate.setHours(0, 0, 0, 0);
          endDate.setDate(endDate.getDate() + 1); // 明天00:00（包含今天）
          startDate = new Date(endDate.getTime() - 45 * 24 * 60 * 60 * 1000);
          granularity = 'day';
          break;
        case 'weekly':
          // 显示最近16周（每周一个点）
          endDate = new Date(now);
          endDate.setHours(0, 0, 0, 0);
          startDate = new Date(endDate.getTime() - 16 * 7 * 24 * 60 * 60 * 1000);
          granularity = 'week';
          break;
        case 'monthly':
          // 显示最近18个月（每月一个点）
          endDate = new Date(now.getFullYear(), now.getMonth() + 2, 1); // 下下月1号
          startDate = new Date(now.getFullYear() - 1, now.getMonth() - 6, 1);
          granularity = 'month';
          break;
        default:
          startDate = new Date(now.getTime() - 45 * 24 * 60 * 60 * 1000);
          endDate = new Date(now);
          granularity = 'day';
      }
      
      console.log(`[趋势图] 初始化 ${trendGranularity} 数据:`, startDate.toISOString(), '-', endDate.toISOString());
      
      loadTimeSeries(formatLocalISO(startDate, tz), formatLocalISO(endDate, tz), granularity);
      
      // 记录当前数据范围（直接保存 Date 对象，避免时区转换问题）
      setTimeSeriesRange({
        start: startDate,
        end: endDate,
      });
    }
  }, [isReady, selectedTab, trendGranularity, loadTimeSeries, birthData?.timezone, formatLocalISO]);

  // 处理图表可视范围变化 - 动态加载更多数据
  const handleVisibleRangeChange = useCallback(async (range: {
    from: Date;
    to: Date;
    needsMoreBefore: boolean;
    needsMoreAfter: boolean;
  }) => {
    // 严格的前置条件检查
    if (isLoadingMoreData) {
      console.log('[趋势图] 跳过：正在加载中');
      return;
    }
    if (!birthData || !timeSeriesRange) {
      console.log('[趋势图] 跳过：缺少出生数据或时间范围');
      return;
    }
    if (trendGranularity === 'yearly') {
      console.log('[趋势图] 跳过：年度视图使用固定数据');
      return;
    }
    
    const tz = birthData.timezone ?? 8;
    let newStart = timeSeriesRange.start;
    let newEnd = timeSeriesRange.end;
    let hasChange = false;
    
    // 计算扩展量（根据粒度不同）
    const extendDays = {
      hourly: 1,      // 扩展 1 天
      daily: 15,      // 扩展 15 天
      weekly: 28,     // 扩展 4 周
      monthly: 180,   // 扩展 6 个月
    }[trendGranularity] || 15;
    
    const extendMs = extendDays * 24 * 60 * 60 * 1000;
    
    // 限制向过去扩展的最大范围（最多 2 年前）
    const minPast = new Date();
    minPast.setFullYear(minPast.getFullYear() - 2);
    
    if (range.needsMoreBefore) {
      const proposedStart = new Date(timeSeriesRange.start.getTime() - extendMs);
      // 不要超过最小限制
      if (proposedStart > minPast) {
        newStart = proposedStart;
        hasChange = true;
        console.log('[趋势图] 向左扩展到:', newStart.toISOString());
      } else if (timeSeriesRange.start > minPast) {
        newStart = minPast;
        hasChange = true;
        console.log('[趋势图] 向左扩展到最大限制:', newStart.toISOString());
      }
    }
    
    if (range.needsMoreAfter) {
      // 向右扩展（不超过当前时间太远，最多到未来 1 年）
      const maxFuture = new Date();
      maxFuture.setFullYear(maxFuture.getFullYear() + 1);
      const proposedEnd = new Date(timeSeriesRange.end.getTime() + extendMs);
      if (proposedEnd < maxFuture) {
        newEnd = proposedEnd;
        hasChange = true;
        console.log('[趋势图] 向右扩展到:', newEnd.toISOString());
      } else if (timeSeriesRange.end < maxFuture) {
        newEnd = maxFuture;
        hasChange = true;
        console.log('[趋势图] 向右扩展到最大限制:', newEnd.toISOString());
      }
    }
    
    // 检查是否有变化
    if (!hasChange) {
      console.log('[趋势图] 跳过：已达到数据边界');
      return;
    }
    
    // 加载扩展后的数据
    setIsLoadingMoreData(true);
    console.log('[趋势图] 开始加载扩展数据...');
    
    const granularityMap: Record<string, 'hour' | 'day' | 'week' | 'month' | 'year'> = {
      hourly: 'hour',
      daily: 'day',
      weekly: 'week',
      monthly: 'month',
    };
    
    try {
      const apiGranularity = granularityMap[trendGranularity] || 'day';

      // 增量加载：只请求新增区间，合并去重，避免每次重算整段范围
      if (range.needsMoreBefore && newStart.getTime() !== timeSeriesRange.start.getTime()) {
        await extendTimeSeries(
          formatLocalISO(newStart, tz),
          formatLocalISO(timeSeriesRange.start, tz),
          apiGranularity,
          'before'
        );
      }
      if (range.needsMoreAfter && newEnd.getTime() !== timeSeriesRange.end.getTime()) {
        await extendTimeSeries(
          formatLocalISO(timeSeriesRange.end, tz),
          formatLocalISO(newEnd, tz),
          apiGranularity,
          'after'
        );
      }
      
      // 更新范围
      setTimeSeriesRange({
        start: newStart,
        end: newEnd,
      });
      console.log('[趋势图] 数据加载成功');
    } catch (err) {
      console.error('[趋势图] 加载失败:', err);
    } finally {
      setIsLoadingMoreData(false);
    }
  }, [isLoadingMoreData, birthData, timeSeriesRange, trendGranularity, extendTimeSeries, formatLocalISO]);

  // 点击趋势图数据点：所有粒度都显示浮窗
  // - 小时粒度：调用 score-breakdown-all，显示分数+因子
  // - 日/周/月/年粒度：调用 active-factors，显示正/负影响因子
  const handleTrendPointClick = useCallback(async (point: { time: string }, dimension: 'overall' | 'career' | 'relationship' | 'health' | 'finance' | 'spiritual' = 'overall', event?: MouseEvent) => {
    if (!birthData) return;
    if (!point?.time) return;

    // 转换粒度格式：hourly -> hour, daily -> day, etc.
    const granularityMap: Record<string, 'hour' | 'day' | 'week' | 'month' | 'year'> = {
      hourly: 'hour',
      daily: 'day',
      weekly: 'week',
      monthly: 'month',
      yearly: 'year',
    };
    const apiGranularity = granularityMap[trendGranularity] || 'day';

    // 构建查询时间
    let queryTime = point.time;
    let displayTime = point.time;
    
    // 年粒度特殊处理：从 label 提取年龄用于显示
    if (apiGranularity === 'year') {
      // 时间已经是 ISO 格式（如 "2020-06-15T12:00:00+08:00"），直接使用
      // 从 ISO 时间中提取年份用于显示
      const yearMatch = point.time.match(/^(\d{4})/);
      if (yearMatch) {
        const year = parseInt(yearMatch[1], 10);
        const age = year - birthData.year;
        displayTime = `${year}年 (${age}岁)`;
      }
    } else if (!queryTime.includes('T')) {
      // 其他粒度：如果时间格式不完整，补充为完整格式
      if (queryTime.match(/^\d{4}$/)) {
        queryTime = `${queryTime}-01-01T12:00:00+08:00`;
      } else if (queryTime.match(/^\d{4}-\d{2}$/)) {
        queryTime = `${queryTime}-15T12:00:00+08:00`;
      } else if (queryTime.match(/^\d{4}-\d{2}-\d{2}$/)) {
        queryTime = `${queryTime}T12:00:00+08:00`;
      }
    }

    // 记录点击位置（用于浮窗定位）
    const clickX = event?.clientX ?? window.innerWidth / 2;
    const clickY = event?.clientY ?? window.innerHeight / 2;
    setBreakdownPosition({ x: clickX, y: clickY });

    setBreakdownDimension(dimension);
    setBreakdownGranularity(apiGranularity);
    setBreakdownOpen(true);
    setBreakdownQueryTime(displayTime); // 显示用的时间
    setBreakdownLoading(true);
    setBreakdownError(null);
    setBreakdownData(null);
    setActiveFactorsData(null);

    const reqId = ++breakdownReqIdRef.current;
    try {
      if (apiGranularity === 'hour') {
        // 小时粒度：使用 score-breakdown-all API
        const res = await apiClient.getScoreBreakdownAll(birthData, queryTime);
        if (reqId !== breakdownReqIdRef.current) return;
        setBreakdownData(res);
      } else {
        // 日/周/月/年粒度：使用 active-factors API
        const res = await apiClient.getActiveFactors(birthData, queryTime, apiGranularity, 'all');
        if (reqId !== breakdownReqIdRef.current) return;
        setActiveFactorsData(res);
      }
    } catch (e) {
      if (reqId !== breakdownReqIdRef.current) return;
      setBreakdownError(e instanceof Error ? e.message : '加载失败');
    } finally {
      if (reqId !== breakdownReqIdRef.current) return;
      setBreakdownLoading(false);
    }
  }, [birthData, trendGranularity]);

  // 处理出生数据提交
  const handleBirthDataSubmit = async (data: BirthData) => {
    await setBirthData(data);
  };

  if (showMiniApp) {
    return (
      <div className="relative">
        <MiniAppDemo />
        <button 
          onClick={() => setShowMiniApp(false)}
          className="fixed top-4 right-4 z-[100] bg-black/10 hover:bg-black/20 text-black px-3 py-1 rounded-full text-xs font-medium backdrop-blur-sm transition-all"
        >
          ← 返回平台
        </button>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-4 md:p-8">
      {/* 标题 */}
      <motion.header
        className="text-center mb-8"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-4xl md:text-5xl font-display font-bold mb-2">
          <span className="bg-gradient-to-r from-[#00D4FF] via-[#A855F7] to-[#FF6B9D] bg-clip-text text-transparent">
            ✦ Star
          </span>
        </h1>
        <p className="text-white/60 text-lg">占星计算验证平台</p>
      </motion.header>

      {/* 错误提示 */}
      <AnimatePresence>
        {error && (
          <motion.div
            className="max-w-2xl mx-auto mb-6 p-4 bg-red-500/20 border border-red-500/30 rounded-lg flex items-center justify-between"
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
          >
            <span className="text-red-300">❌ {error}</span>
            <Button size="sm" variant="light" onPress={clearError}>
              关闭
            </Button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* 主内容区 */}
      {!isReady ? (
        // 未输入出生数据时显示表单
        <motion.div
          className="max-w-md mx-auto"
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
        >
          <BirthDataForm2943KL
            onSubmit={handleBirthDataSubmit}
            loading={loading}
          />
        </motion.div>
      ) : (
        // 已有星盘数据时显示完整界面
        <div className="max-w-7xl mx-auto">
          {/* Tab 导航 */}
          <Tabs
            selectedKey={selectedTab}
            onSelectionChange={(key) => setSelectedTab(key as string)}
            classNames={{
              tabList: "bg-white/5 p-1 rounded-xl",
              cursor: "bg-white/10",
              tab: "text-white/60 data-[selected=true]:text-white",
            }}
            className="mb-6"
          >
            <Tab key="chart" title="🌟 星盘" />
            <Tab key="forecast" title="📅 预测" />
            <Tab key="trend" title="📈 趋势" />
            <Tab key="factors" title="📊 因子" />
            <Tab key="settings" title="⚙️ 设置" />
          </Tabs>

          {/* 加载状态 */}
          {loading && (
            <div className="flex items-center justify-center py-12">
              <Spinner size="lg" color="primary" />
              <span className="ml-3 text-white/60">计算中...</span>
            </div>
          )}

          {/* Tab 内容 */}
          <AnimatePresence mode="wait">
            {/* ==================== 星盘 Tab ==================== */}
            {selectedTab === 'chart' && natalChart && birthData && (
              <motion.div
                key="chart"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* 实时五维运势仪表盘 - 顶部显示 */}
                <RealtimeDimensionDashboard7392WZ
                  birthData={birthData}
                  refreshInterval={60000}
                />
                
                {/* 星盘和详情区域 */}
                <div className="grid lg:grid-cols-2 gap-6">
                {/* 星盘 SVG */}
                <div className="glass-card p-6 flex justify-center">
                  <AstroChartContainer
                    data={natalChart}
                    width={Math.min(600, window.innerWidth - 80)}
                    height={Math.min(600, window.innerWidth - 80)}
                  />
                </div>

                {/* 星盘详情 */}
                <div className="space-y-4">
                  {/* 基本信息卡片 */}
                  <div className="glass-card p-4">
                    <h3 className="text-lg font-medium text-white mb-3">📍 基本信息</h3>
                    <div className="grid grid-cols-2 gap-3 text-sm">
                      <div>
                        <span className="text-white/60">上升点：</span>
                        <span className="text-[#00D4FF]">{formatDegree(natalChart.ascendant)}</span>
                      </div>
                      <div>
                        <span className="text-white/60">天顶：</span>
                        <span className="text-[#FF6B9D]">{formatDegree(natalChart.midheaven)}</span>
                      </div>
                      <div>
                        <span className="text-white/60">主导行星：</span>
                        <span>{natalChart.dominantPlanets.map(p => PLANET_SYMBOLS[p]).join(' ')}</span>
                      </div>
                      <div>
                        <span className="text-white/60">命主星：</span>
                        <span style={{ color: PLANET_COLORS[natalChart.chartRuler] }}>
                          {PLANET_SYMBOLS[natalChart.chartRuler]} {PLANET_NAMES[natalChart.chartRuler]}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* 行星列表 */}
                  <div className="glass-card p-4">
                    <h3 className="text-lg font-medium text-white mb-3">🪐 行星位置</h3>
                    <div className="grid grid-cols-2 gap-2 text-sm max-h-64 overflow-y-auto">
                      {natalChart.planets.map(planet => (
                        <motion.div
                          key={planet.id}
                          className={`flex items-center gap-2 p-2 rounded-lg cursor-pointer transition-colors ${
                            highlightedPlanet === planet.id
                              ? 'bg-white/10'
                              : 'hover:bg-white/5'
                          }`}
                          onClick={() => setHighlightedPlanet(
                            highlightedPlanet === planet.id ? null : planet.id
                          )}
                          whileHover={{ scale: 1.02 }}
                        >
                          <span
                            className="text-lg"
                            style={{ color: PLANET_COLORS[planet.id] }}
                          >
                            {PLANET_SYMBOLS[planet.id]}
                          </span>
                          <div className="flex-1">
                            <div className="text-white/80">{PLANET_NAMES[planet.id]}</div>
                            <div className="text-white/40 text-xs">
                              {planet.signName} {Math.floor(planet.signDegree)}°{Math.floor((planet.signDegree % 1) * 60)}'
                              {planet.retrograde && <span className="text-red-400 ml-1">℞</span>}
                            </div>
                          </div>
                          <div className="text-white/30 text-xs">
                            {planet.house}宫
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  </div>

                  {/* 年限法信息 */}
                  {profection && (
                    <div className="glass-card p-4">
                      <h3 className="text-lg font-medium text-white mb-3">🔮 年限法</h3>
                      <div className="text-sm space-y-2">
                        <div className="flex justify-between">
                          <span className="text-white/60">当前年龄：</span>
                          <span className="text-white">{profection.age}岁</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-white/60">激活宫位：</span>
                          <span className="text-[#00D4FF]">第{profection.house}宫 ({profection.houseName})</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-white/60">年主星：</span>
                          <span className="text-[#ffd700]">
                            {profection.lordSymbol} {profection.lordName}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-white/60">主题：</span>
                          <span className="text-white/80">{profection.houseTheme}</span>
                        </div>
                        {profection.houseKeywords && profection.houseKeywords.length > 0 && (
                          <div className="mt-2 pt-2 border-t border-white/10">
                            <span className="text-white/60">关键词：</span>
                            <div className="flex flex-wrap gap-1 mt-1">
                              {profection.houseKeywords.map((keyword, i) => (
                                <span
                                  key={i}
                                  className="px-2 py-0.5 bg-white/5 rounded text-xs text-white/80"
                                >
                                  {keyword}
                                </span>
                              ))}
                            </div>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
                </div>
              </motion.div>
            )}

            {/* ==================== 预测 Tab ==================== */}
            {selectedTab === 'forecast' && (
              <motion.div
                key="forecast"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* 今日预测 */}
                {dailyForecast && (
                  <div>
                    <h3 className="text-xl font-medium text-white mb-4">☀️ 今日预测</h3>
                    <div className="grid md:grid-cols-5 gap-4 mb-4">
                      {/* 综合运势 */}
                      <ScoreCard5612XY
                        title="综合运势"
                        score={dailyForecast.overallScore}
                        size="lg"
                      />
                      
                      {/* 实时运势 - 新增 */}
                      <div className="glass-card p-4 relative overflow-hidden">
                        <div className="absolute top-0 right-0 w-16 h-16 bg-gradient-to-bl from-[#00D4FF]/20 to-transparent rounded-bl-full" />
                        <div className="flex items-center gap-2 mb-2">
                          <span className="text-lg">⚡</span>
                          <span className="text-sm text-white/60">实时运势</span>
                        </div>
                        {realtimeScore ? (
                          <>
                            <div className="flex items-baseline gap-1">
                              <span className="text-3xl font-bold" style={{
                                color: realtimeScore.score >= 80 ? '#4ADE80' 
                                     : realtimeScore.score >= 60 ? '#00D4FF' 
                                     : realtimeScore.score >= 40 ? '#FFE66D' 
                                     : '#FF6B9D'
                              }}>
                                {Math.round(realtimeScore.score)}
                              </span>
                              <span className="text-white/40 text-sm">/ 100</span>
                            </div>
                            <div className="text-xs text-white/40 mt-1">
                              更新于 {realtimeScore.time}
                            </div>
                            <div className="mt-2 h-1.5 bg-white/10 rounded-full overflow-hidden">
                              <motion.div
                                className="h-full rounded-full"
                                style={{
                                  background: realtimeScore.score >= 80 ? 'linear-gradient(90deg, #4ADE80, #22C55E)'
                                           : realtimeScore.score >= 60 ? 'linear-gradient(90deg, #00D4FF, #0EA5E9)'
                                           : realtimeScore.score >= 40 ? 'linear-gradient(90deg, #FFE66D, #EAB308)'
                                           : 'linear-gradient(90deg, #FF6B9D, #EF4444)',
                                }}
                                initial={{ width: 0 }}
                                animate={{ width: `${realtimeScore.score}%` }}
                                transition={{ duration: 0.8, ease: 'easeOut' }}
                              />
                            </div>
                          </>
                        ) : (
                          <div className="text-white/40 text-sm">加载中...</div>
                        )}
                      </div>
                      
                      {/* 五维度雷达图 */}
                      <div className="glass-card p-4 flex items-center justify-center">
                        <DimensionRadarChart5832XY
                          scores={dailyForecast.dimensions || {
                            career: 50,
                            relationship: 50,
                            health: 50,
                            finance: 50,
                            spiritual: 50,
                          }}
                          size={140}
                          showLabels={true}
                          showValues={false}
                        />
                      </div>
                      
                      {/* 五维度详情 */}
                      <div className="md:col-span-2">
                        <DimensionScoresCard5612XY
                          scores={dailyForecast.dimensions || {
                            career: 50,
                            relationship: 50,
                            health: 50,
                            finance: 50,
                            spiritual: 50,
                          }}
                          layout="horizontal"
                        />
                      </div>
                    </div>
                    <DailyForecastCard7821MN
                      forecast={dailyForecast}
                      isToday={true}
                      isExpanded={expandedForecast === dailyForecast.date}
                      onClick={() => setExpandedForecast(
                        expandedForecast === dailyForecast.date ? null : dailyForecast.date
                      )}
                    />
                  </div>
                )}

                {/* 本周预测 */}
                {weeklyForecast && (
                  <div>
                    <h3 className="text-xl font-medium text-white mb-4">📆 本周预测</h3>
                    <div className="glass-card p-4 mb-4">
                      <p className="text-white/80">{weeklyForecast.overallTheme}</p>
                      <div className="flex flex-wrap gap-4 mt-3 text-sm">
                        <div>
                          <span className="text-white/60">周综合分：</span>
                          <span className="text-cyan-400">{Math.round(weeklyForecast.overallScore)}</span>
                        </div>
                        {weeklyForecast.bestDaysFor?.relationship?.length > 0 && (
                          <div>
                            <span className="text-white/60">最佳关系日：</span>
                            <span className="text-green-400">
                              {weeklyForecast.bestDaysFor.relationship.slice(0, 2).join(', ')}
                            </span>
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {weeklyForecast.dailySummaries?.map((summary, index) => (
                        <motion.div
                          key={summary.date}
                          initial={{ opacity: 0, y: 10 }}
                          animate={{ opacity: 1, y: 0 }}
                          transition={{ delay: index * 0.05 }}
                          className="glass-card p-4"
                        >
                          <div className="flex justify-between items-center mb-2">
                            <span className="text-white/60 text-sm">{summary.dayOfWeek}</span>
                            <span className="text-cyan-400 font-bold">{Math.round(summary.overallScore)}</span>
                          </div>
                          <div className="text-white text-sm">{summary.date}</div>
                          <div className="text-white/60 text-xs mt-1">{summary.keyTheme}</div>
                        </motion.div>
                      ))}
                    </div>
                  </div>
                )}

                {/* 多粒度运势查询 */}
                {birthData && (
                  <MultiGranularityScoreViewer8475QR
                    birthData={birthData}
                    className="mt-8"
                  />
                )}
              </motion.div>
            )}

            {/* ==================== 趋势 Tab ==================== */}
            {selectedTab === 'trend' && (
              <motion.div
                key="trend"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* 分数组成浮窗（点击趋势图数据点触发） */}
                <ScoreBreakdownPopup5932MN
                  open={breakdownOpen}
                  position={breakdownPosition}
                  queryTime={breakdownQueryTime}
                  loading={breakdownLoading}
                  error={breakdownError}
                  data={breakdownData}
                  activeFactorsData={activeFactorsData}
                  granularity={breakdownGranularity}
                  dimension={breakdownDimension}
                  onClose={() => setBreakdownOpen(false)}
                />

                {/* 新增：多粒度趋势图 */}
                <div className="glass-card p-6">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-xl font-medium text-white">📊 多粒度趋势分析</h3>
                    <GranularitySelector4721VW
                      value={trendGranularity}
                      onChange={setTrendGranularity}
                    />
                  </div>
                  
                  {/* 交互式时间序列趋势图 - 支持缩放、拖拽、动态Y轴、五维度切换 */}
                  {trendGranularity === 'yearly' && lifeTrend && lifeTrend.points && lifeTrend.points.length > 0 ? (
                    <InteractiveTrendChart9823EF
                      data={lifeTrend.points.map(p => {
                        // 将年龄转换为实际年份的时间戳（lightweight-charts 需要真实时间）
                        const year = (birthData?.year ?? 1990) + p.age;
                        return {
                          time: `${year}-06-15T12:00:00+08:00`, // 使用年中作为该年的代表时间点
                          value: p.overallScore,
                          label: `${p.age}岁 (${year}年)`,
                          dimensions: p.dimensions,
                        };
                      })}
                      title={`生命趋势 (当前: ${currentAge}岁)`}
                      color="#A855F7"
                      height={320}
                      showDimensions={true}
                      className="bg-white/5 rounded-lg"
                      onPointClick={handleTrendPointClick}
                    />
                  ) : timeSeries && timeSeries.points && timeSeries.points.length > 0 ? (
                    (() => {
                      // 根据粒度选择颜色
                      const colorMap: Record<string, string> = {
                        hourly: '#00D4FF',   // 青色
                        daily: '#4ECDC4',    // 绿色
                        weekly: '#FFE66D',   // 黄色
                        monthly: '#FF9F43',  // 橙色
                        yearly: '#A855F7',   // 紫色
                      };
                      const color = colorMap[trendGranularity] || '#00D4FF';
                      const granularityLabel = { hourly: '小时', daily: '天', weekly: '周', monthly: '月', yearly: '年' }[trendGranularity];
                      
                      return (
                        <InteractiveTrendChart9823EF
                          data={timeSeries.points.map(p => ({
                            time: p.time,
                            value: p.display,
                            label: p.label,
                            dimensions: p.dimensions, // 传递五维度数据
                          }))}
                          title={`${granularityLabel}趋势 (${timeSeries.points.length}个数据点)`}
                          color={color}
                          height={320}
                          showDimensions={true}
                          className="bg-white/5 rounded-lg"
                          onVisibleRangeChange={handleVisibleRangeChange}
                          isLoading={isLoadingMoreData}
                          onPointClick={handleTrendPointClick}
                        />
                      );
                    })()
                  ) : loading ? (
                    <div className="h-64 flex items-center justify-center bg-white/5 rounded-lg">
                      <Spinner size="lg" />
                      <span className="ml-3 text-white/60">加载趋势数据...</span>
                    </div>
                  ) : trendGranularity === 'yearly' && !lifeTrend ? (
                    <div className="h-64 flex items-center justify-center bg-white/5 rounded-lg">
                      <Spinner size="lg" />
                      <span className="ml-3 text-white/60">加载生命趋势数据...</span>
                    </div>
                  ) : (
                    <div className="h-64 flex items-center justify-center bg-white/5 rounded-lg border border-dashed border-white/20">
                      <div className="text-center text-white/50">
                        <p className="text-lg mb-2">📈 多粒度趋势图</p>
                        <p className="text-sm">当前粒度: {
                          { hourly: '小时', daily: '日', weekly: '周', monthly: '月', yearly: '年' }[trendGranularity]
                        }</p>
                        <p className="text-xs mt-2">暂无数据</p>
                      </div>
                    </div>
                  )}
                </div>

                {/* 生命趋势图 */}
                {lifeTrend ? (
                  <div>
                    <h3 className="text-xl font-medium text-white mb-4">📈 生命趋势 (0-80岁)</h3>
                    <LifeTimeline4529PQ
                      data={lifeTrend}
                      currentAge={currentAge}
                      height={280}
                      showDimensions={true}
                      onPointClick={(point) => {
                        console.log('点击年龄:', point.age, point);
                      }}
                    />
                  </div>
                ) : (
                  <div className="glass-card p-6 flex items-center justify-center">
                    <Spinner size="lg" />
                    <span className="ml-3 text-white/60">加载生命趋势...</span>
                  </div>
                )}

                {/* 年限法轮盘 */}
                {profectionMap ? (
                  <ProfectionWheel6183RS
                    profections={profectionMap.profections}
                    currentAge={currentAge}
                    size={350}
                    onAgeClick={(age) => {
                      console.log('点击年龄:', age);
                    }}
                  />
                ) : (
                  <div className="glass-card p-6 flex items-center justify-center">
                    <Spinner size="lg" />
                    <span className="ml-3 text-white/60">加载年限法...</span>
                  </div>
                )}

                {/* 重大行运提示 */}
                <div className="glass-card p-4">
                  <h3 className="text-lg font-medium text-white mb-3">🌟 重大行运节点</h3>
                  <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-3 text-sm">
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#A855F7] font-medium">♄ 土星回归</div>
                      <div className="text-white/60">29-30 岁 / 58-60 岁</div>
                      <div className="text-white/40 text-xs mt-1">人生结构升级</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#00D4FF] font-medium">♅ 天王星对冲</div>
                      <div className="text-white/60">40-42 岁</div>
                      <div className="text-white/40 text-xs mt-1">中年觉醒</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#ffd700] font-medium">♃ 木星回归</div>
                      <div className="text-white/60">12 / 24 / 36 / 48 岁</div>
                      <div className="text-white/40 text-xs mt-1">扩张机遇</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#FF6B9D] font-medium">☊ 北交点回归</div>
                      <div className="text-white/60">18-19 / 37-38 岁</div>
                      <div className="text-white/40 text-xs mt-1">命运节点</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#ff8c00] font-medium">⚷ 凯龙回归</div>
                      <div className="text-white/60">50-51 岁</div>
                      <div className="text-white/40 text-xs mt-1">伤痛治愈</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#4169e1] font-medium">♆ 海王星四分</div>
                      <div className="text-white/60">41 岁</div>
                      <div className="text-white/40 text-xs mt-1">灵性转化</div>
                    </div>
                  </div>
                </div>
              </motion.div>
            )}

            {/* ==================== 因子 Tab ==================== */}
            {selectedTab === 'factors' && (
              <motion.div
                key="factors"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* 编辑模式开关 */}
                <div className="flex items-center justify-between">
                  <h3 className="text-xl font-medium text-white">📊 影响因子分析</h3>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-white/60">编辑模式</span>
                    <Switch
                      isSelected={showFactorEditor}
                      onValueChange={setShowFactorEditor}
                      size="sm"
                    />
                  </div>
                </div>

                <div className="grid lg:grid-cols-2 gap-6">
                  {/* 左侧：当前因子列表 */}
                  <div>
                    <InfluenceFactorsPanel8274TU
                      factors={MOCK_INFLUENCE_FACTORS}
                      editable={showFactorEditor}
                      onWeightChange={(name, weight) => {
                        console.log('权重变更:', name, weight);
                      }}
                    />
                  </div>

                  {/* 右侧：自定义因子编辑器 */}
                  <div>
                    <CustomFactorEditor9456DE
                      factors={customFactors}
                      onAdd={handleAddCustomFactor}
                      onRemove={handleRemoveCustomFactor}
                    />
                  </div>
                </div>

                {/* 因子说明 */}
                <div className="glass-card p-4">
                  <h4 className="text-lg font-medium text-white mb-3">📖 因子权重说明</h4>
                  <div className="grid md:grid-cols-3 gap-4 text-sm">
                    <div>
                      <div className="text-white/80 font-medium mb-2">尊贵度 (Dignity)</div>
                      <ul className="text-white/60 space-y-1 list-disc list-inside">
                        <li>入庙 (Domicile): +3</li>
                        <li>旺相 (Exaltation): +2</li>
                        <li>落陷 (Detriment): -2</li>
                        <li>失势 (Fall): -3</li>
                      </ul>
                    </div>
                    <div>
                      <div className="text-white/80 font-medium mb-2">其他因子</div>
                      <ul className="text-white/60 space-y-1 list-disc list-inside">
                        <li>逆行: -2</li>
                        <li>相位阶段: ×0.8</li>
                        <li>外行星放大: ×1.2</li>
                        <li>年主星加成: +1.0</li>
                      </ul>
                    </div>
                    <div>
                      <div className="text-white/80 font-medium mb-2">时间级别</div>
                      <ul className="text-white/60 space-y-1 list-disc list-inside">
                        <li>年度级: 土星回归、木星回归</li>
                        <li>月度级: 太阳换座、月相</li>
                        <li>周度级: 水星逆行</li>
                        <li>日度级: 月亮换座</li>
                        <li>小时级: 行星时、月空</li>
                      </ul>
                    </div>
                  </div>
                </div>

                {/* 自定义因子格式说明 */}
                <div className="glass-card p-4">
                  <h4 className="text-lg font-medium text-white mb-3">💡 自定义因子格式</h4>
                  <div className="text-sm text-white/60">
                    <p className="mb-2">格式: <code className="text-cosmic-nova bg-black/30 px-2 py-0.5 rounded">Operation=(value*dimension,duration,startTime)</code></p>
                    <p className="mb-2">示例: <code className="text-green-400 bg-black/30 px-2 py-0.5 rounded">AddScore=(2*healthScore,2.5,202517301212)</code></p>
                    <p>含义: 从2025年1月17日30分12秒开始，健康值+2，持续2.5小时</p>
                  </div>
                </div>
              </motion.div>
            )}

            {/* ==================== 设置 Tab ==================== */}
            {selectedTab === 'settings' && (
              <motion.div
                key="settings"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="max-w-2xl space-y-6"
              >
                <div className="glass-card p-6">
                  <h3 className="text-xl font-medium text-white mb-4">👤 当前出生数据</h3>
                  {birthData && (
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <span className="text-white/60">出生日期：</span>
                        <span className="text-white">
                          {birthData.year}年{birthData.month}月{birthData.day}日
                        </span>
                      </div>
                      <div>
                        <span className="text-white/60">出生时间：</span>
                        <span className="text-white">
                          {String(birthData.hour).padStart(2, '0')}:{String(birthData.minute).padStart(2, '0')}
                        </span>
                      </div>
                      <div>
                        <span className="text-white/60">出生地点：</span>
                        <span className="text-white">
                          {birthData.latitude.toFixed(4)}°, {birthData.longitude.toFixed(4)}°
                        </span>
                      </div>
                      <div>
                        <span className="text-white/60">时区：</span>
                        <span className="text-white">UTC{birthData.timezone >= 0 ? '+' : ''}{birthData.timezone}</span>
                      </div>
                    </div>
                  )}
                  <Button
                    className="mt-4"
                    variant="flat"
                    onPress={() => window.location.reload()}
                  >
                    🔄 重新输入
                  </Button>
                </div>

                <div className="glass-card p-6">
                  <h3 className="text-xl font-medium text-white mb-4">🔧 系统信息</h3>
                  <div className="text-sm space-y-2 text-white/60">
                    <p>版本: 1.0.0</p>
                    <p>算法: VSOP87 简化模型 / Placidus 分宫制</p>
                    <p>精度: 行星经度 &lt;1° / 太阳 &lt;0.1°</p>
                    <p>数据来源: 内置天文算法计算</p>
                  </div>
                </div>

                <div className="glass-card p-6">
                  <h3 className="text-xl font-medium text-white mb-4">⚠️ 免责声明</h3>
                  <p className="text-sm text-white/60">
                    本平台基于天文算法进行占星学计算，所有数据均有理论支撑，
                    仅供研究和学习参考。预测结果不构成任何决策建议，
                    请理性看待占星学分析结果。
                  </p>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      )}

      {/* 页脚 */}
      <footer className="text-center text-white/30 text-sm mt-12 space-y-4">
        <div>
          <p>Star 占星计算验证平台 v1.0.0</p>
          <p className="mt-1">数据基于天文算法计算，仅供研究参考</p>
        </div>
        
        {/* MiniApp 入口 */}
        <div className="pt-4 flex justify-center">
          <Button
            variant="flat"
            className="bg-[#24A1DE]/10 hover:bg-[#24A1DE]/20 text-[#24A1DE] border border-[#24A1DE]/30"
            onPress={() => setShowMiniApp(true)}
          >
            <span className="mr-2">✈️</span>
            进入Telegram MiniAPP （演示模式）
          </Button>
        </div>
      </footer>
    </div>
  );
}

export default App;
