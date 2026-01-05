/**
 * FactorLevelFilter6293ZA - 因子时间级别筛选器
 * 支持筛选 年/月/周/日/时 级别的因子
 */

import { motion } from 'framer-motion';

export type FactorTimeLevel = 'yearly' | 'monthly' | 'weekly' | 'daily' | 'hourly';

interface FactorLevelFilterProps {
  selected: FactorTimeLevel[];
  onChange: (levels: FactorTimeLevel[]) => void;
  viewLevel?: FactorTimeLevel; // 当前视图级别，用于控制可见性
  className?: string;
  compact?: boolean;
}

// 因子级别配置
const factorLevels: { value: FactorTimeLevel; label: string; shortLabel: string; color: string; description: string }[] = [
  { value: 'yearly', label: '年度级', shortLabel: '年', color: '#A855F7', description: '土星回归、木星回归等长期因子' },
  { value: 'monthly', label: '月度级', shortLabel: '月', color: '#FF6B9D', description: '行星换座、月相周期等' },
  { value: 'weekly', label: '周度级', shortLabel: '周', color: '#FF9F43', description: '水星逆行、月亮节点等' },
  { value: 'daily', label: '日度级', shortLabel: '日', color: '#4ECDC4', description: '月亮换座、行星日等' },
  { value: 'hourly', label: '小时级', shortLabel: '时', color: '#00D4FF', description: '行星时、月亮空亡等' },
];

// 判断在某视图级别下，某因子级别是否可见
function isLevelVisibleAtView(factorLevel: FactorTimeLevel, viewLevel: FactorTimeLevel): boolean {
  const order: Record<FactorTimeLevel, number> = {
    yearly: 1,
    monthly: 2,
    weekly: 3,
    daily: 4,
    hourly: 5,
  };
  // 大级别因子在小级别视图中可见
  return order[factorLevel] <= order[viewLevel];
}

export function FactorLevelFilter6293ZA({
  selected,
  onChange,
  viewLevel = 'hourly',
  className = '',
  compact = false,
}: FactorLevelFilterProps) {
  const toggleLevel = (level: FactorTimeLevel) => {
    if (selected.includes(level)) {
      onChange(selected.filter((l) => l !== level));
    } else {
      onChange([...selected, level]);
    }
  };

  const selectAll = () => {
    const visibleLevels = factorLevels
      .filter((l) => isLevelVisibleAtView(l.value, viewLevel))
      .map((l) => l.value);
    onChange(visibleLevels);
  };

  const clearAll = () => {
    onChange([]);
  };

  const visibleLevels = factorLevels.filter((l) => isLevelVisibleAtView(l.value, viewLevel));

  if (compact) {
    return (
      <div className={`flex flex-wrap gap-1 ${className}`}>
        {visibleLevels.map((level) => {
          const isSelected = selected.includes(level.value);
          return (
            <motion.button
              key={level.value}
              onClick={() => toggleLevel(level.value)}
              className={`
                px-2 py-0.5 rounded-full text-xs font-medium transition-all
                ${isSelected
                  ? 'text-white shadow-lg'
                  : 'text-celestial-silver/60 bg-white/5 hover:bg-white/10'}
              `}
              style={isSelected ? { backgroundColor: level.color } : undefined}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              {level.shortLabel}
            </motion.button>
          );
        })}
      </div>
    );
  }

  return (
    <div className={`space-y-3 ${className}`}>
      {/* 标题和操作按钮 */}
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-celestial-silver/80">因子级别</span>
        <div className="flex gap-2">
          <button
            onClick={selectAll}
            className="text-xs text-cosmic-nova hover:text-cosmic-nova/80 transition-colors"
          >
            全选
          </button>
          <button
            onClick={clearAll}
            className="text-xs text-celestial-silver/60 hover:text-celestial-silver transition-colors"
          >
            清空
          </button>
        </div>
      </div>

      {/* 级别列表 */}
      <div className="space-y-2">
        {visibleLevels.map((level) => {
          const isSelected = selected.includes(level.value);
          return (
            <motion.button
              key={level.value}
              onClick={() => toggleLevel(level.value)}
              className={`
                w-full flex items-center gap-3 p-2 rounded-lg transition-all
                ${isSelected
                  ? 'bg-white/10 border border-white/20'
                  : 'bg-white/5 border border-transparent hover:bg-white/8'}
              `}
              whileHover={{ x: 4 }}
              whileTap={{ scale: 0.98 }}
            >
              {/* 颜色指示器 */}
              <div
                className={`w-3 h-3 rounded-full transition-all ${isSelected ? 'ring-2 ring-white/50' : ''}`}
                style={{ backgroundColor: level.color }}
              />

              {/* 标签 */}
              <div className="flex-1 text-left">
                <div className="text-sm font-medium">{level.label}</div>
                <div className="text-xs text-celestial-silver/50">{level.description}</div>
              </div>

              {/* 选中状态 */}
              <div
                className={`w-5 h-5 rounded border-2 flex items-center justify-center transition-colors
                  ${isSelected ? 'bg-cosmic-nova border-cosmic-nova' : 'border-celestial-silver/30'}`}
              >
                {isSelected && (
                  <motion.svg
                    width="12"
                    height="12"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="3"
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                  >
                    <polyline points="20 6 9 17 4 12" />
                  </motion.svg>
                )}
              </div>
            </motion.button>
          );
        })}
      </div>

      {/* 当前视图说明 */}
      <div className="text-xs text-celestial-silver/40 text-center pt-2 border-t border-white/10">
        当前视图：{factorLevels.find((l) => l.value === viewLevel)?.label}
        <br />
        显示该级别及更高级别的因子
      </div>
    </div>
  );
}

export default FactorLevelFilter6293ZA;

