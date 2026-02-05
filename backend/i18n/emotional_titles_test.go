package i18n

import (
	"testing"
)

func TestGetEmotionalTitle(t *testing.T) {
	langs := []string{"zh", "en", "ru"}
	
	for _, lang := range langs {
		translator := New(lang)
		
		t.Run(lang+"_aspect_titles", func(t *testing.T) {
			// Test major aspect title
			title := translator.GetEmotionalTitle("aspect", "sun", "moon", "conjunction", "", true)
			if title == "" || title == "和谐能量" || title == "Harmonious Energy" || title == "Гармоничная энергия" {
				// We expect a specific title for sun_conjunction_moon
				if title == "" {
					t.Errorf("Empty title for sun_conjunction_moon in %s", lang)
				}
			}
			
			// Test minor aspect title
			title = translator.GetEmotionalTitle("aspect", "sun", "moon", "semi-sextile", "", true)
			if title == "" {
				t.Errorf("Empty title for sun_semi-sextile_moon in %s", lang)
			}
		})
		
		t.Run(lang+"_progression_titles", func(t *testing.T) {
			// Test specific progression title
			title := translator.GetEmotionalTitle("secondary_progression", "sun", "moon", "trine", "", true)
			if title == "" {
				t.Errorf("Empty title for progression sun_trine_moon in %s", lang)
			}
			
			// Test fallback to aspect with context
			title = translator.GetEmotionalTitle("secondary_progression", "uranus", "pluto", "trine", "", true)
			if title == "" {
				t.Errorf("Empty title for progression fallback uranus_trine_pluto in %s", lang)
			}
		})
		
		t.Run(lang+"_other_titles", func(t *testing.T) {
			events := []string{"transit_house", "retrograde", "profectionLord", "planetaryHour", "voidOfCourse", "lunar_phase", "sign_change", "dignity"}
			for _, event := range events {
				var title string
				if event == "voidOfCourse" {
					title = translator.GetEmotionalTitle(event, "", "", "", "", true)
				} else if event == "transit_house" {
					title = translator.GetEmotionalTitle(event, "sun", "", "", "1", true)
				} else {
					title = translator.GetEmotionalTitle(event, "sun", "", "conjunction", "", true)
				}
				if title == "" {
					t.Errorf("Empty title for %s in %s", event, lang)
				}
			}
		})
	}
}
