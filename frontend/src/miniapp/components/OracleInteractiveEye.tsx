import { useEffect, useState } from "react"

interface OracleInteractiveEyeProps {
  direction: "up" | "down" | "down-right" | "down-left" | "center"
  className?: string
  inverted?: boolean
}

export function OracleInteractiveEye({ direction, className = "", inverted = false }: OracleInteractiveEyeProps) {
  const [eyeOffset, setEyeOffset] = useState({ x: 0, y: 0 })
  const color = inverted ? "white" : "black"
  const innerColor = inverted ? "black" : "white"

  // 定义尖锐眼眶的路径
  const eyePath = "M 10 100 Q 100 45 190 100 Q 100 155 10 100 Z"

  useEffect(() => {
    // 根据direction计算眼球偏移
    const maxMove = 25
    let targetX = 0
    let targetY = 0

    switch (direction) {
      case "up":
        targetX = 0
        targetY = -maxMove // 向上看
        break
      case "down":
        targetX = 0
        targetY = maxMove // 向下看
        break
      case "down-right":
        targetX = maxMove * 0.7 // 向右下角
        targetY = maxMove * 0.7
        break
      case "down-left":
        targetX = -maxMove * 0.7 // 向左下角
        targetY = maxMove * 0.7
        break
      case "center":
      default:
        targetX = 0
        targetY = 0
        break
    }

    setEyeOffset({ x: targetX, y: targetY })
  }, [direction])

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <svg width="100%" height="100%" viewBox="0 0 200 200">
        <defs>
          {/* 尖锐眼眶的裁剪区域 */}
          <clipPath id="interactiveEyeClip">
            <path d={eyePath} />
          </clipPath>
        </defs>

        {/* 装饰性外部线条 */}
        <path d="M 0 100 L 30 100 M 170 100 L 200 100" stroke={color} strokeWidth="1" opacity="0.3" />

        {/* 主眼眶轮廓 */}
        <path d={eyePath} fill="none" stroke={color} strokeWidth="2.5" />

        {/* 内部细轮廓线 */}
        <path d="M 18 100 Q 100 52 182 100 Q 100 148 18 100 Z" fill="none" stroke={color} strokeWidth="1" opacity="0.4" />

        {/* 可移动的圆眼球 */}
        <g clipPath="url(#interactiveEyeClip)">
          <g
            style={{
              transform: `translate(${eyeOffset.x}px, ${eyeOffset.y}px)`,
              transition: "transform 0.3s ease-out",
            }}
          >
            {/* 虹膜/眼球主体 - 填充颜色 */}
            <circle cx="100" cy="100" r="42" fill={color} stroke={color} strokeWidth="1.5" />
            <circle cx="100" cy="100" r="38" fill="none" stroke={innerColor} strokeWidth="0.3" strokeDasharray="2 2" opacity="0.5" />

            {/* 放射状精密线条 - 改为内色以在白色眼球上显示 */}
            {[...Array(48)].map((_, i) => (
              <line
                key={i}
                x1="100"
                y1="100"
                x2={Math.round((100 + Math.cos((i * 7.5 * Math.PI) / 180) * 42) * 100) / 100}
                y2={Math.round((100 + Math.sin((i * 7.5 * Math.PI) / 180) * 42) * 100) / 100}
                stroke={innerColor}
                strokeWidth="0.3"
                opacity="0.15"
              />
            ))}

            {/* 核心瞳孔 - 改为内色(黑色) */}
            <ellipse cx="100" cy="100" rx="8" ry="28" fill={innerColor} />

            {/* 瞳孔内部的细缝 - 改为外色(白色) */}
            <ellipse cx="100" cy="100" rx="1.5" ry="20" fill={color} />
          </g>
        </g>

        {/* 眼角排线装饰 */}
        <g stroke={color} strokeWidth="1" opacity="0.2">
          <line x1="20" y1="95" x2="35" y2="85" />
          <line x1="20" y1="105" x2="35" y2="115" />
          <line x1="180" y1="95" x2="165" y2="85" />
          <line x1="180" y1="105" x2="165" y2="115" />
        </g>
      </svg>
    </div>
  )
}

