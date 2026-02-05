package i18n

import (
	"testing"
)

func TestGetDetailedInterpretation(t *testing.T) {
	langs := []string{"zh", "en", "ru"}
	
	for _, lang := range langs {
		translator := New(lang)
		
		t.Run(lang+"_aspect", func(t *testing.T) {
			// Test major aspect
			text := translator.GetDetailedInterpretation("aspect", "sun", "moon", "conjunction", "", true)
			if text == "" {
				t.Errorf("Empty interpretation for sun_conjunction_moon in %s", lang)
			}
			
			// Test minor aspect (fallback/template)
			text = translator.GetDetailedInterpretation("aspect", "sun", "moon", "semi-sextile", "", true)
			if text == "" {
				t.Errorf("Empty interpretation for sun_semi-sextile_moon in %s", lang)
			}
		})
		
		t.Run(lang+"_progression", func(t *testing.T) {
			// Test specific progression
			text := translator.GetDetailedInterpretation("secondary_progression", "sun", "moon", "trine", "", true)
			if text == "" {
				t.Errorf("Empty interpretation for progression sun_trine_moon in %s", lang)
			}
			
			// Test fallback to aspect
			text = translator.GetDetailedInterpretation("secondary_progression", "uranus", "pluto", "trine", "", true)
			if text == "" {
				t.Errorf("Empty interpretation for progression fallback uranus_trine_pluto in %s", lang)
			}
		})
		
		t.Run(lang+"_transit_house", func(t *testing.T) {
			text := translator.GetDetailedInterpretation("transit_house", "sun", "", "", "1", true)
			if text == "" {
				t.Errorf("Empty interpretation for sun_house_1 in %s", lang)
			}
		})
		
		t.Run(lang+"_other_events", func(t *testing.T) {
			events := []string{"retrograde", "profectionLord", "planetaryHour", "voidOfCourse", "lunar_phase", "sign_change", "dignity"}
			for _, event := range events {
				var text string
				if event == "voidOfCourse" {
					text = translator.GetDetailedInterpretation(event, "", "", "", "", true)
				} else {
					text = translator.GetDetailedInterpretation(event, "sun", "", "conjunction", "", true)
				}
				if text == "" {
					t.Errorf("Empty interpretation for %s in %s", event, lang)
				}
			}
		})
	}
}
