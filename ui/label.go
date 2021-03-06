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
	text                     string
	rect                     *Rect
	Image                    *ebiten.Image
	Dirty, Visible, drawRect bool
	bg, fg                   color.Color
}

func NewLabel(text string, rect []int, bg, fg color.Color) *Label {
	return &Label{
		text:     text,
		rect:     NewRect(rect),
		Image:    nil,
		Dirty:    true,
		Visible:  true,
		drawRect: false,
		bg:       bg,
		fg:       fg}
}

func (l *Label) SetBg(value color.Color) {
	if l.bg == value {
		return
	}
	l.bg = value
	l.Dirty = true
}

func (l *Label) SetRect(value bool) {
	if l.drawRect == value {
		return
	}
	l.drawRect = value
	l.Dirty = true
}

func (l *Label) SetFg(value color.Color) {
	if l.fg == value {
		return
	}
	l.fg = value
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
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return mplusFont
}

func (l *Label) getFontSize() int {
	percent := 0.85
	w, h := l.rect.Size()
	var sz = l.rect.GetLowestSize()
	fontSize := percent * float64(sz)
	fnt := l.getFont(fontSize)
	defer fnt.Close()
	bound := text.BoundString(fnt, l.text)
	for w < bound.Max.X || h < bound.Max.Y {
		fontSize = percent * float64(sz)
		fnt = l.getFont(fontSize)
		defer fnt.Close()
		bound = text.BoundString(fnt, l.text)
		percent -= 0.01
	}
	return int(fontSize)
}

func (l *Label) Layout() *ebiten.Image {
	if !l.Dirty {
		return l.Image
	}
	w, h := l.rect.Size()
	image := ebiten.NewImage(w, h)
	image.Fill(l.bg)
	if l.drawRect {
		ebitenutil.DrawRect(image, 0, 0, float64(w), float64(h), l.fg)
		ebitenutil.DrawRect(image, 2, 2, float64(w)-4, float64(h)-4, l.bg)
	}
	fnt := l.getFont(float64(l.getFontSize()))
	defer fnt.Close()
	b := text.BoundString(fnt, l.text)
	x := (l.rect.W - b.Max.X) / 2
	y := l.rect.H - (l.rect.H-b.Dy())/2
	text.Draw(image, l.text, fnt, x, y, l.fg)
	l.Dirty = false
	return image
}

func (l *Label) Update(dt int) {}
func (l *Label) Draw(surface *ebiten.Image) {
	if l.Dirty {
		l.Image = l.Layout()
	}
	if l.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := l.rect.Pos()
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
