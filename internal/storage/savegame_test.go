package storage

import (
	"os"
	"testing"

	"snake/internal/engine"
)

func TestSaveAndLoadBonusItem(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_save_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Override config directory for test
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := engine.DefaultConfig()
	cfg.Width = 24
	cfg.Height = 18
	cfg.Seed = 7777

	g := engine.New(cfg)
	g.BonusPoint = engine.Point{X: 12, Y: 9}
	g.BonusType = engine.ItemPoison
	g.BonusActive = true
	g.BonusExpires = 150

	if err := SaveActiveGame(g); err != nil {
		t.Fatalf("Failed to save active game: %v", err)
	}

	loaded, err := LoadSavedGame()
	if err != nil {
		t.Fatalf("Failed to load saved game: %v", err)
	}

	if loaded.BonusPoint != g.BonusPoint {
		t.Errorf("Expected BonusPoint %+v, got %+v", g.BonusPoint, loaded.BonusPoint)
	}
	if loaded.BonusType != g.BonusType {
		t.Errorf("Expected BonusType %v, got %v", g.BonusType, loaded.BonusType)
	}
	if loaded.BonusActive != g.BonusActive {
		t.Errorf("Expected BonusActive %v, got %v", g.BonusActive, loaded.BonusActive)
	}
	if loaded.BonusExpires != g.BonusExpires {
		t.Errorf("Expected BonusExpires %d, got %d", g.BonusExpires, loaded.BonusExpires)
	}

	if err := ClearSaveGame(); err != nil {
		t.Fatalf("Failed to clear save game: %v", err)
	}
	if HasSaveGame() {
		t.Errorf("Expected save game to be cleared")
	}
}
