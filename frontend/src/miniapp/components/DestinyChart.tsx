import { useMemo, useState, useEffect, useRef, useCallback } from "react"
import { Tabs, Tab } from "@heroui/react"
import { TextReveal } from "./TextReveal"
import { StarfieldEffect } from "./StarfieldEffect"
import type { TimeSeriesResponse, TimeGranularity, DimensionScores } from "../../types"
import { cn } from "../../utils/cn"
import { InteractiveTrendChart9823EF, type VisibleRangeChange } from "../../components/timeline/InteractiveTrendChart9823EF"

interface DestinyChartProps {
  mode: "today" | "life"
  showChart?: boolean
  selectedHour: number | null
  onSelectHour: (hour: number | null) => void
  onCardSelect?: () => void
  timeSeries?: TimeSeriesResponse | null
  granularity?: TimeGranularity
  onGranularityChange?: (g: TimeGranularity) => void
  selectedDimension?: string
  onDimensionChange?: (d: string) => void
  onVisibleRangeChange?: (range: VisibleRangeChange) => void
  isLoading?: boolean
}

const CARD_FILES = [
  "0603205238985_00_93f05f199fe56df4354be531f1666b53.jpg",
  "0603205238985_02_5d4153da20d3d91dc0b84dd27247c1a5.jpg",
  "0603205238985_03_5c4ec1905aeae9bec5909153315e0269.jpg"
]

const getRandomCards = (): string[] => {
  return [...CARD_FILES].sort(() => Math.random() - 0.5).slice(0, 3)
}

export function DestinyChart({ 
  mode, 
  showChart = true, 
  selectedHour, 
  onSelectHour, 
  onCardSelect,
  timeSeries,
  granularity = 'day',
  onGranularityChange,
  onVisibleRangeChange,
  isLoading
}: DestinyChartProps) {
  const [cards, setCards] = useState<string[]>([])
  const [selectedCard, setSelectedCard] = useState<number | null>(null)
  const [isFlipping, setIsFlipping] = useState(false)
  const [isFlipped, setIsFlipped] = useState(false)
  const [isOthersHiding, setIsOthersHiding] = useState(false)
  const [isCentering, setIsCentering] = useState(false)
  const [isShaking, setIsShaking] = useState(false)
  const [isShrinking, setIsShrinking] = useState(false)
  const [isBlurring, setIsBlurring] = useState(false)
  const [cardName, setCardName] = useState<string>("")
  const [cardReading, setCardReading] = useState<string>("")
  const [showReading, setShowReading] = useState(false)
  const [readingMode, setReadingMode] = useState<"idle" | "fadein" | "fadeaway">("idle")
  const [hintMode, setHintMode] = useState<"idle" | "fadein" | "fadeaway">("idle")
  const [warpSpeed, setWarpSpeed] = useState(1)
  const [colorMode, setColorMode] = useState<"normal" | "rainbow" | "fading">("normal")
  
  const cardNames = ["The Wheel of Fortune", "Temperance", "Death"]
  const cardReadings = ["The wheel turns in your favor", "Harmony awaits your next move", "Embrace the change within"]
  
  const granularities: { id: TimeGranularity, label: string }[] = [
    { id: 'hour', label: 'H' },
    { id: 'day', label: 'D' },
    { id: 'week', label: 'W' },
    { id: 'month', label: 'M' },
    { id: 'year', label: 'Y' },
  ]

  useEffect(() => {
    setCards(getRandomCards())
  }, [])
  
  const speedAnimationRef = useRef<number | null>(null)
  
  const handleCardClick = (index: number) => {
    if (selectedCard !== null || isFlipping) return
    setSelectedCard(index)
    setCardName(cardNames[index])
    setCardReading(cardReadings[index])
    setHintMode("fadeaway")
    setIsOthersHiding(true)
    setIsCentering(true)
    
    setTimeout(() => {
      setIsShaking(true)
      setColorMode("rainbow")
      const startTime = Date.now()
      const duration = 2000
      const animateSpeed = () => {
        const elapsed = Date.now() - startTime
        const progress = Math.min(elapsed / duration, 1)
        const currentSpeed = 1 + 9 * (progress * progress)
        setWarpSpeed(currentSpeed)
        if (progress < 1) {
          speedAnimationRef.current = requestAnimationFrame(animateSpeed)
        } else {
          setIsShaking(false)
          setTimeout(() => {
            setIsFlipping(true)
            setIsFlipped(true)
            setColorMode("fading")
            setTimeout(() => {
              setIsFlipping(false)
              setTimeout(() => {
                setShowReading(true)
                setReadingMode("fadein")
              }, 1500)
            }, 500)
          }, 10)
        }
      }
      speedAnimationRef.current = requestAnimationFrame(animateSpeed)
    }, 500)
  }
  
  const handleReadingComplete = useCallback(() => {
    if (readingMode === "fadein") {
      setTimeout(() => setReadingMode("fadeaway"), 5000)
    } else if (readingMode === "fadeaway") {
      setIsBlurring(true)
      setTimeout(() => {
        setIsShrinking(true)
      }, 300)
      setTimeout(() => {
        onCardSelect?.()
      }, 800)
    }
  }, [readingMode, onCardSelect])

  const chartData = useMemo(() => {
    if (!timeSeries?.points || timeSeries.points.length === 0) return []
    return timeSeries.points.map(p => ({
      time: p.time,
      value: p.display,
      label: p.label,
      dimensions: p.dimensions as DimensionScores
    }))
  }, [timeSeries])

  return (
    <div className="w-full select-none freya-mode">
      <div 
        className={`relative overflow-hidden ${showChart ? 'box-frame !p-0' : 'border border-black'}`} 
        style={{ height: showChart ? 'auto' : '263px' }}
      >
        {!showChart && (
          <>
            <div className="absolute inset-0 z-10">
              <StarfieldEffect speedMultiplier={warpSpeed} colorMode={colorMode} />
            </div>
            <div className="absolute z-30 left-0 right-0 text-center top-5">
              {hintMode !== "fadeaway" ? (
                <span className="text-[10px] tracking-widest uppercase opacity-40">Draw Card to Unlock Today's Destiny</span>
              ) : !showReading && (
                <TextReveal text="Draw Card to Unlock Today's Destiny" mode="fadeaway" className="text-[10px] tracking-widest uppercase opacity-40" />
              )}
            </div>
            <div className="absolute z-30 left-0 right-0 text-center bottom-5 flex flex-col items-center gap-1">
              {showReading && (
                <>
                  <TextReveal text={cardName} mode={readingMode} className="text-sm font-serif tracking-wider" />
                  <TextReveal text={cardReading} mode={readingMode} className="text-[10px] italic opacity-60" startDelay={1000} onComplete={handleReadingComplete} />
                </>
              )}
            </div>
            <div className="absolute z-20 inset-0 flex items-center justify-center">
              <div className="flex gap-4">
                {cards.map((card, i) => {
                  const isSelected = selectedCard === i
                  const isOther = selectedCard !== null && !isSelected
                  return (
                    <div 
                      key={i} 
                      onClick={() => handleCardClick(i)}
                      className={`relative w-16 h-24 border border-black bg-white transition-all duration-500 cursor-pointer ${isSelected ? 'scale-125 z-50' : 'scale-100 z-10'} ${isOther ? 'opacity-0 scale-50' : 'opacity-100'}`}
                      style={{
                        transform: isFlipped && isSelected ? 'rotateY(180deg) scale(1.2)' : isShaking && isSelected ? 'rotate(5deg)' : ''
                      }}
                    >
                      <div className="absolute inset-0 flex items-center justify-center backface-hidden" style={{ backfaceVisibility: 'hidden' }}>
                        <div className="w-4 h-4 border border-black rounded-full" />
                      </div>
                      <div className="absolute inset-0 backface-hidden" style={{ backfaceVisibility: 'hidden', transform: 'rotateY(180deg)' }}>
                        <div className="w-full h-full bg-black/5 flex items-center justify-center text-[8px] uppercase font-serif">Revealed</div>
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>
          </>
        )}
        
        {showChart && (
          <div className="p-0">
            {/* 粒度切换 - Tab 形式，针对手机优化点击区域 */}
            <div className="flex justify-between items-center px-4 py-3 border-b border-black/5">
               <Tabs 
                  aria-label="Granularity"
                  variant="solid"
                  size="md"
                  radius="none"
                  selectedKey={granularity}
                  onSelectionChange={(key) => onGranularityChange?.(key as TimeGranularity)}
                  classNames={{
                    tabList: "bg-black/5 p-1 gap-0.5 rounded-none",
                    cursor: "bg-white shadow-none rounded-none border-[0.5px] border-black/10",
                    tab: "px-4 h-9",
                    tabContent: "group-data-[selected=true]:text-black group-data-[selected=true]:font-bold text-black/30 text-[11px] uppercase tracking-widest"
                  }}
                >
                  {granularities.map(g => (
                    <Tab key={g.id} title={g.label} />
                  ))}
                </Tabs>
                <div className="text-[10px] uppercase tracking-[0.25em] opacity-20 font-sans font-bold">Flow</div>
            </div>

            <div className="p-0">
              <InteractiveTrendChart9823EF
                data={chartData}
                aspectRatio={2}
                showDimensions={true}
                color="#000000"
                isLoading={isLoading}
                onVisibleRangeChange={onVisibleRangeChange}
                onPointClick={(p, dim) => {
                  const index = timeSeries?.points.findIndex(tp => tp.time === p.time);
                  if (index !== undefined && index !== -1) {
                    onSelectHour(index);
                  }
                }}
              />
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
