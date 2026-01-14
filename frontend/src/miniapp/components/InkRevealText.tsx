import React, { useEffect, useState } from "react"
import { cn } from "../../utils/cn"

interface InkRevealTextProps {
  text: string
  className?: string
  delay?: number
}

export const InkRevealText: React.FC<InkRevealTextProps> = ({ text, className, delay = 0 }) => {
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), delay * 1000)
    return () => clearTimeout(timer)
  }, [delay])

  return (
    <span
      className={cn(
        "inline-block",
        isVisible ? "ink-effect" : "opacity-0",
        className
      )}
    >
      {text}
    </span>
  )
}

