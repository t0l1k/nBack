package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	w, h := r.rect.Size()
	image := ebiten.NewImage(w, h)
	image.Fill(r.bg)
	rows := r.getRows()
	boxWidth := r.rect.W / rows
	boxHeight := int(float64(r.rect.GetLowestSize()) * 0.05)
	keys := make([]int, 0)
	for k := range getApp().db.todayData {
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
		x, y := r.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(r.Image, op)
	}
}

func (r *ResultLbls) Resize(rect []int) {
	r.rect = ui.NewRect(rect)
	r.Dirty = true
}

type ResultPlot struct {
	rect          *ui.Rect
	Image         *ebiten.Image
	Dirty, Visibe bool
	bg, fg        color.RGBA
}

func NewResultPlot(rect []int) *ResultPlot {
	return &ResultPlot{
		rect:   ui.NewRect(rect),
		bg:     color.RGBA{0, 64, 0, 0},
		fg:     color.RGBA{255, 255, 0, 255},
		Dirty:  true,
		Visibe: true,
	}
}
func (r *ResultPlot) Layout() *ebiten.Image {
	if !r.Dirty {
		return r.Image
	}
	xArr, yArr, lvlValues, percents, colors := getApp().db.todayData.PlotData()
	fmt.Println(xArr, yArr, lvlValues, percents, colors)
	axisXMax := xArr.Len()
	axisYMax := getApp().db.todayData.getMax() + 2
	w0, h0 := r.rect.Size()
	image := ebiten.NewImage(w0, h0)
	bg := r.bg
	fg := r.fg
	image.Fill(bg)
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
	ebitenutil.DrawLine(image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
	xTicks := xArr.Len()
	gridWidth := 0
	lastW := 0
	for i := 1; i < xTicks+1; i++ {
		x := axisXMax * i / xTicks
		x1, y1 := xPos(float64(x)), axisRect.Bottom()
		x2, y2 := xPos(float64(x)), axisRect.Bottom()+margin/4
		ebitenutil.DrawLine(image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
		x1, y1 = xPos(float64(x)), axisRect.Bottom()
		x2, y2 = xPos(float64(x)), axisRect.Top()
		ebitenutil.DrawLine(image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
		gridWidth = int(xPos(float64(x))) - int(xPos(float64(lastW)))
		lastW = x
		if i%5 == 0 || i == 1 || i == xTicks {
			xL, yL := int(xPos(float64(x))-float64(margin)/2), axisRect.Bottom()+int(float64(margin)*0.1)
			w, h = margin, margin
			lbl := ui.NewLabel(strconv.Itoa(x), []int{xL, yL, w, h})
			lbl.SetBg(bg)
			lbl.Draw(image)
		}
	}
	if gridWidth > margin {
		gridWidth = margin
	}
	// y axis
	x1, y1 = axisRect.BottomLeft()
	x2, y2 = axisRect.TopLeft()
	ebitenutil.DrawLine(image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
	yTicks := axisYMax
	for i := 1; i < yTicks+1; i++ {
		y = axisYMax * i / yTicks
		x1, y1 := axisRect.Left(), yPos(float64(y))
		x2, y2 := axisRect.Left()-margin/4, yPos(float64(y))
		ebitenutil.DrawLine(image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
		x1, y1 = axisRect.Left(), yPos(float64(y))
		x2, y2 = axisRect.Right(), yPos(float64(y))
		ebitenutil.DrawLine(image, float64(x1), float64(y1), float64(x2), float64(y2), fg)
		boxSize := int(float64(axisRect.GetLowestSize()) * 0.05)
		xL, yL := axisRect.Left()-int(float64(boxSize)*1.2), int(yPos(float64(y))-float64(boxSize)/2)
		w, h = boxSize, boxSize
		lbl := ui.NewLabel(strconv.Itoa(y), []int{xL, yL, w, h})
		lbl.SetBg(bg)
		lbl.Draw(image)
	}
	{
		boxSize := margin * 7
		xL, yL := axisRect.Right()/2-boxSize/2, axisRect.Top()-int(float64(boxSize)/4.5)
		w, h = boxSize, boxSize/4
		lbl := ui.NewLabel("Daily results", []int{xL, yL, w, h})
		lbl.SetBg(bg)
		lbl.Draw(image)
	}
	{
		boxSize := margin
		xL, yL := axisRect.Left()+int(float64(boxSize)*0.2), axisRect.Top()-boxSize
		w, h = int(float64(boxSize)*1.5), boxSize
		lbl := ui.NewLabel("Level", []int{xL, yL, w, h})
		lbl.SetBg(bg)
		lbl.Draw(image)
	}
	{
		boxSize := margin
		xL, yL := axisRect.Right()-boxSize*3, axisRect.Bottom()-boxSize
		w, h = boxSize*3, boxSize
		lbl := ui.NewLabel("Game number", []int{xL, yL, w, h})
		lbl.SetBg(bg)
		lbl.Draw(image)
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
	{
		points := zip(xArr, yArr)
		var results1 []float64
		for e := points.Front(); e != nil; e = e.Next() {
			x := e.Value.(*list.List).Front().Value
			y := e.Value.(*list.List).Back().Value
			xx := xPos(float64(axisXMax) * float64(x.(int)) / float64(xArr.Len()))
			yy := yPos(float64(y.(int)))
			results1 = append(results1, xx, yy)
			fmt.Printf("x:%v y:%v\n", xx, yy)
		}
		fmt.Println(results1)
		for i, j := 0, 1; j < len(results1)-2; i, j = i+2, j+2 {
			x1, y1, x2, y2 := results1[i], results1[j], results1[i+2], results1[j+2]
			fmt.Println(i, j, x1, y1, x2, y2)
			ebitenutil.DrawLine(image, x1, y1, x2, y2, color.RGBA{0, 255, 0, 255})
		}
	}
	{
		points := zip(xArr, lvlValues)
		var results1 []float64
		for e := points.Front(); e != nil; e = e.Next() {
			x := e.Value.(*list.List).Front().Value
			y := e.Value.(*list.List).Back().Value
			xx := xPos(float64(axisXMax) * float64(x.(int)) / float64(xArr.Len()))
			yy := yPos(y.(float64))
			results1 = append(results1, xx, yy)
			fmt.Printf("x:%v y:%v\n", xx, yy)
		}
		fmt.Println(results1)
		var perc []int
		for e := percents.Front(); e != nil; e = e.Next() {
			perc = append(perc, e.Value.(int))
		}
		var clrs []color.RGBA
		for e := colors.Front(); e != nil; e = e.Next() {
			clrs = append(clrs, e.Value.(color.RGBA))
		}
		for i, j := 0, 1; j < len(results1)-2; i, j = i+2, j+2 {
			x1, y1, x2, y2 := results1[i], results1[j], results1[i+2], results1[j+2]
			fmt.Println(i, j, x1, y1, x2, y2)
			ebitenutil.DrawLine(image, x1, y1, x2, y2, color.RGBA{0, 0, 255, 255})
		}
		k := 0
		for i, j := 0, 1; j < len(results1); i, j = i+2, j+2 {
			x1, y1 := results1[i], results1[j]
			fmt.Println(i, j, x1, y1, perc[k])
			boxSize := gridWidth
			lbl := ui.NewLabel(strconv.Itoa(perc[k]), []int{int(x1) - boxSize/2, int(y1) - boxSize/2, boxSize, boxSize})
			lbl.SetBg(clrs[k])
			lbl.Draw(image)
			k++
		}
	}
	r.Dirty = false
	return image
}
func (r *ResultPlot) Update(dt int) {}
func (r *ResultPlot) Draw(surface *ebiten.Image) {
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

func (r *ResultPlot) Resize(rect []int) {
	r.rect = ui.NewRect(rect)
	r.Dirty = true
}
