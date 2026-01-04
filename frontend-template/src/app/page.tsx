'use client';

import { useState } from 'react';
import { useAstroData, type BirthDataInput } from '@/lib/api';

export default function Home() {
  const { chart, daily, loading, error, initialize } = useAstroData();
  const [formData, setFormData] = useState({
    name: '',
    date: '',
    latitude: 39.9042,
    longitude: 116.4074,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const birthData: BirthDataInput = {
      name: formData.name || 'Anonymous',
      date: new Date(formData.date).toISOString(),
      latitude: formData.latitude,
      longitude: formData.longitude,
      timezone: 'Asia/Shanghai',
    };
    
    await initialize(birthData);
  };

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8 bg-gradient-to-r from-purple-400 to-pink-500 bg-clip-text text-transparent">
          Star 占星系统
        </h1>
        
        {!chart ? (
          <div className="glass-card p-8">
            <h2 className="text-xl mb-6 text-center">输入出生信息</h2>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm text-gray-400 mb-1">姓名</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-2"
                  placeholder="你的名字"
                />
              </div>
              
              <div>
                <label className="block text-sm text-gray-400 mb-1">出生日期时间</label>
                <input
                  type="datetime-local"
                  value={formData.date}
                  onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                  className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-2"
                  required
                />
              </div>
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm text-gray-400 mb-1">纬度</label>
                  <input
                    type="number"
                    step="0.0001"
                    value={formData.latitude}
                    onChange={(e) => setFormData({ ...formData, latitude: parseFloat(e.target.value) })}
                    className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-2"
                  />
                </div>
                <div>
                  <label className="block text-sm text-gray-400 mb-1">经度</label>
                  <input
                    type="number"
                    step="0.0001"
                    value={formData.longitude}
                    onChange={(e) => setFormData({ ...formData, longitude: parseFloat(e.target.value) })}
                    className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-2"
                  />
                </div>
              </div>
              
              <button
                type="submit"
                disabled={loading}
                className="w-full bg-gradient-to-r from-purple-500 to-pink-500 text-white py-3 rounded-lg font-semibold hover:opacity-90 disabled:opacity-50"
              >
                {loading ? '计算中...' : '生成星盘'}
              </button>
            </form>
            
            {error && (
              <p className="mt-4 text-red-400 text-center">{error.message}</p>
            )}
          </div>
        ) : (
          <div className="space-y-6">
            {/* 星盘摘要 */}
            <div className="glass-card p-6">
              <h2 className="text-2xl font-bold mb-4">{chart.birthData.name}</h2>
              
              <div className="grid grid-cols-4 gap-4 mb-4">
                {chart.planets?.slice(0, 4).map((planet: any) => (
                  <div key={planet.id} className="bg-white/5 rounded-lg p-3 text-center">
                    <div className="text-2xl mb-1">{planet.symbol}</div>
                    <div className="text-sm text-gray-400">{planet.name}</div>
                    <div className="text-sm">{planet.signName}</div>
                  </div>
                ))}
              </div>
              
              <div className="flex gap-2 flex-wrap">
                {chart.dominantPlanets?.map((id: string) => {
                  const planet = chart.planets?.find((p: any) => p.id === id);
                  return (
                    <span key={id} className="px-3 py-1 bg-purple-500/20 rounded-full text-sm">
                      {planet?.symbol} {planet?.name}
                    </span>
                  );
                })}
              </div>
            </div>
            
            {/* 每日预测 */}
            {daily && (
              <div className="glass-card p-6">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-xl font-bold">今日运势</h3>
                  <span className={`text-3xl font-bold ${
                    daily.overallScore >= 70 ? 'text-green-400' :
                    daily.overallScore >= 50 ? 'text-yellow-400' : 'text-red-400'
                  }`}>
                    {daily.overallScore}分
                  </span>
                </div>
                
                <p className="text-gray-300 mb-4">{daily.theme}</p>
                
                <div className="grid grid-cols-5 gap-2">
                  {Object.entries(daily.dimensions || {}).map(([key, value]) => {
                    const labels: Record<string, string> = {
                      career: '💼 事业',
                      relationship: '💕 关系',
                      health: '🏃 健康',
                      finance: '💰 财务',
                      spiritual: '🧘 灵性',
                    };
                    return (
                      <div key={key} className="bg-white/5 rounded-lg p-2 text-center">
                        <div className="text-xs text-gray-400">{labels[key]}</div>
                        <div className="font-bold">{Math.round(value as number)}</div>
                      </div>
                    );
                  })}
                </div>
                
                {daily.factors && (
                  <div className="mt-4 pt-4 border-t border-white/10">
                    <div className="text-sm text-gray-400">
                      影响因子调整: 
                      <span className={daily.factors.totalAdjustment >= 0 ? 'text-green-400' : 'text-red-400'}>
                        {' '}{daily.factors.totalAdjustment >= 0 ? '+' : ''}{daily.factors.totalAdjustment.toFixed(1)}
                      </span>
                    </div>
                    <p className="text-xs text-gray-500 mt-1">{daily.factors.summary}</p>
                  </div>
                )}
              </div>
            )}
            
            {/* 重新输入按钮 */}
            <button
              onClick={() => window.location.reload()}
              className="w-full bg-white/10 hover:bg-white/20 py-3 rounded-lg transition-colors"
            >
              重新输入
            </button>
          </div>
        )}
        
        <footer className="mt-12 text-center text-gray-500 text-sm">
          <p>Powered by Star API</p>
        </footer>
      </div>
    </main>
  );
}

