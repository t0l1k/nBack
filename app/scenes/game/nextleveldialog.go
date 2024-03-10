package scene_game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

type nextLevelDialog struct {
	eui.DrawableBase
	messageLbl, timeLbl *eui.Text
	timer               *eui.Timer
}

func newNextLevelDialog(delay int) *nextLevelDialog {
	n := &nextLevelDialog{}
	n.messageLbl = eui.NewText("")
	n.Add(n.messageLbl)
	n.timeLbl = eui.NewText("")
	n.messageLbl.Visible = false
	n.timeLbl.Visible = false
	n.Add(n.timeLbl)
	n.timer = eui.NewTimer(delay)
	n.messageLbl.Visible = false
	n.timeLbl.Visible = false
	n.Visible = false
	return n
}

func (n *nextLevelDialog) show(msg string, col color.Color) {
	n.Bg(col)
	n.messageLbl.Bg(col)
	n.timeLbl.Bg(col)
	n.messageLbl.SetText(msg)
	n.timeLbl.SetText(n.timer.String())
	n.Visible = true
	n.messageLbl.Visible = true
	n.timeLbl.Visible = true
	n.timer.On()
}

func (n *nextLevelDialog) Update(dt int) {
	n.timer.Update(dt)
	n.timeLbl.SetText(n.timer.String())
	if n.timer.IsDone() {
		n.messageLbl.Visible = false
		n.timeLbl.Visible = false
		n.Visible = false
	}
}

func (n *nextLevelDialog) Draw(surface *ebiten.Image) {
	for _, v := range n.GetContainer() {
		v.Draw(surface)
	}
}

func (n *nextLevelDialog) Resize(rect []int) {
	n.Rect(eui.NewRect(rect))
	n.SpriteBase.Resize(rect)
	x, y := n.GetRect().Pos()
	w, h := n.GetRect().Size()
	n.messageLbl.Resize([]int{x, y, w, h / 2})
	n.timeLbl.Resize([]int{x, y + h/2, w, h / 2})
	n.ImageReset()
}
