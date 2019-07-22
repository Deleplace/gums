package gums

// 2 opponents, using minimax.

func simulate(depth int) *State {
	s := InitialState()
	currentPlayer := Green
	for k := 0; k < 64; k++ {
		canMove, move, _ := Choose(s, currentPlayer, depth)
		if canMove {
			s.Play(currentPlayer, move)
		} else {
			// TODO: detect when neither can move?
		}
		currentPlayer = currentPlayer.Opponent()
	}
	return s
}
