/**
 * ScoreBreakdownPopup5932MN - 轻量分数组成浮窗
 * - 小时粒度：显示分数 + 因子贡献值
 * - 日/周/月/年粒度：显示正面/负面影响因子列表
 * 跟随点击位置展示，类似右键菜单
 */
import { useEffect, useRef } from 'react';
import type { ScoreBreakdownAllResponse, ScoreBreakdownResponse, ActiveFactorsResponse } from '../../types';
import { DIMENSION_NAMES } from '../../utils/astro';

type DimensionType = 'overall' | 'career' | 'relationship' | 'health' | 'finance' | 'spiritual';
type GranularityType = 'hour' | 'day' | 'week' | 'month' | 'year';

function pickHourBreakdown(data: ScoreBreakdownAllResponse | null): ScoreBreakdownResponse | null {
  if (!data?.breakdown) return null;
  return (data.breakdown.hour ?? null) as ScoreBreakdownResponse | null;
}

// 格式化时间显示
// 使用 UTC 时间方法，与 lightweight-charts 图表显示保持一致
function formatTimeLabel(timeStr: string, granularity: GranularityType): string {
  // 如果已经是格式化好的显示时间（如 "2020年 (30岁)"），直接返回
  if (timeStr.includes('岁') || timeStr.includes('年 (')) {
    return timeStr;
  }
  
  const d = new Date(timeStr);
  if (isNaN(d.getTime())) return timeStr;
  
  // 使用 UTC 时间（lightweight-charts 使用 Unix 时间戳，默认按 UTC 显示）
  const year = d.getUTCFullYear();
  const month = d.getUTCMonth() + 1;
  const day = d.getUTCDate();
  const hour = d.getUTCHours();
  
  switch (granularity) {
    case 'hour':
      return `${year}年${month}月${day}日${hour}时`;
    case 'day':
      return `${year}年${month}月${day}日`;
    case 'week':
      return `${year}年${month}月${day}日 周`;
    case 'month':
      return `${year}年${month}月`;
    case 'year':
      return `${year}年`;
    default:
      return `${year}年${month}月${day}日${hour}时`;
  }
}

// 粒度中文名
const GRANULARITY_NAMES: Record<GranularityType, string> = {
  hour: '小时',
  day: '日',
  week: '周',
  month: '月',
  year: '年',
};

export interface ScoreBreakdownPopup5932MNProps {
  open: boolean;
  position: { x: number; y: number }; // 点击位置
  queryTime: string | null;
  loading: boolean;
  error: string | null;
  // 小时粒度使用
  data: ScoreBreakdownAllResponse | null;
  // 日/周/月/年粒度使用
  activeFactorsData: ActiveFactorsResponse | null;
  granularity: GranularityType;
  dimension: DimensionType;
  onClose: () => void;
}

export function ScoreBreakdownPopup5932MN({
  open,
  position,
  queryTime,
  loading,
  error,
  data,
  activeFactorsData,
  granularity,
  dimension,
  onClose,
}: ScoreBreakdownPopup5932MNProps) {
  const popupRef = useRef<HTMLDivElement>(null);

  // 点击外部关闭
  useEffect(() => {
    if (!open) return;
    const handleClick = (e: MouseEvent) => {
      if (popupRef.current && !popupRef.current.contains(e.target as Node)) {
        onClose();
      }
    };
    // 延迟添加监听，避免立即触发
    const timer = setTimeout(() => {
      document.addEventListener('click', handleClick);
    }, 100);
    return () => {
      clearTimeout(timer);
      document.removeEventListener('click', handleClick);
    };
  }, [open, onClose]);

  if (!open) return null;

  const isHourly = granularity === 'hour';
  const dimLabel = dimension === 'overall' ? '综合' : (DIMENSION_NAMES[dimension as keyof typeof DIMENSION_NAMES] ?? dimension);

  // 计算浮窗位置（避免超出屏幕）
  const popupWidth = 280;
  const popupHeight = isHourly ? 180 : 260;
  let left = position.x + 10;
  let top = position.y + 10;
  
  if (typeof window !== 'undefined') {
    if (left + popupWidth > window.innerWidth - 20) {
      left = position.x - popupWidth - 10;
    }
    if (top + popupHeight > window.innerHeight - 20) {
      top = position.y - popupHeight - 10;
    }
    if (left < 10) left = 10;
    if (top < 10) top = 10;
  }

  // ==================== 小时粒度：显示分数 + 因子贡献值 ====================
  const renderHourlyContent = () => {
    const current = pickHourBreakdown(data);
    if (!current) return <div className="text-white/50 text-sm">暂无数据</div>;

    // 获取分数
    let score: number | null = null;
    if (dimension === 'overall') {
      score = current.overallScore;
    } else {
      const dimData = current.dimensions?.find(d => d.dimension === dimension);
      score = dimData?.finalScore ?? null;
    }

    // 根据选中的维度收集对应的因子
    const allFactors: { name: string; value: number }[] = [];
    
    if (dimension === 'overall') {
      if (current?.factorsByLevel) {
        Object.values(current.factorsByLevel).forEach(factors => {
          factors?.forEach(f => {
            allFactors.push({ name: f.name, value: f.adjustment });
          });
        });
      }
    } else {
      const dimData = current?.dimensions?.find(d => d.dimension === dimension);
      if (dimData?.factors) {
        dimData.factors.forEach(f => {
          allFactors.push({ name: f.name, value: f.adjustment });
        });
      }
    }

    allFactors.sort((a, b) => Math.abs(b.value) - Math.abs(a.value));
    const topFactors = allFactors.slice(0, 5);

    return (
      <div className="space-y-2">
        {/* 分数 */}
        <div className="flex items-baseline">
          <span className="text-white/70 text-sm">{dimLabel}分：</span>
          <span className="text-cosmic-nova text-xl font-bold ml-1">
            {score !== null ? score.toFixed(1) : '--'}
          </span>
        </div>
        
        {/* 时间 */}
        <div className="text-white/50 text-xs">
          时间：{queryTime ? formatTimeLabel(queryTime, granularity) : '--'}
        </div>

        {/* 主要影响因素 */}
        {topFactors.length > 0 ? (
          <div className="pt-1">
            <div className="text-white/60 text-xs mb-1">主要影响：</div>
            <div className="text-xs text-white/80 space-y-0.5">
              {topFactors.map((f, i) => (
                <div key={i} className="flex justify-between">
                  <span className="truncate mr-2">{f.name}</span>
                  <span className={`flex-shrink-0 ${f.value >= 0 ? 'text-green-400' : 'text-rose-400'}`}>
                    {f.value >= 0 ? '+' : ''}{f.value.toFixed(1)}
                  </span>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className="text-white/40 text-xs pt-1">无影响因子</div>
        )}
      </div>
    );
  };

  // ==================== 日/周/月/年粒度：显示正面/负面影响因子 ====================
  const renderActiveFactorsContent = () => {
    if (!activeFactorsData) return <div className="text-white/50 text-sm">暂无数据</div>;

    const { factors, positiveCount, negativeCount } = activeFactorsData;
    
    // 分离正面和负面因子
    const positiveFactors = factors.filter(f => f.isPositive);
    const negativeFactors = factors.filter(f => !f.isPositive);

    // 按影响强度排序（maxStrength * |baseValue|）
    positiveFactors.sort((a, b) => (b.maxStrength * Math.abs(b.baseValue)) - (a.maxStrength * Math.abs(a.baseValue)));
    negativeFactors.sort((a, b) => (b.maxStrength * Math.abs(b.baseValue)) - (a.maxStrength * Math.abs(a.baseValue)));

    // 最多显示前5个
    const topPositive = positiveFactors.slice(0, 5);
    const topNegative = negativeFactors.slice(0, 5);

    return (
      <div className="space-y-2">
        {/* 时间范围 */}
        <div>
          <div className="text-white/70 text-sm font-medium">
            {GRANULARITY_NAMES[granularity]}度影响因子
          </div>
          <div className="text-white/50 text-xs mt-0.5">
            {queryTime ? formatTimeLabel(queryTime, granularity) : '--'}
          </div>
        </div>

        {/* 正面影响 */}
        {topPositive.length > 0 && (
          <div>
            <div className="flex items-center text-green-400 text-xs mb-1">
              <span className="mr-1">✓</span>
              <span>正面影响 ({positiveCount})</span>
            </div>
            <div className="text-xs text-white/80 space-y-0.5 pl-3">
              {topPositive.map((f, i) => (
                <div key={i} className="truncate">{f.name}</div>
              ))}
              {positiveCount > 5 && (
                <div className="text-white/40">+{positiveCount - 5} 更多...</div>
              )}
            </div>
          </div>
        )}

        {/* 负面影响 */}
        {topNegative.length > 0 && (
          <div>
            <div className="flex items-center text-rose-400 text-xs mb-1">
              <span className="mr-1">✗</span>
              <span>负面影响 ({negativeCount})</span>
            </div>
            <div className="text-xs text-white/80 space-y-0.5 pl-3">
              {topNegative.map((f, i) => (
                <div key={i} className="truncate">{f.name}</div>
              ))}
              {negativeCount > 5 && (
                <div className="text-white/40">+{negativeCount - 5} 更多...</div>
              )}
            </div>
          </div>
        )}

        {/* 无因子时 */}
        {positiveCount === 0 && negativeCount === 0 && (
          <div className="text-white/40 text-xs">该时段无活跃影响因子</div>
        )}
      </div>
    );
  };

  return (
    <div
      ref={popupRef}
      className="fixed z-[100] glass-card border border-white/20 rounded-lg shadow-xl"
      style={{
        left,
        top,
        width: popupWidth,
        maxWidth: 'calc(100vw - 40px)',
      }}
    >
      <div className="p-3">
        {loading ? (
          <div className="text-white/60 text-sm">加载中...</div>
        ) : error ? (
          <div className="text-red-300 text-sm">{error}</div>
        ) : isHourly ? (
          renderHourlyContent()
        ) : (
          renderActiveFactorsContent()
        )}
      </div>
    </div>
  );
}

export default ScoreBreakdownPopup5932MN;
