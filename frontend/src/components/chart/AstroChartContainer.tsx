import React, { useEffect, useRef, useId } from 'react';
// @ts-ignore - UMD format library, may export as default or named
import AstroChartLib from '@astrodraw/astrochart';
import type { NatalChart } from '../../types';

interface AstroChartContainerProps {
  data: NatalChart;
  width?: number | string;
  height?: number | string;
}

// 深色宇宙主题配置
const COSMIC_DARK_THEME = {
  // 背景透明，由外层容器控制
  COLOR_BACKGROUND: 'transparent',
  
  // 行星符号 - 亮白色
  POINTS_COLOR: '#E0E0FF',
  POINTS_TEXT_SIZE: 10,
  POINTS_STROKE: 2.2,
  
  // 星座符号 - 亮白色
  SIGNS_COLOR: '#FFFFFF',
  SIGNS_STROKE: 2,
  
  // 圆圈和线条 - 柔和的蓝紫色
  CIRCLE_COLOR: 'rgba(180, 200, 255, 0.6)',
  CIRCLE_STRONG: 2.5,
  LINE_COLOR: 'rgba(180, 200, 255, 0.5)',
  
  // 轴点标签 - 青色
  SYMBOL_AXIS_FONT_COLOR: '#00D4FF',
  SYMBOL_AXIS_STROKE: 2.5,
  
  // 宫位和相位线粗细
  CUSPS_STROKE: 1.8,
  CUSPS_FONT_COLOR: 'rgba(220, 230, 255, 0.9)',
  
  // 星座颜色数组（必须使用 COLORS_SIGNS，库有 bug 不会读取单独的 COLOR_ARIES 等）
  // 顺序：白羊、金牛、双子、巨蟹、狮子、处女、天秤、天蝎、射手、摩羯、水瓶、双鱼
  COLORS_SIGNS: [
    '#EF4444', // 白羊 - 火红
    '#22C55E', // 金牛 - 土绿
    '#EAB308', // 双子 - 风黄
    '#3B82F6', // 巨蟹 - 水蓝
    '#F97316', // 狮子 - 火橙
    '#84CC16', // 处女 - 土青绿
    '#EC4899', // 天秤 - 风粉
    '#DC2626', // 天蝎 - 水深红
    '#A855F7', // 射手 - 火紫
    '#71717A', // 摩羯 - 土灰
    '#06B6D4', // 水瓶 - 风青
    '#8B5CF6', // 双鱼 - 水紫
  ],
  
  // 相位颜色 - 与参考图一致（红=紧张，绿=和谐）
  ASPECTS: {
    conjunction: { degree: 0, orbit: 10, color: '#FFD700' },  // 金色 - 合相
    sextile: { degree: 60, orbit: 6, color: '#22C55E' },      // 绿色 - 六分相（和谐）
    square: { degree: 90, orbit: 8, color: '#DC2626' },       // 红色 - 四分相（紧张）
    trine: { degree: 120, orbit: 8, color: '#22C55E' },       // 绿色 - 三分相（和谐）
    opposition: { degree: 180, orbit: 10, color: '#DC2626' }, // 红色 - 对分相（紧张）
  },
  
  // 符号缩放
  SYMBOL_SCALE: 1.15,
  
  // 显示尊贵度
  SHOW_DIGNITIES_TEXT: true,
  
  // 碰撞半径
  COLLISION_RADIUS: 14,
};

export const AstroChartContainer: React.FC<AstroChartContainerProps> = ({ data, width = 600, height = 600 }) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartId = useId().replace(/:/g, '-'); // React useId 生成唯一ID，替换冒号为连字符

  useEffect(() => {
    if (!chartRef.current || !data) return;

    // 清空容器
    chartRef.current.innerHTML = '';

    // 确保容器有唯一的 id
    if (!chartRef.current.id) {
      chartRef.current.id = `astrochart-${chartId}`;
    }

    try {
      // 映射后端数据到 AstroChart 格式
      // AstroChart 需要: { planets: { "Sun": [longitude], ... }, cusps: [...] }
      const planets: Record<string, number[]> = {};
      
      // 行星名称映射表
      const planetNameMap: Record<string, string> = {
        'sun': 'Sun',
        'moon': 'Moon',
        'mercury': 'Mercury',
        'venus': 'Venus',
        'mars': 'Mars',
        'jupiter': 'Jupiter',
        'saturn': 'Saturn',
        'uranus': 'Uranus',
        'neptune': 'Neptune',
        'pluto': 'Pluto',
        'northNode': 'NNode',  // AstroChart 使用 NNode
        'chiron': 'Chiron',
      };

      data.planets.forEach(p => {
        const name = planetNameMap[p.id] || p.id.charAt(0).toUpperCase() + p.id.slice(1);
        planets[name] = [p.longitude];
      });

      // 宫位角度（12个宫位）
      const cusps = data.houses.map(h => h.cusp);

      // 计算实际尺寸
      const chartWidth = typeof width === 'number' ? width : 600;
      const chartHeight = typeof height === 'number' ? height : 600;

      // 创建 AstroChart 实例
      // UMD 格式库，尝试多种方式访问 Chart 类
      // @ts-ignore
      const Chart = (AstroChartLib as any)?.Chart || 
                    (AstroChartLib as any)?.default?.Chart || 
                    (AstroChartLib as any)?.default ||
                    (window as any)?.astrochart?.Chart ||
                    AstroChartLib;
      
      if (!Chart || typeof Chart !== 'function') {
        throw new Error(`无法找到 AstroChart.Chart 类。库对象键: ${JSON.stringify(Object.keys(AstroChartLib || {}))}`);
      }

      // 使用深色宇宙主题
      const chart = new Chart(chartRef.current.id, chartWidth, chartHeight, COSMIC_DARK_THEME);

      // 绘制本命盘 (Radix)
      const radix = chart.radix({
        planets: planets,
        cusps: cusps,
      });
      
      // 绘制相位线（内圈红绿线）
      radix.aspects();
      
    } catch (error) {
      console.error('[AstroChartContainer] 渲染错误:', error);
      // 显示错误信息
      if (chartRef.current) {
        chartRef.current.innerHTML = `<div style="color: #FF6B9D; padding: 20px; background: rgba(0,0,0,0.5); border-radius: 12px; border: 1px solid rgba(255,107,157,0.3);">
          <h3 style="margin: 0 0 10px 0; color: #FF6B9D;">⚠️ 星盘渲染失败</h3>
          <p style="margin: 0; color: rgba(255,255,255,0.8);">${error instanceof Error ? error.message : String(error)}</p>
          <p style="font-size: 12px; margin-top: 10px; color: rgba(255,255,255,0.5);">请查看浏览器控制台获取详细信息</p>
        </div>`;
      }
    }
  }, [data, chartId, width, height]);

  return (
    <div 
      ref={chartRef} 
      style={{ 
        width: typeof width === 'number' ? `${width}px` : width, 
        height: typeof height === 'number' ? `${height}px` : height,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        background: 'radial-gradient(ellipse at center, rgba(30, 30, 60, 0.8) 0%, rgba(15, 15, 35, 0.9) 50%, rgba(10, 10, 25, 1) 100%)',
        borderRadius: '50%',
        boxShadow: '0 0 40px rgba(0, 180, 255, 0.15), 0 0 80px rgba(168, 85, 247, 0.1), inset 0 0 60px rgba(0, 0, 0, 0.3)',
      }} 
    />
  );
};

