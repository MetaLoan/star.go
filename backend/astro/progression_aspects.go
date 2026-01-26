package astro

import (
	"fmt"
	"math"
	"star/models"
	"time"
)

// ProgressionType defines the type of progression
type ProgressionType string

const (
	SecondaryProgression ProgressionType = "secondary" // 1 day = 1 year
	TertiaryProgression  ProgressionType = "tertiary"  // 1 day = 1 month
)

// ProgressionAspectEvent represents an aspect between progressed and natal planets
type ProgressionAspectEvent struct {
	ProgressionType  ProgressionType  `json:"progressionType"`  // secondary or tertiary
	ProgressedPlanet models.PlanetID  `json:"progressedPlanet"` // Progressed planet
	NatalPlanet      models.PlanetID  `json:"natalPlanet"`      // Natal planet
	AspectType       string           `json:"aspectType"`       // conjunction, sextile, etc.
	AspectAngle      float64          `json:"aspectAngle"`      // 0, 60, 90, 120, 180
	Orb              float64          `json:"orb"`              // Current orb
	IsApplying       bool             `json:"isApplying"`       // Approaching exact
	IsExact          bool             `json:"isExact"`          // Within 1 degree
	ExactDate        time.Time        `json:"exactDate"`        // When aspect is exact
	StartDate        time.Time        `json:"startDate"`        // When aspect enters orb
	EndDate          time.Time        `json:"endDate"`          // When aspect exits orb
	IsPositive       bool             `json:"isPositive"`       // Harmonious or challenging
	Intensity        string           `json:"intensity"`        // high, medium, low
	Title            string           `json:"title"`            // Event title
	Description      string           `json:"description"`      // Event description
	Theme            string           `json:"theme"`            // Psychological theme
	Advice           string           `json:"advice"`           // Practical advice
}

// ProgressionAspectConfig configuration for aspect calculation
type ProgressionAspectConfig struct {
	Angle      float64
	Name       string
	Orb        float64
	IsPositive bool
}

// ProgressionAspects defines aspects to check in progressions
var ProgressionAspects = []ProgressionAspectConfig{
	{0, "conjunction", 1.5, true},    // Conjunction - tighter orb for progressions
	{60, "sextile", 1.0, true},       // Sextile
	{90, "square", 1.5, false},       // Square
	{120, "trine", 1.0, true},        // Trine
	{180, "opposition", 1.5, false},  // Opposition
}

// CalculateSecondaryProgressions calculates secondary progressed positions
// Secondary progression: 1 day after birth = 1 year of life
func CalculateSecondaryProgressions(chart *models.NatalChart, targetDate time.Time) []models.PlanetPosition {
	birthDate := chart.BirthData.ToTime()
	
	// Calculate years from birth
	duration := targetDate.Sub(birthDate)
	yearsFromBirth := duration.Hours() / (365.25 * 24)
	
	// Progressed date: birth date + (years as days)
	progressedDate := birthDate.AddDate(0, 0, int(yearsFromBirth))
	
	// Get planet positions for progressed date
	progressedJd := DateToJulianDay(progressedDate)
	return GetPlanetPositionsUnified(progressedJd)
}

// CalculateTertiaryProgressions calculates tertiary progressed positions
// Tertiary progression: 1 day after birth = 1 month of life
func CalculateTertiaryProgressions(chart *models.NatalChart, targetDate time.Time) []models.PlanetPosition {
	birthDate := chart.BirthData.ToTime()
	
	// Calculate months from birth
	duration := targetDate.Sub(birthDate)
	monthsFromBirth := duration.Hours() / (30.44 * 24) // Average days per month
	
	// Progressed date: birth date + (months as days)
	progressedDate := birthDate.AddDate(0, 0, int(monthsFromBirth))
	
	// Get planet positions for progressed date
	progressedJd := DateToJulianDay(progressedDate)
	return GetPlanetPositionsUnified(progressedJd)
}

// CalculateProgressionToNatalAspects calculates aspects between progressed and natal planets
func CalculateProgressionToNatalAspects(progressed, natal []models.PlanetPosition, progType ProgressionType, targetDate time.Time, chart *models.NatalChart) []ProgressionAspectEvent {
	var events []ProgressionAspectEvent
	
	// Key progressed planets (Moon is most important in secondary)
	keyPlanets := []models.PlanetID{
		models.Sun, models.Moon, models.Mercury, models.Venus, models.Mars,
	}
	
	for _, progPlanet := range progressed {
		// Check if this is a key planet
		isKey := false
		for _, k := range keyPlanets {
			if progPlanet.ID == k {
				isKey = true
				break
			}
		}
		if !isKey {
			continue
		}
		
		// Check aspects to each natal planet
		for _, natalPlanet := range natal {
			for _, aspectDef := range ProgressionAspects {
				// Calculate angle between progressed and natal
				diff := math.Abs(progPlanet.Longitude - natalPlanet.Longitude)
				if diff > 180 {
					diff = 360 - diff
				}
				
				orb := math.Abs(diff - aspectDef.Angle)
				if orb <= aspectDef.Orb {
					event := createProgressionAspectEvent(
						progPlanet, natalPlanet, aspectDef, orb,
						progType, targetDate, chart,
					)
					events = append(events, event)
				}
			}
		}
	}
	
	return events
}

// createProgressionAspectEvent creates a progression aspect event
func createProgressionAspectEvent(
	progPlanet, natalPlanet models.PlanetPosition,
	aspectDef ProgressionAspectConfig,
	orb float64,
	progType ProgressionType,
	targetDate time.Time,
	chart *models.NatalChart,
) ProgressionAspectEvent {
	progInfo := GetPlanetInfo(progPlanet.ID)
	natalInfo := GetPlanetInfo(natalPlanet.ID)
	
	progName := ""
	natalName := ""
	if progInfo != nil {
		progName = progInfo.Name
	}
	if natalInfo != nil {
		natalName = natalInfo.Name
	}
	
	// Determine if applying or separating
	isApplying := orb > 0.5 && orb < aspectDef.Orb
	isExact := orb < 0.5
	
	// Calculate dates based on progression type and planet speed
	startDate, exactDate, endDate := estimateProgressionDates(progPlanet.ID, orb, aspectDef.Orb, progType, targetDate)
	
	// Intensity based on planets involved
	intensity := "medium"
	if progPlanet.ID == models.Moon || progPlanet.ID == models.Sun {
		intensity = "high"
	}
	
	// Create title based on progression type
	progTypeLabel := "SP"
	if progType == TertiaryProgression {
		progTypeLabel = "TP"
	}
	
	title := fmt.Sprintf("%s %s %s natal %s", progTypeLabel, progName, aspectDef.Name, natalName)
	
	return ProgressionAspectEvent{
		ProgressionType:  progType,
		ProgressedPlanet: progPlanet.ID,
		NatalPlanet:      natalPlanet.ID,
		AspectType:       aspectDef.Name,
		AspectAngle:      aspectDef.Angle,
		Orb:              orb,
		IsApplying:       isApplying,
		IsExact:          isExact,
		ExactDate:        exactDate,
		StartDate:        startDate,
		EndDate:          endDate,
		IsPositive:       aspectDef.IsPositive,
		Intensity:        intensity,
		Title:            title,
		Description:      generateProgressionDescription(progPlanet.ID, natalPlanet.ID, aspectDef, progType),
		Theme:            generateProgressionTheme(progPlanet.ID, natalPlanet.ID, aspectDef.IsPositive),
		Advice:           generateProgressionAdvice(progPlanet.ID, natalPlanet.ID, aspectDef.IsPositive),
	}
}

// estimateProgressionDates estimates start, exact, and end dates for progression aspect
func estimateProgressionDates(planet models.PlanetID, currentOrb, maxOrb float64, progType ProgressionType, now time.Time) (time.Time, time.Time, time.Time) {
	// Planet speeds in degrees per day (in progressed time)
	speeds := map[models.PlanetID]float64{
		models.Moon:    13.0,   // ~13 degrees per day
		models.Sun:     1.0,    // ~1 degree per day
		models.Mercury: 1.5,    // variable
		models.Venus:   1.2,    // ~1.2 degrees per day
		models.Mars:    0.5,    // ~0.5 degrees per day
	}
	
	speed := speeds[planet]
	if speed == 0 {
		speed = 0.5 // Default
	}
	
	// In secondary progression, 1 day = 1 year
	// In tertiary progression, 1 day = 1 month
	var multiplier float64
	switch progType {
	case SecondaryProgression:
		multiplier = 365.25 // 1 degree of motion = ~1 year
	case TertiaryProgression:
		multiplier = 30.44 // 1 degree of motion = ~1 month
	}
	
	// Time to move through orb (in real days)
	daysPerDegree := multiplier / speed
	
	// Estimate dates
	daysToExact := currentOrb * daysPerDegree
	daysFromStart := (maxOrb - currentOrb) * daysPerDegree
	daysToEnd := maxOrb * daysPerDegree
	
	startDate := now.AddDate(0, 0, -int(daysFromStart))
	exactDate := now.AddDate(0, 0, int(daysToExact))
	endDate := now.AddDate(0, 0, int(daysToEnd))
	
	return startDate, exactDate, endDate
}

// generateProgressionDescription generates description for progression aspect
func generateProgressionDescription(progPlanet, natalPlanet models.PlanetID, aspectDef ProgressionAspectConfig, progType ProgressionType) string {
	progInfo := GetPlanetInfo(progPlanet)
	natalInfo := GetPlanetInfo(natalPlanet)
	
	if progInfo == nil || natalInfo == nil {
		return ""
	}
	
	typeLabel := "Secondary"
	if progType == TertiaryProgression {
		typeLabel = "Tertiary"
	}
	
	if aspectDef.IsPositive {
		return fmt.Sprintf("%s progressed %s forms harmonious %s with your natal %s, bringing opportunities for growth and positive developments.",
			typeLabel, progInfo.Name, aspectDef.Name, natalInfo.Name)
	}
	return fmt.Sprintf("%s progressed %s forms challenging %s with your natal %s, indicating a period of tension that promotes personal evolution.",
		typeLabel, progInfo.Name, aspectDef.Name, natalInfo.Name)
}

// generateProgressionTheme generates theme for progression aspect
func generateProgressionTheme(progPlanet, natalPlanet models.PlanetID, isPositive bool) string {
	// Key themes for planet combinations
	themes := map[models.PlanetID]map[models.PlanetID]string{
		models.Moon: {
			models.Sun:     "Emotional alignment with life purpose",
			models.Mercury: "Emotional intelligence and communication",
			models.Venus:   "Emotional fulfillment in relationships",
			models.Mars:    "Emotional drive and assertion",
			models.Jupiter: "Emotional expansion and optimism",
			models.Saturn:  "Emotional maturity and responsibility",
		},
		models.Sun: {
			models.Moon:    "Identity and emotional integration",
			models.Mercury: "Self-expression and communication",
			models.Venus:   "Creative self-expression and values",
			models.Mars:    "Willpower and action",
			models.Jupiter: "Personal growth and expansion",
			models.Saturn:  "Life purpose and structure",
		},
		models.Venus: {
			models.Mars:    "Love and passion integration",
			models.Jupiter: "Relationship expansion and joy",
			models.Saturn:  "Commitment and relationship maturity",
		},
	}
	
	if planetThemes, ok := themes[progPlanet]; ok {
		if theme, ok := planetThemes[natalPlanet]; ok {
			return theme
		}
	}
	
	if isPositive {
		return "Harmonious development and integration"
	}
	return "Growth through challenge and adjustment"
}

// generateProgressionAdvice generates advice for progression aspect
func generateProgressionAdvice(progPlanet, natalPlanet models.PlanetID, isPositive bool) string {
	progInfo := GetPlanetInfo(progPlanet)
	
	if progInfo == nil {
		return "Pay attention to this developmental phase."
	}
	
	if isPositive {
		return fmt.Sprintf("This is a favorable period for %s-related matters. Take advantage of the supportive energy.", progInfo.Name)
	}
	return fmt.Sprintf("Navigate this period with awareness. The %s energy requires patience and conscious effort.", progInfo.Name)
}

// GetSecondaryProgressionEvents returns secondary progression aspects active on a date
func GetSecondaryProgressionEvents(chart *models.NatalChart, targetDate time.Time) []ProgressionAspectEvent {
	// Calculate secondary progressed positions
	progressedPositions := CalculateSecondaryProgressions(chart, targetDate)
	
	// Calculate aspects to natal
	return CalculateProgressionToNatalAspects(
		progressedPositions,
		chart.Planets,
		SecondaryProgression,
		targetDate,
		chart,
	)
}

// GetTertiaryProgressionEvents returns tertiary progression aspects active on a date
func GetTertiaryProgressionEvents(chart *models.NatalChart, targetDate time.Time) []ProgressionAspectEvent {
	// Calculate tertiary progressed positions
	progressedPositions := CalculateTertiaryProgressions(chart, targetDate)
	
	// Calculate aspects to natal
	return CalculateProgressionToNatalAspects(
		progressedPositions,
		chart.Planets,
		TertiaryProgression,
		targetDate,
		chart,
	)
}

// CalculateProgressionFactors calculates influence factors from progressions
func CalculateProgressionFactors(chart *models.NatalChart, date time.Time, weight float64, progType ProgressionType) []models.InfluenceFactor {
	var factors []models.InfluenceFactor
	var events []ProgressionAspectEvent
	
	switch progType {
	case SecondaryProgression:
		events = GetSecondaryProgressionEvents(chart, date)
	case TertiaryProgression:
		events = GetTertiaryProgressionEvents(chart, date)
	}
	
	for _, e := range events {
		// Base value based on positivity and orb
		baseValue := 0.5
		if e.IsPositive {
			baseValue = 0.7
		} else {
			baseValue = 0.3
		}
		
		// Orb factor: closer to exact = stronger
		orbFactor := 1.0 - (e.Orb / 1.5) // Max orb is ~1.5
		baseValue *= (0.5 + 0.5*orbFactor)
		
		// Get dimension impact based on planets
		dimensionImpact := getProgressionDimensionImpact(e.ProgressedPlanet, e.NatalPlanet)
		
		// Time level based on progression type
		timeLevel := models.TimeLevelYearly
		if progType == TertiaryProgression {
			timeLevel = models.TimeLevelMonthly
		}
		
		// Lifecycle
		lifecycle := &models.FactorLifecycle{
			StartTime: e.StartDate,
			PeakTime:  e.ExactDate,
			EndTime:   e.EndDate,
			Duration:  e.EndDate.Sub(e.StartDate).Hours(),
		}
		
		// Current strength
		strength := CalculateWaveStrength(lifecycle.StartTime, lifecycle.PeakTime, lifecycle.EndTime, date)
		
		factorType := models.InfluenceFactorType("secondary")
		typeLabel := "SP"
		if progType == TertiaryProgression {
			factorType = models.InfluenceFactorType("tertiary")
			typeLabel = "TP"
		}
		
		factor := models.InfluenceFactor{
			ID:              fmt.Sprintf("%s_%d_%d_%s_%s", typeLabel, e.ProgressedPlanet, e.NatalPlanet, e.AspectType, date.Format("20060102")),
			Type:            factorType,
			Name:            e.Title,
			Description:     e.Description,
			TimeLevel:       timeLevel,
			Lifecycle:       lifecycle,
			BaseValue:       baseValue,
			Weight:          weight,
			CurrentStrength: strength,
			Adjustment:      baseValue * weight * strength,
			DimensionImpact: dimensionImpact,
			SourcePlanet:    e.ProgressedPlanet,
			IsPositive:      e.IsPositive,
			AstroReason:     fmt.Sprintf("%s %s forming %s to natal %s", typeLabel, GetPlanetInfo(e.ProgressedPlanet).Name, e.AspectType, GetPlanetInfo(e.NatalPlanet).Name),
		}
		
		factors = append(factors, factor)
	}
	
	return factors
}

// getProgressionDimensionImpact returns dimension impact based on planets
func getProgressionDimensionImpact(progPlanet, natalPlanet models.PlanetID) models.DimensionImpact {
	// Combine dimension impacts of both planets
	progImpact := GetPlanetDimensionImpact(progPlanet)
	natalImpact := GetPlanetDimensionImpact(natalPlanet)
	
	return models.DimensionImpact{
		Career:       (progImpact.Career + natalImpact.Career) / 2,
		Relationship: (progImpact.Relationship + natalImpact.Relationship) / 2,
		Health:       (progImpact.Health + natalImpact.Health) / 2,
		Finance:      (progImpact.Finance + natalImpact.Finance) / 2,
		Spiritual:    (progImpact.Spiritual + natalImpact.Spiritual) / 2,
	}
}
