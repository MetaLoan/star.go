import { motion } from "framer-motion"

interface FunctionButtonsProps {
  onStoneClick: () => void
  onTarotClick: () => void
  onEchoClick: () => void
  onConnectClick: () => void
  disabled?: boolean
  echoGenerating?: boolean
  echoProgress?: number
}

export function FunctionButtons({
  onStoneClick,
  onTarotClick,
  onEchoClick,
  onConnectClick,
  disabled = false,
  echoGenerating = false,
  echoProgress = 0,
}: FunctionButtonsProps) {
  const buttons = [
    { id: "stone", label: "Stone", icon: "💎", onClick: onStoneClick },
    { id: "tarot", label: "Tarot", icon: "🎴", onClick: onTarotClick },
    { id: "echo", label: "Echo", icon: "🎵", onClick: onEchoClick },
    { id: "connect", label: "Connect", icon: "🔗", onClick: onConnectClick },
  ]

  return (
    <div className="flex flex-wrap gap-2 mb-4">
      {buttons.map((btn) => (
        <button
          key={btn.id}
          onClick={btn.onClick}
          disabled={disabled}
          className="relative px-4 py-2 border border-black dark:border-white text-xs font-light hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition-colors disabled:opacity-50 disabled:cursor-not-allowed group"
        >
          <span className="mr-2">{btn.icon}</span>
          {btn.label}
          
          {btn.id === "echo" && echoGenerating && (
            <div className="absolute inset-0 bg-black/10 dark:bg-white/10 pointer-events-none overflow-hidden">
              <motion.div 
                className="absolute inset-y-0 left-0 bg-black/20 dark:bg-white/20"
                initial={{ width: 0 }}
                animate={{ width: `${echoProgress}%` }}
              />
            </div>
          )}
        </button>
      ))}
    </div>
  )
}

