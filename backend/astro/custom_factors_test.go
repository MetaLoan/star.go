package astro

import (
	"fmt"
	"math"
	"star/models"
	"testing"
	"time"
)

// TestCustomFactorParsing 测试自定义因子解析
func TestCustomFactorParsing(t *testing.T) {
	testCases := []struct {
		input       string
		shouldParse bool
		operation   string
		dimension   string
		value       float64
		duration    float64
	}{
		{
			input:       "AddScore=(10*health,2.5,202601051200)",
			shouldParse: true,
			operation:   "AddScore",
			dimension:   "health",
			value:       10,
			duration:    2.5,
		},
		{
			input:       "SubScore=(5*career,1.0,202601051000)",
			shouldParse: true,
			operation:   "SubScore",
			dimension:   "career",
			value:       5,
			duration:    1.0,
		},
		{
			input:       "MulScore=(1.5*finance,3,202601050800)",
			shouldParse: true,
			operation:   "MulScore",
			dimension:   "finance",
			value:       1.5,
			duration:    3,
		},
		{
			input:       "InvalidOp=(10,1,202601050000)",
			shouldParse: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			factor, err := ParseCustomFactor(tc.input)

			if tc.shouldParse {
				if err != nil {
					t.Errorf("Expected parsing to succeed, but got error: %v", err)
					return
				}
				if factor.Operation != tc.operation {
					t.Errorf("Expected operation %s, got %s", tc.operation, factor.Operation)
				}
				if factor.Dimension != tc.dimension {
					t.Errorf("Expected dimension %s, got %s", tc.dimension, factor.Dimension)
				}
				if factor.Value != tc.value {
					t.Errorf("Expected value %f, got %f", tc.value, factor.Value)
				}
				if factor.Duration != tc.duration {
					t.Errorf("Expected duration %f, got %f", tc.duration, factor.Duration)
				}
			} else {
				if err == nil {
					t.Errorf("Expected parsing to fail, but it succeeded")
				}
			}
		})
	}
}

// TestCustomFactorInfluence 测试自定义因子对分数的实际影响
func TestCustomFactorInfluence(t *testing.T) {
	// 创建测试用的本命盘数据
	testBirthData := models.BirthData{
		Year:      1990,
		Month:     6,
		Day:       15,
		Hour:      14,
		Minute:    30,
		Latitude:  39.9042,  // 北京
		Longitude: 116.4074,
		Timezone:  8,
	}

	// 初始化 Swiss Ephemeris（需要先运行）
	InitSwissEphemeris("")

	// 计算本命盘
	chart := CalculateNatalChart(testBirthData)

	// 清除之前的因子
	testUserID := "test_user_001"
	ClearCustomFactors(testUserID)

	// 设置测试时间：2026年1月5日 12:00
	testTime := time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC)

	// 1. 计算没有自定义因子时的分数（使用空用户ID）
	scoreWithoutFactor := CalculateUnifiedHourlyScoreWithUser(chart, testTime, "")

	fmt.Println("\n========== 自定义因子影响测试 ==========")
	fmt.Println("\n【无因子时的分数】")
	fmt.Printf("  综合分: %.2f\n", scoreWithoutFactor.Overall)
	fmt.Printf("  事业: %.2f\n", scoreWithoutFactor.Dimensions["career"])
	fmt.Printf("  关系: %.2f\n", scoreWithoutFactor.Dimensions["relationship"])
	fmt.Printf("  健康: %.2f\n", scoreWithoutFactor.Dimensions["health"])
	fmt.Printf("  财务: %.2f\n", scoreWithoutFactor.Dimensions["finance"])
	fmt.Printf("  灵性: %.2f\n", scoreWithoutFactor.Dimensions["spiritual"])

	// 2. 添加健康维度 +20 的因子（持续3小时，从11:00开始）
	factorDef := "AddScore=(20*health,3.0,202601051100)"
	factor, err := AddCustomFactor(testUserID, factorDef)
	if err != nil {
		t.Fatalf("Failed to add custom factor: %v", err)
	}

	fmt.Println("\n【添加的自定义因子】")
	fmt.Printf("  定义: %s\n", factorDef)
	fmt.Printf("  操作: %s\n", factor.Operation)
	fmt.Printf("  维度: %s\n", factor.Dimension)
	fmt.Printf("  值: %.2f\n", factor.Value)
	fmt.Printf("  开始时间: %s\n", factor.StartTime.Format("2006-01-05 15:04:05"))
	fmt.Printf("  结束时间: %s\n", factor.EndTime.Format("2006-01-05 15:04:05"))
	fmt.Printf("  持续时长: %.1f 小时\n", factor.Duration)

	// 3. 计算有自定义因子时的分数
	scoreWithFactor := CalculateUnifiedHourlyScoreWithUser(chart, testTime, testUserID)

	fmt.Println("\n【有因子时的分数】")
	fmt.Printf("  综合分: %.2f (变化: %+.2f)\n", scoreWithFactor.Overall, scoreWithFactor.Overall-scoreWithoutFactor.Overall)
	fmt.Printf("  事业: %.2f (变化: %+.2f)\n", scoreWithFactor.Dimensions["career"], scoreWithFactor.Dimensions["career"]-scoreWithoutFactor.Dimensions["career"])
	fmt.Printf("  关系: %.2f (变化: %+.2f)\n", scoreWithFactor.Dimensions["relationship"], scoreWithFactor.Dimensions["relationship"]-scoreWithoutFactor.Dimensions["relationship"])
	fmt.Printf("  健康: %.2f (变化: %+.2f)\n", scoreWithFactor.Dimensions["health"], scoreWithFactor.Dimensions["health"]-scoreWithoutFactor.Dimensions["health"])
	fmt.Printf("  财务: %.2f (变化: %+.2f)\n", scoreWithFactor.Dimensions["finance"], scoreWithFactor.Dimensions["finance"]-scoreWithoutFactor.Dimensions["finance"])
	fmt.Printf("  灵性: %.2f (变化: %+.2f)\n", scoreWithFactor.Dimensions["spiritual"], scoreWithFactor.Dimensions["spiritual"]-scoreWithoutFactor.Dimensions["spiritual"])

	// 验证健康分数确实发生了变化
	healthDiff := scoreWithFactor.Dimensions["health"] - scoreWithoutFactor.Dimensions["health"]
	if math.Abs(healthDiff) < 0.01 {
		t.Errorf("健康维度分数没有变化！预期有变化，实际差值: %.4f", healthDiff)
	} else {
		fmt.Printf("\n✅ 健康维度分数变化: %+.2f (符合预期)\n", healthDiff)
	}

	// 4. 测试因子在不同时间点的强度（正弦曲线）
	fmt.Println("\n【因子强度在时间轴上的变化（正弦曲线）】")
	fmt.Println("  时间           | 强度   | 健康分变化")
	fmt.Println("  -------------- | ------ | ----------")

	testTimes := []time.Time{
		time.Date(2026, 1, 5, 11, 0, 0, 0, time.UTC),  // 开始
		time.Date(2026, 1, 5, 11, 30, 0, 0, time.UTC), // 1/4
		time.Date(2026, 1, 5, 12, 30, 0, 0, time.UTC), // 峰值
		time.Date(2026, 1, 5, 13, 30, 0, 0, time.UTC), // 3/4
		time.Date(2026, 1, 5, 14, 0, 0, 0, time.UTC),  // 结束
		time.Date(2026, 1, 5, 15, 0, 0, 0, time.UTC),  // 结束后
	}

	for _, tt := range testTimes {
		// 计算当前强度
		lifecycle := &models.FactorLifecycle{
			StartTime: factor.StartTime,
			PeakTime:  factor.StartTime.Add(time.Duration(factor.Duration/2) * time.Hour),
			EndTime:   factor.EndTime,
			Duration:  factor.Duration,
		}
		strength := CalculateFactorStrength(lifecycle, tt)

		// 计算有因子时的分数
		scoreAt := CalculateUnifiedHourlyScoreWithUser(chart, tt, testUserID)
		baseScoreAt := CalculateUnifiedHourlyScoreWithUser(chart, tt, "")
		healthChange := scoreAt.Dimensions["health"] - baseScoreAt.Dimensions["health"]

		fmt.Printf("  %s | %.4f | %+.2f\n", tt.Format("15:04"), strength, healthChange)
	}

	// 5. 测试因子结束后是否停止影响
	fmt.Println("\n【因子结束后的验证】")
	afterEndTime := time.Date(2026, 1, 5, 15, 0, 0, 0, time.UTC)
	scoreAfterEnd := CalculateUnifiedHourlyScoreWithUser(chart, afterEndTime, testUserID)
	baseAfterEnd := CalculateUnifiedHourlyScoreWithUser(chart, afterEndTime, "")
	afterEndDiff := math.Abs(scoreAfterEnd.Dimensions["health"] - baseAfterEnd.Dimensions["health"])

	if afterEndDiff < 0.01 {
		fmt.Printf("  ✅ 因子结束后，健康分数差值: %.4f (接近0，符合预期)\n", afterEndDiff)
	} else {
		t.Errorf("因子结束后分数仍有差异: %.4f", afterEndDiff)
	}

	// 清理
	ClearCustomFactors(testUserID)

	fmt.Println("\n========== 测试完成 ==========")
}

// TestMultipleFactorsCombined 测试多个因子组合影响
func TestMultipleFactorsCombined(t *testing.T) {
	// 初始化
	testUserID := "test_user_002"
	ClearCustomFactors(testUserID)
	InitSwissEphemeris("")

	testBirthData := models.BirthData{
		Year: 1990, Month: 6, Day: 15, Hour: 14, Minute: 30,
		Latitude: 39.9042, Longitude: 116.4074, Timezone: 8,
	}
	chart := CalculateNatalChart(testBirthData)
	testTime := time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC)

	// 基准分数
	baseScore := CalculateUnifiedHourlyScoreWithUser(chart, testTime, "")

	// 添加多个因子
	AddCustomFactor(testUserID, "AddScore=(10*career,2,202601051100)")
	AddCustomFactor(testUserID, "AddScore=(15*health,2,202601051100)")
	AddCustomFactor(testUserID, "SubScore=(5*finance,2,202601051100)")

	combinedScore := CalculateUnifiedHourlyScoreWithUser(chart, testTime, testUserID)

	fmt.Println("\n========== 多因子组合测试 ==========")
	fmt.Println("\n【添加的因子】")
	fmt.Println("  1. 事业 +10（2小时）")
	fmt.Println("  2. 健康 +15（2小时）")
	fmt.Println("  3. 财务 -5（2小时）")

	fmt.Println("\n【分数变化】")
	fmt.Printf("  事业: %.2f → %.2f (变化: %+.2f)\n",
		baseScore.Dimensions["career"],
		combinedScore.Dimensions["career"],
		combinedScore.Dimensions["career"]-baseScore.Dimensions["career"])
	fmt.Printf("  健康: %.2f → %.2f (变化: %+.2f)\n",
		baseScore.Dimensions["health"],
		combinedScore.Dimensions["health"],
		combinedScore.Dimensions["health"]-baseScore.Dimensions["health"])
	fmt.Printf("  财务: %.2f → %.2f (变化: %+.2f)\n",
		baseScore.Dimensions["finance"],
		combinedScore.Dimensions["finance"],
		combinedScore.Dimensions["finance"]-baseScore.Dimensions["finance"])

	// 验证所有因子都有影响
	careerDiff := combinedScore.Dimensions["career"] - baseScore.Dimensions["career"]
	healthDiff := combinedScore.Dimensions["health"] - baseScore.Dimensions["health"]
	financeDiff := combinedScore.Dimensions["finance"] - baseScore.Dimensions["finance"]

	if careerDiff <= 0 {
		t.Errorf("事业因子没有正向影响: %+.2f", careerDiff)
	}
	if healthDiff <= 0 {
		t.Errorf("健康因子没有正向影响: %+.2f", healthDiff)
	}
	if financeDiff >= 0 {
		t.Errorf("财务因子没有负向影响: %+.2f", financeDiff)
	}

	// 清理
	ClearCustomFactors(testUserID)
	fmt.Println("\n========== 测试完成 ==========")
}

