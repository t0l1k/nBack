package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type SceneOptions struct {
	ui.ContainerDefault
	topBar                                                          *TopBarOpt
	lblSelectGame, lblSelectOptios                                  *ui.Label
	btnAppOpt, btnGameOpt, btnClassicGame, btnJaeggiGame, btnManual *ui.Button
	btnMoves, btnThreePigs, btnUglyDuck, btnModals                  *ui.Button
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{}
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(nil, nil)
	s.topBar.btnReset.Visible = false
	s.topBar.btnApply.Visible = false
	s.Add(s.topBar)

	s.lblSelectGame = ui.NewLabel("Выбрать игру", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblSelectGame)

	s.lblSelectOptios = ui.NewLabel("Настройки", rect, ui.GetTheme().Get("warning color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblSelectOptios)

	s.btnAppOpt = ui.NewButton("Настройки приложения", rect, ui.GetTheme().Get("warning color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewAppOpt())
	})
	s.Add(s.btnAppOpt)
	s.btnGameOpt = ui.NewButton("Все настройки игры", rect, ui.GetTheme().Get("warning color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewGameOpt())
	})
	s.Add(s.btnGameOpt)

	s.btnClassicGame = ui.NewButton("Классический nBack(BrainWorkshop)", rect, ui.GetTheme().Get("regular color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewClassicOpt())
	})
	s.Add(s.btnClassicGame)

	s.btnJaeggiGame = ui.NewButton("Jaeggi nBack mode", rect, ui.Fuchsia, ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewJaeggiOpt())
	})
	s.Add(s.btnJaeggiGame)

	s.btnManual = ui.NewButton("Играть на ручнике", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewManualOpt())
	})
	s.Add(s.btnManual)

	s.btnMoves = ui.NewButton("Настройка числа ходов(сложность)", rect, ui.GetTheme().Get("warning color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewMovesOpt())
	})
	s.Add(s.btnMoves)

	s.btnThreePigs = ui.NewButton("Играть спасти трех поросят", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewOptThreePigs())
	})
	s.Add(s.btnThreePigs)

	s.btnUglyDuck = ui.NewButton("Играть преобразить гадкого утенка", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewOptUglyDuck())
	})
	s.Add(s.btnUglyDuck)

	s.btnModals = ui.NewButton("Настройка модальностей", rect, ui.GetTheme().Get("warning color"), ui.GetTheme().Get("fg"), func(b *ui.Button) {
		ui.Push(NewOptModals())
	})
	s.Add(s.btnModals)

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
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*6
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.lblSelectGame.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnClassicGame.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnJaeggiGame.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnManual.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnUglyDuck.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnThreePigs.Resize([]int{x, y, w1, hTop - 2})
	y += hTop * 2
	s.lblSelectOptios.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnAppOpt.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnModals.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnMoves.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.btnGameOpt.Resize([]int{x, y, w1, hTop - 2})
}

func (s *SceneOptions) Close() {
	for _, v := range s.Container {
		v.Close()
	}
}
