package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	X, Y float64
	Img  *ebiten.Image
}