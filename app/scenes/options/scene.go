package options

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
)

var (
	strPosModalKey = "Клавиша для модальности Позиции"
	strColModalKey = "Клавиша для модальности Цвета"
	strNumModalKey = "Клавиша для модальности Цифры"
	strAriModalKey = "Клавиша для модальности Арифметика"
)

type SceneOptions struct {
	eui.SceneBase
	topbar                                                 *eui.TopBar
	btnApply, btnReset                                     *eui.Button
	optRestDelay                                           *eui.ComboBox
	optPosModKey, optColModKey, optNumModKey, optAriModKey *InputKey
	restDelay                                              int
}

func NewSceneOptions() *SceneOptions {
	s := &SceneOptions{}
	s.topbar = eui.NewTopBar("Настройки нназад", nil)
	s.Add(s.topbar)
	dt := []interface{}{1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}
	s.restDelay = dt[1].(int)
	s.optRestDelay = eui.NewComboBox("Обязательная пауза после сессии", dt, 1, func(c *eui.ComboBox) {
		s.restDelay = c.Value().(int)
	})
	s.Add(s.optRestDelay)
	s.optPosModKey = NewInputKey(strPosModalKey)
	s.Add(s.optPosModKey)
	s.optColModKey = NewInputKey(strColModalKey)
	s.Add(s.optColModKey)
	s.optNumModKey = NewInputKey(strNumModalKey)
	s.Add(s.optNumModKey)
	s.optAriModKey = NewInputKey(strAriModalKey)
	s.Add(s.optAriModKey)
	s.btnApply = eui.NewButton("Применить", func(b *eui.Button) {
		appOpt := eui.GetUi().GetSettings()
		appOpt.Set(app.PositionKeypress, s.optPosModKey.Value())
		appOpt.Set(app.ColorKeypress, s.optColModKey.Value())
		appOpt.Set(app.NumberKeypress, s.optNumModKey.Value())
		appOpt.Set(app.AriphmeticsKeypress, s.optAriModKey.Value())
		appOpt.Set(app.RestDuration, s.restDelay)
	})
	s.Add(s.btnApply)
	s.btnReset = eui.NewButton("Обнулить", func(b *eui.Button) {
		app.SetConf()
		s.resetOpt()
	})
	s.Add(s.btnReset)
	return s
}

func (s *SceneOptions) Entered() {
	s.Resize()
	s.resetOpt()
}

func (s *SceneOptions) resetOpt() {
	appOpt := eui.GetUi().GetSettings()
	s.optRestDelay.SetValue(appOpt.Get(app.RestDuration))
	s.optPosModKey.SetValue(appOpt.Get(app.PositionKeypress).(ebiten.Key))
	s.optColModKey.SetValue(appOpt.Get(app.ColorKeypress).(ebiten.Key))
	s.optNumModKey.SetValue(appOpt.Get(app.NumberKeypress).(ebiten.Key))
	s.optAriModKey.SetValue(appOpt.Get(app.AriphmeticsKeypress).(ebiten.Key))
}

func (s *SceneOptions) Resize() {
	w0, h0 := eui.GetUi().Size()
	h1 := int(float64(h0) * 0.068)
	w1 := w0 - w0/5
	x, y := 0, 0
	s.topbar.Resize([]int{x, y, w0, h1})
	y += h1
	s.optRestDelay.Resize([]int{x, y, w1, h1})
	y += h1
	s.optPosModKey.Resize([]int{x, y, w1, h1})
	y += h1
	s.optColModKey.Resize([]int{x, y, w1, h1})
	y += h1
	s.optNumModKey.Resize([]int{x, y, w1, h1})
	y += h1
	s.optAriModKey.Resize([]int{x, y, w1, h1})
	y = h0 - h1
	s.btnApply.Resize([]int{x, y, w1 / 2, h1})
	s.btnReset.Resize([]int{x + w1/2, y, w1 / 2, h1})
}
