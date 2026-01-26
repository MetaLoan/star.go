#!/usr/bin/env python3
"""
生命周期准确性验证测试脚本
验证因子的入相/出相时间准确
"""

import requests
import json
import sys
from datetime import datetime, timedelta
from dateutil import parser as date_parser

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


def get_factors(target_date):
    """获取指定日期的因子"""
    resp = requests.post(f"{API_BASE}/api/calc/daily", json={
        "birthData": BIRTH_DATA,
        "targetDate": target_date,
        "withFactors": True
    }, timeout=10)
    
    data = resp.json()
    factors = data.get("factors", {}).get("factors", [])
    return factors


def test_lifecycle_structure():
    """测试1: 验证生命周期结构"""
    print("\n📋 测试1: 生命周期结构验证")
    print("-" * 40)
    
    try:
        factors = get_factors("2026-01-15T12:00:00+08:00")
        
        # 检查 lifecycle 字段
        required_fields = ["startTime", "peakTime", "endTime", "duration"]
        
        factors_with_lifecycle = 0
        incomplete_lifecycle = 0
        
        for f in factors:
            lifecycle = f.get("lifecycle")
            if lifecycle:
                factors_with_lifecycle += 1
                for field in required_fields:
                    if field not in lifecycle:
                        incomplete_lifecycle += 1
                        break
        
        coverage = factors_with_lifecycle / len(factors) * 100 if factors else 0
        print(f"  有 lifecycle 的因子: {factors_with_lifecycle}/{len(factors)} ({coverage:.1f}%)")
        print(f"  结构不完整的: {incomplete_lifecycle}")
        
        if factors_with_lifecycle == 0:
            return False, "没有因子包含 lifecycle"
        
        return True, f"{factors_with_lifecycle} 个因子有 lifecycle"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_remaining_days_accuracy():
    """测试2: 验证 remainingDays 准确性（跨日期对比）"""
    print("\n📋 测试2: remainingDays 跨日期验证")
    print("-" * 40)
    
    try:
        # 查询两个日期，间隔3天
        date1 = "2026-01-10T12:00:00+08:00"
        date2 = "2026-01-13T12:00:00+08:00"
        interval_days = 3
        
        factors1 = get_factors(date1)
        factors2 = get_factors(date2)
        
        # 构建 ID 到因子的映射
        factors1_by_id = {f["id"]: f for f in factors1}
        factors2_by_id = {f["id"]: f for f in factors2}
        
        # 找到共同的因子
        common_ids = set(factors1_by_id.keys()) & set(factors2_by_id.keys())
        print(f"  日期1因子数: {len(factors1)}")
        print(f"  日期2因子数: {len(factors2)}")
        print(f"  ID 匹配数: {len(common_ids)}")
        
        # 验证 remainingDays 差值
        correct_count = 0
        tested_count = 0
        
        for fid in common_ids:
            f1 = factors1_by_id[fid]
            f2 = factors2_by_id[fid]
            
            rd1 = f1.get("remainingDays", 0)
            rd2 = f2.get("remainingDays", 0)
            
            # 只测试有效的 remainingDays
            if rd1 > 0 and rd2 >= 0:
                tested_count += 1
                diff = rd1 - rd2
                
                if abs(diff - interval_days) < 0.2:
                    correct_count += 1
                else:
                    print(f"  ⚠️ {f1.get('name')}: 差值 {diff:.2f} (预期 {interval_days})")
        
        if tested_count == 0:
            print("  没有可验证的因子对")
            return True, "无可验证因子对"
        
        accuracy = correct_count / tested_count * 100
        print(f"  验证通过: {correct_count}/{tested_count} ({accuracy:.1f}%)")
        
        if accuracy >= 80:
            return True, f"准确率 {accuracy:.1f}%"
        else:
            return False, f"准确率不足: {accuracy:.1f}%"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_peak_time_consistency():
    """测试3: 验证 peakTime 一致性"""
    print("\n📋 测试3: peakTime 一致性验证")
    print("-" * 40)
    
    try:
        # 查询两个日期
        date1 = "2026-01-10T12:00:00+08:00"
        date2 = "2026-01-11T12:00:00+08:00"
        
        factors1 = get_factors(date1)
        factors2 = get_factors(date2)
        
        # 构建 ID 到因子的映射
        factors1_by_id = {f["id"]: f for f in factors1}
        factors2_by_id = {f["id"]: f for f in factors2}
        
        common_ids = set(factors1_by_id.keys()) & set(factors2_by_id.keys())
        
        consistent_count = 0
        tested_count = 0
        
        for fid in common_ids:
            f1 = factors1_by_id[fid]
            f2 = factors2_by_id[fid]
            
            lc1 = f1.get("lifecycle", {})
            lc2 = f2.get("lifecycle", {})
            
            peak1 = lc1.get("peakTime", "")
            peak2 = lc2.get("peakTime", "")
            
            if peak1 and peak2:
                tested_count += 1
                
                # 解析时间并比较
                try:
                    t1 = date_parser.parse(peak1[:25])
                    t2 = date_parser.parse(peak2[:25])
                    diff_hours = abs((t2 - t1).total_seconds()) / 3600
                    
                    if diff_hours < 1:  # 差异小于1小时
                        consistent_count += 1
                except:
                    pass
        
        if tested_count == 0:
            return True, "无可验证因子对"
        
        consistency = consistent_count / tested_count * 100
        print(f"  peakTime 一致的因子: {consistent_count}/{tested_count} ({consistency:.1f}%)")
        
        if consistency >= 70:
            return True, f"一致性 {consistency:.1f}%"
        else:
            return False, f"一致性不足: {consistency:.1f}%"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_lifecycle_time_order():
    """测试4: 验证生命周期时间顺序"""
    print("\n📋 测试4: 时间顺序验证")
    print("-" * 40)
    
    try:
        factors = get_factors("2026-01-15T12:00:00+08:00")
        
        valid_count = 0
        invalid_count = 0
        
        for f in factors:
            lifecycle = f.get("lifecycle", {})
            
            start = lifecycle.get("startTime", "")
            peak = lifecycle.get("peakTime", "")
            end = lifecycle.get("endTime", "")
            
            if start and peak and end:
                try:
                    t_start = date_parser.parse(start[:25])
                    t_peak = date_parser.parse(peak[:25])
                    t_end = date_parser.parse(end[:25])
                    
                    # 检查顺序：start <= peak <= end
                    if t_start <= t_peak <= t_end:
                        valid_count += 1
                    else:
                        invalid_count += 1
                        print(f"  ⚠️ 时间顺序错误: {f.get('name')}")
                except:
                    pass
        
        print(f"  时间顺序正确: {valid_count}")
        print(f"  时间顺序错误: {invalid_count}")
        
        if invalid_count == 0:
            return True, f"所有 {valid_count} 个因子时间顺序正确"
        else:
            return False, f"{invalid_count} 个因子时间顺序错误"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_duration_validity():
    """测试5: 验证持续时间有效性"""
    print("\n📋 测试5: 持续时间验证")
    print("-" * 40)
    
    try:
        factors = get_factors("2026-01-15T12:00:00+08:00")
        
        valid_count = 0
        invalid_count = 0
        
        for f in factors:
            lifecycle = f.get("lifecycle", {})
            duration = lifecycle.get("duration", 0)
            
            if duration:
                # 持续时间应该是正数，且在合理范围内
                # 注意：慢行星（如冥王星、凯龙星）相位可持续数年，约 30000 小时
                if 0.1 <= duration <= 50000:
                    valid_count += 1
                else:
                    invalid_count += 1
                    print(f"  ⚠️ 异常持续时间: {f.get('name')} = {duration}小时")
        
        print(f"  有效持续时间: {valid_count}")
        print(f"  异常持续时间: {invalid_count}")
        
        if invalid_count == 0:
            return True, f"所有 {valid_count} 个因子持续时间有效"
        else:
            return False, f"{invalid_count} 个因子持续时间异常"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_id_format():
    """测试6: 验证因子 ID 格式"""
    print("\n📋 测试6: 因子 ID 格式验证")
    print("-" * 40)
    
    try:
        factors = get_factors("2026-01-15T12:00:00+08:00")
        
        # 检查相位因子的 ID 格式
        aspect_factors = [f for f in factors if f.get("type") == "aspectPhase"]
        
        valid_format_count = 0
        for f in aspect_factors:
            fid = f.get("id", "")
            # 预期格式: aspectPhase_Name_YYYYMMDD_HH
            parts = fid.split("_")
            if len(parts) >= 4:
                # 检查日期部分
                date_part = parts[-2]
                hour_part = parts[-1]
                if len(date_part) == 8 and date_part.isdigit() and len(hour_part) == 2:
                    valid_format_count += 1
        
        coverage = valid_format_count / len(aspect_factors) * 100 if aspect_factors else 0
        print(f"  相位因子数: {len(aspect_factors)}")
        print(f"  符合格式的: {valid_format_count} ({coverage:.1f}%)")
        
        if coverage >= 90:
            return True, f"ID 格式正确率 {coverage:.1f}%"
        else:
            return False, f"ID 格式正确率不足: {coverage:.1f}%"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_remaining_days_calculation():
    """测试7: 验证 remainingDays 计算逻辑"""
    print("\n📋 测试7: remainingDays 计算验证")
    print("-" * 40)
    
    try:
        query_date = "2026-01-15T12:00:00+08:00"
        factors = get_factors(query_date)
        
        query_time = date_parser.parse(query_date)
        
        correct_count = 0
        tested_count = 0
        
        for f in factors:
            lifecycle = f.get("lifecycle", {})
            remaining = f.get("remainingDays", 0)
            end_time_str = lifecycle.get("endTime", "")
            
            if end_time_str and remaining > 0:
                try:
                    end_time = date_parser.parse(end_time_str[:25])
                    
                    # 计算预期的 remainingDays
                    expected = (end_time - query_time).total_seconds() / 86400
                    
                    if expected > 0:
                        tested_count += 1
                        # 允许小误差
                        if abs(remaining - expected) < 0.1:
                            correct_count += 1
                        else:
                            print(f"  ⚠️ {f.get('name')[:30]}: {remaining:.2f} vs {expected:.2f}")
                except:
                    pass
        
        if tested_count == 0:
            return True, "无可验证因子"
        
        accuracy = correct_count / tested_count * 100
        print(f"  计算正确: {correct_count}/{tested_count} ({accuracy:.1f}%)")
        
        if accuracy >= 80:
            return True, f"计算准确率 {accuracy:.1f}%"
        else:
            return False, f"计算准确率不足: {accuracy:.1f}%"
    except Exception as e:
        return False, f"请求失败: {e}"


def main():
    print("=" * 60)
    print("生命周期准确性验证测试")
    print(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"API 地址: {API_BASE}")
    print("=" * 60)
    
    result = TestResult()
    
    # 运行所有测试
    tests = [
        ("生命周期结构验证", test_lifecycle_structure),
        ("remainingDays 跨日期验证", test_remaining_days_accuracy),
        ("peakTime 一致性验证", test_peak_time_consistency),
        ("时间顺序验证", test_lifecycle_time_order),
        ("持续时间验证", test_duration_validity),
        ("因子 ID 格式验证", test_id_format),
        ("remainingDays 计算验证", test_remaining_days_calculation),
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
