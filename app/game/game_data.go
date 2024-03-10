package game

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
)

type GameData struct {
	DtBeg, DtEnd                      string // format "2006-01-02 15:04:05.000"
	Modalities                        []*Modality
	Id, Level, TryUp, TryDown         int
	Moves, TotalMoves, TotalTime      int
	Score, Percent, Advance, Fallback int
	done, ResetOnWrong, CheckIn       bool
	Duration                          time.Duration
	MoveTime                          float64
	up, down, same                    int
}

func NewGame(id int, mods []*Modality, level, tryUp, tryDown, totalMoves, advance, fallback, totalTime int, moveTime float64, resetOnError, checkIn bool) *GameData {
	g := &GameData{
		Id:           id,
		Modalities:   mods,
		Level:        level,
		TryUp:        tryUp,
		TryDown:      tryDown,
		TotalMoves:   totalMoves,
		Advance:      advance,
		Fallback:     fallback,
		TotalTime:    totalTime,
		CheckIn:      checkIn,
		MoveTime:     moveTime,
		ResetOnWrong: resetOnError,
	}
	return g
}

func (g *GameData) CheckNextLevel(conf GameConf) (result string, col color.Color) {
	correct, wrong := g.getModalsResult()
	g.calcPercent(correct, wrong)
	level := g.Level
	tryUp := g.TryUp
	tryDown := g.TryDown
	percent := g.Percent
	stepUp := g.Advance
	stepDown := g.Fallback
	theme := eui.GetUi().GetTheme()
	if percent >= stepUp {
		tryUp--
		if tryUp == 0 {
			level++
			g.up++
			tryUp = conf.Get(ThresholdAdvanceSessions).(int)
			tryDown = conf.Get(ThresholdFallbackSessions).(int)
			result = fmt.Sprintf("Уровень(%v) пройден отлично, вверх на(%v)!", level-1, level)
			col = theme.Get(app.ColorCorrect)
		} else {
			result = fmt.Sprintf("Играть уровень(%v) снова, ещё успешных попыток(%v)!", level, tryUp)
			col = theme.Get(app.ColorNeutral)
			g.same++
		}
	} else if percent >= stepDown && percent < stepUp {
		tryUp = conf.Get(ThresholdAdvanceSessions).(int)
		result = fmt.Sprintf("Играть уровень(%v) снова!", level)
		col = theme.Get(app.ColorNeutral)
		g.same++
	} else if percent < stepDown {
		tryUp = conf.Get(ThresholdAdvanceSessions).(int)
		if tryDown > 1 {
			tryDown--
			result = fmt.Sprintf("Играть уровень(%v) снова! Попыток осталось(%v)", level, tryDown)
			col = theme.Get(app.ColorWrong)
		} else {
			if level > 1 {
				level--
				tryDown = conf.Get(ThresholdFallbackSessions).(int)
				g.down++
			}
			result = fmt.Sprintf("Уровень вниз(%v)!", level)
			col = theme.Get(app.ColorMissed)
		}
	}
	g.Level = level
	g.TryUp = tryUp
	g.TryDown = tryDown
	g.TotalMoves = g.getTotalMoves(conf)
	for _, v := range g.Modalities {
		v.ResetResults()
	}
	return result, col
}

func (g GameData) getTotalMoves(conf GameConf) int {
	trials := conf.Get(Trials).(int)
	factor := conf.Get(TrialsFactor).(int)
	exponent := conf.Get(TrialsExponent).(int)
	return trials + factor*int(math.Pow(float64(g.Level), float64(exponent)))
}

func (g *GameData) FillField(conf GameConf) {
	for _, v := range g.Modalities {
		v.AddField(newField(conf, g.Level, g.TotalMoves, v.GetSym()))
	}
}

func (g *GameData) IsDone() bool { return g.done }

// Игра на количество ходов, вычислить процент достаточно. Игра на время с проверкой после указаного числа ходов, и от результата переход вверх или вниз, или повтор уровня, уже за каждый нейтральный и правильный ответ начисляются очки(уровень текущей игры), в конце суммируется, процент вычисляется только при накоплении ходов равному числу ходов, что маловероятно, но возможно.
func (g *GameData) SetGameDone(moves int) {
	g.DtEnd = time.Now().Format("2006-01-02 15:04:05.000")
	correct, wrong, missed := 0.0, 0.0, 0.0
	if g.MoveTime > 0 && g.CheckIn {
		res := g.GetModalitiesScore()
		movs := g.GetModalitiesMoves()
		var (
			move   int
			result map[ModalType][]int
			arr    []int
		)
		result = make(map[ModalType][]int)
		for k, values := range res {
			arr = make([]int, 3)
			for i, value := range values {
				g.Score += value
				switch movs[k][i] {
				case AddCorrect:
					arr[0] += 1
					correct += 1
				case AddWrong:
					wrong += 1
					arr[1] += 1
				case AddMissed:
					missed += 1
					arr[2] += 1
				}
			}
			result[k] = arr
			move = len(values)
		}
		i := 0
		for k, v := range result {
			mod := g.Modalities[i]
			if mod.GetSym() == k {
				mod.correct = v[0]
				mod.wrong = v[1]
				mod.missed = v[2]
			}
			i++
		}
		g.Moves = move
		if g.TotalMoves != moves {
			g.Percent = g.Fallback
		}
	} else {
		correct, wrong = g.getModalsResult()
		g.calcPercent(correct, wrong)
		g.Moves = moves
	}
	g.done = true
}

func (g *GameData) getModalsResult() (correct, wrong float64) {
	for _, v := range g.Modalities {
		correct += float64(v.correct)
		wrong += float64(v.wrong + v.missed)
	}
	return correct, wrong
}

func (g *GameData) calcPercent(aa, bb float64) {
	var (
		i, j float64
	)
	if aa == 0 && bb == 0 {
		i, j = 1, 0
	} else if aa == 0 && bb > 0 {
		i, j = 0, 1
	} else {
		i, j = aa, bb
	}
	g.Percent = int(i * 100 / (i + j))
	if g.ResetOnWrong && g.Percent <= g.Advance {
		g.Percent = 0
	}
}

func (g GameData) SetupNext() GameData { return g }

func (g *GameData) GetModalities() []*Modality {
	return g.Modalities
}

func (g *GameData) GetModalitiesScore() (moves map[ModalType][]int) {
	moves = make(map[ModalType][]int)
	for _, v := range g.Modalities {
		moves[v.GetSym()] = v.score
	}
	return moves
}

func (g *GameData) GetModalitiesMoves() (moves map[ModalType][]MoveType) {
	moves = make(map[ModalType][]MoveType)
	for _, v := range g.Modalities {
		moves[v.GetSym()] = v.GetMovesStatus()
	}
	return moves
}

func (g *GameData) IsContainMod(mod ModalType) bool {
	for _, v := range g.Modalities {
		if v.GetSym() == mod {
			return true
		}
	}
	return false
}

func (g *GameData) GetModalityValues(mod ModalType) []int {
	for _, v := range g.Modalities {
		if v.GetSym() == mod {
			return v.GetField()
		}
	}
	return nil
}

func (g *GameData) GameMode() (result ModalType) {
	switch len(g.Modalities) {
	case 1:
		result = ModalType(g.Modalities[0].String() + strconv.Itoa(g.Level))
	default:
		for _, v := range g.Modalities {
			result += v.GetSym()
		}
		result += ModalType(strconv.Itoa(g.Level))
	}
	return result
}

func (g *GameData) ShortResultStringWithColors() (str string, bg, fg color.Color) {
	if g.TotalTime > 0 && g.CheckIn {
		str = fmt.Sprintf("#%v %v %v", g.Id, g.GameMode(), g.Score)
	} else {
		str = fmt.Sprintf("#%v %v %v%%", g.Id, g.GameMode(), g.Percent)
	}
	fg, bg = g.getResultColors()
	return str, bg, fg
}

func (g *GameData) getResultColors() (fg, bg color.Color) {
	theme := eui.GetUi().GetTheme()
	clrNeutral := theme.Get(app.ColorNeutral)
	clrCorrect := theme.Get(app.ColorCorrect)
	clrWrong := theme.Get(app.ColorWrong)
	clrMissed := theme.Get(app.ColorMissed)
	fg = eui.White
	if g.Percent >= g.Advance {
		bg = clrCorrect
	} else if g.Percent >= g.Fallback && g.Percent < g.Advance {
		bg = clrNeutral
	} else if g.Percent < g.Fallback {
		if g.TryDown > 1 {
			bg = clrWrong
		} else {
			bg = clrMissed
		}
	}
	return fg, bg
}

func (g *GameData) LastGameFullResult() (str string) {
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
	s1 := ""
	for _, v := range g.Modalities {
		s1 += "[" + string(v.sym) + ":" + strconv.Itoa(v.correct) + "(" + strconv.Itoa(v.wrong) + "-" + strconv.Itoa(v.missed) + ")] "
	}
	if g.TotalTime > 0 && g.CheckIn {
		str = fmt.Sprintf("#%v %v очков:%v %v ходов(%v) проверка %v ходе [в:%v н:%v э:%v] %v", g.Id, g.GameMode(), g.Score, s1, g.Moves, g.TotalMoves, g.up, g.down, g.same, gameDuration)
	} else {
		str = fmt.Sprintf("#%v %v %v%% %v moves(%v/%v) %v", g.Id, g.GameMode(), g.Percent, s1, g.Moves, g.TotalMoves, gameDuration)
	}
	return str
}
