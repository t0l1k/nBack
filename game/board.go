package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

type status int

const (
	Neutral status = iota
	Regular
	Correct
	Warning
	Error
)

type Board struct {
	rect                               *ui.Rect
	inGame, userMoved, reset           bool
	gameCount, level, totalMoves, move int
	grid                               *ui.GridView
	field                              []*Cell
	container                          []ui.Drawable
	moveNumber                         int
	arr, moves                         []int
	countCorrect, countWrong           int
	dtBeg, dtEnd                       time.Time
	moveStatus                         status
	pref                               *ui.Preferences
	theme                              *ui.Theme
}

func NewBoard(rect []int, pref *ui.Preferences, theme *ui.Theme) *Board {
	rand.Seed(time.Now().UnixNano())
	b := &Board{
		rect:   ui.NewRect(rect),
		inGame: false,
		pref:   pref,
		theme:  theme,
	}
	b.grid = ui.NewGridView(rect, 3, (*b.theme)["game bg"], (*b.theme)["game fg"])
	b.Add(b.grid)
	b.field = b.initCells()
	b.Resize(rect)
	return b
}

func (b *Board) Reset(gameCount, level int) {
	b.inGame = true
	b.reset = false
	b.userMoved = false
	b.grid.Visibe = true
	b.setFieldVisible(true)
	b.gameCount = gameCount
	b.level = level
	b.moves = make([]int, 0)
	b.move = 0
	b.totalMoves = TotalMoves(b.level)
	b.arr = getArr(b.level, b.totalMoves, b.pref)
	b.countCorrect, b.countWrong = 0, 0
	b.dtBeg = time.Now()
	b.MakeMove()
}

func (b *Board) CheckUserMove() {
	b.userMoved = true
	b.moveStatus = Regular
	log.Printf("User Moved %v", b)
}

func (b *Board) CheckMoveRegular() {
	s := fmt.Sprintf("Check regular Move %v", b)
	if len(b.moves) > b.level+1 {
		var i, j int
		if b.inGame {
			i = b.move - b.level - 2
			j = b.move - 1
		} else {
			i = b.move - b.level - 1
			j = b.move
		}
		s += fmt.Sprintf("%v", b.moves[i:j])
		aa := b.moves[i:j]
		if aa[0] == aa[len(aa)-1] && b.userMoved {
			b.moveStatus = Correct
			b.countCorrect += 1
			s += " correct answer!"
		} else if aa[0] == aa[len(aa)-1] && !b.userMoved {
			b.moveStatus = Error
			b.countWrong += 1
			s += " missed the answer!"
		} else if aa[0] != aa[len(aa)-1] && b.userMoved {
			b.moveStatus = Warning
			b.countWrong += 1
			s += fmt.Sprintf(" there was no repeat %v steps back!", b.level)
		}
	} else {
		if b.userMoved {
			b.moveStatus = Warning
			b.countWrong += 1
			s += "error! went early."
		} else {
			s += " analyze early"
		}
	}
	if b.countWrong > 0 && (*b.pref)["reset on first wrong"].(bool) {
		b.reset = true
	}
	b.userMoved = false
	log.Println(s)
}
func (b *Board) MakeMove() {
	b.moveStatus = Neutral
	if b.move == b.totalMoves || b.reset {
		b.inGame = false
		b.CheckMoveRegular()
		b.grid.Visibe = false
		b.setFieldVisible(false)
		b.dtEnd = time.Now()
		return
	}
	b.moveNumber = b.arr[b.move]
	b.moves = append(b.moves, b.moveNumber)
	b.move += 1
}

func (b *Board) setFieldVisible(value bool) {
	for _, cell := range b.field {
		cell.Visibe = value
	}
}

func (b *Board) ShowActiveCell() {
	b.field[b.moveNumber].SetActive(true)
}

func (b *Board) HideActiveCell() {
	b.field[b.moveNumber].SetActive(false)
}

func (b *Board) IsShowActiveCell() bool {
	return b.field[b.moveNumber].Active
}

func (b *Board) getPercent() int {
	if b.reset {
		return 0
	}
	var (
		aa, bb, i, j float64
	)
	aa, bb = float64(b.countCorrect), float64(b.countWrong)
	if aa == 0 && bb == 0 {
		i, j = 1, 0
	} else if aa == 0 && bb > 0 {
		i, j = 0, 1
	} else {
		i, j = aa, bb
	}
	return int(i * 100 / (i + j))
}

func (b *Board) initCells() (field []*Cell) {
	dim := (*b.pref)["grid size"].(int)
	cellBg := (*b.theme)["game bg"]
	cellFg := (*b.theme)["game fg"]
	cellActiveColor := (*b.theme)["game active color"]
	for i := 0; i < dim*dim; i++ {
		isCenter := false
		aX := i % dim
		aY := i / dim
		if aX == dim/2 && aY == dim/2 && !(*b.pref)["use center cell"].(bool) && dim%2 != 0 {
			isCenter = true
		}
		cell := NewCell([]int{0, 0, 1, 1}, isCenter, cellBg, cellFg, cellActiveColor)
		field = append(field, cell)
		b.Add(cell)
	}
	return field
}

func (b *Board) Layout() *ebiten.Image {
	return nil
}

func (b *Board) Add(item ui.Drawable) {
	b.container = append(b.container, item)
}
func (b *Board) Update(dt int) {
	for _, value := range b.container {
		value.Update(dt)
	}
}
func (b *Board) Draw(surface *ebiten.Image) {
	for _, value := range b.container {
		value.Draw(surface)
	}
}

func (b *Board) String() string {
	return fmt.Sprintf("#%v nB%v %v/%v", b.gameCount, b.level, b.move, b.totalMoves)
}

func (b *Board) Resize(rect []int) {
	b.rect = ui.NewRect(rect)
	b.grid.Resize(rect)
	b.resizeCells()
}

func (b *Board) resizeCells() {
	dim := (*b.pref)["grid size"].(int)
	x, y := b.rect.Pos()
	cellSize, _ := b.rect.Size()
	cellSize /= dim
	for i := 0; i < dim*dim; i++ {
		aX := i % dim
		aY := i / dim
		cellX := aX*cellSize + x
		cellY := aY*cellSize + y
		b.field[i].Resize([]int{cellX, cellY, cellSize, cellSize})
	}
}

func TotalMoves(level int) int {
	return (*ui.GetPreferences())["trials"].(int) +
		(*ui.GetPreferences())["trials factor"].(int)*
			int(math.Pow(float64(level), float64((*ui.GetPreferences())["trials exponent"].(int))))
}
