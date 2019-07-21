package gums

import (
	"math"
)

func Choose(s *State, player Player, depth int) (position, float64, bool) {
	moves := s.PossibleMoves(player)
	if len(moves) == 0 {
		return position{}, 0, false
	}

	bestdiff := math.Inf(-1)
	bestmove := position{}
	ok := false
	for _, move := range moves {
		t := s.Play(player, move)
		if depth == 1 {
			g, r := Eval(&t)
			diff := g - r
			if player == Red {
				diff = -diff
			}
			if diff > bestdiff {
				bestdiff = diff
				bestmove = move
				ok = true
			}
		} else {
			_, oppodiff, oppook := Choose(&t, player.Opponent(), depth-1)
			if !oppook {
				// Then what? player plays again?
				continue
			}
			if player == Red {
				oppodiff = -oppodiff
			}
			if oppodiff > bestdiff {
				bestdiff = oppodiff
				bestmove = move
				ok = true
			}
		}
	}
	return bestmove, bestdiff, ok
}
