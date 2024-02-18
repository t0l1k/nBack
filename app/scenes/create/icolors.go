package create

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app/game"
)

type Colors struct {
	eui.DrawableBase
	bg, fg color.Color
}

func NewColorsBar(bg, fg color.Color) *Colors {
	c := &Colors{bg: bg, fg: fg}
	return c
}

func (c *Colors) Setup() { c.Layout() }

func (c *Colors) Layout() {
	c.SpriteBase.Layout()
	c.Image().Fill(c.bg)
	sz := float32(c.GetRect().W) / float32(len(game.Colors))
	for i, clr := range game.Colors {
		vector.DrawFilledRect(c.Image(), sz*float32(i), 0, float32(sz), float32(c.GetRect().H), clr, true)
	}
	c.Dirty = false
}

func (c *Colors) Resize(r []int) {
	c.Rect(eui.NewRect(r))
	c.SpriteBase.Resize(r)
	c.ImageReset()
}
