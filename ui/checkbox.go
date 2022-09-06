package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Checkbox struct {
	rect                  *Rect
	Image                 *ebiten.Image
	Dirty, Visible, focus bool
	bg, fg                color.Color
	text                  string
	checked               bool
	mouseDown             bool
	onCheckChanged        func(c *Checkbox)
}

func NewCheckbox(text string, rect []int, bg, fg color.Color, f func(c *Checkbox)) *Checkbox {
	return &Checkbox{
		text:           text,
		rect:           NewRect(rect),
		Image:          nil,
		Dirty:          true,
		Visible:        true,
		focus:          false,
		checked:        false,
		mouseDown:      false,
		bg:             bg,
		fg:             fg,
		onCheckChanged: f,
	}
}

func (c *Checkbox) Layout() {
	w, h := c.rect.Size()
	if c.Image == nil {
		c.Image = ebiten.NewImage(w, h)
	} else {
		c.Image.Clear()
	}
	boxWidth, boxHeight := h, h
	m := int(float64(c.rect.GetLowestSize()) * 0.1)
	lbl := NewLabel(c.text, []int{boxHeight + m, m, w - h - m*2, h - m*2}, c.bg, c.fg)
	defer lbl.Close()
	var (
		bg, fg color.Color
		rect   bool
	)
	if c.focus && !c.mouseDown {
		bg = c.fg
		fg = c.bg
		rect = false
	} else if !c.focus && !c.mouseDown {
		bg = c.bg
		fg = c.fg
	} else if c.focus && c.mouseDown {
		bg = c.fg
		fg = c.bg
		rect = true
		m = 2
	}
	ebitenutil.DrawRect(c.Image, 0, 0, float64(w), float64(h), fg)
	ebitenutil.DrawRect(c.Image, float64(m), float64(m), float64(w-m*2), float64(h-m*2), bg)
	lbl.SetRect(rect)
	lbl.SetBg(bg)
	lbl.SetFg(fg)
	lbl.Draw(c.Image)
	if c.checked {
		ebitenutil.DrawRect(c.Image, 0, 0, float64(boxWidth), float64(boxHeight), fg)
		ebitenutil.DrawRect(c.Image, float64(m), float64(m), float64(boxWidth-m*2), float64(boxHeight-m*2), bg)
		ebitenutil.DrawLine(c.Image, float64(m*3), float64(m*3), float64(boxWidth-m*3), float64(boxHeight-m*3), fg)
		ebitenutil.DrawLine(c.Image, float64(boxWidth-m*3), float64(m*3), float64(m*3), float64(boxHeight-m*3), fg)
	} else {
		ebitenutil.DrawRect(c.Image, 0, 0, float64(boxWidth), float64(boxHeight), fg)
		ebitenutil.DrawRect(c.Image, float64(m), float64(m), float64(boxWidth-m*2), float64(boxHeight-m*2), bg)
	}
	c.Dirty = false
}

func (c *Checkbox) Checked() bool { return c.checked }
func (c *Checkbox) SetChecked(value bool) {
	if c.checked == value {
		return
	}
	c.checked = value
	c.Dirty = true
}

func (c *Checkbox) SetFocus(value bool) {
	if c.focus == value {
		return
	}
	c.focus = value
	c.Dirty = true
}

func (c *Checkbox) SetMouseDown(value bool) {
	if c.mouseDown == value {
		return
	}
	c.mouseDown = value
	c.Dirty = true
}

func (c *Checkbox) Update(dt int) {
	x, y := ebiten.CursorPosition()
	if c.rect.InRect(x, y) {
		c.SetFocus(true)
	} else {
		c.SetFocus(false)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if c.rect.InRect(x, y) {
			c.SetMouseDown(true)
		} else {
			c.SetMouseDown(false)
		}
	} else {
		if c.mouseDown {
			c.checked = !c.checked
			if c.onCheckChanged != nil {
				c.onCheckChanged(c)
			}
		}
		c.SetMouseDown(false)
	}
}

func (c *Checkbox) Draw(surface *ebiten.Image) {
	if c.Dirty {
		c.Layout()
	}
	if c.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := c.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(c.Image, op)
	}
}

func (c *Checkbox) Resize(rect []int) {
	c.rect = NewRect(rect)
	c.Dirty = true
	c.Image = nil
}

func (c *Checkbox) Close() {
	c.Image.Dispose()
}
