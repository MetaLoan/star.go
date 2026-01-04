/**
 * 分数卡片组件
 * 组件命名规范：ScoreCard + 5612 + XY
 */

import { motion } from 'framer-motion';
import type { Dimension } from '../../types';
import { DIMENSION_NAMES, DIMENSION_COLORS, DIMENSION_ICONS, getScoreLevel } from '../../utils/astro';

interface ScoreCardProps {
  title: string;
  score: number;
  maxScore?: number;
  dimension?: Dimension;
  showLevel?: boolean;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export function ScoreCard5612XY({
  title,
  score,
  maxScore = 100,
  dimension,
  showLevel = true,
  size = 'md',
  className = '',
}: ScoreCardProps) {
  const percentage = Math.min((score / maxScore) * 100, 100);
  const { label, color } = getScoreLevel(score);
  
  const barColor = dimension ? DIMENSION_COLORS[dimension] : color;
  const icon = dimension ? DIMENSION_ICONS[dimension] : '⭐';
  const displayTitle = dimension ? DIMENSION_NAMES[dimension] : title;
  
  const sizeClasses = {
    sm: 'p-2',
    md: 'p-4',
    lg: 'p-6',
  };
  
  const fontSizes = {
    sm: 'text-xs',
    md: 'text-sm',
    lg: 'text-base',
  };
  
  const scoreFontSizes = {
    sm: 'text-lg',
    md: 'text-2xl',
    lg: 'text-4xl',
  };

  return (
    <motion.div
      className={`glass-card ${sizeClasses[size]} ${className}`}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ scale: 1.02 }}
      transition={{ duration: 0.2 }}
    >
      {/* 标题行 */}
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <span className="text-base">{icon}</span>
          <span className={`${fontSizes[size]} text-white/80 font-medium`}>
            {displayTitle}
          </span>
        </div>
        {showLevel && (
          <span
            className={`${fontSizes[size]} px-2 py-0.5 rounded-full`}
            style={{ backgroundColor: `${color}20`, color }}
          >
            {label}
          </span>
        )}
      </div>
      
      {/* 分数 */}
      <div className="flex items-baseline gap-1 mb-3">
        <motion.span
          className={`${scoreFontSizes[size]} font-bold text-white`}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          key={score}
        >
          {Math.round(score)}
        </motion.span>
        <span className="text-white/40 text-sm">/ {maxScore}</span>
      </div>
      
      {/* 进度条 */}
      <div className="score-bar">
        <motion.div
          className="score-bar-fill"
          style={{ backgroundColor: barColor }}
          initial={{ width: 0 }}
          animate={{ width: `${percentage}%` }}
          transition={{ duration: 0.8, ease: 'easeOut' }}
        />
      </div>
    </motion.div>
  );
}

// 维度分数组
interface DimensionScoresCardProps {
  scores: {
    career: number;
    relationship: number;
    health: number;
    finance: number;
    spiritual: number;
  };
  layout?: 'horizontal' | 'vertical' | 'grid';
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export function DimensionScoresCard5612XY({
  scores,
  layout = 'grid',
  size = 'sm',
  className = '',
}: DimensionScoresCardProps) {
  const dimensions: Dimension[] = ['career', 'relationship', 'health', 'finance', 'spiritual'];
  
  const layoutClasses = {
    horizontal: 'flex flex-wrap gap-4',
    vertical: 'flex flex-col gap-3',
    grid: 'grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3',
  };

  return (
    <div className={`${layoutClasses[layout]} ${className}`}>
      {dimensions.map((dimension, index) => (
        <motion.div
          key={dimension}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: index * 0.1 }}
          className="flex-1 min-w-[120px]"
        >
          <ScoreCard5612XY
            title=""
            score={scores[dimension]}
            dimension={dimension}
            showLevel={false}
            size={size}
          />
        </motion.div>
      ))}
    </div>
  );
}

export default ScoreCard5612XY;

