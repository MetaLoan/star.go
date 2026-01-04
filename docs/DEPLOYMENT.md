# 部署指南

## 1. 环境要求

### 1.1 后端 (Go)

| 依赖 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 编译运行 |
| Git | - | 代码管理 |

### 1.2 前端 (Node.js)

| 依赖 | 版本 | 说明 |
|------|------|------|
| Node.js | 18+ | 运行环境 |
| npm/pnpm | - | 包管理 |
| Git | - | 代码管理 |

---

## 2. 本地开发

### 2.1 克隆代码

```bash
git clone <repository-url>
cd Star
```

### 2.2 启动后端

```bash
# 进入 Go 项目目录
cd Star.go

# 编译
go build -o bin/star-api .

# 运行
./bin/star-api

# 或者一步到位
go run .
```

后端服务运行在 `http://localhost:8080`

### 2.3 启动前端

```bash
# 进入 Next.js 项目目录
cd Star

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务运行在 `http://localhost:3000`

### 2.4 环境变量

创建 `.env.local` 文件：

```bash
# Star/.env.local

# API 地址（客户端）
NEXT_PUBLIC_API_URL=http://localhost:8080

# API 地址（服务端）
API_URL=http://localhost:8080
```

---

## 3. 生产构建

### 3.1 后端构建

```bash
cd Star.go

# 构建生产版本
go build -ldflags="-s -w" -o bin/star-api .

# 或使用 Makefile
make build
```

### 3.2 前端构建

```bash
cd Star

# 构建生产版本
npm run build

# 启动生产服务器
npm start
```

---

## 4. Docker 部署

### 4.1 后端 Dockerfile

```dockerfile
# Star.go/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o star-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/star-api .

EXPOSE 8080
CMD ["./star-api"]
```

### 4.2 前端 Dockerfile

```dockerfile
# Star/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app

ENV NODE_ENV production

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/public ./public

EXPOSE 3000
CMD ["npm", "start"]
```

### 4.3 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build:
      context: ./Star.go
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  web:
    build:
      context: ./Star
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://api:8080
      - API_URL=http://api:8080
    depends_on:
      api:
        condition: service_healthy
    restart: unless-stopped
```

### 4.4 运行 Docker

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

---

## 5. 云平台部署

### 5.1 Vercel (前端)

1. 连接 GitHub 仓库
2. 选择 Star 目录
3. 配置环境变量：
   - `NEXT_PUBLIC_API_URL`: 后端 API 地址
4. 部署

```bash
# 或使用 CLI
cd Star
npx vercel
```

### 5.2 Railway/Render (后端)

1. 连接 GitHub 仓库
2. 选择 Star.go 目录
3. 配置：
   - 构建命令: `go build -o bin/star-api .`
   - 启动命令: `./bin/star-api`
   - 端口: 8080
4. 部署

### 5.3 自建服务器

#### 使用 systemd

```ini
# /etc/systemd/system/star-api.service
[Unit]
Description=Star API Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/star
ExecStart=/opt/star/star-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# 启用并启动服务
sudo systemctl enable star-api
sudo systemctl start star-api

# 查看状态
sudo systemctl status star-api
```

#### 使用 PM2 (前端)

```bash
# 安装 PM2
npm install -g pm2

# 构建并启动
cd Star
npm run build
pm2 start npm --name "star-web" -- start

# 保存配置
pm2 save
pm2 startup
```

---

## 6. 反向代理

### 6.1 Nginx

```nginx
# /etc/nginx/sites-available/star
server {
    listen 80;
    server_name star.example.com;

    # 前端
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # 后端 API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 健康检查
    location /health {
        proxy_pass http://localhost:8080/health;
    }
}
```

### 6.2 HTTPS (Let's Encrypt)

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d star.example.com

# 自动续期
sudo certbot renew --dry-run
```

---

## 7. 监控与日志

### 7.1 日志配置

**后端日志**:
```bash
# 输出到文件
./star-api 2>&1 | tee /var/log/star-api.log

# 使用 journalctl (systemd)
journalctl -u star-api -f
```

**前端日志**:
```bash
# PM2 日志
pm2 logs star-web
```

### 7.2 健康检查

```bash
# 定时检查
curl -s http://localhost:8080/health | jq '.status'

# 监控脚本
#!/bin/bash
STATUS=$(curl -s http://localhost:8080/health | jq -r '.status')
if [ "$STATUS" != "ok" ]; then
    echo "API unhealthy!"
    # 发送告警
fi
```

### 7.3 性能监控

推荐工具：
- **Prometheus + Grafana**: 指标收集和可视化
- **Uptime Kuma**: 简单的可用性监控
- **Sentry**: 错误追踪

---

## 8. 备份与恢复

### 8.1 数据备份

当前版本使用内存存储，无需数据备份。

未来数据库版本需要：
```bash
# PostgreSQL
pg_dump star_db > backup.sql

# 恢复
psql star_db < backup.sql
```

### 8.2 代码备份

```bash
# Git 仓库
git push origin main

# 定期归档
tar -czf star-backup-$(date +%Y%m%d).tar.gz Star Star.go
```

---

## 9. 安全配置

### 9.1 防火墙

```bash
# UFW
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 9.2 CORS 配置

后端已配置 CORS，如需修改：

```go
// api/routes.go
func RegisterRoutes(r *gin.Engine) {
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://star.example.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
}
```

### 9.3 速率限制

```go
// 使用 gin-ratelimit
import "github.com/gin-contrib/ratelimit"

r.Use(ratelimit.RateLimiter(
    time.Second, // 时间窗口
    100,         // 每窗口请求数
))
```

---

## 10. 常见问题

### 10.1 端口冲突

```bash
# 查找占用端口的进程
lsof -i :8080
lsof -i :3000

# 终止进程
kill -9 <PID>
```

### 10.2 依赖问题

```bash
# 清理 Go 模块缓存
go clean -modcache
go mod tidy

# 清理 Node 模块
rm -rf node_modules package-lock.json
npm install
```

### 10.3 内存不足

```bash
# 限制 Node.js 内存
export NODE_OPTIONS="--max-old-space-size=512"

# Go 程序内存分析
go tool pprof http://localhost:8080/debug/pprof/heap
```

### 10.4 时区问题

```bash
# 容器中设置时区
TZ=Asia/Shanghai

# Docker
docker run -e TZ=Asia/Shanghai ...
```

---

## 11. 版本更新

### 11.1 后端更新

```bash
cd Star.go
git pull
go build -o bin/star-api .
sudo systemctl restart star-api
```

### 11.2 前端更新

```bash
cd Star
git pull
npm install
npm run build
pm2 restart star-web
```

### 11.3 零停机更新

使用蓝绿部署或滚动更新策略。

---

## 12. 联系支持

- 问题反馈: GitHub Issues
- 文档更新: Pull Request

