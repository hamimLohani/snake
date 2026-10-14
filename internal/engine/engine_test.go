package engine

import (
	"reflect"
	"testing"
)

func TestSnakeMovement(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.InitialSnakeLen = 3
	cfg.Seed = 42

	game := New(cfg)
	// Initial head is at (5, 5), body at (4, 5), (3, 5)
	initialHead := game.Snake[0]

	// Move right
	events := game.Tick()
	if game.GameOver {
		t.Fatalf("Game ended unexpectedly: %v", events)
	}

	newHead := game.Snake[0]
	if newHead.X != initialHead.X+1 || newHead.Y != initialHead.Y {
		t.Errorf("Expected head at (%d, %d), got (%d, %d)", initialHead.X+1, initialHead.Y, newHead.X, newHead.Y)
	}
	if len(game.Snake) != 3 {
		t.Errorf("Expected snake len 3, got %d", len(game.Snake))
	}
}

func TestWallCollision(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 6
	cfg.Height = 6
	cfg.InitialSnakeLen = 2
	cfg.Seed = 100

	game := New(cfg)
	// Start moving right towards wall at X=6
	for i := 0; i < 10; i++ {
		game.Tick()
		if game.GameOver {
			break
		}
	}

	if !game.GameOver {
		t.Errorf("Expected game over on wall hit, but game is still running")
	}
	if game.GameOverReason != "Crashed into wall" {
		t.Errorf("Expected crash reason 'Crashed into wall', got '%s'", game.GameOverReason)
	}
}

func TestSelfCollision(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.InitialSnakeLen = 5
	cfg.Seed = 200

	game := New(cfg)
	// Make a quick loop: Up, Left, Down
	game.QueueDirection(DirUp)
	game.Tick()
	game.QueueDirection(DirLeft)
	game.Tick()
	game.QueueDirection(DirDown)
	game.Tick()

	if !game.GameOver {
		t.Errorf("Expected self collision game over, got running game")
	}
	if game.GameOverReason != "Bit own tail" {
		t.Errorf("Expected crash reason 'Bit own tail', got '%s'", game.GameOverReason)
	}
}

func TestFoodConsumptionAndGrowth(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.InitialSnakeLen = 3
	cfg.Seed = 300

	game := New(cfg)
	// Artificially place normal food directly in front of snake
	head := game.Snake[0]
	game.Food = Point{X: head.X + 1, Y: head.Y}
	game.FoodType = ItemFood

	initialLen := len(game.Snake)
	game.Tick() // Moves into food

	if game.Score <= 0 {
		t.Errorf("Expected score increase, got %d", game.Score)
	}
	// On next tick, growth is completed
	game.Tick()
	if len(game.Snake) != initialLen+1 {
		t.Errorf("Expected snake length %d after eating food, got %d", initialLen+1, len(game.Snake))
	}
}

func TestZenModeWrapAround(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 6
	cfg.Height = 6
	cfg.Mode = ModeZen
	cfg.InitialSnakeLen = 2
	cfg.Seed = 400

	game := New(cfg)
	// Advance through the right boundary
	for i := 0; i < 10; i++ {
		game.Tick()
		if game.GameOver {
			t.Fatalf("Zen mode should never trigger Game Over on wall collision")
		}
	}
}

func TestDeterminism(t *testing.T) {
	seed := int64(123456789)
	cfg1 := DefaultConfig()
	cfg1.Seed = seed
	g1 := New(cfg1)

	cfg2 := DefaultConfig()
	cfg2.Seed = seed
	g2 := New(cfg2)

	// Verify initial food is identical
	if g1.Food != g2.Food || g1.FoodType != g2.FoodType {
		t.Fatalf("Initial food mismatch with same seed: %+v vs %+v", g1.Food, g2.Food)
	}

	// Sequence of moves
	moves := []Direction{
		DirRight, DirDown, DirDown, DirLeft, DirUp, DirRight, DirDown, DirDown,
	}

	for i, d := range moves {
		g1.QueueDirection(d)
		g2.QueueDirection(d)
		g1.Tick()
		g2.Tick()

		if !reflect.DeepEqual(g1.Snake, g2.Snake) {
			t.Fatalf("Step %d: Snake mismatch.\ng1: %+v\ng2: %+v", i, g1.Snake, g2.Snake)
		}
		if g1.Food != g2.Food {
			t.Fatalf("Step %d: Food mismatch.\ng1: %+v\ng2: %+v", i, g1.Food, g2.Food)
		}
		if g1.Score != g2.Score {
			t.Fatalf("Step %d: Score mismatch.\ng1: %d, g2: %d", i, g1.Score, g2.Score)
		}
	}
}

func TestGhostPowerUp(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.InitialSnakeLen = 5
	cfg.Seed = 500

	game := New(cfg)
	game.ActivePowerUp = &ActivePowerUp{
		Type:          ItemGhost,
		ExpiresAtTick: 100,
	}

	// Loop back onto itself
	game.QueueDirection(DirUp)
	game.Tick()
	game.QueueDirection(DirLeft)
	game.Tick()
	game.QueueDirection(DirDown)
	game.Tick()

	if game.GameOver {
		t.Errorf("Ghost power-up should allow passing through own body without death")
	}
}
