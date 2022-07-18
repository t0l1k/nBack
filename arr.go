package main

import (
	"log"
	"math/rand"
	"time"
)

func getArr(level, moves int, pref *Setting) (arr []int) {
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

func genArr(moves int, pref *Setting) (a []int) {
	dim := pref.GridSize
	for len(a) < moves {
		num := rand.Intn((dim * dim) - 1)
		if num != (dim*dim-1)/2 && !pref.Usecentercell || pref.Usecentercell {
			a = append(a, num)
		}
	}
	return a
}

func checkRR(a []int, level int, pref *Setting) (bool, int) {
	RR := pref.RR
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
