package data

import (
	"math"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
)

type GamesData struct {
	Data []*GameData
	id   int
}

func NewGameJaeggiPos3x3Rulez() *GamesData {
	g := &GamesData{id: 0}
	eui.GetUi().GetSettings().Set(app.Trials, 20)
	g.Setup([]string{Pos}, 1, 1, 90, 75, false, false, 1.5)
	return g
}

func NewGameJaeggiPos3x3ColRulez() *GamesData {
	g := &GamesData{id: 0}
	eui.GetUi().GetSettings().Set(app.Trials, 20)
	g.Setup([]string{Pos, Col}, 1, 1, 90, 75, false, false, 1.5)
	return g
}

func NewGameDevelRulez() *GamesData {
	g := &GamesData{id: 0}
	g.Setup([]string{Pos, Col}, 1, 1, 90, 75, false, false, 2.5)
	return g
}

func (g *GamesData) Setup(mods []string, level, lives, adv, fall int, reset, manual bool, moveTime float64) {
	var modals []*Modality
	for _, mod := range mods {
		modals = append(modals, NewModality(mod))
	}
	g.id = len(g.Data)
	gData := NewGame(
		g.id,
		modals,
		level,
		lives,
		g.getTotalMoves(level),
		adv,
		fall,
		reset,
		manual,
		moveTime,
	)
	g.Data = append(g.Data, gData)
}

func (g *GamesData) NewGame(level, lives int) {
	lastGame := g.Last().SetupNext()
	for _, v := range lastGame.Modalities {
		v.Reset()
	}
	g.id = len(g.Data)
	gData := NewGame(
		g.id,
		lastGame.Modalities,
		level,
		lives,
		g.getTotalMoves(level),
		lastGame.Advance,
		lastGame.Fallback,
		lastGame.ResetOnError,
		lastGame.Manual,
		lastGame.MoveTime,
	)
	g.Data = append(g.Data, gData)
}

func (g *GamesData) Last() *GameData {
	return g.Data[g.id]
}

func (g *GamesData) getTotalMoves(level int) int {
	conf := eui.GetUi().GetSettings()
	trials := conf.Get(app.Trials).(int)
	factor := conf.Get(app.TrialsFactor).(int)
	exponent := conf.Get(app.TrialsExponent).(int)
	return trials + factor*int(math.Pow(float64(level), float64(exponent)))
}
