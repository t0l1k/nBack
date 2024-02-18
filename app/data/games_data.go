package data

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/game"
)

type GamesData struct {
	id    int
	Games []*game.GameData
	Conf  *game.GameConf
}

func NewGamesData(conf *game.GameConf) *GamesData {
	var modals []*game.Modality
	for _, mod := range conf.Get(game.Modals).(game.ModalType) {
		modals = append(modals, game.NewModality(game.ModalType(mod)))
	}
	g := &GamesData{id: 0, Conf: conf}
	level := g.Conf.Get(game.DefaultLevel).(int)
	tryUp := g.Conf.Get(game.ThresholdAdvanceSessions).(int)
	tryDown := g.Conf.Get(game.ThresholdFallbackSessions).(int)
	stepUp := g.Conf.Get(game.ThresholdAdvance).(int)
	stepDown := g.Conf.Get(game.ThresholdFallback).(int)
	moveTime := g.Conf.Get(game.MoveTime).(float64)
	resetOnWrong := g.Conf.Get(game.ResetOnFirstWrong).(bool)
	g.id = len(g.Games)
	gData := game.NewGame(
		g.id,
		modals,
		level,
		tryUp,
		tryDown,
		g.getTotalMoves(level),
		stepUp,
		stepDown,
		moveTime,
		resetOnWrong,
	)
	g.Games = append(g.Games, gData)
	gData.FillField(*g.Conf)
	return g
}

func (g *GamesData) NewGame(level, tryUp, tryDown int) {
	lastGame := g.Last().SetupNext()
	for _, v := range lastGame.Modalities {
		v.Reset()
	}
	g.id = len(g.Games)
	gData := game.NewGame(
		g.id,
		lastGame.Modalities,
		level,
		tryUp,
		tryDown,
		g.getTotalMoves(level),
		lastGame.Advance,
		lastGame.Fallback,
		lastGame.MoveTime,
		lastGame.ResetOnWrong,
	)
	gData.FillField(*g.Conf)
	g.Games = append(g.Games, gData)
}

func (g *GamesData) Id() int {
	return g.id
}

func (g *GamesData) Last() *game.GameData {
	return g.Games[g.id]
}

func (g GamesData) getTotalMoves(level int) int {
	trials := g.Conf.Get(game.Trials).(int)
	factor := g.Conf.Get(game.TrialsFactor).(int)
	exponent := g.Conf.Get(game.TrialsExponent).(int)
	return trials + factor*int(math.Pow(float64(level), float64(exponent)))
}

func (g *GamesData) calcGameFor(id int) (level, tryUp, tryDown int, result string, col color.Color) {
	level = g.Games[id].Level
	tryUp = g.Games[id].TryUp
	tryDown = g.Games[id].TryDown
	percent := g.Games[id].Percent
	stepUp := g.Games[id].Advance
	stepDown := g.Games[id].Fallback
	theme := eui.GetUi().GetTheme()
	if percent >= stepUp {
		tryUp--
		if tryUp == 0 {
			level++
			tryUp = g.Conf.Get(game.ThresholdAdvanceSessions).(int)
			tryDown = g.Conf.Get(game.ThresholdFallbackSessions).(int)
			result = fmt.Sprintf("Уровень(%v) пройден отлично, вверх на(%v)!", level-1, level)
			col = theme.Get(app.ColorCorrect)
		} else {
			result = fmt.Sprintf("Играть уровень(%v) снова, ещё успешных попыток(%v)!", level, tryUp)
			col = theme.Get(app.ColorNeutral)
		}
	} else if percent >= stepDown && percent < stepUp {
		tryUp = g.Conf.Get(game.ThresholdAdvanceSessions).(int)
		result = fmt.Sprintf("Играть уровень(%v) снова!", level)
		col = theme.Get(app.ColorNeutral)
	} else if percent < stepDown {
		tryUp = g.Conf.Get(game.ThresholdAdvanceSessions).(int)
		if tryDown > 1 {
			tryDown--
			result = fmt.Sprintf("Играть уровень(%v) снова! Попыток осталось(%v)", level, tryDown)
			col = theme.Get(app.ColorWrong)
		} else {
			if level > 1 {
				level--
				tryDown = g.Conf.Get(game.ThresholdFallbackSessions).(int)
			}
			result = fmt.Sprintf("Уровень вниз(%v)!", level)
			col = theme.Get(app.ColorMissed)
		}
	}
	return level, tryUp, tryDown, result, col
}

func (g *GamesData) NextLevel() (level, tryUp, tryDown int, result string, col color.Color) {
	level, tryUp, tryDown, result, col = g.calcGameFor(g.Last().Id)
	return level, tryUp, tryDown, result, col
}

func (g *GamesData) PrevGame() (level, tryUp, tryDown int, result string, col color.Color) {
	level, tryUp, tryDown, result, col = g.calcGameFor(g.Last().Id - 1)
	return level, tryUp, tryDown, result, col
}

func (g GamesData) String() (result string) {
	var (
		max       int
		avg       float64
		durration time.Duration
	)
	for _, v := range g.Games {
		if !v.IsDone() {
			continue
		}
		durration += v.Duration
		avg += float64(v.Level)
		if v.Level > max {
			max = v.Level
		}
	}
	avg = avg / float64(len(g.Games))
	mSec := durration.Milliseconds() / 1e3
	sec := durration.Seconds()
	m := int(sec / 60)
	seconds := int(sec) % 60
	gameDuration := fmt.Sprintf("%02v:%02v.%03v", m, seconds, int(mSec))
	result = fmt.Sprintf("%v #%v max:%v avg:%0.2v [%v]", time.Now().Format("2006-01-02"), g.id, max, avg, gameDuration)
	return result
}
