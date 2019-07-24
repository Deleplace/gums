package gums

var coef_0 = [W][W]float64{
	{50, -5, 5, 5, 5, 5, -5, 50},
	{-5, -15, 1, 1, 1, 1, -15, -5},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{-5, -15, 1, 1, 1, 1, -15, -5},
	{50, -5, 5, 5, 5, 5, -5, 50},
}

var coef_1 = [W][W]float64{
	{50, -5, 5, 5, 5, 5, -5, 50},
	{-5, -15, -4, -4, -4, -4, -15, -5},
	{5, -4, 1, 1, 1, 1, -4, 5},
	{5, -4, 1, 1, 1, 1, -4, 5},
	{5, -4, 1, 1, 1, 1, -4, 5},
	{5, -4, 1, 1, 1, 1, -4, 5},
	{-5, -15, -4, -4, -4, -4, -15, -5},
	{50, -5, 5, 5, 5, 5, -5, 50},
}

// Eval computes a "desirability" of a given state, for each player.
// This is not the same as the final score, where each cell counts for 1.
func Eval(s State, coef *[W][W]float64) (green float64, red float64) {
	for i := 0; i < W; i++ {
		for j := 0; j < W; j++ {
			c := s.At(i, j)
			switch c {
			case Green:
				green += coef[i][j]
			case Red:
				red += coef[i][j]
			}
		}
	}
	return
}

func Desirability(s State, player PlayerColor, eval evaluator) float64 {
	g, r := eval(s)
	switch player {
	case Green:
		return g - r
	case Red:
		return r - g
	default:
		panic(player)
	}
}

type evaluator func(s State) (green float64, red float64)

func newEvaluator(coef *[W][W]float64) evaluator {
	return func(s State) (green float64, red float64) {
		return Eval(s, coef)
	}
}
