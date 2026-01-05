/**
 * ScoreBreakdownDrawer4827LX - 分数组成浮窗（抽屉）
 * 点击趋势图数据点后，通过 /api/calc/score-breakdown-all 展示分数组成
 */
import type { ScoreBreakdownAllResponse, ScoreBreakdownResponse } from '../../types';
import { DIMENSION_NAMES } from '../../utils/astro';

type BreakdownKey = 'hour' | 'day' | 'month' | 'year';

const KEY_LABEL: Record<BreakdownKey, string> = {
  hour: '小时',
  day: '日',
  month: '月',
  year: '年',
};

const LEVEL_LABEL: Record<string, string> = {
  yearly: '年度',
  monthly: '月度',
  weekly: '周度',
  daily: '日度',
  hourly: '小时',
  custom: '自定义',
};

function pickBreakdown(data: ScoreBreakdownAllResponse | null, key: BreakdownKey): ScoreBreakdownResponse | null {
  if (!data?.breakdown) return null;
  return (data.breakdown[key] ?? null) as ScoreBreakdownResponse | null;
}

type DimensionType = 'overall' | 'career' | 'relationship' | 'health' | 'finance' | 'spiritual';

export interface ScoreBreakdownDrawer4827LXProps {
  open: boolean;
  queryTime: string | null;
  loading: boolean;
  error: string | null;
  data: ScoreBreakdownAllResponse | null;
  selectedKey: BreakdownKey;
  onSelectKey: (k: BreakdownKey) => void;
  selectedDimension?: DimensionType; // 用户点击时选中的维度
  onClose: () => void;
}

export function ScoreBreakdownDrawer4827LX({
  open,
  queryTime,
  loading,
  error,
  data,
  selectedKey,
  onSelectKey,
  selectedDimension = 'overall',
  onClose,
}: ScoreBreakdownDrawer4827LXProps) {
  if (!open) return null;

  const current = pickBreakdown(data, selectedKey);
  
  // 找到选中维度的分解数据
  const highlightedDim = current?.dimensions?.find(d => d.dimension === selectedDimension);
  const dimensionLabel = selectedDimension === 'overall' ? '综合' : (DIMENSION_NAMES[selectedDimension as keyof typeof DIMENSION_NAMES] ?? selectedDimension);

  return (
    <div className="fixed inset-0 z-50">
      {/* 遮罩 */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
      />

      {/* 右侧抽屉 */}
      <div className="absolute right-0 top-0 h-full w-full max-w-[520px] glass-card border-l border-white/10 overflow-hidden">
        <div className="p-4 border-b border-white/10 flex items-start justify-between gap-3">
          <div>
            <div className="text-white/90 font-medium">分数组成</div>
            <div className="text-xs text-white/50 mt-1">
              {queryTime ? `时间：${new Date(queryTime).toLocaleString()}` : '未选择时间点'}
            </div>
          </div>
          <button
            className="px-2 py-1 rounded-lg bg-white/5 hover:bg-white/10 text-white/70 text-sm"
            onClick={onClose}
          >
            关闭
          </button>
        </div>

        {/* 粒度切换 */}
        <div className="p-4 pb-0">
          <div className="flex gap-2">
            {(['hour', 'day', 'month', 'year'] as BreakdownKey[]).map((k) => (
              <button
                key={k}
                onClick={() => onSelectKey(k)}
                className={`px-3 py-1.5 rounded-lg text-sm transition-all border ${
                  selectedKey === k
                    ? 'bg-cosmic-nova/20 border-cosmic-nova text-white'
                    : 'bg-white/5 border-white/10 text-white/60 hover:text-white/80 hover:bg-white/8'
                }`}
              >
                {KEY_LABEL[k]}
              </button>
            ))}
          </div>
          <div className="text-[11px] text-white/35 mt-2">
            注：该接口返回 hour/day/month/year（不含 week）。
          </div>
        </div>

        <div className="p-4 overflow-y-auto h-[calc(100%-120px)]">
          {loading && (
            <div className="p-4 bg-white/5 rounded-lg text-white/70 text-sm">
              正在加载分数组成...
            </div>
          )}

          {!loading && error && (
            <div className="p-4 bg-red-500/15 border border-red-500/20 rounded-lg text-red-200 text-sm">
              {error}
            </div>
          )}

          {!loading && !error && !current && (
            <div className="p-4 bg-white/5 rounded-lg text-white/50 text-sm">
              暂无数据
            </div>
          )}

          {!loading && !error && current && (
            <div className="space-y-4">
              {/* 概览 - 根据选中维度显示 */}
              <div className="bg-white/5 rounded-lg p-4">
                <div className="flex items-center justify-between">
                  <div className="text-white/80 text-sm">
                    当前粒度：{KEY_LABEL[selectedKey]}
                    {selectedDimension !== 'overall' && (
                      <span className="ml-2 px-2 py-0.5 bg-cosmic-nova/20 text-cosmic-nova rounded text-xs">
                        {dimensionLabel}
                      </span>
                    )}
                  </div>
                  <div className="text-white/40 text-xs">{current.meta?.dataSource}</div>
                </div>
                <div className="mt-3 flex items-baseline justify-between">
                  <div>
                    <div className="text-xs text-white/50">
                      {selectedDimension === 'overall' ? '综合分' : `${dimensionLabel}分`}
                    </div>
                    <div className="text-3xl font-bold text-cosmic-nova">
                      {selectedDimension === 'overall' 
                        ? (Number.isFinite(current.overallScore) ? current.overallScore.toFixed(2) : '--')
                        : (highlightedDim ? highlightedDim.finalScore.toFixed(2) : '--')}
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-xs text-white/50">
                      {selectedDimension === 'overall' ? '原始分' : '分解详情'}
                    </div>
                    {selectedDimension === 'overall' ? (
                      <>
                        <div className="text-white/80 font-medium">
                          {Number.isFinite(current.overallRaw) ? current.overallRaw.toFixed(2) : '--'}
                        </div>
                        <div className="text-xs text-white/40 mt-1">
                          因子：{current.meta?.totalFactorCount ?? 0}（+{current.meta?.positiveFactors ?? 0}/-{current.meta?.negativeFactors ?? 0}）
                        </div>
                      </>
                    ) : highlightedDim ? (
                      <div className="text-white/70 text-sm space-y-0.5">
                        <div>基础: {highlightedDim.baseScore.toFixed(2)}</div>
                        <div>相位: {highlightedDim.aspectScore >= 0 ? '+' : ''}{highlightedDim.aspectScore.toFixed(2)}</div>
                        <div>因子: {highlightedDim.factorScore >= 0 ? '+' : ''}{highlightedDim.factorScore.toFixed(2)}</div>
                      </div>
                    ) : null}
                  </div>
                </div>
              </div>

              {/* 五维度 */}
              <div className="bg-white/5 rounded-lg p-4">
                <div className="text-white/80 text-sm mb-3">五维度分解</div>
                <div className="space-y-2">
                  {/* 如果选中了特定维度，把它放在最前面并高亮 */}
                  {current.dimensions?.slice().sort((a, b) => {
                    if (a.dimension === selectedDimension) return -1;
                    if (b.dimension === selectedDimension) return 1;
                    return 0;
                  }).map((d) => {
                    const isHighlighted = d.dimension === selectedDimension && selectedDimension !== 'overall';
                    return (
                      <div 
                        key={d.dimension} 
                        className={`p-3 rounded-lg border ${
                          isHighlighted 
                            ? 'bg-cosmic-nova/10 border-cosmic-nova/50' 
                            : 'bg-white/5 border-white/10'
                        }`}
                      >
                        <div className="flex items-center justify-between">
                          <div className={`text-sm ${isHighlighted ? 'text-cosmic-nova font-medium' : 'text-white/80'}`}>
                            {DIMENSION_NAMES[d.dimension as keyof typeof DIMENSION_NAMES] ?? d.dimension}
                            {isHighlighted && <span className="ml-2 text-xs">← 选中</span>}
                          </div>
                          <div className={`font-semibold ${isHighlighted ? 'text-cosmic-nova' : 'text-white'}`}>
                            {d.finalScore.toFixed(2)}
                          </div>
                        </div>
                        <div className="grid grid-cols-3 gap-2 text-xs text-white/50 mt-2">
                          <div>基础：{d.baseScore.toFixed(2)}</div>
                          <div>相位：{d.aspectScore.toFixed(2)}</div>
                          <div>因子：{d.factorScore.toFixed(2)}</div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>

              {/* 因子（按时间级别分组） */}
              <div className="bg-white/5 rounded-lg p-4">
                <div className="text-white/80 text-sm mb-3">因子（按时间级别）</div>
                <div className="space-y-3">
                  {Object.entries(current.factorsByLevel ?? {}).map(([level, factors]) => {
                    const label = LEVEL_LABEL[level] ?? level;
                    const top = [...(factors ?? [])]
                      .sort((a, b) => Math.abs(b.adjustment) - Math.abs(a.adjustment))
                      .slice(0, 6);
                    return (
                      <div key={level} className="p-3 rounded-lg bg-white/5 border border-white/10">
                        <div className="flex items-center justify-between">
                          <div className="text-white/70 text-sm">{label}</div>
                          <div className="text-white/40 text-xs">{(factors ?? []).length} 个</div>
                        </div>
                        <div className="mt-2 space-y-1">
                          {top.length === 0 ? (
                            <div className="text-xs text-white/40">无</div>
                          ) : (
                            top.map((f) => (
                              <div key={f.id} className="flex items-center justify-between text-xs">
                                <div className="text-white/70 truncate pr-2">{f.name}</div>
                                <div className={f.adjustment >= 0 ? 'text-green-300' : 'text-rose-300'}>
                                  {f.adjustment >= 0 ? '+' : ''}{f.adjustment.toFixed(2)}
                                </div>
                              </div>
                            ))
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default ScoreBreakdownDrawer4827LX;


