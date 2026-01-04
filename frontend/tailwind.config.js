import { heroui } from "@heroui/react";

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 宇宙主题色彩系统
        cosmic: {
          nova: '#00D4FF',      // 主色 - 星光蓝
          aurora: '#FF6B9D',    // 强调色 - 极光粉
          nebula: '#A855F7',    // 神秘色 - 星云紫
          void: '#0A0A0F',      // 背景色 - 虚空黑
          dust: '#1a1a2e',      // 次背景 - 星尘蓝
        },
        // 行星颜色
        planet: {
          sun: '#ffd700',
          moon: '#c0c0c0',
          mercury: '#b5651d',
          venus: '#ff69b4',
          mars: '#dc143c',
          jupiter: '#daa520',
          saturn: '#8b7355',
          uranus: '#40e0d0',
          neptune: '#4169e1',
          pluto: '#800080',
        },
        // 相位颜色
        aspect: {
          conjunction: '#ffd700',
          trine: '#22c55e',
          sextile: '#3b82f6',
          square: '#ef4444',
          opposition: '#f97316',
        },
        // 维度颜色
        dimension: {
          career: '#3b82f6',
          relationship: '#ec4899',
          health: '#22c55e',
          finance: '#f59e0b',
          spiritual: '#a855f7',
        },
      },
      fontFamily: {
        display: ['Space Grotesk', 'system-ui', 'sans-serif'],
        body: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      backgroundImage: {
        'cosmic-gradient': 'radial-gradient(ellipse at top, #1a1a2e 0%, #0a0a0f 100%)',
        'aurora-gradient': 'linear-gradient(135deg, #00D4FF 0%, #A855F7 50%, #FF6B9D 100%)',
      },
      animation: {
        'glow': 'glow 2s ease-in-out infinite alternate',
        'float': 'float 6s ease-in-out infinite',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        glow: {
          '0%': { boxShadow: '0 0 5px rgba(0, 212, 255, 0.3)' },
          '100%': { boxShadow: '0 0 20px rgba(0, 212, 255, 0.6)' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        },
      },
    },
  },
  darkMode: "class",
  plugins: [heroui()],
}

