package data

import (
	"strconv"
)

type GameConfValue int

const (
	DefaultLevel GameConfValue = iota
	MoveTime
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
)

type GameConf map[GameConfValue]interface{}

func DefaultSettings() *GameConf {
	gc := NewGameConf()
	gc.Set(DefaultLevel, 1)
	gc.Set(MoveTime, 1.5)
	gc.Set(RandomRepition, 30)
	gc.Set(GridSize, 3)
	gc.Set(ShowGrid, true)
	gc.Set(UseCenterCell, false)
	gc.Set(ShowCrossHair, true)
	gc.Set(ResetOnFirstWrong, false)
	gc.Set(ThresholdAdvance, 90)
	gc.Set(ThresholdFallback, 75)
	gc.Set(ThresholdFallbackSessions, 1)
	gc.Set(Trials, 20)
	gc.Set(TrialsFactor, 1)
	gc.Set(TrialsExponent, 2)
	gc.Set(MaxNumber, 10)
	gc.Set(UseAddSub, true)
	gc.Set(UseMulDiv, false)
	return &gc
}

func NewGameConf() GameConf {
	return make(GameConf)
}

func (g GameConf) Get(set GameConfValue) (value interface{}) {
	return g[set]
}

func (g GameConf) Set(set GameConfValue, value interface{}) {
	g[set] = value
}

func (g GameConf) GameConf(gDt *GameData) (result []string) {
	result = append(result, "Уровень следующий:"+strconv.Itoa(gDt.Level))
	result = append(result, "Ходов:"+strconv.Itoa(gDt.TotalMoves))
	result = append(result, "Время хода:"+strconv.FormatFloat(g.Get(MoveTime).(float64), 'f', 2, 64)+" секунд")
	result = append(result, "Переход вверх:"+strconv.Itoa(g.Get(ThresholdAdvance).(int)))
	result = append(result, "Переход вниз:"+strconv.Itoa(g.Get(ThresholdFallback).(int)))
	result = append(result, "Попыток:"+strconv.Itoa(g.Get(ThresholdFallbackSessions).(int)))
	result = append(result, "Показать прицел:"+strconv.FormatBool(g.Get(ShowCrossHair).(bool)))
	if gDt.IsContainMod(Pos) {
		result = append(result, "Размер сетки:"+strconv.Itoa(g.Get(GridSize).(int)))
		result = append(result, "Показать сетку:"+strconv.FormatBool(g.Get(ShowGrid).(bool)))
	}
	if gDt.IsContainMod(Sym) {
		result = append(result, "Наибольшее число:"+strconv.Itoa(g.Get(MaxNumber).(int)))
	}
	if gDt.IsContainMod(Ari) {
		result = append(result, "Наибольшее число:"+strconv.Itoa(g.Get(MaxNumber).(int)))
		result = append(result, "Сложение/Вычитание:"+strconv.FormatBool(g.Get(UseAddSub).(bool)))
		result = append(result, "Умножение/Деление:"+strconv.FormatBool(g.Get(UseMulDiv).(bool)))
	}
	result = append(result, "Процент повторов:"+strconv.Itoa(g.Get(RandomRepition).(int)))
	return result
}
