package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
)

type GameOpt struct {
	ui.ContainerDefault
	topBar                                    *TopBarOpt
	pref                                      *ui.Preferences
	optCenterCell, optResetOnWrong, optManual *ui.Checkbox
	optShowGrid, optShowCross                 *ui.Checkbox
	optRR                                     *ui.Combobox
	optGridSize, optDefLevel, optManualAdv    *ui.Combobox
	optAdv, optFall, optFallSessions          *ui.Combobox
	optTrials, optFactor, optExponent         *ui.Combobox
	optTmNextCell, optTmShowCell              *ui.Combobox
	optGameType                               *ui.Combobox
	optMaxSym, optMaxAriphmetic               *ui.Combobox
}

func NewGameOpt() *GameOpt {
	s := &GameOpt{}
	rect := []int{0, 0, 1, 1}
	s.pref = LoadPreferences()
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)

	// opt for game feedback resetOnWrong RR pause
	s.optCenterCell = ui.NewCheckbox(ui.GetLocale().Get("optcc"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("use center cell", s.optCenterCell.Checked())
		log.Printf("Use center cell: %v", s.pref.Get("use center cell").(bool))
	})
	s.Add(s.optCenterCell)

	s.optResetOnWrong = ui.NewCheckbox(ui.GetLocale().Get("optreset"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("reset on first wrong", s.optResetOnWrong.Checked())
		log.Printf("Reset on wrong: %v", s.pref.Get("reset on first wrong").(bool))
	})
	s.Add(s.optResetOnWrong)

	s.optManual = ui.NewCheckbox(ui.GetLocale().Get("optmanual"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("manual mode", s.optManual.Checked())
		log.Printf("Manual: %v", s.pref.Get("manual mode").(bool))
	})
	s.Add(s.optManual)

	s.optShowGrid = ui.NewCheckbox(ui.GetLocale().Get("optgrid"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("show grid", s.optShowGrid.Checked())
		log.Printf("Show Grid: %v", s.pref.Get("show grid").(bool))
	})
	s.Add(s.optShowGrid)

	s.optShowCross = ui.NewCheckbox(ui.GetLocale().Get("optcross"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("show crosshair", s.optShowCross.Checked())
		log.Printf("Show crosshair: %v", s.pref.Get("show crosshair").(bool))
	})
	s.Add(s.optShowCross)

	lvls := []interface{}{2, 3, 4, 5}
	idx := 1
	s.optGridSize = ui.NewCombobox(ui.GetLocale().Get("optgridsz"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), lvls, idx, func(c *ui.Combobox) {
		s.pref.Set("grid size", s.optGridSize.Value().(int))
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
		if i == s.pref.Get("random repition").(float64) {
			idx = j
		}
	}
	s.optRR = ui.NewCombobox(ui.GetLocale().Get("optrr"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), rrData, idx, func(c *ui.Combobox) {
		s.pref.Set("random repition", s.optRR.Value().(float64))
	})
	s.Add(s.optRR)

	values, _ := data.GetDb().ReadAllGamesScore(0, "", "")
	max := values.Max
	if max == 0 {
		max = 1
	}
	current := 0
	var arr []interface{}
	for i := 1; i <= max; i++ {
		arr = append(arr, i)
		if s.pref.Get("default level") == i {
			current = i - 1
		}
	}
	s.optDefLevel = ui.NewCombobox(ui.GetLocale().Get("optdeflev"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arr, current, func(c *ui.Combobox) {
		s.pref.Set("default level", s.optDefLevel.Value().(int))
	})
	s.Add(s.optDefLevel)

	arrAdvManual := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optManualAdv = ui.NewCombobox(ui.GetLocale().Get("optdeflevadv"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrAdvManual, idx, func(b *ui.Combobox) {
		s.pref.Set("manual advance", s.optManualAdv.Value().(int))
	})
	s.Add(s.optManualAdv)

	{
		var arrAdv []interface{}
		for i, j := 5, 0; i <= 100; i, j = i+5, j+1 {
			arrAdv = append(arrAdv, i)
			if s.pref.Get("threshold advance") == int(i) {
				idx = j
			}
		}
		s.optAdv = ui.NewCombobox(ui.GetLocale().Get("optadv"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrAdv, idx, func(b *ui.Combobox) {
			s.pref.Set("threshold advance", s.optAdv.Value().(int))
		})
		s.Add(s.optAdv)
	}
	{
		var arrFall []interface{}
		for i, j := 5, 0; i <= 100; i, j = i+5, j+1 {
			arrFall = append(arrFall, i)
			if s.pref.Get("threshold fallback").(int) == int(i) {
				idx = j
			}
		}
		s.optFall = ui.NewCombobox(ui.GetLocale().Get("optfall"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrFall, idx, func(b *ui.Combobox) {
			s.pref.Set("threshold fallback", s.optFall.Value().(int))
		})
		s.Add(s.optFall)
	}

	arrFallSessions := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 3
	s.optFallSessions = ui.NewCombobox(ui.GetLocale().Get("optgmadv"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrFallSessions, idx, func(b *ui.Combobox) {
		s.pref.Set("threshold fallback sessions", s.optFallSessions.Value().(int))
	})
	s.Add(s.optFallSessions)

	arrTrials := []interface{}{5, 10, 20, 30, 50}
	idx = 0
	s.optTrials = ui.NewCombobox(ui.GetLocale().Get("optmv"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrTrials, idx, func(b *ui.Combobox) {
		s.pref.Set("trials", s.optTrials.Value().(int))
	})
	s.Add(s.optTrials)

	arrFactor := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optFactor = ui.NewCombobox(ui.GetLocale().Get("optfc"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrFactor, idx, func(b *ui.Combobox) {
		s.pref.Set("trials factor", s.optFactor.Value().(int))
	})
	s.Add(s.optFactor)

	arrExp := []interface{}{1, 2, 3}
	idx = 1
	s.optExponent = ui.NewCombobox(ui.GetLocale().Get("optexp"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrExp, idx, func(b *ui.Combobox) {
		s.pref.Set("trials exponent", s.optExponent.Value().(int))
	})
	s.Add(s.optExponent)

	var arrTimeNextCell []interface{}
	for i, j = 1.0, 0; i <= 5; i, j = i+0.5, j+1 {
		arrTimeNextCell = append(arrTimeNextCell, i)
		if s.pref.Get("time to next cell").(float64) == i {
			idx = j
		}
	}
	s.optTmNextCell = ui.NewCombobox(ui.GetLocale().Get("opttmnc"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrTimeNextCell, idx, func(b *ui.Combobox) {
		tmnc := s.pref.Get("time to next cell").(float64)
		tmsc := s.pref.Get("time to show cell").(float64)
		if tmnc-0.5 > tmsc {
			s.pref.Set("time to next cell", s.optTmNextCell.Value().(float64))
		}
	})
	s.Add(s.optTmNextCell)

	arrShow := []interface{}{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5}
	idx = 0
	s.optTmShowCell = ui.NewCombobox(ui.GetLocale().Get("opttmsc"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrShow, idx, func(b *ui.Combobox) {
		tmnc := s.pref.Get("time to next cell").(float64)
		value := s.optTmShowCell.Value().(float64)
		if value < tmnc {
			s.pref.Set("time to show cell", value)
		}
	})
	s.Add(s.optTmShowCell)

	gamesType := []interface{}{game.Pos, game.Col, game.Sym, game.Ari}
	idx = 0
	s.optGameType = ui.NewCombobox(s.getGameType(), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), gamesType, idx, func(b *ui.Combobox) {
		s.pref.Set("game type", s.optGameType.Value().(string))
		s.optGameType.SetText(s.getGameType())
	})
	s.Add(s.optGameType)

	arrMaxSymbols := []interface{}{10, 20, 50, 100, 200, 500, 1000}
	s.optMaxSym = ui.NewCombobox(ui.GetLocale().Get("optmaxsym"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrMaxSymbols, 3, func(c *ui.Combobox) {
		s.pref.Set("symbols count", s.optMaxSym.Value().(int))
	})
	s.Add(s.optMaxSym)

	s.optMaxAriphmetic = ui.NewCombobox(ui.GetLocale().Get("optmaxari"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrMaxSymbols, 1, func(c *ui.Combobox) {
		s.pref.Set("ariphmetic max", s.optMaxAriphmetic.Value().(int))
	})
	s.Add(s.optMaxAriphmetic)
	return s
}

func (s *GameOpt) getGameType() string {
	result := ui.GetLocale().Get("optgmtp") + " "
	tp := ui.GetPreferences().Get("game type").(string)
	switch tp {
	case game.Pos:
		result += ui.GetLocale().Get("optpos")
	case game.Col:
		result += ui.GetLocale().Get("optcol")
	case game.Sym:
		result += ui.GetLocale().Get("optsym")
	case game.Ari:
		result += ui.GetLocale().Get("optari")
	}
	return result
}

func (s *GameOpt) Setup(sets *ui.Preferences) {
	s.optResetOnWrong.SetChecked(sets.Get("reset on first wrong").(bool))
	s.optRR.SetValue(sets.Get("random repition").(float64))
	s.optTmNextCell.SetValue(sets.Get("time to next cell").(float64))
	s.optTmShowCell.SetValue(sets.Get("time to show cell").(float64))
	s.optManual.SetChecked(sets.Get("manual mode").(bool))
	s.optDefLevel.SetValue(sets.Get("default level").(int))
	s.optManualAdv.SetValue(sets.Get("manual advance").(int))
	s.optTrials.SetValue(sets.Get("trials").(int))
	s.optFactor.SetValue(sets.Get("trials factor").(int))
	s.optExponent.SetValue(sets.Get("trials exponent").(int))
	s.optAdv.SetValue(sets.Get("threshold advance").(int))
	s.optFall.SetValue(sets.Get("threshold fallback").(int))
	s.optFallSessions.SetValue(sets.Get("threshold fallback sessions").(int))
	s.optGameType.SetValue(sets.Get("game type").(string))
	s.optShowCross.SetChecked(sets.Get("show crosshair").(bool))
	s.optGridSize.SetValue(sets.Get("grid size").(int))
	s.optShowGrid.SetChecked(sets.Get("show grid").(bool))
	s.optCenterCell.SetChecked(sets.Get("use center cell").(bool))
	s.optMaxSym.SetValue(sets.Get("symbols count").(int))
	s.optMaxAriphmetic.SetValue(sets.Get("ariphmetic max").(int))
}

func (s *GameOpt) Reset(b *ui.Button) {
	s.pref = ui.GetUi().ApplyPreferences(NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *GameOpt) Apply(b *ui.Button) {
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	ui.Pop()
}

func (r *GameOpt) Update(dt int) {
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *GameOpt) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *GameOpt) Entered() {
	s.Setup(LoadPreferences())
	s.Resize()
}

func (s *GameOpt) Resize() {
	s.topBar.Resize()
	w, h := ui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/20
	x, y := rect.CenterX()-w1/2, hTop
	s.optRR.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optResetOnWrong.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optDefLevel.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optManual.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optManualAdv.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optTmNextCell.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optTmShowCell.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optAdv.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optFall.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optFallSessions.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optTrials.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optFactor.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optExponent.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optGameType.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optShowGrid.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optShowCross.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optGridSize.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optCenterCell.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optMaxSym.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optMaxAriphmetic.Resize([]int{x, y, w1, h1 - 2})
}

func (r *GameOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
