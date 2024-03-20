// +build OMIT

package main

import (
	"fmt"
)

func Sqrt(x float64) (float64, error) {
	return 0, nil
}

func main() {
	for _, number := range []float64{2, -2} {
		if root, err := Sqrt(number); err == nil {
			fmt.Println(root)
		} else {
			fmt.Println(err)
		}
	}
}
