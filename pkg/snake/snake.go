package snake

import (
	"errors"
	"image"
)

type Direction int

const (
	DirNorth = 1
	DirSouth = 2
	DirWest  = 3
	DirEast  = 4
)

type Game interface {
	Step(dir Direction) error
	IdleRounds() int
	GetRounds() int
	NextFood() image.Point
	GetHead() image.Point
	GetNeck() image.Point
	GetTail() image.Point
	GetLength() int
	GetSnake() []image.Point
	ToImage() image.Image
	At(x, y int) float64
}

var ErrLost = errors.New("lost")

// var ErrDirectionInvalid = errors.New("direction Invalid, using last direction")
