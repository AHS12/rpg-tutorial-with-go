package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	X, Y, Dx, Dy float64 //Dx(delta x) change in x, Dy(delta y) change in Y
	Img  *ebiten.Image
}