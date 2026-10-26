package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"snake/internal/engine"
)

func TestExtreme_SaveAndLoadFullFidelity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_extreme_storage")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := engine.DefaultConfig()
	cfg.Width = 28
	cfg.Height = 16
	cfg.Mode = engine.ModeSurvival
	cfg.BaseTickRateMs = 75
	cfg.WithObstacles = true
	cfg.ObstacleDensity = 0.05
	cfg.Seed = 1234567

	g := engine.New(cfg)
	g.Score = 420
	g.ScoreMultiplier = 3
	g.TickCount = 250
	g.GrowthPending = 2
	g.Bounds = engine.Bounds{MinX: 3, MaxX: 25, MinY: 2, MaxY: 14} // Contracted in survival
	g.ActivePowerUp = &engine.ActivePowerUp{
		Type:          engine.ItemMultiplier,
		ExpiresAtTick: 300,
	}

	if err := SaveActiveGame(g); err != nil {
		t.Fatalf("Failed to save active game: %v", err)
	}

	if !HasSaveGame() {
		t.Errorf("Expected HasSaveGame to return true")
	}

	loaded, err := LoadSavedGame()
	if err != nil {
		t.Fatalf("Failed to load saved game: %v", err)
	}

	if loaded.Score != 420 {
		t.Errorf("Expected Score 420, got %d", loaded.Score)
	}
	if loaded.ScoreMultiplier != 3 {
		t.Errorf("Expected ScoreMultiplier 3, got %d", loaded.ScoreMultiplier)
	}
	if loaded.ActivePowerUp == nil || loaded.ActivePowerUp.Type != engine.ItemMultiplier {
		t.Errorf("Expected active powerup ItemMultiplier, got %+v", loaded.ActivePowerUp)
	}
	if loaded.Config.BaseTickRateMs != 75 {
		t.Errorf("Expected BaseTickRateMs 75, got %d", loaded.Config.BaseTickRateMs)
	}
	if !loaded.Config.WithObstacles {
		t.Errorf("Expected WithObstacles to be true")
	}
	if loaded.Bounds != g.Bounds {
		t.Errorf("Expected Bounds %+v, got %+v", g.Bounds, loaded.Bounds)
	}
	if loaded.Config.Width != 28 || loaded.Config.Height != 16 {
		t.Errorf("Expected Config dimensions 28x16, got %dx%d", loaded.Config.Width, loaded.Config.Height)
	}
}

func TestExtreme_CorruptedSaveHandling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_corrupt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	saveDir := filepath.Join(tmpDir, ".config", "snake")
	_ = os.MkdirAll(saveDir, 0755)
	savePath := filepath.Join(saveDir, "save.json")

	// 1. Zero-byte save file
	_ = os.WriteFile(savePath, []byte{}, 0644)
	if HasSaveGame() {
		t.Errorf("HasSaveGame should be false for 0-byte file")
	}

	// 2. Corrupted JSON file
	_ = os.WriteFile(savePath, []byte("NOT_JSON_DATA{{{"), 0644)
	_, err = LoadSavedGame()
	if err == nil {
		t.Errorf("Expected error loading corrupted JSON save file, got nil")
	}
}

func TestExtreme_StatsFoodEatenUnderflow(t *testing.T) {
	st := NewStats()

	// Short snakes (< 3 segments, e.g. after shrink powerup) must not cause underflow
	st.RecordGame(ScoreEntry{
		Date:        time.Now(),
		Mode:        "classic",
		Score:       10,
		Length:      2,
		DurationSec: 30,
	})

	if st.TotalFoodEaten < 0 {
		t.Errorf("TotalFoodEaten must not be negative, got %d", st.TotalFoodEaten)
	}
	if st.TotalFoodEaten != 0 {
		t.Errorf("TotalFoodEaten for length 2 should add 0 food, got %d", st.TotalFoodEaten)
	}

	// Normal game adds Length - 3
	st.RecordGame(ScoreEntry{
		Date:        time.Now(),
		Mode:        "classic",
		Score:       50,
		Length:      8,
		DurationSec: 45,
	})
	if st.TotalFoodEaten != 5 {
		t.Errorf("TotalFoodEaten should be 5, got %d", st.TotalFoodEaten)
	}
}

func TestExtreme_StatsAchievementsAndCSVExport(t *testing.T) {
	st := NewStats()

	// Score 150 points with length 26 to unlock multiple achievements
	_, unlocked1 := st.RecordGame(ScoreEntry{
		Date:        time.Now(),
		Mode:        "classic",
		Score:       150,
		Length:      26,
		DurationSec: 350,
	})

	if len(unlocked1) == 0 {
		t.Errorf("Expected achievements to unlock")
	}

	// Repeating game should NOT duplicate achievements
	_, unlocked2 := st.RecordGame(ScoreEntry{
		Date:        time.Now(),
		Mode:        "classic",
		Score:       150,
		Length:      26,
		DurationSec: 350,
	})

	if len(unlocked2) != 0 {
		t.Errorf("Achievements must not unlock multiple times: %+v", unlocked2)
	}

	// CSV Export
	csvData, err := st.ExportCSV()
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}
	if !strings.Contains(csvData, "Rank,Date,Mode,Score,Length,DurationSec,Seed") {
		t.Errorf("CSV missing header: %s", csvData)
	}
	if !strings.Contains(csvData, "classic,150,26,350") {
		t.Errorf("CSV missing expected row data: %s", csvData)
	}

	// JSON Export
	jsonData, err := st.ExportJSON()
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}
	if !strings.Contains(jsonData, `"total_games": 2`) {
		t.Errorf("JSON missing total_games: %s", jsonData)
	}
}

func TestExtreme_ReplayGzipCompression(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_replay_extreme")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	rep := ReplayFile{
		Version:       "1.0.0",
		Date:          time.Now(),
		Mode:          engine.ModeClassic,
		Seed:          99999,
		Width:         24,
		Height:        20,
		WithObstacles: false,
		FinalScore:    280,
		FinalLength:   31,
		TotalTicks:    400,
		Steps: []ReplayStep{
			{Tick: 10, Direction: engine.DirDown},
			{Tick: 25, Direction: engine.DirLeft},
			{Tick: 40, Direction: engine.DirUp},
			{Tick: 60, Direction: engine.DirRight},
		},
	}

	repPath := filepath.Join(tmpDir, "game.snk")
	if err := SaveReplay(repPath, rep); err != nil {
		t.Fatalf("Failed to save replay: %v", err)
	}

	loaded, err := LoadReplay(repPath)
	if err != nil {
		t.Fatalf("Failed to load replay: %v", err)
	}

	if loaded.FinalScore != rep.FinalScore || len(loaded.Steps) != len(rep.Steps) {
		t.Errorf("Loaded replay does not match: expected %+v, got %+v", rep, loaded)
	}
}

func TestExtreme_ConfigDefaultsAndPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_config_extreme")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// First load creates default config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Theme != "dark" || cfg.GridWidth != 24 || cfg.GridHeight != 20 {
		t.Errorf("Unexpected default config values: %+v", cfg)
	}

	// Modify and save
	cfg.Theme = "corporate"
	cfg.GridWidth = 36
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload
	reloaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}
	if reloaded.Theme != "corporate" || reloaded.GridWidth != 36 {
		t.Errorf("Config modification did not persist: %+v", reloaded)
	}
}
