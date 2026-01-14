import { useEffect, useState, useCallback, useRef } from "react"
import { X, Music, ChevronLeft } from "lucide-react"
import { InkRevealText } from "./InkRevealText"
import { EchoEffect } from "./EchoEffect"
import { AudioPlayer } from "./AudioPlayer"
import type { EchoResult } from "../../lib/miniapp/oracle-types"

type Phase = "choice" | "library" | "transitioning" | "playing" | "result" | "exiting"

const echoLibrary = [
  {
    id: "deep-sleep",
    title: "Deep Sleep",
    subtitle: "Delta waves for restorative rest",
    frequency: "432Hz + White Noise",
    duration: "8:00",
    description: "Gentle delta wave frequencies tuned to 432Hz harmonize with your natural sleep cycle. The layered white noise masks environmental disturbances, guiding you into deep restorative sleep within minutes.",
  },
  {
    id: "focus-flow",
    title: "Focus Flow",
    subtitle: "Gamma waves for concentration",
    frequency: "40Hz Binaural",
    duration: "25:00",
    description: "40Hz gamma binaural beats enhance neural synchronization and cognitive clarity. Ideal for deep work sessions, creative problem-solving, or whenever you need laser-sharp mental focus.",
  },
  {
    id: "anxiety-shield",
    title: "Anxiety Shield",
    subtitle: "Calming frequencies",
    frequency: "528Hz + Singing Bowl",
    duration: "12:00",
    description: "The 528Hz 'miracle tone' combined with Tibetan singing bowls creates a protective sonic cocoon. Reduces cortisol levels and activates your parasympathetic nervous system for instant calm.",
  },
  {
    id: "energy-boost",
    title: "Energy Boost",
    subtitle: "Morning activation",
    frequency: "Beta Wave 15Hz",
    duration: "10:00",
    description: "Energizing beta frequencies stimulate alertness and mental clarity. Perfect for replacing your morning coffee or overcoming afternoon fatigue. Awakens body and mind naturally.",
  },
]

interface EchoUnifiedProps {
  isVisible: boolean
  recommendation: string
  frequency: number
  trackName: string
  isGenerating?: boolean
  externalProgress?: number
  onAccept: () => void
  onReject: () => void
  onCancel: () => void
  onComplete: (result: EchoResult) => void
  onExit: () => void
}

export function EchoUnified({
  isVisible,
  recommendation,
  frequency,
  trackName,
  isGenerating = false,
  externalProgress = 0,
  onAccept,
  onReject,
  onCancel,
  onComplete,
  onExit
}: EchoUnifiedProps) {
  const [phase, setPhase] = useState<Phase>("choice")
  const [animationState, setAnimationState] = useState<'hidden' | 'visible' | 'exiting'>('hidden')
  const [playProgress] = useState(0)
  const [generateProgress, setGenerateProgress] = useState(0)
  const [playCountMap, setPlayCountMap] = useState<Record<string, number>>({
    "deep-sleep": 5,
    "focus-flow": 3,
    "anxiety-shield": 1,
    "energy-boost": 0,
  })
  const timeoutsRef = useRef<number[]>([])
  const progressIntervalRef = useRef<number | null>(null)
  const generateIntervalRef = useRef<number | null>(null)

  const PLAY_DURATION = 30

  const schedule = useCallback((fn: () => void, ms: number) => {
    const id = window.setTimeout(fn, ms)
    timeoutsRef.current.push(id)
    return id
  }, [])

  const clearScheduled = useCallback(() => {
    timeoutsRef.current.forEach((id) => window.clearTimeout(id))
    timeoutsRef.current = []
    if (progressIntervalRef.current) {
      window.clearInterval(progressIntervalRef.current)
      progressIntervalRef.current = null
    }
    if (generateIntervalRef.current) {
      window.clearInterval(generateIntervalRef.current)
      generateIntervalRef.current = null
    }
  }, [])

  useEffect(() => {
    if (isVisible) {
      if (animationState === 'hidden') {
        clearScheduled()
        if (isGenerating) {
          setPhase("transitioning")
          setGenerateProgress(externalProgress)
        } else {
          setPhase("choice")
          setGenerateProgress(0)
        }
        setAnimationState('hidden')
        schedule(() => setAnimationState('visible'), 50)
      }
    } else {
      setAnimationState('hidden')
      clearScheduled()
    }
  }, [isVisible, isGenerating, clearScheduled, schedule])

  useEffect(() => {
    if (isGenerating && phase === "transitioning") {
      setGenerateProgress(externalProgress)
    }
  }, [externalProgress, isGenerating, phase])

  const handleAccept = useCallback(() => {
    setPhase("transitioning")
    onAccept()
  }, [onAccept])

  const handleReject = useCallback(() => {
    setAnimationState('exiting')
    schedule(() => onReject(), 500)
  }, [onReject, schedule])

  const handleExit = useCallback(() => {
    clearScheduled()
    setPhase("exiting")
    schedule(() => onExit(), 500)
  }, [onExit, schedule, clearScheduled])

  const handleComplete = useCallback(() => {
    setPhase("exiting")
    const result: EchoResult = {
      frequency,
      duration: PLAY_DURATION,
      trackName,
      purpose: "磁场调整",
      recommendation
    }
    schedule(() => onComplete(result), 800)
  }, [frequency, trackName, recommendation, onComplete, schedule])

  const handleShowLibrary = useCallback(() => {
    setPhase("library")
  }, [])

  const handleViewLater = useCallback(() => {
    setAnimationState('exiting')
    schedule(() => onExit(), 500)
  }, [schedule, onExit])

  const handleCancel = useCallback(() => {
    clearScheduled()
    setAnimationState('exiting')
    schedule(() => onCancel(), 500)
  }, [clearScheduled, schedule, onCancel])

  const handleBackToChoice = useCallback(() => {
    setPhase("choice")
    setAnimationState('visible')
  }, [])

  const handlePlayComplete = useCallback((audioId: string) => {
    setPlayCountMap((prev) => ({
      ...prev,
      [audioId]: Math.min((prev[audioId] || 0) + 1, 5),
    }))
  }, [])

  if (!isVisible) return null

  const isChoicePhase = phase === "choice"
  const isTransitioning = phase === "transitioning"
  const isChoiceOrTransitioning = isChoicePhase || isTransitioning
  const isLibraryPhase = phase === "library"
  const isPlaying = phase === "playing"
  const isResult = phase === "result"

  return (
    <div className={`fixed z-[100] inset-0 transition-all duration-1000 ${phase === "exiting" ? "opacity-0" : "opacity-100"} pointer-events-none`}>
      <div 
        className={`absolute left-0 right-0 transition-all duration-[1200ms] cubic-bezier(0.4, 0, 0.2, 1) ${(isChoiceOrTransitioning || isLibraryPhase) ? "bg-white/80 dark:bg-black/80 backdrop-blur-md" : "bg-white dark:bg-black"}`}
        style={{
          bottom: (isChoiceOrTransitioning || isLibraryPhase) ? "64px" : "0",
          top: isLibraryPhase ? "0" : (isChoiceOrTransitioning ? "calc(100% - 384px)" : "0"),
          pointerEvents: (isChoiceOrTransitioning && !isLibraryPhase) ? "none" : "auto",
        }}
      />

      {(isPlaying || isResult || isLibraryPhase) && (
        <button 
          onClick={isLibraryPhase ? handleBackToChoice : handleExit} 
          className="absolute top-6 right-6 z-[110] p-2 border border-black/30 dark:border-white/30 bg-transparent hover:bg-black/10 dark:hover:bg-white/10 transition-colors pointer-events-auto"
        >
          {isLibraryPhase ? <ChevronLeft className="w-5 h-5 text-black dark:text-white" /> : <X className="w-5 h-5 text-black dark:text-white" />}
        </button>
      )}

      {isLibraryPhase && (
        <div 
          className="absolute inset-0 overflow-y-auto pointer-events-auto"
          style={{
            paddingTop: "80px",
            paddingBottom: "100px",
          }}
        >
          <div className="px-6 max-w-screen-lg mx-auto">
            <div className="mb-8">
              <h1 className="text-3xl font-light mb-2 font-serif text-black dark:text-white">
                <InkRevealText text="Echo" />
              </h1>
              <p className="text-sm opacity-60 font-light font-sans text-black dark:text-white">
                <InkRevealText text="Soul resonance & healing frequencies" />
              </p>
            </div>

            <section className="space-y-3">
              {echoLibrary.map((audio) => (
                <AudioPlayer 
                  key={audio.id} 
                  {...audio} 
                  playCount={playCountMap[audio.id] || 0}
                  onPlayComplete={() => handlePlayComplete(audio.id)}
                />
              ))}
            </section>
          </div>
        </div>
      )}

      {!isLibraryPhase && (
        <div 
          className="absolute inset-0 pointer-events-none"
          style={{
            display: 'flex', 
            flexDirection: 'column', 
            alignItems: 'center', 
            justifyContent: 'center',
            transform: isChoiceOrTransitioning ? "translateY(calc(50vh - 224px - 40px))" : "translateY(0)", 
            opacity: (animationState === 'visible' || !isChoiceOrTransitioning) ? 1 : 0,
            filter: (animationState === 'visible' || !isChoiceOrTransitioning) ? 'blur(0px)' : 'blur(10px)',
            transition: "transform 1.2s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.8s ease-out, filter 0.8s ease-out",
            willChange: "transform, opacity, filter",
          }}
        >
          <div 
            className="relative pointer-events-auto"
            style={{
              transition: "transform 1.2s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.8s ease-out",
              opacity: phase === "exiting" ? 0 : 1,
              transform: isChoiceOrTransitioning ? "scale(0.5)" : "scale(1)",
              zIndex: 20,
              width: isChoiceOrTransitioning ? '200px' : '100%',
              height: isChoiceOrTransitioning ? '200px' : '100%',
              maxWidth: isChoiceOrTransitioning ? '200px' : '400px',
              maxHeight: isChoiceOrTransitioning ? '200px' : '400px',
            }}
          >
            <EchoEffect isPlaying={isPlaying || isChoiceOrTransitioning} />
          </div>

          {(isChoicePhase || phase === "transitioning") && (
            <div 
              className="absolute w-full flex flex-col items-center pointer-events-none px-6" 
              style={{ top: "calc(50% + 60px)", opacity: animationState === 'visible' ? 1 : 0 }}
            >
              <h3 className="text-lg font-light tracking-wide mb-4 font-serif text-black dark:text-white">
                <InkRevealText text="是否进入 Echo 调频" />
              </h3>
              
              {phase === "transitioning" ? (
                <div className="flex gap-4 pointer-events-auto">
                  <button 
                    disabled
                    className="relative w-48 py-3 text-sm font-light border border-black dark:border-white bg-white dark:bg-black text-black dark:text-white cursor-wait overflow-hidden"
                  >
                    <div 
                      className="absolute inset-0 z-0 laser-gradient"
                      style={{
                        width: `${generateProgress}%`,
                        transition: 'width 0.1s ease-out',
                      }}
                    />
                    <span className="relative z-10 mix-blend-difference text-white dark:text-black">
                      Generating - {generateProgress}%
                    </span>
                  </button>
                  <button 
                    onClick={handleCancel}
                    className="w-16 py-3 text-sm font-light border border-black dark:border-white bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-colors"
                  >
                    取消
                  </button>
                </div>
              ) : (
                <div className="flex gap-4 pointer-events-auto">
                  <button 
                    onClick={handleAccept} 
                    className="w-32 py-3 text-sm font-light border border-black dark:border-white bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-colors shadow-sm"
                  >
                    接受
                  </button>
                  <button 
                    onClick={handleReject} 
                    className="w-32 py-3 text-sm font-light border border-black dark:border-white bg-black dark:bg-white text-white dark:text-black hover:bg-white dark:hover:bg-black hover:text-black dark:hover:text-white transition-colors shadow-sm"
                  >
                    拒绝
                  </button>
                </div>
              )}
            </div>
          )}

        </div>
      )}

      {isPlaying && (
        <div className="fixed bottom-32 left-0 right-0 flex flex-col items-center pointer-events-none z-[110] px-12">
          <div className="w-full max-w-md">
            <div className="h-[2px] bg-black/10 dark:bg-white/10 rounded-full overflow-hidden">
              <div 
                className="h-full bg-black/40 dark:bg-white/40 transition-all duration-100"
                style={{ width: `${playProgress * 100}%` }}
              />
            </div>
            <div className="flex justify-between mt-2 text-xs text-black/40 dark:text-white/40 font-mono">
              <span>{Math.floor(playProgress * PLAY_DURATION)}s</span>
              <span>{PLAY_DURATION}s</span>
            </div>
          </div>
          <p className="text-xs text-black/50 dark:text-white/50 mt-4 tracking-widest font-sans uppercase">
            {trackName} · {frequency}Hz
          </p>
        </div>
      )}

      {(isChoicePhase || phase === "transitioning") && (
        <div 
          className="fixed left-0 right-0 flex justify-center z-[110]"
          style={{ 
            bottom: "79px",
            opacity: animationState === 'visible' ? 1 : 0,
            transition: "opacity 0.8s ease-out",
          }}
        >
          <button 
            onClick={phase === "transitioning" ? handleViewLater : handleShowLibrary}
            className="flex items-center gap-2 px-6 py-2 text-xs text-black/60 dark:text-white/60 hover:text-black dark:hover:text-white transition-colors pointer-events-auto font-sans"
          >
            <Music className="w-4 h-4" />
            <span>{phase === "transitioning" ? "稍后在 Echo 中查看" : "查看所有历史 Echo"}</span>
          </button>
        </div>
      )}

      {isResult && (
        <div 
          className="fixed bottom-32 left-0 right-0 flex flex-col items-center pointer-events-none z-[110] px-6"
          style={{
            opacity: 1,
            transition: "opacity 0.8s ease-out",
          }}
        >
          <InkRevealText text="磁场调整完成" className="text-xl font-light mb-4 font-serif text-black dark:text-white" />
          <p className="text-xs text-black/60 dark:text-white/60 max-w-md text-center leading-relaxed mb-6 font-sans">
            {frequency}Hz 频率已帮助你调整当前能量场，建议保持平静状态片刻
          </p>
          <button 
            onClick={handleComplete} 
            className="px-12 py-3 text-sm font-light border border-black dark:border-white bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-colors pointer-events-auto"
          >
            Continue
          </button>
        </div>
      )}
    </div>
  )
}

