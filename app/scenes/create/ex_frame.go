package create

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
)

const (
	btnJ = "Jaeggi"
	btnB = "Brainworkshop"
	bntQ = "Гадкий утёнок"
	btnP = "Три поросенка"
)

type ExamplesFrame struct {
	eui.DrawableBase
	layout *eui.ContainerBase
	bg, fg color.Color
	th     float32
}

func NewExamplesFrame(f func(b *eui.Button)) *ExamplesFrame {
	e := &ExamplesFrame{}
	e.Visible = true
	e.layout = &eui.ContainerBase{}
	strs := []string{btnJ, btnB, bntQ, btnP}
	lbl := eui.NewText("Примеры настроек")
	e.layout.Add(lbl)
	for _, v := range strs {
		b := eui.NewButton(v, f)
		e.layout.Add(b)
	}
	return e
}

func (e *ExamplesFrame) setup(bg, fg color.Color) {
	e.bg = bg
	e.fg = fg
}

func (e *ExamplesFrame) Layout() {
	e.SpriteBase.Layout()
	e.Image().Fill(e.bg)
	w, h := e.GetRect().Size()
	vector.StrokeRect(e.Image(), 0, 0, float32(w), float32(h), e.th, e.fg, true)

	for _, v := range e.layout.GetContainer() {
		v.Draw(e.Image())
	}
	e.Dirty = false
}

func (e *ExamplesFrame) Update(dt int) {
	for _, v := range e.layout.GetContainer() {
		v.Update(dt)
	}
}

func (v *ExamplesFrame) Draw(surface *ebiten.Image) {
	if !v.Visible {
		return
	}
	if v.Dirty {
		v.Layout()
		for _, c := range v.layout.GetContainer() {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x, y := v.GetRect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(v.Image(), op)
	for _, v := range v.layout.GetContainer() {
		v.Draw(surface)
	}
}

func (e *ExamplesFrame) Resize(r []int) {
	e.Rect(eui.NewRect(r))
	e.SpriteBase.Resize(r)
	x0, y0 := e.GetRect().Pos()
	w0, h0 := e.GetRect().Size()
	e.th = float32(e.GetRect().GetLowestSize()) * 0.02
	rect := eui.NewRect([]int{x0 + int(e.th), y0 + int(e.th), w0 - int(e.th)*2, h0 - int(e.th)*2})
	sz := float32(rect.H) / 5
	for i, value := range e.layout.GetContainer() {
		switch v := value.(type) {
		case *eui.Text:
			v.Resize([]int{rect.X, rect.Y, rect.W, int(sz)})
		case *eui.Button:
			v.Resize([]int{rect.X, rect.Y + int(sz*float32(i)), rect.W, int(sz)})
		}
	}
	e.ImageReset()
}
