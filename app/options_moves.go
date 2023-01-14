package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
)

type MovesOpt struct {
	ui.ContainerDefault
	topBar                            *TopBarOpt
	lblResult                         *ui.Label
	optDefLevel                       *ui.Combobox
	optTrials, optFactor, optExponent *ui.Combobox
	pref                              *ui.Preferences
	dirty                             bool
}

func NewMovesOpt() *MovesOpt {
	s := &MovesOpt{}
	s.pref = LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)

	s.lblResult = ui.NewLabel("Moves:", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
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
		s.dirty = true
	})
	s.Add(s.optDefLevel)

	arrTrials := []interface{}{5, 10, 20, 30, 50}
	idx := 0
	s.optTrials = ui.NewCombobox(ui.GetLocale().Get("optmv"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrTrials, idx, func(b *ui.Combobox) {
		s.pref.Set("trials", s.optTrials.Value().(int))
		s.dirty = true
	})
	s.Add(s.optTrials)

	arrFactor := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	idx = 0
	s.optFactor = ui.NewCombobox(ui.GetLocale().Get("optfc"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrFactor, idx, func(b *ui.Combobox) {
		s.pref.Set("trials factor", s.optFactor.Value().(int))
		s.dirty = true
	})
	s.Add(s.optFactor)

	arrExp := []interface{}{1, 2, 3}
	idx = 1
	s.optExponent = ui.NewCombobox(ui.GetLocale().Get("optexp"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrExp, idx, func(b *ui.Combobox) {
		s.pref.Set("trials exponent", s.optExponent.Value().(int))
		s.dirty = true
	})
	s.Add(s.optExponent)
	return s
}

func (s *MovesOpt) Setup(sets *ui.Preferences) {
	level := sets.Get("default level").(int)
	s.optDefLevel.SetValue(level)
	s.optTrials.SetValue(sets.Get("trials").(int))
	s.optFactor.SetValue(sets.Get("trials factor").(int))
	s.optExponent.SetValue(sets.Get("trials exponent").(int))
	moves := game.TotalMoves(level)
	s.lblResult.SetText(fmt.Sprintf("Ходов на уровне %v будет %v", level, strconv.Itoa(moves)))
}

func (s *MovesOpt) Reset(b *ui.Button) {
	s.pref = ui.GetUi().ApplyPreferences(NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *MovesOpt) Apply(b *ui.Button) {
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	ui.Pop()
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
	s.Setup(LoadPreferences())
	s.Resize()
}

func (s *MovesOpt) Resize() {
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
