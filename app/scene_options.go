package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type SceneOptions struct {
	ui.ContainerDefault
	topBar                                           *TopBarOpt
	btnAppOpt, btnGameOpt, btnClassicGame, btnManual *ui.Button
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{}
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(nil, nil)
	s.topBar.btnReset.Disable()
	s.topBar.btnApply.Disable()
	s.Add(s.topBar)
	s.btnAppOpt = ui.NewButton("App Options", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewAppOpt())
	})
	s.Add(s.btnAppOpt)
	s.btnGameOpt = ui.NewButton("All Game Options", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewGameOpt())
	})
	s.Add(s.btnGameOpt)

	s.btnClassicGame = ui.NewButton("Классический нназад(BrainWorkshop)", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewClassicOpt())
	})
	s.Add(s.btnClassicGame)

	s.btnManual = ui.NewButton("Игра на ручние", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewManualOpt())
	})
	s.Add(s.btnManual)

	return s
}

func (s *SceneOptions) Entered() {
	s.Resize()
	log.Println("Entered SceneOptions")
}

func (s *SceneOptions) Update(dt int) {
	for _, value := range s.Container {
		value.Update(dt)
	}
}

func (s *SceneOptions) Draw(surface *ebiten.Image) {
	surface.Fill(ui.GetTheme().Get("game bg"))
	for _, value := range s.Container {
		value.Draw(surface)
	}
}

func (s *SceneOptions) Resize() {
	s.topBar.Resize()
	w, h := ui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*4
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.btnAppOpt.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnClassicGame.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnManual.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnGameOpt.Resize([]int{x, y, w1, hTop - 2})
}

func (s *SceneOptions) Close() {
	for _, v := range s.Container {
		v.Close()
	}
}
