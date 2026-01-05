/**
 * FactorLifecycleChart7184BC - 因子生命周期曲线
 * 显示因子的正弦曲线强度变化：入相→峰值→离相
 */

import { useMemo } from 'react';
import { motion } from 'framer-motion';

interface FactorLifecycle {
  startTime: string;
  peakTime: string;
  endTime: string;
  duration: number;
}

interface FactorLifecycleChartProps {
  lifecycle: FactorLifecycle;
  currentTime?: Date;
  factorName?: string;
  baseValue?: number;
  color?: string;
  width?: number;
  height?: number;
  showLabels?: boolean;
  className?: string;
}

export function FactorLifecycleChart7184BC({
  lifecycle,
  currentTime = new Date(),
  factorName = '因子',
  baseValue = 1,
  color = '#00D4FF',
  width = 300,
  height = 100,
  showLabels = true,
  className = '',
}: FactorLifecycleChartProps) {
  const startTime = new Date(lifecycle.startTime);
  const endTime = new Date(lifecycle.endTime);
  // const peakTime = new Date(lifecycle.peakTime); // 峰值时间在正弦曲线中间

  // 计算当前进度和强度
  const { progress, currentStrength, phase } = useMemo(() => {
    const now = currentTime.getTime();
    const start = startTime.getTime();
    const end = endTime.getTime();
    const duration = end - start;

    if (now < start) {
      return { progress: 0, currentStrength: 0, phase: 'before' as const };
    }
    if (now > end) {
      return { progress: 1, currentStrength: 0, phase: 'after' as const };
    }

    const elapsed = now - start;
    const prog = elapsed / duration;
    const strength = Math.sin(Math.PI * prog);
    
    return {
      progress: prog,
      currentStrength: strength,
      phase: prog < 0.5 ? 'rising' as const : 'falling' as const,
    };
  }, [currentTime, startTime, endTime]);

  // 生成正弦曲线路径
  const curvePath = useMemo(() => {
    const padding = 20;
    const chartWidth = width - padding * 2;
    const chartHeight = height - padding * 2;
    const points: string[] = [];

    for (let i = 0; i <= 100; i++) {
      const x = padding + (i / 100) * chartWidth;
      const y = padding + chartHeight - Math.sin(Math.PI * (i / 100)) * chartHeight;
      points.push(`${i === 0 ? 'M' : 'L'} ${x} ${y}`);
    }

    return points.join(' ');
  }, [width, height]);

  // 当前位置
  const currentX = 20 + progress * (width - 40);
  const currentY = 20 + (height - 40) - currentStrength * (height - 40);

  // 格式化时间
  const formatTime = (date: Date) => {
    return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`;
  };

  // 格式化持续时间
  const formatDuration = (hours: number) => {
    if (hours < 1) return `${Math.round(hours * 60)}分钟`;
    if (hours < 24) return `${hours.toFixed(1)}小时`;
    if (hours < 168) return `${(hours / 24).toFixed(1)}天`;
    if (hours < 720) return `${(hours / 168).toFixed(1)}周`;
    return `${(hours / 720).toFixed(1)}月`;
  };

  const phaseLabels = {
    before: '未开始',
    rising: '增强中',
    falling: '减弱中',
    after: '已结束',
  };

  const phaseColors = {
    before: '#666',
    rising: '#4ECDC4',
    falling: '#FF9F43',
    after: '#666',
  };

  return (
    <div className={`${className}`}>
      {/* 标题 */}
      {showLabels && (
        <div className="flex items-center justify-between mb-2">
          <span className="text-sm font-medium">{factorName}</span>
          <span
            className="text-xs px-2 py-0.5 rounded-full"
            style={{ backgroundColor: `${phaseColors[phase]}30`, color: phaseColors[phase] }}
          >
            {phaseLabels[phase]}
          </span>
        </div>
      )}

      {/* 曲线图 */}
      <div className="relative" style={{ width, height }}>
        <svg width={width} height={height}>
          <defs>
            {/* 渐变填充 */}
            <linearGradient id={`lifecycle-gradient-${factorName}`} x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor={color} stopOpacity="0.1" />
              <stop offset={`${progress * 100}%`} stopColor={color} stopOpacity="0.4" />
              <stop offset={`${progress * 100}%`} stopColor={color} stopOpacity="0.1" />
              <stop offset="100%" stopColor={color} stopOpacity="0.05" />
            </linearGradient>
            {/* 发光效果 */}
            <filter id="currentGlow">
              <feGaussianBlur stdDeviation="3" result="blur" />
              <feMerge>
                <feMergeNode in="blur" />
                <feMergeNode in="SourceGraphic" />
              </feMerge>
            </filter>
          </defs>

          {/* 背景网格 */}
          <line x1="20" y1={height - 20} x2={width - 20} y2={height - 20} stroke="rgba(255,255,255,0.1)" />
          <line x1="20" y1="20" x2="20" y2={height - 20} stroke="rgba(255,255,255,0.1)" />

          {/* 峰值线 */}
          <line
            x1={width / 2}
            y1="20"
            x2={width / 2}
            y2={height - 20}
            stroke="rgba(255,255,255,0.2)"
            strokeDasharray="4,4"
          />

          {/* 填充区域 */}
          <path
            d={`${curvePath} L ${width - 20} ${height - 20} L 20 ${height - 20} Z`}
            fill={`url(#lifecycle-gradient-${factorName})`}
          />

          {/* 曲线 */}
          <motion.path
            d={curvePath}
            fill="none"
            stroke={color}
            strokeWidth="2"
            initial={{ pathLength: 0 }}
            animate={{ pathLength: 1 }}
            transition={{ duration: 1 }}
          />

          {/* 当前位置指示器 */}
          {phase !== 'before' && phase !== 'after' && (
            <>
              {/* 垂直线 */}
              <line
                x1={currentX}
                y1={currentY}
                x2={currentX}
                y2={height - 20}
                stroke={color}
                strokeWidth="1"
                strokeDasharray="2,2"
                opacity="0.5"
              />
              {/* 当前点 */}
              <motion.circle
                cx={currentX}
                cy={currentY}
                r={6}
                fill={color}
                stroke="white"
                strokeWidth="2"
                filter="url(#currentGlow)"
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: 0.5, type: 'spring' }}
              />
            </>
          )}
        </svg>

        {/* 时间标签 */}
        {showLabels && (
          <>
            <div className="absolute left-5 bottom-1 text-xs text-celestial-silver/50">
              {formatTime(startTime)}
            </div>
            <div className="absolute left-1/2 -translate-x-1/2 bottom-1 text-xs text-celestial-silver/50">
              峰值
            </div>
            <div className="absolute right-5 bottom-1 text-xs text-celestial-silver/50">
              {formatTime(endTime)}
            </div>
          </>
        )}

        {/* 当前强度标签 */}
        {phase !== 'before' && phase !== 'after' && (
          <motion.div
            className="absolute text-xs font-bold px-2 py-0.5 rounded bg-black/50"
            style={{
              left: currentX,
              top: currentY - 25,
              transform: 'translateX(-50%)',
              color,
            }}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.7 }}
          >
            {(currentStrength * baseValue).toFixed(2)}
          </motion.div>
        )}
      </div>

      {/* 底部信息 */}
      {showLabels && (
        <div className="flex justify-between text-xs text-celestial-silver/50 mt-2">
          <span>持续时间: {formatDuration(lifecycle.duration)}</span>
          <span>当前强度: {(currentStrength * 100).toFixed(0)}%</span>
        </div>
      )}
    </div>
  );
}

// 迷你版本
export function FactorLifecycleChartMini7184BC({
  lifecycle,
  currentTime,
  color = '#00D4FF',
  className = '',
}: Pick<FactorLifecycleChartProps, 'lifecycle' | 'currentTime' | 'color' | 'className'>) {
  return (
    <FactorLifecycleChart7184BC
      lifecycle={lifecycle}
      currentTime={currentTime}
      color={color}
      width={120}
      height={40}
      showLabels={false}
      className={className}
    />
  );
}

export default FactorLifecycleChart7184BC;

