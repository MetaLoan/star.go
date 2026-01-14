import { useState, useEffect, useRef } from "react"
import { cn } from "../../utils/cn"

interface TextRevealProps {
  text: string
  className?: string
  mode: "fadein" | "fadeaway" | "idle"
  delayPerChar?: number // 每个字符的延迟（毫秒）
  startDelay?: number // 动画开始前的延迟（毫秒）
  onComplete?: () => void
}

export function TextReveal({ 
  text, 
  className = "", 
  mode, 
  delayPerChar = 50,
  startDelay = 0,
  onComplete 
}: TextRevealProps) {
  // 使用 ref 保存 onComplete 避免依赖变化导致无限循环
  const onCompleteRef = useRef(onComplete)
  onCompleteRef.current = onComplete
  
  // 追踪是否已经完成过动画，避免重复触发
  const hasCompletedRef = useRef(false)
  
  // 追踪动画是否已开始（用于 startDelay）
  const [hasStarted, setHasStarted] = useState(startDelay === 0)
  
  // fadein 模式：使用和 InkRevealText 一样的计数器方式
  const [revealedCount, setRevealedCount] = useState(mode === "fadein" ? 0 : text.length)
  
  // fadeaway 模式：使用状态数组
  const [charStates, setCharStates] = useState<("visible" | "fading" | "gone")[]>(
    text.split("").map(() => "visible" as const)
  )

  // startDelay 处理
  useEffect(() => {
    if (startDelay > 0 && mode === "fadein") {
      setHasStarted(false)
      const timer = setTimeout(() => {
        setHasStarted(true)
      }, startDelay)
      return () => clearTimeout(timer)
    } else {
      setHasStarted(true)
    }
  }, [startDelay, mode, text])

  // fadein 模式
  useEffect(() => {
    if (mode !== "fadein") return
    
    hasCompletedRef.current = false
    setRevealedCount(0)
  }, [mode, text])
  
  useEffect(() => {
    if (mode !== "fadein" || !hasStarted) return
    
    if (revealedCount < text.length) {
      const timer = setTimeout(() => {
        setRevealedCount(prev => prev + 1)
      }, delayPerChar)
      return () => clearTimeout(timer)
    } else if (revealedCount === text.length && !hasCompletedRef.current) {
      hasCompletedRef.current = true
      onCompleteRef.current?.()
    }
  }, [revealedCount, text.length, delayPerChar, mode, hasStarted])

  // fadeaway 模式
  useEffect(() => {
    if (mode !== "fadeaway") return
    
    hasCompletedRef.current = false
    setCharStates(text.split("").map(() => "visible" as const))
    
    let currentIndex = 0
    const interval = window.setInterval(() => {
      setCharStates(prev => {
        const newStates = [...prev]
        
        if (currentIndex < text.length) {
          newStates[currentIndex] = "fading"
          if (currentIndex > 0) {
            newStates[currentIndex - 1] = "gone"
          }
        }
        
        return newStates
      })
      
      currentIndex++
      
      if (currentIndex > text.length) {
        setCharStates(prev => prev.map(() => "gone" as const))
        clearInterval(interval)
        if (!hasCompletedRef.current) {
          hasCompletedRef.current = true
          onCompleteRef.current?.()
        }
      }
    }, delayPerChar)
    
    return () => clearInterval(interval)
  }, [mode, text, delayPerChar])

  // idle 模式
  useEffect(() => {
    if (mode === "idle") {
      setRevealedCount(text.length)
      setCharStates(text.split("").map(() => "visible" as const))
    }
  }, [mode, text])

  if (mode === "fadein") {
    return (
      <span className={cn("inline-block", className)}>
        {text.split("").map((char, i) => (
          <span
            key={i}
            className={cn("inline-block", i < revealedCount && "ink-effect")}
            style={{
              opacity: i < revealedCount ? 1 : 0,
            }}
          >
            {char === " " ? "\u00A0" : char}
          </span>
        ))}
      </span>
    )
  }

  return (
    <span className={className}>
      {text.split("").map((char, index) => {
        const state = charStates[index]
        
        let style: React.CSSProperties = {
          display: "inline-block",
          whiteSpace: char === " " ? "pre" : "normal",
        }
        
        if (state === "fading") {
          style = {
            ...style,
            opacity: 0.3,
            filter: "blur(4px)",
            transform: "translateY(-3px)",
            transition: "all 0.15s ease-out"
          }
        } else if (state === "gone") {
          style = {
            ...style,
            opacity: 0,
            filter: "blur(8px)",
            transform: "translateY(-8px)",
            transition: "all 0.15s ease-out"
          }
        } else {
          style = {
            ...style,
            opacity: 1,
            filter: "blur(0px)",
            transform: "scale(1)"
          }
        }
        
        return (
          <span
            key={index}
            style={style}
          >
            {char === " " ? "\u00A0" : char}
          </span>
        )
      })}
    </span>
  )
}

