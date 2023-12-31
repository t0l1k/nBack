package data

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
)

type GameData struct {
	DtBeg, DtEnd               string // format "2006-01-02 15:04:05.000"
	Modalities                 []*Modality
	Id, Level, Lives           int
	Moves, TotalMoves          int
	Percent, Advance, Fallback int
	Manual, ResetOnError, done bool
	Duration                   time.Duration
	MoveTime                   float64
}

func NewGame(id int, mods []*Modality, level, lives, totalMoves, advance, fallback int, resetOnError, manual bool, moveTime float64) *GameData {
	g := &GameData{
		Id:           id,
		Modalities:   mods,
		Level:        level,
		Lives:        lives,
		TotalMoves:   totalMoves,
		ResetOnError: resetOnError,
		Manual:       manual,
		Advance:      advance,
		Fallback:     fallback,
		MoveTime:     moveTime,
	}
	return g
}

func (g *GameData) IsDone() bool { return g.done }

func (g *GameData) SetGameDone(moves int) {
	g.DtEnd = time.Now().Format("2006-01-02 15:04:05.000")
	g.Moves = moves
	g.calcPercent()
	g.done = true
}

func (g *GameData) calcPercent() {
	var (
		aa, bb, i, j   float64
		correct, error float64
	)

	for _, v := range g.Modalities {
		correct += float64(v.correct)
		error += float64(v.wrong + v.missed)
	}

	aa, bb = float64(correct), float64(error)
	if aa == 0 && bb == 0 {
		i, j = 1, 0
	} else if aa == 0 && bb > 0 {
		i, j = 0, 1
	} else {
		i, j = aa, bb
	}
	g.Percent = int(i * 100 / (i + j))
}

func (g GameData) SetupNext() GameData { return g }

func (g *GameData) NextLevel() (level, lives int, result string, col color.Color) {
	level = g.Level
	lives = g.Lives
	conf := eui.GetUi().GetSettings()
	if g.Percent >= g.Advance {
		level++
		lives = conf.Get(app.ThresholdFallbackSessions).(int)
		result = fmt.Sprintf("Уровень(%v) пройден отлично, вверх на(%v)!", level-1, level)
		col = conf.Get(app.ColorCorrect).(color.Color)
	} else if g.Percent >= g.Fallback && g.Percent < g.Advance {
		result = fmt.Sprintf("Играть уровень(%v) снова!", level)
		col = conf.Get(app.ColorNeutral).(color.Color)
	} else if g.Percent < g.Fallback {
		if lives > 1 {
			lives--
			result = fmt.Sprintf("Играть уровень(%v) снова! Попыток осталось(%v)", level, lives)
			col = conf.Get(app.ColorWrong).(color.Color)
		} else {
			if level > 1 {
				level--
				lives = conf.Get(app.ThresholdFallbackSessions).(int)
			}
			result = fmt.Sprintf("Уровень вниз(%v)!", level)
			col = conf.Get(app.ColorMissed).(color.Color)
		}
	}
	return level, lives, result, col
}

func (g *GameData) GetModalities() []*Modality {
	return g.Modalities
}

func (g *GameData) IsContainMod(mod string) bool {
	for _, v := range g.Modalities {
		if v.GetSym() == mod {
			return true
		}
	}
	return false
}

func (g *GameData) GetModalityValues(mod string) []int {
	for _, v := range g.Modalities {
		if v.GetSym() == mod {
			return v.GetField()
		}
	}
	return nil
}

func (g *GameData) GameMode() (result string) {
	switch len(g.Modalities) {
	case 1:
		result = g.Modalities[0].String() + strconv.Itoa(g.Level)
	default:
		for _, v := range g.Modalities {
			result += v.GetSym()
		}
		result += strconv.Itoa(g.Level)
	}
	return result
}

func (g *GameData) ShortResultStringWithColors() (str string, bg, fg color.Color) {
	str = fmt.Sprintf("#%v %v %v%%", g.Id, g.GameMode(), g.Percent)
	conf := eui.GetUi().GetSettings()
	clrNeutral := conf.Get(app.ColorNeutral).(color.Color)
	clrCorrect := conf.Get(app.ColorCorrect).(color.Color)
	clrWrong := conf.Get(app.ColorWrong).(color.Color)
	clrMissed := conf.Get(app.ColorMissed).(color.Color)
	fg = eui.Black
	if g.Percent >= g.Advance {
		bg = clrCorrect
	} else if g.Percent >= g.Fallback && g.Percent < g.Advance {
		bg = clrNeutral
	} else if g.Percent < g.Fallback {
		if g.Lives > 1 {
			bg = clrWrong
		} else {
			bg = clrMissed
		}
	}
	return str, bg, fg
}

func (g *GameData) String() string {
	dtFormat := "2006-01-02 15:04:05.000"
	dtBeg, err := time.Parse(dtFormat, g.DtBeg)
	if err != nil {
		panic(err)
	}
	dtEnd, err := time.Parse(dtFormat, g.DtEnd)
	if err != nil {
		panic(err)
	}
	g.Duration = dtEnd.Sub(dtBeg)
	mSec := g.Duration.Milliseconds() / 1e3
	sec := g.Duration.Seconds()
	m := int(sec / 60)
	seconds := int(sec) % 60
	gameDuration := fmt.Sprintf("%02v:%02v.%03v", m, seconds, int(mSec))
	s1 := "[Modality Correct-(Wrong-Missed)]"
	for _, v := range g.Modalities {
		s1 += "[" + v.sym + ":" + strconv.Itoa(v.correct) + "(" + strconv.Itoa(v.wrong) + "-" + strconv.Itoa(v.missed) + ")] "
	}
	return fmt.Sprintf("#%v %v score: %v%% %v moves(%v/%v) %v", g.Id, g.GameMode(), g.Percent, s1, g.Moves, g.TotalMoves, gameDuration)
}