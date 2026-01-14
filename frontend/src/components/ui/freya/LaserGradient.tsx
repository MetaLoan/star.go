import type React from "react"
import { cn } from "../../../utils/cn"

export interface LaserGradientProps extends React.HTMLAttributes<HTMLDivElement> {
  animated?: boolean
  speed?: "slow" | "normal" | "fast"
}

export function LaserGradient({
  className,
  animated = true,
  speed = "normal",
  children,
  ...props
}: LaserGradientProps) {
  const speeds = {
    slow: "animate-laser duration-[12s]",
    normal: "animate-laser duration-[8s]",
    fast: "animate-laser duration-[4s]",
  }

  return (
    <div className={cn("laser-gradient", animated && speeds[speed], className)} {...props}>
      {children}
    </div>
  )
}

export function LaserBorder({ className, children, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={cn("relative p-[2px] rounded-xl overflow-hidden", className)} {...props}>
      <div className="absolute inset-0 laser-gradient" />
      <div className="relative bg-white dark:bg-black rounded-xl">{children}</div>
    </div>
  )
}

