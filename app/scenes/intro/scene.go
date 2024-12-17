package intro

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/data"
	scene_game "github.com/t0l1k/nBack/app/scenes/game"
)

type SceneIntro struct {
	eui.SceneBase
	topbar                                           *eui.TopBar
	lblIntro, lblMotto, lblSw, lblResults, lblHelper *eui.Text
	gamesData                                        *data.GamesData
	restStopwatch                                    *eui.Stopwatch
	restDuration                                     int
	warningDuration                                  time.Duration
	btnStart                                         *eui.Button
	listShort, listConf                              *eui.ListView
	movesLine                                        *MovesLine
	movesIcon                                        *eui.Icon
	plot                                             *eui.Plot
}

func NewSceneIntro(gdata *data.GamesData, text string) *SceneIntro {
	s := &SceneIntro{}
	s.gamesData = gdata
	s.topbar = eui.NewTopBar(text, nil)
	s.topbar.SetTitleCoverArea(0.9)
	s.topbar.SetStopwatchCoverArea(0.1)
	s.Add(s.topbar)
	s.listShort = eui.NewListView()
	s.Add(s.listShort)
	s.listConf = eui.NewListView()
	s.Add(s.listConf)
	s.topbar.SetUseStopwatch()
	s.topbar.SetShowStoppwatch(true)
	s.lblIntro = eui.NewText("")
	s.Add(s.lblIntro)
	s.lblMotto = eui.NewText("")
	s.Add(s.lblMotto)
	s.lblSw = eui.NewText("")
	s.Add(s.lblSw)
	s.lblResults = eui.NewText("")
	s.Add(s.lblResults)
	s.lblHelper = eui.NewText("[Modality Correct-(Wrong-Missed)]")
	s.Add(s.lblHelper)
	s.btnStart = eui.NewButton("Начать новую сессию", func(b *eui.Button) {
		s.playNewGame()
	})
	s.Add(s.btnStart)
	s.movesIcon = eui.NewIcon(nil)
	s.Add(s.movesIcon)
	s.movesIcon.Visible(false)
	s.lblIntro.Visible(false)
	s.lblMotto.Visible(false)
	s.lblHelper.Visible(false)
	s.movesLine = NewMovesLine()
	s.restStopwatch = eui.NewStopwatch()
	xArr := []int{0}
	yArr := []int{0}
	vArr := []int{0}
	s.plot = eui.NewPlot(xArr, yArr, vArr, "Score", "Level", "Game")
	s.Add(s.plot)
	return s
}

func (s *SceneIntro) Entered() {
	s.Resize()
	if s.gamesData.Last().IsDone() {
		s.lblIntro.Visible(true)
		s.lblMotto.Visible(true)
		s.lblHelper.Visible(true)
		level, tryUp, tryDown, mottoStr, colorStr := s.gamesData.NextLevel()
		s.lblIntro.SetText(s.gamesData.Last().LastGameFullResult())
		s.lblIntro.Bg(colorStr)
		s.lblMotto.SetText(mottoStr)
		s.lblMotto.Bg(colorStr)
		s.lblHelper.Bg(colorStr)
		s.movesLine.Setup(s.gamesData.Last().GetModalitiesMoves(), s.gamesData.Last().GetModalitiesScore())
		s.movesIcon.SetIcon(s.movesLine.Image())
		s.movesIcon.Visible(true)

		log.Println("new game", level, tryDown)
		s.warningDuration = s.gamesData.Last().Duration / 2
		s.gamesData.NewGame(level, tryUp, tryDown)
	} else {
		if len(s.gamesData.Games) > 1 {
			_, _, _, mottoStr, colorStr := s.gamesData.PrevGame()
			s.lblMotto.Visible(true)
			s.lblMotto.SetText(mottoStr)
			s.lblMotto.Bg(colorStr)
		}
	}

	rect := s.plot.GetRect().GetArr()
	s.plot = nil
	xArr, yArr, vArr := s.gamesData.GetPlotData()
	s.plot = eui.NewPlot(xArr, yArr, vArr, "Score", "Game", "Level")
	s.plot.Resize(rect)
	s.Add(s.plot)

	s.lblResults.SetText(s.gamesData.String())
	var (
		strs     []string
		bgs, fgs []color.Color
	)
	for i := len(s.gamesData.Games) - 1; i >= 0; i-- {
		v := s.gamesData.Games[i]
		if v.IsDone() {
			str, bg, fg := v.ShortResultStringWithColors()
			strs = append(strs, str)
			bgs = append(bgs, bg)
			fgs = append(fgs, fg)
			fmt.Println(v.LastGameFullResult())
		}
	}
	s.listShort.Reset()
	s.listShort.SetListViewTextWithBgFgColors(strs, bgs, fgs)
	s.listConf.Reset()
	s.listConf.SetupListViewText(s.gamesData.Conf.GameConf(s.gamesData.Games[s.gamesData.Id()]), 30, 1, colors.Teal, colors.Yellow)

	conf := eui.GetUi().GetSettings()
	s.restDuration = conf.Get(app.RestDuration).(int)
	s.restStopwatch.Start()
}

func (s *SceneIntro) Update(dt int) {
	for _, v := range s.GetContainer() {
		v.Update(dt)
	}
	s.lblSw.SetText(s.restStopwatch.StringShort())
	if s.restStopwatch.Duration() < time.Duration(s.restDuration)*time.Second {
		if s.lblSw.GetBg() != colors.Red && s.warningDuration > 0 {
			s.lblSw.Bg(colors.Red)
			s.btnStart.Visible(false)
		}
	} else if s.restStopwatch.Duration() < s.warningDuration {
		if s.lblSw.GetBg() != colors.Orange {
			s.lblSw.Bg(colors.Orange)
			s.btnStart.Visible(true)
		}
	} else if s.restStopwatch.Duration() > s.warningDuration {
		if s.lblSw.GetBg() != colors.Blue {
			s.lblSw.Bg(colors.Blue)
			s.btnStart.Visible(true)
		}
	}
}

func (s *SceneIntro) playNewGame() {
	log.Println("new session start")
	sc := scene_game.New()
	sc.Setup(*s.gamesData.Conf, s.gamesData.Last())
	eui.GetUi().Push(sc)
}

func (s *SceneIntro) Draw(surface *ebiten.Image) {
	for _, v := range s.GetContainer() {
		v.Draw(surface)
	}
}

func (s *SceneIntro) Resize() {
	w0, h0 := eui.GetUi().Size()
	w1 := int(float64(w0) * 0.68)
	h1 := int(float64(h0) * 0.068)
	s.topbar.Resize([]int{0, 0, w0, h1})
	x, y := (w0-w1)/2, h1
	s.lblResults.Resize([]int{x, y, w1, h1})
	y += h1
	s.plot.Resize([]int{x, y, w1, h1*7 + h1/2})
	y = h0 - h1
	s.btnStart.Resize([]int{x, y, w1, h1})
	w := h1 * 2
	h := h1
	x = (w0 - h) / 2
	y -= h
	s.lblSw.Resize([]int{x, y, w, h})

	x = (w0 - w1) / 2
	y -= h1
	s.movesIcon.Resize([]int{x, y, w1, h1})
	s.movesLine.Resize([]int{x, y, w1, h1})

	y -= h1
	s.lblIntro.Resize([]int{x, y, w1, h1})
	y -= h1 / 2
	s.lblHelper.Resize([]int{x, y, w1, h1 / 2})

	y -= h1
	s.lblMotto.Resize([]int{x, y, w1, h1})

	x, y = w0-(w0-w1)/2, h1
	w, h = (w0-w1)/2, h0-h1*2
	s.listShort.Resize([]int{x, y, w, h})

	x = 0
	s.listConf.Resize([]int{x, y, w, h})
}
