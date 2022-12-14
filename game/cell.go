package game

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	ui "github.com/t0l1k/eui"
)

type Cell struct {
	rect                                      *ui.Rect
	Image                                     *ebiten.Image
	Dirty, Visibe, DrawRect, IsCenter, Active bool
	bg, fg, activeColor                       color.Color
	margin                                    float64
	sym                                       int
	text                                      string
	ariphmetic                                bool
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
		sym:         0,
	}
}

func (c *Cell) Layout() {
	w, h := c.rect.Size()
	if c.Image == nil {
		c.Image = ebiten.NewImage(w, h)
	} else {
		c.Image.Clear()
	}
	bg := c.bg
	if c.Active {
		bg = c.activeColor
	}
	m := float64(w) * c.margin
	if c.DrawRect && c.sym == 0 {
		ebitenutil.DrawRect(c.Image, m, m, float64(w)-m*2, float64(h)-m*2, bg)
	}
	if c.sym > 0 && c.Active && !c.ariphmetic {
		res := fmt.Sprintf("%v", strconv.Itoa(c.sym))
		l := ui.NewLabel(res, []int{0, 0, w, h}, c.bg, c.activeColor)
		defer l.Close()
		l.Draw(c.Image)
	} else if c.sym > 0 && c.Active && c.ariphmetic {
		l := ui.NewLabel(c.text, []int{0, 0, w, h}, c.bg, c.activeColor)
		defer l.Close()
		l.Draw(c.Image)
	}

	if c.IsCenter {
		m := float64(c.rect.H) * 0.4
		x1 := float64(c.rect.W) / 2
		y1 := m
		x2 := x1
		y2 := float64(c.rect.H) - m
		ebitenutil.DrawLine(c.Image, x1, y1, x2, y2, c.fg)
		ebitenutil.DrawLine(c.Image, y1, x1, y2, x2, c.fg)
	}
	c.Dirty = false
}

func (c *Cell) SetActiveColor(value color.Color) {
	if c.activeColor == value {
		return
	}
	c.activeColor = value
	c.Dirty = true
}

func (c *Cell) SetSymbol(value int) {
	if c.sym == value {
		return
	}
	c.sym = value
	c.Dirty = true
}

func (c *Cell) SetAriphmetic(value int) {
	c.ariphmetic = true
	max := ui.GetPreferences().Get("ariphmetic max").(int)
	c.sym = value
	oper := NewOperation()
	oper.Rand()
	a, b := oper.Get(getNum(value), value, max)
	c.text = fmt.Sprintln(a, oper, b)
	c.Dirty = true
}

func (c *Cell) SetActive(value bool) {
	c.Active = value
	c.Dirty = true
}
func (c *Cell) Update(dt int) {}
func (c *Cell) Draw(surface *ebiten.Image) {
	if c.Dirty {
		c.Layout()
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
	c.Image = nil
}

func (c Cell) String() string {
	return fmt.Sprintf("Cell: %v", c.rect)
}

func (c *Cell) Close() {
	c.Image.Dispose()
}
