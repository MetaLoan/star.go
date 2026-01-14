import { useEffect, useState } from "react"
import { OracleEyeIcon } from "./OracleEyeIcon"

interface ChatMessageProps {
  content: string
  isUser: boolean
  timestamp?: string
  contextHint?: string
}

export function ChatMessage({ content, isUser, timestamp, contextHint }: ChatMessageProps) {
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    // 延迟一帧后开始动画，确保初始状态已渲染
    const timer = requestAnimationFrame(() => {
      setIsVisible(true)
    })
    return () => cancelAnimationFrame(timer)
  }, [])

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"} mb-6 gap-3`}>
      {/* Oracle 头像 - 只在非用户消息时显示，无动画保持稳定 */}
      {!isUser && (
        <div className="shrink-0">
          <OracleEyeIcon className="w-10 h-10" />
        </div>
      )}

      <div className={`max-w-[85%] ${isUser ? "text-right" : "text-left"}`}>
        {/* Context hint for Oracle messages */}
        {!isUser && contextHint && (
          <div 
            className="text-[10px] opacity-30 mb-2 tracking-wider  transition-all duration-500 ml-1"
            style={{
              opacity: isVisible ? 0.3 : 0,
              transform: isVisible ? "translateY(0)" : "translateY(8px)",
            }}
          >
            {contextHint}
          </div>
        )}

        {/* Message bubble - 整体淡入动画 */}
        <div
          className={`inline-block border px-4 py-3 transition-all duration-700 ease-out ${
            isUser ? "border-black bg-black text-white dark:border-white dark:bg-white dark:text-black" : "border-black bg-white text-black dark:border-white dark:bg-black dark:text-white"
          }`}
          style={{
            opacity: isVisible ? 1 : 0,
            transform: isVisible ? "translateY(0)" : "translateY(12px)",
          }}
        >
          <p className="text-sm leading-relaxed font-light whitespace-pre-wrap break-words font-sans">
            {content}
          </p>
        </div>

        {/* Timestamp */}
        {timestamp && (
          <div 
            className="text-[10px] opacity-20 mt-2 ml-1 transition-all duration-500 delay-300"
            style={{
              opacity: isVisible ? 0.2 : 0,
            }}
          >
            {timestamp}
          </div>
        )}
      </div>
    </div>
  )
}

