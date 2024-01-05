package data

import (
	"fmt"
	"image/color"
	"math"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
)

type GamesData struct {
	id    int
	Games []*GameData
	Conf  *GameConf
}

func NewGamesData() *GamesData {
	g := &GamesData{id: 0, Conf: DefaultSettings()}
	return g
}

func (g *GamesData) Setup(mods []string, level, lives, adv, fall int, moveTime float64) {
	var modals []*Modality
	for _, mod := range mods {
		modals = append(modals, NewModality(mod))
	}
	g.id = len(g.Games)
	gData := NewGame(
		g.id,
		modals,
		level,
		lives,
		g.getTotalMoves(level),
		adv,
		fall,
		moveTime,
	)
	g.Games = append(g.Games, gData)
}

func (g *GamesData) NewGame(level, lives int) {
	lastGame := g.Last().SetupNext()
	for _, v := range lastGame.Modalities {
		v.Reset()
	}
	g.id = len(g.Games)
	gData := NewGame(
		g.id,
		lastGame.Modalities,
		level,
		lives,
		g.getTotalMoves(level),
		lastGame.Advance,
		lastGame.Fallback,
		lastGame.MoveTime,
	)
	g.Games = append(g.Games, gData)
}

func (g *GamesData) Last() *GameData {
	return g.Games[g.id]
}

func (g *GamesData) getTotalMoves(level int) int {
	trials := g.Conf.Get(Trials).(int)
	factor := g.Conf.Get(TrialsFactor).(int)
	exponent := g.Conf.Get(TrialsExponent).(int)
	return trials + factor*int(math.Pow(float64(level), float64(exponent)))
}

func (g *GamesData) NextLevel() (level, lives int, result string, col color.Color) {
	level = g.Games[g.Last().Id].Level
	lives = g.Games[g.Last().Id].Lives
	percent := g.Games[g.Last().Id].Percent
	adv := g.Games[g.Last().Id].Advance
	fall := g.Games[g.Last().Id].Fallback
	theme := eui.GetUi().GetTheme()
	if percent >= adv {
		level++
		lives = g.Conf.Get(ThresholdFallbackSessions).(int)
		result = fmt.Sprintf("Уровень(%v) пройден отлично, вверх на(%v)!", level-1, level)
		col = theme.Get(app.ColorCorrect)
	} else if percent >= fall && percent < adv {
		result = fmt.Sprintf("Играть уровень(%v) снова!", level)
		col = theme.Get(app.ColorNeutral)
	} else if percent < fall {
		if lives > 1 {
			lives--
			result = fmt.Sprintf("Играть уровень(%v) снова! Попыток осталось(%v)", level, lives)
			col = theme.Get(app.ColorWrong)
		} else {
			if level > 1 {
				level--
				lives = g.Conf.Get(ThresholdFallbackSessions).(int)
			}
			result = fmt.Sprintf("Уровень вниз(%v)!", level)
			col = theme.Get(app.ColorMissed)
		}
	}
	return level, lives, result, col
}
