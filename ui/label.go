package ui

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Label struct {
	text                    string
	rect                    *Rect
	Image                   *ebiten.Image
	Dirty, Visibe, DrawRect bool
	bg, fg                  color.RGBA
}

func NewLabel(text string, rect []int) *Label {
	return &Label{
		text:     text,
		rect:     NewRect(rect),
		Image:    nil,
		Dirty:    true,
		Visibe:   true,
		DrawRect: false,
		bg:       color.RGBA{0, 128, 0, 255},
		fg:       color.RGBA{255, 255, 0, 255}}
}

func (l *Label) SetBg(value color.RGBA) {
	if l.bg == value {
		return
	}
	l.bg = value
	l.Dirty = true
}
func (l *Label) SetText(value string) {
	if l.text == value {
		return
	}
	l.text = value
	l.Dirty = true
}

func (*Label) getFont(size float64) font.Face {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return mplusFont
}

func (l *Label) getFontSize() int {
	percent := 0.85
	w, _ := l.rect.GetSize()
	var sz = l.rect.GetLowestSize()
	fontSize := percent * float64(sz)
	fnt := l.getFont(fontSize)
	defer fnt.Close()
	bound, _ := font.BoundString(fnt, l.text)
	for w < bound.Max.X.Ceil() {
		fontSize = percent * float64(sz)
		fnt = l.getFont(fontSize)
		defer fnt.Close()
		bound, _ = font.BoundString(fnt, l.text)
		percent -= 0.01
	}
	return int(fontSize)
}

func (l *Label) Layout() *ebiten.Image {
	if !l.Dirty {
		return l.Image
	}
	w, h := l.rect.GetSize()
	image := ebiten.NewImage(w, h)
	image.Fill(l.bg)
	if l.DrawRect {
		ebitenutil.DrawRect(image, 0, 0, float64(w), float64(h), l.fg)
		ebitenutil.DrawRect(image, 2, 2, float64(w)-4, float64(h)-4, l.bg)
	}
	fnt := l.getFont(float64(l.getFontSize()))
	defer fnt.Close()
	bound, _ := font.BoundString(fnt, l.text)
	wF := (bound.Max.X - bound.Min.X).Ceil()
	hF := (bound.Max.Y - bound.Min.Y).Ceil()
	x := (w - wF) / 2
	y := h - (h-hF)/2
	text.Draw(image, l.text, fnt, x, y, l.fg)
	l.Dirty = false
	return image
}

func (l *Label) Update(dt int) {}
func (l *Label) Draw(surface *ebiten.Image) {
	if l.Dirty {
		l.Image = l.Layout()
	}
	if l.Visibe {
		op := &ebiten.DrawImageOptions{}
		x, y := l.rect.GetPos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(l.Image, op)
	}
}

func (l *Label) Resize(rect []int) {
	l.rect = NewRect(rect)
	l.Dirty = true
}

func (l Label) String() string {
	return fmt.Sprintf("%v %v", l.text, l.rect)
}
