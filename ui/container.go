package ui

import "github.com/hajimehoshi/ebiten/v2"

type Container interface {
	Add(item Drawable)
	Update(dt int)
	Draw(surface *ebiten.Image)
}
