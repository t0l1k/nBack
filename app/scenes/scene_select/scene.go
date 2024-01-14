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
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 2)
			conf.Set(data.ThresholdAdvance, 80)
			conf.Set(data.ThresholdFallback, 50)
			conf.Set(data.ThresholdFallbackSessions, 3)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Pos}, conf)
			return g
		}(),
		"Single nBack Colors rulez brainworkshop": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 2)
			conf.Set(data.ThresholdAdvance, 80)
			conf.Set(data.ThresholdFallback, 50)
			conf.Set(data.ThresholdFallbackSessions, 3)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Col}, conf)
			return g
		}(),
		"Single nBack Numbers rulez brainworkshop": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 2)
			conf.Set(data.ThresholdAdvance, 80)
			conf.Set(data.ThresholdFallback, 50)
			conf.Set(data.ThresholdFallbackSessions, 3)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Sym}, conf)
			return g
		}(),
		"Single nBack Ariphmetics rulez brainworkshop": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 2)
			conf.Set(data.ThresholdAdvance, 80)
			conf.Set(data.ThresholdFallback, 50)
			conf.Set(data.ThresholdFallbackSessions, 3)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Ari}, conf)
			return g
		}(),
		"Dual nBack Position(3x3), Colors rulez brainworkshop": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 2)
			conf.Set(data.ThresholdAdvance, 80)
			conf.Set(data.ThresholdFallback, 50)
			conf.Set(data.ThresholdFallbackSessions, 3)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 2.5)
			g := data.NewGamesData([]string{data.Pos, data.Col}, conf)
			return g
		}(),
		"Triple nBack Position(3x3), Numbers, Colors,rulez brainworkshop": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 2)
			conf.Set(data.ThresholdAdvance, 80)
			conf.Set(data.ThresholdFallback, 50)
			conf.Set(data.ThresholdFallbackSessions, 3)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 3.0)
			g := data.NewGamesData([]string{data.Pos, data.Sym, data.Col}, conf)
			return g
		}(),
		"Single nBack Position(3x3) Jaeggi Rulez": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, false)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Pos}, conf)
			return g
		}(),
		"Single nBack Numbers Jaeggi Rulez": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, false)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Sym}, conf)
			return g
		}(),
		"Single nBack Colors Jaeggi Rulez": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, false)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Col}, conf)
			return g
		}(),
		"Single nBack Ariphmetics Jaeggi Rulez": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, false)
			conf.Set(data.MoveTime, 2.0)
			g := data.NewGamesData([]string{data.Ari}, conf)
			return g
		}(),
		"Dual nBack Position(3x3), Color Jaeggi Rulez": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, false)
			conf.Set(data.MoveTime, 2.5)
			g := data.NewGamesData([]string{data.Pos, data.Col}, conf)
			return g
		}(),
		"Dual nBack Position(3x3), Numbers Jaeggi Rulez": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, false)
			conf.Set(data.MoveTime, 2.5)
			g := data.NewGamesData([]string{data.Pos, data.Sym}, conf)
			return g
		}(),
		"Triple nBack Position(3x3), Number, Color Jaeggi Rulez": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, false)
			conf.Set(data.MoveTime, 3.0)
			g := data.NewGamesData([]string{data.Pos, data.Sym, data.Col}, conf)
			return g
		}(),
		"Гадкий утёнок позиции(3x3) легко": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 0)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.5)
			g := data.NewGamesData([]string{data.Pos}, conf)
			return g
		}(),
		"Гадкий утёнок цифры легко(ход 1 сек)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 0)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.0)
			g := data.NewGamesData([]string{data.Pos}, conf)
			return g
		}(),
		"Гадкий утёнок цифры цвет легко(ход 2 сек)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 20)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 0)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 2.0)
			g := data.NewGamesData([]string{data.Sym, data.Col}, conf)
			return g
		}(),
		// "Три поросёнка позиции(3x3) легко",
		"Devel test set0 pos (move 1 sec)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 5)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.0)
			g := data.NewGamesData([]string{data.Pos}, conf)
			return g
		}(),
		"Devel test set1 sym (move 1 sec)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 5)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 1.0)
			g := data.NewGamesData([]string{data.Sym}, conf)
			return g
		}(),
		"Devel test set2 dual pos/sym (move 2 sec)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 5)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 2.0)
			g := data.NewGamesData([]string{data.Pos, data.Sym}, conf)
			return g
		}(),
		"Devel test set3 dual pos/col (move 2 sec)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 5)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 2.0)
			g := data.NewGamesData([]string{data.Pos, data.Col}, conf)
			return g
		}(),
		"Devel test set4 dual sym/col (move 2 sec)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 5)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 2.0)
			g := data.NewGamesData([]string{data.Sym, data.Col}, conf)
			return g
		}(),
		"Devel test set5 triple pos/sym/col (move 3 sec)": func() *data.GamesData {
			conf := data.DefaultSettings()
			conf.Set(data.Trials, 5)
			conf.Set(data.TrialsFactor, 1)
			conf.Set(data.TrialsExponent, 1)
			conf.Set(data.ThresholdAdvance, 90)
			conf.Set(data.ThresholdFallback, 75)
			conf.Set(data.ThresholdFallbackSessions, 1)
			conf.Set(data.GridSize, 3)
			conf.Set(data.ShowGrid, true)
			conf.Set(data.MoveTime, 3.0)
			g := data.NewGamesData([]string{data.Pos, data.Sym, data.Col}, conf)
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
			log.Println("selected profile:", b.GetText())
		}
	}
}

func (s *SceneSelectGame) Entered() {
	for k, v := range s.lst {
		for _, v1 := range v.Games {
			if v1.IsDone() {
				fmt.Println(k, v1.LastGameFullResult())
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
