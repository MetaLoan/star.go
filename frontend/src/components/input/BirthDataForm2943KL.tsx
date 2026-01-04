/**
 * 出生数据表单组件
 * 组件命名规范：BirthDataForm + 2943 + KL
 */

import { useState, useCallback } from 'react';
import { motion } from 'framer-motion';
import type { BirthData } from '../../types';

interface BirthDataFormProps {
  onSubmit: (data: BirthData) => void;
  initialData?: Partial<BirthData>;
  loading?: boolean;
  className?: string;
}

// 常用城市预设
const CITY_PRESETS = [
  { name: '北京', lat: 39.9042, lng: 116.4074, tz: 8 },
  { name: '上海', lat: 31.2304, lng: 121.4737, tz: 8 },
  { name: '香港', lat: 22.3193, lng: 114.1694, tz: 8 },
  { name: '台北', lat: 25.0330, lng: 121.5654, tz: 8 },
  { name: '东京', lat: 35.6762, lng: 139.6503, tz: 9 },
  { name: '纽约', lat: 40.7128, lng: -74.0060, tz: -5 },
  { name: '洛杉矶', lat: 34.0522, lng: -118.2437, tz: -8 },
  { name: '伦敦', lat: 51.5074, lng: -0.1278, tz: 0 },
];

// 时区选项
const TIMEZONE_OPTIONS = [
  { value: -8, label: 'UTC-8 太平洋' },
  { value: -5, label: 'UTC-5 东部' },
  { value: 0, label: 'UTC+0 格林威治' },
  { value: 8, label: 'UTC+8 北京' },
  { value: 9, label: 'UTC+9 东京' },
];

// 表单输入框样式
const inputClass = `
  w-full px-3 py-2.5 
  bg-white/5 border border-white/10 rounded-lg
  text-white text-sm
  focus:outline-none focus:border-cyan-400/50 focus:ring-1 focus:ring-cyan-400/30
  transition-all duration-200
`;

const labelClass = "block text-xs text-white/50 mb-1.5 font-medium";

export function BirthDataForm2943KL({
  onSubmit,
  initialData,
  loading = false,
  className = '',
}: BirthDataFormProps) {
  const [formData, setFormData] = useState({
    year: initialData?.year || 1990,
    month: initialData?.month || 6,
    day: initialData?.day || 15,
    hour: initialData?.hour || 12,
    minute: initialData?.minute || 0,
    latitude: initialData?.latitude || 39.9042,
    longitude: initialData?.longitude || 116.4074,
    timezone: initialData?.timezone || 8,
  });

  const handleChange = useCallback((field: string, value: number) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  }, []);

  const handleCitySelect = useCallback((city: typeof CITY_PRESETS[0]) => {
    setFormData(prev => ({
      ...prev,
      latitude: city.lat,
      longitude: city.lng,
      timezone: city.tz,
    }));
  }, []);

  const handleSubmit = useCallback((e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      year: formData.year,
      month: formData.month,
      day: formData.day,
      hour: formData.hour,
      minute: formData.minute,
      second: 0,
      latitude: formData.latitude,
      longitude: formData.longitude,
      timezone: formData.timezone,
    });
  }, [formData, onSubmit]);

  // 生成月份天数
  const daysInMonth = new Date(formData.year, formData.month, 0).getDate();

  return (
    <motion.form
      className={`glass-card p-6 ${className}`}
      onSubmit={handleSubmit}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
    >
      {/* 标题 */}
      <h2 className="text-xl font-semibold text-white mb-6 flex items-center gap-2">
        <span className="text-2xl">🌟</span>
        输入出生信息
      </h2>

      {/* 出生日期 */}
      <div className="mb-5">
        <div className="text-xs text-white/50 mb-2 font-medium">出生日期</div>
        <div className="grid grid-cols-3 gap-2">
          <div>
            <select
              value={formData.year}
              onChange={(e) => handleChange('year', Number(e.target.value))}
              className={inputClass}
            >
              {Array.from({ length: 100 }, (_, i) => 2024 - i).map(year => (
                <option key={year} value={year}>{year}年</option>
              ))}
            </select>
          </div>
          <div>
            <select
              value={formData.month}
              onChange={(e) => handleChange('month', Number(e.target.value))}
              className={inputClass}
            >
              {Array.from({ length: 12 }, (_, i) => i + 1).map(month => (
                <option key={month} value={month}>{month}月</option>
              ))}
            </select>
          </div>
          <div>
            <select
              value={formData.day}
              onChange={(e) => handleChange('day', Number(e.target.value))}
              className={inputClass}
            >
              {Array.from({ length: daysInMonth }, (_, i) => i + 1).map(day => (
                <option key={day} value={day}>{day}日</option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* 出生时间 */}
      <div className="mb-5">
        <div className="text-xs text-white/50 mb-2 font-medium">出生时间</div>
        <div className="grid grid-cols-2 gap-2">
          <div>
            <select
              value={formData.hour}
              onChange={(e) => handleChange('hour', Number(e.target.value))}
              className={inputClass}
            >
              {Array.from({ length: 24 }, (_, i) => i).map(hour => (
                <option key={hour} value={hour}>{String(hour).padStart(2, '0')}时</option>
              ))}
            </select>
          </div>
          <div>
            <select
              value={formData.minute}
              onChange={(e) => handleChange('minute', Number(e.target.value))}
              className={inputClass}
            >
              {Array.from({ length: 60 }, (_, i) => i).map(minute => (
                <option key={minute} value={minute}>{String(minute).padStart(2, '0')}分</option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* 快速选择城市 */}
      <div className="mb-5">
        <div className="text-xs text-white/50 mb-2 font-medium">快速选择城市</div>
        <div className="flex flex-wrap gap-1.5">
          {CITY_PRESETS.map(city => (
            <button
              key={city.name}
              type="button"
              onClick={() => handleCitySelect(city)}
              className={`
                px-3 py-1.5 text-xs rounded-md transition-all
                ${formData.latitude === city.lat && formData.longitude === city.lng
                  ? 'bg-cyan-500/30 text-cyan-300 border border-cyan-400/50'
                  : 'bg-white/5 text-white/70 border border-white/10 hover:bg-white/10'
                }
              `}
            >
              {city.name}
            </button>
          ))}
        </div>
      </div>

      {/* 经纬度 */}
      <div className="mb-5">
        <div className="text-xs text-white/50 mb-2 font-medium">出生地点坐标</div>
        <div className="grid grid-cols-2 gap-2">
          <div>
            <input
              type="number"
              value={formData.latitude}
              onChange={(e) => handleChange('latitude', parseFloat(e.target.value) || 0)}
              step="0.0001"
              placeholder="纬度"
              className={inputClass}
            />
            <div className="text-[10px] text-white/30 mt-1">纬度 (-90 ~ 90)</div>
          </div>
          <div>
            <input
              type="number"
              value={formData.longitude}
              onChange={(e) => handleChange('longitude', parseFloat(e.target.value) || 0)}
              step="0.0001"
              placeholder="经度"
              className={inputClass}
            />
            <div className="text-[10px] text-white/30 mt-1">经度 (-180 ~ 180)</div>
          </div>
        </div>
      </div>

      {/* 时区 */}
      <div className="mb-6">
        <label className={labelClass}>时区</label>
        <select
          value={formData.timezone}
          onChange={(e) => handleChange('timezone', Number(e.target.value))}
          className={inputClass}
        >
          {TIMEZONE_OPTIONS.map(tz => (
            <option key={tz.value} value={tz.value}>{tz.label}</option>
          ))}
        </select>
      </div>

      {/* 提交按钮 */}
      <button
        type="submit"
        disabled={loading}
        className={`
          w-full py-3 px-4 rounded-lg font-medium text-white
          bg-gradient-to-r from-cyan-500 to-purple-500
          hover:from-cyan-400 hover:to-purple-400
          active:scale-[0.98] transition-all duration-200
          disabled:opacity-50 disabled:cursor-not-allowed
          shadow-lg shadow-cyan-500/20
        `}
      >
        {loading ? (
          <span className="flex items-center justify-center gap-2">
            <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            计算中...
          </span>
        ) : (
          <span>✨ 生成星盘</span>
        )}
      </button>
    </motion.form>
  );
}

export default BirthDataForm2943KL;
