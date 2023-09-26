package today

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/ui/scene/game"
	"github.com/t0l1k/nBack/ui/scene/options"
	"github.com/t0l1k/nBack/ui/scene/plot"
	"github.com/t0l1k/nBack/ui/scene/score"
)

type SceneToday struct {
	lblName, lblPeriodResult, lblDt, lblHelper                  *eui.Label
	btnScore, btnStart, btnQuit, btnPlot, btnFullScreen, btnOpt *eui.Button
	lblsResult                                                  *eui.List
	plotResult                                                  *plot.ResultPlot
	toggleResults                                               bool
	rect                                                        *eui.Rect
	eui.ContainerDefault
}

func NewSceneToday() *SceneToday {
	s := &SceneToday{
		rect: eui.NewRect([]int{0, 0, 1, 1}),
	}
	rect := []int{0, 0, 1, 1}
	s.btnStart = eui.NewButton(eui.GetLocale().Get("btnStart"), rect, eui.GetTheme().Get("game active color"), eui.GetTheme().Get("game bg"), func(b *eui.Button) { eui.Push(game.NewSceneGame()) })
	s.Add(s.btnStart)
	s.btnScore = eui.NewButton(eui.GetLocale().Get("btnScore"), rect, eui.GetTheme().Get("error color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { eui.Push(score.NewSceneScore()) })
	s.Add(s.btnScore)
	str := "<"
	if eui.GetUi().IsMainScene() {
		str = "x"
	}
	s.btnQuit = eui.NewButton(str, rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { eui.Pop() })
	s.Add(s.btnQuit)
	s.lblName = eui.NewLabel(eui.GetLocale().Get("AppName"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblName)
	s.lblPeriodResult = eui.NewLabel(data.GetDb().TodayData.String(), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblPeriodResult)
	s.lblDt = eui.NewLabel(eui.GetLocale().Get("lblUpTm"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblDt)
	s.lblsResult = eui.NewList(nil, nil, rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), s.getRows())
	s.Add(s.lblsResult)
	s.lblHelper = eui.NewLabel(eui.GetLocale().Get("btnHelper"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblHelper)
	s.plotResult = plot.NewResultPlot(rect)
	s.plotResult.Visible = false
	s.Add(s.plotResult)
	s.toggleResults = false
	s.btnPlot = eui.NewButton(eui.GetLocale().Get("btnPlot"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { s.togglePlot() })
	s.Add(s.btnPlot)
	s.btnFullScreen = eui.NewButton("[ ]", rect, eui.GetTheme().Get("regular color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { eui.GetUi().ToggleFullscreen() })
	s.Add(s.btnFullScreen)
	s.btnOpt = eui.NewButton(eui.GetLocale().Get("btnOpt"), rect, eui.GetTheme().Get("warning color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { eui.Push(options.NewSceneOptions()) })
	s.Add(s.btnOpt)
	return s
}

func (s *SceneToday) Entered() {
	data.GetDb().ReadTodayGames(time.Now().Format("2006-01-02"))
	s.lblPeriodResult.SetText(data.GetDb().TodayData.String())
	a, b := data.GetDb().TodayData.ListShortStr()
	s.lblsResult.SetList(a, b)
	s.Resize()
	log.Println("Entered SceneToday")
	log.Println(data.GetDb().TodayData.LongStr())
}

func (s *SceneToday) Update(dt int) {
	s.lblDt.SetText(eui.GetUi().UpdateUpTime())
	for _, value := range s.Container {
		value.Update(dt)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		eui.Push(game.NewSceneGame())
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		s.togglePlot()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		eui.Push(score.NewSceneScore())
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyO) {
		eui.Push(options.NewSceneOptions())
	}
}

func (s *SceneToday) togglePlot() {
	s.toggleResults = !s.toggleResults
	if s.toggleResults {
		s.plotResult.Dirty = true
		s.plotResult.Visible = true
		s.lblsResult.Visible = false
	} else {
		s.plotResult.Visible = false
		s.lblsResult.Visible = true
		s.lblsResult.Dirty = true
	}
}

func (s *SceneToday) Draw(surface *ebiten.Image) {
	for _, value := range s.Container {
		value.Draw(surface)
	}
}

func (s *SceneToday) Resize() {
	w, h := eui.GetUi().GetScreenSize()
	s.rect = eui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(s.rect.H)*0.05), int(float64(s.rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, int(float64(s.rect.W)*0.20)
	s.lblName.Resize([]int{x, y, w, h})
	x = w + h
	s.btnScore.Resize([]int{x, y, w, h})
	x = s.rect.Right() - w
	s.lblDt.Resize([]int{x, y, w, h})
	x -= w
	s.btnOpt.Resize([]int{x, y, w, h})
	w = int(float64(s.rect.H) * 0.05)
	x = x - w
	s.btnFullScreen.Resize([]int{x, y, w, h})

	w = int(float64(s.rect.H) * 0.8)
	x = (s.rect.W - w) / 2
	y = int(float64(h) * 1.2)
	s.btnStart.Resize([]int{x, y, w, h})
	w = int(float64(s.rect.W) * 0.9)
	x, y = (s.rect.W-w)/2, int(float64(h)*2.4)
	s.lblPeriodResult.Resize([]int{x, y, w, h})
	y = int(float64(h) * 3.6)
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.75)
	x = (s.rect.W - w) / 2
	s.lblsResult.SetRows(s.getRows())
	s.lblsResult.Resize([]int{x, y, w, h})
	s.plotResult.Resize([]int{x, y, w, h})
	w = (s.rect.W - w) / 2
	h = w
	x = s.rect.W - w
	s.btnPlot.Resize([]int{x, y, w, h})
	w, h = s.rect.Right(), int(float64(s.rect.H)*0.05)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h))
	s.lblHelper.Resize([]int{x, y, w, h})
}

func (s *SceneToday) Close() {
	for _, v := range s.Container {
		v.Close()
	}
	data.GetDb().Close()
}

func (s *SceneToday) getRows() (rows int) {
	switch w := s.rect.Right(); {
	case w <= 360:
		rows = 2
	case w <= 640:
		rows = 3
	case w <= 800:
		rows = 4
	case w <= 1024:
		rows = 5
	default:
		rows = 6
	}
	return
}
