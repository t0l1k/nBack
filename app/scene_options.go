package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui"
)

type SceneOptions struct {
	rect                                                              *ui.Rect
	container                                                         []ui.Drawable
	lblName                                                           *ui.Label
	optTheme                                                          *OptTheme
	btnQuit, btnReset, btnApply                                       *ui.Button
	optFullScr, optCenterCell, optFeeback, optResetOnWrong, optManual *ui.Checkbox
	optShowGrid, optShowCross                                         *ui.Checkbox
	optRR, optPause                                                   *ui.Combobox
	optGridSize, optDefLevel, optManualAdv                            *ui.Combobox
	optAdv, optFall, optFallSessions                                  *ui.Combobox
	optTrials, optFactor, optExponent                                 *ui.Combobox
	optTmNextCell, optTmShowCell                                      *ui.Combobox
	optGameType, optLang                                              *ui.Combobox
	newSets                                                           *ui.Preferences
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{
		rect: ui.NewRect([]int{0, 0, 1, 1}),
	}
	s.newSets = LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.btnQuit = ui.NewButton("<", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetUi().Pop() })
	s.Add(s.btnQuit)

	s.lblName = ui.NewLabel(ui.GetLocale().Get("btnOpt"), rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblName)

	s.optTheme = NewOptTheme(rect)
	s.Add(s.optTheme)

	s.btnReset = ui.NewButton(ui.GetLocale().Get("btnReset"), rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], s.Reset)
	s.Add(s.btnReset)

	s.btnApply = ui.NewButton(ui.GetLocale().Get("btnSave"), rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], s.Apply)
	s.Add(s.btnApply)

	// opt app fullscreen lang
	s.optFullScr = ui.NewCheckbox(ui.GetLocale().Get("optfs"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["fullscreen"] = s.optFullScr.Checked()
		log.Printf("fullscreen checked: %v", (*s.newSets)["fullscreen"].(bool))
	})
	s.Add(s.optFullScr)

	// opt for game feedback resetOnWrong RR pause
	s.optCenterCell = ui.NewCheckbox(ui.GetLocale().Get("optcc"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["use center cell"] = s.optCenterCell.Checked()
		log.Printf("Use center cell: %v", (*s.newSets)["use center cell"].(bool))
	})
	s.Add(s.optCenterCell)

	s.optFeeback = ui.NewCheckbox(ui.GetLocale().Get("optfeedback"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["feedback on user move"] = s.optFeeback.Checked()
		log.Printf("Feedback on mpve: %v", (*s.newSets)["feedback on user move"].(bool))
	})
	s.Add(s.optFeeback)

	s.optResetOnWrong = ui.NewCheckbox(ui.GetLocale().Get("optreset"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["reset on first wrong"] = s.optResetOnWrong.Checked()
		log.Printf("Reset on wrong: %v", (*s.newSets)["reset on first wrong"].(bool))
	})
	s.Add(s.optResetOnWrong)

	s.optManual = ui.NewCheckbox(ui.GetLocale().Get("optmanual"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["manual mode"] = s.optManual.Checked()
		log.Printf("Manual: %v", (*s.newSets)["manual mode"].(bool))
	})
	s.Add(s.optManual)

	s.optShowGrid = ui.NewCheckbox(ui.GetLocale().Get("optgrid"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["show grid"] = s.optShowGrid.Checked()
		log.Printf("Show Grid: %v", (*s.newSets)["show grid"].(bool))
	})
	s.Add(s.optShowGrid)

	s.optShowCross = ui.NewCheckbox(ui.GetLocale().Get("optcross"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], func(c *ui.Checkbox) {
		(*s.newSets)["show crosshair"] = s.optShowCross.Checked()
		log.Printf("Show crosshair: %v", (*s.newSets)["show crosshair"].(bool))
	})
	s.Add(s.optShowCross)

	lvls := []interface{}{2, 3, 4, 5}
	idx := 1
	s.optGridSize = ui.NewCombobox(ui.GetLocale().Get("optgridsz"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], lvls, idx, func(c *ui.Combobox) {
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
	s.optRR = ui.NewCombobox(ui.GetLocale().Get("optrr"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], rrData, idx, func(c *ui.Combobox) { (*s.newSets)["random repition"] = s.optRR.Value().(float64) })
	s.Add(s.optRR)

	arrPauses := []interface{}{3, 5, 10, 15, 20, 30, 45, 60, 90, 180}
	s.optPause = ui.NewCombobox(ui.GetLocale().Get("optpause"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrPauses, 2, func(c *ui.Combobox) { (*s.newSets)["pause to rest"] = s.optPause.Value().(int) })
	s.Add(s.optPause)

	values, _ := data.GetDb().ReadAllGamesScore()
	max := values.Max
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
	s.optDefLevel = ui.NewCombobox(ui.GetLocale().Get("optdeflev"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arr, current, func(c *ui.Combobox) {
		(*s.newSets)["default level"] = s.optDefLevel.Value().(int)
	})
	s.Add(s.optDefLevel)

	arrAdvManual := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optManualAdv = ui.NewCombobox(ui.GetLocale().Get("optdeflevadv"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrAdvManual, idx, func(b *ui.Combobox) {
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
		s.optAdv = ui.NewCombobox(ui.GetLocale().Get("optadv"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrAdv, idx, func(b *ui.Combobox) { (*s.newSets)["threshold advance"] = s.optAdv.Value().(int) })
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
		s.optFall = ui.NewCombobox(ui.GetLocale().Get("optfall"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrFall, idx, func(b *ui.Combobox) { (*s.newSets)["threshold fallback"] = s.optFall.Value().(int) })
		s.Add(s.optFall)
	}

	arrFallSessions := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 3
	s.optFallSessions = ui.NewCombobox(ui.GetLocale().Get("optgmadv"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrFallSessions, idx, func(b *ui.Combobox) { (*s.newSets)["threshold fallback sessions"] = s.optFallSessions.Value().(int) })
	s.Add(s.optFallSessions)

	arrTrials := []interface{}{5, 10, 20, 30, 50}
	idx = 0
	s.optTrials = ui.NewCombobox(ui.GetLocale().Get("optmv"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrTrials, idx, func(b *ui.Combobox) { (*s.newSets)["trials"] = s.optTrials.Value().(int) })
	s.Add(s.optTrials)

	arrFactor := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optFactor = ui.NewCombobox(ui.GetLocale().Get("optfc"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrFactor, idx, func(b *ui.Combobox) { (*s.newSets)["trials factor"] = s.optFactor.Value().(int) })
	s.Add(s.optFactor)

	arrExp := []interface{}{1, 2, 3}
	idx = 1
	s.optExponent = ui.NewCombobox(ui.GetLocale().Get("optexp"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrExp, idx, func(b *ui.Combobox) { (*s.newSets)["trials exponent"] = s.optExponent.Value().(int) })
	s.Add(s.optExponent)

	var arrTimeNextCell []interface{}
	for i, j = 1.5, 0; i <= 5; i, j = i+0.5, j+1 {
		arrTimeNextCell = append(arrTimeNextCell, i)
		if (*s.newSets)["time to next cell"].(float64) == i {
			idx = j
		}
	}
	s.optTmNextCell = ui.NewCombobox(ui.GetLocale().Get("opttmnc"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrTimeNextCell, idx, func(b *ui.Combobox) {
		(*s.newSets)["time to next cell"] = s.optTmNextCell.Value().(float64)
	})
	s.Add(s.optTmNextCell)

	arrShow := []interface{}{0.5, 1.0}
	idx = 0
	s.optTmShowCell = ui.NewCombobox(ui.GetLocale().Get("opttmsc"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], arrShow, idx, func(b *ui.Combobox) { (*s.newSets)["time to show cell"] = s.optTmShowCell.Value().(float64) })
	s.Add(s.optTmShowCell)

	gamesType := []interface{}{game.Pos, game.Col, game.Sym}
	idx = 0
	s.optGameType = ui.NewCombobox(s.getGameType(), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], gamesType, idx, func(b *ui.Combobox) {
		(*s.newSets)["game type"] = s.optGameType.Value().(string)
		s.optGameType.SetText(s.getGameType())
	})
	s.Add(s.optGameType)

	langs := []interface{}{"en", "ru"}
	idx = 0
	for i, v := range langs {
		prefLang := ui.GetPreferences().Get("lang")
		if prefLang == v {
			idx = i
			break
		}
	}
	s.optLang = ui.NewCombobox(ui.GetLocale().Get("optlang"), rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"], langs, idx, func(b *ui.Combobox) {
		(*s.newSets)["lang"] = s.optLang.Value().(string)
	})
	s.Add(s.optLang)
	return s
}

func (s *SceneOptions) getGameType() string {
	result := ui.GetLocale().Get("optgmtp") + " "
	tp := ui.GetPreferences().Get("game type").(string)
	switch tp {
	case game.Pos:
		result += ui.GetLocale().Get("optpos")
	case game.Col:
		result += ui.GetLocale().Get("optcol")
	case game.Sym:
		result += ui.GetLocale().Get("optsym")
	}
	return result
}
func (s *SceneOptions) Setup(sets *ui.Preferences) {
	s.optLang.SetValue((*sets)["lang"].(string))
	s.optFullScr.SetChecked((*sets)["fullscreen"].(bool))
	s.optPause.SetValue((*sets)["pause to rest"].(int))
	s.optFeeback.SetChecked((*sets)["feedback on user move"].(bool))
	s.optResetOnWrong.SetChecked((*sets)["reset on first wrong"].(bool))
	s.optRR.SetValue((*sets)["random repition"].(float64))
	s.optTmNextCell.SetValue((*sets)["time to next cell"].(float64))
	s.optTmShowCell.SetValue((*sets)["time to show cell"].(float64))
	s.optManual.SetChecked((*sets)["manual mode"].(bool))
	s.optDefLevel.SetValue((*sets)["default level"].(int))
	s.optManualAdv.SetValue((*sets)["manual advance"].(int))
	s.optTrials.SetValue((*sets)["trials"].(int))
	s.optFactor.SetValue((*sets)["trials factor"].(int))
	s.optExponent.SetValue((*sets)["trials exponent"].(int))
	s.optAdv.SetValue((*sets)["threshold advance"].(int))
	s.optFall.SetValue((*sets)["threshold fallback"].(int))
	s.optFallSessions.SetValue((*sets)["threshold fallback sessions"].(int))
	s.optGameType.SetValue((*sets)["game type"].(string))
	s.optShowCross.SetChecked((*sets)["show crosshair"].(bool))
	s.optGridSize.SetValue((*sets)["grid size"].(int))
	s.optShowGrid.SetChecked((*sets)["show grid"].(bool))
	s.optCenterCell.SetChecked((*sets)["use center cell"].(bool))
}

func (s *SceneOptions) Reset(b *ui.Button) {
	s.Setup(NewPref())
	log.Println("Reset All Options to Defaults")
}

func (s *SceneOptions) Apply(b *ui.Button) {
	sets := ApplyPreferences(s.newSets)
	data.GetDb().InsertSettings(sets)
	log.Println("Apply Settings")
	ui.GetUi().Pop()
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
	w, h := ui.GetUi().GetScreenSize()
	s.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(s.rect.H)*0.05), int(float64(s.rect.H)/20)
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
	x, y = 0, int(float64(cellHeight)*1.1)
	rect = []int{x, y, cellWidth, cellHeight}
	s.optFullScr.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, cellWidth, cellHeight}
	s.optLang.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, cellWidth, cellHeight}
	s.optPause.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, cellWidth, cellHeight}
	s.optFeeback.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, cellWidth, cellHeight}
	s.optResetOnWrong.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, cellWidth, cellHeight}
	s.optRR.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	h2 := float64(cellWidth) / 2.1
	rect = []int{x, y, int(h2), cellHeight}
	s.optTmNextCell.Resize(rect)
	x = int(h2 * 1.1)
	rect = []int{x, y, int(h2), cellHeight}
	s.optTmShowCell.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	h3 := float64(cellWidth) / 3.2
	rect = []int{x, y, int(h3), cellHeight}
	s.optDefLevel.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optManual.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optManualAdv.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, int(h3), cellHeight}
	s.optAdv.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optFall.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optFallSessions.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, int(h3), cellHeight}
	s.optTrials.Resize(rect)
	x = int(h3 * 1.1)
	rect = []int{x, y, int(h3), cellHeight}
	s.optFactor.Resize(rect)
	x = int(h3 * 2.2)
	rect = []int{x, y, int(h3), cellHeight}
	s.optExponent.Resize(rect)

	cellWidth, cellHeight = w, h
	x, y = 0, int(float64(cellHeight)*1.1)+y
	rect = []int{x, y, cellWidth, cellHeight}
	s.optGameType.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	h2 = float64(cellWidth) / 2.1
	rect = []int{x, y, int(h2), cellHeight}
	s.optShowCross.Resize(rect)
	x = int(h2 * 1.1)
	rect = []int{x, y, int(h2), cellHeight}
	s.optShowGrid.Resize(rect)

	x, y = 0, int(float64(cellHeight)*1.1)+y
	h2 = float64(cellWidth) / 2.1
	rect = []int{x, y, int(h2), cellHeight}
	s.optGridSize.Resize(rect)
	x = int(h2 * 1.1)
	rect = []int{x, y, int(h2), cellHeight}
	s.optCenterCell.Resize(rect)
}

func (s *SceneOptions) Quit() {
	for _, v := range s.container {
		v.Close()
	}
}