package app

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
)

type JaeggiOpt struct {
	ui.ContainerDefault
	topBar                   *TopBarOpt
	lblResult                *ui.Label
	optDefLevel, optGameType *ui.Combobox
	pref                     *ui.Preferences
}

func NewJaeggiOpt() *JaeggiOpt {
	s := &JaeggiOpt{}
	s.pref = LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)
	s.lblResult = ui.NewLabel("Jaeggi NBack rulez", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblResult)
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
		s.lblResult.SetText(fmt.Sprintf("Выбрать игру Jaeggi NBack mode уровень:%v, %v, ходов:%v", s.pref.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))
	})
	s.Add(s.optDefLevel)

	gamesType := []interface{}{game.Pos, game.Col, game.Sym, game.Ari}
	idx := 0
	s.optGameType = ui.NewCombobox(s.getGameType(), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), gamesType, idx, func(b *ui.Combobox) {
		s.pref.Set("game type", s.optGameType.Value().(string))
		s.optGameType.SetText(s.getGameType())
		s.lblResult.SetText(fmt.Sprintf("Выбрать игру Jaeggi NBack mode уровень:%v, %v, ходов:%v", s.pref.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))

	})
	s.Add(s.optGameType)

	return s
}

func (s *JaeggiOpt) getGameType() string {
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

func (s *JaeggiOpt) Setup(sets *ui.Preferences) {
	sets.Set("trials", 20)
	sets.Set("trials factor", 1)
	sets.Set("trials exponent", 1)
	s.optDefLevel.SetValue(sets.Get("default level").(int))
	s.optGameType.SetValue(sets.Get("game type").(string))
	s.lblResult.SetText(fmt.Sprintf("Выбрать игру Jaeggi NBack mode уровень:%v, %v, ходов:%v", sets.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))
}

func (s *JaeggiOpt) Reset(b *ui.Button) {
	s.pref = ui.GetUi().ApplyPreferences(NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *JaeggiOpt) Apply(b *ui.Button) {
	s.pref.Set("time to next cell", 1.5)
	s.pref.Set("time to show cell", 0.65)
	s.pref.Set("trials", 20) //20 classic = trials+factor*level**exponent
	s.pref.Set("trials factor", 1)
	s.pref.Set("trials exponent", 1)
	s.pref.Set("threshold advance", 90)
	s.pref.Set("threshold fallback", 75)
	s.pref.Set("threshold fallback sessions", 1)
	s.pref.Set("grid size", 3)
	s.pref.Set("show grid", false)
	s.pref.Set("reset on first wrong", false)
	s.pref.Set("manual mode", false)
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	ui.Pop()
}

func (r *JaeggiOpt) Update(dt int) {
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *JaeggiOpt) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *JaeggiOpt) Entered() {
	s.Setup(LoadPreferences())
	s.Resize()
}

func (s *JaeggiOpt) Resize() {
	s.topBar.Resize()
	w, h := ui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*6
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.lblResult.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optDefLevel.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optGameType.Resize([]int{x, y, w1, hTop - 2})
}

func (r *JaeggiOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
