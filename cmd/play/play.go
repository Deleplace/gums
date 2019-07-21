package main

import (
	"fmt"

	"github.com/Deleplace/gums"
)

func main() {
	s := gums.InitialState()
	fmt.Println(&s)
	/*
		moves := (&s).PossibleMoves(gums.Green)
		for _, move := range moves {
			t := (&s).Play(gums.Green, move)
			fmt.Println()
			fmt.Println(&t)
			fmt.Println(t.Score())
			fmt.Println(gums.Eval(&t))
		}
		fmt.Println()
	*/
	currentPlayer := gums.Green

	for k := 0; k < 30; k++ {
		move, diff, ok := gums.Choose(&s, currentPlayer, 7)
		fmt.Printf("%v's turn: evaluates %.2f %v", currentPlayer, diff, ok)
		t := (&s).Play(currentPlayer, move)
		fmt.Println()
		fmt.Println(&t)
		fmt.Println(t.Score())
		fmt.Println(gums.Eval(&t))
		s = t
		currentPlayer = currentPlayer.Opponent()
	}
}
