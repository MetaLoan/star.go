import { useEffect, useState } from "react"
import { ChatMessage } from "./ChatMessage"
import { FunctionEmbed } from "./FunctionEmbed"
import { AudioPlayer } from "./AudioPlayer"
import type { Message } from "../../lib/miniapp/oracle-types"

interface MessageRendererProps {
  message: Message
  onFunctionComplete?: (functionName: string, result: any) => void
  onFunctionCancel?: () => void
  onChoiceAccept?: () => void
  onChoiceReject?: () => void
}

export function MessageRenderer({ message, onFunctionComplete, onFunctionCancel }: MessageRendererProps) {
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    const timer = requestAnimationFrame(() => {
      setIsVisible(true)
    })
    return () => cancelAnimationFrame(timer)
  }, [])

  if (message.type === "text") {
    return (
      <ChatMessage 
        content={message.content}
        isUser={message.isUser}
        timestamp={message.timestamp}
        contextHint={message.contextHint}
      />
    )
  }

  if (message.type === "function-trigger") {
    return (
      <div className="space-y-4">
        {message.oracleResponse && (
          <ChatMessage 
            content={message.oracleResponse}
            isUser={false}
            timestamp={message.timestamp}
            contextHint="Function activated"
          />
        )}
      </div>
    )
  }

  if (message.type === "function-embed") {
    return (
      <div className="space-y-4 my-6">
        {message.oracleResponse && (
          <ChatMessage 
            content={message.oracleResponse}
            isUser={false}
            timestamp={message.timestamp}
            contextHint="Active"
          />
        )}
        
        <div 
          className="border border-black/10 dark:border-white/10 bg-black/[0.02] dark:bg-white/[0.02] p-6 rounded-sm transition-all duration-700"
          style={{
            opacity: isVisible ? 1 : 0,
            transform: isVisible ? "translateY(0)" : "translateY(20px)",
          }}
        >
          <FunctionEmbed
            functionName={message.functionName}
            onComplete={(functionName, result) => {
              onFunctionComplete?.(functionName, result)
            }}
            onCancel={onFunctionCancel}
          />
        </div>
      </div>
    )
  }

  if (message.type === "function-result") {
    return (
      <ChatMessage 
        content={message.summary}
        isUser={false}
        timestamp={message.timestamp}
        contextHint="Result"
      />
    )
  }

  if (message.type === "system") {
    return (
      <div className="flex flex-col items-center justify-center my-8 space-y-2">
        <p className="text-[10px] opacity-40 font-light tracking-widest text-center max-w-[80%] uppercase font-sans">
          {message.content}
        </p>
        {message.timestamp && (
          <span className="text-[9px] opacity-20 uppercase tracking-widest font-sans">{message.timestamp}</span>
        )}
      </div>
    )
  }

  if (message.type === "echo-audio") {
    return (
      <div 
        className="transition-all duration-700"
        style={{
          opacity: isVisible ? 1 : 0,
          transform: isVisible ? "translateY(0)" : "translateY(20px)",
        }}
      >
        <AudioPlayer
          title={message.title}
          subtitle={message.subtitle}
          frequency={message.frequency}
          duration={message.duration}
          description={message.description}
          playCount={0}
        />
      </div>
    )
  }

  return null
}

