package data

type GameProfiles map[string]*GamesData

func DefalutGameProfiles() *GameProfiles {
	p := make(GameProfiles)
	p.AddGameProfile("Example for novice triple pos/sym/col (move 3 sec)", func() *GameConf {
		conf := DefaultSettings()
		conf.Set(Trials, 10)
		conf.Set(TrialsFactor, 1)
		conf.Set(TrialsExponent, 1)
		conf.Set(ThresholdAdvance, 90)
		conf.Set(ThresholdFallback, 75)
		conf.Set(ThresholdFallbackSessions, 1)
		conf.Set(GridSize, 3)
		conf.Set(ShowGrid, true)
		conf.Set(MoveTime, 3.0)
		conf.Set(Modals, Pos+Sym+Col)
		return conf
	}())
	return &p
}

func (p GameProfiles) AddGameProfile(name string, value *GameConf) {
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
	return result
}
