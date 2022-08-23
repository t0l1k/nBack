package ui

import "github.com/hajimehoshi/ebiten/v2"

type Drawable interface {
	Layout()
	Update(dt int)
	Draw(surface *ebiten.Image)
	Resize([]int)
	Close()
}
