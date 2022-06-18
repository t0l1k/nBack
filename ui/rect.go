package ui

import "fmt"

type Rect struct {
	X, Y, W, H int
}

func NewRect(arr []int) *Rect {
	return &Rect{
		X: arr[0],
		Y: arr[1],
		W: arr[2],
		H: arr[3],
	}
}

func (r Rect) GetPos() (int, int) {
	return r.X, r.Y
}

func (r Rect) GetSize() (int, int) {
	return r.W, r.H
}

func (r Rect) GetLowestSize() int {
	result := r.W
	if r.W > r.H {
		result = r.H
	}
	return result
}

func (r Rect) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v]", r.X, r.Y, r.W, r.H)
}
