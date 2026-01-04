/**
 * Star 占星计算平台 - 主应用组件
 */

import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Tabs, Tab, Spinner, Button, Switch } from '@heroui/react';
import { useAstroData } from './hooks/useAstroData';
import { NatalChartSVG3847AB } from './components/chart/NatalChartSVG3847AB';
import { BirthDataForm2943KL } from './components/input/BirthDataForm2943KL';
import { ScoreCard5612XY, DimensionScoresCard5612XY } from './components/ui/ScoreCard5612XY';
import { DailyForecastCard7821MN } from './components/forecast/DailyForecastCard7821MN';
import { LifeTimeline4529PQ } from './components/timeline/LifeTimeline4529PQ';
import { ProfectionWheel6183RS } from './components/timeline/ProfectionWheel6183RS';
import { InfluenceFactorsPanel8274TU } from './components/factors/InfluenceFactorsPanel8274TU';
import type { PlanetID, BirthData, InfluenceFactor } from './types';
import { PLANET_NAMES, PLANET_SYMBOLS, PLANET_COLORS, formatDegree } from './utils/astro';

// 模拟影响因子数据（后续从 API 获取）
const MOCK_INFLUENCE_FACTORS: InfluenceFactor[] = [
  { id: '1', type: 'dignity', name: '太阳入庙狮子', value: 3, weight: 1, adjustment: 3, description: '太阳在狮子座获得入庙尊贵', isPositive: true },
  { id: '2', type: 'dignity', name: '金星入旺双鱼', value: 2, weight: 1, adjustment: 2, description: '金星在双鱼座获得旺相尊贵', isPositive: true },
  { id: '3', type: 'retrograde', name: '水星逆行', value: -2, weight: 1, adjustment: -2, description: '水星逆行期间沟通需谨慎', isPositive: false },
  { id: '4', type: 'aspectPhase', name: '木星三分太阳', value: 1.5, weight: 0.8, adjustment: 1.2, description: '木星与太阳形成和谐相位', isPositive: true },
  { id: '5', type: 'aspectPhase', name: '土星四分月亮', value: -1.2, weight: 0.8, adjustment: -0.96, description: '土星与月亮形成紧张相位', isPositive: false },
  { id: '6', type: 'lunarPhase', name: '月亮上弦', value: 0.5, weight: 0.7, adjustment: 0.35, description: '月相处于上弦阶段，适合行动', isPositive: true },
  { id: '7', type: 'profectionLord', name: '年主星木星', value: 1.0, weight: 1, adjustment: 1, description: '今年由木星主管，带来扩张机遇', isPositive: true },
];

function App() {
  const {
    birthData,
    natalChart,
    dailyForecast,
    weeklyForecast,
    lifeTrend,
    profection,
    profectionMap,
    currentAge,
    loading,
    error,
    isReady,
    setBirthData,
    refreshWeekly,
    loadLifeTrend,
    loadProfectionMap,
    clearError,
  } = useAstroData();

  const [selectedTab, setSelectedTab] = useState('chart');
  const [highlightedPlanet, setHighlightedPlanet] = useState<PlanetID | null>(null);
  const [expandedForecast, setExpandedForecast] = useState<string | null>(null);
  const [showFactorEditor, setShowFactorEditor] = useState(false);

  // 加载周预测（当有出生数据时）
  useEffect(() => {
    if (isReady && !weeklyForecast) {
      refreshWeekly();
    }
  }, [isReady, weeklyForecast, refreshWeekly]);

  // 加载趋势数据（当切换到趋势 tab 时）
  useEffect(() => {
    if (isReady && selectedTab === 'trend') {
      if (!lifeTrend) {
        loadLifeTrend(0, 80);
      }
      if (!profectionMap) {
        loadProfectionMap(0, 80);
      }
    }
  }, [isReady, selectedTab, lifeTrend, profectionMap, loadLifeTrend, loadProfectionMap]);

  // 处理出生数据提交
  const handleBirthDataSubmit = async (data: BirthData) => {
    await setBirthData(data);
  };

  return (
    <div className="min-h-screen p-4 md:p-8">
      {/* 标题 */}
      <motion.header
        className="text-center mb-8"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-4xl md:text-5xl font-display font-bold mb-2">
          <span className="bg-gradient-to-r from-[#00D4FF] via-[#A855F7] to-[#FF6B9D] bg-clip-text text-transparent">
            ✦ Star
          </span>
        </h1>
        <p className="text-white/60 text-lg">占星计算验证平台</p>
      </motion.header>

      {/* 错误提示 */}
      <AnimatePresence>
        {error && (
          <motion.div
            className="max-w-2xl mx-auto mb-6 p-4 bg-red-500/20 border border-red-500/30 rounded-lg flex items-center justify-between"
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
          >
            <span className="text-red-300">❌ {error}</span>
            <Button size="sm" variant="light" onPress={clearError}>
              关闭
            </Button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* 主内容区 */}
      {!isReady ? (
        // 未输入出生数据时显示表单
        <motion.div
          className="max-w-md mx-auto"
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
        >
          <BirthDataForm2943KL
            onSubmit={handleBirthDataSubmit}
            loading={loading}
          />
        </motion.div>
      ) : (
        // 已有星盘数据时显示完整界面
        <div className="max-w-7xl mx-auto">
          {/* Tab 导航 */}
          <Tabs
            selectedKey={selectedTab}
            onSelectionChange={(key) => setSelectedTab(key as string)}
            classNames={{
              tabList: "bg-white/5 p-1 rounded-xl",
              cursor: "bg-white/10",
              tab: "text-white/60 data-[selected=true]:text-white",
            }}
            className="mb-6"
          >
            <Tab key="chart" title="🌟 星盘" />
            <Tab key="forecast" title="📅 预测" />
            <Tab key="trend" title="📈 趋势" />
            <Tab key="factors" title="📊 因子" />
            <Tab key="settings" title="⚙️ 设置" />
          </Tabs>

          {/* 加载状态 */}
          {loading && (
            <div className="flex items-center justify-center py-12">
              <Spinner size="lg" color="primary" />
              <span className="ml-3 text-white/60">计算中...</span>
            </div>
          )}

          {/* Tab 内容 */}
          <AnimatePresence mode="wait">
            {/* ==================== 星盘 Tab ==================== */}
            {selectedTab === 'chart' && natalChart && (
              <motion.div
                key="chart"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="grid lg:grid-cols-2 gap-6"
              >
                {/* 星盘 SVG */}
                <div className="glass-card p-6 flex justify-center">
                  <NatalChartSVG3847AB
                    chart={natalChart}
                    size={Math.min(400, window.innerWidth - 80)}
                    showAspects={true}
                    showHouses={true}
                    highlightPlanet={highlightedPlanet}
                    onPlanetClick={setHighlightedPlanet}
                  />
                </div>

                {/* 星盘详情 */}
                <div className="space-y-4">
                  {/* 基本信息卡片 */}
                  <div className="glass-card p-4">
                    <h3 className="text-lg font-medium text-white mb-3">📍 基本信息</h3>
                    <div className="grid grid-cols-2 gap-3 text-sm">
                      <div>
                        <span className="text-white/60">上升点：</span>
                        <span className="text-[#00D4FF]">{formatDegree(natalChart.ascendant)}</span>
                      </div>
                      <div>
                        <span className="text-white/60">天顶：</span>
                        <span className="text-[#FF6B9D]">{formatDegree(natalChart.midheaven)}</span>
                      </div>
                      <div>
                        <span className="text-white/60">主导行星：</span>
                        <span>{natalChart.dominantPlanets.map(p => PLANET_SYMBOLS[p]).join(' ')}</span>
                      </div>
                      <div>
                        <span className="text-white/60">命主星：</span>
                        <span style={{ color: PLANET_COLORS[natalChart.chartRuler] }}>
                          {PLANET_SYMBOLS[natalChart.chartRuler]} {PLANET_NAMES[natalChart.chartRuler]}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* 行星列表 */}
                  <div className="glass-card p-4">
                    <h3 className="text-lg font-medium text-white mb-3">🪐 行星位置</h3>
                    <div className="grid grid-cols-2 gap-2 text-sm max-h-64 overflow-y-auto">
                      {natalChart.planets.map(planet => (
                        <motion.div
                          key={planet.id}
                          className={`flex items-center gap-2 p-2 rounded-lg cursor-pointer transition-colors ${
                            highlightedPlanet === planet.id
                              ? 'bg-white/10'
                              : 'hover:bg-white/5'
                          }`}
                          onClick={() => setHighlightedPlanet(
                            highlightedPlanet === planet.id ? null : planet.id
                          )}
                          whileHover={{ scale: 1.02 }}
                        >
                          <span
                            className="text-lg"
                            style={{ color: PLANET_COLORS[planet.id] }}
                          >
                            {PLANET_SYMBOLS[planet.id]}
                          </span>
                          <div className="flex-1">
                            <div className="text-white/80">{PLANET_NAMES[planet.id]}</div>
                            <div className="text-white/40 text-xs">
                              {planet.signName} {Math.floor(planet.signDegree)}°{Math.floor((planet.signDegree % 1) * 60)}'
                              {planet.retrograde && <span className="text-red-400 ml-1">℞</span>}
                            </div>
                          </div>
                          <div className="text-white/30 text-xs">
                            {planet.house}宫
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  </div>

                  {/* 年限法信息 */}
                  {profection && (
                    <div className="glass-card p-4">
                      <h3 className="text-lg font-medium text-white mb-3">🔮 年限法</h3>
                      <div className="text-sm space-y-2">
                        <div className="flex justify-between">
                          <span className="text-white/60">当前年龄：</span>
                          <span className="text-white">{profection.age}岁</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-white/60">激活宫位：</span>
                          <span className="text-[#00D4FF]">第{profection.house}宫 ({profection.houseName})</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-white/60">年主星：</span>
                          <span className="text-[#ffd700]">
                            {profection.lordSymbol} {profection.lordName}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-white/60">主题：</span>
                          <span className="text-white/80">{profection.houseTheme}</span>
                        </div>
                        {profection.houseKeywords && profection.houseKeywords.length > 0 && (
                          <div className="mt-2 pt-2 border-t border-white/10">
                            <span className="text-white/60">关键词：</span>
                            <div className="flex flex-wrap gap-1 mt-1">
                              {profection.houseKeywords.map((keyword, i) => (
                                <span
                                  key={i}
                                  className="px-2 py-0.5 bg-white/5 rounded text-xs text-white/80"
                                >
                                  {keyword}
                                </span>
                              ))}
                            </div>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </motion.div>
            )}

            {/* ==================== 预测 Tab ==================== */}
            {selectedTab === 'forecast' && (
              <motion.div
                key="forecast"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* 今日预测 */}
                {dailyForecast && (
                  <div>
                    <h3 className="text-xl font-medium text-white mb-4">☀️ 今日预测</h3>
                    <div className="grid md:grid-cols-3 gap-4 mb-4">
                      <ScoreCard5612XY
                        title="综合运势"
                        score={dailyForecast.overallScore}
                        size="lg"
                      />
                      <div className="md:col-span-2">
                        <DimensionScoresCard5612XY
                          scores={dailyForecast.dimensions || {
                            career: 50,
                            relationship: 50,
                            health: 50,
                            finance: 50,
                            spiritual: 50,
                          }}
                          layout="horizontal"
                        />
                      </div>
                    </div>
                    <DailyForecastCard7821MN
                      forecast={dailyForecast}
                      isToday={true}
                      isExpanded={expandedForecast === dailyForecast.date}
                      onClick={() => setExpandedForecast(
                        expandedForecast === dailyForecast.date ? null : dailyForecast.date
                      )}
                    />
                  </div>
                )}

                {/* 本周预测 */}
                {weeklyForecast && (
                  <div>
                    <h3 className="text-xl font-medium text-white mb-4">📆 本周预测</h3>
                    <div className="glass-card p-4 mb-4">
                      <p className="text-white/80">{weeklyForecast.overallTheme}</p>
                      <div className="flex flex-wrap gap-4 mt-3 text-sm">
                        <div>
                          <span className="text-white/60">周综合分：</span>
                          <span className="text-cyan-400">{Math.round(weeklyForecast.overallScore)}</span>
                        </div>
                        {weeklyForecast.bestDaysFor?.relationship?.length > 0 && (
                          <div>
                            <span className="text-white/60">最佳关系日：</span>
                            <span className="text-green-400">
                              {weeklyForecast.bestDaysFor.relationship.slice(0, 2).join(', ')}
                            </span>
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {weeklyForecast.dailySummaries?.map((summary, index) => (
                        <motion.div
                          key={summary.date}
                          initial={{ opacity: 0, y: 10 }}
                          animate={{ opacity: 1, y: 0 }}
                          transition={{ delay: index * 0.05 }}
                          className="glass-card p-4"
                        >
                          <div className="flex justify-between items-center mb-2">
                            <span className="text-white/60 text-sm">{summary.dayOfWeek}</span>
                            <span className="text-cyan-400 font-bold">{Math.round(summary.overallScore)}</span>
                          </div>
                          <div className="text-white text-sm">{summary.date}</div>
                          <div className="text-white/60 text-xs mt-1">{summary.keyTheme}</div>
                        </motion.div>
                      ))}
                    </div>
                  </div>
                )}
              </motion.div>
            )}

            {/* ==================== 趋势 Tab ==================== */}
            {selectedTab === 'trend' && (
              <motion.div
                key="trend"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* 生命趋势图 */}
                {lifeTrend ? (
                  <div>
                    <h3 className="text-xl font-medium text-white mb-4">📈 生命趋势 (0-80岁)</h3>
                    <LifeTimeline4529PQ
                      data={lifeTrend}
                      currentAge={currentAge}
                      height={280}
                      showDimensions={true}
                      onPointClick={(point) => {
                        console.log('点击年龄:', point.age, point);
                      }}
                    />
                  </div>
                ) : (
                  <div className="glass-card p-6 flex items-center justify-center">
                    <Spinner size="lg" />
                    <span className="ml-3 text-white/60">加载生命趋势...</span>
                  </div>
                )}

                {/* 年限法轮盘 */}
                {profectionMap ? (
                  <ProfectionWheel6183RS
                    profections={profectionMap.profections}
                    currentAge={currentAge}
                    size={350}
                    onAgeClick={(age) => {
                      console.log('点击年龄:', age);
                    }}
                  />
                ) : (
                  <div className="glass-card p-6 flex items-center justify-center">
                    <Spinner size="lg" />
                    <span className="ml-3 text-white/60">加载年限法...</span>
                  </div>
                )}

                {/* 重大行运提示 */}
                <div className="glass-card p-4">
                  <h3 className="text-lg font-medium text-white mb-3">🌟 重大行运节点</h3>
                  <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-3 text-sm">
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#A855F7] font-medium">♄ 土星回归</div>
                      <div className="text-white/60">29-30 岁 / 58-60 岁</div>
                      <div className="text-white/40 text-xs mt-1">人生结构升级</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#00D4FF] font-medium">♅ 天王星对冲</div>
                      <div className="text-white/60">40-42 岁</div>
                      <div className="text-white/40 text-xs mt-1">中年觉醒</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#ffd700] font-medium">♃ 木星回归</div>
                      <div className="text-white/60">12 / 24 / 36 / 48 岁</div>
                      <div className="text-white/40 text-xs mt-1">扩张机遇</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#FF6B9D] font-medium">☊ 北交点回归</div>
                      <div className="text-white/60">18-19 / 37-38 岁</div>
                      <div className="text-white/40 text-xs mt-1">命运节点</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#ff8c00] font-medium">⚷ 凯龙回归</div>
                      <div className="text-white/60">50-51 岁</div>
                      <div className="text-white/40 text-xs mt-1">伤痛治愈</div>
                    </div>
                    <div className="bg-white/5 rounded-lg p-3">
                      <div className="text-[#4169e1] font-medium">♆ 海王星四分</div>
                      <div className="text-white/60">41 岁</div>
                      <div className="text-white/40 text-xs mt-1">灵性转化</div>
                    </div>
                  </div>
                </div>
              </motion.div>
            )}

            {/* ==================== 因子 Tab ==================== */}
            {selectedTab === 'factors' && (
              <motion.div
                key="factors"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* 编辑模式开关 */}
                <div className="flex items-center justify-between">
                  <h3 className="text-xl font-medium text-white">📊 影响因子分析</h3>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-white/60">编辑模式</span>
                    <Switch
                      isSelected={showFactorEditor}
                      onValueChange={setShowFactorEditor}
                      size="sm"
                    />
                  </div>
                </div>

                <InfluenceFactorsPanel8274TU
                  factors={MOCK_INFLUENCE_FACTORS}
                  editable={showFactorEditor}
                  onWeightChange={(name, weight) => {
                    console.log('权重变更:', name, weight);
                  }}
                />

                {/* 因子说明 */}
                <div className="glass-card p-4">
                  <h4 className="text-lg font-medium text-white mb-3">📖 因子权重说明</h4>
                  <div className="grid md:grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-white/80 font-medium mb-2">尊贵度 (Dignity)</div>
                      <ul className="text-white/60 space-y-1 list-disc list-inside">
                        <li>入庙 (Domicile): +3</li>
                        <li>旺相 (Exaltation): +2</li>
                        <li>落陷 (Detriment): -2</li>
                        <li>失势 (Fall): -3</li>
                      </ul>
                    </div>
                    <div>
                      <div className="text-white/80 font-medium mb-2">其他因子</div>
                      <ul className="text-white/60 space-y-1 list-disc list-inside">
                        <li>逆行: -2</li>
                        <li>相位阶段: ×0.8</li>
                        <li>外行星放大: ×1.2</li>
                        <li>年主星加成: +1.0</li>
                      </ul>
                    </div>
                  </div>
                </div>
              </motion.div>
            )}

            {/* ==================== 设置 Tab ==================== */}
            {selectedTab === 'settings' && (
              <motion.div
                key="settings"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="max-w-2xl space-y-6"
              >
                <div className="glass-card p-6">
                  <h3 className="text-xl font-medium text-white mb-4">👤 当前出生数据</h3>
                  {birthData && (
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <span className="text-white/60">出生日期：</span>
                        <span className="text-white">
                          {birthData.year}年{birthData.month}月{birthData.day}日
                        </span>
                      </div>
                      <div>
                        <span className="text-white/60">出生时间：</span>
                        <span className="text-white">
                          {String(birthData.hour).padStart(2, '0')}:{String(birthData.minute).padStart(2, '0')}
                        </span>
                      </div>
                      <div>
                        <span className="text-white/60">出生地点：</span>
                        <span className="text-white">
                          {birthData.latitude.toFixed(4)}°, {birthData.longitude.toFixed(4)}°
                        </span>
                      </div>
                      <div>
                        <span className="text-white/60">时区：</span>
                        <span className="text-white">UTC{birthData.timezone >= 0 ? '+' : ''}{birthData.timezone}</span>
                      </div>
                    </div>
                  )}
                  <Button
                    className="mt-4"
                    variant="flat"
                    onPress={() => window.location.reload()}
                  >
                    🔄 重新输入
                  </Button>
                </div>

                <div className="glass-card p-6">
                  <h3 className="text-xl font-medium text-white mb-4">🔧 系统信息</h3>
                  <div className="text-sm space-y-2 text-white/60">
                    <p>版本: 1.0.0</p>
                    <p>算法: VSOP87 简化模型 / Placidus 分宫制</p>
                    <p>精度: 行星经度 &lt;1° / 太阳 &lt;0.1°</p>
                    <p>数据来源: 内置天文算法计算</p>
                  </div>
                </div>

                <div className="glass-card p-6">
                  <h3 className="text-xl font-medium text-white mb-4">⚠️ 免责声明</h3>
                  <p className="text-sm text-white/60">
                    本平台基于天文算法进行占星学计算，所有数据均有理论支撑，
                    仅供研究和学习参考。预测结果不构成任何决策建议，
                    请理性看待占星学分析结果。
                  </p>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      )}

      {/* 页脚 */}
      <footer className="text-center text-white/30 text-sm mt-12">
        <p>Star 占星计算验证平台 v1.0.0</p>
        <p className="mt-1">数据基于天文算法计算，仅供研究参考</p>
      </footer>
    </div>
  );
}

export default App;
