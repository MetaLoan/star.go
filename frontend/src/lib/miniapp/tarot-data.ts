import type { TarotCardData } from "./oracle-types";

export const tarotCards: TarotCardData[] = [
  {
    id: "00",
    name: "The Fool",
    image: "/card/00.jpg",
    keywords: ["Innocence", "New Beginnings", "Free Spirit"],
    meaning: "The Fool represents the beginning of a journey, a leap of faith into the unknown.",
    uprightMeaning: "新的开始即将到来，放下恐惧，勇敢迈出第一步。保持童真与好奇心，相信宇宙的引导。",
    reversedMeaning: "过于冲动或鲁貌，需要更谨慎地思考。可能在逃避责任或害怕承诺。"
  },
  {
    id: "01",
    name: "The Magician",
    image: "/card/01.jpg",
    keywords: ["Manifestation", "Resourcefulness", "Power"],
    meaning: "The Magician indicates that you have the tools and resources to manifest your desires.",
    uprightMeaning: "你拥有实现目标所需的一切资源。现在是将想法付诸行动的最佳时机，发挥你的创造力。",
    reversedMeaning: "才能被浪费或误用，可能存在欺骗或操控。需要审视自己的真实意图。"
  },
  {
    id: "02",
    name: "The High Priestess",
    image: "/card/02.jpg",
    keywords: ["Intuition", "Sacred Knowledge", "Divine Feminine"],
    meaning: "The High Priestess suggests a time to look inward and trust your intuition.",
    uprightMeaning: "倾听内心的声音，答案就在你的直觉中。保持神秘感，不必急于揭示一切。",
    reversedMeaning: "忽视直觉或过度依赖理性。可能有隐藏的信息尚未浮出水面。"
  }
];

export const getRandomCard = (): TarotCardData => {
  return tarotCards[Math.floor(Math.random() * tarotCards.length)];
};

