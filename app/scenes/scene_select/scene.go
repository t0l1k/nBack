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
	lst                                                        map[string]*data.GamesData
}

func NewSceneSelectGame() *SceneSelectGame {
	s := &SceneSelectGame{}
	s.topBar = eui.NewTopBar("нНазад")
	s.topBar.SetShowStopwatch()
	s.Add(s.topBar)
	s.lst = map[string]*data.GamesData{
		"Single nBack Position(3x3) rulez brainworkshop": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 2)
			g.Conf.Set(data.ThresholdFallbackSessions, 3)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, true)
			g.Setup([]string{data.Pos}, 1, 3, 80, 50, 1.5)
			return g
		}(),
		"Single nBack Colors rulez brainworkshop": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 2)
			g.Conf.Set(data.ThresholdFallbackSessions, 3)
			g.Setup([]string{data.Col}, 1, 3, 80, 50, 1.5)
			return g
		}(),
		"Single nBack Numbers rulez brainworkshop": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 2)
			g.Conf.Set(data.ThresholdFallbackSessions, 3)
			g.Setup([]string{data.Sym}, 1, 3, 80, 50, 1.5)
			return g
		}(),
		"Single nBack Ariphmetics rulez brainworkshop": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 2)
			g.Conf.Set(data.ThresholdFallbackSessions, 3)
			g.Setup([]string{data.Ari}, 1, 3, 80, 50, 1.5)
			return g
		}(),
		"Dual nBack Position(3x3), Colors rulez brainworkshop": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 2)
			g.Conf.Set(data.ThresholdFallbackSessions, 3)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, true)
			g.Setup([]string{data.Pos, data.Col}, 1, 3, 80, 50, 2.5)
			return g
		}(),
		"Triple nBack Position(3x3), Colors, Numbers rulez brainworkshop": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 2)
			g.Conf.Set(data.ThresholdFallbackSessions, 3)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, true)
			g.Setup([]string{data.Pos, data.Sym, data.Col}, 1, 3, 80, 50, 3.0)
			return g
		}(),
		"Single nBack Position(3x3) Jaeggi Rulez": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Pos}, 1, 1, 90, 75, 1.5)
			return g
		}(),
		"Single nBack Numbers Jaeggi Rulez": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Sym}, 1, 1, 90, 75, 1.5)
			return g
		}(),
		"Single nBack Colors Jaeggi Rulez": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Col}, 1, 1, 90, 75, 1.5)
			return g
		}(),
		"Single nBack Ariphmetics Jaeggi Rulez": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Ari}, 1, 1, 90, 75, 1.5)
			return g
		}(),
		"Dual nBack Position(3x3), Color Jaeggi Rulez": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Pos, data.Col}, 1, 1, 90, 75, 2.0)
			return g
		}(),
		"Dual nBack Position(3x3), Numbers Jaeggi Rulez": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Pos, data.Sym}, 1, 1, 90, 75, 2.0)
			return g
		}(),
		"Triple nBack Position(3x3), Number, Color Jaeggi Rulez": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Pos, data.Sym, data.Col}, 1, 1, 90, 75, 3)
			return g
		}(),
		"Гадкий утёнок позиции(3x3) легко": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Pos}, 1, 3, 90, 0, 1.5)
			return g
		}(),
		"Гадкий утёнок цифры легко": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 20)
			g.Conf.Set(data.TrialsFactor, 1)
			g.Conf.Set(data.TrialsExponent, 1)
			g.Conf.Set(data.ThresholdFallbackSessions, 1)
			g.Conf.Set(data.GridSize, 3)
			g.Conf.Set(data.ShowGrid, false)
			g.Setup([]string{data.Sym}, 1, 1, 90, 0, 1.5)
			return g
		}(),

		// "Три поросёнка позиции(3x3) легко",
		"Devel test set sym": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 5)
			g.Setup([]string{data.Sym}, 1, 1, 90, 75, 2.0)
			return g
		}(),
		"Devel test set2 dual pos/sym": func() *data.GamesData {
			g := data.NewGamesData()
			g.Conf.Set(data.Trials, 5)
			g.Setup([]string{data.Sym, data.Pos}, 1, 1, 90, 75, 2.0)
			return g
		}(),
	}
	s.listGames = eui.NewListView()
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(eui.ButtonBg)
	fg := theme.Get(eui.ButtonFg)
	var lst []string
	for k := range s.lst {
		lst = append(lst, k)
	}
	s.listGames.SetupListViewButtons(lst, 30, 1, bg, fg, s.btnsLogic)
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
	for k, v := range s.lst {
		if k == b.GetText() {
			sc := scene_intro.NewSceneIntro(v, k)
			eui.GetUi().Push(sc)
			log.Println("new session started", b.GetText())
		}
	}
}

func (s *SceneSelectGame) Entered() {
	for k, v := range s.lst {
		for _, v1 := range v.Games {
			if v1.IsDone() {
				fmt.Println(k, v1)
			}
		}
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
