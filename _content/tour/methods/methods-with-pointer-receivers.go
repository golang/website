//go:build OMIT
// +build OMIT

package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	v := &Vertex{3, 4}
	fmt.Printf("缩放前：%+v，绝对值：%v\n", v, v.Abs())
	v.Scale(5)
	fmt.Printf("缩放后：%+v，绝对值：%v\n", v, v.Abs())
}
