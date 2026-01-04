/**
 * 影响因子面板组件
 * 组件命名规范：InfluenceFactorsPanel + 8274 + TU
 * 支持运营调整因子权重
 */

import { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Button, Slider, Tooltip } from '@heroui/react';
import type { InfluenceFactor, InfluenceFactorGroup } from '../../types';

interface InfluenceFactorsPanelProps {
  factors: InfluenceFactor[];
  onWeightChange?: (name: string, weight: number) => void;
  editable?: boolean;
  className?: string;
}

// 因子类别配置
const CATEGORY_CONFIG: Record<string, { icon: string; color: string; label: string }> = {
  dignity: { icon: '👑', color: '#ffd700', label: '尊贵度' },
  retrograde: { icon: '℞', color: '#ef4444', label: '逆行' },
  aspectPhase: { icon: '📐', color: '#3b82f6', label: '相位阶段' },
  aspectOrb: { icon: '🎯', color: '#22c55e', label: '容许度' },
  outerPlanet: { icon: '🪐', color: '#a855f7', label: '外行星' },
  profectionLord: { icon: '🔮', color: '#ec4899', label: '年主星' },
  lunarPhase: { icon: '🌙', color: '#c0c0c0', label: '月相' },
  planetaryHour: { icon: '⏰', color: '#f59e0b', label: '行星时' },
  personal: { icon: '👤', color: '#06b6d4', label: '个人因子' },
  custom: { icon: '⚙️', color: '#71717a', label: '自定义' },
};

// 将因子按类型分组
function groupFactors(factors: InfluenceFactor[]): InfluenceFactorGroup[] {
  const groups: Record<string, InfluenceFactor[]> = {};
  
  factors.forEach(factor => {
    const category = factor.type || 'other';
    if (!groups[category]) {
      groups[category] = [];
    }
    groups[category].push(factor);
  });

  return Object.entries(groups).map(([category, factors]) => ({
    category,
    factors,
    totalWeight: factors.reduce((sum, f) => sum + f.weight * f.value, 0),
  }));
}

export function InfluenceFactorsPanel8274TU({
  factors,
  onWeightChange,
  editable = false,
  className = '',
}: InfluenceFactorsPanelProps) {
  const [expandedCategory, setExpandedCategory] = useState<string | null>(null);
  const [localWeights, setLocalWeights] = useState<Record<string, number>>({});

  const groupedFactors = groupFactors(factors);

  // 计算总影响值
  const totalInfluence = factors.reduce((sum, f) => sum + f.weight * f.value, 0);
  const positiveInfluence = factors.filter(f => f.isPositive).reduce((sum, f) => sum + f.weight * f.value, 0);
  const negativeInfluence = factors.filter(f => !f.isPositive).reduce((sum, f) => sum + Math.abs(f.weight * f.value), 0);

  const handleWeightChange = useCallback((name: string, weight: number) => {
    setLocalWeights(prev => ({ ...prev, [name]: weight }));
    onWeightChange?.(name, weight);
  }, [onWeightChange]);

  const getFactorWeight = (factor: InfluenceFactor) => {
    return localWeights[factor.name] ?? factor.weight;
  };

  return (
    <div className={`glass-card p-4 ${className}`}>
      <h3 className="text-lg font-medium text-white mb-4 flex items-center justify-between">
        <span className="flex items-center gap-2">
          <span>📊</span>
          影响因子分析
        </span>
        {editable && (
          <span className="text-xs text-white/40 bg-white/5 px-2 py-1 rounded">
            可编辑模式
          </span>
        )}
      </h3>

      {/* 总体概览 */}
      <div className="grid grid-cols-3 gap-3 mb-6">
        <div className="bg-white/5 rounded-lg p-3 text-center">
          <div className="text-2xl font-bold text-white">
            {totalInfluence > 0 ? '+' : ''}{totalInfluence.toFixed(1)}
          </div>
          <div className="text-xs text-white/60">净影响</div>
        </div>
        <div className="bg-green-500/10 rounded-lg p-3 text-center">
          <div className="text-2xl font-bold text-green-400">
            +{positiveInfluence.toFixed(1)}
          </div>
          <div className="text-xs text-white/60">积极因素</div>
        </div>
        <div className="bg-red-500/10 rounded-lg p-3 text-center">
          <div className="text-2xl font-bold text-red-400">
            -{negativeInfluence.toFixed(1)}
          </div>
          <div className="text-xs text-white/60">挑战因素</div>
        </div>
      </div>

      {/* 因子类别列表 */}
      <div className="space-y-2">
        {groupedFactors.map((group) => {
          const config = CATEGORY_CONFIG[group.category] || CATEGORY_CONFIG.custom;
          const isExpanded = expandedCategory === group.category;

          return (
            <div key={group.category} className="bg-white/5 rounded-lg overflow-hidden">
              {/* 类别头部 */}
              <motion.button
                className="w-full flex items-center justify-between p-3 hover:bg-white/5 transition-colors"
                onClick={() => setExpandedCategory(isExpanded ? null : group.category)}
              >
                <div className="flex items-center gap-3">
                  <span className="text-lg" style={{ color: config.color }}>
                    {config.icon}
                  </span>
                  <span className="font-medium text-white">{config.label}</span>
                  <span className="text-sm text-white/40">
                    ({group.factors.length})
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <span
                    className={`text-sm font-medium ${
                      group.totalWeight >= 0 ? 'text-green-400' : 'text-red-400'
                    }`}
                  >
                    {group.totalWeight >= 0 ? '+' : ''}{group.totalWeight.toFixed(2)}
                  </span>
                  <motion.span
                    className="text-white/40"
                    animate={{ rotate: isExpanded ? 180 : 0 }}
                  >
                    ▼
                  </motion.span>
                </div>
              </motion.button>

              {/* 展开的因子列表 */}
              <AnimatePresence>
                {isExpanded && (
                  <motion.div
                    initial={{ height: 0, opacity: 0 }}
                    animate={{ height: 'auto', opacity: 1 }}
                    exit={{ height: 0, opacity: 0 }}
                    className="border-t border-white/5"
                  >
                    <div className="p-3 space-y-3">
                      {group.factors.map((factor, index) => (
                        <motion.div
                          key={factor.name}
                          className="flex items-start gap-3 text-sm"
                          initial={{ opacity: 0, x: -10 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.05 }}
                        >
                          {/* 正负指示 */}
                          <span
                            className={`w-5 h-5 rounded-full flex items-center justify-center text-xs ${
                              factor.isPositive
                                ? 'bg-green-500/20 text-green-400'
                                : 'bg-red-500/20 text-red-400'
                            }`}
                          >
                            {factor.isPositive ? '+' : '-'}
                          </span>

                          {/* 因子信息 */}
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <span className="text-white">{factor.name}</span>
                              {factor.dimension && (
                                <Tooltip content={`影响维度: ${factor.dimension}`}>
                                  <span className="text-white/30 cursor-help">ⓘ</span>
                                </Tooltip>
                              )}
                            </div>
                            <div className="text-white/50 text-xs mt-0.5">
                              {factor.description}
                            </div>

                            {/* 可编辑权重 */}
                            {editable && (
                              <div className="mt-2 flex items-center gap-3">
                                <span className="text-xs text-white/40">权重:</span>
                                <Slider
                                  size="sm"
                                  step={0.1}
                                  minValue={0}
                                  maxValue={2}
                                  value={getFactorWeight(factor)}
                                  onChange={(value) => 
                                    handleWeightChange(factor.name, value as number)
                                  }
                                  className="flex-1 max-w-32"
                                />
                                <span className="text-xs text-white/60 w-8">
                                  ×{getFactorWeight(factor).toFixed(1)}
                                </span>
                              </div>
                            )}
                          </div>

                          {/* 数值显示 */}
                          <div className="text-right">
                            <div
                              className={`font-medium ${
                                factor.isPositive ? 'text-green-400' : 'text-red-400'
                              }`}
                            >
                              {factor.value >= 0 ? '+' : ''}{factor.value.toFixed(2)}
                            </div>
                            <div className="text-white/40 text-xs">
                              ×{getFactorWeight(factor).toFixed(1)}
                            </div>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          );
        })}
      </div>

      {/* 编辑模式操作按钮 */}
      {editable && Object.keys(localWeights).length > 0 && (
        <div className="mt-4 flex justify-end gap-2">
          <Button
            size="sm"
            variant="flat"
            onPress={() => setLocalWeights({})}
          >
            重置
          </Button>
          <Button
            size="sm"
            className="bg-cosmic-nova text-white"
            onPress={() => {
              // TODO: 保存权重到后端
              console.log('保存权重:', localWeights);
            }}
          >
            保存更改
          </Button>
        </div>
      )}

      {/* 说明文字 */}
      <div className="mt-4 text-xs text-white/40 border-t border-white/5 pt-3">
        💡 影响因子基于行星尊贵度、逆行状态、相位强度等天文参数计算，
        所有数据均有占星学理论支撑。
      </div>
    </div>
  );
}

export default InfluenceFactorsPanel8274TU;

