package renderer

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the visual styling and glyphs for the game.
type Theme struct {
	Name            string
	HeadGlyph       func(dirString string) string
	BodyGlyph       string
	FoodGlyph       string
	SlowMoGlyph     string
	GhostGlyph      string
	MultiplierGlyph string
	ShrinkGlyph     string
	PoisonGlyph     string
	HeadStyle       lipgloss.Style
	BodyStyle       lipgloss.Style
	FoodStyle       lipgloss.Style
	SlowMoStyle     lipgloss.Style
	GhostStyle      lipgloss.Style
	MultiplierStyle lipgloss.Style
	ShrinkStyle     lipgloss.Style
	PoisonStyle     lipgloss.Style
	BorderStyle     lipgloss.Style
	TitleStyle      lipgloss.Style
	ScoreStyle      lipgloss.Style
	HighScoreStyle  lipgloss.Style
	HUDLabelStyle   lipgloss.Style
	ModalStyle      lipgloss.Style
}

// GetTheme returns a theme by name, defaulting to "dark".
func GetTheme(name string) Theme {
	switch name {
	case "light":
		return LightTheme()
	case "high-contrast":
		return HighContrastTheme()
	case "colorblind":
		return ColorblindSafeTheme()
	case "corporate":
		return CorporateTheme()
	case "dark":
		fallthrough
	default:
		return DarkTheme()
	}
}

// DarkTheme represents the premier modern neon cyber aesthetic.
func DarkTheme() Theme {
	return Theme{
		Name: "dark",
		HeadGlyph: func(dir string) string {
			switch dir {
			case "UP":
				return "▲ "
			case "DOWN":
				return "▼ "
			case "LEFT":
				return "◀ "
			case "RIGHT":
				return "▶ "
			default:
				return "● "
			}
		},
		BodyGlyph:       "● ",
		FoodGlyph:       "★ ",
		SlowMoGlyph:     "⏱ ",
		GhostGlyph:      "👻",
		MultiplierGlyph: "⚡",
		ShrinkGlyph:     "✂ ",
		PoisonGlyph:     "☠ ",
		HeadStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true),
		BodyStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#00DF89")),
		FoodStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0055")).Bold(true),
		SlowMoStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Bold(true),
		GhostStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#D500F9")).Bold(true),
		MultiplierStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD600")).Bold(true),
		ShrinkStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6D00")).Bold(true),
		PoisonStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0055")).Bold(true),
		BorderStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#7928CA")).Border(lipgloss.RoundedBorder()).Padding(0, 0),
		TitleStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true),
		ScoreStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Bold(true),
		HighScoreStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD600")).Bold(true),
		HUDLabelStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")),
		ModalStyle: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FF0055")).
			Padding(1, 2).
			Bold(true),
	}
}

// LightTheme represents a crisp, clean aesthetic for bright terminals.
func LightTheme() Theme {
	return Theme{
		Name: "light",
		HeadGlyph: func(dir string) string {
			return "◆ "
		},
		BodyGlyph:       "■ ",
		FoodGlyph:       "● ",
		SlowMoGlyph:     "S ",
		GhostGlyph:      "G ",
		MultiplierGlyph: "M ",
		ShrinkGlyph:     "K ",
		PoisonGlyph:     "X ",
		HeadStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#005F73")).Bold(true),
		BodyStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#0A9396")),
		FoodStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#AE2012")).Bold(true),
		SlowMoStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#0077B6")).Bold(true),
		GhostStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#7209B7")).Bold(true),
		MultiplierStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#CA6702")).Bold(true),
		ShrinkStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#BB3E03")).Bold(true),
		PoisonStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#AE2012")).Bold(true),
		BorderStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#94D2BD")).Border(lipgloss.NormalBorder()),
		TitleStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#005F73")).Bold(true),
		ScoreStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#0A9396")).Bold(true),
		HighScoreStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#EE9B00")).Bold(true),
		HUDLabelStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")),
		ModalStyle: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#AE2012")).
			Padding(1, 2),
	}
}

// HighContrastTheme provides pure black & white and max accessibility.
func HighContrastTheme() Theme {
	return Theme{
		Name: "high-contrast",
		HeadGlyph: func(dir string) string {
			return "@ "
		},
		BodyGlyph:       "O ",
		FoodGlyph:       "* ",
		SlowMoGlyph:     "S ",
		GhostGlyph:      "G ",
		MultiplierGlyph: "X ",
		ShrinkGlyph:     "- ",
		PoisonGlyph:     "! ",
		HeadStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#000000")).Bold(true),
		BodyStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		FoodStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Blink(true),
		SlowMoStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		GhostStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		MultiplierStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		ShrinkStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		PoisonStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Blink(true),
		BorderStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Border(lipgloss.ThickBorder()),
		TitleStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		ScoreStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		HighScoreStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		HUDLabelStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		ModalStyle: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Padding(1, 2),
	}
}

// ColorblindSafeTheme uses blue/orange/yellow accessible palettes (Okabe-Ito).
func ColorblindSafeTheme() Theme {
	return Theme{
		Name: "colorblind",
		HeadGlyph: func(dir string) string {
			return "▶ "
		},
		BodyGlyph:       "■ ",
		FoodGlyph:       "★ ",
		SlowMoGlyph:     "❄ ",
		GhostGlyph:      "◇ ",
		MultiplierGlyph: "+ ",
		ShrinkGlyph:     "/ ",
		PoisonGlyph:     "✖ ",
		HeadStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#0072B2")).Bold(true), // Blue
		BodyStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#56B4E9")),             // Sky Blue
		FoodStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#D55E00")).Bold(true), // Vermilion
		SlowMoStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#009E73")).Bold(true), // Bluish green
		GhostStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#CC79A7")).Bold(true), // Reddish purple
		MultiplierStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#F0E442")).Bold(true), // Yellow
		ShrinkStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#E69F00")).Bold(true), // Orange
		PoisonStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#D55E00")).Bold(true), // Vermilion hazard
		BorderStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#0072B2")).Border(lipgloss.RoundedBorder()),
		TitleStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#0072B2")).Bold(true),
		ScoreStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#56B4E9")).Bold(true),
		HighScoreStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#F0E442")).Bold(true),
		HUDLabelStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#999999")),
		ModalStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#D55E00")).
			Padding(1, 2),
	}
}

// CorporateTheme blends inconspicuously into an enterprise terminal workspace.
func CorporateTheme() Theme {
	return Theme{
		Name: "corporate",
		HeadGlyph: func(dir string) string {
			return "# "
		},
		BodyGlyph:       "= ",
		FoodGlyph:       "$ ",
		SlowMoGlyph:     "s ",
		GhostGlyph:      "g ",
		MultiplierGlyph: "x ",
		ShrinkGlyph:     "- ",
		PoisonGlyph:     "! ",
		HeadStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")).Bold(true),
		BodyStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#707070")),
		FoodStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#C0C0C0")).Bold(true),
		SlowMoStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#8A8A8A")),
		GhostStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#8A8A8A")),
		MultiplierStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")),
		ShrinkStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#8A8A8A")),
		PoisonStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#C0392B")).Bold(true),
		BorderStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Border(lipgloss.NormalBorder()),
		TitleStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")),
		ScoreStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#D0D0D0")),
		HighScoreStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#B0B0B0")),
		HUDLabelStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")),
		ModalStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#707070")).
			Padding(1, 2),
	}
}
