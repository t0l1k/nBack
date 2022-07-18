package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

type SceneOptions struct {
	name                                                              string
	rect                                                              *ui.Rect
	container                                                         []ui.Drawable
	lblName                                                           *ui.Label
	btnQuit                                                           *ui.Button
	optTheme                                                          *OptTheme
	btnReset, btnApply                                                *ui.Button
	optFullScr, optCenterCell, optFeeback, optResetOnWrong, optManual *ui.Checkbox
	optRR, optPause                                                   *ui.Combobox
	optGridSize, optDefLevel, optManualAdv                            *ui.Combobox
	optAdv, optFall, optFallSessions                                  *ui.Combobox
	optTrials, optFactor, optExponent                                 *ui.Combobox
	optTmNextCell, optTmShowCell                                      *ui.Combobox
	newSets                                                           *Setting
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{
		rect:    getApp().rect,
		newSets: getApp().preferences.Load(),
	}
	rect := []int{0, 0, 1, 1}
	s.btnQuit = ui.NewButton("<", rect, getApp().theme.CorrectColor, getApp().theme.Fg, func(b *ui.Button) { getApp().Pop() })
	s.Add(s.btnQuit)
	s.name = "N-Back Options"
	s.lblName = ui.NewLabel(s.name, rect, getApp().theme.CorrectColor, getApp().theme.Fg)
	s.Add(s.lblName)
	s.optTheme = NewOptTheme(rect)
	s.Add(s.optTheme)
	s.btnReset = ui.NewButton("Reset", rect, getApp().theme.CorrectColor, getApp().theme.Fg, s.Reset)
	s.Add(s.btnReset)
	s.btnApply = ui.NewButton("Apply", rect, getApp().theme.CorrectColor, getApp().theme.Fg, s.Apply)
	s.Add(s.btnApply)
	s.optFullScr = ui.NewCheckbox("Fullscreen on app start", rect, app.theme.Bg, app.theme.Fg, func(c *ui.Checkbox) {
		s.newSets.FullScreen = s.optFullScr.Checked()
		log.Printf("fullscreen checked: %v", s.newSets.FullScreen)
	})
	s.Add(s.optFullScr)
	s.optCenterCell = ui.NewCheckbox("Use center cell", rect, app.theme.Bg, app.theme.Fg, func(c *ui.Checkbox) {
		s.newSets.Usecentercell = s.optCenterCell.Checked()
		log.Printf("Use center cell: %v", s.newSets.Usecentercell)
	})
	s.Add(s.optCenterCell)
	s.optFeeback = ui.NewCheckbox("Feedback on move", rect, app.theme.Bg, app.theme.Fg, func(c *ui.Checkbox) {
		s.newSets.FeedbackOnUserMove = s.optFeeback.Checked()
		log.Printf("Feedback on mpve: %v", s.newSets.FeedbackOnUserMove)
	})
	s.Add(s.optFeeback)

	s.optResetOnWrong = ui.NewCheckbox("Reset on wrong", rect, app.theme.Bg, app.theme.Fg, func(c *ui.Checkbox) {
		s.newSets.ResetOnFirstWrong = s.optResetOnWrong.Checked()
		log.Printf("Reset on wrong: %v", s.newSets.ResetOnFirstWrong)
	})
	s.Add(s.optResetOnWrong)

	s.optManual = ui.NewCheckbox("Manual", rect, app.theme.Bg, app.theme.Fg, func(c *ui.Checkbox) {
		s.newSets.Manual = s.optManual.Checked()
		log.Printf("Manual: %v", s.newSets.Manual)
	})
	s.Add(s.optManual)

	data := []interface{}{2, 3, 4, 5}
	idx := 1
	s.optGridSize = ui.NewCombobox("Grid size", rect, getApp().theme.Bg, getApp().theme.Fg, data, idx, func(c *ui.Combobox) {
		s.newSets.GridSize = s.optGridSize.Value().(int)
		log.Println("Grid Size changed")
	})
	s.Add(s.optGridSize)

	var (
		rrData []interface{}
		i      float64
		j      int
	)
	for i, j = 5, 0; i < 50; i, j = i+0.5, j+1 {
		rrData = append(rrData, i)
		if i == s.newSets.RR {
			idx = j
		}
	}
	s.optRR = ui.NewCombobox("Random Repition", rect, getApp().theme.Bg, getApp().theme.Fg, rrData, idx, func(c *ui.Combobox) { s.newSets.RR = s.optRR.Value().(float64) })
	s.Add(s.optRR)

	arrPauses := []interface{}{3, 5, 10, 15, 20, 30, 45, 60, 90, 180}
	s.optPause = ui.NewCombobox("Pause after game", rect, getApp().theme.Bg, getApp().theme.Fg, arrPauses, 2, func(c *ui.Combobox) { s.newSets.PauseRest = s.optPause.Value().(int) })
	s.Add(s.optPause)

	values, _ := getApp().db.ReadAllGamesScore()
	max := values.max
	if max == 0 {
		max = 1
	}
	current := 0
	var arr []interface{}
	for i := 1; i <= max; i++ {
		arr = append(arr, i)
		if s.newSets.DefaultLevel == i {
			current = i - 1
		}
	}
	s.optDefLevel = ui.NewCombobox("Default level", rect, getApp().theme.Bg, getApp().theme.Fg, arr, current, func(c *ui.Combobox) {
		s.newSets.DefaultLevel = s.optDefLevel.Value().(int)
	})
	s.Add(s.optDefLevel)

	arrAdvManual := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optManualAdv = ui.NewCombobox("Manual advance", rect, getApp().theme.Bg, getApp().theme.Fg, arrAdvManual, idx, func(b *ui.Combobox) {
		s.newSets.ManualAdv = s.optManualAdv.Value().(int)
	})
	s.Add(s.optManualAdv)

	{
		var arrAdv []interface{}
		for i, j := 5, 0; i <= 100; i, j = i+5, j+1 {
			arrAdv = append(arrAdv, i)
			if s.newSets.ThresholdAdvance == int(i) {
				idx = j
			}
		}
		s.optAdv = ui.NewCombobox("Advance", rect, getApp().theme.Bg, getApp().theme.Fg, arrAdv, idx, func(b *ui.Combobox) { s.newSets.ThresholdAdvance = s.optAdv.Value().(int) })
		s.Add(s.optAdv)
	}
	{
		var arrFall []interface{}
		for i, j := 5, 0; i <= 100; i, j = i+5, j+1 {
			arrFall = append(arrFall, i)
			if s.newSets.ThresholdFallback == int(i) {
				idx = j
			}
		}
		s.optFall = ui.NewCombobox("Fallback", rect, getApp().theme.Bg, getApp().theme.Fg, arrFall, idx, func(b *ui.Combobox) { s.newSets.ThresholdFallback = s.optFall.Value().(int) })
		s.Add(s.optFall)
	}

	arrFallSessions := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 3
	s.optFallSessions = ui.NewCombobox("Fallback sessions", rect, getApp().theme.Bg, getApp().theme.Fg, arrFallSessions, idx, func(b *ui.Combobox) { s.newSets.ThresholdFallbackSessions = s.optFallSessions.Value().(int) })
	s.Add(s.optFallSessions)

	arrTrials := []interface{}{5, 10, 20, 30, 50}
	idx = 0
	s.optTrials = ui.NewCombobox("Trials", rect, getApp().theme.Bg, getApp().theme.Fg, arrTrials, idx, func(b *ui.Combobox) { s.newSets.Trials = s.optTrials.Value().(int) })
	s.Add(s.optTrials)

	arrFactor := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optFactor = ui.NewCombobox("Factor", rect, getApp().theme.Bg, getApp().theme.Fg, arrFactor, idx, func(b *ui.Combobox) { s.newSets.TrialsFactor = s.optFactor.Value().(int) })
	s.Add(s.optFactor)

	arrExp := []interface{}{1, 2, 3}
	idx = 1
	s.optExponent = ui.NewCombobox("Exponent", rect, getApp().theme.Bg, getApp().theme.Fg, arrExp, idx, func(b *ui.Combobox) { s.newSets.TrialsExponent = s.optExponent.Value().(int) })
	s.Add(s.optExponent)

	var arrTimeNextCell []interface{}
	for i, j = 1.5, 0; i <= 5; i, j = i+0.5, j+1 {
		arrTimeNextCell = append(arrTimeNextCell, i)
		if s.newSets.TimeToNextCell == i {
			idx = j
		}
	}
	s.optTmNextCell = ui.NewCombobox("Time to next cell", rect, getApp().theme.Bg, getApp().theme.Fg, arrTimeNextCell, idx, func(b *ui.Combobox) {
		s.newSets.TimeToNextCell = s.optTmNextCell.Value().(float64)
	})
	s.Add(s.optTmNextCell)

	arrShow := []interface{}{0.5, 1.0}
	idx = 0
	s.optTmShowCell = ui.NewCombobox("Time to show cell", rect, getApp().theme.Bg, getApp().theme.Fg, arrShow, idx, func(b *ui.Combobox) { s.newSets.TimeShowCell = s.optTmShowCell.Value().(float64) })
	s.Add(s.optTmShowCell)
	return s
}

func (s *SceneOptions) Setup() {
	sets := getApp().preferences.Load()
	s.optFullScr.SetChecked(sets.FullScreen)
	s.optCenterCell.SetChecked(sets.Usecentercell)
	s.optFeeback.SetChecked(sets.FeedbackOnUserMove)
	s.optResetOnWrong.SetChecked(sets.ResetOnFirstWrong)
	s.optManual.SetChecked(sets.Manual)
	s.optRR.SetValue(sets.RR)
	s.optPause.SetValue(sets.PauseRest)
	s.optGridSize.SetValue(sets.GridSize)
	s.optDefLevel.SetValue(sets.DefaultLevel)
	s.optManualAdv.SetValue(sets.ManualAdv)
	s.optAdv.SetValue(sets.ThresholdAdvance)
	s.optFall.SetValue(sets.ThresholdFallback)
	s.optFallSessions.SetValue(sets.ThresholdFallbackSessions)
	s.optTrials.SetValue(sets.Trials)
	s.optFactor.SetValue(sets.TrialsFactor)
	s.optExponent.SetValue(sets.TrialsExponent)
	s.optTmNextCell.SetValue(sets.TimeToNextCell)
	s.optTmShowCell.SetValue(sets.TimeShowCell)
}

func (s *SceneOptions) Reset(b *ui.Button) {
	getApp().preferences.Reset()
	s.Setup()
	log.Println("Reset All Options to Defaults")
}

func (s *SceneOptions) Apply(b *ui.Button) {
	sets := getApp().preferences.Apply(s.newSets)
	getApp().db.InsertSettings(sets)
	log.Println("Apply Settings")
	getApp().Pop()
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
	surface.Fill(getApp().theme.GameBg)
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneOptions) Resize() {
	s.rect = getApp().rect
	x, y, w, h := 0, 0, int(float64(getApp().rect.H)*0.05), int(float64(getApp().rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, int(float64(getApp().rect.W)*0.20)
	s.lblName.Resize([]int{x, y, w, h})
	s.btnReset.Resize([]int{s.rect.W - w*2, y, w, h})
	s.btnApply.Resize([]int{s.rect.W - w, y, w, h})
	y = s.rect.H - (h * 3)
	w, h1 := s.rect.W, h*3
	rect := []int{0, y, w, h1}
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
	s.optManualAdv.Resize(rect)

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
