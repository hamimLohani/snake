package engine

import (
	"fmt"
	"time"
)

// ActivePowerUp tracks an ongoing temporary bonus.
type ActivePowerUp struct {
	Type          ItemType
	ExpiresAtTick uint64
}

// Bounds defines the playable box on the grid.
type Bounds struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

// Game implements pure deterministic Snake mechanics.
type Game struct {
	Config           Config
	Snake            []Point
	Direction        Direction
	DirQueue         []Direction
	Food             Point
	FoodType         ItemType
	FoodExpiresAt    uint64
	Score            int
	Highscore        int
	TickCount        uint64
	GrowthPending    int
	GameOver         bool
	GameOverReason   string
	Bounds           Bounds
	Obstacles        map[Point]bool
	ActivePowerUp    *ActivePowerUp
	ScoreMultiplier  int
	TimeRemainingSec int // For Time Attack
	RNG              *PRNG
	LastEvents       []Event
}

// New creates a new game engine instance.
func New(cfg Config) *Game {
	if cfg.Width <= 4 {
		cfg.Width = 24
	}
	if cfg.Height <= 4 {
		cfg.Height = 20
	}
	if cfg.InitialSnakeLen <= 0 {
		cfg.InitialSnakeLen = 3
	}
	if cfg.BaseTickRateMs <= 0 {
		cfg.BaseTickRateMs = 110
	}
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}

	rng := NewPRNG(cfg.Seed)

	startX := cfg.Width / 2
	startY := cfg.Height / 2

	snake := make([]Point, cfg.InitialSnakeLen)
	for i := 0; i < cfg.InitialSnakeLen; i++ {
		snake[i] = Point{X: startX - i, Y: startY}
	}

	g := &Game{
		Config:          cfg,
		Snake:           snake,
		Direction:       DirRight,
		DirQueue:        make([]Direction, 0, 4),
		Score:           0,
		ScoreMultiplier: 1,
		Bounds: Bounds{
			MinX: 0,
			MaxX: cfg.Width,
			MinY: 0,
			MaxY: cfg.Height,
		},
		Obstacles:        make(map[Point]bool),
		RNG:              rng,
		LastEvents:       make([]Event, 0),
		TimeRemainingSec: 60, // 60s for time attack
	}

	g.spawnFood()
	return g
}

// QueueDirection pushes a requested direction change to the input buffer.
func (g *Game) QueueDirection(dir Direction) bool {
	if dir == DirNone || g.GameOver {
		return false
	}

	// Reference against the last planned direction (either in queue or current)
	current := g.Direction
	if len(g.DirQueue) > 0 {
		current = g.DirQueue[len(g.DirQueue)-1]
	}

	// Cannot reverse directly into opposite direction
	if dir == current || dir == current.Opposite() {
		return false
	}

	// Keep buffer reasonable (max 2 queued turns)
	if len(g.DirQueue) < 2 {
		g.DirQueue = append(g.DirQueue, dir)
		return true
	}
	return false
}

// Tick advances the game state by one unit of time.
func (g *Game) Tick() []Event {
	g.LastEvents = g.LastEvents[:0]
	if g.GameOver {
		return g.LastEvents
	}

	g.TickCount++

	// 1. Process queued direction change
	if len(g.DirQueue) > 0 {
		g.Direction = g.DirQueue[0]
		g.DirQueue = g.DirQueue[1:]
	}

	// 2. Check active power-up expiration
	if g.ActivePowerUp != nil && g.TickCount >= g.ActivePowerUp.ExpiresAtTick {
		g.LastEvents = append(g.LastEvents, Event{
			Type:    EventPowerUpExpire,
			Message: fmt.Sprintf("%s wore off", g.ActivePowerUp.Type.String()),
		})
		if g.ActivePowerUp.Type == ItemMultiplier {
			g.ScoreMultiplier = 1
		}
		g.ActivePowerUp = nil
	}

	// 3. Mode-specific tick adjustments
	if g.Config.Mode == ModeSurvival {
		// Shrink bounds every 50 ticks, minimum 6x6
		if g.TickCount%50 == 0 && (g.Bounds.MaxX-g.Bounds.MinX > 6) && (g.Bounds.MaxY-g.Bounds.MinY > 6) {
			g.Bounds.MinX++
			g.Bounds.MaxX--
			g.Bounds.MinY++
			g.Bounds.MaxY--
			g.LastEvents = append(g.LastEvents, Event{
				Type:    EventWallShrink,
				Message: "Warning! The walls are closing in!",
			})
		}
	} else if g.Config.Mode == ModeTimeAttack {
		// Ticks are ~100ms, every 10 ticks = 1 second
		if g.TickCount%10 == 0 && g.TimeRemainingSec > 0 {
			g.TimeRemainingSec--
			if g.TimeRemainingSec <= 0 {
				g.GameOver = true
				g.GameOverReason = "Time is up!"
				g.LastEvents = append(g.LastEvents, Event{
					Type:    EventGameOver,
					Message: g.GameOverReason,
					Score:   g.Score,
				})
				return g.LastEvents
			}
		}
	}

	// 4. Calculate next head position
	dx, dy := g.Direction.Delta()
	head := Point{X: g.Snake[0].X + dx, Y: g.Snake[0].Y + dy}

	// 5. Handle boundaries & wrapping
	if g.Config.Mode == ModeZen {
		// Zen wraps around boundaries
		width := g.Bounds.MaxX - g.Bounds.MinX
		height := g.Bounds.MaxY - g.Bounds.MinY
		if head.X < g.Bounds.MinX {
			head.X = g.Bounds.MaxX - 1
		} else if head.X >= g.Bounds.MaxX {
			head.X = g.Bounds.MinX
		}
		if head.Y < g.Bounds.MinY {
			head.Y = g.Bounds.MaxY - 1
		} else if head.Y >= g.Bounds.MaxY {
			head.Y = g.Bounds.MinY
		}
		_ = width
		_ = height
	} else {
		// Wall collision checks
		if head.X < g.Bounds.MinX || head.X >= g.Bounds.MaxX || head.Y < g.Bounds.MinY || head.Y >= g.Bounds.MaxY {
			g.GameOver = true
			g.GameOverReason = "Crashed into wall"
			g.LastEvents = append(g.LastEvents, Event{
				Type:     EventWallCollision,
				Position: head,
				Score:    g.Score,
				Message:  g.GameOverReason,
			})
			g.LastEvents = append(g.LastEvents, Event{
				Type:    EventGameOver,
				Score:   g.Score,
				Message: g.GameOverReason,
			})
			return g.LastEvents
		}

		// Obstacle collision checks
		if g.Obstacles[head] {
			g.GameOver = true
			g.GameOverReason = "Crashed into obstacle"
			g.LastEvents = append(g.LastEvents, Event{
				Type:     EventWallCollision,
				Position: head,
				Score:    g.Score,
				Message:  g.GameOverReason,
			})
			g.LastEvents = append(g.LastEvents, Event{
				Type:    EventGameOver,
				Score:   g.Score,
				Message: g.GameOverReason,
			})
			return g.LastEvents
		}
	}

	// 6. Self collision check (unless Zen or Ghost powerup is active)
	isGhost := g.ActivePowerUp != nil && g.ActivePowerUp.Type == ItemGhost
	if g.Config.Mode != ModeZen && !isGhost {
		// Tail moves unless growth is pending
		bodyToCheck := g.Snake
		if g.GrowthPending == 0 && len(g.Snake) > 1 {
			bodyToCheck = g.Snake[:len(g.Snake)-1]
		}
		for _, seg := range bodyToCheck {
			if seg == head {
				g.GameOver = true
				g.GameOverReason = "Bit own tail"
				g.LastEvents = append(g.LastEvents, Event{
					Type:     EventSelfCollision,
					Position: head,
					Score:    g.Score,
					Message:  g.GameOverReason,
				})
				g.LastEvents = append(g.LastEvents, Event{
					Type:    EventGameOver,
					Score:   g.Score,
					Message: g.GameOverReason,
				})
				return g.LastEvents
			}
		}
	}

	// 7. Advance snake
	g.Snake = append([]Point{head}, g.Snake...)

	// 8. Did we eat an item?
	if head == g.Food {
		g.consumeFood()
	}

	if g.GrowthPending > 0 {
		g.GrowthPending--
	} else {
		// Pop tail
		g.Snake = g.Snake[:len(g.Snake)-1]
	}

	if head != g.Food {
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventMove,
			Position: head,
		})
	}

	return g.LastEvents
}

// consumeFood handles points, power-ups, and spawning next food.
func (g *Game) consumeFood() {
	item := g.FoodType
	basePoints := 10
	if g.Config.Mode == ModeTimeAttack {
		basePoints = 25
		g.TimeRemainingSec += 3 // Bonus time on eat
	}

	switch item {
	case ItemFood:
		earned := basePoints * g.ScoreMultiplier
		g.Score += earned
		g.GrowthPending += 1
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatFood,
			Position: g.Food,
			Score:    earned,
			Message:  fmt.Sprintf("+%d pts", earned),
		})
	case ItemSlowMo:
		g.ActivePowerUp = &ActivePowerUp{Type: ItemSlowMo, ExpiresAtTick: g.TickCount + 60}
		g.Score += 20 * g.ScoreMultiplier
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatPowerUp,
			Position: g.Food,
			Message:  "Slow-Mo Activated! (50% speed for 60 ticks)",
		})
	case ItemGhost:
		g.ActivePowerUp = &ActivePowerUp{Type: ItemGhost, ExpiresAtTick: g.TickCount + 50}
		g.Score += 20 * g.ScoreMultiplier
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatPowerUp,
			Position: g.Food,
			Message:  "Ghost Mode Activated! (Pass through self)",
		})
	case ItemMultiplier:
		g.ActivePowerUp = &ActivePowerUp{Type: ItemMultiplier, ExpiresAtTick: g.TickCount + 70}
		g.ScoreMultiplier = 3
		g.Score += 30
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatPowerUp,
			Position: g.Food,
			Message:  "3x Multiplier Activated!",
		})
	case ItemShrink:
		// Shrink cuts snake length in half (minimum 3)
		if len(g.Snake) > 5 {
			newLen := len(g.Snake) / 2
			if newLen < 3 {
				newLen = 3
			}
			g.Snake = g.Snake[:newLen]
		}
		g.Score += 15 * g.ScoreMultiplier
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatPowerUp,
			Position: g.Food,
			Message:  "Snake Shrunk!",
		})
	}

	g.spawnFood()
}

// spawnFood places food on an open cell.
func (g *Game) spawnFood() {
	totalCells := (g.Bounds.MaxX - g.Bounds.MinX) * (g.Bounds.MaxY - g.Bounds.MinY)
	if len(g.Snake)+len(g.Obstacles) >= totalCells {
		// Perfect win condition (board completely filled)
		g.GameOver = true
		g.GameOverReason = "Victory! Board fully cleared!"
		g.LastEvents = append(g.LastEvents, Event{
			Type:    EventGameOver,
			Score:   g.Score,
			Message: g.GameOverReason,
		})
		return
	}

	// Fast rejection sampling for food coordinates
	occupied := make(map[Point]bool, len(g.Snake)+len(g.Obstacles))
	for _, p := range g.Snake {
		occupied[p] = true
	}
	for o := range g.Obstacles {
		occupied[o] = true
	}

	width := g.Bounds.MaxX - g.Bounds.MinX
	height := g.Bounds.MaxY - g.Bounds.MinY

	for attempts := 0; attempts < 1000; attempts++ {
		cand := Point{
			X: g.Bounds.MinX + g.RNG.IntN(width),
			Y: g.Bounds.MinY + g.RNG.IntN(height),
		}
		if !occupied[cand] {
			g.Food = cand
			// Roll for food type (15% chance of power-up in Classic/Survival)
			roll := g.RNG.Float64()
			if roll < 0.85 {
				g.FoodType = ItemFood
			} else if roll < 0.89 {
				g.FoodType = ItemSlowMo
			} else if roll < 0.93 {
				g.FoodType = ItemGhost
			} else if roll < 0.97 {
				g.FoodType = ItemMultiplier
			} else {
				g.FoodType = ItemShrink
			}
			return
		}
	}

	// Fallback scan if grid is nearly full
	for y := g.Bounds.MinY; y < g.Bounds.MaxY; y++ {
		for x := g.Bounds.MinX; x < g.Bounds.MaxX; x++ {
			p := Point{X: x, Y: y}
			if !occupied[p] {
				g.Food = p
				g.FoodType = ItemFood
				return
			}
		}
	}
}

// TickDuration calculates the current tick interval in milliseconds.
func (g *Game) TickDuration() time.Duration {
	base := g.Config.BaseTickRateMs

	// Slow-mo halves game speed
	if g.ActivePowerUp != nil && g.ActivePowerUp.Type == ItemSlowMo {
		return time.Duration(base*2) * time.Millisecond
	}

	if !g.Config.AdaptiveSpeed {
		return time.Duration(base) * time.Millisecond
	}

	// Adaptive speed curve: increases slightly every 50 points, capped at 40ms minimum
	speedBonus := (g.Score / 50) * 4
	dur := base - speedBonus
	if dur < 40 {
		dur = 40
	}
	return time.Duration(dur) * time.Millisecond
}

// IsSnakeHead checks if p is the head.
func (g *Game) IsSnakeHead(p Point) bool {
	return len(g.Snake) > 0 && g.Snake[0] == p
}

// IsSnakeBody checks if p is part of the snake body (excluding head).
func (g *Game) IsSnakeBody(p Point) bool {
	for i := 1; i < len(g.Snake); i++ {
		if g.Snake[i] == p {
			return true
		}
	}
	return false
}
