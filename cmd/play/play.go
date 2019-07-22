package main

import (
	"fmt"
	"math/rand"

	"github.com/Deleplace/gums"
)

func main() {
	// Fix randomness source.
	// The only randomness use is the shuffling of the possible moves, to make
	// them equiprobable.
	// There is no map iteration (non-deterministic).
	// There is no concurrency (which may be non-deterministic).
	//
	// Games are reproducible, given the seed.
	rand.Seed(420004)

	s := gums.InitialState()
	fmt.Println(&s)
	currentPlayer := gums.Green

	for k := 0; k < 64; k++ {
		canMove, move, diff := gums.Choose(&s, currentPlayer, 5)
		fmt.Printf("%v's turn: evaluates %.2f %v \n", currentPlayer, diff, canMove)
		if canMove {
			t := (&s).Play(currentPlayer, move)
			fmt.Println()
			fmt.Println(&t)
			fmt.Println(t.Score())
			fmt.Println(gums.Eval(&t))
			s = t
		} else {
			// TODO: detect when neither can move?
		}
		currentPlayer = currentPlayer.Opponent()
	}
	g, r := s.Score()
	fmt.Printf("\nFinal score %d %d \n", g, r)
}
