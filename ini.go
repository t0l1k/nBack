package main

import "image/color"

type Setting struct {
	timeToNextCell, timeShowCell, rr                                 float64
	defaultLevel, manualAdv                                          int
	manual                                                           bool
	thresholdAdvance, thresholdFallback, thresholdFallbackSessions   int
	trials, trialsFactor, trialsExponent                             int
	feedbackOnUserMove, usecentercell, resetOnFirstWrong, fullScreen bool
	pauseRest, gridSize                                              int
}

func NewSettings() *Setting {
	s := &Setting{}
	return s
}

func (s *Setting) Reset() {
	s.timeToNextCell = 2.0
	s.timeShowCell = 0.5
	s.trials = 5 //20 = classic = trials*factor+level**exponent
	s.trialsFactor = 1
	s.trialsExponent = 2
	s.thresholdAdvance = 80
	s.thresholdFallback = 50
	s.thresholdFallbackSessions = 3
	s.defaultLevel = 1 // Level in manul mode and first game level today
	s.manual = false
	s.manualAdv = 3 // games with 100% to next level in manual mode, 0 same level
	s.resetOnFirstWrong = true
	s.rr = 12.5 // Random Repition
	s.usecentercell = false
	s.feedbackOnUserMove = true
	s.gridSize = 3
	s.pauseRest = 5
	s.fullScreen = false
}

func (s *Setting) Apply(value *Setting) {
	s.timeToNextCell = value.timeToNextCell
	s.timeShowCell = value.timeShowCell
	s.trials = value.trials
	s.trialsFactor = value.trialsFactor
	s.trialsExponent = value.trialsExponent
	s.thresholdAdvance = value.thresholdAdvance
	s.thresholdFallback = value.thresholdFallback
	s.thresholdFallbackSessions = value.thresholdFallbackSessions
	s.defaultLevel = value.defaultLevel
	s.manual = value.manual
	s.manualAdv = value.manualAdv
	s.resetOnFirstWrong = value.resetOnFirstWrong
	s.rr = value.rr
	s.usecentercell = value.usecentercell
	s.feedbackOnUserMove = value.feedbackOnUserMove
	s.gridSize = value.gridSize
	s.pauseRest = value.pauseRest
	s.fullScreen = value.fullScreen
	getApp().db.InsertSettings(s)
}

func (s *Setting) Load() *Setting {
	if s.defaultLevel == 0 {
		s.Reset()
	}
	return s
}

type Theme struct {
	bg, fg, gameBg, gameFg, gameActiveColor, regular, correct, error, warning color.RGBA
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
		bg:              gray,
		fg:              white,
		gameBg:          black,
		gameFg:          gray,
		gameActiveColor: yellow,
		regular:         blue,
		correct:         green,
		warning:         orange,
		error:           red,
	}
}
