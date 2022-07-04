package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

type SceneOptions struct {
	name                                                   string
	rect                                                   *ui.Rect
	container                                              []ui.Drawable
	lblName                                                *ui.Label
	optTheme                                               *OptTheme
	btnReset, btnApply                                     *ui.Button
	optFullScr, optCenterCell, optFeeback, optResetOnWrong *ui.Checkbox
	optGridSize, optRR, optPause                           *ui.Button
	optDefLevel, optAdvManual                              *ui.Button
	optManual                                              *ui.Checkbox
	optAdv, optFall, optFallSessions                       *ui.Button
	optTrials, optFactor, optExponent                      *ui.Button
	optTmNextCell, optTmShowCell                           *ui.Button
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{
		rect: getApp().rect,
	}
	rect := []int{0, 0, 1, 1}
	s.name = "N-Back Options"
	s.lblName = ui.NewLabel(s.name, rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblName)
	s.optTheme = NewOptTheme(rect)
	s.Add(s.optTheme)
	s.btnReset = ui.NewButton("Reset", rect, getApp().theme.correct, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Pressed Reset") })
	s.Add(s.btnReset)
	s.btnApply = ui.NewButton("Apply", rect, getApp().theme.correct, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Pressed Apply") })
	s.Add(s.btnApply)
	s.optFullScr = ui.NewCheckbox("Fullscreen on app start", rect, app.theme.bg, app.theme.fg, func(c *ui.Checkbox) {
		getApp().preferences.fullScreen = s.optFullScr.Checked()
		log.Printf("fullscreen checked: %v", getApp().preferences.fullScreen)
	})
	s.Add(s.optFullScr)
	s.optCenterCell = ui.NewCheckbox("Use center cell", rect, app.theme.bg, app.theme.fg, func(c *ui.Checkbox) {
		getApp().preferences.usecentercell = s.optCenterCell.Checked()
		log.Printf("Use center cell: %v", getApp().preferences.usecentercell)
	})
	s.Add(s.optCenterCell)
	s.optFeeback = ui.NewCheckbox("Feedback on move", rect, app.theme.bg, app.theme.fg, func(c *ui.Checkbox) {
		getApp().preferences.feedbackOnUserMove = s.optFeeback.Checked()
		log.Printf("Feedback on mpve: %v", getApp().preferences.feedbackOnUserMove)
	})
	s.Add(s.optFeeback)
	s.optGridSize = ui.NewButton(fmt.Sprintf("Grid size: %v", getApp().preferences.gridSize), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Pressed grid size") })
	s.Add(s.optGridSize)
	s.optResetOnWrong = ui.NewCheckbox("Reset on wrong", rect, app.theme.bg, app.theme.fg, func(c *ui.Checkbox) {
		getApp().preferences.resetOnFirstWrong = s.optResetOnWrong.Checked()
		log.Printf("Reset on wrong: %v", getApp().preferences.resetOnFirstWrong)
	})
	s.Add(s.optResetOnWrong)
	s.optRR = ui.NewButton(fmt.Sprintf("Random Repition: %v", getApp().preferences.rr), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Random Repition") })
	s.Add(s.optRR)
	s.optPause = ui.NewButton(fmt.Sprintf("Pause after game: %v", getApp().preferences.pauseRest), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Pause to rest") })
	s.Add(s.optPause)
	s.optDefLevel = ui.NewButton(fmt.Sprintf("Default level: %v", getApp().preferences.defaultLevel), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Default level") })
	s.Add(s.optDefLevel)
	s.optManual = ui.NewCheckbox("Manual", rect, app.theme.bg, app.theme.fg, func(c *ui.Checkbox) {
		getApp().preferences.manual = s.optManual.Checked()
		log.Printf("Manual: %v", getApp().preferences.manual)
	})
	s.Add(s.optManual)
	s.optAdvManual = ui.NewButton(fmt.Sprintf("Manual advance: %v", getApp().preferences.manualAdv), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Manual advance") })
	s.Add(s.optAdvManual)

	s.optAdv = ui.NewButton(fmt.Sprintf("Advance: %v", getApp().preferences.thresholdAdvance), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Advance pecent") })
	s.Add(s.optAdv)
	s.optFall = ui.NewButton(fmt.Sprintf("Fallback: %v", getApp().preferences.thresholdFallback), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Fallback pecent") })
	s.Add(s.optFall)
	s.optFallSessions = ui.NewButton(fmt.Sprintf("Fallback sessions: %v", getApp().preferences.thresholdFallbackSessions), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Fallback sessions") })
	s.Add(s.optFallSessions)

	s.optTrials = ui.NewButton(fmt.Sprintf("Trials: %v", getApp().preferences.trials), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Trials") })
	s.Add(s.optTrials)

	s.optFactor = ui.NewButton(fmt.Sprintf("Factor: %v", getApp().preferences.trialsFactor), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Factor") })
	s.Add(s.optFactor)

	s.optExponent = ui.NewButton(fmt.Sprintf("Exponent: %v", getApp().preferences.trialsExponent), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Expponent") })
	s.Add(s.optExponent)

	s.optTmNextCell = ui.NewButton(fmt.Sprintf("Time to next cell: %v", getApp().preferences.timeToNextCell), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Time to next cell") })
	s.Add(s.optTmNextCell)

	s.optTmShowCell = ui.NewButton(fmt.Sprintf("Time to show cell: %v", getApp().preferences.timeShowCell), rect, getApp().theme.bg, getApp().theme.fg, func(b *ui.Button) { fmt.Println("Time to show cell") })
	s.Add(s.optTmShowCell)
	return s
}

func (s *SceneOptions) Setup() {
	s.optFullScr.SetChecked(getApp().preferences.fullScreen)
	s.optFeeback.SetChecked(getApp().preferences.feedbackOnUserMove)
	s.optCenterCell.SetChecked(getApp().preferences.usecentercell)
	s.optResetOnWrong.SetChecked(getApp().preferences.resetOnFirstWrong)
	s.optManual.SetChecked(getApp().preferences.manual)
}

func (s *SceneOptions) Entered() {
	s.Setup()
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
	s.btnReset.Resize([]int{s.rect.W - w*2, y, w, h})
	s.btnApply.Resize([]int{s.rect.W - w, y, w, h})
	y = s.rect.H - (h * 3)
	w, h1 := s.rect.W, h*3
	rect := []int{x, y, w, h1}
	s.optTheme.Resize(rect)
	cellWidth, cellHeight := w, h
	x, y = 0, int(float64(cellHeight)*1.2)
	rect = []int{x, y, cellWidth, cellHeight}
	s.optFullScr.Resize(rect)

	x, y = 0, int(float64(cellHeight)*2.4)
	h3 := float64(cellWidth) / 3.2
	rect = []int{x, y, int(h3), cellHeight}
	s.optCenterCell.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optFeeback.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optGridSize.Resize(rect)

	x, y = 0, int(float64(cellHeight)*3.6)
	rect = []int{x, y, int(h3), cellHeight}
	s.optResetOnWrong.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optRR.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optPause.Resize(rect)

	x, y = 0, int(float64(cellHeight)*4.8)
	rect = []int{x, y, int(h3), cellHeight}
	s.optDefLevel.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optManual.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optAdvManual.Resize(rect)

	x, y = 0, int(float64(cellHeight)*6.0)
	rect = []int{x, y, int(h3), cellHeight}
	s.optAdv.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optFall.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optFallSessions.Resize(rect)

	x, y = 0, int(float64(cellHeight)*7.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optTrials.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optFactor.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optExponent.Resize(rect)

	x, y = 0, int(float64(cellHeight)*8.4)
	h2 := float64(cellWidth) / 2.1
	rect = []int{x, y, int(h2), cellHeight}
	s.optTmNextCell.Resize(rect)
	x = int(h2 * 1.1)
	rect = []int{x, y, int(h2), cellHeight}
	s.optTmShowCell.Resize(rect)
}

func (s *SceneOptions) Quit() {}
