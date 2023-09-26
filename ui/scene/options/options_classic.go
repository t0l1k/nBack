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

type ClassicOpt struct {
	eui.ContainerDefault
	topBar                   *TopBarOpt
	lblResult                *eui.Label
	optDefLevel, optGameType *eui.Combobox
	pref                     *eui.Preferences
}

func NewClassicOpt() *ClassicOpt {
	s := &ClassicOpt{}
	s.pref = app.LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)
	s.lblResult = eui.NewLabel("Classic NBack rulez", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
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
	s.optDefLevel = eui.NewCombobox(eui.GetLocale().Get("optdeflev"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arr, current, func(c *eui.Combobox) {
		s.pref.Set("default level", s.optDefLevel.Value().(int))
		s.lblResult.SetText(fmt.Sprintf("Выбрать игру классический NBack уровень:%v, %v, ходов:%v", s.pref.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))
	})
	s.Add(s.optDefLevel)

	gamesType := []interface{}{game.Pos, game.Col, game.Sym, game.Ari}
	idx := 0
	s.optGameType = eui.NewCombobox(s.getGameType(), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), gamesType, idx, func(b *eui.Combobox) {
		s.pref.Set("game type", s.optGameType.Value().(string))
		s.optGameType.SetText(s.getGameType())
		s.lblResult.SetText(fmt.Sprintf("Выбрать игру классический NBack уровень:%v, %v, ходов:%v", s.pref.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))

	})
	s.Add(s.optGameType)

	return s
}

func (s *ClassicOpt) getGameType() string {
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

func (s *ClassicOpt) Setup(sets *eui.Preferences) {
	sets.Set("trials", 20)
	sets.Set("trials factor", 1)
	sets.Set("trials exponent", 2)
	s.optDefLevel.SetValue(sets.Get("default level").(int))
	s.optGameType.SetValue(sets.Get("game type").(string))
	s.lblResult.SetText(fmt.Sprintf("Выбрать игру классический NBack уровень:%v, %v, ходов:%v", sets.Get("default level").(int), s.getGameType(), game.TotalMoves(s.pref.Get("default level").(int))))
}

func (s *ClassicOpt) Reset(b *eui.Button) {
	s.pref = eui.GetUi().ApplyPreferences(app.NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *ClassicOpt) Apply(b *eui.Button) {
	s.pref.Set("time to next cell", 3.0)
	s.pref.Set("time to show cell", 0.75)
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
	eui.Pop()
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
	s.Setup(app.LoadPreferences())
	s.Resize()
}

func (s *ClassicOpt) Resize() {
	s.topBar.Resize()
	w, h := eui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := eui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*6
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.lblResult.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optDefLevel.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optGameType.Resize([]int{x, y, w1, hTop - 2})
}

func (r *ClassicOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
