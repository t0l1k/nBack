package main

import "image/color"

type Setting struct {
	timeToNextCell, timeShowCell                                     float64
	defaultLevel                                                     int
	manual                                                           bool
	manualAdv                                                        int
	thresholdAdvance, thresholdFallback, thresholdFallbackSessions   int
	trials, trialsFactor, trialsExponent                             int
	rr                                                               float64
	feedbackOnUserMove, usecentercell, resetOnFirstWrong, fullScreen bool
	pauseRest                                                        int
	gridSize                                                         int
}

func NewSettings() *Setting {
	return &Setting{
		timeToNextCell:            2.0,
		timeShowCell:              0.5,
		trials:                    5, //20 = classic = trials*factor+level**exponent
		trialsFactor:              1,
		trialsExponent:            2,
		thresholdAdvance:          80,
		thresholdFallback:         50,
		thresholdFallbackSessions: 3,
		defaultLevel:              1, // Level in manul mode and first game level today
		manual:                    false,
		manualAdv:                 3, // games with 100% to next level in manual mode, 0 same level
		resetOnFirstWrong:         true,
		rr:                        12.5, // Random Repition
		usecentercell:             false,
		feedbackOnUserMove:        true,
		gridSize:                  3,
		pauseRest:                 5,
		fullScreen:                false,
	}
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
