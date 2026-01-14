import { useState, useEffect, useRef, useMemo } from "react"
import { InkRevealText } from "./InkRevealText"

interface AudioPlayerProps {
  title: string
  subtitle: string
  frequency: string
  duration: string
  description?: string
  playCount?: number
  onPlayComplete?: () => void
}

export function AudioPlayer({ 
  title, 
  subtitle, 
  frequency, 
  duration, 
  description,
  playCount = 0,
  onPlayComplete
}: AudioPlayerProps) {
  const [isPlaying, setIsPlaying] = useState(false)
  const [progress, setProgress] = useState(0)
  const [isExpanded, setIsExpanded] = useState(false)
  const progressTimerRef = useRef<number | null>(null)

  // 解析时长为秒数
  const totalSeconds = useMemo(() => {
    const parts = duration.split(":")
    if (parts.length === 2) {
      return parseInt(parts[0]) * 60 + parseInt(parts[1])
    }
    return 0
  }, [duration])

  // 计算剩余时长
  const remainingTime = useMemo(() => {
    const remaining = Math.ceil(totalSeconds * (1 - progress / 100))
    const mins = Math.floor(remaining / 60)
    const secs = remaining % 60
    return `${mins}:${secs.toString().padStart(2, "0")}`
  }, [totalSeconds, progress])

  // 模拟播放进度
  useEffect(() => {
    if (isPlaying) {
      progressTimerRef.current = window.setInterval(() => {
        setProgress((prev) => {
          if (prev >= 100) {
            setIsPlaying(false)
            // 播放完成，触发回调
            if (onPlayComplete) {
              onPlayComplete()
            }
            return 0
          }
          return prev + 0.5
        })
      }, 100)
    } else {
      if (progressTimerRef.current) {
        window.clearInterval(progressTimerRef.current)
        progressTimerRef.current = null
      }
    }

    return () => {
      if (progressTimerRef.current) {
        window.clearInterval(progressTimerRef.current)
      }
    }
  }, [isPlaying, onPlayComplete])

  const togglePlay = () => {
    setIsPlaying((prev) => !prev)
  }

  const handleRowClick = () => {
    setIsExpanded((prev) => !prev)
  }

  const handlePlayClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()
    setIsExpanded(true)
    togglePlay()
  }

  const totalBars = 32

  return (
    <div
      className="border border-black dark:border-white bg-white dark:bg-black hover:bg-black/[0.02] dark:hover:bg-white/[0.02] transition-colors cursor-pointer select-none"
      onClick={handleRowClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault()
          handleRowClick()
        }
      }}
    >
      <div className="flex items-center gap-3 p-4">
        <div className="shrink-0 w-32 sm:w-40 overflow-hidden">
          <div className="text-sm font-normal whitespace-nowrap overflow-hidden">
            <span className="inline-block">
              <InkRevealText text={title} />
            </span>
          </div>
          <div className="text-xs font-light whitespace-nowrap overflow-hidden">
            <span className="inline-block opacity-60">{frequency}</span>
          </div>
          <div className="w-20 h-1 mt-2 bg-black/10 dark:bg-white/10 overflow-hidden">
            <div
              className="h-full laser-gradient transition-all duration-500 ease-out"
              style={{
                width: `${Math.min(playCount, 5) * 20}%`,
              }}
            />
          </div>
        </div>

        <div className="flex-1 min-w-0 relative h-4 overflow-hidden">
          <div
            className="absolute inset-0 flex items-center justify-center transition-all duration-500 ease-out"
            style={{
              transform: isPlaying ? "translateX(100%)" : "translateX(0)",
              opacity: isPlaying ? 0 : 1,
            }}
          >
            <span className="flex items-center gap-1.5 text-xs opacity-40">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                <circle cx="12" cy="12" r="10" />
                <path d="M12 6v6l4 2" />
              </svg>
              {duration}
            </span>
          </div>

          <div
            className="absolute inset-0 flex items-center gap-[1px] transition-all duration-500 ease-out"
            style={{
              transform: isPlaying ? "translateX(0)" : "translateX(-100%)",
              opacity: isPlaying ? 1 : 0,
            }}
          >
            {Array.from({ length: totalBars }).map((_, i) => {
              const height = Math.sin(i * 0.4) * 35 + 55
              const barProgress = (i / totalBars) * 100
              const isActive = barProgress < progress

              return (
                <div
                  key={i}
                  className={`flex-1 rounded-full transition-all duration-150 ${
                    isActive
                      ? "laser-gradient"
                      : "bg-black/20 dark:bg-white/20"
                  }`}
                  style={{ height: `${height}%` }}
                />
              )
            })}
          </div>
        </div>

        <div
          className="shrink-0 text-xs font-mono tabular-nums transition-all duration-500 ease-out overflow-hidden"
          style={{
            width: isPlaying ? "40px" : "0px",
            opacity: isPlaying ? 1 : 0,
            marginLeft: isPlaying ? "0" : "-12px",
          }}
        >
          {remainingTime}
        </div>

        <button
          onClick={handlePlayClick}
          aria-label={isPlaying ? "Pause" : "Play"}
          className={`w-8 h-8 border border-black dark:border-white flex items-center justify-center transition-all shrink-0 ${
            isPlaying
              ? "bg-black text-white dark:bg-white dark:text-black"
              : "bg-white text-black dark:bg-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black"
          }`}
        >
          {isPlaying ? (
            <svg width="10" height="10" viewBox="0 0 12 12" fill="currentColor">
              <rect width="4" height="12" />
              <rect x="8" width="4" height="12" />
            </svg>
          ) : (
            <svg width="10" height="10" viewBox="0 0 12 12" fill="currentColor">
              <polygon points="0,0 0,12 12,6" />
            </svg>
          )}
        </button>
      </div>

      <div
        className="grid transition-all duration-300 ease-out"
        style={{
          gridTemplateRows: isExpanded ? "1fr" : "0fr",
        }}
      >
        <div className="overflow-hidden">
          <div className="px-4 pb-4 pt-2">
            <p className="text-xs opacity-50 font-light leading-relaxed">
              {description || subtitle}
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

