package gums

import (
	"math"
)

// Choose returns the best move *for this player*, the desirability of the outcome
// *for this player*, and whether a move *by this player* is possible at all.
// The desirability is computed even if this player can't move
func Choose(s *State, player Player, depth int) (canMove bool, move position, desirability float64) {
	if depth < 1 {
		panic(depth)
	}

	moves := s.PossibleMoves(player)
	if len(moves) == 0 {
		opponentMoves := s.PossibleMoves(player.Opponent())
		switch {
		case len(opponentMoves) == 0:
			// Game over!
			// Who won, for real?
			g, r := s.Score()
			switch {
			case g == r:
				desirability = 0
			case g < r:
				desirability = minusInfinity
			case g > r:
				desirability = infinity
			}
			if player == Red {
				desirability = -desirability
			}
			return false, position{}, desirability
		case depth == 1:
			// End of recursion. Let's provide an estimate.
			desi := Desirability(s, player)
			return false, position{}, desi
		default:
			// The opponent may move
			oppoCanMove, _, oppoDesi := Choose(s, player.Opponent(), depth-1)
			if !oppoCanMove {
				panic("inconsistent")
			}
			desi := -oppoDesi
			return false, position{}, desi
		}
	}

	bestdesi := minusInfinity
	bestmove := position{}
	ok := false
	for _, move := range moves {
		t := *s
		t.Play(player, move)
		if depth == 1 {
			desi := Desirability(&t, player)
			if !ok || desi > bestdesi {
				bestdesi = desi
				bestmove = move
				ok = true
			}
		} else {
			_, _, oppodesi := Choose(&t, player.Opponent(), depth-1)
			desi := -oppodesi
			if !ok || desi > bestdesi {
				bestdesi = desi
				bestmove = move
				ok = true
			}
		}
	}
	if !ok {
		panic("inconsistent")
	}
	return ok, bestmove, bestdesi
}

var infinity, minusInfinity = math.Inf(1), math.Inf(-1)
