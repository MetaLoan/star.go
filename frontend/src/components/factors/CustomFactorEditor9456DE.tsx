/**
 * CustomFactorEditor9456DE - 自定义因子编辑器
 * 支持添加、编辑、删除自定义影响因子
 * 格式: Operation=(value*dimension,duration,startTime)
 * 例如: AddScore=(2*healthScore,2.5,202517301212)
 */

import { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

// 维度配置
const DIMENSIONS = [
  { key: 'career', label: '事业', icon: '💼', color: '#FF6B6B' },
  { key: 'relationship', label: '关系', icon: '💕', color: '#FF9F43' },
  { key: 'health', label: '健康', icon: '💪', color: '#4ECDC4' },
  { key: 'finance', label: '财务', icon: '💰', color: '#4FC3F7' },
  { key: 'spiritual', label: '灵性', icon: '🔮', color: '#A855F7' },
] as const;

// 操作类型
const OPERATIONS = [
  { key: 'AddScore', label: '加分', description: '在指定维度上增加分数' },
  { key: 'SubScore', label: '减分', description: '在指定维度上减少分数' },
  { key: 'MulScore', label: '倍增', description: '将指定维度分数乘以倍数' },
  { key: 'SetScore', label: '设定', description: '将指定维度设置为固定值' },
] as const;

type DimensionKey = typeof DIMENSIONS[number]['key'];
type OperationKey = typeof OPERATIONS[number]['key'];

interface CustomFactor {
  id: string;
  operation: OperationKey;
  value: number;
  dimension: DimensionKey;
  duration: number; // 小时
  startTime: string; // 格式: YYYYMMDDHHmm
  name?: string;
}

interface CustomFactorEditorProps {
  factors: CustomFactor[];
  onAdd: (factor: Omit<CustomFactor, 'id'>) => void;
  onRemove: (id: string) => void;
  // onUpdate?: (id: string, factor: Partial<CustomFactor>) => void; // 预留
  className?: string;
}

// 时间格式化辅助函数
function formatTimeToInput(time: string): string {
  // 从 202517301212 转换为 datetime-local 格式
  if (time.length !== 12) return '';
  const year = time.slice(0, 4);
  const month = time.slice(4, 6);
  const day = time.slice(6, 8);
  const hour = time.slice(8, 10);
  const minute = time.slice(10, 12);
  return `${year}-${month}-${day}T${hour}:${minute}`;
}

function formatInputToTime(input: string): string {
  // 从 datetime-local 格式转换为 202517301212
  return input.replace(/[-T:]/g, '');
}

export function CustomFactorEditor9456DE({
  factors,
  onAdd,
  onRemove,
  className = '',
}: CustomFactorEditorProps) {
  const [isAdding, setIsAdding] = useState(false);
  // const [editingId, setEditingId] = useState<string | null>(null); // 预留编辑功能

  // 新因子表单状态
  const [newFactor, setNewFactor] = useState<Omit<CustomFactor, 'id'>>({
    operation: 'AddScore',
    value: 1,
    dimension: 'career',
    duration: 1,
    startTime: formatInputToTime(new Date().toISOString().slice(0, 16)),
    name: '',
  });

  // 处理添加
  const handleAdd = useCallback(() => {
    onAdd(newFactor);
    setNewFactor({
      operation: 'AddScore',
      value: 1,
      dimension: 'career',
      duration: 1,
      startTime: formatInputToTime(new Date().toISOString().slice(0, 16)),
      name: '',
    });
    setIsAdding(false);
  }, [newFactor, onAdd]);

  // 生成因子定义字符串
  const getFactorDefinition = (factor: Omit<CustomFactor, 'id'>) => {
    return `${factor.operation}=(${factor.value}*${factor.dimension}Score,${factor.duration},${factor.startTime})`;
  };

  // 解析因子定义字符串 - 预留功能
  // const parseFactorDefinition = (def: string): Omit<CustomFactor, 'id'> | null => {
  //   const match = def.match(/^(\w+)=\(([\d.-]+)\*(\w+)Score,([\d.]+),(\d+)\)$/);
  //   if (!match) return null;
  //   const [, operation, value, dimension, duration, startTime] = match;
  //   return { operation, value: parseFloat(value), dimension, duration: parseFloat(duration), startTime };
  // };

  return (
    <div className={`${className}`}>
      {/* 标题 */}
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-lg font-medium">自定义因子</h3>
          <p className="text-xs text-celestial-silver/50">添加临时影响因子进行测试</p>
        </div>
        <motion.button
          onClick={() => setIsAdding(!isAdding)}
          className="px-3 py-1.5 rounded-lg bg-cosmic-nova text-white text-sm font-medium"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          {isAdding ? '取消' : '+ 添加'}
        </motion.button>
      </div>

      {/* 添加表单 */}
      <AnimatePresence>
        {isAdding && (
          <motion.div
            className="glass-card p-4 mb-4"
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
          >
            <div className="grid grid-cols-2 gap-4">
              {/* 名称 */}
              <div className="col-span-2">
                <label className="block text-xs text-celestial-silver/60 mb-1">名称（可选）</label>
                <input
                  type="text"
                  value={newFactor.name || ''}
                  onChange={(e) => setNewFactor({ ...newFactor, name: e.target.value })}
                  placeholder="如：测试增益"
                  className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm focus:border-cosmic-nova/50 focus:outline-none"
                />
              </div>

              {/* 操作类型 */}
              <div>
                <label className="block text-xs text-celestial-silver/60 mb-1">操作</label>
                <select
                  value={newFactor.operation}
                  onChange={(e) => setNewFactor({ ...newFactor, operation: e.target.value as OperationKey })}
                  className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm focus:border-cosmic-nova/50 focus:outline-none appearance-none cursor-pointer"
                >
                  {OPERATIONS.map((op) => (
                    <option key={op.key} value={op.key} className="bg-cosmic-void">
                      {op.label}
                    </option>
                  ))}
                </select>
              </div>

              {/* 数值 */}
              <div>
                <label className="block text-xs text-celestial-silver/60 mb-1">数值</label>
                <input
                  type="number"
                  value={newFactor.value}
                  onChange={(e) => setNewFactor({ ...newFactor, value: parseFloat(e.target.value) || 0 })}
                  step="0.1"
                  className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm focus:border-cosmic-nova/50 focus:outline-none"
                />
              </div>

              {/* 维度 */}
              <div>
                <label className="block text-xs text-celestial-silver/60 mb-1">维度</label>
                <div className="flex gap-1">
                  {DIMENSIONS.map((dim) => (
                    <button
                      key={dim.key}
                      onClick={() => setNewFactor({ ...newFactor, dimension: dim.key })}
                      className={`flex-1 px-2 py-2 rounded-lg text-center transition-all ${
                        newFactor.dimension === dim.key
                          ? 'ring-2'
                          : 'bg-white/5 hover:bg-white/10'
                      }`}
                      style={newFactor.dimension === dim.key ? {
                        backgroundColor: `${dim.color}30`,
                        borderColor: dim.color,
                        // @ts-ignore
                        '--tw-ring-color': dim.color,
                      } : undefined}
                      title={dim.label}
                    >
                      <span className="text-lg">{dim.icon}</span>
                    </button>
                  ))}
                </div>
              </div>

              {/* 持续时间 */}
              <div>
                <label className="block text-xs text-celestial-silver/60 mb-1">持续时间（小时）</label>
                <input
                  type="number"
                  value={newFactor.duration}
                  onChange={(e) => setNewFactor({ ...newFactor, duration: parseFloat(e.target.value) || 1 })}
                  min="0.1"
                  step="0.5"
                  className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm focus:border-cosmic-nova/50 focus:outline-none"
                />
              </div>

              {/* 开始时间 */}
              <div className="col-span-2">
                <label className="block text-xs text-celestial-silver/60 mb-1">开始时间</label>
                <input
                  type="datetime-local"
                  value={formatTimeToInput(newFactor.startTime)}
                  onChange={(e) => setNewFactor({ ...newFactor, startTime: formatInputToTime(e.target.value) })}
                  className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm focus:border-cosmic-nova/50 focus:outline-none"
                />
              </div>
            </div>

            {/* 预览 */}
            <div className="mt-4 p-3 rounded-lg bg-black/30 font-mono text-xs">
              <span className="text-celestial-silver/50">定义: </span>
              <span className="text-cosmic-nova">{getFactorDefinition(newFactor)}</span>
            </div>

            {/* 提交按钮 */}
            <div className="mt-4 flex justify-end gap-2">
              <button
                onClick={() => setIsAdding(false)}
                className="px-4 py-2 rounded-lg bg-white/5 text-sm hover:bg-white/10 transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleAdd}
                className="px-4 py-2 rounded-lg bg-cosmic-nova text-white text-sm font-medium hover:bg-cosmic-nova/80 transition-colors"
              >
                添加因子
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* 因子列表 */}
      <div className="space-y-2">
        {factors.length > 0 ? (
          factors.map((factor, index) => {
            const dimInfo = DIMENSIONS.find((d) => d.key === factor.dimension);
            const opInfo = OPERATIONS.find((o) => o.key === factor.operation);

            return (
              <motion.div
                key={factor.id}
                className="glass-card p-3"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ delay: index * 0.05 }}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    {/* 维度图标 */}
                    <div
                      className="w-10 h-10 rounded-full flex items-center justify-center"
                      style={{ backgroundColor: `${dimInfo?.color}30` }}
                    >
                      <span className="text-xl">{dimInfo?.icon}</span>
                    </div>

                    {/* 信息 */}
                    <div>
                      <div className="font-medium">
                        {factor.name || `${opInfo?.label} ${dimInfo?.label}`}
                      </div>
                      <div className="text-xs text-celestial-silver/50">
                        {opInfo?.label} {factor.value > 0 ? '+' : ''}{factor.value} • {factor.duration}小时
                      </div>
                    </div>
                  </div>

                  {/* 操作按钮 */}
                  <div className="flex items-center gap-2">
                    {/* 状态标签 */}
                    {(() => {
                      const now = new Date();
                      const start = new Date(formatTimeToInput(factor.startTime));
                      const end = new Date(start.getTime() + factor.duration * 3600000);
                      
                      if (now < start) {
                        return (
                          <span className="text-xs px-2 py-0.5 rounded-full bg-yellow-500/20 text-yellow-400">
                            未开始
                          </span>
                        );
                      } else if (now > end) {
                        return (
                          <span className="text-xs px-2 py-0.5 rounded-full bg-gray-500/20 text-gray-400">
                            已结束
                          </span>
                        );
                      } else {
                        return (
                          <span className="text-xs px-2 py-0.5 rounded-full bg-green-500/20 text-green-400">
                            进行中
                          </span>
                        );
                      }
                    })()}

                    {/* 删除按钮 */}
                    <button
                      onClick={() => onRemove(factor.id)}
                      className="p-1.5 rounded-lg text-red-400 hover:bg-red-500/20 transition-colors"
                      title="删除"
                    >
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
                      </svg>
                    </button>
                  </div>
                </div>

                {/* 定义字符串 */}
                <div className="mt-2 px-2 py-1 rounded bg-black/30 font-mono text-xs text-celestial-silver/50 overflow-x-auto">
                  {getFactorDefinition(factor)}
                </div>
              </motion.div>
            );
          })
        ) : (
          <div className="text-center py-8 text-celestial-silver/50">
            <div className="text-4xl mb-2">⚡</div>
            <div>暂无自定义因子</div>
            <div className="text-xs mt-1">点击"添加"创建临时影响因子</div>
          </div>
        )}
      </div>
    </div>
  );
}

export default CustomFactorEditor9456DE;

