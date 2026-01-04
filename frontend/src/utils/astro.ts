/**
 * 占星学工具函数
 */

import type { PlanetID, ZodiacID, AspectType, Dimension } from '../types';

// ==================== 行星配置 ====================

export const PLANET_SYMBOLS: Record<PlanetID, string> = {
  sun: '☉',
  moon: '☽',
  mercury: '☿',
  venus: '♀',
  mars: '♂',
  jupiter: '♃',
  saturn: '♄',
  uranus: '♅',
  neptune: '♆',
  pluto: '♇',
  northNode: '☊',
  chiron: '⚷',
};

export const PLANET_NAMES: Record<PlanetID, string> = {
  sun: '太阳',
  moon: '月亮',
  mercury: '水星',
  venus: '金星',
  mars: '火星',
  jupiter: '木星',
  saturn: '土星',
  uranus: '天王星',
  neptune: '海王星',
  pluto: '冥王星',
  northNode: '北交点',
  chiron: '凯龙',
};

export const PLANET_COLORS: Record<PlanetID, string> = {
  sun: '#ffd700',
  moon: '#c0c0c0',
  mercury: '#b5651d',
  venus: '#ff69b4',
  mars: '#dc143c',
  jupiter: '#daa520',
  saturn: '#8b7355',
  uranus: '#40e0d0',
  neptune: '#4169e1',
  pluto: '#800080',
  northNode: '#9932cc',
  chiron: '#ff8c00',
};

// ==================== 星座配置 ====================

export const ZODIAC_SYMBOLS: Record<ZodiacID, string> = {
  aries: '♈',
  taurus: '♉',
  gemini: '♊',
  cancer: '♋',
  leo: '♌',
  virgo: '♍',
  libra: '♎',
  scorpio: '♏',
  sagittarius: '♐',
  capricorn: '♑',
  aquarius: '♒',
  pisces: '♓',
};

export const ZODIAC_NAMES: Record<ZodiacID, string> = {
  aries: '白羊座',
  taurus: '金牛座',
  gemini: '双子座',
  cancer: '巨蟹座',
  leo: '狮子座',
  virgo: '处女座',
  libra: '天秤座',
  scorpio: '天蝎座',
  sagittarius: '射手座',
  capricorn: '摩羯座',
  aquarius: '水瓶座',
  pisces: '双鱼座',
};

export const ZODIAC_COLORS: Record<ZodiacID, string> = {
  aries: '#ef4444',
  taurus: '#22c55e',
  gemini: '#eab308',
  cancer: '#3b82f6',
  leo: '#f97316',
  virgo: '#84cc16',
  libra: '#ec4899',
  scorpio: '#dc2626',
  sagittarius: '#a855f7',
  capricorn: '#71717a',
  aquarius: '#06b6d4',
  pisces: '#8b5cf6',
};

// ==================== 相位配置 ====================

export const ASPECT_SYMBOLS: Record<AspectType, string> = {
  conjunction: '☌',
  sextile: '⚹',
  square: '□',
  trine: '△',
  opposition: '☍',
};

export const ASPECT_NAMES: Record<AspectType, string> = {
  conjunction: '合相',
  sextile: '六分相',
  square: '四分相',
  trine: '三分相',
  opposition: '对分相',
};

export const ASPECT_COLORS: Record<AspectType, string> = {
  conjunction: '#ffd700',
  sextile: '#3b82f6',
  square: '#ef4444',
  trine: '#22c55e',
  opposition: '#f97316',
};

export const ASPECT_ANGLES: Record<AspectType, number> = {
  conjunction: 0,
  sextile: 60,
  square: 90,
  trine: 120,
  opposition: 180,
};

// ==================== 维度配置 ====================

export const DIMENSION_NAMES: Record<Dimension, string> = {
  career: '事业',
  relationship: '关系',
  health: '健康',
  finance: '财务',
  spiritual: '灵性',
};

export const DIMENSION_COLORS: Record<Dimension, string> = {
  career: '#3b82f6',
  relationship: '#ec4899',
  health: '#22c55e',
  finance: '#f59e0b',
  spiritual: '#a855f7',
};

export const DIMENSION_ICONS: Record<Dimension, string> = {
  career: '💼',
  relationship: '❤️',
  health: '🌿',
  finance: '💰',
  spiritual: '✨',
};

// ==================== 宫位配置 ====================

export const HOUSE_NAMES: string[] = [
  '命宫',     // 1
  '财帛宫',   // 2
  '兄弟宫',   // 3
  '田宅宫',   // 4
  '子女宫',   // 5
  '奴仆宫',   // 6
  '夫妻宫',   // 7
  '疾厄宫',   // 8
  '迁移宫',   // 9
  '官禄宫',   // 10
  '福德宫',   // 11
  '玄秘宫',   // 12
];

export const HOUSE_THEMES: string[] = [
  '自我身份',     // 1
  '资源价值',     // 2
  '沟通学习',     // 3
  '家庭根基',     // 4
  '创造表达',     // 5
  '服务健康',     // 6
  '关系合作',     // 7
  '转化共享',     // 8
  '探索智慧',     // 9
  '事业成就',     // 10
  '愿景社群',     // 11
  '内省超越',     // 12
];

// ==================== 工具函数 ====================

/**
 * 格式化经度为度分秒格式
 */
export function formatDegree(longitude: number): string {
  const sign = Math.floor(longitude / 30);
  const degree = Math.floor(longitude % 30);
  const minuteValue = Math.floor((longitude % 1) * 60);
  
  const signs: ZodiacID[] = [
    'aries', 'taurus', 'gemini', 'cancer', 'leo', 'virgo',
    'libra', 'scorpio', 'sagittarius', 'capricorn', 'aquarius', 'pisces'
  ];
  
  return `${degree}°${minuteValue}' ${ZODIAC_SYMBOLS[signs[sign]]}`;
}

/**
 * 获取经度对应的星座
 */
export function getSignFromLongitude(longitude: number): ZodiacID {
  const signs: ZodiacID[] = [
    'aries', 'taurus', 'gemini', 'cancer', 'leo', 'virgo',
    'libra', 'scorpio', 'sagittarius', 'capricorn', 'aquarius', 'pisces'
  ];
  return signs[Math.floor(longitude / 30) % 12];
}

/**
 * 计算两点之间的相位角度
 */
export function calculateAspectAngle(long1: number, long2: number): number {
  let diff = Math.abs(long1 - long2);
  if (diff > 180) diff = 360 - diff;
  return diff;
}

/**
 * 分数等级描述
 */
export function getScoreLevel(score: number): {
  level: 'excellent' | 'good' | 'neutral' | 'challenging' | 'difficult';
  label: string;
  color: string;
} {
  if (score >= 80) return { level: 'excellent', label: '极佳', color: '#22c55e' };
  if (score >= 65) return { level: 'good', label: '良好', color: '#84cc16' };
  if (score >= 45) return { level: 'neutral', label: '平稳', color: '#eab308' };
  if (score >= 30) return { level: 'challenging', label: '挑战', color: '#f97316' };
  return { level: 'difficult', label: '困难', color: '#ef4444' };
}

/**
 * 格式化日期显示
 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

/**
 * 格式化时间显示
 */
export function formatDateTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * 将极坐标转换为笛卡尔坐标（用于 SVG 绘制）
 */
export function polarToCartesian(
  centerX: number,
  centerY: number,
  radius: number,
  angleInDegrees: number
): { x: number; y: number } {
  const angleInRadians = ((angleInDegrees - 90) * Math.PI) / 180;
  return {
    x: centerX + radius * Math.cos(angleInRadians),
    y: centerY + radius * Math.sin(angleInRadians),
  };
}

/**
 * 生成 SVG 圆弧路径
 */
export function describeArc(
  x: number,
  y: number,
  radius: number,
  startAngle: number,
  endAngle: number
): string {
  const start = polarToCartesian(x, y, radius, endAngle);
  const end = polarToCartesian(x, y, radius, startAngle);
  const largeArcFlag = endAngle - startAngle <= 180 ? '0' : '1';
  
  return [
    'M', start.x, start.y,
    'A', radius, radius, 0, largeArcFlag, 0, end.x, end.y,
  ].join(' ');
}

