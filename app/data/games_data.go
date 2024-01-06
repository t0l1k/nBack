package data

import (
	"fmt"
	"image/color"
	"math"
	"time"

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

func (g *GamesData) calcGameFor(id int) (level, lives int, result string, col color.Color) {
	level = g.Games[id].Level
	lives = g.Games[id].Lives
	percent := g.Games[id].Percent
	adv := g.Games[id].Advance
	fall := g.Games[id].Fallback
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

func (g *GamesData) NextLevel() (level, lives int, result string, col color.Color) {
	level, lives, result, col = g.calcGameFor(g.Last().Id)
	return level, lives, result, col
}

func (g *GamesData) PrevGame() (level, lives int, result string, col color.Color) {
	level, lives, result, col = g.calcGameFor(g.Last().Id - 1)
	return level, lives, result, col
}

func (g GamesData) String() (result string) {
	var (
		max, avg  int
		durration time.Duration
	)
	for _, v := range g.Games {
		if !v.IsDone() {
			continue
		}
		durration += v.Duration
		avg += v.Level
		if v.Level > max {
			max = v.Level
		}
	}
	avg = avg / len(g.Games)
	mSec := durration.Milliseconds() / 1e3
	sec := durration.Seconds()
	m := int(sec / 60)
	seconds := int(sec) % 60
	gameDuration := fmt.Sprintf("%02v:%02v.%03v", m, seconds, int(mSec))
	result = fmt.Sprintf("%v #%v max:%v avg:%v [%v]", time.Now().Format("2006-01-02"), g.id, max, avg, gameDuration)
	return result
}
