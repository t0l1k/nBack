package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawRect(surface *ebiten.Image, rect *Rect, color color.Color) {
	x1, y1 := rect.TopLeft()
	x2, y2 := rect.TopRight()
	ebitenutil.DrawLine(surface, float64(x1), float64(y1), float64(x2), float64(y2), color)
	x1, y1 = rect.TopRight()
	x2, y2 = rect.BottomRight()
	ebitenutil.DrawLine(surface, float64(x1), float64(y1), float64(x2), float64(y2), color)
	x1, y1 = rect.BottomRight()
	x2, y2 = rect.BottomLeft()
	ebitenutil.DrawLine(surface, float64(x1), float64(y1), float64(x2), float64(y2), color)
	x1, y1 = rect.BottomLeft()
	x2, y2 = rect.TopLeft()
	ebitenutil.DrawLine(surface, float64(x1), float64(y1), float64(x2), float64(y2), color)
}

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

func GetTip(center Point, percent, lenght, width, height float64) (tip Point) {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	sine := math.Sin(radians)
	cosine := math.Cos(radians)
	tip.X = center.X + lenght*sine - width
	tip.Y = center.Y + lenght*cosine - height
	return tip
}

func GetAngle(percent float64) float64 {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	angle := (radians * -180 / math.Pi)
	return angle
}
