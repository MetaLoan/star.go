/**
 * 人生趋势时间线组件
 * 组件命名规范：LifeTimeline + 4529 + PQ
 */

import { useMemo, useState } from 'react';
import { motion } from 'framer-motion';
import type { LifeTrend, LifeTrendPoint, Dimension } from '../../types';
import { DIMENSION_COLORS, DIMENSION_NAMES } from '../../utils/astro';

interface LifeTimelineProps {
  data: LifeTrend;
  currentAge: number;
  height?: number;
  showDimensions?: boolean;
  onPointClick?: (point: LifeTrendPoint) => void;
  className?: string;
}

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

  const dimensions: Array<Dimension | 'overall'> = ['overall', 'career', 'relationship', 'health', 'finance', 'spiritual'];

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
      {/* 维度选择器 */}
      {showDimensions && (
        <div className="flex flex-wrap gap-2 mb-4">
          {dimensions.map(dim => (
            <button
              key={dim}
              className={`px-3 py-1 rounded-full text-sm transition-all ${
                selectedDimension === dim
                  ? 'bg-white/20 text-white'
                  : 'bg-white/5 text-white/60 hover:bg-white/10'
              }`}
              onClick={() => setSelectedDimension(dim)}
            >
              {dim === 'overall' ? '综合' : DIMENSION_NAMES[dim]}
            </button>
          ))}
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

      {/* 悬停信息提示 */}
      {hoveredPoint && (
        <motion.div
          className="absolute bg-cosmic-dust/90 backdrop-blur-md rounded-lg p-3 text-sm border border-white/10 shadow-lg"
          style={{
            left: xScale(hoveredPoint.age),
            top: padding.top,
            transform: 'translateX(-50%)',
          }}
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <div className="font-medium text-white mb-1">
            {hoveredPoint.age}岁 ({hoveredPoint.year}年)
          </div>
          <div className="text-white/60">
            综合分数: <span className="text-cosmic-nova">{Math.round(hoveredPoint.overallScore)}</span>
          </div>
          {hoveredPoint.isMajorTransit && hoveredPoint.majorTransits && hoveredPoint.majorTransits.length > 0 && (
            <div className="mt-1 text-pink-400 text-xs">
              🌟 {hoveredPoint.majorTransits[0]}
            </div>
          )}
        </motion.div>
      )}

      {/* 峰值/谷值标记 */}
      {data.summary && (
        <div className="flex justify-between mt-4 text-xs text-white/60">
          <div className="flex items-center gap-4">
            <span className="flex items-center gap-1">
              <span className="w-2 h-2 rounded-full bg-green-500"></span>
              高峰年: {data.summary.peakYears?.slice(0, 3).map(y => `${y}年`).join(', ') || '无'}
            </span>
            <span className="flex items-center gap-1">
              <span className="w-2 h-2 rounded-full bg-orange-500"></span>
              挑战年: {data.summary.challengeYears?.slice(0, 3).map(y => `${y}年`).join(', ') || '无'}
            </span>
          </div>
        </div>
      )}
    </div>
  );
}

export default LifeTimeline4529PQ;

