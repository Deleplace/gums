package gums

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkSimulate(b *testing.B) {
	var finalState State
	for i := 0; i < b.N; i++ {
		// Games are reproducible, given the seed.
		rand.Seed(42)
		finalState = simulate(4)
	}
	fmt.Println(finalState)
	fmt.Println(finalState.Score())
}
