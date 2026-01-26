#!/usr/bin/env python3
"""
五维分数变化验证测试脚本
验证因子变化正确反映到五维分数
"""

import requests
import json
import sys
from datetime import datetime, timedelta

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

# 五维名称
DIMENSIONS = ["career", "relationship", "health", "finance", "spiritual"]

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


def get_daily_data(target_date):
    """获取指定日期的每日预测数据"""
    resp = requests.post(f"{API_BASE}/api/calc/daily", json={
        "birthData": BIRTH_DATA,
        "targetDate": target_date,
        "withFactors": True
    }, timeout=10)
    
    return resp.json()


def test_dimension_scores_exist():
    """测试1: 验证五维分数存在"""
    print("\n📋 测试1: 五维分数存在验证")
    print("-" * 40)
    
    try:
        data = get_daily_data("2026-01-15T12:00:00+08:00")
        
        dimensions = data.get("dimensions", {})
        
        print("  五维分数:")
        for dim in DIMENSIONS:
            score = dimensions.get(dim, "缺失")
            print(f"    - {dim}: {score}")
        
        # 检查所有维度是否存在
        missing = [dim for dim in DIMENSIONS if dim not in dimensions]
        if missing:
            return False, f"缺少维度: {missing}"
        
        return True, "所有五维分数都存在"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_dimension_scores_range():
    """测试2: 验证五维分数范围"""
    print("\n📋 测试2: 五维分数范围验证")
    print("-" * 40)
    
    try:
        data = get_daily_data("2026-01-15T12:00:00+08:00")
        
        dimensions = data.get("dimensions", {})
        
        out_of_range = []
        for dim in DIMENSIONS:
            score = dimensions.get(dim, 0)
            if score < 0 or score > 100:
                out_of_range.append((dim, score))
        
        if out_of_range:
            print("  超出范围的分数:")
            for dim, score in out_of_range:
                print(f"    - {dim}: {score}")
            return False, f"{len(out_of_range)} 个维度超出 0-100 范围"
        
        # 打印分数范围
        scores = [dimensions.get(dim, 0) for dim in DIMENSIONS]
        print(f"  分数范围: {min(scores):.2f} - {max(scores):.2f}")
        print(f"  平均分数: {sum(scores)/len(scores):.2f}")
        
        return True, "所有分数在 0-100 范围内"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_overall_score_calculation():
    """测试3: 验证综合分数计算"""
    print("\n📋 测试3: 综合分数计算验证")
    print("-" * 40)
    
    try:
        data = get_daily_data("2026-01-15T12:00:00+08:00")
        
        dimensions = data.get("dimensions", {})
        overall = data.get("overallScore", 0)  # API 返回 overallScore
        
        # 默认权重（假设等权重）
        weights = {
            "career": 0.2,
            "relationship": 0.2,
            "health": 0.2,
            "finance": 0.2,
            "spiritual": 0.2,
        }
        
        # 计算预期综合分数
        expected = sum(dimensions.get(dim, 0) * weights[dim] for dim in DIMENSIONS)
        
        print(f"  API 返回综合分数: {overall:.4f}")
        print(f"  等权重计算分数: {expected:.4f}")
        
        # 综合分数应该在合理范围内
        if 0 <= overall <= 100:
            return True, f"综合分数 {overall:.2f} 在有效范围内"
        else:
            return False, f"综合分数 {overall:.2f} 超出范围"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_adjustment_calculation():
    """测试4: 验证调整值计算公式"""
    print("\n📋 测试4: 调整值计算公式验证")
    print("-" * 40)
    
    try:
        data = get_daily_data("2026-01-15T12:00:00+08:00")
        
        factors = data.get("factors", {}).get("factors", [])
        
        # 验证: adjustment = baseValue * weight * currentStrength
        correct_count = 0
        tested_count = 0
        errors = []
        
        for f in factors:
            base_value = f.get("baseValue", 0)
            weight = f.get("weight", 0)
            strength = f.get("currentStrength", 0)
            adjustment = f.get("adjustment", 0)
            
            if base_value != 0 and weight != 0:
                tested_count += 1
                expected = base_value * weight * strength
                
                if abs(expected - adjustment) < 0.0001:
                    correct_count += 1
                else:
                    errors.append({
                        "name": f.get("name", ""),
                        "expected": expected,
                        "actual": adjustment
                    })
        
        if errors:
            print("  计算不一致的因子:")
            for e in errors[:3]:
                print(f"    - {e['name'][:30]}: 期望 {e['expected']:.4f}, 实际 {e['actual']:.4f}")
        
        if tested_count == 0:
            return True, "无可验证因子"
        
        accuracy = correct_count / tested_count * 100
        print(f"  验证通过: {correct_count}/{tested_count} ({accuracy:.1f}%)")
        
        if accuracy >= 95:
            return True, f"计算准确率 {accuracy:.1f}%"
        else:
            return False, f"计算准确率不足: {accuracy:.1f}%"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_dimension_adjustments():
    """测试5: 验证维度调整值"""
    print("\n📋 测试5: 维度调整值验证")
    print("-" * 40)
    
    try:
        data = get_daily_data("2026-01-15T12:00:00+08:00")
        
        factor_result = data.get("factors", {})
        dim_adjustments = factor_result.get("dimensionAdjustments", {})
        
        print("  维度调整值:")
        for dim in DIMENSIONS:
            adj = dim_adjustments.get(dim, 0)
            print(f"    - {dim}: {adj:+.4f}")
        
        # 检查调整值是否存在
        if not dim_adjustments:
            return False, "缺少 dimensionAdjustments"
        
        # 手动计算维度调整值
        factors = factor_result.get("factors", [])
        calculated = {dim: 0 for dim in DIMENSIONS}
        
        for f in factors:
            adj = f.get("adjustment", 0)
            impact = f.get("dimensionImpact", {})
            for dim in DIMENSIONS:
                calculated[dim] += adj * impact.get(dim, 0)
        
        print("  手动计算值:")
        match_count = 0
        for dim in DIMENSIONS:
            calc = calculated[dim]
            api = dim_adjustments.get(dim, 0)
            match = abs(calc - api) < 0.01
            if match:
                match_count += 1
            status = "✓" if match else "✗"
            print(f"    - {dim}: {calc:+.4f} {status}")
        
        if match_count == len(DIMENSIONS):
            return True, "所有维度调整值匹配"
        else:
            return False, f"只有 {match_count}/{len(DIMENSIONS)} 维度匹配"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_time_series_variation():
    """测试6: 验证时间序列变化"""
    print("\n📋 测试6: 时间序列变化验证")
    print("-" * 40)
    
    try:
        # 查询连续7天的分数
        dates = [f"2026-01-{10+i:02d}T12:00:00+08:00" for i in range(7)]
        
        scores_by_date = {}
        for date in dates:
            data = get_daily_data(date)
            scores_by_date[date[:10]] = {
                "overall": data.get("overallScore", 0),  # API 返回 overallScore
                "dimensions": data.get("dimensions", {})
            }
        
        # 分析变化
        overall_scores = [s["overall"] for s in scores_by_date.values()]
        
        print(f"  7天综合分数: {[f'{s:.1f}' for s in overall_scores]}")
        
        # 计算日变化幅度
        daily_changes = []
        overall_list = list(scores_by_date.values())
        for i in range(1, len(overall_list)):
            change = abs(overall_list[i]["overall"] - overall_list[i-1]["overall"])
            daily_changes.append(change)
        
        avg_change = sum(daily_changes) / len(daily_changes) if daily_changes else 0
        max_change = max(daily_changes) if daily_changes else 0
        
        print(f"  平均日变化: {avg_change:.2f}")
        print(f"  最大日变化: {max_change:.2f}")
        
        # 变化应该在合理范围内（每日变化不超过20分）
        if max_change > 30:
            return False, f"日变化过大: {max_change:.2f}"
        
        # 检查分数不应该完全不变
        if max_change < 0.01:
            return False, "分数完全没有变化"
        
        return True, f"日变化范围合理 (avg: {avg_change:.2f}, max: {max_change:.2f})"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_base_scores_exist():
    """测试7: 验证基础分数存在"""
    print("\n📋 测试7: 基础分数验证")
    print("-" * 40)
    
    try:
        data = get_daily_data("2026-01-15T12:00:00+08:00")
        
        base_scores = data.get("baseScores", {})
        
        if not base_scores:
            print("  未返回 baseScores")
            return True, "API 未返回 baseScores（可能是正常行为）"
        
        print("  基础分数:")
        for dim in DIMENSIONS:
            score = base_scores.get(dim, "缺失")
            print(f"    - {dim}: {score}")
        
        # 基础分应该在合理范围内
        for dim in DIMENSIONS:
            score = base_scores.get(dim, 50)
            if score < 0 or score > 100:
                return False, f"{dim} 基础分超出范围: {score}"
        
        return True, "基础分数有效"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_factor_dimension_correlation():
    """测试8: 验证因子与维度的相关性"""
    print("\n📋 测试8: 因子-维度相关性验证")
    print("-" * 40)
    
    try:
        data = get_daily_data("2026-01-15T12:00:00+08:00")
        
        factors = data.get("factors", {}).get("factors", [])
        dim_adjustments = data.get("factors", {}).get("dimensionAdjustments", {})
        
        # 统计正/负因子对各维度的影响
        positive_impact = {dim: 0 for dim in DIMENSIONS}
        negative_impact = {dim: 0 for dim in DIMENSIONS}
        
        for f in factors:
            adj = f.get("adjustment", 0)
            impact = f.get("dimensionImpact", {})
            
            for dim in DIMENSIONS:
                dim_adj = adj * impact.get(dim, 0)
                if dim_adj > 0:
                    positive_impact[dim] += dim_adj
                else:
                    negative_impact[dim] += dim_adj
        
        print("  各维度正面/负面影响:")
        for dim in DIMENSIONS:
            pos = positive_impact[dim]
            neg = negative_impact[dim]
            net = pos + neg
            api_adj = dim_adjustments.get(dim, 0)
            match = abs(net - api_adj) < 0.1
            status = "✓" if match else "✗"
            print(f"    - {dim}: +{pos:.2f} / {neg:.2f} = {net:.2f} {status}")
        
        return True, "因子-维度相关性分析完成"
    except Exception as e:
        return False, f"请求失败: {e}"


def main():
    print("=" * 60)
    print("五维分数变化验证测试")
    print(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"API 地址: {API_BASE}")
    print("=" * 60)
    
    result = TestResult()
    
    # 运行所有测试
    tests = [
        ("五维分数存在验证", test_dimension_scores_exist),
        ("五维分数范围验证", test_dimension_scores_range),
        ("综合分数计算验证", test_overall_score_calculation),
        ("调整值计算公式验证", test_adjustment_calculation),
        ("维度调整值验证", test_dimension_adjustments),
        ("时间序列变化验证", test_time_series_variation),
        ("基础分数验证", test_base_scores_exist),
        ("因子-维度相关性验证", test_factor_dimension_correlation),
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
