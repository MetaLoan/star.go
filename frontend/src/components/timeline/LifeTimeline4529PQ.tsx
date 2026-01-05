/**
 * 人生趋势时间线组件
 * 组件命名规范：LifeTimeline + 4529 + PQ
 * 
 * 数据说明：
 * - 年度分数由12个月分平均聚合
 * - 月分由该月所有日分平均聚合
 * - 日分由24个小时分平均聚合
 * - 小时分是基于 Swiss Ephemeris 的原始计算
 */

import { useMemo, useState } from 'react';
import { motion } from 'framer-motion';
import type { LifeTrend, LifeTrendPoint, Dimension } from '../../types';
import { DIMENSION_COLORS, DIMENSION_NAMES, DIMENSION_ICONS } from '../../utils/astro';

interface LifeTimelineProps {
  data: LifeTrend;
  currentAge: number;
  height?: number;
  showDimensions?: boolean;
  onPointClick?: (point: LifeTrendPoint) => void;
  className?: string;
}

// 维度选项配置
const DIMENSION_OPTIONS: Array<{ id: Dimension | 'overall'; label: string; icon: string }> = [
  { id: 'overall', label: '综合', icon: '📊' },
  { id: 'career', label: '事业', icon: '💼' },
  { id: 'relationship', label: '关系', icon: '❤️' },
  { id: 'health', label: '健康', icon: '🏃' },
  { id: 'finance', label: '财务', icon: '💰' },
  { id: 'spiritual', label: '灵性', icon: '🧘' },
];

export function LifeTimeline4529PQ({
  data,
  currentAge,
  height = 200,
  showDimensions = false,
  onPointClick,
  className = '',
}: LifeTimelineProps) {
  const [hoveredAge, setHoveredAge] = useState<number | null>(null);
  const [selectedDimension, setSelectedDimension] = useState<Dimension | 'overall'>('overall');

  // 计算图表尺寸
  const padding = { top: 20, right: 20, bottom: 30, left: 40 };
  const chartWidth = 800;
  const chartHeight = height - padding.top - padding.bottom;

  // 计算数据范围和比例
  const { minScore, maxScore, xScale, yScale } = useMemo(() => {
    let min = 100, max = 0;
    
    data.points.forEach(point => {
      const score = selectedDimension === 'overall' 
        ? point.overallScore 
        : point.dimensions[selectedDimension];
      min = Math.min(min, score);
      max = Math.max(max, score);
    });

    // 添加 padding
    min = Math.max(0, min - 5);
    max = Math.min(100, max + 5);

    const xScale = (age: number) => {
      const minAge = data.points[0]?.age || 0;
      const maxAge = data.points[data.points.length - 1]?.age || 80;
      return padding.left + ((age - minAge) / (maxAge - minAge)) * (chartWidth - padding.left - padding.right);
    };

    const yScale = (score: number) => {
      return padding.top + chartHeight - ((score - min) / (max - min)) * chartHeight;
    };

    return { minScore: min, maxScore: max, xScale, yScale };
  }, [data.points, selectedDimension, chartWidth, chartHeight, padding]);

  // 生成 SVG 路径
  const linePath = useMemo(() => {
    if (data.points.length === 0) return '';

    const points = data.points.map(point => {
      const score = selectedDimension === 'overall' 
        ? point.overallScore 
        : point.dimensions[selectedDimension];
      return `${xScale(point.age)},${yScale(score)}`;
    });

    return `M ${points.join(' L ')}`;
  }, [data.points, selectedDimension, xScale, yScale]);

  // 生成渐变填充区域路径
  const areaPath = useMemo(() => {
    if (data.points.length === 0) return '';

    const points = data.points.map(point => {
      const score = selectedDimension === 'overall' 
        ? point.overallScore 
        : point.dimensions[selectedDimension];
      return `${xScale(point.age)},${yScale(score)}`;
    });

    const firstX = xScale(data.points[0].age);
    const lastX = xScale(data.points[data.points.length - 1].age);
    const baseY = yScale(minScore);

    return `M ${firstX},${baseY} L ${points.join(' L ')} L ${lastX},${baseY} Z`;
  }, [data.points, selectedDimension, xScale, yScale, minScore]);

  // 当前年龄线位置
  const currentAgeX = xScale(currentAge);

  // 悬停点信息
  const hoveredPoint = hoveredAge !== null 
    ? data.points.find(p => p.age === hoveredAge) 
    : null;

  const lineColor = selectedDimension === 'overall' 
    ? '#00D4FF' 
    : DIMENSION_COLORS[selectedDimension];

  return (
    <div className={`glass-card p-4 ${className}`}>
      {/* 维度选择器 - 增强样式 */}
      {showDimensions && (
        <div className="flex flex-wrap gap-2 mb-4">
          {DIMENSION_OPTIONS.map(dim => {
            const isSelected = selectedDimension === dim.id;
            const dimColor = dim.id === 'overall' ? '#00D4FF' : DIMENSION_COLORS[dim.id];
            return (
              <button
                key={dim.id}
                className={`px-3 py-1.5 rounded-full text-sm transition-all flex items-center gap-1.5 ${
                  isSelected
                    ? 'text-white shadow-lg'
                    : 'text-white/60 hover:text-white/80'
                }`}
                style={{
                  backgroundColor: isSelected ? `${dimColor}30` : 'rgba(255,255,255,0.05)',
                  border: isSelected ? `1px solid ${dimColor}` : '1px solid transparent',
                  boxShadow: isSelected ? `0 0 12px ${dimColor}40` : 'none',
                }}
                onClick={() => setSelectedDimension(dim.id)}
              >
                <span>{dim.icon}</span>
                <span>{dim.label}</span>
              </button>
            );
          })}
        </div>
      )}

      {/* 图表 */}
      <svg
        width="100%"
        height={height}
        viewBox={`0 0 ${chartWidth} ${height}`}
        className="overflow-visible"
        preserveAspectRatio="xMidYMid meet"
      >
        <defs>
          {/* 渐变填充 */}
          <linearGradient id="areaGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor={lineColor} stopOpacity="0.3" />
            <stop offset="100%" stopColor={lineColor} stopOpacity="0" />
          </linearGradient>

          {/* 发光效果 */}
          <filter id="lineGlow">
            <feGaussianBlur stdDeviation="2" result="coloredBlur" />
            <feMerge>
              <feMergeNode in="coloredBlur" />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>

        {/* Y 轴刻度线 */}
        {[0, 25, 50, 75, 100].map(score => {
          if (score < minScore || score > maxScore) return null;
          const y = yScale(score);
          return (
            <g key={score}>
              <line
                x1={padding.left}
                y1={y}
                x2={chartWidth - padding.right}
                y2={y}
                stroke="rgba(255,255,255,0.05)"
                strokeDasharray="4 4"
              />
              <text
                x={padding.left - 8}
                y={y}
                fill="rgba(255,255,255,0.4)"
                fontSize="10"
                textAnchor="end"
                dominantBaseline="central"
              >
                {score}
              </text>
            </g>
          );
        })}

        {/* X 轴刻度 */}
        {data.points
          .filter((_, i) => i % 10 === 0)
          .map(point => {
            const x = xScale(point.age);
            return (
              <g key={point.age}>
                <line
                  x1={x}
                  y1={padding.top}
                  x2={x}
                  y2={height - padding.bottom}
                  stroke="rgba(255,255,255,0.05)"
                />
                <text
                  x={x}
                  y={height - padding.bottom + 15}
                  fill="rgba(255,255,255,0.4)"
                  fontSize="10"
                  textAnchor="middle"
                >
                  {point.age}岁
                </text>
              </g>
            );
          })}

        {/* 填充区域 */}
        <motion.path
          d={areaPath}
          fill="url(#areaGradient)"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.5 }}
        />

        {/* 趋势线 */}
        <motion.path
          d={linePath}
          fill="none"
          stroke={lineColor}
          strokeWidth="2"
          filter="url(#lineGlow)"
          initial={{ pathLength: 0 }}
          animate={{ pathLength: 1 }}
          transition={{ duration: 1.5, ease: 'easeOut' }}
        />

        {/* 当前年龄线 */}
        <line
          x1={currentAgeX}
          y1={padding.top}
          x2={currentAgeX}
          y2={height - padding.bottom}
          stroke="#FF6B9D"
          strokeWidth="2"
          strokeDasharray="4 4"
        />
        <text
          x={currentAgeX}
          y={padding.top - 5}
          fill="#FF6B9D"
          fontSize="10"
          textAnchor="middle"
        >
          现在
        </text>

        {/* 交互层 */}
          {data.points.map((point) => {
          const x = xScale(point.age);
          const score = selectedDimension === 'overall' 
            ? point.overallScore 
            : point.dimensions[selectedDimension];
          const y = yScale(score);
          const isHovered = hoveredAge === point.age;

          return (
            <g
              key={point.age}
              onMouseEnter={() => setHoveredAge(point.age)}
              onMouseLeave={() => setHoveredAge(null)}
              onClick={() => onPointClick?.(point)}
              style={{ cursor: onPointClick ? 'pointer' : 'default' }}
            >
              {/* 悬停区域 */}
              <rect
                x={x - 5}
                y={padding.top}
                width={10}
                height={chartHeight}
                fill="transparent"
              />

              {/* 数据点 */}
              {(isHovered || point.age === currentAge) && (
                <motion.circle
                  cx={x}
                  cy={y}
                  r={isHovered ? 6 : 4}
                  fill={lineColor}
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                />
              )}

              {/* 重大事件标记 */}
              {point.isMajorTransit && (
                <circle
                  cx={x}
                  cy={y}
                  r={3}
                  fill="#A855F7"
                  opacity={0.8}
                />
              )}
            </g>
          );
        })}
      </svg>

      {/* 悬停信息提示 - 增强版 */}
      {hoveredPoint && (
        <motion.div
          className="absolute bg-cosmic-dust/95 backdrop-blur-md rounded-xl p-4 text-sm border border-white/20 shadow-2xl z-20"
          style={{
            left: Math.min(Math.max(xScale(hoveredPoint.age), 120), chartWidth - 120),
            top: padding.top + 10,
            transform: 'translateX(-50%)',
            minWidth: '200px',
          }}
          initial={{ opacity: 0, y: -10, scale: 0.95 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
        >
          {/* 头部 */}
          <div className="font-bold text-white mb-3 flex items-center justify-between">
            <span>{hoveredPoint.age}岁</span>
            <span className="text-white/50 text-xs">{hoveredPoint.year}年</span>
          </div>
          
          {/* 综合分数 */}
          <div className="flex items-center justify-between mb-3 pb-2 border-b border-white/10">
            <span className="text-white/60">综合分数</span>
            <span className="text-xl font-bold text-cosmic-nova">{Math.round(hoveredPoint.overallScore)}</span>
          </div>
          
          {/* 五维度分数 */}
          <div className="space-y-1.5">
            {(Object.keys(hoveredPoint.dimensions) as Dimension[]).map((dim) => {
              const score = hoveredPoint.dimensions[dim];
              const icon = DIMENSION_ICONS[dim];
              const name = DIMENSION_NAMES[dim];
              const color = DIMENSION_COLORS[dim];
              return (
                <div key={dim} className="flex items-center gap-2">
                  <span className="text-xs">{icon}</span>
                  <span className="text-xs text-white/60 w-8">{name}</span>
                  <div className="flex-1 h-1.5 bg-white/10 rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full"
                      style={{ 
                        width: `${Math.min(score, 100)}%`,
                        backgroundColor: color,
                      }}
                    />
                  </div>
                  <span className="text-xs w-6 text-right" style={{ color }}>{Math.round(score)}</span>
                </div>
              );
            })}
          </div>
          
          {/* 重大行运 */}
          {hoveredPoint.isMajorTransit && hoveredPoint.majorTransits && hoveredPoint.majorTransits.length > 0 && (
            <div className="mt-3 pt-2 border-t border-white/10">
              <div className="text-pink-400 text-xs flex items-center gap-1">
                <span>🌟</span>
                <span>{hoveredPoint.majorTransits[0]}</span>
              </div>
            </div>
          )}
          
          {/* 数据来源说明 */}
          <div className="mt-2 text-[10px] text-white/30 text-center">
            年度分数 = 12个月平均
          </div>
        </motion.div>
      )}

      {/* 图例和峰值/谷值标记 */}
      <div className="flex flex-wrap justify-between items-center mt-4 gap-4">
        {/* 当前选中维度说明 */}
        <div className="flex items-center gap-2 text-xs text-white/60">
          <span>当前显示:</span>
          <span 
            className="px-2 py-0.5 rounded-full"
            style={{ 
              backgroundColor: `${lineColor}20`,
              border: `1px solid ${lineColor}`,
              color: lineColor,
            }}
          >
            {selectedDimension === 'overall' ? '综合' : DIMENSION_NAMES[selectedDimension]}
          </span>
          <span className="text-white/40">| 数据来源: Swiss Ephemeris</span>
        </div>
        
        {/* 峰值/谷值标记 */}
        {data.summary && (
          <div className="flex items-center gap-4 text-xs text-white/60">
            <span className="flex items-center gap-1">
              <span className="w-2 h-2 rounded-full bg-green-500 shadow-[0_0_6px_#22c55e]"></span>
              高峰年: {data.summary.peakYears?.slice(0, 3).map(y => `${y}年`).join(', ') || '无'}
            </span>
            <span className="flex items-center gap-1">
              <span className="w-2 h-2 rounded-full bg-orange-500 shadow-[0_0_6px_#f97316]"></span>
              挑战年: {data.summary.challengeYears?.slice(0, 3).map(y => `${y}年`).join(', ') || '无'}
            </span>
          </div>
        )}
      </div>
      
      {/* 计算说明 */}
      <div className="mt-3 text-[10px] text-white/30 text-center">
        年度分数聚合逻辑: 12个月分 → 平均 → 年分 | 月分 = 该月日分平均 | 日分 = 24小时分平均
      </div>
    </div>
  );
}

export default LifeTimeline4529PQ;

