import React from 'react';
import { InkRevealText } from './InkRevealText';
import { PageHeaderDescStyle } from './PageHeaderDescStyle';

interface PageHeaderStyleProps {
  title: string;
  subtitle?: string;
}

export const PageHeaderStyle: React.FC<PageHeaderStyleProps> = ({ title, subtitle }) => {
  return (
    <header className="space-y-2">
      <h1 className="text-4xl font-serif font-light tracking-tight">
        <InkRevealText text={title} />
      </h1>
      {subtitle && (
        <PageHeaderDescStyle text={subtitle} />
      )}
    </header>
  );
};

