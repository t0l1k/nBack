package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
)

type AppOpt struct {
	ui.ContainerDefault
	topBar                 *TopBarOpt
	optFullScr, optFeeback *ui.Checkbox
	optLang, optPause      *ui.Combobox
	optTheme               *OptTheme
	pref                   *ui.Preferences
}

func NewAppOpt() *AppOpt {
	s := &AppOpt{}
	s.pref = LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)
	s.optTheme = NewOptTheme(rect)
	s.Add(s.optTheme)
	s.optFullScr = ui.NewCheckbox(ui.GetLocale().Get("optfs"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("fullscreen", s.optFullScr.Checked())
		log.Printf("fullscreen checked: %v", s.pref.Get("fullscreen").(bool))
	})
	s.Add(s.optFullScr)

	langs := []interface{}{"en", "ru"}
	idx := 0
	for i, v := range langs {
		prefLang := ui.GetPreferences().Get("lang")
		if prefLang == v {
			idx = i
			break
		}
	}
	s.optLang = ui.NewCombobox(ui.GetLocale().Get("optlang"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), langs, idx, func(b *ui.Combobox) {
		s.pref.Set("lang", s.optLang.Value().(string))
	})
	s.Add(s.optLang)

	arrPauses := []interface{}{3, 5, 10, 15, 20, 30, 45, 60, 90, 180}
	s.optPause = ui.NewCombobox(ui.GetLocale().Get("optpause"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), arrPauses, 2, func(c *ui.Combobox) {
		s.pref.Set("pause to rest", s.optPause.Value().(int))
	})
	s.Add(s.optPause)

	s.optFeeback = ui.NewCheckbox(ui.GetLocale().Get("optfeedback"), rect, ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), func(c *ui.Checkbox) {
		s.pref.Set("feedback on user move", s.optFeeback.Checked())
		log.Printf("Feedback on mpve: %v", s.pref.Get("feedback on user move").(bool))
	})
	s.Add(s.optFeeback)

	return s
}

func (s *AppOpt) Setup(sets *ui.Preferences) {
	s.optLang.SetValue(sets.Get("lang").(string))
	s.optFullScr.SetChecked(sets.Get("fullscreen").(bool))
	s.optPause.SetValue(sets.Get("pause to rest").(int))
	s.optFeeback.SetChecked(sets.Get("feedback on user move").(bool))
}

func (s *AppOpt) Reset(b *ui.Button) {
	s.pref = ui.GetUi().ApplyPreferences(NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *AppOpt) Apply(b *ui.Button) {
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	ui.Pop()
}

func (r *AppOpt) Update(dt int) {
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *AppOpt) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *AppOpt) Entered() {
	s.Setup(LoadPreferences())
	s.Resize()
}

func (s *AppOpt) Resize() {
	s.topBar.Resize()
	w, h := ui.GetUi().GetScreenSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/2-hTop*6
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.optFullScr.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optLang.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optPause.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optFeeback.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optTheme.Resize([]int{x, y, w1, hTop*4 - 2})
}

func (r *AppOpt) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
