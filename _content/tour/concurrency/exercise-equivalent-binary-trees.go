//go:build nobuild || OMIT

package main

import "golang.org/x/tour/tree"

// Walk walks the sorted tree t sending all values
// in order from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int)

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool

func main() {
}
