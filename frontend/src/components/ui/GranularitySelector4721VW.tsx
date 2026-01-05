/**
 * GranularitySelector4721VW - 时间粒度选择器
 * 支持 小时/天/周/月/年 五级切换，类似K线时间切换
 */

import { motion } from 'framer-motion';

export type TimeGranularity = 'hourly' | 'daily' | 'weekly' | 'monthly' | 'yearly';

interface GranularitySelectorProps {
  value: TimeGranularity;
  onChange: (granularity: TimeGranularity) => void;
  className?: string;
  disabled?: TimeGranularity[];
  size?: 'sm' | 'md' | 'lg';
}

const granularityOptions: { value: TimeGranularity; label: string; shortLabel: string; icon: string }[] = [
  { value: 'hourly', label: '小时', shortLabel: '时', icon: '⏱️' },
  { value: 'daily', label: '日', shortLabel: '日', icon: '📅' },
  { value: 'weekly', label: '周', shortLabel: '周', icon: '📆' },
  { value: 'monthly', label: '月', shortLabel: '月', icon: '🗓️' },
  { value: 'yearly', label: '年', shortLabel: '年', icon: '📊' },
];

export function GranularitySelector4721VW({
  value,
  onChange,
  className = '',
  disabled = [],
  size = 'md',
}: GranularitySelectorProps) {
  const sizeStyles = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-3 py-1.5 text-sm',
    lg: 'px-4 py-2 text-base',
  };

  return (
    <div className={`inline-flex rounded-lg bg-white/5 p-1 ${className}`}>
      {granularityOptions.map((option) => {
        const isActive = value === option.value;
        const isDisabled = disabled.includes(option.value);

        return (
          <motion.button
            key={option.value}
            onClick={() => !isDisabled && onChange(option.value)}
            disabled={isDisabled}
            className={`
              relative ${sizeStyles[size]} rounded-md font-medium transition-colors
              ${isDisabled ? 'opacity-30 cursor-not-allowed' : 'cursor-pointer'}
              ${isActive ? 'text-white' : 'text-celestial-silver/60 hover:text-celestial-silver'}
            `}
            whileHover={!isDisabled ? { scale: 1.02 } : undefined}
            whileTap={!isDisabled ? { scale: 0.98 } : undefined}
          >
            {isActive && (
              <motion.div
                layoutId="granularity-active"
                className="absolute inset-0 rounded-md bg-gradient-to-r from-cosmic-nova/80 to-cosmic-aurora/80"
                transition={{ type: 'spring', bounce: 0.2, duration: 0.4 }}
              />
            )}
            <span className="relative z-10 flex items-center gap-1">
              <span className="hidden sm:inline">{option.icon}</span>
              <span className="hidden md:inline">{option.label}</span>
              <span className="md:hidden">{option.shortLabel}</span>
            </span>
          </motion.button>
        );
      })}
    </div>
  );
}

// 紧凑版本 - 仅图标
export function GranularitySelectorCompact4721VW({
  value,
  onChange,
  className = '',
  disabled = [],
}: Omit<GranularitySelectorProps, 'size'>) {
  return (
    <div className={`inline-flex rounded-lg bg-white/5 p-0.5 ${className}`}>
      {granularityOptions.map((option) => {
        const isActive = value === option.value;
        const isDisabled = disabled.includes(option.value);

        return (
          <motion.button
            key={option.value}
            onClick={() => !isDisabled && onChange(option.value)}
            disabled={isDisabled}
            title={option.label}
            className={`
              relative w-8 h-8 rounded-md flex items-center justify-center
              ${isDisabled ? 'opacity-30 cursor-not-allowed' : 'cursor-pointer'}
              ${isActive ? 'text-white' : 'text-celestial-silver/60 hover:text-celestial-silver'}
            `}
            whileHover={!isDisabled ? { scale: 1.1 } : undefined}
            whileTap={!isDisabled ? { scale: 0.9 } : undefined}
          >
            {isActive && (
              <motion.div
                layoutId="granularity-compact-active"
                className="absolute inset-0 rounded-md bg-cosmic-nova/80"
                transition={{ type: 'spring', bounce: 0.2, duration: 0.4 }}
              />
            )}
            <span className="relative z-10 text-sm">{option.icon}</span>
          </motion.button>
        );
      })}
    </div>
  );
}

export default GranularitySelector4721VW;

