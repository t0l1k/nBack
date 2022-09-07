package game

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
	s.btnQuit = ui.NewButton("<", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetApp().Pop() })
	s.Add(s.btnQuit)
	s.lblName = ui.NewLabel(fmt.Sprintf("%v %v", ui.GetLocale().Get("scrName"), s.scorePeriod), rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblName)
	s.lblPeriodResult = ui.NewLabel("", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblPeriodResult)
	s.plotScore = NewScorePlot(rect)
	s.Add(s.plotScore)
	s.lblDt = ui.NewLabel(" ", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
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
	s.lblDt.SetText(ui.GetApp().UpdateUpTime())
	for _, value := range s.container {
		value.Update(dt)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		ui.GetApp().Push(NewSceneGame())
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		s.scorePeriod++
		if s.scorePeriod > all {
			s.scorePeriod = week
		}
		s.lblName.SetText(fmt.Sprintf("%v %v", ui.GetLocale().Get("scrName"), s.scorePeriod))
	}
}

func (s *SceneScore) Draw(surface *ebiten.Image) {
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneScore) Resize() {
	w, h := ui.GetApp().GetScreenSize()
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

func (s *SceneScore) Quit() {
	for _, v := range s.container {
		v.Close()
	}
}
