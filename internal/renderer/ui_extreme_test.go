package renderer

import (
	"strings"
	"testing"

	"snake/internal/engine"
	"snake/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func TestExtreme_ModalBoxesFormatting(t *testing.T) {
	reasons := []string{
		"Bit own tail",
		"Crashed into wall",
		"Crashed into obstacle",
		"Time is up!",
		"Victory! Board fully cleared!",
		"An extremely long custom game over reason that exceeds the box width by a lot",
	}

	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	for _, reason := range reasons {
		g := newTestGame()
		g.GameOver = true
		g.GameOverReason = reason
		g.Score = 9999
		model := NewUIModel(g, cfg, stats, "dark", false)
		model.TermWidth = 80
		model.TermHeight = 24

		view := model.View()
		if !strings.Contains(view, "GAME OVER") {
			t.Errorf("View missing 'GAME OVER' for reason '%s'", reason)
		}
		if !strings.Contains(view, "Final Score:") {
			t.Errorf("View missing 'Final Score:' for reason '%s'", reason)
		}

		// Ensure box doesn't throw runtime panic or render blank
		if len(view) < 100 {
			t.Errorf("View suspiciously short for reason '%s':\n%s", reason, view)
		}
	}
}

func TestExtreme_BossScreensRendering(t *testing.T) {
	screens := []string{"htop", "build", "git", "default"}
	for _, screen := range screens {
		rendered := RenderBossScreen(screen, 80, 24)
		if len(rendered) < 50 {
			t.Errorf("Boss screen '%s' rendered empty or too short", screen)
		}
		switch screen {
		case "htop":
			if !strings.Contains(rendered, "Load average") {
				t.Errorf("htop screen missing Load average")
			}
		case "build":
			if !strings.Contains(rendered, "Built target") && !strings.Contains(rendered, "PASS") {
				t.Errorf("build screen missing compile logs")
			}
		case "git":
			if !strings.Contains(rendered, "commit") {
				t.Errorf("git screen missing commit logs")
			}
		}
	}
}

func TestExtreme_AllThemesRendering(t *testing.T) {
	themeNames := []string{"dark", "light", "high-contrast", "colorblind", "corporate"}
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	for _, name := range themeNames {
		th := GetTheme(name)
		if th.Name != name {
			t.Errorf("GetTheme(%s) returned theme named %s", name, th.Name)
		}

		g := newTestGame()
		// Place all item types on board to exercise glyph styles
		g.FoodItems = []engine.FoodItem{
			{Point: engine.Point{X: 1, Y: 1}, Type: engine.ItemFood},
			{Point: engine.Point{X: 2, Y: 1}, Type: engine.ItemSlowMo},
			{Point: engine.Point{X: 3, Y: 1}, Type: engine.ItemGhost},
			{Point: engine.Point{X: 4, Y: 1}, Type: engine.ItemMultiplier},
			{Point: engine.Point{X: 5, Y: 1}, Type: engine.ItemShrink},
			{Point: engine.Point{X: 6, Y: 1}, Type: engine.ItemPoison},
		}

		model := NewUIModel(g, cfg, stats, name, false)
		view := model.View()

		if len(view) < 100 {
			t.Errorf("Theme '%s' view rendered empty", name)
		}
	}
}

func TestExtreme_ScreenReaderMode(t *testing.T) {
	cfg := storage.DefaultConfig()
	cfg.ScreenReader = true
	stats := storage.NewStats()
	g := newTestGame()

	model := NewUIModel(g, cfg, stats, "dark", false)
	view := model.View()

	if !strings.Contains(view, "Snake CLI. Mode:") || !strings.Contains(view, "Score: 0") {
		t.Errorf("Screen reader view invalid:\n%s", view)
	}
}

func TestExtreme_ResponsiveSmallTerminal(t *testing.T) {
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()
	g := newTestGame()

	model := NewUIModel(g, cfg, stats, "dark", true)
	// Small terminal pane (e.g. split pane 60x18)
	model.TermWidth = 60
	model.TermHeight = 18

	view := model.View()
	if !strings.Contains(view, "MAIN MENU & SETTINGS") {
		t.Errorf("Small terminal view missing menu header")
	}

	// Update with resize
	m2, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 15})
	view2 := m2.View()
	if len(view2) < 50 {
		t.Errorf("Very narrow terminal view rendered empty")
	}
}
