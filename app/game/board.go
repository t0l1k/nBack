package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/data"
)

type Board struct {
	eui.View
	layout                   *eui.GridLayoutRightDown
	gData                    *data.GameData
	cells                    []*cell
	Move, LastMove, TestMove int
	show                     bool
}

func New() *Board {
	g := &Board{}
	g.SetupView()
	g.layout = eui.NewGridLayoutRightDown(1, 1)
	conf := eui.GetUi().GetSettings()
	g.Bg(conf.Get(app.GameColorBg).(color.Color))
	return g
}

func (g *Board) Setup(gameData *data.GameData) {
	g.gData = gameData
	dim := 1
	if g.gData.IsContainMod(data.Pos) {
		conf := eui.GetUi().GetSettings()
		dim = conf.Get(app.GridSize).(int)
	}
	g.layout.SetRows(dim)
	g.layout.SetColumns(dim)
	if len(g.cells) == 0 {
		for i := 0; i < dim*dim; i++ {
			conf := eui.GetUi().GetSettings()
			showCrosshair := conf.Get(app.ShowCrossHair).(bool)
			useCenterCell := conf.Get(app.UseCenterCell).(bool)
			isCenter := false
			aX := i % dim
			aY := i / dim
			if aX == dim/2 && aY == dim/2 && !useCenterCell && dim%2 != 0 && showCrosshair {
				isCenter = true
			}
			cell := newCell(isCenter)
			g.cells = append(g.cells, cell)
			g.layout.Add(cell)
		}
	}
	g.Add(g.layout)
	for _, v := range g.gData.Modalities {
		v.AddField(newField(gameData.Level, gameData.TotalMoves, v.GetSym()))
	}
	for _, cell := range g.cells {
		cell.Setup(g.gData.GetModalities())
	}
}

func (g *Board) MakeMove() {
	g.LastMove = g.Move
	level := g.gData.Level
	if g.Move >= level {
		g.TestMove = g.Move - level
	}
	if g.Move == 0 {
		g.TestMove = -1
	} else {
		prevIdx := func() (idx int) {
			if len(g.cells) == 1 {
				return 0
			}
			if g.gData.IsContainMod(data.Pos) {
				idx = g.gData.GetModalityValues(data.Pos)[g.Move-1]
			}
			return idx
		}()
		g.cells[prevIdx].SetInactive()
	}
	idx := func() (idx int) {
		if len(g.cells) == 1 {
			return 0
		}
		if g.gData.IsContainMod(data.Pos) {
			idx = g.gData.GetModalityValues(data.Pos)[g.LastMove]
		}
		return idx
	}()
	g.cells[idx].SetActive(g.Move)
}

func (g *Board) NextMove() {
	g.MakeMove()
	g.Move++
}

func (g *Board) IsVisible() bool {
	return g.show
}

func (g *Board) Visible(value bool) {
	if g.show != value {
		g.show = value
	}
	for _, v := range g.cells {
		v.Visible(value)
	}
}

func (g *Board) Update(dt int) {
	for _, cell := range g.layout.Container {
		cell.Update(dt)
	}
}

func (g *Board) Draw(surface *ebiten.Image) {
	for _, cell := range g.layout.Container {
		cell.Draw(surface)
	}
}
