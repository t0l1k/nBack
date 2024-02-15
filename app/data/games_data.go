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
	lives := g.Conf.Get(game.ThresholdFallbackSessions).(int)
	adv := g.Conf.Get(game.ThresholdAdvance).(int)
	fall := g.Conf.Get(game.ThresholdFallback).(int)
	moveTime := g.Conf.Get(game.MoveTime).(float64)
	g.id = len(g.Games)
	gData := game.NewGame(
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
	return g
}

func (g *GamesData) NewGame(level, lives int) {
	lastGame := g.Last().SetupNext()
	for _, v := range lastGame.Modalities {
		v.Reset()
	}
	g.id = len(g.Games)
	gData := game.NewGame(
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

func (g *GamesData) Id() int {
	return g.id
}

func (g *GamesData) Last() *game.GameData {
	return g.Games[g.id]
}

func (g *GamesData) getTotalMoves(level int) int {
	trials := g.Conf.Get(game.Trials).(int)
	factor := g.Conf.Get(game.TrialsFactor).(int)
	exponent := g.Conf.Get(game.TrialsExponent).(int)
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
		lives = g.Conf.Get(game.ThresholdFallbackSessions).(int)
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
				lives = g.Conf.Get(game.ThresholdFallbackSessions).(int)
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
