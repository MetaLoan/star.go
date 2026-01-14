import { FateStoneEmbedded } from "./FateStoneEmbedded"
import type { FunctionName, StoneResult } from "../../lib/miniapp/oracle-types"

interface FunctionEmbedProps {
  functionName: FunctionName
  onComplete: (functionName: FunctionName, result: any) => void
  onCancel?: () => void
}

export function FunctionEmbed({ functionName, onComplete, onCancel }: FunctionEmbedProps) {
  const handleStoneComplete = (result: StoneResult) => {
    onComplete("stone", result)
  }

  switch (functionName) {
    case "stone":
      return (
        <FateStoneEmbedded 
          onComplete={handleStoneComplete}
          onCancel={onCancel}
          embedded={true}
        />
      )
    
    case "tarot":
      return <div className="text-xs opacity-40 italic py-4">Tarot coming soon...</div>
    
    case "echo":
      return <div className="text-xs opacity-40 italic py-4">Echo coming soon...</div>
    
    case "connect":
      return <div className="text-xs opacity-40 italic py-4">Connect coming soon...</div>
    
    default:
      return null
  }
}

