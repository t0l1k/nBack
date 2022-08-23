package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GridView struct {
	grid                     *Point
	rect                     *Rect
	Image                    *ebiten.Image
	Dirty, Visible, DrawRect bool
	bg, fg                   color.Color
}

func NewGridView(rect []int, grid *Point, bg, fg color.Color) *GridView {
	return &GridView{
		grid:     grid,
		rect:     NewRect(rect),
		Image:    nil,
		Dirty:    true,
		Visible:  false,
		DrawRect: false,
		bg:       bg,
		fg:       fg}
}

func (g *GridView) Layout() {
	w, h := g.rect.Size()
	if g.Image == nil {
		g.Image = ebiten.NewImage(w, h)
	} else {
		g.Image.Clear()
	}
	gridX := int(float64(g.rect.W) / g.grid.X)
	gridY := int(float64(g.rect.H) / g.grid.Y)
	fmt.Println(g.rect, g.grid, gridX, gridY)
	if g.DrawRect {
		ebitenutil.DrawRect(g.Image, 0, 0, float64(w), float64(h), g.fg)
		ebitenutil.DrawRect(g.Image, 2, 2, float64(w)-4, float64(h)-4, g.bg)
	}
	for y := gridY; y <= g.rect.H-gridY; y += gridY {
		ebitenutil.DrawLine(g.Image, 0, float64(y), float64(g.rect.W), float64(y), g.fg)
	}
	for x := gridX; x <= g.rect.W-gridX; x += gridX {
		ebitenutil.DrawLine(g.Image, float64(x), 0, float64(x), float64(g.rect.H), g.fg)
	}
	g.Dirty = false
}
func (g *GridView) SetGrid(grid *Point) {
	g.grid = grid
	g.Dirty = true
}

func (g *GridView) Update(dt int) {}
func (g *GridView) Draw(surface *ebiten.Image) {
	if g.Dirty {
		g.Layout()
	}
	if g.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := g.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(g.Image, op)
	}
}

func (g *GridView) Resize(rect []int) {
	g.rect = NewRect(rect)
	g.Dirty = true
	g.Image = nil
}

func (g *GridView) Close() {
	g.Image.Dispose()
}
