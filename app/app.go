package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
)

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
	SymbolKeypress
	AudKeypress
	AppLang
)

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("nBack")
	k := 30
	w, h := 16*k, 9*k
	u.SetSize(w, h)
	SetDefaultConf()
	setTheme()
	return u
}

func setTheme() {
	u := eui.GetUi()
	u.GetTheme().Set(eui.SceneBg, colors.Black)
	u.GetTheme().Set(LabelColorDefault, colors.Silver)
	u.GetTheme().Set(GameColorBg, colors.Black)
	u.GetTheme().Set(GameColorActiveBg, colors.Teal)
	u.GetTheme().Set(GameColorFg, colors.Navy)
	u.GetTheme().Set(GameColorFgCrosshair, colors.Silver)
	u.GetTheme().Set(ColorNeutral, colors.Green)
	u.GetTheme().Set(ColorCorrect, colors.Blue)
	u.GetTheme().Set(ColorWrong, colors.Orange)
	u.GetTheme().Set(ColorMissed, colors.Red)
}

func SetDefaultConf() {
	u := eui.GetUi()
	u.GetSettings().Set(eui.UiFullscreen, false)
	u.GetSettings().Set(RestDuration, 3)
	u.GetSettings().Set(PositionKeypress, ebiten.KeyA)
	u.GetSettings().Set(ColorKeypress, ebiten.KeyC)
	u.GetSettings().Set(SymbolKeypress, ebiten.KeyS)
	u.GetSettings().Set(AudKeypress, ebiten.KeyR)
	u.GetSettings().Set(AppLang, "ru")
}
