package gums

var coef = [W][W]float64{
	{10, 2, 5, 5, 5, 5, 2, 10},
	{2, 1, 1, 1, 1, 1, 1, 2},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{5, 1, 1, 1, 1, 1, 1, 5},
	{2, 1, 1, 1, 1, 1, 1, 2},
	{10, 2, 5, 5, 5, 5, 2, 10},
}

// Eval computes a "desirability" of a given state, for each player.
// This is not the same as the final score, where each cell counts for 1.
func Eval(s State) (green float64, red float64) {
	for i, row := range s {
		for j, c := range row {
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

func Desirability(s State, player Player) float64 {
	g, r := Eval(s)
	switch player {
	case Green:
		return g - r
	case Red:
		return r - g
	default:
		panic(player)
	}
}
