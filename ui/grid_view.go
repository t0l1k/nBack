package ui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GridView struct {
	spacing                 int
	rect                    *Rect
	Image                   *ebiten.Image
	Dirty, Visibe, DrawRect bool
	bg, fg                  color.RGBA
}

func NewGridView(rect *Rect, spacing int) *GridView {
	return &GridView{
		spacing:  spacing,
		rect:     rect,
		Image:    nil,
		Dirty:    true,
		Visibe:   false,
		DrawRect: true,
		bg:       color.RGBA{64, 0, 0, 255},
		fg:       color.RGBA{255, 255, 0, 255}}
}

func (g *GridView) Layout() *ebiten.Image {
	if !g.Dirty {
		return g.Image
	}
	w, h := g.rect.GetSize()
	image := ebiten.NewImage(w, h)
	spacing := int(g.rect.GetLowestSize()/g.spacing + 1)
	if g.DrawRect {
		ebitenutil.DrawRect(image, 0, 0, float64(w), float64(h), g.fg)
		ebitenutil.DrawRect(image, 2, 2, float64(w)-4, float64(h)-4, g.bg)
	}
	for y := spacing; y <= g.rect.H; y += spacing {
		ebitenutil.DrawLine(image, 0, float64(y), float64(g.rect.W), float64(y), g.fg)
	}
	for x := spacing; x <= g.rect.W; x += spacing {
		ebitenutil.DrawLine(image, float64(x), 0, float64(x), float64(g.rect.H), g.fg)
	}
	g.Dirty = false
	log.Printf("Grid created. %v %v %v %v", w, h, g.rect, spacing)
	return image
}

func (g *GridView) Update(dt int) {}
func (g *GridView) Draw(surface *ebiten.Image) {
	if g.Dirty {
		g.Image = g.Layout()
	}
	if g.Visibe {
		op := &ebiten.DrawImageOptions{}
		x, y := g.rect.GetPos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(g.Image, op)
	}
}
