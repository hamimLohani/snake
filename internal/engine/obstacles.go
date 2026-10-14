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
					pt := Point{X: p.X + dx, Y: p.Y + dy}
					if !isSafeZone(pt.X, pt.Y) && pt.X < width && pt.Y < height {
						obstacles[pt] = true
					}
				}
			}
		}

	case 1:
		// Horizontal barrier slots
		y1 := height / 4
		y2 := (height * 3) / 4
		for x := 3; x < width-3; x++ {
			if x != width/2 && x != width/2-1 && x != width/2+1 {
				p1 := Point{X: x, Y: y1}
				p2 := Point{X: x, Y: y2}
				if !isSafeZone(p1.X, p1.Y) {
					obstacles[p1] = true
				}
				if !isSafeZone(p2.X, p2.Y) {
					obstacles[p2] = true
				}
			}
		}

	case 2:
		// Cross corners
		length := 3
		for i := 0; i < length; i++ {
			obstacles[Point{X: 2 + i, Y: 2}] = true
			obstacles[Point{X: 2, Y: 2 + i}] = true

			obstacles[Point{X: width - 3 - i, Y: 2}] = true
			obstacles[Point{X: width - 3, Y: 2 + i}] = true

			obstacles[Point{X: 2 + i, Y: height - 3}] = true
			obstacles[Point{X: 2, Y: height - 3 - i}] = true

			obstacles[Point{X: width - 3 - i, Y: height - 3}] = true
			obstacles[Point{X: width - 3, Y: height - 3 - i}] = true
		}

	case 3:
		// Random scatter clusters
		numClusters := int(float64(width*height) * density / 4)
		if numClusters < 2 {
			numClusters = 2
		}
		for c := 0; c < numClusters; c++ {
			cx := 2 + rng.IntN(width-4)
			cy := 2 + rng.IntN(height-4)
			for i := 0; i < 3; i++ {
				ox := cx + rng.IntN(2)
				oy := cy + rng.IntN(2)
				if !isSafeZone(ox, oy) && ox >= 1 && ox < width-1 && oy >= 1 && oy < height-1 {
					obstacles[Point{X: ox, Y: oy}] = true
				}
			}
		}
	}

	return obstacles
}
