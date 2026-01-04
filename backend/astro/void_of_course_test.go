package astro

import (
	"testing"
	"time"
)

// TestVoidOfCourse 测试月亮空亡计算
func TestVoidOfCourse(t *testing.T) {
	testDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	jd := DateToJulianDay(testDate)

	t.Log("========== 月亮空亡测试 ==========")
	t.Logf("测试日期: %s", testDate.Format("2006-01-02 15:04"))

	info := CalculateVoidOfCourse(jd, nil)
	
	t.Logf("月亮空亡: %v", info.IsVoid)
	t.Logf("持续时间: %.2f 小时", info.Duration)
	t.Logf("下一个星座: %s", info.NextSign)
	t.Logf("最后相位: %s", info.LastAspect)
	t.Logf("影响因子: %.2f", info.Influence)

	// 验证返回值合理性
	if info.Duration < 0 || info.Duration > 72 {
		t.Errorf("持续时间不合理: %.2f", info.Duration)
	}

	if info.Influence > 0 || info.Influence < -15 {
		t.Errorf("影响因子不合理: %.2f", info.Influence)
	}
}

// TestVoidOfCourseMultipleDates 测试多个日期的月亮空亡
func TestVoidOfCourseMultipleDates(t *testing.T) {
	testDates := []time.Time{
		time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 15, 6, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 21, 18, 0, 0, 0, time.UTC),
		time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 21, 12, 0, 0, 0, time.UTC),
	}

	t.Log("========== 多日期月亮空亡测试 ==========")
	
	for _, date := range testDates {
		jd := DateToJulianDay(date)
		info := CalculateVoidOfCourse(jd, nil)
		
		status := "正常"
		if info.IsVoid {
			status = "空亡"
		}
		
		t.Logf("%s: %s, 持续%.1f小时, 影响%.1f", 
			date.Format("2006-01-02"),
			status,
			info.Duration,
			info.Influence)
	}
}

// TestIsVoidOfCourseMoon 测试简化版月亮空亡检测
func TestIsVoidOfCourseMoon(t *testing.T) {
	testDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	jd := DateToJulianDay(testDate)

	isVoid := IsVoidOfCourseMoon(jd)
	t.Logf("简化版月亮空亡: %v", isVoid)

	// 验证函数不崩溃即可
}

