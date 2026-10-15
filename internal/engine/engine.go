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
	BonusPoint       Point
	BonusType        ItemType
	BonusActive      bool
	BonusExpires     uint64
	FoodItems        []FoodItem
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

	var obstacles map[Point]bool
	if cfg.WithObstacles {
		density := cfg.ObstacleDensity
		if density <= 0 {
			density = 0.05
		}
		obstacles = GenerateProceduralObstacles(cfg.Width, cfg.Height, cfg.Seed, density)
	} else {
		obstacles = make(map[Point]bool)
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
		Obstacles:        obstacles,
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

	// 2b. Check item on-board expiration (5s timer for special items)
	var remainingItems []FoodItem
	expiredAny := false
	for _, item := range g.FoodItems {
		if item.ExpiresAtTick > 0 && g.TickCount >= item.ExpiresAtTick {
			expiredAny = true
		} else {
			remainingItems = append(remainingItems, item)
		}
	}
	if expiredAny {
		g.FoodItems = remainingItems
		g.spawnNeededItems()
	} else if g.BonusActive && g.TickCount >= g.BonusExpires {
		g.BonusActive = false
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
	ateItem := false
	eatenIndex := -1
	var eatenItem FoodItem
	for i, item := range g.FoodItems {
		if item.Point == head {
			eatenIndex = i
			eatenItem = item
			break
		}
	}

	if eatenIndex >= 0 {
		ateItem = true
		g.FoodItems = append(g.FoodItems[:eatenIndex], g.FoodItems[eatenIndex+1:]...)
		g.applyItemEffect(eatenItem)
		g.spawnNeededItems()
	} else if head == g.Food {
		g.consumeFood()
		ateItem = true
	} else if g.BonusActive && head == g.BonusPoint {
		g.consumeBonus()
		ateItem = true
	}

	if g.GrowthPending > 0 {
		g.GrowthPending--
	} else {
		// Pop tail
		g.Snake = g.Snake[:len(g.Snake)-1]
	}

	if !ateItem {
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventMove,
			Position: head,
		})
	}

	return g.LastEvents
}

// TicksForDuration converts real-world duration to engine ticks based on BaseTickRateMs.
func (g *Game) TicksForDuration(d time.Duration) uint64 {
	ms := g.Config.BaseTickRateMs
	if ms <= 0 {
		ms = 110
	}
	ticks := uint64(d.Milliseconds()) / uint64(ms)
	if ticks < 1 {
		ticks = 1
	}
	return ticks
}

// TargetFoodCount returns the target count of permanent foods on board (single food, respawns immediately).
func (g *Game) TargetFoodCount() int {
	return 1
}

// TargetSpecialCount returns target count of special featured items on board (powerup + poison hazard).
func (g *Game) TargetSpecialCount() int {
	return 2
}

// ItemAt returns the item type and true if an item exists at pt.
func (g *Game) ItemAt(pt Point) (ItemType, bool) {
	for _, item := range g.FoodItems {
		if item.Point == pt {
			return item.Type, true
		}
	}
	if pt == g.Food {
		return g.FoodType, true
	}
	if g.BonusActive && pt == g.BonusPoint {
		return g.BonusType, true
	}
	return ItemFood, false
}

// syncFoodFields syncs legacy Food and Bonus fields for backward compatibility.
func (g *Game) syncFoodFields() {
	foundFood := false
	for _, item := range g.FoodItems {
		if item.Type == ItemFood {
			g.Food = item.Point
			g.FoodType = ItemFood
			foundFood = true
			break
		}
	}
	if !foundFood && len(g.FoodItems) > 0 {
		g.Food = g.FoodItems[0].Point
		g.FoodType = g.FoodItems[0].Type
	}

	foundBonus := false
	for _, item := range g.FoodItems {
		if item.Type != ItemFood {
			g.BonusPoint = item.Point
			g.BonusType = item.Type
			g.BonusActive = true
			g.BonusExpires = item.ExpiresAtTick
			foundBonus = true
			break
		}
	}
	if !foundBonus {
		g.BonusActive = false
	}
}

// applyItemEffect processes score, growth, and 5-second duration effects.
func (g *Game) applyItemEffect(item FoodItem) {
	ticks5s := g.TicksForDuration(5 * time.Second)
	switch item.Type {
	case ItemFood:
		basePoints := 10
		if g.Config.Mode == ModeTimeAttack {
			basePoints = 25
			g.TimeRemainingSec += 3
		}
		earned := basePoints * g.ScoreMultiplier
		g.Score += earned
		g.GrowthPending += 1
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatFood,
			Position: item.Point,
			Score:    earned,
			Message:  fmt.Sprintf("+%d pts", earned),
		})
	case ItemSlowMo:
		g.ActivePowerUp = &ActivePowerUp{Type: ItemSlowMo, ExpiresAtTick: g.TickCount + ticks5s}
		g.Score += 20 * g.ScoreMultiplier
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatPowerUp,
			Position: item.Point,
			Score:    20 * g.ScoreMultiplier,
			Message:  "Slow-Mo Activated! (5s)",
		})
	case ItemGhost:
		g.ActivePowerUp = &ActivePowerUp{Type: ItemGhost, ExpiresAtTick: g.TickCount + ticks5s}
		g.Score += 20 * g.ScoreMultiplier
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatPowerUp,
			Position: item.Point,
			Score:    20 * g.ScoreMultiplier,
			Message:  "Ghost Mode Activated! (5s)",
		})
	case ItemMultiplier:
		g.ActivePowerUp = &ActivePowerUp{Type: ItemMultiplier, ExpiresAtTick: g.TickCount + ticks5s}
		g.ScoreMultiplier = 3
		g.Score += 30
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatPowerUp,
			Position: item.Point,
			Score:    30,
			Message:  "3x Multiplier Activated! (5s)",
		})
	case ItemShrink:
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
			Position: item.Point,
			Score:    15 * g.ScoreMultiplier,
			Message:  "Snake Shrunk!",
		})
	case ItemPoison:
		penalty := 30
		g.Score -= penalty
		if g.Score < 0 {
			g.Score = 0
		}
		g.GrowthPending += 1
		g.ActivePowerUp = &ActivePowerUp{Type: ItemPoison, ExpiresAtTick: g.TickCount + ticks5s}
		g.LastEvents = append(g.LastEvents, Event{
			Type:     EventEatHazard,
			Position: item.Point,
			Score:    -penalty,
			Message:  "⚠️ Poison Ingested! -30 pts & Speed Surge (5s)!",
		})
	}
}

// consumeFood handles points, power-ups, and spawning next food (legacy/direct calls).
func (g *Game) consumeFood() {
	g.applyItemEffect(FoodItem{Point: g.Food, Type: g.FoodType})
	g.spawnFood()
}

// consumeBonus handles bonus items (legacy/direct calls).
func (g *Game) consumeBonus() {
	g.BonusActive = false
	g.applyItemEffect(FoodItem{Point: g.BonusPoint, Type: g.BonusType})
}

// findOpenCell finds an empty grid coordinate not occupied by snake, obstacles, or food items.
func (g *Game) findOpenCell() (Point, bool) {
	totalCells := (g.Bounds.MaxX - g.Bounds.MinX) * (g.Bounds.MaxY - g.Bounds.MinY)
	occupiedCount := len(g.Snake) + len(g.Obstacles) + len(g.FoodItems) + 2
	if occupiedCount >= totalCells {
		return Point{}, false
	}

	occupied := make(map[Point]bool, occupiedCount)
	for _, p := range g.Snake {
		occupied[p] = true
	}
	for o := range g.Obstacles {
		occupied[o] = true
	}
	for _, f := range g.FoodItems {
		occupied[f.Point] = true
	}
	occupied[g.Food] = true
	if g.BonusActive {
		occupied[g.BonusPoint] = true
	}

	width := g.Bounds.MaxX - g.Bounds.MinX
	height := g.Bounds.MaxY - g.Bounds.MinY

	for attempts := 0; attempts < 500; attempts++ {
		cand := Point{
			X: g.Bounds.MinX + g.RNG.IntN(width),
			Y: g.Bounds.MinY + g.RNG.IntN(height),
		}
		if !occupied[cand] {
			return cand, true
		}
	}

	// Fallback scan
	for y := g.Bounds.MinY; y < g.Bounds.MaxY; y++ {
		for x := g.Bounds.MinX; x < g.Bounds.MaxX; x++ {
			p := Point{X: x, Y: y}
			if !occupied[p] {
				return p, true
			}
		}
	}

	return Point{}, false
}

// spawnFood places food on open cells and ensures required food/special count.
func (g *Game) spawnFood() {
	g.spawnNeededItems()
}

// spawnBonusItem spawns an extra bonus item alongside permanent food.
func (g *Game) spawnBonusItem() {
	g.spawnNeededItems()
}

// spawnNeededItems maintains multiple permanent foods and special items (with poison hazard).
func (g *Game) spawnNeededItems() {
	totalCells := (g.Bounds.MaxX - g.Bounds.MinX) * (g.Bounds.MaxY - g.Bounds.MinY)
	if len(g.Snake)+len(g.Obstacles) >= totalCells {
		g.GameOver = true
		g.GameOverReason = "Victory! Board fully cleared!"
		g.LastEvents = append(g.LastEvents, Event{
			Type:    EventGameOver,
			Score:   g.Score,
			Message: g.GameOverReason,
		})
		return
	}

	targetFood := g.TargetFoodCount()
	targetSpecial := g.TargetSpecialCount()

	currentFood := 0
	currentSpecial := 0
	hasPoison := false
	for _, item := range g.FoodItems {
		if item.Type == ItemFood {
			currentFood++
		} else {
			currentSpecial++
			if item.Type == ItemPoison {
				hasPoison = true
			}
		}
	}

	// 1. Replenish permanent foods (multiple on board simultaneously)
	for currentFood < targetFood {
		pt, ok := g.findOpenCell()
		if !ok {
			break
		}
		g.FoodItems = append(g.FoodItems, FoodItem{
			Point: pt,
			Type:  ItemFood,
		})
		currentFood++
	}

	// 2. Replenish special items (guarantee 1 poison hazard food + 1 powerup, each for 5s)
	ticks5s := g.TicksForDuration(5 * time.Second)
	for currentSpecial < targetSpecial {
		pt, ok := g.findOpenCell()
		if !ok {
			break
		}

		var itemType ItemType
		if !hasPoison && g.Config.Mode != ModeZen {
			itemType = ItemPoison
			hasPoison = true
		} else {
			roll := g.RNG.Float64()
			if roll < 0.25 {
				itemType = ItemSlowMo
			} else if roll < 0.50 {
				itemType = ItemGhost
			} else if roll < 0.75 {
				itemType = ItemMultiplier
			} else {
				itemType = ItemShrink
			}
		}

		g.FoodItems = append(g.FoodItems, FoodItem{
			Point:         pt,
			Type:          itemType,
			ExpiresAtTick: g.TickCount + ticks5s, // 5 seconds on board
		})
		currentSpecial++
	}

	g.syncFoodFields()
}

// TickDuration calculates the current tick interval in milliseconds.
func (g *Game) TickDuration() time.Duration {
	base := g.Config.BaseTickRateMs

	dur := base
	if g.Config.AdaptiveSpeed {
		speedBonus := (g.Score / 50) * 4
		dur = base - speedBonus
		if dur < 40 {
			dur = 40
		}
	}

	if g.ActivePowerUp != nil {
		switch g.ActivePowerUp.Type {
		case ItemSlowMo:
			return time.Duration(dur*2) * time.Millisecond
		case ItemPoison:
			poisonDur := dur * 65 / 100
			if poisonDur < 30 {
				poisonDur = 30
			}
			return time.Duration(poisonDur) * time.Millisecond
		}
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
