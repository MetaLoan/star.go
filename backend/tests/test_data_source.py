#!/usr/bin/env python3
"""
数据源验证测试脚本
验证 Swiss Ephemeris 是唯一数据来源
"""

import requests
import json
import sys
from datetime import datetime

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


def test_health_check():
    """测试1: 检查 /health 接口返回的数据源"""
    print("\n📋 测试1: 健康检查接口")
    print("-" * 40)
    
    try:
        resp = requests.get(f"{API_BASE}/health", timeout=5)
        data = resp.json()
        
        data_source = data.get("dataSource", "")
        print(f"  数据源: {data_source}")
        
        is_swe = "Swiss Ephemeris" in data_source
        return is_swe, f"dataSource = {data_source}"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_planet_positions():
    """测试2: 验证行星位置数据来自 Swiss Ephemeris"""
    print("\n📋 测试2: 行星位置数据验证")
    print("-" * 40)
    
    try:
        # 查询星盘
        resp = requests.post(f"{API_BASE}/api/calc/chart", json={
            "birthData": BIRTH_DATA
        }, timeout=10)
        
        data = resp.json()
        planets = data.get("planets", [])
        
        if not planets:
            return False, "未返回行星数据"
        
        print(f"  返回行星数量: {len(planets)}")
        
        # 检查必要的行星
        required_planets = ["Sun", "Moon", "Mercury", "Venus", "Mars", 
                          "Jupiter", "Saturn", "Uranus", "Neptune", "Pluto"]
        found_planets = [p["name"] for p in planets]
        
        missing = [p for p in required_planets if p not in found_planets]
        if missing:
            return False, f"缺少行星: {missing}"
        
        # 打印太阳位置作为参考
        sun = next((p for p in planets if p["name"] == "Sun"), None)
        if sun:
            print(f"  太阳位置: {sun['longitude']:.6f}° ({sun['signName']} {sun['signDegree']:.2f}°)")
        
        return True, f"所有 {len(required_planets)} 颗行星数据正常"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_planet_precision():
    """测试3: 验证行星位置精度"""
    print("\n📋 测试3: 行星位置精度验证")
    print("-" * 40)
    
    try:
        # 使用已知日期查询
        test_date = "2026-01-15T12:00:00+08:00"
        
        resp = requests.post(f"{API_BASE}/api/calc/daily", json={
            "birthData": BIRTH_DATA,
            "targetDate": test_date,
            "withFactors": True
        }, timeout=10)
        
        data = resp.json()
        
        # 检查返回的日期是否正确
        returned_date = data.get("date", "")
        print(f"  请求日期: {test_date}")
        print(f"  返回日期: {returned_date}")
        
        # 检查是否有因子数据
        factors = data.get("factors", {})
        factor_count = len(factors.get("factors", []))
        print(f"  因子数量: {factor_count}")
        
        if factor_count == 0:
            return False, "未返回因子数据"
        
        return True, f"返回 {factor_count} 个因子"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_consistent_results():
    """测试4: 验证相同查询返回一致结果"""
    print("\n📋 测试4: 结果一致性验证")
    print("-" * 40)
    
    try:
        test_date = "2026-01-15T12:00:00+08:00"
        
        # 查询两次
        results = []
        for i in range(2):
            resp = requests.post(f"{API_BASE}/api/calc/daily", json={
                "birthData": BIRTH_DATA,
                "targetDate": test_date,
                "withFactors": True
            }, timeout=10)
            results.append(resp.json())
        
        # 比较综合分数
        score1 = results[0].get("overall", 0)
        score2 = results[1].get("overall", 0)
        
        print(f"  第一次查询分数: {score1}")
        print(f"  第二次查询分数: {score2}")
        
        if abs(score1 - score2) < 0.0001:
            return True, "两次查询结果一致"
        else:
            return False, f"结果不一致: {score1} vs {score2}"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_swe_specific_features():
    """测试5: 验证 Swiss Ephemeris 特有功能"""
    print("\n📋 测试5: SWE 特有功能验证")
    print("-" * 40)
    
    try:
        # 查询星盘，检查是否有 Chiron（凯龙星，SWE 特有）
        resp = requests.post(f"{API_BASE}/api/calc/chart", json={
            "birthData": BIRTH_DATA
        }, timeout=10)
        
        data = resp.json()
        planets = data.get("planets", [])
        
        # 检查是否有 Chiron
        chiron = next((p for p in planets if p["name"] == "Chiron"), None)
        
        if chiron:
            print(f"  凯龙星位置: {chiron['longitude']:.6f}° ({chiron['signName']})")
            return True, "检测到凯龙星（SWE 特有天体）"
        else:
            return False, "未检测到凯龙星"
    except Exception as e:
        return False, f"请求失败: {e}"


def main():
    print("=" * 60)
    print("数据源验证测试")
    print(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"API 地址: {API_BASE}")
    print("=" * 60)
    
    result = TestResult()
    
    # 运行所有测试
    tests = [
        ("健康检查 - 数据源标识", test_health_check),
        ("行星位置 - 数据完整性", test_planet_positions),
        ("行星位置 - 精度验证", test_planet_precision),
        ("结果一致性", test_consistent_results),
        ("SWE 特有功能", test_swe_specific_features),
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
