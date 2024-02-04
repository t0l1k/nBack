package game

import (
	"time"

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
	theme := eui.GetUi().GetTheme()
	g.Bg(theme.Get(app.GameColorBg))
	return g
}

func (g *Board) Setup(conf data.GameConf, gameData *data.GameData) {
	g.gData = gameData
	dim := 1
	if g.gData.IsContainMod(data.Pos) {
		dim = conf.Get(data.GridSize).(int)
	}
	g.layout.SetRows(dim)
	g.layout.SetColumns(dim)
	if len(g.cells) == 0 {
		for i := 0; i < dim*dim; i++ {
			showCrosshair := conf.Get(data.ShowCrossHair).(bool)
			useCenterCell := conf.Get(data.UseCenterCell).(bool)
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
	for _, v := range g.gData.Modalities {
		v.AddField(newField(conf, gameData.Level, gameData.TotalMoves, v.GetSym()))
	}
	for _, cell := range g.cells {
		cell.Setup(conf, g.gData.GetModalities())
	}
	g.gData.DtBeg = time.Now().Format("2006-01-02 15:04:05.000")
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
	g.Visible(true)
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
	for _, cell := range g.layout.GetContainer() {
		cell.Update(dt)
	}
}

func (g *Board) Draw(surface *ebiten.Image) {
	for _, cell := range g.layout.GetContainer() {
		cell.Draw(surface)
	}
}

func (c *Board) Resize(rect []int) {
	c.View.Resize(rect)
	c.layout.Resize(rect)
}
