package engine

import (
	"math/rand/v2"
)

// PRNG wraps math/rand/v2 with a deterministic seed.
type PRNG struct {
	r *rand.Rand
}

// NewPRNG creates a deterministic PRNG from an int64 seed.
func NewPRNG(seed int64) *PRNG {
	// PCG needs two uint64 seed values. We use seed and seed XOR gold ratio.
	s1 := uint64(seed)
	s2 := uint64(seed) ^ 0x9e3779b97f4a7c15
	pcg := rand.NewPCG(s1, s2)
	return &PRNG{
		r: rand.New(pcg),
	}
}

// IntN returns a non-negative pseudo-random number in [0, n).
func (p *PRNG) IntN(n int) int {
	if n <= 0 {
		return 0
	}
	return p.r.IntN(n)
}

// Float64 returns a pseudo-random number in [0.0, 1.0).
func (p *PRNG) Float64() float64 {
	return p.r.Float64()
}
