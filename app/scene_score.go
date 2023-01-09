package app

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
)

type SceneScore struct {
	ui.ContainerDefault
	lblName, lblPeriodResult, lblDt *ui.Label
	btnQuit                         *ui.Button
	plotScore                       *ScorePlot
	scorePeriod                     data.Period
	rect                            *ui.Rect
	mn                              int
}

func NewSceneScore() *SceneScore {
	s := &SceneScore{
		rect:        ui.NewRect([]int{0, 0, 1, 1}),
		scorePeriod: data.All,
	}
	rect := []int{0, 0, 1, 1}
	s.btnQuit = ui.NewButton("<", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) { ui.Pop() })
	s.Add(s.btnQuit)
	s.lblName = ui.NewLabel(fmt.Sprintf("%v %v", ui.GetLocale().Get("scrName"), s.scorePeriod), rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblName)
	s.lblPeriodResult = ui.NewLabel("", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblPeriodResult)
	s.plotScore = NewScorePlot(rect)
	s.Add(s.plotScore)
	s.lblDt = ui.NewLabel(" ", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblDt)
	return s
}

func (s *SceneScore) Entered() {
	s.checkPeriod(s.mn)
	s.Resize()
	log.Println("Entered SceneScore")
}

func (s *SceneScore) Update(dt int) {
	for _, value := range s.Container {
		value.Update(dt)
	}
	s.lblDt.SetText(ui.GetUi().UpdateUpTime())
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		ui.Push(NewSceneGame())
	} else if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		s.mn = 0
		s.scorePeriod.Next()
		s.plotScore.SetPeriod(s.scorePeriod)
		s.lblName.SetText(fmt.Sprintf("%v %v", ui.GetLocale().Get("scrName"), s.scorePeriod))
		s.checkPeriod(s.mn)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		if s.checkPeriod(s.mn + 1) {
			s.mn++
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		if s.checkPeriod(s.mn - 1) {
			s.mn--
		}
	}
}

func (s *SceneScore) checkPeriod(mn int) bool {
	var (
		from, to string
		result   bool = true
	)
	switch s.scorePeriod {
	case data.Week:
		from, to, result = data.NextWeek(mn)
	case data.Month:
		from, to, result = data.NextMonth(mn)
	case data.Year:
		from, to, result = data.NextYear(mn)
	case data.All:
	}
	if !result {
		return result
	}
	s.updateData(from, to)
	s.plotScore.Dirty = true
	return result
}

func (s *SceneScore) updateData(from, to string) {
	data.GetDb().ReadAllGamesForScoresByDays(s.scorePeriod.Len(), from, to)
	_, str := data.GetDb().ReadAllGamesScore(s.scorePeriod.Len(), from, to)
	s.lblPeriodResult.SetText(str)
}

func (s *SceneScore) Draw(surface *ebiten.Image) {
	for _, value := range s.Container {
		value.Draw(surface)
	}
}

func (s *SceneScore) Resize() {
	w, h := ui.GetUi().GetScreenSize()
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

func (s *SceneScore) Close() {
	for _, v := range s.Container {
		v.Close()
	}
}
