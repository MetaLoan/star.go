package astro

import (
	"math"
	"star/models"
	"testing"
	"time"
)

// 精度验证测试
// 使用已知天文事件验证计算精度

// TestSunPositionAccuracy 测试太阳位置精度
func TestSunPositionAccuracy(t *testing.T) {
	testCases := []struct {
		name         string
		date         time.Time
		expectedLong float64 // 期望黄经
		tolerance    float64 // 容许误差（度）
	}{
		{
			name:         "2024年春分",
			date:         time.Date(2024, 3, 20, 3, 6, 0, 0, time.UTC),
			expectedLong: 0.0,
			tolerance:    0.5,
		},
		{
			name:         "2024年夏至",
			date:         time.Date(2024, 6, 20, 20, 51, 0, 0, time.UTC),
			expectedLong: 90.0,
			tolerance:    0.5,
		},
		{
			name:         "2024年秋分",
			date:         time.Date(2024, 9, 22, 12, 44, 0, 0, time.UTC),
			expectedLong: 180.0,
			tolerance:    0.5,
		},
		{
			name:         "2024年冬至",
			date:         time.Date(2024, 12, 21, 9, 21, 0, 0, time.UTC),
			expectedLong: 270.0,
			tolerance:    0.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jd := DateToJulianDay(tc.date)
			sunPos := CalculateSunPosition(jd)

			diff := math.Abs(sunPos.Longitude - tc.expectedLong)
			if diff > 180 {
				diff = 360 - diff
			}

			t.Logf("%s: 计算黄经=%.4f°, 期望=%.1f°, 误差=%.4f°",
				tc.name, sunPos.Longitude, tc.expectedLong, diff)

			if diff > tc.tolerance {
				t.Errorf("太阳位置误差过大: %.4f° > %.1f°", diff, tc.tolerance)
			}
		})
	}
}

// TestMoonPositionAccuracy 测试月亮位置精度（新月/满月）
func TestMoonPositionAccuracy(t *testing.T) {
	testCases := []struct {
		name      string
		date      time.Time
		isNewMoon bool
		tolerance float64
	}{
		{
			name:      "2024年1月新月",
			date:      time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC),
			isNewMoon: true,
			tolerance: 5.0,
		},
		{
			name:      "2024年1月满月",
			date:      time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC),
			isNewMoon: false,
			tolerance: 5.0,
		},
		{
			name:      "2024年6月新月",
			date:      time.Date(2024, 6, 6, 12, 38, 0, 0, time.UTC),
			isNewMoon: true,
			tolerance: 5.0,
		},
		{
			name:      "2024年6月满月",
			date:      time.Date(2024, 6, 22, 1, 8, 0, 0, time.UTC),
			isNewMoon: false,
			tolerance: 5.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jd := DateToJulianDay(tc.date)
			sunPos := CalculateSunPosition(jd)
			moonPos := CalculateMoonPosition(jd)

			diff := math.Abs(moonPos.Longitude - sunPos.Longitude)
			if diff > 180 {
				diff = 360 - diff
			}

			var expectedDiff float64
			if tc.isNewMoon {
				expectedDiff = 0.0
			} else {
				expectedDiff = 180.0
			}

			actualError := math.Abs(diff - expectedDiff)

			t.Logf("%s: 日月角距=%.4f°, 期望=%.1f°, 误差=%.4f°",
				tc.name, diff, expectedDiff, actualError)

			if actualError > tc.tolerance {
				t.Errorf("月亮位置误差过大: %.4f° > %.1f°", actualError, tc.tolerance)
			}
		})
	}
}

// TestPlanetPositionAccuracy 测试行星位置精度
// 注：简化算法精度有限，主要验证计算流程正确性
// 如需高精度，建议集成 VSOP87 或 Swiss Ephemeris
func TestPlanetPositionAccuracy(t *testing.T) {
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	jd := DateToJulianDay(testDate)

	// 主要验证行星位置计算不崩溃，并输出结果供人工核查
	planets := []models.PlanetID{
		models.Mercury, models.Venus, models.Mars,
		models.Jupiter, models.Saturn, models.Uranus,
		models.Neptune, models.Pluto,
	}

	t.Log("========== 行星位置计算结果（简化算法） ==========")
	t.Log("日期: 2024-01-01 00:00 UTC")
	t.Log("注：简化开普勒算法，精度约 ±30°，后续可集成精确星历表")
	t.Log("")

	for _, planet := range planets {
		pos := CalculatePlanetPosition(planet, jd)

		// 验证位置在有效范围
		if pos.Longitude < 0 || pos.Longitude >= 360 {
			t.Errorf("%s 黄经超出范围: %.2f°", planet, pos.Longitude)
		}

		// 验证星座分配正确
		expectedSign := int(pos.Longitude / 30)
		actualSign := getZodiacIndex(pos.Sign)
		if expectedSign != actualSign {
			t.Errorf("%s 星座分配错误: 期望=%d, 实际=%d", planet, expectedSign, actualSign)
		}

		t.Logf("%-10s: %7.2f° (%s %5.1f°) %s",
			pos.Name, pos.Longitude, pos.SignName, pos.SignDegree,
			func() string {
				if pos.Retrograde {
					return "℞"
				}
				return ""
			}())
	}
}

// getZodiacIndex 获取星座索引
func getZodiacIndex(sign models.ZodiacID) int {
	zodiacOrder := []models.ZodiacID{
		models.Aries, models.Taurus, models.Gemini, models.Cancer,
		models.Leo, models.Virgo, models.Libra, models.Scorpio,
		models.Sagittarius, models.Capricorn, models.Aquarius, models.Pisces,
	}
	for i, z := range zodiacOrder {
		if z == sign {
			return i
		}
	}
	return -1
}

// TestJulianDayCalculation 测试儒略日计算
func TestJulianDayCalculation(t *testing.T) {
	testCases := []struct {
		date       time.Time
		expectedJD float64
	}{
		{
			date:       time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedJD: 2451545.0,
		},
		{
			date:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedJD: 2460310.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.date.Format("2006-01-02"), func(t *testing.T) {
			jd := DateToJulianDay(tc.date)
			diff := math.Abs(jd - tc.expectedJD)

			t.Logf("日期=%s, 计算JD=%.6f, 期望JD=%.1f, 差=%.6f",
				tc.date.Format("2006-01-02 15:04"), jd, tc.expectedJD, diff)

			if diff > 0.001 {
				t.Errorf("儒略日误差过大: %.6f", diff)
			}
		})
	}
}

// TestHouseCuspCalculation 测试宫位计算
func TestHouseCuspCalculation(t *testing.T) {
	testDate := time.Date(1990, 6, 15, 6, 30, 0, 0, time.UTC)
	lat := 39.9042
	lon := 116.4074

	jd := DateToJulianDay(testDate)
	houses, asc, mc := CalculateHouses(jd, lat, lon)

	t.Logf("测试日期: %s", testDate.Format("2006-01-02 15:04 MST"))
	t.Logf("位置: 纬度=%.4f, 经度=%.4f", lat, lon)
	t.Logf("上升点: %.2f°, 中天: %.2f°", asc, mc)

	for _, h := range houses {
		sign := getZodiacFromLongitude(h.Cusp)
		degree := math.Mod(h.Cusp, 30)
		t.Logf("第%d宫: %.2f° (%s %.1f°)", h.House, h.Cusp, sign, degree)
	}

	if asc < 0 || asc >= 360 {
		t.Errorf("上升点超出范围: %.2f", asc)
	}
}

func getZodiacFromLongitude(longitude float64) string {
	signs := []string{
		"白羊", "金牛", "双子", "巨蟹", "狮子", "处女",
		"天秤", "天蝎", "射手", "摩羯", "水瓶", "双鱼",
	}
	index := int(longitude/30) % 12
	if index < 0 {
		index += 12
	}
	return signs[index]
}

// TestAspectCalculation 测试相位计算
func TestAspectCalculation(t *testing.T) {
	planets := []models.PlanetPosition{
		{ID: models.Sun, Longitude: 0.0},
		{ID: models.Moon, Longitude: 120.0},
		{ID: models.Mars, Longitude: 90.0},
		{ID: models.Venus, Longitude: 60.0},
		{ID: models.Jupiter, Longitude: 180.0},
	}

	aspects := CalculateAspects(planets)

	t.Log("========== 相位计算测试 ==========")
	t.Logf("找到 %d 个相位", len(aspects))

	// 期望的相位
	expectedAspects := []struct {
		p1, p2 models.PlanetID
		kind   models.AspectType
	}{
		{models.Sun, models.Moon, models.Trine},         // 0° - 120° = 120° (三分相)
		{models.Sun, models.Mars, models.Square},        // 0° - 90° = 90° (四分相)
		{models.Sun, models.Venus, models.Sextile},      // 0° - 60° = 60° (六分相)
		{models.Sun, models.Jupiter, models.Opposition}, // 0° - 180° = 180° (对分相)
	}

	for _, exp := range expectedAspects {
		found := false
		for _, asp := range aspects {
			// 检查两种顺序
			match1 := asp.Planet1 == exp.p1 && asp.Planet2 == exp.p2
			match2 := asp.Planet1 == exp.p2 && asp.Planet2 == exp.p1
			if (match1 || match2) && asp.AspectType == exp.kind {
				found = true
				t.Logf("✓ %s-%s: %s (容许度=%.2f°, 强度=%.2f)",
					exp.p1, exp.p2, asp.AspectType, asp.Orb, asp.Strength)
				break
			}
		}
		if !found {
			t.Errorf("✗ 未找到期望的相位: %s-%s (%s)", exp.p1, exp.p2, exp.kind)
		}
	}
}

// BenchmarkSunPosition 性能测试
func BenchmarkSunPosition(b *testing.B) {
	jd := 2451545.0
	for i := 0; i < b.N; i++ {
		CalculateSunPosition(jd + float64(i))
	}
}

// BenchmarkMoonPosition 性能测试
func BenchmarkMoonPosition(b *testing.B) {
	jd := 2451545.0
	for i := 0; i < b.N; i++ {
		CalculateMoonPosition(jd + float64(i))
	}
}

// BenchmarkNatalChart 性能测试
func BenchmarkNatalChart(b *testing.B) {
	birthData := models.BirthData{
		Name:      "Test",
		Year:      1990,
		Month:     6,
		Day:       15,
		Hour:      14,
		Minute:    30,
		Latitude:  39.9042,
		Longitude: 116.4074,
		Timezone:  8,
	}

	for i := 0; i < b.N; i++ {
		CalculateNatalChart(birthData)
	}
}
