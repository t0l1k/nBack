package score

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/ui/scene/game"
	"github.com/t0l1k/nBack/ui/scene/plot"
)

type SceneScore struct {
	eui.ContainerDefault
	lblName, lblPeriodResult, lblDt *eui.Label
	btnQuit                         *eui.Button
	plotScore                       *plot.ScorePlot
	plotResult                      *plot.ResultPlot
	scorePeriod                     data.Period
	rect                            *eui.Rect
	mn                              int
}

func NewSceneScore() *SceneScore {
	s := &SceneScore{
		rect:        eui.NewRect([]int{0, 0, 1, 1}),
		scorePeriod: data.All,
	}
	rect := []int{0, 0, 1, 1}
	s.btnQuit = eui.NewButton("<", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { eui.Pop() })
	s.Add(s.btnQuit)
	s.lblName = eui.NewLabel(fmt.Sprintf("%v %v", eui.GetLocale().Get("scrName"), s.scorePeriod), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblName)
	s.lblPeriodResult = eui.NewLabel("", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblPeriodResult)
	s.plotScore = plot.NewScorePlot(rect)
	s.Add(s.plotScore)
	s.lblDt = eui.NewLabel(" ", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblDt)
	s.plotResult = plot.NewResultPlot(rect)
	s.plotResult.Visible = false
	s.Add(s.plotResult)
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
	s.lblDt.SetText(eui.GetUi().UpdateUpTime())
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		eui.Push(game.NewSceneGame())
	} else if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		s.mn = 0
		s.scorePeriod.Next()
		s.plotScore.SetPeriod(s.scorePeriod)
		s.plotResult.SetPeriod(s.scorePeriod)
		s.lblName.SetText(fmt.Sprintf("%v %v", eui.GetLocale().Get("scrName"), s.scorePeriod))
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
	case data.Day:
		from, to, result = data.NextDay(mn)
		s.plotResult.Visible = true
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
	s.plotResult.Dirty = true
	return result
}

func (s *SceneScore) updateData(from, to string) {
	data.GetDb().ReadAllGamesForScoresByDays(s.scorePeriod.Len(), from, to)
	_, str := data.GetDb().ReadAllGamesScore(s.scorePeriod.Len(), from, to)
	if s.scorePeriod == data.Day {
		data.GetDb().ReadTodayGames(from)
		str = data.GetDb().TodayData.String()
		if len(str) <= 10 {
			str = from
		}
	}
	s.lblPeriodResult.SetText(str)
}

func (s *SceneScore) Draw(surface *ebiten.Image) {
	for _, value := range s.Container {
		value.Draw(surface)
	}
}

func (s *SceneScore) Resize() {
	w, h := eui.GetUi().GetScreenSize()
	s.rect = eui.NewRect([]int{0, 0, w, h})
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
	s.plotResult.Resize([]int{x, y, w, h})
}

func (s *SceneScore) Close() {
	for _, v := range s.Container {
		v.Close()
	}
}
