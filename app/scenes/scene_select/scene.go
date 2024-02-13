package scene_select

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app/data"
	"github.com/t0l1k/nBack/app/scenes/create"
	"github.com/t0l1k/nBack/app/scenes/intro"
)

type SceneSelectGame struct {
	eui.SceneBase
	topBar                                  *eui.TopBar
	listGames                               *eui.ListView
	btnSel, btnCrt, btnPrgs, btnTut, btnOpt *eui.Button
	btnsLayout                              *eui.BoxLayout
	profiles                                *data.GameProfiles
}

func NewSceneSelectGame() *SceneSelectGame {
	s := &SceneSelectGame{}
	s.topBar = eui.NewTopBar("нНазад", nil)
	s.topBar.SetShowStopwatch()
	s.Add(s.topBar)
	s.profiles = data.DefalutGameProfiles()
	s.listGames = eui.NewListView()
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(eui.ButtonBg)
	fg := theme.Get(eui.ButtonFg)
	s.listGames.SetupListViewButtons(s.profiles.GetProfilesName(), 30, 1, bg, fg, s.btnsLogic)
	s.Add(s.listGames)
	s.listGames.Bg(eui.Blue)
	s.btnsLayout = eui.NewHLayout()
	s.btnSel = eui.NewButton("Играть", s.btnsLogic)
	s.btnSel.Disable()
	s.btnSel.Bg(eui.YellowGreen)
	s.btnsLayout.Add(s.btnSel)
	s.btnCrt = eui.NewButton("Создать", s.btnsLogic)
	s.btnsLayout.Add(s.btnCrt)
	s.btnPrgs = eui.NewButton("Итоги", s.btnsLogic)
	s.btnsLayout.Add(s.btnPrgs)
	s.btnTut = eui.NewButton("Обучение", s.btnsLogic)
	s.btnsLayout.Add(s.btnTut)
	s.btnOpt = eui.NewButton("Настройки", s.btnsLogic)
	s.btnsLayout.Add(s.btnOpt)
	s.Resize()
	return s
}

func (s *SceneSelectGame) btnsLogic(b *eui.Button) {
	fmt.Println("selected", b.GetText())
	for name, game := range s.profiles.GetGameProfiles() {
		if name == b.GetText() {
			sc := intro.NewSceneIntro(game, name)
			eui.GetUi().Push(sc)
			log.Println("selected profile:", b.GetText())
		}
	}

	if b.GetText() == "Создать" {
		sc := create.NewSceneCreateGame(s.profiles)
		eui.GetUi().Push(sc)
		log.Println("Выбрана сцена создание профиля")
	}
}

func (s *SceneSelectGame) Entered() {
	s.listGames.Reset()
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(eui.ButtonBg)
	fg := theme.Get(eui.ButtonFg)
	s.listGames.SetupListViewButtons(s.profiles.GetProfilesName(), 30, 1, bg, fg, s.btnsLogic)
	for k, v := range s.profiles.GetGameProfiles() {
		for _, v1 := range v.Games {
			if v1.IsDone() {
				fmt.Println(k, v1.LastGameFullResult())
			}
		}
	}
}

func (s *SceneSelectGame) Update(dt int) {
	for _, v := range s.GetContainer() {
		v.Update(dt)
	}
	for _, v := range s.btnsLayout.GetContainer() {
		v.Update(dt)
	}
}

func (s *SceneSelectGame) Draw(surface *ebiten.Image) {
	for _, v := range s.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range s.btnsLayout.GetContainer() {
		v.Draw(surface)
	}
}

func (s *SceneSelectGame) Resize() {
	w0, h0 := eui.GetUi().Size()
	x, y := 0, 0
	h1 := int(float64(h0) * 0.05)
	s.topBar.Resize([]int{x, y, w0, h1})
	w := int(float64(w0) * 0.99)
	x = (w0 - w) / 2
	y += h1
	h2 := int(float64(h0) * 0.79)
	s.listGames.Resize([]int{x, y, w, h2})
	s.listGames.Itemsize(h1)
	x = 0
	h3 := int(float64(h0) * 0.15)
	y = h0 - h3
	s.btnsLayout.Resize([]int{x, y, w0, h3})
}
