package intro

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/game"
)

type MovesLine struct {
	eui.DrawableBase
	moves  map[game.ModalType][]game.MoveType
	scores map[game.ModalType][]int
}

func NewMovesLine() *MovesLine {
	c := &MovesLine{}
	return c
}

func (c *MovesLine) Setup(moves map[game.ModalType][]game.MoveType, scores map[game.ModalType][]int) {
	c.moves = moves
	c.scores = scores
	c.Layout()
}

func (c *MovesLine) Layout() {
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(app.GameColorActiveBg)
	regularColor := theme.Get(app.GameColorFg)
	correctColor := theme.Get(app.ColorCorrect)
	wrongColor := theme.Get(app.ColorWrong)
	missedColor := theme.Get(app.ColorMissed)
	c.SpriteBase.Layout()
	c.Image().Fill(bg)
	var (
		col                                  color.Color
		x1, y1, x2, y2, cellSizeW, cellSizeH float32
	)
	w0, h0 := c.GetRect().Size()
	i := 0
	for k, values := range c.moves {
		cellSizeW = float32(w0) / float32(len(values)+1)
		cellSizeH = float32(h0) / float32(len(c.moves))
		y1 = float32(i) * cellSizeH
		y2 = float32(i)*cellSizeH + cellSizeH
		lblModName := eui.NewText(k.String())
		lblModName.Resize([]int{0, int(y1), int(cellSizeW), int(cellSizeH)})
		lblModName.Bg(regularColor)
		lblModName.Draw(c.Image())
		for j, value := range values {
			switch value {
			case game.AddRegular:
				col = regularColor
			case game.AddCorrect:
				col = correctColor
			case game.AddWrong:
				col = wrongColor
			case game.AddMissed:
				col = missedColor
			}
			x1 = cellSizeW * float32(j+1)
			x2 = cellSizeW * float32(j+2)
			x := x1 + (x2-x1)/2
			vector.StrokeLine(c.Image(), x, y1, x, y2, cellSizeW-2, col, true)
			score := c.scores[k][j]
			fmt.Println("sc:", k, j, score, len(c.moves), len(values))
			lblScore := eui.NewText(strconv.Itoa(score))
			lblScore.Resize([]int{int(x1), int(y1), int(cellSizeW), int(cellSizeH)})
			lblScore.Bg(col)
			lblScore.Draw(c.Image())
		}
		i++
		vector.StrokeRect(c.Image(), 0, y1, float32(w0), cellSizeH, 1, regularColor, true)
	}
	c.Dirty = false
}

func (c *MovesLine) Resize(r []int) {
	c.Rect(eui.NewRect(r))
	c.SpriteBase.Resize(r)
	c.ImageReset()
}
