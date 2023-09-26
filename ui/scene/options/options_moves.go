package options

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui/app"
)

type MovesOpt struct {
	eui.ContainerDefault
	topBar                            *TopBarOpt
	lblResult                         *eui.Label
	optDefLevel                       *eui.Combobox
	optTrials, optFactor, optExponent *eui.Combobox
	pref                              *eui.Preferences
	dirty                             bool
}

func NewMovesOpt() *MovesOpt {
	s := &MovesOpt{}
	s.pref = app.LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)

	s.lblResult = eui.NewLabel("Moves:", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
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
		s.dirty = true
	})
	s.Add(s.optDefLevel)

	arrTrials := []interface{}{5, 10, 20, 30, 50}
	idx := 0
	s.optTrials = eui.NewCombobox(eui.GetLocale().Get("optmv"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arrTrials, idx, func(b *eui.Combobox) {
		s.pref.Set("trials", s.optTrials.Value().(int))
		s.dirty = true
	})
	s.Add(s.optTrials)

	arrFactor := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optFactor = eui.NewCombobox(eui.GetLocale().Get("optfc"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arrFactor, idx, func(b *eui.Combobox) {
		s.pref.Set("trials factor", s.optFactor.Value().(int))
		s.dirty = true
	})
	s.Add(s.optFactor)

	arrExp := []interface{}{1, 2, 3}
	idx = 1
	s.optExponent = eui.NewCombobox(eui.GetLocale().Get("optexp"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arrExp, idx, func(b *eui.Combobox) {
		s.pref.Set("trials exponent", s.optExponent.Value().(int))
		s.dirty = true
	})
	s.Add(s.optExponent)
	return s
}

func (s *MovesOpt) Setup(sets *eui.Preferences) {
	level := sets.Get("default level").(int)
	s.optDefLevel.SetValue(level)
	s.optTrials.SetValue(sets.Get("trials").(int))
	s.optFactor.SetValue(sets.Get("trials factor").(int))
	s.optExponent.SetValue(sets.Get("trials exponent").(int))
	moves := game.TotalMoves(level)
	s.lblResult.SetText(fmt.Sprintf("Ходов на уровне %v будет %v", level, strconv.Itoa(moves)))
}

func (s *MovesOpt) Reset(b *eui.Button) {
	s.pref = eui.GetUi().ApplyPreferences(app.NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *MovesOpt) Apply(b *eui.Button) {
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	eui.Pop()
}

func (r *MovesOpt) Update(dt int) {
	if r.dirty {
		r.Setup(r.pref)
	}
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *MovesOpt) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *MovesOpt) Entered() {
	s.Setup(app.LoadPreferences())
	s.Resize()
}

func (s *MovesOpt) Resize() {
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
	s.optTrials.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optFactor.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optExponent.Resize([]int{x, y, w1, hTop - 2})

}

func (r *MovesOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
