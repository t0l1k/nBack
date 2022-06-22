package main

import "image/color"

type Theme struct {
	bg, fg, active, regular, correct, error, warning color.RGBA
}

func NewTheme() *Theme {
	return &Theme{
		bg:      color.RGBA{0, 0, 0, 255},
		fg:      color.RGBA{255, 255, 255, 255},
		active:  color.RGBA{255, 255, 0, 255},
		regular: color.RGBA{0, 0, 128, 255},
		correct: color.RGBA{0, 128, 0, 255},
		warning: color.RGBA{255, 128, 0, 255},
		error:   color.RGBA{255, 0, 0, 255},
	}
}
