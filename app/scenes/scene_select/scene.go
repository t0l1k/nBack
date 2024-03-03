package scene_select

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app/data"
	"github.com/t0l1k/nBack/app/db"
	"github.com/t0l1k/nBack/app/scenes/create"
	"github.com/t0l1k/nBack/app/scenes/intro"
	"github.com/t0l1k/nBack/app/scenes/options"
	"github.com/t0l1k/nBack/app/scenes/result"
	"github.com/t0l1k/nBack/app/scenes/tutor"
)

const (
	bSelect  = "Выбрать"
	bCreate  = "Создать"
	bResult  = "Итоги"
	bTutor   = "Помощь"
	bOptions = "Настройки"
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
	s.listGames.SetupListViewButtons(s.profiles.GetProfilesName(), 30, 1, bg, fg, s.btnsSelectGameLogic)
	s.Add(s.listGames)
	s.listGames.Bg(eui.Blue)
	s.btnsLayout = eui.NewHLayout()
	s.btnSel = eui.NewButton(bSelect, s.btnsLogic)
	s.btnSel.Disable()
	s.btnSel.Bg(eui.YellowGreen)
	s.btnsLayout.Add(s.btnSel)
	s.btnCrt = eui.NewButton(bCreate, s.btnsLogic)
	s.btnsLayout.Add(s.btnCrt)
	s.btnPrgs = eui.NewButton(bResult, s.btnsLogic)
	s.btnsLayout.Add(s.btnPrgs)
	s.btnTut = eui.NewButton(bTutor, s.btnsLogic)
	s.btnsLayout.Add(s.btnTut)
	s.btnOpt = eui.NewButton(bOptions, s.btnsLogic)
	s.btnsLayout.Add(s.btnOpt)
	s.Resize()
	return s
}

func (s *SceneSelectGame) btnsSelectGameLogic(b *eui.Button) {
	for name, game := range s.profiles.GetGameProfiles() {
		if name == b.GetText() {
			sc := intro.NewSceneIntro(game, name)
			eui.GetUi().Push(sc)
			log.Println("selected profile:", b.GetText())
		}
	}
}

func (s *SceneSelectGame) btnsLogic(b *eui.Button) {
	var sc eui.Sceneer
	switch b.GetText() {
	case bCreate:
		sc = create.NewSceneCreateGame(s.profiles)
		log.Println("Выбрана сцена создание профиля")
	case bResult:
		sc = result.NewSceneResults()
		log.Println("Выбрана сцена изучения итогов")
	case bTutor:
		sc = tutor.NewSceneTutor()
		log.Println("Выбрана сцена чтения помощи")
	case bOptions:
		sc = options.NewSceneOptions()
		log.Println("Выбрана сцена настроек")
	}
	eui.GetUi().Push(sc)
}

func (s *SceneSelectGame) Entered() {
	s.loadConf()

	s.listGames.Reset()
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(eui.ButtonBg)
	fg := theme.Get(eui.ButtonFg)
	s.listGames.SetupListViewButtons(s.profiles.GetProfilesName(), 30, 1, bg, fg, s.btnsSelectGameLogic)
	for k, v := range s.profiles.GetGameProfiles() {
		for _, v1 := range v.Games {
			if v1.IsDone() {
				fmt.Println(k, v1.LastGameFullResult())
			}
		}
	}
}

func (*SceneSelectGame) loadConf() {
	if conf := db.GetDb().GetFromDbAppConfData(); conf != nil {
		log.Print("___Get Conf from DB___:", conf)
		for k, v := range *conf {
			fmt.Println(k, v)
		}
	} else {
		log.Println("___New Conf___")
		db.GetDb().InsertAppConf()
	}
	ebiten.SetFullscreen(eui.GetUi().GetSettings().Get(eui.UiFullscreen).(bool))
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
