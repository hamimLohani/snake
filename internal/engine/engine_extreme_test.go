package engine

import (
	"testing"
)

// TestExtreme_BoardDimensions tests various extreme board dimensions from small to large.
func TestExtreme_BoardDimensions(t *testing.T) {
	sizes := []struct {
		w, h int
	}{
		{5, 5},
		{6, 6},
		{7, 10},
		{10, 7},
		{12, 12},
		{50, 10},
		{10, 50},
		{80, 40},
	}

	for _, s := range sizes {
		cfg := DefaultConfig()
		cfg.Width = s.w
		cfg.Height = s.h
		cfg.InitialSnakeLen = 3
		cfg.Seed = 12345

		g := New(cfg)
		if g.Bounds.MaxX != s.w || g.Bounds.MaxY != s.h {
			t.Errorf("Size %dx%d: expected bounds (%d, %d), got (%d, %d)",
				s.w, s.h, s.w, s.h, g.Bounds.MaxX, g.Bounds.MaxY)
		}
		if len(g.Snake) != 3 {
			t.Errorf("Size %dx%d: expected snake len 3, got %d", s.w, s.h, len(g.Snake))
		}

		// Run 20 ticks
		for i := 0; i < 20; i++ {
			g.Tick()
			if g.GameOver {
				break
			}
		}
	}
}

// TestExtreme_FullBoardVictory verifies clean Game Over when all open cells are filled.
func TestExtreme_FullBoardVictory(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 6
	cfg.Height = 6
	cfg.InitialSnakeLen = 3
	cfg.Seed = 999

	g := New(cfg)
	// Artificially fill entire board except head
	g.Snake = make([]Point, 35)
	idx := 0
	for y := 0; y < 6; y++ {
		for x := 0; x < 6; x++ {
			if x == 0 && y == 0 {
				continue // Open cell
			}
			g.Snake[idx] = Point{X: x, Y: y}
			idx++
		}
	}
	g.Snake = append([]Point{{X: 0, Y: 0}}, g.Snake...)
	g.Obstacles = make(map[Point]bool)
	g.FoodItems = nil

	g.spawnNeededItems()

	if !g.GameOver {
		t.Errorf("Expected Game Over when board is completely filled")
	}
	if g.GameOverReason != "Victory! Board fully cleared!" {
		t.Errorf("Expected reason 'Victory! Board fully cleared!', got '%s'", g.GameOverReason)
	}
}

// TestExtreme_ZenModeInfiniteGrowth verifies Zen mode continues spawning food even if snake length exceeds cell count.
func TestExtreme_ZenModeInfiniteGrowth(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.Mode = ModeZen
	cfg.Seed = 777

	g := New(cfg)
	// Create snake with 150 segments (overlapping on 10x10 board with 100 cells)
	g.Snake = make([]Point, 150)
	for i := 0; i < 150; i++ {
		// All coiled into 4 cells: (0,0), (0,1), (1,1), (1,0)
		g.Snake[i] = Point{X: i % 2, Y: (i / 2) % 2}
	}

	g.FoodItems = nil
	g.spawnNeededItems()

	if g.GameOver {
		t.Fatalf("Zen mode must NEVER trigger Game Over, got reason: %s", g.GameOverReason)
	}
	if len(g.FoodItems) == 0 {
		t.Errorf("Food must spawn in Zen mode when open cells exist despite snake length > cell count")
	}
}

// TestExtreme_SurvivalModeContraction tests independent axis shrinking and food pruning.
func TestExtreme_SurvivalModeContraction(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 20
	cfg.Height = 10
	cfg.Mode = ModeSurvival
	cfg.Seed = 42

	g := New(cfg)
	// Place food at the outer perimeter (0, 0)
	g.FoodItems = []FoodItem{
		{Point: Point{X: 0, Y: 0}, Type: ItemFood},
		{Point: Point{X: 10, Y: 5}, Type: ItemPoison, ExpiresAtTick: 500},
	}
	g.syncFoodFields()

	// Safely orbit snake around center for 50 ticks
	for i := 0; i < 50; i++ {
		head := g.Snake[0]
		if g.Direction == DirRight && head.X >= 15 {
			g.QueueDirection(DirDown)
		} else if g.Direction == DirDown && head.Y >= 7 {
			g.QueueDirection(DirLeft)
		} else if g.Direction == DirLeft && head.X <= 5 {
			g.QueueDirection(DirUp)
		} else if g.Direction == DirUp && head.Y <= 3 {
			g.QueueDirection(DirRight)
		}
		g.Tick()
		if g.GameOver {
			t.Fatalf("Unexpected death at tick %d: %s", g.TickCount, g.GameOverReason)
		}
	}

	// MinX should be 1, MaxX should be 19; MinY should be 1, MaxY should be 9
	if g.Bounds.MinX != 1 || g.Bounds.MaxX != 19 || g.Bounds.MinY != 1 || g.Bounds.MaxY != 9 {
		t.Errorf("Expected bounds (1, 19, 1, 9), got %+v", g.Bounds)
	}

	// The food at (0, 0) was outside bounds; verify all FoodItems are strictly within new bounds
	for _, f := range g.FoodItems {
		if f.Point.X < g.Bounds.MinX || f.Point.X >= g.Bounds.MaxX ||
			f.Point.Y < g.Bounds.MinY || f.Point.Y >= g.Bounds.MaxY {
			t.Errorf("Food item %+v is outside shrinking bounds %+v", f, g.Bounds)
		}
	}

	// Ensure there is still permanent food reachable
	hasFood := false
	for _, f := range g.FoodItems {
		if f.Type == ItemFood {
			hasFood = true
			break
		}
	}
	if !hasFood {
		t.Errorf("Permanent food should be replenished after shrinking")
	}
}

// TestExtreme_ActivePowerUpOverwrite tests that multiplier resets when another power-up is eaten.
func TestExtreme_ActivePowerUpOverwrite(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Mode = ModeZen
	cfg.BaseTickRateMs = 100
	g := New(cfg)

	// 1. Eat Multiplier
	g.applyItemEffect(FoodItem{Point: Point{X: 1, Y: 1}, Type: ItemMultiplier})
	if g.ScoreMultiplier != 3 {
		t.Fatalf("Expected Multiplier 3, got %d", g.ScoreMultiplier)
	}

	// 2. While Multiplier is active, eat Poison
	g.applyItemEffect(FoodItem{Point: Point{X: 2, Y: 2}, Type: ItemPoison})
	if g.ScoreMultiplier != 1 {
		t.Errorf("Expected ScoreMultiplier reset to 1 when Poison overwrites Multiplier, got %d", g.ScoreMultiplier)
	}
	if g.ActivePowerUp == nil || g.ActivePowerUp.Type != ItemPoison {
		t.Fatalf("Expected ActivePowerUp to be ItemPoison")
	}

	// 3. Advance until Poison wears off
	for g.ActivePowerUp != nil {
		g.Tick()
	}
	if g.ScoreMultiplier != 1 {
		t.Errorf("ScoreMultiplier must remain 1 after poison wears off, got %d", g.ScoreMultiplier)
	}
}

// TestExtreme_DirectionQueueRapid180 tests the direction queue anti-180 and rapid turning mechanics.
func TestExtreme_DirectionQueueRapid180(t *testing.T) {
	cfg := DefaultConfig()
	g := New(cfg)
	// Snake is moving DirRight

	// Direct opposite should be rejected
	if g.QueueDirection(DirLeft) {
		t.Errorf("Queueing DirLeft while moving DirRight must be rejected")
	}

	// Rapid double turn: Up then Left (L-turn)
	if !g.QueueDirection(DirUp) {
		t.Errorf("Queueing DirUp should succeed")
	}
	if !g.QueueDirection(DirLeft) {
		t.Errorf("Queueing DirLeft after DirUp should succeed (L-turn)")
	}

	// Third turn should be rejected (max queue length 2)
	if g.QueueDirection(DirDown) {
		t.Errorf("Queueing 3rd turn should be rejected (buffer full)")
	}

	// First tick consumes DirUp
	g.Tick()
	if g.Direction != DirUp {
		t.Errorf("Expected Direction DirUp, got %v", g.Direction)
	}

	// Second tick consumes DirLeft
	g.Tick()
	if g.Direction != DirLeft {
		t.Errorf("Expected Direction DirLeft, got %v", g.Direction)
	}
}

// TestExtreme_TailVacatingRule tests moving into current tail position.
func TestExtreme_TailVacatingRule(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.InitialSnakeLen = 4

	// Scenario 1: Tail vacates normally (growth pending == 0) -> SAFE
	g := New(cfg)
	// Snake at (5,5), (4,5), (4,4), (5,4) - a 2x2 loop
	g.Snake = []Point{
		{X: 5, Y: 5},
		{X: 4, Y: 5},
		{X: 4, Y: 4},
		{X: 5, Y: 4}, // tail
	}
	g.Direction = DirUp // New head will be (5, 4), which is current tail
	g.GrowthPending = 0

	g.Tick()
	if g.GameOver {
		t.Errorf("Moving into vacating tail should be SAFE when no growth is pending")
	}

	// Scenario 2: Growth IS pending -> Tail does NOT vacate -> CRASH
	g2 := New(cfg)
	g2.Snake = []Point{
		{X: 5, Y: 5},
		{X: 4, Y: 5},
		{X: 4, Y: 4},
		{X: 5, Y: 4}, // tail
	}
	g2.Direction = DirUp
	g2.GrowthPending = 1 // Growth pending from previously eaten apple

	g2.Tick()
	if !g2.GameOver {
		t.Errorf("Moving into tail when growth is pending must cause Game Over (tail has not vacated)")
	}
	if g2.GameOverReason != "Bit own tail" {
		t.Errorf("Expected crash reason 'Bit own tail', got '%s'", g2.GameOverReason)
	}
}

// TestExtreme_ProceduralObstaclesSafetyZoneAllPatterns tests all obstacle patterns across 100 seeds.
func TestExtreme_ProceduralObstaclesSafetyZoneAllPatterns(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		w, h := 24, 20
		obstacles := GenerateProceduralObstacles(w, h, seed, 0.10)
		centerX := w / 2
		centerY := h / 2

		for p := range obstacles {
			// Bounds check
			if p.X < 0 || p.X >= w || p.Y < 0 || p.Y >= h {
				t.Fatalf("Seed %d: obstacle (%d, %d) is outside board bounds (%d, %d)", seed, p.X, p.Y, w, h)
			}
			// Safety zone check: x in [centerX-4, centerX+2] and y in [centerY-2, centerY+2]
			if p.X >= centerX-4 && p.X <= centerX+2 && p.Y >= centerY-2 && p.Y <= centerY+2 {
				t.Fatalf("Seed %d: obstacle (%d, %d) spawned inside safety clearance zone (center %d, %d)",
					seed, p.X, p.Y, centerX, centerY)
			}
		}
	}
}
