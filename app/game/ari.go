package game

import (
	"math/rand"
)

const (
	add Operation = iota
	sub
	mul
	div
)

type Operation int

func NewOperation() Operation {
	return Operation(0)
}

func (o Operation) Get(a, c, max int) (int, int) {
	b := 0
	switch o {
	case add:
		b = c - a
	case sub:
		for {
			if c > 1 {
				a = getNum(max)
			} else if c == 1 {
				a = 2
			}
			b = a - c
			if (a-b == c || b-a == c) && (a < max || b < max) && (a > 0 && b > 0) {
				break
			}
		}
	case mul:
		var arr []int
		for i := 1; i <= c; i++ {
			if c%i == 0 {
				arr = append(arr, i)
			}
		}
		if len(arr) > 1 {
			a = arr[rand.Intn(len(arr)-1)]
			for {
				b = c / a
				if b*a == c {
					break
				}
			}
		} else {
			a, b = 1, 1
		}
	case div:
		if c == 1 {
			return 1, 1
		}
		var arr []int
		for i := 1; i <= c; i++ {
			if c%i == 0 {
				arr = append(arr, i)
			}
		}
		idx := rand.Intn(len(arr) - 1)
		a = arr[idx]
		for {
			b = c / a
			if a/b == c {
				break
			}
			idx++
			if idx > len(arr) {
				idx = 0
			}
			a = arr[idx]
		}
	}
	return a, b
}

func (o *Operation) Rand(conf GameConf) {
	adds := conf.Get(UseAddSub).(bool)
	muls := conf.Get(UseMulDiv).(bool)
	if adds && !muls {
		*o = Operation(rand.Intn(2))
	} else if muls && !adds {
		*o = Operation(rand.Intn(2) + 2)
	} else {
		*o = Operation(rand.Intn(4))
	}
}

func (o Operation) String() string {
	s := ""
	switch o {
	case add:
		s = "+"
	case sub:
		s = "-"
	case mul:
		s = "*"
	case div:
		s = "/"
	}
	return s
}

func getNum(max int) int {
	n := 0
	for {
		n = rand.Intn(max) + 1
		if n <= max && n >= 0 {
			break
		}
	}
	return n
}
