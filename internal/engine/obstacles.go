package engine

// GenerateProceduralObstacles places obstacles on the board deterministically based on seed.
// It ensures that the central region around the snake spawn point remains completely free.
func GenerateProceduralObstacles(width, height int, seed int64, density float64) map[Point]bool {
	obstacles := make(map[Point]bool)
	if density <= 0.0 {
		return obstacles
	}

	rng := NewPRNG(seed ^ 0x5a5a5a5a)

	centerX := width / 2
	centerY := height / 2

	// Reserved safety zone around snake starting positions
	isSafeZone := func(x, y int) bool {
		return x >= centerX-4 && x <= centerX+2 && y >= centerY-2 && y <= centerY+2
	}

	addObstacle := func(pt Point) {
		if !isSafeZone(pt.X, pt.Y) && pt.X >= 0 && pt.X < width && pt.Y >= 0 && pt.Y < height {
			obstacles[pt] = true
		}
	}

	// Layout patterns based on seed
	pattern := rng.IntN(4)

	switch pattern {
	case 0:
		// Symmetric pillars
		margin := 4
		pillars := []Point{
			{X: margin, Y: margin},
			{X: width - margin - 1, Y: margin},
			{X: margin, Y: height - margin - 1},
			{X: width - margin - 1, Y: height - margin - 1},
		}
		for _, p := range pillars {
			for dx := 0; dx <= 1; dx++ {
				for dy := 0; dy <= 1; dy++ {
					addObstacle(Point{X: p.X + dx, Y: p.Y + dy})
				}
			}
		}

	case 1:
		// Horizontal barrier slots
		y1 := height / 4
		y2 := (height * 3) / 4
		for x := 3; x < width-3; x++ {
			if x != width/2 && x != width/2-1 && x != width/2+1 {
				addObstacle(Point{X: x, Y: y1})
				addObstacle(Point{X: x, Y: y2})
			}
		}

	case 2:
		// Cross corners
		length := 3
		for i := 0; i < length; i++ {
			addObstacle(Point{X: 2 + i, Y: 2})
			addObstacle(Point{X: 2, Y: 2 + i})

			addObstacle(Point{X: width - 3 - i, Y: 2})
			addObstacle(Point{X: width - 3, Y: 2 + i})

			addObstacle(Point{X: 2 + i, Y: height - 3})
			addObstacle(Point{X: 2, Y: height - 3 - i})

			addObstacle(Point{X: width - 3 - i, Y: height - 3})
			addObstacle(Point{X: width - 3, Y: height - 3 - i})
		}

	case 3:
		// Random scatter clusters
		numClusters := int(float64(width*height) * density / 4)
		if numClusters < 2 {
			numClusters = 2
		}
		spanX := width - 4
		spanY := height - 4
		if spanX <= 0 {
			spanX = 1
		}
		if spanY <= 0 {
			spanY = 1
		}
		for c := 0; c < numClusters; c++ {
			cx := 2 + rng.IntN(spanX)
			cy := 2 + rng.IntN(spanY)
			for i := 0; i < 3; i++ {
				ox := cx + rng.IntN(2)
				oy := cy + rng.IntN(2)
				addObstacle(Point{X: ox, Y: oy})
			}
		}
	}

	return obstacles
}
