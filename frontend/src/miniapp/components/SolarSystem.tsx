import { useEffect, useState } from "react"
import { cn } from "../../utils/cn"

const PLANETS = [
  { name: 'sun', speed: 0, size: 85, orbit: 0, color: '#ffffff' },
  { name: 'mercury', speed: 0.2, size: 5, orbit: 85, color: '#d1d1d1' },
  { name: 'venus', speed: 0.6, size: 15, orbit: 170, color: '#e5e5e5' },
  { name: 'earth', speed: 1, size: 20, orbit: 255, color: '#ffffff' },
  { name: 'mars', speed: 1.9, size: 10, orbit: 340, color: '#cccccc' },
  { name: 'jupiter', speed: 11.9, size: 50, orbit: 425, color: '#ebebeb' },
  { name: 'saturn', speed: 29.5, size: 35, orbit: 510, color: '#d9d9d9' },
  { name: 'uranus', speed: 84, size: 30, orbit: 595, color: '#f0f0f0' },
  { name: 'neptune', speed: 164.8, size: 25, orbit: 680, color: '#c0c0c0' },
]

export function SolarSystem({ active }: { active: boolean }) {
  const [, setCounters] = useState({ years: 0, hours: 0, days: 0 })

  useEffect(() => {
    if (!active) return
    const timer = setInterval(() => {
      setCounters(prev => ({
        years: prev.years + 1,
        hours: prev.hours + Math.round(8760/200),
        days: prev.days + Math.round(365/20)
      }))
    }, 100)
    return () => clearInterval(timer)
  }, [active])

  return (
    <div className={cn(
      "absolute inset-0 overflow-hidden transition-opacity duration-1000",
      active ? "opacity-100" : "opacity-0 pointer-events-none"
    )}>
      <div className="absolute inset-0 flex items-end justify-center scale-[0.5] md:scale-[0.8] origin-bottom">
        <div className="relative w-0 h-0 preserve-3d" style={{ transformStyle: 'preserve-3d' }}>
          {PLANETS.map((planet, i) => (
            <div 
              key={planet.name}
              className={cn(
                "absolute rounded-full",
                planet.name !== 'sun' && "border-white"
              )}
              style={{
                width: planet.orbit * 2,
                height: planet.orbit * 2,
                left: -planet.orbit,
                top: -planet.orbit,
                borderWidth: planet.name !== 'sun' ? '0.2px' : 0,
                borderStyle: planet.name !== 'sun' ? 'solid' : 'none',
                animation: planet.speed > 0 ? `solar-orbit ${planet.speed * 5}s linear infinite` : 'none',
                zIndex: 10 - i,
                transformStyle: 'preserve-3d'
              }}
            >
              <div 
                className="absolute rounded-full"
                style={{
                  width: planet.size,
                  height: planet.size,
                  left: planet.name === 'sun' ? '0' : '50%',
                  top: planet.name === 'sun' ? '0' : '0',
                  transform: 'translate(-50%, -50%)',
                  background: planet.name === 'sun' 
                    ? 'radial-gradient(circle, #fff 0%, #aaa 100%)' 
                    : `linear-gradient(135deg, ${planet.color}, #333)`,
                  boxShadow: planet.name === 'sun' 
                    ? '0 0 60px #fff, 0 0 20px #888' 
                    : `0 0 10px ${planet.color}44`,
                }}
              >
                {planet.name === 'saturn' && (
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[180%] h-[20%] border border-white/20 rounded-[50%] rotate-[25deg]" />
                )}
                {planet.name === 'earth' && (
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[160%] h-[160%] animate-[solar-orbit_2s_linear_infinite]" style={{ animation: 'solar-orbit 2s linear infinite' }}>
                    <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 w-1.5 h-1.5 bg-white rounded-full shadow-[0_0_5px_#fff]" />
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      <style>{`
        @keyframes solar-orbit {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  )
}

