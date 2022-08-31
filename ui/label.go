package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Label struct {
	text                     string
	rect                     *Rect
	Image                    *ebiten.Image
	Dirty, Visible, drawRect bool
	bg, fg                   color.Color
}

func NewLabel(text string, rect []int, bg, fg color.Color) *Label {
	return &Label{
		text:     text,
		rect:     NewRect(rect),
		Image:    nil,
		Dirty:    true,
		Visible:  true,
		drawRect: false,
		bg:       bg,
		fg:       fg,
	}
}

func (l *Label) SetBg(value color.Color) {
	if l.bg == value {
		return
	}
	l.bg = value
	l.Dirty = true
}

func (l *Label) SetRect(value bool) {
	if l.drawRect == value {
		return
	}
	l.drawRect = value
	l.Dirty = true
}

func (l *Label) SetFg(value color.Color) {
	if l.fg == value {
		return
	}
	l.fg = value
	l.Dirty = true
}

func (l *Label) SetText(value string) {
	if l.text == value {
		return
	}
	l.text = value
	l.Dirty = true
}

func (l *Label) Layout() {
	w, h := l.rect.Size()
	if l.Image == nil {
		l.Image = ebiten.NewImage(w, h)
	} else {
		l.Image.Clear()
	}
	l.Image.Fill(l.bg)
	if l.drawRect {
		ebitenutil.DrawRect(l.Image, 0, 0, float64(w), float64(h), l.fg)
		ebitenutil.DrawRect(l.Image, 2, 2, float64(w)-4, float64(h)-4, l.bg)
	}
	fntSize := GetFonts().calcFontSize(l.text, l.rect)
	font := GetFonts().get(fntSize)
	b := text.BoundString(font, l.text)
	x := (l.rect.W - b.Max.X) / 2
	y := l.rect.H - (l.rect.H+b.Min.Y)/2
	text.Draw(l.Image, l.text, font, x, y, l.fg)
	l.Dirty = false
}

func (l *Label) Update(dt int) {}
func (l *Label) Draw(surface *ebiten.Image) {
	if l.Dirty {
		l.Layout()
	}
	if l.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := l.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(l.Image, op)
	}
}

func (l *Label) Resize(rect []int) {
	l.rect = NewRect(rect)
	l.Dirty = true
	l.Image = nil
}

func (l Label) String() string {
	return fmt.Sprintf("%v %v", l.text, l.rect)
}

func (l *Label) Close() {
	l.Image.Dispose()
}
