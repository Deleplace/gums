package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/Deleplace/gums"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	// Cell width in pixels
	d = 40

	delay = 450
)

func main() {
	// Randomness source.
	// The only randomness use is the shuffling of the possible moves, to make
	// them equiprobable.
	// There is no map iteration (non-deterministic).
	// There is no concurrency (which may be non-deterministic).

	seed := time.Now().UnixNano()
	// Games are reproducible, given the seed.
	fmt.Println("Seed =", seed)
	rand.Seed(420004)

	human := &humanPlayer{}
	players := map[gums.PlayerColor]gums.Player{
		gums.Green: human,
		gums.Red:   gums.NewMinimax(6),
	}
	s := gums.InitialState()
	fmt.Println(s)
	currentColor := gums.Green
	k := 0
	var mPause sync.Mutex
	paused := false
	passed := 0

	drawBoard := func(screen *ebiten.Image) error {
		var doMove bool
		var legalMoves []gums.Position
		mPause.Lock()
		isPaused := paused
		mPause.Unlock()

		switch {
		case isPaused:
			doMove = false
		case k == 64:
			fmt.Println("Board full, end of game")
			mPause.Lock()
			paused = true
			mPause.Unlock()
		case passed == 2:
			fmt.Println("Both passed, end of game")
			mPause.Lock()
			paused = true
			mPause.Unlock()
		case currentColor == gums.Green:
			// Process input, if any
			// Human choice
			legalMoves = s.PossibleMoves(currentColor)
			switch {
			case len(legalMoves) == 0:
				human.pass = true
				doMove = true
			case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
				x, y := ebiten.CursorPosition()
				i, j := y/d, x/d
				fmt.Printf("Clicked cell (%d, %d) \n", i, j)
				choice := gums.Position{I: i, J: j}
				in := func(move gums.Position, moves []gums.Position) bool {
					for _, x := range moves {
						if x == move {
							return true
						}
					}
					return false
				}
				if in(choice, legalMoves) {
					//human.choice <- choice
					human.last = choice
					doMove = true
				} else {
					fmt.Println("***\n*** Illegal choice", choice, "\n***")
					doMove = false
				}
			default:
				// No click yet
				doMove = false
			}
		case currentColor == gums.Red:
			// Red AI
			doMove = true
		default:
			panic("no")
		}

		// Logic
		if doMove {
			player := players[currentColor]
			canMove, move := player.Choose(currentColor, s)
			if canMove {
				fmt.Println("Player", currentColor, "plays", move)
				t := s.Play(currentColor, move)
				fmt.Println()
				fmt.Println(t)
				fmt.Println(t.Score())
				s = t
				passed = 0
			} else {
				fmt.Println("Player", currentColor, "passes")
				passed++
			}

			if currentColor == gums.Green {
				fmt.Println("Pausing after human move")
				mPause.Lock()
				paused = true
				mPause.Unlock()
				time.AfterFunc(delay*time.Millisecond, func() {
					mPause.Lock()
					paused = false
					mPause.Unlock()
				})
			}

			currentColor = currentColor.Opponent()
			k++
		}

		// Render
		// fmt.Println("Rendering board")
		screen.Fill(darkGreen)
		for m := 0; m <= gums.W; m++ {
			drawLineH(screen, m, 0, gums.W, lightStroke)
			drawLineV(screen, 0, m, gums.W, lightStroke)
		}
		for i := 0; i < gums.W; i++ {
			for j := 0; j < gums.W; j++ {
				switch s.At(i, j) {
				case gums.Empty:
				case gums.Green:
					drawDisc(screen, i, j, uGreen)
				case gums.Red:
					drawDisc(screen, i, j, uRed)
				}
			}
		}
		for _, p := range legalMoves {
			drawDisc(screen, p.I, p.J, uLightGreen)
		}

		return nil
	}

	ebiten.Run(drawBoard,
		d*gums.W+1,
		d*gums.W+1,
		1.0,
		"Gums")

	g, r := s.Score()
	fmt.Printf("\nFinal score %d %d \n", g, r)
}

type humanPlayer struct {
	last gums.Position
	pass bool
}

func (h *humanPlayer) Choose(c gums.PlayerColor, s gums.State) (canMove bool, move gums.Position) {
	if h.pass {
		h.pass = false
		return false, gums.Position{}
	}
	return true, h.last
}

func drawLineH(screen *ebiten.Image, i, j, w int, c color.Color) {
	for x, y, last := d*j, d*i, d*(j+w); x <= last; x++ {
		screen.Set(x, y, c)
	}
}

func drawLineV(screen *ebiten.Image, i, j, w int, c color.Color) {
	for x, y, last := d*j, d*i, d*(i+w); y <= last; y++ {
		screen.Set(x, y, c)
	}
}

func drawDisc(screen *ebiten.Image, i, j int, src *ebiten.Image) {
	screen.DrawTriangles(
		[]ebiten.Vertex{
			{
				DstX:   float32(d*j + 3),
				DstY:   float32(d*i + 3),
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   float32(d*j + 3),
				DstY:   float32(d*(i+1) - 3),
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   float32(d*(j+1) - 3),
				DstY:   float32(d*i + 3),
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   float32(d*(j+1) - 3),
				DstY:   float32(d*(i+1) - 3),
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
		},
		[]uint16{0, 1, 2, 1, 2, 3},
		src,
		nil,
	)
}

func uniform(c color.Color) *ebiten.Image {
	src, err := ebiten.NewImage(1, 1, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	src.Fill(c)
	return src
}

var (
	green       = color.RGBA{0xee, 0xee, 0xee, 0xff}
	darkGreen   = color.RGBA{0, 0x7f, 0, 0xff}
	lightGreen  = color.RGBA{0x11, 0x99, 0x11, 0xff}
	red         = color.RGBA{0x11, 0x11, 0x11, 0xff}
	lightStroke = color.RGBA{0x20, 0x40, 0x20, 0xff}

	uGreen      = uniform(green)
	uLightGreen = uniform(lightGreen)
	uRed        = uniform(red)
)
