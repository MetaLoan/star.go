import { useEffect, useState, useRef, useCallback } from "react"
import { X } from "lucide-react"
import { InkRevealText } from "./InkRevealText"
import { TarotSpread } from "./TarotSpread"
import type { TarotResult, TarotCardData } from "../../lib/miniapp/oracle-types"

type Phase = "choice" | "transitioning" | "interaction" | "result" | "exiting"

interface TarotUnifiedProps {
  isVisible: boolean
  onAccept: () => void
  onReject: () => void
  onComplete: (result: TarotResult) => void
  onExit: () => void
}

interface FlyingCard {
  data: TarotCardData
  fromRect: DOMRect
  targetSlot: number
}

export function TarotUnified({ 
  isVisible, 
  onAccept, 
  onReject, 
  onComplete, 
  onExit 
}: TarotUnifiedProps) {
  const FLY_DURATION_MS = 1200
  const [phase, setPhase] = useState<Phase>("choice")
  const [animationState, setAnimationState] = useState<'hidden' | 'entering' | 'visible' | 'exiting'>('hidden')
  
  const [drawnCards, setDrawnCards] = useState<TarotCardData[]>([])
  const [excludedCardIds, setExcludedCardIds] = useState<string[]>([])
  const [flyingCard, setFlyingCard] = useState<FlyingCard | null>(null)
  const [flippedSlots, setFlippedSlots] = useState<boolean[]>([false, false, false])
  const [cardOrientations, setCardOrientations] = useState<('upright' | 'reversed')[]>([])
  
  const slotLabels = ["过去", "现在", "未来"]
  const [showExitConfirm, setShowExitConfirm] = useState(false)
  const slotRefs = useRef<Array<HTMLDivElement | null>>([])
  const timeoutsRef = useRef<number[]>([])

  const schedule = useCallback((fn: () => void, ms: number) => {
    const id = window.setTimeout(fn, ms)
    timeoutsRef.current.push(id)
    return id
  }, [])

  const clearScheduled = useCallback(() => {
    timeoutsRef.current.forEach((id) => window.clearTimeout(id))
    timeoutsRef.current = []
  }, [])

  useEffect(() => {
    if (isVisible) {
      clearScheduled()
      setPhase("choice")
      setDrawnCards([])
      setExcludedCardIds([])
      setFlyingCard(null)
      setFlippedSlots([false, false, false])
      setCardOrientations([])
      setShowExitConfirm(false)
      setAnimationState('hidden')
      schedule(() => setAnimationState('visible'), 50)
    }
    return () => {
      clearScheduled()
    }
  }, [isVisible, clearScheduled, schedule])

  const handleAccept = useCallback(() => {
    setPhase("transitioning")
    onAccept()
    schedule(() => setPhase("interaction"), 1200)
  }, [onAccept, schedule])

  const handleReject = useCallback(() => {
    setAnimationState('exiting')
    schedule(() => onReject(), 500)
  }, [onReject, schedule])

  const handleExit = useCallback(() => {
    if (showExitConfirm) {
      setPhase("exiting")
      schedule(() => onExit(), 500)
    } else {
      setShowExitConfirm(true)
      schedule(() => setShowExitConfirm(false), 3000)
    }
  }, [showExitConfirm, onExit, schedule])

  const handleCardSelect = useCallback((card: TarotCardData, rect: DOMRect) => {
    if (drawnCards.length >= 3 || flyingCard) return
    
    const targetSlot = drawnCards.length
    setFlyingCard({ data: card, fromRect: rect, targetSlot })
    setExcludedCardIds(prev => [...prev, card.id])
    
    const orientation: 'upright' | 'reversed' = Math.random() > 0.5 ? 'upright' : 'reversed'
    setCardOrientations(prev => [...prev, orientation])

    schedule(() => {
      setDrawnCards(prev => [...prev, card])
      setFlyingCard(null)
      
      schedule(() => {
        setFlippedSlots(prev => {
          const next = [...prev]
          next[targetSlot] = true
          return next
        })
        
        if (drawnCards.length === 2) {
          schedule(() => setPhase("result"), 800)
        }
      }, 200)
    }, FLY_DURATION_MS)
  }, [drawnCards.length, flyingCard, schedule])

  const handleFinalComplete = useCallback(() => {
    if (drawnCards.length === 0) return
    setPhase("exiting")
    
    const cardNames = drawnCards.map((card, i) => {
      const orientation = cardOrientations[i] === 'reversed' ? '逆位' : '正位'
      return `${slotLabels[i]}·${card.name}(${orientation})`
    }).join('、')
    
    const fullSummary = drawnCards.map((card, i) => {
      const isReversed = cardOrientations[i] === 'reversed'
      const interpretation = isReversed ? card.reversedMeaning : card.uprightMeaning
      return `【${slotLabels[i]}·${card.name}·${isReversed ? '逆位' : '正位'}】${interpretation}`
    }).join('\n\n')
    
    const result: TarotResult = {
      cards: drawnCards,
      orientations: cardOrientations,
      name: cardNames,
      meaning: drawnCards[0]?.meaning || '',
      summary: fullSummary,
      interpretation: `三张牌阵揭示了你的过去、现在与未来。`
    }
    schedule(() => onComplete(result), 800)
  }, [drawnCards, cardOrientations, onComplete, schedule])

  if (!isVisible) return null

  const isChoicePhase = phase === "choice"
  const isTransitioning = phase === "transitioning"
  const isInteraction = phase === "interaction"

  const getSlotMetrics = (slotIndex: number, startW: number, startH: number) => {
    const slotEl = slotRefs.current[slotIndex]
    if (slotEl) {
      const rect = slotEl.getBoundingClientRect()
      const centerX = rect.left + rect.width / 2
      const centerY = rect.top + rect.height / 2
      const targetScale = rect.width / startW
      return { x: centerX, y: centerY, scale: targetScale }
    }

    const screenWidth = typeof window !== 'undefined' ? window.innerWidth : 800
    const slotWidth = 96
    const gap = 24
    const totalWidth = slotWidth * 3 + gap * 2
    const startX = (screenWidth - totalWidth) / 2
    const x = startX + slotIndex * (slotWidth + gap) + slotWidth / 2
    const y = 148
    const scale = slotWidth / startW
    return { x, y, scale }
  }

  return (
    <>
    <div className={`fixed z-[100] inset-0 transition-all duration-1000 ${phase === "exiting" ? "opacity-0" : "opacity-100"} pointer-events-none`}>
      <div 
        className={`absolute left-0 right-0 transition-all duration-[1200ms] cubic-bezier(0.4, 0, 0.2, 1) ${isChoicePhase ? "bg-white/80 dark:bg-black/80 backdrop-blur-md" : "bg-white dark:bg-black"}`}
        style={{
          bottom: isChoicePhase ? "64px" : "0",
          top: isChoicePhase ? "calc(100% - 384px)" : "0",
          pointerEvents: isChoicePhase ? "none" : "auto",
        }}
      />

      {(isInteraction || phase === "result") && (
        <button onClick={handleExit} className="absolute top-6 right-6 z-[110] p-2 border border-black/30 dark:border-white/30 hover:bg-black/10 dark:hover:bg-white/10 transition-colors pointer-events-auto">
          <X className="w-5 h-5 text-black dark:text-white" />
        </button>
      )}

      {(isInteraction || phase === "result") && (
        <div 
          className="fixed left-0 right-0 flex justify-center gap-3 px-4 pointer-events-none z-[110]"
          style={{
            top: phase === "result" ? "60px" : "50px",
            transition: "top 0.8s cubic-bezier(0.4, 0, 0.2, 1)",
          }}
        >
          {[0, 1, 2].map(i => {
            const isReversed = cardOrientations[i] === 'reversed'
            const orientationLabel = isReversed ? '逆位' : '正位'
            return (
              <div key={i} className="flex flex-col items-center flex-1" style={{ maxWidth: phase === "result" ? "128px" : "96px" }}>
                <div
                  ref={(el) => { slotRefs.current[i] = el }}
                  className="border-[0.5px] border-black/10 dark:border-white/10 bg-black/[0.02] dark:bg-white/[0.02] rounded flex items-center justify-center relative overflow-hidden w-full"
                  style={{ 
                    perspective: '600px',
                    aspectRatio: '2/3',
                    transition: "max-width 0.8s cubic-bezier(0.4, 0, 0.2, 1)",
                  }}
                >
                  {drawnCards[i] && (
                    <div 
                      className="w-full h-full relative"
                      style={{
                        transformStyle: 'preserve-3d',
                        transition: 'transform 0.6s cubic-bezier(0.4, 0, 0.2, 1)',
                        transform: flippedSlots[i] ? 'rotateY(180deg)' : 'rotateY(0deg)',
                      }}
                    >
                      <div 
                        className="absolute inset-0 bg-white dark:bg-black flex items-center justify-center p-1 rounded"
                        style={{ backfaceVisibility: 'hidden' }}
                      >
                        <div className="w-full h-full border border-black/30 dark:border-white/30 rounded-sm flex items-center justify-center">
                          <svg width="32" height="32" viewBox="0 0 200 200" className="opacity-60 text-black dark:text-white">
                            <path d="M 10 100 Q 100 45 190 100 Q 100 155 10 100 Z" fill="none" stroke="currentColor" strokeWidth="3" />
                            <ellipse cx="100" cy="100" rx="8" ry="26" fill="currentColor" />
                          </svg>
                        </div>
                      </div>
                      <div 
                        className="absolute inset-0"
                        style={{ 
                          backfaceVisibility: 'hidden',
                          transform: `rotateY(180deg) ${isReversed ? 'rotate(180deg)' : ''}`,
                        }}
                      >
                        <img src={drawnCards[i].image} alt={drawnCards[i].name} className="w-full h-full object-cover grayscale rounded" />
                      </div>
                    </div>
                  )}
                </div>
                {phase === "result" && drawnCards[i] && (
                  <p 
                    className="text-[10px] text-black/50 dark:text-white/50 mt-2 tracking-wider font-sans"
                    style={{
                      opacity: flippedSlots[i] ? 1 : 0,
                      transition: 'opacity 0.5s ease-out 0.3s',
                    }}
                  >
                    {slotLabels[i]}·{orientationLabel}
                  </p>
                )}
              </div>
            )
          })}
        </div>
      )}

      {isInteraction && drawnCards.length < 3 && (
        <div 
          className="fixed left-0 right-0 flex justify-center pointer-events-none z-[105]"
          style={{
            top: "210px",
            transition: "opacity 0.5s ease-out",
          }}
        >
          <p className="text-xs text-black/40 dark:text-white/40 tracking-widest font-sans uppercase">
            {3 - drawnCards.length} cards remaining
          </p>
        </div>
      )}

      {phase === "result" && (
        <div 
          className="fixed left-0 right-0 flex flex-col items-center pointer-events-none px-6 z-[105]"
          style={{
            top: "300px",
            opacity: 1,
            transition: "opacity 0.8s ease-out",
          }}
        >
          <InkRevealText text="基础卡牌解读" className="text-lg font-light mb-4 font-serif text-black dark:text-white" />
          
          <div className="w-full max-w-md space-y-4 mb-6">
            {drawnCards.map((card, i) => {
              const isReversed = cardOrientations[i] === 'reversed'
              const interpretation = isReversed ? card.reversedMeaning : card.uprightMeaning
              return (
                <div key={card.id} className="text-left">
                  <p className="text-xs font-medium text-black dark:text-white mb-1">
                    {slotLabels[i]} · {card.name} · {isReversed ? '逆位' : '正位'}
                  </p>
                  <p className="text-xs text-black/60 dark:text-white/60 leading-relaxed font-sans">
                    {interpretation}
                  </p>
                </div>
              )
            })}
          </div>
          
          <button 
            onClick={handleFinalComplete} 
            className="px-12 py-3 border border-black dark:border-white text-sm font-light hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-all pointer-events-auto text-black dark:text-white bg-white dark:bg-black"
          >
            Continue
          </button>
        </div>
      )}

      <div 
        className="absolute inset-0 pointer-events-none"
        style={{
          display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
          transform: isChoicePhase ? "translateY(calc(50vh - 224px - 40px))" : (isInteraction || isTransitioning) ? "translateY(calc(50vh - 144px - 35px))" : "translateY(0)", 
          opacity: animationState === 'visible' ? 1 : 0,
          filter: animationState === 'visible' ? 'blur(0px)' : 'blur(10px)',
          transition: "transform 1.2s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.8s ease-out, filter 0.8s ease-out",
          willChange: "transform, opacity, filter",
        }}
      >
        <div 
          className="relative pointer-events-auto"
          style={{
            transition: "transform 1.2s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.8s ease-out",
            opacity: phase === "exiting" || phase === "result" ? 0 : 1,
            transform: isChoicePhase ? "scale(0.33)" : "scale(1)",
            zIndex: 20,
            width: '100%',
            height: isChoicePhase ? '320px' : '100%',
            pointerEvents: phase === "result" ? "none" : "auto",
          }}
        >
          {(isChoicePhase || isTransitioning || isInteraction) && (
            <TarotSpread 
              key={isChoicePhase ? "shuffling" : "active"} 
              onSelect={handleCardSelect} 
              disabled={!!flyingCard || drawnCards.length >= 3} 
              initialPhase={isChoicePhase ? "shuffling" : "gathering"}
              excludedCardIds={excludedCardIds}
            />
          )}
        </div>

        {isChoicePhase && (
          <div className="absolute w-full flex flex-col items-center pointer-events-none" style={{ top: "calc(50% + 60px)", opacity: animationState === 'visible' ? 1 : 0 }}>
            <h3 className="text-lg font-light tracking-wide mb-4 font-serif text-black dark:text-white"><InkRevealText text="是否进入塔罗占卜" /></h3>
            <div className="flex gap-4 pointer-events-auto">
              <button onClick={handleAccept} className="w-32 border border-black dark:border-white py-3 text-sm font-light bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-colors shadow-sm">接受</button>
              <button onClick={handleReject} className="w-32 border border-black dark:border-white py-3 text-sm font-light bg-black dark:bg-white text-white dark:text-black hover:bg-white dark:hover:bg-black hover:text-black dark:hover:text-white transition-colors shadow-sm">拒绝</button>
            </div>
          </div>
        )}
      </div>

    </div>

    {flyingCard && (() => {
      const current = flyingCard
      const startX = flyingCard.fromRect.left
      const startY = flyingCard.fromRect.top
      const startW = flyingCard.fromRect.width
      const startH = flyingCard.fromRect.height
      const slot = getSlotMetrics(current.targetSlot, startW, startH)
      const deltaX = slot.x - (startX + startW / 2)
      const deltaY = slot.y - (startY + startH / 2)

      return (
        <>
          <div className="fixed inset-0 pointer-events-none z-[300]">
            <div
              className="absolute pointer-events-none"
              style={{
                left: startX,
                top: startY,
                width: startW,
                height: startH,
                animation: `fly-to-slot-${current.targetSlot} ${FLY_DURATION_MS}ms cubic-bezier(0.4, 0, 0.2, 1) forwards`,
                transformOrigin: 'center center',
              }}
            >
              <div className="w-full h-full border border-black dark:border-white rounded bg-white dark:bg-black flex items-center justify-center p-3 overflow-hidden shadow-xl">
                <div className="w-full h-full border border-black/30 dark:border-white/30 rounded-sm flex flex-col items-center justify-center relative">
                  <svg width="80" height="80" viewBox="0 0 200 200" className="opacity-90 text-black dark:text-white">
                    <defs>
                      <clipPath id="flyingEyeClip">
                        <path d="M 10 100 Q 100 45 190 100 Q 100 155 10 100 Z" />
                      </clipPath>
                    </defs>
                    <path d="M 0 100 L 30 100 M 170 100 L 200 100" stroke="currentColor" strokeWidth="1.5" opacity="0.5" />
                    <path d="M 10 100 Q 100 45 190 100 Q 100 155 10 100 Z" fill="none" stroke="currentColor" strokeWidth="2.5" />
                    <path d="M 18 100 Q 100 52 182 100 Q 100 148 18 100 Z" fill="none" stroke="currentColor" strokeWidth="1.2" opacity="0.5" />
                    <g clipPath="url(#flyingEyeClip)">
                      <circle cx="100" cy="100" r="42" fill="none" stroke="currentColor" strokeWidth="2" />
                      <circle cx="100" cy="100" r="38" fill="none" stroke="currentColor" strokeWidth="1" strokeDasharray="3 3" />
                      <ellipse cx="100" cy="100" rx="10" ry="32" fill="currentColor" />
                      <ellipse cx="100" cy="100" rx="2" ry="24" fill="currentColor" className="text-white dark:text-black" />
                    </g>
                  </svg>
                </div>
              </div>
            </div>
          </div>
          <style>{`
            @keyframes fly-to-slot-${current.targetSlot} {
              0% {
                transform: translate(0px, 0px);
              }
              35% {
                transform: translate(0px, -100px);
              }
              100% {
                transform: translate(${deltaX}px, ${deltaY}px);
                opacity: 1;
              }
            }
          `}</style>
        </>
      )
    })()}
    </>
  )
}

