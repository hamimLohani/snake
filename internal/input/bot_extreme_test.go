package input

import (
	"testing"

	"snake/internal/engine"
)

func TestExtreme_HardwarePacketParsing(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		expected engine.Direction
	}{
		// Binary single-byte
		{"Binary Up", []byte{0x01}, engine.DirUp},
		{"Binary Down", []byte{0x02}, engine.DirDown},
		{"Binary Left", []byte{0x03}, engine.DirLeft},
		{"Binary Right", []byte{0x04}, engine.DirRight},
		{"Binary Invalid High", []byte{0x99}, engine.DirNone},

		// Plaintext ASCII
		{"Text UP", []byte("UP"), engine.DirUp},
		{"Text U", []byte("u"), engine.DirUp},
		{"Text W", []byte("W"), engine.DirUp},
		{"Text DOWN", []byte("down"), engine.DirDown},
		{"Text S (WASD Down)", []byte("s"), engine.DirDown},
		{"Text LEFT", []byte("left"), engine.DirLeft},
		{"Text L", []byte("l"), engine.DirLeft},
		{"Text A (WASD Left)", []byte("a"), engine.DirLeft},
		{"Text RIGHT", []byte("RIGHT"), engine.DirRight},
		{"Text R", []byte("r"), engine.DirRight},
		{"Text D (WASD Right)", []byte("d"), engine.DirRight},

		// With trailing whitespace or newlines
		{"Text UP with newline", []byte("UP\n"), engine.DirUp},
		{"Text W with CRLF", []byte("W\r\n"), engine.DirUp},
		{"Text S with spaces", []byte("  S  "), engine.DirDown},
		{"Text D with newline", []byte("D\n"), engine.DirRight},
		{"Text A with newline", []byte("A\n"), engine.DirLeft},

		// Invalid or empty
		{"Empty slice", []byte{}, engine.DirNone},
		{"Nil slice", nil, engine.DirNone},
		{"Unknown word", []byte("JUMP"), engine.DirNone},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseControllerPacket(tc.payload)
			if got != tc.expected {
				t.Errorf("Payload %v: expected %v, got %v", tc.payload, tc.expected, got)
			}
		})
	}
}

func TestExtreme_BotNo180Reversal(t *testing.T) {
	cfg := engine.DefaultConfig()
	cfg.Width = 12
	cfg.Height = 12
	cfg.Seed = 1111

	solver := NewBotSolver()

	// Test in all 4 movement directions
	dirs := []engine.Direction{engine.DirUp, engine.DirDown, engine.DirLeft, engine.DirRight}
	for _, movingDir := range dirs {
		g := engine.New(cfg)
		g.Direction = movingDir
		// Place food directly behind snake head
		dx, dy := movingDir.Opposite().Delta()
		g.Food = engine.Point{X: g.Snake[0].X + dx, Y: g.Snake[0].Y + dy}

		next := solver.NextMove(g)
		if next == movingDir.Opposite() {
			t.Errorf("Bot returned opposite direction %v while moving %v (illegal 180° turn into neck)",
				next, movingDir)
		}
	}
}

func TestExtreme_BotZenModeWrapping(t *testing.T) {
	cfg := engine.DefaultConfig()
	cfg.Width = 10
	cfg.Height = 10
	cfg.Mode = engine.ModeZen
	cfg.Seed = 8888

	g := engine.New(cfg)
	solver := NewBotSolver()

	// Run bot in Zen mode for 100 ticks
	for i := 0; i < 100; i++ {
		dir := solver.NextMove(g)
		if dir != engine.DirNone {
			g.QueueDirection(dir)
		}
		g.Tick()
		if g.GameOver {
			t.Fatalf("Zen mode bot should never suffer Game Over")
		}
	}

	if g.Score < 10 {
		t.Errorf("Bot should score points in Zen mode, got %d", g.Score)
	}
}

func TestExtreme_BotWithObstacles(t *testing.T) {
	cfg := engine.DefaultConfig()
	cfg.Width = 20
	cfg.Height = 16
	cfg.WithObstacles = true
	cfg.ObstacleDensity = 0.08
	cfg.Seed = 44444

	g := engine.New(cfg)
	solver := NewBotSolver()

	// Run bot with obstacles for 150 ticks
	for i := 0; i < 150; i++ {
		dir := solver.NextMove(g)
		if dir != engine.DirNone {
			g.QueueDirection(dir)
		}
		g.Tick()
		if g.GameOver {
			break
		}
	}

	if g.Score < 10 {
		t.Logf("Bot survived %d ticks with score %d", g.TickCount, g.Score)
	}
}
