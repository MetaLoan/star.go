import React from 'react';
import { InkRevealText } from './InkRevealText';

interface PageHeaderDescStyleProps {
  text: string;
  className?: string;
  delay?: number;
}

export const PageHeaderDescStyle: React.FC<PageHeaderDescStyleProps> = ({ text, className, delay = 0.5 }) => {
  return (
    <p className={`text-[10px] uppercase tracking-[0.2em] opacity-30 ${className || ''}`}>
      <InkRevealText text={text} delay={delay} />
    </p>
  );
};

