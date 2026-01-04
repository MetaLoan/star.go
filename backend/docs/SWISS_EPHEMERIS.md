# Swiss Ephemeris 高精度星历集成

## 概述

Star 项目支持可选集成 Swiss Ephemeris 库以获得更高精度的天文计算。默认情况下，系统使用内置的纯 Go 算法，精度足以满足大多数占星应用需求。

## 精度对比

| 计算项 | 内置算法 | Swiss Ephemeris |
|--------|----------|-----------------|
| 太阳位置 | < 0.1° | < 0.001° |
| 月亮位置 | < 0.5° | < 0.01° |
| 内行星 | < 1° | < 0.01° |
| 外行星 | < 2° | < 0.01° |
| 新月/满月时间 | < 6小时 | < 1分钟 |

## 何时需要集成 Swiss Ephemeris

- 需要极高精度的天文事件预测（如日月食）
- 需要计算历史久远或未来遥远的日期
- 需要恒星位置或小行星数据
- 专业天文馆或科研应用

## 安装步骤

### 1. 安装 Swiss Ephemeris C 库

#### macOS

```bash
# 使用 Homebrew
brew install swiss-ephemeris

# 或从源码编译
git clone https://github.com/aloistr/swisseph.git
cd swisseph
make
sudo make install
```

#### Linux

```bash
# Ubuntu/Debian
sudo apt-get install libswe-dev

# 或从源码
git clone https://github.com/aloistr/swisseph.git
cd swisseph
make
sudo make install
```

### 2. 下载星历表文件（可选）

Swiss Ephemeris 可以使用内置的 Moshier 算法（无需额外文件），或使用更精确的 JPL 星历表：

```bash
# 创建星历表目录
mkdir -p /usr/local/share/ephe

# 下载基础星历表
cd /usr/local/share/ephe
wget https://www.astro.com/ftp/swisseph/ephe/sepl_18.se1
wget https://www.astro.com/ftp/swisseph/ephe/semo_18.se1
```

### 3. 编译 Star 后端（启用 Swiss Ephemeris）

```bash
cd /Users/leo/Documents/Star/backend

# 使用 swe 构建标签编译
go build -tags swe -o bin/star-api .
```

### 4. 配置星历表路径

在代码中初始化：

```go
import "star/astro"

func main() {
    // 设置星历表路径
    astro.InitSwissEphemeris("/usr/local/share/ephe")
    defer astro.CloseSwissEphemeris()
    
    // ... 其他代码
}
```

## 使用方式

启用 Swiss Ephemeris 后，系统会自动使用高精度计算：

```go
// 高精度行星位置
pos := astro.CalculatePlanetPositionSwe(models.Mars, jd)

// 高精度宫位
houses, asc, mc := astro.CalculateHousesSwe(jd, lat, lon)
```

## 回退机制

如果 Swiss Ephemeris 计算失败（如超出星历表范围），系统会自动回退到内置算法，确保服务稳定性。

## 许可证说明

Swiss Ephemeris 采用双重许可：
- **GPL 2.0+**：开源项目可免费使用
- **商业许可**：闭源商业项目需购买许可

请根据您的项目性质选择合适的许可。

## 验证安装

```bash
# 运行精度测试
cd /Users/leo/Documents/Star/backend
go test -tags swe -v ./astro/... -run Accuracy
```

## 故障排除

### 链接错误 `library 'swe' not found`

确保 C 库正确安装：
```bash
# 检查库文件
ls -la /usr/local/lib/libswe*

# 检查 pkg-config
pkg-config --libs swisseph
```

### CGO 编译问题

确保设置正确的环境变量：
```bash
export CGO_ENABLED=1
export CGO_LDFLAGS="-L/usr/local/lib -lswe"
export CGO_CFLAGS="-I/usr/local/include"
```

