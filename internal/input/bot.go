package input

import (
	"snake/internal/engine"
)

// BotSolver implements pathfinding and survival heuristics for Snake.
type BotSolver struct{}

// NewBotSolver creates a new solver.
func NewBotSolver() *BotSolver {
	return &BotSolver{}
}

// NextMove calculates the next optimal Direction given the current engine state.
func (b *BotSolver) NextMove(g *engine.Game) engine.Direction {
	if g == nil || len(g.Snake) == 0 || g.GameOver {
		return engine.DirNone
	}

	head := g.Snake[0]
	tail := g.Snake[len(g.Snake)-1]

	// 1. Try finding safe path to food (support multiple foods and beneficial power-ups)
	var foodTargets []engine.Point
	for _, item := range g.FoodItems {
		if item.Type != engine.ItemPoison {
			foodTargets = append(foodTargets, item.Point)
		}
	}
	if len(foodTargets) == 0 {
		foodTargets = append(foodTargets, g.Food)
	}

	var shortestPathToFood []engine.Point
	for _, target := range foodTargets {
		path := b.findPathBFS(g, head, target)
		if len(path) > 0 {
			if len(shortestPathToFood) == 0 || len(path) < len(shortestPathToFood) {
				shortestPathToFood = path
			}
		}
	}

	if len(shortestPathToFood) > 0 {
		nextStep := shortestPathToFood[0]
		// Verify safety: if we take this step, can we still reach our tail?
		sim := b.simulateStep(g, nextStep)
		if sim != nil && (len(b.findPathBFS(sim, sim.Snake[0], tail)) > 0 || len(g.Snake) < 5) {
			return b.directionFromPoints(head, nextStep)
		}
	}

	// 2. Fallback: Path to tail (survival / stall tactic)
	pathToTail := b.findPathBFS(g, head, tail)
	if len(pathToTail) > 0 {
		return b.directionFromPoints(head, pathToTail[0])
	}

	// 3. Fallback: Any valid neighbor with maximum open space (flood fill heuristic)
	bestDir := engine.DirNone
	maxSpace := -1

	directions := []engine.Direction{engine.DirUp, engine.DirDown, engine.DirLeft, engine.DirRight}
	for _, dir := range directions {
		if dir == g.Direction.Opposite() {
			continue
		}
		dx, dy := dir.Delta()
		cand := engine.Point{X: head.X + dx, Y: head.Y + dy}

		if b.isCellSafe(g, cand) {
			space := b.countReachableCells(g, cand)
			if space > maxSpace {
				maxSpace = space
				bestDir = dir
			}
		}
	}

	if bestDir != engine.DirNone {
		return bestDir
	}

	// Final default: keep moving forward if possible
	return g.Direction
}

func (b *BotSolver) isCellSafe(g *engine.Game, p engine.Point) bool {
	isGhost := g.ActivePowerUp != nil && g.ActivePowerUp.Type == engine.ItemGhost
	if g.Config.Mode == engine.ModeZen || isGhost {
		// In Zen mode or Ghost, walls, obstacles, and body overlap are safe
		// Only avoid poison hazard
		if g.BonusActive && g.BonusType == engine.ItemPoison && p == g.BonusPoint {
			return false
		}
		for _, item := range g.FoodItems {
			if item.Type == engine.ItemPoison && item.Point == p {
				return false
			}
		}
		return true
	}

	// Wall check
	if p.X < g.Bounds.MinX || p.X >= g.Bounds.MaxX || p.Y < g.Bounds.MinY || p.Y >= g.Bounds.MaxY {
		return false
	}
	// Obstacle check
	if g.Obstacles[p] {
		return false
	}
	// Poison hazard check (avoid negative impact food)
	if g.BonusActive && g.BonusType == engine.ItemPoison && p == g.BonusPoint {
		return false
	}
	for _, item := range g.FoodItems {
		if item.Type == engine.ItemPoison && item.Point == p {
			return false
		}
	}
	// Body check (tail is safe because it will vacate on next step)
	tail := g.Snake[len(g.Snake)-1]
	for i := 0; i < len(g.Snake)-1; i++ {
		if g.Snake[i] == p && p != tail {
			return false
		}
	}
	return true
}

func (b *BotSolver) findPathBFS(g *engine.Game, start, target engine.Point) []engine.Point {
	if start == target {
		return nil
	}

	queue := []engine.Point{start}
	parent := make(map[engine.Point]engine.Point)
	visited := make(map[engine.Point]bool)
	visited[start] = true

	directions := []engine.Direction{engine.DirUp, engine.DirDown, engine.DirLeft, engine.DirRight}
	isGhost := g.ActivePowerUp != nil && g.ActivePowerUp.Type == engine.ItemGhost

	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == target {
			found = true
			break
		}

		for _, dir := range directions {
			dx, dy := dir.Delta()
			neighbor := engine.Point{X: curr.X + dx, Y: curr.Y + dy}

			if g.Config.Mode == engine.ModeZen || isGhost {
				if neighbor.X < g.Bounds.MinX {
					neighbor.X = g.Bounds.MaxX - 1
				} else if neighbor.X >= g.Bounds.MaxX {
					neighbor.X = g.Bounds.MinX
				}
				if neighbor.Y < g.Bounds.MinY {
					neighbor.Y = g.Bounds.MaxY - 1
				} else if neighbor.Y >= g.Bounds.MaxY {
					neighbor.Y = g.Bounds.MinY
				}
			}

			// If neighbor is target, it is allowed even if occupied
			if neighbor == target || b.isCellSafe(g, neighbor) {
				if !visited[neighbor] {
					visited[neighbor] = true
					parent[neighbor] = curr
					queue = append(queue, neighbor)
				}
			}
		}
	}

	if !found {
		return nil
	}

	// Reconstruct path
	var path []engine.Point
	curr := target
	for curr != start {
		path = append([]engine.Point{curr}, path...)
		curr = parent[curr]
	}
	return path
}

func (b *BotSolver) countReachableCells(g *engine.Game, start engine.Point) int {
	queue := []engine.Point{start}
	visited := make(map[engine.Point]bool)
	visited[start] = true
	count := 0

	directions := []engine.Direction{engine.DirUp, engine.DirDown, engine.DirLeft, engine.DirRight}
	for len(queue) > 0 && count < 64 {
		curr := queue[0]
		queue = queue[1:]
		count++

		for _, dir := range directions {
			dx, dy := dir.Delta()
			neighbor := engine.Point{X: curr.X + dx, Y: curr.Y + dy}
			if !visited[neighbor] && b.isCellSafe(g, neighbor) {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return count
}

func (b *BotSolver) simulateStep(g *engine.Game, nextHead engine.Point) *engine.Game {
	newSnake := make([]engine.Point, len(g.Snake))
	copy(newSnake, g.Snake)
	newSnake = append([]engine.Point{nextHead}, newSnake[:len(newSnake)-1]...)

	sim := *g
	sim.Snake = newSnake
	return &sim
}

func (b *BotSolver) directionFromPoints(from, to engine.Point) engine.Direction {
	// Handle toroidal screen wrapping
	if to.X == from.X-1 || to.X > from.X+1 {
		return engine.DirLeft
	}
	if to.X == from.X+1 || to.X < from.X-1 {
		return engine.DirRight
	}
	if to.Y == from.Y-1 || to.Y > from.Y+1 {
		return engine.DirUp
	}
	if to.Y == from.Y+1 || to.Y < from.Y-1 {
		return engine.DirDown
	}
	return engine.DirNone
}
