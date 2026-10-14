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
	Snake            []engine.Point   `json:"snake"`
	Direction        engine.Direction `json:"direction"`
	Food             engine.Point     `json:"food"`
	FoodType         engine.ItemType  `json:"food_type"`
	Score            int              `json:"score"`
	TickCount        uint64           `json:"tick_count"`
	GrowthPending    int              `json:"growth_pending"`
	Bounds           engine.Bounds    `json:"bounds"`
	Mode             engine.Mode      `json:"mode"`
	Seed             int64            `json:"seed"`
	TimeRemainingSec int              `json:"time_remaining_sec"`
	SavedAt          time.Time        `json:"saved_at"`
}

// SavePath returns ~/.config/snake/save.json.
func SavePath() string {
	return filepath.Join(ConfigDir(), "save.json")
}

// HasSaveGame returns true if a valid savegame exists.
func HasSaveGame() bool {
	_, err := os.Stat(SavePath())
	return err == nil
}

// SaveActiveGame writes the engine state to disk.
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
		Score:            g.Score,
		TickCount:        g.TickCount,
		GrowthPending:    g.GrowthPending,
		Bounds:           g.Bounds,
		Mode:             g.Config.Mode,
		Seed:             g.Config.Seed,
		TimeRemainingSec: g.TimeRemainingSec,
		SavedAt:          time.Now(),
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(SavePath(), data, 0644)
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

	cfg := engine.Config{
		Width:            state.Bounds.MaxX,
		Height:           state.Bounds.MaxY,
		Mode:             state.Mode,
		Seed:             state.Seed,
		InitialSnakeLen:  len(state.Snake),
		AdaptiveSpeed:    true,
		BaseTickRateMs:   110,
	}

	g := engine.New(cfg)
	g.Snake = state.Snake
	g.Direction = state.Direction
	g.Food = state.Food
	g.FoodType = state.FoodType
	g.Score = state.Score
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
