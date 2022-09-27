package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Notification struct {
	delay   int
	elapsed int
	lbl     *Label
	Show    bool
}

func NewNotification(text string, delay int, rect []int, bg, fg color.Color) *Notification {
	return &Notification{
		delay:   delay,
		elapsed: 0,
		lbl:     NewLabel(text, rect, bg, fg),
		Show:    true,
	}
}

func (n *Notification) Update(dt int) {
	if !n.Show {
		return
	}
	if n.elapsed > n.delay*1000 {
		n.Show = false
		n.lbl = nil
	}
	n.elapsed += dt
}

func (n *Notification) Draw(surface *ebiten.Image) {
	if !n.Show {
		return
	}
	n.lbl.Draw(surface)
}
