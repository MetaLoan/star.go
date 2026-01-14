import type React from "react"
import { forwardRef } from "react"
import { cn } from "../../../utils/cn"

export interface FreyaButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "ghost" | "laser"
  size?: "sm" | "md" | "lg"
  isLoading?: boolean
}

const FreyaButton = forwardRef<HTMLButtonElement, FreyaButtonProps>(
  ({ className, variant = "primary", size = "md", isLoading, children, disabled, ...props }, ref) => {
    return (
      <button
        ref={ref}
        disabled={disabled || isLoading}
        className={cn(
          // 基础样式
          "relative inline-flex items-center justify-center gap-2",
          "font-sans font-medium tracking-wide",
          "rounded-xl transition-all duration-300",
          "disabled:opacity-50 disabled:cursor-not-allowed",
          "active:scale-[0.98]",

          // 尺寸变体
          {
            "px-4 py-2 text-sm": size === "sm",
            "px-6 py-3 text-base": size === "md",
            "px-8 py-4 text-lg": size === "lg",
          },

          // 样式变体
          {
            "bg-black text-white dark:bg-white dark:text-black hover:opacity-90": variant === "primary",
            "bg-gray-100 text-black dark:bg-gray-800 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700": variant === "secondary",
            "hover:bg-gray-100/50 dark:hover:bg-gray-800/50": variant === "ghost",
            "laser-gradient text-white shadow-lg hover:shadow-xl": variant === "laser",
          },

          className,
        )}
        {...props}
      >
        {isLoading && (
          <div className="absolute inset-0 flex items-center justify-center bg-inherit rounded-xl">
            <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
          </div>
        )}
        <span className={cn("transition-opacity", isLoading && "opacity-0")}>{children}</span>
      </button>
    )
  },
)

FreyaButton.displayName = "FreyaButton"

export { FreyaButton }

