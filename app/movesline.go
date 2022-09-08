package app

import (
	"container/list"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/ui"
)

type MovesLine struct {
	rect           *ui.Rect
	Image          *ebiten.Image
	Dirty, Visible bool
	bg, fg         color.Color
}

func NewMovesLine(rect []int) *MovesLine {
	return &MovesLine{
		rect:    ui.NewRect(rect),
		bg:      (*ui.GetTheme())["game bg"],
		fg:      (*ui.GetTheme())["game fg"],
		Dirty:   false,
		Visible: false,
	}
}
func (p *MovesLine) Layout() {
	w0, h0 := p.rect.Size()
	if p.Image == nil {
		p.Image = ebiten.NewImage(w0, h0)
	} else {
		p.Image.Clear()
	}
	count := data.GetDb().TodayGamesCount
	if count < 1 || !p.Visible {
		return
	}
	bg := p.bg
	fg := p.fg
	red, g, b, a := fg.RGBA()
	a /= 3
	fg2 := color.RGBA{uint8(red), uint8(g), uint8(b), uint8(a)}
	p.Image.Fill(bg)
	xArr, colors := data.GetDb().TodayData[count].MovesColor()
	var yArr list.List
	for i := 0; i < xArr.Len(); i++ {
		yArr.PushBack(1)
	}
	axisXMax := xArr.Len()
	axisYMax := 1
	margin := int(float64(p.rect.GetLowestSize()) * 0.2)
	x, y := margin, margin
	w, h := w0-margin*2, h0-margin*2
	axisRect := ui.NewRect([]int{x, y, w, h})

	lerp := func(t, inStart, inEnd, outStart, outEnd float64) float64 {
		return outStart + (t-inStart)/(inEnd-inStart)*(outEnd-outStart)
	}
	xPos := func(x float64) float64 {
		return math.Round(lerp(x, 0, float64(axisXMax), float64(axisRect.Left()), float64(axisRect.Right())))
	}
	yPos := func(y float64) float64 {
		return math.Round(lerp(y, 0, float64(axisYMax), float64(axisRect.Bottom()), float64(axisRect.Top())))
	}
	x1, y1 := axisRect.BottomLeft()
	x2, y2 := axisRect.BottomRight()
	ebitenutil.DrawLine(p.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg2)
	x1, y1 = axisRect.TopLeft()
	x2, y2 = axisRect.TopRight()
	y1 += margin / 4
	y2 += margin / 4
	ebitenutil.DrawLine(p.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg2)
	xTicks := xArr.Len()
	for i := 1; i < xTicks+1; i++ {
		x := (float64(axisXMax) * float64(i-1) / float64(xArr.Len()-1))
		x1, y1 := xPos(float64(x)), axisRect.Bottom()
		x2, y2 := xPos(float64(x)), axisRect.Bottom()+margin/4
		ebitenutil.DrawLine(p.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
		x1, y1 = xPos(float64(x)), axisRect.Bottom()
		x2, y2 = xPos(float64(x)), axisRect.Top()
		ebitenutil.DrawLine(p.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg2)
	}

	{
		boxSize := margin * 2
		xL, yL := axisRect.Right()-boxSize*3, axisRect.Bottom()-boxSize
		w, h = boxSize*3, boxSize
		lbl := ui.NewLabel(ui.GetLocale().Get("optmv"), []int{xL, yL, w, h}, fg2, fg)
		defer lbl.Close()
		lbl.SetBg(bg)
		lbl.Draw(p.Image)
	}

	zip := func(a, b list.List) *list.List {
		if a.Len() != b.Len() {
			panic("len(a) != len(b)")
		}
		r := list.New()
		for e, j := a.Front(), b.Front(); e != nil; e, j = e.Next(), j.Next() {
			a := list.New()
			a.PushBack(e.Value)
			a.PushBack(j.Value)
			r.PushBack(a)
		}
		return r
	}
	{ // parse data
		points := zip(xArr, yArr)
		var results1, results2 []float64
		for e := points.Front(); e != nil; e = e.Next() {
			x := e.Value.(*list.List).Front().Value
			y := e.Value.(*list.List).Back().Value
			xx := xPos(float64(axisXMax) * float64(x.(int)-1) / float64(xArr.Len()-1))
			yy := yPos(float64(y.(int)))
			yy2 := yPos(0)
			results1 = append(results1, xx, yy)
			results2 = append(results2, xx, yy2)
		}
		var clrs []color.Color
		for e := colors.Front(); e != nil; e = e.Next() {
			clrs = append(clrs, e.Value.(color.Color))
		}
		k := 0
		for i, j := 0, 1; j < len(results2); i, j = i+2, j+2 { // moves line
			x1, y1, x2, y2 := results1[i], results1[j], results2[i], results2[j]
			ebitenutil.DrawLine(p.Image, x1, y1, x2, y2, clrs[k])
			k++
		}
	}
	p.Dirty = false
}

func (p *MovesLine) Update(dt int) {}
func (p *MovesLine) Draw(surface *ebiten.Image) {
	if p.Dirty {
		p.Layout()
	}
	if p.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := p.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(p.Image, op)
	}
}

func (p *MovesLine) Resize(rect []int) {
	p.rect = ui.NewRect(rect)
	p.Dirty = true
	p.Image = nil
}
func (p *MovesLine) Close() {
	p.Image.Dispose()
}
