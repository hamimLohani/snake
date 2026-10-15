package engine

// Point represents a 2D coordinate on the grid.
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Direction represents snake movement orientation.
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

// String returns the string representation of a Direction.
func (d Direction) String() string {
	switch d {
	case DirUp:
		return "UP"
	case DirDown:
		return "DOWN"
	case DirLeft:
		return "LEFT"
	case DirRight:
		return "RIGHT"
	default:
		return "NONE"
	}
}

// Opposite returns the reverse direction.
func (d Direction) Opposite() Direction {
	switch d {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	default:
		return DirNone
	}
}

// Delta returns the X and Y change for the direction.
func (d Direction) Delta() (dx, dy int) {
	switch d {
	case DirUp:
		return 0, -1
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	case DirRight:
		return 1, 0
	default:
		return 0, 0
	}
}

// Mode represents the game mode.
type Mode string

const (
	ModeClassic    Mode = "classic"
	ModeTimeAttack Mode = "time-attack"
	ModeSurvival   Mode = "survival"
	ModeZen        Mode = "zen"
)

// ItemType defines types of edible items on the board.
type ItemType int

const (
	ItemFood ItemType = iota
	ItemSlowMo
	ItemGhost
	ItemMultiplier
	ItemShrink
	ItemPoison
)

// String returns the name of the item.
func (it ItemType) String() string {
	switch it {
	case ItemFood:
		return "Food"
	case ItemSlowMo:
		return "Slow-Mo"
	case ItemGhost:
		return "Ghost"
	case ItemMultiplier:
		return "Multiplier"
	case ItemShrink:
		return "Shrink"
	case ItemPoison:
		return "Poison"
	default:
		return "Unknown"
	}
}

// EventType defines events emitted by the engine during a tick.
type EventType int

const (
	EventNone EventType = iota
	EventMove
	EventEatFood
	EventEatPowerUp
	EventEatHazard
	EventPowerUpExpire
	EventWallCollision
	EventSelfCollision
	EventGameOver
	EventWallShrink
)

// Event records an occurrence during a game tick.
type Event struct {
	Type     EventType
	Position Point
	Score    int
	Message  string
}

// Config defines the initialization parameters for a game engine.
type Config struct {
	Width            int
	Height           int
	Mode             Mode
	Seed             int64
	InitialSnakeLen  int
	AdaptiveSpeed    bool
	BaseTickRateMs   int
	WithObstacles    bool
	ObstacleDensity  float64
}

// DefaultConfig provides sensible defaults.
func DefaultConfig() Config {
	return Config{
		Width:           24,
		Height:          20,
		Mode:            ModeClassic,
		Seed:            0,
		InitialSnakeLen: 3,
		AdaptiveSpeed:   true,
		BaseTickRateMs:  110,
		WithObstacles:   false,
		ObstacleDensity: 0.05,
	}
}
