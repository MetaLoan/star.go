#!/usr/bin/env python3
"""
V2 Astro API 多粒度与事件类型接口测试脚本

测试目标：
1. 覆盖 5 种粒度（hour/day/week/month/year）
2. 校验 SubSlots 数量、分数范围、事件过滤逻辑
3. 按事件类型校验必填字段与取值范围
4. 发现不合理点并整理成列表

运行方式：
    先启动服务: cd backend && go run main.go
    再运行测试: python3 backend/tests/test_v2_astro_api.py
"""

import requests
import json
import sys
from datetime import datetime
from collections import Counter

# 配置
API_BASE = "http://localhost:8080"
V2_ENDPOINT = f"{API_BASE}/api/v2/astro"

# 固定出生数据
BIRTH_DATA = {
    "year": 1990,
    "month": 5,
    "day": 15,
    "hour": 10,
    "minute": 30,
    "second": 0,
    "latitude": 39.9042,
    "longitude": 116.4074,
    "timezone": 8
}

# 固定查询时间 (2026年1月16日)
QUERY_TIME = "2026-01-16T14:00:00+08:00"

# 粒度配置
GRANULARITIES = ["hour", "day", "week", "month", "year"]

# 有效事件类型
VALID_EVENT_TYPES = {
    "aspect", "transit_house", "progression", "planetary_hour",
    "void_of_course", "retrograde", "lunar_phase", "sign_change", "dignity"
}

# 有效相位阶段
VALID_PHASES = {"approaching", "exact", "separating", "active", ""}

# 有效时间层级
VALID_TIME_LEVELS = {"hourly", "daily", "weekly", "monthly", "yearly", ""}

# 五维字段
DIMENSIONS = ["career", "relationship", "health", "finance", "spiritual"]

# 尊贵度类型
VALID_DIGNITY_TYPES = {"domicile", "exaltation", "detriment", "fall"}


class TestResult:
    """测试结果收集器"""
    def __init__(self):
        self.passed = 0
        self.failed = 0
        self.issues = []  # 不合理点列表
        self.results = []
    
    def add(self, name: str, passed: bool, message: str = ""):
        self.results.append({
            "name": name,
            "passed": passed,
            "message": message
        })
        if passed:
            self.passed += 1
        else:
            self.failed += 1
    
    def add_issue(self, granularity: str, event_type: str, description: str):
        """记录不合理点"""
        self.issues.append({
            "granularity": granularity,
            "event_type": event_type,
            "description": description
        })
    
    def summary(self):
        total = self.passed + self.failed
        print("\n" + "=" * 70)
        print(f"测试结果: {self.passed}/{total} 通过")
        print("=" * 70)
        
        for r in self.results:
            status = "PASS" if r["passed"] else "FAIL"
            print(f"  [{status}] {r['name']}")
            if r["message"] and not r["passed"]:
                print(f"          {r['message']}")
        
        if self.issues:
            print("\n" + "-" * 70)
            print(f"发现 {len(self.issues)} 个不合理点:")
            print("-" * 70)
            for i, issue in enumerate(self.issues, 1):
                print(f"  {i}. [{issue['granularity']}][{issue['event_type']}] {issue['description']}")
        
        return self.failed == 0


def call_v2_api(granularity: str, language: str = "en", query_time: str = QUERY_TIME):
    """调用 V2 API"""
    payload = {
        "birth": BIRTH_DATA,
        "time": query_time,
        "granularity": granularity,
        "language": language
    }
    
    try:
        resp = requests.post(V2_ENDPOINT, json=payload, timeout=60)
        resp.raise_for_status()
        return resp.json()
    except requests.RequestException as e:
        return {"error": str(e)}


def check_score_range(score: float, name: str, result: TestResult, granularity: str):
    """检查分数是否在 [0, 100] 范围内"""
    if score is None:
        result.add_issue(granularity, "scores", f"{name} 为 None")
        return False
    if not (0 <= score <= 100):
        result.add_issue(granularity, "scores", f"{name}={score} 超出 [0,100] 范围")
        return False
    return True


def check_scores(scores: dict, result: TestResult, granularity: str) -> bool:
    """检查分数结构和范围"""
    if not scores:
        result.add_issue(granularity, "scores", "scores 为空")
        return False
    
    all_valid = True
    
    # 检查 overall
    if "overall" not in scores:
        result.add_issue(granularity, "scores", "缺少 overall 字段")
        all_valid = False
    else:
        if not check_score_range(scores["overall"], "overall", result, granularity):
            all_valid = False
    
    # 检查五维
    for dim in DIMENSIONS:
        if dim not in scores:
            result.add_issue(granularity, "scores", f"缺少 {dim} 字段")
            all_valid = False
        else:
            if not check_score_range(scores[dim], dim, result, granularity):
                all_valid = False
    
    return all_valid


def check_subslots(subslots: list, granularity: str, result: TestResult, query_time: str) -> bool:
    """检查 SubSlots 数量和结构"""
    
    # 预期数量
    expected = {
        "hour": 0,
        "day": 24,
        "week": 7,
        "month": None,  # 依赖当月天数
        "year": 12
    }
    
    actual_count = len(subslots) if subslots else 0
    
    if granularity == "hour":
        if actual_count != 0:
            result.add_issue(granularity, "subSlots", f"hour 粒度应无 SubSlots，实际有 {actual_count}")
            return False
        return True
    
    if granularity == "month":
        # 解析查询时间确定当月天数
        try:
            dt = datetime.fromisoformat(query_time.replace("Z", "+00:00"))
            import calendar
            days_in_month = calendar.monthrange(dt.year, dt.month)[1]
            if actual_count != days_in_month:
                result.add_issue(granularity, "subSlots", 
                    f"month 应有 {days_in_month} 个 SubSlots（当月天数），实际 {actual_count}")
                return False
        except:
            if actual_count < 28 or actual_count > 31:
                result.add_issue(granularity, "subSlots", f"month SubSlots 数量异常: {actual_count}")
                return False
        return True
    
    exp = expected.get(granularity)
    if exp is not None and actual_count != exp:
        result.add_issue(granularity, "subSlots", f"应有 {exp} 个 SubSlots，实际 {actual_count}")
        return False
    
    # 检查 SubSlots 时间单调性
    if subslots:
        prev_time = None
        for i, sub in enumerate(subslots):
            st = sub.get("startTime")
            if st and prev_time:
                try:
                    curr = datetime.fromisoformat(st.replace("Z", "+00:00"))
                    if curr <= prev_time:
                        result.add_issue(granularity, "subSlots", 
                            f"SubSlots[{i}].startTime 非递增: {prev_time} -> {curr}")
                        return False
                    prev_time = curr
                except:
                    pass
            elif st:
                try:
                    prev_time = datetime.fromisoformat(st.replace("Z", "+00:00"))
                except:
                    pass
            
            # 检查子槽分数
            sub_scores = sub.get("scores", {})
            for dim in ["overall"] + DIMENSIONS:
                val = sub_scores.get(dim)
                if val is not None and (val < 0 or val > 100):
                    result.add_issue(granularity, "subSlots", 
                        f"SubSlots[{i}].scores.{dim}={val} 超出范围")
    
    return True


def check_events(events: list, granularity: str, result: TestResult) -> bool:
    """检查事件列表"""
    if not events:
        # 空事件列表是允许的
        return True
    
    # 粒度对应的允许 timeLevel
    allowed_levels = {
        "hour": {"hourly", "daily", "weekly", "monthly", "yearly"},
        "day": {"daily", "weekly", "monthly", "yearly"},
        "week": {"weekly", "monthly", "yearly"},
        "month": {"monthly", "yearly"},
        "year": {"yearly"}
    }
    
    allowed = allowed_levels.get(granularity, set())
    all_valid = True
    event_ids = []
    
    for i, event in enumerate(events):
        event_id = event.get("eventId", "")
        event_type = event.get("type", "")
        time_level = event.get("timeLevel", "")
        
        # 检查 eventId 非空
        if not event_id:
            result.add_issue(granularity, event_type or "unknown", f"events[{i}] 缺少 eventId")
            all_valid = False
        event_ids.append(event_id)
        
        # 检查事件类型
        if event_type and event_type not in VALID_EVENT_TYPES:
            result.add_issue(granularity, event_type, f"events[{i}] 未知事件类型: {event_type}")
            all_valid = False
        
        # 检查 timeLevel 与粒度匹配
        if time_level and time_level not in allowed:
            result.add_issue(granularity, event_type, 
                f"events[{i}] timeLevel={time_level} 不应出现在 {granularity} 粒度（应为 {allowed}）")
            all_valid = False
        
        # 检查 intensity 范围
        intensity = event.get("intensity")
        if intensity is not None and (intensity < 0 or intensity > 1):
            result.add_issue(granularity, event_type, f"events[{i}] intensity={intensity} 超出 [0,1]")
            all_valid = False
        
        # 检查 isPositive 存在
        if "isPositive" not in event:
            result.add_issue(granularity, event_type, f"events[{i}] 缺少 isPositive")
            all_valid = False
        
        # 检查 phase
        phase = event.get("phase", "")
        if phase and phase not in VALID_PHASES:
            result.add_issue(granularity, event_type, f"events[{i}] 未知 phase: {phase}")
            all_valid = False
        
        # 检查 impact 五维完整
        impact = event.get("impact", {})
        for dim in DIMENSIONS:
            if dim not in impact:
                result.add_issue(granularity, event_type, f"events[{i}] impact 缺少 {dim}")
                all_valid = False
        
        # 检查 impactDelta 五维完整
        impact_delta = event.get("impactDelta", {})
        for dim in DIMENSIONS:
            if dim not in impact_delta:
                result.add_issue(granularity, event_type, f"events[{i}] impactDelta 缺少 {dim}")
                all_valid = False
        
        # 按事件类型做额外校验
        all_valid = check_event_by_type(event, i, granularity, result) and all_valid
    
    # 检查 eventId 重复
    id_counts = Counter(event_ids)
    duplicates = [eid for eid, cnt in id_counts.items() if cnt > 1 and eid]
    if duplicates:
        result.add_issue(granularity, "events", f"eventId 重复: {duplicates[:3]}")
        all_valid = False
    
    return all_valid


def check_event_by_type(event: dict, idx: int, granularity: str, result: TestResult) -> bool:
    """按事件类型做额外校验"""
    event_type = event.get("type", "")
    all_valid = True
    
    if event_type == "aspect":
        # 检查必填字段
        for field in ["primaryPlanet", "secondaryPlanet", "aspect"]:
            if not event.get(field):
                result.add_issue(granularity, event_type, f"events[{idx}] aspect 事件缺少 {field}")
                all_valid = False
    
    elif event_type == "planetary_hour":
        if not event.get("primaryPlanet"):
            result.add_issue(granularity, event_type, f"events[{idx}] planetary_hour 缺少 primaryPlanet")
            all_valid = False
        # 小时级事件在非 hour 粒度不应出现（已在 timeLevel 检查中覆盖）
    
    elif event_type == "void_of_course":
        pp = event.get("primaryPlanet", "")
        if pp and pp.lower() != "moon":
            result.add_issue(granularity, event_type, f"events[{idx}] void_of_course primaryPlanet 应为 moon，实际 {pp}")
            all_valid = False
    
    elif event_type == "retrograde":
        if not event.get("primaryPlanet"):
            result.add_issue(granularity, event_type, f"events[{idx}] retrograde 缺少 primaryPlanet")
            all_valid = False
    
    elif event_type == "lunar_phase":
        aspect = event.get("aspect", "")
        if aspect and aspect.lower() not in {"new", "full", "firstquarter", "lastquarter", "first_quarter", "last_quarter", "waxing", "waning"}:
            # 宽松检查，仅警告
            pass
    
    elif event_type == "sign_change":
        aspect = event.get("aspect", "")
        if not aspect:
            result.add_issue(granularity, event_type, f"events[{idx}] sign_change 缺少 aspect（新星座）")
            all_valid = False
    
    elif event_type == "dignity":
        aspect = event.get("aspect", "")
        if aspect and aspect.lower() not in VALID_DIGNITY_TYPES:
            result.add_issue(granularity, event_type, f"events[{idx}] dignity aspect={aspect} 不在有效类型中")
            all_valid = False
    
    return all_valid


def check_delta(delta: dict, granularity: str, result: TestResult) -> bool:
    """检查 delta 结构"""
    if granularity == "hour":
        # hour 粒度 delta 可选
        return True
    
    if not delta:
        result.add_issue(granularity, "delta", f"{granularity} 粒度应有 delta，实际为空")
        return False
    
    # 检查 dimensions
    dims = delta.get("dimensions", {})
    for dim in ["overall"] + DIMENSIONS:
        if dim not in dims:
            result.add_issue(granularity, "delta", f"delta.dimensions 缺少 {dim}")
            return False
    
    return True


def check_guidance(guidance: dict, granularity: str, result: TestResult) -> bool:
    """检查 guidance 结构"""
    if not guidance:
        # guidance 可选，但最好有
        return True
    
    # 检查必要字段存在
    for field in ["summary", "dos", "donts", "focus"]:
        if field not in guidance:
            result.add_issue(granularity, "guidance", f"guidance 缺少 {field}")
            return False
    
    # focus 应为五维之一或空
    focus = guidance.get("focus", "")
    if focus and focus not in DIMENSIONS:
        result.add_issue(granularity, "guidance", f"guidance.focus={focus} 不在五维中")
        return False
    
    return True


def check_meta(meta: dict, events_count: int, granularity: str, result: TestResult) -> bool:
    """检查 meta 结构"""
    if not meta:
        result.add_issue(granularity, "meta", "meta 为空")
        return False
    
    # eventCount 与实际事件数一致
    reported = meta.get("eventCount", -1)
    if reported != events_count:
        result.add_issue(granularity, "meta", f"meta.eventCount={reported} 与实际 {events_count} 不一致")
        return False
    
    # cached 为布尔值
    if "cached" not in meta:
        result.add_issue(granularity, "meta", "meta 缺少 cached")
        return False
    
    return True


def test_granularity(granularity: str, result: TestResult, language: str = "en"):
    """测试单个粒度"""
    print(f"\n>>> 测试粒度: {granularity} (language={language})")
    
    data = call_v2_api(granularity, language)
    
    if "error" in data:
        result.add(f"{granularity}: API 调用", False, data["error"])
        return
    
    slot = data.get("slot", {})
    meta = data.get("meta", {})
    
    # 1. 检查基本结构
    if not slot:
        result.add(f"{granularity}: slot 存在", False, "slot 为空")
        return
    result.add(f"{granularity}: slot 存在", True)
    
    # 2. 检查 granularity 字段
    returned_gran = slot.get("granularity", "")
    if returned_gran != granularity:
        result.add(f"{granularity}: granularity 匹配", False, f"返回 {returned_gran}")
    else:
        result.add(f"{granularity}: granularity 匹配", True)
    
    # 3. 检查分数
    scores = slot.get("scores", {})
    scores_valid = check_scores(scores, result, granularity)
    result.add(f"{granularity}: scores 有效", scores_valid)
    
    # 4. 检查 SubSlots
    subslots = slot.get("subSlots", [])
    subslots_valid = check_subslots(subslots, granularity, result, QUERY_TIME)
    result.add(f"{granularity}: subSlots 有效", subslots_valid)
    
    # 5. 检查事件
    events = slot.get("events", [])
    events_valid = check_events(events, granularity, result)
    result.add(f"{granularity}: events 有效", events_valid, f"共 {len(events)} 个事件")
    
    # 6. 检查 delta
    delta = slot.get("delta", {})
    delta_valid = check_delta(delta, granularity, result)
    result.add(f"{granularity}: delta 有效", delta_valid)
    
    # 7. 检查 guidance
    guidance = slot.get("guidance", {})
    guidance_valid = check_guidance(guidance, granularity, result)
    result.add(f"{granularity}: guidance 有效", guidance_valid)
    
    # 8. 检查 meta
    meta_valid = check_meta(meta, len(events), granularity, result)
    result.add(f"{granularity}: meta 有效", meta_valid)
    
    # 9. 时间范围检查
    start_time = slot.get("startTime", "")
    end_time = slot.get("endTime", "")
    if start_time and end_time:
        try:
            st = datetime.fromisoformat(start_time.replace("Z", "+00:00"))
            et = datetime.fromisoformat(end_time.replace("Z", "+00:00"))
            if st >= et:
                result.add_issue(granularity, "time", f"startTime >= endTime: {start_time} >= {end_time}")
                result.add(f"{granularity}: 时间范围", False)
            else:
                result.add(f"{granularity}: 时间范围", True)
        except Exception as e:
            result.add(f"{granularity}: 时间范围", False, str(e))
    else:
        result.add(f"{granularity}: 时间范围", False, "缺少 startTime 或 endTime")
    
    # 统计事件类型分布
    type_counts = Counter(e.get("type", "unknown") for e in events)
    print(f"    事件类型分布: {dict(type_counts)}")
    
    # 统计 timeLevel 分布
    level_counts = Counter(e.get("timeLevel", "unknown") for e in events)
    print(f"    timeLevel 分布: {dict(level_counts)}")


def test_language_variation(result: TestResult):
    """测试不同语言的响应"""
    print("\n>>> 测试语言变化 (day 粒度)")
    
    languages = ["zh", "en", "ru"]
    titles_by_lang = {}
    
    for lang in languages:
        data = call_v2_api("day", lang)
        if "error" in data:
            result.add(f"语言 {lang}: API 调用", False, data["error"])
            continue
        
        events = data.get("slot", {}).get("events", [])
        if events:
            first_title = events[0].get("title", "")
            titles_by_lang[lang] = first_title
            print(f"    {lang}: 首个事件标题 = {first_title[:50]}...")
    
    # 检查至少有两种语言返回了不同内容（如果有事件的话）
    unique_titles = set(titles_by_lang.values())
    if len(titles_by_lang) >= 2:
        if len(unique_titles) >= 2:
            result.add("语言: 多语言标题差异", True)
        else:
            result.add("语言: 多语言标题差异", False, "所有语言返回相同标题")
    else:
        result.add("语言: 多语言标题差异", True, "事件不足以比较")


def main():
    print("=" * 70)
    print("V2 Astro API 多粒度与事件类型测试")
    print("=" * 70)
    print(f"API: {V2_ENDPOINT}")
    print(f"查询时间: {QUERY_TIME}")
    print(f"出生数据: {BIRTH_DATA['year']}-{BIRTH_DATA['month']}-{BIRTH_DATA['day']}")
    
    # 健康检查
    try:
        health = requests.get(f"{API_BASE}/health", timeout=5)
        if health.status_code != 200:
            print("\n[ERROR] 服务未启动或不可用")
            sys.exit(1)
        print("\n服务状态: OK")
    except:
        print("\n[ERROR] 无法连接服务，请先启动: cd backend && go run main.go")
        sys.exit(1)
    
    result = TestResult()
    
    # 测试所有粒度
    for gran in GRANULARITIES:
        test_granularity(gran, result)
    
    # 测试语言变化
    test_language_variation(result)
    
    # 输出总结
    success = result.summary()
    
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
