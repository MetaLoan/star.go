import { useEffect, useRef } from "react"

interface StarfieldEffectProps {
  className?: string
  speedMultiplier?: number
  colorMode?: "normal" | "rainbow" | "fading"
}

export function StarfieldEffect({ className, speedMultiplier = 1, colorMode = "normal" }: StarfieldEffectProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const animationRef = useRef<number>(0)
  const speedRef = useRef(speedMultiplier)
  const colorModeRef = useRef(colorMode)
  const fadeProgressRef = useRef(0)
  const starsRef = useRef<{x: number, y: number, z: number, hue: number}[]>([])
  
  useEffect(() => {
    speedRef.current = speedMultiplier
  }, [speedMultiplier])
  
  useEffect(() => {
    colorModeRef.current = colorMode
    if (colorMode === "fading") {
      fadeProgressRef.current = 0
    }
  }, [colorMode])

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const resizeCanvas = () => {
      if (canvas.parentElement) {
        canvas.width = canvas.parentElement.offsetWidth
        canvas.height = canvas.parentElement.offsetHeight
      }
    }
    resizeCanvas()

    const starCount = 400
    const baseSpeed = 2

    if (starsRef.current.length === 0) {
      for (let i = 0; i < starCount; i++) {
        starsRef.current.push({
          x: (Math.random() - 0.5) * canvas.width,
          y: (Math.random() - 0.5) * canvas.height,
          z: Math.random() * canvas.width,
          hue: Math.random() * 360
        })
      }
    }

    ctx.fillStyle = 'white'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    function draw() {
      const currentSpeed = baseSpeed * speedRef.current
      const currentColorMode = colorModeRef.current
      const isBoost = speedRef.current > 1.5
      
      let colorProgress = 0
      let fadeOpacity = 1
      if (currentColorMode === "rainbow") {
        colorProgress = Math.min(1, Math.max(0, (speedRef.current - 1) / 9))
      } else if (currentColorMode === "fading") {
        fadeProgressRef.current = Math.min(1, fadeProgressRef.current + 0.012)
        colorProgress = 1
        fadeOpacity = 1 - fadeProgressRef.current * 0.8
      }
      
      const trailAlpha = isBoost ? 0.12 : 0.3
      ctx!.fillStyle = `rgba(255, 255, 255, ${trailAlpha})`
      ctx!.fillRect(0, 0, canvas!.width, canvas!.height)

      for (const star of starsRef.current) {
        star.z -= currentSpeed
        if (currentColorMode === "rainbow" || (currentColorMode === "fading" && colorProgress > 0)) {
          star.hue = (star.hue + 0.3) % 360
        }

        if (star.z <= 0) {
          star.z = canvas!.width
          star.x = (Math.random() - 0.5) * canvas!.width
          star.y = (Math.random() - 0.5) * canvas!.height
          star.hue = Math.random() * 360
        }

        const k = 128 / star.z
        const x = star.x * k + canvas!.width / 2
        const y = star.y * k + canvas!.height / 2

        const size = Math.max(0, (1 - star.z / canvas!.width) * 2)
        const brightness = 1 - star.z / canvas!.width

        const shapeProgress = currentColorMode === "rainbow" 
          ? Math.min(1, Math.max(0, (speedRef.current - 1) / 2))
          : (currentColorMode === "fading" ? 1 : 0)
        
        if (shapeProgress > 0) {
          const trailOpacity = 0.7
          ctx!.strokeStyle = `rgba(0, 0, 0, ${trailOpacity * fadeOpacity * shapeProgress})`
          ctx!.lineWidth = Math.max(0.5, size * 0.8)
          ctx!.beginPath()
          ctx!.moveTo(x, y)
          ctx!.lineTo(
            star.x * (k + 0.1) + canvas!.width / 2,
            star.y * (k + 0.1) + canvas!.height / 2
          )
          ctx!.stroke()
        }
        
        if (shapeProgress < 1 && currentColorMode !== "fading") {
          const dotOpacity = brightness * (1 - shapeProgress) * fadeOpacity
          ctx!.fillStyle = `rgba(0, 0, 0, ${dotOpacity})`
          ctx!.beginPath()
          ctx!.arc(x, y, size, 0, Math.PI * 2)
          ctx!.fill()
        }
      }

      animationRef.current = window.requestAnimationFrame(draw)
    }

    draw()

    const handleResize = () => {
      resizeCanvas()
      ctx.fillStyle = 'white'
      ctx.fillRect(0, 0, canvas.width, canvas.height)
    }

    window.addEventListener('resize', handleResize)

    return () => {
      window.cancelAnimationFrame(animationRef.current)
      window.removeEventListener('resize', handleResize)
    }
  }, [])

  return (
    <canvas
      ref={canvasRef}
      className={`w-full h-full ${className || ""}`}
      style={{ background: 'white' }}
    />
  )
}

