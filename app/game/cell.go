package game

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/data"
)

type cell struct {
	eui.Icon
	modals                        []*data.Modality
	active, show, center          bool
	bg, bgActive, fg, fgCrosshair color.Color
	text                          string
}

func newCell(center bool) *cell {
	c := &cell{}
	c.SetupIcon(nil)
	conf := eui.GetUi().GetSettings()
	c.bg = conf.Get(app.GameColorBg).(color.Color)
	c.bgActive = conf.Get(app.GameColorActiveBg).(color.Color)
	c.fg = conf.Get(app.GameColorFg).(color.Color)
	c.fgCrosshair = conf.Get(app.GameColorFgCrosshair).(color.Color)
	c.Bg(c.bg)
	c.Fg(c.fg)
	c.center = center
	return c
}

func (c *cell) Setup(modals []*data.Modality) {
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
	rect := c.GetRect()
	lbl.Resize(rect.GetArr())
	lbl.Layout()
	srcImg := lbl.GetImage()
	c.SetIcon(srcImg)
	c.Layout()
	if c.center {
		dstImage := c.SetCrossHair()
		c.GetImage().DrawImage(dstImage, nil)
		op := &ebiten.DrawImageOptions{}
		op.Blend = ebiten.BlendDestinationAtop
		c.GetImage().DrawImage(srcImg, op)
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
	return crosshair.Image
}

func (c *cell) SetActive(idx int) {
	c.active = true
	for _, v := range c.modals {
		switch v.GetSym() {
		case data.Col:
			c.bgActive = Colors[v.GetField()[idx]]
		case data.Sym:
			c.text = strconv.Itoa(v.GetField()[idx])
		case data.Ari:
			conf := eui.GetUi().GetSettings()
			maxNum := conf.Get(app.MaxNumber).(int)
			oper := NewOperation()
			oper.Rand()
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
