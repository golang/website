//go:build ignore || OMIT
// +build ignore OMIT

package main

import "golang.org/x/tour/pic"

type Image struct{}

func main() {
	m := Image{}
	pic.ShowImage(m)
}
