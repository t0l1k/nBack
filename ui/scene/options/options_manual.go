package options

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui/app"
)

type ManualOpt struct {
	eui.ContainerDefault
	topBar                                 *TopBarOpt
	lblResult                              *eui.Label
	optDefLevel, optGameType, optManualAdv *eui.Combobox
	optResetOnWrong, optManual             *eui.Checkbox
	pref                                   *eui.Preferences
}

func NewManualOpt() *ManualOpt {
	s := &ManualOpt{}
	s.pref = app.LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)

	s.lblResult = eui.NewLabel("Manual mode", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblResult)

	gamesType := []interface{}{game.Pos, game.Col, game.Sym, game.Ari}
	idx := 0
	s.optGameType = eui.NewCombobox(s.getGameType(), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), gamesType, idx, func(b *eui.Combobox) {
		s.pref.Set("game type", s.optGameType.Value().(string))
		s.optGameType.SetText(s.getGameType())
		s.lblResult.SetText(fmt.Sprintf("Выбрать играть на ручнике уровень:%v, %v, ходов:%v", s.pref.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))

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
	s.optDefLevel = eui.NewCombobox(eui.GetLocale().Get("optdeflev"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arr, current, func(c *eui.Combobox) {
		s.pref.Set("default level", s.optDefLevel.Value().(int))
		s.lblResult.SetText(fmt.Sprintf("Выбрать играть на ручнике уровень:%v, %v, ходов:%v", s.pref.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))

	})
	s.Add(s.optDefLevel)

	arrAdvManual := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optManualAdv = eui.NewCombobox(eui.GetLocale().Get("optdeflevadv"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arrAdvManual, idx, func(b *eui.Combobox) {
		s.pref.Set("manual advance", s.optManualAdv.Value().(int))
	})
	s.Add(s.optManualAdv)

	s.optManual = eui.NewCheckbox(eui.GetLocale().Get("optmanual"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), func(c *eui.Checkbox) {
		s.pref.Set("manual mode", s.optManual.Checked())
		log.Printf("Manual: %v", s.pref.Get("manual mode").(bool))
	})
	s.Add(s.optManual)

	s.optResetOnWrong = eui.NewCheckbox(eui.GetLocale().Get("optreset"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), func(c *eui.Checkbox) {
		s.pref.Set("reset on first wrong", s.optResetOnWrong.Checked())
		log.Printf("Reset on wrong: %v", s.pref.Get("reset on first wrong").(bool))
	})
	s.Add(s.optResetOnWrong)

	return s
}

func (s *ManualOpt) getGameType() string {
	result := eui.GetLocale().Get("optgmtp") + " "
	tp := eui.GetPreferences().Get("game type").(string)
	switch tp {
	case game.Pos:
		result += eui.GetLocale().Get("optpos")
	case game.Col:
		result += eui.GetLocale().Get("optcol")
	case game.Sym:
		result += eui.GetLocale().Get("optsym")
	case game.Ari:
		result += eui.GetLocale().Get("optari")
	}
	return result
}

func (s *ManualOpt) Setup(sets *eui.Preferences) {
	s.optGameType.SetValue(sets.Get("game type").(string))
	s.optDefLevel.SetValue(sets.Get("default level").(int))
	sets.Set("manual advance", 0)
	s.optManualAdv.SetValue(sets.Get("manual advance").(int))
	sets.Set("manual mode", true)
	s.optManual.SetChecked(sets.Get("manual mode").(bool))
	s.optResetOnWrong.SetChecked(sets.Get("reset on first wrong").(bool))
	s.lblResult.SetText(fmt.Sprintf("Выбрать играть на ручнике уровень:%v, %v, ходов:%v", sets.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))

}

func (s *ManualOpt) Reset(b *eui.Button) {
	s.pref = eui.GetUi().ApplyPreferences(app.NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *ManualOpt) Apply(b *eui.Button) {
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	eui.Pop()
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
	s.Setup(app.LoadPreferences())
	s.Resize()
}

func (s *ManualOpt) Resize() {
	s.topBar.Resize()
	w, h := eui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := eui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*6
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.lblResult.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
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
