import { useEffect, useState, useRef, useCallback } from "react"
import { X } from "lucide-react"
import { InkRevealText } from "./InkRevealText"
import { Stone3D } from "./Stone3D"
import type { StoneResult } from "../../lib/miniapp/oracle-types"

type Phase = "choice" | "transitioning" | "interaction" | "result" | "exiting"

interface StoneUnifiedProps {
  isVisible: boolean
  onAccept: () => void
  onReject: () => void
  onComplete: (result: StoneResult) => void
  onExit: () => void
}

const ORACLE_MESSAGES = [
  { result: "YES" as const, message: "The path ahead glows with certainty" },
  { result: "NO" as const, message: "The shadows whisper caution" },
  { result: "WAIT" as const, message: "Time holds the answer you seek" },
  { result: "SILENCE" as const, message: "Some truths reveal themselves in stillness" },
  { result: "RELEASE" as const, message: "Let go of what no longer serves you" },
]

export function StoneUnified({ 
  isVisible, 
  onAccept, 
  onReject, 
  onComplete, 
  onExit 
}: StoneUnifiedProps) {
  const [phase, setPhase] = useState<Phase>("choice")
  const [showExitConfirm, setShowExitConfirm] = useState(false)
  
  const [isHolding, setIsHolding] = useState(false)
  const [progress, setProgress] = useState(0)
  const [revealed, setRevealed] = useState(false)
  const [isShaking, setIsShaking] = useState(false)
  const [oracle, setOracle] = useState(ORACLE_MESSAGES[0])
  const [, setUserChoice] = useState<"follow" | "rebel" | null>(null)
  
  const holdStartRef = useRef<number | null>(null)
  const animationRef = useRef<number | null>(null)
  const progressRef = useRef(0)

  const [animationState, setAnimationState] = useState<'hidden' | 'entering' | 'visible' | 'exiting'>('hidden')

  useEffect(() => {
    if (isVisible) {
      setPhase("choice")
      setProgress(0)
      setRevealed(false)
      setIsHolding(false)
      setUserChoice(null)
      setShowExitConfirm(false)
      setAnimationState('hidden')
      
      const timer = setTimeout(() => {
        setAnimationState('visible')
      }, 50)
      return () => clearTimeout(timer)
    }
  }, [isVisible])

  const handleAccept = useCallback(() => {
    setPhase("transitioning")
    onAccept()
    
    setTimeout(() => {
      setPhase("interaction")
    }, 1200)
  }, [onAccept])

  const handleReject = useCallback(() => {
    setAnimationState('exiting')
    setTimeout(() => {
      onReject()
    }, 500)
  }, [onReject])

  const handleExit = useCallback(() => {
    if (showExitConfirm) {
      setPhase("exiting")
      setTimeout(() => {
        onExit()
      }, 500)
    } else {
      setShowExitConfirm(true)
      setTimeout(() => setShowExitConfirm(false), 3000)
    }
  }, [showExitConfirm, onExit])

  const handleHoldStart = useCallback(() => {
    if (phase !== "interaction" || revealed) return
    
    setIsHolding(true)
    holdStartRef.current = Date.now()
    progressRef.current = 0
    
    const animate = () => {
      if (!holdStartRef.current) return
      
      const elapsed = Date.now() - holdStartRef.current
      const newProgress = Math.min(elapsed / 3000, 1)
      progressRef.current = newProgress
      setProgress(newProgress)
      
      if (newProgress > 0.8 && !isShaking) {
        setIsShaking(true)
      }
      
      if (newProgress >= 1) {
        setIsHolding(false)
        setRevealed(true)
        setIsShaking(false)
        holdStartRef.current = null
        
        const randomOracle = ORACLE_MESSAGES[Math.floor(Math.random() * ORACLE_MESSAGES.length)]
        setOracle(randomOracle)
        setPhase("result")
        return
      }
      
      animationRef.current = requestAnimationFrame(animate)
    }
    
    animationRef.current = requestAnimationFrame(animate)
  }, [phase, revealed, isShaking])

  const handleHoldEnd = useCallback(() => {
    if (!isHolding) return
    
    setIsHolding(false)
    setIsShaking(false)
    holdStartRef.current = null
    
    if (animationRef.current) {
      cancelAnimationFrame(animationRef.current)
    }
    
    const resetProgress = () => {
      setProgress((prev) => {
        if (prev <= 0) return 0
        requestAnimationFrame(resetProgress)
        return prev - 0.02
      })
    }
    resetProgress()
  }, [isHolding])

  const handleChoice = useCallback((choice: "follow" | "rebel") => {
    setUserChoice(choice)
    setPhase("exiting")
    
    const result: StoneResult = {
      result: oracle.result,
      message: oracle.message,
      choice
    }
    
    setTimeout(() => {
      onComplete(result)
    }, 800)
  }, [oracle, onComplete])

  useEffect(() => {
    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current)
      }
    }
  }, [])

  if (!isVisible) return null

  const isChoicePhase = phase === "choice"
  const isTransitioning = phase === "transitioning"

  return (
    <div 
      className={`fixed z-[100] inset-0 transition-all duration-1000 ${
        phase === "exiting" ? "opacity-0" : "opacity-100"
      } pointer-events-none`}
    >
      <div 
        className={`absolute left-0 right-0 transition-all duration-[1200ms] cubic-bezier(0.4, 0, 0.2, 1) ${
          isChoicePhase ? "bg-white/80 dark:bg-black/80 backdrop-blur-md" : "bg-white dark:bg-black"
        }`}
        style={{
          bottom: isChoicePhase ? "64px" : "0",
          top: isChoicePhase ? "calc(100% - 384px)" : "0",
          opacity: 1,
          pointerEvents: isChoicePhase ? "none" : "auto",
        }}
      />

      {(phase === "interaction" || phase === "result") && (
        <button
          onClick={handleExit}
          className="absolute top-6 right-6 z-50 p-2 border border-black/20 dark:border-white/20 hover:bg-black/10 dark:hover:bg-white/10 transition-colors pointer-events-auto"
          aria-label="Exit"
        >
          <X className="w-5 h-5 text-black dark:text-white" />
        </button>
      )}

      {showExitConfirm && (
        <div className="absolute top-20 right-6 z-50 border border-black dark:border-white bg-white dark:bg-black p-4 max-w-xs pointer-events-auto">
          <p className="text-xs font-light mb-2">确定要退出吗？</p>
          <p className="text-[10px] opacity-60 font-light">退出后将无法继续此次占卜</p>
        </div>
      )}

      <div 
        className="absolute inset-0 pointer-events-none"
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          transform: isChoicePhase 
            ? "translateY(calc(50vh - 224px - 40px))" 
            : "translateY(0)", 
          opacity: animationState === 'visible' ? 1 : 0,
          filter: animationState === 'visible' ? 'blur(0px)' : 'blur(10px)',
          transition: "transform 1.2s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.8s ease-out, filter 0.8s ease-out",
          willChange: "transform, opacity, filter",
        }}
      >
        <div 
          className={`absolute pointer-events-auto`}
          style={{
            transition: "all 1.2s cubic-bezier(0.4, 0, 0.2, 1)",
            opacity: isChoicePhase 
              ? (animationState === 'visible' || animationState === 'entering' ? 1 : 0)
              : (phase === "exiting" ? 0 : 1),
            transform: isChoicePhase ? "scale(1)" : "scale(1.67)",
            zIndex: 20,
          }}
          onMouseDown={phase === "interaction" ? handleHoldStart : undefined}
          onMouseUp={handleHoldEnd}
          onMouseLeave={handleHoldEnd}
          onTouchStart={phase === "interaction" ? handleHoldStart : undefined}
          onTouchEnd={handleHoldEnd}
        >
          <Stone3D
            size={120}
            isHolding={isHolding}
            progress={progress}
            revealed={revealed}
            isShaking={isShaking}
          />
        </div>

        <div 
          className="absolute w-full flex flex-col items-center pointer-events-none"
          style={{
            top: "calc(50% + 60px)",
            opacity: (isTransitioning || !isChoicePhase) ? 0 : (animationState === 'visible' || animationState === 'entering' ? 1 : 0),
            transform: (isTransitioning || !isChoicePhase) ? "translateY(40px)" : "translateY(0)",
            transition: "all 0.8s cubic-bezier(0.4, 0, 0.2, 1)",
          }}
        >
          <div className="flex flex-col items-center mb-2">
            <h3 className="text-lg font-light tracking-wide font-serif">
              {(animationState === 'visible' || animationState === 'entering') && <InkRevealText text="是否进入Fate Stone" />}
            </h3>
          </div>
          <div className="w-px h-5 bg-black/20 dark:bg-white/20 mb-2" />
          <div className="flex gap-4 w-full justify-center px-6 max-w-sm pointer-events-auto">
            <button
              onClick={handleAccept}
              className="w-32 border border-black dark:border-white py-3 text-sm font-light bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-colors shadow-sm"
            >
              接受
            </button>
            <button
              onClick={handleReject}
              className="w-32 border border-black dark:border-white py-3 text-sm font-light bg-black dark:bg-white text-white dark:text-black hover:bg-white dark:hover:bg-black hover:text-black dark:hover:text-white transition-colors shadow-sm"
            >
              拒绝
            </button>
          </div>
        </div>

        <div 
          className="absolute w-full text-center"
          style={{ 
            top: "calc(50% + 160px)",
            opacity: isTransitioning ? 1 : 0, 
            transition: "opacity 1s",
            pointerEvents: "none"
          }}
        >
          <p className="text-sm opacity-40 font-light animate-pulse font-serif">
            正在连接命运之石...
          </p>
        </div>

        <div 
          className="absolute w-full text-center" 
          style={{ 
            top: "calc(50% + 180px)",
            opacity: phase === "interaction" ? 0.6 : 0,
            transform: phase === "interaction" ? "translateY(0)" : "translateY(20px)",
            transition: "all 1s cubic-bezier(0.4, 0, 0.2, 1) 0.5s"
          }}
        >
          <p className="text-sm font-light font-serif">
            长按石头，直到它揭示答案
          </p>
        </div>

        {phase === "result" && (
          <div 
            className="absolute w-full flex flex-col items-center pointer-events-auto"
            style={{ top: "calc(50% + 140px)" }}
          >
            <p className="text-2xl font-light tracking-widest mb-2 font-serif text-black dark:text-white">
              {oracle.result}
            </p>
            <p className="text-sm opacity-60 font-light mb-8 text-center max-w-xs font-serif text-black dark:text-white">
              {oracle.message}
            </p>

            <div className="flex gap-4">
              <button
                onClick={() => handleChoice("follow")}
                className="px-8 py-3 border border-black dark:border-white text-sm font-light hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-all text-black dark:text-white"
              >
                Follow
              </button>
              <button
                onClick={() => handleChoice("rebel")}
                className="px-8 py-3 border border-black dark:border-white bg-black dark:bg-white text-white dark:text-black text-sm font-light hover:bg-white dark:hover:bg-black hover:text-black dark:hover:text-white transition-all"
              >
                No Way
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

