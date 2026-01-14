import { cn } from "../../utils/cn"

export type NavItem = "destiny" | "echo" | "oracle" | "connect"

interface BottomNavProps {
  activeTab: NavItem
  onTabChange: (tab: NavItem) => void
}

export function BottomNav({ activeTab, onTabChange }: BottomNavProps) {
  const items: { id: NavItem; name: string; icon: string }[] = [
    { id: "destiny", name: "Destiny", icon: "insights" },
    { id: "echo", name: "Echo", icon: "graphic_eq" },
    { id: "oracle", name: "Oracle", icon: "visibility" },
    { id: "connect", name: "Connect", icon: "all_inclusive" },
  ]

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-black/5 z-50 h-16 freya-mode">
      <div className="flex items-center justify-around h-full px-2 max-w-screen-sm mx-auto">
        {items.map((item) => (
          <button
            key={item.id}
            onClick={() => onTabChange(item.id)}
            className={cn(
              "flex flex-col items-center justify-center gap-1 min-w-[60px] transition-all",
              activeTab === item.id ? "opacity-100 scale-110" : "opacity-30 hover:opacity-60"
            )}
          >
            <span className="material-symbols-outlined text-[20px]" style={{ fontVariationSettings: "'wght' 200" }}>{item.icon}</span>
            <span className="text-[8px] uppercase tracking-[0.2em] font-sans">{item.name}</span>
          </button>
        ))}
      </div>
    </nav>
  )
}

