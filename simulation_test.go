package gums

import (
	"math/rand"
	"testing"
)

func BenchmarkSimulate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Games are reproducible, given the seed.
		rand.Seed(42)
		finalState := simulate(4)
		_ = finalState
		// fmt.Println(finalState.Score())
	}
}
