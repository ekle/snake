package snake

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"golang.org/x/image/draw"
)

var _ Game = (*game)(nil)

type game struct {
	rnd           *rand.Rand
	size          image.Point
	snake         []image.Point
	food          image.Point
	lastDirection Direction
	lost          bool
}

var (
	colorBG        = color.RGBA{R: 0, G: 0, B: 0, A: 128}
	colorSnake     = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	colorSnakeHead = color.RGBA{R: 0, G: 128, B: 0, A: 255}
	colorFood      = color.RGBA{R: 255, G: 0, B: 0, A: 255}
)

func (g *game) ToImage() image.Image {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: g.size,
	})
	draw.Draw(img, img.Bounds(), &image.Uniform{C: colorBG}, image.Point{}, draw.Src)
	for _, s := range g.snake {
		// fmt.Println("S", s.X, s.Y)
		img.Set(s.X, s.Y, colorSnake)
	}
	img.Set(g.snake[0].X, g.snake[0].Y, colorSnakeHead)
	// fmt.Println("F", g.food.X, g.food.Y)
	img.Set(g.food.X, g.food.Y, colorFood)
	return img
}

func NewGame(size image.Point, seed int64) Game {
	g := &game{
		snake: []image.Point{
			image.Pt(size.X/2, size.Y/2+1),
			image.Pt(size.X/2, size.Y/2),
		},
		size:          size,
		rnd:           rand.New(rand.NewSource(seed)),
		lastDirection: DirSouth,
	}
	g.placeFood()
	return g
}
func (g *game) Step(dir Direction) error {
	if g.lost {
		return ErrLost
	}
	try := g.snake[0]
	// fmt.Println("old:", try)
	switch dir {
	case DirNorth:
		try.Y -= 1
	case DirSouth:
		try.Y += 1
	case DirWest:
		try.X -= 1
	case DirEast:
		try.X += 1
	}
	if !try.In(image.Rectangle{Max: g.size}) {
		g.lost = true
		return ErrLost // hit wall
	}
	for _, s := range g.snake {
		if try.Eq(s) {
			g.lost = true
			return ErrLost // hit self
		}
	}
	g.snake = append([]image.Point{try}, g.snake...)

	if try.Eq(g.food) {
		g.placeFood()
		return nil
	}
	g.snake = g.snake[:len(g.snake)-1] // shorten snake
	return nil
}

func (g *game) NextFood() image.Point {
	return g.food
}

func (g *game) GetHead() image.Point {
	return g.snake[0]
}

func (g *game) GetLength() int {
	return len(g.snake)
}

func (g *game) GetSnake() []image.Point {
	return g.snake
}

func (g *game) placeFood() {
outer:
	for {
		newFood := image.Point{
			X: g.rnd.Intn(g.size.X),
			Y: g.rnd.Intn(g.size.Y),
		}
		for _, s := range g.snake {
			if s.Eq(newFood) {
				continue outer
			}
		}
		fmt.Println("placed food", newFood)
		g.food = newFood
		return
	}
}
