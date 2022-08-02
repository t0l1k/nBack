package game

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
	optShowGrid, optShowCross                                         *ui.Checkbox
	optRR, optPause                                                   *ui.Combobox
	optGridSize, optDefLevel, optManualAdv                            *ui.Combobox
	optAdv, optFall, optFallSessions                                  *ui.Combobox
	optTrials, optFactor, optExponent                                 *ui.Combobox
	optTmNextCell, optTmShowCell                                      *ui.Combobox
	newSets                                                           *ui.Preferences
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{
		rect: ui.NewRect([]int{0, 0, 1, 1}),
	}
	s.newSets = LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.btnQuit = ui.NewButton("<", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetApp().Pop() })
	s.Add(s.btnQuit)
	s.name = "N-Back Options"
	s.lblName = ui.NewLabel(s.name, rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblName)
	s.optTheme = NewOptTheme(rect)
	s.Add(s.optTheme)
	s.btnReset = ui.NewButton("Reset", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], s.Reset)
	s.Add(s.btnReset)
	s.btnApply = ui.NewButton("Apply", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], s.Apply)
	s.Add(s.btnApply)
	s.optFullScr = ui.NewCheckbox("Fullscreen on app start", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["fullscreen"] = s.optFullScr.Checked()
		log.Printf("fullscreen checked: %v", (*s.newSets)["fullscreen"].(bool))
	})
	s.Add(s.optFullScr)
	s.optCenterCell = ui.NewCheckbox("Use center cell", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["use center cell"] = s.optCenterCell.Checked()
		log.Printf("Use center cell: %v", (*s.newSets)["use center cell"].(bool))
	})
	s.Add(s.optCenterCell)
	s.optFeeback = ui.NewCheckbox("Feedback on move", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["feedback on user move"] = s.optFeeback.Checked()
		log.Printf("Feedback on mpve: %v", (*s.newSets)["feedback on user move"].(bool))
	})
	s.Add(s.optFeeback)

	s.optResetOnWrong = ui.NewCheckbox("Reset on wrong", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["reset on first wrong"] = s.optResetOnWrong.Checked()
		log.Printf("Reset on wrong: %v", (*s.newSets)["reset on first wrong"].(bool))
	})
	s.Add(s.optResetOnWrong)

	s.optManual = ui.NewCheckbox("Manual", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["manual mode"] = s.optManual.Checked()
		log.Printf("Manual: %v", (*s.newSets)["manual mode"].(bool))
	})
	s.Add(s.optManual)

	s.optShowGrid = ui.NewCheckbox("Show Grid", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["show grid"] = s.optShowGrid.Checked()
		log.Printf("Show Grid: %v", (*s.newSets)["show grid"].(bool))
	})
	s.Add(s.optShowGrid)

	s.optShowCross = ui.NewCheckbox("Show Crosshair", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["show crosshair"] = s.optShowCross.Checked()
		log.Printf("Show crosshair: %v", (*s.newSets)["show crosshair"].(bool))
	})
	s.Add(s.optShowCross)

	data := []interface{}{2, 3, 4, 5}
	idx := 1
	s.optGridSize = ui.NewCombobox("Grid size", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], data, idx, func(c *ui.Combobox) {
		(*s.newSets)["grid size"] = s.optGridSize.Value().(int)
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
		if i == (*s.newSets)["random repition"].(float64) {
			idx = j
		}
	}
	s.optRR = ui.NewCombobox("Random Repition", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], rrData, idx, func(c *ui.Combobox) { (*s.newSets)["random repition"] = s.optRR.Value().(float64) })
	s.Add(s.optRR)

	arrPauses := []interface{}{3, 5, 10, 15, 20, 30, 45, 60, 90, 180}
	s.optPause = ui.NewCombobox("Pause after game", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrPauses, 2, func(c *ui.Combobox) { (*s.newSets)["pause to rest"] = s.optPause.Value().(int) })
	s.Add(s.optPause)

	values, _ := getDb().ReadAllGamesScore()
	max := values.max
	if max == 0 {
		max = 1
	}
	current := 0
	var arr []interface{}
	for i := 1; i <= max; i++ {
		arr = append(arr, i)
		if (*s.newSets)["default level"] == i {
			current = i - 1
		}
	}
	s.optDefLevel = ui.NewCombobox("Default level", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arr, current, func(c *ui.Combobox) {
		(*s.newSets)["default level"] = s.optDefLevel.Value().(int)
	})
	s.Add(s.optDefLevel)

	arrAdvManual := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optManualAdv = ui.NewCombobox("Manual advance", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrAdvManual, idx, func(b *ui.Combobox) {
		(*s.newSets)["manual advance"] = s.optManualAdv.Value().(int)
	})
	s.Add(s.optManualAdv)

	{
		var arrAdv []interface{}
		for i, j := 5, 0; i <= 100; i, j = i+5, j+1 {
			arrAdv = append(arrAdv, i)
			if (*s.newSets)["threshold advance"] == int(i) {
				idx = j
			}
		}
		s.optAdv = ui.NewCombobox("Advance", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrAdv, idx, func(b *ui.Combobox) { (*s.newSets)["threshold advance"] = s.optAdv.Value().(int) })
		s.Add(s.optAdv)
	}
	{
		var arrFall []interface{}
		for i, j := 5, 0; i <= 100; i, j = i+5, j+1 {
			arrFall = append(arrFall, i)
			if (*s.newSets)["threshold fallback"].(int) == int(i) {
				idx = j
			}
		}
		s.optFall = ui.NewCombobox("Fallback", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrFall, idx, func(b *ui.Combobox) { (*s.newSets)["threshold fallback"] = s.optFall.Value().(int) })
		s.Add(s.optFall)
	}

	arrFallSessions := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 3
	s.optFallSessions = ui.NewCombobox("Fallback sessions", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrFallSessions, idx, func(b *ui.Combobox) { (*s.newSets)["threshold fallback sessions"] = s.optFallSessions.Value().(int) })
	s.Add(s.optFallSessions)

	arrTrials := []interface{}{5, 10, 20, 30, 50}
	idx = 0
	s.optTrials = ui.NewCombobox("Trials", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrTrials, idx, func(b *ui.Combobox) { (*s.newSets)["trials"] = s.optTrials.Value().(int) })
	s.Add(s.optTrials)

	arrFactor := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optFactor = ui.NewCombobox("Factor", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrFactor, idx, func(b *ui.Combobox) { (*s.newSets)["trials factor"] = s.optFactor.Value().(int) })
	s.Add(s.optFactor)

	arrExp := []interface{}{1, 2, 3}
	idx = 1
	s.optExponent = ui.NewCombobox("Exponent", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrExp, idx, func(b *ui.Combobox) { (*s.newSets)["trials exponent"] = s.optExponent.Value().(int) })
	s.Add(s.optExponent)

	var arrTimeNextCell []interface{}
	for i, j = 1.5, 0; i <= 5; i, j = i+0.5, j+1 {
		arrTimeNextCell = append(arrTimeNextCell, i)
		if (*s.newSets)["time to next cell"].(float64) == i {
			idx = j
		}
	}
	s.optTmNextCell = ui.NewCombobox("Time to next cell", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrTimeNextCell, idx, func(b *ui.Combobox) {
		(*s.newSets)["time to next cell"] = s.optTmNextCell.Value().(float64)
	})
	s.Add(s.optTmNextCell)

	arrShow := []interface{}{0.5, 1.0}
	idx = 0
	s.optTmShowCell = ui.NewCombobox("Time to show cell", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrShow, idx, func(b *ui.Combobox) { (*s.newSets)["time to show cell"] = s.optTmShowCell.Value().(float64) })
	s.Add(s.optTmShowCell)
	return s
}

func (s *SceneOptions) Setup(sets *ui.Preferences) {
	s.optFullScr.SetChecked((*sets)["fullscreen"].(bool))
	s.optCenterCell.SetChecked((*sets)["use center cell"].(bool))
	s.optFeeback.SetChecked((*sets)["feedback on user move"].(bool))
	s.optResetOnWrong.SetChecked((*sets)["reset on first wrong"].(bool))
	s.optManual.SetChecked((*sets)["manual mode"].(bool))
	s.optShowGrid.SetChecked((*sets)["show grid"].(bool))
	s.optShowCross.SetChecked((*sets)["show crosshair"].(bool))
	s.optRR.SetValue((*sets)["random repition"].(float64))
	s.optPause.SetValue((*sets)["pause to rest"].(int))
	s.optGridSize.SetValue((*sets)["grid size"].(int))
	s.optDefLevel.SetValue((*sets)["default level"].(int))
	s.optManualAdv.SetValue((*sets)["manual advance"].(int))
	s.optAdv.SetValue((*sets)["threshold advance"].(int))
	s.optFall.SetValue((*sets)["threshold fallback"].(int))
	s.optFallSessions.SetValue((*sets)["threshold fallback sessions"].(int))
	s.optTrials.SetValue((*sets)["trials"].(int))
	s.optFactor.SetValue((*sets)["trials factor"].(int))
	s.optExponent.SetValue((*sets)["trials exponent"].(int))
	s.optTmNextCell.SetValue((*sets)["time to next cell"].(float64))
	s.optTmShowCell.SetValue((*sets)["time to show cell"].(float64))
}

func (s *SceneOptions) Reset(b *ui.Button) {
	s.Setup(NewPref())
	log.Println("Reset All Options to Defaults")
}

func (s *SceneOptions) Apply(b *ui.Button) {
	sets := ApplyPreferences(s.newSets)
	getDb().InsertSettings(sets)
	log.Println("Apply Settings")
	ui.GetApp().Pop()
}

func (s *SceneOptions) Entered() {
	s.Setup(LoadPreferences())
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
	surface.Fill((*ui.GetTheme())["game bg"])
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneOptions) Resize() {
	w, h := ui.GetApp().GetScreenSize()
	s.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(s.rect.H)*0.05), int(float64(s.rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, int(float64(s.rect.W)*0.20)
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

	x, y = 0, int(float64(cellHeight)*9.6)
	h2 = float64(cellWidth) / 2.1
	rect = []int{x, y, int(h2), cellHeight}
	s.optShowGrid.Resize(rect)
	x = int(h2 * 1.1)
	rect = []int{x, y, int(h2), cellHeight}
	s.optShowCross.Resize(rect)
}

func (s *SceneOptions) Quit() {}
