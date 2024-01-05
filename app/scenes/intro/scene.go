package scene_intro

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/data"
	scene_game "github.com/t0l1k/nBack/app/scenes/game"
)

type SceneIntro struct {
	eui.SceneBase
	topbar                            *eui.TopBar
	lblIntro, lblResult, lblStopwatch *eui.Text
	gamesData                         *data.GamesData
	restStopwatch                     *eui.Stopwatch
	restDuration                      int
	warningDuration                   time.Duration
	btnStart                          *eui.Button
	listShort                         *eui.ListView
}

func NewSceneIntro(gdata *data.GamesData, text string) *SceneIntro {
	s := &SceneIntro{}
	s.gamesData = gdata
	s.topbar = eui.NewTopBar(text)
	s.Add(s.topbar)
	s.listShort = eui.NewListView()
	s.Add(s.listShort)
	s.topbar.SetShowStopwatch()
	s.lblIntro = eui.NewText("Нажать <Space> для старта игры!")
	s.Add(s.lblIntro)
	s.lblResult = eui.NewText("")
	s.Add(s.lblResult)
	s.lblStopwatch = eui.NewText("")
	s.Add(s.lblStopwatch)
	s.btnStart = eui.NewButton("Начать новую сессию", func(b *eui.Button) {
		s.playNewGame()
	})
	s.Add(s.btnStart)
	s.restStopwatch = eui.NewStopwatch()
	return s
}

func (s *SceneIntro) Entered() {
	s.Resize()
	if s.gamesData.Last().IsDone() {
		level, lives, resultStr, colorStr := s.gamesData.NextLevel()
		s.lblIntro.SetText(s.gamesData.Last().String())
		s.lblIntro.Bg(colorStr)
		s.lblResult.SetText(resultStr)
		s.lblResult.Bg(colorStr)
		log.Println("new game", level, lives)
		s.warningDuration = s.gamesData.Last().Duration / 2
		s.gamesData.NewGame(level, lives)
	}
	s.restStopwatch.Start()
	var (
		strs     []string
		bgs, fgs []color.Color
	)
	for _, v := range s.gamesData.Games {
		if v.IsDone() {
			str, bg, fg := v.ShortResultStringWithColors()
			strs = append(strs, str)
			bgs = append(bgs, bg)
			fgs = append(fgs, fg)
			fmt.Println(v)
		}
	}
	s.listShort.Reset()
	s.listShort.SetListViewTextWithBgFgColors(strs, bgs, fgs)
	conf := eui.GetUi().GetSettings()
	s.restDuration = conf.Get(app.RestDuration).(int)
}

func (s *SceneIntro) Update(dt int) {
	for _, v := range s.Container {
		v.Update(dt)
	}
	s.lblStopwatch.SetText(s.restStopwatch.StringShort())
	if s.restStopwatch.Duration() < time.Duration(s.restDuration)*time.Second {
		if s.lblStopwatch.GetBg() != eui.Red && s.warningDuration > 0 {
			s.lblStopwatch.Bg(eui.Red)
		}
	} else if s.restStopwatch.Duration() < s.warningDuration {
		if s.lblStopwatch.GetBg() != eui.Orange {
			s.lblStopwatch.Bg(eui.Orange)
		}
	} else if s.restStopwatch.Duration() > s.warningDuration {
		if s.lblStopwatch.GetBg() != eui.Blue {
			s.lblStopwatch.Bg(eui.Blue)
		}
	}
}

func (s *SceneIntro) playNewGame() {
	sc := scene_game.New()
	sc.Setup(*s.gamesData.Conf, s.gamesData.Last())
	eui.GetUi().Push(sc)
	log.Println("new session started")
}

func (s *SceneIntro) Resize() {
	w0, h0 := eui.GetUi().Size()
	w1 := int(float64(w0) * 0.68)
	h1 := int(float64(h0) * 0.068)
	s.topbar.Resize([]int{0, 0, w0, h1 / 2})
	x, y := (w0-w1)/2, h0/2-h1
	s.lblIntro.Resize([]int{x, y, w1, h1})
	y += h1 + h1/2
	s.lblResult.Resize([]int{x, y, w1, h1})
	y += h1 + h1/2
	s.btnStart.Resize([]int{x, y, w1, h1})
	w := h1 * 2
	h := h1 * 2
	x = (w0 - h) / 2
	y = h0/2 - h*2
	s.lblStopwatch.Resize([]int{x, y, w, h})
	x, y = w0-h1*3, h1
	w, h = h1*3, h0-h1*2
	s.listShort.Resize([]int{x, y, w, h})
}
