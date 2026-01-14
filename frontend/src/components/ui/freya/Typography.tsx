import type React from "react"
import { forwardRef, useEffect, useState } from "react"
import { cn } from "../../../utils/cn"

// Primary Page Title
export interface PageTitleProps extends React.HTMLAttributes<HTMLHeadingElement> {
  inkEffect?: boolean
}

export const PageTitle = forwardRef<HTMLHeadingElement, PageTitleProps>(
  ({ className, inkEffect = true, children, ...props }, ref) => {
    const [isVisible, setIsVisible] = useState(!inkEffect)

    useEffect(() => {
      if (inkEffect) {
        const timer = setTimeout(() => setIsVisible(true), 100)
        return () => clearTimeout(timer)
      }
    }, [inkEffect])

    return (
      <h1
        ref={ref}
        className={cn(
          "font-serif text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight",
          "text-black dark:text-white",
          isVisible && inkEffect && "ink-effect",
          className,
        )}
        {...props}
      >
        {children}
      </h1>
    )
  },
)

PageTitle.displayName = "PageTitle"

// Primary Page Subtitle
export const PageSubtitle = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <p
        ref={ref}
        className={cn("font-sans text-lg md:text-xl text-gray-500 dark:text-gray-400 leading-relaxed", "mt-4", className)}
        {...props}
      >
        {children}
      </p>
    )
  },
)

PageSubtitle.displayName = "PageSubtitle"

// Section Title
export const SectionTitle = forwardRef<HTMLHeadingElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <h2
        ref={ref}
        className={cn("font-serif text-2xl md:text-3xl font-semibold tracking-tight", "text-black dark:text-white", className)}
        {...props}
      >
        {children}
      </h2>
    )
  },
)

SectionTitle.displayName = "SectionTitle"

// Article Title
export const ArticleTitle = forwardRef<HTMLHeadingElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <h3
        ref={ref}
        className={cn("font-serif text-xl md:text-2xl font-semibold", "text-black dark:text-white", className)}
        {...props}
      >
        {children}
      </h3>
    )
  },
)

ArticleTitle.displayName = "ArticleTitle"

// Detail Text
export const DetailText = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <p ref={ref} className={cn("font-sans text-base text-black/80 dark:text-white/80 leading-relaxed", className)} {...props}>
        {children}
      </p>
    )
  },
)

DetailText.displayName = "DetailText"

// Emphasis Text
export const EmphasisText = forwardRef<HTMLSpanElement, React.HTMLAttributes<HTMLSpanElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <span ref={ref} className={cn("font-sans font-semibold text-black dark:text-white laser-text", className)} {...props}>
        {children}
      </span>
    )
  },
)

EmphasisText.displayName = "EmphasisText"

// Laser Text
export const LaserText = forwardRef<HTMLSpanElement, React.HTMLAttributes<HTMLSpanElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <span ref={ref} className={cn("laser-text font-serif font-bold", className)} {...props}>
        {children}
      </span>
    )
  },
)

LaserText.displayName = "LaserText"

export interface InkTextProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: "normal" | "emphasis"
  delay?: number
  duration?: number
}

export const InkText = forwardRef<HTMLSpanElement, InkTextProps>(
  ({ className, children, variant = "normal", delay = 0, duration = 0.8, ...props }, ref) => {
    const [isVisible, setIsVisible] = useState(false)

    useEffect(() => {
      const timer = setTimeout(() => setIsVisible(true), delay * 1000)
      return () => clearTimeout(timer)
    }, [delay])

    return (
      <span
        ref={ref}
        className={cn(
          "inline-block",
          variant === "emphasis" ? "font-sans font-semibold text-black dark:text-white laser-text" : "font-sans text-black dark:text-white",
          isVisible && "ink-effect",
          className,
        )}
        style={{
          animationDuration: `${duration}s`,
        }}
        {...props}
      >
        {children}
      </span>
    )
  },
)

InkText.displayName = "InkText"

