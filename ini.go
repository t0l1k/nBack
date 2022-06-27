package main

import "image/color"

type Setting struct {
	timeToNextCell, timeShowCell                                     int
	defaultLevel                                                     int
	manual                                                           bool
	manualAdv                                                        int
	thresholdAdvance, thresholdFallback, thresholdFallbackSessions   int
	trials, trialsFactor, trialsExponent                             int
	rr                                                               float64
	feedbackOnUserMove, usecentercell, resetOnFirstWrong, fullScreen bool
	pauseRest                                                        int
}

func NewSettings() *Setting {
	return &Setting{
		timeToNextCell:            2000,
		timeShowCell:              500,
		defaultLevel:              1, // Level in manul mode and first game level today
		manual:                    false,
		manualAdv:                 3, // games with 100% to next level in manual mode, 0 same level
		thresholdAdvance:          80,
		thresholdFallback:         50,
		thresholdFallbackSessions: 3,
		trials:                    5, //20 = classic = trials*factor+level**exponent
		trialsFactor:              1,
		trialsExponent:            2,
		rr:                        12.5, // Random Repition
		usecentercell:             false,
		resetOnFirstWrong:         false,
		fullScreen:                false,
		pauseRest:                 5000,
		feedbackOnUserMove:        true,
	}
}

type Theme struct {
	bg, fg, active, regular, correct, error, warning color.RGBA
}

func NewTheme() *Theme {
	// black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	gray := color.RGBA{192, 192, 192, 255}
	return &Theme{
		bg:      gray,
		fg:      white,
		active:  color.RGBA{255, 255, 0, 255},
		regular: color.RGBA{0, 0, 192, 255},
		correct: color.RGBA{0, 192, 0, 255},
		warning: color.RGBA{255, 128, 0, 255},
		error:   color.RGBA{255, 0, 0, 255},
	}
}
