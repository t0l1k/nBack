package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

type SceneToday struct {
	name                            string
	lblName, lblPeriodResult, lblDt *ui.Label
	lblsResult                      *ResultLbls
	plotResult                      *ResultPlot
	toggleResults                   bool
	rect                            *ui.Rect
	container                       []ui.Drawable
}

func NewSceneToday() *SceneToday {
	s := &SceneToday{
		rect: getApp().rect,
	}
	rect := []int{0, 0, 1, 1}
	s.name = "Games for Today"
	s.lblName = ui.NewLabel(s.name, rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblName)
	s.lblPeriodResult = ui.NewLabel(getApp().db.todayData.String(), rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblPeriodResult)
	s.lblDt = ui.NewLabel("up: 00:00 ", rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblDt)
	s.lblsResult = NewResultLbls(rect)
	s.Add(s.lblsResult)
	s.plotResult = NewResultPlot(rect)
	s.plotResult.Visibe = false
	s.Add(s.plotResult)
	s.toggleResults = false
	return s
}

func (s *SceneToday) Entered() {
	getApp().db.ReadTodayGames()
	s.lblPeriodResult.SetText(getApp().db.todayData.String())
	s.Resize()
	log.Println("Entered SceneToday")
}
func (s *SceneToday) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}
func (s *SceneToday) Update(dt int) {
	s.updateDt()
	for _, value := range s.container {
		value.Update(dt)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		getApp().Push(NewSceneGame())
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		s.toggleResults = !s.toggleResults
		if s.toggleResults {
			s.plotResult.Dirty = true
			s.plotResult.Visibe = true
			s.lblsResult.Visibe = false
		} else {
			s.plotResult.Visibe = false
			s.lblsResult.Visibe = true
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		getApp().Push(NewSceneScore())
	}
}

func (s *SceneToday) updateDt() {
	durration := time.Since(getApp().startDt)
	d := durration.Round(time.Second)
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	sec := d / time.Second
	result := ""
	if hours > 0 {
		result = fmt.Sprintf("%02v:%02v:%02v", int(hours), int(minutes), int(sec))
	} else {
		result = fmt.Sprintf("%02v:%02v", int(minutes), int(sec))
	}
	ss := fmt.Sprintf("up: %v", result)
	s.lblDt.SetText(ss)
}

func (s *SceneToday) Draw(surface *ebiten.Image) {
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneToday) Resize() {
	s.rect = getApp().rect
	x, y, w, h := 0, 0, int(float64(getApp().rect.W)*0.25), int(float64(getApp().rect.H)*0.05)
	s.lblName.Resize([]int{x, y, w, h})
	x = s.rect.Right() - w
	s.lblDt.Resize([]int{x, y, w, h})
	w, h = int(float64(getApp().rect.W)*0.9), int(float64(getApp().rect.H)*0.1)
	x, y = (s.rect.W-w)/2, int(float64(h)*0.8)
	s.lblPeriodResult.Resize([]int{x, y, w, h})
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.75)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*1.05)
	s.lblsResult.Resize([]int{x, y, w, h})
	s.plotResult.Resize([]int{x, y, w, h})
}

func (s *SceneToday) Quit() {}
