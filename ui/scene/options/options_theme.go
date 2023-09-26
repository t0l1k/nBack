package options

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

type OptTheme struct {
	rect          *eui.Rect
	Image         *ebiten.Image
	Dirty, Visibe bool
	bg, fg        color.Color
}

func NewOptTheme(rect []int) *OptTheme {
	return &OptTheme{
		rect:   eui.NewRect(rect),
		bg:     eui.GetTheme().Get("game bg"),
		fg:     eui.GetTheme().Get("fg"),
		Dirty:  true,
		Visibe: true,
	}
}
func (r *OptTheme) Layout() {
	w, h := r.rect.Size()
	cellWidth, cellHeight := w, h/3
	if r.Image == nil {
		r.Image = ebiten.NewImage(w, h)
	} else {
		r.Image.Clear()
	}
	r.Image.Fill(r.bg)
	x, y := 0, 0
	rect := []int{x, y, cellWidth / 2, cellHeight}
	bgLbl := eui.NewLabel("app color background", rect, eui.GetTheme().Get("bg"), eui.GetTheme().Get("fg"))
	bgLbl.SetRect(true)
	bgLbl.Draw(r.Image)
	x, y = cellWidth/2, 0
	rect = []int{x, y, cellWidth / 2, cellHeight}
	fgLbl := eui.NewLabel("app color foreground", rect, eui.GetTheme().Get("fg"), eui.GetTheme().Get("bg"))
	fgLbl.SetRect(true)
	fgLbl.Draw(r.Image)

	x, y = 0, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameBgLbl := eui.NewLabel("game background", rect, eui.GetTheme().Get("game bg"), eui.GetTheme().Get("game fg"))
	gameBgLbl.SetRect(true)
	gameBgLbl.Draw(r.Image)
	x, y = cellWidth/3, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameFgLbl := eui.NewLabel("game foreground", rect, eui.GetTheme().Get("game fg"), eui.GetTheme().Get("game bg"))
	gameFgLbl.SetRect(true)
	gameFgLbl.Draw(r.Image)
	x, y = cellWidth/3*2, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gameActiveLbl := eui.NewLabel("game active cell color", rect, eui.GetTheme().Get("game active color"), eui.GetTheme().Get("game bg"))
	gameActiveLbl.SetRect(true)
	gameActiveLbl.Draw(r.Image)

	x, y = 0, cellHeight*2
	w, h = cellWidth/4, cellHeight
	rect = []int{x, y, w, h}
	regularLbl := eui.NewLabel("color regular", rect, eui.GetTheme().Get("regular color"), eui.GetTheme().Get("fg"))
	regularLbl.SetRect(true)
	regularLbl.Draw(r.Image)
	x = cellWidth / 4
	rect = []int{x, y, w, h}
	correctLbl := eui.NewLabel("color correct", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	correctLbl.SetRect(true)
	correctLbl.Draw(r.Image)
	x = cellWidth / 4 * 2
	rect = []int{x, y, w, h}
	warningLbl := eui.NewLabel("color warning", rect, eui.GetTheme().Get("warning color"), eui.GetTheme().Get("fg"))
	warningLbl.SetRect(true)
	warningLbl.Draw(r.Image)
	x = cellWidth / 4 * 3
	rect = []int{x, y, w, h}
	errorLbl := eui.NewLabel("color error", rect, eui.GetTheme().Get("error color"), eui.GetTheme().Get("fg"))
	errorLbl.SetRect(true)
	errorLbl.Draw(r.Image)
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
	r.rect = eui.NewRect(rect)
	r.Dirty = true
	r.Image = nil
}

func (r *OptTheme) Close() {
	r.Image.Dispose()
}
