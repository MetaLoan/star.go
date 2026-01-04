/**
 * UnifiedTimeline8456GH - 统一时间序列组件
 * 支持 H/D/W/M/Y 多粒度切换，显示原始分数和标准化分数
 */

import { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import type { TimeSeriesPoint } from '../../types';

interface UnifiedTimelineProps {
  data: TimeSeriesPoint[];
  className?: string;
  onGranularityChange?: (granularity: Granularity) => void;
}

type Granularity = 'H' | 'D' | 'W' | 'M' | 'Y';

const granularityLabels: Record<Granularity, string> = {
  H: '小时',
  D: '每日',
  W: '每周',
  M: '每月',
  Y: '每年',
};

// 获取分数颜色
function getScoreColor(score: number): string {
  if (score >= 80) return '#4ECDC4';
  if (score >= 60) return '#4FC3F7';
  if (score >= 40) return '#FFE66D';
  if (score >= 20) return '#FF9F43';
  return '#FF6B6B';
}

export function UnifiedTimeline8456GH({
  data,
  className = '',
  onGranularityChange,
}: UnifiedTimelineProps) {
  const [granularity, setGranularity] = useState<Granularity>('D');
  const [showRaw, setShowRaw] = useState(false);
  const [selectedPoint, setSelectedPoint] = useState<TimeSeriesPoint | null>(null);

  // 计算统计信息
  const stats = useMemo(() => {
    if (data.length === 0) return null;

    const scores = data.map((d) => d.normalizedScore);
    const rawScores = data.map((d) => d.rawScore);
    const mean = scores.reduce((a, b) => a + b, 0) / scores.length;
    const variance =
      scores.reduce((sum, s) => sum + Math.pow(s - mean, 2), 0) / scores.length;
    const stdDev = Math.sqrt(variance);

    return {
      mean: mean.toFixed(1),
      stdDev: stdDev.toFixed(1),
      min: Math.min(...scores).toFixed(1),
      max: Math.max(...scores).toFixed(1),
      rawMin: Math.min(...rawScores).toFixed(1),
      rawMax: Math.max(...rawScores).toFixed(1),
      trend: scores[scores.length - 1] - scores[0] > 0 ? '上升' : '下降',
    };
  }, [data]);

  const handleGranularityChange = (g: Granularity) => {
    setGranularity(g);
    onGranularityChange?.(g);
  };

  if (data.length === 0) {
    return (
      <div className={`glass-card p-6 ${className}`}>
        <div className="text-center text-celestial-silver/60">
          暂无时间序列数据
        </div>
      </div>
    );
  }

  return (
    <div className={`glass-card p-6 ${className}`}>
      {/* 控制栏 */}
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-bold">📈 统一时间序列</h3>
        <div className="flex items-center gap-4">
          {/* 粒度切换 */}
          <div className="flex gap-1 bg-white/5 rounded-lg p-1">
            {(Object.keys(granularityLabels) as Granularity[]).map((g) => (
              <button
                key={g}
                onClick={() => handleGranularityChange(g)}
                className={`px-3 py-1 rounded text-sm font-medium transition-all ${
                  granularity === g
                    ? 'bg-cosmic-nova/30 text-cosmic-nova'
                    : 'text-celestial-silver/60 hover:text-celestial-silver'
                }`}
              >
                {g}
              </button>
            ))}
          </div>

          {/* 显示模式切换 */}
          <button
            onClick={() => setShowRaw(!showRaw)}
            className={`px-3 py-1 rounded-lg text-sm font-medium transition-all ${
              showRaw
                ? 'bg-cosmic-aurora/20 text-cosmic-aurora'
                : 'bg-white/5 text-celestial-silver/60'
            }`}
          >
            {showRaw ? '原始值' : '标准化'}
          </button>
        </div>
      </div>

      {/* 统计信息 */}
      {stats && (
        <div className="grid grid-cols-5 gap-4 mb-6">
          <StatCard label="平均值" value={stats.mean} />
          <StatCard label="标准差" value={stats.stdDev} />
          <StatCard label="最小值" value={stats.min} />
          <StatCard label="最大值" value={stats.max} />
          <StatCard label="趋势" value={stats.trend} />
        </div>
      )}

      {/* 图表区域 */}
      <div className="relative h-64 mb-4">
        <svg
          viewBox={`0 0 ${data.length * 20 + 40} 200`}
          className="w-full h-full"
          preserveAspectRatio="none"
        >
          {/* 网格线 */}
          {[0, 25, 50, 75, 100].map((y) => (
            <g key={y}>
              <line
                x1="40"
                y1={180 - y * 1.6}
                x2={data.length * 20 + 40}
                y2={180 - y * 1.6}
                stroke="rgba(255,255,255,0.1)"
                strokeDasharray="4"
              />
              <text
                x="35"
                y={185 - y * 1.6}
                fill="rgba(255,255,255,0.4)"
                fontSize="10"
                textAnchor="end"
              >
                {y}
              </text>
            </g>
          ))}

          {/* 数据线 */}
          <motion.polyline
            fill="none"
            stroke="url(#lineGradient)"
            strokeWidth="2"
            points={data
              .map((point, i) => {
                const score = showRaw
                  ? Math.min(100, Math.max(0, point.rawScore))
                  : point.normalizedScore;
                const x = 50 + i * 20;
                const y = 180 - score * 1.6;
                return `${x},${y}`;
              })
              .join(' ')}
            initial={{ pathLength: 0 }}
            animate={{ pathLength: 1 }}
            transition={{ duration: 1 }}
          />

          {/* 填充区域 */}
          <motion.polygon
            fill="url(#areaGradient)"
            points={[
              `50,180`,
              ...data.map((point, i) => {
                const score = showRaw
                  ? Math.min(100, Math.max(0, point.rawScore))
                  : point.normalizedScore;
                return `${50 + i * 20},${180 - score * 1.6}`;
              }),
              `${50 + (data.length - 1) * 20},180`,
            ].join(' ')}
            initial={{ opacity: 0 }}
            animate={{ opacity: 0.3 }}
            transition={{ delay: 0.5 }}
          />

          {/* 数据点 */}
          {data.map((point, i) => {
            const score = showRaw
              ? Math.min(100, Math.max(0, point.rawScore))
              : point.normalizedScore;
            const x = 50 + i * 20;
            const y = 180 - score * 1.6;

            return (
              <motion.circle
                key={i}
                cx={x}
                cy={y}
                r={selectedPoint?.timestamp === point.timestamp ? 6 : 4}
                fill={getScoreColor(point.normalizedScore)}
                stroke="rgba(255,255,255,0.5)"
                strokeWidth="1"
                style={{ cursor: 'pointer' }}
                whileHover={{ scale: 1.5 }}
                onClick={() => setSelectedPoint(point)}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: i * 0.02 }}
              />
            );
          })}

          {/* 渐变定义 */}
          <defs>
            <linearGradient id="lineGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#00D4FF" />
              <stop offset="100%" stopColor="#A855F7" />
            </linearGradient>
            <linearGradient id="areaGradient" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stopColor="#00D4FF" stopOpacity="0.4" />
              <stop offset="100%" stopColor="#00D4FF" stopOpacity="0" />
            </linearGradient>
          </defs>
        </svg>
      </div>

      {/* 选中点详情 */}
      {selectedPoint && (
        <motion.div
          className="p-4 rounded-lg bg-white/5 border border-cosmic-nova/30"
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <div className="flex items-center justify-between mb-3">
            <div className="font-medium">{selectedPoint.timestamp}</div>
            <button
              onClick={() => setSelectedPoint(null)}
              className="text-celestial-silver/60 hover:text-celestial-silver"
            >
              ✕
            </button>
          </div>
          <div className="grid grid-cols-2 gap-4 mb-3">
            <div>
              <div className="text-sm text-celestial-silver/60">标准化分数</div>
              <div
                className="text-2xl font-bold"
                style={{ color: getScoreColor(selectedPoint.normalizedScore) }}
              >
                {selectedPoint.normalizedScore.toFixed(1)}
              </div>
            </div>
            <div>
              <div className="text-sm text-celestial-silver/60">原始分数</div>
              <div className="text-2xl font-bold">
                {selectedPoint.rawScore.toFixed(1)}
              </div>
            </div>
          </div>

          {/* 因子分解 */}
          {selectedPoint.factors && selectedPoint.factors.length > 0 && (
            <div>
              <div className="text-sm text-celestial-silver/60 mb-2">
                影响因子
              </div>
              <div className="space-y-1 max-h-32 overflow-y-auto">
                {selectedPoint.factors.map((factor, i) => (
                  <div
                    key={i}
                    className="flex justify-between text-sm py-1 border-b border-white/5"
                  >
                    <span>{factor.name}</span>
                    <span
                      className={
                        factor.adjustment > 0 ? 'text-green-400' : 'text-red-400'
                      }
                    >
                      {factor.adjustment > 0 ? '+' : ''}
                      {factor.adjustment.toFixed(1)}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </motion.div>
      )}
    </div>
  );
}

// 统计卡片
function StatCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="p-3 rounded-lg bg-white/5 text-center">
      <div className="text-sm text-celestial-silver/60">{label}</div>
      <div className="text-lg font-bold text-cosmic-nova">{value}</div>
    </div>
  );
}

export default UnifiedTimeline8456GH;

