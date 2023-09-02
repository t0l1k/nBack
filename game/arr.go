package game

import (
	"log"
	"math/rand"
	"time"

	ui "github.com/t0l1k/eui"
)

func getArr(level, moves int, pref *ui.Preferences) (arr []int) {
	start := time.Now()
	pause := 3
	count := 0
	check := false
	max := 0
	best := make([]int, 0)
	elapsed := time.Since(start)
	for (int(elapsed.Seconds()) < pause) && count < 100000 && !check {
		arr = genArr(moves, pref)
		var percent int
		check, percent = checkRR(arr, level, pref)
		if percent > max {
			max = int(percent)
			best = arr
		}
		elapsed = time.Since(start)
		count += 1
	}
	if !check {
		log.Printf("RR:%v elapsed time:%v count:%v", max, time.Since(start), count)
		return best
	}
	log.Printf("RR:%v elapsed time:%v count:%v", max, time.Since(start), count)
	return arr
}

func genArr(moves int, pref *ui.Preferences) (a []int) {
	dim := pref.Get("grid size").(int)
	center := (dim*dim - 1) / 2
	num := 0
	for len(a) < moves {
		gt := pref.Get("game type").(string)
		if gt == Pos {
			num = rand.Intn((dim * dim))
			if num != center && !pref.Get("use center cell").(bool) || pref.Get("use center cell").(bool) {
				a = append(a, num)
			}
		} else if gt == Col {
			num = rand.Intn(len(Colors))
			a = append(a, num)
		} else if pref.Get("game type").(string) == Sym {
			num = rand.Intn(pref.Get("symbols count").(int)-1) + 1
			a = append(a, num)
		} else if gt == Ari {
			num = rand.Intn(pref.Get("ariphmetic max").(int)-1) + 1
			a = append(a, num)
		}
	}
	return a
}

func checkRR(a []int, level int, pref *ui.Preferences) (bool, int) {
	RR := pref.Get("random repition").(float64)
	count := 0
	for i, v := range a {
		nextMove := i + level
		if nextMove > len(a)-1 {
			break
		}
		if v == a[nextMove] {
			count += 1
		}
	}
	perc := 100 * float64(count) / float64(len(a))
	return perc > RR && perc < 80, int(perc)
}
