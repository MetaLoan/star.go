export type FunctionName = "stone" | "tarot" | "echo" | "connect";

export interface IntentDetection {
  intent: FunctionName | null;
  confidence: number;
}

export interface StoneResult {
  result: string;
  message: string;
  choice?: "follow" | "rebel";
}

export interface TarotCardData {
  id: string;
  name: string;
  image: string;
  keywords: string[];
  meaning: string;
  uprightMeaning: string;
  reversedMeaning: string;
}

export interface TarotResult {
  cards: TarotCardData[];
  orientations: ('upright' | 'reversed')[];
  name: string;      // 牌面名称摘要
  meaning: string;   // 牌面含义摘要
  summary: string;   // 完整解读
  interpretation?: string;
}

export interface EchoResult {
  frequency: number;      // 频率 Hz
  duration: number;       // 播放时长（秒）
  trackName: string;      // 音轨名称
  purpose: string;        // 用途描述
  recommendation: string; // 推荐理由
}

export interface BaseMessage {
  id: number;
  timestamp: string;
  isUser: boolean;
}

export interface TextMessage extends BaseMessage {
  type: "text";
  content: string;
  contextHint?: string;
}

export interface FunctionTriggerMessage extends BaseMessage {
  type: "function-trigger";
  functionName: FunctionName;
  triggerText: string;
  oracleResponse?: string;
}

export interface FunctionResultMessage extends BaseMessage {
  type: "function-result";
  functionName: FunctionName;
  result: any;
  summary?: string;
}

export interface FunctionEmbedMessage extends BaseMessage {
  type: "function-embed";
  functionName: FunctionName;
  oracleResponse?: string;
}

export interface ChoiceMessage extends BaseMessage {
  type: "choice";
  prompt: string;
  choices: { label: string; value: string }[];
}

export interface SystemMessage extends BaseMessage {
  type: "system";
  content: string;
}

export interface EchoAudioMessage extends BaseMessage {
  type: "echo-audio";
  title: string;
  subtitle: string;
  frequency: string;
  duration: string;
  description?: string;
}

export type Message =
  | TextMessage
  | FunctionTriggerMessage
  | FunctionResultMessage
  | FunctionEmbedMessage
  | ChoiceMessage
  | SystemMessage
  | EchoAudioMessage;

