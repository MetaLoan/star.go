// 意图识别系统
import type { FunctionName, IntentDetection } from "./oracle-types"

// 意图识别关键词映射
const intentKeywords: Record<FunctionName, string[]> = {
  stone: ["石头", "stone", "命运石", "fate stone", "占卜", "询问", "yes or no", "yes/no", "是否", "能不能", "应该"],
  tarot: ["塔罗", "tarot", "抽牌", "牌", "卡牌", "占卜", "预测", "塔罗牌"],
  echo: ["音频", "echo", "声音", "频率", "显化", "manifestation", "录音", "播放"],
  connect: ["连接", "connect", "关系", "缘分", "共振", "resonance", "分析", "匹配"]
}

/**
 * 检测用户输入中的意图
 */
export function detectIntent(text: string): IntentDetection {
  const lowerText = text.toLowerCase().trim()
  
  // 检查直接命令（如 /stone）
  if (lowerText.startsWith("/")) {
    const command = lowerText.slice(1).split(" ")[0] as FunctionName
    if (["stone", "tarot", "echo", "connect"].includes(command)) {
      return { intent: command, confidence: 1.0 }
    }
  }
  
  // 关键词匹配
  for (const [functionName, keywords] of Object.entries(intentKeywords)) {
    if (keywords.some(keyword => lowerText.includes(keyword))) {
      return { 
        intent: functionName as FunctionName, 
        confidence: 0.8 
      }
    }
  }
  
  return { intent: null, confidence: 0 }
}

/**
 * 生成Oracle对功能触发的响应
 */
export function generateOracleResponse(functionName: FunctionName, _userText: string): string {
  const responses: Record<FunctionName, string[]> = {
    stone: [
      "The Fate Stone awaits your question. Hold it close and let your intention flow through your touch.",
      "I sense you seek guidance. The stone will reveal the answer when you are ready.",
      "Place your question in the stone's embrace. It will respond when your focus is true."
    ],
    tarot: [
      "The cards are ready. Let your intuition guide you to the one that calls.",
      "The Tarot deck whispers. Draw a card and let it reveal what the universe wants you to know.",
      "The cards await your touch. Choose the one that resonates with your energy."
    ],
    echo: [
      "The frequencies are aligned. Let the sound guide your intention.",
      "I hear your call. The Echo chamber is ready to amplify your voice.",
      "The resonance chamber awaits. Speak your truth or let the frequencies carry you."
    ],
    connect: [
      "Two souls, one resonance. Let me analyze the connection between you.",
      "I see the threads that bind. The resonance portal is opening.",
      "The magnetic field between you is strong. Let me reveal its nature."
    ]
  }
  
  const options = responses[functionName]
  return options[Math.floor(Math.random() * options.length)]
}

/**
 * 生成Oracle对功能结果的总结
 */
export function generateResultSummary(functionName: FunctionName, result: any): string {
  switch (functionName) {
    case "stone":
      const stoneResult = result as { result: string; message: string }
      return `The stone has spoken: "${stoneResult.result}". ${stoneResult.message} This answer resonates with the energy you've been carrying.`
    
    case "tarot":
      return `The card reveals its wisdom. This guidance aligns with your current path.`
    
    case "echo":
      return `The frequencies have been aligned. Your intention is now amplified.`
    
    case "connect":
      return `The resonance analysis is complete. The connection between you has been revealed.`
    
    default:
      return "The reading is complete. May this guidance serve you well."
  }
}

