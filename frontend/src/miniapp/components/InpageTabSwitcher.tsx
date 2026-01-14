import React from 'react';
import { cn } from '../../utils/cn';

interface TabOption {
  id: string;
  label: string;
}

interface InpageTabSwitcherProps {
  options: TabOption[];
  activeId: string;
  onChange: (id: any) => void;
  className?: string;
}

export const InpageTabSwitcher: React.FC<InpageTabSwitcherProps> = ({ 
  options, 
  activeId, 
  onChange,
  className 
}) => {
  return (
    <div className={cn("flex border-b border-black/5", className)}>
      {options.map((option) => (
        <button
          key={option.id}
          onClick={() => onChange(option.id)}
          className={cn(
            "flex-1 py-4 text-[10px] uppercase tracking-widest transition-all",
            activeId === option.id 
              ? "opacity-100 font-bold" 
              : "opacity-20 hover:opacity-40"
          )}
        >
          {option.label}
        </button>
      ))}
    </div>
  );
};

