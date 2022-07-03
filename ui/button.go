package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	rect                  *Rect
	Image                 *ebiten.Image
	Dirty, Visible, focus bool
	bg, fg                color.Color
	text                  string
	mouseDown             bool
	onPressed             func(b *Button)
	lbl                   *Label
}

func NewButton(text string, rect []int, bg, fg color.Color, f func(b *Button)) *Button {
	return &Button{
		text:      text,
		rect:      NewRect(rect),
		Image:     nil,
		Dirty:     true,
		Visible:   true,
		focus:     false,
		mouseDown: false,
		bg:        bg,
		fg:        fg,
		onPressed: f,
		lbl:       NewLabel(text, rect, bg, fg),
	}
}

func (b *Button) Layout() *ebiten.Image {
	if !b.Dirty {
		return b.Image
	}
	w, h := b.rect.Size()
	image := ebiten.NewImage(w, h)
	b.Dirty = false
	return image
}

func (b *Button) Update(dt int) {
	x, y := ebiten.CursorPosition()
	if b.rect.InRect(x, y) {
		b.focus = true
		b.lbl.SetBg(b.fg)
		b.lbl.SetFg(b.bg)
	} else {
		b.focus = false
		b.lbl.SetBg(b.bg)
		b.lbl.SetFg(b.fg)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if b.focus {
			b.mouseDown = true
			b.lbl.SetRect(true)
		} else {
			b.mouseDown = false
			b.lbl.SetRect(false)
		}
	} else {
		if b.mouseDown {
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.mouseDown = false
		b.lbl.SetRect(false)
	}
	b.lbl.Update(dt)
}
func (b *Button) Draw(surface *ebiten.Image) {
	if b.Dirty {
		b.Image = b.Layout()
	}
	if b.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := b.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(b.Image, op)
	}
	b.lbl.Draw(surface)
}

func (b *Button) Resize(rect []int) {
	b.rect = NewRect(rect)
	b.lbl.Resize(rect)
	b.Dirty = true
}
