package main

import "image/color"

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
