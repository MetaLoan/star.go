import { motion } from "framer-motion"

interface OracleThinkingProps {
  className?: string
}

export function OracleThinking({ className = "" }: OracleThinkingProps) {
  return (
    <div className={`flex justify-start mb-6 ${className}`}>
      <div className="flex gap-1 items-center px-4 py-3 border border-black/10 bg-black/[0.02] dark:border-white/10 dark:bg-white/[0.02]">
        <motion.div
          animate={{ opacity: [0.2, 1, 0.2] }}
          transition={{ duration: 1.5, repeat: Number.POSITIVE_INFINITY, delay: 0 }}
          className="w-1.5 h-1.5 bg-black/40 dark:bg-white/40 rounded-full"
        />
        <motion.div
          animate={{ opacity: [0.2, 1, 0.2] }}
          transition={{ duration: 1.5, repeat: Number.POSITIVE_INFINITY, delay: 0.3 }}
          className="w-1.5 h-1.5 bg-black/40 dark:bg-white/40 rounded-full"
        />
        <motion.div
          animate={{ opacity: [0.2, 1, 0.2] }}
          transition={{ duration: 1.5, repeat: Number.POSITIVE_INFINITY, delay: 0.6 }}
          className="w-1.5 h-1.5 bg-black/40 dark:bg-white/40 rounded-full"
        />
      </div>
    </div>
  )
}

