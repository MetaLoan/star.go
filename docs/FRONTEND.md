# 前端功能文档

## 1. 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Next.js | 15.x | React 全栈框架 |
| React | 19.x | UI 库 |
| TypeScript | 5.x | 类型安全 |
| HeroUI | 2.x | UI 组件库 |
| TailwindCSS | 3.x | 样式框架 |
| Framer Motion | - | 动画库 |

---

## 2. 目录结构

```
src/
├── app/                      # Next.js App Router
│   ├── layout.tsx           # 根布局
│   ├── page.tsx             # 首页（主应用）
│   ├── globals.css          # 全局样式
│   └── providers.tsx        # 全局 Provider
│
├── components/               # React 组件
│   ├── chart/               # 星盘组件
│   ├── timeline/            # 时间线组件
│   ├── forecast/            # 预测组件
│   ├── factors/             # 影响因子组件
│   └── input/               # 输入组件
│
└── lib/                      # 核心库
    ├── api/                 # API 客户端
    ├── astro/               # 占星计算
    ├── services/            # 业务服务
    └── types/               # 类型定义
```

---

## 3. 组件详解

### 3.1 星盘组件

#### NatalChartSVG

**路径**: `src/components/chart/NatalChartSVG.tsx`

**功能**: 渲染 SVG 格式的本命盘图。

**Props**:
```typescript
interface NatalChartSVGProps {
  planets: PlanetPosition[];
  houses: HouseCusp[];
  aspects: AspectData[];
  ascendant: number;
  size?: number;  // 默认 600
}
```

**功能特点**:
- 12 宫位显示（Placidus 分宫制）
- 12 星座环
- 行星位置标记
- 相位连线（可配置颜色）
- 上升点/中天点标记
- 响应式大小

**使用示例**:
```tsx
<NatalChartSVG
  planets={chart.planets}
  houses={chart.houses}
  aspects={chart.aspects}
  ascendant={chart.ascendant}
  size={500}
/>
```

#### ChartInfo

**路径**: `src/components/chart/ChartInfo.tsx`

**功能**: 显示星盘详细信息面板。

**Props**:
```typescript
interface ChartInfoProps {
  chart: ChartData;
}
```

**显示内容**:
- 行星列表（位置、星座、宫位）
- 相位列表
- 元素/模式平衡
- 主导行星
- 图形相位

---

### 3.2 时间线组件

#### LifeTimeline

**路径**: `src/components/timeline/LifeTimeline.tsx`

**功能**: 显示 80 年人生趋势线。

**Props**:
```typescript
interface LifeTimelineProps {
  trendData: LifeTrendData;
  currentYear: number;
}
```

**功能特点**:
- SVG 渐变趋势线
- 悬停显示详情
- 维度切换（综合/事业/关系/健康/财务/灵性）
- 重大行运标记
- 当前年份指示器
- 周期信息展示

**交互**:
- 悬停显示年份详情
- 点击维度切换显示
- 缩放和平移（待实现）

#### UnifiedTimeline

**路径**: `src/components/timeline/UnifiedTimeline.tsx`

**功能**: 统一时间序列图，支持 H/D/W/M/Y 切换。

**Props**:
```typescript
interface UnifiedTimelineProps {
  birthData: {
    name: string;
    date: string;
    latitude: number;
    longitude: number;
    timezone: string;
  };
  apiUrl?: string;
}
```

**功能特点**:
- **粒度切换**: H(小时)/D(日)/W(周)/M(月)/Y(年)
- **原始/标准化切换**: 查看原始分数或归一化显示
- **统计信息**: 均值、标准差、趋势、波动率
- **因子分解**: 显示分数来源
- **维度分数**: 五维度雷达图

**使用示例**:
```tsx
<UnifiedTimeline
  birthData={{
    name: "张三",
    date: "1990-05-15T10:30:00Z",
    latitude: 39.9042,
    longitude: 116.4074,
    timezone: "Asia/Shanghai"
  }}
  apiUrl="http://localhost:8080"
/>
```

#### ProfectionWheel

**路径**: `src/components/timeline/ProfectionWheel.tsx`

**功能**: 年限法轮盘可视化。

**Props**:
```typescript
interface ProfectionWheelProps {
  chart: ChartData;
  currentAge: number;
}
```

**功能特点**:
- 12 宫位轮盘
- 当前年份高亮
- 年主星显示
- 宫位主题提示

---

### 3.3 预测组件

#### DailyForecastView

**路径**: `src/components/forecast/DailyForecastView.tsx`

**功能**: 每日运势详情展示。

**Props**:
```typescript
interface DailyForecastViewProps {
  forecast: ProcessedDailyForecast;
}
```

**显示内容**:
- 综合评分（圆形进度条）
- 当日主题
- 五维度分数条
- 月亮信息（星座、月相、月空亡）
- 24 小时趋势
- 活跃相位
- 最佳时机建议

**维度标签**:
| 键 | 显示 |
|----|----|
| career | 事业 |
| relationship | 关系 |
| health | 健康 |
| finance | 财务 |
| spiritual | 灵性 |

#### WeeklyForecastView

**路径**: `src/components/forecast/WeeklyForecastView.tsx`

**功能**: 每周运势概览。

**Props**:
```typescript
interface WeeklyForecastViewProps {
  forecast: ProcessedWeeklyForecast;
}
```

**显示内容**:
- 周度综合评分
- 周主题
- 7 天概览
- 关键日期
- 最佳日期推荐
- 周内行运

---

### 3.4 影响因子组件

#### InfluenceFactorsPanel

**路径**: `src/components/factors/InfluenceFactorsPanel.tsx`

**功能**: 影响因子可视化和配置。

**Props**:
```typescript
interface InfluenceFactorsPanelProps {
  factors: FactorResult;
  onWeightsChange?: (weights: FactorWeights) => void;
}
```

**功能特点**:
- 因子列表（正面/负面分类）
- 总调整值显示
- 权重滑块配置
- 维度调整预览

**因子类型显示**:
| 类型 | 中文 |
|------|------|
| dignity | 尊贵度 |
| retrograde | 逆行 |
| aspectPhase | 相位阶段 |
| aspectOrb | 容许度 |
| outerPlanet | 外行星 |
| profectionLord | 年主星 |
| lunarPhase | 月相 |
| planetaryHour | 行星时 |

---

### 3.5 输入组件

#### BirthDataForm

**路径**: `src/components/input/BirthDataForm.tsx`

**功能**: 出生数据输入表单。

**Props**:
```typescript
interface BirthDataFormProps {
  onSubmit: (data: BirthDataInput) => void;
  loading?: boolean;
}
```

**字段**:
- 姓名（可选）
- 出生日期
- 出生时间
- 出生地点（经纬度或地名搜索）
- 时区

**验证**:
- 日期格式验证
- 经纬度范围验证
- 必填字段检查

---

## 4. 页面结构

### 4.1 首页 (page.tsx)

首页是主应用页面，包含完整的占星功能。

**布局结构**:
```
┌────────────────────────────────────────────────────┐
│                     Header                          │
├────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────┐  ┌──────────────────────────────┐ │
│  │             │  │                               │ │
│  │  出生数据    │  │         星盘 SVG             │ │
│  │  表单       │  │                               │ │
│  │             │  │                               │ │
│  └─────────────┘  └──────────────────────────────┘ │
│                                                     │
├────────────────────────────────────────────────────┤
│                   Tabs 导航                         │
│  [概览] [每日] [每周] [人生趋势] [年限法]           │
├────────────────────────────────────────────────────┤
│                                                     │
│                    Tab 内容                         │
│                                                     │
└────────────────────────────────────────────────────┘
```

**Tab 内容**:
1. **概览**: 星盘信息、当前状态
2. **每日**: 每日预测、影响因子
3. **每周**: 每周预测
4. **人生趋势**: 统一时间线、人生趋势
5. **年限法**: 年限法轮盘、未来年份预览

### 4.2 状态管理

使用 `useAstroData` Hook 管理全局状态：

```typescript
const {
  birthData,      // 出生数据
  chart,          // 本命盘
  daily,          // 每日预测
  weekly,         // 每周预测
  lifeTrend,      // 人生趋势
  loading,        // 加载状态
  error,          // 错误信息
  initialize,     // 初始化函数
} = useAstroData();
```

---

## 5. API 客户端

### 5.1 客户端配置

**路径**: `src/lib/api/client.ts`

```typescript
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const calcApi = {
  calculateChart: (birthData) => fetcher('/api/calc/chart', birthData),
  calculateDaily: (birthData, date, withFactors) => ...,
  calculateWeekly: (birthData, date, withFactors) => ...,
  calculateLifeTrend: (birthData, startYear, endYear) => ...,
  calculateTimeSeries: (birthData, start, end, granularity) => ...,
};
```

### 5.2 React Hooks

**路径**: `src/lib/api/hooks.ts`

```typescript
// 本命盘计算
export function useChartCalculation() {
  const [state, setState] = useState<UseApiState<ChartData>>(...);
  const calculate = useCallback(async (birthData) => {...}, []);
  return { ...state, calculate };
}

// 每日预测
export function useDailyForecast() {...}

// 每周预测
export function useWeeklyForecast() {...}

// 人生趋势
export function useLifeTrend() {...}

// 综合数据（一次获取所有）
export function useAstroData() {...}
```

---

## 6. 样式系统

### 6.1 主题配置

使用 TailwindCSS 自定义主题：

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        cosmic: {
          nova: '#00D4FF',
          aurora: '#FF6B9D',
          nebula: '#A855F7',
          void: '#0A0A0F',
        }
      },
      fontFamily: {
        display: ['Space Grotesk', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      }
    }
  }
}
```

### 6.2 全局样式

**路径**: `src/app/globals.css`

```css
/* 玻璃拟态卡片 */
.glass-card {
  @apply bg-white/5 backdrop-blur-md border border-white/10 rounded-xl;
}

/* 宇宙背景 */
.cosmic-bg {
  background: radial-gradient(ellipse at top, #1a1a2e 0%, #0a0a0f 100%);
}

/* 发光效果 */
.glow {
  box-shadow: 0 0 20px rgba(0, 212, 255, 0.3);
}
```

---

## 7. 动画效果

使用 Framer Motion 实现动画：

```tsx
import { motion, AnimatePresence } from 'framer-motion';

// 页面切换动画
<AnimatePresence mode="wait">
  <motion.div
    key={activeTab}
    initial={{ opacity: 0, y: 20 }}
    animate={{ opacity: 1, y: 0 }}
    exit={{ opacity: 0, y: -20 }}
    transition={{ duration: 0.3 }}
  >
    {content}
  </motion.div>
</AnimatePresence>

// 进入动画
<motion.div
  initial={{ scale: 0.9, opacity: 0 }}
  animate={{ scale: 1, opacity: 1 }}
  transition={{ delay: 0.1 }}
>
  <Card />
</motion.div>
```

---

## 8. 响应式设计

### 8.1 断点

| 断点 | 宽度 | 用途 |
|------|------|------|
| sm | 640px | 手机横屏 |
| md | 768px | 平板竖屏 |
| lg | 1024px | 平板横屏/笔记本 |
| xl | 1280px | 桌面 |
| 2xl | 1536px | 大屏 |

### 8.2 响应式布局

```tsx
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  <Card />
  <Card />
  <Card />
</div>
```

---

## 9. 性能优化

### 9.1 代码分割

使用 Next.js 动态导入：

```tsx
import dynamic from 'next/dynamic';

const HeavyComponent = dynamic(
  () => import('@/components/HeavyComponent'),
  { loading: () => <Spinner /> }
);
```

### 9.2 数据缓存

使用 React Query 或 SWR 缓存 API 响应（待实现）。

### 9.3 虚拟化

大列表使用虚拟滚动（待实现）。

---

## 10. 可访问性

### 10.1 键盘导航

- Tab 键切换焦点
- Enter 键激活按钮
- Escape 键关闭模态框

### 10.2 屏幕阅读器

- 语义化 HTML 标签
- ARIA 标签
- alt 文本

### 10.3 颜色对比

- 确保文字与背景对比度 >= 4.5:1
- 不仅依赖颜色传达信息

---

## 11. 错误处理

### 11.1 错误边界

```tsx
<ErrorBoundary fallback={<ErrorFallback />}>
  <App />
</ErrorBoundary>
```

### 11.2 API 错误

```tsx
const { data, error, loading } = useAstroData();

if (error) {
  return <Alert color="danger">{error.message}</Alert>;
}

if (loading) {
  return <Spinner />;
}
```

---

## 12. 测试（待实现）

### 12.1 单元测试

```typescript
// components/__tests__/NatalChartSVG.test.tsx
describe('NatalChartSVG', () => {
  it('renders planets correctly', () => {...});
  it('draws aspect lines', () => {...});
});
```

### 12.2 集成测试

```typescript
// pages/__tests__/Home.test.tsx
describe('Home Page', () => {
  it('loads chart on form submit', () => {...});
});
```

### 12.3 E2E 测试

使用 Playwright 或 Cypress。

---

## 13. 开发指南

### 13.1 添加新组件

1. 在 `src/components/` 下创建目录
2. 创建组件文件
3. 导出到 `index.ts`
4. 添加到文档

### 13.2 添加新 API

1. 在 `src/lib/api/client.ts` 添加方法
2. 创建对应的 Hook
3. 更新类型定义

### 13.3 样式规范

- 使用 TailwindCSS 类
- 自定义样式放入 `globals.css`
- 组件特定样式使用 CSS Modules

