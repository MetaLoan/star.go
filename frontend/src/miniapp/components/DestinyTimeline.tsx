import { useEffect, useRef, useMemo } from "react"
import { InkRevealText } from "./InkRevealText"
import type { TimeSeriesResponse, TimeSeriesPoint } from "../../types"

interface DestinyTimelineProps {
  selectedHour: number | null
  onSelectHour: (hour: number | null) => void
  timeSeries?: TimeSeriesResponse | null
}

export function DestinyTimeline({ selectedHour, onSelectHour, timeSeries }: DestinyTimelineProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  
  const points = useMemo(() => {
    return timeSeries?.points || []
  }, [timeSeries])

  const selectedPoint = useMemo(() => {
    if (selectedHour === null || !points.length) return null
    return points.find((_, i) => i === selectedHour)
  }, [selectedHour, points])

  // 提取点的信息
  const detailInfo = useMemo(() => {
    if (!selectedPoint) return null
    
    // 从因素中提取警告或建议
    const factors = selectedPoint.factors || []
    const warningFactor = factors.find(f => !f.isPositive)
    const positiveFactor = factors.find(f => f.isPositive)
    
    return {
      type: warningFactor ? "warning" : "do",
      title: warningFactor ? "Static Noise" : "Energy Peak",
      content: warningFactor ? warningFactor.description : (positiveFactor ? positiveFactor.description : "Cosmic alignment is stable."),
    }
  }, [selectedPoint])

  useEffect(() => {
    if (selectedHour !== null && containerRef.current) {
      const el = containerRef.current.querySelector(`[data-hour="${selectedHour}"]`)
      if (el) {
        const rect = el.getBoundingClientRect()
        const cRect = containerRef.current.getBoundingClientRect()
        containerRef.current.scrollTo({
          left: containerRef.current.scrollLeft + (rect.left - cRect.left - cRect.width / 2 + rect.width / 2),
          behavior: 'smooth'
        })
      }
    }
  }, [selectedHour])

  return (
    <div className="mt-2 pb-12 freya-mode">
      {detailInfo && (
        <div className="relative mb-4">
          <div className="max-w-md mx-auto p-4 box-frame">
            <div className="flex items-center gap-3 mb-2">
              <span className="text-[8px] uppercase tracking-widest px-2 py-1 bg-black/5">{detailInfo.type}</span>
              <h3 className="text-sm font-serif"><InkRevealText text={detailInfo.title} /></h3>
            </div>
            <p className="text-xs opacity-60 leading-relaxed"><InkRevealText text={detailInfo.content} /></p>
          </div>
          <div className="flex justify-center mt-2">
            <div className="w-0 h-0 border-l-[6px] border-l-transparent border-r-[6px] border-r-transparent border-t-[8px] border-t-black" />
          </div>
        </div>
      )}

      <div ref={containerRef} className="relative overflow-x-auto scrollbar-hide">
        <div className="relative inline-flex items-start gap-0 py-10" style={{ minWidth: '100%', paddingLeft: 'calc(50% - 40px)', paddingRight: 'calc(50% - 40px)' }}>
          <div className="absolute left-0 right-0 top-10 h-[0.5px] bg-black/10" />
          {points.map((item, i) => {
            const isSelected = selectedHour === i
            return (
              <div key={i} className="relative flex flex-col items-center shrink-0" style={{ width: '80px' }}>
                <button
                  data-hour={i}
                  onClick={() => onSelectHour(isSelected ? null : i)}
                  className={`absolute w-4 h-4 rounded-full flex items-center justify-center transition-all top-0 -translate-y-1/2 ${
                    isSelected ? 'border border-black scale-125' : 'border border-black/5 hover:scale-110'
                  }`}
                >
                  <div className={`w-1.5 h-1.5 rounded-full ${isSelected ? 'bg-black' : 'bg-black/10'}`} />
                </button>
                <span className={`absolute top-4 text-[9px] font-mono ${isSelected ? 'opacity-100' : 'opacity-30'}`}>{item.label}</span>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}

