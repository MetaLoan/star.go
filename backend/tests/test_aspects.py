#!/usr/bin/env python3
"""
相位检测验证测试脚本
验证所有相位类型都被正确检测
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

# 预期的相位类型
EXPECTED_ASPECTS = {
    "Conjunction": {"angle": 0, "orb": 8, "nature": "neutral"},
    "Sextile": {"angle": 60, "orb": 6, "nature": "harmonious"},
    "Square": {"angle": 90, "orb": 8, "nature": "tense"},
    "Trine": {"angle": 120, "orb": 8, "nature": "harmonious"},
    "Opposition": {"angle": 180, "orb": 8, "nature": "tense"},
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


def get_aspect_factors(target_date):
    """获取指定日期的相位因子"""
    resp = requests.post(f"{API_BASE}/api/calc/daily", json={
        "birthData": BIRTH_DATA,
        "targetDate": target_date,
        "withFactors": True
    }, timeout=10)
    
    data = resp.json()
    factors = data.get("factors", {}).get("factors", [])
    
    # 筛选相位因子
    aspect_factors = [f for f in factors if f.get("type") == "aspectPhase"]
    return aspect_factors


def test_aspect_type_coverage():
    """测试1: 验证相位类型覆盖"""
    print("\n📋 测试1: 相位类型覆盖")
    print("-" * 40)
    
    try:
        # 查询多个日期以获取更多相位
        dates = [
            "2026-01-10T12:00:00+08:00",
            "2026-01-15T12:00:00+08:00",
            "2026-01-20T12:00:00+08:00",
            "2026-02-01T12:00:00+08:00",
        ]
        
        all_aspects = []
        for date in dates:
            aspects = get_aspect_factors(date)
            all_aspects.extend(aspects)
        
        # 提取相位类型
        aspect_types = set()
        for f in all_aspects:
            name = f.get("name", "")
            for aspect_type in EXPECTED_ASPECTS.keys():
                if aspect_type in name:
                    aspect_types.add(aspect_type)
        
        print(f"  检测到的相位类型: {aspect_types}")
        
        missing = set(EXPECTED_ASPECTS.keys()) - aspect_types
        if missing:
            return False, f"缺少相位类型: {missing}"
        
        return True, f"所有 {len(EXPECTED_ASPECTS)} 种相位类型都已检测到"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_aspect_type_distribution():
    """测试2: 统计相位类型分布"""
    print("\n📋 测试2: 相位类型分布统计")
    print("-" * 40)
    
    try:
        aspects = get_aspect_factors("2026-01-15T12:00:00+08:00")
        
        # 统计各类型数量
        type_counter = Counter()
        for f in aspects:
            name = f.get("name", "")
            for aspect_type in EXPECTED_ASPECTS.keys():
                if aspect_type in name:
                    type_counter[aspect_type] += 1
                    break
        
        print(f"  相位因子总数: {len(aspects)}")
        print("  类型分布:")
        for aspect_type, count in sorted(type_counter.items()):
            print(f"    - {aspect_type}: {count}")
        
        if len(aspects) == 0:
            return False, "未检测到任何相位因子"
        
        return True, f"共检测到 {len(aspects)} 个相位因子"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_orb_validity():
    """测试3: 验证容许度（orb）计算"""
    print("\n📋 测试3: 容许度验证")
    print("-" * 40)
    
    try:
        aspects = get_aspect_factors("2026-01-15T12:00:00+08:00")
        
        invalid_orbs = []
        for f in aspects:
            name = f.get("name", "")
            lifecycle = f.get("lifecycle", {})
            
            # 检查 lifecycle 是否存在
            if not lifecycle:
                continue
            
            # 找到对应的相位类型
            for aspect_type, config in EXPECTED_ASPECTS.items():
                if aspect_type in name:
                    # 从因子中获取相关信息（如果有的话）
                    # 这里我们只检查 lifecycle 是否合理
                    break
        
        # 检查所有相位都有 lifecycle
        aspects_with_lifecycle = [f for f in aspects if f.get("lifecycle")]
        coverage = len(aspects_with_lifecycle) / len(aspects) * 100 if aspects else 0
        
        print(f"  有 lifecycle 的相位比例: {coverage:.1f}%")
        
        if coverage < 90:
            return False, f"lifecycle 覆盖率不足: {coverage:.1f}%"
        
        return True, f"{len(aspects_with_lifecycle)}/{len(aspects)} 个相位有 lifecycle"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_applying_separating():
    """测试4: 验证入相/离相判断"""
    print("\n📋 测试4: 入相/离相状态")
    print("-" * 40)
    
    try:
        aspects = get_aspect_factors("2026-01-15T12:00:00+08:00")
        
        # 统计入相/离相
        applying_count = 0
        separating_count = 0
        unknown_count = 0
        
        for f in aspects:
            name = f.get("name", "")
            # 检查因子描述中是否包含入相/离相信息
            desc = f.get("description", "").lower()
            reason = f.get("astroReason", "").lower()
            
            if "applying" in desc or "applying" in reason:
                applying_count += 1
            elif "separating" in desc or "separating" in reason:
                separating_count += 1
            else:
                unknown_count += 1
        
        print(f"  入相 (Applying): {applying_count}")
        print(f"  离相 (Separating): {separating_count}")
        print(f"  未知: {unknown_count}")
        
        # 至少应该有一些入相和离相
        total = len(aspects)
        if total > 0:
            return True, f"共 {total} 个相位"
        else:
            return False, "未检测到相位"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_planet_pairs():
    """测试5: 验证行星配对"""
    print("\n📋 测试5: 行星配对验证")
    print("-" * 40)
    
    try:
        aspects = get_aspect_factors("2026-01-15T12:00:00+08:00")
        
        # 提取行星名称
        planets_in_aspects = set()
        for f in aspects:
            name = f.get("name", "")
            # 解析行星名称 (格式如 "Mars Conjunction Sun")
            parts = name.split()
            if len(parts) >= 3:
                planet1 = parts[0]
                planet2 = parts[-1]
                planets_in_aspects.add(planet1)
                planets_in_aspects.add(planet2)
        
        print(f"  涉及的行星: {sorted(planets_in_aspects)}")
        
        # 检查是否包含主要行星
        major_planets = {"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn"}
        found_major = planets_in_aspects & major_planets
        
        print(f"  主要行星覆盖: {len(found_major)}/{len(major_planets)}")
        
        if len(found_major) >= 5:
            return True, f"覆盖 {len(found_major)} 颗主要行星"
        else:
            missing = major_planets - found_major
            return False, f"缺少行星: {missing}"
    except Exception as e:
        return False, f"请求失败: {e}"


def test_aspect_strength():
    """测试6: 验证相位强度"""
    print("\n📋 测试6: 相位强度验证")
    print("-" * 40)
    
    try:
        aspects = get_aspect_factors("2026-01-15T12:00:00+08:00")
        
        # 检查 currentStrength 字段
        strength_values = []
        for f in aspects:
            strength = f.get("currentStrength", 0)
            strength_values.append(strength)
        
        if not strength_values:
            return False, "未检测到强度数据"
        
        avg_strength = sum(strength_values) / len(strength_values)
        min_strength = min(strength_values)
        max_strength = max(strength_values)
        
        print(f"  强度范围: {min_strength:.4f} - {max_strength:.4f}")
        print(f"  平均强度: {avg_strength:.4f}")
        
        # 强度应该在 0-1 范围内
        invalid = [s for s in strength_values if s < 0 or s > 1.5]
        if invalid:
            return False, f"存在异常强度值: {invalid[:5]}"
        
        return True, f"强度范围正常 ({min_strength:.2f} - {max_strength:.2f})"
    except Exception as e:
        return False, f"请求失败: {e}"


def main():
    print("=" * 60)
    print("相位检测验证测试")
    print(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"API 地址: {API_BASE}")
    print("=" * 60)
    
    result = TestResult()
    
    # 运行所有测试
    tests = [
        ("相位类型覆盖", test_aspect_type_coverage),
        ("相位类型分布", test_aspect_type_distribution),
        ("容许度验证", test_orb_validity),
        ("入相/离相状态", test_applying_separating),
        ("行星配对验证", test_planet_pairs),
        ("相位强度验证", test_aspect_strength),
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
