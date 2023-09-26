package options

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

type TopBarOpt struct {
	eui.ContainerDefault
	lblName                     *eui.Label
	btnQuit, btnReset, btnApply *eui.Button
}

func NewTopBarOpt(fnReset, fnApply func(b *eui.Button)) *TopBarOpt {
	s := &TopBarOpt{}
	rect := []int{0, 0, 1, 1}
	s.btnQuit = eui.NewButton("<", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { eui.Pop() })
	s.Add(s.btnQuit)
	s.lblName = eui.NewLabel(eui.GetLocale().Get("btnOpt"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblName)
	s.btnReset = eui.NewButton(eui.GetLocale().Get("btnReset"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), fnReset)
	s.Add(s.btnReset)
	s.btnApply = eui.NewButton(eui.GetLocale().Get("btnSave"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), fnApply)
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
	w, h := eui.GetUi().GetScreenSize()
	rect := eui.NewRect([]int{0, 0, w, h})
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
