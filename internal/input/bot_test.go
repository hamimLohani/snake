package input

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"snake/internal/engine"
	"snake/internal/storage"
)

func TestBotCanScorePoints(t *testing.T) {
	cfg := engine.DefaultConfig()
	cfg.Width = 15
	cfg.Height = 15
	cfg.Seed = 1337

	game := engine.New(cfg)
	solver := NewBotSolver()

	// Run bot for 200 ticks
	for i := 0; i < 200; i++ {
		dir := solver.NextMove(game)
		if dir != engine.DirNone {
			game.QueueDirection(dir)
		}
		game.Tick()
		if game.GameOver {
			break
		}
	}

	if game.Score < 20 {
		t.Errorf("Expected bot to achieve score >= 20 in 200 ticks, got %d", game.Score)
	}
}

func TestReplayDeterminism(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_replay_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	seed := int64(987654321)
	cfg := engine.DefaultConfig()
	cfg.Width = 16
	cfg.Height = 16
	cfg.Seed = seed

	game1 := engine.New(cfg)
	solver := NewBotSolver()

	var recordedSteps []storage.ReplayStep

	// Run 100 ticks and record every direction change
	for i := 0; i < 100; i++ {
		dir := solver.NextMove(game1)
		if dir != engine.DirNone {
			if game1.QueueDirection(dir) {
				recordedSteps = append(recordedSteps, storage.ReplayStep{
					Tick:      game1.TickCount,
					Direction: dir,
				})
			}
		}
		game1.Tick()
	}

	rep := storage.ReplayFile{
		Version:     "1.0.0",
		Mode:        cfg.Mode,
		Seed:        seed,
		Width:       cfg.Width,
		Height:      cfg.Height,
		FinalScore:  game1.Score,
		FinalLength: len(game1.Snake),
		TotalTicks:  game1.TickCount,
		Steps:       recordedSteps,
	}

	replayPath := filepath.Join(tmpDir, "test.snk")
	if err := storage.SaveReplay(replayPath, rep); err != nil {
		t.Fatalf("Failed to save replay: %v", err)
	}

	loadedRep, err := storage.LoadReplay(replayPath)
	if err != nil {
		t.Fatalf("Failed to load replay: %v", err)
	}

	// Now execute loaded replay on a fresh engine
	game2 := engine.New(cfg)
	stepMap := make(map[uint64]engine.Direction, len(loadedRep.Steps))
	for _, s := range loadedRep.Steps {
		stepMap[s.Tick] = s.Direction
	}

	for i := uint64(0); i < loadedRep.TotalTicks; i++ {
		if d, ok := stepMap[game2.TickCount]; ok {
			game2.QueueDirection(d)
		}
		game2.Tick()
	}

	// Assert 100% determinism
	if game1.Score != game2.Score {
		t.Fatalf("Replay determinism mismatch on score: %d vs %d", game1.Score, game2.Score)
	}
	if !reflect.DeepEqual(game1.Snake, game2.Snake) {
		t.Fatalf("Replay determinism mismatch on snake body.\ngame1: %+v\ngame2: %+v", game1.Snake, game2.Snake)
	}
	if game1.Food != game2.Food {
		t.Fatalf("Replay determinism mismatch on food position: %+v vs %+v", game1.Food, game2.Food)
	}
}

func TestBotAvoidsPoisonHazard(t *testing.T) {
	cfg := engine.DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.Seed = 5555

	game := engine.New(cfg)
	solver := NewBotSolver()

	// Snake head at (5, 5), body at (4, 5), (3, 5), moving Right
	head := game.Snake[0]
	// Place poison directly in front of head at (6, 5)
	game.BonusPoint = engine.Point{X: head.X + 1, Y: head.Y}
	game.BonusType = engine.ItemPoison
	game.BonusActive = true

	// Place safe food above at (5, 4)
	game.Food = engine.Point{X: head.X, Y: head.Y - 1}

	nextMove := solver.NextMove(game)
	if nextMove == engine.DirRight {
		t.Errorf("Bot should NOT move Right into poison hazard at (6, 5)")
	}
	if nextMove != engine.DirUp {
		t.Logf("Bot navigated safely with move: %v", nextMove)
	}
}
