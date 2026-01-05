/**
 * InteractiveTrendChart9823EF - 交互式趋势图
 * 类似股票K线的交互体验：
 * - 动态Y轴（自适应数据范围）
 * - 拖拽平移
 * - 缩放
 * - 十字准线
 * - 支持五维度数据展示
 */

import { useEffect, useRef, useCallback, useState } from 'react';
import { createChart, AreaSeries } from 'lightweight-charts';
import type { IChartApi, ISeriesApi, AreaData, Time } from 'lightweight-charts';
import type { DimensionScores, Dimension } from '../../types';
import { DIMENSION_NAMES, DIMENSION_COLORS } from '../../utils/astro';

interface TrendDataPoint {
  time: string; // ISO 时间或时间戳
  value: number;
  label?: string;
  dimensions?: DimensionScores; // 新增：五维度数据
}

// 可视范围变化事件（导出供外部使用）
export interface VisibleRangeChange {
  from: Date;
  to: Date;
  needsMoreBefore: boolean; // 是否需要加载左侧更多数据
  needsMoreAfter: boolean;  // 是否需要加载右侧更多数据
}

interface InteractiveTrendChartProps {
  data: TrendDataPoint[];
  title?: string;
  color?: string;
  height?: number;
  showVolume?: boolean;
  showDimensions?: boolean; // 新增：是否显示维度选择器
  className?: string;
  onPointHover?: (point: TrendDataPoint | null) => void;
  onPointClick?: (point: TrendDataPoint, dimension: Dimension | 'overall', event?: MouseEvent) => void;
  onVisibleRangeChange?: (range: VisibleRangeChange) => void; // 新增：可视范围变化回调
  isLoading?: boolean; // 新增：是否正在加载更多数据
}

// 颜色主题
const CHART_COLORS = {
  background: 'transparent',
  textColor: 'rgba(255, 255, 255, 0.6)',
  gridColor: 'rgba(255, 255, 255, 0.1)',
  crosshairColor: 'rgba(255, 255, 255, 0.5)',
};

// 维度配置
const DIMENSION_OPTIONS: Array<{ id: Dimension | 'overall'; label: string; color: string }> = [
  { id: 'overall', label: '综合', color: '#00D4FF' },
  { id: 'career', label: '事业', color: DIMENSION_COLORS.career },
  { id: 'relationship', label: '关系', color: DIMENSION_COLORS.relationship },
  { id: 'health', label: '健康', color: DIMENSION_COLORS.health },
  { id: 'finance', label: '财务', color: DIMENSION_COLORS.finance },
  { id: 'spiritual', label: '灵性', color: DIMENSION_COLORS.spiritual },
];

// timeScale 右侧留白（影响 logicalRange.to 的判断，需要统一扣除）
const TIME_SCALE_RIGHT_OFFSET = 5;
// 避免初始化/fitContent/数据刷新时误触发加载更多
const IGNORE_RANGE_CHANGE_MS = 500;

export function InteractiveTrendChart9823EF({
  data,
  title,
  color = '#00D4FF',
  height = 300,
  showDimensions = false,
  className = '',
  onPointHover,
  onPointClick,
  onVisibleRangeChange,
  isLoading = false,
}: InteractiveTrendChartProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const seriesRef = useRef<ISeriesApi<'Area'> | null>(null);
  const tooltipRef = useRef<HTMLDivElement>(null);
  const onVisibleRangeChangeRef = useRef<typeof onVisibleRangeChange>(onVisibleRangeChange);
  const onPointClickRef = useRef<typeof onPointClick>(onPointClick);
  const [selectedDimension, setSelectedDimension] = useState<Dimension | 'overall'>('overall');
  const selectedDimensionRef = useRef<Dimension | 'overall'>(selectedDimension); // 用于点击回调
  const lastRangeCheckRef = useRef<number>(0); // 防抖：上次检查时间
  const lastClickEventRef = useRef<MouseEvent | null>(null); // 记录最后一次点击的 MouseEvent
  const dataRangeRef = useRef<{ minTime: number; maxTime: number } | null>(null); // 当前数据范围
  const dataLengthRef = useRef<number>(0); // 当前数据长度（避免闭包问题）
  const ignoreRangeChangeUntilRef = useRef<number>(0); // 忽略可视范围变化回调直到某个时间点
  const hasUserInteractedRef = useRef<boolean>(false); // 用户是否已经进行缩放/拖拽等交互（避免在交互中 fitContent 导致断言错误）

  // 保持可视范围回调为最新引用（避免 chart 初始化 effect 不重跑导致回调闭包陈旧）
  useEffect(() => {
    onVisibleRangeChangeRef.current = onVisibleRangeChange;
  }, [onVisibleRangeChange]);

  useEffect(() => {
    onPointClickRef.current = onPointClick;
  }, [onPointClick]);

  // 同步 selectedDimension 到 ref（供点击回调使用）
  useEffect(() => {
    selectedDimensionRef.current = selectedDimension;
  }, [selectedDimension]);

  // 时间戳（秒） -> 原始数据点（用于 hover/click 定位）
  const timeToPointRef = useRef<Map<number, TrendDataPoint>>(new Map());
  useEffect(() => {
    const map = new Map<number, TrendDataPoint>();
    for (const p of data) {
      // 所有时间应该都是 ISO 格式，使用 Unix 秒数作为 key
      const ms = Date.parse(p.time);
      if (Number.isFinite(ms)) {
        map.set(Math.floor(ms / 1000), p);
      }
    }
    timeToPointRef.current = map;
  }, [data]);

  // 用户交互状态：正在交互时跳过可能冲突的操作
  const isInteractingRef = useRef(false);
  const interactingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  
  // 监听用户交互
  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    
    const startInteraction = () => {
      hasUserInteractedRef.current = true;
      isInteractingRef.current = true;
      // 清除之前的超时
      if (interactingTimeoutRef.current) {
        clearTimeout(interactingTimeoutRef.current);
      }
    };
    
    const endInteraction = () => {
      // 延迟结束交互状态，避免快速连续操作
      if (interactingTimeoutRef.current) {
        clearTimeout(interactingTimeoutRef.current);
      }
      interactingTimeoutRef.current = setTimeout(() => {
        isInteractingRef.current = false;
      }, 300);
    };
    
    el.addEventListener('wheel', startInteraction, { passive: true });
    el.addEventListener('mousedown', startInteraction);
    el.addEventListener('touchstart', startInteraction, { passive: true });
    el.addEventListener('mouseup', endInteraction);
    el.addEventListener('mouseleave', endInteraction);
    el.addEventListener('touchend', endInteraction);
    
    return () => {
      el.removeEventListener('wheel', startInteraction);
      el.removeEventListener('mousedown', startInteraction);
      el.removeEventListener('touchstart', startInteraction);
      el.removeEventListener('mouseup', endInteraction);
      el.removeEventListener('mouseleave', endInteraction);
      el.removeEventListener('touchend', endInteraction);
      if (interactingTimeoutRef.current) {
        clearTimeout(interactingTimeoutRef.current);
      }
    };
  }, []);

  // 获取当前选中维度的颜色
  const currentColor = showDimensions 
    ? DIMENSION_OPTIONS.find(d => d.id === selectedDimension)?.color || color
    : color;

  // 根据选中维度获取数据值
  const getValueForDimension = useCallback((point: TrendDataPoint): number => {
    if (selectedDimension === 'overall' || !point.dimensions) {
      return point.value;
    }
    return point.dimensions[selectedDimension] || point.value;
  }, [selectedDimension]);

  // 转换数据格式
  const convertData = useCallback((rawData: TrendDataPoint[]): AreaData[] => {
    // 首先按时间排序（lightweight-charts 要求数据必须时间升序）
    const sortedData = [...rawData].sort((a, b) => {
      // 仅在“看起来是日期/时间”的情况下解析，避免年龄等字符串被 Date.parse 误解析
      const timeA = (a.time.includes('-') || a.time.includes('T')) ? (Date.parse(a.time) || 0) : 0;
      const timeB = (b.time.includes('-') || b.time.includes('T')) ? (Date.parse(b.time) || 0) : 0;
      return timeA - timeB;
    });
    
    // 去重：同一时间点只保留最后一个（避免重复 time 造成内部范围异常）
    const deduped: TrendDataPoint[] = [];
    const seen = new Map<number, TrendDataPoint>();
    for (const p of sortedData) {
      if (p.time.includes('-') || p.time.includes('T')) {
        const ms = Date.parse(p.time);
        if (Number.isFinite(ms)) {
          seen.set(ms, p);
          continue;
        }
      }
      // 非可解析时间（例如年龄、纯文本），直接保留
      deduped.push(p);
    }
    if (seen.size > 0) {
      // 将可解析的时间点按时间升序拼接到前面
      const byTime = Array.from(seen.entries())
        .sort((a, b) => a[0] - b[0])
        .map(([, p]) => p);
      deduped.unshift(...byTime);
    }

    return deduped.map((point, index) => {
      // 尝试解析时间
      let time: Time;
      
      if (point.time.includes('T') || point.time.includes('-')) {
        // 统一使用 Date.parse 处理时区（后端返回 RFC3339，浏览器可直接解析）
        const ms = Date.parse(point.time);
        time = Number.isFinite(ms) ? (Math.floor(ms / 1000) as Time) : (index as Time);
      } else if (point.time.includes('岁')) {
        // 年龄格式 - 使用索引作为时间
        time = index as Time;
      } else if (/^\d{1,2}:\d{2}$/.test(point.time)) {
        // 时间格式如 "11:00" - 使用索引
        time = index as Time;
      } else {
        // 其他格式 - 使用索引
        time = index as Time;
      }
      
      // 使用选中维度的值
      const value = getValueForDimension(point);
      return { time, value };
    });
  }, [getValueForDimension]);

  // 初始化图表
  useEffect(() => {
    if (!containerRef.current) return;

    // 创建图表
    const chart = createChart(containerRef.current, {
      width: containerRef.current.clientWidth,
      height,
      layout: {
        background: { color: CHART_COLORS.background },
        textColor: CHART_COLORS.textColor,
      },
      grid: {
        vertLines: { color: CHART_COLORS.gridColor },
        horzLines: { color: CHART_COLORS.gridColor },
      },
      crosshair: {
        mode: 1, // Normal
        vertLine: {
          color: CHART_COLORS.crosshairColor,
          width: 1,
          style: 2, // Dashed
          labelBackgroundColor: color,
        },
        horzLine: {
          color: CHART_COLORS.crosshairColor,
          width: 1,
          style: 2,
          labelBackgroundColor: color,
        },
      },
      rightPriceScale: {
        borderColor: CHART_COLORS.gridColor,
        scaleMargins: {
          top: 0.1,
          bottom: 0.1,
        },
        autoScale: true, // 动态Y轴
      },
      timeScale: {
        borderColor: CHART_COLORS.gridColor,
        timeVisible: true,
        secondsVisible: false,
        rightOffset: TIME_SCALE_RIGHT_OFFSET,
        barSpacing: 10,
        minBarSpacing: 2,
        fixLeftEdge: false,
        fixRightEdge: false,
      },
      handleScroll: {
        mouseWheel: true,
        pressedMouseMove: true,
        horzTouchDrag: true,
        vertTouchDrag: false,
      },
      handleScale: {
        axisPressedMouseMove: true,
        mouseWheel: true,
        pinch: true,
      },
    });

    // 创建面积图系列 (lightweight-charts v5+ API)
    const areaSeries = chart.addSeries(AreaSeries, {
      lineColor: currentColor,
      topColor: `${currentColor}80`,
      bottomColor: `${currentColor}10`,
      lineWidth: 2,
      crosshairMarkerVisible: true,
      crosshairMarkerRadius: 6,
      crosshairMarkerBorderColor: '#ffffff',
      crosshairMarkerBackgroundColor: currentColor,
    });

    chartRef.current = chart;
    seriesRef.current = areaSeries;

    // 设置数据
    if (data.length > 0) {
      const chartData = convertData(data);
      areaSeries.setData(chartData);
      // 初始渲染时才 fitContent，避免与用户缩放/拖拽并发导致断言错误
      if (!hasUserInteractedRef.current) {
        try {
          chart.timeScale().fitContent();
        } catch {
          // ignore fitContent errors
        }
        ignoreRangeChangeUntilRef.current = Date.now() + IGNORE_RANGE_CHANGE_MS;
      }
      
      // 记录当前数据范围和长度
      dataLengthRef.current = chartData.length;
      if (chartData.length > 0) {
        // 时间已统一为 UTCTimestamp（秒）
        const times = chartData
          .map(d => (typeof d.time === 'number' ? d.time : 0))
          .filter(t => Number.isFinite(t) && t > 0);
        if (times.length > 0) {
          dataRangeRef.current = {
            minTime: Math.min(...times),
            maxTime: Math.max(...times),
          };
        }
      }
    }

    // 订阅可视范围变化事件 - 用于自动加载更多数据
    chart.timeScale().subscribeVisibleLogicalRangeChange((logicalRange) => {
      const handler = onVisibleRangeChangeRef.current;
      if (!handler) return;
      if (!logicalRange || !dataRangeRef.current) return;
      if (Date.now() < ignoreRangeChangeUntilRef.current) return;

      // 防抖：300ms 内不重复触发（避免频繁请求）
      const now = Date.now();
      if (now - lastRangeCheckRef.current < 300) return;
      lastRangeCheckRef.current = now;

      const { minTime, maxTime } = dataRangeRef.current;
      const dataLength = dataLengthRef.current;
      if (dataLength <= 0) return;

      // 计算可视范围边界（逻辑索引）
      const visibleFrom = logicalRange.from;
      const visibleTo = logicalRange.to;
      // 注意：rightOffset 会让 visibleTo 在“正常展示”时也偏大，需要扣除
      const effectiveVisibleTo = visibleTo - TIME_SCALE_RIGHT_OFFSET;
      // 边界缓冲：滚到边缘“附近”就触发，避免必须拖到负数/超长才触发
      const edgeBuffer = Math.max(3, Math.min(12, Math.floor(dataLength * 0.08)));

      // 检查是否需要加载更多数据
      // 左边缘：可视范围靠近数据起点
      const needsMoreBefore = visibleFrom < edgeBuffer;
      // 右边缘：可视范围靠近数据终点（扣除 rightOffset 影响）
      const needsMoreAfter = effectiveVisibleTo > (dataLength - 1 - edgeBuffer);

      if (needsMoreBefore || needsMoreAfter) {
        // 计算实际时间范围
        const fromDate = new Date(minTime * 1000);
        const toDate = new Date(maxTime * 1000);

        handler({
          from: fromDate,
          to: toDate,
          needsMoreBefore,
          needsMoreAfter,
        });
      }
    });

    // 订阅十字准线移动事件
    chart.subscribeCrosshairMove((param) => {
      if (param.time && param.seriesData.size > 0) {
        const price = param.seriesData.get(areaSeries);
        if (price && tooltipRef.current) {
          const t = typeof param.time === 'number' ? param.time : 0;
          const point = timeToPointRef.current.get(t);
          if (point) {
            tooltipRef.current.style.display = 'block';
            const displayValue = getValueForDimension(point);
            const dimensionLabel = selectedDimension === 'overall' ? '综合' : DIMENSION_NAMES[selectedDimension];
            // 格式化时间显示（使用 UTC，与图表底部保持一致）
            let timeDisplay = point.label || point.time;
            if (point.time.includes('T')) {
              const d = new Date(point.time);
              if (!isNaN(d.getTime())) {
                const h = d.getUTCHours().toString().padStart(2, '0');
                const m = d.getUTCMinutes().toString().padStart(2, '0');
                timeDisplay = `${h}:${m}`;
              }
            }
            tooltipRef.current.innerHTML = `
              <div class="font-bold" style="color: ${currentColor}">${displayValue.toFixed(1)}</div>
              <div class="text-xs text-white/60">${timeDisplay}</div>
              ${showDimensions ? `<div class="text-xs text-white/40">${dimensionLabel}</div>` : ''}
            `;
            onPointHover?.(point);
          }
        }
      } else {
        if (tooltipRef.current) {
          tooltipRef.current.style.display = 'none';
        }
        onPointHover?.(null);
      }
    });

    // 订阅点击事件（点击数据点弹出详情）
    // 捕获容器上的原生点击事件（用于获取鼠标位置）
    const handleContainerClick = (e: MouseEvent) => {
      lastClickEventRef.current = e;
    };
    containerRef.current?.addEventListener('click', handleContainerClick);

    chart.subscribeClick((param) => {
      const handler = onPointClickRef.current;
      if (!handler) return;
      if (!param.time) return;
      const t = typeof param.time === 'number' ? param.time : 0;
      const point = timeToPointRef.current.get(t);
      if (point) {
        handler(point, selectedDimensionRef.current, lastClickEventRef.current ?? undefined);
      }
    });

    // 监听窗口大小变化
    const handleResize = () => {
      if (containerRef.current && chartRef.current) {
        chartRef.current.applyOptions({
          width: containerRef.current.clientWidth,
        });
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      containerRef.current?.removeEventListener('click', handleContainerClick);
      chart.remove();
      chartRef.current = null;
      seriesRef.current = null;
    };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [height, currentColor, selectedDimension, showDimensions]);

  // 更新数据（保持当前视图位置，不自动 fitContent）
  useEffect(() => {
    if (seriesRef.current && data.length > 0) {
      const chartData = convertData(data);
      
      // 更新数据 - 始终执行，不管用户是否在交互
      try {
        seriesRef.current.setData(chartData);
      } catch {
        // ignore setData errors
        return;
      }
      
      // 更新数据范围和长度
      const previousDataRange = dataRangeRef.current;
      dataLengthRef.current = chartData.length;
      
      if (chartData.length > 0) {
        const times = chartData
          .map(d => (typeof d.time === 'number' ? d.time : 0))
          .filter(t => Number.isFinite(t));
        if (times.length > 0) {
          dataRangeRef.current = {
            minTime: Math.min(...times),
            maxTime: Math.max(...times),
          };
        }
      }
      
      // 仅在用户未交互过时，允许根据数据范围变化进行一次适配
      // 用户交互后，让图表自己管理范围，避免冲突
      if (!hasUserInteractedRef.current) {
        if (!previousDataRange ||
            previousDataRange.minTime !== dataRangeRef.current?.minTime ||
            previousDataRange.maxTime !== dataRangeRef.current?.maxTime) {
          // 延迟执行 fitContent，避免和当前渲染冲突
          requestAnimationFrame(() => {
            try {
              chartRef.current?.timeScale().fitContent();
            } catch {
              // ignore fitContent errors
            }
          });
          ignoreRangeChangeUntilRef.current = Date.now() + IGNORE_RANGE_CHANGE_MS;
        }
      }
    }
  }, [data, convertData]);

  return (
    <div className={`relative ${className}`}>
      {/* 标题栏 */}
      <div className="flex items-center justify-between mb-2 px-2">
        {/* 标题 */}
        {title && (
          <div className="text-sm font-medium text-white/80">
            {title}
          </div>
        )}
        
        {/* 维度选择器 */}
        {showDimensions && (
          <div className="flex gap-1">
            {DIMENSION_OPTIONS.map((dim) => (
              <button
                key={dim.id}
                onClick={() => setSelectedDimension(dim.id)}
                className={`px-2 py-1 text-xs rounded-lg transition-all ${
                  selectedDimension === dim.id
                    ? 'text-white'
                    : 'text-white/50 hover:text-white/70'
                }`}
                style={{
                  backgroundColor: selectedDimension === dim.id ? `${dim.color}30` : 'transparent',
                  border: selectedDimension === dim.id ? `1px solid ${dim.color}` : '1px solid transparent',
                }}
              >
                {dim.label}
              </button>
            ))}
          </div>
        )}

        {/* 提示框 */}
        <div
          ref={tooltipRef}
          className="glass-card px-3 py-2 hidden"
        />
      </div>
      
      {/* 图表容器 */}
      <div className="relative">
        <div
          ref={containerRef}
          className="w-full rounded-lg overflow-hidden"
          style={{ height }}
        />
        
        {/* 加载中指示器 */}
        {isLoading && (
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div className="flex items-center gap-2 bg-black/60 backdrop-blur-sm px-4 py-2 rounded-lg">
              <div className="w-4 h-4 border-2 border-white/30 border-t-white/80 rounded-full animate-spin" />
              <span className="text-xs text-white/80">加载更多数据...</span>
            </div>
          </div>
        )}
      </div>
      
      {/* 操作提示 */}
      <div className="absolute bottom-2 left-2 z-10 text-xs text-white/40">
        拖拽平移 | 滚轮缩放 | 双击重置
      </div>
    </div>
  );
}

export default InteractiveTrendChart9823EF;

