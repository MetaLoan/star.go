/**
 * FactorGanttChart9527XZ - 影响因子横道图（甘特图）
 * 以天为横坐标，展示多个因子的生命周期时间线
 */

import { useMemo, useState } from 'react';
import { motion } from 'framer-motion';

interface FactorLifecycle {
  startTime: string;
  peakTime: string;
  endTime: string;
  duration: number;
}

interface InfluenceFactor {
  id: string;
  name: string;
  type: string;
  value: number;
  adjustment: number;
  isPositive: boolean;
  lifecycle?: FactorLifecycle;
  remainingDays?: number;
}

interface FactorGanttChartProps {
  factors: InfluenceFactor[];
  startDate?: Date;
  endDate?: Date;
  currentTime?: Date;
  daysToShow?: number;
  height?: number;
  rowHeight?: number;
  showLegend?: boolean;
  onFactorClick?: (factor: InfluenceFactor) => void;
  className?: string;
}

// 因子类型颜色映射（覆盖所有26种因子类型）
const FACTOR_TYPE_COLORS: Record<string, string> = {
  // ===== 基础因子 =====
  aspectPhase: '#00D4FF',      // 青色 - 相位
  aspectOrb: '#00B4D8',        // 深青色 - 相位容许度
  dignity: '#FFD700',          // 金色 - 尊贵度
  retrograde: '#FF6B9D',       // 粉色 - 逆行
  lunarPhase: '#A855F7',       // 紫色 - 月相
  planetaryHour: '#4ECDC4',    // 绿色 - 行星时
  profectionLord: '#FF9F43',   // 橙色 - 年主星
  voidOfCourse: '#666666',     // 灰色 - 月空
  outerPlanet: '#6366F1',      // 靛蓝色 - 外行星
  personal: '#EC4899',         // 玫红色 - 个人
  custom: '#8B5CF6',           // 紫罗兰 - 自定义
  
  // ===== 日月食与交点 =====
  eclipse: '#DC2626',          // 深红色 - 日月食
  lunarNode: '#7C3AED',        // 紫色 - 月交点
  
  // ===== 行星状态 =====
  combustion: '#F97316',       // 橙红色 - 燃烧
  station: '#FBBF24',          // 琥珀色 - 停滞
  reception: '#10B981',        // 翠绿色 - 互容
  
  // ===== 恒星与特殊点 =====
  fixedStar: '#F0E68C',        // 卡其色 - 恒星
  arabicPart: '#20B2AA',       // 浅海绿 - 阿拉伯点
  midpoint: '#9CA3AF',         // 灰色 - 中点
  antiscion: '#78716C',        // 石灰色 - 反生点
  
  // ===== 界限与分度 =====
  term: '#D4A574',             // 褐色 - 界限
  decan: '#C4A484',            // 浅褐色 - 十度面
  
  // ===== 推运技术 =====
  solarArc: '#FF4500',         // 橙红色 - 太阳弧
  primary: '#8B0000',          // 暗红色 - 主限推进
  firdaria: '#4B0082',         // 靛青色 - 法达
  zodiacal: '#00CED1',         // 暗绿松色 - 黄道释放
};

// 格式化日期为简短形式
function formatDateShort(date: Date): string {
  return `${date.getMonth() + 1}/${date.getDate()}`;
}

// 格式化日期为完整形式
function formatDateFull(date: Date): string {
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hours = String(date.getHours()).padStart(2, '0');
  const mins = String(date.getMinutes()).padStart(2, '0');
  return `${month}/${day} ${hours}:${mins}`;
}

export function FactorGanttChart9527XZ({
  factors,
  startDate,
  endDate,
  currentTime = new Date(),
  daysToShow = 14,
  height,
  rowHeight = 36,
  showLegend = true,
  onFactorClick,
  className = '',
}: FactorGanttChartProps) {
  const [hoveredFactor, setHoveredFactor] = useState<string | null>(null);
  const [tooltipData, setTooltipData] = useState<{
    factor: InfluenceFactor;
    x: number;
    y: number;
  } | null>(null);

  // 计算时间范围
  const { timeRange, dayMarkers, totalDays } = useMemo(() => {
    let rangeStart: Date;
    let rangeEnd: Date;

    if (startDate && endDate) {
      rangeStart = new Date(startDate);
      rangeEnd = new Date(endDate);
    } else {
      // 默认显示从3天前到未来11天（共14天）
      rangeStart = new Date(currentTime);
      rangeStart.setDate(rangeStart.getDate() - 3);
      rangeStart.setHours(0, 0, 0, 0);
      
      rangeEnd = new Date(rangeStart);
      rangeEnd.setDate(rangeEnd.getDate() + daysToShow);
    }

    // 生成日期标记
    const markers: { date: Date; label: string; isToday: boolean }[] = [];
    const current = new Date(rangeStart);
    const today = new Date(currentTime);
    today.setHours(0, 0, 0, 0);

    while (current <= rangeEnd) {
      const isToday = current.getTime() === today.getTime();
      markers.push({
        date: new Date(current),
        label: formatDateShort(current),
        isToday,
      });
      current.setDate(current.getDate() + 1);
    }

    const days = (rangeEnd.getTime() - rangeStart.getTime()) / (24 * 60 * 60 * 1000);

    return {
      timeRange: { start: rangeStart, end: rangeEnd },
      dayMarkers: markers,
      totalDays: days,
    };
  }, [startDate, endDate, currentTime, daysToShow]);

  // 过滤有生命周期的因子并排序
  const validFactors = useMemo(() => {
    return factors
      .filter(f => f.lifecycle && f.lifecycle.startTime && f.lifecycle.endTime)
      .sort((a, b) => {
        // 先按类型排序，再按开始时间排序
        if (a.type !== b.type) {
          return a.type.localeCompare(b.type);
        }
        return new Date(a.lifecycle!.startTime).getTime() - 
               new Date(b.lifecycle!.startTime).getTime();
      });
  }, [factors]);

  // 计算图表尺寸
  const chartWidth = 800;
  const labelWidth = 180;
  const timelineWidth = chartWidth - labelWidth;
  const computedHeight = height || Math.max(300, validFactors.length * rowHeight + 80);

  // 将时间转换为X坐标
  const timeToX = (time: Date): number => {
    const totalMs = timeRange.end.getTime() - timeRange.start.getTime();
    const elapsedMs = time.getTime() - timeRange.start.getTime();
    return labelWidth + (elapsedMs / totalMs) * timelineWidth;
  };

  // 当前时间线位置
  const nowX = timeToX(currentTime);
  const isNowInRange = nowX >= labelWidth && nowX <= chartWidth;

  // 计算因子条形
  const factorBars = useMemo(() => {
    return validFactors.map((factor, index) => {
      const lifecycle = factor.lifecycle!;
      const start = new Date(lifecycle.startTime);
      const end = new Date(lifecycle.endTime);
      const peak = new Date(lifecycle.peakTime);

      // 裁剪到可视范围
      const clippedStart = start < timeRange.start ? timeRange.start : start;
      const clippedEnd = end > timeRange.end ? timeRange.end : end;

      // 如果完全不在范围内，跳过
      if (clippedStart >= timeRange.end || clippedEnd <= timeRange.start) {
        return null;
      }

      const x1 = timeToX(clippedStart);
      const x2 = timeToX(clippedEnd);
      const peakX = timeToX(peak);
      const y = 60 + index * rowHeight;
      const barHeight = rowHeight - 8;

      // 计算当前强度（基于正弦曲线）
      const totalDuration = end.getTime() - start.getTime();
      const elapsed = currentTime.getTime() - start.getTime();
      let progress = elapsed / totalDuration;
      progress = Math.max(0, Math.min(1, progress));
      const strength = Math.sin(Math.PI * progress);

      return {
        factor,
        x1,
        x2,
        peakX,
        y,
        barHeight,
        width: x2 - x1,
        strength,
        isActive: currentTime >= start && currentTime <= end,
        isFuture: currentTime < start,
        isPast: currentTime > end,
        isClippedStart: start < timeRange.start,
        isClippedEnd: end > timeRange.end,
      };
    }).filter(Boolean);
  }, [validFactors, timeRange, currentTime, rowHeight]);

  // 处理鼠标悬停
  const handleMouseEnter = (factor: InfluenceFactor, event: React.MouseEvent) => {
    setHoveredFactor(factor.id);
    setTooltipData({
      factor,
      x: event.clientX,
      y: event.clientY,
    });
  };

  // 处理鼠标移动 - 让tooltip跟随鼠标
  const handleMouseMove = (factor: InfluenceFactor, event: React.MouseEvent) => {
    setTooltipData({
      factor,
      x: event.clientX,
      y: event.clientY,
    });
  };

  const handleMouseLeave = () => {
    setHoveredFactor(null);
    setTooltipData(null);
  };

  // 获取因子类型标签
  const getTypeLabel = (type: string): string => {
    const labels: Record<string, string> = {
      aspectPhase: '相位',
      dignity: '尊贵',
      retrograde: '逆行',
      lunarPhase: '月相',
      planetaryHour: '行星时',
      profectionLord: '年主星',
      voidOfCourse: '月空',
      custom: '自定义',
    };
    return labels[type] || type;
  };

  if (validFactors.length === 0) {
    return (
      <div className={`glass-card p-6 ${className}`}>
        <div className="text-center text-white/50 py-8">
          <p className="text-lg mb-2">📊 影响因子横道图</p>
          <p className="text-sm">暂无带生命周期的因子数据</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`glass-card p-4 ${className}`}>
      {/* 标题 */}
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-medium text-white">📊 影响因子时间线</h3>
        <div className="text-sm text-white/50">
          {formatDateShort(timeRange.start)} - {formatDateShort(timeRange.end)}
          <span className="ml-2 text-white/30">({validFactors.length} 个因子)</span>
        </div>
      </div>

      {/* 图例 */}
      {showLegend && (
        <div className="flex flex-wrap gap-3 mb-4 pb-3 border-b border-white/10">
          {Object.entries(FACTOR_TYPE_COLORS).map(([type, color]) => {
            const count = validFactors.filter(f => f.type === type).length;
            if (count === 0) return null;
            return (
              <div key={type} className="flex items-center gap-1.5 text-xs">
                <div
                  className="w-3 h-3 rounded"
                  style={{ backgroundColor: color }}
                />
                <span className="text-white/60">{getTypeLabel(type)} ({count})</span>
              </div>
            );
          })}
        </div>
      )}

      {/* 甘特图主体 */}
      <div className="overflow-x-auto">
        <svg
          width={chartWidth}
          height={computedHeight}
          className="select-none"
        >
          <defs>
            {/* 渐变定义 */}
            {Object.entries(FACTOR_TYPE_COLORS).map(([type, color]) => (
              <linearGradient
                key={type}
                id={`factor-gradient-${type}`}
                x1="0%"
                y1="0%"
                x2="100%"
                y2="0%"
              >
                <stop offset="0%" stopColor={color} stopOpacity="0.6" />
                <stop offset="50%" stopColor={color} stopOpacity="1" />
                <stop offset="100%" stopColor={color} stopOpacity="0.6" />
              </linearGradient>
            ))}
            {/* 发光效果 */}
            <filter id="glow">
              <feGaussianBlur stdDeviation="2" result="blur" />
              <feMerge>
                <feMergeNode in="blur" />
                <feMergeNode in="SourceGraphic" />
              </feMerge>
            </filter>
          </defs>

          {/* 背景网格 */}
          <g className="grid">
            {dayMarkers.map((marker, i) => {
              const x = labelWidth + (i / totalDays) * timelineWidth;
              return (
                <g key={i}>
                  <line
                    x1={x}
                    y1={50}
                    x2={x}
                    y2={computedHeight - 20}
                    stroke={marker.isToday ? 'rgba(0, 212, 255, 0.3)' : 'rgba(255, 255, 255, 0.1)'}
                    strokeWidth={marker.isToday ? 2 : 1}
                    strokeDasharray={marker.isToday ? undefined : '4,4'}
                  />
                  <text
                    x={x}
                    y={40}
                    textAnchor="middle"
                    className="text-xs"
                    fill={marker.isToday ? '#00D4FF' : 'rgba(255, 255, 255, 0.5)'}
                    fontWeight={marker.isToday ? 'bold' : 'normal'}
                  >
                    {marker.label}
                  </text>
                  {marker.isToday && (
                    <text
                      x={x}
                      y={25}
                      textAnchor="middle"
                      className="text-xs"
                      fill="#00D4FF"
                    >
                      今天
                    </text>
                  )}
                </g>
              );
            })}
          </g>

          {/* 因子标签列 */}
          <g className="labels">
            {factorBars.map((bar, i) => {
              if (!bar) return null;
              const { factor, y, barHeight } = bar;
              const color = FACTOR_TYPE_COLORS[factor.type] || '#888';
              const isHovered = hoveredFactor === factor.id;

              return (
                <g key={factor.id}>
                  {/* 行背景 */}
                  <rect
                    x={0}
                    y={y - 2}
                    width={labelWidth - 8}
                    height={barHeight + 4}
                    fill={isHovered ? 'rgba(255, 255, 255, 0.05)' : 'transparent'}
                    rx={4}
                  />
                  {/* 类型指示器 */}
                  <rect
                    x={4}
                    y={y + 4}
                    width={4}
                    height={barHeight - 8}
                    fill={color}
                    rx={2}
                  />
                  {/* 因子名称 */}
                  <text
                    x={14}
                    y={y + barHeight / 2 + 1}
                    dominantBaseline="middle"
                    className="text-xs"
                    fill={bar.isActive ? '#fff' : 'rgba(255, 255, 255, 0.6)'}
                    fontWeight={bar.isActive ? 'bold' : 'normal'}
                  >
                    {factor.name.length > 20
                      ? factor.name.slice(0, 18) + '...'
                      : factor.name}
                  </text>
                  {/* 正/负指示 */}
                  <text
                    x={labelWidth - 12}
                    y={y + barHeight / 2 + 1}
                    dominantBaseline="middle"
                    textAnchor="end"
                    className="text-xs"
                    fill={factor.isPositive ? '#4ADE80' : '#FF6B9D'}
                  >
                    {factor.isPositive ? '+' : '-'}
                    {Math.abs(factor.adjustment).toFixed(1)}
                  </text>
                </g>
              );
            })}
          </g>

          {/* 因子时间条 */}
          <g className="bars">
            {factorBars.map((bar) => {
              if (!bar) return null;
              const { factor, x1, x2, peakX, y, barHeight, width, strength, isActive } = bar;
              const color = FACTOR_TYPE_COLORS[factor.type] || '#888';
              const isHovered = hoveredFactor === factor.id;

              return (
                <motion.g
                  key={factor.id}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.3 }}
                  onMouseEnter={(e) => handleMouseEnter(factor, e as unknown as React.MouseEvent)}
                  onMouseMove={(e) => handleMouseMove(factor, e as unknown as React.MouseEvent)}
                  onMouseLeave={handleMouseLeave}
                  onClick={() => onFactorClick?.(factor)}
                  style={{ cursor: onFactorClick ? 'pointer' : 'default' }}
                >
                  {/* 条形背景 */}
                  <rect
                    x={x1}
                    y={y + 2}
                    width={Math.max(2, width)}
                    height={barHeight - 4}
                    fill={`url(#factor-gradient-${factor.type})`}
                    opacity={bar.isPast ? 0.3 : bar.isFuture ? 0.5 : 0.8}
                    rx={4}
                    filter={isHovered ? 'url(#glow)' : undefined}
                  />

                  {/* 峰值标记 */}
                  {peakX >= x1 && peakX <= x2 && (
                    <circle
                      cx={peakX}
                      cy={y + barHeight / 2}
                      r={isActive ? 5 : 3}
                      fill={color}
                      stroke="white"
                      strokeWidth={1}
                    />
                  )}

                  {/* 当前进度指示器（如果活跃） */}
                  {isActive && (
                    <motion.line
                      x1={nowX}
                      y1={y + 4}
                      x2={nowX}
                      y2={y + barHeight - 4}
                      stroke="#fff"
                      strokeWidth={2}
                      initial={{ opacity: 0 }}
                      animate={{ opacity: [0.5, 1, 0.5] }}
                      transition={{ duration: 1.5, repeat: Infinity }}
                    />
                  )}

                  {/* 左边界裁剪指示 */}
                  {bar.isClippedStart && (
                    <polygon
                      points={`${x1},${y + barHeight / 2} ${x1 + 6},${y + 4} ${x1 + 6},${y + barHeight - 4}`}
                      fill="rgba(255, 255, 255, 0.3)"
                    />
                  )}

                  {/* 右边界裁剪指示 */}
                  {bar.isClippedEnd && (
                    <polygon
                      points={`${x2},${y + barHeight / 2} ${x2 - 6},${y + 4} ${x2 - 6},${y + barHeight - 4}`}
                      fill="rgba(255, 255, 255, 0.3)"
                    />
                  )}
                </motion.g>
              );
            })}
          </g>

          {/* 当前时间线 */}
          {isNowInRange && (
            <g>
              <motion.line
                x1={nowX}
                y1={50}
                x2={nowX}
                y2={computedHeight - 20}
                stroke="#FF6B9D"
                strokeWidth={2}
                strokeDasharray="4,2"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
              />
              <motion.circle
                cx={nowX}
                cy={50}
                r={4}
                fill="#FF6B9D"
                initial={{ scale: 0 }}
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ duration: 2, repeat: Infinity }}
              />
            </g>
          )}
        </svg>
      </div>

      {/* 悬停提示 */}
      {tooltipData && (
        <motion.div
          className="fixed z-50 bg-black/90 border border-white/20 rounded-lg p-3 shadow-xl pointer-events-none"
          style={{
            left: tooltipData.x + 10,
            top: tooltipData.y + 10,
            maxWidth: 300,
          }}
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
        >
          <div className="text-sm">
            <div className="flex items-center gap-2 mb-2">
              <div
                className="w-2 h-2 rounded-full"
                style={{ backgroundColor: FACTOR_TYPE_COLORS[tooltipData.factor.type] }}
              />
              <span className="font-medium text-white">{tooltipData.factor.name}</span>
            </div>
            <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
              <span className="text-white/50">类型:</span>
              <span className="text-white/80">{getTypeLabel(tooltipData.factor.type)}</span>
              
              <span className="text-white/50">调整值:</span>
              <span className={tooltipData.factor.isPositive ? 'text-green-400' : 'text-red-400'}>
                {tooltipData.factor.isPositive ? '+' : ''}{tooltipData.factor.adjustment.toFixed(2)}
              </span>
              
              {tooltipData.factor.lifecycle && (
                <>
                  <span className="text-white/50">开始:</span>
                  <span className="text-white/80">
                    {formatDateFull(new Date(tooltipData.factor.lifecycle.startTime))}
                  </span>
                  
                  <span className="text-white/50">峰值:</span>
                  <span className="text-white/80">
                    {formatDateFull(new Date(tooltipData.factor.lifecycle.peakTime))}
                  </span>
                  
                  <span className="text-white/50">结束:</span>
                  <span className="text-white/80">
                    {formatDateFull(new Date(tooltipData.factor.lifecycle.endTime))}
                  </span>
                  
                  <span className="text-white/50">持续:</span>
                  <span className="text-white/80">
                    {tooltipData.factor.lifecycle.duration < 24
                      ? `${tooltipData.factor.lifecycle.duration.toFixed(1)}小时`
                      : `${(tooltipData.factor.lifecycle.duration / 24).toFixed(1)}天`}
                  </span>
                </>
              )}
              
              {tooltipData.factor.remainingDays !== undefined && (
                <>
                  <span className="text-white/50">剩余:</span>
                  <span className="text-cyan-400">
                    {tooltipData.factor.remainingDays.toFixed(2)}天
                  </span>
                </>
              )}
            </div>
          </div>
        </motion.div>
      )}
    </div>
  );
}

export default FactorGanttChart9527XZ;
