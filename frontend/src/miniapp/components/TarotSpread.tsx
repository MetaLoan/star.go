import { useEffect, useRef, useState } from "react"
import { tarotCards } from "../../lib/miniapp/tarot-data"
import type { TarotCardData } from "../../lib/miniapp/oracle-types"
import { InkRevealText } from "./InkRevealText"

interface TarotSpreadProps {
  onSelect: (card: TarotCardData, rect: DOMRect) => void
  disabled?: boolean
  initialPhase?: Phase
  excludedCardIds?: string[]
}

type Phase = "shuffling" | "gathering" | "expanding" | "scrolling" | "selected"

export function TarotSpread({ onSelect, disabled, initialPhase = "shuffling", excludedCardIds = [] }: TarotSpreadProps) {
  const [phase, setPhase] = useState<Phase>(initialPhase)
  const [cards] = useState<TarotCardData[]>(() => {
    return [...tarotCards]
      .sort(() => Math.random() - 0.5)
      .slice(0, 11)
  })
  const cardRefs = useRef<Array<HTMLDivElement | null>>([])
  const [frozenTransforms, setFrozenTransforms] = useState<string[] | null>(null)
  const [gatherArmed, setGatherArmed] = useState(false)

  useEffect(() => {
    if (phase !== "gathering") return
    const raf = requestAnimationFrame(() => setGatherArmed(true))
    const t = window.setTimeout(() => setPhase("expanding"), 600)
    return () => {
      cancelAnimationFrame(raf)
      window.clearTimeout(t)
    }
  }, [phase])

  useEffect(() => {
    if (phase !== "expanding") return
    const t = window.setTimeout(() => setPhase("scrolling"), 800)
    return () => window.clearTimeout(t)
  }, [phase])

  const handleDeckClick = () => {
    if (phase !== "shuffling") return
    const transforms = cardRefs.current.map((el) => {
      if (!el) return "none"
      const t = window.getComputedStyle(el).transform
      return t && t !== "none" ? t : "none"
    })
    setFrozenTransforms(transforms)
    setGatherArmed(false)
    setPhase("gathering")
  }

  const handleCardClick = (card: TarotCardData, rect: DOMRect) => {
    if (phase !== "scrolling" || disabled) return
    onSelect(card, rect)
  }

  return (
    <div className="relative w-full h-full flex flex-col items-center justify-center">
      {phase === "shuffling" && (
        <div 
          className="relative perspective-1000 cursor-pointer"
          onClick={handleDeckClick}
          style={{ width: "calc(192px + 240px)", height: "calc(288px + 80px)", padding: "40px 120px", perspective: '1000px' }}
        >
          {cards.map((card, index) => {
            const isOdd = index % 2 === 1
            return (
              <div
                key={card.id}
                ref={(el) => { cardRefs.current[index] = el }}
                className={`absolute w-48 h-72 ${isOdd ? "animate-shuffle-card-1" : "animate-shuffle-card-2"}`}
                style={{
                  top: "40px",
                  left: "120px",
                  animationDelay: `${index * 0.1}s`,
                  zIndex: 10 + index,
                }}
              >
                <CardBack index={index} />
              </div>
            )
          })}
        </div>
      )}

      {(phase === "gathering" || phase === "expanding" || phase === "scrolling" || phase === "selected") && (
        <div className="relative w-full h-72 flex items-center justify-center overflow-hidden">
          <div 
            className={`relative w-48 h-72 ${phase === "scrolling" ? "animate-scroll-left" : ""}`}
            style={{
              animationPlayState: (phase === "selected" || disabled) ? "paused" : "running",
            }}
          >
            {[...cards, ...cards, ...cards, ...cards].map((card, index) => {
              if (excludedCardIds.includes(card.id)) {
                return null
              }
              
              const cardIndex = index % 11
              const isFirstRound = index < 11
              
              let translateX = 0
              let extraTransform = ""
              let zIndex = index
              let transition = "transform 500ms ease-out, opacity 500ms ease-out"
              let opacity = 1
              let visibility: "visible" | "hidden" = "visible"

              if (phase === "gathering") {
                if (!isFirstRound) {
                  opacity = 0
                  translateX = 0
                  extraTransform = "scale(0.85)"
                } else {
                  const frozen = frozenTransforms?.[index] ?? "none"
                  translateX = 0
                  extraTransform = gatherArmed ? "scale(0.85)" : frozen
                  zIndex = 20 - index
                  transition = gatherArmed
                    ? `transform 500ms cubic-bezier(0.22, 1, 0.36, 1) ${index * 30}ms, opacity 500ms cubic-bezier(0.22, 1, 0.36, 1) ${index * 30}ms`
                    : "none"
                }
              } else if (phase === "expanding") {
                if (!isFirstRound) {
                  opacity = 0
                  translateX = (cardIndex - 5) * 48 + Math.floor(index / 11) * (11 * 48)
                  extraTransform = "scale(0.85)"
                } else {
                  translateX = (index - 5) * 48
                  extraTransform = "scale(0.85)"
                  zIndex = index
                  transition = `transform 700ms cubic-bezier(0.34, 1.56, 0.64, 1) ${index * 50}ms, opacity 700ms cubic-bezier(0.34, 1.56, 0.64, 1) ${index * 50}ms`
                }
              } else {
                translateX = (cardIndex - 5) * 48 + Math.floor(index / 11) * (11 * 48)
                opacity = 1
                extraTransform = "scale(0.85)"
                zIndex = index
              }

              const transform = phase === "gathering" && !gatherArmed && isFirstRound
                ? extraTransform 
                : `translateX(${translateX}px) ${extraTransform}`

              return (
                <div
                  key={`card-${index}`}
                  className="absolute inset-0"
                  style={{
                    transform,
                    zIndex,
                    transition,
                    opacity,
                    visibility,
                    cursor: (phase === "scrolling" && !disabled) ? "pointer" : "default",
                    willChange: "transform, opacity",
                  }}
                  onClick={
                    (phase === "scrolling" && !disabled)
                      ? (e) => {
                          const target = e.currentTarget as HTMLDivElement
                          const rect = target.getBoundingClientRect()
                          target.style.visibility = 'hidden'
                          handleCardClick(card, rect)
                        }
                      : undefined
                  }
                >
                  <CardBack index={cardIndex} isInteractive={phase === "scrolling"} />
                </div>
              )
            })}
          </div>
        </div>
      )}

      <div className="absolute bottom-[-15px] left-0 right-0 text-center pointer-events-none z-20">
        {(phase === "gathering" || phase === "expanding") && (
          <p className="text-xs opacity-60 tracking-widest uppercase font-sans">
            ...
          </p>
        )}
        {phase === "scrolling" && !disabled && (
          <p className="text-xs opacity-100 font-bold uppercase font-sans">
            <InkRevealText text="Choose your card" />
          </p>
        )}
      </div>

      <style>{`
        @keyframes shuffle-card-1 {
          0%, 40% { transform: translateX(0) rotate(0deg) rotateY(0deg); z-index: 10; }
          55% { transform: translateX(120px) rotate(25deg) rotateY(40deg); z-index: 50; }
          70%, 100% { transform: translateX(0) rotate(0deg) rotateY(0deg); z-index: 30; }
        }
        @keyframes shuffle-card-2 {
          0%, 40% { transform: translateX(0) rotate(0deg) rotateY(0deg); z-index: 11; }
          55% { transform: translateX(-120px) rotate(-25deg) rotateY(-40deg); z-index: 50; }
          70%, 100% { transform: translateX(0) rotate(0deg) rotateY(0deg); z-index: 31; }
        }
        @keyframes scroll-left {
          0% { transform: translateX(0); }
          100% { transform: translateX(-528px); }
        }
        .animate-scroll-left { animation: scroll-left 15s linear infinite; }
        .animate-shuffle-card-1 { animation: shuffle-card-1 2.5s ease-in-out infinite; }
        .animate-shuffle-card-2 { animation: shuffle-card-2 2.5s ease-in-out infinite; }
      `}</style>
    </div>
  )
}

function CardBack({ index, isInteractive = false }: { index: number; isInteractive?: boolean }) {
  return (
    <div className={`w-full h-full border border-black dark:border-white rounded bg-white dark:bg-black flex items-center justify-center p-3 overflow-hidden shadow-lg transition-shadow ${isInteractive ? "hover:shadow-xl" : ""}`}>
      <div className="w-full h-full border border-black/20 dark:border-white/20 rounded-sm flex flex-col items-center justify-center relative">
        <svg width="80" height="80" viewBox="0 0 200 200" className="opacity-90">
          <defs>
            <clipPath id={`eyeClip-${index}`}>
              <path d="M 10 100 Q 100 45 190 100 Q 100 155 10 100 Z" />
            </clipPath>
          </defs>
          <path d="M 0 100 L 30 100 M 170 100 L 200 100" stroke="currentColor" strokeWidth="1.5" opacity="0.5" className="text-black dark:text-white" />
          <path d="M 10 100 Q 100 45 190 100 Q 100 155 10 100 Z" fill="none" stroke="currentColor" strokeWidth="2.5" className="text-black dark:text-white" />
          <path d="M 18 100 Q 100 52 182 100 Q 100 148 18 100 Z" fill="none" stroke="currentColor" strokeWidth="1.2" opacity="0.5" className="text-black dark:text-white" />
          <g clipPath={`url(#eyeClip-${index})`}>
            <circle cx="100" cy="100" r="42" fill="none" stroke="currentColor" strokeWidth="2" className="text-black dark:text-white" />
            <circle cx="100" cy="100" r="38" fill="none" stroke="currentColor" strokeWidth="1" strokeDasharray="3 3" className="text-black dark:text-white" />
            {[...Array(48)].map((_, i) => (
              <line key={i} x1="100" y1="100" x2={100 + Math.cos(i * 7.5 * Math.PI / 180) * 42} y2={100 + Math.sin(i * 7.5 * Math.PI / 180) * 42} stroke="currentColor" strokeWidth="0.8" opacity="0.3" className="text-black dark:text-white" />
            ))}
            <ellipse cx="100" cy="100" rx="10" ry="32" fill="currentColor" className="text-black dark:text-white" />
            <ellipse cx="100" cy="100" rx="2" ry="24" fill="currentColor" className="text-white dark:text-black" />
          </g>
          <g stroke="currentColor" strokeWidth="1.2" opacity="0.5" className="text-black dark:text-white">
            <line x1="20" y1="95" x2="35" y2="85" /><line x1="20" y1="105" x2="35" y2="115" />
            <line x1="180" y1="95" x2="165" y2="85" /><line x1="180" y1="105" x2="165" y2="115" />
          </g>
        </svg>
      </div>
    </div>
  )
}

