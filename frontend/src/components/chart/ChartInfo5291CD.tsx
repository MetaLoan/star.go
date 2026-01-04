/**
 * ChartInfo5291CD - 星盘详情信息面板
 * 显示行星列表、相位、元素/模式平衡等详细信息
 */

import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { NatalChart, PlanetPosition, ZodiacID, AspectData } from '../../types';

interface ChartInfoProps {
  chart: NatalChart;
  className?: string;
}

// 元素颜色
const elementColors: Record<string, string> = {
  fire: '#FF6B6B',
  earth: '#4ECDC4',
  air: '#FFE66D',
  water: '#4FC3F7',
};

// 元素对应星座
const elementSigns: Record<string, ZodiacID[]> = {
  fire: ['aries', 'leo', 'sagittarius'],
  earth: ['taurus', 'virgo', 'capricorn'],
  air: ['gemini', 'libra', 'aquarius'],
  water: ['cancer', 'scorpio', 'pisces'],
};

// 模式对应星座
const modalitySigns: Record<string, ZodiacID[]> = {
  cardinal: ['aries', 'cancer', 'libra', 'capricorn'],
  fixed: ['taurus', 'leo', 'scorpio', 'aquarius'],
  mutable: ['gemini', 'virgo', 'sagittarius', 'pisces'],
};

type TabType = 'planets' | 'aspects' | 'elements' | 'houses';

export function ChartInfo5291CD({ chart, className = '' }: ChartInfoProps) {
  const [activeTab, setActiveTab] = useState<TabType>('planets');

  // 计算元素分布
  const elementDistribution = calculateElementDistribution(chart.planets);
  const modalityDistribution = calculateModalityDistribution(chart.planets);

  const tabs: { id: TabType; label: string }[] = [
    { id: 'planets', label: '行星' },
    { id: 'aspects', label: '相位' },
    { id: 'elements', label: '元素' },
    { id: 'houses', label: '宫位' },
  ];

  return (
    <div className={`glass-card p-4 ${className}`}>
      {/* Tab 导航 */}
      <div className="flex gap-2 mb-4 border-b border-white/10 pb-2">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
              activeTab === tab.id
                ? 'bg-cosmic-nova/20 text-cosmic-nova'
                : 'text-celestial-silver/60 hover:text-celestial-silver'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Tab 内容 */}
      <AnimatePresence mode="wait">
        <motion.div
          key={activeTab}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -10 }}
          transition={{ duration: 0.2 }}
          className="min-h-[300px]"
        >
          {activeTab === 'planets' && <PlanetList planets={chart.planets} />}
          {activeTab === 'aspects' && <AspectList aspects={chart.aspects} />}
          {activeTab === 'elements' && (
            <ElementsPanel
              elements={elementDistribution}
              modalities={modalityDistribution}
            />
          )}
          {activeTab === 'houses' && <HousesList houses={chart.houses} />}
        </motion.div>
      </AnimatePresence>
    </div>
  );
}

// 行星列表
function PlanetList({ planets }: { planets: PlanetPosition[] }) {
  return (
    <div className="space-y-2">
      {planets.map((planet) => (
        <div
          key={planet.id}
          className="flex items-center justify-between p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors"
        >
          <div className="flex items-center gap-3">
            <span className="text-xl">{planet.symbol}</span>
            <div>
              <div className="font-medium">{planet.name}</div>
              <div className="text-sm text-celestial-silver/60">
                {planet.signSymbol} {planet.signName} {planet.signDegree.toFixed(1)}°
              </div>
            </div>
          </div>
          <div className="text-right">
            <div className="text-sm text-celestial-silver/60">
              第 {planet.house} 宫
            </div>
            {planet.retrograde && (
              <div className="text-xs text-cosmic-aurora">℞ 逆行</div>
            )}
            {planet.dignityScore !== 0 && (
              <div
                className={`text-xs ${
                  planet.dignityScore > 0 ? 'text-green-400' : 'text-red-400'
                }`}
              >
                {planet.dignityScore > 0 ? '+' : ''}{planet.dignityScore}
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}

// 相位列表
function AspectList({ aspects }: { aspects: AspectData[] }) {
  const aspectColors: Record<string, string> = {
    conjunction: '#FFD700',
    sextile: '#4FC3F7',
    square: '#FF6B6B',
    trine: '#4ECDC4',
    opposition: '#FF6B6B',
  };

  return (
    <div className="space-y-2 max-h-[350px] overflow-y-auto">
      {aspects
        .filter((a) => a.strength > 0.5)
        .sort((a, b) => b.strength - a.strength)
        .map((aspect, index) => (
          <div
            key={index}
            className="flex items-center justify-between p-2 rounded-lg bg-white/5"
          >
            <div className="flex items-center gap-2">
              <span
                className="w-3 h-3 rounded-full"
                style={{ backgroundColor: aspectColors[aspect.aspectType] || '#888' }}
              />
              <span>{aspect.planet1Symbol}</span>
              <span className="text-celestial-silver/60">{aspect.aspectSymbol}</span>
              <span>{aspect.planet2Symbol}</span>
            </div>
            <div className="text-right">
              <div className="text-sm">{aspect.aspectName}</div>
              <div className="text-xs text-celestial-silver/60">
                {aspect.orb.toFixed(1)}° {aspect.applying ? '入相' : '离相'}
              </div>
            </div>
          </div>
        ))}
    </div>
  );
}

// 元素和模式面板
function ElementsPanel({
  elements,
  modalities,
}: {
  elements: Record<string, number>;
  modalities: Record<string, number>;
}) {
  const elementNames: Record<string, string> = {
    fire: '火象',
    earth: '土象',
    air: '风象',
    water: '水象',
  };

  const modalityNames: Record<string, string> = {
    cardinal: '基本',
    fixed: '固定',
    mutable: '变动',
  };

  const total = Object.values(elements).reduce((a, b) => a + b, 0);

  return (
    <div className="space-y-6">
      {/* 元素分布 */}
      <div>
        <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
          元素分布
        </h4>
        <div className="grid grid-cols-2 gap-3">
          {Object.entries(elements).map(([element, count]) => (
            <div
              key={element}
              className="p-3 rounded-lg bg-white/5"
              style={{ borderLeft: `3px solid ${elementColors[element]}` }}
            >
              <div className="flex justify-between items-center">
                <span className="font-medium">{elementNames[element]}</span>
                <span className="text-lg font-bold">{count}</span>
              </div>
              <div className="mt-2 h-1.5 bg-white/10 rounded-full overflow-hidden">
                <div
                  className="h-full rounded-full transition-all"
                  style={{
                    width: `${(count / total) * 100}%`,
                    backgroundColor: elementColors[element],
                  }}
                />
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* 模式分布 */}
      <div>
        <h4 className="text-sm font-medium text-celestial-silver/60 mb-3">
          模式分布
        </h4>
        <div className="grid grid-cols-3 gap-3">
          {Object.entries(modalities).map(([modality, count]) => (
            <div key={modality} className="p-3 rounded-lg bg-white/5 text-center">
              <div className="text-2xl font-bold text-cosmic-nova">{count}</div>
              <div className="text-sm text-celestial-silver/60">
                {modalityNames[modality]}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

// 宫位列表
function HousesList({ houses }: { houses: NatalChart['houses'] }) {
  const houseThemes: Record<number, string> = {
    1: '自我身份',
    2: '资源价值',
    3: '沟通学习',
    4: '家庭根基',
    5: '创造表达',
    6: '服务健康',
    7: '关系合作',
    8: '转化共享',
    9: '探索智慧',
    10: '事业成就',
    11: '愿景社群',
    12: '内省超越',
  };

  return (
    <div className="grid grid-cols-2 gap-2 max-h-[350px] overflow-y-auto">
      {houses.map((house) => (
        <div
          key={house.house}
          className="p-3 rounded-lg bg-white/5 hover:bg-white/10 transition-colors"
        >
          <div className="flex justify-between items-start">
            <div className="text-lg font-bold text-cosmic-nova">
              {house.house}宫
            </div>
            <div className="text-sm text-celestial-silver/60">
              {house.cusp.toFixed(1)}°
            </div>
          </div>
          <div className="text-sm mt-1">{house.signName}</div>
          <div className="text-xs text-celestial-silver/60 mt-1">
            {houseThemes[house.house]}
          </div>
        </div>
      ))}
    </div>
  );
}

// 计算元素分布
function calculateElementDistribution(
  planets: PlanetPosition[]
): Record<string, number> {
  const distribution: Record<string, number> = {
    fire: 0,
    earth: 0,
    air: 0,
    water: 0,
  };

  planets.forEach((planet) => {
    for (const [element, signs] of Object.entries(elementSigns)) {
      if (signs.includes(planet.sign)) {
        distribution[element]++;
        break;
      }
    }
  });

  return distribution;
}

// 计算模式分布
function calculateModalityDistribution(
  planets: PlanetPosition[]
): Record<string, number> {
  const distribution: Record<string, number> = {
    cardinal: 0,
    fixed: 0,
    mutable: 0,
  };

  planets.forEach((planet) => {
    for (const [modality, signs] of Object.entries(modalitySigns)) {
      if (signs.includes(planet.sign)) {
        distribution[modality]++;
        break;
      }
    }
  });

  return distribution;
}

export default ChartInfo5291CD;

