/**
 * DimensionRadarChart5832XY - 五维度雷达图
 * 展示事业/关系/健康/财务/灵性五个维度的分数
 */

import { useMemo } from 'react';
import { motion } from 'framer-motion';
import type { DimensionScores } from '../../types';

interface DimensionRadarChartProps {
  scores: DimensionScores;
  baseScores?: DimensionScores; // 可选：本命盘基础分作为对比
  size?: number;
  showLabels?: boolean;
  showValues?: boolean;
  animated?: boolean;
  className?: string;
}

// 维度配置
const dimensions = [
  { key: 'career', label: '事业', icon: '💼', color: '#FF6B6B', angle: -90 },
  { key: 'relationship', label: '关系', icon: '💕', color: '#FF9F43', angle: -18 },
  { key: 'spiritual', label: '灵性', icon: '🔮', color: '#A855F7', angle: 54 },
  { key: 'finance', label: '财务', icon: '💰', color: '#4ECDC4', angle: 126 },
  { key: 'health', label: '健康', icon: '💪', color: '#4FC3F7', angle: 198 },
] as const;

// 角度转弧度
const toRad = (deg: number) => (deg * Math.PI) / 180;

// 极坐标转笛卡尔坐标
function polarToCartesian(centerX: number, centerY: number, radius: number, angleInDegrees: number) {
  const rad = toRad(angleInDegrees);
  return {
    x: centerX + radius * Math.cos(rad),
    y: centerY + radius * Math.sin(rad),
  };
}

export function DimensionRadarChart5832XY({
  scores,
  baseScores,
  size = 200,
  showLabels = true,
  showValues = true,
  animated = true,
  className = '',
}: DimensionRadarChartProps) {
  const center = size / 2;
  const maxRadius = size * 0.38;
  const labelRadius = size * 0.48;

  // 生成雷达图路径
  const radarPath = useMemo(() => {
    const points = dimensions.map((dim) => {
      const score = scores[dim.key as keyof DimensionScores] || 50;
      const radius = (score / 100) * maxRadius;
      return polarToCartesian(center, center, radius, dim.angle);
    });
    return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ') + ' Z';
  }, [scores, center, maxRadius]);

  // 生成基础分路径（如果提供）
  const baseRadarPath = useMemo(() => {
    if (!baseScores) return null;
    const points = dimensions.map((dim) => {
      const score = baseScores[dim.key as keyof DimensionScores] || 50;
      const radius = (score / 100) * maxRadius;
      return polarToCartesian(center, center, radius, dim.angle);
    });
    return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ') + ' Z';
  }, [baseScores, center, maxRadius]);

  // 生成网格圆
  const gridCircles = [20, 40, 60, 80, 100];

  return (
    <div className={`relative ${className}`} style={{ width: size, height: size }}>
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        <defs>
          {/* 渐变填充 */}
          <linearGradient id="radarGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#00D4FF" stopOpacity="0.6" />
            <stop offset="100%" stopColor="#A855F7" stopOpacity="0.6" />
          </linearGradient>
          <linearGradient id="baseGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#FFE66D" stopOpacity="0.3" />
            <stop offset="100%" stopColor="#FF9F43" stopOpacity="0.3" />
          </linearGradient>
          {/* 发光效果 */}
          <filter id="radarGlow" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur stdDeviation="3" result="blur" />
            <feMerge>
              <feMergeNode in="blur" />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>

        {/* 背景网格 */}
        <g className="grid">
          {/* 同心圆 */}
          {gridCircles.map((value) => (
            <circle
              key={value}
              cx={center}
              cy={center}
              r={(value / 100) * maxRadius}
              fill="none"
              stroke="rgba(255,255,255,0.1)"
              strokeWidth="1"
              strokeDasharray={value === 100 ? 'none' : '2,4'}
            />
          ))}
          {/* 轴线 */}
          {dimensions.map((dim) => {
            const point = polarToCartesian(center, center, maxRadius, dim.angle);
            return (
              <line
                key={dim.key}
                x1={center}
                y1={center}
                x2={point.x}
                y2={point.y}
                stroke="rgba(255,255,255,0.15)"
                strokeWidth="1"
              />
            );
          })}
        </g>

        {/* 基础分雷达图（底层） */}
        {baseRadarPath && (
          <motion.path
            d={baseRadarPath}
            fill="url(#baseGradient)"
            stroke="#FFE66D"
            strokeWidth="1"
            strokeDasharray="4,2"
            initial={animated ? { opacity: 0 } : undefined}
            animate={animated ? { opacity: 1 } : undefined}
            transition={{ duration: 0.5 }}
          />
        )}

        {/* 当前分数雷达图 */}
        <motion.path
          d={radarPath}
          fill="url(#radarGradient)"
          stroke="#00D4FF"
          strokeWidth="2"
          filter="url(#radarGlow)"
          initial={animated ? { opacity: 0, scale: 0.5 } : undefined}
          animate={animated ? { opacity: 1, scale: 1 } : undefined}
          transition={{ duration: 0.6, type: 'spring', bounce: 0.3 }}
          style={{ transformOrigin: `${center}px ${center}px` }}
        />

        {/* 数据点 */}
        {dimensions.map((dim, index) => {
          const score = scores[dim.key as keyof DimensionScores] || 50;
          const radius = (score / 100) * maxRadius;
          const point = polarToCartesian(center, center, radius, dim.angle);

          return (
            <motion.circle
              key={dim.key}
              cx={point.x}
              cy={point.y}
              r={4}
              fill={dim.color}
              stroke="white"
              strokeWidth="2"
              filter="url(#radarGlow)"
              initial={animated ? { scale: 0 } : undefined}
              animate={animated ? { scale: 1 } : undefined}
              transition={{ delay: index * 0.1, duration: 0.3 }}
            />
          );
        })}
      </svg>

      {/* 标签 */}
      {showLabels && (
        <div className="absolute inset-0 pointer-events-none">
          {dimensions.map((dim) => {
            const score = scores[dim.key as keyof DimensionScores] || 50;
            const point = polarToCartesian(center, center, labelRadius, dim.angle);

            return (
              <motion.div
                key={dim.key}
                className="absolute flex flex-col items-center"
                style={{
                  left: point.x,
                  top: point.y,
                  transform: 'translate(-50%, -50%)',
                }}
                initial={animated ? { opacity: 0, y: 10 } : undefined}
                animate={animated ? { opacity: 1, y: 0 } : undefined}
                transition={{ delay: 0.3 }}
              >
                <span className="text-lg">{dim.icon}</span>
                <span className="text-xs text-celestial-silver/80">{dim.label}</span>
                {showValues && (
                  <span
                    className="text-xs font-bold"
                    style={{ color: dim.color }}
                  >
                    {score.toFixed(0)}
                  </span>
                )}
              </motion.div>
            );
          })}
        </div>
      )}

      {/* 中心分数 */}
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
        <motion.div
          className="text-center"
          initial={animated ? { scale: 0 } : undefined}
          animate={animated ? { scale: 1 } : undefined}
          transition={{ delay: 0.4, type: 'spring' }}
        >
          <div className="text-2xl font-bold text-cosmic-nova">
            {Math.round(
              (scores.career + scores.relationship + scores.health + scores.finance + scores.spiritual) / 5
            )}
          </div>
          <div className="text-xs text-celestial-silver/60">平均分</div>
        </motion.div>
      </div>
    </div>
  );
}

// 迷你版本
export function DimensionRadarChartMini5832XY({
  scores,
  size = 80,
  className = '',
}: Pick<DimensionRadarChartProps, 'scores' | 'size' | 'className'>) {
  return (
    <DimensionRadarChart5832XY
      scores={scores}
      size={size}
      showLabels={false}
      showValues={false}
      animated={false}
      className={className}
    />
  );
}

export default DimensionRadarChart5832XY;

