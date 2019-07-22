package gums

import (
	"math/rand"
	"strings"
)

type Content int

const (
	Empty Content = iota
	Green
	Red
)

type Player = Content

// The grid is 8x8 == 64 cells
const W = 8

type State [W][W]Content

func (s State) CanPlay(player Player, pos position) bool {
	if player != Green && player != Red {
		panic(player)
	}
	if s[pos.i][pos.j] != Empty {
		return false
	}
	for _, dir := range directions {
		if s.CanCapture(player, pos, dir) {
			return true
		}
	}
	return false
}

func (s State) CanCapture(player Player, pos position, dir direction) bool {
	adjpos, ok := pos.move(dir)
	if !ok {
		return false
	}
	if s[adjpos.i][adjpos.j] != player.Opponent() {
		return false
	}

	nextpos, ok := adjpos.move(dir)
	if !ok {
		return false
	}
	switch s[nextpos.i][nextpos.j] {
	case Empty:
		return false
	case player:
		return true
	case player.Opponent():
		return s.CanCapture(player, adjpos, dir)
	default:
		panic("ouch")
	}
}

func (s State) Capture(player Player, pos position, dir direction) (State, int) {
	adjpos, ok := pos.move(dir)
	if !ok {
		panic("bug")
	}
	switch s[adjpos.i][adjpos.j] {
	default:
		panic("bug")
	case player:
		// Sandwich closed
		return s, 0
	case Empty:
		// No sandwich??
		panic("bug")
	case player.Opponent():
		rec, n := s.Capture(player, adjpos, dir)
		rec[adjpos.i][adjpos.j] = player
		return rec, 1 + n
	}
}

func (s State) Play(player Player, pos position) State {
	t := s
	t[pos.i][pos.j] = player

	for _, dir := range directions {
		if s.CanCapture(player, pos, dir) {
			t, _ = t.Capture(player, pos, dir)
		}
	}
	return t
}

func (s State) PossibleMoves(player Player) (options []position) {
	for i, row := range s {
		for j := range row {
			pos := makepos(i, j)
			if s.CanPlay(player, pos) {
				options = append(options, pos)
			}
		}
	}

	// Shuffle :)
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

type position struct{ i, j int }

func makepos(i, j int) position {
	return position{
		i: i,
		j: j,
	}
}

type direction struct{ di, dj int }

var directions = []direction{
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
	{-1, -1},
	{0, -1},
	{1, -1},
}

func (pos position) move(dir direction) (position, bool) {
	newpos := position{
		i: pos.i + dir.di,
		j: pos.j + dir.dj,
	}
	ok := newpos.i >= 0 && newpos.i < W &&
		newpos.j >= 0 && newpos.j < W
	return newpos, ok
}

func (player Player) Opponent() Player {
	return 3 - player
}

func InitialState() State {
	return State{
		3: [W]Content{3: Green, 4: Red},
		4: [W]Content{3: Red, 4: Green},
	}
}

func (c Content) String() string {
	return []string{
		".",
		"G",
		"R",
	}[c]
}

func (s State) String() string {
	var sb strings.Builder
	for _, row := range s {
		for _, c := range row {
			sb.WriteString(c.String() + " ")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (s State) Score() (green int, red int) {
	var n [3]int
	for _, row := range s {
		for _, c := range row {
			n[c]++
		}
	}
	if total := n[Empty] + n[Green] + n[Red]; total != W*W {
		panic(total)
	}
	return n[Green], n[Red]
}
