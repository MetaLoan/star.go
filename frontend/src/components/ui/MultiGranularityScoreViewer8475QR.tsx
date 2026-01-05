/**
 * 多时间粒度分数查看器组件
 * 组件命名规范：MultiGranularityScoreViewer + 8475 + QR
 * 
 * 功能：
 * - 展示小时/日/月/年四个时间粒度的五维值和综合值
 * - 支持回溯查询历史数据
 * - 支持查询未来预测数据
 * - 动画效果展示分数变化
 */

import { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Button, Spinner } from '@heroui/react';
import type { BirthData } from '../../types';
import { DIMENSION_NAMES, DIMENSION_COLORS, DIMENSION_ICONS } from '../../utils/astro';

// 时间粒度类型
type TimeGranularity = 'hourly' | 'daily' | 'monthly' | 'yearly';

// 维度数据类型
interface DimensionData {
  career: number;
  relationship: number;
  health: number;
  finance: number;
  spiritual: number;
}

// 分数数据类型
interface ScoreData {
  overall: number;
  dimensions: DimensionData;
  label: string;
  time: string;
}

// 多粒度数据
interface MultiGranularityData {
  hourly: ScoreData | null;
  daily: ScoreData | null;
  monthly: ScoreData | null;
  yearly: ScoreData | null;
}

interface MultiGranularityScoreViewerProps {
  birthData: BirthData;
  className?: string;
}

// 粒度配置
const GRANULARITY_CONFIG: Record<TimeGranularity, {
  label: string;
  icon: string;
  color: string;
  apiGranularity: string;
}> = {
  hourly: { label: '小时', icon: '⏰', color: '#00D4FF', apiGranularity: 'hour' },
  daily: { label: '日', icon: '📅', color: '#4ECDC4', apiGranularity: 'day' },
  monthly: { label: '月', icon: '🗓️', color: '#FFE66D', apiGranularity: 'month' },
  yearly: { label: '年', icon: '📆', color: '#A855F7', apiGranularity: 'year' },
};

// 格式化本地 ISO 时间
const formatLocalISO = (date: Date, timezone: number = 8): string => {
  const offsetHours = Math.floor(Math.abs(timezone));
  const offsetMins = Math.round((Math.abs(timezone) % 1) * 60);
  const sign = timezone >= 0 ? '+' : '-';
  const pad = (n: number) => n.toString().padStart(2, '0');

  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}${sign}${pad(offsetHours)}:${pad(offsetMins)}`;
};

// 获取时间范围
const getTimeRange = (date: Date, granularity: TimeGranularity, timezone: number): { start: string; end: string } => {
  const start = new Date(date);
  const end = new Date(date);

  switch (granularity) {
    case 'hourly':
      start.setMinutes(0, 0, 0);
      end.setHours(end.getHours() + 1);
      end.setMinutes(0, 0, 0);
      break;
    case 'daily':
      start.setHours(0, 0, 0, 0);
      end.setDate(end.getDate() + 1);
      end.setHours(0, 0, 0, 0);
      break;
    case 'monthly':
      start.setDate(1);
      start.setHours(0, 0, 0, 0);
      end.setMonth(end.getMonth() + 1, 1);
      end.setHours(0, 0, 0, 0);
      break;
    case 'yearly':
      start.setMonth(0, 1);
      start.setHours(0, 0, 0, 0);
      end.setFullYear(end.getFullYear() + 1, 0, 1);
      end.setHours(0, 0, 0, 0);
      break;
  }

  return {
    start: formatLocalISO(start, timezone),
    end: formatLocalISO(end, timezone),
  };
};

// 格式化显示时间
const formatDisplayTime = (date: Date, granularity: TimeGranularity): string => {
  const pad = (n: number) => n.toString().padStart(2, '0');
  switch (granularity) {
    case 'hourly':
      return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${pad(date.getHours())}:00`;
    case 'daily':
      return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`;
    case 'monthly':
      return `${date.getFullYear()}年${date.getMonth() + 1}月`;
    case 'yearly':
      return `${date.getFullYear()}年`;
  }
};

// 调整时间
const adjustDate = (date: Date, granularity: TimeGranularity, delta: number): Date => {
  const newDate = new Date(date);
  switch (granularity) {
    case 'hourly':
      newDate.setHours(newDate.getHours() + delta);
      break;
    case 'daily':
      newDate.setDate(newDate.getDate() + delta);
      break;
    case 'monthly':
      newDate.setMonth(newDate.getMonth() + delta);
      break;
    case 'yearly':
      newDate.setFullYear(newDate.getFullYear() + delta);
      break;
  }
  return newDate;
};

// 维度分数条
const DimensionBar = ({
  dimension,
  score,
  compact = false,
}: {
  dimension: keyof DimensionData;
  score: number;
  compact?: boolean;
}) => {
  const color = DIMENSION_COLORS[dimension];
  const icon = DIMENSION_ICONS[dimension];
  const name = DIMENSION_NAMES[dimension];

  return (
    <div className={`flex items-center gap-2 ${compact ? 'py-1' : 'py-2'}`}>
      <span className={compact ? 'text-sm' : 'text-lg'}>{icon}</span>
      <span className={`${compact ? 'text-xs w-10' : 'text-sm w-12'} text-white/70`}>{name}</span>
      <div className="flex-1 h-2 bg-white/10 rounded-full overflow-hidden">
        <motion.div
          className="h-full rounded-full"
          style={{ backgroundColor: color }}
          initial={{ width: 0 }}
          animate={{ width: `${Math.min(score, 100)}%` }}
          transition={{ duration: 0.8, ease: 'easeOut' }}
        />
      </div>
      <motion.span
        key={score}
        initial={{ opacity: 0, x: 10 }}
        animate={{ opacity: 1, x: 0 }}
        className={`${compact ? 'text-xs w-8' : 'text-sm w-10'} text-right font-medium`}
        style={{ color }}
      >
        {Math.round(score)}
      </motion.span>
    </div>
  );
};

// 综合分数显示
const OverallScore = ({ score, label, color }: { score: number; label: string; color: string }) => {
  return (
    <div className="text-center">
      <motion.div
        key={score}
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        className="text-4xl font-bold mb-1"
        style={{ color }}
      >
        {Math.round(score)}
      </motion.div>
      <div className="text-xs text-white/50">{label}综合分</div>
    </div>
  );
};

// 单个粒度卡片
const GranularityCard = ({
  granularity,
  data,
  isLoading,
  isSelected,
  onClick,
}: {
  granularity: TimeGranularity;
  data: ScoreData | null;
  isLoading: boolean;
  isSelected: boolean;
  onClick: () => void;
}) => {
  const config = GRANULARITY_CONFIG[granularity];

  return (
    <motion.div
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
      onClick={onClick}
      className={`cursor-pointer rounded-2xl p-4 transition-all duration-300 ${
        isSelected
          ? 'ring-2 ring-white/30 bg-white/10'
          : 'bg-white/5 hover:bg-white/8'
      }`}
      style={{
        borderLeft: `4px solid ${config.color}`,
      }}
    >
      {/* 头部 */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className="text-xl">{config.icon}</span>
          <span className="font-medium text-white">{config.label}数据</span>
        </div>
        {isLoading ? (
          <Spinner size="sm" color="primary" />
        ) : data ? (
          <span
            className="text-2xl font-bold"
            style={{ color: config.color }}
          >
            {Math.round(data.overall)}
          </span>
        ) : (
          <span className="text-white/30">--</span>
        )}
      </div>

      {/* 五维度迷你条 */}
      {data && !isLoading && (
        <div className="space-y-0.5">
          {(Object.keys(data.dimensions) as (keyof DimensionData)[]).map((dim) => (
            <DimensionBar
              key={dim}
              dimension={dim}
              score={data.dimensions[dim]}
              compact={true}
            />
          ))}
        </div>
      )}

      {/* 时间标签 */}
      {data && (
        <div className="mt-2 text-xs text-white/40 text-right">
          {data.label}
        </div>
      )}
    </motion.div>
  );
};

// 详细展示面板
const DetailPanel = ({
  granularity,
  data,
  selectedDate,
  onPrev,
  onNext,
  onDateChange,
}: {
  granularity: TimeGranularity;
  data: ScoreData | null;
  selectedDate: Date;
  onPrev: () => void;
  onNext: () => void;
  onDateChange: (date: Date) => void;
}) => {
  const config = GRANULARITY_CONFIG[granularity];
  const displayTime = formatDisplayTime(selectedDate, granularity);
  
  // 判断是否是未来
  const isFuture = selectedDate > new Date();

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="glass-card p-6 rounded-2xl"
      style={{
        background: `linear-gradient(135deg, ${config.color}10, transparent)`,
        borderTop: `3px solid ${config.color}`,
      }}
    >
      {/* 时间导航 */}
      <div className="flex items-center justify-between mb-6">
        <Button
          isIconOnly
          variant="flat"
          size="sm"
          onPress={onPrev}
          className="bg-white/10 hover:bg-white/20"
        >
          ◀
        </Button>
        
        <div className="text-center flex-1">
          <div className="flex items-center justify-center gap-2 mb-1">
            <span className="text-2xl">{config.icon}</span>
            <span className="text-xl font-bold text-white">{displayTime}</span>
            {isFuture && (
              <span className="px-2 py-0.5 bg-purple-500/20 text-purple-300 text-xs rounded-full">
                预测
              </span>
            )}
          </div>
          <div className="text-xs text-white/50">
            {config.label}级别数据 | 数据来源: Swiss Ephemeris
          </div>
        </div>
        
        <Button
          isIconOnly
          variant="flat"
          size="sm"
          onPress={onNext}
          className="bg-white/10 hover:bg-white/20"
        >
          ▶
        </Button>
      </div>

      {/* 综合分数 */}
      {data && (
        <div className="grid md:grid-cols-[1fr_2fr] gap-6">
          {/* 左侧：综合分 */}
          <div className="flex flex-col items-center justify-center p-6 bg-white/5 rounded-xl">
            <OverallScore
              score={data.overall}
              label={config.label}
              color={config.color}
            />
            <div className="mt-4 w-full">
              <div className="h-3 bg-white/10 rounded-full overflow-hidden">
                <motion.div
                  className="h-full rounded-full"
                  style={{
                    background: `linear-gradient(90deg, ${config.color}80, ${config.color})`,
                  }}
                  initial={{ width: 0 }}
                  animate={{ width: `${Math.min(data.overall, 100)}%` }}
                  transition={{ duration: 1, ease: 'easeOut' }}
                />
              </div>
            </div>
            <div className="mt-3 text-xs text-white/50 text-center">
              {data.overall >= 80 ? '🌟 运势极佳' : 
               data.overall >= 60 ? '✨ 运势良好' : 
               data.overall >= 40 ? '💫 运势平稳' : '🌙 需要注意'}
            </div>
          </div>

          {/* 右侧：五维度 */}
          <div className="space-y-2">
            <h3 className="text-sm font-medium text-white/80 mb-3">五维度分数</h3>
            {(Object.keys(data.dimensions) as (keyof DimensionData)[]).map((dim) => (
              <DimensionBar
                key={dim}
                dimension={dim}
                score={data.dimensions[dim]}
              />
            ))}
          </div>
        </div>
      )}

      {/* 快速跳转 */}
      <div className="mt-6 pt-4 border-t border-white/10 flex items-center justify-center gap-2 flex-wrap">
        <span className="text-xs text-white/40 mr-2">快速跳转:</span>
        <Button
          size="sm"
          variant="flat"
          className="bg-white/5 hover:bg-white/10 text-xs"
          onPress={() => onDateChange(new Date())}
        >
          今天
        </Button>
        <Button
          size="sm"
          variant="flat"
          className="bg-white/5 hover:bg-white/10 text-xs"
          onPress={() => onDateChange(adjustDate(new Date(), granularity, -1))}
        >
          上一{config.label.replace('小', '')}
        </Button>
        <Button
          size="sm"
          variant="flat"
          className="bg-white/5 hover:bg-white/10 text-xs"
          onPress={() => onDateChange(adjustDate(new Date(), granularity, 1))}
        >
          下一{config.label.replace('小', '')}
        </Button>
        <Button
          size="sm"
          variant="flat"
          className="bg-white/5 hover:bg-white/10 text-xs"
          onPress={() => onDateChange(adjustDate(new Date(), granularity, -7))}
        >
          前7{config.label.replace('小', '')}
        </Button>
        <Button
          size="sm"
          variant="flat"
          className="bg-white/5 hover:bg-white/10 text-xs"
          onPress={() => onDateChange(adjustDate(new Date(), granularity, 7))}
        >
          后7{config.label.replace('小', '')}
        </Button>
      </div>
    </motion.div>
  );
};

// 主组件
export function MultiGranularityScoreViewer8475QR({
  birthData,
  className = '',
}: MultiGranularityScoreViewerProps) {
  const [selectedGranularity, setSelectedGranularity] = useState<TimeGranularity>('daily');
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const [data, setData] = useState<MultiGranularityData>({
    hourly: null,
    daily: null,
    monthly: null,
    yearly: null,
  });
  const [loading, setLoading] = useState<Record<TimeGranularity, boolean>>({
    hourly: false,
    daily: false,
    monthly: false,
    yearly: false,
  });
  const [detailData, setDetailData] = useState<ScoreData | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  // 获取指定粒度的数据
  const fetchGranularityData = useCallback(async (
    granularity: TimeGranularity,
    targetDate: Date
  ): Promise<ScoreData | null> => {
    try {
      const { start, end } = getTimeRange(targetDate, granularity, birthData.timezone);
      const config = GRANULARITY_CONFIG[granularity];

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
          start,
          end,
          granularity: config.apiGranularity,
        }),
      });

      if (!response.ok) throw new Error('Failed to fetch data');

      const result = await response.json();
      
      if (result.points && result.points.length > 0) {
        const point = result.points[0];
        return {
          overall: point.display || 50,
          dimensions: point.dimensions || {
            career: 50,
            relationship: 50,
            health: 50,
            finance: 50,
            spiritual: 50,
          },
          label: formatDisplayTime(targetDate, granularity),
          time: point.time,
        };
      }
      return null;
    } catch (err) {
      console.error(`获取${granularity}数据失败:`, err);
      return null;
    }
  }, [birthData]);

  // 加载所有粒度的当前数据
  const loadAllCurrentData = useCallback(async () => {
    const now = new Date();
    const granularities: TimeGranularity[] = ['hourly', 'daily', 'monthly', 'yearly'];

    setLoading({
      hourly: true,
      daily: true,
      monthly: true,
      yearly: true,
    });

    const results = await Promise.all(
      granularities.map(g => fetchGranularityData(g, now))
    );

    setData({
      hourly: results[0],
      daily: results[1],
      monthly: results[2],
      yearly: results[3],
    });

    setLoading({
      hourly: false,
      daily: false,
      monthly: false,
      yearly: false,
    });
  }, [fetchGranularityData]);

  // 加载详细数据
  const loadDetailData = useCallback(async () => {
    setDetailLoading(true);
    const result = await fetchGranularityData(selectedGranularity, selectedDate);
    setDetailData(result);
    setDetailLoading(false);
  }, [fetchGranularityData, selectedGranularity, selectedDate]);

  // 初始加载
  useEffect(() => {
    loadAllCurrentData();
  }, [loadAllCurrentData]);

  // 当选择的粒度或日期变化时，加载详细数据
  useEffect(() => {
    loadDetailData();
  }, [loadDetailData]);

  // 处理粒度选择
  const handleGranularitySelect = (granularity: TimeGranularity) => {
    setSelectedGranularity(granularity);
    setSelectedDate(new Date()); // 重置到当前时间
  };

  // 处理时间导航
  const handlePrev = () => {
    setSelectedDate(prev => adjustDate(prev, selectedGranularity, -1));
  };

  const handleNext = () => {
    setSelectedDate(prev => adjustDate(prev, selectedGranularity, 1));
  };

  const handleDateChange = (date: Date) => {
    setSelectedDate(date);
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className={`space-y-6 ${className}`}
    >
      {/* 标题 */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-bold text-white flex items-center gap-2">
          <span>📊</span>
          <span>多粒度运势查询</span>
        </h2>
        <Button
          size="sm"
          variant="flat"
          className="bg-white/10"
          onPress={loadAllCurrentData}
        >
          🔄 刷新全部
        </Button>
      </div>

      {/* 四粒度概览卡片 */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {(['hourly', 'daily', 'monthly', 'yearly'] as TimeGranularity[]).map((g) => (
          <GranularityCard
            key={g}
            granularity={g}
            data={data[g]}
            isLoading={loading[g]}
            isSelected={selectedGranularity === g}
            onClick={() => handleGranularitySelect(g)}
          />
        ))}
      </div>

      {/* 详细展示面板 */}
      <AnimatePresence mode="wait">
        {detailLoading ? (
          <motion.div
            key="loading"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="glass-card p-12 flex items-center justify-center"
          >
            <Spinner size="lg" color="primary" />
            <span className="ml-3 text-white/60">加载数据中...</span>
          </motion.div>
        ) : detailData ? (
          <DetailPanel
            key={`${selectedGranularity}-${selectedDate.toISOString()}`}
            granularity={selectedGranularity}
            data={detailData}
            selectedDate={selectedDate}
            onPrev={handlePrev}
            onNext={handleNext}
            onDateChange={handleDateChange}
          />
        ) : (
          <motion.div
            key="no-data"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="glass-card p-12 text-center text-white/50"
          >
            暂无数据
          </motion.div>
        )}
      </AnimatePresence>

      {/* 说明 */}
      <div className="text-xs text-white/40 text-center space-y-1">
        <p>💡 点击上方卡片切换时间粒度，使用 ◀ ▶ 按钮或快速跳转浏览不同时间段</p>
        <p>📡 所有数据基于 Swiss Ephemeris 天文算法实时计算</p>
      </div>
    </motion.div>
  );
}

export default MultiGranularityScoreViewer8475QR;

