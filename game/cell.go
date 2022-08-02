package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/t0l1k/nBack/ui"
)

type Cell struct {
	rect                                      *ui.Rect
	Image                                     *ebiten.Image
	Dirty, Visibe, DrawRect, IsCenter, Active bool
	bg, fg, activeColor                       color.Color
	margin                                    float64
}

func NewCell(rect []int, isCenter bool, bg, fg, activeColor color.Color) *Cell {
	return &Cell{
		rect:        ui.NewRect(rect),
		Image:       nil,
		IsCenter:    isCenter,
		Dirty:       true,
		Visibe:      false,
		DrawRect:    true,
		Active:      false,
		margin:      0.05,
		bg:          bg,
		fg:          fg,
		activeColor: activeColor,
	}
}

func (c *Cell) Layout() *ebiten.Image {
	w, h := c.rect.Size()
	image := ebiten.NewImage(w, h)
	bg := c.bg
	if c.Active {
		bg = c.activeColor
	}
	m := float64(w) * c.margin
	if c.DrawRect {
		ebitenutil.DrawRect(image, m, m, float64(w)-m*2, float64(h)-m*2, bg)

	}
	if c.IsCenter {
		m := float64(c.rect.H) * 0.4
		x1 := float64(c.rect.W) / 2
		y1 := m
		x2 := x1
		y2 := float64(c.rect.H) - m
		ebitenutil.DrawLine(image, x1, y1, x2, y2, c.fg)
		ebitenutil.DrawLine(image, y1, x1, y2, x2, c.fg)
	}
	c.Dirty = false
	return image
}

func (c *Cell) SetActiveColor(value color.Color) {
	if c.activeColor == value {
		return
	}
	c.activeColor = value
	c.Dirty = true
}

func (c *Cell) SetActive(value bool) {
	c.Active = value
	c.Dirty = true
}
func (c *Cell) Update(dt int) {}
func (c *Cell) Draw(surface *ebiten.Image) {
	if c.Dirty {
		c.Image = c.Layout()
	}
	if c.Visibe {
		op := &ebiten.DrawImageOptions{}
		x, y := c.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(c.Image, op)
	}
}

func (c *Cell) Resize(rect []int) {
	c.rect = ui.NewRect(rect)
	c.Dirty = true
}

func (c Cell) String() string {
	return fmt.Sprintf("Cell: %v", c.rect)
}
