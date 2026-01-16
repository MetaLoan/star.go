#!/bin/bash

# ============================================
# Star API 服务重启工具
# 用于重启后端、前端和公网隧道服务
# ============================================

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

# 显示标题
echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║       🌟 Star API 服务重启工具 🌟          ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""

# ==================== 停止现有服务 ====================
print_info "正在停止现有服务..."

# 停止后端
if pgrep -f "star-api" > /dev/null 2>&1; then
    pkill -f "star-api" 2>/dev/null || true
    print_success "后端服务已停止"
else
    print_warning "后端服务未运行"
fi

# 停止前端
if pgrep -f "vite" > /dev/null 2>&1; then
    pkill -f "vite" 2>/dev/null || true
    print_success "前端服务已停止"
else
    print_warning "前端服务未运行"
fi

# 停止 ngrok
if pgrep -f "ngrok" > /dev/null 2>&1; then
    pkill -f "ngrok" 2>/dev/null || true
    print_success "ngrok 隧道已停止"
else
    print_warning "ngrok 隧道未运行"
fi

sleep 2

# ==================== 启动后端 ====================
print_info "正在启动后端服务..."

if [ ! -f "$BACKEND_DIR/bin/star-api" ]; then
    print_warning "后端二进制文件不存在，正在编译..."
    cd "$BACKEND_DIR"
    if [ -f "build.sh" ]; then
        chmod +x build.sh
        ./build.sh > /dev/null 2>&1
    else
        CGO_ENABLED=1 go build -tags swe -o bin/star-api . 2>/dev/null
    fi
fi

cd "$BACKEND_DIR"
nohup ./bin/star-api > server.log 2>&1 &
BACKEND_PID=$!

# 等待后端启动
sleep 3
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    print_success "后端服务已启动 (PID: $BACKEND_PID)"
else
    print_error "后端服务启动失败，请检查 $BACKEND_DIR/server.log"
    cat "$BACKEND_DIR/server.log" | tail -10
    exit 1
fi

# ==================== 启动前端 ====================
print_info "正在启动前端服务..."

cd "$FRONTEND_DIR"

# 检查依赖
if [ ! -d "node_modules" ]; then
    print_warning "正在安装前端依赖..."
    npm install > /dev/null 2>&1
fi

nohup npm run dev > /tmp/frontend.log 2>&1 &
FRONTEND_PID=$!

# 等待前端启动
sleep 3
if lsof -ti:5173 > /dev/null 2>&1; then
    print_success "前端服务已启动 (PID: $FRONTEND_PID)"
else
    print_warning "前端服务启动中，请稍候..."
fi

# ==================== 启动 ngrok ====================
print_info "正在启动公网隧道..."

cd "$SCRIPT_DIR"
nohup ngrok http 8080 --log=stdout > /tmp/ngrok.log 2>&1 &
NGROK_PID=$!

# 等待 ngrok 启动
sleep 5

# 获取公网 URL
NGROK_URL=""
for i in {1..10}; do
    NGROK_URL=$(curl -s http://127.0.0.1:4040/api/tunnels 2>/dev/null | python3 -c "import sys, json; data = json.load(sys.stdin); tunnels = data.get('tunnels', []); print(tunnels[0]['public_url'] if tunnels else '')" 2>/dev/null || echo "")
    if [ -n "$NGROK_URL" ]; then
        break
    fi
    sleep 1
done

if [ -n "$NGROK_URL" ]; then
    print_success "公网隧道已启动 (PID: $NGROK_PID)"
else
    print_error "公网隧道启动失败，请检查 /tmp/ngrok.log"
fi

# ==================== 显示服务信息 ====================
echo ""
echo -e "${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                    🎉 服务启动完成 🎉                          ║${NC}"
echo -e "${CYAN}╠════════════════════════════════════════════════════════════════╣${NC}"
echo -e "${CYAN}║                                                                ║${NC}"
echo -e "${CYAN}║  ${GREEN}📡 后端 API 服务${NC}                                            ${CYAN}║${NC}"
echo -e "${CYAN}║     本地地址: ${YELLOW}http://localhost:8080${NC}                          ${CYAN}║${NC}"
if [ -n "$NGROK_URL" ]; then
echo -e "${CYAN}║     公网地址: ${YELLOW}${NGROK_URL}${NC}  ${CYAN}║${NC}"
else
echo -e "${CYAN}║     公网地址: ${RED}启动失败${NC}                                       ${CYAN}║${NC}"
fi
echo -e "${CYAN}║                                                                ║${NC}"
echo -e "${CYAN}║  ${GREEN}🖥️  前端服务${NC}                                                 ${CYAN}║${NC}"
echo -e "${CYAN}║     本地地址: ${YELLOW}http://localhost:5173${NC}                          ${CYAN}║${NC}"
echo -e "${CYAN}║                                                                ║${NC}"
echo -e "${CYAN}║  ${GREEN}🔧 管理面板${NC}                                                  ${CYAN}║${NC}"
echo -e "${CYAN}║     ngrok:    ${YELLOW}http://localhost:4040${NC}                          ${CYAN}║${NC}"
echo -e "${CYAN}║                                                                ║${NC}"
echo -e "${CYAN}╠════════════════════════════════════════════════════════════════╣${NC}"
echo -e "${CYAN}║  ${GREEN}🔌 常用 API 端点${NC}                                             ${CYAN}║${NC}"
echo -e "${CYAN}║     健康检查:   GET  /health                                   ║${NC}"
echo -e "${CYAN}║     每日预测:   POST /api/calc/daily                           ║${NC}"
echo -e "${CYAN}║     时间序列:   POST /api/calc/time-series                     ║${NC}"
echo -e "${CYAN}║     分数解释:   POST /api/calc/score-explain                   ║${NC}"
echo -e "${CYAN}║     活跃因子:   POST /api/calc/active-factors                  ║${NC}"
echo -e "${CYAN}║                                                                ║${NC}"
echo -e "${CYAN}║  ${GREEN}📊 监控仪表板${NC} ⭐ 新增                                        ${CYAN}║${NC}"
echo -e "${CYAN}║     Web界面:    ${YELLOW}http://localhost:8080/api/monitor/dashboard${NC} ${CYAN}║${NC}"
echo -e "${CYAN}║     API统计:    GET  /api/monitor/summary                      ║${NC}"
echo -e "${CYAN}║                                                                ║${NC}"
echo -e "${CYAN}╠════════════════════════════════════════════════════════════════╣${NC}"
echo -e "${CYAN}║  ${YELLOW}⏹  停止服务: pkill -f 'star-api|vite|ngrok'${NC}                  ${CYAN}║${NC}"
echo -e "${CYAN}║  ${YELLOW}📋 查看日志: tail -f backend/server.log${NC}                      ${CYAN}║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# 保存服务信息到文件
cat > "$SCRIPT_DIR/.service-info" << EOF
# Star API 服务信息 - $(date '+%Y-%m-%d %H:%M:%S')
BACKEND_LOCAL=http://localhost:8080
BACKEND_PUBLIC=$NGROK_URL
FRONTEND_LOCAL=http://localhost:5173
NGROK_ADMIN=http://localhost:4040
BACKEND_PID=$BACKEND_PID
FRONTEND_PID=$FRONTEND_PID
NGROK_PID=$NGROK_PID
EOF

print_success "服务信息已保存到 .service-info"
