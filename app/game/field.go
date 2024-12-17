package game

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/t0l1k/eui/colors"
)

var Colors = []color.Color{colors.Blue, colors.Aqua, colors.Green, colors.Olive, colors.Yellow, colors.Red, colors.Purple, colors.Orange, colors.White, colors.Gray}

type field struct {
	curModal                                       ModalType
	level, totalMoves, randomRepition, dim, maxNum int
	useCenter                                      bool
}

func newField(conf GameConf, level, totalMoves int, sym ModalType) []int {
	beginDt := time.Now()
	f := &field{level: level, totalMoves: totalMoves}
	f.curModal = sym
	f.randomRepition = conf.Get(RandomRepition).(int)
	f.maxNum = conf.Get(MaxNumber).(int)
	f.dim = conf.Get(GridSize).(int)
	f.useCenter = conf.Get(UseCenterCell).(bool)
	percent, max := 0, 0
	best := make([]int, 0)
	count := 0
	for max < f.randomRepition && ((count < 10000 || max < 13) && count < 100000) {
		result := f.generate()
		percent = f.checkRR(result)
		if percent > max {
			max = percent
			best = result
		}
		count++
	}
	log.Printf("generated modality %v field for level %v, moves %v, count %v, RR percent %v, arr:%v %v", f.curModal, f.level, f.totalMoves, count, max, best, time.Since(beginDt))
	return best
}

func (f *field) generate() (result []int) {
	rand.Seed(time.Now().UnixNano())
	center := (f.dim*f.dim - 1) / 2
	num := 0
	for len(result) < f.totalMoves {
		if f.curModal == Pos {
			num = rand.Intn(f.dim * f.dim)
			if num != center && !f.useCenter || f.useCenter {
				result = append(result, num)
			}
		} else if f.curModal == Col {
			num = rand.Intn(len(Colors))
			result = append(result, num)
		} else if f.curModal == Sym || f.curModal == Ari {
			num = rand.Intn(f.maxNum-1) + 1
			result = append(result, num)
		}
	}
	return result
}

func (f *field) checkRR(moves []int) int {
	moveCount := 0
	for i, v := range moves {
		nextMove := i + f.level
		if nextMove > len(moves)-1 {
			break
		}
		if v == moves[nextMove] {
			moveCount++
		}
	}
	percent := 100 * float64(moveCount) / float64(len(moves))
	return int(percent)
}
