import { InkRevealText } from "./InkRevealText"

interface LifePoint {
  month: string
  value: number
  type: "actual" | "follow" | "rebel"
}

const lifeData: LifePoint[] = [
  { month: "Jan", value: 65, type: "actual" },
  { month: "Feb", value: 70, type: "actual" },
  { month: "Mar", value: 68, type: "actual" },
  { month: "Apr", value: 75, type: "follow" },
  { month: "May", value: 78, type: "follow" },
  { month: "Jun", value: 82, type: "follow" },
]

export function LifeChart() {
  const chartHeight = 200

  return (
    <div className="w-full freya-mode">
      <div className="mb-6">
        <InkRevealText text="Life Trajectory" className="text-sm font-light opacity-60 mb-2" />
        <p className="text-xs opacity-40 leading-relaxed italic font-serif">
          <InkRevealText text="Your choices shape the path ahead. Each decision creates a new branch in your timeline." />
        </p>
      </div>

      <div className="relative border border-black/5 p-6" style={{ height: chartHeight + 60 }}>
        <svg className="w-full h-full" viewBox="0 0 300 200" preserveAspectRatio="none">
          {[25, 50, 75].map((val) => (
            <line key={val} x1="0" y1={200 - (val / 100) * 200} x2="300" y2={200 - (val / 100) * 200} stroke="black" strokeWidth="0.5" opacity="0.05" />
          ))}
          <path
            d={lifeData.map((point, i) => `${i === 0 ? "M" : "L"} ${(i / (lifeData.length - 1)) * 300} ${200 - (point.value / 100) * 200}`).join(" ")}
            stroke="black" strokeWidth="1" fill="none"
          />
          {lifeData.map((point, i) => {
            const x = (i / (lifeData.length - 1)) * 300
            const y = 200 - (point.value / 100) * 200
            return (
              <g key={i}>
                <circle cx={x} cy={y} r="2" fill="black" />
                {point.type === "follow" && (
                  <circle cx={x} cy={y} r="5" fill="none" stroke="black" strokeWidth="0.5" strokeDasharray="1 1" className="animate-spin" style={{ animationDuration: '3s' }} />
                )}
              </g>
            )
          })}
        </svg>
        <div className="flex justify-between mt-2 px-2 text-[8px] font-mono opacity-30">
          {lifeData.map((point, i) => <span key={i}>{point.month}</span>)}
        </div>
      </div>

      <div className="mt-8 border border-black/10 p-6 text-center">
        <InkRevealText text="Unlock 3-month forecast" className="text-sm font-serif mb-2" />
        <p className="text-[10px] opacity-40 mb-4 font-sans uppercase tracking-widest">See the full butterfly effect of your choices</p>
        <button className="px-8 py-3 bg-black text-white text-[10px] uppercase tracking-widest">Upgrade to Pro</button>
      </div>
    </div>
  )
}

