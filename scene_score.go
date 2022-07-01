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
	name                     string
	lblName, lblPeriodResult *ui.Label
	plotScore                *ScorePlot
	scorePeriod              period
	rect                     *ui.Rect
	container                []ui.Drawable
}

func NewSceneScore() *SceneScore {
	s := &SceneScore{
		rect:        getApp().rect,
		scorePeriod: all,
	}
	rect := []int{0, 0, 1, 1}
	s.name = fmt.Sprintf("Games for the period %v", s.scorePeriod)
	s.lblName = ui.NewLabel(s.name, rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblName)
	s.lblPeriodResult = ui.NewLabel(getApp().db.todayData.String(), rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblPeriodResult)
	s.plotScore = NewScorePlot(rect)
	s.Add(s.plotScore)
	return s
}

func (s *SceneScore) Entered() {
	getApp().db.ReadAllGamesForScoresByDays()
	_, str := getApp().db.ReadAllGamesScore()
	s.lblPeriodResult.SetText(str)
	s.Resize()
	log.Println("Entered SceneScore")
}

func (s *SceneScore) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}

func (s *SceneScore) Update(dt int) {
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
	s.rect = getApp().rect
	x, y, w, h := 0, 0, int(float64(getApp().rect.W)*0.25), int(float64(getApp().rect.H)*0.05)
	s.lblName.Resize([]int{x, y, w, h})
	w, h = int(float64(getApp().rect.W)*0.9), int(float64(getApp().rect.H)*0.1)
	x, y = (s.rect.W-w)/2, int(float64(h)*0.8)
	s.lblPeriodResult.Resize([]int{x, y, w, h})
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.75)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*1.05)
	s.plotScore.Resize([]int{x, y, w, h})
}

func (s *SceneScore) Quit() {}
