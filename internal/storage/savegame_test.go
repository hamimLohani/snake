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
	if len(loaded.FoodItems) != len(g.FoodItems) {
		t.Errorf("Expected FoodItems count %d, got %d", len(g.FoodItems), len(loaded.FoodItems))
	}

	if err := ClearSaveGame(); err != nil {
		t.Fatalf("Failed to clear save game: %v", err)
	}
	if HasSaveGame() {
		t.Errorf("Expected save game to be cleared")
	}
}

func TestSaveAndLoadMultipleFoodItems(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_save_multi_food_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := engine.DefaultConfig()
	g := engine.New(cfg)
	g.FoodItems = []engine.FoodItem{
		{Point: engine.Point{X: 3, Y: 4}, Type: engine.ItemFood},
		{Point: engine.Point{X: 7, Y: 8}, Type: engine.ItemFood},
		{Point: engine.Point{X: 11, Y: 12}, Type: engine.ItemSlowMo, ExpiresAtTick: 120},
		{Point: engine.Point{X: 15, Y: 16}, Type: engine.ItemPoison, ExpiresAtTick: 120},
	}

	if err := SaveActiveGame(g); err != nil {
		t.Fatalf("Failed to save active game with multiple foods: %v", err)
	}

	loaded, err := LoadSavedGame()
	if err != nil {
		t.Fatalf("Failed to load saved game: %v", err)
	}

	if len(loaded.FoodItems) != 4 {
		t.Fatalf("Expected 4 FoodItems loaded, got %d", len(loaded.FoodItems))
	}
	if loaded.FoodItems[0].Type != engine.ItemFood || loaded.FoodItems[3].Type != engine.ItemPoison {
		t.Errorf("Mismatch in loaded FoodItems types: %+v", loaded.FoodItems)
	}
}
