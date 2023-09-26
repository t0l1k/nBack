package options

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

type SceneOptions struct {
	eui.ContainerDefault
	topBar                                                          *TopBarOpt
	lblSelectGame, lblSelectOptios                                  *eui.Label
	btnAppOpt, btnGameOpt, btnClassicGame, btnJaeggiGame, btnManual *eui.Button
	btnMoves, btnThreePigs, btnUglyDuck, btnModals                  *eui.Button
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{}
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(nil, nil)
	s.topBar.btnReset.Visible = false
	s.topBar.btnApply.Visible = false
	s.Add(s.topBar)

	s.lblSelectGame = eui.NewLabel("Выбрать игру", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblSelectGame)

	s.lblSelectOptios = eui.NewLabel("Настройки", rect, eui.GetTheme().Get("warning color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblSelectOptios)

	s.btnAppOpt = eui.NewButton("Настройки приложения", rect, eui.GetTheme().Get("warning color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewAppOpt())
	})
	s.Add(s.btnAppOpt)
	s.btnGameOpt = eui.NewButton("Все настройки игры", rect, eui.GetTheme().Get("warning color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewGameOpt())
	})
	s.Add(s.btnGameOpt)

	s.btnClassicGame = eui.NewButton("Классический nBack(BrainWorkshop)", rect, eui.GetTheme().Get("regular color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewClassicOpt())
	})
	s.Add(s.btnClassicGame)

	s.btnJaeggiGame = eui.NewButton("Jaeggi nBack mode", rect, eui.Fuchsia, eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewJaeggiOpt())
	})
	s.Add(s.btnJaeggiGame)

	s.btnManual = eui.NewButton("Играть на ручнике", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewManualOpt())
	})
	s.Add(s.btnManual)

	s.btnMoves = eui.NewButton("Настройка числа ходов(сложность)", rect, eui.GetTheme().Get("warning color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewMovesOpt())
	})
	s.Add(s.btnMoves)

	s.btnThreePigs = eui.NewButton("Играть спасти трех поросят", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewOptThreePigs())
	})
	s.Add(s.btnThreePigs)

	s.btnUglyDuck = eui.NewButton("Играть преобразить гадкого утенка", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewOptUglyDuck())
	})
	s.Add(s.btnUglyDuck)

	s.btnModals = eui.NewButton("Настройка модальностей", rect, eui.GetTheme().Get("warning color"), eui.GetTheme().Get("fg"), func(b *eui.Button) {
		eui.Push(NewOptModals())
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
	surface.Fill(eui.GetTheme().Get("game bg"))
	for _, value := range s.Container {
		value.Draw(surface)
	}
}

func (s *SceneOptions) Resize() {
	s.topBar.Resize()
	w, h := eui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := eui.NewRect([]int{0, hTop, w, h - hTop})
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
