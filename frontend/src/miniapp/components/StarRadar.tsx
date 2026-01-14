import { useState, useEffect, useMemo } from "react"
import { InkRevealText } from "./InkRevealText"

const baseDimensions = [
  { name: "Passion", value: 75, angle: 0 },
  { name: "Love", value: 60, angle: 72 },
  { name: "Wealth", value: 85, angle: 144 },
  { name: "Intellect", value: 70, angle: 216 },
  { name: "Spirit", value: 90, angle: 288 },
]

export function StarRadar() {
  const [time, setTime] = useState(0)
  
  const dimensions = useMemo(() => {
    return baseDimensions.map((d, i) => {
      const t = time * 0.1
      const jitter = Math.sin(t + i) * 2
      return { ...d, value: Math.max(0, Math.min(100, d.value + jitter)) }
    })
  }, [time])
  
  useEffect(() => {
    const timer = setInterval(() => setTime(prev => prev + 1), 50)
    return () => clearInterval(timer)
  }, [])
  
  const size = 300
  const center = size / 2
  const maxRadius = size / 2 - 40

  const polarToCartesian = (angle: number, radius: number) => {
    const rad = ((angle - 90) * Math.PI) / 180
    return {
      x: center + radius * Math.cos(rad),
      y: center + radius * Math.sin(rad),
    }
  }

  const dataPoints = dimensions.map((d) => polarToCartesian(d.angle, (d.value / 100) * maxRadius))
  const pathData = dataPoints.map((p, i) => `${i === 0 ? "M" : "L"} ${p.x} ${p.y}`).join(" ") + " Z"

  return (
    <div className="w-full freya-mode">
      <div className="box-frame p-6 relative flex flex-col items-center">
        <svg width={size} height={size} className="overflow-visible">
          {[0.25, 0.5, 0.75, 1].map((scale) => (
            <polygon
              key={scale}
              points={dimensions.map(d => {
                const p = polarToCartesian(d.angle, maxRadius * scale)
                return `${p.x},${p.y}`
              }).join(" ")}
              fill="none"
              stroke="black"
              strokeWidth="0.5"
              opacity="0.05"
            />
          ))}
          {dimensions.map((d) => {
            const endpoint = polarToCartesian(d.angle, maxRadius)
            return <line key={d.name} x1={center} y1={center} x2={endpoint.x} y2={endpoint.y} stroke="black" strokeWidth="0.5" opacity="0.1" />
          })}
          <path d={pathData} fill="rgba(0,0,0,0.03)" stroke="black" strokeWidth="1" />
          {dataPoints.map((p, i) => <circle key={i} cx={p.x} cy={p.y} r="1.5" fill="black" />)}
          {dimensions.map((d) => {
            const pos = polarToCartesian(d.angle, maxRadius + 20)
            return <text key={d.name} x={pos.x} y={pos.y} textAnchor="middle" dominantBaseline="middle" className="text-[8px] uppercase tracking-widest fill-black font-light">{d.name}</text>
          })}
        </svg>

        <div className="mt-8 space-y-3 w-full">
          <div className="flex gap-3 text-[10px]">
            <span className="opacity-30 uppercase tracking-tighter shrink-0">Alert</span>
            <span className="opacity-60 font-serif italic"><InkRevealText text="Mercury retrograde affecting Intellect dimension" /></span>
          </div>
          <div className="flex gap-3 text-[10px]">
            <span className="opacity-30 uppercase tracking-tighter shrink-0">Peak</span>
            <span className="opacity-60 font-serif italic"><InkRevealText text="Spirit energy at highest point this month" /></span>
          </div>
        </div>
      </div>
    </div>
  )
}

