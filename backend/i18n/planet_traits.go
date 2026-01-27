package i18n

import "star/models"

// PlanetTraits holds thematic keywords and descriptions for each planet,
// used for Tier 2/3 aspect interpretation fallbacks.
type PlanetTraits struct {
	Themes   []string // career, relationship, health, finance, spiritual
	Keywords []string // energy descriptors
	Domain   string   // short domain description
}

// PlanetTraitMap maps PlanetID to traits. Keys use string(PlanetID).
var PlanetTraitMap = map[string]PlanetTraits{
	"sun": {
		Themes:   []string{"identity", "vitality", "career", "health"},
		Keywords: []string{"confidence", "will", "leadership", "self-expression"},
		Domain:   "identity and vitality",
	},
	"moon": {
		Themes:   []string{"emotions", "family", "relationship", "home"},
		Keywords: []string{"feelings", "nurturing", "intuition", "security"},
		Domain:   "emotions and inner life",
	},
	"mercury": {
		Themes:   []string{"communication", "career", "learning"},
		Keywords: []string{"thought", "ideas", "logic", "expression"},
		Domain:   "mind and communication",
	},
	"venus": {
		Themes:   []string{"relationship", "finance", "beauty"},
		Keywords: []string{"love", "harmony", "values", "pleasure"},
		Domain:   "love and values",
	},
	"mars": {
		Themes:   []string{"action", "career", "health"},
		Keywords: []string{"drive", "courage", "assertion", "passion"},
		Domain:   "action and desire",
	},
	"jupiter": {
		Themes:   []string{"growth", "finance", "spiritual"},
		Keywords: []string{"expansion", "luck", "wisdom", "opportunity"},
		Domain:   "expansion and growth",
	},
	"saturn": {
		Themes:   []string{"career", "structure", "responsibility"},
		Keywords: []string{"discipline", "maturity", "limits", "reality"},
		Domain:   "structure and responsibility",
	},
	"uranus": {
		Themes:   []string{"change", "spiritual", "innovation"},
		Keywords: []string{"awakening", "freedom", "sudden change", "rebellion"},
		Domain:   "sudden change and innovation",
	},
	"neptune": {
		Themes:   []string{"spiritual", "creativity", "dreams"},
		Keywords: []string{"inspiration", "compassion", "imagination", "transcendence"},
		Domain:   "inspiration and transcendence",
	},
	"pluto": {
		Themes:   []string{"transformation", "power", "spiritual"},
		Keywords: []string{"rebirth", "depth", "intensity", "regeneration"},
		Domain:   "transformation and depth",
	},
	"northNode": {
		Themes:   []string{"growth", "spiritual", "direction"},
		Keywords: []string{"evolution", "destiny", "soul direction"},
		Domain:   "soul growth and direction",
	},
	"chiron": {
		Themes:   []string{"healing", "wounds", "spiritual"},
		Keywords: []string{"wounding", "healing", "wisdom through pain"},
		Domain:   "healing and wounded wisdom",
	},
}

// GetPlanetTraits returns traits for a planet, or empty traits if unknown.
func GetPlanetTraits(p models.PlanetID) PlanetTraits {
	if t, ok := PlanetTraitMap[string(p)]; ok {
		return t
	}
	return PlanetTraits{
		Themes:   []string{"growth"},
		Keywords: []string{"energy"},
		Domain:   "personal development",
	}
}
