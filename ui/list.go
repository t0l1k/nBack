package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type List struct {
	rect                     *Rect
	Image                    *ebiten.Image
	Dirty, Visible, focus    bool
	bg, fg                   color.Color
	a                        []string
	b                        []color.Color
	rows                     int
	mouseDown, mouseUp, drag bool
	dX, dY, lastY, lastJ     int
	boxWidth, boxHeight      int
	jump, lastDiff           int
	end                      bool
}

func NewList(a []string, b []color.Color, rect []int, bg, fg color.Color, rows int) *List {
	return &List{
		a:         a,
		b:         b,
		rect:      NewRect(rect),
		Image:     nil,
		Dirty:     true,
		Visible:   true,
		focus:     false,
		mouseDown: false,
		mouseUp:   false,
		drag:      false,
		bg:        bg,
		fg:        fg,
		rows:      rows,
		jump:      0,
	}
}

func (l *List) SetList(a []string, b []color.Color) {
	l.a = a
	l.b = b
	l.Dirty = true
}

func (l *List) SetRows(rows int) {
	l.rows = rows
	l.Dirty = true
}

func (l *List) Layout() *ebiten.Image {
	if !l.Dirty {
		return l.Image
	}
	w, h := l.rect.Size()
	image := ebiten.NewImage(w, h)
	m := 2.0
	if l.focus {
		ebitenutil.DrawRect(image, 0, 0, float64(w), float64(h), l.fg)
		ebitenutil.DrawRect(image, m, m, float64(w)-m*2, float64(h)-m*2, l.bg)
	}
	boxWidth := l.boxWidth - int(m)/l.rows
	boxHeight := l.boxHeight - int(m)/l.rows
	j := l.jump * l.rows
	for i := 0; i < len(l.a); i++ {
		x := i%l.rows*boxWidth + int(m)
		y := i/l.rows*boxHeight + int(m)
		if y+boxHeight > l.rect.Bottom() {
			break
		}
		lbl := NewLabel(l.a[j], []int{x, y, boxWidth, boxHeight}, l.bg, l.fg)
		lbl.SetBg(l.b[j])
		lbl.Draw(image)
		if j < len(l.a)-1 {
			j++
		} else {
			if y+boxHeight <= l.rect.H {
				l.end = true
			} else if y >= l.rect.H {
				l.end = false
			}
			break
		}
	}
	l.Dirty = false
	return image
}

func (l *List) setFocus(value bool) {
	if l.focus == value {
		return
	}
	l.focus = value
	l.Dirty = true
}

func (l *List) Update(dt int) {
	x, y := ebiten.CursorPosition()
	if l.rect.InRect(x, y) {
		l.setFocus(true)
	} else {
		l.setFocus(false)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if l.focus && !l.mouseDown {
			l.mouseDown = true
			l.mouseUp = false
		}
		if l.mouseDown && l.lastY != y && !l.drag {
			if len(l.a)/l.rows*l.boxHeight > l.rect.H {
				l.drag = true
				l.dX = x
				l.dY = y
			}
		}
	} else {
		if l.mouseDown {
			l.mouseDown = false
			l.drag = false
			l.lastJ = 0
			l.lastDiff = 0
		}
	}
	if l.drag {
		diff := l.dY - l.lastY
		j := diff / l.boxHeight
		if j != l.lastJ && diff-l.lastDiff > 0 {
			if !l.end {
				l.jump++
				l.Dirty = true
			}
		}
		if j != l.lastJ && diff-l.lastDiff < 0 {
			if l.jump > 0 {
				l.jump--
				l.Dirty = true
			}
		}
		if l.lastJ <= 0 {
			l.lastJ = 0
		}
		l.lastJ = j
		l.lastDiff = diff
	}
	l.lastY = y
}

func (l *List) Draw(surface *ebiten.Image) {
	if l.Dirty {
		l.Image = l.Layout()
	}
	if l.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := l.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(l.Image, op)
	}
}

func (l *List) Resize(rect []int) {
	l.rect = NewRect(rect)
	l.boxWidth = l.rect.W / l.rows
	l.boxHeight = int(float64(l.rect.GetLowestSize()) * 0.05)
	l.jump = 0
	l.end = false
	l.Dirty = true
}
