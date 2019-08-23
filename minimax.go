package gums

import (
	"math"
)

type Minimax struct {
	depth int
	eval  evaluator
}

func NewMinimax(depth int) Minimax {
	return Minimax{
		depth: depth,
		eval:  newEvaluator(&coef_1),
	}
}

// Choose returns the best move *for this player*, and whether a move *by this player* is possible at all.
func (mm Minimax) Choose(color PlayerColor, s State) (canMove bool, move Position) {
	ok, move, _ := mm.Dig(s, color, mm.depth)
	return ok, move
}

// Dig returns the best move *for this player color*, the desirability of the outcome
// *for this player*, and whether a move *by this player* is possible at all.
// The desirability is computed even if this player can't move
func (mm Minimax) Dig(s State, player PlayerColor, depth int) (canMove bool, move Position, desirability float64) {
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
			return false, Position{}, desirability
		case depth == 1:
			// End of recursion. Let's provide an estimate.
			desi := Desirability(s, player, mm.eval)
			return false, Position{}, desi
		default:
			// The opponent may move
			oppoCanMove, _, oppoDesi := mm.Dig(s, player.Opponent(), depth-1)
			if !oppoCanMove {
				panic("inconsistent")
			}
			desi := -oppoDesi
			return false, Position{}, desi
		}
	}

	bestdesi := minusInfinity
	bestmove := Position{}
	ok := false
	for _, move := range moves {
		t := s.Play(player, move)
		if depth == 1 {
			desi := Desirability(t, player, mm.eval)
			if !ok || desi > bestdesi {
				bestdesi = desi
				bestmove = move
				ok = true
			}
		} else {
			_, _, oppodesi := mm.Dig(t, player.Opponent(), depth-1)
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
