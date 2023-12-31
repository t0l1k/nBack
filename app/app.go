package app

import "github.com/t0l1k/eui"

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("nBack")
	k := 30
	w, h := 16*k, 9*k
	u.SetSize(w, h)
	setConf(u)
	return u
}

func setConf(u *eui.Ui) {
	u.GetTheme().Set(eui.SceneBg, eui.Black)
	u.GetSettings().Set(DefaultLevel, 1)
	u.GetSettings().Set(MoveTime, 1.5)
	u.GetSettings().Set(ShowCellPercent, 0.65)
	u.GetSettings().Set(LabelColorDefault, eui.Silver)
	u.GetSettings().Set(ColorNeutral, eui.Green)
	u.GetSettings().Set(ColorCorrect, eui.Blue)
	u.GetSettings().Set(ColorWrong, eui.Orange)
	u.GetSettings().Set(ColorMissed, eui.Red)
	u.GetSettings().Set(RestDuration, 5)
	u.GetSettings().Set(RandomRepition, 30)
	u.GetSettings().Set(GridSize, 3)
	u.GetSettings().Set(ShowGrid, false)
	u.GetSettings().Set(UseCenterCell, false)
	u.GetSettings().Set(ShowCrossHair, true)
	u.GetSettings().Set(ResetOnFirstWrong, false)
	u.GetSettings().Set(ThresholdAdvance, 90)
	u.GetSettings().Set(ThresholdFallback, 75)
	u.GetSettings().Set(ThresholdFallbackSessions, 1)
	u.GetSettings().Set(Trials, 5)
	u.GetSettings().Set(TrialsFactor, 1)
	u.GetSettings().Set(TrialsExponent, 1)
	u.GetSettings().Set(MaxNumber, 10)
	u.GetSettings().Set(UseAddSub, true)
	u.GetSettings().Set(UseMulDiv, false)
	u.GetSettings().Set(GameColorBg, eui.Black)
	u.GetSettings().Set(GameColorActiveBg, eui.Blue)
	u.GetSettings().Set(GameColorFg, eui.Navy)
	u.GetSettings().Set(GameColorFgCrosshair, eui.Silver)
}

const (
	DefaultLevel eui.SettingName = iota + 100
	MoveTime
	ShowCellPercent
	LabelColorDefault
	ColorNeutral
	ColorCorrect
	ColorWrong
	ColorMissed
	RestDuration
	RandomRepition
	GridSize
	ShowGrid
	UseCenterCell
	ShowCrossHair
	ResetOnFirstWrong
	ThresholdAdvance
	ThresholdFallback
	ThresholdFallbackSessions
	Trials
	TrialsFactor
	TrialsExponent
	MaxNumber
	UseAddSub
	UseMulDiv
	GameColorBg
	GameColorActiveBg
	GameColorFg
	GameColorFgCrosshair
)
