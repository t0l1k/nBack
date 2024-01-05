package data

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

func NewGameConf() GameConf {
	return make(GameConf)
}

func (t GameConf) Get(set GameConfValue) (value interface{}) {
	return t[set]
}

func (t GameConf) Set(set GameConfValue, value interface{}) {
	t[set] = value
}

func DefaultSettings() *GameConf {
	gc := NewGameConf()
	gc.Set(DefaultLevel, 1)
	gc.Set(MoveTime, 1.5)
	gc.Set(RandomRepition, 30)
	gc.Set(GridSize, 3)
	gc.Set(ShowGrid, false)
	gc.Set(UseCenterCell, false)
	gc.Set(ShowCrossHair, true)
	gc.Set(ResetOnFirstWrong, false)
	gc.Set(ThresholdAdvance, 90)
	gc.Set(ThresholdFallback, 75)
	gc.Set(ThresholdFallbackSessions, 1)
	gc.Set(Trials, 5)
	gc.Set(TrialsFactor, 1)
	gc.Set(TrialsExponent, 1)
	gc.Set(MaxNumber, 10)
	gc.Set(UseAddSub, true)
	gc.Set(UseMulDiv, false)
	return &gc
}
