import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Star 占星系统',
  description: '占星时间系统 - 人生节律与趋势分析',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="zh-CN">
      <body className="bg-gray-950 text-white min-h-screen">
        {children}
      </body>
    </html>
  );
}

