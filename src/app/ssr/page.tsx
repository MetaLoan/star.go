/**
 * SSR Demo Page
 * 服务端渲染示例页面
 * 
 * 这个页面演示如何在服务端获取和渲染占星数据
 */

import { Suspense } from 'react';
import { Card, CardBody, Chip, Spinner } from '@heroui/react';
import { calculateChartServer, calculateDailyServer, type BirthDataInput } from '@/lib/api/server';

// 示例出生数据（实际使用时从 URL 参数或数据库获取）
const DEMO_BIRTH_DATA: BirthDataInput = {
  name: 'Demo User',
  date: '1990-06-15T10:30:00Z',
  latitude: 39.9042,
  longitude: 116.4074,
  timezone: 'Asia/Shanghai',
};

// 服务端数据获取
async function getAstroData() {
  try {
    const [chart, daily] = await Promise.all([
      calculateChartServer(DEMO_BIRTH_DATA),
      calculateDailyServer(DEMO_BIRTH_DATA, undefined, true),
    ]);
    
    return { chart, daily, error: null };
  } catch (error) {
    console.error('SSR fetch error:', error);
    return { chart: null, daily: null, error: 'Failed to load data' };
  }
}

// 星盘摘要组件
function ChartSummary({ chart }: { chart: any }) {
  if (!chart) return null;
  
  const sun = chart.planets?.find((p: any) => p.id === 'sun');
  const moon = chart.planets?.find((p: any) => p.id === 'moon');
  
  return (
    <Card className="glass-card">
      <CardBody>
        <h2 className="text-xl font-display text-cosmic-nova mb-4">{chart.birthData.name}</h2>
        
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <div className="text-xs text-white/50 mb-1">太阳星座</div>
            <div className="text-lg">{sun?.symbol} {sun?.signName}</div>
          </div>
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <div className="text-xs text-white/50 mb-1">月亮星座</div>
            <div className="text-lg">{moon?.symbol} {moon?.signName}</div>
          </div>
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <div className="text-xs text-white/50 mb-1">上升</div>
            <div className="text-lg">{Math.floor(chart.ascendant / 30) + 1}宫</div>
          </div>
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <div className="text-xs text-white/50 mb-1">盘主星</div>
            <div className="text-lg">{chart.chartRuler}</div>
          </div>
        </div>
        
        <div className="flex flex-wrap gap-2">
          {chart.dominantPlanets?.slice(0, 3).map((planetId: string) => {
            const planet = chart.planets?.find((p: any) => p.id === planetId);
            return (
              <Chip key={planetId} variant="flat" size="sm">
                {planet?.symbol} {planet?.name}
              </Chip>
            );
          })}
        </div>
      </CardBody>
    </Card>
  );
}

// 每日预测组件
function DailyForecast({ daily }: { daily: any }) {
  if (!daily) return null;
  
  const scoreColor = daily.overallScore >= 70 ? 'success' : 
                     daily.overallScore >= 50 ? 'warning' : 'danger';
  
  return (
    <Card className="glass-card">
      <CardBody>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-display text-cosmic-nova">今日运势</h3>
          <Chip color={scoreColor} size="lg" variant="flat">
            {daily.overallScore}分
          </Chip>
        </div>
        
        <p className="text-white/70 mb-4">{daily.theme}</p>
        
        <div className="grid grid-cols-5 gap-2">
          {Object.entries(daily.dimensions || {}).map(([key, value]) => {
            const labels: Record<string, string> = {
              career: '💼',
              relationship: '💕',
              health: '🏃',
              finance: '💰',
              spiritual: '🧘',
            };
            return (
              <div key={key} className="bg-white/5 rounded-lg p-2 text-center">
                <div className="text-lg mb-1">{labels[key]}</div>
                <div className="text-sm font-bold">{Math.round(value as number)}</div>
              </div>
            );
          })}
        </div>
        
        {daily.factors && (
          <div className="mt-4 pt-4 border-t border-white/10">
            <div className="flex items-center gap-2 mb-2">
              <span className="text-sm text-white/50">影响因子:</span>
              <Chip 
                size="sm" 
                color={daily.factors.totalAdjustment >= 0 ? 'success' : 'warning'}
                variant="flat"
              >
                {daily.factors.totalAdjustment >= 0 ? '+' : ''}{daily.factors.totalAdjustment.toFixed(1)}
              </Chip>
            </div>
            <p className="text-xs text-white/50">{daily.factors.summary}</p>
          </div>
        )}
      </CardBody>
    </Card>
  );
}

// 加载占位符
function LoadingSkeleton() {
  return (
    <div className="flex items-center justify-center min-h-[200px]">
      <Spinner size="lg" color="primary" />
    </div>
  );
}

// 主页面（Server Component）
export default async function SSRPage() {
  const { chart, daily, error } = await getAstroData();
  
  return (
    <main className="min-h-screen relative">
      <div className="starfield" />
      
      <div className="container mx-auto px-4 py-12">
        <header className="text-center mb-12">
          <h1 className="text-3xl font-display text-cosmic-nova mb-2">
            SSR 渲染示例
          </h1>
          <p className="text-white/50">服务端渲染占星数据</p>
          <Chip className="mt-2" variant="flat" color="primary">SSR</Chip>
        </header>
        
        {error ? (
          <Card className="glass-card">
            <CardBody>
              <p className="text-red-400 text-center">{error}</p>
            </CardBody>
          </Card>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Suspense fallback={<LoadingSkeleton />}>
              <ChartSummary chart={chart} />
            </Suspense>
            
            <Suspense fallback={<LoadingSkeleton />}>
              <DailyForecast daily={daily} />
            </Suspense>
          </div>
        )}
        
        <div className="mt-12 text-center">
          <p className="text-xs text-white/30">
            此页面数据在服务端获取并渲染，适用于 SEO 和首屏加载优化
          </p>
        </div>
      </div>
    </main>
  );
}

