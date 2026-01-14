import { useEffect, useRef, useCallback } from "react"

interface EchoEffectProps {
  isPlaying?: boolean
  className?: string
}

export function EchoEffect({ isPlaying = true, className = "" }: EchoEffectProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const animationRef = useRef<number | null>(null)

  const initEffect = useCallback(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    // 设置画布尺寸
    const updateSize = () => {
      const container = canvas.parentElement
      if (container) {
        canvas.width = container.clientWidth
        canvas.height = container.clientHeight
      }
    }
    updateSize()

    // 圆环配置 - 减少粒子数量，缩小半径
    const config = {
      circleCount: 7,           // 减少圆环数量
      radius: Math.min(canvas.width, canvas.height) * 0.28,  // 缩小半径
      echo: 30,                 // 减少波动幅度
      speed: 0.5,
      segments: 180             // 减少每个圆环的粒子数
    }

    // 初始化圆环
    class Circle {
      r: number
      expanding: boolean
      max: number
      min: number
      val: number

      constructor(index: number) {
        this.r = config.radius - index * config.radius / config.circleCount
        this.expanding = index % 2 === 1
        this.max = Math.random() * config.echo
        this.min = -Math.random() * config.echo
        this.val = Math.random() * (this.max - this.min) + this.min
      }
    }

    const circles: Circle[] = []
    for (let i = 0; i < config.circleCount; i++) {
      circles.push(new Circle(i))
    }

    // 绘制函数
    const draw = () => {
      if (!canvas || !ctx) return

      // 清空画布 - 完全透明背景
      ctx.clearRect(0, 0, canvas.width, canvas.height)

      const centerX = canvas.width / 2
      const centerY = canvas.height / 2

      // 绘制每个圆环 - 纯黑粒子
      circles.forEach((circle) => {
        for (let i = 0; i < config.segments; i++) {
          const angle = i * Math.PI * 2 / config.segments
          const distortion = circle.val * Math.cos(i / 2)
          const x = Math.cos(angle) * (circle.r - distortion)
          const y = Math.sin(angle) * (circle.r - distortion)

          ctx.fillStyle = 'rgba(0, 0, 0, 0.8)'  // 纯黑粒子
          ctx.fillRect(centerX + x, centerY - y, 1.2, 1.2)
        }

        // 更新圆环状态
        if (isPlaying) {
          circle.val = circle.expanding ? circle.val + config.speed : circle.val - config.speed
          
          if (circle.val < circle.min) {
            circle.expanding = true
            circle.max = Math.random() * config.echo
          }
          if (circle.val > circle.max) {
            circle.expanding = false
            circle.min = -Math.random() * config.echo
          }
        }
      })

      animationRef.current = requestAnimationFrame(draw)
    }

    draw()

    // 监听窗口大小变化
    const handleResize = () => {
      updateSize()
      config.radius = Math.min(canvas.width, canvas.height) * 0.28
    }
    window.addEventListener('resize', handleResize)

    return () => {
      window.removeEventListener('resize', handleResize)
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current)
      }
    }
  }, [isPlaying])

  useEffect(() => {
    const cleanup = initEffect()
    return cleanup
  }, [initEffect])

  return (
    <canvas
      ref={canvasRef}
      className={`w-full h-full ${className}`}
      style={{ background: 'transparent' }}
    />
  )
}

