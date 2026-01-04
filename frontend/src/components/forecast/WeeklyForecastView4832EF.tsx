/**
 * WeeklyForecastView4832EF - 每周预测展示组件
 * 显示周度趋势、每日分数和关键事件
 */

import { motion } from 'framer-motion';
import type { WeeklyForecast, DailySummary } from '../../types';

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
  // 使用 dailySummaries 或 dailyForecasts
  const dailyData = forecast.dailySummaries || forecast.dailyForecasts || [];
  const weekStart = forecast.weekStart || forecast.startDate;
  const weekEnd = forecast.weekEnd || forecast.endDate;
  const avgScore = forecast.weeklyAverage || forecast.overallScore || 50;
  const themes = forecast.weeklyThemes || (forecast.overallTheme ? [forecast.overallTheme] : []);

  // 计算最高最低分
  const scores = dailyData.map((d) => d.overallScore || 0);
  const maxScore = scores.length > 0 ? Math.max(...scores) : 50;
  const minScore = scores.length > 0 ? Math.min(...scores) : 50;

  return (
    <div className={`glass-card p-6 ${className}`}>
      {/* 周度概览 */}
      <div className="mb-6">
        <h3 className="text-lg font-bold mb-2">
          📅 {weekStart} - {weekEnd}
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
      {dailyData.length > 0 && (
        <div className="mb-6">
          <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
            周度趋势
          </h4>
          <div className="flex items-end gap-2 h-32">
            {dailyData.map((day, index) => {
              const date = new Date(day.date);
              const weekday = weekdayNames[date.getDay()];
              const score = day.overallScore || 50;
              const height = (score / 100) * 100;
              const style = getScoreStyle(score);

              return (
                <motion.div
                  key={day.date || index}
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
                      {score.toFixed(0)}
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
      )}

      {/* 主题关键词 */}
      {themes.length > 0 && (
        <div className="mb-6">
          <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
            本周主题
          </h4>
          <div className="flex flex-wrap gap-2">
            {themes.map((theme, index) => (
              <span
                key={index}
                className="px-3 py-1 rounded-full bg-cosmic-nova/20 text-cosmic-nova text-sm"
              >
                {theme}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* 每日详情 */}
      {dailyData.length > 0 && (
        <div>
          <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
            每日详情
          </h4>
          <div className="space-y-3 max-h-[400px] overflow-y-auto">
            {dailyData.map((day) => (
              <DayCard key={day.date} day={day} />
            ))}
          </div>
        </div>
      )}

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
                {keyDate.event || keyDate.reason || '重要事件'}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

// 每日卡片组件
function DayCard({ day }: { day: DailySummary }) {
  const date = new Date(day.date);
  const weekday = weekdayNames[date.getDay()];
  const dateStr = `${date.getMonth() + 1}/${date.getDate()}`;
  const score = day.overallScore || 50;
  const style = getScoreStyle(score);

  // 获取月相和月亮星座信息
  const moonPhase = day.moonPhase || '';
  const moonSign = day.moonSign || '';

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
            {score.toFixed(0)}
          </div>
          <div>
            <div className="font-medium">
              {weekday} {dateStr}
            </div>
            <div className="text-xs text-celestial-silver/60">{style.label}</div>
          </div>
        </div>
        {(moonSign || moonPhase) && (
          <div className="text-right">
            <div className="text-sm">{moonSign}</div>
            <div className="text-xs text-celestial-silver/60">{moonPhase}</div>
          </div>
        )}
      </div>

      {/* 五维度分数 */}
      {day.dimensions && (
        <div className="grid grid-cols-5 gap-2">
          {Object.entries(day.dimensions).map(([dim, scoreVal]) => (
            <DimensionBadge key={dim} dimension={dim} score={scoreVal as number} />
          ))}
        </div>
      )}

      {/* 每日主题 */}
      {day.theme && (
        <div className="mt-3 text-sm text-celestial-silver/80">
          {day.theme}
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
      <div className="text-lg">{dimensionIcons[dimension] || '📊'}</div>
      <div className="text-xs font-bold" style={{ color: style.color }}>
        {score.toFixed(0)}
      </div>
      <div className="text-xs text-celestial-silver/60">
        {dimensionNames[dimension] || dimension}
      </div>
    </div>
  );
}

export default WeeklyForecastView4832EF;
