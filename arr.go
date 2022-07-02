package main

import (
	"log"
	"math/rand"
	"time"
)

func getArr(level, moves int) (arr []int) {
	start := time.Now()
	pause := 3
	count := 0
	check := false
	max := 0
	best := make([]int, 0)
	elapsed := time.Since(start)
	for (int(elapsed.Seconds()) < pause) && count < 100000 && !check {
		arr = genArr(moves)
		var percent int
		check, percent = checkRR(arr, level)
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

func genArr(moves int) (a []int) {
	dim := getApp().preferences.gridSize
	for len(a) < moves {
		num := rand.Intn((dim * dim) - 1)
		if num != (dim*dim-1)/2 && !getApp().preferences.usecentercell || getApp().preferences.usecentercell {
			a = append(a, num)
		}
	}
	return a
}

func checkRR(a []int, level int) (bool, int) {
	RR := getApp().preferences.rr
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
