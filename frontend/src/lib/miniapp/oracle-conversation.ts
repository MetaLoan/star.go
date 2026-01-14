// Oracle 对话状态管理

export type ConversationState = 
  | "normal"              // 正常对话
  | "stone_triggered"     // 用户触发了STONE
  | "tarot_triggered"     // 用户触发了TAROT
  | "echo_triggered"      // 用户触发了ECHO
  | "connect_triggered"   // 用户触发了CONNECT
  | "waiting_for_question" // 等待用户描述疑惑
  | "question_received"   // 已收到疑惑
  | "suggesting_stone"    // 正在建议石头
  | "suggesting_tarot"    // 正在建议塔罗
  | "suggesting_echo"     // 正在建议Echo
  | "suggesting_connect"  // 正在建议Connect
  | "waiting_for_choice"  // 等待用户选择接受/拒绝
  | "stone_fullscreen"    // 全屏石头交互中
  | "tarot_fullscreen"    // 全屏塔罗交互中
  | "echo_fullscreen"     // 全屏Echo交互中
  | "connect_fullscreen"  // 全屏Connect交互中
  | "stone_completed"     // 石头交互完成
  | "tarot_completed"     // 塔罗交互完成
  | "echo_completed"      // Echo交互完成
  | "connect_completed"   // Connect交互完成
  | "waiting_for_suggestion_confirmation" // 等待用户确认建议

export interface ConversationContext {
  state: ConversationState
  userQuestion?: string  // 用户的疑惑
  stoneResult?: {
    result: string
    message: string
    choice?: "follow" | "rebel"
  }
}

// 生成Oracle的深度解读
export function generateDeepInterpretation(
  result: { result: string; message: string; choice?: "follow" | "rebel" },
  userQuestion?: string
): string {
  const choiceAnalysis = result.choice === "follow" 
    ? "你选择了遵循石头的指引，这显示了你内心的开放和信任。"
    : "你选择了拒绝石头的答案，这反映了你内心的独立和坚持。"
  
  const resultMeanings: Record<string, string> = {
    "YES": "肯定的答案通常意味着时机已到，但你的选择揭示了更深层的矛盾。",
    "NO": "否定的答案可能是保护，但你的反应说明你并不完全接受这个结果。",
    "WAIT": "等待的提示需要耐心，但你的选择显示你更倾向于主动行动。",
    "SILENCE": "沉默是答案，但你的内心选择暴露了你对不确定性的焦虑。",
    "RELEASE": "放手的建议需要勇气，但你的选择表明你还在紧紧抓住某些东西。"
  }
  
  const resultMeaning = resultMeanings[result.result] || "这个答案有其深意。"
  
  const questionContext = userQuestion 
    ? `关于"${userQuestion.slice(0, 30)}${userQuestion.length > 30 ? '...' : ''}"这个问题，`
    : ""
  
  return `${questionContext}石头告诉你：${result.message}。${resultMeaning}${choiceAnalysis}这也侧面反映了你当前内心的挣扎和需要面对的挑战。`
}

