/**
 * 实时五维值仪表盘组件
 * 组件命名规范：RealtimeDimensionDashboard + 7392 + WZ
 * 
 * 功能：
 * - 实时显示五个维度的分数（事业/关系/健康/财务/灵性）
 * - 每分钟自动刷新
 * - 动画效果展示分数变化
 * - 显示数据来源信息
 */

import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import type { BirthData } from '../../types';
import { DIMENSION_NAMES, DIMENSION_COLORS, DIMENSION_ICONS } from '../../utils/astro';

interface DimensionData {
  career: number;
  relationship: number;
  health: number;
  finance: number;
  spiritual: number;
}

interface RealtimeData {
  overall: number;
  dimensions: DimensionData;
  time: string;
  timestamp: Date;
}

interface RealtimeDimensionDashboardProps {
  birthData: BirthData;
  refreshInterval?: number; // 刷新间隔（毫秒），默认60000（1分钟）
  className?: string;
}

// 单个维度卡片
const DimensionCard = ({
  dimension,
  score,
  index,
}: {
  dimension: keyof DimensionData;
  score: number;
  index: number;
}) => {
  const color = DIMENSION_COLORS[dimension];
  const icon = DIMENSION_ICONS[dimension];
  const name = DIMENSION_NAMES[dimension];
  const percentage = Math.min(score, 100);

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9, y: 20 }}
      animate={{ opacity: 1, scale: 1, y: 0 }}
      transition={{ delay: index * 0.1, type: 'spring', stiffness: 200 }}
      className="relative group"
    >
      <div
        className="glass-card p-4 rounded-2xl border border-white/10 hover:border-white/20 
                   transition-all duration-300 hover:scale-105 cursor-default"
        style={{
          background: `linear-gradient(135deg, ${color}10, transparent)`,
        }}
      >
        {/* 背景光晕 */}
        <div
          className="absolute inset-0 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"
          style={{
            background: `radial-gradient(circle at center, ${color}15, transparent 70%)`,
          }}
        />

        {/* 图标和名称 */}
        <div className="flex items-center gap-2 mb-3 relative z-10">
          <span className="text-2xl">{icon}</span>
          <span className="text-sm font-medium text-white/80">{name}</span>
        </div>

        {/* 分数显示 */}
        <div className="flex items-baseline gap-1 mb-3 relative z-10">
          <motion.span
            key={score}
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-3xl font-bold"
            style={{ color }}
          >
            {Math.round(score)}
          </motion.span>
          <span className="text-white/40 text-sm">/100</span>
        </div>

        {/* 进度条 */}
        <div className="h-2 bg-white/10 rounded-full overflow-hidden relative z-10">
          <motion.div
            className="h-full rounded-full"
            style={{ backgroundColor: color }}
            initial={{ width: 0 }}
            animate={{ width: `${percentage}%` }}
            transition={{ duration: 1, ease: 'easeOut' }}
          />
        </div>

        {/* 状态指示 */}
        <div className="mt-2 flex items-center gap-1 relative z-10">
          <div
            className="w-2 h-2 rounded-full animate-pulse"
            style={{ backgroundColor: color }}
          />
          <span className="text-xs text-white/50">
            {score >= 80 ? '极佳' : score >= 60 ? '良好' : score >= 40 ? '一般' : '待提升'}
          </span>
        </div>
      </div>
    </motion.div>
  );
};

// 综合分数圆环
const OverallScoreRing = ({ score }: { score: number }) => {
  const radius = 60;
  const strokeWidth = 8;
  const circumference = 2 * Math.PI * radius;
  const progress = (score / 100) * circumference;
  const gradientId = `overall-gradient-${Date.now()}`;

  return (
    <div className="relative flex items-center justify-center">
      <svg width="150" height="150" className="transform -rotate-90">
        <defs>
          <linearGradient id={gradientId} x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#00D4FF" />
            <stop offset="50%" stopColor="#A855F7" />
            <stop offset="100%" stopColor="#FF6B9D" />
          </linearGradient>
        </defs>
        {/* 背景环 */}
        <circle
          cx="75"
          cy="75"
          r={radius}
          fill="none"
          stroke="rgba(255,255,255,0.1)"
          strokeWidth={strokeWidth}
        />
        {/* 进度环 */}
        <motion.circle
          cx="75"
          cy="75"
          r={radius}
          fill="none"
          stroke={`url(#${gradientId})`}
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          strokeDasharray={circumference}
          initial={{ strokeDashoffset: circumference }}
          animate={{ strokeDashoffset: circumference - progress }}
          transition={{ duration: 1.5, ease: 'easeOut' }}
        />
      </svg>
      {/* 中心分数 */}
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <motion.span
          key={score}
          initial={{ opacity: 0, scale: 0.5 }}
          animate={{ opacity: 1, scale: 1 }}
          className="text-4xl font-bold bg-gradient-to-r from-[#00D4FF] via-[#A855F7] to-[#FF6B9D] bg-clip-text text-transparent"
        >
          {Math.round(score)}
        </motion.span>
        <span className="text-xs text-white/50">综合运势</span>
      </div>
    </div>
  );
};

export function RealtimeDimensionDashboard7392WZ({
  birthData,
  refreshInterval = 60000,
  className = '',
}: RealtimeDimensionDashboardProps) {
  const [realtimeData, setRealtimeData] = useState<RealtimeData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 格式化本地 ISO 时间
  const formatLocalISO = (date: Date, timezone: number = 8) => {
    const offsetHours = Math.floor(Math.abs(timezone));
    const offsetMins = Math.round((Math.abs(timezone) % 1) * 60);
    const sign = timezone >= 0 ? '+' : '-';
    const pad = (n: number) => n.toString().padStart(2, '0');

    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}${sign}${pad(offsetHours)}:${pad(offsetMins)}`;
  };

  // 获取实时数据
  const fetchRealtimeData = useCallback(async () => {
    try {
      setError(null);
      
      // 计算当前小时的时间范围
      const nowUtc = Date.now();
      const userTimezoneOffset = birthData.timezone * 60 * 60 * 1000;
      const userLocalTime = new Date(nowUtc + userTimezoneOffset + new Date().getTimezoneOffset() * 60 * 1000);

      const start = new Date(userLocalTime);
      start.setMinutes(0, 0, 0);
      const end = new Date(start);
      end.setHours(end.getHours() + 1);

      const startStr = formatLocalISO(start, birthData.timezone);
      const endStr = formatLocalISO(end, birthData.timezone);

      const response = await fetch('http://localhost:8080/api/calc/time-series', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          birthData: {
            year: birthData.year,
            month: birthData.month,
            day: birthData.day,
            hour: birthData.hour,
            minute: birthData.minute,
            latitude: birthData.latitude,
            longitude: birthData.longitude,
            timezone: birthData.timezone,
          },
          start: startStr,
          end: endStr,
          granularity: 'hour',
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to fetch realtime data');
      }

      const data = await response.json();
      
      if (data.points && data.points.length > 0) {
        const point = data.points[0];
        const now = new Date();
        
        setRealtimeData({
          overall: point.display || 50,
          dimensions: point.dimensions || {
            career: 50,
            relationship: 50,
            health: 50,
            finance: 50,
            spiritual: 50,
          },
          time: now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
          timestamp: now,
        });
      }
      
      setLoading(false);
    } catch (err) {
      console.error('获取实时数据失败:', err);
      setError(err instanceof Error ? err.message : '获取数据失败');
      setLoading(false);
    }
  }, [birthData]);

  // 初始加载和定时刷新
  useEffect(() => {
    fetchRealtimeData();
    
    const interval = setInterval(fetchRealtimeData, refreshInterval);
    return () => clearInterval(interval);
  }, [fetchRealtimeData, refreshInterval]);

  // 加载状态
  if (loading && !realtimeData) {
    return (
      <div className={`glass-card p-8 ${className}`}>
        <div className="flex items-center justify-center gap-3">
          <div className="w-8 h-8 border-2 border-[#00D4FF] border-t-transparent rounded-full animate-spin" />
          <span className="text-white/60">正在获取实时数据...</span>
        </div>
      </div>
    );
  }

  // 错误状态
  if (error && !realtimeData) {
    return (
      <div className={`glass-card p-8 ${className}`}>
        <div className="text-center">
          <span className="text-red-400">❌ {error}</span>
          <button
            onClick={fetchRealtimeData}
            className="ml-4 text-[#00D4FF] hover:underline"
          >
            重试
          </button>
        </div>
      </div>
    );
  }

  if (!realtimeData) return null;

  const dimensions = Object.keys(realtimeData.dimensions) as (keyof DimensionData)[];

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className={`glass-card p-6 rounded-3xl ${className}`}
    >
      {/* 标题栏 */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className="w-3 h-3 rounded-full bg-green-500 animate-pulse" />
          <h2 className="text-xl font-bold text-white">⚡ 实时五维运势</h2>
        </div>
        <div className="flex items-center gap-2 text-xs text-white/50">
          <span>🔄 {realtimeData.time} 更新</span>
          <span className="px-2 py-1 bg-[#00D4FF]/10 text-[#00D4FF] rounded-full">
            Swiss Ephemeris
          </span>
        </div>
      </div>

      {/* 主内容区 */}
      <div className="grid lg:grid-cols-[auto_1fr] gap-8 items-center">
        {/* 左侧：综合分数圆环 */}
        <div className="flex justify-center">
          <OverallScoreRing score={realtimeData.overall} />
        </div>

        {/* 右侧：五维分数卡片 */}
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
          {dimensions.map((dimension, index) => (
            <DimensionCard
              key={dimension}
              dimension={dimension}
              score={realtimeData.dimensions[dimension]}
              index={index}
            />
          ))}
        </div>
      </div>

      {/* 底部信息 */}
      <div className="mt-6 pt-4 border-t border-white/10 flex items-center justify-between text-xs text-white/40">
        <div className="flex items-center gap-4">
          <span>📡 数据来源: Swiss Ephemeris (高精度)</span>
          <span>🔄 每分钟自动刷新</span>
        </div>
        <div className="flex items-center gap-2">
          <span>
            {realtimeData.overall >= 80 ? '🌟' : realtimeData.overall >= 60 ? '✨' : realtimeData.overall >= 40 ? '💫' : '🌙'}
          </span>
          <span>
            当前运势
            {realtimeData.overall >= 80 ? '极佳' : realtimeData.overall >= 60 ? '良好' : realtimeData.overall >= 40 ? '平稳' : '需注意'}
          </span>
        </div>
      </div>
    </motion.div>
  );
}

export default RealtimeDimensionDashboard7392WZ;

