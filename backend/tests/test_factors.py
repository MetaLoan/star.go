#!/usr/bin/env python3
"""
影响因子完整性验证测试脚本
验证所有因子类型都被正确生成
"""

import requests
import json
import sys
from datetime import datetime
from collections import Counter

# 配置
API_BASE = "http://localhost:8080"
BIRTH_DATA = {
    "year": 1990,
    "month": 1,
    "day": 9,
    "hour": 10,
    "minute": 30,
    "second": 0,
    "latitude": 39.9042,
    "longitude": 116.4074,
    "timezone": 8
}

# 预期的因子类型
EXPECTED_FACTOR_TYPES = {
    "dignity": "尊贵度因子 - 行星落座状态",
    "retrograde": "逆行因子 - 行星逆行状态",
    "aspectPhase": "相位因子 - 行运-本命相位",
    "lunarPhase": "月相因子 - 月亮相位",
    "planetaryHour": "行星时因子 - 当前行星时主星",
}

class TestResult:
    def __init__(self):
        self.passed = 0
        self.failed = 0
        self.results = []
    
    def add(self, name, passed, message=""):
        self.results.append({
            "name": name,
            "passed": passed,
            "message": message
        })
        if passed:
            self.passed += 1
        else:
            self.failed += 1
    
    def summary(self):
        total = self.passed + self.failed
        print("\n" + "=" * 60)
        print(f"测试结果: {self.passed}/{total} 通过")
        print("=" * 60)
        for r in self.results:
            status = "✅ PASS" if r["passed"] else "❌ FAIL"
            print(f"  {status}: {r['name']}")
            if r["message"]:
                print(f"         {r['message']}")
        return self.failed == 0


def get_all_factors(target_date):
    """获取指定日期的所有因子"""
    resp = requests.post(f"{API_BASE}/api/calc/daily", json={
        "birthData": BIRTH_DATA,
        "targetDate": target_date,
        "withFactors": True
    }, timeout=10)
    
    data = resp.json()
    factors = data.get("factors", {}).get("factors", [])
    return factors, data


def test_factor_type_coverage():
    """测试1: 验证因子类型覆盖"""
    print("\n📋 测试1: 因子类型覆盖")
    print("-" * 40)
    
    try:
        factors, _ = get_all_factors("2026-01-15T12:00:00+08:00")
        
        # 统计因子类型
        type_counter = Counter()
        for f in factors:
            factor_type = f.get("type", "unknown")
            type_counter[factor_type] += 1
        
        print(f"  因子总数: {len(factors)}")
        print("  类型分布:")
        for factor_type, count in sorted(type_counter.items()):
            desc = EXPECTED_FACTOR_TYPES.get(factor_type, "未知类型")
            print(f"    - {factor_type}: {count} ({desc})")
        
        # 检查是否包含所有预期类型
        found_types = set(type_counter.keys())
        expected_types = set(EXPECTED_FACTOR_TYPES.keys())
        missing = expected_types - found_types
        
        if missing:
            return False, f"缺少因子类型: {missing}"
        
        return True, f"所有 {len(expected_types)} 种因子类型都已生成"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_factor_count():
    """测试2: 验证因子数量合理性"""
    print("\n📋 测试2: 因子数量合理性")
    print("-" * 40)
    
    try:
        factors, _ = get_all_factors("2026-01-15T12:00:00+08:00")
        
        count = len(factors)
        print(f"  因子总数: {count}")
        
        # 预期范围：20-60
        if count < 10:
            return False, f"因子数量过少: {count} (预期 > 10)"
        elif count > 100:
            return False, f"因子数量过多: {count} (预期 < 100)"
        
        return True, f"因子数量在合理范围内: {count}"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_essential_factors():
    """测试3: 验证必要因子存在"""
    print("\n📋 测试3: 必要因子验证")
    print("-" * 40)
    
    try:
        factors, _ = get_all_factors("2026-01-15T12:00:00+08:00")
        
        # 检查太阳和月亮相关因子
        sun_factors = [f for f in factors if "Sun" in f.get("name", "")]
        moon_factors = [f for f in factors if "Moon" in f.get("name", "")]
        
        print(f"  太阳相关因子: {len(sun_factors)}")
        print(f"  月亮相关因子: {len(moon_factors)}")
        
        if len(sun_factors) == 0:
            return False, "缺少太阳相关因子"
        if len(moon_factors) == 0:
            return False, "缺少月亮相关因子"
        
        return True, f"太阳因子 {len(sun_factors)} 个, 月亮因子 {len(moon_factors)} 个"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_factor_structure():
    """测试4: 验证因子结构完整性"""
    print("\n📋 测试4: 因子结构完整性")
    print("-" * 40)
    
    try:
        factors, _ = get_all_factors("2026-01-15T12:00:00+08:00")
        
        # 必要字段
        required_fields = ["id", "type", "name", "baseValue", "weight", 
                          "currentStrength", "adjustment", "isPositive"]
        
        missing_fields_count = {}
        for f in factors:
            for field in required_fields:
                if field not in f:
                    missing_fields_count[field] = missing_fields_count.get(field, 0) + 1
        
        if missing_fields_count:
            print("  缺失字段统计:")
            for field, count in missing_fields_count.items():
                print(f"    - {field}: {count} 个因子缺失")
            return False, f"存在缺失字段"
        
        print(f"  所有 {len(factors)} 个因子结构完整")
        return True, f"所有必要字段都存在"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_positive_negative_balance():
    """测试5: 验证正负因子平衡"""
    print("\n📋 测试5: 正负因子平衡")
    print("-" * 40)
    
    try:
        factors, data = get_all_factors("2026-01-15T12:00:00+08:00")
        
        # 统计正负因子
        positive = [f for f in factors if f.get("isPositive", False)]
        negative = [f for f in factors if not f.get("isPositive", False)]
        
        print(f"  正面因子: {len(positive)}")
        print(f"  负面因子: {len(negative)}")
        
        # 检查 factors 结构中的统计
        factor_result = data.get("factors", {})
        positive_factors = factor_result.get("positiveFactors", [])
        negative_factors = factor_result.get("negativeFactors", [])
        
        print(f"  API返回正面因子列表: {len(positive_factors)}")
        print(f"  API返回负面因子列表: {len(negative_factors)}")
        
        # 应该有正负两种因子
        if len(positive) == 0:
            return False, "没有正面因子"
        if len(negative) == 0:
            return False, "没有负面因子"
        
        ratio = len(positive) / (len(positive) + len(negative))
        print(f"  正面因子比例: {ratio:.1%}")
        
        return True, f"正/负比例: {len(positive)}:{len(negative)}"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_dimension_impact():
    """测试6: 验证维度影响分配"""
    print("\n📋 测试6: 维度影响分配")
    print("-" * 40)
    
    try:
        factors, _ = get_all_factors("2026-01-15T12:00:00+08:00")
        
        # 检查 dimensionImpact 字段
        dimensions = ["career", "relationship", "health", "finance", "spiritual"]
        
        factors_with_impact = 0
        for f in factors:
            impact = f.get("dimensionImpact", {})
            if impact:
                factors_with_impact += 1
                # 检查所有维度是否存在
                for dim in dimensions:
                    if dim not in impact:
                        print(f"  警告: 因子 {f.get('name')} 缺少 {dim} 维度")
        
        coverage = factors_with_impact / len(factors) * 100 if factors else 0
        print(f"  有维度影响的因子: {factors_with_impact}/{len(factors)} ({coverage:.1f}%)")
        
        if coverage < 80:
            return False, f"维度影响覆盖率不足: {coverage:.1f}%"
        
        return True, f"维度影响覆盖率: {coverage:.1f}%"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_total_adjustment():
    """测试7: 验证调整值计算"""
    print("\n📋 测试7: 调整值计算验证")
    print("-" * 40)
    
    try:
        factors, data = get_all_factors("2026-01-15T12:00:00+08:00")
        
        # 计算所有因子的 adjustment 总和
        total_adjustment = sum(f.get("adjustment", 0) for f in factors)
        
        # 获取 API 返回的 totalAdjustment
        factor_result = data.get("factors", {})
        api_total = factor_result.get("totalAdjustment", 0)
        
        print(f"  手动计算总调整值: {total_adjustment:.4f}")
        print(f"  API返回总调整值: {api_total:.4f}")
        
        # 允许小误差
        if abs(total_adjustment - api_total) > 0.01:
            return False, f"调整值不匹配: {total_adjustment:.4f} vs {api_total:.4f}"
        
        return True, f"调整值验证通过: {api_total:.4f}"
    except Exception as e:
        return False, f"请求失败: {e}"


def main():
    print("=" * 60)
    print("影响因子完整性验证测试")
    print(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"API 地址: {API_BASE}")
    print("=" * 60)
    
    result = TestResult()
    
    # 运行所有测试
    tests = [
        ("因子类型覆盖", test_factor_type_coverage),
        ("因子数量合理性", test_factor_count),
        ("必要因子验证", test_essential_factors),
        ("因子结构完整性", test_factor_structure),
        ("正负因子平衡", test_positive_negative_balance),
        ("维度影响分配", test_dimension_impact),
        ("调整值计算验证", test_total_adjustment),
    ]
    
    for name, test_func in tests:
        try:
            passed, message = test_func()
            result.add(name, passed, message)
        except Exception as e:
            result.add(name, False, f"异常: {e}")
    
    # 输出总结
    success = result.summary()
    
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())
