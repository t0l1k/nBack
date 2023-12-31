package scene_select

import (
	"fmt"
	"log"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app/data"
	scene_intro "github.com/t0l1k/nBack/app/scenes/intro"
)

type SceneSelectGame struct {
	eui.SceneBase
	topBar                                                     *eui.TopBar
	listGames                                                  *eui.ListView
	btnSelect, btnCreate, btnProgress, btnTutorial, btnOptions *eui.Button
	btnsLayout                                                 *eui.BoxLayout
	lst                                                        []string
}

func NewSceneSelectGame() *SceneSelectGame {
	s := &SceneSelectGame{}
	s.topBar = eui.NewTopBar("нНазад")
	s.topBar.SetShowStopwatch()
	s.Add(s.topBar)
	s.lst = []string{
		"Single nBack Position(3x3) rulez brainworkshop",
		"Single nBack Position(3x3) Jaeggi Rulez",
		"Single nBack Numbers Jaeggi Rulez",
		"Single nBack Colors Jaeggi Rulez",
		"Single nBack Ariphmetics Jaeggi Rulez",
		"Dual nBack Position(3x3), Color Jaeggi Rulez",
		"Гадкий утёнок позиции(3x3) легко",
		"Три поросёнка позиции(3x3) легко",
		"Devel test set"}
	s.listGames = eui.NewListView()
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(eui.ButtonBg)
	fg := theme.Get(eui.ButtonFg)
	s.listGames.SetupListViewButtons(s.lst, 30, 1, bg, fg, s.btnsLogic)
	s.Add(s.listGames)
	s.listGames.Bg(eui.Blue)
	s.btnsLayout = eui.NewHLayout()
	s.btnSelect = eui.NewButton("Играть", s.btnsLogic)
	s.btnSelect.Disable()
	s.btnSelect.Bg(eui.YellowGreen)
	s.btnsLayout.Add(s.btnSelect)
	s.btnCreate = eui.NewButton("Создать", s.btnsLogic)
	s.btnsLayout.Add(s.btnCreate)
	s.btnProgress = eui.NewButton("Итоги", s.btnsLogic)
	s.btnsLayout.Add(s.btnProgress)
	s.btnTutorial = eui.NewButton("Обучение", s.btnsLogic)
	s.btnsLayout.Add(s.btnTutorial)
	s.btnOptions = eui.NewButton("Настройки", s.btnsLogic)
	s.btnsLayout.Add(s.btnOptions)
	s.Add(s.btnsLayout)
	s.Resize()
	return s
}

func (s *SceneSelectGame) btnsLogic(b *eui.Button) {
	fmt.Println("selected", b.GetText())
	if b.GetText() == s.lst[0] {
		sc := scene_intro.NewSceneIntro(data.NewGamePos3x3BRRulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[1] {
		sc := scene_intro.NewSceneIntro(data.NewGameJaeggiPos3x3Rulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[2] {
		sc := scene_intro.NewSceneIntro(data.NewGameJaeggiSymRulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[3] {
		sc := scene_intro.NewSceneIntro(data.NewGameJaeggiColRulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[4] {
		sc := scene_intro.NewSceneIntro(data.NewGameJaeggiAriRulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[5] {
		sc := scene_intro.NewSceneIntro(data.NewGameJaeggiPos3x3ColRulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[6] {
		sc := scene_intro.NewSceneIntro(data.NewGameUngleDuckPos3x3Rulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[7] {
		sc := scene_intro.NewSceneIntro(data.NewGameDevelRulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
	} else if b.GetText() == s.lst[8] {
		sc := scene_intro.NewSceneIntro(data.NewGameDevelRulez())
		eui.GetUi().Push(sc)
		log.Println("new session started", b.GetText())
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
