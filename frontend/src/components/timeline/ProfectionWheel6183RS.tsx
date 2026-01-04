/**
 * 年限法轮盘组件
 * 组件命名规范：ProfectionWheel + 6183 + RS
 */

import { useMemo } from 'react';
import { motion } from 'framer-motion';
import type { AnnualProfection } from '../../types';
import {
  PLANET_SYMBOLS,
  PLANET_COLORS,
  PLANET_NAMES,
  HOUSE_NAMES,
  HOUSE_THEMES,
  polarToCartesian,
} from '../../utils/astro';

interface ProfectionWheelProps {
  profections: AnnualProfection[];
  currentAge: number;
  size?: number;
  onAgeClick?: (age: number) => void;
  className?: string;
}

export function ProfectionWheel6183RS({
  profections,
  currentAge,
  size = 400,
  onAgeClick,
  className = '',
}: ProfectionWheelProps) {
  const center = size / 2;
  const outerRadius = size * 0.45;
  const innerRadius = size * 0.25;
  const labelRadius = size * 0.35;

  // 每个宫位占 30 度
  const houseAngle = 30;

  // 渲染宫位扇形
  const renderHouseSegments = useMemo(() => {
    return Array.from({ length: 12 }, (_, i) => {
      const houseNumber = i + 1;
      const startAngle = i * houseAngle - 90; // 从顶部开始
      const endAngle = startAngle + houseAngle;
      
      // 计算扇形路径
      const outerStart = polarToCartesian(center, center, outerRadius, startAngle);
      const outerEnd = polarToCartesian(center, center, outerRadius, endAngle);
      const innerStart = polarToCartesian(center, center, innerRadius, startAngle);
      const innerEnd = polarToCartesian(center, center, innerRadius, endAngle);
      
      const path = [
        `M ${outerStart.x} ${outerStart.y}`,
        `A ${outerRadius} ${outerRadius} 0 0 1 ${outerEnd.x} ${outerEnd.y}`,
        `L ${innerEnd.x} ${innerEnd.y}`,
        `A ${innerRadius} ${innerRadius} 0 0 0 ${innerStart.x} ${innerStart.y}`,
        'Z',
      ].join(' ');
      
      // 宫位标签位置
      const labelAngle = startAngle + houseAngle / 2;
      const labelPos = polarToCartesian(center, center, labelRadius, labelAngle);
      
      // 当前年龄是否在此宫位
      const isCurrentHouse = profections.find(
        p => p.age === currentAge && p.house === houseNumber
      );
      
      // 颜色渐变（根据宫位性质）
      const houseColors = [
        '#ef4444', // 1 - 命宫 (火)
        '#22c55e', // 2 - 财帛 (土)
        '#eab308', // 3 - 兄弟 (风)
        '#3b82f6', // 4 - 田宅 (水)
        '#f97316', // 5 - 子女 (火)
        '#84cc16', // 6 - 奴仆 (土)
        '#ec4899', // 7 - 夫妻 (风)
        '#dc2626', // 8 - 疾厄 (水)
        '#a855f7', // 9 - 迁移 (火)
        '#71717a', // 10 - 官禄 (土)
        '#06b6d4', // 11 - 福德 (风)
        '#8b5cf6', // 12 - 玄秘 (水)
      ];

      return (
        <g key={houseNumber}>
          {/* 宫位扇形 */}
          <motion.path
            d={path}
            fill={isCurrentHouse ? `${houseColors[i]}40` : 'rgba(255,255,255,0.03)'}
            stroke={houseColors[i]}
            strokeWidth={isCurrentHouse ? 2 : 1}
            strokeOpacity={isCurrentHouse ? 1 : 0.3}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: i * 0.05 }}
            whileHover={{ fill: `${houseColors[i]}20` }}
          />
          
          {/* 宫位编号 */}
          <text
            x={labelPos.x}
            y={labelPos.y}
            fill={isCurrentHouse ? houseColors[i] : 'rgba(255,255,255,0.6)'}
            fontSize={isCurrentHouse ? 14 : 12}
            fontWeight={isCurrentHouse ? 'bold' : 'normal'}
            textAnchor="middle"
            dominantBaseline="central"
          >
            {houseNumber}
          </text>
        </g>
      );
    });
  }, [center, outerRadius, innerRadius, labelRadius, profections, currentAge]);

  // 渲染年龄标记
  const renderAgeMarkers = useMemo(() => {
    // 显示当前年龄前后各 6 年
    const visibleProfections = profections.filter(
      p => Math.abs(p.age - currentAge) <= 6
    );

    return visibleProfections.map((prof, index) => {
      // 计算年龄对应的角度（基于激活宫位）
      const baseAngle = (prof.house - 1) * houseAngle - 90;
      // 在宫位内偏移（避免重叠）
      const offsetAngle = baseAngle + houseAngle / 2;
      
      // 半径（当前年龄在外圈，其他在内圈）
      const markerRadius = prof.age === currentAge 
        ? outerRadius + 15 
        : outerRadius + 8 + (index % 2) * 8;
      
      const pos = polarToCartesian(center, center, markerRadius, offsetAngle);
      const isCurrent = prof.age === currentAge;

      return (
        <motion.g
          key={prof.age}
          initial={{ scale: 0, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ delay: 0.5 + index * 0.03 }}
          style={{ cursor: onAgeClick ? 'pointer' : 'default' }}
          onClick={() => onAgeClick?.(prof.age)}
        >
          {/* 年龄气泡 */}
          <circle
            cx={pos.x}
            cy={pos.y}
            r={isCurrent ? 16 : 12}
            fill={isCurrent ? PLANET_COLORS[prof.lordOfYear] : 'rgba(255,255,255,0.1)'}
            stroke={PLANET_COLORS[prof.lordOfYear]}
            strokeWidth={isCurrent ? 2 : 1}
          />
          
          {/* 年龄文字 */}
          <text
            x={pos.x}
            y={pos.y}
            fill={isCurrent ? '#fff' : 'rgba(255,255,255,0.8)'}
            fontSize={isCurrent ? 10 : 8}
            fontWeight={isCurrent ? 'bold' : 'normal'}
            textAnchor="middle"
            dominantBaseline="central"
          >
            {prof.age}
          </text>
        </motion.g>
      );
    });
  }, [profections, currentAge, center, outerRadius, onAgeClick]);

  // 当前年限法信息
  const currentProfection = profections.find(p => p.age === currentAge);

  return (
    <div className={`glass-card p-4 ${className}`}>
      <h3 className="text-lg font-medium text-white mb-4 flex items-center gap-2">
        <span>🔮</span>
        年限法轮盘
      </h3>

      <div className="flex flex-col lg:flex-row gap-6">
        {/* SVG 轮盘 */}
        <div className="flex justify-center">
          <svg
            width={size}
            height={size}
            viewBox={`0 0 ${size} ${size}`}
            className="overflow-visible"
          >
            <defs>
              <radialGradient id="wheelBg" cx="50%" cy="50%" r="50%">
                <stop offset="0%" stopColor="rgba(26, 26, 46, 0.5)" />
                <stop offset="100%" stopColor="rgba(10, 10, 15, 0.8)" />
              </radialGradient>
            </defs>

            {/* 背景圆 */}
            <circle
              cx={center}
              cy={center}
              r={outerRadius}
              fill="url(#wheelBg)"
            />

            {/* 宫位分区 */}
            {renderHouseSegments}

            {/* 内圆 */}
            <circle
              cx={center}
              cy={center}
              r={innerRadius}
              fill="rgba(10, 10, 15, 0.9)"
              stroke="rgba(255,255,255,0.1)"
              strokeWidth="1"
            />

            {/* 中心年龄显示 */}
            {currentProfection && (
              <g>
                <text
                  x={center}
                  y={center - 15}
                  fill="white"
                  fontSize="24"
                  fontWeight="bold"
                  textAnchor="middle"
                >
                  {currentAge}岁
                </text>
                <text
                  x={center}
                  y={center + 10}
                  fill={PLANET_COLORS[currentProfection.lordOfYear]}
                  fontSize="28"
                  textAnchor="middle"
                >
                  {PLANET_SYMBOLS[currentProfection.lordOfYear]}
                </text>
                <text
                  x={center}
                  y={center + 35}
                  fill="rgba(255,255,255,0.6)"
                  fontSize="12"
                  textAnchor="middle"
                >
                  第{currentProfection.house}宫
                </text>
              </g>
            )}

            {/* 年龄标记 */}
            {renderAgeMarkers}
          </svg>
        </div>

        {/* 当前年限法详情 */}
        {currentProfection && (
          <motion.div
            className="flex-1 space-y-4"
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.3 }}
          >
            {/* 激活宫位 */}
            <div className="bg-white/5 rounded-lg p-4">
              <div className="text-sm text-white/60 mb-1">激活宫位</div>
              <div className="text-xl font-medium text-white">
                第{currentProfection.house}宫 · {currentProfection.houseName}
              </div>
              <div className="text-sm text-white/50 mt-1">
                {currentProfection.houseTheme}
              </div>
            </div>

            {/* 年主星 */}
            <div className="bg-white/5 rounded-lg p-4">
              <div className="text-sm text-white/60 mb-1">年主星</div>
              <div 
                className="text-xl font-medium flex items-center gap-2"
                style={{ color: PLANET_COLORS[currentProfection.lordOfYear] }}
              >
                <span className="text-2xl">{currentProfection.lordSymbol}</span>
                <span>{currentProfection.lordName}</span>
              </div>
            </div>

            {/* 年度关键词 */}
            {currentProfection.houseKeywords && currentProfection.houseKeywords.length > 0 && (
              <div className="bg-white/5 rounded-lg p-4">
                <div className="text-sm text-white/60 mb-2">关键词</div>
                <div className="flex flex-wrap gap-2">
                  {currentProfection.houseKeywords.map((keyword, i) => (
                    <span
                      key={i}
                      className="px-3 py-1 bg-purple-500/20 text-purple-300 rounded-full text-sm"
                    >
                      {keyword}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </motion.div>
        )}
      </div>
    </div>
  );
}

export default ProfectionWheel6183RS;

