package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

type period int

const (
	week period = iota
	month
	year
	all
)

func (p period) String() string {
	s := "all"
	switch p {
	case week:
		s = "week"
	case month:
		s = "month"
	case year:
		s = "year"
	default:
		s = "all"
	}
	return s
}

type SceneScore struct {
	name                            string
	lblName, lblPeriodResult, lblDt *ui.Label
	btnQuit                         *ui.Button
	plotScore                       *ScorePlot
	scorePeriod                     period
	rect                            *ui.Rect
	container                       []ui.Drawable
}

func NewSceneScore() *SceneScore {
	s := &SceneScore{
		rect:        ui.NewRect([]int{0, 0, 1, 1}),
		scorePeriod: all,
	}
	rect := []int{0, 0, 1, 1}
	s.btnQuit = ui.NewButton("<", rect, getTheme().CorrectColor, getTheme().Fg, func(b *ui.Button) { getApp().Pop() })
	s.Add(s.btnQuit)
	s.name = fmt.Sprintf("Games for the period %v", s.scorePeriod)
	s.lblName = ui.NewLabel(s.name, rect, getTheme().CorrectColor, getTheme().Fg)
	s.Add(s.lblName)
	s.lblPeriodResult = ui.NewLabel(getDb().todayData.String(), rect, getTheme().CorrectColor, getTheme().Fg)
	s.Add(s.lblPeriodResult)
	s.plotScore = NewScorePlot(rect)
	s.Add(s.plotScore)
	s.lblDt = ui.NewLabel("up: 00:00 ", rect, getTheme().CorrectColor, getTheme().Fg)
	s.Add(s.lblDt)
	return s
}

func (s *SceneScore) Entered() {
	getDb().ReadAllGamesForScoresByDays()
	_, str := getDb().ReadAllGamesScore()
	s.lblPeriodResult.SetText(str)
	s.Resize()
	log.Println("Entered SceneScore")
}

func (s *SceneScore) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}

func (s *SceneScore) Update(dt int) {
	s.lblDt.SetText(getApp().updateUpTime())
	for _, value := range s.container {
		value.Update(dt)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		getApp().Push(NewSceneGame())
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		s.scorePeriod++
		if s.scorePeriod > all {
			s.scorePeriod = week
		}
		s.lblName.SetText(fmt.Sprintf("Games for the period %v", s.scorePeriod))
	}
}

func (s *SceneScore) Draw(surface *ebiten.Image) {
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneScore) Resize() {
	w, h := getApp().GetScreenSize()
	s.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(s.rect.H)*0.05), int(float64(s.rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, int(float64(s.rect.W)*0.20)
	s.lblName.Resize([]int{x, y, w, h})
	x = s.rect.Right() - w
	s.lblDt.Resize([]int{x, y, w, h})
	w = int(float64(s.rect.W) * 0.9)
	x, y = (s.rect.W-w)/2, int(float64(h)*1.2)
	s.lblPeriodResult.Resize([]int{x, y, w, h})
	y = int(float64(h) * 2.4)
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.85)
	x = (s.rect.W - w) / 2
	s.plotScore.Resize([]int{x, y, w, h})
}

func (s *SceneScore) Quit() {}
