package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"snake/pkg/snake"
)

var (
	DefaultDirection snake.Direction = snake.DirSouth
	DefaultTick                      = time.Millisecond * 200
)

func main() {
	go func() {
		w := app.NewWindow(app.Title("SNAKE"))

		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var btnReset widget.Clickable
	var ops op.Ops
	var tag = new(bool)

	key.InputOp{Keys: key.NameDownArrow, Tag: tag}.Add(&ops)

	game := snake.NewGame(image.Pt(30, 30), 0)

	th := material.NewTheme(gofont.Collection())

	var direction snake.Direction = DefaultDirection
	var lastDirection snake.Direction = DefaultDirection
	tick := time.NewTicker(DefaultTick)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			err := game.Step(direction)
			if err != nil {
				fmt.Println(err)
			}
			lastDirection = direction
			w.Invalidate()
		case e := <-w.Events():
			// fmt.Printf("\ngot event: %T %s\n", e, e)
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
				// log.Println("system.FrameEvent", len(e.Queue.Events(nil)))
				// x := e.Queue.Events(key.Event{})
				// fmt.Println(x)
				gtx := layout.NewContext(&ops, e)
				// op.InvalidateOp{At: gtx.Now.Add(DefaultTick / 10)}.Add(&ops)

				// Create a clip area the size of the window.
				areaStack := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
				// Register for all pointer inputs on the current clip area.
				// pointer.InputOp{
				//	Types: pointer.Enter | pointer.Leave | pointer.Drag | pointer.Press | pointer.Release | pointer.Scroll | pointer.Move,
				//	Tag:   w, // Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
				// }.Add(gtx.Ops)
				// Register for keyboard input on the current clip area.
				key.InputOp{
					Tag: w, // Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
				}.Add(gtx.Ops)
				key.InputOp{
					Tag:  w, // Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
					Keys: key.Set(strings.Join([]string{key.NameDownArrow, key.NameUpArrow, key.NameLeftArrow, key.NameRightArrow}, "|")),
				}.Add(gtx.Ops)
				// Request keyboard focus to the current clip area.
				key.FocusOp{
					Tag: w, // Focus the input area with our window as the tag.
				}.Add(gtx.Ops)
				// Pop the clip area to finalize it.
				areaStack.Pop()

				for _, event := range gtx.Events(w) {
					// Perform event handling here instead of in the outer type switch.
					// log.Printf("%#+v", event)
					switch k := event.(type) {
					default:
						log.Printf("UNKNOWN key event: %T \n", e)
					case key.Event:
						if k.State != key.Press {
							continue
						}
						switch k.Name {
						case key.NameDownArrow:
							if lastDirection != snake.DirNorth {
								direction = snake.DirSouth
							}
						case key.NameUpArrow:
							if lastDirection != snake.DirSouth {
								direction = snake.DirNorth
							}
						case key.NameLeftArrow:
							if lastDirection != snake.DirEast {
								direction = snake.DirWest
							}
						case key.NameRightArrow:
							if lastDirection != snake.DirWest {
								direction = snake.DirEast
							}
						}
					}
				}

				if btnReset.Clicked() {
					direction = DefaultDirection
					lastDirection = DefaultDirection
					game = snake.NewGame(image.Pt(30, 30), 0)
				}

				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.Button(th, &btnReset, "reset").Layout),

					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return widget.Image{
							Src: paint.NewImageOp(game.ToImage()),
							Fit: widget.Contain,
						}.Layout(gtx)
					}),
				)
				e.Frame(gtx.Ops)
			}
		}
	}
	return nil
}
