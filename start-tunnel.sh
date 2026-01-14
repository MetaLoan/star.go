#!/bin/bash

# Freya API 公网隧道启动脚本
# 使用 ngrok 创建公网访问地址

set -e

echo "🚀 启动 Freya API 公网隧道..."

# 检查后端是否运行
if ! pgrep -f star-api > /dev/null; then
    echo "⚠️  后端 API 未运行，正在启动..."
    cd "$(dirname "$0")/backend"
    ./bin/star-api &
    sleep 3
fi

# 停止现有 ngrok
pkill -f ngrok 2>/dev/null || true
sleep 1

# 启动 ngrok
echo "🌐 正在创建公网隧道..."
ngrok http 8080 --log=stdout > /tmp/ngrok.log 2>&1 &

# 等待 ngrok 启动
sleep 5

# 获取公网 URL
NGROK_URL=$(curl -s http://localhost:4040/api/tunnels | jq -r '.tunnels[0].public_url')

if [ "$NGROK_URL" != "null" ] && [ -n "$NGROK_URL" ]; then
    echo ""
    echo "==========================================="
    echo "✅ 公网隧道已开启！"
    echo "==========================================="
    echo ""
    echo "🌍 公网地址："
    echo "   $NGROK_URL"
    echo ""
    echo "📡 API 端点："
    echo "   $NGROK_URL/health"
    echo "   $NGROK_URL/api/calc/daily"
    echo "   $NGROK_URL/api/calc/time-series"
    echo ""
    echo "🔧 管理面板："
    echo "   http://localhost:4040"
    echo ""
    echo "⏹  停止隧道：pkill -f ngrok"
    echo "==========================================="
else
    echo "❌ 启动失败，请检查 /tmp/ngrok.log"
    cat /tmp/ngrok.log | tail -20
    exit 1
fi

