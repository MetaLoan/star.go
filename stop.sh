#!/bin/bash

# ============================================
# Star API 服务停止工具
# ============================================

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo ""
echo -e "${CYAN}🛑 正在停止 Star API 服务...${NC}"
echo ""

# 停止后端
if pgrep -f "star-api" > /dev/null 2>&1; then
    pkill -f "star-api" 2>/dev/null
    echo -e "${GREEN}✅ 后端服务已停止${NC}"
else
    echo -e "${YELLOW}⚠️  后端服务未运行${NC}"
fi

# 停止前端
if pgrep -f "vite" > /dev/null 2>&1; then
    pkill -f "vite" 2>/dev/null
    echo -e "${GREEN}✅ 前端服务已停止${NC}"
else
    echo -e "${YELLOW}⚠️  前端服务未运行${NC}"
fi

# 停止 ngrok
if pgrep -f "ngrok" > /dev/null 2>&1; then
    pkill -f "ngrok" 2>/dev/null
    echo -e "${GREEN}✅ ngrok 隧道已停止${NC}"
else
    echo -e "${YELLOW}⚠️  ngrok 隧道未运行${NC}"
fi

echo ""
echo -e "${GREEN}🎉 所有服务已停止${NC}"
echo ""
