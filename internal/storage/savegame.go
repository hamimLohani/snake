package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"snake/internal/engine"
)

// SaveState stores the serializable snapshot of an in-progress game.
type SaveState struct {
	Snake            []engine.Point        `json:"snake"`
	Direction        engine.Direction      `json:"direction"`
	Food             engine.Point          `json:"food"`
	FoodType         engine.ItemType       `json:"food_type"`
	BonusPoint       engine.Point          `json:"bonus_point,omitempty"`
	BonusType        engine.ItemType       `json:"bonus_type,omitempty"`
	BonusActive      bool                  `json:"bonus_active,omitempty"`
	BonusExpires     uint64                `json:"bonus_expires,omitempty"`
	FoodItems        []engine.FoodItem     `json:"food_items,omitempty"`
	Score            int                   `json:"score"`
	ScoreMultiplier  int                   `json:"score_multiplier,omitempty"`
	TickCount        uint64                `json:"tick_count"`
	GrowthPending    int                   `json:"growth_pending"`
	Bounds           engine.Bounds         `json:"bounds"`
	Mode             engine.Mode           `json:"mode"`
	Seed             int64                 `json:"seed"`
	TimeRemainingSec int                   `json:"time_remaining_sec"`
	SavedAt          time.Time             `json:"saved_at"`
	BaseTickRateMs   int                   `json:"base_tick_rate_ms,omitempty"`
	WithObstacles    bool                  `json:"with_obstacles,omitempty"`
	ObstacleDensity  float64               `json:"obstacle_density,omitempty"`
	ActivePowerUp    *engine.ActivePowerUp `json:"active_power_up,omitempty"`
	ConfigWidth      int                   `json:"config_width,omitempty"`
	ConfigHeight     int                   `json:"config_height,omitempty"`
}

// SavePath returns ~/.config/snake/save.json.
func SavePath() string {
	return filepath.Join(ConfigDir(), "save.json")
}

// HasSaveGame returns true if a valid, non-empty savegame exists.
func HasSaveGame() bool {
	info, err := os.Stat(SavePath())
	if err != nil || info.Size() == 0 {
		return false
	}
	return true
}

// SaveActiveGame writes the engine state to disk atomically.
func SaveActiveGame(g *engine.Game) error {
	if g == nil || g.GameOver {
		// Clean up existing save if game is already over
		_ = ClearSaveGame()
		return nil
	}

	state := SaveState{
		Snake:            g.Snake,
		Direction:        g.Direction,
		Food:             g.Food,
		FoodType:         g.FoodType,
		BonusPoint:       g.BonusPoint,
		BonusType:        g.BonusType,
		BonusActive:      g.BonusActive,
		BonusExpires:     g.BonusExpires,
		FoodItems:        g.FoodItems,
		Score:            g.Score,
		ScoreMultiplier:  g.ScoreMultiplier,
		TickCount:        g.TickCount,
		GrowthPending:    g.GrowthPending,
		Bounds:           g.Bounds,
		Mode:             g.Config.Mode,
		Seed:             g.Config.Seed,
		TimeRemainingSec: g.TimeRemainingSec,
		SavedAt:          time.Now(),
		BaseTickRateMs:   g.Config.BaseTickRateMs,
		WithObstacles:    g.Config.WithObstacles,
		ObstacleDensity:  g.Config.ObstacleDensity,
		ActivePowerUp:    g.ActivePowerUp,
		ConfigWidth:      g.Config.Width,
		ConfigHeight:     g.Config.Height,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Atomic write using temp file and rename
	tmpFile := filepath.Join(dir, "save.tmp")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, SavePath())
}

// LoadSavedGame restores the saved game state into an engine.
func LoadSavedGame() (*engine.Game, error) {
	data, err := os.ReadFile(SavePath())
	if err != nil {
		return nil, err
	}

	var state SaveState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	width := state.ConfigWidth
	if width <= 0 {
		width = state.Bounds.MaxX
	}
	height := state.ConfigHeight
	if height <= 0 {
		height = state.Bounds.MaxY
	}
	baseSpeed := state.BaseTickRateMs
	if baseSpeed <= 0 {
		baseSpeed = 110
	}

	cfg := engine.Config{
		Width:           width,
		Height:          height,
		Mode:            state.Mode,
		Seed:            state.Seed,
		InitialSnakeLen: len(state.Snake),
		AdaptiveSpeed:   true,
		BaseTickRateMs:  baseSpeed,
		WithObstacles:   state.WithObstacles,
		ObstacleDensity: state.ObstacleDensity,
	}

	g := engine.New(cfg)
	g.Snake = state.Snake
	g.Direction = state.Direction
	g.Food = state.Food
	g.FoodType = state.FoodType
	g.BonusPoint = state.BonusPoint
	g.BonusType = state.BonusType
	g.BonusActive = state.BonusActive
	g.BonusExpires = state.BonusExpires
	if len(state.FoodItems) > 0 {
		g.FoodItems = state.FoodItems
	}
	g.Score = state.Score
	if state.ScoreMultiplier > 0 {
		g.ScoreMultiplier = state.ScoreMultiplier
	} else {
		g.ScoreMultiplier = 1
	}
	g.ActivePowerUp = state.ActivePowerUp
	g.TickCount = state.TickCount
	g.GrowthPending = state.GrowthPending
	g.Bounds = state.Bounds
	g.TimeRemainingSec = state.TimeRemainingSec

	return g, nil
}

// ClearSaveGame removes the save file.
func ClearSaveGame() error {
	err := os.Remove(SavePath())
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
