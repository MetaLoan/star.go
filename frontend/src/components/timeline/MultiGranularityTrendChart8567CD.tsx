/**
 * MultiGranularityTrendChart8567CD - 多粒度趋势图
 * 支持 小时/天/周/月/年 多级切换查看，类似K线图
 * 支持因子级别筛选、因子详情查看、时间范围选择
 */

import { useState, useMemo, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { GranularitySelector4721VW, type TimeGranularity } from '../ui/GranularitySelector4721VW';
import { FactorLevelFilter6293ZA, type FactorTimeLevel } from '../factors/FactorLevelFilter6293ZA';
import { DimensionRadarChart5832XY } from '../chart/DimensionRadarChart5832XY';
import { FactorLifecycleChart7184BC } from '../factors/FactorLifecycleChart7184BC';
import type { DimensionScores, InfluenceFactor } from '../../types';

// 趋势点数据
interface TrendDataPoint {
  time: string;
  timestamp: number;
  overallScore: number;
  dimensions: DimensionScores;
  factors: InfluenceFactor[];
  isMajorEvent?: boolean;
  eventLabel?: string;
}

interface MultiGranularityTrendChartProps {
  data: TrendDataPoint[];
  granularity: TimeGranularity;
  onGranularityChange: (g: TimeGranularity) => void;
  // selectedTime?: Date; // 预留：选中时间点高亮
  onTimeSelect?: (time: Date) => void;
  className?: string;
  height?: number;
  showFactorPanel?: boolean;
}

// 颜色配置
const COLORS = {
  overall: '#00D4FF',
  career: '#FF6B6B',
  relationship: '#FF9F43',
  health: '#4ECDC4',
  finance: '#4FC3F7',
  spiritual: '#A855F7',
  grid: 'rgba(255,255,255,0.1)',
  majorEvent: '#FF6B9D',
};

export function MultiGranularityTrendChart8567CD({
  data,
  granularity,
  onGranularityChange,
  onTimeSelect,
  className = '',
  height = 300,
  showFactorPanel = true,
}: MultiGranularityTrendChartProps) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);
  const [selectedFactorLevels, setSelectedFactorLevels] = useState<FactorTimeLevel[]>([
    'yearly', 'monthly', 'weekly', 'daily', 'hourly',
  ]);
  const [showDimensions, setShowDimensions] = useState(false);
  const [expandedFactor, setExpandedFactor] = useState<InfluenceFactor | null>(null);

  // 图表尺寸
  const chartPadding = { top: 30, right: 20, bottom: 40, left: 50 };
  const chartWidth = 800;
  const chartHeight = height;

  // 计算数据范围
  const { minScore, maxScore, dataPoints } = useMemo(() => {
    if (data.length === 0) {
      return { minScore: 0, maxScore: 100, dataPoints: [] };
    }

    let min = 100, max = 0;
    data.forEach(d => {
      if (d.overallScore < min) min = d.overallScore;
      if (d.overallScore > max) max = d.overallScore;
    });

    // 添加一些边距
    const padding = (max - min) * 0.1 || 10;
    min = Math.max(0, min - padding);
    max = Math.min(100, max + padding);

    return { minScore: min, maxScore: max, dataPoints: data };
  }, [data]);

  // 坐标转换
  const getX = useCallback((index: number) => {
    const usableWidth = chartWidth - chartPadding.left - chartPadding.right;
    return chartPadding.left + (index / (dataPoints.length - 1 || 1)) * usableWidth;
  }, [dataPoints.length]);

  const getY = useCallback((score: number) => {
    const usableHeight = chartHeight - chartPadding.top - chartPadding.bottom;
    const normalized = (score - minScore) / (maxScore - minScore || 1);
    return chartPadding.top + usableHeight * (1 - normalized);
  }, [minScore, maxScore, chartHeight]);

  // 生成趋势线路径
  const trendPath = useMemo(() => {
    if (dataPoints.length === 0) return '';
    return dataPoints
      .map((d, i) => `${i === 0 ? 'M' : 'L'} ${getX(i)} ${getY(d.overallScore)}`)
      .join(' ');
  }, [dataPoints, getX, getY]);

  // 生成填充区域
  const fillPath = useMemo(() => {
    if (dataPoints.length === 0) return '';
    const baseline = chartHeight - chartPadding.bottom;
    return trendPath + ` L ${getX(dataPoints.length - 1)} ${baseline} L ${getX(0)} ${baseline} Z`;
  }, [trendPath, dataPoints.length, getX, chartHeight]);

  // 当前悬停的数据点
  const hoveredPoint = hoveredIndex !== null ? dataPoints[hoveredIndex] : null;

  // 过滤因子
  const filteredFactors = useMemo(() => {
    if (!hoveredPoint) return [];
    return hoveredPoint.factors.filter(f => {
      const level = f.timeLevel?.level as FactorTimeLevel;
      return selectedFactorLevels.includes(level);
    });
  }, [hoveredPoint, selectedFactorLevels]);

  // 格式化时间标签
  const formatTimeLabel = (time: string, g: TimeGranularity) => {
    const date = new Date(time);
    switch (g) {
      case 'hourly':
        return `${date.getHours()}:00`;
      case 'daily':
        return `${date.getMonth() + 1}/${date.getDate()}`;
      case 'weekly':
        return `第${Math.ceil(date.getDate() / 7)}周`;
      case 'monthly':
        return `${date.getMonth() + 1}月`;
      case 'yearly':
        return `${date.getFullYear()}`;
      default:
        return time;
    }
  };

  // 映射粒度到因子级别
  const granularityToFactorLevel: Record<TimeGranularity, FactorTimeLevel> = {
    hourly: 'hourly',
    daily: 'daily',
    weekly: 'weekly',
    monthly: 'monthly',
    yearly: 'yearly',
  };

  return (
    <div className={`${className}`}>
      {/* 顶部控制栏 */}
      <div className="flex flex-wrap items-center justify-between gap-4 mb-4">
        <GranularitySelector4721VW
          value={granularity}
          onChange={onGranularityChange}
        />
        
        <div className="flex items-center gap-4">
          <FactorLevelFilter6293ZA
            selected={selectedFactorLevels}
            onChange={setSelectedFactorLevels}
            viewLevel={granularityToFactorLevel[granularity]}
            compact
          />
          
          <button
            onClick={() => setShowDimensions(!showDimensions)}
            className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
              showDimensions
                ? 'bg-cosmic-nova text-white'
                : 'bg-white/5 text-celestial-silver/60 hover:text-celestial-silver'
            }`}
          >
            📊 维度
          </button>
        </div>
      </div>

      {/* 主图表区域 */}
      <div className="flex gap-4">
        {/* 趋势图 */}
        <div className="flex-1 relative glass-card p-4">
          <svg
            width="100%"
            height={chartHeight}
            viewBox={`0 0 ${chartWidth} ${chartHeight}`}
            preserveAspectRatio="none"
            className="overflow-visible"
          >
            <defs>
              {/* 渐变填充 */}
              <linearGradient id="trendGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" stopColor={COLORS.overall} stopOpacity="0.4" />
                <stop offset="100%" stopColor={COLORS.overall} stopOpacity="0.05" />
              </linearGradient>
              {/* 发光效果 */}
              <filter id="lineGlow">
                <feGaussianBlur stdDeviation="2" result="blur" />
                <feMerge>
                  <feMergeNode in="blur" />
                  <feMergeNode in="SourceGraphic" />
                </feMerge>
              </filter>
            </defs>

            {/* 网格线 */}
            <g className="grid">
              {/* 水平网格线 */}
              {[0, 25, 50, 75, 100].map((v) => {
                const y = getY(minScore + (maxScore - minScore) * (v / 100));
                return (
                  <g key={v}>
                    <line
                      x1={chartPadding.left}
                      y1={y}
                      x2={chartWidth - chartPadding.right}
                      y2={y}
                      stroke={COLORS.grid}
                      strokeDasharray={v === 50 ? 'none' : '4,4'}
                    />
                    <text
                      x={chartPadding.left - 10}
                      y={y}
                      textAnchor="end"
                      alignmentBaseline="middle"
                      fill="rgba(255,255,255,0.5)"
                      fontSize="10"
                    >
                      {Math.round(minScore + (maxScore - minScore) * (v / 100))}
                    </text>
                  </g>
                );
              })}
            </g>

            {/* 填充区域 */}
            <motion.path
              d={fillPath}
              fill="url(#trendGradient)"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.5 }}
            />

            {/* 趋势线 */}
            <motion.path
              d={trendPath}
              fill="none"
              stroke={COLORS.overall}
              strokeWidth="2"
              filter="url(#lineGlow)"
              initial={{ pathLength: 0 }}
              animate={{ pathLength: 1 }}
              transition={{ duration: 1 }}
            />

            {/* 重大事件标记 */}
            {dataPoints.map((d, i) => d.isMajorEvent && (
              <g key={`event-${i}`}>
                <line
                  x1={getX(i)}
                  y1={chartPadding.top}
                  x2={getX(i)}
                  y2={chartHeight - chartPadding.bottom}
                  stroke={COLORS.majorEvent}
                  strokeWidth="1"
                  strokeDasharray="4,4"
                  opacity="0.5"
                />
                <circle
                  cx={getX(i)}
                  cy={getY(d.overallScore)}
                  r={6}
                  fill={COLORS.majorEvent}
                  stroke="white"
                  strokeWidth="2"
                />
              </g>
            ))}

            {/* 交互层 - 数据点 */}
            {dataPoints.map((d, i) => (
              <g
                key={i}
                onMouseEnter={() => setHoveredIndex(i)}
                onMouseLeave={() => setHoveredIndex(null)}
                onClick={() => onTimeSelect?.(new Date(d.time))}
                className="cursor-pointer"
              >
                <rect
                  x={getX(i) - 10}
                  y={chartPadding.top}
                  width={20}
                  height={chartHeight - chartPadding.top - chartPadding.bottom}
                  fill="transparent"
                />
                <circle
                  cx={getX(i)}
                  cy={getY(d.overallScore)}
                  r={hoveredIndex === i ? 6 : 3}
                  fill={hoveredIndex === i ? 'white' : COLORS.overall}
                  stroke={hoveredIndex === i ? COLORS.overall : 'none'}
                  strokeWidth="2"
                  className="transition-all"
                />
              </g>
            ))}

            {/* X轴标签 */}
            <g className="x-axis">
              {dataPoints.map((d, i) => {
                // 根据数据量决定显示间隔
                const step = Math.max(1, Math.floor(dataPoints.length / 10));
                if (i % step !== 0 && i !== dataPoints.length - 1) return null;
                
                return (
                  <text
                    key={i}
                    x={getX(i)}
                    y={chartHeight - 10}
                    textAnchor="middle"
                    fill="rgba(255,255,255,0.5)"
                    fontSize="10"
                  >
                    {formatTimeLabel(d.time, granularity)}
                  </text>
                );
              })}
            </g>
          </svg>

          {/* 悬停提示 */}
          <AnimatePresence>
            {hoveredPoint && hoveredIndex !== null && (
              <motion.div
                className="absolute glass-card p-3 z-10 min-w-[200px]"
                style={{
                  left: getX(hoveredIndex) + 10,
                  top: getY(hoveredPoint.overallScore) - 50,
                }}
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.9 }}
              >
                <div className="text-sm font-medium mb-2">
                  {new Date(hoveredPoint.time).toLocaleString()}
                </div>
                <div className="text-2xl font-bold text-cosmic-nova mb-2">
                  {hoveredPoint.overallScore.toFixed(1)}
                </div>
                
                {/* 维度分数 */}
                <div className="grid grid-cols-2 gap-1 text-xs">
                  <div style={{ color: COLORS.career }}>事业: {hoveredPoint.dimensions.career.toFixed(1)}</div>
                  <div style={{ color: COLORS.relationship }}>关系: {hoveredPoint.dimensions.relationship.toFixed(1)}</div>
                  <div style={{ color: COLORS.health }}>健康: {hoveredPoint.dimensions.health.toFixed(1)}</div>
                  <div style={{ color: COLORS.finance }}>财务: {hoveredPoint.dimensions.finance.toFixed(1)}</div>
                  <div style={{ color: COLORS.spiritual }}>灵性: {hoveredPoint.dimensions.spiritual.toFixed(1)}</div>
                </div>

                {/* 活跃因子数量 */}
                <div className="mt-2 pt-2 border-t border-white/10 text-xs text-celestial-silver/60">
                  活跃因子: {filteredFactors.length} 个
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        {/* 右侧面板 */}
        {showFactorPanel && (
          <motion.div
            className="w-80 glass-card p-4 overflow-y-auto"
            style={{ maxHeight: chartHeight + 50 }}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
          >
            {hoveredPoint ? (
              <>
                {/* 雷达图 */}
                {showDimensions && (
                  <div className="mb-4 flex justify-center">
                    <DimensionRadarChart5832XY
                      scores={hoveredPoint.dimensions}
                      size={160}
                    />
                  </div>
                )}

                {/* 因子列表 */}
                <div className="space-y-2">
                  <div className="text-sm font-medium mb-2">
                    影响因子 ({filteredFactors.length})
                  </div>
                  {filteredFactors.length > 0 ? (
                    filteredFactors.map((factor, i) => (
                      <motion.div
                        key={i}
                        className={`p-2 rounded-lg cursor-pointer transition-colors ${
                          expandedFactor === factor
                            ? 'bg-white/15 border border-cosmic-nova/50'
                            : 'bg-white/5 hover:bg-white/10'
                        }`}
                        onClick={() => setExpandedFactor(expandedFactor === factor ? null : factor)}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: i * 0.05 }}
                      >
                        <div className="flex items-center justify-between">
                          <span className="text-sm">{factor.name}</span>
                          <span className={`text-xs px-1.5 py-0.5 rounded ${
                            factor.value > 0 ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'
                          }`}>
                            {factor.value > 0 ? '+' : ''}{factor.value.toFixed(2)}
                          </span>
                        </div>
                        
                        {/* 展开的生命周期图 */}
                        <AnimatePresence>
                          {expandedFactor === factor && factor.lifecycle && (
                            <motion.div
                              className="mt-2"
                              initial={{ height: 0, opacity: 0 }}
                              animate={{ height: 'auto', opacity: 1 }}
                              exit={{ height: 0, opacity: 0 }}
                            >
                              <FactorLifecycleChart7184BC
                                lifecycle={factor.lifecycle}
                                currentTime={new Date(hoveredPoint.time)}
                                factorName={factor.name}
                                baseValue={Math.abs(factor.value)}
                                color={factor.value > 0 ? '#4ECDC4' : '#FF6B6B'}
                                width={260}
                                height={80}
                              />
                            </motion.div>
                          )}
                        </AnimatePresence>
                      </motion.div>
                    ))
                  ) : (
                    <div className="text-sm text-celestial-silver/50 text-center py-4">
                      当前筛选条件下无因子
                    </div>
                  )}
                </div>
              </>
            ) : (
              <div className="flex items-center justify-center h-full text-celestial-silver/50">
                悬停查看详情
              </div>
            )}
          </motion.div>
        )}
      </div>
    </div>
  );
}

export default MultiGranularityTrendChart8567CD;

