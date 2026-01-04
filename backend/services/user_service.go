package services

import (
	"errors"
	"star/astro"
	"star/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// 用户存储 (内存存储，可替换为数据库)
var (
	userStore = make(map[string]*models.User)
	userMutex sync.RWMutex
)

// CreateUser 创建用户
func CreateUser(name string, birthData models.BirthData) *models.User {
	userMutex.Lock()
	defer userMutex.Unlock()

	// 生成唯一ID
	id := "user_" + uuid.New().String()[:8]

	// 计算本命盘
	natalChart := astro.CalculateNatalChart(birthData)

	now := time.Now()
	user := &models.User{
		ID:         id,
		Name:       name,
		BirthData:  birthData,
		NatalChart: natalChart,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	userStore[id] = user
	return user
}

// GetUser 获取用户
func GetUser(id string) (*models.User, error) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	user, ok := userStore[id]
	if !ok {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// GetAllUsers 获取所有用户
func GetAllUsers() []*models.User {
	userMutex.RLock()
	defer userMutex.RUnlock()

	users := make([]*models.User, 0, len(userStore))
	for _, user := range userStore {
		users = append(users, user)
	}
	return users
}

// UpdateUser 更新用户
func UpdateUser(id, name string, birthData models.BirthData) (*models.User, error) {
	userMutex.Lock()
	defer userMutex.Unlock()

	user, ok := userStore[id]
	if !ok {
		return nil, errors.New("用户不存在")
	}

	// 更新信息
	user.Name = name
	user.BirthData = birthData
	user.NatalChart = astro.CalculateNatalChart(birthData)
	user.UpdatedAt = time.Now()

	return user, nil
}

// DeleteUser 删除用户
func DeleteUser(id string) error {
	userMutex.Lock()
	defer userMutex.Unlock()

	if _, ok := userStore[id]; !ok {
		return errors.New("用户不存在")
	}

	delete(userStore, id)
	return nil
}

// GetUserSnapshot 获取用户当前状态快照
func GetUserSnapshot(user *models.User) *models.UserSnapshot {
	now := time.Now()

	// 计算年龄
	age := now.Year() - user.BirthData.ToTime().Year()
	if now.YearDay() < user.BirthData.ToTime().YearDay() {
		age--
	}

	// 确保本命盘存在
	chart := user.NatalChart
	if chart == nil {
		chart = astro.CalculateNatalChart(user.BirthData)
	}

	// 计算每日预测
	dailyForecast := astro.CalculateDailyForecast(chart, now, true)

	// 计算年限法
	profection := astro.CalculateAnnualProfection(chart, age)

	// 获取当前活跃行运
	activeTransits := astro.GetCurrentTransits(chart, now)

	// 计算推运
	progressedChart := astro.CalculateProgressions(chart, now.Format("2006-01-02"))

	return &models.UserSnapshot{
		User:            user,
		CurrentDate:     now,
		Age:             age,
		DailyForecast:   dailyForecast,
		Profection:      profection,
		ActiveTransits:  activeTransits,
		ProgressedChart: progressedChart,
	}
}

// BatchGetUserSnapshots 批量获取用户快照
func BatchGetUserSnapshots(ids []string) []*models.UserSnapshot {
	snapshots := make([]*models.UserSnapshot, 0, len(ids))

	for _, id := range ids {
		user, err := GetUser(id)
		if err != nil {
			continue
		}
		snapshots = append(snapshots, GetUserSnapshot(user))
	}

	return snapshots
}

// GetAgentContext 获取智能体上下文
func GetAgentContext() *models.AgentContext {
	users := GetAllUsers()
	now := time.Now()

	// 构建用户上下文
	userContexts := make([]models.AgentUserContext, 0, len(users))
	for _, user := range users {
		chart := user.NatalChart
		if chart == nil {
			chart = astro.CalculateNatalChart(user.BirthData)
		}

		// 计算年龄
		age := now.Year() - user.BirthData.ToTime().Year()
		if now.YearDay() < user.BirthData.ToTime().YearDay() {
			age--
		}

		// 获取本命盘信息
		var sunSign, moonSign, ascendant models.ZodiacID
		for _, p := range chart.Planets {
			if p.ID == models.Sun {
				sunSign = p.Sign
			}
			if p.ID == models.Moon {
				moonSign = p.Sign
			}
		}
		ascZodiac := astro.GetZodiacByLongitude(chart.Ascendant)
		if ascZodiac != nil {
			ascendant = ascZodiac.ID
		}

		// 当前月亮信息
		transitPositions := astro.GetTransitPositions(now)
		var currentMoonSign models.ZodiacID
		for _, p := range transitPositions {
			if p.ID == models.Moon {
				currentMoonSign = p.Sign
			}
		}

		// 年限法
		profection := astro.CalculateAnnualProfection(chart, age)

		// 活跃行运
		activeTransits := astro.GetCurrentTransits(chart, now)
		majorTransitsActive := make([]string, 0)
		for _, t := range activeTransits {
			if t.Intensity > 70 {
				majorTransitsActive = append(majorTransitsActive, t.Interpretation.Theme)
			}
		}

		userContexts = append(userContexts, models.AgentUserContext{
			ID:   user.ID,
			Name: user.Name,
			CurrentState: map[string]interface{}{
				"sunSign":            sunSign,
				"moonSign":           moonSign,
				"ascendant":          ascendant,
				"currentMoonSign":    currentMoonSign,
				"profectionHouse":    profection.House,
				"profectionTheme":    profection.HouseTheme,
				"lordOfYear":         profection.LordOfYear,
				"majorTransitsActive": majorTransitsActive,
			},
		})
	}

	// 获取全局行运
	transitPositions := astro.GetTransitPositions(now)
	var sunSign, moonSign models.ZodiacID
	retrogradePlanets := make([]models.PlanetID, 0)

	for _, p := range transitPositions {
		if p.ID == models.Sun {
			sunSign = p.Sign
		}
		if p.ID == models.Moon {
			moonSign = p.Sign
		}
		if p.Retrograde {
			retrogradePlanets = append(retrogradePlanets, p.ID)
		}
	}

	// 计算月相
	moonPhase := astro.CalculateMoonPhase(
		transitPositions[0].Longitude, // Sun
		transitPositions[1].Longitude, // Moon
	)

	return &models.AgentContext{
		Users:       userContexts,
		CurrentDate: now,
		GlobalTransits: models.GlobalTransits{
			SunSign:           sunSign,
			MoonSign:          moonSign,
			MoonPhase:         moonPhase.Phase,
			RetrogradePlanets: retrogradePlanets,
		},
	}
}

// HandleAgentQuery 处理智能体查询
func HandleAgentQuery(userID, query string) *models.AgentQueryResponse {
	user, err := GetUser(userID)
	if err != nil {
		return &models.AgentQueryResponse{
			Response: "未找到该用户",
		}
	}

	// 获取用户快照
	snapshot := GetUserSnapshot(user)

	// 生成响应
	response := generateQueryResponse(query, snapshot)

	return &models.AgentQueryResponse{
		Response: response,
		Data: map[string]interface{}{
			"dailyForecast": snapshot.DailyForecast,
			"profection":    snapshot.Profection,
		},
	}
}

// generateQueryResponse 生成查询响应
func generateQueryResponse(query string, snapshot *models.UserSnapshot) string {
	// 基于规则的响应生成
	// 可以替换为 LLM 调用

	forecast := snapshot.DailyForecast

	return "今天整体运势分数为 " + formatScore(forecast.OverallScore) +
		"，主题是「" + forecast.OverallTheme + "」。" +
		"月亮目前在" + forecast.MoonSign.Name + "，" +
		"月相为" + forecast.MoonPhase.Name + "。" +
		"今年是" + snapshot.Profection.HouseName + "年，" +
		"主题是「" + snapshot.Profection.HouseTheme + "」。"
}

// formatScore 格式化分数
func formatScore(score float64) string {
	return string(rune(int(score))) + "分"
}

