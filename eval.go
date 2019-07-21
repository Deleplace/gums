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

func Eval(s *State) (green float64, red float64) {
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
