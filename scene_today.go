package main

import (
	"image/color"
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
	rect := []int{0, 0, 1, 1}
	s.name = "Games for Today"
	s.lblName = ui.NewLabel(s.name, rect)
	s.Add(s.lblName)
	s.lblPeriodResult = ui.NewLabel(getApp().db.todayData.String(), rect)
	s.Add(s.lblPeriodResult)
	s.lblsResult = NewResultLbls(rect)
	s.Add(s.lblsResult)
	s.Resize()
	log.Println("Eneterd SceneToday")
}
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
	surface.Fill(color.RGBA{0, 0, 0, 255})
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneToday) Resize() {
	s.rect = getApp().rect
	x, y, w, h := 0, 0, int(float64(getApp().rect.W)*0.3), int(float64(getApp().rect.H)*0.05)
	s.lblName.Resize([]int{x, y, w, h})
	w, h = int(float64(getApp().rect.W)*0.9), int(float64(getApp().rect.H)*0.1)
	x, y = (s.rect.W-w)/2, int(float64(h)*0.8)
	s.lblPeriodResult.Resize([]int{x, y, w, h})
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.75)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*1.05)
	s.lblsResult.Resize([]int{x, y, w, h})
}

func (s *SceneToday) Quit() {}
