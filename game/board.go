package game

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

const (
	Pos string = "p"
	Col string = "c"
	Sym string = "s"
)

type Status string

const (
	Neutral Status = "nil"
	Regular Status = "regular"
	Correct Status = "correct"
	Warning Status = "missed"
	Error   Status = "wrong"
)

var Colors = []color.Color{ui.Blue, ui.Aqua, ui.Green, ui.Olive, ui.Yellow, ui.Red, ui.Purple, ui.Orange, ui.White, ui.Gray}

type Board struct {
	rect                                  *ui.Rect
	InGame, UserMoved, reset              bool
	gameCount, level, TotalMoves, Move    int
	grid                                  *ui.GridView
	field                                 []*Cell
	container                             []ui.Drawable
	moveValue                             int
	arr, moves                            []int
	CountCorrect, CountWrong, CountMissed int
	DtBeg, DtEnd                          time.Time
	MoveStatus                            Status
	MovesStatus                           map[int]Status
	pref                                  *ui.Preferences
	theme                                 *ui.Theme
}

func NewBoard(rect []int, pref *ui.Preferences, theme *ui.Theme) *Board {
	rand.Seed(time.Now().UnixNano())
	b := &Board{
		rect:   ui.NewRect(rect),
		InGame: false,
		pref:   pref,
		theme:  theme,
	}
	if b.pref.Get("show grid").(bool) && b.pref.Get("game type").(string) == Pos {
		gridSz := b.pref.Get("grid size").(int)
		b.grid = ui.NewGridView(rect, ui.NewPoint(float64(gridSz), float64(gridSz)), (*b.theme)["game bg"], (*b.theme)["game fg"])
		b.Add(b.grid)
	}
	b.field = b.initCells()
	b.Resize(rect)
	return b
}

func (b *Board) Reset(gameCount, level int) {
	b.InGame = true
	b.reset = false
	b.UserMoved = false
	if b.pref.Get("show grid").(bool) && b.pref.Get("game type").(string) == Pos {
		b.grid.Visible = true
	}
	b.setFieldVisible(true)
	b.gameCount = gameCount
	b.level = level
	b.moves = make([]int, 0)
	b.Move = 0
	b.MovesStatus = make(map[int]Status)
	b.MoveStatus = Neutral
	b.TotalMoves = TotalMoves(b.level)
	b.arr = getArr(b.level, b.TotalMoves, b.pref)
	b.CountCorrect, b.CountWrong, b.CountMissed = 0, 0, 0
	b.DtBeg = time.Now()
	b.MakeMove()
}

func (b *Board) CheckUserMove() {
	b.UserMoved = true
	b.MoveStatus = Regular
	log.Printf("User Moved %v", b)
}

func (b *Board) CheckMoveRegular() {
	s := fmt.Sprintf("Check regular Move %v", b)
	mv := b.Move - 1
	if len(b.moves) > b.level+1 {
		var i, j int
		if b.InGame {
			i = b.Move - b.level - 2
			j = b.Move - 1
		} else {
			i = b.Move - b.level - 1
			j = b.Move
		}
		mv = j
		s += fmt.Sprintf("%v", b.moves[i:j])
		aa := b.moves[i:j]
		if aa[0] == aa[len(aa)-1] && b.UserMoved {
			b.MovesStatus[mv] = Correct
			b.MoveStatus = Correct
			b.CountCorrect += 1
			s += " correct answer!"
		} else if aa[0] == aa[len(aa)-1] && !b.UserMoved {
			b.MovesStatus[mv] = Error
			b.MoveStatus = Error
			b.CountMissed += 1
			s += " missed the answer!"
		} else if aa[0] != aa[len(aa)-1] && b.UserMoved {
			b.MovesStatus[mv] = Warning
			b.MoveStatus = Warning
			b.CountWrong += 1
			s += fmt.Sprintf(" there was no repeat %v steps back!", b.level)
		}
	} else {
		if b.UserMoved {
			b.MovesStatus[mv] = Warning
			b.MoveStatus = Warning
			b.CountWrong += 1
			s += "error! went early."
		} else {
			s += " analyze early"
		}
	}
	if (b.CountWrong > 0 || b.CountMissed > 0) && (*b.pref)["reset on first wrong"].(bool) {
		b.reset = true
	}
	b.UserMoved = false
	log.Println(s)
}
func (b *Board) MakeMove() {
	b.MoveStatus = Neutral
	if b.Move == b.TotalMoves || b.reset {
		b.InGame = false
		b.CheckMoveRegular()
		if b.pref.Get("show grid").(bool) && b.pref.Get("game type").(string) == Pos {
			b.grid.Visible = false
		}
		b.setFieldVisible(false)
		b.DtEnd = time.Now()
		return
	}
	b.moveValue = b.arr[b.Move]
	b.moves = append(b.moves, b.moveValue)
	b.Move += 1
	b.MovesStatus[b.Move] = Regular
}

func (b *Board) setFieldVisible(value bool) {
	for _, cell := range b.field {
		cell.Visibe = value
	}
}

func (b *Board) ShowActiveCell() {
	var mv int
	if (*b.pref)["game type"].(string) == Pos {
		mv = b.moveValue
	} else if (*b.pref)["game type"].(string) == Col {
		mv = 0
		b.field[mv].SetActiveColor(Colors[b.moveValue])
	} else if (*b.pref)["game type"].(string) == Sym {
		mv = 0
		b.field[mv].SetSymbol(b.moveValue)
	}
	b.field[mv].SetActive(true)
}

func (b *Board) HideActiveCell() {
	var mv int
	if (*b.pref)["game type"].(string) == Pos {
		mv = b.moveValue
	} else if (*b.pref)["game type"].(string) == Col {
		mv = 0
	} else if (*b.pref)["game type"].(string) == Sym {
		mv = 0
	}
	b.field[mv].SetActive(false)
}

func (b *Board) IsShowActiveCell() bool {
	var mv int
	if (*b.pref)["game type"].(string) == Pos {
		mv = b.moveValue
	} else if (*b.pref)["game type"].(string) == Col {
		mv = 0
	} else if (*b.pref)["game type"].(string) == Sym {
		mv = 0
	}
	return b.field[mv].Active
}

func (b *Board) GetPercent() int {
	if b.reset {
		return 0
	}
	var (
		aa, bb, i, j float64
	)
	aa, bb = float64(b.CountCorrect), float64(b.CountWrong+b.CountMissed)
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
	var dim int
	if (*b.pref)["game type"].(string) == Pos {
		dim = (*b.pref)["grid size"].(int)
	} else {
		dim = 1
	}
	cellBg := (*b.theme)["game bg"]
	cellFg := (*b.theme)["game fg"]
	cellActiveColor := (*b.theme)["game active color"]
	showCrossHair := b.pref.Get("show crosshair").(bool)
	for i := 0; i < dim*dim; i++ {
		isCenter := false
		aX := i % dim
		aY := i / dim
		if aX == dim/2 && aY == dim/2 && !(*b.pref)["use center cell"].(bool) && dim%2 != 0 && showCrossHair {
			isCenter = true
		}
		cell := NewCell([]int{0, 0, 1, 1}, isCenter, cellBg, cellFg, cellActiveColor)
		field = append(field, cell)
		b.Add(cell)
	}
	return field
}

func (b *Board) Layout() {
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
	tp := b.pref.Get("game type").(string)
	return fmt.Sprintf("#%v %vB%v %v/%v", b.gameCount, tp, b.level, b.Move, b.TotalMoves)
}

func (b *Board) Resize(rect []int) {
	b.rect = ui.NewRect(rect)
	if b.pref.Get("show grid").(bool) && b.pref.Get("game type").(string) == Pos {
		b.grid.Resize(rect)
	}
	b.resizeCells()
}

func (b *Board) resizeCells() {
	var dim int
	if (*b.pref)["game type"].(string) == Pos {
		dim = (*b.pref)["grid size"].(int)
	} else {
		dim = 1
	}
	x, y := b.rect.Pos()
	cellSize, _ := b.rect.Size()
	cellSize /= dim
	for i, v := range b.field {
		aX := i % dim
		aY := i / dim
		cellX := aX*cellSize + x
		cellY := aY*cellSize + y
		v.Resize([]int{cellX, cellY, cellSize, cellSize})
	}
}

func (b *Board) Close() {
	for _, v := range b.container {
		v.Close()
	}
	for _, v := range b.field {
		v.Close()
	}
}

func TotalMoves(level int) int {
	return (*ui.GetPreferences())["trials"].(int) +
		(*ui.GetPreferences())["trials factor"].(int)*
			int(math.Pow(float64(level), float64((*ui.GetPreferences())["trials exponent"].(int))))
}
