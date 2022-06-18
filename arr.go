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
	elapsed := time.Now().Sub(start)
	for (int(elapsed.Seconds()) < pause) && count < 100000 && !check {
		arr = genArr(moves)
		var percent int
		check, percent = checkRR(arr, level)
		if percent > max {
			max = int(percent)
			best = arr
		}
		elapsed = time.Now().Sub(start)
		count += 1
	}
	if !check {
		log.Printf("RR:%v generate:[%v] elapsed time:%v count:%v", max, best, time.Now().Sub(start), count)
		return best
	}
	log.Printf("RR:%v generate:[%v] elapsed time:%v count:%v", max, arr, time.Now().Sub(start), count)
	return arr
}

func genArr(moves int) (a []int) {
	fieldSize := 3
	for len(a) < moves {
		num := rand.Intn((fieldSize * fieldSize) - 1)
		if num != (fieldSize*fieldSize-1)/2 {
			a = append(a, num)
		}
	}
	log.Println("gen:", a)
	return a
}

func checkRR(a []int, level int) (bool, int) {
	RR := 12.5
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
