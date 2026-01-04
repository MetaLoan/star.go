# Star 占星系统 - 独立前端模板

这是一个可以独立部署的前端模板，通过 API 与 Star 后端通信。

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置 API 地址

创建 `.env.local` 文件：

```bash
NEXT_PUBLIC_API_URL=https://your-api-server.com
API_URL=https://your-api-server.com
```

### 3. 启动开发服务器

```bash
npm run dev
```

### 4. 构建生产版本

```bash
npm run build
npm start
```

## 项目结构

```
frontend-template/
├── src/
│   ├── app/           # Next.js App Router
│   ├── components/    # UI 组件
│   └── lib/
│       └── api/       # API 客户端
├── .env.local         # 环境变量
└── package.json
```

## 使用 API 客户端

### CSR（客户端渲染）

```typescript
import { useAstroData } from '@/lib/api';

function MyComponent() {
  const { chart, daily, initialize, loading } = useAstroData();
  
  const handleSubmit = async (birthData) => {
    await initialize(birthData);
  };
  
  return <ChartView chart={chart} />;
}
```

### SSR（服务端渲染）

```typescript
import { calculateChartServer } from '@/lib/api/server';

export default async function Page() {
  const chart = await calculateChartServer(birthData);
  return <ChartView chart={chart} />;
}
```

## 部署

### Vercel

```bash
vercel deploy
```

### Docker

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build
CMD ["npm", "start"]
```

## API 文档

详见后端仓库的 `/docs/API-CLIENT.md`。

