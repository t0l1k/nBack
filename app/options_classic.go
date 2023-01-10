package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
)

type ClassicOpt struct {
	ui.ContainerDefault
	topBar                   *TopBarOpt
	optDefLevel, optGameType *ui.Combobox
	pref                     *ui.Preferences
}

func NewClassicOpt() *ClassicOpt {
	s := &ClassicOpt{}
	s.pref = LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)

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

	gamesType := []interface{}{game.Pos, game.Col, game.Sym, game.Ari}
	idx := 0
	s.optGameType = ui.NewCombobox(s.getGameType(), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), gamesType, idx, func(b *ui.Combobox) {
		s.pref.Set("game type", s.optGameType.Value().(string))
		s.optGameType.SetText(s.getGameType())
	})
	s.Add(s.optGameType)

	return s
}

func (s *ClassicOpt) getGameType() string {
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

func (s *ClassicOpt) Setup(sets *ui.Preferences) {
	s.optGameType.SetValue(sets.Get("game type").(string))
	s.optDefLevel.SetValue(sets.Get("default level").(int))
}

func (s *ClassicOpt) Reset(b *ui.Button) {
	s.pref = ui.GetUi().ApplyPreferences(NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *ClassicOpt) Apply(b *ui.Button) {
	s.pref.Set("time to next cell", 3.0)
	s.pref.Set("time to show cell", 2.5)
	s.pref.Set("trials", 20) //20 classic = trials+factor*level**exponent
	s.pref.Set("trials factor", 1)
	s.pref.Set("trials exponent", 2)
	s.pref.Set("threshold advance", 80)
	s.pref.Set("threshold fallback", 50)
	s.pref.Set("threshold fallback sessions", 3)
	s.pref.Set("reset on first wrong", false)
	s.pref.Set("manual mode", false)
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	ui.Pop()
}

func (r *ClassicOpt) Update(dt int) {
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *ClassicOpt) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *ClassicOpt) Entered() {
	s.Setup(LoadPreferences())
	s.Resize()
}

func (s *ClassicOpt) Resize() {
	s.topBar.Resize()
	w, h := ui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*2
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.optDefLevel.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optGameType.Resize([]int{x, y, w1, hTop - 2})
}

func (r *ClassicOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
