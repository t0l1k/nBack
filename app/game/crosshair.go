package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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
	w0, h0 := c.Rect.Size()
	if c.Image == nil {
		c.Image = ebiten.NewImage(w0, h0)
	} else {
		c.Image.Clear()
	}
	c.Image.Fill(c.bg)
	var (
		x1, y1, x2, y2, m float32
	)
	m = float32(c.Rect.GetLowestSize()) * c.lenght
	x1 = float32(c.Rect.W) / 2
	y1 = m
	x2 = x1
	y2 = float32(c.Rect.H) - m
	vector.StrokeLine(c.Image, x1, y1, x2, y2, c.thickness, c.fg, true)
	vector.StrokeLine(c.Image, y1, x1, y2, x2, c.thickness, c.fg, true)
	c.Dirty = false
}

func (c *Crosshair) Resize(r []int) {
	c.Rect = eui.NewRect(r)
	c.Dirty = true
}
