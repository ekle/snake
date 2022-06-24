package snake

import (
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
	problem       image.Point
	lastDirection Direction
	lost          bool
}

var (
	colorBG        = color.RGBA{R: 0, G: 0, B: 0, A: 128}
	colorSnake     = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	colorSnakeHead = color.RGBA{R: 0, G: 128, B: 0, A: 255}
	colorFood      = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	colorProblem   = color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

func (g *game) ToImage() image.Image {
	src := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: g.size.Add(image.Pt(2, 2)),
	})
	draw.Draw(src, image.Rectangle{Min: image.Pt(1, 1), Max: g.size.Add(image.Pt(1, 1))}, &image.Uniform{C: colorBG}, image.Point{}, draw.Src)
	if g.rnd == nil {
		return src // game not running
	}
	for _, s := range g.snake {
		// fmt.Println("S", s.X, s.Y)
		src.Set(s.X+1, s.Y+1, colorSnake)
	}
	src.Set(g.snake[0].X+1, g.snake[0].Y+1, colorSnakeHead)
	// fmt.Println("F", g.food.X, g.food.Y)
	src.Set(g.food.X+1, g.food.Y+1, colorFood)
	src.Set(g.problem.X+1, g.problem.Y+1, colorProblem)
	// resize.Re
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X*10, src.Bounds().Max.Y*10))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
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
		problem:       image.Pt(-100, -100),
	}
	g.placeFood()
	return g
}

func (g *game) At(x, y int) float64 {
	head := g.snake[0]
	if head.X == x && head.Y == y {
		return 0.5
	}
	if g.food.X == x && g.food.Y == y {
		return 1
	}
	for _, p := range g.snake {
		if p.X == x && p.Y == y {
			return 0
		}
	}
	return -1
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
		g.problem = try
		g.lost = true
		return ErrLost // hit wall
	}
	for _, s := range g.snake {
		if try.Eq(s) {
			g.problem = try
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

func (g *game) GetNeck() image.Point {
	return g.snake[1]
}

func (g *game) GetTail() image.Point {
	return g.snake[len(g.snake)-1]
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
		// fmt.Println("placed food", newFood)
		g.food = newFood
		return
	}
}
