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
