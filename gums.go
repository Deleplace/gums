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

// A State fits in 128 bits:
// - filled has a 1 if a cell is occupied
// - p has a 0 if occupied by Green, 1 if occupied by Red
type State struct{ filled, p uint64 }

func (s *State) At(i, j int) Content {
	k := flatten(i, j)
	switch {
	case !hasbit(s.filled, k):
		return Empty
	case !hasbit(s.p, k):
		return Green
	case hasbit(s.p, k):
		return Red
	default:
		panic(s)
	}
}

func (s *State) Set(i, j int, c Content) {
	k := flatten(i, j)
	switch c {
	case Empty:
		s.filled = clearbit(s.filled, k)
	case Green:
		s.filled = setbit(s.filled, k)
		s.p = clearbit(s.p, k)
	case Red:
		s.filled = setbit(s.filled, k)
		s.p = setbit(s.p, k)
	default:
		panic(c)
	}
}

func flatten(i, j int) uint {
	return uint(8*i + j)
}

func setbit(x uint64, k uint) uint64 {
	return x | (1 << k)
}

func clearbit(x uint64, k uint) uint64 {
	return x & ^(1 << k)
}

func hasbit(x uint64, k uint) bool {
	return x&(1<<k) != 0
}

func (s *State) CanPlay(player Player, pos position) bool {
	if player != Green && player != Red {
		panic(player)
	}
	if s.At(pos.i, pos.j) != Empty {
		return false
	}
	for _, dir := range directions {
		if s.CanCapture(player, pos, dir) {
			return true
		}
	}
	return false
}

func (s *State) CanCapture(player Player, pos position, dir direction) bool {
	adjpos, ok := pos.move(dir)
	if !ok {
		return false
	}
	if s.At(adjpos.i, adjpos.j) != player.Opponent() {
		return false
	}

	nextpos, ok := adjpos.move(dir)
	if !ok {
		return false
	}
	switch s.At(nextpos.i, nextpos.j) {
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

func (s *State) Capture(player Player, pos position, dir direction) int {
	adjpos, ok := pos.move(dir)
	if !ok {
		panic("bug")
	}
	switch s.At(adjpos.i, adjpos.j) {
	default:
		panic("bug")
	case player:
		// Sandwich closed
		return 0
	case Empty:
		// No sandwich??
		panic("bug")
	case player.Opponent():
		n := s.Capture(player, adjpos, dir)
		s.Set(adjpos.i, adjpos.j, player)
		return 1 + n
	}
}

func (s *State) Play(player Player, pos position) {
	s.Set(pos.i, pos.j, player)

	for _, dir := range directions {
		if s.CanCapture(player, pos, dir) {
			s.Capture(player, pos, dir)
		}
	}
}

func (s *State) PossibleMoves(player Player) (options []position) {
	for i := 0; i < W; i++ {
		for j := 0; j < W; j++ {
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

func InitialState() *State {
	s := &State{}
	s.Set(3, 3, Green)
	s.Set(3, 4, Red)
	s.Set(4, 3, Red)
	s.Set(4, 4, Green)
	return s
	// return State{
	// 	3: [W]Content{3: Green, 4: Red},
	// 	4: [W]Content{3: Red, 4: Green},
	// }
}

func (c Content) String() string {
	return []string{
		".",
		"G",
		"R",
	}[c]
}

func (s *State) String() string {
	var sb strings.Builder
	for i := 0; i < W; i++ {
		for j := 0; j < W; j++ {
			c := s.At(i, j)
			sb.WriteString(c.String() + " ")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (s *State) Score() (green int, red int) {
	var n [3]int
	for i := 0; i < W; i++ {
		for j := 0; j < W; j++ {
			c := s.At(i, j)
			n[c]++
		}
	}
	if total := n[Empty] + n[Green] + n[Red]; total != W*W {
		panic(total)
	}
	return n[Green], n[Red]
}
