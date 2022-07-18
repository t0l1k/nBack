package main

import (
	"image/color"
	"math"
)

type Setting struct {
	TimeToNextCell, TimeShowCell, RR                                 float64
	DefaultLevel, ManualAdv                                          int
	Manual                                                           bool
	ThresholdAdvance, ThresholdFallback, ThresholdFallbackSessions   int
	Trials, TrialsFactor, TrialsExponent                             int
	FeedbackOnUserMove, Usecentercell, ResetOnFirstWrong, FullScreen bool
	PauseRest, GridSize                                              int
}

func NewSettings() *Setting {
	s := &Setting{}
	return s
}

func (s *Setting) Reset() {
	s.TimeToNextCell = 2.0
	s.TimeShowCell = 0.5
	s.Trials = 5 //20 = classic = trials+factor*level**exponent
	s.TrialsFactor = 1
	s.TrialsExponent = 2
	s.ThresholdAdvance = 80
	s.ThresholdFallback = 50
	s.ThresholdFallbackSessions = 3
	s.DefaultLevel = 1 // Level in manul mode and first game level today
	s.Manual = false
	s.ManualAdv = 3 // games with 100% to next level in manual mode, 0 same level
	s.ResetOnFirstWrong = true
	s.RR = 12.5 // Random Repition
	s.Usecentercell = false
	s.FeedbackOnUserMove = true
	s.GridSize = 3
	s.PauseRest = 5
	s.FullScreen = false
}

func (s *Setting) Apply(value *Setting) *Setting {
	s.TimeToNextCell = value.TimeToNextCell
	s.TimeShowCell = value.TimeShowCell
	s.Trials = value.Trials
	s.TrialsFactor = value.TrialsFactor
	s.TrialsExponent = value.TrialsExponent
	s.ThresholdAdvance = value.ThresholdAdvance
	s.ThresholdFallback = value.ThresholdFallback
	s.ThresholdFallbackSessions = value.ThresholdFallbackSessions
	s.DefaultLevel = value.DefaultLevel
	s.Manual = value.Manual
	s.ManualAdv = value.ManualAdv
	s.ResetOnFirstWrong = value.ResetOnFirstWrong
	s.RR = value.RR
	s.Usecentercell = value.Usecentercell
	s.FeedbackOnUserMove = value.FeedbackOnUserMove
	s.GridSize = value.GridSize
	s.PauseRest = value.PauseRest
	s.FullScreen = value.FullScreen
	return s
}

func (s *Setting) Load() *Setting {
	if s.DefaultLevel == 0 {
		s.Reset()
	}
	return s
}

func (s *Setting) TotalMoves(level int) int {
	return s.Trials + s.TrialsFactor*int(math.Pow(float64(level), float64(s.TrialsExponent)))
}

type Theme struct {
	Bg, Fg, GameBg, GameFg, GameActiveColor, RegularColor, CorrectColor, ErrorColor, WarningColor color.RGBA
}

func NewTheme() *Theme {
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	gray := color.RGBA{192, 192, 192, 255}
	yellow := color.RGBA{255, 255, 0, 255}
	blue := color.RGBA{0, 0, 192, 255}
	green := color.RGBA{0, 192, 0, 255}
	orange := color.RGBA{255, 165, 0, 255}
	red := color.RGBA{255, 0, 0, 255}
	return &Theme{
		Bg:              gray,
		Fg:              white,
		GameBg:          black,
		GameFg:          gray,
		GameActiveColor: yellow,
		RegularColor:    blue,
		CorrectColor:    green,
		WarningColor:    orange,
		ErrorColor:      red,
	}
}
