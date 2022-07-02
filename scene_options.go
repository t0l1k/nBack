package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

type SceneOptions struct {
	name      string
	rect      *ui.Rect
	container []ui.Drawable
	lblName   *ui.Label
	optGame   *OptGame
	optTheme  *OptTheme
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{
		rect: getApp().rect,
	}
	rect := []int{0, 0, 1, 1}
	s.name = "N-Back Options"
	s.lblName = ui.NewLabel(s.name, rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblName)
	s.optGame = (NewOptGame)(rect)
	s.Add(s.optGame)
	s.optTheme = NewOptTheme(rect)
	s.Add(s.optTheme)
	return s
}

func (s *SceneOptions) Entered() {
	s.Resize()
	log.Println("Entered SceneOptions")
}

func (s *SceneOptions) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}

func (s *SceneOptions) Update(dt int) {
	for _, value := range s.container {
		value.Update(dt)
	}
}

func (s *SceneOptions) Draw(surface *ebiten.Image) {
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneOptions) Resize() {
	s.rect = getApp().rect
	x, y, w, h := 0, 0, int(float64(getApp().rect.W)*0.25), int(float64(getApp().rect.H)*0.05)
	s.lblName.Resize([]int{x, y, w, h})
	y = s.rect.H - (h * 3)
	w, h1 := s.rect.W, h*3
	s.optTheme.Resize([]int{x, y, w, h1})
	y = h
	w, h2 := s.rect.W, h*9
	s.optGame.Resize([]int{x, y, w, h2})
}

func (s *SceneOptions) Quit() {}
