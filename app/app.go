package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("nBack")
	k := 30
	w, h := 16*k, 9*k
	u.SetSize(w, h)
	SetConf()
	setTheme()
	return u
}

func setTheme() {
	u := eui.GetUi()
	u.GetTheme().Set(eui.SceneBg, eui.Black)
	u.GetTheme().Set(LabelColorDefault, eui.Silver)
	u.GetTheme().Set(GameColorBg, eui.Black)
	u.GetTheme().Set(GameColorActiveBg, eui.Teal)
	u.GetTheme().Set(GameColorFg, eui.Navy)
	u.GetTheme().Set(GameColorFgCrosshair, eui.Silver)
	u.GetTheme().Set(ColorNeutral, eui.Green)
	u.GetTheme().Set(ColorCorrect, eui.Blue)
	u.GetTheme().Set(ColorWrong, eui.Orange)
	u.GetTheme().Set(ColorMissed, eui.Red)
}

func SetConf() {
	u := eui.GetUi()
	u.GetSettings().Set(eui.UiFullscreen, false)
	u.GetSettings().Set(RestDuration, 3)
	u.GetSettings().Set(PositionKeypress, ebiten.KeyA)
	u.GetSettings().Set(ColorKeypress, ebiten.KeyC)
	u.GetSettings().Set(NumberKeypress, ebiten.KeyS)
	u.GetSettings().Set(AriphmeticsKeypress, ebiten.KeyR)
}

const (
	LabelColorDefault eui.ThemeValue = iota + 200
	ColorNeutral
	ColorCorrect
	ColorWrong
	ColorMissed
	GameColorBg
	GameColorActiveBg
	GameColorFg
	GameColorFgCrosshair
)

const (
	RestDuration eui.SettingName = iota + 100
	PositionKeypress
	ColorKeypress
	NumberKeypress
	AriphmeticsKeypress
)
