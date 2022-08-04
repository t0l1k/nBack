package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Combobox struct {
	text                       string
	rect, upRect, downRect     *Rect
	Image                      *ebiten.Image
	Dirty, Visible, focus      bool
	bg, fg                     color.Color
	mouseUpDown, mouseDownDown bool
	data                       []interface{}
	current                    int
	onChange                   func(c *Combobox)
}

func NewCombobox(text string, rect []int, bg, fg color.Color, data []interface{}, current int, f func(c *Combobox)) *Combobox {
	return &Combobox{
		text:          text,
		rect:          NewRect(rect),
		Image:         nil,
		Dirty:         true,
		Visible:       true,
		focus:         false,
		mouseUpDown:   false,
		mouseDownDown: false,
		bg:            bg,
		fg:            fg,
		data:          data,
		current:       current,
		onChange:      f,
	}
}

func (c *Combobox) Layout() *ebiten.Image {
	if !c.Dirty {
		return c.Image
	}
	w, h := c.rect.Size()
	image := ebiten.NewImage(w, h)
	boxHeight := h
	m := 1
	x, y, w, h := boxHeight*2, 0, w-h*2, h
	lblText := NewLabel(c.text, []int{x, y, w, h}, c.bg, c.fg)
	var result string
	switch value := c.data[c.current].(type) {
	case int:
		result = strconv.Itoa(value)
	case float64:
		result = fmt.Sprintf("%.1f", value)
	case string:
		result = fmt.Sprintf("%v", value)
	}
	lblValue := NewLabel(result, []int{0, 0, boxHeight, h}, c.bg, c.fg)
	x, y, w, h = boxHeight+m, m, boxHeight-m*2, (boxHeight-m*2)/2
	btnUpRect := []int{x, y, w, h}
	x, y, w, h = boxHeight+m, (boxHeight/2)+m, boxHeight-m*2, (boxHeight-m*2)/2
	btnDownRect := []int{x, y, w, h}
	btnUp := NewLabel("\u25b2", btnUpRect, c.bg, c.fg)
	btnDown := NewLabel("\u25bc", btnDownRect, c.bg, c.fg)
	var (
		bg, fg           color.Color
		rectUp, rectDown bool
	)
	if c.focus && !c.mouseUpDown && !c.mouseDownDown {
		bg = c.fg
		fg = c.bg
		rectUp = false
		rectDown = false
	} else if !c.focus && !c.mouseUpDown {
		bg = c.bg
		fg = c.fg
	} else if c.focus && c.mouseUpDown {
		bg = c.fg
		fg = c.bg
		rectUp = true
		rectDown = false
		m = 2
	} else if c.focus && c.mouseDownDown {
		bg = c.fg
		fg = c.bg
		rectUp = false
		rectDown = true
		m = 2
	}
	lblText.SetRect(true)
	lblText.SetBg(bg)
	lblText.SetFg(fg)
	lblText.Draw(image)
	lblValue.SetRect(true)
	lblValue.SetBg(bg)
	lblValue.SetFg(fg)
	lblValue.Draw(image)
	btnUp.SetRect(rectUp)
	btnUp.SetBg(bg)
	btnUp.SetFg(fg)
	btnUp.Draw(image)
	btnDown.SetRect(rectDown)
	btnDown.SetBg(bg)
	btnDown.SetFg(fg)
	btnDown.Draw(image)
	c.Dirty = false
	return image
}

func (c *Combobox) Value() interface{} { return c.data[c.current] }

func (c *Combobox) SetValue(value interface{}) {
	for i, v := range c.data {
		if v == value {
			c.current = i
			c.Dirty = true
			break
		}
	}
}

func (c *Combobox) SetFocus(value bool) {
	if c.focus == value {
		return
	}
	c.focus = value
	c.Dirty = true
}

func (c *Combobox) SetMouseUpDown(value bool) {
	if c.mouseUpDown == value {
		return
	}
	c.mouseUpDown = value
	c.Dirty = true
}
func (c *Combobox) SetMouseDownDown(value bool) {
	if c.mouseDownDown == value {
		return
	}
	c.mouseDownDown = value
	c.Dirty = true
}

func (c *Combobox) Update(dt int) {
	x, y := ebiten.CursorPosition()
	if c.rect.InRect(x, y) {
		c.SetFocus(true)
	} else {
		c.SetFocus(false)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if c.upRect.InRect(x, y) {
			c.SetMouseUpDown(true)
		} else {
			c.SetMouseUpDown(false)
		}
		if c.downRect.InRect(x, y) {
			c.SetMouseDownDown(true)
		} else {
			c.SetMouseDownDown(false)
		}
	} else {
		if c.mouseUpDown {
			if c.current < len(c.data)-1 {
				c.current++
			}
			if c.onChange != nil {
				c.onChange(c)
			}
		}
		if c.mouseDownDown {
			if c.current > 0 {
				c.current--
			}
			if c.onChange != nil {
				c.onChange(c)
			}
		}
		c.SetMouseUpDown(false)
		c.SetMouseDownDown(false)
	}
}

func (c *Combobox) Draw(surface *ebiten.Image) {
	if c.Dirty {
		c.Image = c.Layout()
	}
	if c.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := c.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(c.Image, op)
	}
}

func (c *Combobox) Resize(rect []int) {
	c.rect = NewRect(rect)
	boxHeight := c.rect.H
	x, y, w, h := c.rect.X+boxHeight, c.rect.Y, boxHeight, (boxHeight)/2
	btnUpRect := []int{x, y, w, h}
	x, y, w, h = c.rect.X+boxHeight, c.rect.Y+(boxHeight/2), boxHeight, (boxHeight)/2
	btnDownRect := []int{x, y, w, h}
	c.upRect = NewRect(btnUpRect)
	c.downRect = NewRect(btnDownRect)
	c.Dirty = true
}
