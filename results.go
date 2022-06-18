package main

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

type ResultLbls struct {
	rect          *ui.Rect
	Image         *ebiten.Image
	Dirty, Visibe bool
	bg, fg        color.RGBA
}

func NewResultLbls(rect []int) *ResultLbls {
	return &ResultLbls{
		rect:   ui.NewRect(rect),
		bg:     color.RGBA{0, 64, 0, 255},
		fg:     color.RGBA{255, 255, 0, 255},
		Dirty:  true,
		Visibe: true,
	}
}
func (r *ResultLbls) getRows() (result int) {
	w, _ := getApp().GetScreenSize()
	if w <= 640 {
		result = 2
	} else if w <= 800 {
		result = 3
	} else if w <= 1024 {
		result = 4
	} else {
		result = 5
	}
	return
}
func (r *ResultLbls) Layout() *ebiten.Image {
	if !r.Dirty {
		return r.Image
	}
	w, h := r.rect.GetSize()
	image := ebiten.NewImage(w, h)
	image.Fill(r.bg)
	rows := r.getRows()
	boxWidth := r.rect.W / rows
	boxHeight := int(float64(r.rect.GetLowestSize()) * 0.05)
	keys := make([]int, 0)
	for k, _ := range getApp().db.todayData {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for i, v := range keys {
		x := i % rows
		y := i / rows
		str := getApp().db.todayData[v].ShortStr()
		l := ui.NewLabel(str, []int{x * boxWidth, y * boxHeight, boxWidth, boxHeight})
		l.SetBg(getApp().db.todayData[v].BgColor())
		l.Draw(image)
	}
	r.Dirty = false
	return image
}
func (r *ResultLbls) Update(dt int) {}
func (r *ResultLbls) Draw(surface *ebiten.Image) {
	if r.Dirty {
		r.Image = r.Layout()
	}
	if r.Visibe {
		op := &ebiten.DrawImageOptions{}
		x, y := r.rect.GetPos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(r.Image, op)
	}
}
