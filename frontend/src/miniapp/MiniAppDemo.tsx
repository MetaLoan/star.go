import React, { useState, useEffect, useRef, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { MessageRenderer } from './components/MessageRenderer';
import { OracleInteractiveEye } from './components/OracleInteractiveEye';
import { OracleThinking } from './components/OracleThinking';
import { StoneUnified } from './components/StoneUnified';
import { TarotUnified } from './components/TarotUnified';
import { EchoUnified } from './components/EchoUnified';
import { ConnectUnified } from './components/ConnectUnified';
import { FunctionButtons } from './components/FunctionButtons';
import { BottomNav, type NavItem } from './components/BottomNav';
import { DestinyChart } from './components/DestinyChart';
import { DestinyTimeline } from './components/DestinyTimeline';
import { StarRadar } from './components/StarRadar';
import { LifeChart } from './components/LifeChart';
import { InkRevealText } from './components/InkRevealText';
import { PageHeaderStyle } from './components/PageHeaderStyle';
import { InpageTabSwitcher } from './components/InpageTabSwitcher';
import { useAstroData } from '../hooks/useAstroData';
import { detectIntent } from '../lib/miniapp/intent-detection';
import { generateDeepInterpretation, type ConversationState } from '../lib/miniapp/oracle-conversation';
import type { Message, TextMessage, FunctionResultMessage, SystemMessage, TarotResult, EchoResult, EchoAudioMessage, StoneResult } from '../lib/miniapp/oracle-types';

const initialMessages: Message[] = [
  {
    id: 1,
    type: "text",
    content: "I sense your presence. I can see your star chart shows Mars in retrograde. Tell me what weighs on your heart.",
    isUser: false,
    contextHint: "Context-aware",
    timestamp: "Just now",
  },
]

const MiniAppDemo: React.FC = () => {
  const { timeSeries, loadTimeSeries, extendTimeSeries, isReady, birthData, setBirthData, loading } = useAstroData();
  const [activeTab, setActiveTab] = useState<NavItem>("destiny");
  const [granularity, setGranularity] = useState<'hour' | 'day' | 'week' | 'month' | 'year'>('day');
  const [selectedDimension, setSelectedDimension] = useState<string>('overall');
  
  // 初始化演示数据
  useEffect(() => {
    if (!birthData) {
      setBirthData({
        year: 1995,
        month: 10,
        day: 24,
        hour: 10,
        minute: 30,
        latitude: 31.23,
        longitude: 121.47,
        timezone: 8
      });
    }
  }, [birthData, setBirthData]);

  // 加载时间序列数据
  useEffect(() => {
    if (isReady && activeTab === "destiny") {
      const now = new Date();
      const start = new Date(now);
      const end = new Date(now);
      
      if (granularity === 'hour') {
        start.setHours(now.getHours() - 12);
        end.setHours(now.getHours() + 12);
      } else if (granularity === 'day') {
        start.setDate(now.getDate() - 15);
        end.setDate(now.getDate() + 15);
      } else if (granularity === 'week') {
        start.setMonth(now.getMonth() - 3);
        end.setMonth(now.getMonth() + 3);
      } else if (granularity === 'month') {
        start.setFullYear(now.getFullYear() - 1);
        end.setFullYear(now.getFullYear() + 1);
      } else {
        start.setFullYear(now.getFullYear() - 5);
        end.setFullYear(now.getFullYear() + 5);
      }
      
      loadTimeSeries(start.toISOString(), end.toISOString(), granularity);
    }
  }, [isReady, activeTab, granularity, loadTimeSeries]);

  // 处理图表可视范围变化（动态加载数据）
  const handleVisibleRangeChange = useCallback(async (range: any) => {
    if (!isReady || loading) return;

    if (range.needsMoreBefore) {
      // 向前加载更多（例如加载前 30 个点对应的时长）
      const firstPoint = timeSeries?.points[0];
      if (!firstPoint) return;
      const firstTime = new Date(firstPoint.time);
      const newStart = new Date(firstTime);
      
      if (granularity === 'hour') newStart.setHours(firstTime.getHours() - 24);
      else if (granularity === 'day') newStart.setDate(firstTime.getDate() - 30);
      else if (granularity === 'week') newStart.setDate(firstTime.getDate() - 7 * 12);
      else if (granularity === 'month') newStart.setMonth(firstTime.getMonth() - 12);
      else newStart.setFullYear(firstTime.getFullYear() - 5);

      await extendTimeSeries(newStart.toISOString(), firstTime.toISOString(), granularity, 'before');
    } else if (range.needsMoreAfter) {
      // 向后加载更多
      const lastPoint = timeSeries?.points[timeSeries.points.length - 1];
      if (!lastPoint) return;
      const lastTime = new Date(lastPoint.time);
      const newEnd = new Date(lastTime);
      
      if (granularity === 'hour') newEnd.setHours(lastTime.getHours() + 24);
      else if (granularity === 'day') newEnd.setDate(lastTime.getDate() + 30);
      else if (granularity === 'week') newEnd.setDate(lastTime.getDate() + 7 * 12);
      else if (granularity === 'month') newEnd.setMonth(lastTime.getMonth() + 12);
      else newEnd.setFullYear(lastTime.getFullYear() + 5);

      await extendTimeSeries(lastTime.toISOString(), newEnd.toISOString(), granularity, 'after');
    }
  }, [isReady, loading, timeSeries, granularity, extendTimeSeries]);

  const [messages, setMessages] = useState<Message[]>(initialMessages)
  const [input, setInput] = useState("")
  const [eyeDirection, setEyeDirection] = useState<"down" | "down-right" | "down-left" | "center">("center")
  const [isThinking, setIsThinking] = useState(false)
  const [conversationState, setConversationState] = useState<ConversationState>("normal")
  const [userQuestion, setUserQuestion] = useState<string>("")
  const [selectedFunction, setSelectedFunction] = useState<"stone" | "tarot" | "echo" | "connect" | null>(null)
  const [destinyMode, setDestinyMode] = useState<"today" | "life">("today")
  const [cardRevealed, setCardRevealed] = useState(false)
  const [selectedHour, setSelectedHour] = useState<number | null>(null)
  
  const [echoRecommendation] = useState({
    frequency: 432,
    trackName: "舒缓旋律·磁场调整",
    recommendation: "根据你的星盘分析，建议通过432Hz频率进行能量调整。"
  })
  
  const [isEchoGenerating, setIsEchoGenerating] = useState(false)
  const [echoGenerateProgress, setEchoGenerateProgress] = useState(0)
  const scrollContainerRef = useRef<HTMLDivElement>(null)
  const messageIdCounterRef = useRef<number>(1000)

  const generateMessageId = () => {
    messageIdCounterRef.current += 1
    return messageIdCounterRef.current
  }

  const scrollToBottom = () => {
    setTimeout(() => {
      if (scrollContainerRef.current) {
        scrollContainerRef.current.scrollTo({ top: scrollContainerRef.current.scrollHeight, behavior: "smooth" })
      }
    }, 100)
  }

  // --- Functions for Oracle Interaction ---
  const handleFunctionButtonClick = (func: "stone" | "tarot" | "echo" | "connect") => {
    setSelectedFunction(func)
    setConversationState("waiting_for_question")
    setIsThinking(true)
    setTimeout(() => {
      setIsThinking(false)
      const questionMessage: TextMessage = {
        id: generateMessageId(),
        type: "text",
        content: `I see you seek the guidance of ${func}. What weighs on your heart?`,
        isUser: false,
        timestamp: "Just now",
      }
      setMessages((prev) => [...prev, questionMessage])
      scrollToBottom()
    }, 1000)
  }

  const sendMessage = () => {
    if (!input.trim()) return
    const userMessage: TextMessage = { id: generateMessageId(), type: "text", content: input, isUser: true, timestamp: "Just now" }
    setMessages((prev) => [...prev, userMessage])
    setInput("")
    scrollToBottom()
    
    setIsThinking(true)
    setTimeout(() => {
      setIsThinking(false)
      const aiResponse: TextMessage = {
        id: generateMessageId(),
        type: "text",
        content: "The cosmos has heard your inquiry. Your path is unfolding as it should.",
        isUser: false,
        timestamp: "Just now",
      }
      setMessages((prev) => [...prev, aiResponse])
      scrollToBottom()
    }, 2000)
  }

  return (
    <div className="freya-mode min-h-screen bg-white text-black font-sans selection:bg-black/5">
      {/* 顶部镭射装饰条 */}
      <div className="h-[2px] w-full laser-gradient fixed top-0 left-0 z-[60]" />

      <main className="pb-24 pt-12 max-w-screen-sm mx-auto px-6">
        {/* DESTINY PAGE */}
        {activeTab === "destiny" && (
          <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="space-y-12">
            <PageHeaderStyle 
              title="Destiny Trending" 
              subtitle="Your today's fortune curve and life fortune curve" 
            />

            <InpageTabSwitcher 
              options={[
                { id: 'today', label: 'Today' },
                { id: 'life', label: 'Life' }
              ]}
              activeId={destinyMode}
              onChange={setDestinyMode}
            />

            <DestinyChart 
              mode={destinyMode} 
              showChart={cardRevealed || destinyMode === "life"} 
              selectedHour={selectedHour} 
              onSelectHour={setSelectedHour} 
              onCardSelect={() => setCardRevealed(true)}
              timeSeries={timeSeries}
              granularity={granularity}
              onGranularityChange={setGranularity}
              selectedDimension={selectedDimension}
              onDimensionChange={setSelectedDimension}
              onVisibleRangeChange={handleVisibleRangeChange}
              isLoading={loading}
            />
            
            {destinyMode === "life" && (
              <div className="pt-12">
                <LifeChart />
              </div>
            )}

            {(cardRevealed || destinyMode === "life") && (
              <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
                <DestinyTimeline 
                  selectedHour={selectedHour} 
                  onSelectHour={setSelectedHour} 
                  timeSeries={timeSeries}
                />
                <div className="pt-12 space-y-4">
                  <PageHeaderStyle title="Stars" />
                  <StarRadar />
                </div>
              </motion.div>
            )}
          </motion.div>
        )}

        {/* ORACLE PAGE */}
        {activeTab === "oracle" && (
          <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="flex flex-col h-[calc(100vh-160px)]">
            <header className="flex items-center gap-4 mb-8">
              <OracleInteractiveEye direction={eyeDirection} className="w-12 h-12" />
              <div>
                <h1 className="text-3xl font-serif font-light"><InkRevealText text="Freya Oracle" /></h1>
                <p className="text-[9px] uppercase tracking-widest opacity-30 font-sans">Your all-knowing cosmic companion</p>
              </div>
            </header>

            <div ref={scrollContainerRef} className="flex-1 overflow-y-auto space-y-6 hide-scrollbar">
              {messages.map(m => <MessageRenderer key={m.id} message={m} />)}
              {isThinking && <OracleThinking />}
            </div>

            <div className="mt-6 space-y-4">
              <FunctionButtons 
                onStoneClick={() => handleFunctionButtonClick("stone")}
                onTarotClick={() => handleFunctionButtonClick("tarot")}
                onEchoClick={() => handleFunctionButtonClick("echo")}
                onConnectClick={() => handleFunctionButtonClick("connect")}
              />
              <div className="flex gap-2 border border-black/10 p-1">
                <input 
                  type="text" 
                  value={input} 
                  onChange={e => setInput(e.target.value)} 
                  onKeyDown={e => e.key === 'Enter' && sendMessage()}
                  placeholder="Ask the oracle..." 
                  className="flex-1 px-4 py-3 text-sm font-light bg-transparent focus:outline-none placeholder:opacity-20"
                />
                <button onClick={sendMessage} className="px-6 py-2 bg-black text-white text-[10px] uppercase tracking-widest">Send</button>
              </div>
            </div>
          </motion.div>
        )}

        {/* ECHO PAGE */}
        {activeTab === "echo" && (
          <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-12">
            <PageHeaderStyle 
              title="Echo Chambers" 
              subtitle="Soul resonance & healing frequencies" 
            />
            <div className="grid gap-6">
              <div className="box-frame p-8 flex flex-col items-center justify-center space-y-6">
                <div className="w-32 h-32 rounded-full border border-black/5 flex items-center justify-center relative">
                  <div className="absolute inset-0 laser-gradient opacity-10 rounded-full animate-pulse" />
                  <span className="material-symbols-outlined text-4xl opacity-20">graphic_eq</span>
                </div>
                <div className="text-center">
                  <h3 className="font-serif text-lg">Harmonic Calibration</h3>
                  <p className="text-[10px] opacity-40 uppercase tracking-widest mt-1">432Hz Resonance</p>
                </div>
                <button className="px-12 py-3 border border-black text-[10px] uppercase tracking-widest hover:bg-black hover:text-white transition-all">Begin Session</button>
              </div>
            </div>
          </motion.div>
        )}

        {/* CONNECT PAGE */}
        {activeTab === "connect" && (
          <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-12">
            <PageHeaderStyle 
              title="Resonance Portal" 
              subtitle="Connect souls across the cosmic field" 
            />
            <div className="box-frame aspect-square flex flex-col items-center justify-center space-y-8">
              <div className="relative w-48 h-48 flex items-center justify-center">
                <div className="absolute inset-0 border border-black/5 rounded-full animate-[spin_20s_linear_infinite]" />
                <div className="absolute inset-4 border border-black/5 rounded-full animate-[spin_15s_linear_infinite_reverse]" />
                <span className="material-symbols-outlined text-5xl opacity-10">all_inclusive</span>
              </div>
              <p className="text-xs opacity-40 font-serif italic">Waiting for connection...</p>
              <button className="px-10 py-3 bg-black text-white text-[10px] uppercase tracking-widest">Search Souls</button>
            </div>
          </motion.div>
        )}
      </main>

      <BottomNav activeTab={activeTab} onTabChange={setActiveTab} />

      {/* Fullscreen Overlays */}
      <AnimatePresence>
        {selectedFunction === "stone" && (
          <StoneUnified isVisible={true} onAccept={() => {}} onReject={() => setSelectedFunction(null)} onComplete={r => { console.log(r); setSelectedFunction(null); }} onExit={() => setSelectedFunction(null)} />
        )}
        {selectedFunction === "tarot" && (
          <TarotUnified isVisible={true} onAccept={() => {}} onReject={() => setSelectedFunction(null)} onComplete={r => { console.log(r); setSelectedFunction(null); }} onExit={() => setSelectedFunction(null)} />
        )}
      </AnimatePresence>
    </div>
  )
}

export default MiniAppDemo;
