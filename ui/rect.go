package ui

import (
	"fmt"
)

type Point struct {
	X, Y float64
}

func (p Point) Get() (float64, float64) {
	return p.X, p.Y
}

func (p Point) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", p.X, p.Y)
}

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
func (r Rect) InRect(x, y int) bool {
	return r.Left() <= x && r.Right() >= x && r.Top() <= y && r.Bottom() >= y
}

func (r Rect) Pos() (int, int) {
	return r.X, r.Y
}

func (r Rect) Size() (int, int) {
	return r.W, r.H
}

func (r Rect) Left() int {
	return r.X
}

func (r Rect) Right() int {
	return r.X + r.W
}

func (r Rect) Top() int {
	return r.Y
}

func (r Rect) Bottom() int {
	return r.Y + r.H
}
func (r Rect) CenterX() int {
	return (r.Right() - r.X) / 2
}

func (r Rect) CenterY() int {
	return (r.Bottom() - r.Y) / 2
}
func (r Rect) Center() (int, int) {
	return r.CenterX(), r.CenterY()
}

func (r Rect) TopLeft() (int, int) {
	return r.X, r.Y
}

func (r Rect) TopRight() (int, int) {
	return r.X + r.W, r.Y
}

func (r Rect) BottomLeft() (int, int) {
	return r.X, r.Y + r.H
}

func (r Rect) BottomRight() (int, int) {
	return r.X + r.W, r.Y + r.H
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
