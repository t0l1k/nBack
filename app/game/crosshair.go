package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
)

type Crosshair struct {
	eui.DrawableBase
	bg, fg            color.Color
	thickness, lenght float32
}

func NewCrosshair(bg, fg color.Color, thickness, lenght float32) *Crosshair {
	c := &Crosshair{bg: bg, fg: fg, thickness: thickness, lenght: lenght}
	return c
}

func (c *Crosshair) Layout() {
	c.SpriteBase.Layout()
	c.Image().Fill(c.bg)
	var (
		x1, y1, x2, y2, m float32
	)
	m = float32(c.GetRect().GetLowestSize()) * c.lenght
	x1 = float32(c.GetRect().W) / 2
	y1 = m
	x2 = float32(c.GetRect().W) / 2
	y2 = float32(c.GetRect().H) - m
	vector.StrokeLine(c.Image(), x1, y1, x2, y2, c.thickness, c.fg, true)
	x1 = m
	y1 = float32(c.GetRect().H) / 2
	x2 = float32(c.GetRect().W) - m
	y2 = float32(c.GetRect().H) / 2
	vector.StrokeLine(c.Image(), x1, y1, x2, y2, c.thickness, c.fg, true)
	c.Dirty = false
}

func (c *Crosshair) Resize(r []int) {
	c.Rect(eui.NewRect(r))
	c.SpriteBase.Resize(r)
	c.ImageReset()
}
