package game

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
)

type cell struct {
	eui.Icon
	modals                        []*Modality
	active, show, center          bool
	bg, bgActive, fg, fgCrosshair color.Color
	text                          string
	conf                          GameConf
}

func newCell(center bool) *cell {
	c := &cell{}
	theme := eui.GetUi().GetTheme()
	c.bg = theme.Get(app.GameColorBg)
	c.bgActive = theme.Get(app.GameColorActiveBg)
	c.fg = theme.Get(app.GameColorFg)
	c.fgCrosshair = theme.Get(app.GameColorFgCrosshair)
	c.Bg(c.bg)
	c.Fg(c.fg)
	c.center = center
	c.Icon.Visible = true
	return c
}

func (c *cell) Setup(conf GameConf, modals []*Modality) {
	c.conf = conf
	c.modals = modals
	c.show = true
}

func (c *cell) SetText() {
	lbl := eui.NewText(c.text)
	if c.active {
		lbl.Bg(c.bgActive)
		lbl.Fg(c.fg)
	}
	if !c.active || !c.show {
		lbl.Bg(c.bg)
		lbl.Fg(c.bg)
	}
	defer lbl.Close()
	lbl.Resize(c.GetRect().GetArr())
	lbl.Layout()
	srcImg := lbl.Image()
	c.SetIcon(srcImg)
	c.Layout()
	if c.center {
		dstImage := c.SetCrossHair()
		c.Image().DrawImage(dstImage, nil)
		op := &ebiten.DrawImageOptions{}
		op.Blend = ebiten.BlendDestinationAtop
		c.Image().DrawImage(srcImg, op)
	}
}

func (c *cell) SetCrossHair() *ebiten.Image {
	r, g, b, _ := c.bg.RGBA() // фон прозрачный
	a := 1
	bg := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	crosshair := NewCrosshair(bg, c.fgCrosshair, 3, 0.34)
	defer crosshair.Close()
	rect := c.GetRect()
	crosshair.Resize(rect.GetArr())
	crosshair.Layout()
	return crosshair.Image()
}

func (c *cell) SetActive(idx int) {
	c.active = true
	for _, v := range c.modals {
		switch v.GetSym() {
		case Col:
			c.bgActive = Colors[v.GetField()[idx]]
		case Sym:
			c.text = strconv.Itoa(v.GetField()[idx])
		case Ari:
			maxNum := c.conf.Get(MaxNumber).(int)
			oper := NewOperation()
			oper.Rand(c.conf)
			a, b := oper.Get(getNum(v.GetField()[idx]), v.GetField()[idx], maxNum)
			c.text = fmt.Sprintln(a, oper, b)
		}
	}
	c.SetText()
}

func (c *cell) SetInactive() {
	c.active = false
	c.SetText()
}

func (c *cell) IsVisible() bool {
	return c.show
}

func (c *cell) Visible(value bool) {
	if value {
		if value && !c.show {
			c.show = true
		}
	}
	if !value {
		if c.show {
			c.show = false
		}
	}
	c.SetText()
}
