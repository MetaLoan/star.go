/**
 * 每日预测卡片组件
 * 组件命名规范：DailyForecastCard + 7821 + MN
 */

import { motion } from 'framer-motion';
import type { DailyForecast, Dimension } from '../../types';
import { 
  DIMENSION_NAMES, 
  DIMENSION_COLORS, 
  DIMENSION_ICONS,
  getScoreLevel,
  formatDate,
} from '../../utils/astro';

interface DailyForecastCardProps {
  forecast: DailyForecast;
  isToday?: boolean;
  isExpanded?: boolean;
  onClick?: () => void;
  className?: string;
}

export function DailyForecastCard7821MN({
  forecast,
  isToday = false,
  isExpanded = false,
  onClick,
  className = '',
}: DailyForecastCardProps) {
  const { label, color } = getScoreLevel(forecast.overallScore);
  const dimensions: Dimension[] = ['career', 'relationship', 'health', 'finance', 'spiritual'];
  
  // 兼容两种字段名
  const getDimensionScore = (dim: Dimension): number => {
    if (forecast.dimensions) {
      return forecast.dimensions[dim] || 50;
    }
    return 50;
  };

  return (
    <motion.div
      className={`glass-card p-4 cursor-pointer ${isToday ? 'ring-2 ring-cyan-400/50' : ''} ${className}`}
      onClick={onClick}
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
      layout
    >
      {/* 头部 */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          {isToday && (
            <span className="px-2 py-0.5 bg-cyan-400/20 text-cyan-400 text-xs rounded-full">
              今日
            </span>
          )}
          <span className="text-white/80 font-medium">
            {formatDate(forecast.date)}
          </span>
          {forecast.dayOfWeek && (
            <span className="text-white/50 text-sm">{forecast.dayOfWeek}</span>
          )}
        </div>
        
        {/* 总分 */}
        <div className="flex items-center gap-2">
          <motion.span
            className="text-2xl font-bold"
            style={{ color }}
            initial={{ scale: 0.8 }}
            animate={{ scale: 1 }}
          >
            {Math.round(forecast.overallScore)}
          </motion.span>
          <span
            className="text-xs px-2 py-0.5 rounded-full"
            style={{ backgroundColor: `${color}20`, color }}
          >
            {label}
          </span>
        </div>
      </div>

      {/* 主题 */}
      {forecast.overallTheme && (
        <div className="text-sm text-white/70 mb-3">
          🌟 {forecast.overallTheme}
        </div>
      )}

      {/* 维度分数条 */}
      {forecast.dimensions && (
        <div className="space-y-2 mb-3">
          {dimensions.map((dimension, index) => (
            <motion.div
              key={dimension}
              className="flex items-center gap-2"
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.05 }}
            >
              <span className="text-sm w-16 flex items-center gap-1">
                <span>{DIMENSION_ICONS[dimension]}</span>
                <span className="text-white/60 text-xs">{DIMENSION_NAMES[dimension]}</span>
              </span>
              <div className="flex-1 h-1.5 bg-white/10 rounded-full overflow-hidden">
                <motion.div
                  className="h-full rounded-full"
                  style={{ backgroundColor: DIMENSION_COLORS[dimension] }}
                  initial={{ width: 0 }}
                  animate={{ width: `${getDimensionScore(dimension)}%` }}
                  transition={{ duration: 0.5, delay: index * 0.05 }}
                />
              </div>
              <span className="text-xs text-white/40 w-8 text-right">
                {Math.round(getDimensionScore(dimension))}
              </span>
            </motion.div>
          ))}
        </div>
      )}

      {/* 月相信息 */}
      {forecast.moonPhase && (
        <div className="flex items-center gap-2 text-sm text-white/60 mb-2">
          <span>🌙</span>
          <span>{forecast.moonPhase.name}</span>
          {forecast.moonSign && (
            <>
              <span className="text-white/40">|</span>
              <span>月亮在 {forecast.moonSign.name}</span>
            </>
          )}
        </div>
      )}

      {/* 展开内容 */}
      {isExpanded && (
        <motion.div
          initial={{ opacity: 0, height: 0 }}
          animate={{ opacity: 1, height: 'auto' }}
          exit={{ opacity: 0, height: 0 }}
          className="mt-4 pt-4 border-t border-white/10"
        >
          {/* 小时预测 */}
          {forecast.hourlyBreakdown && forecast.hourlyBreakdown.length > 0 && (
            <div className="mb-3">
              <h4 className="text-sm font-medium text-white/80 mb-2">小时预测</h4>
              <div className="grid grid-cols-6 gap-1">
                {forecast.hourlyBreakdown.filter((_, i) => i % 4 === 0).map((hour) => (
                  <div 
                    key={hour.hour} 
                    className="text-xs text-center p-1 rounded bg-white/5"
                  >
                    <div className="text-white/40">{hour.hour}:00</div>
                    <div className="text-cyan-400">{Math.round(hour.score)}</div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* 活跃相位 */}
          {forecast.activeAspects && forecast.activeAspects.length > 0 && (
            <div className="mb-3">
              <h4 className="text-sm font-medium text-white/80 mb-2">主要相位</h4>
              <div className="space-y-1">
                {forecast.activeAspects.slice(0, 3).map((asp, i) => (
                  <div key={i} className="text-xs text-white/60">
                    {asp.interpretation}
                  </div>
                ))}
              </div>
            </div>
          )}
        </motion.div>
      )}
    </motion.div>
  );
}

export default DailyForecastCard7821MN;
