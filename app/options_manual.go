package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
)

type ManualOpt struct {
	ui.ContainerDefault
	topBar                                 *TopBarOpt
	optDefLevel, optGameType, optManualAdv *ui.Combobox
	optResetOnWrong, optManual             *ui.Checkbox
	pref                                   *ui.Preferences
}

func NewManualOpt() *ManualOpt {
	s := &ManualOpt{}
	s.pref = LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)

	gamesType := []interface{}{game.Pos, game.Col, game.Sym, game.Ari}
	idx := 0
	s.optGameType = ui.NewCombobox(s.getGameType(), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), gamesType, idx, func(b *ui.Combobox) {
		s.pref.Set("game type", s.optGameType.Value().(string))
		s.optGameType.SetText(s.getGameType())
	})
	s.Add(s.optGameType)

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

	s.optManual = ui.NewCheckbox(ui.GetLocale().Get("optmanual"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("manual mode", s.optManual.Checked())
		log.Printf("Manual: %v", s.pref.Get("manual mode").(bool))
	})
	s.Add(s.optManual)

	s.optResetOnWrong = ui.NewCheckbox(ui.GetLocale().Get("optreset"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("reset on first wrong", s.optResetOnWrong.Checked())
		log.Printf("Reset on wrong: %v", s.pref.Get("reset on first wrong").(bool))
	})
	s.Add(s.optResetOnWrong)

	return s
}

func (s *ManualOpt) getGameType() string {
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

func (s *ManualOpt) Setup(sets *ui.Preferences) {
	s.optGameType.SetValue(sets.Get("game type").(string))
	s.optDefLevel.SetValue(sets.Get("default level").(int))
	s.optManualAdv.SetValue(sets.Get("manual advance").(int))
	s.optManual.SetChecked(sets.Get("manual mode").(bool))
	s.optResetOnWrong.SetChecked(sets.Get("reset on first wrong").(bool))
}

func (s *ManualOpt) Reset(b *ui.Button) {
	s.pref = ui.GetUi().ApplyPreferences(NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *ManualOpt) Apply(b *ui.Button) {
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	ui.Pop()
}

func (r *ManualOpt) Update(dt int) {
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *ManualOpt) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *ManualOpt) Entered() {
	s.Setup(LoadPreferences())
	s.Resize()
}

func (s *ManualOpt) Resize() {
	s.topBar.Resize()
	w, h := ui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*6
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.optGameType.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optDefLevel.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optManual.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optManualAdv.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optResetOnWrong.Resize([]int{x, y, w1, hTop - 2})

}

func (r *ManualOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
