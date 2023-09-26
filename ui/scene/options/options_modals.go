package options

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui/app"
)

type OptModals struct {
	eui.ContainerDefault
	topBar                                   *TopBarOpt
	lblResult                                *eui.Label
	optGameType                              *eui.Combobox
	optGridSize, optMaxSym, optMaxAriphmetic *eui.Combobox
	optShowCross, optCenterCell, optShowGrid *eui.Checkbox
	optColors                                *eui.Icon
	pref                                     *eui.Preferences
}

func NewOptModals() *OptModals {
	s := &OptModals{}
	s.pref = app.LoadPreferences()
	rect := []int{0, 0, 1, 1}
	s.topBar = NewTopBarOpt(s.Reset, s.Apply)
	s.Add(s.topBar)

	s.lblResult = eui.NewLabel("Настройка модальности", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblResult)

	gamesType := []interface{}{game.Pos, game.Col, game.Sym, game.Ari}
	idx := 0
	s.optGameType = eui.NewCombobox(s.getGameType(), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), gamesType, idx, func(b *eui.Combobox) {
		s.pref.Set("game type", s.optGameType.Value().(string))
		s.selectWhatOptShow()
		s.optGameType.SetText(s.getGameType())
		s.lblResult.SetText(fmt.Sprintf("Настроить %v для уровня:%v, ходов:%v", s.getGameType(), s.pref.Get("default level").(int), game.TotalMoves(s.pref.Get("default level").(int))))
	})
	s.Add(s.optGameType)

	s.optShowCross = eui.NewCheckbox(eui.GetLocale().Get("optcross"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), func(c *eui.Checkbox) {
		s.pref.Set("show crosshair", s.optShowCross.Checked())
		log.Printf("Show crosshair: %v", s.pref.Get("show crosshair").(bool))
	})
	s.Add(s.optShowCross)

	arrMaxSymbols := []interface{}{10, 20, 50, 100, 200, 500, 1000}
	s.optMaxSym = eui.NewCombobox(eui.GetLocale().Get("optmaxsym"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arrMaxSymbols, 3, func(c *eui.Combobox) {
		s.pref.Set("symbols count", s.optMaxSym.Value().(int))
	})
	s.Add(s.optMaxSym)

	s.optMaxAriphmetic = eui.NewCombobox(eui.GetLocale().Get("optmaxari"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), arrMaxSymbols, 1, func(c *eui.Combobox) {
		s.pref.Set("ariphmetic max", s.optMaxAriphmetic.Value().(int))
	})
	s.Add(s.optMaxAriphmetic)

	s.optCenterCell = eui.NewCheckbox(eui.GetLocale().Get("optcc"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), func(c *eui.Checkbox) {
		s.pref.Set("use center cell", s.optCenterCell.Checked())
		log.Printf("Use center cell: %v", s.pref.Get("use center cell").(bool))
	})
	s.Add(s.optCenterCell)

	s.optShowGrid = eui.NewCheckbox(eui.GetLocale().Get("optgrid"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), func(c *eui.Checkbox) {
		s.pref.Set("show grid", s.optShowGrid.Checked())
		log.Printf("Show Grid: %v", s.pref.Get("show grid").(bool))
	})
	s.Add(s.optShowGrid)

	lvls := []interface{}{2, 3, 4, 5, 6, 7, 8, 9}
	idx = 1
	s.optGridSize = eui.NewCombobox(eui.GetLocale().Get("optgridsz"), rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"), lvls, idx, func(c *eui.Combobox) {
		s.pref.Set("grid size", s.optGridSize.Value().(int))
		log.Println("Grid Size changed")
	})
	s.Add(s.optGridSize)

	s.optColors = eui.NewIcon(nil, rect)
	s.Add(s.optColors)
	s.optColors.SetIcon(s.getColorsIcon())
	return s
}

func (s *OptModals) getColorsIcon() *ebiten.Image {
	var w0, h0 int
	sz := len(game.Colors)
	w0, h0 = sz*sz, sz
	cellWidth, cellHeight := w0, h0
	image := ebiten.NewImage(w0, h0)
	y := 0
	w, h := cellWidth/len(game.Colors), cellHeight
	for i, v := range game.Colors {
		cellX := i % sz * w
		ebitenutil.DrawRect(image, float64(cellX), float64(y), float64(w), float64(h), v)
	}
	eui.DrawRect(image, eui.NewRect([]int{0, 0, w0, h0}), eui.White)
	return image
}

func (s *OptModals) selectWhatOptShow() {
	tp := eui.GetPreferences().Get("game type").(string)
	switch tp {
	case game.Ari:
		s.optShowCross.Visible = true
		s.optMaxAriphmetic.Visible = true
		s.optMaxSym.Visible = false
		s.optCenterCell.Visible = false
		s.optShowGrid.Visible = false
		s.optGridSize.Visible = false
		s.optColors.Visible = false
	case game.Sym:
		s.optShowCross.Visible = true
		s.optMaxAriphmetic.Visible = false
		s.optMaxSym.Visible = true
		s.optCenterCell.Visible = false
		s.optShowGrid.Visible = false
		s.optGridSize.Visible = false
		s.optColors.Visible = false
	case game.Col:
		s.optShowCross.Visible = true
		s.optMaxAriphmetic.Visible = false
		s.optMaxSym.Visible = false
		s.optCenterCell.Visible = false
		s.optShowGrid.Visible = false
		s.optGridSize.Visible = false
		s.optColors.Visible = true
	case game.Pos:
		s.optShowCross.Visible = true
		s.optMaxAriphmetic.Visible = false
		s.optMaxSym.Visible = false
		s.optCenterCell.Visible = true
		s.optShowGrid.Visible = true
		s.optGridSize.Visible = true
		s.optColors.Visible = false
	}
}

func (s *OptModals) getGameType() string {
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

func (s *OptModals) Setup(sets *eui.Preferences) {
	s.selectWhatOptShow()
	s.optGameType.SetValue(sets.Get("game type").(string))
	s.optShowCross.SetChecked(sets.Get("show crosshair").(bool))
	s.optGridSize.SetValue(sets.Get("grid size").(int))
	s.optShowGrid.SetChecked(sets.Get("show grid").(bool))
	s.optCenterCell.SetChecked(sets.Get("use center cell").(bool))
	s.optMaxSym.SetValue(sets.Get("symbols count").(int))
	s.optMaxAriphmetic.SetValue(sets.Get("ariphmetic max").(int))
	s.lblResult.SetText(fmt.Sprintf("Настроить %v для уровня:%v, ходов:%v", s.getGameType(), s.pref.Get("default level").(int), game.TotalMoves(s.pref.Get("default level").(int))))
}

func (s *OptModals) Reset(b *eui.Button) {
	s.pref = eui.GetUi().ApplyPreferences(app.NewPref())
	s.Setup(s.pref)
	log.Println("Reset All Options to Defaults")
}

func (s *OptModals) Apply(b *eui.Button) {
	data.GetDb().InsertSettings(s.pref)
	log.Println("Apply Settings")
	eui.Pop()
}

func (r *OptModals) Update(dt int) {
	for _, value := range r.Container {
		value.Update(dt)
	}
}

func (r *OptModals) Draw(surface *ebiten.Image) {
	for _, value := range r.Container {
		value.Draw(surface)
	}
}

func (s *OptModals) Entered() {
	s.Setup(app.LoadPreferences())
	s.Resize()
}

func (s *OptModals) Resize() {
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
	s.optShowCross.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optMaxSym.Resize([]int{x, y, w1, hTop - 2})
	s.optMaxAriphmetic.Resize([]int{x, y, w1, hTop - 2})
	s.optCenterCell.Resize([]int{x, y, w1, hTop - 2})
	s.optColors.Resize([]int{x, y, w1, hTop - 2})
	s.optColors.SetIcon(s.getColorsIcon())
	y += hTop
	s.optShowGrid.Resize([]int{x, y, w1, hTop - 2})
	y += hTop
	s.optGridSize.Resize([]int{x, y, w1, hTop - 2})
}

func (r *OptModals) Close() {
	for _, v := range r.Container {
		v.Close()
	}
}
