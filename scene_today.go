package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

type SceneToday struct {
	name                     string
	lblName, lblPeriodResult *ui.Label
	lblsResult               *ResultLbls
	rect                     *ui.Rect
	container                []ui.Drawable
}

func NewSceneToday() *SceneToday {
	return &SceneToday{
		rect: getApp().rect,
	}
}

func (s *SceneToday) Entered() {
	getApp().db.ReadTodayGames()
	x, y, w, h := 0, 0, int(float64(getApp().rect.W)*0.3), int(float64(getApp().rect.H)*0.05)
	s.name = "Games for Today"
	s.lblName = ui.NewLabel(s.name, []int{x, y, w, h})
	s.Add(s.lblName)
	w, h = int(float64(getApp().rect.W)*0.9), int(float64(getApp().rect.H)*0.1)
	x, y = (s.rect.W-w)/2, int(float64(h)*0.8)
	s.lblPeriodResult = ui.NewLabel(getApp().db.todayData.String(), []int{x, y, w, h})
	s.Add(s.lblPeriodResult)
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.75)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*1.05)
	s.lblsResult = NewResultLbls([]int{x, y, w, h})
	s.Add(s.lblsResult)
	log.Println("Eneterd SceneToday")
}
func (s *SceneToday) Quit() {}
func (s *SceneToday) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}
func (s *SceneToday) Update(dt int) {
	for _, value := range s.container {
		value.Update(dt)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		getApp().Push(NewSceneGame())
	}
}
func (s *SceneToday) Draw(surface *ebiten.Image) {
	for _, value := range s.container {
		value.Draw(surface)
	}
}
