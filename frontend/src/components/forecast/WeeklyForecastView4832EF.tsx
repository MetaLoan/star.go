/**
 * WeeklyForecastView4832EF - 每周预测展示组件
 * 显示周度趋势、每日分数和关键事件
 */

import { motion } from 'framer-motion';
import type { WeeklyForecast, DailyForecast } from '../../types';

interface WeeklyForecastViewProps {
  forecast: WeeklyForecast;
  className?: string;
}

// 星期名称
const weekdayNames = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];

// 获取分数等级和颜色
function getScoreStyle(score: number): { color: string; label: string } {
  if (score >= 80) return { color: '#4ECDC4', label: '极佳' };
  if (score >= 60) return { color: '#4FC3F7', label: '良好' };
  if (score >= 40) return { color: '#FFE66D', label: '平稳' };
  if (score >= 20) return { color: '#FF9F43', label: '谨慎' };
  return { color: '#FF6B6B', label: '低迷' };
}

export function WeeklyForecastView4832EF({
  forecast,
  className = '',
}: WeeklyForecastViewProps) {
  const maxScore = Math.max(...forecast.dailyForecasts.map((d) => d.overallScore));
  const minScore = Math.min(...forecast.dailyForecasts.map((d) => d.overallScore));
  const avgScore = forecast.weeklyAverage;

  return (
    <div className={`glass-card p-6 ${className}`}>
      {/* 周度概览 */}
      <div className="mb-6">
        <h3 className="text-lg font-bold mb-2">
          📅 {forecast.weekStart} - {forecast.weekEnd}
        </h3>
        <div className="flex items-center gap-4">
          <div className="flex-1">
            <div className="text-sm text-celestial-silver/60 mb-1">周平均分</div>
            <div
              className="text-3xl font-bold"
              style={{ color: getScoreStyle(avgScore).color }}
            >
              {avgScore.toFixed(0)}
            </div>
          </div>
          <div className="text-right">
            <div className="text-sm text-celestial-silver/60">
              最高 {maxScore.toFixed(0)} / 最低 {minScore.toFixed(0)}
            </div>
            <div className="text-sm mt-1">
              波动 {(maxScore - minScore).toFixed(0)}
            </div>
          </div>
        </div>
      </div>

      {/* 周度趋势图 */}
      <div className="mb-6">
        <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
          周度趋势
        </h4>
        <div className="flex items-end gap-2 h-32">
          {forecast.dailyForecasts.map((day, index) => {
            const date = new Date(day.date);
            const weekday = weekdayNames[date.getDay()];
            const height = (day.overallScore / 100) * 100;
            const style = getScoreStyle(day.overallScore);

            return (
              <motion.div
                key={day.date}
                className="flex-1 flex flex-col items-center"
                initial={{ scaleY: 0 }}
                animate={{ scaleY: 1 }}
                transition={{ delay: index * 0.1 }}
              >
                <div className="w-full relative" style={{ height: '100px' }}>
                  <motion.div
                    className="absolute bottom-0 w-full rounded-t-lg"
                    style={{
                      height: `${height}%`,
                      backgroundColor: style.color,
                      opacity: 0.8,
                    }}
                    whileHover={{ opacity: 1 }}
                  />
                  <div className="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-bold">
                    {day.overallScore.toFixed(0)}
                  </div>
                </div>
                <div className="text-xs text-celestial-silver/60 mt-2">
                  {weekday}
                </div>
              </motion.div>
            );
          })}
        </div>
      </div>

      {/* 主题关键词 */}
      <div className="mb-6">
        <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
          本周主题
        </h4>
        <div className="flex flex-wrap gap-2">
          {forecast.weeklyThemes.map((theme, index) => (
            <span
              key={index}
              className="px-3 py-1 rounded-full bg-cosmic-nova/20 text-cosmic-nova text-sm"
            >
              {theme}
            </span>
          ))}
        </div>
      </div>

      {/* 每日详情 */}
      <div>
        <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
          每日详情
        </h4>
        <div className="space-y-3 max-h-[400px] overflow-y-auto">
          {forecast.dailyForecasts.map((day) => (
            <DayCard key={day.date} day={day} />
          ))}
        </div>
      </div>

      {/* 重要提醒 */}
      {forecast.keyDates && forecast.keyDates.length > 0 && (
        <div className="mt-6 p-4 rounded-lg bg-cosmic-aurora/10 border border-cosmic-aurora/30">
          <h4 className="text-sm font-medium text-cosmic-aurora mb-2">
            ⚠️ 重要日期
          </h4>
          <ul className="space-y-1">
            {forecast.keyDates.map((keyDate, index) => (
              <li key={index} className="text-sm">
                <span className="font-medium">{keyDate.date}</span>:{' '}
                {keyDate.description}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

// 每日卡片组件
function DayCard({ day }: { day: DailyForecast }) {
  const date = new Date(day.date);
  const weekday = weekdayNames[date.getDay()];
  const dateStr = `${date.getMonth() + 1}/${date.getDate()}`;
  const style = getScoreStyle(day.overallScore);

  return (
    <motion.div
      className="p-4 rounded-lg bg-white/5 hover:bg-white/10 transition-colors"
      whileHover={{ scale: 1.01 }}
    >
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-3">
          <div
            className="w-10 h-10 rounded-full flex items-center justify-center text-lg font-bold"
            style={{ backgroundColor: `${style.color}20`, color: style.color }}
          >
            {day.overallScore.toFixed(0)}
          </div>
          <div>
            <div className="font-medium">
              {weekday} {dateStr}
            </div>
            <div className="text-xs text-celestial-silver/60">{style.label}</div>
          </div>
        </div>
        <div className="text-right">
          <div className="text-sm">{day.moonInfo.signName}</div>
          <div className="text-xs text-celestial-silver/60">
            {day.moonInfo.phaseName}
          </div>
        </div>
      </div>

      {/* 五维度分数 */}
      <div className="grid grid-cols-5 gap-2">
        {Object.entries(day.dimensionScores).map(([dim, score]) => (
          <DimensionBadge key={dim} dimension={dim} score={score} />
        ))}
      </div>

      {/* 每日主题 */}
      {day.dailyTheme && (
        <div className="mt-3 text-sm text-celestial-silver/80">
          {day.dailyTheme}
        </div>
      )}
    </motion.div>
  );
}

// 维度徽章
function DimensionBadge({
  dimension,
  score,
}: {
  dimension: string;
  score: number;
}) {
  const dimensionNames: Record<string, string> = {
    career: '事业',
    relationship: '关系',
    health: '健康',
    finance: '财务',
    spiritual: '灵性',
  };

  const dimensionIcons: Record<string, string> = {
    career: '💼',
    relationship: '💕',
    health: '💪',
    finance: '💰',
    spiritual: '🔮',
  };

  const style = getScoreStyle(score);

  return (
    <div className="text-center">
      <div className="text-lg">{dimensionIcons[dimension]}</div>
      <div className="text-xs font-bold" style={{ color: style.color }}>
        {score.toFixed(0)}
      </div>
      <div className="text-xs text-celestial-silver/60">
        {dimensionNames[dimension]}
      </div>
    </div>
  );
}

export default WeeklyForecastView4832EF;

