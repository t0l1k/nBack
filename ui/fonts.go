package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/t0l1k/nBack/res"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func init() {
	fontsInstance = GetFonts()
}

var fontsInstance *Fonts = nil

func GetFonts() (f *Fonts) {
	if fontsInstance == nil {
		f = &Fonts{}
		log.Printf("Fonts init done")
	} else {
		f = fontsInstance
	}
	return f
}

type Fonts map[int]font.Face

func (f Fonts) addFont(size int) {
	tt, err := opentype.Parse(res.RobotoRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	f[size] = mplusFont
	log.Println("Added font size", size)
}

func (f Fonts) get(size int) font.Face {
	for k, v := range f {
		if k == size {
			return v
		}
	}
	f.addFont(size)
	return f[size]
}

func (f Fonts) calcFontSize(txt string, rect *Rect) int {
	var fontSize float64
	percent := 0.85
	w, h := rect.Size()
	sz := rect.GetLowestSize()
	for {
		fontSize = percent * float64(sz)
		fnt := f.get(int(fontSize))
		defer fnt.Close()
		bound := text.BoundString(fnt, txt)
		if w > bound.Max.X && h > bound.Max.Y {
			break
		}
		percent -= 0.01
	}
	return int(fontSize)
}
