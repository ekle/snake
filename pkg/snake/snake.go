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
	NextFood() image.Point
	GetHead() image.Point
	GetLength() int
	GetSnake() []image.Point
	ToImage() image.Image
}

var ErrLost = errors.New("lost")

// var ErrDirectionInvalid = errors.New("direction Invalid, using last direction")