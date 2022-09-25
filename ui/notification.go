package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Notification struct {
	text                     string
	rect                     *Rect
	Image                    *ebiten.Image
	Dirty, Visible, drawRect bool
	bg, fg                   color.Color
	delay                    int
	elapsed                  int
}

func NewNotification(text string, rect []int, bg, fg color.Color) *Notification {
	return &Notification{
		text:     text,
		rect:     NewRect(rect),
		Image:    nil,
		Dirty:    true,
		Visible:  true,
		drawRect: false,
		bg:       bg,
		fg:       fg,
		delay:    1,
		elapsed:  0,
	}
}

func (n *Notification) SetBg(value color.Color) {
	if n.bg == value {
		return
	}
	n.bg = value
	n.Dirty = true
}

func (n *Notification) SetRect(value bool) {
	if n.drawRect == value {
		return
	}
	n.drawRect = value
	n.Dirty = true
}

func (n *Notification) SetFg(value color.Color) {
	if n.fg == value {
		return
	}
	n.fg = value
	n.Dirty = true
}

func (n *Notification) SetText(value string) {
	if n.text == value {
		return
	}
	n.text = value
	n.Dirty = true
	n.Visible = true
	n.elapsed = 0
}

func (n *Notification) Layout() {
	w, h := n.rect.Size()
	if n.Image == nil {
		n.Image = ebiten.NewImage(w, h)
	} else {
		n.Image.Clear()
	}
	n.Image.Fill(n.bg)
	if n.drawRect {
		ebitenutil.DrawRect(n.Image, 0, 0, float64(w), float64(h), n.fg)
		ebitenutil.DrawRect(n.Image, 2, 2, float64(w)-4, float64(h)-4, n.bg)
	}
	fntSize := GetFonts().calcFontSize(n.text, n.rect)
	font := GetFonts().get(fntSize)
	b := text.BoundString(font, n.text)
	x := (n.rect.W - b.Max.X) / 2
	y := n.rect.H - (n.rect.H+b.Min.Y)/2
	text.Draw(n.Image, n.text, font, x, y, n.fg)
	n.Dirty = false
}

func (n *Notification) Update(dt int) {
	if n.elapsed > n.delay*1000 {
		n.Visible = false
	}
	n.elapsed += dt
}

func (n *Notification) Draw(surface *ebiten.Image) {
	if n.Dirty {
		n.Layout()
	}
	if n.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := n.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(n.Image, op)
	}
}

func (n *Notification) Resize(rect []int) {
	n.rect = NewRect(rect)
	n.Dirty = true
	n.Image = nil
}

func (n Notification) String() string {
	return fmt.Sprintf("%v %v", n.text, n.rect)
}

func (n *Notification) Close() {
	n.Image.Dispose()
}
