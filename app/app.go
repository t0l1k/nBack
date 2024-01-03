package app

import "github.com/t0l1k/eui"

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("nBack")
	k := 30
	w, h := 16*k, 9*k
	u.SetSize(w, h)
	setConf(u)
	setTheme(u)
	return u
}

func setTheme(u *eui.Ui) {
	u.GetTheme().Set(eui.SceneBg, eui.Black)
	u.GetTheme().Set(LabelColorDefault, eui.Silver)
	u.GetTheme().Set(GameColorBg, eui.Black)
	u.GetTheme().Set(GameColorActiveBg, eui.Yellow)
	u.GetTheme().Set(GameColorFg, eui.Navy)
	u.GetTheme().Set(GameColorFgCrosshair, eui.Silver)
	u.GetTheme().Set(ColorNeutral, eui.Green)
	u.GetTheme().Set(ColorCorrect, eui.Blue)
	u.GetTheme().Set(ColorWrong, eui.Orange)
	u.GetTheme().Set(ColorMissed, eui.Red)
}

func setConf(u *eui.Ui) {
	u.GetSettings().Set(RestDuration, 5)
	u.GetSettings().Set(ShowCellPercent, 0.65)
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
	ShowCellPercent
)
