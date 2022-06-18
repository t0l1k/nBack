package ui

import "github.com/hajimehoshi/ebiten/v2"

type Drawable interface {
	Layout() *ebiten.Image
	Update(dt int)
	Draw(surface *ebiten.Image)
}
