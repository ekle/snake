package main

import (
	"errors"
	"image"
	"log"
	"math"
	"os"
	"sort"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"snake/pkg/neuron"
	"snake/pkg/snake"
)

const MaxGames = 2000
const MaxShow = 25
const maxRounds = 1000

type RunningGame struct {
	net    *neuron.Net
	game   snake.Game
	out    []neuron.Neuron
	rounds int
}

func createNetWithGame() RunningGame {
	size := 5
	game := snake.NewGame(image.Pt(size, size), int64(time.Now().UnixNano()))
	input := []*neuron.Input{
		// neuron.NewInput(func() float64 { return float64(size) / float64(size) }),                  // size X
		// neuron.NewInput(func() float64 { return float64(size) / float64(size) }),                  // size Y
		// neuron.NewInput(func() float64 { return float64(game.NextFood().X) / float64(size) }), // food X
		// neuron.NewInput(func() float64 { return float64(game.NextFood().Y) / float64(size) }), // food Y
		// neuron.NewInput(func() float64 { return float64(game.GetHead().X) / float64(size) }), // head X
		// neuron.NewInput(func() float64 { return float64(game.GetHead().Y) / float64(size) }), // head Y
		// neuron.NewInput(func() float64 { return float64(game.GetNeck().X) / float64(size) }),  // head X
		// neuron.NewInput(func() float64 { return float64(game.GetNeck().Y) / float64(size) }),  // head Y
		// neuron.NewInput(func() float64 { return float64(game.GetTail().X) / float64(size) }),  // head X
		// neuron.NewInput(func() float64 { return float64(game.GetTail().Y) / float64(size) }),  // head Y
		neuron.NewInput(func() float64 { return float64(game.GetLength()) }), // length of snake
		// neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.GetNeck()).X) / float64(size) }),
		// neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.GetNeck()).Y) / float64(size) }),
		// neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.NextFood()).X) / float64(size) }),
		// neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.NextFood()).Y) / float64(size) }),
		// neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.GetTail()).X) / float64(size) }),
		neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.GetNeck()).X) }),
		neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.GetNeck()).Y) }),
		neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.NextFood()).X) }),
		neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.NextFood()).Y) }),
		neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.GetTail()).X) }),
		neuron.NewInput(func() float64 { return float64(game.GetHead().Sub(game.GetTail()).Y) }),
	}
	input = nil
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			i := i
			j := j
			input = append(input, neuron.NewInput(func() float64 {
				return game.At(i, j)
				// img := game.ToImage().At(i, j)
				// a, b, c, d := img.RGBA()
				// return float64(a * b * c * d)
			}))
		}
	}

	net, out := neuron.NewNet(input, len(input)*2, len(input)*2, 4)
	// net, out := neuron.NewNet(input, 4)
	return RunningGame{net: net, game: game, out: out}
}

func main() {

	go func() {
		w := app.NewWindow(app.Title("SNAKE AI"))

		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	go learn()
	app.Main()
}

var refresh = make(chan struct{})

func loop(w *app.Window) error {

	// var btnReset widget.Clickable

	var ops op.Ops
	var games = make([]RunningGame, MaxShow)
	for {
		select {
		case <-refresh:
			w.Invalidate()
		case e := <-w.Events():
			switch e := e.(type) {
			default:
				log.Printf("UNKNOWN event: %T \n", e)
			case system.DestroyEvent:
				return e.Err
			case pointer.Event:
				log.Println("pointer.Event")
			case key.Event:
				log.Println("key.Event")
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				op.InvalidateOp{At: gtx.Now.Add(time.Millisecond * 100)}.Add(&ops)
				var fc []layout.FlexChild

				var rows []layout.FlexChild

				for k, g := range games {
					if g.net == nil {
						if len(results) < (k + 1) {
							continue
						}
						running := createNetWithGame()
						if err := running.net.Load(results[k].net); err != nil {
							log.Fatal(err)
						}
						games[k] = running
						g = running
					}
				}

				perLine := int(math.Sqrt(float64(MaxShow)))

				for k, g := range games {
					if g.net == nil {
						continue
					}
					g := g // local copy
					err := oneStep(&g)
					if err != nil || g.rounds > (g.game.GetLength()*10) {
						games[k].net = nil // reset
					}
					fc = append(fc, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return widget.Image{
							Src: paint.NewImageOp(g.game.ToImage()),
							Fit: widget.Contain,
						}.Layout(gtx)
					}))
					if len(fc) >= perLine {
						lfc := fc
						rows = append(rows, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, lfc...)
						}))
						fc = nil
					}
				}
				if len(fc) > 0 {
					lfc := fc
					rows = append(rows, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, lfc...)
					}))
					fc = nil
				}

				layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
				e.Frame(gtx.Ops)

			}

		}
	}
}

// var games = make([]snake.Game, MaxGames)

func learn() {
	// add initial net
	initial := createNetWithGame()
	rawNet, err := initial.net.Save()
	if err != nil {
		log.Fatal(err)
	}
	results = append(results, result{
		net:  rawNet,
		sort: 0,
	})
	results = append(results, result{
		net:  rawNet,
		sort: 1,
	})
	generation := 0
	for {
		generation++
		sort.Slice(results, func(i, j int) bool {
			return results[i].sort > results[j].sort
		})
		if len(results) > (MaxShow * 2) {
			results = results[:MaxShow*2]
		}

		n := 0

		best := results[0]
		log.Printf("starting generation %d finished. best is %d atm with sort %d\n", generation, best.length, best.sort)
		var wg sync.WaitGroup
		for i := 0; i < MaxGames; i++ {
			n++
			if n >= len(results)/10 {
				n = 0
			}
			best := results[n]
			running := createNetWithGame()

			if err := running.net.Load(best.net); err != nil {
				log.Fatal(err)
			}

			if i >= 1 { // do not randomize first
				x := (float64(1) / (MaxGames * 5)) * float64(i)
				running.net.Randomize(x)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				res := play(running)
				results = append(results, res)
			}()
		}
		wg.Wait()
	}
}

func play(r RunningGame) result {
	var err error
	for err == nil {
		err = oneStep(&r)
	}
	sort := r.game.GetLength()
	if errors.Is(err, ErrOutOfRounds) {
		sort--
	}
	sort *= 1000
	sort -= r.rounds
	rawNet, err := r.net.Save()
	if err != nil {
		panic("could not save net")
	}
	return result{
		net: rawNet,
		// sort:   (r.game.GetLength() * 1000000) - r.rounds,
		sort:   sort,
		length: r.game.GetLength(),
	}
}

var ErrOutOfRounds = errors.New("out of rounds")

func oneStep(r *RunningGame) error {
	r.rounds++
	if r.rounds > maxRounds {
		// fmt.Print(r.rounds, maxRounds)
		return ErrOutOfRounds
	}
	r.net.Calc()
	dir := toDirection(r.out)
	return r.game.Step(dir)
}

func toDirection(out []neuron.Neuron) snake.Direction {
	var dir snake.Direction
	var min float64 = 0
	for k, v := range out {
		if c := v.Read(); c > min {
			min = c
			switch k {
			case 0:
				dir = snake.DirNorth
			case 1:
				dir = snake.DirSouth
			case 2:
				dir = snake.DirWest
			case 3:
				dir = snake.DirEast
			default:
				panic("cannot happen")
			}
		}
	}
	return dir
}

type result struct {
	net    []byte
	sort   int
	length int
}

var results []result
