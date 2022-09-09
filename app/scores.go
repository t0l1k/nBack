package app

import (
	"container/list"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/ui"
)

type ScorePlot struct {
	rect          *ui.Rect
	Image         *ebiten.Image
	Dirty, Visibe bool
	bg, fg        color.Color
}

func NewScorePlot(rect []int) *ScorePlot {
	return &ScorePlot{
		rect:   ui.NewRect(rect),
		bg:     (*ui.GetTheme())["bg"],
		fg:     (*ui.GetTheme())["fg"],
		Dirty:  true,
		Visibe: true,
	}
}
func (r *ScorePlot) Layout() {
	xArr, yArr, avgArr, strsArr := data.GetDb().ScoresData.PlotData()
	axisXMax := xArr.Len()
	score, _ := data.GetDb().ReadAllGamesScore()
	axisYMax := score.Max + 1
	w0, h0 := r.rect.Size()
	if r.Image == nil {
		r.Image = ebiten.NewImage(w0, h0)
	} else {
		r.Image.Clear()
	}
	bg := r.bg
	fg := r.fg
	red, g, b, a := fg.RGBA()
	a /= 3
	fg2 := color.RGBA{uint8(red), uint8(g), uint8(b), uint8(a)}
	r.Image.Fill(bg)
	margin := int(float64(r.rect.GetLowestSize()) * 0.05)
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
	// x axis
	x1, y1 := axisRect.BottomLeft()
	x2, y2 := axisRect.BottomRight()
	ebitenutil.DrawLine(r.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
	xTicks := xArr.Len()
	gridWidth := 0
	lastW := 0
	for i := 1; i < xTicks+1; i++ {
		x := axisXMax * i / xTicks
		x1, y1 := xPos(float64(x)), axisRect.Bottom()
		x2, y2 := xPos(float64(x)), axisRect.Bottom()+margin/4
		ebitenutil.DrawLine(r.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
		x1, y1 = xPos(float64(x)), axisRect.Bottom()
		x2, y2 = xPos(float64(x)), axisRect.Top()
		ebitenutil.DrawLine(r.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg2)
		gridWidth = int(xPos(float64(x))) - int(xPos(float64(lastW)))
		lastW = x
		if i%5 == 0 || i == 1 || i == xTicks {
			xL, yL := int(xPos(float64(x))-float64(margin)/2), axisRect.Bottom()+int(float64(margin)*0.1)
			w, h = margin, margin
			lbl := ui.NewLabel(strconv.Itoa(x), []int{xL, yL, w, h}, bg, fg)
			defer lbl.Close()
			lbl.SetBg(bg)
			lbl.Draw(r.Image)
		}
	}
	if gridWidth > margin*2 {
		gridWidth = margin * 2
	}
	{
		boxSize := margin
		xL, yL := axisRect.Right()-boxSize*3, axisRect.Bottom()-boxSize
		w, h = boxSize*3, boxSize
		lbl := ui.NewLabel(ui.GetLocale().Get("lblDays"), []int{xL, yL, w, h}, bg, fg)
		defer lbl.Close()
		lbl.SetBg(bg)
		lbl.Draw(r.Image)
	}
	// y axis
	x1, y1 = axisRect.BottomLeft()
	x2, y2 = axisRect.TopLeft()
	ebitenutil.DrawLine(r.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
	yTicks := axisYMax
	for i := 1; i < yTicks+1; i++ {
		y = axisYMax * i / yTicks
		x1, y1 := axisRect.Left(), yPos(float64(y))
		x2, y2 := axisRect.Left()-margin/4, yPos(float64(y))
		ebitenutil.DrawLine(r.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
		x1, y1 = axisRect.Left(), yPos(float64(y))
		x2, y2 = axisRect.Right(), yPos(float64(y))
		ebitenutil.DrawLine(r.Image, float64(x1), float64(y1), float64(x2), float64(y2), fg2)
		boxSize := int(float64(axisRect.GetLowestSize()) * 0.05)
		xL, yL := axisRect.Left()-int(float64(boxSize)*1.2), int(yPos(float64(y))-float64(boxSize)/2)
		w, h = boxSize, boxSize
		lbl := ui.NewLabel(strconv.Itoa(y), []int{xL, yL, w, h}, bg, fg)
		defer lbl.Close()
		lbl.SetBg(bg)
		lbl.Draw(r.Image)
	}
	{
		boxSize := margin
		xL, yL := axisRect.Left()+int(float64(boxSize)*0.2), axisRect.Top()-boxSize
		w, h = int(float64(boxSize)*1.5), boxSize
		lbl := ui.NewLabel(ui.GetLocale().Get("lblLevel"), []int{xL, yL, w, h}, bg, fg)
		defer lbl.Close()
		lbl.SetBg(bg)
		lbl.Draw(r.Image)
	}
	{
		boxSize := margin * 7
		xL, yL := axisRect.Right()/2-boxSize/2, axisRect.Top()-int(float64(boxSize)/4.5)
		w, h = boxSize, boxSize/3
		lbl := ui.NewLabel(ui.GetLocale().Get("btnScore"), []int{xL, yL, w, h}, bg, fg)
		defer lbl.Close()
		lbl.SetBg(bg)
		lbl.Draw(r.Image)
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

	{ // label
		points := zip(xArr, yArr)
		var results1, results2 []float64
		for e := points.Front(); e != nil; e = e.Next() {
			x := e.Value.(*list.List).Front().Value
			y := e.Value.(*list.List).Back().Value
			xx := xPos(float64(axisXMax) * float64(x.(int)) / float64(xArr.Len()))
			yy := yPos(float64(y.(int)))
			yy2 := yPos(0)
			results1 = append(results1, xx, yy)
			results2 = append(results2, xx, yy2)
		}
		var strs []string
		for e := strsArr.Front(); e != nil; e = e.Next() {
			strs = append(strs, e.Value.(string))
		}
		k := 0
		for i, j := 0, 1; j < len(results1); i, j = i+2, j+2 {
			x1, y1 := results2[i], results2[j]
			var x, y, w, h, boxSize float64
			boxSize = float64(gridWidth) / 2
			x, y = 0, 0
			w, h = results2[j]-results1[j], boxSize
			rect := []int{int(x), int(y), int(w), int(h)}
			lbl := ui.NewLabel(strs[k], rect, (*ui.GetTheme())["correct color"], fg)
			defer lbl.Close()
			lbl.Layout()
			w1, h1 := lbl.Image.Size()
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(-float64(w1)/2, -float64(h1)/2)
			count := -90
			op.GeoM.Rotate(float64(count%360) * 2 * math.Pi / 360)
			op.GeoM.Translate(x1, y1-float64(w1)/2)
			r.Image.DrawImage(lbl.Image, &op)
			k++
		}
	}

	{ // parse data - max line
		points := zip(xArr, yArr)
		var results1 []float64
		xx := xPos(float64(axisXMax) * float64(0) / float64(xArr.Len()))
		yy := yPos(float64(0))
		results1 = append(results1, xx, yy)
		for e := points.Front(); e != nil; e = e.Next() {
			x := e.Value.(*list.List).Front().Value
			y := e.Value.(*list.List).Back().Value
			xx := xPos(float64(axisXMax) * float64(x.(int)) / float64(xArr.Len()))
			yy := yPos(float64(y.(int)))
			results1 = append(results1, xx, yy)
		}
		for i, j := 0, 1; j < len(results1)-2; i, j = i+2, j+2 {
			x1, y1, x2, y2 := results1[i], results1[j], results1[i+2], results1[j+2]
			ebitenutil.DrawLine(r.Image, x1, y1, x2, y2, (*ui.GetTheme())["regular color"])
		}
	}
	{ // parse data - average line
		points := zip(xArr, avgArr)
		var results1 []float64
		xx := xPos(float64(axisXMax) * float64(0) / float64(xArr.Len()))
		yy := yPos(float64(0))
		results1 = append(results1, xx, yy)
		for e := points.Front(); e != nil; e = e.Next() {
			x := e.Value.(*list.List).Front().Value
			y := e.Value.(*list.List).Back().Value
			xx := xPos(float64(axisXMax) * float64(x.(int)) / float64(xArr.Len()))
			yy := yPos(float64(y.(float64)))
			results1 = append(results1, xx, yy)
		}
		for i, j := 0, 1; j < len(results1)-2; i, j = i+2, j+2 {
			x1, y1, x2, y2 := results1[i], results1[j], results1[i+2], results1[j+2]
			ebitenutil.DrawLine(r.Image, x1, y1, x2, y2, (*ui.GetTheme())["warning color"])
		}
	}

	r.Dirty = false
}
func (r *ScorePlot) Update(dt int) {}
func (r *ScorePlot) Draw(surface *ebiten.Image) {
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

func (r *ScorePlot) Resize(rect []int) {
	r.rect = ui.NewRect(rect)
	r.Dirty = true
	r.Image = nil
}

func (r *ScorePlot) Close() {
	r.Image.Dispose()
}