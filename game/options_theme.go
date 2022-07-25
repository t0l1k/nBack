package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

type OptTheme struct {
	rect          *ui.Rect
	Image         *ebiten.Image
	Dirty, Visibe bool
	bg, fg        color.Color
}

func NewOptTheme(rect []int) *OptTheme {
	return &OptTheme{
		rect:   ui.NewRect(rect),
		bg:     (*ui.GetTheme())["bg"],
		fg:     (*ui.GetTheme())["fg"],
		Dirty:  true,
		Visibe: true,
	}
}
func (r *OptTheme) Layout() *ebiten.Image {
	if !r.Dirty {
		return r.Image
	}
	w, h := r.rect.Size()
	cellWidth, cellHeight := w, h/3
	image := ebiten.NewImage(w, h)
	image.Fill(r.bg)
	x, y := 0, 0
	rect := []int{x, y, cellWidth / 2, cellHeight}
	bgLbl := ui.NewLabel("background", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"])
	bgLbl.SetRect(true)
	bgLbl.Draw(image)
	x, y = cellWidth/2, 0
	rect = []int{x, y, cellWidth / 2, cellHeight}
	fgLbl := ui.NewLabel("foreground", rect, (*ui.GetTheme())["fg"], (*ui.GetTheme())["bg"])
	fgLbl.SetRect(true)
	fgLbl.Draw(image)

	x, y = 0, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameBgLbl := ui.NewLabel("game background", rect, (*ui.GetTheme())["game bg"], (*ui.GetTheme())["game fg"])
	gameBgLbl.SetRect(true)
	gameBgLbl.Draw(image)
	x, y = cellWidth/3, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameFgLbl := ui.NewLabel("game foreground", rect, (*ui.GetTheme())["game fg"], (*ui.GetTheme())["game bg"])
	gameFgLbl.SetRect(true)
	gameFgLbl.Draw(image)
	x, y = cellWidth/3*2, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameActiveLbl := ui.NewLabel("game active color", rect, (*ui.GetTheme())["game active color"], (*ui.GetTheme())["game bg"])
	gameActiveLbl.SetRect(true)
	gameActiveLbl.Draw(image)

	x, y = 0, cellHeight*2
	w, h = cellWidth/4, cellHeight
	rect = []int{x, y, w, h}
	regularLbl := ui.NewLabel("color regular", rect, (*ui.GetTheme())["regular color"], (*ui.GetTheme())["fg"])
	regularLbl.SetRect(true)
	regularLbl.Draw(image)
	x = cellWidth / 4
	rect = []int{x, y, w, h}
	correctLbl := ui.NewLabel("color correct", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	correctLbl.SetRect(true)
	correctLbl.Draw(image)
	x = cellWidth / 4 * 2
	rect = []int{x, y, w, h}
	warningLbl := ui.NewLabel("color warning", rect, (*ui.GetTheme())["warning color"], (*ui.GetTheme())["fg"])
	warningLbl.SetRect(true)
	warningLbl.Draw(image)
	x = cellWidth / 4 * 3
	rect = []int{x, y, w, h}
	errorLbl := ui.NewLabel("color error", rect, (*ui.GetTheme())["error color"], (*ui.GetTheme())["fg"])
	errorLbl.SetRect(true)
	errorLbl.Draw(image)
	r.Dirty = false
	return image
}
func (r *OptTheme) Update(dt int) {}
func (r *OptTheme) Draw(surface *ebiten.Image) {
	if r.Dirty {
		r.Image = r.Layout()
	}
	if r.Visibe {
		op := &ebiten.DrawImageOptions{}
		x, y := r.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(r.Image, op)
	}
}

func (r *OptTheme) Resize(rect []int) {
	r.rect = ui.NewRect(rect)
	r.Dirty = true
}
