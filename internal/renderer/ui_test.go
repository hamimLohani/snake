package renderer

import (
	"strings"
	"testing"

	"snake/internal/engine"
	"snake/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func newTestGame() *engine.Game {
	return engine.New(engine.Config{
		Width:           20,
		Height:          10,
		Mode:            engine.ModeClassic,
		Seed:            42,
		InitialSnakeLen: 3,
		BaseTickRateMs:  110,
	})
}

func TestNewUIModel_StartInMenu(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	// Default start in menu with no save game
	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false

	if !model.InMenu {
		t.Errorf("Expected InMenu to be true")
	}
	if !model.Paused {
		t.Errorf("Expected Paused to be true")
	}
	if !model.IsInitialStart {
		t.Errorf("Expected IsInitialStart to be true")
	}
	if model.MenuIndex != 0 {
		t.Errorf("Expected MenuIndex to be 0, got %d", model.MenuIndex)
	}

	view := model.View()
	if !strings.Contains(view, "MAIN MENU & SETTINGS") {
		t.Errorf("Expected view to contain 'MAIN MENU & SETTINGS', got:\n%s", view)
	}
	if !strings.Contains(view, "START GAME") {
		t.Errorf("Expected view to contain 'START GAME'")
	}
	if !strings.Contains(view, "Mode:") {
		t.Errorf("Expected view to contain 'Mode:'")
	}
	if !strings.Contains(view, "Board Size:") {
		t.Errorf("Expected view to contain 'Board Size:'")
	}
}

func TestNewUIModel_StartInMenu_WithSavedGame(t *testing.T) {
	game := newTestGame()
	game.Score = 150
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = true

	view := model.View()
	if !strings.Contains(view, "RESUME SAVED GAME (Score: 150)") {
		t.Errorf("Expected view to offer resuming saved game, got:\n%s", view)
	}
	if !strings.Contains(view, "Start New Game") {
		t.Errorf("Expected view to offer Start New Game")
	}
}

func TestNewUIModel_DirectStart(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	// Direct start without menu
	model := NewUIModel(game, cfg, stats, "dark", false)

	if model.InMenu {
		t.Errorf("Expected InMenu to be false")
	}
	if model.Paused {
		t.Errorf("Expected Paused to be false")
	}
	if model.IsInitialStart {
		t.Errorf("Expected IsInitialStart to be false")
	}

	view := model.View()
	if strings.Contains(view, "MAIN MENU & SETTINGS") {
		t.Errorf("Expected view NOT to contain 'MAIN MENU & SETTINGS' during direct start")
	}
	if !strings.Contains(view, "SCORE:") {
		t.Errorf("Expected view to contain HUD 'SCORE:'")
	}
}

func TestUIModel_MenuNavigationAndAdjustment(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false

	// Press Down to navigate to Mode
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := newModel.(UIModel)
	if m.MenuIndex != 1 {
		t.Fatalf("Expected MenuIndex 1 after Down key, got %d", m.MenuIndex)
	}

	// Press Right to change Mode (from Classic to Time Attack)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = newModel.(UIModel)
	if m.SelectedMode != 1 {
		t.Errorf("Expected SelectedMode 1, got %d", m.SelectedMode)
	}

	// Press Down to Difficulty
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newModel.(UIModel)
	if m.MenuIndex != 2 {
		t.Fatalf("Expected MenuIndex 2 after Down key, got %d", m.MenuIndex)
	}

	// Press Right to change Difficulty (from Normal=1 to Hard=2)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = newModel.(UIModel)
	if m.SelectedDifficulty != 2 {
		t.Errorf("Expected SelectedDifficulty 2, got %d", m.SelectedDifficulty)
	}

	// Press Down to Board Size
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newModel.(UIModel)
	if m.MenuIndex != 3 {
		t.Fatalf("Expected MenuIndex 3 after Down key, got %d", m.MenuIndex)
	}

	// Cycle Board Size to LARGE (index 3)
	m.SelectedSize = 3 // BoardSizePresets[3] = LARGE (36x20)

	// Press Enter to start game (or Space)
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = newModel.(UIModel)

	if m.InMenu {
		t.Errorf("Expected InMenu to be false after starting")
	}
	if m.Paused {
		t.Errorf("Expected Paused to be false after starting")
	}
	if m.IsInitialStart {
		t.Errorf("Expected IsInitialStart to be false after starting")
	}
	if m.Game.Config.Mode != engine.ModeTimeAttack {
		t.Errorf("Expected game mode to be Time Attack, got %s", m.Game.Config.Mode)
	}
	if m.Game.Config.BaseTickRateMs != 75 { // Hard = 75ms
		t.Errorf("Expected tick speed to be 75ms for Hard, got %d", m.Game.Config.BaseTickRateMs)
	}
	if m.Game.Config.Width != 36 || m.Game.Config.Height != 20 {
		t.Errorf("Expected game grid to be 36x20, got %dx%d", m.Game.Config.Width, m.Game.Config.Height)
	}
	if cmd == nil {
		t.Errorf("Expected ticker command to be returned on start")
	}
}

func TestUIModel_MenuEnterStart(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false

	// Cursor is at index 0 [ ▶ START GAME ]. Pressing Enter starts the game!
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := newModel.(UIModel)

	if m.InMenu {
		t.Errorf("Expected InMenu to be false after Enter on START GAME")
	}
	if m.Paused {
		t.Errorf("Expected Paused to be false after Enter on START GAME")
	}
	if cmd == nil {
		t.Errorf("Expected ticker command after Enter on START GAME")
	}
}

func TestUIModel_QuitFromMenu(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false

	// Press 'q' to quit
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Errorf("Expected quit command on 'q' press")
	}
}

func TestUIModel_SampleBoardPreview_Rendering(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false
	view := model.View()

	if !strings.Contains(view, "SAMPLE BOARD PREVIEW") {
		t.Errorf("Expected view to contain 'SAMPLE BOARD PREVIEW', got:\n%s", view)
	}
	if !strings.Contains(view, "CLASSIC (24×20)") {
		t.Errorf("Expected view to contain 'CLASSIC (24×20)' in sample preview")
	}
	if !strings.Contains(view, "DARK") {
		t.Errorf("Expected view to contain 'DARK' theme tag in sample preview")
	}
	if !strings.Contains(view, "Head:") || !strings.Contains(view, "Body:") || !strings.Contains(view, "Food:") {
		t.Errorf("Expected view to contain glyph legend (Head, Body, Food)")
	}
}

func TestUIModel_SampleBoardPreview_ThemeChanges(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false

	// Navigate to Theme (Row 5 in start menu with no save: 0:Start, 1:Mode, 2:Diff, 3:Size, 4:Obs, 5:Theme)
	for i := 0; i < 5; i++ {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = newModel.(UIModel)
	}
	if model.MenuIndex != 5 {
		t.Fatalf("Expected MenuIndex 5, got %d", model.MenuIndex)
	}

	// themes: 0: dark, 1: light, 2: high-contrast, 3: colorblind, 4: corporate
	expectedThemes := []string{"LIGHT", "HIGH-CONTRAST", "COLORBLIND", "CORPORATE", "DARK"}
	for _, expected := range expectedThemes {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
		model = newModel.(UIModel)
		view := model.View()

		if !strings.Contains(view, expected) {
			t.Errorf("Expected preview to contain theme tag '%s'", expected)
		}
	}
}

func TestUIModel_SampleBoardPreview_SizeChanges(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false

	// Navigate to Board Size (Row 3: 0:Start, 1:Mode, 2:Diff, 3:Size)
	for i := 0; i < 3; i++ {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = newModel.(UIModel)
	}
	if model.MenuIndex != 3 {
		t.Fatalf("Expected MenuIndex 3, got %d", model.MenuIndex)
	}

	// Board sizes: initial is Classic (idx 2). Right cycles to Large (idx 3).
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(UIModel)
	view := model.View()
	if !strings.Contains(view, "LARGE (36×20)") {
		t.Errorf("Expected preview to show LARGE (36×20)")
	}

	// Cycle to Small (idx 0)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(UIModel)
	view = model.View()
	if !strings.Contains(view, "SMALL (20×12)") {
		t.Errorf("Expected preview to show SMALL (20×12)")
	}

	// Cycle to Medium (idx 1)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(UIModel)
	view = model.View()
	if !strings.Contains(view, "MEDIUM (28×16)") {
		t.Errorf("Expected preview to show MEDIUM (28×16)")
	}
}

func TestUIModel_SampleBoardPreview_ObstaclesToggle(t *testing.T) {
	game := newTestGame()
	cfg := storage.DefaultConfig()
	stats := storage.NewStats()

	model := NewUIModel(game, cfg, stats, "dark", true)
	model.HasSaveGame = false

	// Navigate to Obstacles (Row 4: 0:Start, 1:Mode, 2:Diff, 3:Size, 4:Obs)
	for i := 0; i < 4; i++ {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = newModel.(UIModel)
	}
	if model.MenuIndex != 4 {
		t.Fatalf("Expected MenuIndex 4, got %d", model.MenuIndex)
	}

	// Initially obstacles are false
	viewBefore := model.View()
	if strings.Contains(viewBefore, "Obs:██") {
		t.Errorf("Expected preview NOT to contain 'Obs:██' when obstacles are OFF")
	}

	// Press Right to toggle Obstacles ON
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(UIModel)
	viewAfter := model.View()
	if !strings.Contains(viewAfter, "Obs:██") {
		t.Errorf("Expected preview to contain 'Obs:██' when obstacles are ON")
	}
}
