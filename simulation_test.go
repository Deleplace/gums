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
		g := Minimax{depth: 4}
		r := Minimax{depth: 4}
		rand.Seed(42)
		finalState = fight(g, r)
	}
	fmt.Println(finalState)
	fmt.Println(finalState.Score())
}

func TestDepth(t *testing.T) {
	rand.Seed(42)
	g := Minimax{depth: 4, eval: newEvaluator(&coef_0)}
	r := Minimax{depth: 4, eval: newEvaluator(&coef_1)}

	K := 50
	ng, nr := 0, 0
	for k := 0; k < K; k++ {
		finalState := fight(g, r)
		a, b := finalState.Score()
		if a > b {
			ng++
		}
		if b > a {
			nr++
		}
	}
	fmt.Println(ng, nr)

	ng, nr = 0, 0
	for k := 0; k < K; k++ {
		finalState := fight(r, g)
		a, b := finalState.Score()
		if a > b {
			ng++
		}
		if b > a {
			nr++
		}
	}
	fmt.Println(ng, nr)
}
