package game

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app/data"
)

var Colors = []color.Color{eui.Blue, eui.Aqua, eui.Green, eui.Olive, eui.Yellow, eui.Red, eui.Purple, eui.Orange, eui.White, eui.Gray}

type field struct {
	curModal                           string
	level, totalMoves, rr, dim, maxNum int
	useCenter                          bool
}

func newField(conf data.GameConf, level, totalMoves int, sym string) []int {
	beginDt := time.Now()
	f := &field{level: level, totalMoves: totalMoves}
	f.curModal = sym
	f.rr = conf.Get(data.RandomRepition).(int)
	f.maxNum = conf.Get(data.MaxNumber).(int)
	f.dim = conf.Get(data.GridSize).(int)
	f.useCenter = conf.Get(data.UseCenterCell).(bool)
	check := false
	percent, max := 0, 0
	best := make([]int, 0)
	count := 0
	for !check {
		result := f.generate()
		check, percent = f.checkRR(result)
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
	center := (f.dim*f.dim - 1) / 2
	num := 0
	for len(result) < f.totalMoves {
		if f.curModal == data.Pos {
			num = rand.Intn(f.dim * f.dim)
			if num != center && !f.useCenter || f.useCenter {
				result = append(result, num)
			}
		} else if f.curModal == data.Col {
			num = rand.Intn(len(Colors))
			result = append(result, num)
		} else if f.curModal == data.Sym || f.curModal == data.Ari {
			num = rand.Intn(f.maxNum-1) + 1
			result = append(result, num)
		}
	}
	return result
}

func (f *field) checkRR(arr []int) (bool, int) {
	count := 0
	for i, v := range arr {
		nextMove := i + f.level
		if nextMove > len(arr)-1 {
			break
		}
		if v == arr[nextMove] {
			count++
		}
	}
	perc := 100 * float64(count) / float64(len(arr))
	return perc > float64(f.rr), int(perc)
}
