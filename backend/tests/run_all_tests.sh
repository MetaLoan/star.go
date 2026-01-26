#!/bin/bash
#
# 数据流准确性测试 - 一键运行脚本
# 运行所有测试并生成报告
#

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 配置
API_BASE="${API_BASE:-http://localhost:8080}"
REPORT_FILE="test_report_$(date +%Y%m%d_%H%M%S).md"

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  数据流准确性测试套件${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "测试时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo -e "API 地址: ${API_BASE}"
echo ""

# 检查 API 是否可用
echo -e "${YELLOW}检查 API 服务...${NC}"
if curl -s "${API_BASE}/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ API 服务正常${NC}"
else
    echo -e "${RED}✗ API 服务不可用，请先启动服务${NC}"
    echo "  运行: cd ../backend && ./bin/star-api"
    exit 1
fi
echo ""

# 检查 Python 依赖
echo -e "${YELLOW}检查 Python 依赖...${NC}"
python3 -c "import requests" 2>/dev/null || {
    echo -e "${YELLOW}安装 requests 库...${NC}"
    pip3 install requests -q
}
python3 -c "import dateutil" 2>/dev/null || {
    echo -e "${YELLOW}安装 python-dateutil 库...${NC}"
    pip3 install python-dateutil -q
}
echo -e "${GREEN}✓ Python 依赖就绪${NC}"
echo ""

# 初始化报告
cat > "$REPORT_FILE" << EOF
# 数据流准确性测试报告

**测试时间**: $(date '+%Y-%m-%d %H:%M:%S')  
**API 地址**: ${API_BASE}

---

EOF

# 测试列表
TESTS=(
    "test_data_source.py:数据源验证"
    "test_aspects.py:相位检测验证"
    "test_factors.py:因子完整性验证"
    "test_lifecycle.py:生命周期验证"
    "test_scores.py:五维分数验证"
)

TOTAL=0
PASSED=0
FAILED=0

# 运行每个测试
for test_info in "${TESTS[@]}"; do
    IFS=':' read -r test_file test_name <<< "$test_info"
    
    echo -e "${BLUE}============================================${NC}"
    echo -e "${BLUE}运行: ${test_name}${NC}"
    echo -e "${BLUE}============================================${NC}"
    
    TOTAL=$((TOTAL + 1))
    
    # 添加到报告
    echo "## ${test_name}" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo '```' >> "$REPORT_FILE"
    
    # 运行测试并捕获输出
    if python3 "$test_file" 2>&1 | tee -a "$REPORT_FILE"; then
        echo -e "${GREEN}✓ ${test_name} 通过${NC}"
        PASSED=$((PASSED + 1))
        echo '```' >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
        echo "**结果**: ✅ 通过" >> "$REPORT_FILE"
    else
        echo -e "${RED}✗ ${test_name} 失败${NC}"
        FAILED=$((FAILED + 1))
        echo '```' >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
        echo "**结果**: ❌ 失败" >> "$REPORT_FILE"
    fi
    
    echo "" >> "$REPORT_FILE"
    echo "---" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo ""
done

# 生成总结
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  测试总结${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "总测试数: ${TOTAL}"
echo -e "通过: ${GREEN}${PASSED}${NC}"
echo -e "失败: ${RED}${FAILED}${NC}"
echo ""

# 添加总结到报告
cat >> "$REPORT_FILE" << EOF
## 总结

| 指标 | 数值 |
|------|------|
| 总测试数 | ${TOTAL} |
| 通过 | ${PASSED} |
| 失败 | ${FAILED} |
| 通过率 | $(echo "scale=1; ${PASSED}*100/${TOTAL}" | bc)% |

EOF

# 最终结果
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}============================================${NC}"
    echo -e "${GREEN}  所有测试通过！${NC}"
    echo -e "${GREEN}============================================${NC}"
    echo "**最终结果**: ✅ 所有测试通过" >> "$REPORT_FILE"
else
    echo -e "${RED}============================================${NC}"
    echo -e "${RED}  有 ${FAILED} 个测试失败${NC}"
    echo -e "${RED}============================================${NC}"
    echo "**最终结果**: ❌ ${FAILED} 个测试失败" >> "$REPORT_FILE"
fi

echo ""
echo -e "测试报告已保存到: ${YELLOW}${REPORT_FILE}${NC}"
echo ""

exit $FAILED
