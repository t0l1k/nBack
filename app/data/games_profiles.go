package data

import (
	"sort"

	"github.com/t0l1k/nBack/app/game"
)

type GameProfiles map[string]*GamesData

func DefalutGameProfiles() *GameProfiles {
	p := make(GameProfiles)
	p.AddGameProfile("Example for novice triple pos/sym/col (move 3 sec)", func() *game.GameConf {
		conf := game.DefaultSettings()
		conf.Set(game.Trials, 10)
		conf.Set(game.TrialsFactor, 1)
		conf.Set(game.TrialsExponent, 1)
		conf.Set(game.ThresholdAdvance, 90)
		conf.Set(game.ThresholdFallback, 75)
		conf.Set(game.ThresholdFallbackSessions, 1)
		conf.Set(game.GridSize, 3)
		conf.Set(game.ShowGrid, true)
		conf.Set(game.MoveTime, 3.0)
		conf.Set(game.Modals, game.Pos+game.Sym+game.Col)
		return conf
	}())
	return &p
}

func (p GameProfiles) AddGameProfile(name string, value *game.GameConf) {
	p[name] = NewGamesData(value)
}

func (p GameProfiles) GetGameProfiles() map[string]*GamesData {
	return p
}

func (p GameProfiles) GetGamesData(name string) *GamesData {
	return p[name]
}

func (p GameProfiles) GetProfilesName() (result []string) {
	for k := range p {
		result = append(result, k)
	}
	sort.Strings(sort.StringSlice(result))
	return result
}
