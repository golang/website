package main

import (
	"fmt"
)

// The index function returns the index of the first occurrence of v in s,
// or -1 if not present.
func index[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

func main() {
	s := []int{1, 3, 5, 2, 4}
	fmt.Println(index(s, 3))
	fmt.Println(index(s, 6))
}
