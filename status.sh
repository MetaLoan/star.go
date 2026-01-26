#!/bin/bash

# ============================================
# Star API 服务状态查看工具
# ============================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║       🌟 Star API 服务状态 🌟              ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""

# 检查后端
echo -e "${CYAN}📡 后端 API 服务${NC}"
if pgrep -f "star-api" > /dev/null 2>&1; then
    BACKEND_PID=$(pgrep -f "star-api" | head -1)
    echo -e "   状态: ${GREEN}运行中${NC} (PID: $BACKEND_PID)"
    echo -e "   本地: ${YELLOW}http://localhost:8080${NC}"
    
    # 健康检查
    HEALTH=$(curl -s http://localhost:8080/health 2>/dev/null)
    if [ -n "$HEALTH" ]; then
        VERSION=$(echo "$HEALTH" | python3 -c "import sys, json; print(json.load(sys.stdin).get('version', 'unknown'))" 2>/dev/null || echo "unknown")
        DATA_SOURCE=$(echo "$HEALTH" | python3 -c "import sys, json; print(json.load(sys.stdin).get('dataSource', 'unknown'))" 2>/dev/null || echo "unknown")
        echo -e "   版本: ${GREEN}$VERSION${NC}"
        echo -e "   数据源: ${GREEN}$DATA_SOURCE${NC}"
    fi
else
    echo -e "   状态: ${RED}未运行${NC}"
fi
echo ""

# 检查前端
echo -e "${CYAN}🖥️  前端服务${NC}"
if pgrep -f "vite" > /dev/null 2>&1; then
    FRONTEND_PID=$(pgrep -f "vite" | head -1)
    echo -e "   状态: ${GREEN}运行中${NC} (PID: $FRONTEND_PID)"
    echo -e "   本地: ${YELLOW}http://localhost:5173${NC}"
else
    echo -e "   状态: ${RED}未运行${NC}"
fi
echo ""

# 检查 ngrok
echo -e "${CYAN}🌐 公网隧道 (ngrok)${NC}"
if pgrep -f "ngrok" > /dev/null 2>&1; then
    NGROK_PID=$(pgrep -f "ngrok" | head -1)
    echo -e "   状态: ${GREEN}运行中${NC} (PID: $NGROK_PID)"
    
    # 获取公网 URL
    NGROK_URL=$(curl -s http://127.0.0.1:4040/api/tunnels 2>/dev/null | python3 -c "import sys, json; data = json.load(sys.stdin); tunnels = data.get('tunnels', []); print(tunnels[0]['public_url'] if tunnels else '')" 2>/dev/null || echo "")
    if [ -n "$NGROK_URL" ]; then
        echo -e "   公网: ${YELLOW}$NGROK_URL${NC}"
    fi
    echo -e "   管理: ${YELLOW}http://localhost:4040${NC}"
else
    echo -e "   状态: ${RED}未运行${NC}"
fi
echo ""

# 显示服务信息文件（如果存在）
if [ -f "$SCRIPT_DIR/.service-info" ]; then
    echo -e "${CYAN}📋 上次启动信息${NC}"
    echo -e "   $(head -1 "$SCRIPT_DIR/.service-info" | sed 's/# //')"
fi
echo ""
