# 因子维度有符号影响系统设计文档

## 一、问题背景

### 旧系统的问题

在原有系统中，一个因子对所有维度的影响方向是一致的：

```
水星逆行 = -2.0
DimensionImpact = {career: 0.3, relationship: 0.2, ...}
结果：
  career: -2.0 × 0.3 = -0.6 ❌
  relationship: -2.0 × 0.2 = -0.4 ❌
  spiritual: -2.0 × 0.25 = -0.5 ❌
```

**问题**：所有维度都是负影响，曲线同向波动。

### 占星学理论

根据占星学理论，同一天体事件在不同维度的影响可以是不同方向的：

- **水星逆行**
  - 事业：❌ 负面（沟通障碍、计划延误）
  - 关系：❌ 负面（误解增多）
  - 灵性：✅ 正面（反思、复盘的好时机）

- **土星过境**
  - 事业：✅ 正面（责任、成就、建立结构）
  - 健康：❌ 负面（压力、慢性疾病）
  - 灵性：✅ 正面（业力清理、成熟）

## 二、新系统设计

### 核心概念：有符号维度影响

```go
type SignedDimensionImpact struct {
    Career       float64  // -1.0 到 +1.0
    Relationship float64
    Health       float64
    Finance      float64
    Spiritual    float64
}
```

**值域说明**：
- 正值：该维度受到正向影响
- 负值：该维度受到负向影响
- 0值：该维度不受影响或影响中性
- 绝对值：影响强度

### 新计算公式

```go
baseAdjustment = abs(BaseValue) × Weight × CurrentStrength

dimensionAdjustment = {
    career:       baseAdjustment × signedImpact.Career,
    relationship: baseAdjustment × signedImpact.Relationship,
    health:       baseAdjustment × signedImpact.Health,
    finance:      baseAdjustment × signedImpact.Finance,
    spiritual:    baseAdjustment × signedImpact.Spiritual,
}
```

### 示例计算

**水星逆行**：

```
BaseValue = 2.0 (强度)
Weight = 1.0
CurrentStrength = 1.0
baseAdjustment = abs(2.0) × 1.0 × 1.0 = 2.0

signedImpact = {
    career: -0.40,       // 负影响
    relationship: -0.35,  // 负影响
    spiritual: +0.30,     // 正影响！
}

结果：
  career: 2.0 × (-0.40) = -0.80 ❌
  relationship: 2.0 × (-0.35) = -0.70 ❌
  spiritual: 2.0 × 0.30 = +0.60 ✅
```

✅ **成功实现不同维度的不同方向影响！**

## 三、因子类型的默认影响模式

### 逆行因子

| 维度 | 水星逆行 | 金星逆行 | 火星逆行 | 木星逆行 | 土星逆行 |
|------|---------|---------|---------|---------|---------|
| 事业 | -0.40 | -0.15 | -0.35 | -0.25 | -0.20 |
| 关系 | -0.35 | -0.40 | -0.30 | -0.10 | -0.15 |
| 健康 | -0.10 | -0.05 | -0.25 | -0.10 | -0.20 |
| 财务 | -0.25 | -0.25 | -0.20 | -0.30 | -0.15 |
| 灵性 | +0.30 | +0.25 | +0.20 | +0.35 | +0.40 |

**规律**：所有逆行对灵性维度都是正面的（内省、深化）。

### 尊贵度因子

| 维度 | 入庙 | 旺相 | 失势 | 落陷 |
|------|------|------|------|------|
| 事业 | +0.30 | +0.35 | -0.25 | -0.30 |
| 关系 | +0.25 | +0.30 | -0.30 | -0.35 |
| 健康 | +0.25 | +0.30 | -0.20 | -0.25 |
| 财务 | +0.25 | +0.30 | -0.25 | -0.30 |
| 灵性 | +0.20 | +0.25 | -0.15 | -0.20 |

### 特殊事件因子

| 因子 | Career | Relationship | Health | Finance | Spiritual |
|------|--------|--------------|--------|---------|-----------|
| 月空亡 | -0.30 | -0.10 | -0.05 | -0.35 | +0.20 |
| 日月食 | -0.20 | -0.15 | -0.25 | -0.10 | +0.30 |
| 燃烧 | -0.25 | -0.20 | -0.15 | -0.15 | -0.05 |
| 互容 | +0.30 | +0.35 | +0.20 | +0.25 | +0.20 |
| 停滞 | +0.20 | +0.10 | -0.10 | +0.15 | +0.25 |

## 四、相位的复杂性处理

相位因子需要考虑三个层面：

1. **相位类型**：合/刑/冲/拱/六
2. **行星特性**：涉及哪两颗行星
3. **维度修正系数**

### 相位类型修正

```go
AspectTypeImpactModifiers = {
    "square": {
        "relationship": 0.7,  // 四分相对关系压力大
        "spiritual": 1.2,     // 但促进灵性成长
    },
    "trine": {
        "career": 1.1,
        "relationship": 1.2,  // 三分相关系和谐
    },
    "opposition": {
        "relationship": 0.6,  // 对冲对关系挑战大
        "spiritual": 1.1,     // 需要平衡与整合
    },
}
```

### 行星特性调整

```go
// 水星相位：放大事业和关系影响
if planet == Mercury {
    impact.Career *= 1.2
    impact.Relationship *= 1.1
}

// 金星相位：放大关系和财务影响
if planet == Venus {
    impact.Relationship *= 1.3
    impact.Finance *= 1.1
}
```

## 五、实现细节

### 文件结构

```
backend/astro/
├── factor_dimension_impacts.go  # 新增：有符号影响映射
├── score_calculator.go          # 修改：使用新的计算逻辑
└── score_breakdown.go           # 修改：全因子接口使用新逻辑
```

### 核心函数

```go
// 获取因子的有符号维度影响
func GetFactorDimensionImpact(factor *models.InfluenceFactor) SignedDimensionImpact

// 计算维度调整值（使用有符号影响）
func calculateDimensionAdjustment(factor *models.InfluenceFactor) models.DimensionScoresV2
```

### 映射数据

1. **FactorTypeDefaultImpacts** - 因子类型默认影响
2. **RetrogradeImpactsByPlanet** - 行星逆行专属影响
3. **AspectTypeImpactModifiers** - 相位类型修正系数
4. **DignityImpactsByType** - 尊贵度类型影响

## 六、效果预期

### 曲线独立性

现在五个维度的曲线可以独立变化：

```
事业曲线：  ↗️ ↘️ ↗️ ↘️ ↗️
关系曲线：  ↘️ ↗️ ↘️ ↗️ ↘️
健康曲线：  ↗️ ↗️ ↘️ ↗️ ↗️
财务曲线：  ↘️ ↘️ ↗️ ↗️ ↘️
灵性曲线：  ↗️ ↗️ ↗️ ↘️ ↗️
```

### API 响应示例

```json
{
  "overall": {
    "positiveCount": 5,
    "negativeCount": 3,
    "netAdjustment": 2.3
  },
  "dimensions": {
    "career": {
      "positiveTotal": 3.2,
      "negativeTotal": -1.5,
      "netAdjustment": 1.7  // 正
    },
    "spiritual": {
      "positiveTotal": 4.5,
      "negativeTotal": -0.8,
      "netAdjustment": 3.7  // 强正
    }
  }
}
```

## 七、扩展性

### 添加新因子类型

1. 在 `FactorTypeDefaultImpacts` 中定义默认影响
2. 如需细化，添加专门的映射表（如逆行的行星映射）
3. 在 `GetFactorDimensionImpact` 中添加特殊处理逻辑

### 调整现有影响值

直接修改映射表中的数值，无需改动计算逻辑。

### 个性化定制

未来可以支持用户级别的影响值定制：

```go
UserFactorImpactOverrides[userID][factorType] = SignedDimensionImpact{...}
```

## 八、理论依据

### 占星学文献

1. **《Christian Astrology》** (William Lilly, 1647)
   - 行星的多面性：同一行星在不同宫位表现不同
   
2. **《Planets in Transit》** (Robert Hand)
   - 行运的复杂影响：机遇与挑战并存
   
3. **《Saturn: A New Look at an Old Devil》** (Liz Greene)
   - 土星的双面性：限制即成长

### 现代心理占星

"没有绝对的吉星凶星，只有能量的不同表达方式。
挑战促进成长，舒适可能导致停滞。"

## 九、测试验证

### 测试案例

1. **水星逆行期间**
   - 事业曲线 ↘️
   - 灵性曲线 ↗️
   
2. **金星入庙**
   - 所有维度曲线 ↗️
   
3. **土星四分月亮**
   - 健康曲线 ↘️
   - 灵性曲线 ↗️（压力促进成熟）

### 验证方法

```bash
# 查询水星逆行期间的分数
curl -X POST /api/calc/total-factors \
  -d '{"birthData": {...}, "queryTime": "2026-01-15T12:00:00+08:00"}'

# 检查dimensions.spiritual是否为正
# 检查dimensions.career是否为负
```

## 十、总结

这个新系统实现了占星学理论中"同一事件在不同领域有不同影响"的核心机制，使得：

1. ✅ 五个维度曲线可以独立波动
2. ✅ 符合占星学理论（如逆行对灵性有利）
3. ✅ 易于扩展和调整
4. ✅ 保持代码清晰和可维护性

这是对原系统的重大改进，解决了"一刀切"导致所有维度同向波动的问题。
