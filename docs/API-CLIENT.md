# API 客户端使用指南

本文档说明如何使用 Star 系统的 API 客户端，支持前后端分离部署。

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                     前端应用                                 │
├─────────────────────────────────────────────────────────────┤
│  CSR (Client)          │        SSR (Server)                │
│  ┌─────────────────┐   │   ┌─────────────────────────┐     │
│  │  useAstroData   │   │   │  calculateChartServer   │     │
│  │  useChartCalc   │   │   │  calculateDailyServer   │     │
│  │  useDailyFcast  │   │   │  fetchUserServer        │     │
│  └────────┬────────┘   │   └────────────┬────────────┘     │
│           │            │                │                   │
│           └────────────┼────────────────┘                   │
│                        │                                    │
│                        ▼                                    │
│              ┌─────────────────┐                            │
│              │   API Client    │                            │
│              │  /lib/api/      │                            │
│              └────────┬────────┘                            │
└───────────────────────┼─────────────────────────────────────┘
                        │ HTTP
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                     后端 API                                 │
├─────────────────────────────────────────────────────────────┤
│  /api/calc/chart      计算本命盘（无需用户）                  │
│  /api/calc/daily      计算每日预测                           │
│  /api/calc/weekly     计算每周预测                           │
│  /api/calc/life-trend 计算人生趋势                           │
│  /api/users/          用户管理                               │
│  /api/agent/          智能体接口                             │
└─────────────────────────────────────────────────────────────┘
```

## 环境配置

### 前端独立部署时

在 `.env.local` 中配置 API 地址：

```bash
# 客户端 API 地址（CSR）
NEXT_PUBLIC_API_URL=https://api.your-domain.com

# 服务端 API 地址（SSR）
API_URL=https://api.your-domain.com
```

### 同一服务部署时

不需要配置，默认使用相对路径。

## CSR（客户端渲染）使用

### 1. 基础 API 客户端

```typescript
import { calcApi, userApi, forecastApi } from '@/lib/api/client';

// 计算本命盘
const chart = await calcApi.calculateChart({
  name: '张三',
  date: '1990-06-15T10:30:00Z',
  latitude: 39.9042,
  longitude: 116.4074,
  timezone: 'Asia/Shanghai',
});

// 计算每日预测（带影响因子）
const daily = await calcApi.calculateDaily(birthData, undefined, true);

// 计算人生趋势
const trend = await calcApi.calculateLifeTrend(birthData, 1990, 2070);
```

### 2. React Hooks（推荐）

```typescript
import { useAstroData, useDailyForecast, useChartCalculation } from '@/lib/api/hooks';

function MyComponent() {
  // 综合数据 Hook
  const { chart, daily, weekly, lifeTrend, loading, initialize } = useAstroData();
  
  // 或单独使用
  const { data: chartData, calculate, loading } = useChartCalculation();
  
  const handleSubmit = async (birthData) => {
    await initialize(birthData);
  };
  
  if (loading) return <Spinner />;
  
  return <ChartView chart={chart} />;
}
```

### 3. 综合数据 Hook

`useAstroData` 提供一站式数据管理：

```typescript
const {
  birthData,       // 出生数据
  chart,           // 本命盘
  daily,           // 每日预测
  weekly,          // 每周预测
  lifeTrend,       // 人生趋势
  loading,         // 加载状态
  error,           // 错误信息
  initialize,      // 初始化（输入出生数据后调用）
  refreshDaily,    // 刷新每日预测
  refreshWeekly,   // 刷新每周预测
  reset,           // 重置所有数据
} = useAstroData();
```

## SSR（服务端渲染）使用

### 服务端数据获取

```typescript
// app/user/[id]/page.tsx
import { fetchUserServer, calculateDailyServer } from '@/lib/api/server';

export default async function UserPage({ params }: { params: { id: string } }) {
  // 服务端获取数据
  const user = await fetchUserServer(params.id);
  const daily = await calculateDailyServer(user.birthData);
  
  return (
    <div>
      <h1>{user.name}</h1>
      <DailyView data={daily} />
    </div>
  );
}
```

### 混合渲染

```typescript
// 服务端获取初始数据
export default async function Page() {
  const initialData = await calculateChartServer(DEMO_DATA);
  
  return (
    // 客户端组件接收服务端数据作为初始值
    <ClientComponent initialChart={initialData} />
  );
}
```

## API 端点详解

### 无需用户的计算接口

这些接口直接接收出生数据，不需要创建用户：

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/calc/chart` | POST | 计算本命盘 |
| `/api/calc/daily` | POST | 计算每日预测 |
| `/api/calc/weekly` | POST | 计算每周预测 |
| `/api/calc/life-trend` | POST | 计算人生趋势 |

### 请求格式

```typescript
// 本命盘
POST /api/calc/chart
{
  "name": "张三",
  "date": "1990-06-15T10:30:00Z",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "timezone": "Asia/Shanghai"
}

// 每日预测
POST /api/calc/daily
{
  "birthData": { ... },
  "date": "2026-01-04T00:00:00Z",  // 可选
  "withFactors": true              // 是否包含影响因子
}
```

### 用户相关接口

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/users` | POST | 创建用户 |
| `/api/users/[id]` | GET | 获取用户 |
| `/api/users/[id]/snapshot` | GET | 获取实时快照 |
| `/api/users/[id]/forecast` | GET | 获取预测 |

## 独立前端部署

### 1. 创建新的前端项目

```bash
npx create-next-app@latest my-astro-frontend
cd my-astro-frontend
```

### 2. 安装依赖

```bash
npm install @heroui/react framer-motion
```

### 3. 复制 API 客户端

将 `src/lib/api/` 目录复制到新项目：

```
my-astro-frontend/
├── src/
│   └── lib/
│       └── api/
│           ├── client.ts    # API 客户端
│           ├── hooks.ts     # React Hooks
│           ├── server.ts    # SSR 工具
│           └── index.ts     # 导出
```

### 4. 配置环境变量

```bash
# .env.local
NEXT_PUBLIC_API_URL=https://your-api-server.com
API_URL=https://your-api-server.com
```

### 5. 使用 API

```typescript
import { useAstroData } from '@/lib/api';

export default function App() {
  const { chart, daily, initialize, loading } = useAstroData();
  
  // ... 你的 UI 逻辑
}
```

## 类型定义

所有 API 响应类型都在 `client.ts` 中导出：

```typescript
import type {
  BirthDataInput,
  NatalChartResponse,
  DailyForecastResponse,
  WeeklyForecastResponse,
  LifeTrendResponse,
  UserSnapshotResponse,
} from '@/lib/api/client';
```

## 错误处理

```typescript
import { calcApi } from '@/lib/api/client';

try {
  const chart = await calcApi.calculateChart(birthData);
} catch (error) {
  if (error instanceof Error) {
    console.error('API Error:', error.message);
  }
}

// 或使用 Hook
const { error, loading, data } = useChartCalculation();

if (error) {
  return <ErrorMessage message={error.message} />;
}
```

## CORS 配置

如果前后端分离部署，需要在后端配置 CORS。

在 `next.config.js` 中：

```javascript
module.exports = {
  async headers() {
    return [
      {
        source: '/api/:path*',
        headers: [
          { key: 'Access-Control-Allow-Origin', value: '*' },
          { key: 'Access-Control-Allow-Methods', value: 'GET,POST,OPTIONS' },
          { key: 'Access-Control-Allow-Headers', value: 'Content-Type' },
        ],
      },
    ];
  },
};
```

## 最佳实践

1. **使用 Hooks**: 优先使用 `useAstroData` 等 Hooks，它们处理了加载状态和错误
2. **类型安全**: 使用导出的 TypeScript 类型
3. **SSR 优化**: 对 SEO 重要的页面使用服务端渲染
4. **缓存策略**: 本命盘数据可以缓存，每日预测需要实时获取
5. **错误边界**: 在 UI 层面添加错误边界处理 API 失败

