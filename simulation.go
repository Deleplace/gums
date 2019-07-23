package gums

// 2 opponents, using minimax.

func fight(g, r Player) State {
	s := InitialState()
	players := [...]Player{
		Green: g,
		Red:   r,
	}
	currentColor := Green
	currentPlayer := g

	for k := 0; k < 64; k++ {
		canMove, move := currentPlayer.Choose(currentColor, s)
		if canMove {
			t := s.Play(currentColor, move)
			s = t
		} else {
			// TODO: detect when neither can move?
		}
		currentColor = currentColor.Opponent()
		currentPlayer = players[currentColor]
	}
	return s
}
