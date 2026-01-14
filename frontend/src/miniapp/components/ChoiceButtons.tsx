import { useState } from "react"

interface ChoiceButtonsProps {
  onAccept: () => void
  onReject: () => void
}

export function ChoiceButtons({ onAccept, onReject }: ChoiceButtonsProps) {
  const [selected, setSelected] = useState<"accept" | "reject" | null>(null)

  const handleAccept = () => {
    setSelected("accept")
    setTimeout(() => {
      onAccept()
    }, 300)
  }

  const handleReject = () => {
    setSelected("reject")
    setTimeout(() => {
      onReject()
    }, 300)
  }

  return (
    <div className="flex gap-4 w-full max-w-sm mx-auto mt-4">
      <button
        onClick={handleAccept}
        disabled={selected !== null}
        className={`flex-1 border border-black dark:border-white py-3 text-sm font-light transition-all ${
          selected === "accept"
            ? "bg-black text-white dark:bg-white dark:text-black"
            : selected === null
            ? "hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black"
            : "opacity-50 cursor-not-allowed"
        }`}
      >
        接受
      </button>
      <button
        onClick={handleReject}
        disabled={selected !== null}
        className={`flex-1 border border-black dark:border-white py-3 text-sm font-light transition-all ${
          selected === "reject"
            ? "bg-black text-white dark:bg-white dark:text-black"
            : selected === null
            ? "hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black"
            : "opacity-50 cursor-not-allowed"
        }`}
      >
        拒绝
      </button>
    </div>
  )
}

