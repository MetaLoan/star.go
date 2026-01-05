/**
 * 星盘 SVG 组件
 * 组件命名规范：NatalChartSVG + 3847 + AB
 * 自主开发的星盘可视化组件
 */

import { useMemo } from 'react';
import { motion } from 'framer-motion';
import type { NatalChart, PlanetPosition, PlanetID, AspectType } from '../../types';
import {
  PLANET_SYMBOLS,
  PLANET_COLORS,
  ZODIAC_SYMBOLS,
  ZODIAC_COLORS,
  ASPECT_COLORS,
} from '../../utils/astro';

interface NatalChartSVGProps {
  chart: NatalChart;
  size?: number;
  showAspects?: boolean;
  showHouses?: boolean;
  showGuideLines?: boolean;
  highlightPlanet?: PlanetID | null;
  onPlanetClick?: (planet: PlanetID) => void;
}

// 星座顺序（从白羊座开始）
const ZODIAC_ORDER: Array<keyof typeof ZODIAC_SYMBOLS> = [
  'aries', 'taurus', 'gemini', 'cancer', 'leo', 'virgo',
  'libra', 'scorpio', 'sagittarius', 'capricorn', 'aquarius', 'pisces',
];

// 极坐标转笛卡尔坐标
function polarToCartesian(
  centerX: number,
  centerY: number,
  radius: number,
  angleInDegrees: number
): { x: number; y: number } {
  const angleInRadians = (angleInDegrees * Math.PI) / 180;
  return {
    x: centerX + radius * Math.cos(angleInRadians),
    y: centerY - radius * Math.sin(angleInRadians),
  };
}

export function NatalChartSVG3847AB({
  chart,
  size = 500,
  showAspects = true,
  showHouses = true,
  showGuideLines = true,
  highlightPlanet = null,
  onPlanetClick,
}: NatalChartSVGProps) {
  // 增加 padding 用于显示外部标签
  const padding = size * 0.08;
  const viewBoxSize = size + padding * 2;
  const center = viewBoxSize / 2;
  const outerRadius = size * 0.44;
  const zodiacRingWidth = size * 0.07;
  const zodiacInnerRadius = outerRadius - zodiacRingWidth;
  const houseRingWidth = size * 0.042;
  const innerRadius = zodiacInnerRadius - houseRingWidth;
  const planetRadius = innerRadius * 0.72;
  const aspectRadius = innerRadius * 0.48;
  const guideRadius1 = innerRadius * 0.85;
  const guideRadius2 = innerRadius * 0.60;
  const guideRadius3 = innerRadius * 0.35;

  // 黄道经度转换为 SVG 显示角度
  const longitudeToAngle = (longitude: number): number => {
    const angle = 180 + chart.ascendant - longitude;
    return ((angle % 360) + 360) % 360;
  };

  // 渲染辅助线（同心圆和角度分割线）
  const renderGuideLines = useMemo(() => {
    if (!showGuideLines) return null;
    
    const lines = [];
    
    // 同心圆引导线 - 更明显
    const concentricCircles = [guideRadius1, guideRadius2, guideRadius3];
    const circleOpacities = [0.12, 0.10, 0.08];
    concentricCircles.forEach((r, i) => {
      lines.push(
        <circle
          key={`concentric-${i}`}
          cx={center}
          cy={center}
          r={r}
          fill="none"
          stroke={`rgba(150, 180, 255, ${circleOpacities[i]})`}
          strokeWidth="1"
          strokeDasharray="4 8"
        />
      );
    });
    
    // 30度分割线（细线）- 更明显
    for (let i = 0; i < 12; i++) {
      const angle = i * 30;
      const start = polarToCartesian(center, center, guideRadius3, angle);
      const end = polarToCartesian(center, center, innerRadius, angle);
      lines.push(
        <line
          key={`guide-30-${i}`}
          x1={start.x}
          y1={start.y}
          x2={end.x}
          y2={end.y}
          stroke="rgba(150, 180, 255, 0.08)"
          strokeWidth="1"
        />
      );
    }
    
    return <g className="guide-lines">{lines}</g>;
  }, [showGuideLines, center, innerRadius, guideRadius1, guideRadius2, guideRadius3]);

  // 渲染主轴线（ASC-DSC, MC-IC）
  const renderAxisLines = useMemo(() => {
    // ASC 在左侧 (180度)，DSC 在右侧 (0度)
    const ascAngle = 180;
    const dscAngle = 0;
    
    // MC 通常在顶部附近，IC 在底部
    // 根据 chart.midheaven 计算 MC 角度
    const mcAngle = longitudeToAngle(chart.midheaven);
    const icAngle = (mcAngle + 180) % 360;
    
    return (
      <g className="axis-lines">
        {/* ASC-DSC 轴线光晕 */}
        <line
          x1={polarToCartesian(center, center, innerRadius, ascAngle).x}
          y1={polarToCartesian(center, center, innerRadius, ascAngle).y}
          x2={polarToCartesian(center, center, innerRadius, dscAngle).x}
          y2={polarToCartesian(center, center, innerRadius, dscAngle).y}
          stroke="rgba(0, 212, 255, 0.15)"
          strokeWidth="6"
        />
        {/* ASC-DSC 轴线 */}
        <line
          x1={polarToCartesian(center, center, innerRadius, ascAngle).x}
          y1={polarToCartesian(center, center, innerRadius, ascAngle).y}
          x2={polarToCartesian(center, center, innerRadius, dscAngle).x}
          y2={polarToCartesian(center, center, innerRadius, dscAngle).y}
          stroke="rgba(0, 212, 255, 0.6)"
          strokeWidth="2"
          strokeDasharray="10 5"
        />
        
        {/* MC-IC 轴线光晕 */}
        <line
          x1={polarToCartesian(center, center, innerRadius, mcAngle).x}
          y1={polarToCartesian(center, center, innerRadius, mcAngle).y}
          x2={polarToCartesian(center, center, innerRadius, icAngle).x}
          y2={polarToCartesian(center, center, innerRadius, icAngle).y}
          stroke="rgba(255, 107, 157, 0.15)"
          strokeWidth="6"
        />
        {/* MC-IC 轴线 */}
        <line
          x1={polarToCartesian(center, center, innerRadius, mcAngle).x}
          y1={polarToCartesian(center, center, innerRadius, mcAngle).y}
          x2={polarToCartesian(center, center, innerRadius, icAngle).x}
          y2={polarToCartesian(center, center, innerRadius, icAngle).y}
          stroke="rgba(255, 107, 157, 0.6)"
          strokeWidth="2"
          strokeDasharray="10 5"
        />
        
        {/* ASC 箭头标记 */}
        <polygon
          points={`
            ${center - innerRadius - 6},${center}
            ${center - innerRadius + 10},${center - 6}
            ${center - innerRadius + 10},${center + 6}
          `}
          fill="rgba(0, 212, 255, 0.8)"
        />
        
        {/* MC 箭头标记 */}
        {(() => {
          const mcPos = polarToCartesian(center, center, innerRadius, mcAngle);
          const arrowAngle = mcAngle * Math.PI / 180;
          const arrowLen = 10;
          const arrowWidth = 6;
          // 计算箭头方向（指向外侧）
          const tipX = mcPos.x + Math.cos(arrowAngle) * 6;
          const tipY = mcPos.y - Math.sin(arrowAngle) * 6;
          const baseX = mcPos.x - Math.cos(arrowAngle) * arrowLen;
          const baseY = mcPos.y + Math.sin(arrowAngle) * arrowLen;
          const perpX = Math.sin(arrowAngle) * arrowWidth;
          const perpY = Math.cos(arrowAngle) * arrowWidth;
          
          return (
            <polygon
              points={`${tipX},${tipY} ${baseX + perpX},${baseY + perpY} ${baseX - perpX},${baseY - perpY}`}
              fill="rgba(255, 107, 157, 0.8)"
            />
          );
        })()}
      </g>
    );
  }, [chart.ascendant, chart.midheaven, center, innerRadius]);

  // 渲染元素三角形（火、土、风、水）
  const renderElementTriangles = useMemo(() => {
    if (!showGuideLines) return null;
    
    const elementGroups = {
      fire: [0, 4, 8],      // 白羊(0)、狮子(4)、射手(8)
      earth: [1, 5, 9],     // 金牛(1)、处女(5)、摩羯(9)
      air: [2, 6, 10],      // 双子(2)、天秤(6)、水瓶(10)
      water: [3, 7, 11],    // 巨蟹(3)、天蝎(7)、双鱼(11)
    };
    
    const elementColors = {
      fire: { fill: 'rgba(239, 68, 68, 0.18)', stroke: 'rgba(239, 68, 68, 0.4)' },
      earth: { fill: 'rgba(34, 197, 94, 0.15)', stroke: 'rgba(34, 197, 94, 0.35)' },
      air: { fill: 'rgba(234, 179, 8, 0.15)', stroke: 'rgba(234, 179, 8, 0.35)' },
      water: { fill: 'rgba(59, 130, 246, 0.18)', stroke: 'rgba(59, 130, 246, 0.4)' },
    };
    
    return (
      <g className="element-triangles">
        {Object.entries(elementGroups).map(([element, signs]) => {
          const points = signs.map(signIndex => {
            // 星座中心位置（每个星座15度）
            const longitude = signIndex * 30 + 15;
            const angle = longitudeToAngle(longitude);
            return polarToCartesian(center, center, guideRadius2, angle);
          });
          
          const colors = elementColors[element as keyof typeof elementColors];
          
          return (
            <polygon
              key={`triangle-${element}`}
              points={points.map(p => `${p.x},${p.y}`).join(' ')}
              fill={colors.fill}
              stroke={colors.stroke}
              strokeWidth="1.5"
            />
          );
        })}
      </g>
    );
  }, [showGuideLines, chart.ascendant, center, guideRadius2]);

  // 渲染星座环背景扇形
  const renderZodiacBackground = useMemo(() => {
    return ZODIAC_ORDER.map((sign, index) => {
      const signStartLongitude = index * 30;
      const startAngle = longitudeToAngle(signStartLongitude);
      const endAngle = longitudeToAngle(signStartLongitude + 30);
      
      const outerStart = polarToCartesian(center, center, outerRadius - 2, startAngle);
      const outerEnd = polarToCartesian(center, center, outerRadius - 2, endAngle);
      const innerStart = polarToCartesian(center, center, zodiacInnerRadius + 1, startAngle);
      const innerEnd = polarToCartesian(center, center, zodiacInnerRadius + 1, endAngle);
      
      let angleDiff = startAngle - endAngle;
      if (angleDiff < 0) angleDiff += 360;
      const largeArc = angleDiff > 180 ? 1 : 0;
      
      const path = `
        M ${outerStart.x} ${outerStart.y}
        A ${outerRadius - 2} ${outerRadius - 2} 0 ${largeArc} 0 ${outerEnd.x} ${outerEnd.y}
        L ${innerEnd.x} ${innerEnd.y}
        A ${zodiacInnerRadius + 1} ${zodiacInnerRadius + 1} 0 ${largeArc} 1 ${innerStart.x} ${innerStart.y}
        Z
      `;
      
      return (
        <path
          key={`zodiac-bg-${sign}`}
          d={path}
          fill={`${ZODIAC_COLORS[sign]}30`}
          stroke={`${ZODIAC_COLORS[sign]}20`}
          strokeWidth="0.5"
        />
      );
    });
  }, [chart.ascendant, center, outerRadius, zodiacInnerRadius]);

  // 渲染星座符号
  const renderZodiacSymbols = useMemo(() => {
    return ZODIAC_ORDER.map((sign, index) => {
      const signStartLongitude = index * 30;
      const midLongitude = signStartLongitude + 15;
      const midAngle = longitudeToAngle(midLongitude);
      const startAngle = longitudeToAngle(signStartLongitude);
      
      const symbolPos = polarToCartesian(center, center, outerRadius - zodiacRingWidth / 2, midAngle);
      const lineStart = polarToCartesian(center, center, zodiacInnerRadius, startAngle);
      const lineEnd = polarToCartesian(center, center, outerRadius, startAngle);
      const signColor = ZODIAC_COLORS[sign];
      
      return (
        <g key={sign}>
          {/* 星座分隔线 */}
          <line
            x1={lineStart.x}
            y1={lineStart.y}
            x2={lineEnd.x}
            y2={lineEnd.y}
            stroke="rgba(200, 220, 255, 0.4)"
            strokeWidth="1"
          />
          {/* 星座符号背景光晕 */}
          <circle
            cx={symbolPos.x}
            cy={symbolPos.y}
            r={size * 0.035}
            fill={signColor}
            opacity={0.2}
          />
          {/* 星座符号背景 */}
          <circle
            cx={symbolPos.x}
            cy={symbolPos.y}
            r={size * 0.026}
            fill={signColor}
            opacity={0.95}
          />
          {/* 星座符号 */}
          <text
            x={symbolPos.x}
            y={symbolPos.y + 1}
            fill="white"
            fontSize={size * 0.026}
            textAnchor="middle"
            dominantBaseline="central"
            className="select-none"
            fontWeight="bold"
            style={{
              textShadow: '0 1px 2px rgba(0,0,0,0.5)',
            }}
          >
            {ZODIAC_SYMBOLS[sign]}
          </text>
        </g>
      );
    });
  }, [chart.ascendant, center, outerRadius, zodiacRingWidth, zodiacInnerRadius, size]);

  // 渲染宫位
  const renderHouses = useMemo(() => {
    if (!showHouses) return null;
    
    return chart.houses.map((house, index) => {
      const angle = longitudeToAngle(house.cusp);
      const nextHouse = chart.houses[(index + 1) % 12];
      const nextAngle = longitudeToAngle(nextHouse.cusp);
      
      const lineStart = polarToCartesian(center, center, innerRadius, angle);
      const lineEnd = polarToCartesian(center, center, zodiacInnerRadius, angle);
      
      let midAngle: number;
      let angleDiff = ((angle - nextAngle + 360) % 360);
      midAngle = (angle - angleDiff / 2 + 360) % 360;
      
      const numPos = polarToCartesian(center, center, innerRadius + houseRingWidth / 2, midAngle);
      const isAngular = [1, 4, 7, 10].includes(house.house);
      
      return (
        <g key={`house-${house.house}`}>
          {/* 宫位分隔线 */}
          <line
            x1={lineStart.x}
            y1={lineStart.y}
            x2={lineEnd.x}
            y2={lineEnd.y}
            stroke={isAngular ? 'rgba(200, 220, 255, 0.6)' : 'rgba(150, 180, 255, 0.25)'}
            strokeWidth={isAngular ? 1.5 : 0.8}
          />
          {/* 宫位数字 */}
          <text
            x={numPos.x}
            y={numPos.y}
            fill={isAngular ? 'rgba(200, 220, 255, 0.85)' : 'rgba(180, 200, 255, 0.6)'}
            fontSize={isAngular ? size * 0.024 : size * 0.020}
            textAnchor="middle"
            dominantBaseline="central"
            className="select-none"
            fontWeight={isAngular ? '600' : '400'}
          >
            {house.house}
          </text>
        </g>
      );
    });
  }, [chart.houses, chart.ascendant, showHouses, center, innerRadius, zodiacInnerRadius, houseRingWidth, size]);

  // 处理行星重叠
  const processedPlanets = useMemo(() => {
    const result: Array<{ planet: PlanetPosition; angle: number; displayAngle: number; radius: number }> = 
      chart.planets.map(planet => ({
        planet,
        angle: longitudeToAngle(planet.longitude),
        displayAngle: longitudeToAngle(planet.longitude),
        radius: planetRadius,
      }));
    
    result.sort((a, b) => a.displayAngle - b.displayAngle);
    
    const minGap = 15;
    
    for (let iteration = 0; iteration < 8; iteration++) {
      result.sort((a, b) => a.displayAngle - b.displayAngle);
      
      for (let i = 0; i < result.length; i++) {
        const current = result[i];
        const next = result[(i + 1) % result.length];
        
        let diff = ((next.displayAngle - current.displayAngle + 360) % 360);
        if (diff === 0) diff = 360;
        
        if (diff < minGap) {
          const adjustment = (minGap - diff) / 2 + 0.5;
          current.displayAngle = ((current.displayAngle - adjustment + 360) % 360);
          next.displayAngle = ((next.displayAngle + adjustment) % 360);
        }
      }
    }
    
    return result;
  }, [chart.planets, chart.ascendant, planetRadius]);

  // 创建行星角度映射
  const planetAngleMap = useMemo(() => {
    const map = new Map<PlanetID, { displayAngle: number; radius: number }>();
    processedPlanets.forEach(({ planet, displayAngle, radius }) => {
      map.set(planet.id, { displayAngle, radius });
    });
    return map;
  }, [processedPlanets]);

  // 渲染行星连接线（从行星位置到实际经度位置的引导线）
  const renderPlanetGuideLines = useMemo(() => {
    return processedPlanets.map(({ planet, angle, displayAngle }) => {
      // 如果显示角度和实际角度相差较大，画一条引导线
      const angleDiff = Math.abs(displayAngle - angle);
      if (angleDiff > 5 && angleDiff < 355) {
        const planetPos = polarToCartesian(center, center, planetRadius - size * 0.02, displayAngle);
        const actualPos = polarToCartesian(center, center, innerRadius - size * 0.01, angle);
        
        return (
          <line
            key={`guide-${planet.id}`}
            x1={planetPos.x}
            y1={planetPos.y}
            x2={actualPos.x}
            y2={actualPos.y}
            stroke={PLANET_COLORS[planet.id]}
            strokeWidth="1"
            strokeOpacity="0.3"
            strokeDasharray="2 2"
          />
        );
      }
      return null;
    });
  }, [processedPlanets, center, planetRadius, innerRadius, size]);

  // 渲染相位线
  const renderAspects = useMemo(() => {
    if (!showAspects) return null;
    
    return chart.aspects.map((aspect, index) => {
      const p1Info = planetAngleMap.get(aspect.planet1);
      const p2Info = planetAngleMap.get(aspect.planet2);
      
      if (!p1Info || !p2Info) return null;
      
      const pos1 = polarToCartesian(center, center, aspectRadius, p1Info.displayAngle);
      const pos2 = polarToCartesian(center, center, aspectRadius, p2Info.displayAngle);
      
      const isHighlighted = 
        highlightPlanet === aspect.planet1 || highlightPlanet === aspect.planet2;
      
      const aspectColor = ASPECT_COLORS[aspect.aspectType as AspectType] || '#888';
      
      // 根据相位类型设置不同的线条样式
      const isTenseAspect = aspect.aspectType === 'square' || aspect.aspectType === 'opposition';
      const isHarmoniousAspect = aspect.aspectType === 'trine' || aspect.aspectType === 'sextile';
      
      return (
        <motion.line
          key={`aspect-${index}`}
          x1={pos1.x}
          y1={pos1.y}
          x2={pos2.x}
          y2={pos2.y}
          stroke={aspectColor}
          strokeWidth={isHighlighted ? 3 : 2}
          strokeOpacity={isHighlighted ? 1 : 0.75}
          strokeDasharray={isTenseAspect ? '6 3' : isHarmoniousAspect ? 'none' : '3 3'}
          initial={{ pathLength: 0, opacity: 0 }}
          animate={{ pathLength: 1, opacity: isHighlighted ? 1 : 0.75 }}
          transition={{ duration: 0.4, delay: index * 0.015 }}
          style={{
            filter: isHighlighted ? `drop-shadow(0 0 4px ${aspectColor})` : 'none',
          }}
        />
      );
    });
  }, [chart.aspects, planetAngleMap, showAspects, center, aspectRadius, highlightPlanet]);

  // 格式化行星度数（只显示在星座内的度数 0-29°）
  const formatDegree = (longitude: number): string => {
    const degreeInSign = Math.floor(longitude % 30);
    return `${degreeInSign}°`;
  };

  // 渲染行星
  const renderPlanets = useMemo(() => {
    return processedPlanets.map(({ planet, displayAngle }, index) => {
      const pos = polarToCartesian(center, center, planetRadius, displayAngle);
      const isHighlighted = highlightPlanet === planet.id;
      const retrograde = planet.retrograde;
      const color = PLANET_COLORS[planet.id] || '#fff';
      const degree = formatDegree(planet.longitude);
      
      // 计算度数标签位置（在行星符号右上方）
      const degreeOffsetX = size * 0.026;
      const degreeOffsetY = -size * 0.018;
      
      return (
        <motion.g
          key={planet.id}
          initial={{ scale: 0, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.3, delay: index * 0.04 }}
          style={{ cursor: onPlanetClick ? 'pointer' : 'default' }}
          onClick={() => onPlanetClick?.(planet.id)}
        >
          {/* 行星光晕背景 */}
          <circle
            cx={pos.x}
            cy={pos.y}
            r={size * 0.028}
            fill={color}
            opacity={isHighlighted ? 0.35 : 0.15}
          />
          
          {isHighlighted && (
            <circle
              cx={pos.x}
              cy={pos.y}
              r={size * 0.04}
              fill={color}
              opacity={0.2}
            />
          )}
          
          {/* 行星符号 */}
          <text
            x={pos.x}
            y={pos.y + 1}
            fill={color}
            fontSize={isHighlighted ? size * 0.044 : size * 0.036}
            textAnchor="middle"
            dominantBaseline="central"
            className="select-none"
            style={{
              filter: isHighlighted ? `drop-shadow(0 0 8px ${color})` : `drop-shadow(0 0 4px ${color})`,
              fontWeight: 'bold',
            }}
          >
            {PLANET_SYMBOLS[planet.id]}
          </text>
          
          {/* 度数标注（上标形式） */}
          <text
            x={pos.x + degreeOffsetX}
            y={pos.y + degreeOffsetY}
            fill="rgba(220, 230, 255, 0.85)"
            fontSize={size * 0.016}
            textAnchor="start"
            dominantBaseline="central"
            className="select-none"
            style={{
              fontFamily: 'system-ui, sans-serif',
              fontWeight: '500',
            }}
          >
            {degree}
          </text>
          
          {/* 逆行标记 */}
          {retrograde && (
            <text
              x={pos.x + size * 0.024}
              y={pos.y + size * 0.016}
              fill="#ff6b6b"
              fontSize={size * 0.013}
              textAnchor="middle"
              dominantBaseline="central"
              className="select-none"
              fontWeight="bold"
              style={{
                fontStyle: 'italic',
              }}
            >
              ℞
            </text>
          )}
        </motion.g>
      );
    });
  }, [processedPlanets, center, planetRadius, size, highlightPlanet, onPlanetClick]);

  return (
    <svg
      width={viewBoxSize}
      height={viewBoxSize}
      viewBox={`0 0 ${viewBoxSize} ${viewBoxSize}`}
      className="drop-shadow-xl"
    >
      <defs>
        {/* 深空背景渐变 - 更亮更有层次 */}
        <radialGradient id="chartBg" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#252547" />
          <stop offset="40%" stopColor="#1a1a35" />
          <stop offset="70%" stopColor="#141428" />
          <stop offset="100%" stopColor="#0f0f1e" />
        </radialGradient>
        
        {/* 中心区域渐变 - 更亮 */}
        <radialGradient id="centerBg" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#2a2a4a" />
          <stop offset="50%" stopColor="#1e1e3a" />
          <stop offset="100%" stopColor="#16162d" />
        </radialGradient>
        
        {/* 星空点缀渐变 */}
        <radialGradient id="starGlow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#ffffff" stopOpacity="0.8" />
          <stop offset="100%" stopColor="#ffffff" stopOpacity="0" />
        </radialGradient>
        
        {/* 外圈光晕渐变 */}
        <radialGradient id="outerGlow" cx="50%" cy="50%" r="50%">
          <stop offset="85%" stopColor="transparent" />
          <stop offset="95%" stopColor="rgba(0, 212, 255, 0.15)" />
          <stop offset="100%" stopColor="rgba(0, 212, 255, 0.05)" />
        </radialGradient>
        
        {/* 内圈光晕 */}
        <filter id="glowFilter" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="3" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>
      
      {/* 最外层光晕 - 多层次 */}
      <circle
        cx={center}
        cy={center}
        r={outerRadius + 15}
        fill="none"
        stroke="rgba(0, 180, 220, 0.08)"
        strokeWidth="20"
      />
      <circle
        cx={center}
        cy={center}
        r={outerRadius + 6}
        fill="none"
        stroke="rgba(100, 200, 255, 0.2)"
        strokeWidth="8"
      />
      
      {/* 星座环背景 */}
      <circle
        cx={center}
        cy={center}
        r={outerRadius}
        fill="url(#chartBg)"
      />
      
      {/* 星座扇形背景 */}
      {renderZodiacBackground}
      
      {/* 星座环外圈 - 双线效果 */}
      <circle
        cx={center}
        cy={center}
        r={outerRadius}
        fill="none"
        stroke="rgba(100, 200, 255, 0.7)"
        strokeWidth="2.5"
      />
      <circle
        cx={center}
        cy={center}
        r={outerRadius - 1}
        fill="none"
        stroke="rgba(0, 212, 255, 0.3)"
        strokeWidth="1"
      />
      
      {/* 星座环内圈 */}
      <circle
        cx={center}
        cy={center}
        r={zodiacInnerRadius}
        fill="none"
        stroke="rgba(150, 180, 255, 0.4)"
        strokeWidth="1.5"
      />
      
      {/* 内圈背景 */}
      <circle
        cx={center}
        cy={center}
        r={innerRadius}
        fill="url(#centerBg)"
        stroke="rgba(150, 180, 255, 0.25)"
        strokeWidth="1.5"
      />
      
      {/* 辅助线 */}
      {renderGuideLines}
      
      {/* 元素三角形 */}
      {renderElementTriangles}
      
      {/* 主轴线 */}
      {renderAxisLines}
      
      {/* 中心装饰 - 更精致 */}
      <circle
        cx={center}
        cy={center}
        r={aspectRadius * 0.18}
        fill="none"
        stroke="rgba(150, 180, 255, 0.15)"
        strokeWidth="1"
      />
      <circle
        cx={center}
        cy={center}
        r={aspectRadius * 0.10}
        fill="none"
        stroke="rgba(200, 220, 255, 0.2)"
        strokeWidth="1"
      />
      <circle
        cx={center}
        cy={center}
        r={5}
        fill="rgba(200, 220, 255, 0.6)"
      />
      <circle
        cx={center}
        cy={center}
        r={2}
        fill="rgba(255, 255, 255, 0.9)"
      />
      
      {/* 星座符号 */}
      {renderZodiacSymbols}
      
      {/* 宫位 */}
      {renderHouses}
      
      {/* 行星引导线 */}
      {renderPlanetGuideLines}
      
      {/* 相位线 */}
      {renderAspects}
      
      {/* 行星 */}
      {renderPlanets}
      
      {/* 轴点标记 */}
      <g className="select-none">
        {/* ASC 标签 */}
        <text
          x={center - outerRadius - 22}
          y={center}
          fill="#00D4FF"
          fontSize={size * 0.030}
          textAnchor="end"
          dominantBaseline="central"
          fontWeight="600"
          style={{ 
            filter: 'drop-shadow(0 0 6px rgba(0, 212, 255, 0.6))',
            letterSpacing: '0.5px',
          }}
        >
          ASC
        </text>
        
        {/* DSC 标签 */}
        <text
          x={center + outerRadius + 22}
          y={center}
          fill="#00D4FF"
          fontSize={size * 0.030}
          textAnchor="start"
          dominantBaseline="central"
          fontWeight="600"
          style={{ 
            filter: 'drop-shadow(0 0 6px rgba(0, 212, 255, 0.6))',
            letterSpacing: '0.5px',
          }}
        >
          DSC
        </text>
        
        {/* MC 标记（根据实际位置） */}
        {(() => {
          const mcAngle = longitudeToAngle(chart.midheaven);
          const mcPos = polarToCartesian(center, center, outerRadius + 24, mcAngle);
          return (
            <text
              x={mcPos.x}
              y={mcPos.y}
              fill="#FF6B9D"
              fontSize={size * 0.030}
              textAnchor="middle"
              dominantBaseline="central"
              fontWeight="600"
              style={{ 
                filter: 'drop-shadow(0 0 6px rgba(255, 107, 157, 0.6))',
                letterSpacing: '0.5px',
              }}
            >
              MC
            </text>
          );
        })()}
        
        {/* IC 标记 */}
        {(() => {
          const icAngle = (longitudeToAngle(chart.midheaven) + 180) % 360;
          const icPos = polarToCartesian(center, center, outerRadius + 24, icAngle);
          return (
            <text
              x={icPos.x}
              y={icPos.y}
              fill="#FF6B9D"
              fontSize={size * 0.030}
              textAnchor="middle"
              dominantBaseline="central"
              fontWeight="600"
              style={{ 
                filter: 'drop-shadow(0 0 6px rgba(255, 107, 157, 0.6))',
                letterSpacing: '0.5px',
              }}
            >
              IC
            </text>
          );
        })()}
      </g>
    </svg>
  );
}

export default NatalChartSVG3847AB;
