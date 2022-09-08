package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/t0l1k/nBack/game"
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
		bg:     (*ui.GetTheme())["game bg"],
		fg:     (*ui.GetTheme())["fg"],
		Dirty:  true,
		Visibe: true,
	}
}
func (r *OptTheme) Layout() {
	w, h := r.rect.Size()
	cellWidth, cellHeight := w, h/4
	if r.Image == nil {
		r.Image = ebiten.NewImage(w, h)
	} else {
		r.Image.Clear()
	}
	r.Image.Fill(r.bg)
	x, y := 0, 0
	rect := []int{x, y, cellWidth / 2, cellHeight}
	bgLbl := ui.NewLabel("app color background", rect, (*ui.GetTheme())["bg"], (*ui.GetTheme())["fg"])
	bgLbl.SetRect(true)
	bgLbl.Draw(r.Image)
	x, y = cellWidth/2, 0
	rect = []int{x, y, cellWidth / 2, cellHeight}
	fgLbl := ui.NewLabel("app color foreground", rect, (*ui.GetTheme())["fg"], (*ui.GetTheme())["bg"])
	fgLbl.SetRect(true)
	fgLbl.Draw(r.Image)

	x, y = 0, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameBgLbl := ui.NewLabel("game background", rect, (*ui.GetTheme())["game bg"], (*ui.GetTheme())["game fg"])
	gameBgLbl.SetRect(true)
	gameBgLbl.Draw(r.Image)
	x, y = cellWidth/3, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameFgLbl := ui.NewLabel("game foreground", rect, (*ui.GetTheme())["game fg"], (*ui.GetTheme())["game bg"])
	gameFgLbl.SetRect(true)
	gameFgLbl.Draw(r.Image)
	x, y = cellWidth/3*2, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameActiveLbl := ui.NewLabel("game active cell color", rect, (*ui.GetTheme())["game active color"], (*ui.GetTheme())["game bg"])
	gameActiveLbl.SetRect(true)
	gameActiveLbl.Draw(r.Image)

	x, y = 0, cellHeight*2
	w, h = cellWidth/4, cellHeight
	rect = []int{x, y, w, h}
	regularLbl := ui.NewLabel("color regular", rect, (*ui.GetTheme())["regular color"], (*ui.GetTheme())["fg"])
	regularLbl.SetRect(true)
	regularLbl.Draw(r.Image)
	x = cellWidth / 4
	rect = []int{x, y, w, h}
	correctLbl := ui.NewLabel("color correct", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	correctLbl.SetRect(true)
	correctLbl.Draw(r.Image)
	x = cellWidth / 4 * 2
	rect = []int{x, y, w, h}
	warningLbl := ui.NewLabel("color warning", rect, (*ui.GetTheme())["warning color"], (*ui.GetTheme())["fg"])
	warningLbl.SetRect(true)
	warningLbl.Draw(r.Image)
	x = cellWidth / 4 * 3
	rect = []int{x, y, w, h}
	errorLbl := ui.NewLabel("color error", rect, (*ui.GetTheme())["error color"], (*ui.GetTheme())["fg"])
	errorLbl.SetRect(true)
	errorLbl.Draw(r.Image)

	x, y = 0, cellHeight*3
	w, h = cellWidth/len(game.Colors), cellHeight
	sz := len(game.Colors)
	for i, v := range game.Colors {
		cellX := i % sz * w
		ebitenutil.DrawRect(r.Image, float64(cellX), float64(y), float64(w), float64(h), v)
	}

	r.Dirty = false
}
func (r *OptTheme) Update(dt int) {}
func (r *OptTheme) Draw(surface *ebiten.Image) {
	if r.Dirty {
		r.Layout()
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
	r.Image = nil
}

func (r *OptTheme) Close() {
	r.Image.Dispose()
}
