package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawCircle(surface *ebiten.Image, x, y, radius float64, color color.Color, fill bool) {
	minAngle := math.Acos(1 - 1/radius)
	for angle := float64(0); angle <= 360; angle += minAngle {
		xDelta := radius * math.Cos(angle)
		yDelta := radius * math.Sin(angle)
		x1 := math.Round(x + xDelta)
		y1 := math.Round(y + yDelta)
		if fill {
			if y1 < y {
				for y2 := y1; y2 <= y; y2++ {
					surface.Set(int(x1), int(y2), color)
				}
			} else {
				for y2 := y1; y2 > y; y2-- {
					surface.Set(int(x1), int(y2), color)
				}
			}
		}
		surface.Set(int(x1), int(y1), color)
	}
}
