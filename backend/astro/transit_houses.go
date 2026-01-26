package astro

import (
	"fmt"
	"math"
	"star/models"
	"time"
)

// TransitHouseEvent represents a transit planet passing through a natal house
type TransitHouseEvent struct {
	Planet      models.PlanetID `json:"planet"`      // Transit planet
	House       int             `json:"house"`       // Natal house number (1-12)
	HouseTheme  string          `json:"houseTheme"`  // Theme of the house
	EntryTime   time.Time       `json:"entryTime"`   // When planet enters house
	ExitTime    time.Time       `json:"exitTime"`    // When planet exits house
	DurationDays float64        `json:"durationDays"` // Duration in days
	IsPositive  bool            `json:"isPositive"`  // Overall positive/negative
	Intensity   string          `json:"intensity"`   // high, medium, low
	Title       string          `json:"title"`       // Event title
	Description string          `json:"description"` // Event description
	Theme       string          `json:"theme"`       // Combined theme
	Advice      string          `json:"advice"`      // Practical advice
}

// HouseThemes defines themes for each house
var HouseThemes = map[int]struct {
	Name        string
	Keywords    []string
	Theme       string
	PositivePlanets []models.PlanetID
}{
	1:  {"1st House - Self", []string{"identity", "appearance", "initiative"}, "Self-expression and new beginnings", []models.PlanetID{models.Sun, models.Jupiter, models.Venus}},
	2:  {"2nd House - Resources", []string{"money", "values", "possessions"}, "Financial matters and self-worth", []models.PlanetID{models.Jupiter, models.Venus}},
	3:  {"3rd House - Communication", []string{"learning", "siblings", "short trips"}, "Communication and mental activity", []models.PlanetID{models.Mercury, models.Jupiter}},
	4:  {"4th House - Home", []string{"family", "roots", "security"}, "Home, family and emotional foundation", []models.PlanetID{models.Moon, models.Venus, models.Jupiter}},
	5:  {"5th House - Creativity", []string{"romance", "children", "pleasure"}, "Creative expression and joy", []models.PlanetID{models.Sun, models.Venus, models.Jupiter}},
	6:  {"6th House - Health", []string{"work", "service", "routine"}, "Daily health and work habits", []models.PlanetID{models.Mercury, models.Mars}},
	7:  {"7th House - Partnership", []string{"marriage", "contracts", "others"}, "Relationships and partnerships", []models.PlanetID{models.Venus, models.Jupiter}},
	8:  {"8th House - Transformation", []string{"intimacy", "shared resources", "rebirth"}, "Deep transformation and shared resources", []models.PlanetID{models.Pluto, models.Mars}},
	9:  {"9th House - Expansion", []string{"travel", "philosophy", "higher learning"}, "Higher education and spiritual growth", []models.PlanetID{models.Jupiter, models.Sun}},
	10: {"10th House - Career", []string{"profession", "reputation", "achievement"}, "Career and public standing", []models.PlanetID{models.Sun, models.Jupiter, models.Saturn}},
	11: {"11th House - Community", []string{"friends", "groups", "hopes"}, "Social connections and future goals", []models.PlanetID{models.Jupiter, models.Uranus, models.Venus}},
	12: {"12th House - Spirituality", []string{"subconscious", "retreat", "karma"}, "Inner world and spiritual matters", []models.PlanetID{models.Neptune, models.Jupiter}},
}

// TransitPlanetDurations approximate days for each planet to transit one house
var TransitPlanetDurations = map[models.PlanetID]float64{
	models.Moon:    2.5,   // ~2.5 days per sign
	models.Sun:     30,    // ~30 days per sign
	models.Mercury: 21,    // ~21 days (varies with retro)
	models.Venus:   25,    // ~25 days (varies with retro)
	models.Mars:    45,    // ~45 days per sign
	models.Jupiter: 365,   // ~1 year per sign
	models.Saturn:  912,   // ~2.5 years per sign
	models.Uranus:  2555,  // ~7 years per sign
	models.Neptune: 5110,  // ~14 years per sign
	models.Pluto:   7300,  // ~20 years per sign
}

// CalculateTransitHouseEvents calculates transit planets through natal houses
func CalculateTransitHouseEvents(chart *models.NatalChart, date time.Time) []TransitHouseEvent {
	var events []TransitHouseEvent

	// Get natal house cusps
	natalHouses := chart.Houses

	// Get current transit positions
	jd := DateToJulianDay(date)
	transitPositions := GetPlanetPositionsUnified(jd)

	// Main planets to consider for transit house events
	mainPlanets := map[models.PlanetID]bool{
		models.Sun: true, models.Moon: true, models.Mercury: true,
		models.Venus: true, models.Mars: true, models.Jupiter: true,
		models.Saturn: true, models.Uranus: true, models.Neptune: true,
		models.Pluto: true,
	}
	
	// Check each transit planet
	for _, transit := range transitPositions {
		// Only process main planets
		if !mainPlanets[transit.ID] {
			continue
		}

		// Find which natal house the transit planet is in
		house := GetPlanetHouse(transit.Longitude, natalHouses)
		
		// Calculate entry and exit times
		entryTime, exitTime := estimateHouseTransitTimes(transit.ID, transit.Longitude, house, natalHouses, date)
		
		// Create event
		event := createTransitHouseEvent(transit, house, entryTime, exitTime, chart)
		events = append(events, event)
	}

	return events
}

// estimateHouseTransitTimes estimates when a planet entered and will exit a house
func estimateHouseTransitTimes(planetID models.PlanetID, currentLon float64, house int, houses []models.HouseCusp, now time.Time) (time.Time, time.Time) {
	// Get house cusp longitudes
	houseCusp := houses[house-1].Cusp
	nextHouse := house % 12
	nextCusp := houses[nextHouse].Cusp
	
	// Handle wraparound
	if nextCusp < houseCusp {
		nextCusp += 360
	}
	
	testLon := currentLon
	if testLon < houseCusp {
		testLon += 360
	}
	
	// Calculate how far planet has traveled in this house
	houseSize := nextCusp - houseCusp
	progressInHouse := testLon - houseCusp
	if progressInHouse < 0 {
		progressInHouse += 360
	}
	
	// Get typical duration for this planet
	typicalDuration := TransitPlanetDurations[planetID]
	if typicalDuration == 0 {
		typicalDuration = 30 // Default to 30 days
	}
	
	// Estimate time in house based on house size (average ~30 degrees)
	scaleFactor := houseSize / 30.0
	estimatedDuration := typicalDuration * scaleFactor
	
	// Calculate entry time (past)
	percentComplete := progressInHouse / houseSize
	daysSinceEntry := estimatedDuration * percentComplete
	entryTime := now.AddDate(0, 0, -int(daysSinceEntry))
	
	// Calculate exit time (future)
	daysUntilExit := estimatedDuration * (1 - percentComplete)
	exitTime := now.AddDate(0, 0, int(daysUntilExit))
	
	return entryTime, exitTime
}

// createTransitHouseEvent creates a transit house event
func createTransitHouseEvent(transit models.PlanetPosition, house int, entryTime, exitTime time.Time, chart *models.NatalChart) TransitHouseEvent {
	houseInfo := HouseThemes[house]
	planetInfo := GetPlanetInfo(transit.ID)
	planetName := ""
	if planetInfo != nil {
		planetName = planetInfo.Name
	}
	
	// Check if this is a positive combination
	isPositive := false
	for _, p := range houseInfo.PositivePlanets {
		if p == transit.ID {
			isPositive = true
			break
		}
	}
	
	// Challenging planets in certain houses
	if (transit.ID == models.Saturn || transit.ID == models.Pluto) && 
	   (house == 1 || house == 4 || house == 7 || house == 10) {
		isPositive = false
	}
	if transit.ID == models.Jupiter || transit.ID == models.Venus {
		isPositive = true
	}
	
	// Calculate intensity based on planet type
	intensity := "medium"
	if transit.ID == models.Jupiter || transit.ID == models.Saturn || 
	   transit.ID == models.Uranus || transit.ID == models.Neptune || transit.ID == models.Pluto {
		intensity = "high"
	}
	if transit.ID == models.Moon {
		intensity = "low"
	}
	
	// Duration
	duration := exitTime.Sub(entryTime).Hours() / 24
	
	return TransitHouseEvent{
		Planet:       transit.ID,
		House:        house,
		HouseTheme:   houseInfo.Name,
		EntryTime:    entryTime,
		ExitTime:     exitTime,
		DurationDays: duration,
		IsPositive:   isPositive,
		Intensity:    intensity,
		Title:        fmt.Sprintf("%s in %s", planetName, houseInfo.Name),
		Description:  generateTransitHouseDescription(transit.ID, house, isPositive),
		Theme:        generateTransitHouseTheme(transit.ID, house),
		Advice:       generateTransitHouseAdvice(transit.ID, house),
	}
}

// generateTransitHouseDescription generates description for transit house event
func generateTransitHouseDescription(planet models.PlanetID, house int, isPositive bool) string {
	planetInfo := GetPlanetInfo(planet)
	if planetInfo == nil {
		return ""
	}
	
	houseInfo := HouseThemes[house]
	
	if isPositive {
		return fmt.Sprintf("Transit %s brings beneficial energy to your %s, supporting themes of %s.", 
			planetInfo.Name, houseInfo.Name, houseInfo.Keywords[0])
	}
	return fmt.Sprintf("Transit %s activates your %s, bringing focus to %s and potential challenges.", 
		planetInfo.Name, houseInfo.Name, houseInfo.Keywords[0])
}

// generateTransitHouseTheme generates theme for transit house event
func generateTransitHouseTheme(planet models.PlanetID, house int) string {
	themes := map[models.PlanetID]map[int]string{
		models.Sun: {
			1: "Personal empowerment", 2: "Financial focus", 3: "Mental clarity",
			4: "Family matters", 5: "Creative expression", 6: "Health awareness",
			7: "Partnership focus", 8: "Deep transformation", 9: "Expansion of horizons",
			10: "Career spotlight", 11: "Social connections", 12: "Inner reflection",
		},
		models.Moon: {
			1: "Emotional awareness", 2: "Security needs", 3: "Emotional communication",
			4: "Home comfort", 5: "Playful mood", 6: "Health routines",
			7: "Relationship needs", 8: "Emotional depth", 9: "Seeking meaning",
			10: "Public image", 11: "Friendship needs", 12: "Intuitive time",
		},
		models.Mercury: {
			1: "Self-expression", 2: "Financial planning", 3: "Active learning",
			4: "Family discussions", 5: "Creative ideas", 6: "Work organization",
			7: "Negotiations", 8: "Deep research", 9: "Higher studies",
			10: "Professional communication", 11: "Networking", 12: "Inner dialogue",
		},
		models.Venus: {
			1: "Personal charm", 2: "Financial gains", 3: "Pleasant interactions",
			4: "Home beautification", 5: "Romance and fun", 6: "Work harmony",
			7: "Love relationships", 8: "Intimate bonding", 9: "Cultural appreciation",
			10: "Professional recognition", 11: "Social pleasures", 12: "Spiritual love",
		},
		models.Mars: {
			1: "Personal initiative", 2: "Earning drive", 3: "Assertive communication",
			4: "Home projects", 5: "Competitive spirit", 6: "Work energy",
			7: "Partnership dynamics", 8: "Intense transformation", 9: "Adventure seeking",
			10: "Career ambition", 11: "Group leadership", 12: "Hidden actions",
		},
		models.Jupiter: {
			1: "Personal growth", 2: "Financial expansion", 3: "Learning opportunities",
			4: "Family blessings", 5: "Joy and creativity", 6: "Health improvements",
			7: "Relationship growth", 8: "Transformative gains", 9: "Spiritual expansion",
			10: "Career advancement", 11: "Social opportunities", 12: "Inner wisdom",
		},
		models.Saturn: {
			1: "Self-discipline", 2: "Financial responsibility", 3: "Serious learning",
			4: "Family duties", 5: "Creative discipline", 6: "Work challenges",
			7: "Relationship tests", 8: "Deep restructuring", 9: "Philosophical maturity",
			10: "Career building", 11: "Social responsibilities", 12: "Karmic lessons",
		},
	}
	
	if planetThemes, ok := themes[planet]; ok {
		if theme, ok := planetThemes[house]; ok {
			return theme
		}
	}
	
	return fmt.Sprintf("Transit through house %d", house)
}

// generateTransitHouseAdvice generates advice for transit house event
func generateTransitHouseAdvice(planet models.PlanetID, house int) string {
	// General advice based on planet
	advices := map[models.PlanetID]string{
		models.Sun:     "Focus your energy on this life area. Take initiative and lead.",
		models.Moon:    "Pay attention to your emotional needs in this area.",
		models.Mercury: "Communicate clearly and gather information.",
		models.Venus:   "Seek harmony and appreciate beauty in this domain.",
		models.Mars:    "Take action but avoid impulsive decisions.",
		models.Jupiter: "Embrace opportunities for growth and expansion.",
		models.Saturn:  "Build solid foundations through discipline and patience.",
		models.Uranus:  "Be open to unexpected changes and innovations.",
		models.Neptune: "Trust your intuition but verify facts.",
		models.Pluto:   "Embrace transformation and release what no longer serves you.",
	}
	
	if advice, ok := advices[planet]; ok {
		return advice
	}
	return "Focus attention on this life area."
}

// GetActiveTransitHouseEvents returns transit house events active on a given date
func GetActiveTransitHouseEvents(chart *models.NatalChart, date time.Time, minIntensity string) []TransitHouseEvent {
	allEvents := CalculateTransitHouseEvents(chart, date)
	
	// Filter by intensity
	var filtered []TransitHouseEvent
	for _, e := range allEvents {
		include := false
		switch minIntensity {
		case "low":
			include = true
		case "medium":
			include = e.Intensity == "medium" || e.Intensity == "high"
		case "high":
			include = e.Intensity == "high"
		default:
			include = true
		}
		
		if include {
			filtered = append(filtered, e)
		}
	}
	
	return filtered
}

// CalculateTransitHouseFactors calculates influence factors from transit houses
func CalculateTransitHouseFactors(chart *models.NatalChart, date time.Time, weight float64) []models.InfluenceFactor {
	var factors []models.InfluenceFactor
	
	events := GetActiveTransitHouseEvents(chart, date, "medium")
	
	for _, e := range events {
		// Base value based on positivity
		baseValue := 0.5
		if e.IsPositive {
			baseValue = 0.7
		} else {
			baseValue = 0.3
		}
		
		// Adjust by intensity
		intensityMultiplier := 1.0
		switch e.Intensity {
		case "high":
			intensityMultiplier = 1.3
		case "medium":
			intensityMultiplier = 1.0
		case "low":
			intensityMultiplier = 0.7
		}
		
		// Get dimension impact based on house
		dimensionImpact := getHouseDimensionImpact(e.House)
		
		// Determine time level based on planet speed
		timeLevel := getTransitHouseTimeLevel(e.Planet)
		
		// Calculate lifecycle
		now := date
		lifecycle := &models.FactorLifecycle{
			StartTime: e.EntryTime,
			PeakTime:  e.EntryTime.Add(e.ExitTime.Sub(e.EntryTime) / 2),
			EndTime:   e.ExitTime,
			Duration:  e.DurationDays * 24,
		}
		
		// Calculate current strength
		strength := CalculateWaveStrength(lifecycle.StartTime, lifecycle.PeakTime, lifecycle.EndTime, now)
		
		factor := models.InfluenceFactor{
			ID:              fmt.Sprintf("transit_house_%d_%d_%s", e.Planet, e.House, date.Format("20060102")),
			Type:            models.FactorOuterPlanet,
			Name:            e.Title,
			Description:     e.Description,
			TimeLevel:       timeLevel,
			Lifecycle:       lifecycle,
			BaseValue:       baseValue * intensityMultiplier,
			Weight:          weight,
			CurrentStrength: strength,
			Adjustment:      baseValue * intensityMultiplier * weight * strength,
			DimensionImpact: dimensionImpact,
			SourcePlanet:    e.Planet,
			IsPositive:      e.IsPositive,
			AstroReason:     fmt.Sprintf("Transit %s through natal house %d", GetPlanetInfo(e.Planet).Name, e.House),
		}
		
		factors = append(factors, factor)
	}
	
	return factors
}

// getHouseDimensionImpact returns dimension impact based on house
func getHouseDimensionImpact(house int) models.DimensionImpact {
	// Designed for max 2 dimension labels
	// Single dimension houses: 2, 4, 7, 9, 10, 12
	// Dual dimension houses: 1, 3, 5, 6, 8, 11
	impacts := map[int]models.DimensionImpact{
		1:  {Health: 0.6, Career: 0.4, Relationship: 0, Finance: 0, Spiritual: 0},        // 健康+事业
		2:  {Finance: 0.8, Career: 0, Relationship: 0, Health: 0, Spiritual: 0},          // 财运
		3:  {Career: 0.5, Relationship: 0.5, Health: 0, Finance: 0, Spiritual: 0},        // 事业+关系
		4:  {Relationship: 0.8, Career: 0, Health: 0, Finance: 0, Spiritual: 0},          // 关系
		5:  {Relationship: 0.6, Spiritual: 0.4, Career: 0, Health: 0, Finance: 0},        // 关系+灵性
		6:  {Health: 0.6, Career: 0.4, Relationship: 0, Finance: 0, Spiritual: 0},        // 健康+事业
		7:  {Relationship: 0.9, Career: 0, Health: 0, Finance: 0, Spiritual: 0},          // 关系
		8:  {Spiritual: 0.6, Finance: 0.4, Career: 0, Relationship: 0, Health: 0},        // 灵性+财运
		9:  {Spiritual: 0.8, Career: 0, Relationship: 0, Health: 0, Finance: 0},          // 灵性
		10: {Career: 0.9, Relationship: 0, Health: 0, Finance: 0, Spiritual: 0},          // 事业
		11: {Relationship: 0.6, Spiritual: 0.4, Career: 0, Health: 0, Finance: 0},        // 关系+灵性
		12: {Spiritual: 0.8, Career: 0, Relationship: 0, Health: 0, Finance: 0},          // 灵性
	}
	
	if impact, ok := impacts[house]; ok {
		return impact
	}
	return models.DimensionImpact{Career: 0.5, Relationship: 0, Health: 0, Finance: 0, Spiritual: 0.5}
}

// getTransitHouseTimeLevel returns time level based on planet
func getTransitHouseTimeLevel(planet models.PlanetID) models.FactorTimeLevel {
	switch planet {
	case models.Moon:
		return models.TimeLevelHourly
	case models.Sun, models.Mercury, models.Venus:
		return models.TimeLevelDaily
	case models.Mars:
		return models.TimeLevelWeekly
	case models.Jupiter, models.Saturn:
		return models.TimeLevelMonthly
	case models.Uranus, models.Neptune, models.Pluto:
		return models.TimeLevelYearly
	default:
		return models.TimeLevelDaily
	}
}

// CalculateWaveStrength calculates strength using sine wave lifecycle
func CalculateWaveStrength(start, peak, end, now time.Time) float64 {
	if now.Before(start) || now.After(end) {
		return 0
	}
	
	totalDuration := end.Sub(start).Seconds()
	if totalDuration == 0 {
		return 1.0
	}
	
	elapsed := now.Sub(start).Seconds()
	progress := elapsed / totalDuration
	
	// Sine wave: 0 at start, 1 at peak, 0 at end
	strength := math.Sin(progress * math.Pi)
	
	return strength
}
