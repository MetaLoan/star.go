import { useState, useRef, useEffect } from "react"
import { InkRevealText } from "./InkRevealText"
import { TextReveal } from "./TextReveal"
import { Stone3D } from "./Stone3D"
import type { StoneResult } from "../../lib/miniapp/oracle-types"

const oracles = [
  { result: "YES" as const, message: "The path ahead glows with certainty" },
  { result: "NO" as const, message: "Resistance is the universe's protection" },
  { result: "WAIT" as const, message: "Timing is the hidden dimension of fate" },
  { result: "SILENCE" as const, message: "The silence is the answer" },
  { result: "RELEASE" as const, message: "Let go of what you cannot control" },
]

const holdingMessages = [
  { text: "Hold your breath...", speed: 1, shake: false },
  { text: "Focus on your question...", speed: 1, shake: false },
  { text: "Breathe steadily...", speed: 1, shake: false },
  { text: "I see it! The question deep in your heart!", speed: 2, shake: true },
  { text: "This question seems...", speed: 2, shake: true },
  { text: "Let me answer you...", speed: 1, shake: false },
]

const getRandomWaitTime = () => 3000 + Math.random() * 2000
const unfocusedMessage = "It seems you lack focus. I cannot reveal the answer."

interface FateStoneEmbeddedProps {
  onComplete: (result: StoneResult) => void
  onCancel?: () => void
  embedded?: boolean
}

export function FateStoneEmbedded({ onComplete, onCancel, embedded = false }: FateStoneEmbeddedProps) {
  const [isHolding, setIsHolding] = useState(false)
  const [progress, setProgress] = useState(0)
  const [currentMessageIndex, setCurrentMessageIndex] = useState(0)
  const [messageMode, setMessageMode] = useState<"fadein" | "fadeaway" | "idle">("idle")
  const [isShaking, setIsShaking] = useState(false)
  const [revealed, setRevealed] = useState(false)
  const [oracle, setOracle] = useState(oracles[0])
  const [showChoice, setShowChoice] = useState(false)
  const [isFadingOut, setIsFadingOut] = useState(false)
  const [showUnfocusedMessage, setShowUnfocusedMessage] = useState(false)
  
  const holdTimerRef = useRef<number | null>(null)
  const progressTimerRef = useRef<number | null>(null)
  const messageTimerRef = useRef<number | null>(null)
  const stoneRef = useRef<HTMLDivElement>(null)
  const completedRef = useRef(false)

  const startHold = () => {
    if (revealed) return

    setIsHolding(true)
    setProgress(0)
    setCurrentMessageIndex(0)
    setMessageMode("fadein")
    setShowUnfocusedMessage(false)
    completedRef.current = false

    progressTimerRef.current = window.setInterval(() => {
      setProgress((prev) => Math.min(prev + 2, 100))
    }, 30)
  }
  
  const handleMessageFadeInComplete = () => {
    const waitTime = getRandomWaitTime()
    const isLastMessage = currentMessageIndex === holdingMessages.length - 1
    
    messageTimerRef.current = window.setTimeout(() => {
      if (isLastMessage) {
        completedRef.current = true
        triggerResult()
      } else {
        setMessageMode("fadeaway")
      }
    }, waitTime)
  }
  
  const handleMessageFadeOutComplete = () => {
    const nextIndex = currentMessageIndex + 1
    setCurrentMessageIndex(nextIndex)
    setMessageMode("fadein")
    if (nextIndex < holdingMessages.length) {
      setIsShaking(holdingMessages[nextIndex].shake)
    }
  }

  const triggerResult = () => {
    if (progressTimerRef.current) clearInterval(progressTimerRef.current)
    if (messageTimerRef.current) clearTimeout(messageTimerRef.current)
    
    const randomOracle = oracles[Math.floor(Math.random() * oracles.length)]
    setOracle(randomOracle)
    setRevealed(true)
    setIsHolding(false)
    setProgress(0)
    setMessageMode("idle")
    setIsShaking(false)

    setTimeout(() => {
      setShowChoice(true)
    }, 1500)
  }

  const releaseStone = () => {
    if (!isHolding) return
    
    if (holdTimerRef.current) clearTimeout(holdTimerRef.current)
    if (progressTimerRef.current) clearInterval(progressTimerRef.current)
    if (messageTimerRef.current) clearTimeout(messageTimerRef.current)

    if (completedRef.current || revealed) return

    setShowUnfocusedMessage(true)
    setIsHolding(false)
    setProgress(0)
    setMessageMode("idle")
    setIsShaking(false)
    
    setTimeout(() => {
      setShowUnfocusedMessage(false)
    }, 3000)
  }

  const makeChoice = (choice: "follow" | "rebel") => {
    if (isFadingOut) return
    setIsFadingOut(true)
    
    const result: StoneResult = {
      result: oracle.result,
      message: oracle.message,
      choice
    }
    
    setTimeout(() => {
      onComplete(result)
    }, 600)
  }

  useEffect(() => {
    return () => {
      if (holdTimerRef.current) clearTimeout(holdTimerRef.current)
      if (progressTimerRef.current) clearInterval(progressTimerRef.current)
      if (messageTimerRef.current) clearTimeout(messageTimerRef.current)
    }
  }, [])

  const containerClass = embedded 
    ? "relative w-full h-[400px] min-h-[400px]" 
    : "relative w-full h-full"

  return (
    <div className={containerClass}>
      <div 
        className="absolute inset-0 flex items-center justify-center transition-all duration-700 ease-out"
        style={{
          transform: revealed ? 'translateY(-20px) scale(0.9)' : 'translateY(-20px) scale(1.2)',
          opacity: revealed ? 0.5 : 1,
          filter: revealed ? 'blur(8px)' : 'blur(0px)',
        }}
        ref={stoneRef}
      >
        <div
          onMouseDown={startHold}
          onMouseUp={releaseStone}
          onMouseLeave={releaseStone}
          onTouchStart={startHold}
          onTouchEnd={releaseStone}
          className="relative z-0"
          style={{ 
            cursor: revealed ? 'default' : 'pointer',
            pointerEvents: revealed ? 'none' : 'auto'
          }}
        >
          <Stone3D 
            size={embedded ? 180 : 240}
            isHolding={isHolding}
            progress={progress}
            revealed={revealed}
            isShaking={isShaking}
          />
        </div>
      </div>

      <div className="absolute bottom-0 left-0 right-0 h-[220px] flex flex-col items-center justify-start pointer-events-none z-10">
        <div className="pointer-events-auto flex flex-col items-center w-full">
          {revealed && (
            <div 
              className="text-center space-y-4 transition-all duration-500 ease-out"
              style={{
                opacity: isFadingOut ? 0 : 1,
                transform: isFadingOut ? 'translateY(-20px)' : 'translateY(0)',
              }}
            >
              <div className="text-5xl font-light tracking-wider font-serif">
                <InkRevealText text={oracle.result} />
              </div>
              <div className="text-sm opacity-60 max-w-xs mx-auto leading-relaxed">
                <InkRevealText text={oracle.message} />
              </div>
            </div>
          )}

          {showChoice && (
            <div 
              className="flex gap-4 w-full max-w-sm px-6 mt-6 pb-5 mb-5 transition-all duration-500 ease-out"
              style={{
                opacity: isFadingOut ? 0 : 1,
                transform: isFadingOut ? 'translateY(20px)' : 'translateY(0)',
              }}
            >
              <button
                onClick={() => makeChoice("follow")}
                className="flex-1 border border-black dark:border-white py-3 text-sm font-light hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-all"
              >
                Follow
              </button>
              <button
                onClick={() => makeChoice("rebel")}
                className="flex-1 border border-black dark:border-white bg-black text-white dark:bg-white dark:text-black py-3 text-sm font-light hover:bg-white hover:text-black dark:hover:bg-black dark:hover:text-white transition-all"
              >
                No Way
              </button>
            </div>
          )}
        </div>
      </div>

      {!revealed && (
        <div className="absolute bottom-[-15px] left-0 right-0 text-center pointer-events-none z-20">
          <div className="pointer-events-none">
            {showUnfocusedMessage ? (
              <p className="text-xs opacity-60 font-light tracking-wide italic">
                <InkRevealText text={unfocusedMessage} />
              </p>
            ) : isHolding && messageMode !== "idle" ? (
              <p 
                key={currentMessageIndex}
                className="text-xs opacity-60 font-light tracking-wide"
              >
                <TextReveal 
                  key={`msg-${currentMessageIndex}-${messageMode}`}
                  text={holdingMessages[currentMessageIndex].text} 
                  mode={messageMode}
                  delayPerChar={messageMode === "fadein" 
                    ? Math.round(40 / holdingMessages[currentMessageIndex].speed) 
                    : Math.round(30 / holdingMessages[currentMessageIndex].speed)}
                  onComplete={messageMode === "fadein" ? handleMessageFadeInComplete : handleMessageFadeOutComplete}
                />
              </p>
            ) : !isHolding && (
              <>
                <p className="text-xs opacity-100 font-bold uppercase tracking-widest font-sans">
                  <InkRevealText text="Hold to seek guidance" />
                </p>
                <p className="text-[10px] text-black/40 dark:text-white/40 mt-2 font-normal uppercase tracking-tighter">3 times remaining today</p>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

