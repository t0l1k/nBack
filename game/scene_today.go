package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

type SceneToday struct {
	name                                                        string
	lblName, lblPeriodResult, lblDt, lblHelper                  *ui.Label
	btnScore, btnStart, btnQuit, btnPlot, btnFullScreen, btnOpt *ui.Button
	lblsResult                                                  *ui.List
	plotResult                                                  *ResultPlot
	toggleResults                                               bool
	rect                                                        *ui.Rect
	container                                                   []ui.Drawable
}

func NewSceneToday() *SceneToday {
	s := &SceneToday{
		rect: ui.NewRect([]int{0, 0, 1, 1}),
	}
	rect := []int{0, 0, 1, 1}
	s.btnStart = ui.NewButton("Play", rect, (*ui.GetTheme())["game active color"], (*ui.GetTheme())["game bg"], func(b *ui.Button) { ui.GetApp().Push(NewSceneGame()) })
	s.Add(s.btnStart)
	s.btnScore = ui.NewButton("Score", rect, (*ui.GetTheme())["error color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetApp().Push(NewSceneScore()) })
	s.Add(s.btnScore)
	s.btnQuit = ui.NewButton("<", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetApp().Pop() })
	s.Add(s.btnQuit)
	s.name = "N-Back"
	s.lblName = ui.NewLabel(s.name, rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblName)
	s.lblPeriodResult = ui.NewLabel(getDb().todayData.String(), rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblPeriodResult)
	s.lblDt = ui.NewLabel("up: 00:00 ", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblDt)
	s.lblsResult = ui.NewList(nil, nil, rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], s.getRows())
	s.Add(s.lblsResult)
	s.lblHelper = ui.NewLabel("Press <SPACE> to start the game,<P> plot, <S> score,<F11> toggle fullscreen, <O> Options, <Esc> quit", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblHelper)
	s.plotResult = NewResultPlot(rect)
	s.plotResult.Visibe = false
	s.Add(s.plotResult)
	s.toggleResults = false
	s.btnPlot = ui.NewButton("{P}", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { s.togglePlot() })
	s.Add(s.btnPlot)
	s.btnFullScreen = ui.NewButton("[ ]", rect, (*ui.GetTheme())["regular color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetApp().ToggleFullscreen() })
	s.Add(s.btnFullScreen)
	s.btnOpt = ui.NewButton("Options", rect, (*ui.GetTheme())["warning color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetApp().Push(NewSceneOptions()) })
	s.Add(s.btnOpt)
	return s
}

func (s *SceneToday) Entered() {
	getDb().ReadTodayGames()
	s.lblPeriodResult.SetText(getDb().todayData.String())
	a, b := getDb().todayData.ListShortStr()
	s.lblsResult.SetList(a, b)
	s.Resize()
	log.Println("Entered SceneToday")
}
func (s *SceneToday) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}
func (s *SceneToday) Update(dt int) {
	s.lblDt.SetText(ui.GetApp().UpdateUpTime())
	for _, value := range s.container {
		value.Update(dt)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		ui.GetApp().Push(NewSceneGame())
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		s.togglePlot()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		ui.GetApp().Push(NewSceneScore())
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyO) {
		ui.GetApp().Push(NewSceneOptions())
	}
}

func (s *SceneToday) togglePlot() {
	s.toggleResults = !s.toggleResults
	if s.toggleResults {
		s.plotResult.Dirty = true
		s.plotResult.Visibe = true
		s.lblsResult.Visible = false
	} else {
		s.plotResult.Visibe = false
		s.lblsResult.Visible = true
		s.lblsResult.Dirty = true
	}
}

func (s *SceneToday) Draw(surface *ebiten.Image) {
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneToday) Resize() {
	w, h := ui.GetApp().GetScreenSize()
	s.rect = ui.NewRect([]int{0, 0, w, h})
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

func (s *SceneToday) Quit() {}

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
