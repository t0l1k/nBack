package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type TopBarOpt struct {
	ui.ContainerDefault
	lblName                     *ui.Label
	btnQuit, btnReset, btnApply *ui.Button
}

func NewTopBarOpt(fnReset, fnApply func(b *ui.Button)) *TopBarOpt {
	s := &TopBarOpt{}
	rect := []int{0, 0, 1, 1}
	s.btnQuit = ui.NewButton("<", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) { ui.Pop() })
	s.Add(s.btnQuit)
	s.lblName = ui.NewLabel(ui.GetLocale().Get("btnOpt"), rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblName)
	s.btnReset = ui.NewButton(ui.GetLocale().Get("btnReset"), rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), fnReset)
	s.Add(s.btnReset)
	s.btnApply = ui.NewButton(ui.GetLocale().Get("btnSave"), rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), fnApply)
	s.Add(s.btnApply)
	return s
}

func (r *TopBarOpt) Update(dt int) {
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *TopBarOpt) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *TopBarOpt) Resize() {
	w, h := ui.GetUi().GetScreenSize()
	rect := ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(rect.H)*0.05), int(float64(rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, int(float64(rect.W)*0.20)
	s.lblName.Resize([]int{x, y, w, h})
	s.btnReset.Resize([]int{rect.W - w*2, y, w, h})
	s.btnApply.Resize([]int{rect.W - w, y, w, h})
}

func (r *TopBarOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
